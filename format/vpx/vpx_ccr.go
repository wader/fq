package vpx

// https://www.webmproject.org/vp9/mp4/

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.VPX_CCR,
		Description: "VPX Codec Configuration Record",
		DecodeFn:    vpxCCRDecode,
	})
}

func vpxCCRDecode(d *decode.D, in any) any {
	d.FieldU8("profile")
	d.FieldU8("level", vpxLevelNames)
	d.FieldU4("bit_depth")
	d.FieldU3("chroma_subsampling", vpxChromeSubsamplingNames)
	d.FieldU1("video_full_range_flag")
	d.FieldU8("colour_primaries", format.ISO_23091_2_ColourPrimariesMap)
	d.FieldU8("transfer_characteristics", format.ISO_23091_2_TransferCharacteristicMap)
	d.FieldU8("matrix_coefficients", format.ISO_23091_2_MatrixCoefficients)
	_ = d.FieldU16("codec_initialization_data_size")
	// d.FieldRawLen("codec_initialization_data", int64(initDataSize)*8)

	return nil
}
