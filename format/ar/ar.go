package ar

import (
	"strconv"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var probeFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.AR,
		Description: "Unix archive",
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeAr,
		Dependencies: []decode.Dependency{
			{Names: []string{format.PROBE}, Group: &probeFormat},
		},
	})
}

func decodeAr(d *decode.D, in interface{}) interface{} {
	d.FieldUTF8("signature", 8, d.AssertStr("!<arch>\n"))
	d.FieldArray("files", func(d *decode.D) {
		// TODO: extract? share with tar?
		mapStrToSymU := func(base int) scalar.Mapper {
			return scalar.Fn(func(s scalar.S) (scalar.S, error) {
				ts := strings.Trim(s.ActualStr(), " ")
				if ts != "" {
					n, err := strconv.ParseUint(ts, base, 64)
					if err != nil {
						return s, err
					}
					s.Sym = n
				}
				return s, nil
			})
		}

		for !d.End() {
			d.FieldStruct("file", func(d *decode.D) {
				d.FieldUTF8("identifier", 16, scalar.TrimSpace)
				d.FieldUTF8("modification_timestamp", 12, scalar.TrimSpace, mapStrToSymU(10))
				d.FieldUTF8("owner_id", 6, scalar.TrimSpace, mapStrToSymU(10))
				d.FieldUTF8("group_id", 6, scalar.TrimSpace, mapStrToSymU(10))
				d.FieldUTF8("file_mode", 8, scalar.TrimSpace, mapStrToSymU(8)) // Octal
				sizeS := d.FieldScalarUTF8("file_size", 10, scalar.TrimSpace, mapStrToSymU(10))
				if sizeS.Sym == nil {
					d.Fatalf("could not decode file_size")
				}
				size := int64(sizeS.SymU()) * 8
				d.FieldUTF8("ending_characters", 2)
				dv, _, _ := d.TryFieldFormatLen("data", size, probeFormat, nil)
				if dv == nil {
					d.FieldRawLen("data", size)
				}
				padding := d.AlignBits(16)
				if padding > 0 {
					d.FieldRawLen("padding", int64(padding))
				}
			})
		}
	})

	return nil
}
