package bson

// https://bsonspec.org/spec.html

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed bson.jq
//go:embed bson.md
var bsonFS embed.FS

func init() {
	interp.RegisterFormat(
		format.BSON,
		&decode.Format{
			Description: "Binary JSON",
			DecodeFn:    decodeBSON,
			Functions:   []string{"torepr"},
		})
	interp.RegisterFS(bsonFS)
}

const (
	elementTypeDouble     = 0x01
	elementTypeString     = 0x02
	elementTypeDocument   = 0x03
	elementTypeArray      = 0x04
	elementTypeBinary     = 0x05
	elementTypeUndefined  = 0x06
	elementTypeObjectID   = 0x07
	elementTypeBoolean    = 0x08
	elementTypeDatetime   = 0x09
	elementTypeNull       = 0x0a
	elementTypeRegexp     = 0x0b
	elementTypeJavaScript = 0x0d
	elementTypeInt32      = 0x10
	elementTypeTimestamp  = 0x11
	elementTypeInt64      = 0x12
	elementTypeDecimal128 = 0x13
	elementTypeMinKey     = 0xFF
	elementTypeMaxKey     = 0x7f
)

var elementTypeMap = scalar.UintMap{
	elementTypeDouble:     {Sym: "double", Description: "64-bit binary floating point"},
	elementTypeString:     {Sym: "string", Description: "UTF-8 string"},
	elementTypeDocument:   {Sym: "document", Description: "Embedded document"},
	elementTypeArray:      {Sym: "array", Description: "Array"},
	elementTypeBinary:     {Sym: "binary", Description: "Binary data"},
	elementTypeUndefined:  {Sym: "undefined", Description: "Undefined (deprecated)"},
	elementTypeObjectID:   {Sym: "object_id", Description: "ObjectId"},
	elementTypeBoolean:    {Sym: "boolean", Description: "Boolean"},
	elementTypeDatetime:   {Sym: "datetime", Description: "UTC datetime"},
	elementTypeNull:       {Sym: "null", Description: "Null value"},
	elementTypeRegexp:     {Sym: "regexp", Description: "Regular expression"},
	elementTypeJavaScript: {Sym: "javascript", Description: "JavaScript code"},
	elementTypeInt32:      {Sym: "int32", Description: "32-bit integer"},
	elementTypeTimestamp:  {Sym: "timestamp", Description: "Timestamp"},
	elementTypeInt64:      {Sym: "int64", Description: "64-bit integer"},
	elementTypeDecimal128: {Sym: "decimal128", Description: "128-bit decimal floating point"},
	elementTypeMinKey:     {Sym: "minkey", Description: "Min key"},
	elementTypeMaxKey:     {Sym: "maxkey", Description: "Max key"},
}

func decodeBSONDocument(d *decode.D) {
	size := d.FieldS32("size")
	d.FramedFn((size-4)*8, func(d *decode.D) {
		d.FieldArray("elements", func(d *decode.D) {
			for d.BitsLeft() > 8 {
				d.FieldStruct("element", func(d *decode.D) {
					typ := d.FieldU8("type", elementTypeMap)
					d.FieldUTF8Null("name")
					switch typ {
					case elementTypeDouble:
						d.FieldF64("value")
					case elementTypeString:
						length := d.FieldU32("length")
						d.FieldUTF8NullFixedLen("value", int(length))
					case elementTypeDocument:
						d.FieldStruct("value", decodeBSONDocument)
					case elementTypeArray:
						d.FieldStruct("value", decodeBSONDocument)
					case elementTypeBinary:
						length := d.FieldS32("length")
						d.FieldU8("subtype")
						d.FieldRawLen("value", length*8)
					case elementTypeUndefined:
						//deprecated
					case elementTypeObjectID:
						d.FieldRawLen("value", 12*8)
					case elementTypeBoolean:
						d.FieldU8("value")
					case elementTypeDatetime:
						d.FieldS64("value")
					case elementTypeNull:
						d.FieldValueAny("value", nil)
					case elementTypeRegexp:
						d.FieldUTF8Null("value")
						d.FieldUTF8Null("options")
					case elementTypeJavaScript:
						length := d.FieldS32("length")
						d.FieldUTF8NullFixedLen("value", int(length))
					case elementTypeInt32:
						d.FieldS32("value")
					case elementTypeTimestamp:
						d.FieldU64("value")
					case elementTypeInt64:
						d.FieldS64("value")
					case elementTypeDecimal128:
						// TODO: Parse the IEEE 754 decimal128 value.
						d.FieldRawLen("value", 128)
					case elementTypeMinKey:
						d.FieldValueAny("value", nil)
					case elementTypeMaxKey:
						d.FieldValueAny("value", nil)
					default:
						d.FieldRawLen("value", d.BitsLeft())
					}
				})
			}
		})
		d.FieldU8("terminator", d.UintValidate(0))
	})
}

func decodeBSON(d *decode.D) any {
	d.Endian = decode.LittleEndian

	decodeBSONDocument(d)

	return nil
}
