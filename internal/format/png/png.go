package png

// http://www.libpng.org/pub/png/spec/1.2/PNG-Contents.html
// https://ftp-osl.osuosl.org/pub/libpng/documents/pngext-1.5.0.html

import (
	"fq/internal/decode"
	"fq/internal/format/tiff"
)

var File = &decode.Format{
	Name: "png",
	MIME: "",
	New:  func() decode.Decoder { return &FileDecoder{} },
}

// FileDecoder is a PNG decoder
type FileDecoder struct {
	decode.Common
}

const (
	compressionDeflate = 0
)

var compressionNames = map[uint64]string{
	compressionDeflate: "deflate",
}

// Decode PNG file
func (d *FileDecoder) Decode() {
	d.FieldValidateString("signature", "\x89PNG\r\n\x1a\n")
	for !d.End() {
		d.FieldNoneFn("chunk", func() {
			chunkLength := d.FieldU32("length")

			chunkType := d.FieldStrFn("type", func() (string, string) {
				chunkType := d.UTF8(4)
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
				compressionMethod := d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8)
				_ = compressionMethod

				cb := d.FieldBytesLen("compressed", chunkLength-keywordLen-1)

				switch compressionMethod {
				case compressionDeflate:
					d.FieldDecodeZlib("uncompressed", cb, decode.FormatFn(func(c *decode.Common) {
						c.FieldUTF8("text", c.BitsLeft()/8)
					}))
				}
			case "iCCP":
				profileNameLen := d.PeekFindByte(0, 80)
				d.FieldUTF8("profile_name", profileNameLen-1)
				d.FieldUTF8("null", 1)
				compressionMethod := d.FieldStringMapFn("compression_method", compressionNames, "unknown", d.U8)
				_ = compressionMethod

				cb := d.FieldBytesLen("compressed", chunkLength-profileNameLen-1)

				switch compressionMethod {
				case compressionDeflate:
					d.FieldDecodeZlib("uncompressed", cb, decode.FormatFn(func(c *decode.Common) {
						c.FieldUTF8("text", c.BitsLeft()/8)
					}))
				}
			case "eXIf":
				d.FieldDecodeLen("exit", chunkLength*8, tiff.File)
			default:
				d.FieldBytesLen("data", chunkLength)
			}

			crc := d.FieldU32("crc")

			_ = crc
		})
	}
}
