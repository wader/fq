package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeStringFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	//string is encoded as a long followed by that many bytes of UTF-8 encoded character data.
	//For example, the three-character string "foo" would be encoded as the long value 3 (encoded as hex 06) followed
	//by the UTF-8 encoding of 'f', 'o', and 'o' (the hex bytes 66 6f 6f):
	//06 66 6f 6f
	return func(name string, d *decode.D) {
		length := d.FieldSFn(name+"_len", VarZigZag)
		d.FieldUTF8(name, int(length))
	}, nil
}
