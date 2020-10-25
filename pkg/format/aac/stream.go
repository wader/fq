package aac

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var adts []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:     "aac_stream",
		MIMEs:    []string{"audio/aac"},
		DecodeFn: adtsDecode,
		Deps: []decode.Dep{
			{Names: []string{"adts"}, Formats: &adts},
		},
	})
}

func adtsDecode(d *decode.D) interface{} {
	validFrames := 0

	d.FieldArrayFn("frame", func(d *decode.D) {
		for !d.End() {
			if _, _, errs := d.FieldTryDecode("frame", adts); errs != nil {
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
