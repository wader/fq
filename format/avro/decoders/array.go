package decoders

import (
	"errors"
	"fmt"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeArrayFn(schema schema.SimplifiedSchema) (DecodeFn, error) {
	if schema.Items == nil {
		return nil, errors.New("array schema must have items")
	}

	valueD, err := DecodeFnForSchema(*schema.Items)
	if err != nil {
		return nil, fmt.Errorf("failed getting decode fn for array item: %w", err)
	}

	// Arrays are encoded as a series of blocks. Each block consists of a long count value, followed by that many array
	// items. A block with count zero indicates the end of the array. Each item is encoded per the array's item schema.
	// If a block's count is negative, its absolute value is used, and the count is followed immediately by a long block
	// size indicating the number of bytes in the block. This block size permits fast skipping through data, e.g., when
	// projecting a record to a subset of its fields.
	// For example, the array schema {"type": "array", "items": "long"}
	// an array containing the items 3 and 27 could be encoded as the long value 2 (encoded as hex 04)
	// followed by long values 3 and 27 (encoded as hex 06 36) terminated by zero:
	// 04 06 36 00

	return func(name string, d *decode.D) any {
		var values []any
		d.FieldArray(name, func(d *decode.D) {
			count := int64(-1)
			for count != 0 {
				d.FieldStruct("block", func(d *decode.D) {
					count = d.FieldSintFn("count", VarZigZag)
					if count < 0 {
						d.FieldSintFn("size", VarZigZag)
						count *= -1
					}
					d.FieldArray("data", func(d *decode.D) {
						for i := int64(0); i < count; i++ {
							values = append(values, valueD("entry", d))
						}
					})
				})
			}
		})
		return values
	}, nil
}
