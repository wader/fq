package mpeg

// ISO/IEC 14496-15 AVC file format, 5.3.3.1.2 Syntax
// ISO_IEC_14496-10 AVC

// TODO: PPS
// TODO: use avcLevels
// TODO: nal unescape function?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var avcNALFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_DCR,
		Description: "H.264/AVC Decoder configuration record",
		DecodeFn:    avcDcrDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_AVC_NALU}, Formats: &avcNALFormat},
		},
	})
}

// TODO: share?
func zigzag(n uint64) int64 {
	return int64(n>>1 ^ -(n & 1))
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

func sEV(d *decode.D) int64 { return zigzag(expGolomb(d)) }

func fieldSEV(d *decode.D, name string) int64 {
	return d.FieldSFn(name, func() (int64, decode.DisplayFormat, string) {
		return sEV(d), decode.NumberDecimal, ""
	})
}

const (
	avcNALSequenceParameterSet = 7
	avcNALPictureParameterSet  = 8
)

var avcNALNames = map[uint64]string{
	1:                          "Coded slice of a non-IDR picture",
	2:                          "Coded slice data partition A",
	3:                          "Coded slice data partition B",
	4:                          "Coded slice data partition C",
	5:                          "Coded slice of an IDR picture",
	6:                          "Supplemental enhancement information (SEI)",
	avcNALSequenceParameterSet: "Sequence parameter set",
	avcNALPictureParameterSet:  "Picture parameter set",
	9:                          "Access unit delimiter",
	10:                         "End of sequence",
	11:                         "End of stream",
	12:                         "Filler data",
	13:                         "Sequence parameter set extension",
	14:                         "Prefix NAL unit",
	15:                         "Subset sequence parameter set",
	19:                         "Coded slice of an auxiliary coded picture without partitioning",
	20:                         "Coded slice extension",
}

type avcLevel struct {
	Name         string
	MaxMBPS      uint64
	MaxFS        uint64
	MaxDpbMbs    uint64
	MaxBR        uint64
	MaxCPB       uint64
	MaxVmvR      uint64
	MinCR        uint64
	MaxMvsPer2Mb uint64
}

var avcProfileNames = map[uint64]string{
	// 66: "Constrained Baseline Profile", // (CBP, 66 with constraint set 1)
	66:  "Baseline Profile",
	88:  "Extended Profile",
	77:  "Main Profile",
	100: "High Profile",
	//100: "Constrained High Profile", // (100 with constraint set 4 and 5)
	110: "High 10 Profile",
	122: "High 4:2:2 Profile",
	244: "High 4:4:4 Predictive Profile",
	// 110: "High 10 Intra Profile", // (110 with constraint set 3)
	// 122: "High 4:2:2 Intra Profile", // (122 with constraint set 3)
	// 244: "High 4:4:4 Intra Profile", // (244 with constraint set 3)
	44:  "CAVLC 4:4:4 Intra Profile",
	83:  "Scalable Baseline Profile",
	86:  "Scalable High Profile",
	128: "Stereo High Profile",
	134: "MFC High Profile",
	138: "Multiview Depth High Profile",
	139: "Enhanced Multiview Depth High Profile",
}

// TODO: 1b contraint flag 1?
var avcLevelNames = map[uint64]string{
	10: "1",
	//10:  "1b"
	11: "1.1",
	12: "1.2",
	13: "1.3",
	20: "2",
	21: "2.1",
	22: "2.2",
	30: "3",
	31: "3.1",
	32: "3.2",
	40: "4",
	41: "4.1",
	42: "4.2",
	50: "5",
	51: "5.1",
	52: "5.2",
	60: "6",
	61: "6.1",
	62: "6.2",
}

// TODO: 1b contraint flag 1?
var avcLevels = map[uint64]avcLevel{
	10: {"1", 1485, 99, 396, 64, 175, 64, 2, 0},
	//10:  {"1b", 1485, 99, 396, 128, 350, 64, 2, 0}, //
	11: {"1.1", 3000, 396, 900, 192, 500, 128, 2, 0},
	12: {"1.2", 6000, 396, 2376, 384, 1000, 128, 2, 0},
	13: {"1.3", 11880, 396, 2376, 768, 2000, 128, 2, 0},
	20: {"2", 11880, 396, 2376, 2000, 2000, 128, 2, 0},
	21: {"2.1", 19800, 792, 4752, 4000, 4000, 256, 2, 0},
	22: {"2.2", 20250, 1620, 8100, 4000, 4000, 256, 2, 0},
	30: {"3", 40500, 1620, 8100, 10000, 10000, 256, 2, 32},
	31: {"3.1", 108000, 3600, 18000, 14000, 14000, 512, 4, 16},
	32: {"3.2", 216000, 5120, 20480, 20000, 20000, 512, 4, 16},
	40: {"4", 245760, 8192, 32768, 20000, 25000, 512, 4, 16},
	41: {"4.1", 245760, 8192, 32768, 50000, 62500, 512, 2, 16},
	42: {"4.2", 522240, 8704, 34816, 50000, 62500, 512, 2, 16},
	50: {"5", 589824, 22080, 110400, 135000, 135000, 512, 2, 16},
	51: {"5.1", 983040, 36864, 184320, 240000, 240000, 512, 2, 16},
	52: {"5.2", 2073600, 36864, 184320, 240000, 240000, 512, 2, 16},
	60: {"6", 4177920, 139264, 696320, 240000, 240000, 8192, 2, 16},
	61: {"6.1", 8355840, 139264, 696320, 480000, 480000, 8192, 2, 16},
	62: {"6.2", 16711680, 139264, 696320, 800000, 800000, 8192, 2, 16},
}

func avcDcrParameterSet(d *decode.D, numParamSets uint64) {
	for i := uint64(0); i < numParamSets; i++ {
		d.FieldStructFn("set", func(d *decode.D) {
			paramSetLen := d.FieldU16("length")
			d.FieldDecodeLen("nal", int64(paramSetLen)*8, avcNALFormat)
		})
	}
}

func avcDcrDecode(d *decode.D, in interface{}) interface{} {
	d.FieldU8("configuration_version")
	d.FieldStringMapFn("profile_indication", avcProfileNames, "Unknown", d.U8, decode.NumberDecimal)
	d.FieldU8("profile_compatibility")
	d.FieldStringMapFn("level_indication", avcLevelNames, "Unknown", d.U8, decode.NumberDecimal)
	d.FieldU6("reserved0")
	lengthSizeMinusOne := d.FieldU2("length_size_minus_one")
	d.FieldU3("reserved1")
	numSeqParamSets := d.FieldU5("num_of_sequence_parameter_sets")
	d.FieldArrayFn("sequence_parameter_sets", func(d *decode.D) {
		avcDcrParameterSet(d, numSeqParamSets)
	})
	numPicParamSets := d.FieldU8("num_of_picture_parameter_sets")
	d.FieldArrayFn("picture_parameter_sets", func(d *decode.D) {
		avcDcrParameterSet(d, numPicParamSets)
	})

	if d.BitsLeft() > 0 {
		d.FieldBitBufLen("data", d.BitsLeft())
	}

	// TODO:
	// Compatible extensions to this record will extend it and will not change the configuration version code. Readers
	// should be prepared to ignore unrecognized data beyond the definition of the data they understand (e.g. after
	// the parameter sets in this specification).

	// TODO: something wrong here, seen files with profileIdc = 100 with no bytes after picture_parameter_sets
	// https://github.com/FFmpeg/FFmpeg/blob/069d2b4a50a6eb2f925f36884e6b9bd9a1e54670/libavcodec/h264_ps.c#L333

	// switch profileIdc {
	// case 100, 110, 122, 144:
	// 	d.FieldU6("reserved2")
	// 	d.FieldU6("chroma_format")
	// 	d.FieldU4("reserved3")
	// 	d.FieldU3("bit_depth_luma_minus8")
	// 	d.FieldU5("reserved4")
	// 	d.FieldU3("bit_depth_chroma_minus8")
	// 	numSeqParamSetExt := d.FieldU5("num_of_sequence_parameter_set_ext")
	// 	d.FieldArrayFn("parameter_set_exts", func(d *decode.D) {
	// 		for i := uint64(0); i < numSeqParamSetExt; i++ {
	// 			d.FieldStructFn("parameter_set_ext", func(d *decode.D) {
	// 				paramSetLen := d.FieldU16("length")
	// 				d.FieldBitBufLen("set", int64(paramSetLen)*8)
	// 			})
	// 		}
	// 	})
	// }

	return format.AvcDcrOut{LengthSize: lengthSizeMinusOne + 1}
}
