package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:         format.HEVC_NALU,
		Description:  "H.265/HEVC Network Access Layer Unit",
		DecodeFn:     hevcNALUDecode,
		Dependencies: []decode.Dependency{},
	})
}

var hevcNALNames = map[uint64]string{
	0:  "TRAIL_N",
	1:  "TRAIL_R",
	2:  "TSA_N",
	3:  "TSA_R",
	4:  "STSA_N",
	5:  "STSA_R",
	6:  "RADL_N",
	7:  "RADL_R",
	8:  "RASL_N",
	9:  "RASL_R",
	10: "RSV_VCL_N10",
	12: "RSV_VCL_N12",
	14: "RSV_VCL_N14",
	11: "RSV_VCL_R11",
	13: "RSV_VCL_R13",
	15: "RSV_VCL_R15",
	16: "BLA_W_LP",
	17: "BLA_W_RADL",
	18: "BLA_N_LP",
	19: "IDR_W_RADL",
	20: "IDR_N_LP",
	21: "CRA_NUT",
	22: "RSV_IRAP_VCL22",
	23: "RSV_IRAP_VCL23",
	24: "RSV_VCL24",
	25: "RSV_VCL25",
	26: "RSV_VCL26",
	27: "RSV_VCL27",
	28: "RSV_VCL28",
	29: "RSV_VCL29",
	30: "RSV_VCL30",
	31: "RSV_VCL31",
	32: "VPS_NUT",
	33: "SPS_NUT",
	34: "PPS_NUT",
	35: "AUD_NUT",
	36: "EOS_NUT",
	37: "EOB_NUT",
	38: "FD_NUT",
	39: "PREFIX_SEI_NUT",
	40: "SUFFIX_SEI_NUT",
	41: "RSV_NVCL41",
	42: "RSV_NVCL42",
	43: "RSV_NVCL43",
	44: "RSV_NVCL44",
	45: "RSV_NVCL45",
	46: "RSV_NVCL46",
	47: "RSV_NVCL47",
}

func hevcNALUDecode(d *decode.D, in interface{}) interface{} {
	d.FieldBool("forbidden_zero_bit")
	nalType, _ := d.FieldStringMapFn("nal_unit_type", hevcNALNames, "Unknown", d.U6, decode.NumberDecimal)
	d.FieldU6("nuh_layer_id")
	d.FieldU3("nuh_temporal_id_plus1")
	unescapedBb := decode.MustNewBitBufFromReader(decode.NALUnescapeReader{Reader: d.BitBufRange(d.Pos(), d.BitsLeft())})

	_ = unescapedBb
	_ = nalType

	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
