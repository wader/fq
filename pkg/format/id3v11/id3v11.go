package id3v11

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var Tag = format.MustRegister(&decode.Format{
	Name:      "id3v11",
	New:       func() decode.Decoder { return &TagDecoder{} },
	SkipProbe: true,
})

// TagDecoder is ID3v11 tag decoder
type TagDecoder struct{ decode.Common }

// Decode ID3v1
func (d *TagDecoder) Decode() {
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
