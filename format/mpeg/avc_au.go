package mpeg

// ISO/IEC 14496-15, 5.3.3.1.2 Syntax

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var avcNALUFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.AVC_AU,
		Description: "H.264/AVC Access Unit",
		DecodeFn:    avcAUDecode,
		DecodeInArg: format.AvcAuIn{
			LengthSize: 0,
		},
		RootArray: true,
		RootName:  "access_unit",
		Dependencies: []decode.Dependency{
			{Names: []string{format.AVC_NALU}, Group: &avcNALUFormat},
		},
	})
}

func avcAUDecode(d *decode.D, in any) any {
	avcIn, _ := in.(format.AvcAuIn)

	if avcIn.LengthSize == 0 {
		// TODO: is annexb the correct name?
		annexBDecode(d, nil, avcNALUFormat)
		return nil
	}

	for d.NotEnd() {
		d.FieldStruct("nalu", func(d *decode.D) {
			l := int64(d.FieldU("length", int(avcIn.LengthSize)*8)) * 8
			d.FieldFormatLen("nalu", l, avcNALUFormat, nil)
		})
	}

	return nil
}
