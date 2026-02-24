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
			DefaultInArg: format.AVC_SPS_Info{
				SeparateColourPlaneFlag:      false,
				Log2MaxFrameNum:              4,
				FrameMbsOnlyFlag:             true,
				PicOrderCntType:              0,
				Log2MaxPicOrderCntLsb:        4,
				NalHrdParametersPresent:      false,
				VclHrdParametersPresent:      false,
				InitialCpbRemovalDelayLength: 0,
				CpbRemovalDelayLength:        0,
				DpbOutputDelayLength:         0,
				TimeOffsetLength:             0,
			},
		})
}

const (
	avcSEIBufferingPeriod      = 0
	avcSEIUserDataUnregistered = 5
	avcSEIPicTiming            = 1
)

var seiNames = scalar.UintMapSymStr{
	avcSEIBufferingPeriod:      "buffering_period",
	avcSEIPicTiming:            "pic_timing",
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

type picStructEntry struct {
	name       string
	numClockTS int
}

type picStructMap []picStructEntry

var picStructMapEntries = picStructMap{
	0: {name: "frame", numClockTS: 1},
	1: {name: "top_field", numClockTS: 1},
	2: {name: "bottom_field", numClockTS: 1},
	3: {name: "top_field,bottom_field,in_that_order", numClockTS: 2},
	4: {name: "bottom_field,top_field,in_that_order", numClockTS: 2},
	5: {name: "top_field,bottom_field,top_field_repeated,in_that_order", numClockTS: 3},
	6: {name: "bottom_field,top_field,bottom_field_repeated,in_that_order", numClockTS: 3},
	7: {name: "frame_doubling", numClockTS: 2},
	8: {name: "frame_tripling", numClockTS: 3},
}

func (m picStructMap) MapUint(s scalar.Uint) (scalar.Uint, error) {
	if len(m) < int(s.Actual) {
		s.Sym = m[s.Actual].name
	}
	return s, nil
}

func avcSEIDecode(d *decode.D) any {
	var ai format.AVC_NALU_In
	d.ArgAs(&ai)

	payloadType := d.FieldUintFn("payload_type", func(d *decode.D) uint64 { return ffSum(d) }, seiNames)
	payloadSize := d.FieldUintFn("payload_size", func(d *decode.D) uint64 { return ffSum(d) })

	d.FramedFn(int64(payloadSize)*8, func(d *decode.D) {
		switch payloadType {
		case avcSEIBufferingPeriod:
			d.FieldUintFn("seq_parameter_set_id", uEV)
			if ai.NalHrdParametersPresent {
				d.FieldArray("initial_cpb_removal_delays", func(d *decode.D) {
					d.FieldStruct("initial_cpb_removal_delay", func(d *decode.D) {
						d.FieldU("initial_cpb_removal_delay", int(ai.InitialCpbRemovalDelayLength))
						d.FieldU("initial_cpb_removal_delay_offset", int(ai.InitialCpbRemovalDelayLength))
					})
				})
			}
			if ai.VclHrdParametersPresent {
				d.FieldArray("initial_cpb_removal_delays", func(d *decode.D) {
					d.FieldStruct("initial_cpb_removal_delay", func(d *decode.D) {
						d.FieldU("initial_cpb_removal_delay", int(ai.InitialCpbRemovalDelayLength))
						d.FieldU("initial_cpb_removal_delay_offset", int(ai.InitialCpbRemovalDelayLength))
					})
				})
			}
		case avcSEIUserDataUnregistered:
			d.FieldRawLen("uuid", 16*8, userDataUnregisteredNames)
		case avcSEIPicTiming:
			if ai.NalHrdParametersPresent || ai.VclHrdParametersPresent {
				d.FieldU("cpb_removal_delay", int(ai.CpbRemovalDelayLength))
				d.FieldU("dpb_output_delay", int(ai.DpbOutputDelayLength))
			}
			pic_struct_present_flag := d.FieldBool("pic_struct_present_flag")
			picStruct := 0
			if pic_struct_present_flag {
				picStruct = int(d.FieldU4("pic_struct", picStructMapEntries))
			}
			numClockTS := 0
			if picStruct < len(picStructMapEntries) {
				numClockTS = picStructMapEntries[picStruct].numClockTS
			}
			if numClockTS > 0 {
				d.FieldArray("clocks", func(d *decode.D) {
					for i := 0; i < numClockTS; i++ {
						d.FieldStruct("clock", func(d *decode.D) {
							clock_timestamp_flag := d.FieldBool("clock_timestamp_flag")
							if clock_timestamp_flag {
								d.FieldU2("ct_type")
								d.FieldBool("nuit_field_based_flag")
								d.FieldU5("counting_type")
								full_timestamp_flag := d.FieldBool("full_timestamp_flag")
								d.FieldBool("discontinuity_flag")
								d.FieldBool("cnt_dropped_flag")
								d.FieldU8("nframes")
								d.FieldU2("ct_type")
								if full_timestamp_flag {
									d.FieldU6("seconds_value")
									d.FieldU5("minutes_value")
									d.FieldU5("hours_value")
								} else {
									seconds_flag := d.FieldBool("seconds_flag")
									if seconds_flag {
										d.FieldU6("seconds_value")
										minutes_flag := d.FieldBool("minutes_flag")
										if minutes_flag {
											d.FieldU5("minutes_value")
											hours_flag := d.FieldBool("minutes_flag")
											if hours_flag {
												d.FieldU5("hours_value")
											}
										}
									}
								}
								if ai.TimeOffsetLength > 0 {
									d.FieldS5("time_offset")
								}
							}
						})
					}
				})
			}
		}
		d.FieldRawLen("data", d.BitsLeft())
	})

	d.FieldRawLen("trailing_bits", d.BitsLeft())

	return nil
}
