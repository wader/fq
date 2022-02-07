package jpeg

// https://www.w3.org/Graphics/JPEG/itu-t81.pdf
// TODO: warning on junk before marker?
// TODO: extract photohop to own decoder?

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var exifFormat decode.Group
var iccProfileFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.JPEG,
		Description: "Joint Photographic Experts Group file",
		Groups:      []string{format.PROBE, format.IMAGE},
		DecodeFn:    jpegDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.EXIF}, Group: &exifFormat},
			{Names: []string{format.ICC_PROFILE}, Group: &iccProfileFormat},
		},
	})
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

var markers = scalar.UToScalar{
	SOF0:  {Sym: "SOF0", Description: "Baseline DCT"},
	SOF1:  {Sym: "SOF1", Description: "Extended sequential DCT"},
	SOF2:  {Sym: "SOF2", Description: "Progressive DCT"},
	SOF3:  {Sym: "SOF3", Description: "Lossless (sequential)"},
	SOF5:  {Sym: "SOF5", Description: "Differential sequential DCT"},
	SOF6:  {Sym: "SOF6", Description: "Differential progressive DCT"},
	SOF7:  {Sym: "SOF7", Description: "Differential lossless (sequential)"},
	JPG:   {Sym: "JPG", Description: "Reserved for JPEG extensions"},
	SOF9:  {Sym: "SOF9", Description: "Extended sequential DCT"},
	SOF10: {Sym: "SOF10", Description: "Progressive DCT"},
	SOF11: {Sym: "SOF11", Description: "Lossless (sequential)"},
	SOF13: {Sym: "SOF13", Description: "Differential sequential DCT"},
	SOF14: {Sym: "SOF14", Description: "Differential progressive DCT"},
	SOF15: {Sym: "SOF15", Description: "Differential lossless (sequential)"},
	DHT:   {Sym: "DHT", Description: "Define Huffman table(s)"},
	DAC:   {Sym: "DAC", Description: "Define arithmetic coding conditioning(s)"},
	RST0:  {Sym: "RST0", Description: "Restart with modulo 8 count 0"},
	RST1:  {Sym: "RST1", Description: "Restart with modulo 8 count 1"},
	RST2:  {Sym: "RST2", Description: "Restart with modulo 8 count 2"},
	RST3:  {Sym: "RST3", Description: "Restart with modulo 8 count 3"},
	RST4:  {Sym: "RST4", Description: "Restart with modulo 8 count 4"},
	RST5:  {Sym: "RST5", Description: "Restart with modulo 8 count 5"},
	RST6:  {Sym: "RST6", Description: "Restart with modulo 8 count 6"},
	RST7:  {Sym: "RST7", Description: "Restart with modulo 8 count 7"},
	SOI:   {Sym: "SOI", Description: "Start of image"},
	EOI:   {Sym: "EOI", Description: "End of image true"},
	SOS:   {Sym: "SOS", Description: "Start of scan"},
	DQT:   {Sym: "DQT", Description: "Define quantization table(s)"},
	DNL:   {Sym: "DNL", Description: "Define number of lines"},
	DRI:   {Sym: "DRI", Description: "Define restart interval"},
	DHP:   {Sym: "DHP", Description: "Define hierarchical progression"},
	EXP:   {Sym: "EXP", Description: "Expand reference component(s)"},
	APP0:  {Sym: "APP0", Description: "Reserved for application segments"},
	APP1:  {Sym: "APP1", Description: "Reserved for application segments"},
	APP2:  {Sym: "APP2", Description: "Reserved for application segments"},
	APP3:  {Sym: "APP3", Description: "Reserved for application segments"},
	APP4:  {Sym: "APP4", Description: "Reserved for application segments"},
	APP5:  {Sym: "APP5", Description: "Reserved for application segments"},
	APP6:  {Sym: "APP6", Description: "Reserved for application segments"},
	APP7:  {Sym: "APP7", Description: "Reserved for application segments"},
	APP8:  {Sym: "APP8", Description: "Reserved for application segments"},
	APP9:  {Sym: "APP9", Description: "Reserved for application segments"},
	APP10: {Sym: "APP10", Description: "Reserved for application segments"},
	APP11: {Sym: "APP11", Description: "Reserved for application segments"},
	APP12: {Sym: "APP12", Description: "Reserved for application segments"},
	APP13: {Sym: "APP13", Description: "Reserved for application segments"},
	APP14: {Sym: "APP14", Description: "Reserved for application segments"},
	APP15: {Sym: "APP15", Description: "Reserved for application segments"},
	JPG0:  {Sym: "JPG0", Description: "Reserved for JPEG extensions"},
	JPG1:  {Sym: "JPG1", Description: "Reserved for JPEG extensions"},
	JPG2:  {Sym: "JPG2", Description: "Reserved for JPEG extensions"},
	JPG3:  {Sym: "JPG3", Description: "Reserved for JPEG extensions"},
	JPG4:  {Sym: "JPG4", Description: "Reserved for JPEG extensions"},
	JPG5:  {Sym: "JPG5", Description: "Reserved for JPEG extensions"},
	JPG6:  {Sym: "JPG6", Description: "Reserved for JPEG extensions"},
	JPG7:  {Sym: "JPG7", Description: "Reserved for JPEG extensions"},
	JPG8:  {Sym: "JPG8", Description: "Reserved for JPEG extensions"},
	JPG9:  {Sym: "JPG9", Description: "Reserved for JPEG extensions"},
	JPG10: {Sym: "JPG10", Description: "Reserved for JPEG extensions"},
	JPG11: {Sym: "JPG11", Description: "Reserved for JPEG extensions"},
	JPG12: {Sym: "JPG12", Description: "Reserved for JPEG extensions"},
	JPG13: {Sym: "JPG13", Description: "Reserved for JPEG extensions"},
	COM:   {Sym: "COM", Description: "Comment"},
	TEM:   {Sym: "TEM", Description: "For temporary private use in arithmetic coding"},
}

func jpegDecode(d *decode.D, in interface{}) interface{} {
	d.AssertLeastBytesLeft(2)
	if !bytes.Equal(d.PeekBytes(2), []byte{0xff, SOI}) {
		d.Errorf("no SOI marker")
	}

	var extendedXMP []byte
	soiMarkerFound := false
	eoiMarkerFound := false

	d.FieldArray("segments", func(d *decode.D) {
		inECD := false
		for d.NotEnd() && !eoiMarkerFound {
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
				d.FieldRawLen("entropy_coded_data", ecdLen*8)
				inECD = false
			} else {
				d.FieldStruct("marker", func(d *decode.D) {
					prefixLen := d.PeekFindByte(0xff, -1) + 1
					d.FieldRawLen("prefix", prefixLen*8, d.AssertBitBuf([]byte{0xff}))
					markerCode := d.FieldU8("code", markers)
					_, markerFound := markers[markerCode]

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
						d.FieldArray("frame_components", func(d *decode.D) {
							for i := uint64(0); i < nf; i++ {
								d.FieldStruct("frame_component", func(d *decode.D) {
									d.FieldU8("C")
									d.FieldU4("H")
									d.FieldU4("V")
									d.FieldU8("Tq")
								})
							}
						})
					case COM:
						comLen := d.FieldU16("Lc")
						d.FieldUTF8("Cm", int(comLen)-2)
					case SOS:
						d.FieldU16("Ls")
						ns := d.FieldU8("Ns")
						d.FieldArray("scan_components", func(d *decode.D) {
							for i := uint64(0); i < ns; i++ {
								d.FieldStruct("scan_component", func(d *decode.D) {
									d.FieldU8("Cs")
									d.FieldU4("Td")
									d.FieldU4("Ta")
								})
							}
						})
						d.FieldU8("Ss")
						d.FieldU8("Se")
						d.FieldU4("Ah")
						d.FieldU4("Al")
						inECD = true
					case DQT:
						lQ := int64(d.FieldU16("Lq"))
						// TODO: how to extract n? spec says lq is 2 + sum for i in 1 to n 65+64*Pq(i)
						d.FramedFn(lQ*8-16, func(d *decode.D) {
							d.FieldArray("Qs", func(d *decode.D) {
								for d.NotEnd() {
									d.FieldStruct("Q", func(d *decode.D) {
										pQ := d.FieldU4("Pq")
										qBits := 8
										if pQ != 0 {
											qBits = 16
										}
										d.FieldU4("Tq")
										qK := uint64(0)
										d.FieldArrayLoop("Q", func() bool { return qK < 64 }, func(d *decode.D) {
											d.FieldU("Q", qBits)
											qK++
										})
									})
								}
							})
						})
					case RST0, RST1, RST2, RST3, RST4, RST5, RST6, RST7:
						inECD = true
					case TEM:
					case EOI:
						eoiMarkerFound = true
					default:
						if !markerFound {
							d.Errorf("unknown marker %x", markerCode)
						}

						markerLen := d.FieldU16("length")
						d.FramedFn(int64((markerLen-2)*8), func(d *decode.D) {
							// TODO: map lookup and descriptions?
							app0JFIFPrefix := []byte("JFIF\x00")
							app1ExifPrefix := []byte("Exif\x00\x00")
							extendedXMPPrefix := []byte("http://ns.adobe.com/xmp/extension/\x00")
							app2ICCProfile := []byte("ICC_PROFILE\x00")
							// TODO: other version? generic?
							app13PhotoshopPrefix := []byte("Photoshop 3.0\x00")

							switch {
							case markerCode == APP0 && d.TryHasBytes(app0JFIFPrefix):
								d.FieldUTF8("identifier", len(app0JFIFPrefix))
								d.FieldStruct("version", func(d *decode.D) {
									d.FieldU8("major")
									d.FieldU8("minor")
								})
								d.FieldU8("density_units")
								d.FieldU16("xdensity")
								d.FieldU16("ydensity")
								xThumbnail := d.FieldU8("xthumbnail")
								yThumbnail := d.FieldU8("ythumbnail")
								d.FieldRawLen("data", int64(xThumbnail*yThumbnail)*3*8)
							case markerCode == APP1 && d.TryHasBytes(app1ExifPrefix):
								d.FieldUTF8("exif_prefix", len(app1ExifPrefix))
								d.FieldFormatLen("exif", d.BitsLeft(), exifFormat, nil)
							case markerCode == APP1 && d.TryHasBytes(extendedXMPPrefix):
								d.FieldStruct("extended_xmp_chunk", func(d *decode.D) {
									d.FieldUTF8("signature", len(extendedXMPPrefix))
									d.FieldUTF8("guid", 32)
									fullLength := d.FieldU32("full_length")
									offset := d.FieldU32("offset")
									// TODO: FieldBitsLen? concat bitbuf?
									chunk := d.FieldRawLen("data", d.BitsLeft())
									// TODO: redo this? multi reader?
									chunkBytes := d.MustReadAllBits(chunk)

									if extendedXMP == nil {
										extendedXMP = make([]byte, fullLength)
									}
									copy(extendedXMP[offset:], chunkBytes)
								})
							case markerCode == APP2 && d.TryHasBytes(app2ICCProfile):
								d.FieldUTF8("icc_profile_prefix", len(app2ICCProfile))
								// TODO: support multimarker?
								d.FieldU8("cur_marker")
								d.FieldU8("num_markers")
								d.FieldFormatLen("icc_profile", d.BitsLeft(), iccProfileFormat, nil)
							case markerCode == APP13 && d.TryHasBytes(app13PhotoshopPrefix):
								d.FieldUTF8("identifier", len(app13PhotoshopPrefix))
								signature := d.FieldUTF8("signature", 4)
								switch signature {
								case "8BIM":
									// TODO: description?
									d.FieldU16("block", psImageResourceBlockNames)
									d.FieldRawLen("data", d.BitsLeft())
								default:
								}
							default:
								// TODO: FieldBitsLen?
								d.FieldRawLen("data", d.BitsLeft())
							}
						})
					}
				})
			}
		}
	})

	if !soiMarkerFound {
		d.Errorf("no SOI marker found")
	}

	if extendedXMP != nil {
		d.FieldRootBitBuf("extended_xmp", bitio.NewBitReader(extendedXMP, -1))
	}

	return nil
}
