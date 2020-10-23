package aac

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:      "aac_frame",
		DecodeFn:  aacDecode,
		SkipProbe: true,
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

var ExtenionPayloadIDNames = map[uint64]string{
	FILL:          "FILL",
	FILL_DATA:     "FILL_DATA",
	DATA_ELEMENT:  "DATA_ELEMENT",
	DYNAMIC_RANGE: "DYNAMIC_RANGE",
	SBR_DATA:      "SBR_DATA",
	SBR_DATA_CRC:  "SBR_DATA_CRC",
}

func aacDecode(d *decode.Common) interface{} {
	se, _ := d.FieldStringMapFn("syntax_element", SyntaxElementNames, "", d.U3)
	elementId := d.FieldU4("element_id")

	switch se {
	case FIL:
		filLength := elementId
		if filLength == 15 {
			filLength += d.FieldU8("length_escape")
		}

		d.FieldStringMapFn("type", ExtenionPayloadIDNames, "", d.U4)

		d.SeekRel(int64(filLength)*8 - 4)

	}

	/*
		d.FieldU4("sampling_frequency_index")
		d.FieldU4("channel_configuration")
		d.FieldU1("frame_length_flag")
		d.FieldU1("depends_on_core_coder")
		d.FieldU1("extension_flag")
	*/

	return nil
}
