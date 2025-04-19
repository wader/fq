package mpeg

// TODO: unescape configurable? merge with AVC_NAL? merge with HEVC?

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/mathx"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var avcSPSFormat decode.Group
var avcPPSFormat decode.Group
var avcSEIFormat decode.Group

func init() {
	interp.RegisterFormat(
		format.AVC_NALU,
		&decode.Format{
			Description: "H.264/AVC Network Access Layer Unit",
			DecodeFn:    avcNALUDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AVC_SPS}, Out: &avcSPSFormat},
				{Groups: []*decode.Group{format.AVC_PPS}, Out: &avcPPSFormat},
				{Groups: []*decode.Group{format.AVC_SEI}, Out: &avcSEIFormat},
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

func sEV(d *decode.D) int64 {
	v := expGolomb(d) + 1
	return mathx.ZigZag[uint64, int64](v) - -int64(v&1)
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

var avcNALNames = scalar.UintMap{
	1:                                        {Sym: "slice", Description: "Coded slice of a non-IDR picture"},
	2:                                        {Sym: "dpa", Description: "Coded slice data partition A"},
	3:                                        {Sym: "dpb", Description: "Coded slice data partition B"},
	4:                                        {Sym: "dpc", Description: "Coded slice data partition C"},
	5:                                        {Sym: "idr_slice", Description: "Coded slice of an IDR picture"},
	avcNALSupplementalEnhancementInformation: {Sym: "sei", Description: "Supplemental enhancement information"},
	avcNALSequenceParameterSet:               {Sym: "sps", Description: "Sequence parameter set"},
	avcNALPictureParameterSet:                {Sym: "pps", Description: "Picture parameter set"},
	9:                                        {Sym: "aud", Description: "Access unit delimiter"},
	10:                                       {Sym: "eoseq", Description: "End of sequence"},
	11:                                       {Sym: "eos", Description: "End of stream"},
	12:                                       {Sym: "filler", Description: "Filler data"},
	13:                                       {Sym: "sps_ext", Description: "Sequence parameter set extension"},
	14:                                       {Sym: "prefix", Description: "Prefix NAL unit"},
	15:                                       {Sym: "sub_sps", Description: "Subset sequence parameter set"},
	19:                                       {Sym: "aux_slice", Description: "Coded slice of an auxiliary coded picture without partitioning"},
	20:                                       {Sym: "exten_slice", Description: "Coded slice extension"},
}

var sliceNames = scalar.UintMapSymStr{
	0: "p",
	1: "b",
	2: "i",
	3: "sp",
	4: "si",
	5: "p",
	6: "b",
	7: "i",
	8: "sp",
	9: "si",
}

func findNALUEmulationCode(d *decode.D, maxLen int64) (int64, uint64, error) {
	return d.TryPeekFind(24, 8, maxLen, func(v uint64) bool {
		return v == 0x00_00_03
	})
}

func avcNALUDecode(d *decode.D) any {
	d.FieldBool("forbidden_zero_bit")
	d.FieldU2("nal_ref_idc")
	nalType := d.FieldU5("nal_unit_type", avcNALNames)

	decodeFn := func(d *decode.D) {
		switch nalType {
		case avcNALCodedSliceNonIDR,
			avcNALCodedSlicePartitionA,
			avcNALCodedSlicePartitionB,
			avcNALCodedSlicePartitionC,
			avcNALCodedSliceIDR,
			avcNALCodedSliceAuxWithoutPartition,
			avcNALCodedSliceExtension:
			d.FieldUintFn("first_mb_in_slice", uEV)
			d.FieldUintFn("slice_type", uEV, sliceNames)
			d.FieldUintFn("pic_parameter_set_id", uEV)
			// TODO: if ( separate_colour_plane_flag from SPS ) colour_plane_id; frame_num
		case avcNALSupplementalEnhancementInformation:
			d.Format(&avcSEIFormat, nil)
		case avcNALSequenceParameterSet:
			d.Format(&avcSPSFormat, nil)
		case avcNALPictureParameterSet:
			d.Format(&avcPPSFormat, nil)
		default:
			d.FieldRawLen("data", d.BitsLeft())
		}
	}

	offset, _, _ := findNALUEmulationCode(d, d.BitsLeft())
	if offset < 0 {
		d.FieldStruct("rbsp", decodeFn)
	} else {
		unescapedBR := d.NewBitBufFromReader(&nalUnescapeReader{Reader: bitio.NewIOReader(d.BitBufRange(d.Pos(), d.BitsLeft()))})
		d.FieldStructRootBitBufFn("rbsp", unescapedBR, decodeFn)
	}

	return nil
}
