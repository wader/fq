package ape

// http://wiki.hydrogenaud.io/index.php?title=APE_Tags_Header
// TODO: havent been tested after refactor

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:      "apev2",
		DecodeFn:  apev2Decode,
		SkipProbe: true,
	})
}

func apev2Decode(d *decode.Common) interface{} {
	headerFooterFn := func(d *decode.Common, name string) uint64 {
		var tagCount uint64
		d.FieldStructFn2(name, func(d *decode.Common) {
			d.FieldValidateString("premble", "APETAGEX")
			d.FieldU32LE("version")
			d.FieldU32LE("tag_size")
			tagCount = d.FieldU32LE("item_count")
			d.FieldU32LE("flags")
			d.FieldValidateZeroPadding("reserved", 64)
		})
		return tagCount
	}

	tagCount := headerFooterFn(d, "header")
	d.FieldArrayFn2("tag", func(d *decode.Common) {
		for i := uint64(0); i < tagCount; i++ {
			d.FieldStructFn2("tag", func(d *decode.Common) {
				itemSize := d.FieldU32LE("item_size")
				d.FieldU32LE("item_flags")
				keyLen := d.PeekFindByte(0, -1) - 1
				d.FieldUTF8("key", keyLen)
				d.FieldU8("key_terminator")
				d.FieldUTF8("value", int64(itemSize))
			})
		}
	})

	// TODO: check footer flag
	headerFooterFn(d, "footer")

	return nil
}
