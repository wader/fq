package mpeg

// TODO: unescape configurable? merge with AVC_NAL? merge with HEVC?
// TODO: naming

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var avcSPSFormat []*decode.Format
var avcPPSFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_NALU,
		Description: "H.264/AVC network access layer unit",
		DecodeFn:    avcNALDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_AVC_SPS}, Formats: &avcSPSFormat},
			{Names: []string{format.MPEG_AVC_PPS}, Formats: &avcPPSFormat},
		},
	})
}

func avcNALDecode(d *decode.D, in interface{}) interface{} {
	d.FieldBool("forbidden_zero_bit")
	d.FieldU2("nal_ref_idc")
	nalType, _ := d.FieldStringMapFn("nal_unit_type", avcNALNames, "Unknown", d.U5, decode.NumberDecimal)
	unescapedBb := decode.MustNewBitBufFromReader(decode.NALUnescapeReader{Reader: d.BitBufRange(d.Pos(), int64(d.BitsLeft()))})

	switch nalType {
	case avcNALSequenceParameterSet:
		d.FieldDecodeBitBuf("sps", unescapedBb, avcSPSFormat)
	case avcNALPictureParameterSet:
		d.FieldDecodeBitBuf("pps", unescapedBb, avcPPSFormat)
	}
	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
