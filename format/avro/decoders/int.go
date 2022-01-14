package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeIntFn(schema schema.SimplifiedSchema, sms ...scalar.Mapper) (DecodeFn, error) {
	// Int and long values are written using variable-length zig-zag coding.
	return func(name string, d *decode.D) interface{} {
		return d.FieldSFn(name, VarZigZag)
	}, nil
}
