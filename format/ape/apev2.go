package ape

// http://wiki.hydrogenaud.io/index.php?title=APE_Tags_Header

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var imageFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.APEV2,
		Description: "APEv2 metadata tag",
		DecodeFn:    apev2Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.IMAGE}, Group: &imageFormat},
		},
	})
}

func apev2Decode(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian

	headerFooterFn := func(d *decode.D, name string) uint64 {
		var tagCount uint64
		d.FieldStruct(name, func(d *decode.D) {
			d.FieldUTF8("preamble", 8, d.AssertStr("APETAGEX"))
			d.FieldU32("version")
			d.FieldU32("tag_size")
			tagCount = d.FieldU32("item_count")
			d.FieldU32("flags")
			d.FieldRawLen("reserved", 64, d.BitBufIsZero)
		})
		return tagCount
	}

	tagCount := headerFooterFn(d, "header")
	if tagCount > 1000 {
		d.Fatalf("too many tags %d", tagCount)
	}

	d.FieldArray("tags", func(d *decode.D) {
		for i := uint64(0); i < tagCount; i++ {
			d.FieldStruct("tag", func(d *decode.D) {
				itemSize := d.FieldU32("item_size")
				var binaryItem bool
				d.FieldStruct("item_flags", func(d *decode.D) {
					d.FieldU6("unused0")
					binaryItem = d.FieldBool("binary")
					d.FieldU25("unused1")
				})
				keyLen := d.PeekFindByte(0, -1)
				d.FieldUTF8("key", int(keyLen))
				d.FieldU8("key_terminator")
				if binaryItem {
					d.LenFn(int64(itemSize)*8, func(d *decode.D) {
						d.FieldUTF8Null("filename")
						// assume image if binary
						dv, _, _ := d.FieldTryFormat("value", imageFormat, nil)
						if dv == nil {
							// TODO: framed and unknown instead?
							d.FieldRawLen("value", d.BitsLeft())
						}
					})
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
