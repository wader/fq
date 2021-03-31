package mpeg

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var hevcAUNALFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_HEVC_AU,
		Description: "H.265/HEVC Access Unit",
		DecodeFn:    hevcDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_HEVC_NALU}, Formats: &hevcAUNALFormat},
		},
	})
}

func hevcDecode(d *decode.D, in interface{}) interface{} {
	hevcIn, ok := in.(format.HevcIn)
	if !ok {
		d.Invalid("hevcIn required")
	}

	d.FieldArrayFn("access_unit", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStructFn("nalu", func(d *decode.D) {
				l := d.FieldU("length", int(hevcIn.LengthSize)*8)
				d.FieldDecodeLen("nalu", int64(l)*8, hevcAUNALFormat)
			})
		}
	})

	return nil
}
