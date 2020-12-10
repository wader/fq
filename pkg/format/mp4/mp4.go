package mp4

// TODO: validate structure better? trak/stco etc
// TODO: rename atom -> box?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"strings"
)

var mpegESFormat []*decode.Format
var aacFrameFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MP4,
		Description: "MP4 container",
		Groups:      []string{format.PROBE},
		// TODO: implment MIME()
		MIMEs:    []string{"audio/mp4", "video/mp4"},
		DecodeFn: mp4Decode,
		Deps: []decode.Dep{
			{Names: []string{format.MPEG_ES}, Formats: &mpegESFormat},
			{Names: []string{format.AAC_FRAME}, Formats: &aacFrameFormat},
		},
	})
}

type stsc struct {
	firstChunk      uint32
	samplesPerChunk uint32
}

type track struct {
	id         uint32
	dataFormat string
	stco       []uint64 //
	stsc       []stsc
	stsz       []uint32
}

type decodeContext struct {
	tracks       map[uint32]*track
	currentTrack *track
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
			d.FieldStructArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
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
				t := &track{}
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
			d.FieldUTF8("component_subtype", 4)
			d.FieldUTF8("component_manufacturer", 4)
			d.FieldU32("component_flags")
			d.FieldU32("component_flags_mask")
			d.FieldUTF8("component_name", int(d.BitsLeft()/8))
		},

		"minf": decodeAtoms,
		"dinf": decodeAtoms,
		"dref": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("reference", func() bool { return i < numEntries }, func(d *decode.D) {
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
			d.FieldStructArrayLoopFn("sample_description", func() bool { return i < numEntries }, func(d *decode.D) {
				size := d.FieldU32("size")
				dataFormat := d.FieldUTF8("data_format", 4)
				if ctx.currentTrack != nil {
					ctx.currentTrack.dataFormat = dataFormat
				}
				d.FieldBytesLen("reserved", 6)
				d.FieldU16("data_reference_index")

				switch dataFormat {
				case "mp4a":
					// https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFChap3/qtff3.html#//apple_ref/doc/uid/TP40000939-CH205-SW1
					d.FieldStructFn("data", func(d *decode.D) {
						switch d.FieldU16("version") {
						case 0:
							d.FieldU16("revision_level")
							d.FieldU32("vendor")
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
						}

						// TODO: check for extra 4 zero bytes optionally included in size

						if d.BitsLeft() > 0 {
							decodeAtoms(ctx, d)
						}
					})
				case "avc1":
					d.FieldStructFn("data", func(d *decode.D) {

						d.FieldU16("version")
						d.FieldU16("revision_level")
						d.FieldU32("vendor")
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
							decodeAtoms(ctx, d)
						}
					})

				default:
					d.FieldBytesLen("data", int(size)-16)
				}

				// "Some sample descriptions terminate with four zero bytes that are not otherwise indicated."
				if d.BitsLeft() >= 4 && d.PeekBits(32) == 0 {
					d.FieldU32("zero_terminator")
				}

				i++
			})
		},
		"avcC": func(ctx *decodeContext, d *decode.D) {
			d.FieldBitBufLen("data", d.BitsLeft())
		},
		"esds": func(ctx *decodeContext, d *decode.D) {
			d.FieldU32("version")

			dataFormat := ""
			if ctx.currentTrack != nil {
				dataFormat = ctx.currentTrack.dataFormat
			}
			switch dataFormat {
			case "mp4a":
				//d.FieldDecode("es_desciptor", mpegESFormat)
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
			d.FieldStructArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
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
			d.FieldStructArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
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
				d.FieldStructArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
					size := uint32(d.FieldU32("size"))
					if ctx.currentTrack != nil {
						ctx.currentTrack.stsz = append(ctx.currentTrack.stsz, size)
					}
					i++
				})
			}
		},
		"stco": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			var i uint64
			d.FieldStructArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
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
			d.FieldStructArrayLoopFn("table", func() bool { return i < numEntries }, func(d *decode.D) {
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
			d.FieldStructArrayLoopFn("index_table", func() bool { return i < numEntries }, func(d *decode.D) {
				d.FieldU32("size")
				d.FieldU32("duration")
				d.FieldU32("sap_flags")
				i++
			})
		},
		"udta": decodeAtoms,
		"meta": func(ctx *decodeContext, d *decode.D) {
			// TODO: meta atom sometimes has a 4 byte unknown field? (flag/version?)
			unknown := d.PeekBits(32)
			if unknown == 0 {
				// TODO: rename?
				d.FieldU32("unknown")
			}
			decodeAtoms(ctx, d)
		},
		"ilst":            decodeAtoms,
		"_apple_ilst_box": decodeAtoms,
		"data": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("reserved")
			d.FieldUTF8("data", int(d.BitsLeft()/8))
		},
		"moov": decodeAtoms,
		"moof": decodeAtoms,
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
			d.FieldU32("track_id")
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
				d.FieldU32("default_sample_size")
			}
			if defaultSampleFlagsPresent {
				d.FieldU32("default_sample_flags")
			}
		},
		// Track Fragment Run
		"trun": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			sampleCompositionTimeOffsetsPresent := false
			sampleFlagsPresent := false
			sampleSizePresent := false
			sampleDurationPresent := false
			firstSampleFlagsPresent := false
			dataOffsetPresent := false
			d.FieldStructFn("flags", func(d *decode.D) {
				d.FieldU12("unused0")
				sampleCompositionTimeOffsetsPresent = d.FieldBool("sampleCompositionTimeOffsetsPresent")
				sampleFlagsPresent = d.FieldBool("sample_flags_present")
				sampleSizePresent = d.FieldBool("sample_size_present")
				sampleDurationPresent = d.FieldBool("sample_duration_present")
				d.FieldU5("unused1")
				firstSampleFlagsPresent = d.FieldBool("first_sample_flags_present")
				d.FieldU1("unused2")
				dataOffsetPresent = d.FieldBool("dataOffset_present")
			})
			sampleCount := d.FieldU32("sample_count")
			if dataOffsetPresent {
				d.FieldS32("data_offset")
			}
			if firstSampleFlagsPresent {
				d.FieldU32("first_sample_flags")
			}
			d.FieldArrayFn("sample", func(d *decode.D) {
				for i := uint64(0); i < sampleCount; i++ {
					d.FieldStructFn("sample", func(d *decode.D) {
						if sampleDurationPresent {
							d.FieldU32("sample_duration")
						}
						if sampleSizePresent {
							d.FieldU32("sample_size")
						}
						if sampleFlagsPresent {
							d.FieldU32("sample_flags")
						}
						if sampleCompositionTimeOffsetsPresent {
							d.FieldU32("sample_composition_time_offset")
						}
					})
				}
			})
		},
		"tfdt": func(ctx *decodeContext, d *decode.D) {
			d.FieldU8("version")
			d.FieldU24("flags")
			d.FieldU32("start_time")
		},
	}

	boxSize := d.U32()
	typ := d.UTF8(4)
	d.SeekRel(-8 * 8)
	var dataSize uint64
	switch boxSize {
	case 0:
		// reset of file
		// TODO: FieldU32 with display?
		d.FieldUFn("size", func() (uint64, decode.DisplayFormat, string) { return d.U32(), decode.NumberDecimal, "Rest of file" })
		d.FieldUTF8("type", 4)
		dataSize = uint64(d.Len()-d.Pos()) / 8
		boxSize = dataSize + 8
	case 1:
		// 64 bit length
		d.FieldUFn("size", func() (uint64, decode.DisplayFormat, string) { return d.U32(), decode.NumberDecimal, "Use 64 bit size" })
		d.FieldUTF8("type", 4)
		boxSize = d.FieldU64("size64")
		dataSize = boxSize - 16
	default:
		d.FieldU32("size")
		d.FieldUTF8("type", 4)
		dataSize = boxSize - 8
	}

	if typ[0] == 0xa9 {
		typ = "_apple_ilst_box"
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

func decodeAtoms(ctx *decodeContext, d *decode.D) {
	d.FieldArrayFn("box", func(d *decode.D) {
		for !d.End() {
			d.FieldStructFn("box", func(d *decode.D) {
				decodeAtom(ctx, d)
			})
		}
	})
}

func mp4Decode(d *decode.D) interface{} {
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

	//log.Println("BLA")

	d.FieldArrayFn("track", func(d *decode.D) {
		for _, t := range ctx.tracks {
			d.FieldStructFn("track", func(d *decode.D) {
				d.FieldStrFn("data_format", func() (string, string) { return t.dataFormat, "" })

				sampleCount := uint64(0)

				d.FieldArrayFn("sample", func(d *decode.D) {
					for _, c := range t.stsc {

						cso := t.stco[c.firstChunk-1]

						for csi := uint32(0); csi < c.samplesPerChunk; csi++ {

							stz := uint64(t.stsz[sampleCount])

							// log.Printf("cso*8: %d %#+v\n", cso, cso*8)
							// log.Printf("stz*8: %d %#+v\n", stz, stz*8)

							// if t.dataFormat == "mp4a" {
							//d.FieldDecodeRange("sample", int64(cso)*8, int64(stz)*8, aac.Frame)

							//} else {
							//							d.FieldBytesRange("sample", int64(cso)*8, int(stz))

							// d.FieldStructFn("sample", func(d *decode.D) {

							// 	// d.DecodeRangeFn(int64(cso)*8, int(stz), func(d *decode.D) {
							// 	// 	d.FieldBool("c1")
							// 	// 	d.FieldBool("c2")
							// 	// 	d.FieldU13("fl")
							// 	// 	d.FieldU13("fl")
							// 	// }

							// })

							//}

							d.FieldBitBufRange("sample", int64(cso)*8, int64(stz)*8)

							//d.FieldDecodeRange("sample", int64(cso)*8, int64(stz)*8, aacFrameFormat)

							cso += stz

							sampleCount++
						}
					}
				})
			})
		}
	})

	return nil

}
