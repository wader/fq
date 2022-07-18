package decoders

import (
	"fmt"

	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

func decodeRecordFn(schema schema.SimplifiedSchema) (DecodeFn, error) {
	if len(schema.Fields) == 0 {
		return nil, fmt.Errorf("record must have fields")
	}
	var fieldNames []string
	var fieldDecoders []func(string, *decode.D) any

	for _, f := range schema.Fields {
		fieldNames = append(fieldNames, f.Name)
		fc, err := DecodeFnForSchema(f.Type)
		if err != nil {
			return nil, fmt.Errorf("failed parsing record field %s: %w", f.Name, err)
		}
		fieldDecoders = append(fieldDecoders, fc)
	}

	// A record is encoded by encoding the values of its fields in the order that they are declared. In other words, a
	// record is encoded as just the concatenation of the encodings of its fields. Field values are encoded per their
	// schema. For example, the record schema
	// 	      { "type": "record",
	// 	        "name": "test",
	// 	        "fields" : [
	//	            {"name": "a", "type": "long"},
	//              {"name": "b", "type": "string"}
	//             ]
	// 	      }
	//
	// An instance of this record whose a field has value 27 (encoded as hex 36) and whose b field has value "foo"
	// (encoded as hex bytes 06 66 6f 6f), would be encoded simply as the concatenation of these, namely
	// the hex byte sequence:
	// 36 06 66 6f 6f

	return func(name string, d *decode.D) any {
		val := make(map[string]any)
		d.FieldStruct(name, func(d *decode.D) {
			for i, f := range fieldNames {
				val[f] = fieldDecoders[i](f, d)
			}
		})
		return val
	}, nil
}
