package icc

// http://www.color.org/ICC1-V41.pdf
// https://www.color.org/icc32.pdf

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.ICC_Profile,
		&decode.Format{
			Description: "International Color Consortium profile",
			DecodeFn:    iccProfileDecode,
		})
}

func xyzType(_ int64, d *decode.D) {
	d.FieldFP32("x")
	d.FieldFP32("y")
	d.FieldFP32("z")
}

func textType(_ int64, d *decode.D) {
	d.FieldUTF8NullFixedLen("text", int(d.BitsLeft()/8))
}

func paraType(_ int64, d *decode.D) {
	d.FieldU32("reserved0")
	d.FieldU16("function_type")
	d.FieldU16("reserved1")
	d.FieldRawLen("parameters", d.BitsLeft())
}

func descType(_ int64, d *decode.D) {
	descLen := d.FieldU32("description_length")
	d.FieldUTF8NullFixedLen("description", int(descLen))
	d.FieldU32("language_code")
	localDescLen := d.FieldU32("localizable_description_length")
	d.FieldUTF8NullFixedLen("localizable_description", int(localDescLen))
	d.FieldU16("script_code")
	d.FieldU8("macintosh_description_length")
	d.FieldUTF8NullFixedLen("macintosh_description", 67)
}

func multiLocalizedUnicodeType(tagStart int64, d *decode.D) {
	numberOfNames := d.FieldU32("number_of_names")
	recordSize := d.FieldU32("record_size")
	d.FieldArray("names", func(d *decode.D) {
		for i := uint64(0); i < numberOfNames; i++ {
			d.FieldStruct("name", func(d *decode.D) {
				d.FieldUTF8("language_code", 2)
				d.FieldUTF8("country_code", 2)
				nameLength := d.FieldU32("name_length")
				nameOffset := d.FieldU32("name_offset")
				d.RangeFn(tagStart+int64(nameOffset)*8, int64(nameLength)*8, func(d *decode.D) {
					d.FieldUTF16BE("value", int(nameLength))
				})
			})
			recordPadding := int64(recordSize) - 2 - 2 - 4 - 4
			if recordPadding > 0 {
				d.FieldRawLen("padding", recordPadding)
			}
		}
	})
}

var typeToDecode = map[string]func(tagStart int64, d *decode.D){
	"XYZ":  xyzType,
	"text": textType,
	"para": paraType,
	"desc": descType,
	"mluc": multiLocalizedUnicodeType,
}

func decodeBCDU8(d *decode.D) uint64 {
	n := d.U8()
	return (n>>4)*10 + n&0xf
}

func iccProfileDecode(d *decode.D) any {
	/*
	   0..3 Profile size uInt32Number
	   4..7 CMM Type signature see below
	   8..11 Profile version number see below
	   12..15 Profile/Device Class signature see below
	   16..19 Color space of data (possibly a derived space) [i.e. “the canonical input space”] see below
	   20..23 Profile Connection Space (PCS) [i.e. “the canonical output space”] see below
	   24..35 Date and time this profile was first created dateTimeNumber
	   36..39 ‘acsp’ (61637370h) profile file signature
	   40..43 Primary Platform signature see below
	   44..47 Flags to indicate various options for the CMM such as distributed processing and caching options see below
	   48..51 Device manufacturer of the device for which this profile is created see below
	   52..55 Device model of the device for which this profile is created see below
	   56..63 Device attributes unique to the particular device setup such as media type see below
	   64..67 Rendering Intent see below
	   68..79 The XYZ values of the illuminant of the Profile Connection Space. This must correspond to D50. It is explained in more detail in A.1. XYZNumber
	   80..83 Profile Creator signature see below
	   84..99 Profile ID see below
	   100..127 28 bytes reserved for future expansion - must be set to zeros
	*/

	// TODO: PokeU32()?
	size := d.U32()
	d.SeekRel(-4 * 8)

	d.FramedFn(int64(size)*8, func(d *decode.D) {
		d.FieldStruct("header", func(d *decode.D) {
			d.FieldU32("size")
			d.FieldUTF8NullFixedLen("cmm_type_signature", 4, scalar.ActualTrimSpace)
			d.FieldUintFn("version_major", decodeBCDU8)
			d.FieldUintFn("version_minor", decodeBCDU8)
			d.FieldU16("version_reserved")
			d.FieldUTF8NullFixedLen("device_class_signature", 4, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("color_space", 4, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("connection_space", 4, scalar.ActualTrimSpace)
			d.FieldStruct("timestamp", func(d *decode.D) {
				d.FieldU16("year")
				d.FieldU16("month")
				d.FieldU16("day")
				d.FieldU16("hours")
				d.FieldU16("minutes")
				d.FieldU16("seconds")

			})
			d.FieldUTF8NullFixedLen("file_signature", 4, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("primary_platform", 4, scalar.ActualTrimSpace)
			d.FieldU32("flags")
			d.FieldUTF8NullFixedLen("device_manufacturer", 4, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("device_model", 4, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("device_attribute", 8, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("render_intent", 4, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("xyz_illuminant", 12, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("profile_creator_signature", 4, scalar.ActualTrimSpace)
			d.FieldUTF8NullFixedLen("profile_id", 16, scalar.ActualTrimSpace)
			d.FieldRawLen("reserved", 28*8, d.BitBufIsZero())
		})

		d.FieldStruct("tag_table", func(d *decode.D) {
			tagCount := d.FieldU32("count")
			d.FieldArray("table", func(d *decode.D) {
				for i := uint64(0); i < tagCount; i++ {
					d.FieldStruct("element", func(d *decode.D) {
						d.FieldUTF8NullFixedLen("signature", 4, scalar.ActualTrimSpace)
						offset := d.FieldU32("offset")
						size := d.FieldU32("size")

						d.RangeFn(int64(offset)*8, int64(size)*8, func(d *decode.D) {
							tagStart := d.Pos()
							typ := d.FieldUTF8NullFixedLen("type", 4, scalar.ActualTrimSpace)
							d.FieldU32("reserved")

							if fn, ok := typeToDecode[typ]; ok {
								fn(tagStart, d)
							} else {
								d.FieldRawLen("data", int64(size-4-4)*8)
							}
						})

						// "All tag data is required to start on a 4-byte boundary (relative to the start of the profile data stream)"
						// we can't add this at the start of the element as we don't know how big the previous element in the stream
						// was. instead add alignment after if offset+size does not align and to be sure clamp it if outside buffer.
						alignStart := int64(offset) + int64(size)
						alignBytes := (4 - (int64(offset)+int64(size))%4) % 4
						alignBytes = min(d.Len()/8-alignStart, alignBytes)
						if alignBytes != 0 {
							d.RangeFn(alignStart*8, alignBytes*8, func(d *decode.D) {
								d.FieldRawLen("alignment", d.BitsLeft())
							})
						}
					})
				}
			})
		})
	})

	return nil
}
