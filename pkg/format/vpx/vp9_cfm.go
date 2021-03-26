package vpx

// https://www.webmproject.org/docs/container/#vp9-codec-feature-metadata-codecprivate

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.VP9_CFM,
		Description: "VP9 Codec Feature Metadata",
		DecodeFn:    vp9CFMDecode,
	})
}

func vp9CFMDecode(d *decode.D, in interface{}) interface{} {
	d.FieldArrayFn("features", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStructFn("feature", func(d *decode.D) {
				id, _ := d.FieldStringMapFn("id", vp9FeatureIDNames, "Unknown", d.U8, decode.NumberDecimal)
				l := d.FieldU8("length")
				d.DecodeLenFn(int64(l)*8, func(d *decode.D) {
					switch id {
					case vp9FeatureProfile:
						d.FieldU8("profile")
					case vp9FeatureLevel:
						d.FieldStringMapFn("level", vpxLevelNames, "Unknown", d.U8, decode.NumberDecimal)
					case vp9FeatureBitDepth:
						d.FieldU8("bit_depth")
					case vp9FeatureChromaSubsampling:
						d.FieldStringMapFn("chroma_subsampling", vpxChromeSubsamplingNames, "Unknown", d.U8, decode.NumberDecimal)
					default:
						d.FieldBitBufLen("data", d.BitsLeft())
					}
				})
			})
		}
	})

	return nil
}
