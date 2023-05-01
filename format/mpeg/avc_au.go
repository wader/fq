package mpeg

// ISO/IEC 14496-15, 5.3.3.1.2 Syntax

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var avcNALUFormat decode.Group

func init() {
	interp.RegisterFormat(
		format.AVC_AU,
		&decode.Format{
			Description: "H.264/AVC Access Unit",
			DecodeFn:    avcAUDecode,
			DefaultInArg: format.AVC_AU_In{
				LengthSize: 0,
			},
			RootArray: true,
			RootName:  "access_unit",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AVC_NALU}, Out: &avcNALUFormat},
			},
		})
}

func avcAUDecode(d *decode.D) any {
	var ai format.AVC_AU_In
	d.ArgAs(&ai)

	if ai.LengthSize == 0 {
		// TODO: is annexb the correct name?
		annexBDecode(d, avcNALUFormat)
		return nil
	}

	for d.NotEnd() {
		d.FieldStruct("nalu", func(d *decode.D) {
			l := int64(d.FieldU("length", int(ai.LengthSize)*8)) * 8
			d.FieldFormatLen("nalu", l, &avcNALUFormat, nil)
		})
	}

	return nil
}
