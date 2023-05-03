package decoders

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeBoolFn(sms ...scalar.BoolMapper) (DecodeFn, error) {
	// A boolean is written as a single byte whose value is either 0 (false) or 1 (true).
	return func(name string, d *decode.D) any {
		return d.FieldBoolFn(name, func(d *decode.D) bool {
			return d.U8() >= 1
		}, sms...)
	}, nil
}
