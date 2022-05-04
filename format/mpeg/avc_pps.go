package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.AVC_PPS,
		Description: "H.264/AVC Picture Parameter Set",
		DecodeFn:    avcPPSDecode,
	})
}

func moreRBSPData(d *decode.D) bool {
	l := d.BitsLeft()
	return l >= 8 && d.PeekBits(8) != 0b1000_0000
}

func avcPPSDecode(d *decode.D, in interface{}) interface{} {
	d.FieldUFn("pic_parameter_set_id", uEV)
	d.FieldUFn("seq_parameter_set_id", uEV)
	d.FieldBool("entropy_coding_mode_flag")
	d.FieldBool("bottom_field_pic_order_in_frame_present_flag")
	numSliceGroups := d.FieldUFn("num_slice_groups", uEV, scalar.UAdd(1))
	if numSliceGroups > 1 {
		sliceGroupMapType := d.FieldUFn("slice_group_map_type", uEV)
		switch sliceGroupMapType {
		case 0:
			d.FieldArray("slice_groups", func(d *decode.D) {
				for i := uint64(0); i < numSliceGroups; i++ {
					d.FieldUFn("slice_group", uEV)
				}
			})
		case 2:
			d.FieldArray("slice_groups", func(d *decode.D) {
				for i := uint64(0); i < numSliceGroups; i++ {
					d.FieldStruct("slice_group", func(d *decode.D) {
						d.FieldUFn("top_left", uEV)
						d.FieldUFn("bottom_right", uEV)
					})
				}
			})
		case 3, 4, 5:
			d.FieldArray("slice_groups", func(d *decode.D) {
				for i := uint64(0); i < numSliceGroups; i++ {
					d.FieldStruct("slice_group", func(d *decode.D) {
						d.FieldBool("change_direction_flag")
						d.FieldUFn("change_rate", uEV, scalar.UAdd(1))
					})
				}
			})
		case 6:
			picSizeInMapUnits := d.FieldUFn("pic_size_in_map_units", uEV, scalar.UAdd(1))
			for i := uint64(0); i < picSizeInMapUnits; i++ {
				d.FieldStruct("slice_group", func(d *decode.D) {
					d.FieldBool("id")
				})
			}
		}
	}

	d.FieldUFn("num_ref_idx_l0_default_active", uEV, scalar.UAdd(1))
	d.FieldUFn("num_ref_idx_l1_default_active", uEV, scalar.UAdd(1))
	d.FieldBool("weighted_pred_flag")
	d.FieldU2("weighted_bipred_idc")
	d.FieldSFn("pic_init_qp", sEV, scalar.SAdd(26))
	d.FieldSFn("pic_init_qs", sEV, scalar.SAdd(26))
	d.FieldSFn("chroma_qp_index_offset", sEV)
	d.FieldBool("deblocking_filter_control_present_flag")
	d.FieldBool("constrained_intra_pred_flag")
	d.FieldBool("redundant_pic_cnt_present_flag")

	if moreRBSPData(d) {
		d.FieldBool("transform_8x8_mode_flag")
		picScalingMatrixPresentFlag := d.FieldBool("pic_scaling_matrix_present_flag")
		if picScalingMatrixPresentFlag {
			d.FieldArray("pic_scaling_list", func(d *decode.D) {
				for i := 0; i < 6; i++ {
					d.Bool()
				}
			})
		}
		d.FieldSFn("second_chroma_qp_index_offset", sEV)
	} else {
		d.FieldBool("rbsp_stop_one_bit")
	}

	d.FieldRawLen("rbsp_trailing_bits", d.BitsLeft())

	return nil
}
