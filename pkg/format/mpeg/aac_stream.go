package mpeg

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var aacADTS []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AAC_STREAM,
		Description: "Raw audio data transport stream",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"audio/aac"},
		DecodeFn:    adtsDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_ADTS}, Formats: &aacADTS},
		},
	})
}

func adtsDecode(d *decode.D, in interface{}) interface{} {
	validFrames := 0

	d.FieldArrayFn("frames", func(d *decode.D) {
		for !d.End() {
			if _, _, errs := d.FieldTryDecode("frame", aacADTS); errs != nil {
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
