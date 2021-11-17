package vpx

// https://www.webmproject.org/docs/container/#vp9-codec-feature-metadata-codecprivate

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.VP9_CFM,
		Description: "VP9 Codec Feature Metadata",
		DecodeFn:    vp9CFMDecode,
		RootArray:   true,
		RootName:    "features",
	})
}

func vp9CFMDecode(d *decode.D, in interface{}) interface{} {
	for d.NotEnd() {
		d.FieldStruct("feature", func(d *decode.D) {
			id := d.FieldU8("id", d.MapUToStrSym(vp9FeatureIDNames))
			l := d.FieldU8("length")
			d.LenFn(int64(l)*8, func(d *decode.D) {
				switch id {
				case vp9FeatureProfile:
					d.FieldU8("profile")
				case vp9FeatureLevel:
					d.FieldU8("level", d.MapUToStrSym(vpxLevelNames))
				case vp9FeatureBitDepth:
					d.FieldU8("bit_depth")
				case vp9FeatureChromaSubsampling:
					d.FieldU8("chroma_subsampling", d.MapUToStrSym(vpxChromeSubsamplingNames))
				default:
					d.FieldRawLen("data", d.BitsLeft())
				}
			})
		})
	}

	return nil
}
