package riff

// Audio Definition Model
// https://adm.ebu.io/background/what_is_the_adm.html
// https://tech.ebu.ch/publications/tech3285s7
// https://tech.ebu.ch/publications/tech3285s5

import (
	"github.com/wader/fq/pkg/decode"
)

func chnaDecode(d *decode.D, size int64) {
	d.FieldU16("num_tracks")
	d.FieldU16("num_uids")

	audioIdLen := (size - 4) / 40
	d.FieldStructNArray("audio_ids", "audio_id", int64(audioIdLen), func(d *decode.D) {
		d.FieldU16("track_index")
		d.FieldUTF8("uid", 12)
		d.FieldUTF8("track_format_id_reference", 14)
		d.FieldUTF8("pack_format_id_reference", 11)
		d.FieldRawLen("padding", 8)
	})
}

func axmlDecode(d *decode.D, size int64) {
	// TODO(jmarnell): this chunk is all variable xml, so leave as is?
	d.FieldUTF8("xml", int(size))
}
