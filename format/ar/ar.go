package ar

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var probeGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.AR,
		&decode.Format{
			Description: "Unix archive",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeAr,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Probe}, Out: &probeGroup},
			},
		})
}

func decodeAr(d *decode.D) any {
	d.FieldUTF8("signature", 8, d.StrAssert("!<arch>\n"))
	d.FieldArray("files", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("file", func(d *decode.D) {
				d.FieldUTF8("identifier", 16, scalar.ActualTrimSpace)
				d.FieldUTF8("modification_timestamp", 12, scalar.ActualTrimSpace, scalar.TryStrSymParseUint(10))
				d.FieldUTF8("owner_id", 6, scalar.ActualTrimSpace, scalar.TryStrSymParseUint(10))
				d.FieldUTF8("group_id", 6, scalar.ActualTrimSpace, scalar.TryStrSymParseUint(10))
				d.FieldUTF8("file_mode", 8, scalar.ActualTrimSpace, scalar.TryStrSymParseUint(8)) // Octal
				sizeStr := d.FieldScalarUTF8("file_size", 10, scalar.ActualTrimSpace, scalar.TryStrSymParseUint(10))
				if sizeStr.Sym == nil {
					d.Fatalf("could not decode file_size")
				}
				size := int64(sizeStr.SymUint()) * 8
				d.FieldUTF8("ending_characters", 2)
				d.FieldFormatOrRawLen("data", size, &probeGroup, format.Probe_In{})
				padding := d.AlignBits(16)
				if padding > 0 {
					d.FieldRawLen("padding", int64(padding))
				}
			})
		}
	})

	return nil
}
