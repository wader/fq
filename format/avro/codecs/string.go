package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type StringCodec struct{}

func (c StringCodec) Decode(name string, d *decode.D) {
	length := d.FieldSFn(name+"_len", VarZigZag)
	d.FieldUTF8(name, int(length))
}

func BuildStringCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &StringCodec{}, nil
}
