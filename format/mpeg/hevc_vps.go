package mpeg

// https://www.itu.int/rec/T-REC-H.265

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.HEVC_VPS,
		Description: "H.265/HEVC Video Parameter Set",
		DecodeFn:    hevcVPSDecode,
	})
}

// H.265 page 33
func hevcVPSDecode(d *decode.D, _ any) any {
	d.FieldU4("vps_video_parameter_set_id")
	d.FieldBool("vps_base_layer_internal_flag")
	d.FieldBool("vps_base_layer_available_flag")
	d.FieldU6("vps_max_layers_minus1")
	vpsMaxSubLayersMinus1 := d.FieldU3("vps_max_sub_layers_minus1")
	d.FieldBool("vps_temporal_id_nesting_flag")
	d.FieldU16("vps_reserved_0xffff_16bits")
	profileTierLevelDecode(d, true, vpsMaxSubLayersMinus1)
	vpsSubLayerOrderingInfoPresentFlag := d.FieldBool("vps_sub_layer_ordering_info_present_flag")
	d.FieldArray("vps_sub_layer_ordering_infos", func(d *decode.D) {
		i := vpsMaxSubLayersMinus1
		if vpsSubLayerOrderingInfoPresentFlag {
			i = 0
		}
		for ; i <= vpsMaxSubLayersMinus1; i++ {
			d.FieldStruct("sps_sub_layer_ordering_info", func(d *decode.D) {
				d.FieldUFn("sps_max_dec_pic_buffering_minus1", uEV)
				d.FieldUFn("sps_max_num_reorder_pics", uEV)
				d.FieldUFn("sps_max_latency_increase_plus1", uEV)
			})
		}
	})
	vpsMaxLayerID := d.FieldU6("vps_max_layer_id")
	vpsNumLayerSetsMinus1 := d.FieldUFn("vps_num_layer_sets_minus1", uEV)
	d.FieldArray("layer_id_included_sets_flags", func(d *decode.D) {
		for i := uint64(0); i <= vpsNumLayerSetsMinus1; i++ {
			d.FieldArray("layer_id_included_sets_flags", func(d *decode.D) {
				for j := uint64(0); j <= vpsMaxLayerID; j++ {
					d.FieldBool("layer_id_included_flag_sets_flag")
				}
			})
		}
	})
	vpsTimingInfoPresentFlag := d.FieldBool("vps_timing_info_present_flag")
	if vpsTimingInfoPresentFlag {
		d.FieldU32("vps_num_units_in_tick")
		d.FieldU32("vps_time_scale")
		vpsPocProportionalToTimingFlag := d.FieldBool("vps_poc_proportional_to_timing_flag")
		if vpsPocProportionalToTimingFlag {
			d.FieldUFn("vps_num_ticks_poc_diff_one_minus1", uEV)
		}
		vpsHrdParametersPresentFlag := d.FieldBool("vps_hrd_parameters_present_flag")
		if vpsHrdParametersPresentFlag {
			hevcHrdParameters(d, true, vpsMaxSubLayersMinus1)
		}
	}
	// TODO:

	return nil
}
