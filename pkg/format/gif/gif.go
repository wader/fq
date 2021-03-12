package gif

// https://www.w3.org/Graphics/GIF/spec-gif87.txt
// https://en.wikipedia.org/wiki/GIF

// TODO: local color map
// TODO: bit depth done correct?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var tiffImage []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.GIF,
		Description: "Graphics Interchange Format",
		Groups:      []string{format.PROBE, format.IMAGE},
		MIMEs:       []string{"image/gif"},
		DecodeFn:    gifDecode,
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
		d.FieldArrayFn("global_color_map", func(d *decode.D) {
			for i := 0; i < 1<<bitDepth; i++ {
				d.FieldArrayFn("global_color_map", func(d *decode.D) {
					d.FieldU8("r")
					d.FieldU8("g")
					d.FieldU8("b")

				})
			}
		})
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
					d.FieldU8("function code")
					d.FieldArrayFn("func_data_bytes", func(d *decode.D) {
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
			case 0x2c: /* "," */
				d.FieldStructFn("image", func(d *decode.D) {
					d.FieldU8("image_separator_character")
					d.FieldU16("left")
					d.FieldU16("top")
					d.FieldU16("width")
					d.FieldU16("height")

					d.FieldBool("use_global_color_map")
					d.FieldBool("local_color_follows")
					d.FieldBool("image_sequential")
					d.FieldBool("image_interlaced")
					d.FieldU1("zero")
					d.FieldUFn("bit_depth", func() (uint64, decode.DisplayFormat, string) {
						return d.U3() + 1, decode.NumberDecimal, ""
					})
					d.FieldU8("code_size")

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
