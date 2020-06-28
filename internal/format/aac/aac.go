package aac

import (
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "aac",
	MIME: "",
	New: func(common decode.Common) decode.Decoder {
		return &Decoder{
			Common: common,
		}
	},
	SkipProbe: true,
}

// Decoder is a aac decoder
type Decoder struct {
	decode.Common
}

// Decode aac
func (d *Decoder) Decode(opts decode.Options) {
	d.FieldU5("object_type")

	d.FieldStringMapFn("channels", map[uint64]string{
		0b000: "SCE",
		0b001: "CPE",
		0b010: "CCE",
		0b011: "LFE",
		0b100: "DSE",
		0b101: "PCE",
		0b110: "FIL",
		0b111: "TERM",
	}, "", d.U3)

	d.FieldU4("instance_tag")

}
