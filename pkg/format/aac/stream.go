package aac

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var adts []*decode.Format

var Stream = format.MustRegister(&decode.Format{
	Name:  "aac_stream",
	MIMEs: []string{"audio/aac"},
	New:   func() decode.Decoder { return &StreamDecoder{} },
	Deps: []decode.Dep{
		{Names: []string{"adts"}, Formats: &adts},
	},
})

// StreamDecoder is a adts  decoder
type StreamDecoder struct {
	decode.Common
}

// Decode adts
func (d *StreamDecoder) Decode() {
	validFrames := 0
	for !d.End() {
		if _, _, errs := d.FieldTryDecode("frame", adts); errs != nil {
			break
		}
		validFrames++
	}

	if validFrames == 0 {
		d.Invalid("no frames found")
	}
}
