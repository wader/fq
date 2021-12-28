package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type NullCodec struct{}

func (l NullCodec) Decode(name string, d *decode.D) {
	// null is written as zero bytes.
	d.FieldRawLen(name, 0)
}

func BuildNullCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &NullCodec{}, nil
}
