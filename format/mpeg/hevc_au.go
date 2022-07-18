package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var hevcAUNALFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.HEVC_AU,
		Description: "H.265/HEVC Access Unit",
		DecodeFn:    hevcAUDecode,
		DecodeInArg: format.HevcAuIn{
			LengthSize: 4,
		},
		RootArray: true,
		RootName:  "access_unit",
		Dependencies: []decode.Dependency{
			{Names: []string{format.HEVC_NALU}, Group: &hevcAUNALFormat},
		},
	})
}

func hevcAUDecode(d *decode.D, in any) any {
	hevcIn, ok := in.(format.HevcAuIn)
	if !ok {
		d.Errorf("HevcAuIn required")
	}

	for d.NotEnd() {
		d.FieldStruct("nalu", func(d *decode.D) {
			l := d.FieldU("length", int(hevcIn.LengthSize)*8)
			d.FieldFormatLen("nalu", int64(l)*8, hevcAUNALFormat, nil)
		})
	}

	return nil
}
