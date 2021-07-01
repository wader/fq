package vpx

// https://www.webmproject.org/vp9/mp4/

import (
	"fq/format"
	"fq/pkg/decode"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.VPX_CCR,
		Description: "VPX Codec Configuration Record",
		DecodeFn:    vpxCCRDecode,
	})
}

func vpxCCRDecode(d *decode.D, in interface{}) interface{} {
	d.FieldU8("profile")
	d.FieldStringMapFn("level", vpxLevelNames, "Unknown", d.U8, decode.NumberDecimal)
	d.FieldU4("bit_depth")
	d.FieldStringMapFn("chroma_subsampling", vpxChromeSubsamplingNames, "Unknown", d.U3, decode.NumberDecimal)
	d.FieldU1("video_full_range_flag")
	d.FieldU8("colour_primaries")
	d.FieldU8("transfer_characteristics")
	d.FieldU8("matrix_coefficients")
	_ = d.FieldU16("codec_intialization_data_size")
	// d.FieldBitBufLen("codec_intialization_data", int64(initDataSize)*8)

	return nil
}
