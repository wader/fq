package av1

// https://aomediacodec.github.io/av1-spec/av1-spec.pdf
// https://github.com/ietf-wg-cellar/matroska-specification/blob/master/codec/av1.md
// https://cdn.rawgit.com/AOMediaCodec/av1-isobmff/v1.0.0/

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.AV1_CCR,
		Description: "AV1 Codec Configuration Record",
		DecodeFn:    ccrDecode,
	})
}

func ccrDecode(d *decode.D, in interface{}) interface{} {
	d.FieldU1("marker")
	d.FieldU7("version")
	d.FieldU3("seq_profile")
	d.FieldU5("seq_level_idx_0")
	d.FieldU1("seq_tier_0")
	d.FieldU1("high_bitdepth")
	d.FieldU1("twelve_bit")
	d.FieldU1("monochrome")
	d.FieldU1("chroma_subsampling_x")
	d.FieldU1("chroma_subsampling_y")
	d.FieldU2("chroma_sample_position")
	d.FieldU3("reserved = 0")
	initalPreDelay := d.FieldBool("initial_presentation_delay_present")
	if initalPreDelay {
		d.FieldU4("initial_presentation_delay_minus_one")
	} else {
		d.FieldU4("reserved")
	}
	if d.BitsLeft() > 0 {
		d.FieldBitBufLen("config_obus", d.BitsLeft())
	}

	return nil
}
