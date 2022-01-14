package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeDoubleFn(schema schema.SimplifiedSchema, sms ...scalar.Mapper) (DecodeFn, error) {
	// A double is written as 8 bytes. The double is converted into a 64-bit integer using a method equivalent to Java's
	// doubleToLongBits and then encoded in little-endian format.
	return func(name string, d *decode.D) interface{} {
		return d.FieldF64(name)
	}, nil
}
