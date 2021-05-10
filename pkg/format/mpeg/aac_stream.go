package mpeg

// TODO: mime audio/aac?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var aacADTS []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.AAC_STREAM,
		Description: "Raw audio data transport stream",
		Groups:      []string{format.PROBE},
		DecodeFn:    adtsDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ADTS}, Formats: &aacADTS},
		},
	})
}

func adtsDecode(d *decode.D, in interface{}) interface{} {
	validFrames := 0

	d.FieldArrayFn("frames", func(d *decode.D) {
		for !d.End() {
			if dv, _, _ := d.FieldTryDecode("frame", aacADTS); dv == nil {
				break
			}
			validFrames++
		}
	})

	if validFrames == 0 {
		d.Invalid("no frames found")
	}

	return nil
}
