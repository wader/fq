package mpeg

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.HEVC_NAL,
		Description: "H.265/HEVC sample",
		DecodeFn:    hevcDecode,
	})
}

func hevcDecode(d *decode.D, in interface{}) interface{} {
	hevcIn, ok := in.(format.HevcIn)
	if !ok {
		d.Invalid("avcIn required")
	}

	// TODO: PictureLength?

	d.FieldArrayFn("nals", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStructFn("nal", func(d *decode.D) {
				l := d.FieldU("length", int(hevcIn.LengthSize)*8)
				d.FieldBitBufLen("data", int64(l)*8)
			})
		}
	})

	return nil
}
