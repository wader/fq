package decoders

import (
	"errors"

	"github.com/wader/fq/pkg/scalar"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeFixedFn(schema schema.SimplifiedSchema, sms ...scalar.BitBufMapper) (DecodeFn, error) {
	if schema.Size < 0 {
		return nil, errors.New("fixed size must be greater than or equal to zero")
	}
	size := int64(schema.Size)
	// Fixed instances are encoded using the number of bytes declared in the schema.
	return func(name string, d *decode.D) any {
		r := d.FieldRawLen(name, size*8, sms...)
		val := make([]byte, size)
		if _, err := r.ReadBits(val, size*8); err != nil {
			d.Fatalf("failed to read fixed %s value: %v", name, err)
		}
		return val
	}, nil
}
