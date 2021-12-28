package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type EnumCodec struct{
	symbols []string
}

func (l EnumCodec) Decode(d *decode.D) interface{} {
	value := int(VarZigZag(d))
	if value >= len(l.symbols) {
		d.Fatalf("invalid enum value: %d", value)
	}
	return l.symbols[value]
}

func (l EnumCodec) Type() CodecType {
	return SCALAR
}

func BuildEnumCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &EnumCodec{symbols: schema.Symbols}, nil
}
