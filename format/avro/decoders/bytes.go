package decoders

import (
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type BytesCodec struct{}

func decodeBytesFn(sms ...scalar.BitBufMapper) (DecodeFn, error) {
	// Bytes are encoded as a long followed by that many bytes of data.
	return func(name string, d *decode.D) any {
		var val []byte

		d.FieldStruct(name, func(d *decode.D) {
			length := d.FieldSintFn("length", VarZigZag)
			br := d.FieldRawLen("data", length*8, sms...)

			val = make([]byte, length)
			if _, err := bitio.ReadFull(br, val, length*8); err != nil {
				d.Fatalf("failed to read %s bytes: %v", name, err)
			}
		})

		return val
	}, nil
}
