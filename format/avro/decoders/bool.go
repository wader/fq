package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeBoolFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	//a boolean is written as a single byte whose value is either 0 (false) or 1 (true).
	return func(name string, d *decode.D) {
		d.FieldBoolFn(name, func(d *decode.D) bool {
			return d.U8() >= 1
		})
	}, nil
}
