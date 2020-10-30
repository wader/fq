package id3v11

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:      format.ID3V11,
		DecodeFn:  id3v1Decode,
		SkipProbe: true,
	})
}

func id3v1Decode(d *decode.D) interface{} {
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

	return nil
}
