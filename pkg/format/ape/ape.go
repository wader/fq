package ape

// http://wiki.hydrogenaud.io/index.php?title=APE_Tags_Header

import (
	"fq/pkg/decode"
)

var TagV2 = &decode.Format{
	Name:      "apev2",
	New:       func() decode.Decoder { return &TagV2Decoder{} },
	SkipProbe: true,
}

// TagV2Decoder is APE v2 tag decoder
type TagV2Decoder struct {
	decode.Common
}

// Decode APEv2 tag
func (d *TagV2Decoder) Decode() {
	headerFooterFn := func(name string) uint64 {
		var tagCount uint64
		d.FieldNoneFn(name, func() {
			d.FieldValidateString("premble", "APETAGEX")
			d.FieldU32LE("version")
			d.FieldU32LE("tag_size")
			tagCount = d.FieldU32LE("item_count")
			d.FieldU32LE("flags")
			d.FieldValidateZeroPadding("reserved", 64)
		})
		return tagCount
	}

	tagCount := headerFooterFn("header")

	for i := uint64(0); i < tagCount; i++ {
		d.FieldNoneFn("tag", func() {
			itemSize := d.FieldU32LE("item_size")
			d.FieldU32LE("item_flags")
			keyLen := d.PeekFindByte(0, -1) - 1
			d.FieldUTF8("key", keyLen)
			d.FieldU8("key_terminator")
			d.FieldUTF8("value", itemSize)
		})
	}

	// TODO: check footer flag
	headerFooterFn("footer")
}
