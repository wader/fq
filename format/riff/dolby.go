package riff

// Dolby Metadata, e.g. Atmos, AC3, Dolby Digital [Plus]
// https://tech.ebu.ch/files/live/sites/tech/files/shared/tech/tech3285s6.pdf
// https://github.com/DolbyLaboratories/dbmd-atmos-parser

import (
	"fmt"
	"strings"

	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

func tmp_dbmdDecode(d *decode.D, size int64) any {
	version := d.U32()
	major := (version >> 24) & 0xFF
	minor := (version >> 16) & 0xFF
	patch := (version >> 8) & 0xFF
	build := version & 0xFF
	d.FieldValueStr("version", fmt.Sprintf("%d.%d.%d.%d", major, minor, patch, build))

	d.FieldArray("metadata_segments", func(d *decode.D) {
		for {
			d.FieldStruct("metadata_segment", func(d *decode.D) {
				segmentID := d.FieldU8("metadata_segment_id")

				// TODO(jmarnell): I think I need a loop until, but not creating these empty segments
				// spec says we're done with 0 ID, so I'd like to not make the empty segment(s)
				if segmentID == 0 {
					if d.BitsLeft() > 0 {
						d.SeekRel(d.BitsLeft() * 8)
					}
					return
				}

				segmentSize := d.FieldU16("metadata_segment_size")
				bitsLeft := d.BitsLeft()

				switch segmentID {
				case 1:
					parseDolbyE(d)
				case 3:
					parseDolbyDigital(d)
				case 7:
					parseDolbyDigitalPlus(d)
				case 8:
					parseAudioInfo(d)
				case 9:
					parseDolbyAtmos(d)
				case 10:
					parseDolbyAtmosSupplemental(d)
				default:
					d.FieldRawLen("unknown_segment_raw", int64(segmentSize*8))
				}

				bytesRemaining := (bitsLeft-d.BitsLeft())/8 - int64(segmentSize)
				if bytesRemaining < 0 {
					d.Fatalf("Read too many bytes for segment %d, read %d over, expected %d", segmentID, -bytesRemaining, segmentSize)
				} else if bytesRemaining > 0 {
					d.FieldValueUint("SKIPPED_BYTES", uint64(bytesRemaining))
					d.SeekRel((int64(segmentSize) - bytesRemaining) * 8)
				}

				d.FieldU8("metadata_segment_checksum")
			})
		}
	})

	return nil
}

var compressionDesc = scalar.UintMapDescription{
	0: "none",
	1: "Film, Standard",
	2: "Film, Light",
	3: "Music, Standard",
	4: "Music, Light",
	5: "Speech",
	// TODO(jmarnell): Can I handle rest is "Reserved"?
}

// TODO(jmarnell): Better way to handle "Reserved"?
func mapWithReserved(m map[uint64]string, key uint64) string {
	if val, ok := m[key]; ok {
		return val
	}
	return "Reserved"
}

var bitstreamMode = scalar.UintMapDescription{
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

var binaural = scalar.UintMapDescription{
	0: "bypass",
	1: "near",
	2: "far",
	3: "mid",
	4: "not indicated",
}

var warpMode = scalar.UintMapDescription{
	0: "normal",
	1: "warping",
	2: "downmix Dolby Pro Logic IIx",
	3: "downmix LoRo",
	4: "not indicated (Default warping will be applied.)",
}

var tmp_trimConfigName = scalar.UintMapDescription{
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

var trimType = scalar.UintMapDescription{
	0: "manual",
	1: "automatic",
}

func tmp_parseDolbyE(d *decode.D) {
	d.FieldValueStr("metadata_segment_type", "dolby_e")

	d.FieldU8("program_config")
	d.FieldU8("frame_rate_code")
	d.FieldRawLen("e_SMPTE_time_code", 8*8)
	d.FieldRawLen("e_reserved", 1*8)
	d.FieldRawLen("e_reserved2", 25*8)
	d.FieldRawLen("reserved_for_future_use", 80*8)
}

func tmp_parseDolbyDigital(d *decode.D) {
	d.FieldValueStr("metadata_segment_type", "dolby_digital")

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

func tmp_parseDolbyDigitalPlus(d *decode.D) {
	d.FieldValueStr("metadata_segment_type", "dolby_digital_plus")

	d.FieldU8("program_id")
	programInfo := d.FieldU8("program_info")
	lfeon := programInfo & 0b1_000_000
	bsmod := programInfo & 0b0_111_000
	acmod := programInfo & 0b0_000_111
	d.FieldValueBool("lfe_on", lfeon != 0)
	if bsmod == 0b111 && acmod != 0b001 {
		bsmod = 0b1000
	}
	d.FieldValueStr("bitstream_mode", bitstreamMode[bsmod])

	d.FieldU16LE("ddplus_reserved_a")

	d.FieldU8("surround_config")
	d.FieldU8("dialnorm_info")
	d.FieldU8("langcod")
	d.FieldU8("audio_prod_info")
	d.FieldU8("ext_bsi1_word1")
	d.FieldU8("ext_bsi1_word2")
	d.FieldU8("ext_bsi2_word1")

	d.FieldU24LE("ddplus_reserved_b")

	d.FieldValueStr("compr1_type", mapWithReserved(compressionDesc, d.FieldU8("compr1")))
	d.FieldValueStr("dynrng1_type", mapWithReserved(compressionDesc, d.FieldU8("dynrng1")))

	d.FieldU24LE("ddplus_reserved_c")

	d.FieldU8("ddplus_info1")

	d.FieldU40LE("ddplus_reserved_d")

	d.FieldU16LE("datarate")
	d.FieldRawLen("reserved_for_future_use", 69*8)
}

func tmp_parseAudioInfo(d *decode.D) {
	d.FieldValueStr("metadata_segment_type", "audio_info")

	d.FieldU8("program_id")
	d.FieldUTF8("audio_origin", 32)
	d.FieldU32LE("largest_sample_value")
	d.FieldU32LE("largest_sample_value_2")
	d.FieldU32LE("largest_true_peak_value")
	d.FieldU32LE("largest_true_peak_value_2")
	d.FieldU32LE("dialogue_loudness")
	d.FieldU32LE("dialogue_loudness_2")
	d.FieldU32LE("speech_content")
	d.FieldU32LE("speech_content_2")
	d.FieldUTF8("last_processed_by", 32)
	d.FieldUTF8("last_operation", 32)
	d.FieldUTF8("segment_creation_date", 32)
	d.FieldUTF8("segment_modified_date", 32)
}

func tmp_parseDolbyAtmos(d *decode.D, size uint64) {
	d.FieldValueStr("metadata_segment_type", "dolby_atmos")

	// d.SeekRel(32 * 8)
	str := d.FieldUTF8Null("atmos_dbmd_content_creation_preamble")
	d.SeekRel(int64(32-len(str)-1) * 8)

	str = d.FieldUTF8Null("atmos_dbmd_content_creation_tool")
	d.SeekRel(int64(64-len(str)-1) * 8)

	major := d.U8()
	minor := d.U8()
	micro := d.U8()
	d.FieldValueStr("version", fmt.Sprintf("%d.%d.%d", major, minor, micro))
	d.SeekRel(53 * 8)

	warpBedReserved := d.U8()
	d.FieldValueUint("warp_mode", warpBedReserved&0x7)
	d.FieldValueStr("warp_mode_type", warpMode[warpBedReserved&0x7])

	d.SeekRel(15 * 8)
	d.SeekRel(80 * 8)
}

func tmp_parseDolbyAtmosSupplemental(d *decode.D, size uint64) {
	d.FieldValueStr("metadata_segment_type", "dolby_atmos_supplemental")

	sync := d.FieldU32LE("dasms_sync")
	d.FieldValueBool("dasms_sync_valid", sync == 0xf8726fbd)

	objectCount := int64(d.FieldU16LE("object_count"))
	d.FieldU8LE("reserved")

	i := 0
	d.FieldStructNArray("trim_configs", "trim_config", 9, func(d *decode.D) {
		autoTrimReserved := d.FieldU8LE("auto_trim_reserved")
		autoTrim := autoTrimReserved & 0x01
		d.FieldValueBool("auto_trim", autoTrim == 1)
		d.FieldValueStr("trim_type", trimType[autoTrim])
		d.FieldValueStr("trim_config_name", trimConfigName[uint64(i)])

		//d.SeekRel(14 * 8)
		// d.FieldUTF8("raw", 14)
		str := d.UTF8(14)
		bytes := []byte(str)
		var nonZeroBytes []string
		for _, b := range bytes {
			if b != 0 {
				nonZeroBytes = append(nonZeroBytes, fmt.Sprintf("%d", b))
			}
		}
		// TODO(jmarnell): I think the +3dB trim settings are here.
		//		Would like this at least as an array of numbers, instead of this CSV string
		d.FieldValueStr("trim_defs", strings.Join(nonZeroBytes, ", "))

		i++
	})

	d.FieldStructNArray("objects", "object", objectCount, func(d *decode.D) {
		d.FieldU8LE("object_value")
	})

	d.FieldStructNArray("binaural_render_modes", "binaural_render_mode", objectCount, func(d *decode.D) {
		mode := d.U8LE() & 0x7
		d.FieldValueUint("render_mode", mode)
		d.FieldValueStr("render_mode_type", binaural[mode])
	})
}
