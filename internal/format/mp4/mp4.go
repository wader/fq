package mp4

// TODO: validate structure better? trak/stco etc

import (
	"fq/internal/decode"
	"log"
)

var Register = &decode.Register{
	Name: "mp4",
	MIME: "",
	New: func(common decode.Common) decode.Decoder {
		return &Decoder{
			Common: common,
			tracks: map[uint32]*track{},
		}
	},
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

// Decoder is a mp4 decoder
type Decoder struct {
	decode.Common

	tracks       map[uint32]*track
	currentTrack *track
}

func (d *Decoder) decodeAtom() uint64 {
	boxes := map[string]func(dataSize uint64){
		"ftyp": func(dataSize uint64) {
			d.FieldUTF8("major_brand", 4)
			d.FieldU32("minor_version")
			d.FieldNoneFn("brands", func() {
				numBrands := (dataSize - 8) / 4
				for i := uint64(0); i < numBrands; i++ {
					d.FieldUTF8("brand", 4)
				}
			})
		},
		"moov": d.decodeAtoms,
		"mvhd": func(dataSize uint64) {
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
		"trak": d.decodeAtoms,
		"edts": d.decodeAtoms,
		"elst": func(dataSize uint64) {
			d.FieldU8("version")
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldNoneFn("table", func() {
				for i := uint64(0); i < numEntries; i++ {
					d.FieldU32("track_duration")
					d.FieldU32("media_item")
					d.FieldFP32("media_rate")
				}
			})
		},
		"tref": d.decodeAtoms,
		"tkhd": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			d.FieldU32("creation_time")
			d.FieldU32("modification_time")
			trackID := uint32(d.FieldU32("track_id"))
			d.FieldU32("reserved")
			d.FieldU32("duration")
			d.FieldBytesLen("reserved", 8)
			d.FieldU16("layer")
			// TODO: values
			d.FieldU16("alternate_group")
			d.FieldFP16("volume")
			d.FieldU16("reserved")
			d.FieldBytesLen("matrix_structure", 36)
			d.FieldFP32("track_width")
			d.FieldFP32("track_height")

			if _, ok := d.tracks[trackID]; !ok {
				t := &track{}
				d.tracks[trackID] = t
				d.currentTrack = t
			} else {
				// TODO: dup track id?
			}
		},
		"mdia": d.decodeAtoms,
		"mdhd": func(dataSize uint64) {
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

		"hdlr": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			d.FieldUTF8("component_type", 4)
			d.FieldUTF8("component_subtype", 4)
			d.FieldUTF8("component_manufacturer", 4)
			d.FieldU32("component_flags")
			d.FieldU32("component_flags_mask")
			d.FieldUTF8("component_name", dataSize-24)
		},

		"minf": d.decodeAtoms,
		"dinf": d.decodeAtoms,
		"dref": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldNoneFn("references", func() {
				for i := uint64(0); i < numEntries; i++ {
					size := d.FieldU32("size")
					d.FieldUTF8("type", 4)
					d.FieldU8("version")
					d.FieldU24("flags")
					dataSize := size - 12
					d.FieldBytesLen("data", dataSize)
				}
			})
		},
		"stbl": d.decodeAtoms,
		"stsd": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldNoneFn("table", func() {
				for i := uint64(0); i < numEntries; i++ {
					//size := d.FieldU32("size")
					//dataFormat := d.FieldUTF8("data_format", 4)
					// d.FieldBytesLen("reserved", 6)
					// d.FieldU16("data_reference_index")
					// d.FieldU16("hint_track_version")
					// d.FieldU16("last_compatible_hint_track_version")
					// d.FieldU32("max_packet_size")
					//dataSize := size - 4 - 4
					//d.FieldBytesLen("data", dataSize)

					//d.decodeAtoms(dataSize)
					d.decodeAtom()

					// if d.currentTrack != nil {
					// 	d.currentTrack.dataFormat = dataFormat
					// }
				}
			})
		},
		"stts": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldNoneFn("table", func() {
				for i := uint64(0); i < numEntries; i++ {
					d.FieldU32("count")
					d.FieldU32("duration")
				}
			})
		},
		"stsc": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldNoneFn("table", func() {
				for i := uint64(0); i < numEntries; i++ {
					firstChunk := uint32(d.FieldU32("first_chunk"))
					samplesPerChunk := uint32(d.FieldU32("samples_per_chunk"))
					d.FieldU32("sample_description_id")

					if d.currentTrack != nil {
						d.currentTrack.stsc = append(d.currentTrack.stsc, stsc{
							firstChunk:      firstChunk,
							samplesPerChunk: samplesPerChunk,
						})
					}
				}
			})
		},
		"stsz": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			sampleSize := d.FieldU32("sample_size")
			numEntries := d.FieldU32("num_entries")
			if sampleSize == 0 {
				d.FieldNoneFn("table", func() {
					for i := uint64(0); i < numEntries; i++ {
						size := uint32(d.FieldU32("size"))

						if d.currentTrack != nil {
							d.currentTrack.stsz = append(d.currentTrack.stsz, size)
						}
					}
				})
			}
		},
		"stco": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldNoneFn("table", func() {
				for i := uint64(0); i < numEntries; i++ {
					offset := d.FieldU32("offset")

					if d.currentTrack != nil {
						d.currentTrack.stco = append(d.currentTrack.stco, offset)
					}
				}
			})
		},
		// TODO: refactor: merge with stsco?
		"co64": func(dataSize uint64) {
			d.FieldU8("version")
			// TODO: values
			d.FieldU24("flags")
			numEntries := d.FieldU32("num_entries")
			d.FieldNoneFn("table", func() {
				for i := uint64(0); i < numEntries; i++ {
					offset := d.FieldU64("offset")

					if d.currentTrack != nil {
						d.currentTrack.stco = append(d.currentTrack.stco, offset)
					}
				}
			})
		},
	}

	size := d.U32()
	typ := d.UTF8(4)
	d.SeekRel(-8 * 8)
	d.FieldNoneFn(typ, func() {
		switch size {
		case 0:
			// reset of file
			// TODO: FieldU32 with display?
			d.FieldUFn("size", func() (uint64, decode.Format, string) { return d.U32(), decode.FormatDecimal, "Rest of file" })
			d.FieldUTF8("type", 4)
			size = d.Len() - d.Pos() - (8 * 8)
		case 1:
			// 64 bit length
			d.FieldUFn("size", func() (uint64, decode.Format, string) { return d.U32(), decode.FormatDecimal, "Use 64 bit size" })
			d.FieldUTF8("type", 4)
			d.FieldU64("size64")
		default:
			d.FieldU32("size")
			d.FieldUTF8("type", 4)
		}

		dataLen := size - 8
		if decodeFn, ok := boxes[typ]; ok {
			decodeFn(dataLen)
		} else {
			d.FieldBytesLen("data", dataLen)
		}
	})

	return size
}

func (d *Decoder) decodeAtoms(bytesLeft uint64) {
	for bytesLeft > 0 {
		bytesLeft -= d.decodeAtom()
	}
}

// Decode mp4, mov, qt etc
func (d *Decoder) Decode(opts decode.Options) {
	// TODO: nicer, validate functions without field?
	d.ValidateAtLeastBitsLeft(8 * 16)
	size := d.U32()
	if size < 16 {
		d.Invalid("first box size too small < 16")
	}
	ftyp := d.UTF8(4)
	if ftyp != "ftyp" {
		d.Invalid("no ftyp box found")
	}
	d.SeekRel(-8 * 8)

	d.decodeAtoms(d.BitsLeft() / 8)

	for _, t := range d.tracks {
		d.FieldNoneFn("track", func() {

			d.FieldStrFn("data_format", func() (string, string) { return t.dataFormat, "" })

			sampleCount := uint64(0)

			for _, c := range t.stsc {

				cso := t.stco[c.firstChunk-1]

				for csi := uint32(0); csi < c.samplesPerChunk; csi++ {

					stz := uint64(t.stsz[sampleCount])

					log.Printf("cso*8: %d %#+v\n", cso, cso*8)
					log.Printf("stz*8: %d %#+v\n", stz, stz*8)

					// if t.dataFormat == "mp4a" {
					d.FieldDecodeRange("sample", cso*8, stz*8, []string{"aac_frame"})

					//} else {
					d.FieldBytesRange("sample", cso*8, stz)

					//}

					cso += stz

					sampleCount++
				}
			}
		})
	}

}
