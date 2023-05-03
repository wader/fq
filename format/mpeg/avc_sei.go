package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.AVC_SEI,
		&decode.Format{
			Description: "H.264/AVC Supplemental Enhancement Information",
			DecodeFn:    avcSEIDecode,
		})
}

const (
	avcSEIUserDataUnregistered = 5
)

var seiNames = scalar.UintMapSymStr{
	0:                          "buffering_period",
	1:                          "pic_timing",
	2:                          "pan_scan_rect",
	3:                          "filler_payload",
	4:                          "user_data_registered_itu_t_t35",
	avcSEIUserDataUnregistered: "user_data_unregistered",
	6:                          "recovery_point",
	7:                          "dec_ref_pic_marking_repetition",
	8:                          "spare_pic",
	9:                          "scene_info",
	10:                         "sub_seq_info",
	11:                         "sub_seq_layer_characteristics",
	12:                         "sub_seq_characteristics",
	13:                         "full_frame_freeze",
	14:                         "full_frame_freeze_release",
	15:                         "full_frame_snapshot",
	16:                         "progressive_refinement_segment_start",
	17:                         "progressive_refinement_segment_end",
	18:                         "motion_constrained_slice_group_set",
	19:                         "film_grain_characteristics",
	20:                         "deblocking_filter_display_preference",
	21:                         "stereo_video_info",
	22:                         "post_filter_hint",
	23:                         "tone_mapping_info",
	24:                         "scalability_info",
	25:                         "sub_pic_scalable_layer",
	26:                         "non_required_layer_rep",
	27:                         "priority_layer_info",
	28:                         "layers_not_present",
	29:                         "layer_dependency_change",
	30:                         "scalable_nesting",
	31:                         "base_layer_temporal_hrd",
	32:                         "quality_layer_integrity_check",
	33:                         "redundant_pic_property",
	34:                         "tl0_dep_rep_index",
	35:                         "tl_switching_point",
	36:                         "parallel_decoding_info",
	37:                         "mvc_scalable_nesting",
	38:                         "view_scalability_info",
	39:                         "multiview_scene_info",
	40:                         "multiview_acquisition_info",
	41:                         "non_required_view_component",
	42:                         "view_dependency_change",
	43:                         "operation_points_not_present",
	44:                         "base_view_temporal_hrd",
	45:                         "frame_packing_arrangement",
	46:                         "multiview_view_position",
	47:                         "display_orientation",
	48:                         "mvcd_scalable_nesting",
	49:                         "mvcd_view_scalability_info",
	50:                         "depth_representation_info",
	51:                         "three_dimensional_reference_displays_info",
	52:                         "depth_timing",
	53:                         "depth_sampling_info",
	54:                         "constrained_depth_parameter_set_identifier",
	56:                         "green_metadata",
	137:                        "mastering_display_colour_volume",
	181:                        "alternative_depth_info",
}

var (
	x264Bytes = [16]byte{0xdc, 0x45, 0xe9, 0xbd, 0xe6, 0xd9, 0x48, 0xb7, 0x96, 0x2c, 0xd8, 0x20, 0xd9, 0x23, 0xee, 0xef}
)

var userDataUnregisteredNames = scalar.RawBytesMap{
	{Bytes: x264Bytes[:], Scalar: scalar.BitBuf{Sym: "x264"}},
}

// sum bytes until < 0xff
func ffSum(d *decode.D) uint64 {
	var s uint64
	for {
		b := d.U8()
		s += b
		if b < 0xff {
			break
		}
	}
	return s
}

func avcSEIDecode(d *decode.D) any {
	payloadType := d.FieldUintFn("payload_type", func(d *decode.D) uint64 { return ffSum(d) }, seiNames)
	payloadSize := d.FieldUintFn("payload_size", func(d *decode.D) uint64 { return ffSum(d) })

	d.FramedFn(int64(payloadSize)*8, func(d *decode.D) {
		switch payloadType {
		case avcSEIUserDataUnregistered:
			d.FieldRawLen("uuid", 16*8, userDataUnregisteredNames)
		}
		d.FieldRawLen("data", d.BitsLeft())
	})

	d.FieldRawLen("rbsp_trailing_bits", d.BitsLeft())

	return nil
}
