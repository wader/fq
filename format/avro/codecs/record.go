package codecs

import (
	"fmt"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type RecordCodec struct {
	fields []string
	codecs []Codec
	name   string
}

func (l RecordCodec) Decode(name string, d *decode.D) {
	d.FieldStruct(name, func(d *decode.D) {
		for i, f := range l.fields {
			c := l.codecs[i]
			c.Decode(f, d)
		}
	})
}

func BuildRecordCodec(schema schema.SimplifiedSchema) (Codec, error) {
	var c RecordCodec
	if schema.Fields == nil {
		return c, fmt.Errorf("RecordCodec: no fields")
	}
	c.name = schema.Name

	for _, f := range schema.Fields {
		c.fields = append(c.fields, f.Name)
		fc, err := BuildCodec(f.Type)
		if err != nil {
			return c, fmt.Errorf("RecordCodec: %v", err)
		}
		c.codecs = append(c.codecs, fc)
	}
	return &c, nil
}
