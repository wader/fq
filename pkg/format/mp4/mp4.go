package mp4

// TODO: validate structure better? trak/stco etc
// TODO: rename atom -> box?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"log"
	"strings"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:   format.MP4,
		Groups: []string{format.PROBE},
		// TODO: implment MIME()
		MIMEs:    []string{"audio/mp4", "video/mp4"},
		DecodeFn: mp4Decode,
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
			d.FieldStructArrayLoopFn("brands", func() bool { return i < numBrands }, func(d *decode.D) {
				d.FieldStrFn("brand", func() (string, string) {
					return strings.TrimSpace(d.UTF8(4)), ""
				})
				i++
			})
		},
		"moov": decodeAtoms,
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
				d.FieldU32("media_item")
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
			d.FieldUTF8("component_name", int64(d.BitsLeft()/8))
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
			d.FieldStructArrayLoopFn("reference", func() bool { return i < numEntries }, func(d *decode.D) {
				//size := d.FieldU32("size")
				//dataFormat := d.FieldUTF8("data_format", 4)
				// d.FieldBytesLen("reserved", 6)
				// d.FieldU16("data_reference_index")
				// d.FieldU16("hint_track_version")
				// d.FieldU16("last_compatible_hint_track_version")
				// d.FieldU32("max_packet_size")
				//dataSize := size - 4 - 4
				//d.FieldBytesLen("data", dataSize)

				//decodeAtoms(dataSize)
				decodeAtom(ctx, d)

				// if d.currentTrack != nil {
				// 	d.currentTrack.dataFormat = dataFormat
				// }
				i++
			})
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

	log.Printf("dataSize: %d\n", dataSize)

	if decodeFn, ok := boxes[typ]; ok {
		d.SubLenFn(int64(dataSize*8), func(d *decode.D) { decodeFn(ctx, d) })
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

	log.Println("BLA")

	// for _, t := range d.tracks {
	// 	d.FieldNoneFn("track", func() {

	// 		d.FieldStrFn("data_format", func() (string, string) { return t.dataFormat, "" })

	// 		sampleCount := uint64(0)

	// 		for _, c := range t.stsc {

	// 			cso := t.stco[c.firstChunk-1]

	// 			for csi := uint32(0); csi < c.samplesPerChunk; csi++ {

	// 				stz := uint64(t.stsz[sampleCount])

	// 				// log.Printf("cso*8: %d %#+v\n", cso, cso*8)
	// 				// log.Printf("stz*8: %d %#+v\n", stz, stz*8)

	// 				// if t.dataFormat == "mp4a" {
	// 				d.FieldDecodeRange("sample", int64(cso)*8, int64(stz)*8, aac.Frame)

	// 				//} else {
	// 				d.FieldBytesRange("sample", int64(cso)*8, int64(stz))

	// 				//}

	// 				cso += stz

	// 				sampleCount++

	// 				log.Printf("SAMPLE %d %d", csi, c.samplesPerChunk)
	// 			}

	// 			log.Println("ATTTTTT1")

	// 		}

	// 		log.Println("ATTTTTT2")

	// 	})
	// }

	//log.Println("BLA2")

	return nil

}
