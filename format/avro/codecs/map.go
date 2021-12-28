package codecs

import (
	"errors"
	"fmt"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

type MapCodec struct {
	subCodec Codec
}

func (l MapCodec) Decode(name string, d *decode.D) {
	l.subCodec.Decode(name, d)
}

func BuildMapCodec(s schema.SimplifiedSchema) (Codec, error) {
	if s.Values == nil {
		return nil, errors.New("map schema must have values")
	}

	subSchema := schema.SimplifiedSchema{
		Type: schema.ARRAY,
		Items: &schema.SimplifiedSchema{
			Type: schema.RECORD,
			Fields: []schema.Field{
				{
					Name: "key",
					Type: schema.SimplifiedSchema{Type: schema.STRING},
				},
				{
					Name: "value",
					Type: *s.Values,
				},
			},
		},
	}
	subCodec, err := BuildCodec(subSchema)
	if err != nil {
		return nil, fmt.Errorf("MapCodec: %v", err)
	}

	return &MapCodec{subCodec: subCodec}, nil
}
