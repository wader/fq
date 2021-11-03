package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var annexBHEVCNALUFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.HEVC_ANNEXB,
		Description: "H.265/HEVC Annex B",
		DecodeFn: func(d *decode.D, in interface{}) interface{} {
			return annexBDecode(d, in, annexBHEVCNALUFormat)
		},
		RootArray: true,
		RootName:  "stream",
		Dependencies: []decode.Dependency{
			{Names: []string{format.HEVC_NALU}, Formats: &annexBHEVCNALUFormat},
		},
	})
}
