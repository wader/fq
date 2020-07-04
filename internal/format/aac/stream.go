package aac

import (
	"fq/internal/decode"
)

var Stream = &decode.Register{
	Name: "aac_stream",
	New:  func() decode.Decoder { return &StreamDecoder{} },
}

// StreamDecoder is a adts  decoder
type StreamDecoder struct {
	decode.Common
}

// Decode adts
func (d *StreamDecoder) Decode() {
	for !d.End() {
		if !d.FieldDecode("frame", []string{"adts"}) {
			break
		}
	}
}
