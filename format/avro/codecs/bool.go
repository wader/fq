package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)


type NullCodec struct {}

func (l NullCodec) Decode(d *decode.D) interface{}{
	// null is written as zero bytes.
	return nil
}

func (l NullCodec) Type() CodecType {
	return SCALAR
}

func BuildNullCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &NullCodec{}, nil
}
