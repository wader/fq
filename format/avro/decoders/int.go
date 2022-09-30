package decoders

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeIntFn(sms ...scalar.SintMapper) (DecodeFn, error) {
	// Int and long values are written using variable-length zig-zag coding.
	return func(name string, d *decode.D) any {
		return d.FieldSintFn(name, VarZigZag, sms...)
	}, nil
}
