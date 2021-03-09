package mpeg

// SO/IEC 13818-7 Part 7: Advanced Audio Coding (AAC)

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AAC_FRAME,
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

func aacDecode(d *decode.D, in interface{}) interface{} {
	// TODO: multple blocks
	d.FieldArrayFn("raw_data_blocks", func(d *decode.D) {
		//		for {
		d.FieldStructFn("raw_data_block", func(d *decode.D) {
			se, _ := d.FieldStringMapFn("syntax_element", SyntaxElementNames, "", d.U3)

			switch se {
			case FIL:
				cnt := d.FieldUFn("cnt", func() (uint64, decode.DisplayFormat, string) {
					cnt := d.U4()
					if cnt == 15 {
						return cnt + d.FieldU8("length_escape") - 1, decode.NumberDecimal, ""
					}
					return cnt, decode.NumberDecimal, ""
				})

				d.FieldStructFn("extension_payload", func(d *decode.D) {
					d.DecodeLenFn(int64(cnt)*8, func(d *decode.D) {

						extensionType, _ := d.FieldStringMapFn("extension_type", ExtensionPayloadIDNames, "Unknown", d.U4)
						switch extensionType {
						case FILL:
							d.FieldBitBufLen("other_bits", 8*(int64(cnt)-1)+4)
						}

					})
				})
				// case SCE:
				// 	d.FieldU4("element_instance_tag")
				// 	d.FieldU8("global_gain")
				// case TERM:
			}

			// if d.ByteAlignBits() > 0 {
			// 	d.FieldBitBufLen("byte_align", int64(d.ByteAlignBits()))
			// }
			//return // TODO:
		})
		//		}
	})

	return nil
}
