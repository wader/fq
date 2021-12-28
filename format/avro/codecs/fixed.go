package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type FixedCodec struct {
	size int64
}

func (l FixedCodec) Decode(name string, d *decode.D) {
	d.FieldRawLen(name, l.size*8)
}

func BuildFixedCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &FixedCodec{size: int64(schema.Size)}, nil
}
