package codecs

import (
	"fmt"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type CodecType int

type Codec interface {
	Decode(name string, d *decode.D)
}

func BuildCodec(s schema.SimplifiedSchema) (Codec, error) {
	// TODO support logical types. Right now we will just get the raw type.
	switch s.Type {
	case schema.BOOLEAN:
		return BuildBoolCodec(s)
	case schema.BYTES:
		return BuildBytesCodec(s)
	case schema.DOUBLE:
		return BuildDoubleCodec(s)
	case schema.ENUM:
		return BuildEnumCodec(s)
	case schema.FLOAT:
		return BuildFloatCodec(s)
	case schema.INT:
		return BuildIntCodec(s)
	case schema.LONG:
		return BuildLongCodec(s)
	case schema.NULL:
		return BuildNullCodec(s)
	case schema.RECORD:
		return BuildRecordCodec(s)
	case schema.STRING:
		return BuildStringCodec(s)
	case schema.UNION:
		return BuildUnionCodec(s)
	case schema.MAP:
		return BuildMapCodec(s)
	default:
		return nil, fmt.Errorf("unknown type: %s", s.Type)
	}
}
