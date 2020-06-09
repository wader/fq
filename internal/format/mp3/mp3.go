package mp3

import (
	"fq/internal/decode"
)

// MP3 decoder
type MP3 struct {
	decode.Common
}

// Decode MP3
func (m *MP3) Decode() {
	m.FieldVerifyFn("sync", 0b11111111111, m.U11)
	m.FieldStringMapFn("mpeg_version", map[uint64]string{
		0b00: "MPEG Version 2.5",
		0b01: "reserved",
		0b10: "MPEG Version 2",
		0b11: "MPEG Version 1",
	}, "", m.U2)
	m.FieldStringMapFn("layers", map[uint64]string{
		0b00: "reserved",
		0b01: "Layer III",
		0b10: "Layer II",
		0b11: "Layer I",
	}, "", m.U2)
	m.FieldStringMapFn("protection", map[uint64]string{
		0: "Protected by CRC",
		1: "Not protected",
	}, "", m.U2)
	// TODO: not correct
	kBitRate := m.FieldUFn("bitrate", func() (uint64, decode.Format, string) {
		switch m.U4() {
		case 0b0000:
			return 0, decode.FormatDecimal, "free"
		case 0b0001:
			return 32, decode.FormatDecimal, ""
		case 0b0010:
			return 40, decode.FormatDecimal, ""
		case 0b0011:
			return 48, decode.FormatDecimal, ""
		case 0b0100:
			return 56, decode.FormatDecimal, ""
		case 0b0101:
			return 64, decode.FormatDecimal, ""
		case 0b0110:
			return 80, decode.FormatDecimal, ""
		case 0b0111:
			return 96, decode.FormatDecimal, ""
		case 0b1000:
			return 112, decode.FormatDecimal, ""
		case 0b1001:
			return 128, decode.FormatDecimal, ""
		case 0b1010:
			return 160, decode.FormatDecimal, ""
		case 0b1011:
			return 192, decode.FormatDecimal, ""
		case 0b1100:
			return 224, decode.FormatDecimal, ""
		case 0b1101:
			return 256, decode.FormatDecimal, ""
		case 0b1110:
			return 320, decode.FormatDecimal, ""
		case 0b1111:
			return 0, decode.FormatDecimal, "bad"
		default:
			panic("unreachable")
		}
	})
	sampleRate := m.FieldUFn("sample_rate", func() (uint64, decode.Format, string) {
		switch m.U4() {
		case 0b00:
			return 44100, decode.FormatDecimal, ""
		case 0b01:
			return 48000, decode.FormatDecimal, ""
		case 0b10:
			return 32000, decode.FormatDecimal, ""
		case 0b11:
			return 0, decode.FormatDecimal, "reserved"
		default:
			panic("unreachable")
		}
	})
	padding := m.FieldStringMapFn("padding", map[uint64]string{
		0: "Not padded",
		1: "Padded",
	}, "", m.U1)
	m.FieldU1("private")
	channels := m.FieldStringMapFn("channels", map[uint64]string{
		0b00: "Stereo",
		0b01: "Joint Stereo",
		0b10: "Dual",
		0b11: "Mono",
	}, "", m.U2)
	// Mode extension (only if Joint Stereo is set)
	if channels == 0b1 {
		m.FieldStringMapFn("mode_extension", map[uint64]string{
			0b00: "",
			0b01: "Intensity Stereo",
			0b10: "MS Stereo",
			0b11: "Intensity Stereo,MS Stereo",
		}, "", m.U2)
	}
	m.FieldU1("copyright")
	m.FieldU1("original")
	m.FieldStringMapFn("emphasis", map[uint64]string{
		0b00: "None",
		0b01: "50/15",
		0b10: "reserved",
		0b11: "CCIT J.17",
	}, "", m.U2)

	// frameLen := int((144 * kBitRate * 1000 / sampleRate) + padding)
	//dataLen := frameLen-

}
