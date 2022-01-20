package decoders

import (
	"errors"
	"fmt"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeMapFn(s schema.SimplifiedSchema) (DecodeFn, error) {
	if s.Values == nil {
		return nil, errors.New("map schema must have values")
	}

	// Maps are encoded as a series of blocks. Each block consists of a long count value, followed by that many
	// key/value pairs. A block with count zero indicates the end of the map. Each item is encoded per the map's
	// value schema.
	// If a block's count is negative, its absolute value is used, and the count is followed immediately by a long
	// block size indicating the number of bytes in the block. This block size permits fast skipping through data,
	// e.g., when projecting a record to a subset of its fields.
	// The blocked representation permits one to read and write maps larger than can be buffered in memory, since one
	// can start writing items without knowing the full length of the map.
	// This is the exact same as the array decoder, with the value being a KV record, so we just use the array decoder

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
	subFn, err := DecodeFnForSchema(subSchema)
	if err != nil {
		return nil, fmt.Errorf("decode map: %w", err)
	}
	return func(s string, d *decode.D) interface{} {
		val := make(map[string]interface{})

		rawV := subFn(s, d)
		impl, ok := rawV.([]map[string]interface{})
		if !ok {
			d.Fatalf("decode map: expected array of maps, got %T", rawV)
			return nil
		}
		for entry := range impl {
			rawKey, ok := impl[entry]["key"]
			if !ok {
				d.Fatalf("decode map: expected key in map %v", impl[entry])
			}
			value, ok := impl[entry]["key"]
			if !ok {
				d.Fatalf("decode map: expected value in map %v", impl[entry])
			}
			key, ok := rawKey.(string)
			if !ok {
				d.Fatalf("decode map: expected string key in map %v", impl[entry])
			}
			val[key] = value
		}
		return val
	}, nil
}
