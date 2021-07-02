package mpeg

import (
	"fq/format"
	"fq/format/all/all"
	"fq/pkg/decode"
)

func init() {
	all.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_PPS,
		Description: "H.264/AVC Picture Parameter Set",
		DecodeFn:    avcPPSDecode,
	})
}

// TODO:
// ffmpeg does: get_bits_left(gb) > 0 && show_bits(gb, 8) != 0x80
// we do: there is more than trailing rbsp left
func moreRBSPData(d *decode.D) bool {
	l := d.BitsLeft()
	return l > 8 //|| (l == 8 && d.PeekBits(8) != 1)
}

func avcPPSDecode(d *decode.D, in interface{}) interface{} {
	fieldUEV(d, "pic_parameter_set_id")
	fieldUEV(d, "seq_parameter_set_id")
	d.FieldBool("entropy_coding_mode_flag")
	d.FieldBool("bottom_field_pic_order_in_frame_present_flag")
	numSliceGroupsMinus1 := fieldUEV(d, "num_slice_groups_minus1")
	if numSliceGroupsMinus1 > 0 {
		sliceGroupMapType := fieldUEV(d, "slice_group_map_type")
		switch sliceGroupMapType {
		case 0:
			d.FieldArrayFn("slice_groups", func(d *decode.D) {
				for i := uint64(0); i <= numSliceGroupsMinus1; i++ {
					fieldUEV(d, "slice_group")
				}
			})
		case 2:
			d.FieldArrayFn("slice_groups", func(d *decode.D) {
				for i := uint64(0); i <= numSliceGroupsMinus1; i++ {
					d.FieldStructFn("slice_group", func(d *decode.D) {
						fieldUEV(d, "top_left")
						fieldUEV(d, "bottom_right")
					})
				}
			})
		case 3, 4, 5:
			d.FieldArrayFn("slice_groups", func(d *decode.D) {
				for i := uint64(0); i <= numSliceGroupsMinus1; i++ {
					d.FieldStructFn("slice_group", func(d *decode.D) {
						d.FieldBool("change_direction_flag")
						fieldUEV(d, "change_rate_minus1")
					})
				}
			})
		case 6:
			picSizeInMapUnitsMinus1 := fieldUEV(d, "pic_size_in_map_units_minus1")
			for i := uint64(0); i <= picSizeInMapUnitsMinus1; i++ {
				d.FieldStructFn("slice_group", func(d *decode.D) {
					d.FieldBool("id")
				})
			}
		}
	}

	fieldUEV(d, "num_ref_idx_l0_default_active_minus1")
	fieldUEV(d, "num_ref_idx_l1_default_active_minus1")
	d.FieldBool("weighted_pred_flag")
	d.FieldU2("weighted_bipred_idc")
	fieldSEV(d, "pic_init_qp_minus26") /* relative to 26 */
	fieldSEV(d, "pic_init_qs_minus26") /* relative to 26 */
	fieldSEV(d, "chroma_qp_index_offset")
	d.FieldBool("deblocking_filter_control_present_flag")
	d.FieldBool("constrained_intra_pred_flag")
	d.FieldBool("redundant_pic_cnt_present_flag")

	// TODO: more_data() is there non-zero bits left?
	if moreRBSPData(d) {
		d.FieldBool("transform_8x8_mode_flag")
		picScalingMatrixPresentFlag := d.FieldBool("pic_scaling_matrix_present_flag")
		if picScalingMatrixPresentFlag {
			d.FieldArrayFn("pic_scaling_list", func(d *decode.D) {
				for i := 0; i < 6; i++ {
					d.Bool()
				}
			})
		}
		fieldSEV(d, "second_chroma_qp_index_offset")
	}

	d.FieldBitBufLen("rbsp_trailing_bits", d.BitsLeft())

	return nil
}
