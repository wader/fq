package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var adtsFrame []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.ADTS,
		Description: "Audio Data Transport Stream",
		Groups:      []string{format.PROBE},
		DecodeFn:    adtsDecoder,
		RootArray:   true,
		RootName:    "frames",
		Dependencies: []decode.Dependency{
			{Names: []string{format.ADTS_FRAME}, Formats: &adtsFrame},
		},
	})
}

func adtsDecoder(d *decode.D, in interface{}) interface{} {
	validFrames := 0
	for !d.End() {
		if dv, _, _ := d.FieldTryFormat("frame", adtsFrame, nil); dv == nil {
			break
		}
		validFrames++
	}

	if validFrames == 0 {
		d.Fatal("no valid frames")
	}

	return nil
}
