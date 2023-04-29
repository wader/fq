package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var annexBAVCNALUFormat decode.Group

func init() {
	interp.RegisterFormat(
		format.AvcAnnexb,
		&decode.Format{
			Description: "H.264/AVC Annex B",
			DecodeFn: func(d *decode.D) any {
				return annexBDecode(d, annexBAVCNALUFormat)
			},
			RootArray: true,
			RootName:  "stream",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AvcNalu}, Out: &annexBAVCNALUFormat},
			},
		})
}
