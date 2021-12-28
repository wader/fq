package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type BoolCodec struct{}

func (l BoolCodec) Decode(name string, d *decode.D) {
	d.FieldBoolFn(name, func(d *decode.D) bool {
		return d.U8() >= 1
	})
}

func BuildBoolCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &BoolCodec{}, nil
}
