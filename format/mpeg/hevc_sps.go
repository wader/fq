package mpeg

// https://www.itu.int/rec/T-REC-H.265

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(
		format.HEVC_SPS,
		&decode.Format{
			Description: "H.265/HEVC Sequence Parameter Set",
			DecodeFn:    hevcSPSDecode,
		})
}

func profileLayerDecode(d *decode.D, prefix string, profilePresent bool, levelPresent bool, isSublayer bool) {
	if profilePresent {
		d.FieldU2(prefix + "profile_space")
		d.FieldU1(prefix + "tier_flag")
		generalProfileIdc := d.FieldU5(prefix + "profile_idc")
		var generalProfileCompatibilityFlags [32]bool
		d.FieldArray(prefix+"profile_compatibility_flags", func(d *decode.D) {
			for j := 0; j < 32; j++ {
				generalProfileCompatibilityFlags[j] = d.FieldBool(prefix + "profile_compatibility_flag")
			}
		})
		d.FieldBool(prefix + "progressive_source_flag")
		d.FieldBool(prefix + "interlaced_source_flag")
		d.FieldBool(prefix + "non_packed_constraint_flag")
		d.FieldBool(prefix + "frame_only_constraint_flag")
		if generalProfileIdc == 4 || generalProfileCompatibilityFlags[4] ||
			generalProfileIdc == 5 || generalProfileCompatibilityFlags[5] ||
			generalProfileIdc == 6 || generalProfileCompatibilityFlags[6] ||
			generalProfileIdc == 7 || generalProfileCompatibilityFlags[7] ||
			generalProfileIdc == 8 || generalProfileCompatibilityFlags[8] ||
			generalProfileIdc == 9 || generalProfileCompatibilityFlags[9] ||
			generalProfileIdc == 10 || generalProfileCompatibilityFlags[10] {
			d.FieldBool(prefix + "max_12bit_constraint_flag")
			d.FieldBool(prefix + "max_10bit_constraint_flag")
			d.FieldBool(prefix + "max_8bit_constraint_flag")
			d.FieldBool(prefix + "max_422chroma_constraint_flag")
			d.FieldBool(prefix + "max_420chroma_constraint_flag")
			d.FieldBool(prefix + "max_monochrome_constraint_flag")
			d.FieldBool(prefix + "intra_constraint_flag")
			d.FieldBool(prefix + "one_picture_only_constraint_flag")
			d.FieldBool(prefix + "lower_bit_rate_constraint_flag")
			if generalProfileIdc == 5 || generalProfileCompatibilityFlags[5] ||
				(!isSublayer &&
					(generalProfileIdc == 9 || generalProfileCompatibilityFlags[9] ||
						generalProfileIdc == 10 || generalProfileCompatibilityFlags[10])) {
				d.FieldBool(prefix + "max_14bit_constraint_flag")
				d.FieldU33(prefix + "reserved_zero_33bits")
			} else {
				d.FieldU34(prefix + "reserved_zero_34bits")
			}
		} else {
			d.FieldU43(prefix + "reserved_zero_43bits")
		}
		if (generalProfileIdc >= 1 && generalProfileIdc <= 5) ||
			generalProfileIdc == 9 ||
			generalProfileCompatibilityFlags[1] || generalProfileCompatibilityFlags[2] ||
			generalProfileCompatibilityFlags[3] || generalProfileCompatibilityFlags[4] ||
			generalProfileCompatibilityFlags[5] || generalProfileCompatibilityFlags[9] {
			d.FieldBool(prefix + "inbld_flag")
		} else {
			d.FieldBool(prefix + "reserved_zero_bit")
		}
	}
	if levelPresent {
		d.FieldU8(prefix + "level_idc")
	}
}

// H.265 page 41
func profileTierLevelDecode(d *decode.D, profilePresentFlag bool, maxNumSubLayersMinus1 uint64) {
	profileLayerDecode(d, "general_", profilePresentFlag, true, false)
	subLayerProfilePresentFlags := make([]bool, maxNumSubLayersMinus1)
	subLayerLevelPresentFlags := make([]bool, maxNumSubLayersMinus1)
	d.FieldArray("sub_layer_presents", func(d *decode.D) {
		for i := uint64(0); i < maxNumSubLayersMinus1; i++ {
			d.FieldStruct("sub_layer_present", func(d *decode.D) {
				subLayerProfilePresentFlags[i] = d.FieldBool("sub_layer_profile_present_flag")
				subLayerLevelPresentFlags[i] = d.FieldBool("sub_layer_level_present_flag")
			})
		}
	})
	if maxNumSubLayersMinus1 > 0 {
		d.FieldArray("reserved_zero_2bits", func(d *decode.D) {
			for i := maxNumSubLayersMinus1; i < 8; i++ {
				d.FieldU2("reserved_zero_2bits")
			}
		})
	}
	d.FieldArray("sub_layers", func(d *decode.D) {
		for i := uint64(0); i < maxNumSubLayersMinus1; i++ {
			d.FieldStruct("sub_layer", func(d *decode.D) {
				profileLayerDecode(d, "", subLayerProfilePresentFlags[i], subLayerLevelPresentFlags[i], true)
			})
		}
	})
}

func hevcSubLayerHrdParameters(d *decode.D, subPicHrdParamsPresentFlag bool, cpbCntMinus1 int) {
	for i := 0; i <= cpbCntMinus1; i++ {
		d.FieldStruct("parameters", func(d *decode.D) {
			d.FieldUintFn("bit_rate_value_minus1", uEV)
			d.FieldUintFn("cpb_size_value_minus1", uEV)
			if subPicHrdParamsPresentFlag {
				d.FieldUintFn("cpb_size_du_value_minus1", uEV)
				d.FieldUintFn("bit_rate_du_value_minus1", uEV)
			}
			d.FieldBool("cbr_flag")
		})
	}
}

func hevcHrdParameters(d *decode.D, commonInfPresentFlag bool, maxNumSubLayersMinus1 uint64) {
	var nalHrdParametersPresentFlag bool
	var vclHrdParametersPresentFlag bool
	var subPicHrdParamsPresentFlag bool
	if commonInfPresentFlag {
		nalHrdParametersPresentFlag = d.FieldBool("nal_hrd_parameters_present_flag")
		vclHrdParametersPresentFlag = d.FieldBool("vcl_hrd_parameters_present_flag")
		if nalHrdParametersPresentFlag && vclHrdParametersPresentFlag {
			subPicHrdParamsPresentFlag = d.FieldBool("sub_pic_hrd_params_present_flag")
			if subPicHrdParamsPresentFlag {
				d.FieldU8("tick_divisor_minus2")
				d.FieldU5("sar_wdu_cpb_removal_delay_increment_length_minus1idth")
				d.FieldBool("sub_pic_cpb_params_in_pic_timing_sei_flag")
				d.FieldU5("dpb_output_delay_du_length_minus1")
			}
			d.FieldU8("tick_divisor_minus2")
			d.FieldU8("tick_divisor_minus2")
			if subPicHrdParamsPresentFlag {
				d.FieldU4("cpb_size_du_scale")
			}
			d.FieldU5("initial_cpb_removal_delay_length_minus1")
			d.FieldU5("au_cpb_removal_delay_length_minus1")
			d.FieldU5("dpb_output_delay_length_minus1")
		}
	}
	d.FieldArray("sub_layers", func(d *decode.D) {
		for i := uint64(0); i < maxNumSubLayersMinus1; i++ {
			d.FieldStruct("sub_layer", func(d *decode.D) {
				fixedPicRateGeneralFlag := d.FieldBool("fixed_pic_rate_general_flag")
				var fixedPicRateWithinCvsFlag bool
				if !fixedPicRateGeneralFlag {
					fixedPicRateWithinCvsFlag = d.FieldBool("fixed_pic_rate_within_cvs_flag")
				}
				var lowDelayHrdFlag bool
				if fixedPicRateWithinCvsFlag {
					d.FieldUintFn("elemental_duration_in_tc_minus1", uEV)
				} else {
					lowDelayHrdFlag = d.FieldBool("low_delay_hrd_flag")
				}
				var cpbCntMinus1 int
				if !lowDelayHrdFlag {
					cpbCntMinus1 = int(d.FieldUintFn("cpb_cnt_minus1", uEV))
				}
				if nalHrdParametersPresentFlag {
					hevcSubLayerHrdParameters(d, subPicHrdParamsPresentFlag, cpbCntMinus1)
				}
				if vclHrdParametersPresentFlag {
					hevcSubLayerHrdParameters(d, subPicHrdParamsPresentFlag, cpbCntMinus1)
				}
			})
		}
	})
}

func hevcVuiParameters(d *decode.D, spsMaxSubLayersMinus1 uint64) {
	aspectRatioInfoPresentFlag := d.FieldBool("aspect_ratio_info_present_flag")
	if aspectRatioInfoPresentFlag {
		aspectRatioIdc := d.FieldU8("aspect_ratio_idc", avcAspectRatioIdcMap)
		const extendedSAR = 255
		if aspectRatioIdc == extendedSAR {
			d.FieldU16("sar_width")
			d.FieldU16("sar_height")
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

	d.FieldBool("neutral_chroma_indication_flag")
	d.FieldBool("field_seq_flag")
	d.FieldBool("frame_field_info_present_flag")
	defaultDisplayWindowFlag := d.FieldBool("default_display_window_flag")
	if defaultDisplayWindowFlag {
		d.FieldUintFn("def_disp_win_left_offset", uEV)
		d.FieldUintFn("def_disp_win_right_offset", uEV)
		d.FieldUintFn("def_disp_win_top_offset", uEV)
		d.FieldUintFn("def_disp_win_bottom_offset", uEV)
	}

	vuiTimingInfoPresentFlag := d.FieldBool("vui_timing_info_present_flag")
	if vuiTimingInfoPresentFlag {
		d.FieldU32("vui_num_units_in_tick")
		d.FieldU32("vui_time_scale")
		vuiPocProportionalToTimingFlag := d.FieldBool("vui_poc_proportional_to_timing_flag")
		if vuiPocProportionalToTimingFlag {
			d.FieldUintFn("vui_num_ticks_poc_diff_one_minus1", uEV)
		}
		vuiHrdParametersPresentFlag := d.FieldBool("vui_hrd_parameters_present_flag")
		if vuiHrdParametersPresentFlag {
			hevcHrdParameters(d, true, spsMaxSubLayersMinus1)
		}
	}

	bitstreamRestrictionFlag := d.FieldBool("bitstream_restriction_flag")
	if bitstreamRestrictionFlag {
		d.FieldBool("tiles_fixed_structure_flag")
		d.FieldBool("motion_vectors_over_pic_boundaries_flag")
		d.FieldBool("restricted_ref_pic_lists_flag")
		d.FieldUintFn("min_spatial_segmentation_idc", uEV)
		d.FieldUintFn("max_bytes_per_pic_denom", uEV)
		d.FieldUintFn("max_bits_per_min_cu_denom", uEV)
		d.FieldUintFn("log2_max_mv_length_horizontal", uEV)
		d.FieldUintFn("log2_max_mv_length_vertical", uEV)
	}
}

// H.265 page 34
func hevcSPSDecode(d *decode.D) any {
	d.FieldU4("sps_video_parameter_set_id")
	spsMaxSubLayersMinus1 := d.FieldU3("sps_max_sub_layers_minus1")
	d.FieldBool("sps_temporal_id_nesting_flag")
	profileTierLevelDecode(d, true, spsMaxSubLayersMinus1)
	d.FieldUintFn("sps_seq_parameter_set_id", uEV)
	chromaFormatIdc := d.FieldUintFn("chroma_format_idc", uEV, chromaFormatMap)
	if chromaFormatIdc == 3 {
		d.FieldBool("separate_colour_plane_flag")
	}
	d.FieldUintFn("pic_width_in_luma_samples", uEV)
	d.FieldUintFn("pic_height_in_luma_samples", uEV)
	conformanceWindowFlag := d.FieldBool("conformance_window_flag")
	if conformanceWindowFlag {
		d.FieldUintFn("conf_win_left_offset", uEV)
		d.FieldUintFn("conf_win_right_offset", uEV)
		d.FieldUintFn("conf_win_top_offset", uEV)
		d.FieldUintFn("conf_win_bottom_offset", uEV)
	}
	d.FieldUintFn("bit_depth_luma_minus8", uEV)
	d.FieldUintFn("bit_depth_chroma_minus8", uEV)
	d.FieldUintFn("log2_max_pic_order_cnt_lsb_minus4", uEV)
	spsSubLayerOrderingInfoPresentFlag := d.FieldBool("sps_sub_layer_ordering_info_present_flag")
	d.FieldArray("sps_sub_layer_ordering_infos", func(d *decode.D) {
		i := spsMaxSubLayersMinus1
		if spsSubLayerOrderingInfoPresentFlag {
			i = 0
		}
		for ; i <= spsMaxSubLayersMinus1; i++ {
			d.FieldStruct("sps_sub_layer_ordering_info", func(d *decode.D) {
				d.FieldUintFn("sps_max_dec_pic_buffering_minus1", uEV)
				d.FieldUintFn("sps_max_num_reorder_pics", uEV)
				d.FieldUintFn("sps_max_latency_increase_plus1", uEV)
			})
		}
	})
	d.FieldUintFn("log2_min_luma_coding_block_size_minus3", uEV)
	d.FieldUintFn("log2_diff_max_min_luma_coding_block_size", uEV)
	d.FieldUintFn("log2_min_luma_transform_block_size_minus2", uEV)
	d.FieldUintFn("log2_diff_max_min_luma_transform_block_size", uEV)
	d.FieldUintFn("max_transform_hierarchy_depth_inter", uEV)
	d.FieldUintFn("max_transform_hierarchy_depth_intra", uEV)
	scalingListEnabledFlag := d.FieldBool("scaling_list_enabled_flag")
	if scalingListEnabledFlag {
		spsScalingListDataPresentFlag := d.FieldBool("sps_scaling_list_data_present_flag")
		if spsScalingListDataPresentFlag {
			// TODO: scaling_list_data
			return nil
		}
	}
	d.FieldBool("amp_enabled_flag")
	d.FieldBool("sample_adaptive_offset_enabled_flag")
	pcmEnabledFlag := d.FieldBool("pcm_enabled_flag")
	if pcmEnabledFlag {
		d.FieldU4("pcm_sample_bit_depth_luma_minus1")
		d.FieldU4("pcm_sample_bit_depth_chroma_minus1")
		d.FieldUintFn("log2_min_pcm_luma_coding_block_size_minus3", uEV)
		d.FieldUintFn("log2_diff_max_min_pcm_luma_coding_block_size", uEV)
		d.FieldBool("pcm_loop_filter_disabled_flag")
	}
	numShortTermRefPicSets := d.FieldUintFn("num_short_term_ref_pic_sets", uEV)
	if numShortTermRefPicSets > 0 {
		// TODO
		return nil
	}
	longTermRefPicsPresentFlag := d.FieldBool("long_term_ref_pics_present_flag")
	if longTermRefPicsPresentFlag {
		// TODO
		return nil
	}
	d.FieldBool("sps_temporal_mvp_enabled_flag")
	d.FieldBool("strong_intra_smoothing_enabled_flag")
	vuiParametersPresentFlag := d.FieldBool("vui_parameters_present_flag")
	if vuiParametersPresentFlag {
		d.FieldStruct("vui_parameters", func(d *decode.D) { hevcVuiParameters(d, spsMaxSubLayersMinus1) })
	}
	spsExtensionPresentFlag := d.FieldBool("sps_extension_present_flag")
	if spsExtensionPresentFlag {
		d.FieldU1("sps_range_extension_flag")
		d.FieldU1("sps_multilayer_extension_flag")
		d.FieldU1("sps_3d_extension_flag")
		d.FieldU1("sps_scc_extension_flag")
		d.FieldU4("sps_extension_4bits")
	}

	// TODO

	return nil
}
