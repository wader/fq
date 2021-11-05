package icc

// http://www.color.org/ICC1-V41.pdf
// https://www.color.org/icc32.pdf

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.ICC_PROFILE,
		Description: "International Color Consortium profile",
		DecodeFn:    iccProfileDecode,
	})
}

func xyzType(d *decode.D) {
	d.FieldFP32("X")
	d.FieldFP32("Y")
	d.FieldFP32("Z")
}

func textType(d *decode.D) {
	d.FieldUTF8NullTerminatedLen("text", int(d.BitsLeft()/8))
}

func paraType(d *decode.D) {
	d.FieldU32("reserved0")
	d.FieldU16("function_type")
	d.FieldU16("reserved1")
	d.FieldRawLen("parameters", d.BitsLeft())
}

func descType(d *decode.D) {
	descLen := d.FieldU32("description_length")
	d.FieldUTF8NullTerminatedLen("description", int(descLen))
	d.FieldU32("language_code")
	localDescLen := d.FieldU32("localizable_description_length")
	d.FieldUTF8NullTerminatedLen("localizable_description", int(localDescLen))
	d.FieldU16("script_code")
	d.FieldU8("macintosh_description_length")
	d.FieldUTF8NullTerminatedLen("macintosh_description", 67)
}

var typToDecode = map[string]func(d *decode.D){
	"XYZ ": xyzType,
	"text": textType,
	"para": paraType,
	"desc": descType,
}

func decodeBCDU8(d *decode.D) uint64 {
	n := d.U8()
	return (n>>4)*10 + n&0xf
}

func iccProfileDecode(d *decode.D, in interface{}) interface{} {
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

	d.LenFn(int64(size)*8, func(d *decode.D) {
		d.FieldStruct("header", func(d *decode.D) {
			d.FieldU32("size")
			d.FieldUTF8NullTerminatedLen("cmm_type_signature", 4)
			d.FieldUFn("version_major", decodeBCDU8)
			d.FieldUFn("version_minor", decodeBCDU8)
			d.FieldU16("version_reserved")
			d.FieldUTF8NullTerminatedLen("device_class_signature", 4)
			d.FieldUTF8NullTerminatedLen("color_space", 4)
			d.FieldUTF8NullTerminatedLen("connection_space", 4)
			d.FieldStruct("timestamp", func(d *decode.D) {
				d.FieldU16("year")
				d.FieldU16("month")
				d.FieldU16("day")
				d.FieldU16("hours")
				d.FieldU16("minutes")
				d.FieldU16("seconds")

			})
			d.FieldUTF8NullTerminatedLen("file_signature", 4)
			d.FieldUTF8NullTerminatedLen("primary_platform", 4)
			d.FieldU32("flags")
			d.FieldUTF8NullTerminatedLen("device_manufacturer", 4)
			d.FieldUTF8NullTerminatedLen("device_model", 4)
			d.FieldUTF8NullTerminatedLen("device_attribute", 8)
			d.FieldUTF8NullTerminatedLen("render_intent", 4)
			d.FieldUTF8NullTerminatedLen("xyz_illuminant", 12)
			d.FieldUTF8NullTerminatedLen("profile_creator_signature", 4)
			d.FieldUTF8NullTerminatedLen("profile_id", 16)
			d.FieldRawLen("reserved", 28*8, d.BitBufIsZero)
		})

		d.FieldStruct("tag_table", func(d *decode.D) {
			tagCount := d.FieldU32("count")
			d.FieldArray("table", func(d *decode.D) {
				for i := uint64(0); i < tagCount; i++ {
					d.FieldStruct("element", func(d *decode.D) {
						d.FieldUTF8NullTerminatedLen("signature", 4)
						offset := d.FieldU32("offset")
						size := d.FieldU32("size")

						d.RangeFn(int64(offset)*8, int64(size)*8, func(d *decode.D) {
							typ := d.FieldUTF8NullTerminatedLen("type", 4)
							d.FieldU32("reserved")

							if fn, ok := typToDecode[typ]; ok {
								d.LenFn(int64(size-4-4)*8, fn)
							} else {
								d.FieldRawLen("data", int64(size-4-4)*8)
							}
						})

						// "All tag data is required to start on a 4-byte boundary (relative to the start of the profile data stream)"
						// we can't add this at the start of the element as we don't know how big the previous element in the stream
						// was. instead add alignment after if offset+size does not align and to be sure clamp it if outside buffer.
						paddingStart := int64(offset) + int64(size)
						paddingBytes := (4 - (int64(offset)+int64(size))%4) % 4
						paddingBytes = num.MinInt64(paddingBytes, d.Len()-(paddingStart+paddingBytes))
						if paddingBytes != 0 {
							d.RangeFn(paddingStart*8, paddingBytes*8, func(d *decode.D) {
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
