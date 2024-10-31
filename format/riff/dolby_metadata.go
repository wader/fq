package riff

// Dolby Metadata, e.g. Atmos, AC3, Dolby Digital [Plus]
// https://tech.ebu.ch/files/live/sites/tech/files/shared/tech/tech3285s6.pdf
// https://github.com/DolbyLaboratories/dbmd-atmos-parser

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed dolby_metadata.md
var dolbyMetadataFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Dolby_Metadata,
		&decode.Format{
			Description: "Dolby Metadata (Atmos, AC3, Dolby Digital)",
			DecodeFn:    dbmdDecode,
		},
	)
	interp.RegisterFS(dolbyMetadataFS)
}

func dbmdDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	d.FieldStruct("version", func(d *decode.D) {
		d.FieldU8("major")
		d.FieldU8("minor")
		d.FieldU8("patch")
		d.FieldU8("build")
	})

	d.FieldArray("metadata_segments", func(d *decode.D) {
		seenEnd := false
		for !seenEnd {
			d.FieldStruct("metadata_segment", func(d *decode.D) {
				segmentID := d.FieldU8("id", metadataSegmentTypeMap)

				// TODO(jmarnell): This will always make an empty end segment, I think it would be better to omit it
				if segmentID == metadataSegmentTypeEnd {
					seenEnd = true
					return
				}

				segmentSize := d.FieldU16("size")

				switch segmentID {
				case metadataSegmentTypeDolbyE:
					parseDolbyE(d)
				case metadataSegmentTypeDolbyDigital:
					parseDolbyDigital(d)
				case metadataSegmentTypeDolbyDigitalPlus:
					parseDolbyDigitalPlus(d)
				case metadataSegmentTypeAudioInfo:
					parseAudioInfo(d)
				case metadataSegmentTypeDolbyAtmos:
					parseDolbyAtmos(d)
				case metadataSegmentTypeDolbyAtmosSupplemental:
					parseDolbyAtmosSupplemental(d)
				default:
					d.FieldRawLen("unknown", int64(segmentSize*8))
				}

				// TODO: use this to validate parsing
				d.FieldU8("checksum", scalar.UintHex)
			})
		}
	})

	return nil
}

var compressionDescMap = scalar.UintMapSymStr{
	0: "none",
	1: "film_standard",
	2: "film_light",
	3: "music_standard",
	4: "music_light",
	5: "speech",
}

var downmix5to2DescMap = scalar.UintMap{
	0: {Sym: "not_indicated", Description: "Not indicated (Lo/Ro)"},
	1: {Sym: "loro", Description: "Lo/Ro"},
	2: {Sym: "ltrt_dpl", Description: "Lt/Rt (Dolby Pro Logic)"},
	3: {Sym: "ltrt_dpl2", Description: "Lt/Rt (Dolby Pro Logic II)"},
	4: {Sym: "direct_stereo_render", Description: "Direct stereo render"},
}

var phaseShift5to2DescMap = scalar.UintMap{
	0: {Sym: "no_shift", Description: "Without Phase 90"},
	1: {Sym: "shift_90", Description: "With Phase 90"},
}

var bitstreamModeMap = scalar.UintMapDescription{
	0b000: "main audio service: complete main (CM)",
	0b001: "main audio service: music and effects (ME)",
	0b010: "associated service: visually impaired (VI)",
	0b011: "associated service: hearing impaired (HI)",
	0b100: "associated service: dialogue (D)",
	0b101: "associated service: commentary (C)",
	0b110: "associated service: emergency (E)",
	0b111: "associated service: voice over (VO)",

	0b1000: "associated service: karaoke (K)",
}

var binauralRenderModeMap = scalar.UintMapSymStr{
	0: "bypass",
	1: "near",
	2: "far",
	3: "mid",
	4: "not_indicated",
}

var warpModeMap = scalar.UintMap{
	0: {Sym: "normal", Description: "possibly: Direct render"},
	1: {Sym: "warping", Description: "possibly: Direct render with room balance"},
	2: {Sym: "downmix_dolby_pro_logic_iix", Description: "Dolby Pro Logic IIx"},
	3: {Sym: "downmix_loro", Description: "possibly: Standard (Lo/Ro)"},
	4: {Sym: "not_indicated", Description: "Default warping will be applied"},
}

var trimConfigName = scalar.UintMapDescription{
	0: "2.0",
	1: "5.1",
	2: "7.1",
	3: "2.1.2",
	4: "5.1.2",
	5: "7.1.2",
	6: "2.1.4",
	7: "5.1.4",
	8: "7.1.4",
}

const (
	metadataSegmentTypeEnd                    = 0
	metadataSegmentTypeDolbyE                 = 1
	metadataSegmentTypeDolbyReserved2         = 2
	metadataSegmentTypeDolbyDigital           = 3
	metadataSegmentTypeDolbyReserved4         = 4
	metadataSegmentTypeDolbyReserved5         = 5
	metadataSegmentTypeDolbyReserved6         = 6
	metadataSegmentTypeDolbyDigitalPlus       = 7
	metadataSegmentTypeAudioInfo              = 8
	metadataSegmentTypeDolbyAtmos             = 9
	metadataSegmentTypeDolbyAtmosSupplemental = 10
)

var metadataSegmentTypeMap = scalar.UintMapSymStr{
	metadataSegmentTypeEnd:                    "end",
	metadataSegmentTypeDolbyE:                 "dolby_e_metadata",
	metadataSegmentTypeDolbyReserved2:         "reserved2",
	metadataSegmentTypeDolbyDigital:           "dolby_digital_metadata",
	metadataSegmentTypeDolbyReserved4:         "reserved4",
	metadataSegmentTypeDolbyReserved5:         "reserved5",
	metadataSegmentTypeDolbyReserved6:         "reserved6",
	metadataSegmentTypeDolbyDigitalPlus:       "dolby_digital_plus_metadata",
	metadataSegmentTypeAudioInfo:              "audio_info",
	metadataSegmentTypeDolbyAtmos:             "dolby_atmos",
	metadataSegmentTypeDolbyAtmosSupplemental: "dolby_atmos_supplemental",
}

func parseDolbyE(d *decode.D) {
	d.FieldU8("program_config")
	d.FieldU8("frame_rate_code")
	d.FieldRawLen("e_smpte_time_code", 8*8)
	d.FieldRawLen("e_reserved", 1*8)
	d.FieldRawLen("e_reserved2", 25*8)
	d.FieldRawLen("reserved_for_future_use", 80*8)
}

func parseDolbyDigital(d *decode.D) {
	d.FieldU8("ac3_program_id")
	d.FieldU8("program_info")
	d.FieldU8("datarate_info")
	d.FieldRawLen("reserved", 1*8)
	d.FieldU8("surround_config")
	d.FieldU8("dialnorm_info")
	d.FieldU8("ac3_langcod")
	d.FieldU8("audio_prod_info")
	d.FieldU8("ext_bsi1_word1")
	d.FieldU8("ext_bsi1_word2")
	d.FieldU8("ext_bsi2_word1")
	d.FieldRawLen("reserved2", 3*8)
	d.FieldU8("ac3_compr1")
	d.FieldU8("ac3_dynrng1")
	d.FieldRawLen("reserved_for_future_use", 21*8)
	d.FieldRawLen("program_description_text", 32*8)
}

func parseDolbyDigitalPlus(d *decode.D) {
	d.FieldU8("program_id")
	// TODO: make struct and read U1(?) U1 (lfeon) U3 (bsmod) U3(acmod) fields?
	programInfo := d.FieldU8("program_info")
	lfeon := programInfo & 0b1_000_000
	bsmod := programInfo & 0b0_111_000
	acmod := programInfo & 0b0_000_111
	d.FieldValueBool("lfe_on", lfeon != 0)
	if bsmod == 0b111 && acmod != 0b001 {
		bsmod = 0b1000
	}
	d.FieldValueStr("bitstream_mode", bitstreamModeMap[bsmod])

	d.FieldU16("ddplus_reserved_a")

	d.FieldU8("surround_config")
	d.FieldU8("dialnorm_info")
	d.FieldU8("langcod")
	d.FieldU8("audio_prod_info")
	d.FieldU8("ext_bsi1_word1")
	d.FieldU8("ext_bsi1_word2")
	d.FieldU8("ext_bsi2_word1")

	d.FieldU24("ddplus_reserved_b")

	d.FieldU8("compr1", scalar.UintSym("reserved"), compressionDescMap)
	d.FieldU8("dynrng1", scalar.UintSym("reserved"), compressionDescMap)

	d.FieldU24("ddplus_reserved_c")

	d.FieldU8("ddplus_info1")

	d.FieldU40("ddplus_reserved_d")

	d.FieldU16("datarate")
	d.FieldRawLen("reserved_for_future_use", 69*8)
}

func parseAudioInfo(d *decode.D) {
	d.FieldU8("program_id")
	d.FieldUTF8("audio_origin", 32)
	d.FieldU32("largest_sample_value")
	d.FieldU32("largest_sample_value_2")
	d.FieldU32("largest_true_peak_value")
	d.FieldU32("largest_true_peak_value_2")
	d.FieldU32("dialogue_loudness")
	d.FieldU32("dialogue_loudness_2")
	d.FieldU32("speech_content")
	d.FieldU32("speech_content_2")
	d.FieldUTF8("last_processed_by", 32)
	d.FieldUTF8("last_operation", 32)
	d.FieldUTF8("segment_creation_date", 32)
	d.FieldUTF8("segment_modified_date", 32)
}

func parseDolbyAtmos(d *decode.D) {
	d.FieldUTF8NullFixedLen("atmos_dbmd_content_creation_preamble", 32)
	d.FieldUTF8NullFixedLen("atmos_dbmd_content_creation_tool", 64)
	d.FieldStruct("version", func(d *decode.D) {
		d.FieldU8("major")
		d.FieldU8("minor")
		d.FieldU8("patch")
	})
	// TODO: All these unknowns? (mostly from MediaInfoLib, also Dolby repo)

	d.FieldRawLen("unknown0", 21*8)

	d.FieldRawLen("unknown1", 1)
	d.FieldU3("downmix_5to2", scalar.UintSym("unknown"), downmix5to2DescMap)
	d.FieldRawLen("unknown2", 2)
	d.FieldU2("phaseshift_90deg_5to2", scalar.UintSym("unknown"), phaseShift5to2DescMap)

	d.FieldRawLen("unknown3", 12*8)

	d.FieldRawLen("bed_distribution", 2)
	d.FieldRawLen("reserved0", 3)
	d.FieldU3("warp_mode", warpModeMap)

	d.FieldRawLen("unknown4", 15*8)
	d.FieldRawLen("unknown5", 80*8)
}

func parseDolbyAtmosSupplemental(d *decode.D) {
	d.FieldU32("dasms_sync", d.UintAssert(0xf8726fbd), scalar.UintHex)

	objectCount := int64(d.FieldU16("object_count"))
	d.FieldU8("reserved")

	i := 0
	d.FieldStructNArray("trim_configs", "trim_config", 9, func(d *decode.D) {
		d.FieldRawLen("reserved0", 7)
		trimType := d.FieldU1("type", scalar.UintMapSymStr{
			0: "manual",
			1: "automatic",
		})
		d.FieldValueStr("config_name", trimConfigName[uint64(i)])

		if trimType == 1 {
			d.FieldUTF8("reserved1", 14)
		} else {
			// TODO: Reference MediaInfo's logic and Dolby pdf's
			d.FieldUTF8("manual_trim_raw_config", 14)
		}
		i++
	})

	d.FieldArray("objects", func(d *decode.D) {
		for i := int64(0); i < objectCount; i++ {
			d.FieldU8("object_value")
		}
	})

	d.FieldArray("binaural_render_modes", func(d *decode.D) {
		// TODO: 0x7 mask needed?
		for i := int64(0); i < objectCount; i++ {
			d.FieldU8("render_mode", scalar.UintActualFn(func(a uint64) uint64 { return a & 0x7 }), binauralRenderModeMap)
		}
	})
}
