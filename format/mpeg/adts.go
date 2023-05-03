package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var adtsFrameGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.ADTS,
		&decode.Format{
			Description: "Audio Data Transport Stream",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    adtsDecoder,
			RootArray:   true,
			RootName:    "frames",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.ADTS_Frame}, Out: &adtsFrameGroup},
			},
		})
}

func adtsDecoder(d *decode.D) any {
	validFrames := 0
	for !d.End() {
		if dv, _, _ := d.TryFieldFormat("frame", &adtsFrameGroup, nil); dv == nil {
			break
		}
		validFrames++
	}

	if validFrames == 0 {
		d.Fatalf("no valid frames")
	}

	return nil
}
