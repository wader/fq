package ape

// http://wiki.hydrogenaud.io/index.php?title=APE_Tags_Header
// TODO: havent been tested after refactor

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:     format.APEV2,
		DecodeFn: apev2Decode,
	})
}

func apev2Decode(d *decode.D) interface{} {
	headerFooterFn := func(d *decode.D, name string) uint64 {
		var tagCount uint64
		d.FieldStructFn(name, func(d *decode.D) {
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
	d.FieldArrayFn("tag", func(d *decode.D) {
		for i := uint64(0); i < tagCount; i++ {
			d.FieldStructFn("tag", func(d *decode.D) {
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
