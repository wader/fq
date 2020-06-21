package vorbis

// https://xiph.org/vorbis/doc/Vorbis_I_spec.html
// TODO: setup? more audio?

import (
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "vorbis",
	MIME: "",
	New: func(common decode.Common) decode.Decoder {
		return &Decoder{
			Common: common,
		}
	},
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

// Decoder is a vorbis decoder
type Decoder struct {
	decode.Common
}

// Decode vorbis
func (d *Decoder) Decode(opts decode.Options) {
	packetType := d.FieldUFn("packet_type", func() (uint64, decode.Format, string) {
		var packetType uint64
		packetTypeName := "unknown"
		t := d.U8()
		// 4.2.1. Common header decode
		// "these types are all odd as a packet with a leading single bit of ’0’ is an audio packet"
		if t&1 == 0 {
			packetType = 0
		} else {
			packetType = t
		}
		if n, ok := packetTypeNames[uint(packetType)]; ok {
			packetTypeName = n
		}
		return packetType, decode.FormatDecimal, packetTypeName
	})

	switch packetType {
	case packetTypeIdentification, packetTypeComment, packetTypeSetup:
		d.FieldValidateString("magic", "vorbis")
	}

	switch packetType {
	case packetTypeIdentification:
		d.FieldValidateUFn("vorbis_version", 0, d.U32LE)
		d.FieldU8("audio_channels")
		d.FieldU32LE("audio_sample_rate")
		d.FieldU32LE("bitrate_maximum")
		d.FieldU32LE("bitrate_nominal")
		d.FieldU32LE("bitrate_minimum")
		// TODO: code/comment about 2.1.4. coding bits into byte sequences
		d.FieldUFn("blocksize_1", func() (uint64, decode.Format, string) {
			return 1 << d.U4(), decode.FormatDecimal, ""
		})
		d.FieldUFn("blocksize_0", func() (uint64, decode.Format, string) {
			return 1 << d.U4(), decode.FormatDecimal, ""
		})
		// TODO: warning if blocksize0 > blocksize1
		// TODO: warning if not 64-8192
		d.FieldValidateZeroPadding("padding", 7)
		d.FieldValidateUFn("framing_flag", 1, d.U1)

		// 1   1) [vorbis_version] = read 32 bits as unsigned integer
		// 2   2) [audio_channels] = read 8 bit integer as unsigned
		// 3   3) [audio_sample_rate] = read 32 bits as unsigned integer
		// 4   4) [bitrate_maximum] = read 32 bits as signed integer
		// 5   5) [bitrate_nominal] = read 32 bits as signed integer
		// 6   6) [bitrate_minimum] = read 32 bits as signed integer
		// 7   7) [blocksize_0] = 2 exponent (read 4 bits as unsigned integer)
		// 8   8) [blocksize_1] = 2 exponent (read 4 bits as unsigned integer)
		// 9   9) [framing_flag] = read one bit

	case packetTypeComment:
		nameLen := func(name string) {
			len := d.FieldU32LE(name + "_length")
			d.FieldUTF8(name, len)
		}
		nameLen("vendor")
		userCommentListLength := d.FieldU32LE("user_comment_list_length")
		for i := uint64(0); i < userCommentListLength; i++ {
			nameLen("user_comment")
		}
		d.FieldValidateZeroPadding("padding", 7)
		d.FieldValidateUFn("frame_bit", 1, d.U1)
	}

	if d.BitsLeft() > 0 {
		d.FieldValidateZeroPadding("padding", d.BitsLeft())
	}
}
