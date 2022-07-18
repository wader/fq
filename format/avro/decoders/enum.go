package decoders

import (
	"errors"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/scalar"
)

type EnumMapper struct {
	Symbols []string
}

func (e EnumMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := int(s.ActualS())
	if v < 0 || v >= len(e.Symbols) {
		return s, errors.New("enum value of out range")
	}
	s.Sym = e.Symbols[v]
	return s, nil
}

func decodeEnumFn(schema schema.SimplifiedSchema, sms ...scalar.Mapper) (DecodeFn, error) {
	if len(schema.Symbols) == 0 {
		return nil, errors.New("enum requires symbols")
	}

	// An enum is encoded by an int, representing the zero-based position of the symbol in the schema.
	// For example, consider the enum:
	//	      {"type": "enum", "name": "Foo", "symbols": ["A", "B", "C", "D"] }
	// This would be encoded by an int between zero and three, with zero indicating "A", and 3 indicating "D".
	sms = append([]scalar.Mapper{EnumMapper{Symbols: schema.Symbols}}, sms...)
	return decodeIntFn(sms...)
}
