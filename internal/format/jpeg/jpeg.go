package jpeg

// https://www.w3.org/Graphics/JPEG/itu-t81.pdf

import (
	"fq/internal/decode"
	"log"
)

var Register = &decode.Register{
	Name: "jpeg",
	MIME: "",
	New: func(common decode.Common) decode.Decoder {
		return &Decoder{
			Common: common,
		}
	},
}

type marker struct {
	symbol      string
	description string
	standlone   bool
}

const (
	SOF0  = 0xc0
	SOF1  = 0xc1
	SOF2  = 0xc2
	SOF3  = 0xc3
	SOF5  = 0xc5
	SOF6  = 0xc6
	SOF7  = 0xc7
	JPG   = 0xc8
	SOF9  = 0xc9
	SOF10 = 0xca
	SOF11 = 0xcb
	SOF13 = 0xcd
	SOF14 = 0xce
	SOF15 = 0xcf
	DHT   = 0xc4
	DAC   = 0xcc
	RST0  = 0xd0
	RST1  = 0xd1
	RST2  = 0xd2
	RST3  = 0xd3
	RST4  = 0xd4
	RST5  = 0xd5
	RST6  = 0xd6
	RST7  = 0xd7
	SOI   = 0xd8
	EOI   = 0xd9
	SOS   = 0xda
	DQT   = 0xdb
	DNL   = 0xdc
	DRI   = 0xdd
	DHP   = 0xde
	EXP   = 0xdf
	APP0  = 0xe0
	APP1  = 0xe1
	APP2  = 0xe2
	APP3  = 0xe3
	APP4  = 0xe4
	APP5  = 0xe5
	APP6  = 0xe6
	APP7  = 0xe7
	APP8  = 0xe8
	APP9  = 0xe9
	APP10 = 0xea
	APP11 = 0xeb
	APP12 = 0xec
	APP13 = 0xed
	APP14 = 0xee
	APP15 = 0xef
	JPG0  = 0xf0
	JPG1  = 0xf1
	JPG2  = 0xf2
	JPG3  = 0xf3
	JPG4  = 0xf4
	JPG5  = 0xf5
	JPG6  = 0xf6
	JPG7  = 0xf7
	JPG8  = 0xf8
	JPG9  = 0xf9
	JPG10 = 0xfa
	JPG11 = 0xfb
	JPG12 = 0xfc
	JPG13 = 0xfd
	COM   = 0xfe
	TEM   = 0x01
)

var markers = map[uint]marker{
	SOF0:  {"SOF0", "Baseline DCT", false},
	SOF1:  {"SOF1", "Extended sequential DCT", false},
	SOF2:  {"SOF2", "Progressive DCT", false},
	SOF3:  {"SOF3", "Lossless (sequential)", false},
	SOF5:  {"SOF5", "Differential sequential DCT", false},
	SOF6:  {"SOF6", "Differential progressive DCT", false},
	SOF7:  {"SOF7", "Differential lossless (sequential)", false},
	JPG:   {"JPG", "Reserved for JPEG extensions", false},
	SOF9:  {"SOF9", "Extended sequential DCT", false},
	SOF10: {"SOF10", "Progressive DCT", false},
	SOF11: {"SOF11", "Lossless (sequential)", false},
	SOF13: {"SOF13", "Differential sequential DCT", false},
	SOF14: {"SOF14", "Differential progressive DCT", false},
	SOF15: {"SOF15", "Differential lossless (sequential)", false},
	DHT:   {"DHT", "Define Huffman table(s)", false},
	DAC:   {"DAC", "Define arithmetic coding conditioning(s)", false},
	RST0:  {"RST0", "Restart with modulo 8 count 0", true},
	RST1:  {"RST1", "Restart with modulo 8 count 1", true},
	RST2:  {"RST2", "Restart with modulo 8 count 2", true},
	RST3:  {"RST3", "Restart with modulo 8 count 3", true},
	RST4:  {"RST4", "Restart with modulo 8 count 4", true},
	RST5:  {"RST5", "Restart with modulo 8 count 5", true},
	RST6:  {"RST6", "Restart with modulo 8 count 6", true},
	RST7:  {"RST7", "Restart with modulo 8 count 7", true},
	SOI:   {"SOI", "Start of image ", true},
	EOI:   {"EOI", "End of image true", true},
	SOS:   {"SOS", "Start of scan", false},
	DQT:   {"DQT", "Define quantization table(s)", false},
	DNL:   {"DNL", "Define number of lines", false},
	DRI:   {"DRI", "Define restart interval", false},
	DHP:   {"DHP", "Define hierarchical progression", false},
	EXP:   {"EXP", "Expand reference component(s)", false},
	APP0:  {"APP0", "Reserved for application segments", false},
	APP1:  {"APP1", "Reserved for application segments", false},
	APP2:  {"APP2", "Reserved for application segments", false},
	APP3:  {"APP3", "Reserved for application segments", false},
	APP4:  {"APP4", "Reserved for application segments", false},
	APP5:  {"APP5", "Reserved for application segments", false},
	APP6:  {"APP6", "Reserved for application segments", false},
	APP7:  {"APP7", "Reserved for application segments", false},
	APP8:  {"APP8", "Reserved for application segments", false},
	APP9:  {"APP9", "Reserved for application segments", false},
	APP10: {"APP10", "Reserved for application segments", false},
	APP11: {"APP11", "Reserved for application segments", false},
	APP12: {"APP12", "Reserved for application segments", false},
	APP13: {"APP13", "Reserved for application segments", false},
	APP14: {"APP14", "Reserved for application segments", false},
	APP15: {"APP15", "Reserved for application segments", false},
	JPG0:  {"JPG0", "Reserved for JPEG extensions", false},
	JPG1:  {"JPG1", "Reserved for JPEG extensions", false},
	JPG2:  {"JPG2", "Reserved for JPEG extensions", false},
	JPG3:  {"JPG3", "Reserved for JPEG extensions", false},
	JPG4:  {"JPG4", "Reserved for JPEG extensions", false},
	JPG5:  {"JPG5", "Reserved for JPEG extensions", false},
	JPG6:  {"JPG6", "Reserved for JPEG extensions", false},
	JPG7:  {"JPG7", "Reserved for JPEG extensions", false},
	JPG8:  {"JPG8", "Reserved for JPEG extensions", false},
	JPG9:  {"JPG9", "Reserved for JPEG extensions", false},
	JPG10: {"JPG10", "Reserved for JPEG extensions", false},
	JPG11: {"JPG11", "Reserved for JPEG extensions", false},
	JPG12: {"JPG12", "Reserved for JPEG extensions", false},
	JPG13: {"JPG13", "Reserved for JPEG extensions", false},
	COM:   {"COM", "Comment", false},
	TEM:   {"TEM", "For temporary private use in arithmetic coding", true},
}

// 0x02-BF RES Reserved"},

// Decoder is a jpeg decoder
type Decoder struct {
	decode.Common
}

// Decode jpeg
func (d *Decoder) Decode(opts decode.Options) {
	for !d.End() {
		d.FieldNoneFn("marker", func() {
			d.FieldNoneFn("prefix", func() {
				for d.PeekBits(8) == 0xff {
					d.SeekRel(8)
				}
			})
			var markerCode uint
			var cm *marker
			d.FieldUFn("code", func() (uint64, decode.Format, string) {
				markerCode = uint(d.U8())
				if m, ok := markers[markerCode]; ok {
					cm = &m
					return uint64(markerCode), decode.FormatDecimal, m.symbol
				}
				return uint64(markerCode), decode.FormatDecimal, "RES"
			})

			log.Printf("markerCode: %#+v\n", markerCode)
			if cm.standlone {
				return
			}

			markerLen := d.FieldU16("length")

			log.Printf("markerLen: %#+v\n", markerLen)

			switch markerCode {
			case SOF0, SOF1, SOF2, SOF3, SOF5, SOF6, SOF7, SOF9, SOF10, SOF11:
				d.FieldU8("P")
				d.FieldU16("Y")
				d.FieldU16("X")
				nf := d.FieldU8("Nf")
				for i := uint64(0); i < nf; i++ {
					d.FieldNoneFn("frame_component", func() {
						d.FieldU8("C")
						d.FieldU4("H")
						d.FieldU4("V")
						d.FieldU8("Tq")
					})
				}
			case SOS:
				ns := d.FieldU8("Ns")
				for i := uint64(0); i < ns; i++ {
					d.FieldNoneFn("scan_component", func() {
						d.FieldU8("Cs")
						d.FieldU4("Td")
						d.FieldU4("Ta")
					})
				}
				d.FieldU8("Ss")
				d.FieldU8("Se")
				d.FieldU4("Ah")
				d.FieldU4("Al")
			default:
				d.FieldBytesLen("data", markerLen-2)
			}

		})

	}

}
