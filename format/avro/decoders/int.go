package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeIntFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	// int and long values are written using variable-length zig-zag coding.
	return func(name string, d *decode.D) {
		d.FieldSFn(name, VarZigZag)
	}, nil
}
