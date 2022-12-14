package decoders

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeNullFn(sms ...scalar.AnyMapper) (DecodeFn, error) {
	// null is written as zero bytes.
	return func(name string, d *decode.D) any {
		d.FieldValueAny(name, nil, sms...)
		return nil
	}, nil
}
