package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.AVC_PPS,
		&decode.Format{
			Description: "H.264/AVC Picture Parameter Set",
			DecodeFn:    avcPPSDecode,
		})
}

func moreRBSPData(d *decode.D) bool {
	l := d.BitsLeft()
	return l >= 8 && d.PeekUintBits(8) != 0b1000_0000
}

func avcPPSDecode(d *decode.D) any {
	d.FieldUintFn("pic_parameter_set_id", uEV)
	d.FieldUintFn("seq_parameter_set_id", uEV)
	d.FieldBool("entropy_coding_mode_flag")
	d.FieldBool("bottom_field_pic_order_in_frame_present_flag")
	numSliceGroups := d.FieldUintFn("num_slice_groups", uEV, scalar.UintActualAdd(1))
	if numSliceGroups > 1 {
		sliceGroupMapType := d.FieldUintFn("slice_group_map_type", uEV)
		switch sliceGroupMapType {
		case 0:
			d.FieldArray("slice_groups", func(d *decode.D) {
				for i := uint64(0); i < numSliceGroups; i++ {
					d.FieldUintFn("slice_group", uEV)
				}
			})
		case 2:
			d.FieldArray("slice_groups", func(d *decode.D) {
				for i := uint64(0); i < numSliceGroups; i++ {
					d.FieldStruct("slice_group", func(d *decode.D) {
						d.FieldUintFn("top_left", uEV)
						d.FieldUintFn("bottom_right", uEV)
					})
				}
			})
		case 3, 4, 5:
			d.FieldArray("slice_groups", func(d *decode.D) {
				for i := uint64(0); i < numSliceGroups; i++ {
					d.FieldStruct("slice_group", func(d *decode.D) {
						d.FieldBool("change_direction_flag")
						d.FieldUintFn("change_rate", uEV, scalar.UintActualAdd(1))
					})
				}
			})
		case 6:
			picSizeInMapUnits := d.FieldUintFn("pic_size_in_map_units", uEV, scalar.UintActualAdd(1))
			for i := uint64(0); i < picSizeInMapUnits; i++ {
				d.FieldStruct("slice_group", func(d *decode.D) {
					d.FieldBool("id")
				})
			}
		}
	}

	d.FieldUintFn("num_ref_idx_l0_default_active", uEV, scalar.UintActualAdd(1))
	d.FieldUintFn("num_ref_idx_l1_default_active", uEV, scalar.UintActualAdd(1))
	d.FieldBool("weighted_pred_flag")
	d.FieldU2("weighted_bipred_idc")
	d.FieldSintFn("pic_init_qp", sEV, scalar.SintActualAdd(26))
	d.FieldSintFn("pic_init_qs", sEV, scalar.SintActualAdd(26))
	d.FieldSintFn("chroma_qp_index_offset", sEV)
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
		d.FieldSintFn("second_chroma_qp_index_offset", sEV)
	} else {
		d.FieldBool("rbsp_stop_one_bit")
	}

	d.FieldRawLen("rbsp_trailing_bits", d.BitsLeft())

	return nil
}
