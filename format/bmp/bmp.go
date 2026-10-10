package bmp

// https://en.wikipedia.org/wiki/BMP_file_format
// https://learn.microsoft.com/en-us/windows/win32/gdi/bitmap-header-types

// Decodes the file header and the BITMAPINFOHEADER family of DIB headers
// (size >= 40), which covers essentially all BMPs seen in the wild. The
// 12 byte BITMAPCOREHEADER is not handled specially.

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.BMP,
		&decode.Format{
			Description: "Bitmap image",
			Groups:      []*decode.Group{format.Probe, format.Image},
			DecodeFn:    bmpDecode,
		})
}

const (
	compressionRGB            = 0
	compressionRLE8           = 1
	compressionRLE4           = 2
	compressionBitfields      = 3
	compressionJPEG           = 4
	compressionPNG            = 5
	compressionAlphaBitfields = 6
	compressionCMYK           = 11
	compressionCMYKRLE8       = 12
	compressionCMYKRLE4       = 13
)

var compressionNames = scalar.UintMapSymStr{
	compressionRGB:            "rgb",
	compressionRLE8:           "rle8",
	compressionRLE4:           "rle4",
	compressionBitfields:      "bitfields",
	compressionJPEG:           "jpeg",
	compressionPNG:            "png",
	compressionAlphaBitfields: "alpha_bitfields",
	compressionCMYK:           "cmyk",
	compressionCMYKRLE8:       "cmyk_rle8",
	compressionCMYKRLE4:       "cmyk_rle4",
}

func bmpDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	var bitmapOffset uint64
	var bitsPerPixel uint64
	var colorsUsed uint64

	d.FieldStruct("file_header", func(d *decode.D) {
		d.FieldUTF8("type", 2, d.StrAssert("BM"))
		d.FieldU32("size")
		d.FieldU16("reserved1")
		d.FieldU16("reserved2")
		bitmapOffset = d.FieldU32("bitmap_offset")
	})

	d.FieldStruct("dib_header", func(d *decode.D) {
		headerSize := d.FieldU32("size")
		d.FieldS32("width")
		d.FieldS32("height")
		d.FieldU16("planes")
		bitsPerPixel = d.FieldU16("bits_per_pixel")
		d.FieldU32("compression", compressionNames)
		d.FieldU32("image_size")
		d.FieldS32("x_pixels_per_meter")
		d.FieldS32("y_pixels_per_meter")
		colorsUsed = d.FieldU32("colors_used")
		d.FieldU32("colors_important")
		// v4/v5 headers add bit masks, color space and gamma after the common part
		if headerSize > 40 {
			d.FieldRawLen("rest", int64(headerSize-40)*8)
		}
	})

	if bitsPerPixel <= 8 {
		n := colorsUsed
		if n == 0 {
			n = 1 << bitsPerPixel
		}
		d.FieldArray("color_table", func(d *decode.D) {
			for range int(n) {
				d.FieldStruct("color", func(d *decode.D) {
					d.FieldU8("blue")
					d.FieldU8("green")
					d.FieldU8("red")
					d.FieldU8("reserved")
				})
			}
		})
	}

	if gap := int64(bitmapOffset)*8 - d.Pos(); gap > 0 {
		d.FieldRawLen("gap", gap)
	}

	d.FieldRawLen("pixels", d.BitsLeft())

	return nil
}
