package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)


type BoolCodec struct {}

func (l BoolCodec) Decode(d *decode.D) interface{}{
	// a boolean is written as a single byte whose value is either 0 (false) or 1 (true).
	return d.U8() != 0
}

func (l BoolCodec) Type() CodecType {
	return SCALAR
}

func BuildBoolCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &BoolCodec{}, nil
}
