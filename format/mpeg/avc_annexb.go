package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var annexBAVCNALUGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.AVC_Annexb,
		&decode.Format{
			Description: "H.264/AVC Annex B",
			DecodeFn: func(d *decode.D) any {
				return annexBDecode(d, annexBAVCNALUGroup)
			},
			RootArray: true,
			RootName:  "stream",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AVC_NALU}, Out: &annexBAVCNALUGroup},
			},
		})
}
