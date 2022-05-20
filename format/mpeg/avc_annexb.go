package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var annexBAVCNALUFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.AVC_ANNEXB,
		Description: "H.264/AVC Annex B",
		DecodeFn: func(d *decode.D, in any) any {
			return annexBDecode(d, in, annexBAVCNALUFormat)
		},
		RootArray: true,
		RootName:  "stream",
		Dependencies: []decode.Dependency{
			{Names: []string{format.AVC_NALU}, Group: &annexBAVCNALUFormat},
		},
	})
}
