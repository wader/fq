package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type FloatCodec struct{}

func (l FloatCodec) Decode(d *decode.D) interface{} {
	return d.F32()
}

func (l FloatCodec) Type() CodecType {
	return SCALAR
}

func BuildFloatCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &FloatCodec{}, nil
}
