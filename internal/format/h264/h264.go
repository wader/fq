package h264

import (
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "h264",
	MIME: "",
	New: func(common decode.Common) decode.Decoder {
		return &Decoder{
			Common: common,
		}
	},
	SkipProbe: true,
}

// Decoder is a vorbis decoder
type Decoder struct {
	decode.Common
}

// Decode vorbis
func (d *Decoder) Decode(opts decode.Options) {

}
