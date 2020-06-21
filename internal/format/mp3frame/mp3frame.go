package mp3frame

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3

// TODO: crc
// TODO: same sample decode?

import (
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "mp3frame",
	MIME: "",
	New:  func(common decode.Common) decode.Decoder { return &Decoder{Common: common} },
}

// Decoder is a mp3 frame decoder
type Decoder struct {
	decode.Common
}

// Decode MP3 frame
func (d *Decoder) Decode(opts decode.Options) {
	d.FieldValidateUFn("sync", 0b11111111111, d.U11)

	// v = 3 means version 2.5
	v := d.FieldUFn("mpeg_version", func() (uint64, decode.Format, string) {
		switch d.U2() {
		case 0b00:
			return 3, decode.FormatDecimal, "MPEG Version 2.5"
		case 0b01:
			return 0, decode.FormatDecimal, "reserved"
		case 0b10:
			return 2, decode.FormatDecimal, "MPEG Version 2"
		case 0b11:
			return 1, decode.FormatDecimal, "MPEG Version 1"
		default:
			panic("unreachable")
		}
	})
	l := d.FieldUFn("layer", func() (uint64, decode.Format, string) {
		switch d.U2() {
		case 0b00:
			return 0, decode.FormatDecimal, "reserved"
		case 0b01:
			return 3, decode.FormatDecimal, "Layer III"
		case 0b10:
			return 2, decode.FormatDecimal, "Layer II"
		case 0b11:
			return 1, decode.FormatDecimal, "Layer I"
		default:
			panic("unreachable")
		}
	})
	d.FieldStringMapFn("protection", map[uint64]string{
		0: "Protected by CRC",
		1: "Not protected",
	}, "", d.U1)
	// V1,L1 V1,L2 V1,L3  V2,L1 V2,L2 V2,L3  V2.5,L1 V2.5,L2 V2.5,L3
	var bitRateIndex = map[uint][9]uint{
		0b0001: [...]uint{32, 32, 32, 32, 8, 8, 32, 8, 8},
		0b0010: [...]uint{64, 48, 40, 48, 16, 16, 48, 16, 16},
		0b0011: [...]uint{96, 56, 48, 56, 24, 24, 56, 24, 24},
		0b0100: [...]uint{128, 64, 56, 64, 32, 32, 64, 32, 32},
		0b0101: [...]uint{160, 80, 64, 80, 40, 40, 80, 40, 40},
		0b0110: [...]uint{192, 96, 80, 96, 48, 48, 96, 48, 48},
		0b0111: [...]uint{224, 112, 96, 112, 56, 56, 112, 56, 56},
		0b1000: [...]uint{256, 128, 112, 128, 64, 64, 128, 64, 64},
		0b1001: [...]uint{288, 160, 128, 144, 80, 80, 144, 80, 80},
		0b1010: [...]uint{320, 192, 160, 160, 96, 96, 160, 96, 96},
		0b1011: [...]uint{352, 224, 192, 176, 112, 112, 176, 112, 112},
		0b1100: [...]uint{384, 256, 224, 192, 128, 128, 192, 128, 128},
		0b1101: [...]uint{416, 320, 256, 224, 144, 144, 224, 144, 144},
		0b1110: [...]uint{448, 384, 320, 256, 160, 160, 256, 160, 160},
	}
	bitRate := d.FieldUFn("bitrate", func() (uint64, decode.Format, string) {
		u := d.U4()
		switch u {
		case 0b0000:
			return 0, decode.FormatDecimal, "free"
		case 0b1111:
			return 0, decode.FormatDecimal, "bad"
		default:
			return uint64(bitRateIndex[uint(u)][(v-1)*3+(l-1)]) * 1000, decode.FormatDecimal, ""
		}
	})
	// MPEG1 MPEG2 MPEG2.5
	var sampleRateIndex = map[uint][3]uint{
		0b00: [...]uint{44100, 22050, 11025},
		0b01: [...]uint{48000, 24000, 12000},
		0b10: [...]uint{32000, 16000, 8000},
	}
	sampleRate := d.FieldUFn("sample_rate", func() (uint64, decode.Format, string) {
		u := d.U2()
		switch u {
		case 0b11:
			return 0, decode.FormatDecimal, "reserved"
		default:
			return uint64(sampleRateIndex[uint(u)][v-1]), decode.FormatDecimal, ""
		}
	})
	padding := d.FieldStringMapFn("padding", map[uint64]string{
		0: "Not padded",
		1: "Padded",
	}, "", d.U1)
	d.FieldU1("private")
	d.FieldStringMapFn("channels", map[uint64]string{
		0b00: "Stereo",
		0b01: "Joint Stereo",
		0b10: "Dual",
		0b11: "Mono",
	}, "", d.U2)
	d.FieldStringMapFn("channel_mode", map[uint64]string{
		0b00: "",
		0b01: "Intensity Stereo",
		0b10: "MS Stereo",
		0b11: "Intensity Stereo,MS Stereo",
	}, "", d.U2)
	d.FieldU1("copyright")
	d.FieldU1("original")
	d.FieldStringMapFn("emphasis", map[uint64]string{
		0b00: "None",
		0b01: "50/15",
		0b10: "reserved",
		0b11: "CCIT J.17",
	}, "", d.U2)

	const headerLen = 4
	dataLen := (144 * bitRate / sampleRate) + padding - headerLen
	d.FieldBytes("samples", dataLen)
}
