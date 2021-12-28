package codecs

import "github.com/wader/fq/pkg/decode"


type StringCodec struct {}

func (l StringCodec) Decode(d *decode.D) {
	length := d.FieldSFn("length", VarZigZag)
	d.FieldUTF8("value", int(length))
}

func BuildStringCodec(schema SimplifiedSchema) (Codec, error) {
	return &StringCodec{}, nil
}
