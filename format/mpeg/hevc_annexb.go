package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var annexBHEVCNALUFormat decode.Group

func init() {
	interp.RegisterFormat(
		format.HevcAnnexb,
		&decode.Format{
			Description: "H.265/HEVC Annex B",
			DecodeFn: func(d *decode.D) any {
				return annexBDecode(d, annexBHEVCNALUFormat)
			},
			RootArray: true,
			RootName:  "stream",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.HevcNalu}, Out: &annexBHEVCNALUFormat},
			},
		})
}
