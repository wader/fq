package decoders

import (
	"errors"
	"fmt"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeUnionFn(schema schema.SimplifiedSchema) (DecodeFn, error) {
	if len(schema.UnionTypes) == 0 {
		return nil, errors.New("union must have types")
	}

	var decoders []func(string, *decode.D) any
	for i, t := range schema.UnionTypes {
		decodeFn, err := DecodeFnForSchema(t)
		if err != nil {
			return nil, fmt.Errorf("failed getting decodeFn for union type %d: %w", i, err)
		}
		decoders = append(decoders, decodeFn)
	}

	// A union is encoded by first writing an int value indicating the zero-based position within the union of the
	// schema of its value. The value is then encoded per the indicated schema within the union.
	return func(name string, d *decode.D) any {
		var val any
		d.FieldStruct(name, func(d *decode.D) {
			v := int(d.FieldSintFn("type", VarZigZag))
			if v < 0 || v >= len(decoders) {
				d.Fatalf("invalid union value: %d", v)
			}
			val = decoders[v]("value", d)
		})
		return val
	}, nil
}
