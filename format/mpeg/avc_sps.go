package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.AVC_SPS,
		&decode.Format{
			Description: "H.264/AVC Sequence Parameter Set",
			DecodeFn:    avcSPSDecode,
		})
}

var avcVideoFormatMap = scalar.UintMapSymStr{
	0: "component",
	1: "pal",
	2: "ntsc",
	3: "secam",
	4: "mac",
	5: "unspecified",
	6: "reserved",
	7: "reserved",
}

var avcAspectRatioIdcMap = scalar.UintMapSymStr{
	0:  "unspecified",
	1:  "1:1",
	2:  "12:11",
	3:  "10:11",
	4:  "16:11",
	5:  "40:33",
	6:  "24:11",
	7:  "20:11",
	8:  "32:11",
	9:  "80:33",
	10: "18:11",
	11: "15:11",
	12: "64:33",
	13: "160:99",
	14: "4:3",
	15: "3:2",
	16: "2:1",
}

var chromaFormatMap = scalar.UintMapSymStr{
	0: "monochrome",
	1: "4:2:0",
	2: "4:2:2",
	3: "4:4:4",
}

func avcVuiParameters(d *decode.D) {
	aspectRatioInfoPresentFlag := d.FieldBool("aspect_ratio_info_present_flag")
	if aspectRatioInfoPresentFlag {
		aspectRatioIdc := d.FieldU8("aspect_ratio_idc", avcAspectRatioIdcMap)
		const extendedSAR = 255
		if aspectRatioIdc == extendedSAR {
			d.FieldU16("width")
			d.FieldU16("height")
		}
	}
	overscanInfoPresentFlag := d.FieldBool("overscan_info_present_flag")
	if overscanInfoPresentFlag {
		d.FieldBool("overscan_appropriate_flag")
	}
	videoSignalTypePresentFlag := d.FieldBool("video_signal_type_present_flag")
	if videoSignalTypePresentFlag {
		d.FieldU3("video_format", avcVideoFormatMap)
		d.FieldBool("video_full_range_flag")
		colourDescriptionPresentFlag := d.FieldBool("colour_description_present_flag")
		if colourDescriptionPresentFlag {
			d.FieldU8("colour_primaries", format.ISO_23091_2_ColourPrimariesMap)
			d.FieldU8("transfer_characteristics", format.ISO_23091_2_TransferCharacteristicMap)
			d.FieldU8("matrix_coefficients", format.ISO_23091_2_MatrixCoefficients)
		}
	}
	chromaLocInfoPresentFlag := d.FieldBool("chroma_loc_info_present_flag")
	if chromaLocInfoPresentFlag {
		d.FieldUintFn("chroma_sample_loc_type_top_field", uEV)
		d.FieldUintFn("chroma_sample_loc_type_bottom_field", uEV)
	}

	timingInfoPresentFlag := d.FieldBool("timing_info_present_flag")

	if timingInfoPresentFlag {
		d.FieldU32("num_units_in_tick")
		d.FieldU32("time_scale")
		d.FieldBool("fixed_frame_rate_flag")
	}
	nalHrdParametersPresentFlag := d.FieldBool("nal_hrd_parameters_present_flag")
	if nalHrdParametersPresentFlag {
		d.FieldStruct("nal_hrd_parameters", avcHdrParameters)
	}
	vclHrdParametersPresentFlag := d.FieldBool("vcl_hrd_parameters_present_flag")
	if vclHrdParametersPresentFlag {
		d.FieldStruct("vcl_hrd_parameters", avcHdrParameters)
	}
	if nalHrdParametersPresentFlag || vclHrdParametersPresentFlag {
		d.FieldBool("low_delay_hrd_flag")
	}
	d.FieldBool("pic_struct_present_flag")
	bitstreamRestrictionFlag := d.FieldBool("bitstream_restriction_flag")
	if bitstreamRestrictionFlag {
		d.FieldBool("motion_vectors_over_pic_boundaries_flag")
		d.FieldUintFn("max_bytes_per_pic_denom", uEV)
		d.FieldUintFn("max_bits_per_mb_denom", uEV)
		d.FieldUintFn("log2_max_mv_length_horizontal", uEV)
		d.FieldUintFn("log2_max_mv_length_vertical", uEV)
		d.FieldUintFn("max_num_reorder_frames", uEV)
		d.FieldUintFn("max_dec_frame_buffering", uEV)
	}
}

func avcHdrParameters(d *decode.D) {
	cpbCnt := d.FieldUintFn("cpb_cnt", uEV, scalar.UintActualAdd(1))
	d.FieldU4("bit_rate_scale")
	d.FieldU4("cpb_size_scale")
	d.FieldArray("sched_sels", func(d *decode.D) {
		for i := uint64(0); i < cpbCnt; i++ {
			d.FieldStruct("sched_sel", func(d *decode.D) {
				d.FieldUintFn("bit_rate_value", uEV, scalar.UintActualAdd(1))
				d.FieldUintFn("cpb_size_value", uEV, scalar.UintActualAdd(1))
				d.FieldBool("cbr_flag")
			})
		}
	})
	d.FieldU5("initial_cpb_removal_delay_length", scalar.UintActualAdd(1))
	d.FieldU5("cpb_removal_delay_length", scalar.UintActualAdd(1))
	d.FieldU5("dpb_output_delay_length", scalar.UintActualAdd(1))
	d.FieldU5("time_offset_length")
}

func avcSPSDecode(d *decode.D) any {
	profileIdc := d.FieldU8("profile_idc", avcProfileNames)
	d.FieldBool("constraint_set0_flag")
	d.FieldBool("constraint_set1_flag")
	d.FieldBool("constraint_set2_flag")
	d.FieldBool("constraint_set3_flag")
	d.FieldBool("constraint_set4_flag")
	d.FieldBool("constraint_set5_flag")
	d.FieldU2("reserved_zero_2bits")
	d.FieldU8("level_idc", avcLevelNames)
	d.FieldUintFn("seq_parameter_set_id", uEV)

	switch profileIdc {
	// TODO: ffmpeg has some more (legacy values?)
	case 100, 110, 122, 244, 44, 83, 86, 118, 128, 138, 139, 134, 135:
		chromaFormatIdc := d.FieldUintFn("chroma_format_idc", uEV, chromaFormatMap)
		if chromaFormatIdc == 3 {
			d.FieldBool("separate_colour_plane_flag")
		}

		d.FieldUintFn("bit_depth_luma", uEV, scalar.UintActualAdd(8))
		d.FieldUintFn("bit_depth_chroma", uEV, scalar.UintActualAdd(8))
		d.FieldBool("qpprime_y_zero_transform_bypass_flag")
		seqScalingMatrixPresentFlag := d.FieldBool("seq_scaling_matrix_present_flag")
		if seqScalingMatrixPresentFlag {
			// TODO:
			return nil
		}
	}

	d.FieldUintFn("log2_max_frame_num", uEV, scalar.UintActualAdd(4))

	picOrderCntType := d.FieldUintFn("pic_order_cnt_type", uEV)
	if picOrderCntType == 0 {
		d.FieldUintFn("log2_max_pic_order_cnt_lsb", uEV, scalar.UintActualAdd(4))
	} else if picOrderCntType == 1 {
		d.FieldBool("delta_pic_order_always_zero_flag")
		d.FieldSintFn("offset_for_non_ref_pic", sEV)
		d.FieldSintFn("offset_for_top_to_bottom_field", sEV)
		numRefFramesInPicOrderCntCycle := d.FieldUintFn("num_ref_frames_in_pic_order_cnt_cycle", uEV)
		d.FieldArray("offset_for_ref_frames", func(d *decode.D) {
			for i := uint64(0); i < numRefFramesInPicOrderCntCycle; i++ {
				sEV(d)
			}
		})
	}

	d.FieldUintFn("max_num_ref_frames", uEV)
	d.FieldBool("gaps_in_frame_num_value_allowed_flag")
	d.FieldUintFn("pic_width_in_mbs", uEV, scalar.UintActualAdd(1))
	d.FieldUintFn("pic_height_in_map_units", uEV, scalar.UintActualAdd(1))
	frameMbsOnlyFlag := d.FieldBool("frame_mbs_only_flag")
	if !frameMbsOnlyFlag {
		d.FieldBool("mb_adaptive_frame_field_flag")
	}
	d.FieldBool("direct_8x8_inference_flag")
	frameCroppingFlag := d.FieldBool("frame_cropping_flag")
	if frameCroppingFlag {
		d.FieldUintFn("frame_crop_left_offset", uEV)
		d.FieldUintFn("frame_crop_right_offset", uEV)
		d.FieldUintFn("frame_crop_top_offset", uEV)
		d.FieldUintFn("frame_crop_bottom_offset", uEV)
	}
	vuiParametersPresentFlag := d.FieldBool("vui_parameters_present_flag")
	if vuiParametersPresentFlag {
		d.FieldStruct("vui_parameters", avcVuiParameters)
	}

	d.FieldRawLen("rbsp_trailing_bits", d.BitsLeft())

	return nil
}
