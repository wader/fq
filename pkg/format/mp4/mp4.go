package mp4

// Quicktime file format https://developer.apple.com/standards/qtff-2001.pdf
// FLAC in ISOBMFF https://github.com/xiph/flac/blob/master/doc/isoflac.txt
// https://www.webmproject.org/vp9/mp4/
// TODO: validate structure better? trak/stco etc
// TODO: rename atom -> box?
// TODO: fmp4, default samples sizes etc
// TODO: keep track of structure somehow to detect errors
// TODO: ISO-14496 says mp4 mdat can begin and end with original header/trailer (no used i guess?)
// TODO: heic decode hevc samples (iloc box?)

import (
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/format"
	"sort"
	"strings"
)

var aacFrameFormat []*decode.Format
var av1CCRFormat []*decode.Format
var av1FrameFormat []*decode.Format
var flacFrameFormat []*decode.Format
var flacMetadatablockFormat []*decode.Format
var mp3FrameFormat []*decode.Format
var mpegAVCDCRFormat []*decode.Format
var mpegAVCAUFormat []*decode.Format
var mpegESFormat []*decode.Format
var mpegHEVCDCRFrameFormat []*decode.Format
var mpegHEVCSampleFormat []*decode.Format
var opusPacketFrameFormat []*decode.Format
var vorbisPacketFormat []*decode.Format
var vp9FrameFormat []*decode.Format
var vpxCCRFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MP4,
		Description: "MPEG-4 file",
		Groups:      []string{format.PROBE},
		// TODO: implment MIME()
		MIMEs:    []string{"audio/mp4", "video/mp4"},
		DecodeFn: mp4Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AV1_CCR}, Formats: &av1CCRFormat},
			{Names: []string{format.AV1_FRAME}, Formats: &av1FrameFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &flacFrameFormat},
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &flacMetadatablockFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3FrameFormat},
			{Names: []string{format.AAC_FRAME}, Formats: &aacFrameFormat},
			{Names: []string{format.MPEG_AVC_DCR}, Formats: &mpegAVCDCRFormat},
			{Names: []string{format.MPEG_AVC_AU}, Formats: &mpegAVCAUFormat},
			{Names: []string{format.MPEG_ES}, Formats: &mpegESFormat},
			{Names: []string{format.HEVC_DCR}, Formats: &mpegHEVCDCRFrameFormat},
			{Names: []string{format.HEVC_NAL}, Formats: &mpegHEVCSampleFormat},
			{Names: []string{format.OPUS_PACKET}, Formats: &opusPacketFrameFormat},
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacketFormat},
			{Names: []string{format.VP9_FRAME}, Formats: &vp9FrameFormat},
			{Names: []string{format.VPX_CCR}, Formats: &vpxCCRFormat},
		},
	})
}

type stsc struct {
	firstChunk      uint32
	samplesPerChunk uint32
}

type moof struct {
	offset            int64
	defaultSampleSize uint32
	dataOffset        uint32
	samplesSizes      []uint32
}

type track struct {
	id         uint32
	dataFormat string
	subType    string
	stco       []uint64 //
	stsc       []stsc
	stsz       []uint32
	decodeOpts []decode.Options
	objectType int // if data format is "mp4a"

	moofs       []*moof // for fmp4
	currentMoof *moof
}

type decodeContext struct {
	tracks            map[uint32]*track
	currentTrack      *track
	currentMoofOffset int64
}

func decodeAtom(ctx *decodeContext, d *decode.D) uint64 {
	boxes := map[string]func(ctx *decodeContext, d *decode.D){
		"ftyp": func(ctx *decodeContext, d *decode.D) {
			d.FieldUTF8("major_brand", 4)
			d.FieldU32("minor_version")
			numBrands := d.BitsLeft() / 8 / 4
			var i int64
			d.FieldArrayLoopFn("brands", func() bool { return i < numBrands }, func(d *decode.D) {
				d.FieldStrFn("brand", func() (string, string) {
					return strings.TrimSpace(d.UTF8(4)), ""
				})
				i++
			})
		},
		"mvhd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldUTF8("flags", 3)
			d.FieldU32("creation_time")
			d.FieldU32("modification_time")
			d.FieldU32("time_scale")
			d.FieldU32("duration")
			d.FieldFP32("preferred_rate")
			d.FieldFP16("preferred_volume")
			d.FieldUTF8("reserved", 10)
			d.FieldUTF8("matrix_structure", 36)
			d.FieldU32("preview_time")
			d.FieldU32("preview_duration")
			d.FieldU32("poster_time")
			d.FieldU32("selection_time")
			d.FieldU32("selection_duration")
			d.FieldU32("current_time")
			d.FieldU32("next_track_id")
		},
		"trak": decodeAtoms,
		"edts": decodeAtoms,
		"elst": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("table", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldU32("track_duration")
				d.FieldU32("media_time")
				d.FieldFP32("media_rate")
				i++
			})
		},
		"tref": decodeAtoms,
		"tkhd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			d.FieldU32("creation_time")
			d.FieldU32("modification_time")
			trackID := uint32(d.FieldU32("track_id"))
			d.FieldU32("reserved1")
			d.FieldU32("duration")
			d.FieldBitBufLen("reserved2", 8*8)
			d.FieldU16("layer")
			// TODO: values
			d.FieldU16("alternate_group")
			d.FieldFP16("volume")
			d.FieldU16("reserved3")
			d.FieldBitBufLen("matrix_structure", 36*8)
			d.FieldFP32("track_width")
			d.FieldFP32("track_height")

			if _, ok := ctx.tracks[trackID]; !ok {
				t := &track{id: trackID}
				ctx.tracks[trackID] = t
				ctx.currentTrack = t
			} else {
				// TODO: dup track id?
			}
		},
		"mdia": decodeAtoms,
		"mdhd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			// TODO: timestamps
			d.FieldU32("creation_time")
			d.FieldU32("modification_time")
			d.FieldU32("time_scale")
			d.FieldU32("duration")
			d.FieldU16("language")
			d.FieldU16("quality")
		},

		"hdlr": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			d.FieldUTF8("component_type", 4)
			subType := d.FieldUTF8("component_subtype", 4)
			d.FieldUTF8("component_manufacturer", 4)
			d.FieldU32("component_flags")
			d.FieldU32("component_flags_mask")
			d.FieldUTF8("component_name", int(d.BitsLeft()/8))

			if ctx.currentTrack != nil {
				// component_type seems to be all zero sometimes so can't look for "mhlr"
				switch subType {
				case "vide", "soun":
					ctx.currentTrack.subType = subType
				}
			}
		},

		"minf": decodeAtoms,
		"dinf": decodeAtoms,
		"dref": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("references", "reference", func() bool { return i < numEntries }, func(d *decode.D) {
				size := d.FieldU32("size")
				d.FieldUTF8("type", 4)
				d.FieldU8("version")
				d.FieldU24("flags")
				dataSize := size - 12
				d.FieldBitBufLen("data", int64(dataSize*8))
				i++
			})
		},
		"stbl": decodeAtoms,
		"stsd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldArrayLoopFn("sample_descriptions", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldStructFn("sample_description", func(d *decode.D) {
					// TODO: decode len?
					size := d.FieldU32("size")
					dataFormat := d.FieldUTF8("data_format", 4)
					subType := ""
					if ctx.currentTrack != nil {
						ctx.currentTrack.dataFormat = dataFormat
						subType = ctx.currentTrack.subType
					}

					d.DecodeLenFn(int64(size-8)*8, func(d *decode.D) {
						d.FieldBytesLen("reserved", 6)
						d.FieldU16("data_reference_index")
						version := d.FieldU16("version")
						d.FieldU16("revision_level")
						d.FieldU32("max_packet_size")

						// "Some sample descriptions terminate with four zero bytes that are not otherwise indicated."
						// uses decodeAtomsZeroTerminate

						// Timecode sample
						// TODO: tc64
						// d.DecodeRangeFn(firstBit, nBits, func(d *decode.D) {
						// 	d.FieldStructFn("sample", func(d *decode.D) {
						// 		d.FieldU32("reserved0")
						// 		d.FieldBitBufLen("data", d.BitsLeft())
						// 		// d.FieldU32("flags")
						// 		// d.FieldS32("timescale")
						// 		// d.FieldS32("frame_duration")
						// 		// d.FieldS32("num_frames")
						// 		// d.FieldU8("reserved1")
						// 	})
						// })

						switch subType {
						case "soun":
							// AudioSampleEntry
							// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFChap3/qtff3.html#//apple_ref/doc/uid/TP40000939-CH205-SW1
							switch version {
							case 0:
								d.FieldU16("num_audio_channels")
								d.FieldU16("sample_size")
								d.FieldU16("compression_id")
								d.FieldU16("packet_size")
								d.FieldFP32("sample_rate")
								// case 2:
								// 	d.FieldU16("revision_level")
								// 	d.FieldU32("vendor")
								// 	d.FieldU16("always_3")
								// 	d.FieldU16("always_16")
								// 	d.FieldU16("always_minus_2")
								// 	d.FieldU32("always_0")
								// 	d.FieldU32("always_65536")
								// 	d.FieldU32("size_of_struct_only")
								// 	d.FieldF64("sample_rate")
								// 	d.FieldU32("num_audio_channels")
								// 	d.FieldU32("always_7f000000")
								// 	d.FieldU32("const_bits_per_channel")
								// 	d.FieldU32("format_specific_flags")
								// 	d.FieldU32("const_bytes_per_audio_packet")
								// 	d.FieldU32("const_lpcm_frames_per_audio_packet")

								if d.BitsLeft() > 0 {
									decodeAtomsZeroTerminate(ctx, d)
								}
							default:
								d.FieldBitBufLen("data", d.BitsLeft())
							}
						case "vide":
							// VideoSampleEntry
							switch version {
							case 0:
								d.FieldU32("temporal_quality")
								d.FieldU32("spatial_quality")
								d.FieldU16("width")
								d.FieldU16("height")
								d.FieldFP32("horizontal_resolution")
								d.FieldFP32("vertical_resolution")
								d.FieldU32("data_size")
								d.FieldU16("frame_count")
								d.FieldUTF8("compression_name", 32)
								d.FieldU16("depth")
								d.FieldS16("color_table_id")
								// TODO: if 0 decode ctab
								if d.BitsLeft() > 0 {
									decodeAtomsZeroTerminate(ctx, d)
								}
							default:
								d.FieldBitBufLen("data", d.BitsLeft())
							}
							// case "hint": TODO: Hint entry
						default:
							d.FieldBitBufLen("data", d.BitsLeft())
						}
					})

					i++
				})

			})
		},
		"avcC": func(ctx *decodeContext, d *decode.D) {
			_, dv := d.FieldDecode("value", mpegAVCDCRFormat)
			avcDcrOut, ok := dv.(format.AvcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected AvcDcrOut got %#+v", dv))
			}
			if ctx.currentTrack != nil {
				ctx.currentTrack.decodeOpts = append(ctx.currentTrack.decodeOpts,
					decode.FormatOptions{InArg: format.AvcIn{LengthSize: avcDcrOut.LengthSize}})
			}
		},
		"hvcC": func(ctx *decodeContext, d *decode.D) {
			_, dv := d.FieldDecode("value", mpegHEVCDCRFrameFormat)
			hevcDcrOut, ok := dv.(format.HevcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected HevcDcrOut got %#+v", dv))
			}
			if ctx.currentTrack != nil {
				ctx.currentTrack.decodeOpts = append(ctx.currentTrack.decodeOpts,
					decode.FormatOptions{InArg: format.HevcIn{LengthSize: hevcDcrOut.LengthSize}})
			}
		},
		"dfLa": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			d.FieldArrayFn("metadatablocks", func(d *decode.D) {
				for {
					_, dv := d.FieldDecode("metadatablock", flacMetadatablockFormat)
					flacMetadatablockOut, ok := dv.(format.FlacMetadatablockOut)
					if !ok {
						d.Invalid(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", dv))
					}
					if flacMetadatablockOut.HasStreamInfo {
						if ctx.currentTrack != nil {
							ctx.currentTrack.decodeOpts = append(ctx.currentTrack.decodeOpts,
								decode.FormatOptions{InArg: format.FlacFrameIn{StreamInfo: flacMetadatablockOut.StreamInfo}})
						}
					}
					if flacMetadatablockOut.IsLastBlock {
						return
					}
				}
			})
		},
		"dOps": func(ctx *decodeContext, d *decode.D) {
			d.FieldDecode("value", opusPacketFrameFormat)
		},
		"av1C": func(ctx *decodeContext, d *decode.D) {
			d.FieldDecode("value", av1CCRFormat)
		},
		"vpcC": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldDecode("value", vpxCCRFormat)
		},
		"esds": func(ctx *decodeContext, d *decode.D) {
			d.FieldU32("version")

			dataFormat := ""
			if ctx.currentTrack != nil {
				dataFormat = ctx.currentTrack.dataFormat
			}

			switch dataFormat {
			case "mp4a", "mp4v":
				_, dv := d.FieldDecode("es_descriptor", mpegESFormat)
				mpegEsOut, ok := dv.(format.MpegEsOut)
				if !ok {
					d.Invalid(fmt.Sprintf("expected mpegEsOut got %#+v", dv))
				}

				if ctx.currentTrack != nil && len(mpegEsOut.DecoderConfigs) > 0 {
					ctx.currentTrack.objectType = mpegEsOut.DecoderConfigs[0].ObjectType
				}

			default:
				d.FieldBitBufLen("data", d.BitsLeft())
			}

		},
		"stts": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("table", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldU32("count")
				d.FieldU32("duration")
				i++
			})
		},
		"stsc": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("table", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
				firstChunk := uint32(d.FieldU32("first_chunk"))
				samplesPerChunk := uint32(d.FieldU32("samples_per_chunk"))
				d.FieldU32("sample_description_id")

				if ctx.currentTrack != nil {
					ctx.currentTrack.stsc = append(ctx.currentTrack.stsc, stsc{
						firstChunk:      firstChunk,
						samplesPerChunk: samplesPerChunk,
					})
				}
				i++
			})
		},
		"stsz": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			sampleSize := d.FieldU32("sample_size")
			numEntries := d.FieldU32("num_entries")
			if sampleSize == 0 {
				var i uint64
				d.FieldArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
					size := uint32(d.FieldU32("size"))
					if ctx.currentTrack != nil {
						ctx.currentTrack.stsz = append(ctx.currentTrack.stsz, size)
					}
					i++
				})
			} else {
				if ctx.currentTrack != nil {
					for i := uint64(0); i < numEntries; i++ {
						ctx.currentTrack.stsz = append(ctx.currentTrack.stsz, uint32(sampleSize))
					}
				}
			}
		},
		"stco": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
				offset := d.FieldU32("offset")
				if ctx.currentTrack != nil {
					ctx.currentTrack.stco = append(ctx.currentTrack.stco, offset)
				}
				i++
			})
		},
		// TODO: refactor: merge with stsco?
		"co64": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
				offset := d.FieldU64("offset")
				if ctx.currentTrack != nil {
					ctx.currentTrack.stco = append(ctx.currentTrack.stco, offset)
				}
				i++
			})
		},
		"sidx": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("reference_id")
			d.FieldU32("timescale")
			if version == 0 {
				d.FieldU32("pts")
				d.FieldU32("offset")
			} else {
				d.FieldU64("pts")
				d.FieldU64("offset")
			}
			d.FieldU16("reserved")
			numEntries := d.FieldU16("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("index_table", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldU32("size")
				d.FieldU32("duration")
				d.FieldU32("sap_flags")
				i++
			})
		},
		"udta": decodeAtoms,
		"meta": func(ctx *decodeContext, d *decode.D) {
			// TODO: meta atom sometimes has a 4 byte unknown field? (flag/version?)
			maybeFlags := d.PeekBits(32)
			if maybeFlags == 0 {
				// TODO: rename?
				d.FieldU32("maybe_flags")
			}
			decodeAtoms(ctx, d)
		},
		"ilst":        decodeAtoms,
		"_apple_list": decodeAtoms,
		"_apple_entry": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldUTF8("data", int(d.BitsLeft()/8))
		},
		"data": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("reserved")
			d.FieldUTF8("data", int(d.BitsLeft()/8))
		},
		"moov": decodeAtoms,
		"moof": func(ctx *decodeContext, d *decode.D) {
			ctx.currentMoofOffset = (d.Pos() / 8) - 8
			decodeAtoms(ctx, d)
		},
		// Track Fragment
		"traf": decodeAtoms,
		// Movie Fragment Header
		"mfhd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("sequence_number")
		},
		// Track Fragment Header
		"tfhd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			baseDataOffsetPresent := false
			sampleDescriptionIndexPresent := false
			defaultSampleDurationPresent := false
			defaultSampleSizePresent := false
			defaultSampleFlagsPresent := false
			d.FieldStructFn("flags", func(d *decode.D) {
				d.FieldU7("unused0")
				d.FieldBool("duration_is_empty")
				d.FieldU10("unused1")
				defaultSampleFlagsPresent = d.FieldBool("default_sample_flags_present")
				defaultSampleSizePresent = d.FieldBool("default_sample_size_present")
				defaultSampleDurationPresent = d.FieldBool("default_sample_duration_present")
				d.FieldU1("unused2")
				sampleDescriptionIndexPresent = d.FieldBool("sample_description_index_present")
				baseDataOffsetPresent = d.FieldBool("base_data_offset_present")

			})
			trackID := uint32(d.FieldU32("track_id"))
			m := &moof{}
			ctx.currentTrack = ctx.tracks[trackID]
			if ctx.currentTrack != nil {
				ctx.currentTrack.moofs = append(ctx.currentTrack.moofs, m)
				ctx.currentTrack.currentMoof = m
			}

			if baseDataOffsetPresent {
				d.FieldU64("base_data_offset")
			}
			if sampleDescriptionIndexPresent {
				d.FieldU32("sample_description_index")
			}
			if defaultSampleDurationPresent {
				d.FieldU32("default_sample_duration")
			}
			if defaultSampleSizePresent {
				m.defaultSampleSize = uint32(d.FieldU32("default_sample_size"))
			}
			if defaultSampleFlagsPresent {
				d.FieldU32("default_sample_flags")
			}
		},
		// Track Fragment Run
		"trun": func(ctx *decodeContext, d *decode.D) {
			m := &moof{}
			if ctx.currentTrack != nil && ctx.currentTrack.currentMoof != nil {
				m = ctx.currentTrack.currentMoof
			}
			m.offset = ctx.currentMoofOffset

			d.FieldU8("version")
			sampleCompositionTimeOffsetsPresent := false
			sampleFlagsPresent := false
			sampleSizePresent := false
			sampleDurationPresent := false
			firstSampleFlagsPresent := false
			dataOffsetPresent := false
			d.FieldStructFn("flags", func(d *decode.D) {
				d.FieldU12("unused0")
				sampleCompositionTimeOffsetsPresent = d.FieldBool("sample_composition_time_sffsets_present")
				sampleFlagsPresent = d.FieldBool("sample_flags_present")
				sampleSizePresent = d.FieldBool("sample_size_present")
				sampleDurationPresent = d.FieldBool("sample_duration_present")
				d.FieldU5("unused1")
				firstSampleFlagsPresent = d.FieldBool("first_sample_flags_present")
				d.FieldU1("unused2")
				dataOffsetPresent = d.FieldBool("data_offset_present")
			})
			sampleCount := d.FieldU32("sample_count")
			if dataOffsetPresent {
				m.dataOffset = uint32(d.FieldS32("data_offset"))
			}
			if firstSampleFlagsPresent {
				d.FieldU32("first_sample_flags")
			}

			d.FieldArrayFn("samples", func(d *decode.D) {
				for i := uint64(0); i < sampleCount; i++ {
					sampleSize := m.defaultSampleSize
					d.FieldStructFn("sample", func(d *decode.D) {
						if sampleDurationPresent {
							d.FieldU32("sample_duration")
						}
						if sampleSizePresent {
							sampleSize = uint32(d.FieldU32("sample_size"))
						}
						if sampleFlagsPresent {
							d.FieldU32("sample_flags")
						}
						if sampleCompositionTimeOffsetsPresent {
							d.FieldU32("sample_composition_time_offset")
						}
					})

					m.samplesSizes = append(m.samplesSizes, sampleSize)
				}
			})
		},
		"tfdt": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			if version == 1 {
				d.FieldU64("start_time")
			} else {
				d.FieldU32("start_time")
			}
		},
		"mvex": decodeAtoms,
		"trex": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("track_id")
			d.FieldU32("default_sample_description_index")
			d.FieldU32("default_sample_duration")
			d.FieldU32("default_sample_size")
			d.FieldU4("reserved0")
			d.FieldU2("is_leading")
			d.FieldU2("sample_depends_on")
			d.FieldU2("sample_is_depended_on")
			d.FieldU2("sample_has_redundancy")
			d.FieldU3("sample_padding_value")
			d.FieldU1("sample_is_non_sync_sample")
			d.FieldU16("sample_degradation_priority")
		},
		"mfra": decodeAtoms,
		"tfra": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU26("reserved")
			lengthSizeOfTrafNum := d.FieldU2("length_size_of_traf_num")
			sampleLengthSizeOfTrunNum := d.FieldU2("sample_length_size_of_trun_num")
			lengthSizeOfSampleNum := d.FieldU2("length_size_of_sample_num")
			numEntries := d.FieldU32("number_of_entry")
			d.FieldArrayFn("entries", func(d *decode.D) {
				for i := uint64(0); i < numEntries; i++ {
					d.FieldStructFn("entry", func(d *decode.D) {
						if version == 1 {
							d.FieldU64("time")
							d.FieldU64("moof_offset")
						} else {
							d.FieldU32("time")
							d.FieldU32("moof_offset")
						}
						d.FieldU("traf_number", int(lengthSizeOfTrafNum+1)*8)
						d.FieldU("trun_number", int(sampleLengthSizeOfTrunNum+1)*8)
						d.FieldU("sample_number", int(lengthSizeOfSampleNum+1)*8)
					})
				}
			})
		},
		"mfro": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("mfra_size")
		},
		// TODO: item location
		// HEIC image
		"iloc": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			offsetSize := d.FieldU4("offset_size")
			lengthSize := d.FieldU4("length_size")
			baseOffsetSize := d.FieldU4("base_offset_size")
			d.FieldU4("reserved")
			itemCount := d.FieldU16("item_count")
			d.FieldArrayFn("items", func(d *decode.D) {
				for i := uint64(0); i < itemCount; i++ {
					d.FieldStructFn("item", func(d *decode.D) {
						d.FieldU16("id")
						d.FieldU16("data_reference_index")
						d.FieldU("base_offset", int(baseOffsetSize)*8)
						extentCount := d.FieldU16("extent_count")
						d.FieldArrayFn("extends", func(d *decode.D) {
							for i := uint64(0); i < extentCount; i++ {
								d.FieldStructFn("extent", func(d *decode.D) {
									d.FieldU("offset", int(offsetSize)*8)
									d.FieldU("length", int(lengthSize)*8)
								})
							}
						})
					})
				}
			})
		},
		"infe": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU16("id")
			d.FieldU16("protection_index")
			d.FieldStrZeroTerminated("item_name")
			// TODO: really optional? seems so
			if d.NotEnd() {
				d.FieldStrZeroTerminated("content_type")
			}
			if d.NotEnd() {
				d.FieldStrZeroTerminated("content_encoding")
			}
		},
		"iinf": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			_ = d.FieldU16("entry_count")
			decodeAtoms(ctx, d)
		},
		"iprp": decodeAtoms,
		"ipco": decodeAtoms,
	}

	typeFn := func() (string, string) {
		typ := d.UTF8(4)
		return typ, boxDescriptions[typ]
	}

	boxSize := d.U32()
	typ := d.UTF8(4)
	d.SeekRel(-8 * 8)
	var dataSize uint64
	switch boxSize {
	case 0:
		// reset of file
		// TODO: FieldU32 with display?
		d.FieldUFn("size", func() (uint64, decode.DisplayFormat, string) {
			return d.U32(), decode.NumberDecimal, "Rest of file"
		})
		d.FieldStrFn("type", typeFn)
		dataSize = uint64(d.Len()-d.Pos()) / 8
		boxSize = dataSize + 8
	case 1:
		// 64 bit length
		d.FieldUFn("size", func() (uint64, decode.DisplayFormat, string) {
			return d.U32(), decode.NumberDecimal, "Use 64 bit size"
		})
		d.FieldStrFn("type", typeFn)
		// TODO: 64 bit size zero is rest of file also?
		boxSize = d.FieldU64("size64")
		dataSize = boxSize - 16
	default:
		d.FieldU32("size")
		d.FieldStrFn("type", typeFn)
		dataSize = boxSize - 8
	}

	// TODO: not sure about this
	switch {
	case typ == "\xa9too":
		typ = "_apple_list"
	case typ[0] == 0xa9:
		typ = "_apple_entry"
	}

	if decodeFn, ok := boxes[typ]; ok {
		d.DecodeLenFn(int64(dataSize*8), func(d *decode.D) {
			decodeFn(ctx, d)
		})
	} else {
		d.FieldBitBufLen("data", int64(dataSize*8))
	}

	return boxSize
}

// TODO: ok to merge into decodeAtoms? what about terminator atom? different?
func decodeAtomsZeroTerminate(ctx *decodeContext, d *decode.D) {
	d.FieldStructArrayLoopFn("boxes", "box", d.NotEnd, func(d *decode.D) {
		// "Some sample descriptions terminate with four zero bytes that are not otherwise indicated."
		if d.BitsLeft() == 32 && d.PeekBits(32) == 0 {
			d.FieldU32("zero_terminator")
			return
		}
		decodeAtom(ctx, d)
	})
}

func decodeAtoms(ctx *decodeContext, d *decode.D) {
	d.FieldStructArrayLoopFn("boxes", "box", d.NotEnd, func(d *decode.D) {
		decodeAtom(ctx, d)
	})
}

func mp4Decode(d *decode.D, in interface{}) interface{} {
	ctx := &decodeContext{
		tracks: map[uint32]*track{},
	}

	// TODO: nicer, validate functions without field?
	d.ValidateAtLeastBytesLeft(16)
	size := d.U32()
	if size < 16 {
		d.Invalid("first box size too small < 16")
	}
	ftyp := d.UTF8(4)
	if ftyp != "ftyp" {
		d.Invalid("no ftyp box found")
	}
	d.SeekRel(-8 * 8)

	decodeAtoms(ctx, d)

	// keep track order stable
	var sortedTracks []*track
	for _, t := range ctx.tracks {
		sortedTracks = append(sortedTracks, t)
	}
	sort.Slice(sortedTracks, func(i, j int) bool { return sortedTracks[i].id < sortedTracks[j].id })

	d.FieldArrayFn("tracks", func(d *decode.D) {
		for _, t := range sortedTracks {
			decodeSampleRange := func(d *decode.D, t *track, name string, firstBit int64, nBits int64, opts ...decode.Options) {
				switch t.dataFormat {
				case "fLaC":
					d.FieldDecodeRange("sample", firstBit, nBits, flacFrameFormat, t.decodeOpts...)
				case "Opus":
					d.FieldDecodeRange("sample", firstBit, nBits, opusPacketFrameFormat, t.decodeOpts...)
				case "vp09":
					d.FieldDecodeRange("sample", firstBit, nBits, vp9FrameFormat, t.decodeOpts...)
				case "avc1":
					d.FieldDecodeRange("sample", firstBit, nBits, mpegAVCAUFormat, t.decodeOpts...)
				case "hev1":
					d.FieldDecodeRange("sample", firstBit, nBits, mpegHEVCSampleFormat, t.decodeOpts...)
				case "av01":
					d.FieldDecodeRange("sample", firstBit, nBits, av1FrameFormat, t.decodeOpts...)
				case "mp4a":
					switch t.objectType {
					case format.MPEGObjectTypeMP3:
						d.FieldDecodeRange("sample", firstBit, nBits, mp3FrameFormat, t.decodeOpts...)
					case format.MPEGObjectTypeAAC:
						// TODO: MPEGObjectTypeAACLow, Main etc?
						d.FieldDecodeRange("sample", firstBit, nBits, aacFrameFormat, t.decodeOpts...)
					case format.MPEGObjectTypeVORBIS:
						d.FieldDecodeRange("sample", firstBit, nBits, vorbisPacketFormat, t.decodeOpts...)
					default:
						d.FieldBitBufRange("sample", firstBit, nBits)
					}
				default:
					d.FieldBitBufRange("sample", firstBit, nBits)
				}
			}

			// log.Printf("t.moofs: %#+v\n", t.moofs)

			d.FieldStructFn("track", func(d *decode.D) {
				d.FieldStrFn("data_format", func() (string, string) { return t.dataFormat, "" })

				d.FieldArrayFn("samples", func(d *decode.D) {
					stscIndex := 0
					chunkNr := uint32(0)
					sampleNr := uint64(0)

					for sampleNr < uint64(len(t.stsz)) {
						stscEntry := t.stsc[stscIndex]
						sampleOffset := t.stco[chunkNr]

						for i := uint32(0); i < stscEntry.samplesPerChunk; i++ {
							sampleSize := t.stsz[sampleNr]

							decodeSampleRange(d, t, "sample", int64(sampleOffset)*8, int64(sampleSize)*8, t.decodeOpts...)

							// log.Printf("%s %d/%d %d/%d sample=%d/%d chunk=%d size=%d %d-%d\n", t.dataFormat, stscIndex, len(t.stsc), i, stscEntry.samplesPerChunk, sampleNr, len(t.stsz), chunkNr, sampleSize, sampleOffset, sampleOffset+uint64(sampleSize))

							sampleOffset += uint64(sampleSize)
							sampleNr++

						}

						chunkNr++
						if stscIndex < len(t.stsc)-1 && chunkNr >= t.stsc[stscIndex+1].firstChunk-1 {
							stscIndex++
						}
					}

					for _, m := range t.moofs {
						sampleOffset := m.offset + int64(m.dataOffset)
						for _, sz := range m.samplesSizes {
							// log.Printf("moof sample %s %d-%d\n", t.dataFormat, sampleOffset, int64(sz))

							decodeSampleRange(d, t, "sample", sampleOffset*8, int64(sz)*8, t.decodeOpts...)
							sampleOffset += int64(sz)
						}
					}
				})
			})
		}
	})

	return nil

}
