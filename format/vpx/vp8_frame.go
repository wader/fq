package vpx

// https://tools.ietf.org/html/rfc6386

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: vpx frame?

func init() {
	interp.RegisterFormat(
		format.VP8_Frame,
		&decode.Format{
			Description: "VP8 frame",
			DecodeFn:    vp8Decode,
		})
}

func vp8Decode(d *decode.D) any {
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

	d.FieldStruct("tag", func(d *decode.D) {
		// first_part_size is not contiguous bits
		firstPartSize0 := d.FieldU3("first_part_size0")
		d.FieldU1("show_frame")
		version := d.FieldU3("version")
		keyFrameV := d.FieldBool("frame_type", scalar.BoolMapSymStr{true: "non_key_frame", false: "key_frame"})
		firstPartSize1 := d.FieldU16LE("first_part_size1")

		firstPartSize := firstPartSize0 | firstPartSize1<<3
		d.FieldValueUint("first_part_size", firstPartSize)

		isKeyFrame = !keyFrameV
		if v, ok := versions[version]; ok {
			d.FieldValueStr("reconstruction", v.reconstruction)
			d.FieldValueStr("loop", v.loop)
		}
	})

	if isKeyFrame {
		d.FieldU24("start_code", d.UintValidate(0x9d012a), scalar.UintHex)

		// width and height are not contiguous bits
		width0 := d.FieldU8("width0")
		d.FieldU2("horizontal_scale")
		width1 := d.FieldU6("width1")
		d.FieldValueUint("width", width0|width1<<8)

		height0 := d.FieldU8("height0")
		d.FieldU2("vertical_scale")
		height1 := d.FieldU6("height1")
		d.FieldValueUint("height", height0|height1<<8)
	}

	d.FieldRawLen("data", d.BitsLeft())

	return nil
}
