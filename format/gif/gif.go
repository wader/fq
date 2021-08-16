package gif

// https://www.w3.org/Graphics/GIF/spec-gif87.txt
// https://en.wikipedia.org/wiki/GIF
// https://web.archive.org/web/20160304075538/http://qalle.net/gif89a.php#graphiccontrolextension

// TODO: local color map
// TODO: bit depth done correct?
// TDOO: mime mage/gif

import (
	"bytes"
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.GIF,
		Description: "Graphics Interchange Format",
		Groups:      []string{format.PROBE, format.IMAGE},
		DecodeFn:    gifDecode,
	})
}

const (
	extensionPlainText        = 0x01
	extensionGraphicalControl = 0xf9
	extensionComment          = 0xfe
	extensionApplication      = 0xff
)

var extensionNames = map[uint64]string{
	extensionPlainText:        "PlainText",
	extensionGraphicalControl: "GraphicalControl",
	extensionComment:          "Comment",
	extensionApplication:      "Application",
}

func fieldColorMap(d *decode.D, name string, bitDepth int) {
	d.FieldArrayFn(name, func(d *decode.D) {
		for i := 0; i < 1<<bitDepth; i++ {
			d.FieldArrayFn("color", func(d *decode.D) {
				d.FieldU8("r")
				d.FieldU8("g")
				d.FieldU8("b")
			})
		}
	})
}

func gifDecode(d *decode.D, in interface{}) interface{} {
	d.Endian = decode.LittleEndian

	d.FieldValidateUTF8("header", "GIF89a")
	d.FieldU16("width")
	d.FieldU16("height")
	gcpFollows := d.FieldBool("gcp_follows")
	d.FieldUFn("color_resolution", func() (uint64, decode.DisplayFormat, string) {
		return d.U3() + 1, decode.NumberDecimal, ""
	})
	d.FieldU1("zero")
	bitDepth := d.FieldUFn("bit_depth", func() (uint64, decode.DisplayFormat, string) {
		return d.U3() + 1, decode.NumberDecimal, ""
	})
	d.FieldU8("black_color")
	d.FieldU8("pixel_aspect_ratio")

	if gcpFollows {
		fieldColorMap(d, "global_color_map", int(bitDepth))
	}

	d.FieldArrayFn("blocks", func(d *decode.D) {
	blocks:
		for {
			switch d.PeekBits(8) {
			case 0x3b: /* ";"  */
				break blocks
			case 0x21: /* "!" */
				d.FieldStructFn("extension_block", func(d *decode.D) {
					d.FieldU8("introducer")
					functionCode, _ := d.FieldStringMapFn("function_code", extensionNames, "Unknown", d.U8, decode.NumberHex)

					dataBytes := &bytes.Buffer{}

					d.FieldArrayFn("func_data_bytes", func(d *decode.D) {
						seenTerminator := false
						for !seenTerminator {

							d.FieldStructFn("func_data_byte", func(d *decode.D) {
								byteCount := d.FieldU8("byte_count")
								b := d.FieldBitBufLen("data", int64(byteCount*8))
								if d.PeekBits(8) == 0 {
									d.FieldU8("terminator")
									seenTerminator = true
								}
								decode.MustCopy(dataBytes, b.Copy())
							})
						}
					})

					_ = functionCode

					// TODO: need a FieldStructBitBuf or something
					// switch functionCode {
					// case extensionGraphicalControl:
					// 	d.FieldFormatBitBuf(
					// 		"graphics_control",
					// 		bitio.NewBufferFromBytes(dataBytes.Bytes(), -1),
					// 	)

					// }
				})
			case 0x2c: /* "," */
				d.FieldStructFn("image", func(d *decode.D) {
					d.FieldU8("separator_character")
					d.FieldU16("left")
					d.FieldU16("top")
					d.FieldU16("width")
					d.FieldU16("height")

					localFollows := d.FieldBool("local_color_map_follows")
					d.FieldBool("image_interlaced")
					d.FieldU3("zero")
					d.FieldUFn("bit_depth", func() (uint64, decode.DisplayFormat, string) {
						return d.U3() + 1, decode.NumberDecimal, ""
					})
					d.FieldU8("code_size")

					if localFollows {
						fieldColorMap(d, "local_color_map", int(bitDepth))
					}

					d.FieldArrayFn("image_bytes", func(d *decode.D) {
						seenTerminator := false
						for !seenTerminator {

							d.FieldStructFn("func_data_byte", func(d *decode.D) {
								byteCount := d.FieldU8("byte_count")
								d.FieldBitBufLen("data", int64(byteCount*8))
								if d.PeekBits(8) == 0 {
									d.FieldU8("terminator")
									seenTerminator = true
								}
							})
						}
					})
				})
			}
		}
	})

	d.FieldU8("terminator")

	return nil
}
