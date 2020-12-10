package vpx

// https://storage.googleapis.com/downloads.webmproject.org/docs/vp9/vp9-bitstream-specification-v0.6-20160331-draft.pdf

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

// TODO: vpx frame?

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.VP9_FRAME,
		Description: "VP9 frame",
		DecodeFn:    vp9Decode,
	})
}

func vp9Decode(d *decode.D) interface{} {

	// TODO: header_size at end? even for show_existing_frame?

	d.FieldU2("frame_marker")
	profileLowBit := d.FieldU1("profile_low_bit")
	profileHighBit := d.FieldU1("profile_high_bit")
	profile := int(profileHighBit<<1 + profileLowBit)
	if profile == 3 {
		d.FieldU2("reserved_zero")
	}
	if d.FieldBool("show_existing_frame") {
		d.FieldU2("frame_to_show_map_idx")
		return nil
	}

	d.FieldUFn("frame_type", func() (uint64, decode.DisplayFormat, string) {
		switch d.U1() {
		case 0:
			return 0, decode.NumberDecimal, "key_frame"
		case 1:
			return 0, decode.NumberDecimal, "non_key_frame"
		}
		panic("unreachable")
	})
	d.FieldU1("show_frame")
	d.FieldU1("error_resilient_mode")

	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
