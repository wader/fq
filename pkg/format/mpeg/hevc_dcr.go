package mpeg

// ISO/IEC 14496-15 AVC file format, 5.3.3.1.2 Syntax
// ISO_IEC_14496-10 AVC

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_HEVC_DCR,
		Description: "H.265/HEVC Decoder configuration record",
		DecodeFn:    hevcDcrDecode,
	})
}

func hevcDcrDecode(d *decode.D, in interface{}) interface{} {
	d.FieldU8("configurationVersion")
	d.FieldU2("general_profile_space")
	d.FieldU1("general_tier_flag")
	d.FieldU5("general_profile_idc")
	d.FieldU32("general_profile_compatibility_flags")
	d.FieldU48("general_constraint_indicator_flags")
	d.FieldU8("general_level_idc")
	d.FieldU4("reserved0")
	d.FieldU12("min_spatial_segmentation_idc")
	d.FieldU6("reserved1")
	d.FieldU2("parallelism_type")
	d.FieldU6("reserved2")
	d.FieldU2("chroma_format_idc")
	d.FieldU5("reserved3")
	d.FieldU3("bit_depth_luma_minus8")
	d.FieldU5("reserved4")
	d.FieldU3("bit_depth_chroma_minus8")
	d.FieldU16("avg_frame_rate")
	d.FieldU2("constant_frame_rate")
	d.FieldU3("num_temporal_layers")
	d.FieldU1("temporal_id_nested")
	lengthSizeMinusOne := d.FieldU2("length_size_minus_one")
	numArrays := d.FieldU8("num_of_arrays")
	d.FieldArrayFn("arrays", func(d *decode.D) {
		for i := uint64(0); i < numArrays; i++ {
			d.FieldStructFn("array", func(d *decode.D) {
				d.FieldU1("array_completeness")
				d.FieldU1("reserved0")
				d.FieldU6("nal_unit_type")
				numNals := d.FieldU16("num_nalus")
				d.FieldArrayFn("nals", func(d *decode.D) {
					for i := uint64(0); i < numNals; i++ {
						d.FieldStructFn("nal", func(d *decode.D) {
							nalUnitLength := int64(d.FieldU16("nal_unit_length"))
							d.FieldBitBufLen("data", nalUnitLength*8)
						})
					}
				})
			})
		}
	})

	return format.HevcDcrOut{LengthSize: lengthSizeMinusOne + 1}
}
