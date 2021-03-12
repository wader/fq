package mpeg

// SO/IEC 13818-7 Part 7: Advanced Audio Coding (AAC)
// ISO/IEC 14496-3

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
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
	FILL          = 0x0
	FILL_DATA     = 0x1
	DATA_ELEMENT  = 0x2
	DYNAMIC_RANGE = 0xb
	SBR_DATA      = 0xd
	SBR_DATA_CRC  = 0xe
)

var ExtensionPayloadIDNames = map[uint64]string{
	FILL:          "FILL",
	FILL_DATA:     "FILL_DATA",
	DATA_ELEMENT:  "DATA_ELEMENT",
	DYNAMIC_RANGE: "DYNAMIC_RANGE",
	SBR_DATA:      "SBR_DATA",
	SBR_DATA_CRC:  "SBR_DATA_CRC",
}

const (
	ONLY_LONG_SEQUENCE   = 0x0
	LONG_START_SEQUENCE  = 0x1
	EIGHT_START_SEQUENCE = 0x2
	LONG_STOP_SEQUENCE   = 0x3
)

var windowSequnceNames = map[uint64]string{
	ONLY_LONG_SEQUENCE:   "ONLY_LONG_SEQUENCE",
	LONG_START_SEQUENCE:  "LONG_START_SEQUENCE",
	EIGHT_START_SEQUENCE: "EIGHT_START_SEQUENCE",
	LONG_STOP_SEQUENCE:   "LONG_STOP_SEQUENCE",
}

var windowSequnceNumWindows = map[uint64]int{
	ONLY_LONG_SEQUENCE:   1,
	LONG_START_SEQUENCE:  1,
	EIGHT_START_SEQUENCE: 8,
	LONG_STOP_SEQUENCE:   1,
}

func aacIcsInfo(d *decode.D) {

	d.FieldStructFn("ics_info", func(d *decode.D) {
		d.FieldU1("ics_reserved_bit")
		windowSequence, _ := d.FieldStringMapFn("window_sequence", windowSequnceNames, "", d.U2, decode.NumberDecimal)
		d.FieldU1("window_shape")
		switch windowSequence {
		case EIGHT_START_SEQUENCE:
			d.FieldU4("max_sfb")
			d.FieldU7("scale_factor_grouping")
		default:
			d.FieldU6("max_sfb")
			predictorDataPresent := d.FieldBool("predictor_data_present")
			if predictorDataPresent {

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
	})

}

func aacDecode(d *decode.D, in interface{}) interface{} {
	// TODO: seems tricky to know length of blocks
	// TODO: currently break when length is unknown
	d.FieldArrayFn("raw_data_blocks", func(d *decode.D) {
		seenTerm := false
		for !seenTerm {
			d.FieldStructFn("raw_data_block", func(d *decode.D) {
				se, _ := d.FieldStringMapFn("syntax_element", SyntaxElementNames, "", d.U3, decode.NumberDecimal)

				switch se {
				case FIL:
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
							case FILL:
								d.FieldU4("fill_nibble")
								d.FieldBitBufLen("fill_byte", 8*(int64(cnt)-1))
							}
						})
					})

					if d.ByteAlignBits() > 0 {
						d.FieldBitBufLen("byte_align", int64(d.ByteAlignBits()))
					}

				case SCE:
					d.FieldU4("element_instance_tag")
					d.FieldU8("global_gain")
					aacIcsInfo(d)

					if d.ByteAlignBits() > 0 {
						d.FieldBitBufLen("byte_align", int64(d.ByteAlignBits()))
					}
					seenTerm = true

				default:
					if d.ByteAlignBits() > 0 {
						d.FieldBitBufLen("data", int64(d.ByteAlignBits()))
					}
					fallthrough
				case TERM:
					seenTerm = true
				}

			})
		}
	})

	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
