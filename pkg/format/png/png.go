package png

// http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html
// https://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var iccTag []*decode.Format
var tiffFile []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:     format.PNG,
		Groups:   []string{format.PROBE, format.IMAGE},
		MIMEs:    []string{"image/png"},
		DecodeFn: pngDecode,
		Deps: []decode.Dep{
			{Names: []string{format.ICC}, Formats: &iccTag},
			{Names: []string{format.TIFF}, Formats: &tiffFile},
		},
	})
}

const (
	compressionDeflate = 0
)

var compressionNames = map[uint64]string{
	compressionDeflate: "deflate",
}

func pngDecode(d *decode.D) interface{} {
	d.FieldValidateString("signature", "\x89PNG\r\n\x1a\n")
	d.FieldArrayFn("chunk", func(d *decode.D) {
		for !d.End() {
			d.FieldStructFn("chunk", func(d *decode.D) {
				chunkLength := int64(d.FieldU32("length"))

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
					d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8)
					d.FieldStringMapFn("filter_method", map[uint64]string{
						0: "Adaptive filtering",
					}, "unknown", d.U8)
					d.FieldStringMapFn("interlace_method", map[uint64]string{
						0: "No interlace",
						1: "Adam7 interlace",
					}, "unknown", d.U8)
				case "tEXt":
					// TODO: latin1
					keywordLen := d.PeekFindByte(0, 80)
					d.FieldUTF8("keyword", keywordLen-1)
					d.FieldUTF8("null", 1)
					d.FieldUTF8("text", chunkLength-keywordLen)
				case "zTXt":
					// TODO: latin1
					keywordLen := d.PeekFindByte(0, 80)
					d.FieldUTF8("keyword", keywordLen-1)
					d.FieldUTF8("null", 1)
					compressionMethod, _ := d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8)
					_ = compressionMethod

					switch compressionMethod {
					case compressionDeflate:
						// TODO: make nicer
						d.FieldZlibLen("uncompressed", chunkLength-keywordLen-1, decode.FormatFn(func(c *decode.D) interface{} {
							c.FieldUTF8("text", c.BitsLeft()/8)
							return nil
						}))
					default:
						d.FieldBitBufLen("compressed", (chunkLength-keywordLen-1)*8)
					}
				case "iCCP":
					profileNameLen := d.PeekFindByte(0, 80)
					d.FieldUTF8("profile_name", profileNameLen-1)
					d.FieldUTF8("null", 1)
					compressionMethod, _ := d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8)
					_ = compressionMethod

					switch compressionMethod {
					case compressionDeflate:
						// TODO: make nicer
						d.FieldZlibLen("uncompressed", chunkLength-profileNameLen-1, decode.FormatFn(func(c *decode.D) interface{} {
							c.FieldDecodeLen("icc", c.BitsLeft(), iccTag)
							return nil
						}))
					default:
						d.FieldBitBufLen("compressed", (chunkLength-profileNameLen-1)*8)
					}
				case "eXIf":
					// TODO: decode fail?
					d.FieldDecodeLen("exif", chunkLength*8, tiffFile)
				default:
					d.FieldBitBufLen("data", chunkLength*8)
				}

				crc := d.FieldU32("crc")

				_ = crc
			})
		}
	})

	return nil
}
