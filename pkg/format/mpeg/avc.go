package mpeg

// ISO/IEC 14496-15, 5.3.3.1.2 Syntax

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"log"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC,
		Description: "H.264/AVC sample",
		DecodeFn:    avcDecode,
	})
}

func avcDecode(d *decode.D, in interface{}) interface{} {
	log.Printf("in: %#+v\n", in)
	avcIn, ok := in.(format.AvcIn)
	if !ok {
		d.Invalid("avcIn required")
	}

	// TODO: PictureLength?

	d.FieldArrayFn("nals", func(d *decode.D) {
		for d.NotEnd() {

			d.FieldStructFn("nal", func(d *decode.D) {
				l := d.FieldU("length", int(avcIn.LengthSize)*8)
				d.FieldBitBufLen("data", int64(l)*8)
			})

		}
	})

	return nil
}
