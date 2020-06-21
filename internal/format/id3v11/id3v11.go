package id3v11

import (
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "id3v11",
	MIME: "",
	New:  func(common decode.Common) decode.Decoder { return &Decoder{Common: common} },
}

// Decoder is ID3v1 decoder
type Decoder struct{ decode.Common }

// Decode ID3v1
func (d *Decoder) Decode(opts decode.Options) {
	d.ValidateAtLeastBitsLeft(128 * 8)
	d.FieldValidateString("magic", "TAG+")
	d.FieldUTF8("title", 60)
	d.FieldUTF8("artist", 60)
	d.FieldUTF8("album", 60)
	d.FieldStringMapFn("speed", map[uint64]string{
		0: "unset",
		1: "slow",
		2: "medium",
		3: "fast",
		4: "hardcore",
	}, "Unknown", d.U8)
	d.FieldUTF8("genre", 30)
	d.FieldUTF8("start", 6)
	d.FieldUTF8("stop", 6)
}
