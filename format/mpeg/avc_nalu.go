package mpeg

// TODO: unescape configurable? merge with AVC_NAL? merge with HEVC?

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

func sEV(d *decode.D) int64 {
	v := expGolomb(d) + 1
	return num.ZigZag(v) - -int64(v&1)
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

var avcNALNames = decode.UToScalar{
	1:                                        {Sym: "SLICE", Description: "Coded slice of a non-IDR picture"},
	2:                                        {Sym: "DPA", Description: "Coded slice data partition A"},
	3:                                        {Sym: "DPB", Description: "Coded slice data partition B"},
	4:                                        {Sym: "DPC", Description: "Coded slice data partition C"},
	5:                                        {Sym: "IDR_SLICE", Description: "Coded slice of an IDR picture"},
	avcNALSupplementalEnhancementInformation: {Sym: "SEI", Description: "Supplemental enhancement information"},
	avcNALSequenceParameterSet:               {Sym: "SPS", Description: "Sequence parameter set"},
	avcNALPictureParameterSet:                {Sym: "PPS", Description: "Picture parameter set"},
	9:                                        {Sym: "AUD", Description: "Access unit delimiter"},
	10:                                       {Sym: "EOSEQ", Description: "End of sequence"},
	11:                                       {Sym: "EOS", Description: "End of stream"},
	12:                                       {Sym: "FILLER", Description: "Filler data"},
	13:                                       {Sym: "SPS_EXT", Description: "Sequence parameter set extension"},
	14:                                       {Sym: "PREFIX", Description: "Prefix NAL unit"},
	15:                                       {Sym: "SUB_SPS", Description: "Subset sequence parameter set"},
	19:                                       {Sym: "AUX_SLICE", Description: "Coded slice of an auxiliary coded picture without partitioning"},
	20:                                       {Sym: "EXTEN_SLICE", Description: "Coded slice extension"},
}

var sliceNames = decode.UToStr{
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
	nalType := d.FieldU5("nal_unit_type", d.MapUToScalar(avcNALNames))
	unescapedBb := decode.MustNewBitBufFromReader(d, decode.NALUnescapeReader{Reader: d.BitBufRange(d.Pos(), d.BitsLeft())})

	switch nalType {
	case avcNALCodedSliceNonIDR,
		avcNALCodedSlicePartitionA,
		avcNALCodedSlicePartitionB,
		avcNALCodedSlicePartitionC,
		avcNALCodedSliceIDR,
		avcNALCodedSliceAuxWithoutPartition,
		avcNALCodedSliceExtension:
		d.FieldStruct("slice_header", func(d *decode.D) {
			d.FieldUFn("first_mb_in_slice", uEV)
			d.FieldUFn("slice_type", uEV, d.MapUToStr(sliceNames))
			d.FieldUFn("pic_parameter_set_id", uEV)
			// TODO: if ( separate_colour_plane_flag from SPS ) colour_plane_id; frame_num
		})
	case avcNALSupplementalEnhancementInformation:
		d.FieldFormatBitBuf("sei", unescapedBb, avcSEIFormat, nil)
	case avcNALSequenceParameterSet:
		d.FieldFormatBitBuf("sps", unescapedBb, avcSPSFormat, nil)
	case avcNALPictureParameterSet:
		d.FieldFormatBitBuf("pps", unescapedBb, avcPPSFormat, nil)
	}
	d.FieldRawLen("data", d.BitsLeft())

	return nil
}
