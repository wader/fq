package icc

// http://www.color.org/ICC1-V41.pdf

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:     format.ICC,
		DecodeFn: iccDecode,
	})
}

func xyzType(d *decode.D) {
	d.FieldUTF8("type", 4)
	d.FieldU32("reserved")
	d.FieldFP32("X")
	d.FieldFP32("Y")
	d.FieldFP32("Z")
}

var signatureToDecode = map[string]func(d *decode.D){
	"wtpt": xyzType,
	"bkpt": xyzType,
	"rXYZ": xyzType,
	"gXYZ": xyzType,
	"bXYZ": xyzType,
}

func iccDecode(d *decode.D) interface{} {
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
	d.FieldStructFn("header", func(d *decode.D) {
		d.FieldU32("size")
		d.FieldUTF8("cmm_type_signature", 4)
		d.FieldU32("version")
		d.FieldUTF8("device_class_signature", 4)
		d.FieldUTF8("color_space", 4)
		d.FieldUTF8("connection_space", 4)
		d.FieldStructFn("timestamp", func(d *decode.D) {
			d.FieldU16("year")
			d.FieldU16("month")
			d.FieldU16("day")
			d.FieldU16("hours")
			d.FieldU16("minutes")
			d.FieldU16("seconds")

		})
		d.FieldUTF8("file_signature", 4)
		d.FieldUTF8("primary_platform", 4)
		d.FieldU32("flags")
		d.FieldUTF8("device_manufacturer", 4)
		d.FieldUTF8("device_model", 4)
		d.FieldUTF8("device_attribute", 8)
		d.FieldUTF8("render_intent", 4)
		d.FieldUTF8("xyz_illuminant", 12)
		d.FieldUTF8("profile_creator_signature", 4)
		d.FieldUTF8("profile_id", 16)
		d.FieldValidateZeroPadding("reserved", 28*8)
	})

	d.FieldStructFn("tag_table", func(d *decode.D) {
		tagCount := d.FieldU32("count")
		d.FieldArrayFn("table", func(d *decode.D) {
			for i := uint64(0); i < tagCount; i++ {
				d.FieldStructFn("element", func(d *decode.D) {
					signature := d.FieldUTF8("signature", 4)
					offset := d.FieldU32("offset")
					size := d.FieldU32("size")

					if fn, ok := signatureToDecode[signature]; ok {
						d.SubRangeFn(int64(offset)*8, int64(size)*8, func() { fn(d) })
					} else {
						d.FieldBitBufRange("data", int64(offset)*8, int64(size)*8)
					}
				})
			}
		})
	})

	return nil
}
