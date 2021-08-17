package vpx

// https://tools.ietf.org/html/rfc6386

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

// TODO: vpx frame?

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.VP8_FRAME,
		Description: "VP8 frame",
		DecodeFn:    vp8Decode,
	})
}

func vp8Decode(d *decode.D, in interface{}) interface{} {
	var isKeyFrame bool

	versions := map[uint64]struct {
		reconstruction string
		loop           string
	}{
		0: {"Bicubic", "Normal"},
		1: {"Bilinear", "Simple"},
		2: {"Bilinear", "None"},
		3: {"None", "None"},
	}

	d.FieldStructFn("tag", func(d *decode.D) {
		// first_part_size is not contiguous bits
		firstPartSize0 := d.FieldU3("first_part_size0")
		d.FieldU1("show_frame")
		version := d.FieldU3("version")
		keyFrameV, _ := d.FieldBoolMapFn("frame_type", "non_key_frame", "key_frame", d.Bool)
		firstPartSize1 := d.FieldU16LE("first_part_size1")

		firstPartSize := firstPartSize0 | firstPartSize1<<3
		d.FieldValueU("first_part_size", firstPartSize, "")

		isKeyFrame = !keyFrameV
		if v, ok := versions[version]; ok {
			d.FieldValueStr("reconstruction", v.reconstruction, "")
			d.FieldValueStr("loop", v.loop, "")
		}
	})

	if isKeyFrame {
		d.FieldValidateUFn("start_code", 0x9d012a, d.U24)

		// width and height are not contiguous bits
		width0 := d.FieldU8("width0")
		d.FieldU2("horizontal_scale")
		width1 := d.FieldU6("width1")
		d.FieldValueU("width", width0|width1<<8, "")

		height0 := d.FieldU8("height0")
		d.FieldU2("vertical_scale")
		height1 := d.FieldU6("height1")
		d.FieldValueU("height", height0|height1<<8, "")
	}

	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
