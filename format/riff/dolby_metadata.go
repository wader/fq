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
				if segmentID == 0 {
					seenEnd = true
					return
				}

				segmentSize := d.FieldU16("size")

				switch segmentID {
				case metadataSegmentTypeDolbyEMetadata:
					parseDolbyE(d)
				case metadataSegmentTypeDolbyEDigitaletadata:
					parseDolbyDigital(d)
				case metadataSegmentTypeDolbyDigitalPlusMetadata:
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
	0: {Sym: "normal"},
	1: {Sym: "warping"},
	2: {Sym: "downmix_dolby_pro_logic_iix"},
	3: {Sym: "downmix_loro"},
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
	metadataSegmentTypeEnd                      = 0
	metadataSegmentTypeDolbyEMetadata           = 1
	metadataSegmentTypeDolbyReserved2           = 2
	metadataSegmentTypeDolbyEDigitaletadata     = 3
	metadataSegmentTypeDolbyReserved4           = 4
	metadataSegmentTypeDolbyReserved5           = 5
	metadataSegmentTypeDolbyReserved6           = 6
	metadataSegmentTypeDolbyDigitalPlusMetadata = 7
	metadataSegmentTypeAudioInfo                = 8
	metadataSegmentTypeDolbyAtmos               = 9
	metadataSegmentTypeDolbyAtmosSupplemental   = 10
)

var metadataSegmentTypeMap = scalar.UintMapSymStr{
	metadataSegmentTypeEnd:                      "end",
	metadataSegmentTypeDolbyEMetadata:           "dolby_e_metadata",
	metadataSegmentTypeDolbyReserved2:           "reserved2",
	metadataSegmentTypeDolbyEDigitaletadata:     "dolby_e_digitale_tadata",
	metadataSegmentTypeDolbyReserved4:           "reserved4",
	metadataSegmentTypeDolbyReserved5:           "reserved5",
	metadataSegmentTypeDolbyReserved6:           "reserved6",
	metadataSegmentTypeDolbyDigitalPlusMetadata: "dolby_digital_plus_metadata",
	metadataSegmentTypeAudioInfo:                "audio_info",
	metadataSegmentTypeDolbyAtmos:               "dolby_atmos",
	metadataSegmentTypeDolbyAtmosSupplemental:   "dolby_atmos_supplemental",
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
	// TODO: both these are fixed size null terminated strings?
	d.FieldUTF8NullFixedLen("atmos_dbmd_content_creation_preamble", 32)
	d.FieldUTF8NullFixedLen("atmos_dbmd_content_creation_tool", 64)
	d.FieldStruct("version", func(d *decode.D) {
		d.FieldU8("major")
		d.FieldU8("minor")
		d.FieldU8("micro")
	})

	// TODO: what is this?
	d.FieldRawLen("unknown0", 53*8)

	d.FieldU8("warp_mode", warpModeMap)

	// TODO: what is this?
	d.FieldRawLen("unknown1", 15*8)
	d.FieldRawLen("unknown2", 80*8)
}

func parseDolbyAtmosSupplemental(d *decode.D) {
	d.FieldU32("dasms_sync", d.UintAssert(0xf8726fbd), scalar.UintHex)

	// TODO: wav.go sets LE default i think?
	objectCount := int64(d.FieldU16("object_count"))
	d.FieldU8("reserved")

	i := 0
	d.FieldStructNArray("trim_configs", "trim_config", 9, func(d *decode.D) {
		d.FieldRawLen("reserved", 7)
		d.FieldU1("type", scalar.UintMapSymStr{
			0: "manual",
			1: "automatic",
		})
		d.FieldValueStr("config_name", trimConfigName[uint64(i)])

		// TODO: this is null separted list of def strings?
		d.FieldUTF8("raw", 14)
		// str := d.UTF8(14)
		// bytes := []byte(str)
		// var nonZeroBytes []string
		// for _, b := range bytes {
		// 	if b != 0 {
		// 		nonZeroBytes = append(nonZeroBytes, fmt.Sprintf("%d", b))
		// 	}
		// }
		// TODO(jmarnell): I think the +3dB trim settings are here.
		//		Would like this at least as an array of numbers, instead of this CSV string
		// d.FieldValueStr("trim_defs", strings.Join(nonZeroBytes, ", "))

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
