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

// TODO: share/refactor with avcAUDecode?
func hevcAUDecode(d *decode.D, in any) any {
	hevcIn, _ := in.(format.HevcAuIn)

	if hevcIn.LengthSize == 0 {
		// TODO: is annexb the correct name?
		annexBDecode(d, nil, hevcAUNALFormat)
		return nil
	}

	for d.NotEnd() {
		d.FieldStruct("nalu", func(d *decode.D) {
			l := int64(d.FieldU("length", int(hevcIn.LengthSize)*8)) * 8
			d.FieldFormatLen("nalu", l, hevcAUNALFormat, nil)
		})
	}

	return nil
}
