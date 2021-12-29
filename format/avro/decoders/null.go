package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeNullFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	// null is written as zero bytes.
	return func(name string, d *decode.D) {
		// Is this the best way to represent null in fq?
		d.FieldRawLen(name, 0)
	}, nil
}
