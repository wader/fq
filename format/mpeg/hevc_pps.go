package mpeg

// https://www.itu.int/rec/T-REC-H.265

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(
		format.HEVC_PPS,
		&decode.Format{
			Description: "H.265/HEVC Picture Parameter Set",
			DecodeFn:    hevcPPSDecode,
		})
}

// H.265 page 36
func hevcPPSDecode(d *decode.D) any {
	d.FieldUintFn("pps_pic_parameter_set_id", uEV)
	d.FieldUintFn("pps_seq_parameter_set_id", uEV)
	d.FieldBool("dependent_slice_segments_enabled_flag")
	d.FieldBool("output_flag_present_flag")
	d.FieldU3("num_extra_slice_header_bits")
	d.FieldBool("sign_data_hiding_enabled_flag")
	d.FieldBool("cabac_init_present_flag")
	d.FieldUintFn("num_ref_idx_l0_default_active_minus1", uEV)
	d.FieldUintFn("num_ref_idx_l1_default_active_minus1", uEV)
	d.FieldSintFn("init_qp_minus26", sEV)
	d.FieldBool("constrained_intra_pred_flag")
	d.FieldBool("transform_skip_enabled_flag")
	cuQpDeltaEnabledFlag := d.FieldBool("cu_qp_delta_enabled_flag")
	if cuQpDeltaEnabledFlag {
		d.FieldUintFn("diff_cu_qp_delta_depth", uEV)
	}
	d.FieldSintFn("pps_cb_qp_offset", sEV)
	d.FieldSintFn("pps_cr_qp_offset", sEV)
	d.FieldBool("pps_slice_chroma_qp_offsets_present_flag")
	d.FieldBool("weighted_pred_flag")
	d.FieldBool("weighted_bipred_flag")
	d.FieldBool("transquant_bypass_enabled_flag")
	tilesEnabledFlag := d.FieldBool("tiles_enabled_flag")
	d.FieldBool("entropy_coding_sync_enabled_flag")
	if tilesEnabledFlag {
		numTileColumnsMinus1 := d.FieldUintFn("num_tile_columns_minus1", uEV)
		numTileRowsMinus1 := d.FieldUintFn("num_tile_rows_minus1", uEV)
		uniformSpacingFlag := d.FieldBool("uniform_spacing_flag")
		if !uniformSpacingFlag {
			d.FieldArray("column_widths", func(d *decode.D) {
				for i := uint64(0); i < numTileColumnsMinus1; i++ {
					d.FieldUintFn("column_width", uEV)
				}
			})
			d.FieldArray("row_heights", func(d *decode.D) {
				for i := uint64(0); i < numTileRowsMinus1; i++ {
					d.FieldUintFn("row_height", uEV)
				}
			})
		}
		d.FieldBool("loop_filter_across_tiles_enabled_flag")
	}
	d.FieldBool("pps_loop_filter_across_slices_enabled_flag")
	deblockingFilterControlPresentFlag := d.FieldBool("deblocking_filter_control_present_flag")
	if deblockingFilterControlPresentFlag {
		d.FieldBool("deblocking_filter_override_enabled_flag")
		ppsDeblockingFilterDisabledFlag := d.FieldBool("pps_deblocking_filter_disabled_flag")
		if !ppsDeblockingFilterDisabledFlag {
			d.FieldSintFn("pps_beta_offset_div2", sEV)
			d.FieldSintFn("pps_tc_offset_div2", sEV)
		}
	}
	ppsScalingListDataPresentFlag := d.FieldBool("pps_scaling_list_data_present_flag")
	if ppsScalingListDataPresentFlag {
		// TODO: scaling_list_data
		return nil
	}
	d.FieldBool("lists_modification_present_flag")
	d.FieldUintFn("log2_parallel_merge_level_minus2", uEV)
	d.FieldBool("slice_segment_header_extension_present_flag")
	ppsExtensionPresentFlag := d.FieldBool("pps_extension_present_flag")
	if ppsExtensionPresentFlag {
		d.FieldBool("pps_range_extension_flag")
		d.FieldBool("pps_multilayer_extension_flag")
		d.FieldBool("pps_3d_extension_flag")
		d.FieldBool("pps_scc_extension_flag")
		d.FieldU4("pps_extension_4bits")
	}

	// TODO: extensions

	return nil
}
