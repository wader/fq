package mp3

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3
// https://www.diva-portal.org/smash/get/diva2:830195/FULLTEXT01.pdf

// TODO: crc
// TODO: same sample decode?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var xingHeader []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:     format.MP3_FRAME,
		DecodeFn: frameDecode,
		Deps: []decode.Dep{
			{Names: []string{format.XING_HEADER}, Formats: &xingHeader},
		},
	})
}

func frameDecode(d *decode.D) interface{} {
	d.FieldValidateUFn("sync", 0b11111111111, d.U11)

	// v = 3 means version 2.5
	v := d.FieldUFn("mpeg_version", func() (uint64, decode.DisplayFormat, string) {
		switch d.U2() {
		case 0b00:
			return 3, decode.NumberDecimal, "MPEG Version 2.5"
		case 0b01:
			return 0, decode.NumberDecimal, "reserved"
		case 0b10:
			return 2, decode.NumberDecimal, "MPEG Version 2"
		case 0b11:
			return 1, decode.NumberDecimal, "MPEG Version 1"
		default:
			panic("unreachable")
		}
	})
	l := d.FieldUFn("layer", func() (uint64, decode.DisplayFormat, string) {
		switch d.U2() {
		case 0b00:
			return 0, decode.NumberDecimal, "reserved"
		case 0b01:
			return 3, decode.NumberDecimal, "Layer III"
		case 0b10:
			return 2, decode.NumberDecimal, "Layer II"
		case 0b11:
			return 1, decode.NumberDecimal, "Layer I"
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
	bitRate := d.FieldUFn("bitrate", func() (uint64, decode.DisplayFormat, string) {
		u := d.U4()
		switch u {
		case 0b0000:
			return 0, decode.NumberDecimal, "free"
		case 0b1111:
			return 0, decode.NumberDecimal, "bad"
		default:
			return uint64(bitRateIndex[uint(u)][(v-1)*3+(l-1)]) * 1000, decode.NumberDecimal, ""
		}
	})
	// MPEG1 MPEG2 MPEG2.5
	var sampleRateIndex = map[uint][3]uint{
		0b00: [...]uint{44100, 22050, 11025},
		0b01: [...]uint{48000, 24000, 12000},
		0b10: [...]uint{32000, 16000, 8000},
	}
	sampleRate := d.FieldUFn("sample_rate", func() (uint64, decode.DisplayFormat, string) {
		u := d.U2()
		switch u {
		case 0b11:
			return 0, decode.NumberDecimal, "reserved"
		default:
			return uint64(sampleRateIndex[uint(u)][v-1]), decode.NumberDecimal, ""
		}
	})
	padding, _ := d.FieldStringMapFn("padding", map[uint64]string{
		0: "Not padded",
		1: "Padded",
	}, "", d.U1)
	d.FieldU1("private")
	channelsIndex, _ := d.FieldStringMapFn("channels", map[uint64]string{
		0b00: "Stereo",
		0b01: "Joint Stereo",
		0b10: "Dual",
		0b11: "Mono",
	}, "", d.U2)
	isStereo := channelsIndex != 0b11
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

	d.SubLenFn(int64(dataLen)*8, func(d *decode.D) {
		var sideInfoLen int64
		// [mono/stereo][mpeg version]
		sideInfoIndex := map[bool][4]int64{
			false: {0, 17, 9, 9},   // mono
			true:  {0, 32, 17, 17}, // stereo
		}
		if l == 3 {
			sideInfoLen = sideInfoIndex[isStereo][int(v)]
		}

		if sideInfoLen != 0 {
			d.FieldStructFn("side_info", func(d *decode.D) {
				d.FieldU9("main_data_begin")
				if isStereo {
					d.FieldU3("private_bits")
				} else {
					d.FieldU5("private_bits")
				}
				d.FieldU4("share0")
				if isStereo {
					d.FieldU4("share1")
				}

				granuleNr := 0
				d.FieldStructArrayLoopFn("granule", func() bool { return granuleNr < 2 }, func(d *decode.D) {
					// TODO: array for channels somehow?
					// TODO: tables and interpret values a bit

					d.FieldU12("part2_3_length0")
					if isStereo {
						d.FieldU12("part2_3_length1")
					}
					d.FieldU9("big_values0")
					if isStereo {
						d.FieldU9("big_values1")
					}
					d.FieldU8("global_gain0")
					if isStereo {
						d.FieldU8("global_gain1")
					}
					d.FieldU4("scalefac_compress0")
					if isStereo {
						d.FieldU4("scalefac_compress1")
					}
					d.FieldU1("window_switching0")
					if isStereo {
						d.FieldU1("window_switching1")
					}

					// normal blocks
					d.FieldU5("table_select0_0")
					d.FieldU5("table_select0_1")
					d.FieldU5("table_select0_2")
					if isStereo {
						d.FieldU5("table_select1_0")
						d.FieldU5("table_select1_1")
						d.FieldU5("table_select1_2")
					}
					d.FieldU4("region0_count0")
					if isStereo {
						d.FieldU4("region0_count1")
					}
					d.FieldU3("region1_count0")
					if isStereo {
						d.FieldU3("region1_count1")
					}

					d.FieldU1("preflag0")
					if isStereo {
						d.FieldU1("preflag1")
					}
					d.FieldU1("scalefac_scale0")
					if isStereo {
						d.FieldU1("scalefac_scale1")
					}
					d.FieldU1("count1table_select0")
					if isStereo {
						d.FieldU1("count1table_select1")
					}
					granuleNr++
				})
			})
		}

		d.FieldTryDecode("xing", xingHeader)

		d.FieldBitBufLen("samples", d.BitsLeft())
	})

	return nil
}
