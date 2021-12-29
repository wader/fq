package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeDoubleFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	//a double is written as 8 bytes. The double is converted into a 64-bit integer using a method equivalent to Java's
	//doubleToLongBits and then encoded in little-endian format.
	return func(name string, d *decode.D) {
		d.FieldF64(name)
	}, nil
}
