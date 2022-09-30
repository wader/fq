package decoders

import (
	"fmt"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type DecodeFn func(string, *decode.D) any

func DecodeFnForSchema(s schema.SimplifiedSchema) (DecodeFn, error) {
	var sms []scalar.SintMapper
	mapper := logicalTimeMapperForSchema(s)
	if mapper != nil {
		sms = append(sms, mapper)
	}

	switch s.Type {
	case schema.ARRAY:
		return decodeArrayFn(s)
	case schema.BOOLEAN:
		return decodeBoolFn()
	case schema.BYTES:
		return decodeBytesFn()
	case schema.DOUBLE:
		return decodeDoubleFn()
	case schema.ENUM:
		return decodeEnumFn(s, sms...)
	case schema.FIXED:
		return decodeFixedFn(s)
	case schema.FLOAT:
		return decodeFloatFn()
	case schema.INT:
		return decodeIntFn(sms...)
	case schema.LONG:
		return decodeLongFn(sms...)
	case schema.NULL:
		return decodeNullFn()
	case schema.RECORD:
		return decodeRecordFn(s)
	case schema.STRING:
		return decodeStringFn(s)
	case schema.UNION:
		return decodeUnionFn(s)
	case schema.MAP:
		return decodeMapFn(s)
	default:
		return nil, fmt.Errorf("unknown type: %s", s.Type)
	}
}
