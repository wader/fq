package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type FloatCodec struct{}

func (l FloatCodec) Decode(name string, d *decode.D) {
	d.FieldF32(name)
}

func BuildFloatCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &FloatCodec{}, nil
}
