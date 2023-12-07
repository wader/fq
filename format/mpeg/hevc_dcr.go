package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var hevcHEVCNALUGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.HEVC_DCR,
		&decode.Format{
			Description: "H.265/HEVC Decoder Configuration Record",
			DecodeFn:    hevcDcrDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.HEVC_NALU}, Out: &hevcHEVCNALUGroup},
			},
		})
}

func hevcDcrDecode(d *decode.D) any {
	d.FieldU8("configuration_version")
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
	d.FieldU3("bit_depth_luma", scalar.UintActualAdd(8))
	d.FieldU5("reserved4")
	d.FieldU3("bit_depth_chroma", scalar.UintActualAdd(8))
	d.FieldU16("avg_frame_rate")
	d.FieldU2("constant_frame_rate")
	d.FieldU3("num_temporal_layers")
	d.FieldU1("temporal_id_nested")
	lengthSize := d.FieldU2("length_size", scalar.UintActualAdd(1))
	numArrays := d.FieldU8("num_of_arrays")
	d.FieldArray("arrays", func(d *decode.D) {
		for i := uint64(0); i < numArrays; i++ {
			d.FieldStruct("array", func(d *decode.D) {
				d.FieldU1("array_completeness")
				d.FieldU1("reserved0")
				d.FieldU6("nal_unit_type", hevcNALNames)
				numNals := d.FieldU16("num_nalus")
				d.FieldArray("nals", func(d *decode.D) {
					for i := uint64(0); i < numNals; i++ {
						d.FieldStruct("nal", func(d *decode.D) {
							nalUnitLength := int64(d.FieldU16("nal_unit_length"))
							d.FieldFormatLen("nal", nalUnitLength*8, &hevcHEVCNALUGroup, nil)
						})
					}
				})
			})
		}
	})

	return format.HEVC_DCR_Out{LengthSize: lengthSize}
}
