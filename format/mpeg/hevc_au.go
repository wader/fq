package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var hevcAUNALFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.HEVC_AU,
		Description: "H.265/HEVC Access Unit",
		DecodeFn:    hevcAUDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.HEVC_NALU}, Formats: &hevcAUNALFormat},
		},
	})
}

func hevcAUDecode(d *decode.D, in interface{}) interface{} {
	hevcIn, ok := in.(format.HevcIn)
	if !ok {
		d.Invalid("hevcIn required")
	}

	d.FieldArrayFn("access_unit", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStructFn("nalu", func(d *decode.D) {
				l := d.FieldU("length", int(hevcIn.LengthSize)*8)
				d.FieldFormatLen("nalu", int64(l)*8, hevcAUNALFormat, nil)
			})
		}
	})

	return nil
}
