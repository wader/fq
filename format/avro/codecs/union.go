package codecs

import (
	"errors"
	"fmt"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type UnionCodec struct {
	codecs []Codec
}

func (l UnionCodec) Decode(name string, d *decode.D) {
	// A union is encoded by first writing an int value indicating the zero-based position within the union of the
	// schema of its value. The value is then encoded per the indicated schema within the union.
	d.FieldStruct(name, func(d *decode.D) {
		v := int(d.FieldSFn("type", VarZigZag))
		if v >= len(l.codecs) {
			d.Fatalf("invalid union value: %d", v)
		}
		l.codecs[v].Decode("value", d)
	})
}

func BuildUnionCodec(schema schema.SimplifiedSchema) (Codec, error) {
	var c UnionCodec
	if schema.UnionTypes == nil {
		return nil, errors.New("UnionCodec: no union types")
	}
	for _, t := range schema.UnionTypes {
		tc, err := BuildCodec(t)
		if err != nil {
			return nil, fmt.Errorf("UnionCodec: %v", err)
		}
		c.codecs = append(c.codecs, tc)
	}

	return &c, nil
}
