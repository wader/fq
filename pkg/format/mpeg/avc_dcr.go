package mpeg

// ISO/IEC 14496-15 AVC file format, 5.3.3.1.2 Syntax
// ISO_IEC_14496-10 AVC

// TODO: PPS
// TODO: use avcLevels
// TODO: nal unescape function?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"io"
)

var avcSPSFormat []*decode.Format
var avcPPSFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_DCR,
		Description: "H.264/AVC Decoder configuration record",
		DecodeFn:    avcDcrDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_AVC_SPS}, Formats: &avcSPSFormat},
			{Names: []string{format.MPEG_AVC_PPS}, Formats: &avcPPSFormat},
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

type nalUnescapeReader struct {
	io.Reader
	lastTwoZeros [2]bool
}

func (r nalUnescapeReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	ni := 0
	for i, b := range p[0:n] {
		if r.lastTwoZeros[0] && r.lastTwoZeros[1] && b == 0x03 {
			n--
			r.lastTwoZeros[0] = false
			r.lastTwoZeros[1] = false
			continue
		} else {
			r.lastTwoZeros[1] = r.lastTwoZeros[0]
			r.lastTwoZeros[0] = b == 0
		}
		p[ni] = p[i]
		ni++
	}

	return n, err
}

const (
	avcNALSequenceParameterSet = 7
	avcNALPictureParameterSet  = 8
)

var avcNALNames = map[uint64]string{
	avcNALSequenceParameterSet: "SequenceParameterSet",
	avcNALPictureParameterSet:  "PictureParameterSet",
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

var avcLevels = map[uint64]avcLevel{
	0:  {"1", 1485, 99, 396, 64, 175, 64, 2, 0},
	1:  {"1b", 1485, 99, 396, 128, 350, 64, 2, 0},
	2:  {"1.1", 3000, 396, 900, 192, 500, 128, 2, 0},
	3:  {"1.2", 6000, 396, 2376, 384, 1000, 128, 2, 0},
	4:  {"1.3", 11880, 396, 2376, 768, 2000, 128, 2, 0},
	5:  {"2", 11880, 396, 2376, 2000, 2000, 128, 2, 0},
	6:  {"2.1", 19800, 792, 4752, 4000, 4000, 256, 2, 0},
	7:  {"2.2", 20250, 1620, 8100, 4000, 4000, 256, 2, 0},
	8:  {"3", 40500, 1620, 8100, 10000, 10000, 256, 2, 32},
	9:  {"3.1", 108000, 3600, 18000, 14000, 14000, 512, 4, 16},
	10: {"3.2", 216000, 5120, 20480, 20000, 20000, 512, 4, 16},
	11: {"4", 245760, 8192, 32768, 20000, 25000, 512, 4, 16},
	12: {"4.1", 245760, 8192, 32768, 50000, 62500, 512, 2, 16},
	13: {"4.2", 522240, 8704, 34816, 50000, 62500, 512, 2, 16},
	14: {"5", 589824, 22080, 110400, 135000, 135000, 512, 2, 16},
	15: {"5.1", 983040, 36864, 184320, 240000, 240000, 512, 2, 16},
	16: {"5.2", 2073600, 36864, 184320, 240000, 240000, 512, 2, 16},
	17: {"6", 4177920, 139264, 696320, 240000, 240000, 8192, 2, 16},
	18: {"6.1", 8355840, 139264, 696320, 480000, 480000, 8192, 2, 16},
	19: {"6.2", 16711680, 139264, 696320, 800000, 800000, 8192, 2, 16},
}

func avcDcrParameterSet(d *decode.D, numParamSets uint64) {
	for i := uint64(0); i < numParamSets; i++ {
		d.FieldStructFn("set", func(d *decode.D) {
			paramSetLen := d.FieldU16("length")
			d.DecodeLenFn(int64(paramSetLen)*8, func(d *decode.D) {
				d.FieldBool("forbidden_zero_bit")
				d.FieldU2("nal_ref_idc")
				nalType, _ := d.FieldStringMapFn("nal_unit_type", avcNALNames, "Unknown", d.U5, decode.NumberDecimal)
				unescapedBb := decode.MustNewBitBufFromReader(nalUnescapeReader{Reader: d.BitBufRange(d.Pos(), int64(paramSetLen-1)*8)})

				switch nalType {
				case avcNALSequenceParameterSet:
					d.FieldDecodeBitBuf("nal", unescapedBb, avcSPSFormat)
				case avcNALPictureParameterSet:
					d.FieldDecodeBitBuf("nal", unescapedBb, avcPPSFormat)
				}

				d.FieldBitBufLen("data", d.BitsLeft())

				// 	d.FieldDecodeBitBuf()

				// 	unescapedBb := decode.MustNewBitBufFromReader(nalUnescapeReader{Reader: d.BitBufRange(d.Pos(), int64(paramSetLen-1)*8)})
				// 	d.FieldDecodeBitBuf("unescaped", unescapedBb, decode.FormatFn(func(d *decode.D, in interface{}) interface{} {

				// 		switch nalType {
				// 		case avcNALSequenceParameterSet:
				// 			d.Decode(avcSPSFormat)
				// 		case avcNALPictureParameterSet:
				// 			d.Decode(avcPPSFormat)
				// 		default:
				// 			d.FieldBitBufLen("data", d.BitsLeft())
				// 		}
				// 		return nil
				// 	}))
				// })
			})
		})
	}
}

func avcDcrDecode(d *decode.D, in interface{}) interface{} {
	var profileIdc uint64

	d.FieldU8("configuration_version")
	d.FieldU8("profile_indication")
	d.FieldU8("profile_compatibility")
	d.FieldU8("level_indication")
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

	_ = profileIdc

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
