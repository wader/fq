package png

// http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html
// https://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html
// https://wiki.mozilla.org/APNG_Specification

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"fq/pkg/ranges"
	"hash/crc32"
)

var iccProfileFormat []*decode.Format
var exifFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.PNG,
		Description: "Portable network graphics file",
		Groups:      []string{format.PROBE, format.IMAGE},
		MIMEs:       []string{"image/png"},
		DecodeFn:    pngDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ICC_PROFILE}, Formats: &iccProfileFormat},
			{Names: []string{format.EXIF}, Formats: &exifFormat},
		},
	})
}

const (
	compressionDeflate = 0
)

var compressionNames = map[uint64]string{
	compressionDeflate: "deflate",
}

const (
	disposeOpNone       = 0
	disposeOpBackground = 1
	disposeOpPrevious   = 2
)

var disposeOpNames = map[uint64]string{
	disposeOpNone:       "None",
	disposeOpBackground: "Background",
	disposeOpPrevious:   "Previous",
}

const (
	blendOpNone       = 0
	blendOpBackground = 1
)

var blendOpNames = map[uint64]string{
	blendOpNone:       "Source",
	blendOpBackground: "Over",
}

func pngDecode(d *decode.D, in interface{}) interface{} {
	iEndFound := false

	d.FieldValidateUTF8("signature", "\x89PNG\r\n\x1a\n")
	d.FieldStructArrayLoopFn("chunks", "chunk", func() bool { return d.NotEnd() && !iEndFound }, func(d *decode.D) {
		chunkLength := int(d.FieldU32("length"))
		crcStartPos := d.Pos()
		chunkType := d.FieldStrFn("type", func() (string, string) {
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
			return chunkType, ""
		})

		switch chunkType {
		case "IHDR":
			d.FieldU32("width")
			d.FieldU32("height")
			d.FieldU8("bit_depth")
			d.FieldU8("color_type")
			d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8, decode.NumberDecimal)
			d.FieldStringMapFn("filter_method", map[uint64]string{
				0: "Adaptive filtering",
			}, "unknown", d.U8, decode.NumberDecimal)
			d.FieldStringMapFn("interlace_method", map[uint64]string{
				0: "No interlace",
				1: "Adam7 interlace",
			}, "unknown", d.U8, decode.NumberDecimal)
		case "tEXt":
			// TODO: latin1
			keywordLen := int(d.PeekFindByte(0, 80))
			d.FieldUTF8("keyword", keywordLen-1)
			d.FieldUTF8("null", 1)
			d.FieldUTF8("text", chunkLength-keywordLen)
		case "zTXt":
			// TODO: latin1
			keywordLen := int(d.PeekFindByte(0, 80))
			d.FieldUTF8("keyword", keywordLen-1)
			d.FieldUTF8("null", 1)
			compressionMethod, _ := d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8, decode.NumberDecimal)
			dataLen := (chunkLength - keywordLen - 1) * 8

			switch compressionMethod {
			case compressionDeflate:
				dd := d.FieldStructFn("data", func(d *decode.D) {
					d.FieldDecodeZlibLen("uncompressed", int64(dataLen), decode.FormatFn(func(d *decode.D, in interface{}) interface{} {
						d.FieldUTF8("text", int(d.BitsLeft()/8))
						return nil
					}))
				})
				// TODO: depends on isRoot in postProcess
				dd.Value.Range = ranges.Range{Start: d.Pos() - int64(dataLen), Len: int64(dataLen)}
			default:
				d.FieldBitBufLen("data", int64(dataLen))
			}
		case "iCCP":
			profileNameLen := int(d.PeekFindByte(0, 80))
			d.FieldUTF8("profile_name", profileNameLen-1)
			d.FieldUTF8("null", 1)
			compressionMethod, _ := d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8, decode.NumberDecimal)
			dataLen := (chunkLength - profileNameLen - 1) * 8

			switch compressionMethod {
			case compressionDeflate:
				dd := d.FieldStructFn("data", func(d *decode.D) {
					d.FieldDecodeZlibLen("uncompressed", int64(dataLen), iccProfileFormat)
				})
				dd.Value.Range = ranges.Range{Start: d.Pos() - int64(dataLen), Len: int64(dataLen)}
			default:
				d.FieldBitBufLen("data", int64(dataLen))
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
			df := func() (float64, string) { return float64(d.U32()) / 1000.0, "" }
			d.FieldFloatFn("white_point_x", df)
			d.FieldFloatFn("white_point_y", df)
			d.FieldFloatFn("red_x", df)
			d.FieldFloatFn("red_y", df)
			d.FieldFloatFn("green_x", df)
			d.FieldFloatFn("green_y", df)
			d.FieldFloatFn("blue_x", df)
			d.FieldFloatFn("blue_y", df)
		case "eXIf":
			d.FieldDecodeLen("exif", int64(chunkLength)*8, exifFormat)
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
			d.FieldStringMapFn("dispose_op", disposeOpNames, "Unknown", d.U8, decode.NumberDecimal)
			d.FieldStringMapFn("blend_op", blendOpNames, "Unknown", d.U8, decode.NumberDecimal)
		case "fdAT":
			d.FieldU32("sequence_number")
			d.FieldBitBufLen("data", int64(chunkLength-4)*8)
		default:
			if chunkType == "IEND" {
				iEndFound = true
			} else {
				d.FieldBitBufLen("data", int64(chunkLength)*8)
			}
		}

		chunkCRC := crc32.NewIEEE()
		decode.MustCopy(chunkCRC, d.BitBufRange(crcStartPos, d.Pos()-crcStartPos))
		d.FieldChecksumLen("crc", 32, chunkCRC.Sum(nil), decode.BigEndian)
	})

	return nil
}
