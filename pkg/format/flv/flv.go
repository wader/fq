package flv

// https://www.adobe.com/content/dam/acom/en/devnet/flv/video_file_format_spec_v10.pdf

import (
	"fq/pkg/decode"
)

var File = &decode.Format{
	Name:  "flv",
	MIMEs: []string{"video/x-flv"},
	New:   func() decode.Decoder { return &FileDecoder{} },
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

// FileDecoder is a FLV decoder
type FileDecoder struct{ decode.Common }

// Decode decodes a FLV file
func (d *FileDecoder) Decode() {

	var fieldScriptDataObject func()
	var fieldScriptDataVariable func(name string)

	fieldScriptDataString := func(name string) {
		d.FieldStrFn(name, func() (string, string) {
			len := d.U16()
			return d.UTF8(len), ""
		})
	}
	fieldScriptDataStringLong := func(name string) {
		d.FieldStrFn(name, func() (string, string) {
			len := d.U32()
			return d.UTF8(len), ""
		})
	}

	fieldScriptDataVariable = func(name string) {
		d.FieldNoneFn(name, func() {
			fieldScriptDataString("name")
			fieldScriptDataString("data")
		})
	}

	fieldScriptDataValue := func(name string) uint64 {
		typ := d.FieldStringMapFn("type", typeNames, "Unknown", d.U8)
		if typ == typeECMAArray {
			d.FieldU32("ecma_array_length")
		}

		switch typ {
		case typeNumber:
			d.FieldF64("number")
		case typeBoolean:
			d.FieldU8("boolean")
		case typeString:
			fieldScriptDataString("string")
		case typeObject:
			fieldScriptDataObject()
		case typeMovieClip:
			fieldScriptDataString("path")
		case typeNull:
		case typeUndefined:
		case typeReference:
			d.FieldU16("reference")
		case typeECMAArray:
			d.FieldNoneFn("array", func() {
				for {
					if d.PeekBits(24) == typeObjectEnd { // variableEnd?
						d.FieldU24("end")
						break
					}
					fieldScriptDataVariable("sasdadas")
				}
			})
		case typeStrictArray:
			length := d.FieldU32("length")
			for i := uint64(0); i < length; i++ {
				fieldScriptDataVariable("sasdadas")
			}
		case typeDate:
			d.FieldF64("date_time")
			d.FieldS16("local_data_time_offset")

		case typeLongString:
			fieldScriptDataStringLong("asdsad")

		case typeObjectEnd: // variableEnd also?

		}

		return typ
	}

	fieldScriptDataObject = func() {
		d.FieldNoneFn("object", func() {
			fieldScriptDataString("name")
			fieldScriptDataValue("data")
		})
	}

	d.FieldValidateString("signature", "FLV")
	d.FieldU8("version")
	d.FieldValidateUFn("type_flags_reserved", 0, d.U5)
	d.FieldU1("type_flags_audio")
	d.FieldValidateUFn("type_flags_reserved", 0, d.U1)
	d.FieldU1("type_flags_video")
	dataOffset := d.FieldU32("data_offset")

	d.SeekAbs(dataOffset * 8)

	for !d.End() {
		d.FieldU32("previous_tag_size")

		d.FieldNoneFn("tag", func() {
			tagType := d.FieldStringMapFn("tag_type", tagTypeNames, "unknown", d.U8)
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
}
