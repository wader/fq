package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeStringFn(schema schema.SimplifiedSchema, sms ...scalar.StrMapper) (DecodeFn, error) {
	// String is encoded as a long followed by that many bytes of UTF-8 encoded character data.
	// For example, the three-character string "foo" would be encoded as the long value 3 (encoded as hex 06) followed
	// by the UTF-8 encoding of 'f', 'o', and 'o' (the hex bytes 66 6f 6f):
	// 06 66 6f 6f
	return func(name string, d *decode.D) any {
		var val string
		d.FieldStruct(name, func(d *decode.D) {
			length := d.FieldSintFn("length", VarZigZag)
			val = d.FieldUTF8("data", int(length))
		})
		return val
	}, nil
}
