package decoders

import (
	"fmt"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type DecodeFn func(string, *decode.D) interface{}

func DecodeFnForSchema(s schema.SimplifiedSchema) (DecodeFn, error) {
	var sms []scalar.Mapper
	mapper := logicalMapperForSchema(s)
	if mapper != nil {
		sms = append(sms, mapper)
	}

	switch s.Type {
	case schema.ARRAY:
		return decodeArrayFn(s, sms...)
	case schema.BOOLEAN:
		return decodeBoolFn(s, sms...)
	case schema.BYTES:
		return decodeBytesFn(s, sms...)
	case schema.DOUBLE:
		return decodeDoubleFn(s, sms...)
	case schema.ENUM:
		return decodeEnumFn(s, sms...)
	case schema.FIXED:
		return decodeFixedFn(s, sms...)
	case schema.FLOAT:
		return decodeFloatFn(s, sms...)
	case schema.INT:
		return decodeIntFn(s, sms...)
	case schema.LONG:
		return decodeLongFn(s, sms...)
	case schema.NULL:
		return decodeNullFn(s, sms...)
	case schema.RECORD:
		return decodeRecordFn(s, sms...)
	case schema.STRING:
		return decodeStringFn(s, sms...)
	case schema.UNION:
		return decodeUnionFn(s, sms...)
	case schema.MAP:
		return decodeMapFn(s, sms...)
	default:
		return nil, fmt.Errorf("unknown type: %s", s.Type)
	}
}
