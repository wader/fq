package id3

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.ID3_V11,
		Description: "ID3v1.1 metadata",
		DecodeFn:    id3v11Decode,
	})
}

func id3v11Decode(d *decode.D, in interface{}) interface{} {
	d.ValidateAtLeastBitsLeft(128 * 8)
	d.FieldValidateUTF8("magic", "TAG+")
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

	return nil
}
