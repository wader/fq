package rtmp

// https://rtmp.veriskope.com/pdf/amf0-file-format-specification.pdf

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.AMF0,
		&decode.Format{
			Description: "Action Message Format 0",
			DecodeFn:    amf0Decode,
		})
}

const (
	typeNumber      = 0x0
	typeBoolean     = 0x1
	typeString      = 0x2
	typeObject      = 0x3
	typeMovieClip   = 0x4
	typeNull        = 0x5
	typeUndefined   = 0x6
	typeReference   = 0x7
	typeECMAArray   = 0x8
	typeObjectEnd   = 0x9
	typeStrictArray = 0xa
	typeDate        = 0xb
	typeLongString  = 0xc
	typeUnsupported = 0xd
	typeRecordSet   = 0xe
	typeXML         = 0xf
	typeTypedObject = 0x10
)

var typeNames = scalar.UintMapSymStr{
	typeNumber:      "number",
	typeBoolean:     "boolean",
	typeString:      "string",
	typeObject:      "object",
	typeMovieClip:   "movie_clip",
	typeNull:        "null",
	typeUndefined:   "undefined",
	typeReference:   "reference",
	typeECMAArray:   "ecma_array",
	typeObjectEnd:   "object_end",
	typeStrictArray: "strict_array",
	typeDate:        "date",
	typeLongString:  "long_string",
	typeUnsupported: "unsupported",
	typeRecordSet:   "record_set",
	typeXML:         "xml",
	typeTypedObject: "typed_object",
}

func amf0DecodeValue(d *decode.D) {
	typ := d.FieldU8("type", typeNames)

	switch typ {
	case typeNumber:
		d.FieldF64("value")
	case typeBoolean:
		d.FieldU8("value")
	case typeString:
		l := d.FieldU16("length")
		d.FieldUTF8("value", int(l))
	case typeObject:
		d.FieldArray("value", func(d *decode.D) {
			var typ uint64
			for typ != typeObjectEnd {
				d.FieldStruct("pair", func(d *decode.D) {
					d.FieldStruct("key", func(d *decode.D) {
						l := d.FieldU16("length")
						d.FieldUTF8("value", int(l))
					})
					typ = d.PeekUintBits(8)
					d.FieldStruct("value", amf0DecodeValue)
				})
			}
		})
	case typeNull:
		d.FieldValueAny("value", nil)
	case typeUndefined:
		d.FieldValueAny("value", nil) // TODO: ?
	case typeReference:
		d.FieldU16("value") // TODO: index pointer
	case typeECMAArray:
		d.FieldU32("count")
		d.FieldArray("value", func(d *decode.D) {
			var typ uint64
			for typ != typeObjectEnd {
				d.FieldStruct("entry", func(d *decode.D) {
					d.FieldStruct("key", func(d *decode.D) {
						l := d.FieldU16("length")
						d.FieldUTF8("value", int(l))
					})
					typ = d.PeekUintBits(8)
					d.FieldStruct("value", amf0DecodeValue)
				})
			}
		})
	case typeObjectEnd:
		// nop
	case typeStrictArray:
		length := d.FieldU32("count")
		d.FieldArray("value", func(d *decode.D) {
			for i := uint64(0); i < length; i++ {
				d.FieldStruct("entry", amf0DecodeValue)
			}
		})
	case typeDate:
		d.FieldF64("date_time")
		d.FieldS16("local_data_time_offset")
	case typeLongString:
		l := d.FieldU32("length")
		d.FieldUTF8("value", int(l))
	case typeXML:
		l := d.FieldU16("length")
		d.FieldUTF8("value", int(l))
	case typeTypedObject:
		d.FieldStruct("class_name", func(d *decode.D) {
			l := d.FieldU16("length")
			d.FieldUTF8("value", int(l))
		})
		d.FieldArray("value", func(d *decode.D) {
			var typ uint64
			for typ != typeObjectEnd {
				d.FieldStruct("pair", func(d *decode.D) {
					d.FieldStruct("key", func(d *decode.D) {
						l := d.FieldU16("length")
						d.FieldUTF8("value", int(l))
					})
					typ = d.PeekUintBits(8)
					d.FieldStruct("value", amf0DecodeValue)
				})
			}
		})
	default:
		d.Fatalf("unknown type %d", typ)
	}
}

func amf0Decode(d *decode.D) any {
	amf0DecodeValue(d)
	return nil
}
