package av1

// https://aomediacodec.github.io/av1-spec/av1-spec.pdf
// https://github.com/ietf-wg-cellar/matroska-specification/blob/master/codec/av1.md
// https://cdn.rawgit.com/AOMediaCodec/av1-isobmff/v1.0.0/

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var av1CCRav1OBUGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.AV1_CCR,
		&decode.Format{
			Description: "AV1 Codec Configuration Record",
			DecodeFn:    ccrDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AV1_OBU}, Out: &av1CCRav1OBUGroup},
			},
		})
}

func ccrDecode(d *decode.D) any {
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
		d.FieldU4("initial_presentation_delay", scalar.UintActualAdd(1))
	} else {
		d.FieldU4("reserved")
	}
	d.FieldArray("config_obus", func(d *decode.D) {
		for d.BitsLeft() > 0 {
			d.FieldFormat("config_obu", &av1CCRav1OBUGroup, nil)
		}
	})

	return nil
}
