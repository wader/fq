package mpeg

// https://wiki.multimedia.cx/index.php/MPEG-4_Audio

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(decode.Format{
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

var channelConfigurationNames = decode.UToStr{
	0: "defined in AOT Specifc Config",
	1: "front-center",
	2: "front-left, front-right",
	3: "front-center, front-left, front-right",
	4: "front-center, front-left, front-right, back-center",
	5: "front-center, front-left, front-right, back-left, back-right",
	6: "front-center, front-left, front-right, back-left, back-right, LFE-channel",
	7: "front-center, front-left, front-right, side-left, side-right, back-left, back-right, LFE-channel",
}

func ascDecoder(d *decode.D, in interface{}) interface{} {
	objectType := d.FieldUFn("object_type", func(d *decode.D) uint64 {
		n := d.U5()
		if n == 31 {
			n = 32 + d.U6()
		}
		return n
	}, d.MapUToStrSym(format.MPEGAudioObjectTypeNames))
	d.FieldUScalarFn("sampling_frequency", func(d *decode.D) decode.Scalar {
		v := d.U4()
		if v == 15 {
			return decode.Scalar{Actual: d.U24()}
		}
		if f, ok := frequencyIndexHz[v]; ok {
			return decode.Scalar{Actual: v, Sym: f}
		}
		return decode.Scalar{Description: "invalid"}
	})
	d.FieldU4("channel_configuration", d.MapUToStrSym(channelConfigurationNames))
	// TODO: GASpecificConfig etc
	d.FieldRawLen("var_aot_or_byte_align", d.BitsLeft())

	return format.MPEGASCOut{ObjectType: int(objectType)}
}
