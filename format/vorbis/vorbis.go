package vorbis

// https://xiph.org/vorbis/doc/Vorbis_I_spec.html
// TODO: setup? more audio?
// TODO: end padding? byte align?

import (
	"fmt"
	"fq/format"
	"fq/pkg/decode"
)

var vorbisComment []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.VORBIS_PACKET,
		Description: "Vorbis packet",
		DecodeFn:    vorbisDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.VORBIS_COMMENT}, Formats: &vorbisComment},
		},
	})
}

const (
	packetTypeAudio          = 0
	packetTypeIdentification = 1
	packetTypeComment        = 3
	packetTypeSetup          = 5
)

var packetTypeNames = map[uint]string{
	packetTypeAudio:          "Audio",
	packetTypeIdentification: "Identification",
	packetTypeComment:        "Comment",
	packetTypeSetup:          "Setup",
}

func vorbisDecode(d *decode.D, in interface{}) interface{} {
	packetType := d.FieldUFn("packet_type", func() (uint64, decode.DisplayFormat, string) {
		packetTypeName := "unknown"
		t := d.U8()
		// 4.2.1. Common header decode
		// "these types are all odd as a packet with a leading single bit of ’0’ is an audio packet"
		if t&1 == 0 {
			t = packetTypeAudio
		}
		if n, ok := packetTypeNames[uint(t)]; ok {
			packetTypeName = n
		}
		return t, decode.NumberDecimal, packetTypeName
	})

	switch packetType {
	case packetTypeIdentification, packetTypeSetup, packetTypeComment:
		d.FieldValidateUTF8("magic", "vorbis")
	case packetTypeAudio:
	default:
		d.Invalid(fmt.Sprintf("unknown packet type %d", packetType))
	}

	switch packetType {
	case packetTypeAudio:
	case packetTypeIdentification:
		// 1   1) [vorbis_version] = read 32 bits as unsigned integer
		// 2   2) [audio_channels] = read 8 bit integer as unsigned
		// 3   3) [audio_sample_rate] = read 32 bits as unsigned integer
		// 4   4) [bitrate_maximum] = read 32 bits as signed integer
		// 5   5) [bitrate_nominal] = read 32 bits as signed integer
		// 6   6) [bitrate_minimum] = read 32 bits as signed integer
		// 7   7) [blocksize_0] = 2 exponent (read 4 bits as unsigned integer)
		// 8   8) [blocksize_1] = 2 exponent (read 4 bits as unsigned integer)
		// 9   9) [framing_flag] = read one bit
		d.FieldValidateUFn("vorbis_version", 0, d.U32LE)
		d.FieldU8("audio_channels")
		d.FieldU32LE("audio_sample_rate")
		d.FieldU32LE("bitrate_maximum")
		d.FieldU32LE("bitrate_nominal")
		d.FieldU32LE("bitrate_minimum")
		// TODO: code/comment about 2.1.4. coding bits into byte sequences
		d.FieldUFn("blocksize_1", func() (uint64, decode.DisplayFormat, string) {
			return 1 << d.U4(), decode.NumberDecimal, ""
		})
		d.FieldUFn("blocksize_0", func() (uint64, decode.DisplayFormat, string) {
			return 1 << d.U4(), decode.NumberDecimal, ""
		})
		// TODO: warning if blocksize0 > blocksize1
		// TODO: warning if not 64-8192
		d.FieldValidateZeroPadding("padding0", 7)
		d.FieldValidateUFn("framing_flag", 1, d.U1)

		// if d.BitsLeft() > 0 {
		// 	d.FieldValidateZeroPadding("padding1", int(d.BitsLeft()))
		// }
	case packetTypeSetup:
		d.FieldUFn("vorbis_codebook_count", func() (uint64, decode.DisplayFormat, string) {
			return d.U8() + 1, decode.NumberDecimal, ""
		})
		d.FieldValidateUFn("codecooke_sync", 0x564342, d.U24LE)
		d.FieldU16LE("codebook_dimensions")
		d.FieldU24LE("codebook_entries")

		// d.SeekRel(7)
		// ordered := d.FieldBool("ordered")

		// if ordered {

		// } else {
		// 	d.SeekRel(-2)
		// 	sparse := d.FieldBool("sparse")
		// 	d.SeekRel(1)

		// 	if sparse {

		// 	} else {
		// 		d.SeekRel(-7)
		// 		d.FieldU5("length")

		// 	}
		// }

	case packetTypeComment:
		// TODO: should not be try, FieldDecode?
		d.FieldDecode("comment", vorbisComment)

		// note this uses vorbis bitpacking convention, bits are added LSB first per byte
		d.FieldValidateZeroPadding("padding0", 7)
		d.FieldValidateUFn("frame_bit", 1, d.U1)

		// if d.BitsLeft() > 0 {
		// 	d.FieldValidateZeroPadding("padding1", int(d.BitsLeft()))
		// }
	}

	return nil
}
