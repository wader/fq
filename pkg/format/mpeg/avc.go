package mpeg

// ISO/IEC 14496-15, 5.3.3.1.2 Syntax

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.AVC_NAL,
		Description: "H.264/AVC sample",
		DecodeFn:    avcDecode,
	})
}

func avcDecode(d *decode.D, in interface{}) interface{} {
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
