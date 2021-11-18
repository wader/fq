package png

// http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html
// https://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html
// https://wiki.mozilla.org/APNG_Specification
// TODO: color types

import (
	"compress/zlib"
	"hash/crc32"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var iccProfileFormat decode.Group
var exifFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.PNG,
		Description: "Portable Network Graphics file",
		Groups:      []string{format.PROBE, format.IMAGE},
		DecodeFn:    pngDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ICC_PROFILE}, Group: &iccProfileFormat},
			{Names: []string{format.EXIF}, Group: &exifFormat},
		},
	})
}

const (
	compressionDeflate = 0
)

var compressionNames = decode.UToStr{
	compressionDeflate: "deflate",
}

const (
	disposeOpNone       = 0
	disposeOpBackground = 1
	disposeOpPrevious   = 2
)

var disposeOpNames = decode.UToStr{
	disposeOpNone:       "None",
	disposeOpBackground: "Background",
	disposeOpPrevious:   "Previous",
}

const (
	blendOpNone       = 0
	blendOpBackground = 1
)

var blendOpNames = decode.UToStr{
	blendOpNone:       "Source",
	blendOpBackground: "Over",
}

func pngDecode(d *decode.D, in interface{}) interface{} {
	iEndFound := false

	d.FieldRawLen("signature", 8*8, d.AssertBitBuf([]byte("\x89PNG\r\n\x1a\n")))
	d.FieldStructArrayLoop("chunks", "chunk", func() bool { return d.NotEnd() && !iEndFound }, func(d *decode.D) {
		chunkLength := int(d.FieldU32("length"))
		crcStartPos := d.Pos()
		// TODO: this is a bit weird, use struct?
		chunkType := d.FieldStrFn("type", func(d *decode.D) string {
			chunkType := d.UTF8(4)
			// upper/lower case in chunk type is used to set flags
			d.SeekRel(-4 * 8)
			d.SeekRel(3)
			d.FieldBool("ancillary")
			d.SeekRel(7)
			d.FieldBool("private")
			d.SeekRel(7)
			d.FieldBool("reserved")
			d.SeekRel(7)
			d.FieldBool("safe_to_copy")
			d.SeekRel(4)
			return chunkType
		})

		d.LenFn(int64(chunkLength)*8, func(d *decode.D) {
			switch chunkType {
			case "IHDR":
				d.FieldU32("width")
				d.FieldU32("height")
				d.FieldU8("bit_depth")
				d.FieldU8("color_type")
				d.FieldU8("compression_method", d.MapUToStrSym(compressionNames))
				d.FieldU8("filter_method", d.MapUToStrSym(decode.UToStr{
					0: "Adaptive filtering",
				}))
				d.FieldU8("interlace_method", d.MapUToStrSym(decode.UToStr{
					0: "No interlace",
					1: "Adam7 interlace",
				}))
			case "tEXt":
				d.FieldUTF8Null("keyword")
				d.FieldUTF8("text", int(d.BitsLeft())/8)
			case "zTXt":
				d.FieldUTF8Null("keyword")
				compressionMethod := d.FieldU8("compression_method", d.MapUToStrSym(compressionNames))
				dataLen := d.BitsLeft()

				// TODO: make nicer
				d.FieldRawLen("compressed", dataLen)
				d.SeekRel(-dataLen)

				switch compressionMethod {
				case compressionDeflate:
					d.FieldFormatReaderLen("uncompressed", dataLen, zlib.NewReader, decode.FormatFn(func(d *decode.D, in interface{}) interface{} {
						d.FieldUTF8("text", int(d.BitsLeft()/8))
						return nil
					}))
				default:
					d.FieldRawLen("data", dataLen)
				}
			case "iCCP":
				d.FieldUTF8Null("profile_name")
				compressionMethod := d.FieldU8("compression_method", d.MapUToStrSym(compressionNames))
				dataLen := d.BitsLeft()

				d.FieldRawLen("compressed", dataLen)
				d.SeekRel(-dataLen)

				switch compressionMethod {
				case compressionDeflate:
					d.FieldFormatReaderLen("uncompressed", dataLen, zlib.NewReader, iccProfileFormat)
				default:
					d.FieldRawLen("data", dataLen)
				}
			case "pHYs":
				d.FieldU32("x_pixels_per_unit")
				d.FieldU32("y_pixels_per_unit")
				d.FieldU8("unit")
			case "bKGD":
				d.FieldU16("value")
			case "gAMA":
				d.FieldU32("value")
			case "cHRM":
				df := func(d *decode.D) float64 { return float64(d.U32()) / 1000.0 }
				d.FieldFFn("white_point_x", df)
				d.FieldFFn("white_point_y", df)
				d.FieldFFn("red_x", df)
				d.FieldFFn("red_y", df)
				d.FieldFFn("green_x", df)
				d.FieldFFn("green_y", df)
				d.FieldFFn("blue_x", df)
				d.FieldFFn("blue_y", df)
			case "eXIf":
				d.FieldFormatLen("exif", int64(chunkLength)*8, exifFormat, nil)
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
				d.FieldU8("dispose_op", d.MapUToStrSym(disposeOpNames))
				d.FieldU8("blend_op", d.MapUToStrSym(blendOpNames))
			case "fdAT":
				d.FieldU32("sequence_number")
				d.FieldRawLen("data", int64(chunkLength-4)*8)
			default:
				if chunkType == "IEND" {
					iEndFound = true
				} else {
					d.FieldRawLen("data", int64(chunkLength)*8)
				}
			}
		})

		chunkCRC := crc32.NewIEEE()
		d.MustCopy(chunkCRC, d.BitBufRange(crcStartPos, d.Pos()-crcStartPos))
		d.FieldRawLen("crc", 32, d.ValidateBitBuf(chunkCRC.Sum(nil)), d.RawHex)
	})

	return nil
}
