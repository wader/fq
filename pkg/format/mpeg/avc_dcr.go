package mpeg

// ISO/IEC 14496-15, 5.3.3.1.2 Syntax

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_DCR,
		Description: "H.264/AVC Decoder configuration record",
		DecodeFn:    avcDcrDecode,
	})
}

func avcDcrDecode(d *decode.D, in interface{}) interface{} {
	d.FieldU8("configuration_version")
	profileIdc := d.FieldU8("profile_indication")
	d.FieldU8("profile_compatibility")
	d.FieldU8("level_indication")
	d.FieldU6("reserved0")
	lengthSizeMinusOne := d.FieldU2("length_size_minus_one")
	d.FieldU3("reserved1")
	numSeqParamSets := d.FieldU5("num_of_sequence_parameter_sets")
	d.FieldArrayFn("sequence_parameter_sets", func(d *decode.D) {
		for i := uint64(0); i < numSeqParamSets; i++ {
			d.FieldStructFn("parameter_set", func(d *decode.D) {
				paramSetLen := d.FieldU16("length")
				d.FieldBitBufLen("set", int64(paramSetLen)*8)
			})
		}
	})
	numPicParamSets := d.FieldU8("num_of_picture_parameter_sets")
	d.FieldArrayFn("parameter_sets", func(d *decode.D) {
		for i := uint64(0); i < numPicParamSets; i++ {
			d.FieldStructFn("parameter_set", func(d *decode.D) {
				paramSetLen := d.FieldU16("length")
				d.FieldBitBufLen("set", int64(paramSetLen)*8)
			})
		}
	})

	switch profileIdc {
	case 100, 110, 122, 144:
		d.FieldU6("reserved2")
		d.FieldU6("chroma_format")
		d.FieldU4("reserved3")
		d.FieldU3("bit_depth_luma_minus8")
		d.FieldU5("reserved4")
		d.FieldU3("bit_depth_chroma_minus8")
		numSeqParamSetExt := d.FieldU5("num_of_sequence_parameter_set_ext")
		d.FieldArrayFn("parameter_set_exts", func(d *decode.D) {
			for i := uint64(0); i < numSeqParamSetExt; i++ {
				d.FieldStructFn("parameter_set_ext", func(d *decode.D) {
					paramSetLen := d.FieldU16("length")
					d.FieldBitBufLen("set", int64(paramSetLen)*8)
				})
			}
		})
	}

	return format.AvcDcrOut{LengthSize: lengthSizeMinusOne + 1}
}
