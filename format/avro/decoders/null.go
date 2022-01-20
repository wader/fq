package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeNullFn(schema schema.SimplifiedSchema, sms ...scalar.Mapper) (DecodeFn, error) {
	// null is written as zero bytes.
	return func(name string, d *decode.D) interface{} {
		d.FieldRawLen(name, 0)
		return nil
	}, nil
}
