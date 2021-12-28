package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type IntCodec struct{}

func (l IntCodec) Decode(name string, d *decode.D) {
	// a boolean is written as a single byte whose value is either 0 (false) or 1 (true).
	d.FieldSFn(name, VarZigZag)
}

func BuildIntCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &IntCodec{}, nil
}
