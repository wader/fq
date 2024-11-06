package av1

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.AV1_OBU,
		&decode.Format{
			Description: "AV1 Open Bitstream Unit",
			DecodeFn:    obuDecode,
		})
}

const (
	OBU_SEQUENCE_HEADER        = 1
	OBU_TEMPORAL_DELIMITER     = 2
	OBU_FRAME_HEADER           = 3
	OBU_TILE_GROUP             = 4
	OBU_METADATA               = 5
	OBU_FRAME                  = 6
	OBU_REDUNDANT_FRAME_HEADER = 7
	OBU_TILE_LIST              = 8
	OBU_PADDING                = 15
)

var obuTypeNames = scalar.UintMapSymStr{
	OBU_SEQUENCE_HEADER:        "sequence_header",
	OBU_TEMPORAL_DELIMITER:     "temporal_delimiter",
	OBU_FRAME_HEADER:           "frame_header",
	OBU_TILE_GROUP:             "tile_group",
	OBU_METADATA:               "metadata",
	OBU_FRAME:                  "frame",
	OBU_REDUNDANT_FRAME_HEADER: "redundant_frame_header",
	OBU_TILE_LIST:              "tile_list",
	OBU_PADDING:                "padding",
}

const (
	CP_BT_709       = 1  //	BT.709
	CP_UNSPECIFIED  = 2  //	Unspecified
	CP_BT_470_M     = 4  //	BT.470 System M (historical)
	CP_BT_470_B_G   = 5  //	BT.470 System B, G (historical)
	CP_BT_601       = 6  //	BT.601
	CP_SMPTE_240    = 7  //	SMPTE 240
	CP_GENERIC_FILM = 8  //	Generic film (color filters using illuminant C)
	CP_BT_2020      = 9  //	BT.2020, BT.2100
	CP_XYZ          = 10 //	SMPTE 428 (CIE 1921 XYZ)
	CP_SMPTE_431    = 11 //	SMPTE RP 431-2
	CP_SMPTE_432    = 12 //	SMPTE EG 432-,1
	CP_EBU_3213     = 22 //	EBU Tech. 3213-E
)

var cpTypeNames = scalar.UintMapSymStr{
	CP_BT_709:       "bt_709",
	CP_UNSPECIFIED:  "unspecified",
	CP_BT_470_M:     "bt_470_m",
	CP_BT_470_B_G:   "bt_470_b_g",
	CP_BT_601:       "bt_601",
	CP_SMPTE_240:    "smpte_240",
	CP_GENERIC_FILM: "generic_film",
	CP_BT_2020:      "bt_2020",
	CP_XYZ:          "xyz",
	CP_SMPTE_431:    "smpte_431",
	CP_SMPTE_432:    "smpte_432",
	CP_EBU_3213:     "ebu_3213",
}

const (
	TC_RESERVED_0     = 0  //		For future use
	TC_BT_709         = 1  //		BT.709
	TC_UNSPECIFIED    = 2  //		Unspecified
	TC_RESERVED_3     = 3  //		For future use
	TC_BT_470_M       = 4  //		BT.470 System M (historical)
	TC_BT_470_B_G     = 5  //		BT.470 System B, G (historical)
	TC_BT_601         = 6  //		BT.601
	TC_SMPTE_240      = 7  //		SMPTE 240 M
	TC_LINEAR         = 8  //		Linear
	TC_LOG_100        = 9  //		Logarithmic (100 : 1 range)
	TC_LOG_100_SQRT10 = 10 //		Logarithmic (100 * Sqrt(10) : 1 range)
	TC_IEC_61966      = 11 //		IEC 61966-2-4
	TC_BT_1361        = 12 //		BT.1361
	TC_SRGB           = 13 //		sRGB or sYCC
	TC_BT_2020_10_BIT = 14 //		BT.2020 10-bit systems
	TC_BT_2020_12_BIT = 15 //		BT.2020 12-bit systems
	TC_SMPTE_2084     = 16 //		SMPTE ST 2084, ITU BT.2100 PQ
	TC_SMPTE_428      = 17 //		SMPTE ST 428
	TC_HLG            = 18 //		BT.2100 HLG, ARIB STD-B67
)

var tcTypeNames = scalar.UintMapSymStr{
	TC_RESERVED_0:     "reserved_0",
	TC_BT_709:         "bt_709",
	TC_UNSPECIFIED:    "unspecified",
	TC_RESERVED_3:     "reserved_3",
	TC_BT_470_M:       "bt_470_m",
	TC_BT_470_B_G:     "bt_470_b_g",
	TC_BT_601:         "bt_601",
	TC_SMPTE_240:      "smpte_240",
	TC_LINEAR:         "linear",
	TC_LOG_100:        "log_100",
	TC_LOG_100_SQRT10: "log_100_sqrt10",
	TC_IEC_61966:      "iec_61966",
	TC_BT_1361:        "bt_1361",
	TC_SRGB:           "srgb",
	TC_BT_2020_10_BIT: "bt_2020_10_bit",
	TC_BT_2020_12_BIT: "bt_2020_12_bit",
	TC_SMPTE_2084:     "smpte_2084",
	TC_SMPTE_428:      "smpte_428",
	TC_HLG:            "hlg",
}

const (
	MC_IDENTITY    = 0  //		Identity matrix
	MC_BT_709      = 1  //		BT.709
	MC_UNSPECIFIED = 2  //		Unspecified
	MC_RESERVED_3  = 3  //		For future use
	MC_FCC         = 4  //		US FCC 73.628
	MC_BT_470_B_G  = 5  //		BT.470 System B, G (historical)
	MC_BT_601      = 6  //		BT.601
	MC_SMPTE_240   = 7  //		SMPTE 240 M
	MC_SMPTE_YCGCO = 8  //		YCgCo
	MC_BT_2020_NCL = 9  //		BT.2020 non-constant luminance, BT.2100 YCbCr
	MC_BT_2020_CL  = 10 //		BT.2020 constant luminance
	MC_SMPTE_2085  = 11 //		SMPTE ST 2085 YDzDx
	MC_CHROMAT_NCL = 12 //		Chromaticity-derived non-constant luminance
	MC_CHROMAT_CL  = 13 //		Chromaticity-derived constant luminance
	MC_ICTCP       = 14 //		BT.2100 ICtCp
)

var mcTypeNames = scalar.UintMapSymStr{
	MC_IDENTITY:    "identity",
	MC_BT_709:      "bt_709",
	MC_UNSPECIFIED: "unspecified",
	MC_RESERVED_3:  "reserved_3",
	MC_FCC:         "fcc",
	MC_BT_470_B_G:  "bt_470_b_g",
	MC_BT_601:      "bt_601",
	MC_SMPTE_240:   "smpte_240",
	MC_SMPTE_YCGCO: "smpte_ycgco",
	MC_BT_2020_NCL: "bt_2020_ncl",
	MC_BT_2020_CL:  "bt_2020_cl",
	MC_SMPTE_2085:  "smpte_2085",
	MC_CHROMAT_NCL: "chromat_ncl",
	MC_CHROMAT_CL:  "chromat_cl",
	MC_ICTCP:       "ictcp",
}

const (
	SEQ_PROFILE_MAIN         = 0
	SEQ_PROFILE_HIGH         = 1
	SEQ_PROFILE_PROFESSIONAL = 2
)

var seqProfileNames = scalar.UintMapSymStr{
	SEQ_PROFILE_MAIN:         "main",
	SEQ_PROFILE_HIGH:         "high",
	SEQ_PROFILE_PROFESSIONAL: "professional",
}

// from https://aomediacodec.github.io/av1-spec/#symbols-and-abbreviated-terms
const SELECT_SCREEN_CONTENT_TOOLS = 2

// TODO: ignore empty branch lint warnings
//
//nolint:staticcheck
func obuDecode(d *decode.D) any {
	var obuType uint64
	var obuSize int64
	hasExtension := false
	hasSizeField := false

	d.FieldStruct("header", func(d *decode.D) {
		d.FieldU1("forbidden_bit")
		obuType = d.FieldU4("type", obuTypeNames)
		hasExtension = d.FieldBool("extension_flag")
		hasSizeField = d.FieldBool("has_size_field")
		d.FieldU1("reserved_1bit")
		if hasExtension {
			d.FieldU3("temporal_id")
			d.FieldU2("spatial_id")
			d.FieldU3("extension_header_reserved_3bits")
		}
	})

	if hasSizeField {
		obuSize = int64(d.FieldULEB128("size"))
	} else {
		obuSize = d.BitsLeft() / 8
		if hasExtension {
			obuSize--
		}
	}

	d.FramedFn(obuSize*8, func(d *decode.D) {
		// TODO: this only handles the OBU_SEQUENCE_HEADER case for now
		// fro spec https://aomediacodec.github.io/av1-spec/#general-obu-syntax
		// if ( obu_type != OBU_SEQUENCE_HEADER &&
		// 	obu_type != OBU_TEMPORAL_DELIMITER &&
		// 	OperatingPointIdc != 0 &&
		// 	obu_extension_flag == 1 )

		switch obuType {
		case OBU_SEQUENCE_HEADER:
			seqProfile := d.FieldU3("seq_profile", seqProfileNames)
			d.FieldU1("still_picture")
			reducedStillPictureHeader := d.FieldU1("reduced_still_picture_header")
			if reducedStillPictureHeader == 1 {
				d.FieldU5("seq_level_idx0") // TODO: array as below?
			} else {
				timingInfoPresentFlag := d.FieldU1("timing_info_present_flag")
				if timingInfoPresentFlag == 1 {
					// TODO:
					return
				}
				initialDisplayDelayPresentFlag := d.FieldU1("initial_display_delay_present_flag")
				operatingPointsCntMinus1 := d.FieldU5("operating_points_cnt_minus_1")

				d.FieldArray("operating_points", func(d *decode.D) {
					for i := uint64(0); i <= operatingPointsCntMinus1; i++ {
						d.FieldStruct("operating_point", func(d *decode.D) {
							d.FieldU12("operating_point_idc")
							seqLevelIdx := d.FieldU5("seq_level_idx")

							if seqLevelIdx > 7 {
								d.FieldU1("seq_tier")
							} else {
								// nop
							}
						})
						if initialDisplayDelayPresentFlag == 1 {
							initialDisplayDelayPresentForThisOp := d.FieldU1("seq_tier")
							if initialDisplayDelayPresentForThisOp == 1 {
								d.FieldU4("initial_display_delay_minus_1")
							}
						}
					}
				})

			}

			frameWidthBitsMinus1 := d.FieldU4("frame_width_bits_minus_1")
			frameHeightBitsMinus1 := d.FieldU4("frame_height_bits_minus_1")
			frameWidthMinus1 := d.FieldU("max_frame_width_minus_1", int(frameWidthBitsMinus1)+1)
			frameHeightMinus1 := d.FieldU("max_frame_height_minus_1", int(frameHeightBitsMinus1)+1)
			d.FieldValueUint("frame_width", frameWidthMinus1+1)
			d.FieldValueUint("frame_height", frameHeightMinus1+1)

			var frameIdNumbersPresentFlag uint64 = 0
			if reducedStillPictureHeader == 1 {
				// nop
			} else {
				frameIdNumbersPresentFlag = d.FieldU1("frame_id_numbers_present_flag")
			}
			if frameIdNumbersPresentFlag == 1 {
				d.FieldU4("delta_frame_id_length_minus_2")
				d.FieldU3("additional_frame_id_length_minus_1")
			}
			d.FieldU1("use_128x128_superblock")
			d.FieldU1("enable_filter_intra")
			d.FieldU1("enable_intra_edge_filter")
			if reducedStillPictureHeader == 1 {
				//nop
			} else {
				d.FieldU1("enable_interintra_compound")
				d.FieldU1("enable_masked_compound")
				d.FieldU1("enable_warped_motion")
				d.FieldU1("enable_dual_filter")
				enableOrderHint := d.FieldU1("enable_order_hint")
				if enableOrderHint == 1 {
					d.FieldU1("enable_jnt_comp")
					d.FieldU1("enable_ref_frame_mvs")
				}
				seqChooseScreenContentTools := d.FieldU1("seq_choose_screen_content_tools")
				var seqForceScreenContentTools uint64 = SELECT_SCREEN_CONTENT_TOOLS
				if seqChooseScreenContentTools == 1 {
					// nop
				} else {
					seqForceScreenContentTools = d.FieldU1("seq_force_screen_content_tools")
				}
				if seqForceScreenContentTools > 0 {
					seqChooseIntegerMv := d.FieldU1("seq_choose_integer_mv")
					if seqChooseIntegerMv == 1 {
						// nop
					} else {
						d.FieldU1("seq_force_integer_mv")
					}
				}
				if enableOrderHint == 1 {
					d.FieldU3("order_hint_bits_minus_1")
				}
			}
			d.FieldU1("enable_superres")
			d.FieldU1("enable_cdef")
			d.FieldU1("enable_restoration")
			d.FieldStruct("color_config", func(d *decode.D) {
				// https://aomediacodec.github.io/av1-spec/#color-config-syntax
				highBitdepth := d.FieldU1("high_bitdepth")
				var twelveBit uint64
				var bitDepth uint64 = 0 // TODO: what if seqProfile > 2?
				if seqProfile == 2 && highBitdepth == 1 {
					twelveBit = d.FieldU1("twelve_bit")
					if twelveBit == 1 {
						bitDepth = 12
					} else {
						bitDepth = 10
					}
				} else if seqProfile <= 2 {
					if highBitdepth == 1 {
						bitDepth = 10
					} else {
						bitDepth = 8
					}
				}
				d.FieldValueUint("bit_depth", bitDepth)

				var monoChrome uint64
				if seqProfile == 1 {
					d.FieldValueUint("mono_chrome", 0)
					monoChrome = 0
				} else {
					monoChrome = d.FieldU1("mono_chrome")
				}
				colorDescriptionPresentFlag := d.FieldU1("color_description_present_flag")

				var colorPrimaries uint64 = CP_UNSPECIFIED
				var transferCharacteristics uint64 = TC_UNSPECIFIED
				var matrixCoefficients uint64 = MC_UNSPECIFIED

				if colorDescriptionPresentFlag == 1 {
					colorPrimaries = d.FieldU8("color_primaries", cpTypeNames)
					transferCharacteristics = d.FieldU8("transfer_characteristics", tcTypeNames)
					matrixCoefficients = d.FieldU8("matrix_coefficients", mcTypeNames)
				} else {
					d.FieldValueUint("color_primaries", transferCharacteristics, cpTypeNames)
					d.FieldValueUint("transfer_characteristics", transferCharacteristics, tcTypeNames)
					d.FieldValueUint("matrix_coefficients", matrixCoefficients, mcTypeNames)
				}
				if monoChrome == 1 {
					d.FieldU1("color_range")
					d.FieldValueUint("subsampling_x", 1)
					d.FieldValueUint("subsampling_y", 1)
				} else if colorPrimaries == CP_BT_709 &&
					transferCharacteristics == TC_SRGB &&
					matrixCoefficients == MC_IDENTITY {
					d.FieldValueUint("color_range", 1)
					d.FieldValueUint("subsampling_x", 0)
					d.FieldValueUint("subsampling_y", 0)
					// nop
				} else {
					d.FieldU1("color_range")
					var subsamplingX uint64 = 0
					var subsamplingY uint64 = 0
					if seqProfile == 0 {
						subsamplingX = 1
						subsamplingY = 1
						d.FieldValueUint("subsampling_x", subsamplingX)
						d.FieldValueUint("subsampling_y", subsamplingY)
					} else if seqProfile == 1 {
						d.FieldValueUint("subsampling_x", subsamplingX)
						d.FieldValueUint("subsampling_y", subsamplingY)
					} else {
						if bitDepth == 12 {
							subsamplingX = d.FieldU1("subsampling_x")
							if subsamplingX == 1 {
								subsamplingY = d.FieldU1("subsampling_y")
							} else {
								d.FieldValueUint("subsampling_y", subsamplingY)
							}
						} else {
							subsamplingX = 1
							d.FieldValueUint("subsampling_x", subsamplingX)
							d.FieldValueUint("subsampling_y", subsamplingY)
						}
					}
					if subsamplingX == 1 && subsamplingY == 1 {
						d.FieldU2("chroma_sample_position")
					}
				}
				d.FieldU1("separate_uv_delta_q")
			})
			d.FieldU1("film_grain_params_present")
		default:
		}

		d.FieldRawLen("data", d.BitsLeft())
	})

	return nil
}
