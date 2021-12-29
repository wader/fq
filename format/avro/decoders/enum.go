package decoders

import (
	"errors"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func decodeEnumFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	if len(schema.Symbols) == 0 {
		return nil, errors.New("enum requires symbols")
	}

	//An enum is encoded by a int, representing the zero-based position of the symbol in the schema.
	//For example, consider the enum:
	//	      {"type": "enum", "name": "Foo", "symbols": ["A", "B", "C", "D"] }
	//This would be encoded by an int between zero and three, with zero indicating "A", and 3 indicating "D".
	return func(name string, d *decode.D) {
		d.FieldSFn(name, VarZigZag, scalar.Fn(func(s scalar.S) (scalar.S, error) {
			v := int(s.ActualS())
			if v < 0 || v >= len(schema.Symbols) {
				return s, errors.New("enum value of out range")
			}
			s.Sym = schema.Symbols[v]
			return s, nil
		}))
	}, nil
}
