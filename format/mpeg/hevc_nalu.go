package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var hevcVPSGroup decode.Group
var hevcPPSGroup decode.Group
var hevcSPSGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.HEVC_NALU,
		&decode.Format{
			Description: "H.265/HEVC Network Access Layer Unit",
			DecodeFn:    hevcNALUDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.HEVC_VPS}, Out: &hevcVPSGroup},
				{Groups: []*decode.Group{format.HEVC_PPS}, Out: &hevcPPSGroup},
				{Groups: []*decode.Group{format.HEVC_SPS}, Out: &hevcSPSGroup},
			},
		})
}

const (
	hevcNALNUTVPS = 32
	hevcNALNUTSPS = 33
	hevcNALNUTPPS = 34
)

var hevcNALNames = scalar.UintMapSymStr{
	0:             "TRAIL_N",
	1:             "TRAIL_R",
	2:             "TSA_N",
	3:             "TSA_R",
	4:             "STSA_N",
	5:             "STSA_R",
	6:             "RADL_N",
	7:             "RADL_R",
	8:             "RASL_N",
	9:             "RASL_R",
	10:            "RSV_VCL_N10",
	12:            "RSV_VCL_N12",
	14:            "RSV_VCL_N14",
	11:            "RSV_VCL_R11",
	13:            "RSV_VCL_R13",
	15:            "RSV_VCL_R15",
	16:            "BLA_W_LP",
	17:            "BLA_W_RADL",
	18:            "BLA_N_LP",
	19:            "IDR_W_RADL",
	20:            "IDR_N_LP",
	21:            "CRA_NUT",
	22:            "RSV_IRAP_VCL22",
	23:            "RSV_IRAP_VCL23",
	24:            "RSV_VCL24",
	25:            "RSV_VCL25",
	26:            "RSV_VCL26",
	27:            "RSV_VCL27",
	28:            "RSV_VCL28",
	29:            "RSV_VCL29",
	30:            "RSV_VCL30",
	31:            "RSV_VCL31",
	hevcNALNUTVPS: "VPS_NUT",
	hevcNALNUTSPS: "SPS_NUT",
	hevcNALNUTPPS: "PPS_NUT",
	35:            "AUD_NUT",
	36:            "EOS_NUT",
	37:            "EOB_NUT",
	38:            "FD_NUT",
	39:            "PREFIX_SEI_NUT",
	40:            "SUFFIX_SEI_NUT",
	41:            "RSV_NVCL41",
	42:            "RSV_NVCL42",
	43:            "RSV_NVCL43",
	44:            "RSV_NVCL44",
	45:            "RSV_NVCL45",
	46:            "RSV_NVCL46",
	47:            "RSV_NVCL47",
}

func hevcNALUDecode(d *decode.D) any {
	d.FieldBool("forbidden_zero_bit")
	nalType := d.FieldU6("nal_unit_type", hevcNALNames)
	d.FieldU6("nuh_layer_id")
	d.FieldU3("nuh_temporal_id_plus1")
	unescapedBR := d.NewBitBufFromReader(nalUnescapeReader{Reader: bitio.NewIOReader(d.BitBufRange(d.Pos(), d.BitsLeft()))})

	switch nalType {
	case hevcNALNUTVPS:
		d.FieldFormatBitBuf("vps", unescapedBR, &hevcVPSGroup, nil)
	case hevcNALNUTPPS:
		d.FieldFormatBitBuf("pps", unescapedBR, &hevcPPSGroup, nil)
	case hevcNALNUTSPS:
		d.FieldFormatBitBuf("sps", unescapedBR, &hevcSPSGroup, nil)
	}
	d.FieldRawLen("data", d.BitsLeft())

	return nil
}
