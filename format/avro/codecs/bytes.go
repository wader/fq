package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type BytesCodec struct{}

func (l BytesCodec) Decode(name string, d *decode.D) {
	// What if its a record with a field called name_len?
	// using a struct is probably a better idea. But it makes it less usable
	length := d.FieldSFn(name+"_len", VarZigZag)
	d.FieldRawLen("name", length*8)
}

func BuildBytesCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &BytesCodec{}, nil
}
