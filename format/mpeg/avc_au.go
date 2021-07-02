package mpeg

// ISO/IEC 14496-15, 5.3.3.1.2 Syntax

import (
	"fq/format"
	"fq/format/all/all"
	"fq/pkg/decode"
)

var avcAUNALFormat []*decode.Format

func init() {
	all.MustRegister(&decode.Format{
		Name:        format.AVC_AU,
		Description: "H.264/AVC Access Unit",
		DecodeFn:    avcAUDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AVC_NALU}, Formats: &avcAUNALFormat},
		},
	})
}

func avcAUDecode(d *decode.D, in interface{}) interface{} {
	avcIn, ok := in.(format.AvcIn)
	if !ok {
		d.Invalid("avcIn required")
	}

	d.FieldArrayFn("access_unit", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStructFn("nalu", func(d *decode.D) {
				l := d.FieldU("length", int(avcIn.LengthSize)*8)
				d.FieldDecodeLen("nalu", int64(l)*8, avcAUNALFormat)
			})
		}
	})

	return nil
}
