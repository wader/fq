package vpx

// https://storage.googleapis.com/downloads.webmproject.org/docs/vp9/vp9-bitstream-specification-v0.6-20160331-draft.pdf

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: vpx frame?

const (
	vp9FeatureProfile           = 1
	vp9FeatureLevel             = 2
	vp9FeatureBitDepth          = 3
	vp9FeatureChromaSubsampling = 4
)

var vp9FeatureIDNames = scalar.UintMapSymStr{
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

var vp9ColorSpaceNames = scalar.UintMapSymStr{
	CS_UNKNOWN:   "unknown",
	CS_BT_601:    "bt_601",
	CS_BT_709:    "bt_709",
	CS_SMPTE_170: "smpte_170",
	CS_SMPTE_240: "smpte_240",
	CS_BT_2020:   "bt_2020",
	CS_RESERVED:  "reserved",
	CS_RGB:       "rgb",
}

var vp9ProfilesMap = scalar.UintMapDescription{
	0: "8 bit/sample, chroma subsampling: 4:2:0",
	1: "8 bit, chroma subsampling: 4:2:2, 4:4:0, 4:4:4",
	2: "10–12 bit, chroma subsampling: 4:2:0",
	3: "10–12 bit, chroma subsampling: 4:2:2, 4:4:0, 4:4:4",
}

func init() {
	interp.RegisterFormat(
		format.VP9_Frame,
		&decode.Format{
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
	d.FieldValueUint("bit_depth", uint64(bitDepth))
	colorSpace := d.FieldU3("color_space", vp9ColorSpaceNames)
	_, colorSpaceOk := vp9ColorSpaceNames[colorSpace]
	if !colorSpaceOk || colorSpace != CS_RGB {
		d.FieldU1("color_range")
		if profile == 1 || profile == 3 {
			d.FieldU1("subsampling_x")
			d.FieldU1("subsampling_y")
			d.FieldU1("reserved_zero1")
		} else {
			d.FieldValueUint("subsampling_x", 1)
			d.FieldValueUint("subsampling_y", 1)
		}
	} else {
		d.FieldValueUint("color_range", 1)
		if profile == 1 || profile == 3 {
			d.FieldValueUint("subsampling_x", 0)
			d.FieldValueUint("subsampling_y", 0)
			d.FieldU1("reserved_zero2")
		}
	}
}

func vp9DecodeFrameSize(d *decode.D) {
	d.FieldUintFn("frame_width", func(d *decode.D) uint64 { return d.U16() + 1 })
	d.FieldUintFn("frame_height", func(d *decode.D) uint64 { return d.U16() + 1 })
}

func vp9Decode(d *decode.D) any {

	// TODO: header_size at end? even for show_existing_frame?

	d.FieldU2("frame_marker")
	profileLowBit := d.FieldU1("profile_low_bit")
	profileHighBit := d.FieldU1("profile_high_bit")
	profile := int(profileHighBit<<1 + profileLowBit)
	d.FieldValueUint("profile", uint64(profile), vp9ProfilesMap)
	if profile == 3 {
		d.FieldU1("reserved_zero0")
	}
	showExistingFrame := d.FieldBool("show_existing_frame")
	if showExistingFrame {
		d.FieldU2("frame_to_show_map_idx")
		return nil
	}

	frameType := d.FieldBool("frame_type", scalar.BoolMapSymStr{true: "non_key_frame", false: "key_frame"})
	d.FieldU1("show_frame")
	d.FieldU1("error_resilient_mode")

	if !frameType {
		// is key frame
		vp9DecodeFrameSyncCode(d)
		vp9DecodeColorConfig(d, profile)
		vp9DecodeFrameSize(d)
	}

	d.FieldRawLen("data", d.BitsLeft())

	return nil
}
