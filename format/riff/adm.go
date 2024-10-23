package riff

// Audio Definition Model
// https://adm.ebu.io/background/what_is_the_adm.html
// https://tech.ebu.ch/publications/tech3285s7
// https://tech.ebu.ch/publications/tech3285s5

import (
	"github.com/wader/fq/pkg/decode"
)

func chnaDecode(d *decode.D) {
	d.FieldU16("num_tracks")
	d.FieldU16("num_uids")
	d.FieldArray("audio_ids", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("audio_id", func(d *decode.D) {
				d.FieldU16("track_index")
				d.FieldUTF8("uid", 12)
				d.FieldUTF8("track_format_id_reference", 14)
				d.FieldUTF8("pack_format_id_reference", 11)
				d.FieldRawLen("padding", 8)
			})
		}
	})
}

func axmlDecode(d *decode.D) {
	d.FieldUTF8("xml", int(d.BitsLeft())/8)
}
