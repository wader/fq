package decoders

import (
	"errors"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeFixedFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	if schema.Size < 0 {
		return nil, errors.New("fixed size must be greater than or equal to zero")
	}
	size := int64(schema.Size)
	//Fixed instances are encoded using the number of bytes declared in the schema.
	return func(name string, d *decode.D) {
		d.FieldRawLen(name, size*8)
	}, nil
}
