package mpeg

// https://wiki.multimedia.cx/index.php/MPEG-4_Audio

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_ASC,
		Description: "MPEG-4 Audio specific config",
		DecodeFn:    ascDecoder,
	})
}

var aotNames = map[uint64]string{
	0:  "Null",
	1:  "AAC Main",
	2:  "AAC LC (Low Complexity)",
	3:  "AAC SSR (Scalable Sample Rate)",
	4:  "AAC LTP (Long Term Prediction)",
	5:  "SBR (Spectral Band Replication)",
	6:  "AAC Scalable",
	7:  "TwinVQ",
	8:  "CELP (Code Excited Linear Prediction)",
	9:  "HXVC (Harmonic Vector eXcitation Coding)",
	10: "Reserved",
	11: "Reserved",
	12: "TTSI (Text-To-Speech Interface)",
	13: "Main Synthesis",
	14: "Wavetable Synthesis",
	15: "General MIDI",
	16: "Algorithmic Synthesis and Audio Effects",
	17: "ER (Error Resilient) AAC LC",
	18: "Reserved",
	19: "ER AAC LTP",
	20: "ER AAC Scalable",
	21: "ER TwinVQ",
	22: "ER BSAC (Bit-Sliced Arithmetic Coding)",
	23: "ER AAC LD (Low Delay)",
	24: "ER CELP",
	25: "ER HVXC",
	26: "ER HILN (Harmonic and Individual Lines plus Noise)",
	27: "ER Parametric",
	28: "SSC (SinuSoidal Coding)",
	29: "PS (Parametric Stereo)",
	30: "MPEG Surround",
	31: "(Escape value)",
	32: "Layer-1",
	33: "Layer-2",
	34: "Layer-3",
	35: "DST (Direct Stream Transfer)",
	36: "ALS (Audio Lossless)",
	37: "SLS (Scalable LosslesS)",
	38: "SLS non-core",
	39: "ER AAC ELD (Enhanced Low Delay)",
	40: "SMR (Symbolic Music Representation) Simple",
	41: "SMR Main",
	42: "USAC (Unified Speech and Audio Coding) (no SBR)",
	43: "SAOC (Spatial Audio Object Coding)",
	44: "LD MPEG Surround",
	45: "USAC",
}

var frequencyIndexHz = map[uint64]int{
	0:  96000,
	1:  88200,
	2:  64000,
	3:  48000,
	4:  44100,
	5:  32000,
	6:  24000,
	7:  22050,
	8:  16000,
	9:  12000,
	10: 11025,
	11: 8000,
	12: 7350,
	13: -1,
	14: -1,
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
	d.FieldStringMapFn("object_type", aotNames, "Unknown", func() uint64 {
		n := d.U5()
		if n == 31 {
			n = 32 + d.U6()
		}
		return n
	})
	d.FieldUFn("frequence_index", func() (uint64, decode.DisplayFormat, string) {
		v := d.U4()
		if v == 15 {
			return d.U24(), decode.NumberDecimal, ""
		}
		if f, ok := frequencyIndexHz[v]; ok {
			return uint64(f), decode.NumberDecimal, ""
		}
		return 0, decode.NumberDecimal, "Invalid"
	})
	d.FieldStringMapFn("channel_configuration", channelConfigurationNames, "Reserved", d.U4)
	// TODO: GASpecificConfig etc
	d.FieldBitBufLen("var_aot_or_byte_align", d.BitsLeft())

	return nil
}
