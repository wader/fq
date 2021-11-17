package vpx

// https://www.webmproject.org/vp9/mp4/

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.VPX_CCR,
		Description: "VPX Codec Configuration Record",
		DecodeFn:    vpxCCRDecode,
	})
}

func vpxCCRDecode(d *decode.D, in interface{}) interface{} {
	d.FieldU8("profile")
	d.FieldU8("level", d.MapUToStrSym(vpxLevelNames))
	d.FieldU4("bit_depth")
	d.FieldU3("chroma_subsampling", d.MapUToStrSym(vpxChromeSubsamplingNames))
	d.FieldU1("video_full_range_flag")
	d.FieldU8("colour_primaries")
	d.FieldU8("transfer_characteristics")
	d.FieldU8("matrix_coefficients")
	_ = d.FieldU16("codec_initialization_data_size")
	// d.FieldRawLen("codec_initialization_data", int64(initDataSize)*8)

	return nil
}
