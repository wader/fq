package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type BytesCodec struct{}

func (l BytesCodec) Decode(d *decode.D) interface{} {
	length := d.FieldSFn("length", VarZigZag)
	d.FieldRawLen("value", length*8)
	return nil
}

func (l BytesCodec) Type() CodecType {
	return STRUCT
}

func BuildBytesCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &BytesCodec{}, nil
}
