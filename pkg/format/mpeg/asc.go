package mpeg

// https://wiki.multimedia.cx/index.php/MPEG-4_Audio

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_ASC,
		Description: "MPEG-4 Audio specific config",
		DecodeFn:    ascDecoder,
	})
}

func ascDecoder(d *decode.D) interface{} {
	d.FieldU5("object_type")
	d.FieldU4("frequency_index")
	d.FieldU4("channel_configuration")

	// TODO:
	d.FieldBitBufLen("var_aot_or_byte_align", d.BitsLeft())

	return nil
}
