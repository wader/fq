package vpx

// https://storage.googleapis.com/downloads.webmproject.org/docs/vp9/vp9-bitstream-specification-v0.6-20160331-draft.pdf

import (
	"fq/format"
	"fq/pkg/decode"
)

// TODO: vpx frame?

const (
	vp9FeatureProfile           = 1
	vp9FeatureLevel             = 2
	vp9FeatureBitDepth          = 3
	vp9FeatureChromaSubsampling = 4
)

var vp9FeatureIDNames = map[uint64]string{
	vp9FeatureProfile:           "Profile",
	vp9FeatureLevel:             "Level",
	vp9FeatureBitDepth:          "Bit Depth",
	vp9FeatureChromaSubsampling: "Chroma Subsampling",
}

const (
	CS_UNKNOWN   = 0
	CS_BT_601    = 1
	CS_BT_709    = 2
	CS_SMPTE_170 = 3
	CS_SMPTE_240 = 4
	CS_BT_2020   = 5
	CS_RESERVED  = 6
	CS_RGB       = 7
)

var vp9ColorSpaceNames = map[uint64]string{
	CS_UNKNOWN:   "CS_UNKNOWN",
	CS_BT_601:    "CS_BT_601",
	CS_BT_709:    "CS_BT_709",
	CS_SMPTE_170: "CS_SMPTE_170",
	CS_SMPTE_240: "CS_SMPTE_240",
	CS_BT_2020:   "CS_BT_2020",
	CS_RESERVED:  "CS_RESERVED",
	CS_RGB:       "CS_RGB",
}

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.VP9_FRAME,
		Description: "VP9 frame",
		DecodeFn:    vp9Decode,
	})
}

func vp9DecodeFrameSyncCode(d *decode.D) {
	d.FieldU8("frame_sync_byte_0")
	d.FieldU8("frame_sync_byte_1")
	d.FieldU8("frame_sync_byte_2")
}

func vp9DecodeColorConfig(d *decode.D, profile int) {
	bitDepth := 8
	if profile >= 2 {
		tenOrTwelveBit := d.FieldBool("ten_or_twelve_bit")
		if tenOrTwelveBit {
			bitDepth = 12
		} else {
			bitDepth = 10
		}
	}
	d.FieldValueU("bit_depth", uint64(bitDepth), "")
	colorSpace, colorSpaceOk := d.FieldStringMapFn("color_space", vp9ColorSpaceNames, "Unknown", d.U3, decode.NumberDecimal)
	if !colorSpaceOk || colorSpace != CS_RGB {
		d.FieldU1("color_range")
		if profile == 1 || profile == 3 {
			d.FieldU1("subsampling_x")
			d.FieldU1("subsampling_y")
			d.FieldU1("reserved_zero")
		} else {
			d.FieldValueU("subsampling_x", 1, "")
			d.FieldValueU("subsampling_y", 1, "")
		}
	} else {
		d.FieldValueU("color_range", 1, "")
		if profile == 1 || profile == 3 {
			d.FieldValueU("subsampling_x", 0, "")
			d.FieldValueU("subsampling_y", 0, "")
			d.FieldU1("reserved_zero")
		}
	}
}

func vp9DecodeFrameSize(d *decode.D) {
	d.FieldUFn("frame_width", func() (uint64, decode.DisplayFormat, string) { return d.U16() + 1, decode.NumberDecimal, "" })
	d.FieldUFn("frame_height", func() (uint64, decode.DisplayFormat, string) { return d.U16() + 1, decode.NumberDecimal, "" })
}

func vp9Decode(d *decode.D, in interface{}) interface{} {

	// TODO: header_size at end? even for show_existing_frame?

	d.FieldU2("frame_marker")
	profileLowBit := d.FieldU1("profile_low_bit")
	profileHighBit := d.FieldU1("profile_high_bit")
	profile := int(profileHighBit<<1 + profileLowBit)
	d.FieldValueU("profile", uint64(profile), "")
	if profile == 3 {
		d.FieldU1("reserved_zero")
	}
	showExistingFrame := d.FieldBool("show_existing_frame")
	if showExistingFrame {
		d.FieldU2("frame_to_show_map_idx")
		return nil
	}

	frameType, _ := d.FieldBoolMapFn("frame_type", "non_key_frame", "key_frame", d.Bool)
	d.FieldU1("show_frame")
	d.FieldU1("error_resilient_mode")

	if !frameType {
		// is key frame
		vp9DecodeFrameSyncCode(d)
		vp9DecodeColorConfig(d, profile)
		vp9DecodeFrameSize(d)
	}

	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
