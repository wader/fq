package png

// http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html
// https://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html
// https://wiki.mozilla.org/APNG_Specification

import (
	"compress/zlib"
	"hash/crc32"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var iccProfileGroup decode.Group
var exifGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.PNG,
		&decode.Format{
			Description: "Portable Network Graphics file",
			Groups:      []*decode.Group{format.Probe, format.Image},
			DecodeFn:    pngDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.ICC_Profile}, Out: &iccProfileGroup},
				{Groups: []*decode.Group{format.Exif}, Out: &exifGroup},
			},
		})
}

const (
	compressionDeflate = 0
)

var compressionNames = scalar.UintMapSymStr{
	compressionDeflate: "deflate",
}

const (
	disposeOpNone       = 0
	disposeOpBackground = 1
	disposeOpPrevious   = 2
)

var disposeOpNames = scalar.UintMapSymStr{
	disposeOpNone:       "none",
	disposeOpBackground: "background",
	disposeOpPrevious:   "previous",
}

const (
	blendOpNone       = 0
	blendOpBackground = 1
)

var blendOpNames = scalar.UintMapSymStr{
	blendOpNone:       "source",
	blendOpBackground: "over",
}

const (
	colorTypeGrayscale          = 0
	colorTypeRGB                = 2
	colorTypePalette            = 3
	colorTypeGrayscaleWithAlpha = 4
	colorTypeRGBA               = 6
)

var colorTypeMap = scalar.UintMapSymStr{
	colorTypeGrayscale:          "grayscale",
	colorTypeRGB:                "rgb",
	colorTypePalette:            "palette",
	colorTypeGrayscaleWithAlpha: "grayscale_alpha",
	colorTypeRGBA:               "rgba",
}

func pngDecode(d *decode.D) any {
	iEndFound := false
	var colorType uint64

	d.FieldRawLen("signature", 8*8, d.AssertBitBuf([]byte("\x89PNG\r\n\x1a\n")))
	d.FieldStructArrayLoop("chunks", "chunk", func() bool { return d.NotEnd() && !iEndFound }, func(d *decode.D) {
		chunkLength := d.FieldU32("length")
		crcStartPos := d.Pos()
		chunkType := d.FieldUTF8("type", 4)
		// upper/lower case in chunk type is used for flags
		d.SeekRel(-4 * 8)
		d.SeekRel(2)
		d.FieldBool("ancillary")
		d.SeekRel(7)
		d.FieldBool("private")
		d.SeekRel(7)
		d.FieldBool("reserved")
		d.SeekRel(7)
		d.FieldBool("safe_to_copy")
		d.SeekRel(5)

		d.FramedFn(int64(chunkLength)*8, func(d *decode.D) {
			switch chunkType {
			case "IHDR":
				d.FieldU32("width")
				d.FieldU32("height")
				d.FieldU8("bit_depth")
				colorType = d.FieldU8("color_type", colorTypeMap)
				d.FieldU8("compression_method", compressionNames)
				d.FieldU8("filter_method", scalar.UintMapSymStr{
					0: "adaptive_filtering",
				})
				d.FieldU8("interlace_method", scalar.UintMapSymStr{
					0: "none",
					1: "adam7",
				})
			case "tEXt":
				d.FieldUTF8Null("keyword")
				d.FieldUTF8("text", int(d.BitsLeft())/8)
			case "zTXt":
				d.FieldUTF8Null("keyword")
				compressionMethod := d.FieldU8("compression_method", compressionNames)
				dataLen := d.BitsLeft()

				// TODO: make nicer
				d.FieldRawLen("compressed", dataLen)
				d.SeekRel(-dataLen)

				switch compressionMethod {
				case compressionDeflate:
					d.FieldFormatReaderLen("uncompressed", dataLen, zlib.NewReader, decode.FormatFn(func(d *decode.D) any {
						d.FieldUTF8("text", int(d.BitsLeft()/8))
						return nil
					}))
				default:
					d.FieldRawLen("data", dataLen)
				}
			case "iCCP":
				d.FieldUTF8Null("profile_name")
				compressionMethod := d.FieldU8("compression_method", compressionNames)
				dataLen := d.BitsLeft()

				d.FieldRawLen("compressed", dataLen)
				d.SeekRel(-dataLen)

				switch compressionMethod {
				case compressionDeflate:
					d.FieldFormatReaderLen("uncompressed", dataLen, zlib.NewReader, &iccProfileGroup)
				default:
					d.FieldRawLen("data", dataLen)
				}
			case "pHYs":
				d.FieldU32("x_pixels_per_unit")
				d.FieldU32("y_pixels_per_unit")
				d.FieldU8("unit")
			case "bKGD":
				switch colorType {
				case colorTypePalette:
					d.FieldU8("index")
				case colorTypeGrayscale, colorTypeGrayscaleWithAlpha:
					d.FieldU16("gray")
				case colorTypeRGB, colorTypeRGBA:
					d.FieldU16("r")
					d.FieldU16("g")
					d.FieldU16("b")
				}
			case "gAMA":
				d.FieldU32("value")
			case "cHRM":
				df := func(d *decode.D) float64 { return float64(d.U32()) / 1000.0 }
				d.FieldFltFn("white_point_x", df)
				d.FieldFltFn("white_point_y", df)
				d.FieldFltFn("red_x", df)
				d.FieldFltFn("red_y", df)
				d.FieldFltFn("green_x", df)
				d.FieldFltFn("green_y", df)
				d.FieldFltFn("blue_x", df)
				d.FieldFltFn("blue_y", df)
			case "eXIf":
				d.FieldFormatLen("exif", d.BitsLeft(), &exifGroup, nil)
			case "acTL":
				d.FieldU32("num_frames")
				d.FieldU32("num_plays")
			case "fcTL":
				d.FieldU32("sequence_number")
				d.FieldU32("width")
				d.FieldU32("height")
				d.FieldU32("x_offset")
				d.FieldU32("y_offset")
				d.FieldU16("delay_num")
				d.FieldU16("delay_sep")
				d.FieldU8("dispose_op", disposeOpNames)
				d.FieldU8("blend_op", blendOpNames)
			case "fdAT":
				d.FieldU32("sequence_number")
				d.FieldRawLen("data", d.BitsLeft()-32)
			case "PLTE":
				d.FieldArray("palette", func(d *decode.D) {
					for !d.End() {
						d.FieldStruct("color", func(d *decode.D) {
							d.FieldU8("r")
							d.FieldU8("g")
							d.FieldU8("b")
						})
					}
				})
			case "tRNS":
				switch colorType {
				case colorTypeGrayscale:
					d.FieldU16("alpha")
				case colorTypeRGB:
					d.FieldU16("r")
					d.FieldU16("g")
					d.FieldU16("b")
				case colorTypePalette:
					d.FieldArray("alphas", func(d *decode.D) {
						for !d.End() {
							d.FieldU8("alpha")
						}
					})
				}
			default:
				if chunkType == "IEND" {
					iEndFound = true
				} else {
					d.FieldRawLen("data", d.BitsLeft())
				}
			}
		})

		chunkCRC := crc32.NewIEEE()
		d.Copy(chunkCRC, bitio.NewIOReader(d.BitBufRange(crcStartPos, d.Pos()-crcStartPos)))
		d.FieldU32("crc", d.UintValidateBytes(chunkCRC.Sum(nil)), scalar.UintHex)
	})

	return nil
}
