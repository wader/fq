package icc

// http://www.color.org/ICC1-V41.pdf

import (
	"fq/pkg/decode"
)

var Tag = &decode.Format{
	Name: "icc",
	New:  func() decode.Decoder { return &TagDecoder{} },
}

// TagDecoder is a ICC profile decoder
type TagDecoder struct {
	decode.Common
}

// Decode ICC tag
func (d *TagDecoder) Decode() {

	d.FieldNoneFn("header", func() {
		d.FieldU32("size")
		d.FieldU32("cmm_type_singature")
		d.FieldU32("version")
		d.FieldU32("device_class_signature")
		d.FieldU32("color_space")
		d.FieldU32("connection_space")
		d.FieldBytesLen("timestamp", 12) // TODO
		d.FieldU32("file_signature")
		d.FieldU32("primary_platform")
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

	d.FieldNoneFn("tag_table", func() {
		tagCount := d.FieldU32("count")
		d.FieldNoneFn("table", func() {
			for i := uint64(0); i < tagCount; i++ {
				d.FieldNoneFn("element", func() {
					d.FieldUTF8("signature", 4)
					d.FieldU32("offset")
					d.FieldU32("size")
				})
			}
		})
	})

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

}
