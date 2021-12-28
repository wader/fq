package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type DoubleCodec struct{}

func (l DoubleCodec) Decode(name string, d *decode.D) {
	d.FieldF64(name)
}

func BuildDoubleCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &DoubleCodec{}, nil
}
