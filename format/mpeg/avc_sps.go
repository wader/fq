package mpeg

// ISO/IEC 14496-15 AVC file format, 5.3.3.1.2 Syntax
// ISO_IEC_14496-10 AVC

import (
	"fq/format"
	"fq/format/all/all"
	"fq/pkg/decode"
)

func init() {
	all.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_SPS,
		Description: "H.264/AVC Sequence Parameter Set",
		DecodeFn:    avcSPSDecode,
	})
}

func avcVuiParameters(d *decode.D) {
	aspectRatioInfoPresentFlag := d.FieldBool("aspect_ratio_info_present_flag")
	if aspectRatioInfoPresentFlag {
		aspectRatioIdc := d.FieldU8("aspect_ratio_idc")
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
		d.FieldU3("video_format")
		d.FieldBool("video_full_range_flag")
		colourDescriptionPresentFlag := d.FieldBool("colour_description_present_flag")
		if colourDescriptionPresentFlag {
			d.FieldU8("colour_primaries")
			d.FieldU8("transfer_characteristics")
			d.FieldU8("matrix_coefficients")
		}
	}
	chromaLocInfoPresentFlag := d.FieldBool("chroma_loc_info_present_flag")
	if chromaLocInfoPresentFlag {
		fieldUEV(d, "chroma_sample_loc_type_top_field")
		fieldUEV(d, "chroma_sample_loc_type_bottom_field")
	}

	timingInfoPresentFlag := d.FieldBool("timing_info_present_flag")

	if timingInfoPresentFlag {
		d.FieldU32("num_units_in_tick")
		d.FieldU32("time_scale")
		d.FieldBool("fixed_frame_rate_flag")
	}
	nalHrdParametersPresentFlag := d.FieldBool("nal_hrd_parameters_present_flag")
	if nalHrdParametersPresentFlag {
		d.FieldStructFn("nal_hrd_parameters", avcHdrParameters)
	}
	vclHrdParametersPresentFlag := d.FieldBool("vcl_hrd_parameters_present_flag")
	if vclHrdParametersPresentFlag {
		d.FieldStructFn("vcl_hrd_parameters", avcHdrParameters)
	}
	if nalHrdParametersPresentFlag || vclHrdParametersPresentFlag {
		d.FieldBool("low_delay_hrd_flag")
	}
	d.FieldBool("pic_struct_present_flag")
	bitstreamRestrictionFlag := d.FieldBool("bitstream_restriction_flag")
	if bitstreamRestrictionFlag {
		d.FieldBool("motion_vectors_over_pic_boundaries_flag")
		fieldUEV(d, "max_bytes_per_pic_denom")
		fieldUEV(d, "max_bits_per_mb_denom")
		fieldUEV(d, "log2_max_mv_length_horizontal")
		fieldUEV(d, "log2_max_mv_length_vertical")
		fieldUEV(d, "max_num_reorder_frames")
		fieldUEV(d, "max_dec_frame_buffering")
	}
}

func avcHdrParameters(d *decode.D) {
	cpbCntMinus1 := fieldUEV(d, "cpb_cnt_minus1")
	_ = cpbCntMinus1
	d.FieldU4("bit_rate_scale")
	d.FieldU4("cpb_size_scale")
	d.FieldArrayFn("sched_sels", func(d *decode.D) {
		for i := uint64(0); i <= cpbCntMinus1; i++ {
			d.FieldStructFn("sched_sel", func(d *decode.D) {
				fieldUEV(d, "bit_rate_value_minus1")
				fieldUEV(d, "cpb_size_value_minus1")
				d.FieldBool("cbr_flag")
			})
		}
	})
	d.FieldU5("initial_cpb_removal_delay_length_minus1")
	d.FieldU5("cpb_removal_delay_length_minus1")
	d.FieldU5("dpb_output_delay_length_minus1")
	d.FieldU5("time_offset_length")
}

func avcSPSDecode(d *decode.D, in interface{}) interface{} {
	profileIdc, _ := d.FieldStringMapFn("profile_idc", avcProfileNames, "Unknown", d.U8, decode.NumberDecimal)
	d.FieldBool("constraint_set0_flag")
	d.FieldBool("constraint_set1_flag")
	d.FieldBool("constraint_set2_flag")
	d.FieldBool("constraint_set3_flag")
	d.FieldBool("constraint_set4_flag")
	d.FieldBool("constraint_set5_flag")
	d.FieldU2("reserved_zero_2bits")
	d.FieldStringMapFn("level_idc", avcLevelNames, "Unknown", d.U8, decode.NumberDecimal)
	fieldUEV(d, "seq_parameter_set_id")

	switch profileIdc {
	// TODO: ffmpeg has some more (legacy values?)
	case 100, 110, 122, 244, 44, 83, 86, 118, 128, 138, 139, 134, 135:
		chromaFormatIdc := fieldUEV(d, "chroma_format_idc")
		if chromaFormatIdc == 3 {
			d.FieldBool("separate_colour_plane_flag")
		}

		fieldUEV(d, "bit_depth_luma_minus8")
		fieldUEV(d, "bit_depth_chroma_minus8")
		d.FieldBool("qpprime_y_zero_transform_bypass_flag")
		seqScalingMatrixPresentFlag := d.FieldBool("seq_scaling_matrix_present_flag")
		if seqScalingMatrixPresentFlag {
			// TODO:
		}
	}

	fieldUEV(d, "log2_max_frame_num_minus4")

	picOrderCntType := fieldUEV(d, "pic_order_cnt_type")
	if picOrderCntType == 0 {
		fieldUEV(d, "log2_max_pic_order_cnt_lsb_minus4")
	} else if picOrderCntType == 1 {
		d.FieldBool("delta_pic_order_always_zero_flag")
		fieldSEV(d, "offset_for_non_ref_pic")
		fieldSEV(d, "offset_for_top_to_bottom_field")
		numRefFramesInPicOrderCntCycle := fieldUEV(d, "num_ref_frames_in_pic_order_cnt_cycle")
		d.FieldArrayFn("offset_for_ref_frames", func(d *decode.D) {
			for i := uint64(0); i < numRefFramesInPicOrderCntCycle; i++ {
				sEV(d)
			}
		})
	}

	fieldUEV(d, "max_num_ref_frames")
	d.FieldBool("gaps_in_frame_num_value_allowed_flag")
	fieldUEV(d, "pic_width_in_mbs_minus1")
	fieldUEV(d, "pic_height_in_map_units_minus1")
	frameMbsOnlyFlag := d.FieldBool("frame_mbs_only_flag")
	if !frameMbsOnlyFlag {
		d.FieldBool("mb_adaptive_frame_field_flag")
	}
	d.FieldBool("direct_8x8_inference_flag")
	frameCroppingFlag := d.FieldBool("frame_cropping_flag")
	if frameCroppingFlag {
		fieldUEV(d, "frame_crop_left_offset")
		fieldUEV(d, "frame_crop_right_offset")
		fieldUEV(d, "frame_crop_top_offset")
		fieldUEV(d, "frame_crop_bottom_offset")
	}
	vuiParametersPresentFlag := d.FieldBool("vui_parameters_present_flag")
	if vuiParametersPresentFlag {
		d.FieldStructFn("vui_parameters", avcVuiParameters)
	}

	d.FieldBitBufLen("rbsp_trailing_bits", d.BitsLeft())

	return nil
}
