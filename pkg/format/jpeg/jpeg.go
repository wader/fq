package jpeg

// https://www.w3.org/Graphics/JPEG/itu-t81.pdf
// TODO: exif https://www.exif.org/Exif2-2.PDF

import (
	"fmt"
	"fq/pkg/bitbuf"
	"fq/pkg/decode"
	"fq/pkg/format"
)

var tiffImage []*decode.Format

var File = format.MustRegister(&decode.Format{
	Name:   "jpeg",
	Groups: []string{"image"},
	MIMEs:  []string{"image/jpeg"},
	New:    func() decode.Decoder { return &FileDecoder{} },
	Deps: []decode.Dep{
		{Names: []string{"tiff"}, Formats: &tiffImage},
	},
})

type marker struct {
	symbol      string
	description string
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
	SOF0:  {"SOF0", "Baseline DCT"},
	SOF1:  {"SOF1", "Extended sequential DCT"},
	SOF2:  {"SOF2", "Progressive DCT"},
	SOF3:  {"SOF3", "Lossless (sequential)"},
	SOF5:  {"SOF5", "Differential sequential DCT"},
	SOF6:  {"SOF6", "Differential progressive DCT"},
	SOF7:  {"SOF7", "Differential lossless (sequential)"},
	JPG:   {"JPG", "Reserved for JPEG extensions"},
	SOF9:  {"SOF9", "Extended sequential DCT"},
	SOF10: {"SOF10", "Progressive DCT"},
	SOF11: {"SOF11", "Lossless (sequential)"},
	SOF13: {"SOF13", "Differential sequential DCT"},
	SOF14: {"SOF14", "Differential progressive DCT"},
	SOF15: {"SOF15", "Differential lossless (sequential)"},
	DHT:   {"DHT", "Define Huffman table(s)"},
	DAC:   {"DAC", "Define arithmetic coding conditioning(s)"},
	RST0:  {"RST0", "Restart with modulo 8 count 0"},
	RST1:  {"RST1", "Restart with modulo 8 count 1"},
	RST2:  {"RST2", "Restart with modulo 8 count 2"},
	RST3:  {"RST3", "Restart with modulo 8 count 3"},
	RST4:  {"RST4", "Restart with modulo 8 count 4"},
	RST5:  {"RST5", "Restart with modulo 8 count 5"},
	RST6:  {"RST6", "Restart with modulo 8 count 6"},
	RST7:  {"RST7", "Restart with modulo 8 count 7"},
	SOI:   {"SOI", "Start of image "},
	EOI:   {"EOI", "End of image true"},
	SOS:   {"SOS", "Start of scan"},
	DQT:   {"DQT", "Define quantization table(s)"},
	DNL:   {"DNL", "Define number of lines"},
	DRI:   {"DRI", "Define restart interval"},
	DHP:   {"DHP", "Define hierarchical progression"},
	EXP:   {"EXP", "Expand reference component(s)"},
	APP0:  {"APP0", "Reserved for application segments"},
	APP1:  {"APP1", "Reserved for application segments"},
	APP2:  {"APP2", "Reserved for application segments"},
	APP3:  {"APP3", "Reserved for application segments"},
	APP4:  {"APP4", "Reserved for application segments"},
	APP5:  {"APP5", "Reserved for application segments"},
	APP6:  {"APP6", "Reserved for application segments"},
	APP7:  {"APP7", "Reserved for application segments"},
	APP8:  {"APP8", "Reserved for application segments"},
	APP9:  {"APP9", "Reserved for application segments"},
	APP10: {"APP10", "Reserved for application segments"},
	APP11: {"APP11", "Reserved for application segments"},
	APP12: {"APP12", "Reserved for application segments"},
	APP13: {"APP13", "Reserved for application segments"},
	APP14: {"APP14", "Reserved for application segments"},
	APP15: {"APP15", "Reserved for application segments"},
	JPG0:  {"JPG0", "Reserved for JPEG extensions"},
	JPG1:  {"JPG1", "Reserved for JPEG extensions"},
	JPG2:  {"JPG2", "Reserved for JPEG extensions"},
	JPG3:  {"JPG3", "Reserved for JPEG extensions"},
	JPG4:  {"JPG4", "Reserved for JPEG extensions"},
	JPG5:  {"JPG5", "Reserved for JPEG extensions"},
	JPG6:  {"JPG6", "Reserved for JPEG extensions"},
	JPG7:  {"JPG7", "Reserved for JPEG extensions"},
	JPG8:  {"JPG8", "Reserved for JPEG extensions"},
	JPG9:  {"JPG9", "Reserved for JPEG extensions"},
	JPG10: {"JPG10", "Reserved for JPEG extensions"},
	JPG11: {"JPG11", "Reserved for JPEG extensions"},
	JPG12: {"JPG12", "Reserved for JPEG extensions"},
	JPG13: {"JPG13", "Reserved for JPEG extensions"},
	COM:   {"COM", "Comment"},
	TEM:   {"TEM", "For temporary private use in arithmetic coding"},
}

// FileDecoder is a JPEG decoder
type FileDecoder struct {
	decode.Common
}

// Decode JPEG file
func (d *FileDecoder) Decode() {
	var extendedXMP []byte
	soiMarkerFound := false

	inECD := false
	for !d.End() {
		if inECD {
			ecdLen := int64(0)
			for {
				if d.PeekBits(8) == 0xff && d.PeekBits(16) != 0xff00 {
					break
				}
				d.SeekRel(8)
				ecdLen++
			}
			d.SeekRel(-ecdLen * 8)
			d.FieldBitBufLen("entropy_coded_data", int64(ecdLen)*8)
			inECD = false
		} else {
			d.FieldNoneFn("marker", func() {
				d.FieldNoneFn("prefix", func() {
					for d.PeekBits(8) == 0xff {
						d.SeekRel(8)
					}
				})
				markerFound := false
				markerCode := d.FieldUFn("code", func() (uint64, decode.DisplayFormat, string) {
					n := uint(d.U8())
					if m, ok := markers[n]; ok {
						markerFound = true
						return uint64(n), decode.NumberDecimal, m.symbol
					}
					return uint64(n), decode.NumberDecimal, "RES"
				})

				// RST*, SOI, EOI, TEM does not have a length field. All others have a
				// 2 byte length read as "Lf", "Ls" etc or in the default case as "length".

				// TODO: warning on 0x00?
				switch markerCode {
				case SOI:
					soiMarkerFound = true
				case SOF0, SOF1, SOF2, SOF3, SOF5, SOF6, SOF7, SOF9, SOF10, SOF11:
					d.FieldU16("Lf")
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
				case COM:
					comLen := d.FieldU16("Lc")
					d.FieldUTF8("Cm", int64(comLen)-2)
				case SOS:
					d.FieldU16("Ls")
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
					inECD = true
				case RST0, RST1, RST2, RST3, RST4, RST5, RST6, RST7:
					inECD = true
				case TEM:
				case EOI:
				default:
					if markerFound {
						markerLen := d.FieldU16("length")
						d.SubLenFn(int64((markerLen-2)*8), func() {
							app1ExifPrefix := []byte("Exif\x00\x00")
							extendedXMPPrefix := []byte("http://ns.adobe.com/xmp/extension/\x00")

							switch {
							case markerCode == APP1 && d.TryHasBytes(app1ExifPrefix):
								d.FieldUTF8("exif_prefix", 6)
								d.FieldDecodeLen("exif", d.BitsLeft(), tiffImage)
							case markerCode == APP1 && d.TryHasBytes(extendedXMPPrefix):
								d.FieldNoneFn("extended_xmp_chunk", func() {
									d.FieldUTF8("signature", int64(len(extendedXMPPrefix)))
									d.FieldUTF8("guid", 32)
									fullLength := d.FieldU32("full_length")
									offset := d.FieldU32("offset")
									// TODO: FieldBitsLen? concat bitbuf?
									chunk := d.FieldBytesLen("data", d.BitsLeft()/8)

									if extendedXMP == nil {
										extendedXMP = make([]byte, fullLength)
									}
									copy(extendedXMP[offset:], chunk)
								})
							default:
								// TODO: FieldBitsLen?
								d.FieldBitBufLen("data", d.BitsLeft())
							}
						})

					} else {
						d.Invalid(fmt.Sprintf("unknown marker %x", markerCode))
					}
				}
			})
		}
	}

	if !soiMarkerFound {
		d.Invalid("no SOI marker found")
	}

	if extendedXMP != nil {
		bb, err := bitbuf.NewFromBytes(extendedXMP, 0)
		if err != nil {
			panic(err) // TODO: fixme
		}
		// TODO: bit pos, better bitbhuf api?
		d.FieldBitBufFn("extended_xmp", 0, bb.Len, func() (*bitbuf.Buffer, string) {
			return bb, ""
		})
	}
}
