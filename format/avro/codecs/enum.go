package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type EnumCodec struct {
	symbols []string
}

func (l EnumCodec) Decode(name string, d *decode.D) {
	d.FieldScalarStrFn(name, func(d *decode.D) string {
		value := int(VarZigZag(d))
		if value >= len(l.symbols) {
			d.Fatalf("invalid enum value: %d", value)
		}
		return l.symbols[value]
	})
}

func BuildEnumCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &EnumCodec{symbols: schema.Symbols}, nil
}
