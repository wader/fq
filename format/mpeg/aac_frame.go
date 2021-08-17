package mpeg

// one AAC frame or "raw data block"

// ISO/IEC 13818-7 Part 7: Advanced Audio Coding (AAC)
// ISO/IEC 14496-3

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.AAC_FRAME,
		Description: "Advanced Audio Coding frame",
		DecodeFn:    aacDecode,
	})
}

const (
	SCE  = 0b000
	CPE  = 0b001
	CCE  = 0b010
	LFE  = 0b011
	DSE  = 0b100
	PCE  = 0b101
	FIL  = 0b110
	TERM = 0b111
)

var SyntaxElementNames = map[uint64]string{
	SCE:  "SCE",
	CPE:  "CPE",
	CCE:  "CCE",
	LFE:  "LFE",
	DSE:  "DSE",
	PCE:  "PCE",
	FIL:  "FIL",
	TERM: "TERM",
}

const (
	EXT_FILL          = 0x0
	EXT_FILL_DATA     = 0x1
	EXT_DATA_ELEMENT  = 0x2
	EXT_DYNAMIC_RANGE = 0xb
	EXT_SBR_DATA      = 0xd
	EXT_SBR_DATA_CRC  = 0xe
)

var ExtensionPayloadIDNames = map[uint64]string{
	EXT_FILL:          "EXT_FILL",
	EXT_FILL_DATA:     "EXT_FILL_DATA",
	EXT_DATA_ELEMENT:  "EXT_DATA_ELEMENT",
	EXT_DYNAMIC_RANGE: "EXT_DYNAMIC_RANGE",
	EXT_SBR_DATA:      "EXT_SBR_DATA",
	EXT_SBR_DATA_CRC:  "EXT_SBR_DATA_CRC",
}

const (
	ONLY_LONG_SEQUENCE   = 0x0
	LONG_START_SEQUENCE  = 0x1
	EIGHT_SHORT_SEQUENCE = 0x2
	LONG_STOP_SEQUENCE   = 0x3
)

var windowSequenceNames = map[uint64]string{
	ONLY_LONG_SEQUENCE:   "ONLY_LONG_SEQUENCE",
	LONG_START_SEQUENCE:  "LONG_START_SEQUENCE",
	EIGHT_SHORT_SEQUENCE: "EIGHT_SHORT_SEQUENCE",
	LONG_STOP_SEQUENCE:   "LONG_STOP_SEQUENCE",
}

var windowSequenceNumWindows = map[int]int{
	ONLY_LONG_SEQUENCE:   1,
	LONG_START_SEQUENCE:  1,
	EIGHT_SHORT_SEQUENCE: 8,
	LONG_STOP_SEQUENCE:   1,
}

func aacLTPData(d *decode.D, objectType int, windowSequence int) {
	switch objectType {
	case format.MPEGAudioObjectTypeER_AAC_LD:
		// TODO:
	default:
		d.FieldU11("ltp_lag")
		d.FieldU3("ltp_coef")

		_ = windowSequenceNumWindows[windowSequence]

	}
}

func aacICSInfo(d *decode.D, objectType int) {
	d.FieldU1("ics_reserved_bit")
	windowSequence, _ := d.FieldStringMapFn("window_sequence", windowSequenceNames, "", d.U2, decode.NumberDecimal)
	d.FieldU1("window_shape")
	switch windowSequence {
	case EIGHT_SHORT_SEQUENCE:
		d.FieldU4("max_sfb")
		d.FieldU7("scale_factor_grouping")
	default:
		maxSFB := d.FieldU6("max_sfb")
		predictorDataPresent := d.FieldBool("predictor_data_present")
		if predictorDataPresent {
			switch objectType {
			case format.MPEGAudioObjectTypeMain: // 1
				predictorReset := d.FieldBool("predictor_reset")
				if predictorReset {
					d.FieldU5("predictor_reset_group_number")
				}
				d.FieldU5("predictor_reset_group_number")
				// TODO: min(max_sfb, PRED_SFB_MAX)
				// TODO: array?
				d.FieldBitBufLen("prediction_used", int64(maxSFB))
			default:
				ltpDataPresent := d.FieldBool("ltp_data_present")
				if ltpDataPresent {
					aacLTPData(d, objectType, int(windowSequence))
				}
			}
		}

	}

	// 		;
	// 		if (window_sequence == EIGHT_SHORT_SEQUENCE) {
	// 		max_sfb; scale_factor_grouping;
	// 		} }
	// 		else {
	// 		ltp_data_present;
	// 		if (ltp_data_present) {
	// 		ltp_data(); }
	// 		if (common_window) {
	// 		ltp_data_present;
	// 		LICENSED TO MECON Limited. - RANCHI/BANGALORE,
	// 		FOR INTERNAL USE AT THIS LOCATION ONLY, SUPPLIED BY BOOK SUPPLY BUREAU.
	// 		if (ltp_data_present) {
	// 		ltp_data(); }
	// 		} }
	// 		} }
	// }

}

func aacIndividualChannelStream(d *decode.D, objectType int, commonWindow bool, scaleFlag bool) {
	d.FieldU8("global_gain")
	if !commonWindow && !scaleFlag {
		d.FieldStructFn("ics_info", func(d *decode.D) {
			aacICSInfo(d, objectType)
		})
	}
}

func aacSingleChannelElement(d *decode.D, objectType int) {
	d.FieldU4("element_instance_tag")
	aacIndividualChannelStream(d, objectType, false, false)
}

func aacProgramConfigElement(d *decode.D, ascStartPos int64) {
	d.FieldU4("element_instance_tag")
	d.FieldU2("object_type")
	d.FieldU4("sampling_frequency_index")
	numFrontChannelElements := d.FieldU4("num_front_channel_elements")
	numSideChannelElements := d.FieldU4("num_side_channel_elements")
	numBackChannelElements := d.FieldU4("num_back_channel_elements")
	numLfeChannelElements := d.FieldU2("num_lfe_channel_elements")
	numAssocDataElements := d.FieldU3("num_assoc_data_elements")
	numValidCcElements := d.FieldU4("num_valid_cc_elements")
	monoMixdownPresent := d.FieldBool("mono_mixdown_present")
	if monoMixdownPresent {
		d.FieldU4("mono_mixdown_element_number")
	}
	stereoMixdownPresent := d.FieldBool("stereo_mixdown_present")
	if stereoMixdownPresent {
		d.FieldU4("stereo_mixdown_element_number")
	}
	matrixMixdownIdxPresent := d.FieldBool("matrix_mixdown_idx_present")
	if matrixMixdownIdxPresent {
		d.FieldU2("matrix_mixdown_idx")
		d.FieldBool("pseudo_surround_enable")
	}
	d.FieldArrayFn("front_channel_elements", func(d *decode.D) {
		for i := uint64(0); i < numFrontChannelElements; i++ {
			d.FieldStructFn("front_channel_element", func(d *decode.D) {
				d.FieldBool("is_cpe")
				d.FieldU4("tag_select")
			})
		}
	})
	d.FieldArrayFn("side_channel_elements", func(d *decode.D) {
		for i := uint64(0); i < numSideChannelElements; i++ {
			d.FieldStructFn("side_channel_element", func(d *decode.D) {
				d.FieldBool("is_cpe")
				d.FieldU4("tag_select")
			})
		}
	})
	d.FieldArrayFn("back_channel_elements", func(d *decode.D) {
		for i := uint64(0); i < numBackChannelElements; i++ {
			d.FieldStructFn("back_channel_element", func(d *decode.D) {
				d.FieldBool("is_cpe")
				d.FieldU4("tag_select")
			})
		}
	})
	d.FieldArrayFn("lfe_channel_elements", func(d *decode.D) {
		for i := uint64(0); i < numLfeChannelElements; i++ {
			d.FieldStructFn("lfe_channel_element", func(d *decode.D) {
				d.FieldU4("tag_select")
			})
		}
	})
	d.FieldArrayFn("assoc_data_elements", func(d *decode.D) {
		for i := uint64(0); i < numAssocDataElements; i++ {
			d.FieldStructFn("assoc_data_element", func(d *decode.D) {
				d.FieldU4("tag_select")
			})
		}
	})
	d.FieldArrayFn("valid_cc_elements", func(d *decode.D) {
		for i := uint64(0); i < numValidCcElements; i++ {
			d.FieldStructFn("valid_cc_element", func(d *decode.D) {
				d.FieldU1("cc_element_is_ind_sw")
				d.FieldU4("valid_cc_element_tag_select")
			})
		}
	})

	byteAlignBits := (8 - ((d.Pos() + ascStartPos) & 0x7)) & 0x7
	d.FieldBitBufLen("byte_alignment", byteAlignBits)
	commentFieldBytes := d.FieldU8("comment_field_bytes")
	d.FieldUTF8("comment_field", int(commentFieldBytes))
}

func aacFillElement(d *decode.D) {
	var cnt uint64
	d.FieldStructFn("cnt", func(d *decode.D) {
		count := d.FieldU4("count")
		cnt = count
		if cnt == 15 {
			escCount := d.FieldU8("esc_count")
			cnt += escCount - 1
		}
	})
	d.FieldValueU("payload_length", cnt, "")

	d.FieldStructFn("extension_payload", func(d *decode.D) {
		d.DecodeLenFn(int64(cnt)*8, func(d *decode.D) {

			extensionType, _ := d.FieldStringMapFn("extension_type", ExtensionPayloadIDNames, "Unknown", d.U4, decode.NumberDecimal)

			// d.FieldU("align4", 2)

			switch extensionType {
			case EXT_FILL:
				d.FieldU4("fill_nibble")
				d.FieldBitBufLen("fill_byte", 8*(int64(cnt)-1))
			}
		})
	})
}

func aacDecode(d *decode.D, in interface{}) interface{} {
	var objectType int
	if afi, ok := in.(format.AACFrameIn); ok {
		objectType = afi.ObjectType
	}
	// objectType = 2

	// TODO: seems tricky to know length of blocks
	// TODO: currently break when length is unknown
	d.FieldArrayFn("elements", func(d *decode.D) {
		seenTerm := false
		for !seenTerm {
			d.FieldStructFn("element", func(d *decode.D) {
				se, _ := d.FieldStringMapFn("syntax_element", SyntaxElementNames, "", d.U3, decode.NumberDecimal)

				switch se {
				case FIL:
					aacFillElement(d)

				case SCE:
					aacSingleChannelElement(d, objectType)
					seenTerm = true

				case PCE:
					aacProgramConfigElement(d, 0)
					seenTerm = true

				default:
					fallthrough
				case TERM:
					seenTerm = true
				}
			})
		}

		if d.ByteAlignBits() > 0 {
			d.FieldBitBufLen("byte_align", int64(d.ByteAlignBits()))
		}
	})

	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
