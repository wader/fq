package riff

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
)

func dbmdDecode(d *decode.D) any {
	// TODO(jmarnell): Read as string (and assert?)
	// d.FieldUTF8("chunk_id", 4, d.StrAssert("dbmd"))
	d.FieldU32("chunk_id")

	d.FieldU32("chunk_size")
	// TODO(jmarnell): Should be string formatted to: "'1.15.2.0' corresponds to 0x010F0200"
	d.FieldU32("version")

	d.FieldArray("metadata_segments", func(d *decode.D) {
		for {
			d.FieldStruct("metadata_segment", func(d *decode.D) {
				segmentID := d.FieldU8("metadata_segment_id")
				fmt.Println("segmentID: ", segmentID)
				if segmentID == 0 {
					return
				}

				switch segmentID {
				case 1:
					parseDolbyE(d)
				case 3:
					parseDolbyDigital(d)
				case 7:
					parseDolbyDigitalPlus(d)
				case 8:
					parseAudioInfo(d)
				default:
					d.FieldRawLen("unknown_segment_raw", int64(d.BitsLeft()))
				}
			})
		}
	})

	return nil
}

func parseDolbyE(d *decode.D) {
	d.FieldU8("program_config")
	d.FieldU8("frame_rate_code")
	d.FieldRawLen("e_SMPTE_time_code", 8*8)
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
	d.FieldU8("program_info")
	d.FieldRawLen("ddplus_reserved1", 2*8)
	d.FieldU8("surround_config")
	d.FieldU8("dialnorm_info")
	d.FieldU8("langcod")
	d.FieldU8("audio_prod_info")
	d.FieldU8("ext_bsi1_word1")
	d.FieldU8("ext_bsi1_word2")
	d.FieldU8("ext_bsi2_word1")
	d.FieldRawLen("ddplus_reserved2", 1*8)
	d.FieldRawLen("ddplus_reserved3", 1*8)
	d.FieldRawLen("ddplus_reserved4", 1*8)
	d.FieldU8("compr1")
	d.FieldU8("dynrng1")
	d.FieldRawLen("ddplus_reserved5", 3*8)
	d.FieldU8("ddplus_info1")
	d.FieldRawLen("ddplus_reserved6", 5*8)
	d.FieldU16LE("datarate")
	d.FieldRawLen("reserved_for_future_use", 69*8)
}

func parseAudioInfo(d *decode.D) {
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
