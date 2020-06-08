package mpeg

// https://wiki.multimedia.cx/index.php/MPEG-4_Audio

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.MPEG_ASC,
		Description: "MPEG-4 Audio Specific Config",
		DecodeFn:    ascDecoder,
	})
}

var frequencyIndexHz = map[uint64]int{
	0x0: 96000,
	0x1: 88200,
	0x2: 64000,
	0x3: 48000,
	0x4: 44100,
	0x5: 32000,
	0x6: 24000,
	0x7: 22050,
	0x8: 16000,
	0x9: 12000,
	0xa: 11025,
	0xb: 8000,
	0xc: 7350,
	0xd: -1,
	0xe: -1,
	0xf: -1,
}

var channelConfigurationNames = map[uint64]string{
	0: "Defined in AOT Specifc Config",
	1: "channel: front-center",
	2: "channels: front-left, front-right",
	3: "channels: front-center, front-left, front-right",
	4: "channels: front-center, front-left, front-right, back-center",
	5: "channels: front-center, front-left, front-right, back-left, back-right",
	6: "channels: front-center, front-left, front-right, back-left, back-right, LFE-channel",
	7: "channels: front-center, front-left, front-right, side-left, side-right, back-left, back-right, LFE-channel",
}

func ascDecoder(d *decode.D, in interface{}) interface{} {
	objectType, _ := d.FieldStringMapFn("object_type", format.MPEGAudioObjectTypeNames, "Unknown", func() uint64 {
		n := d.U5()
		if n == 31 {
			n = 32 + d.U6()
		}
		return n
	}, decode.NumberDecimal)
	d.FieldUFn("frequency_index", func() (uint64, decode.DisplayFormat, string) {
		v := d.U4()
		if v == 15 {
			return d.U24(), decode.NumberDecimal, ""
		}
		if f, ok := frequencyIndexHz[v]; ok {
			return uint64(f), decode.NumberDecimal, ""
		}
		return 0, decode.NumberDecimal, "Invalid"
	})
	d.FieldStringMapFn("channel_configuration", channelConfigurationNames, "Reserved", d.U4, decode.NumberDecimal)
	// TODO: GASpecificConfig etc
	d.FieldBitBufLen("var_aot_or_byte_align", d.BitsLeft())

	return format.MPEGASCOut{ObjectType: int(objectType)}
}
