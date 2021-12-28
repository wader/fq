package codecs

import (
	"errors"
	"fmt"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type ArrayCodec struct {
	valueCodec Codec
}

func (l ArrayCodec) Decode(name string, d *decode.D) {
	d.FieldArray(name, func(d *decode.D) {
		count := int64(-1)
		for count != 0 {
			d.FieldStruct(name, func(d *decode.D) {
				count = d.FieldSFn("count", VarZigZag)
				if count < 0 {
					d.FieldSFn("size", VarZigZag)
					count *= -1
				}
				d.FieldArray("entries", func(d *decode.D) {
					for i := int64(0); i < count; i++ {
						l.valueCodec.Decode("entry", d)
					}
				})
			})
		}
	})
}

func BuildArrayCodec(schema schema.SimplifiedSchema) (Codec, error) {
	if schema.Items == nil {
		return nil, errors.New("array schema must have items")
	}

	valueCodec, err := BuildCodec(*schema.Items)
	if err != nil {
		return nil, fmt.Errorf("ArrayCodec: %s", err)
	}

	return &ArrayCodec{valueCodec: valueCodec}, nil
}
