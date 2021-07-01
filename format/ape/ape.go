package ape

// http://wiki.hydrogenaud.io/index.php?title=APE_Tags_Header
// TODO: havent been tested after refactor

import (
	"fq/format"
	"fq/pkg/decode"
)

var imageFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.APEV2,
		Description: "APEv2 metadata tag",
		DecodeFn:    apev2Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.IMAGE}, Formats: &imageFormat},
		},
	})
}

func apev2Decode(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian

	headerFooterFn := func(d *decode.D, name string) uint64 {
		var tagCount uint64
		d.FieldStructFn(name, func(d *decode.D) {
			d.FieldValidateUTF8("preamble", "APETAGEX")
			d.FieldU32("version")
			d.FieldU32("tag_size")
			tagCount = d.FieldU32("item_count")
			d.FieldU32("flags")
			d.FieldValidateZeroPadding("reserved", 64)
		})
		return tagCount
	}

	tagCount := headerFooterFn(d, "header")
	d.FieldArrayFn("tags", func(d *decode.D) {
		for i := uint64(0); i < tagCount; i++ {
			d.FieldStructFn("tag", func(d *decode.D) {
				itemSize := d.FieldU32("item_size")
				var binaryItem bool
				d.FieldStructFn("item_flags", func(d *decode.D) {
					d.FieldU6("unused0")
					binaryItem = d.FieldBool("binary")
					d.FieldU25("unused1")
				})
				keyLen := d.PeekFindByte(0, -1)
				d.FieldUTF8("key", int(keyLen))
				d.FieldU8("key_terminator")
				if binaryItem {
					d.FieldStrNullTerminated("filename")
					// assume image if binary
					dv, _, _ := d.FieldTryDecodeLen("value", int64(itemSize-1)*8, imageFormat)
					if dv == nil {
						d.FieldBitBufLen("value", int64(itemSize-1)*8)
					}
				} else {
					d.FieldUTF8("value", int(itemSize))
				}
			})
		}
	})

	// TODO: check footer flag
	headerFooterFn(d, "footer")

	return nil
}
