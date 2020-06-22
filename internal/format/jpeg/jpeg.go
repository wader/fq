package jpeg

// https://www.w3.org/Graphics/JPEG/itu-t81.pdf

import (
	"fq/internal/decode"
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

var markers = map[uint]marker{
	0xc0: {"SOF0", "Baseline DCT", false},
	0xc1: {"SOF1", "Extended sequential DCT", false},
	0xc2: {"SOF2", "Progressive DCT", false},
	0xc3: {"SOF3", "Lossless (sequential)", false},
	0xc5: {"SOF5", "Differential sequential DCT", false},
	0xc6: {"SOF6", "Differential progressive DCT", false},
	0xc7: {"SOF7", "Differential lossless (sequential)", false},
	0xc8: {"JPG", "Reserved for JPEG extensions", false},
	0xc9: {"SOF9", "Extended sequential DCT", false},
	0xca: {"SOF10", "Progressive DCT", false},
	0xcb: {"SOF11", "Lossless (sequential)", false},
	0xcd: {"SOF13", "Differential sequential DCT", false},
	0xce: {"SOF14", "Differential progressive DCT", false},
	0xcf: {"SOF15", "Differential lossless (sequential)", false},
	0xc4: {"DHT", "Define Huffman table(s)", false},
	0xcc: {"DAC", "Define arithmetic coding conditioning(s)", false},
	0xd0: {"RST0", "Restart with modulo 8 count 0", true},
	0xd1: {"RST1", "Restart with modulo 8 count 1", true},
	0xd2: {"RST2", "Restart with modulo 8 count 2", true},
	0xd3: {"RST3", "Restart with modulo 8 count 3", true},
	0xd4: {"RST4", "Restart with modulo 8 count 4", true},
	0xd5: {"RST5", "Restart with modulo 8 count 5", true},
	0xd6: {"RST6", "Restart with modulo 8 count 6", true},
	0xd7: {"RST7", "Restart with modulo 8 count 7", true},
	0xd8: {"SOI", "Start of image ", true},
	0xd9: {"EOI", "End of image true", true},
	0xda: {"SOS", "Start of scan", false},
	0xdb: {"DQT", "Define quantization table(s)", false},
	0xdc: {"DNL", "Define number of lines", false},
	0xdd: {"DRI", "Define restart interval", false},
	0xde: {"DHP", "Define hierarchical progression", false},
	0xdf: {"EXP", "Expand reference component(s)", false},
	0xe0: {"APP0", "Reserved for application segments", false},
	0xe1: {"APP1", "Reserved for application segments", false},
	0xe2: {"APP2", "Reserved for application segments", false},
	0xe3: {"APP3", "Reserved for application segments", false},
	0xe4: {"APP4", "Reserved for application segments", false},
	0xe5: {"APP5", "Reserved for application segments", false},
	0xe6: {"APP6", "Reserved for application segments", false},
	0xe7: {"APP7", "Reserved for application segments", false},
	0xe8: {"APP8", "Reserved for application segments", false},
	0xe9: {"APP9", "Reserved for application segments", false},
	0xea: {"APPa", "Reserved for application segments", false},
	0xeb: {"APPb", "Reserved for application segments", false},
	0xec: {"APPc", "Reserved for application segments", false},
	0xed: {"APPd", "Reserved for application segments", false},
	0xee: {"APPe", "Reserved for application segments", false},
	0xef: {"APPf", "Reserved for application segments", false},
	0xf0: {"JPG0", "Reserved for JPEG extensions", false},
	0xf1: {"JPG1", "Reserved for JPEG extensions", false},
	0xf2: {"JPG2", "Reserved for JPEG extensions", false},
	0xf3: {"JPG3", "Reserved for JPEG extensions", false},
	0xf4: {"JPG4", "Reserved for JPEG extensions", false},
	0xf5: {"JPG5", "Reserved for JPEG extensions", false},
	0xf6: {"JPG6", "Reserved for JPEG extensions", false},
	0xf7: {"JPG7", "Reserved for JPEG extensions", false},
	0xf8: {"JPG8", "Reserved for JPEG extensions", false},
	0xf9: {"JPG9", "Reserved for JPEG extensions", false},
	0xfa: {"JPGa", "Reserved for JPEG extensions", false},
	0xfb: {"JPGb", "Reserved for JPEG extensions", false},
	0xfc: {"JPGc", "Reserved for JPEG extensions", false},
	0xfd: {"JPGd", "Reserved for JPEG extensions", false},
	0xfe: {"COM", "Comment", false},
	0x01: {"TEM", "For temporary private use in arithmetic coding", true},
}

// 0x02-BF RES Reserved"},

// Decoder is a jpeg decoder
type Decoder struct {
	decode.Common
}

// Decode jpeg
func (d *Decoder) Decode(opts decode.Options) {
	for !d.End() {
		var cm *marker
		d.FieldNoneFn("marker", func() {
			d.FieldNoneFn("prefix", func() {
				for d.PeekBits(8) == 0xff {
					d.SeekRel(8)
				}
			})
			d.FieldUFn("code", func() (uint64, decode.Format, string) {
				n := d.U8()
				if m, ok := markers[uint(n)]; ok {
					cm = &m
					return n, decode.FormatDecimal, m.symbol
				}
				return n, decode.FormatDecimal, "RES"
			})
			markerLen := d.FieldU16("length")

			d.FieldBytesLen("data", markerLen-2)

			_ = cm
		})

		break
	}

}
