package mpeg

// TODO: unescape configurable? merge with AVC_NAL? merge with HEVC?
// TODO: naming

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/decode"
)

var avcSPSFormat []*decode.Format
var avcPPSFormat []*decode.Format
var avcSEIFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.AVC_NALU,
		Description: "H.264/AVC Network Access Layer Unit",
		DecodeFn:    avcNALUDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AVC_SPS}, Formats: &avcSPSFormat},
			{Names: []string{format.AVC_PPS}, Formats: &avcPPSFormat},
			{Names: []string{format.AVC_SEI}, Formats: &avcSEIFormat},
		},
	})
}

// 14496-10 9.1 Parsing process for Exp-Golomb codes
func expGolomb(d *decode.D) uint64 {
	leadingZeroBits := -1
	for b := false; !b; leadingZeroBits++ {
		b = d.Bool()
	}

	var expN uint64
	if leadingZeroBits == 0 {
		expN = 1
	} else {
		expN = 2 << (leadingZeroBits - 1)
	}

	return expN - 1 + d.U(leadingZeroBits)
}

func uEV(d *decode.D) uint64 { return expGolomb(d) }

func fieldUEV(d *decode.D, name string) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return uEV(d), decode.NumberDecimal, ""
	})
}

func sEV(d *decode.D) int64 { return num.ZigZag(expGolomb(d)) }

func fieldSEV(d *decode.D, name string) int64 { //nolint:unparam
	return d.FieldSFn(name, func() (int64, decode.DisplayFormat, string) {
		return sEV(d), decode.NumberDecimal, ""
	})
}

const (
	avcNALCodedSliceNonIDR                   = 1
	avcNALCodedSlicePartitionA               = 2
	avcNALCodedSlicePartitionB               = 3
	avcNALCodedSlicePartitionC               = 4
	avcNALCodedSliceIDR                      = 5
	avcNALSupplementalEnhancementInformation = 6
	avcNALSequenceParameterSet               = 7
	avcNALPictureParameterSet                = 8
	avcNALCodedSliceAuxWithoutPartition      = 19
	avcNALCodedSliceExtension                = 20
)

var avcNALNames = map[uint64]string{
	1:                                        "Coded slice of a non-IDR picture",
	2:                                        "Coded slice data partition A",
	3:                                        "Coded slice data partition B",
	4:                                        "Coded slice data partition C",
	5:                                        "Coded slice of an IDR picture",
	avcNALSupplementalEnhancementInformation: "Supplemental enhancement information (SEI)",
	avcNALSequenceParameterSet:               "Sequence parameter set",
	avcNALPictureParameterSet:                "Picture parameter set",
	9:                                        "Access unit delimiter",
	10:                                       "End of sequence",
	11:                                       "End of stream",
	12:                                       "Filler data",
	13:                                       "Sequence parameter set extension",
	14:                                       "Prefix NAL unit",
	15:                                       "Subset sequence parameter set",
	19:                                       "Coded slice of an auxiliary coded picture without partitioning",
	20:                                       "Coded slice extension",
}

var sliceNames = map[uint64]string{
	0: "P",
	1: "B",
	2: "I",
	3: "SP",
	4: "SI",
	5: "P",
	6: "B",
	7: "I",
	8: "SP",
	9: "SI",
}

func avcNALUDecode(d *decode.D, in interface{}) interface{} {
	d.FieldBool("forbidden_zero_bit")
	d.FieldU2("nal_ref_idc")
	nalType, _ := d.FieldStringMapFn("nal_unit_type", avcNALNames, "Unknown", d.U5, decode.NumberDecimal)
	unescapedBb := decode.MustNewBitBufFromReader(decode.NALUnescapeReader{Reader: d.BitBufRange(d.Pos(), d.BitsLeft())})

	switch nalType {
	case avcNALCodedSliceNonIDR,
		avcNALCodedSlicePartitionA,
		avcNALCodedSlicePartitionB,
		avcNALCodedSlicePartitionC,
		avcNALCodedSliceIDR,
		avcNALCodedSliceAuxWithoutPartition,
		avcNALCodedSliceExtension:
		d.FieldStructFn("slice_header", func(d *decode.D) {
			fieldUEV(d, "first_mb_in_slice")
			d.FieldStringMapFn("slice_type", sliceNames, "Unknown", func() uint64 { return uEV(d) }, decode.NumberDecimal)
			fieldUEV(d, "pic_parameter_set_id")
			// TODO: if ( separate_colour_plane_flag from SPS ) colour_plane_id; frame_num
		})
	case avcNALSupplementalEnhancementInformation:
		d.FieldFormatBitBuf("sei", unescapedBb, avcSEIFormat)
	case avcNALSequenceParameterSet:
		d.FieldFormatBitBuf("sps", unescapedBb, avcSPSFormat)
	case avcNALPictureParameterSet:
		d.FieldFormatBitBuf("pps", unescapedBb, avcPPSFormat)
	}
	d.FieldBitBufLen("data", d.BitsLeft())

	return nil
}
