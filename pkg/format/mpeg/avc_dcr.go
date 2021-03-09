package mpeg

// ISO/IEC 14496-15 AVC file format, 5.3.3.1.2 Syntax
// ISO_IEC_14496-10 AVC

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_DCR,
		Description: "H.264/AVC Decoder configuration record",
		DecodeFn:    avcDcrDecode,
	})
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
		for i := uint64(0); i < numSeqParamSets; i++ {
			d.FieldStructFn("parameter_set", func(d *decode.D) {
				paramSetLen := d.FieldU16("length")

				d.FieldStructFn("set", func(d *decode.D) {
					d.DecodeLenFn(int64(paramSetLen)*8, func(d *decode.D) {

						d.FieldBool("forbidden_zero_bit")
						d.FieldU2("nal_ref_idc")
						nalUnitType := d.FieldU5("nal_unit_type")

						switch nalUnitType {
						case 7:
							profileIdc = d.FieldU8("profile_idc")
							d.FieldU8("level_idc")
							// TODO: more
							d.FieldBitBufLen("data", d.BitsLeft())
						default:
							d.FieldBitBufLen("data", d.BitsLeft())
						}

					})
				})

			})
		}
	})
	numPicParamSets := d.FieldU8("num_of_picture_parameter_sets")
	d.FieldArrayFn("picture_parameter_sets", func(d *decode.D) {
		for i := uint64(0); i < numPicParamSets; i++ {
			d.FieldStructFn("parameter_set", func(d *decode.D) {
				paramSetLen := d.FieldU16("length")
				d.FieldBitBufLen("set", int64(paramSetLen)*8)

			})
		}
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
