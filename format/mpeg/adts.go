package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var adtsFrame decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.ADTS,
		Description: "Audio Data Transport Stream",
		Groups:      []string{format.PROBE},
		DecodeFn:    adtsDecoder,
		RootArray:   true,
		RootName:    "frames",
		Dependencies: []decode.Dependency{
			{Names: []string{format.ADTS_FRAME}, Group: &adtsFrame},
		},
	})
}

func adtsDecoder(d *decode.D) any {
	validFrames := 0
	for !d.End() {
		if dv, _, _ := d.TryFieldFormat("frame", adtsFrame, nil); dv == nil {
			break
		}
		validFrames++
	}

	if validFrames == 0 {
		d.Fatalf("no valid frames")
	}

	return nil
}
