package decoders

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type BytesCodec struct{}

func decodeBytesFn(sms ...scalar.Mapper) (DecodeFn, error) {
	// Bytes are encoded as a long followed by that many bytes of data.
	return func(name string, d *decode.D) interface{} {
		var val []byte

		d.FieldStruct(name, func(d *decode.D) {
			length := d.FieldSFn("length", VarZigZag)
			r := d.FieldRawLen("data", length*8, sms...)

			val = make([]byte, length)
			if _, err := r.ReadBits(val, length*8); err != nil {
				d.Fatalf("failed to read %s bytes: %v", name, err)
			}
		})

		return val
	}, nil
}
