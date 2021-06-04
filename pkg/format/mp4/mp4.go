package mp4

// Quicktime file format https://developer.apple.com/standards/qtff-2001.pdf
// FLAC in ISOBMFF https://github.com/xiph/flac/blob/master/doc/isoflac.txt
// https://www.webmproject.org/vp9/mp4/
// TODO: validate structure better? trak/stco etc
// TODO: fmp4, default samples sizes etc
// TODO: keep track of structure somehow to detect errors
// TODO: ISO-14496 says mp4 mdat can begin and end with original header/trailer (no used i guess?)
// TODO: more metadata
// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/Metadata/Metadata.html#//apple_ref/doc/uid/TP40000939-CH1-SW43
// TODO: split into mov and mp4 decoder?
// TODO: split into mp4_box decoder? needs complex in/out args?
// TODO: fragmented: tracks per fragment? fragment_index in samples?
// TODO: better probe, find first 2 boxes, should be free,ftyp or mdat?
// TODO: mime

import (
	"bytes"
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
var mpegPESPacketSampleFormat []*decode.Format
var opusPacketFrameFormat []*decode.Format
var vorbisPacketFormat []*decode.Format
var vp9FrameFormat []*decode.Format
var vpxCCRFormat []*decode.Format
var jpegFormat []*decode.Format
var id3v2Format []*decode.Format
var protoBufWidevineFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MP4,
		Description: "MPEG-4 file",
		Groups:      []string{format.PROBE},
		DecodeFn:    mp4Decode,
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
			{Names: []string{format.MPEG_HEVC_DCR}, Formats: &mpegHEVCDCRFrameFormat},
			{Names: []string{format.MPEG_HEVC_AU}, Formats: &mpegHEVCSampleFormat},
			{Names: []string{format.MPEG_PES_PACKET}, Formats: &mpegPESPacketSampleFormat},
			{Names: []string{format.OPUS_PACKET}, Formats: &opusPacketFrameFormat},
			{Names: []string{format.VORBIS_PACKET}, Formats: &vorbisPacketFormat},
			{Names: []string{format.VP9_FRAME}, Formats: &vp9FrameFormat},
			{Names: []string{format.VPX_CCR}, Formats: &vpxCCRFormat},
			{Names: []string{format.JPEG}, Formats: &jpegFormat},
			{Names: []string{format.ID3_V2}, Formats: &id3v2Format},
			{Names: []string{format.PROTOBUF_WIDEVINE}, Formats: &protoBufWidevineFormat},
		},
	})
}

type stsc struct {
	firstChunk      uint32
	samplesPerChunk uint32
}

type moof struct {
	offset                        int64
	defaultSampleSize             uint32
	defaultSampleDescriptionIndex uint32
	dataOffset                    uint32
	samplesSizes                  []uint32
}

type sampleDescription struct {
	dataFormat string
}

type track struct {
	id                 uint32
	sampleDescriptions []sampleDescription
	subType            string
	stco               []uint64 //
	stsc               []stsc
	stsz               []uint32
	decodeOpts         []decode.Options
	objectType         int // if data format is "mp4a"

	moofs       []*moof // for fmp4
	currentMoof *moof
}

type decodeContext struct {
	tracks            map[uint32]*track
	currentTrack      *track
	currentMoofOffset int64
}

func decodeBox(ctx *decodeContext, d *decode.D) {

	aliases := map[string]string{
		"styp": "ftyp",
	}

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
		"trak": decodeBoxes,
		"edts": decodeBoxes,
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
		"tref": decodeBoxes,
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
		"mdia": decodeBoxes,
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

		"minf": decodeBoxes,
		"dinf": decodeBoxes,
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
		"stbl": decodeBoxes,
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
						ctx.currentTrack.sampleDescriptions = append(ctx.currentTrack.sampleDescriptions, sampleDescription{
							dataFormat: dataFormat,
						})
						subType = ctx.currentTrack.subType
					}

					d.DecodeLenFn(int64(size-8)*8, func(d *decode.D) {
						d.FieldBytesLen("reserved", 6)
						d.FieldU16("data_reference_index")

						switch subType {
						case "soun", "vide":

							version := d.FieldU16("version")
							d.FieldU16("revision_level")
							d.FieldU32("max_packet_size") // TODO: vendor?

							// "Some sample descriptions terminate with four zero bytes that are not otherwise indicated."
							// uses decodeBoxes

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
									if d.BitsLeft() > 0 {
										decodeBoxes(ctx, d)
									}
								case 1:
									d.FieldU16("num_audio_channels")
									d.FieldU16("sample_size")
									d.FieldU16("compression_id")
									d.FieldU16("packet_size")
									d.FieldFP32("sample_rate")
									d.FieldU32("samples_per_packet")
									d.FieldU32("bytes_per_packet")
									d.FieldU32("bytes_per_frame")
									d.FieldU32("bytes_per_sample")
									if d.BitsLeft() > 0 {
										decodeBoxes(ctx, d)
									}
								case 2:
									d.FieldU16("always_3")
									d.FieldU16("always_16")
									d.FieldU16("always_minus_2")
									d.FieldU32("always_0")
									d.FieldU32("always_65536")
									d.FieldU32("size_of_struct_only")
									d.FieldF64("audio_sample_rate")
									d.FieldU32("num_audio_channels")
									d.FieldU32("always_7f000000")
									d.FieldU32("const_bits_per_channel")
									d.FieldU32("format_specific_flags")
									d.FieldU32("const_bytes_per_audio_packet")
									d.FieldU32("const_lpcm_frames_per_audio_packet")
									if d.BitsLeft() > 0 {
										decodeBoxes(ctx, d)
									}
								default:
									d.FieldBitBufLen("data", d.BitsLeft())
								}
							case "vide":
								// VideoSampleEntry
								// TODO: version 0 and 1 same?
								switch version {
								case 0, 1:
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
										decodeBoxes(ctx, d)
									}
								default:
									d.FieldBitBufLen("data", d.BitsLeft())
								}
								// case "hint": TODO: Hint entry
							default:
								d.FieldBitBufLen("data", d.BitsLeft())
							}
						default:
							d.FieldBitBufLen("data", d.BitsLeft())
						}

					})
				})
				i++
			})
		},
		"avcC": func(ctx *decodeContext, d *decode.D) {
			_, v := d.FieldDecode("value", mpegAVCDCRFormat)
			avcDcrOut, ok := v.(format.AvcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected AvcDcrOut got %#+v", v))
			}
			if ctx.currentTrack != nil {
				ctx.currentTrack.decodeOpts = append(ctx.currentTrack.decodeOpts,
					decode.FormatOptions{InArg: format.AvcIn{LengthSize: avcDcrOut.LengthSize}})
			}
		},
		"hvcC": func(ctx *decodeContext, d *decode.D) {
			_, v := d.FieldDecode("value", mpegHEVCDCRFrameFormat)
			hevcDcrOut, ok := v.(format.HevcDcrOut)
			if !ok {
				d.Invalid(fmt.Sprintf("expected HevcDcrOut got %#+v", v))
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
					_, v := d.FieldDecode("metadatablock", flacMetadatablockFormat)
					flacMetadatablockOut, ok := v.(format.FlacMetadatablockOut)
					if !ok {
						d.Invalid(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
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

			// TODO: some other way to know how to decode?
			dataFormat := ""
			if ctx.currentTrack != nil && len(ctx.currentTrack.sampleDescriptions) > 0 {
				dataFormat = ctx.currentTrack.sampleDescriptions[0].dataFormat
			}

			switch dataFormat {
			case "mp4a", "mp4v":
				_, v := d.FieldDecode("es_descriptor", mpegESFormat)
				mpegEsOut, ok := v.(format.MpegEsOut)
				if !ok {
					d.Invalid(fmt.Sprintf("expected mpegEsOut got %#+v", v))
				}

				if ctx.currentTrack != nil && len(mpegEsOut.DecoderConfigs) > 0 {
					dc := mpegEsOut.DecoderConfigs[0]
					ctx.currentTrack.objectType = dc.ObjectType
					ctx.currentTrack.decodeOpts = append(ctx.currentTrack.decodeOpts,
						decode.FormatOptions{InArg: format.AACFrameIn{ObjectType: dc.ASCObjectType}})
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
			// TODO: bytes_per_sample from audio stsd?
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
		"stss": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldArrayFn("entries", func(d *decode.D) {
				for i := uint64(0); i < numEntries; i++ {
					d.FieldU32("entry")
				}
			})
		},
		"sdtp": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			// TODO: should be count from stsz
			d.FieldArrayFn("entries", func(d *decode.D) {
				for d.NotEnd() {
					d.FieldU8("entry")
				}
			})
		},
		"ctts": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("table", "entry", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldU32("sample_count")
				d.FieldU32("composition_offset")
				i++
			})
		},
		// TODO: refactor: merge with stco?
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
		"udta": decodeBoxes,
		"meta": func(ctx *decodeContext, d *decode.D) {
			// TODO: meta box sometimes has a 4 byte unknown field? (flag/version?)
			maybeFlags := d.PeekBits(32)
			if maybeFlags == 0 {
				// TODO: rename?
				d.FieldU32("maybe_flags")
			}
			decodeBoxes(ctx, d)
		},
		"ilst":        decodeBoxes,
		"_apple_list": decodeBoxes,
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
		"moov": decodeBoxes,
		"moof": func(ctx *decodeContext, d *decode.D) {
			ctx.currentMoofOffset = (d.Pos() / 8) - 8
			decodeBoxes(ctx, d)
		},
		// Track Fragment
		"traf": decodeBoxes,
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
				m.defaultSampleDescriptionIndex = uint32(d.FieldU32("sample_description_index"))
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
		"mvex": decodeBoxes,
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
		"mfra": decodeBoxes,
		"tfra": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU26("reserved")
			lengthSizeOfTrafNum := d.FieldU2("length_size_of_traf_num")
			sampleLengthSizeOfTrunNum := d.FieldU2("sample_length_size_of_trun_num")
			lengthSizeOfSampleNum := d.FieldU2("length_size_of_sample_num")
			numEntries := d.FieldU32("num_entries")
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
			version := d.FieldU8("version")
			d.FieldU24("flags")

			offsetSize := d.FieldU4("offset_size")
			lengthSize := d.FieldU4("length_size")
			baseOffsetSize := d.FieldU4("base_offset_size")
			var indexSize uint64
			switch version {
			case 1, 2:
				indexSize = d.FieldU4("index_size")
			default:
				d.FieldU4("reserved")
			}
			var itemCount uint64
			if version < 2 {
				itemCount = d.FieldU16("item_count")
			} else {
				itemCount = d.FieldU32("item_count")
			}
			d.FieldArrayFn("items", func(d *decode.D) {
				for i := uint64(0); i < itemCount; i++ {
					d.FieldStructFn("item", func(d *decode.D) {
						switch version {
						case 0, 1:
							d.FieldU16("id")
						case 2:
							d.FieldU32("id")
						}
						switch version {
						case 1, 2:
							d.FieldU12("reserved")
							d.FieldU4("construction_method")
						}
						d.FieldU16("data_reference_index")
						d.FieldU("base_offset", int(baseOffsetSize)*8)
						extentCount := d.FieldU16("extent_count")
						d.FieldArrayFn("extends", func(d *decode.D) {
							for i := uint64(0); i < extentCount; i++ {
								d.FieldStructFn("extent", func(d *decode.D) {
									if (version == 1 || version == 2) && indexSize > 0 {
										d.FieldU("index", int(offsetSize)*8)
									}
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
			d.FieldStrNullTerminated("item_name")
			// TODO: really optional? seems so
			if d.NotEnd() {
				d.FieldStrNullTerminated("content_type")
			}
			if d.NotEnd() {
				d.FieldStrNullTerminated("content_encoding")
			}
		},
		"iinf": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			_ = d.FieldU16("entry_count")
			decodeBoxes(ctx, d)
		},
		"iprp": decodeBoxes,
		"ipco": decodeBoxes,
		"ID32": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU1("pad")
			// ISO-639-2/T as 3*5 bit intgers - 0x60
			d.FieldStrFn("langauge", func() (string, string) {
				s := ""
				for i := 0; i < 3; i++ {
					s += fmt.Sprintf("%c", int(d.U5())+0x60)
				}
				return s, ""
			})
			d.FieldDecode("data", id3v2Format)
		},
		"mehd": func(ctx *decodeContext, d *decode.D) {
			flags := d.FieldU24("flags")
			if flags&0b1 != 0 {
				d.FieldU64("fragment_duration")
			} else {
				d.FieldU32("fragment_duration")
			}
		},
		"pssh": func(ctx *decodeContext, d *decode.D) {
			var (
				systemIDCommon    = [16]byte{0x10, 0x77, 0xef, 0xec, 0xc0, 0xb2, 0x4d, 0x02, 0xac, 0xe3, 0x3c, 0x1e, 0x52, 0xe2, 0xfb, 0x4b}
				systemIDWidevine  = [16]byte{0xed, 0xef, 0x8b, 0xa9, 0x79, 0xd6, 0x4a, 0xce, 0xa3, 0xc8, 0x27, 0xdc, 0xd5, 0x1d, 0x21, 0xed}
				systemIDPlayReady = [16]byte{0x9a, 0x04, 0xf0, 0x79, 0x98, 0x40, 0x42, 0x86, 0xab, 0x92, 0xe6, 0x5b, 0xe0, 0x88, 0x5f, 0x95}
			)
			systemIDNames := map[[16]byte]string{
				systemIDCommon:    "Common",
				systemIDWidevine:  "Widevine",
				systemIDPlayReady: "PlayReady",
			}

			version := d.FieldU8("version")
			d.FieldU24("flags")
			systemID, _ := d.FieldStringUUIDMapFn("system_id", systemIDNames, "Unknown", func() []byte { return d.BytesLen(16) })
			switch version {
			case 0:
			case 1:
				kidCount := d.FieldU32("kid_count")
				d.FieldArrayFn("kids", func(d *decode.D) {
					for i := uint64(0); i < kidCount; i++ {
						d.FieldBitBufLen("kid", 16*8)
					}
				})
			}
			dataLen := d.FieldU32("data_size")

			switch {
			case bytes.Equal(systemID, systemIDWidevine[:]):
				d.FieldDecodeLen("data", int64(dataLen)*8, protoBufWidevineFormat)
			case systemID == nil:
				fallthrough
			default:
				d.FieldBitBufLen("data", int64(dataLen)*8)
			}
		},
		"sinf": decodeBoxes,
		"frma": func(ctx *decodeContext, d *decode.D) {
			d.FieldUTF8("format", 4)
		},
		"schm": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldUTF8("encryption_type", 4)
			d.FieldU16("encryption_version")
			if d.BitsLeft() > 0 {
				d.FieldUTF8("uri", int(d.BitsLeft())/8)
			}
		},
		"schi": decodeBoxes,
		"btrt": func(ctx *decodeContext, d *decode.D) {
			d.FieldU32("decoding_buffer_size")
			d.FieldU32("max_bitrate")
			d.FieldU32("avg_bitrate")
		},
		"pasp": func(ctx *decodeContext, d *decode.D) {
			d.FieldU32("h_spacing")
			d.FieldU32("v_spacing")
		},
		"uuid": func(ctx *decodeContext, d *decode.D) {
			var uuidNames = map[[16]byte]string{
				{0xa5, 0xd4, 0x0b, 0x30, 0xe8, 0x14, 0x11, 0xdd, 0xba, 0x2f, 0x08, 0x00, 0x20, 0x0c, 0x9a, 0x66}: "isml_manifest",
				{0xbe, 0x7a, 0xcf, 0xcb, 0x97, 0xa9, 0x42, 0xe8, 0x9c, 0x71, 0x99, 0x94, 0x91, 0xe3, 0xaf, 0xac}: "xmp",
				{0xff, 0xcc, 0x82, 0x63, 0xf8, 0x55, 0x4a, 0x93, 0x88, 0x14, 0x58, 0x7a, 0x02, 0x52, 0x1f, 0xdd}: "spherical",
				{0x55, 0x53, 0x4d, 0x54, 0x21, 0xd2, 0x4f, 0xce, 0xbb, 0x88, 0x69, 0x5c, 0xfa, 0xc9, 0xc7, 0x40}: "psp_usmt",
				{0x6d, 0x1d, 0x9b, 0x05, 0x42, 0xd5, 0x44, 0xe6, 0x80, 0xe2, 0x14, 0x1d, 0xaf, 0xf7, 0x57, 0xb2}: "tfxd",
				{0xd4, 0x80, 0x7e, 0xf2, 0xca, 0x39, 0x46, 0x95, 0x8e, 0x54, 0x26, 0xcb, 0x9e, 0x46, 0xa7, 0x9f}: "tfrf",
				{0x50, 0x52, 0x4f, 0x46, 0x21, 0xd2, 0x4f, 0xce, 0xbb, 0x88, 0x69, 0x5c, 0xfa, 0xc9, 0xc7, 0x40}: "prof",
				{0x6b, 0x68, 0x40, 0xf2, 0x5f, 0x24, 0x4f, 0xc5, 0xba, 0x39, 0xa5, 0x1b, 0xcf, 0x03, 0x23, 0xf3}: "ipod",
			}

			d.FieldStringUUIDMapFn("uuid", uuidNames, "Unknown", func() []byte { return d.BytesLen(16) })
			d.FieldBitBufLen("data", d.BitsLeft())
		},
		"keys": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			entryCount := d.FieldU32("entry_count")
			d.FieldArrayFn("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldStructFn("entry", func(d *decode.D) {
						keySize := d.FieldU32("key_size")
						d.FieldUTF8("key_namespace", 4)
						d.FieldUTF8("key_name", int(keySize)-8)
					})
				}
			})
		},
		"wave": decodeBoxes,
		"saiz": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			flags := d.FieldU24("flags")
			if flags&0b1 != 0 {
				d.FieldU32("aux_info_type")
				d.FieldU32("aux_info_type_parameter")
			}
			defaultSampleInfoSize := d.FieldU8("default_sample_info_size")
			sampleCount := d.FieldU32("sample_count")
			if defaultSampleInfoSize == 0 {
				d.FieldArrayFn("sample_size_info_table", func(d *decode.D) {
					for i := uint64(0); i < sampleCount; i++ {
						d.FieldU8("sample_size")
					}
				})
			}
		},
		"sgpd": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")

			// TODO: version 2?

			d.FieldU32("grouping_type")
			defaultLength := d.FieldU32("default_length")
			entryCount := d.FieldU32("entry_count")
			d.FieldArrayFn("groups", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldBitBufLen("group", int64(defaultLength)*8)
				}
			})
		},
		"sbgp": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")

			d.FieldU32("grouping_type")
			if version == 1 {
				d.FieldU32("grouping_type_parameter")
			}
			entryCount := d.FieldU32("entry_count")
			d.FieldArrayFn("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					d.FieldStructFn("entry", func(d *decode.D) {
						d.FieldU32("sample_count")
						d.FieldU32("group_description_index")
					})
				}
			})
		},
		"saio": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			flags := d.FieldU24("flags")

			if flags&0b1 != 0 {
				d.FieldU32("aux_info_type")
				d.FieldU32("aux_info_type_parameter")
			}
			entryCount := d.FieldU32("entry_count")
			d.FieldArrayFn("entries", func(d *decode.D) {
				for i := uint64(0); i < entryCount; i++ {
					if version == 0 {
						d.FieldU32("offset")
					} else {
						d.FieldU64("offset")
					}
				}
			})
		},
		"senc": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			flags := d.FieldU24("flags")

			sampleCount := d.FieldU32("sample_count")
			d.FieldArrayFn("samples", func(d *decode.D) {
				for i := uint64(0); i < sampleCount; i++ {
					d.FieldStructFn("sample", func(d *decode.D) {
						// TODO: IV_size?
						d.FieldBitBufLen("iv", 8*8)
						if flags&0b10 != 0 {
							subSampleCount := d.FieldU32("sub_sample_count")
							d.FieldArrayFn("subsamples", func(d *decode.D) {
								for i := uint64(0); i < subSampleCount; i++ {
									d.FieldStructFn("subsample", func(d *decode.D) {
										d.FieldU16("bytes_of_clear_data")
										d.FieldU32("bytes_fo_encrypted_data")
									})
								}
							})
						}
					})
				}
			})
		},
		"tenc": func(ctx *decodeContext, d *decode.D) {
			version := d.FieldU8("version")
			d.FieldU24("flags")

			d.FieldU8("reserved0")
			switch version {
			case 0:
				d.FieldU8("reserved1")
			default:
				d.FieldU4("default_crypto_bytes")
				d.FieldU4("default_skip_bytes")
			}

			defaultIsEncrypted := d.FieldU8("default_is_encrypted")
			defaultIVSize := d.FieldU8("default_iv_size")
			d.FieldBitBufLen("default_kid", 8*16)

			if defaultIsEncrypted != 0 && defaultIVSize == 0 {
				d.FieldU8("default_constant_iv_size")
			}
		},
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

	if a, ok := aliases[typ]; ok {
		typ = a
	}

	if decodeFn, ok := boxes[typ]; ok {
		d.DecodeLenFn(int64(dataSize*8), func(d *decode.D) {
			decodeFn(ctx, d)
		})
	} else {
		d.FieldBitBufLen("data", int64(dataSize*8))
	}
}

func decodeBoxes(ctx *decodeContext, d *decode.D) {
	d.FieldStructArrayLoopFn("boxes", "box", func() bool { return d.BitsLeft() >= 8*8 }, func(d *decode.D) {
		decodeBox(ctx, d)
	})

	if d.BitsLeft() > 0 {
		// "Some sample descriptions terminate with four zero bytes that are not otherwise indicated."
		if d.BitsLeft() >= 32 && d.PeekBits(32) == 0 {
			d.FieldU32("zero_terminator")
		}
		if d.BitsLeft() > 0 {
			d.FieldBitBufLen("padding", d.BitsLeft())
		}
	}
}

func mp4Decode(d *decode.D, in interface{}) interface{} {
	ctx := &decodeContext{
		tracks: map[uint32]*track{},
	}

	// TODO: nicer, validate functions without field?
	d.ValidateAtLeastBytesLeft(16)
	size := d.U32()
	if size < 8 {
		d.Invalid("first box size too small < 8")
	}
	firstType := d.UTF8(4)
	switch firstType {
	case "ftyp", "free", "moov":
	default:
		d.Invalid("no ftyp, free or moov box found")
	}

	d.SeekRel(-8 * 8)

	decodeBoxes(ctx, d)

	// keep track order stable
	var sortedTracks []*track
	for _, t := range ctx.tracks {
		sortedTracks = append(sortedTracks, t)
	}
	sort.Slice(sortedTracks, func(i, j int) bool { return sortedTracks[i].id < sortedTracks[j].id })

	d.FieldArrayFn("tracks", func(d *decode.D) {
		for _, t := range sortedTracks {
			decodeSampleRange := func(d *decode.D, t *track, dataFormat string, name string, firstBit int64, nBits int64, opts ...decode.Options) {
				switch dataFormat {
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
				case "mp4v":
					switch t.objectType {
					case format.MPEGObjectTypeMPEG2VideoMain:
						d.FieldDecodeRange("sample", firstBit, nBits, mpegPESPacketSampleFormat, t.decodeOpts...)
					case format.MPEGObjectTypeMJPEG:
						d.FieldDecodeRange("sample", firstBit, nBits, jpegFormat, t.decodeOpts...)
					default:
						d.FieldBitBufRange("sample", firstBit, nBits)
					}
				case "jpeg":
					d.FieldDecodeRange("sample", firstBit, nBits, jpegFormat, t.decodeOpts...)
				default:
					d.FieldBitBufRange("sample", firstBit, nBits)
				}
			}

			d.FieldStructFn("track", func(d *decode.D) {
				// TODO: handle progressive/fragmented mp4 differently somehow?
				if t.moofs == nil && len(t.sampleDescriptions) > 0 {
					d.FieldStrFn("data_format", func() (string, string) { return t.sampleDescriptions[0].dataFormat, "" })
				}

				d.FieldArrayFn("samples", func(d *decode.D) {
					stscIndex := 0
					chunkNr := uint32(0)
					sampleNr := uint64(0)

					for sampleNr < uint64(len(t.stsz)) {
						if int(chunkNr) >= len(t.stco) {
							// TODO: add warning
							break
						}
						stscEntry := t.stsc[stscIndex]
						sampleOffset := t.stco[chunkNr]

						for i := uint32(0); i < stscEntry.samplesPerChunk; i++ {
							if int(sampleNr) >= len(t.stsz) {
								// TODO: add warning
								break
							}

							sampleSize := t.stsz[sampleNr]
							dataFormat := "unknown"
							if len(t.sampleDescriptions) > 0 {
								dataFormat = t.sampleDescriptions[0].dataFormat
							}

							decodeSampleRange(d, t, dataFormat, "sample", int64(sampleOffset)*8, int64(sampleSize)*8, t.decodeOpts...)

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

							dataFormat := "unknown"
							if len(t.sampleDescriptions) > 0 {
								dataFormat = t.sampleDescriptions[0].dataFormat
							}
							if m.defaultSampleDescriptionIndex != 0 && int(m.defaultSampleDescriptionIndex-1) < len(t.sampleDescriptions) {
								dataFormat = t.sampleDescriptions[m.defaultSampleDescriptionIndex-1].dataFormat
							}

							// log.Printf("moof %#+v dataFormat: %#+v\n", m, dataFormat)

							decodeSampleRange(d, t, dataFormat, "sample", sampleOffset*8, int64(sz)*8, t.decodeOpts...)
							sampleOffset += int64(sz)
						}
					}
				})
			})
		}
	})

	return nil

}
