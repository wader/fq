package decoders

import (
	"fmt"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func DecodeFnForSchema(s schema.SimplifiedSchema) (func(string, *decode.D), error) {
	// TODO support logical types. Right now we will just get the raw type.
	switch s.Type {
	case schema.BOOLEAN:
		return decodeBoolFn(s)
	case schema.BYTES:
		return decodeBytesFn(s)
	case schema.DOUBLE:
		return decodeDoubleFn(s)
	case schema.ENUM:
		return decodeEnumFn(s)
	case schema.FIXED:
		return decodeFixedFn(s)
	case schema.FLOAT:
		return decodeFloatFn(s)
	case schema.INT:
		return decodeIntFn(s)
	case schema.LONG:
		return decodeLongFn(s)
	case schema.NULL:
		return decodeNullFn(s)
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
