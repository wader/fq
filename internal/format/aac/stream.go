package aac

import (
	"fq/internal/decode"
)

var Stream = &decode.Format{
	Name: "aac_stream",
	New:  func() decode.Decoder { return &StreamDecoder{} },
}

// StreamDecoder is a adts  decoder
type StreamDecoder struct {
	decode.Common
}

// Decode adts
func (d *StreamDecoder) Decode() {
	validFrames := 0
	for !d.End() {
		if !d.FieldDecode("frame", ADTS) {
			break
		}
		validFrames++
	}

	if validFrames == 0 {
		d.Invalid("no frames found")
	}
}
