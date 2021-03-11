package flv

// https://www.adobe.com/content/dam/acom/en/devnet/flv/video_file_format_spec_v10.pdf

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.FLV,
		Description: "Flash video",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"video/x-flv"},
		DecodeFn:    flvDecode,
	})
}

const (
	audioData        = 8
	videoData        = 9
	scriptDataObject = 18
)

var tagTypeNames = map[uint64]string{
	audioData:        "audioData",
	videoData:        "videoData",
	scriptDataObject: "scriptDataObject",
}

const (
	typeNumber      = 0
	typeBoolean     = 1
	typeString      = 2
	typeObject      = 3
	typeMovieClip   = 4
	typeNull        = 5
	typeUndefined   = 6
	typeReference   = 7
	typeECMAArray   = 8
	typeObjectEnd   = 9
	typeStrictArray = 10
	typeDate        = 11
	typeLongString  = 12
)

var typeNames = map[uint64]string{
	typeNumber:      "Number",
	typeBoolean:     "Boolean",
	typeString:      "String",
	typeObject:      "Object",
	typeMovieClip:   "MovieClip",
	typeNull:        "Null",
	typeUndefined:   "Undefined",
	typeReference:   "Reference",
	typeECMAArray:   "ECMAArray",
	typeObjectEnd:   "ObjectEnd",
	typeStrictArray: "StrictArray",
	typeDate:        "Date",
	typeLongString:  "LongString",
}

func flvDecode(d *decode.D, in interface{}) interface{} {
	var fieldScriptDataObject func()
	var fieldScriptDataVariable func(d *decode.D, name string)

	fieldScriptDataString := func(d *decode.D, name string) {
		d.FieldStrFn(name, func() (string, string) {
			len := d.U16()
			return d.UTF8(int(len)), ""
		})
	}
	fieldScriptDataStringLong := func(d *decode.D, name string) {
		d.FieldStrFn(name, func() (string, string) {
			len := d.U32()
			return d.UTF8(int(len)), ""
		})
	}

	fieldScriptDataVariable = func(d *decode.D, name string) {
		d.FieldStructFn(name, func(d *decode.D) {
			fieldScriptDataString(d, "name")
			fieldScriptDataString(d, "data")
		})
	}

	fieldScriptDataValue := func(d *decode.D, name string) uint64 {
		typ, _ := d.FieldStringMapFn("type", typeNames, "Unknown", d.U8, decode.NumberDecimal)
		if typ == typeECMAArray {
			d.FieldU32("ecma_array_length")
		}

		switch typ {
		case typeNumber:
			d.FieldF64("number")
		case typeBoolean:
			d.FieldU8("boolean")
		case typeString:
			fieldScriptDataString(d, "string")
		case typeObject:
			fieldScriptDataObject()
		case typeMovieClip:
			fieldScriptDataString(d, "path")
		case typeNull:
		case typeUndefined:
		case typeReference:
			d.FieldU16("reference")
		case typeECMAArray:
			d.FieldArrayFn("array", func(d *decode.D) {
				for {
					if d.PeekBits(24) == typeObjectEnd { // variableEnd?
						d.FieldU24("end")
						break
					}
					fieldScriptDataVariable(d, "sasdadas")
				}
			})
		case typeStrictArray:
			length := d.FieldU32("length")
			for i := uint64(0); i < length; i++ {
				fieldScriptDataVariable(d, "sasdadas")
			}
		case typeDate:
			d.FieldF64("date_time")
			d.FieldS16("local_data_time_offset")

		case typeLongString:
			fieldScriptDataStringLong(d, "asdsad")

		case typeObjectEnd: // variableEnd also?

		}

		return typ
	}

	fieldScriptDataObject = func() {
		d.FieldStructFn("object", func(d *decode.D) {
			fieldScriptDataString(d, "name")
			fieldScriptDataValue(d, "data")
		})
	}

	d.FieldValidateUTF8("signature", "FLV")
	d.FieldU8("version")
	d.FieldValidateUFn("type_flags_reserved", 0, d.U5)
	d.FieldU1("type_flags_audio")
	d.FieldValidateUFn("type_flags_reserved", 0, d.U1)
	d.FieldU1("type_flags_video")
	dataOffset := d.FieldU32("data_offset")

	d.SeekAbs(int64(dataOffset) * 8)

	d.FieldArrayFn("tags", func(d *decode.D) {
		for !d.End() {
			d.FieldStructFn("tag", func(d *decode.D) {
				d.FieldU32("previous_tag_size")
				tagType, _ := d.FieldStringMapFn("tag_type", tagTypeNames, "unknown", d.U8, decode.NumberDecimal)
				dataSize := d.FieldU24("data_size")
				d.FieldU24("timestamp")
				d.FieldU8("timestamp_extended")
				d.FieldU24("stream_id")

				switch tagType {
				case audioData, videoData:
					d.SeekRel(int64(dataSize) * 8)
				case scriptDataObject:
					for {
						if d.PeekBits(24) == typeObjectEnd {
							d.FieldU24("end")
							break
						}
						fieldScriptDataObject()
					}
				}
			})
		}
	})

	return nil
}
