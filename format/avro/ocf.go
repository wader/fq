package avro

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/avro/decoders"
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var jsonGroup decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.AVRO_OCF,
		Description: "Avro object container file",
		Groups:      []string{format.PROBE},
		DecodeFn:    avroDecodeOCF,
		Dependencies: []decode.Dependency{
			{Names: []string{format.JSON}, Group: &jsonGroup},
		},
	})
}

func ScalarDescription(description string) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		s.Description = description
		return s, nil
	})
}

type HeaderData struct {
	Schema *schema.SimplifiedSchema
	Codec  string
	Sync   []byte
}

func decodeHeader(d *decode.D) HeaderData {
	var headerData HeaderData

	// Header is encoded in avro so could use avro decoder, but doing it manually so we can
	// keep asserts and treating schema as JSON
	d.FieldRawLen("magic", 4*8, d.AssertBitBuf([]byte{'O', 'b', 'j', 1}))
	var blockCount int64 = -1
	d.FieldStructArrayLoop("meta", "block",
		func() bool { return blockCount != 0 },
		func(d *decode.D) {
			blockCount = d.FieldSFn("count", decoders.VarZigZag)
			// If its negative, then theres another long representing byte size
			if blockCount < 0 {
				blockCount *= -1
				d.FieldSFn("size", decoders.VarZigZag)
			}
			if blockCount == 0 {
				return
			}

			var i int64
			d.FieldStructArrayLoop("entries", "entry", func() bool { return i < blockCount }, func(d *decode.D) {
				keyL := d.FieldSFn("key_len", decoders.VarZigZag)
				key := d.FieldUTF8("key", int(keyL))
				valL := d.FieldSFn("value_len", decoders.VarZigZag)
				if key == "avro.schema" {
					v, _ := d.FieldFormatLen("value", valL*8, jsonGroup, nil)
					s, err := schema.From(v.V.(*scalar.S).Actual)
					headerData.Schema = &s
					if err != nil {
						d.Fatalf("Failed to parse schema: %s", err)
					}
				} else if key == "avro.codec" {
					headerData.Codec = d.FieldUTF8("value", int(valL))
				} else {
					d.FieldUTF8("value", int(valL))
				}
				i++
			})
		})
	if headerData.Schema == nil {
		d.Fatalf("No schema found in header")
	}

	if headerData.Codec == "null" {
		headerData.Codec = ""
	}

	syncbb := d.FieldRawLen("sync", 16*8)
	var err error
	headerData.Sync, err = syncbb.BytesLen(16)
	if err != nil {
		d.Fatalf("unable to read sync bytes: %v", err)
	}
	return headerData
}

func avroDecodeOCF(d *decode.D, in interface{}) interface{} {
	header := decodeHeader(d)

	decodeFn, err := decoders.DecodeFnForSchema(*header.Schema)
	if err != nil {
		d.Fatalf("unable to create codec: %v", err)
	}

	d.FieldStructArrayLoop("blocks", "block", func() bool { return d.NotEnd() }, func(d *decode.D) {
		count := d.FieldSFn("count", decoders.VarZigZag)
		if count <= 0 {
			return
		}
		size := d.FieldSFn("size", decoders.VarZigZag)
		// Currently not supporting encodings.
		if header.Codec != "" {
			d.FieldRawLen("data", size*8, ScalarDescription(header.Codec+" encoded"))
		} else {
			i := int64(0)
			d.FieldArrayLoop("data", func() bool { return i < count }, func(d *decode.D) {
				decodeFn("datum", d)
				i++
			})
		}
		d.FieldRawLen("sync", 16*8, d.AssertBitBuf(header.Sync))
	})

	return nil
}
