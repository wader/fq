package gif

// https://www.w3.org/Graphics/GIF/spec-gif87.txt
// https://en.wikipedia.org/wiki/GIF
// https://web.archive.org/web/20160304075538/http://qalle.net/gif89a.php#graphiccontrolextension

// TODO: local color map
// TODO: bit depth done correct?

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.GIF,
		&decode.Format{
			Description: "Graphics Interchange Format",
			Groups:      []*decode.Group{format.Probe, format.Image},
			DecodeFn:    gifDecode,
		})
}

const (
	extensionPlainText        = 0x01
	extensionGraphicalControl = 0xf9
	extensionComment          = 0xfe
	extensionApplication      = 0xff
)

var extensionNames = scalar.UintMapSymStr{
	extensionPlainText:        "PlainText",
	extensionGraphicalControl: "GraphicalControl",
	extensionComment:          "Comment",
	extensionApplication:      "Application",
}

func fieldColorMap(d *decode.D, name string, bitDepth int) {
	d.FieldArray(name, func(d *decode.D) {
		for i := 0; i < 1<<bitDepth; i++ {
			d.FieldArray("color", func(d *decode.D) {
				d.FieldU8("r")
				d.FieldU8("g")
				d.FieldU8("b")
			})
		}
	})
}

func gifDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldUTF8("header", 6, d.StrAssert("GIF87a", "GIF89a"))

	d.FieldU16("width")
	d.FieldU16("height")
	gcpFollows := d.FieldBool("gcp_follows")
	d.FieldUintFn("color_resolution", func(d *decode.D) uint64 { return d.U3() + 1 })
	d.FieldU1("zero")
	bitDepth := d.FieldUintFn("bit_depth", func(d *decode.D) uint64 { return d.U3() + 1 })
	d.FieldU8("black_color")
	d.FieldU8("pixel_aspect_ratio")

	if gcpFollows {
		fieldColorMap(d, "global_color_map", int(bitDepth))
	}

	d.FieldArray("blocks", func(d *decode.D) {
	blocks:
		for {
			switch d.PeekUintBits(8) {
			case ';':
				break blocks
			case '!': /* "!" */
				d.FieldStruct("extension_block", func(d *decode.D) {
					d.FieldU8("introducer")
					functionCode := d.FieldU8("function_code", extensionNames, scalar.UintHex)

					dataBytes := &bytes.Buffer{}

					d.FieldArray("func_data_bytes", func(d *decode.D) {
						seenTerminator := false
						for !seenTerminator {

							d.FieldStruct("func_data_byte", func(d *decode.D) {
								byteCount := d.FieldU8("byte_count")
								b := d.FieldRawLen("data", int64(byteCount*8))
								if d.PeekUintBits(8) == 0 {
									d.FieldU8("terminator")
									seenTerminator = true
								}
								d.CopyBits(dataBytes, d.CloneReadSeeker(b))
							})
						}
					})

					_ = functionCode

					// TODO: need a FieldStructBitBuf or something
					// switch functionCode {
					// case extensionGraphicalControl:
					// 	d.FieldFormatBitBuf(
					// 		"graphics_control",
					// 		bitio.NewReader(dataBytes.Bytes(), -1),
					// 	)

					// }
				})
			case ',':
				d.FieldStruct("image", func(d *decode.D) {
					d.FieldU8("separator_character")
					d.FieldU16("left")
					d.FieldU16("top")
					d.FieldU16("width")
					d.FieldU16("height")

					localFollows := d.FieldBool("local_color_map_follows")
					d.FieldBool("image_interlaced")
					d.FieldU3("zero")
					d.FieldUintFn("bit_depth", func(d *decode.D) uint64 { return d.U3() + 1 })
					d.FieldU8("code_size")

					if localFollows {
						fieldColorMap(d, "local_color_map", int(bitDepth))
					}

					d.FieldArray("image_bytes", func(d *decode.D) {
						seenTerminator := false
						for !seenTerminator {

							d.FieldStruct("func_data_byte", func(d *decode.D) {
								byteCount := d.FieldU8("byte_count")
								d.FieldRawLen("data", int64(byteCount*8))
								if d.PeekUintBits(8) == 0 {
									d.FieldU8("terminator")
									seenTerminator = true
								}
							})
						}
					})
				})
			default:
				d.Fatalf("unknown block")
			}
		}
	})

	d.FieldU8("terminator")

	return nil
}
