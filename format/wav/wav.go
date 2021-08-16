package wav

// http://soundfile.sapp.org/doc/WaveFormat/
// https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/wavdec.c
// https://tech.ebu.ch/docs/tech/tech3285.pdf
// http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/WAVE/WAVE.html
// TDOO: audio/wav

import (
	"fmt"
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
	"strings"
)

var headerFormat []*decode.Format
var footerFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.WAV,
		ProbeOrder:  10, // after most others (overlap some with webp)
		Description: "WAV file",
		Groups:      []string{format.PROBE},
		DecodeFn:    wavDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ID3V2}, Formats: &headerFormat},
			{Names: []string{format.ID3V1, format.ID3V11}, Formats: &footerFormat},
		},
	})
}

const (
	formatExtensible = 0xfffe
)

// transformed from ffmpeg libavformat/riff.c
var audioFormatName = map[uint64]string{
	0x0001: "PCM",
	0x0002: "ADPCM_MS",
	0x0003: "PCM_FLOAT",
	/* must come after f32le in this list */
	0x0006: "PCM_ALAW",
	0x0007: "PCM_MULAW",
	0x000a: "WMAVOICE",
	0x0010: "ADPCM_IMA_OKI",
	0x0011: "ADPCM_IMA_WAV",
	/* must come after adpcm_ima_wav in this list */
	0x0017: "ADPCM_IMA_OKI",
	0x0020: "ADPCM_YAMAHA",
	0x0022: "TRUESPEECH",
	0x0031: "GSM_MS",
	0x0032: "GSM_MS", /* msn audio */
	0x0038: "AMR_NB", /* rogue format number */
	0x0042: "G723_1",
	0x0045: "ADPCM_G726",
	0x0014: "ADPCM_G726", /* g723 Antex */
	0x0040: "ADPCM_G726", /* g721 Antex */
	0x0050: "MP2",
	0x0055: "MP3",
	0x0057: "AMR_NB",
	0x0058: "AMR_WB",
	/* rogue format number */
	0x0061: "ADPCM_IMA_DK4",
	/* rogue format number */
	0x0062:           "ADPCM_IMA_DK3",
	0x0064:           "ADPCM_G726",
	0x0069:           "ADPCM_IMA_WAV",
	0x0075:           "METASOUND",
	0x0083:           "G729",
	0x00ff:           "AAC",
	0x0111:           "G723_1",
	0x0130:           "SIPR",
	0x0135:           "ACELP_KELVIN",
	0x0160:           "WMAV1",
	0x0161:           "WMAV2",
	0x0162:           "WMAPRO",
	0x0163:           "WMALOSSLESS",
	0x0165:           "XMA1",
	0x0166:           "XMA2",
	0x0200:           "ADPCM_CT",
	0x0215:           "DVAUDIO",
	0x0216:           "DVAUDIO",
	0x0270:           "ATRAC3",
	0x028f:           "ADPCM_G722",
	0x0401:           "IMC",
	0x0402:           "IAC",
	0x0500:           "ON2AVC",
	0x0501:           "ON2AVC",
	0x1500:           "GSM_MS",
	0x1501:           "TRUESPEECH",
	0x1600:           "AAC",
	0x1602:           "AAC_LATM",
	0x2000:           "AC3",
	0x2001:           "DTS",
	0x2048:           "SONIC",
	0x6c75:           "PCM_MULAW",
	0x706d:           "AAC",
	0x4143:           "AAC",
	0x594a:           "XAN_DPCM",
	0x729a:           "G729",
	0xa100:           "G723_1", /* Comverse Infosys Ltd. G723 1 */
	0xa106:           "AAC",
	0xa109:           "SPEEX",
	0xf1ac:           "FLAC",
	('S' << 8) + 'F': "ADPCM_SWF",
	/* HACK/FIXME: Does Vorbis in WAV/AVI have an (in)official ID? */
	('V' << 8) + 'o': "VORBIS",

	formatExtensible: "Extensible",
}

var subFormatNames = map[[16]byte]string{
	{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}: "PCM",
	{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}: "IEEE_FLOAT",
}

func decodeChunk(d *decode.D, expectedChunkID string, stringData bool) int64 { //nolint:unparam
	chunks := map[string]func(d *decode.D){
		"RIFF": func(d *decode.D) {
			d.FieldUTF8("format", 4)
			decodeChunks(d, false)
		},
		"fmt": func(d *decode.D) {
			audioFormat, _ := d.FieldStringMapFn("audio_format", audioFormatName, "Unknown", d.U16LE, decode.NumberDecimal)
			d.FieldU16LE("num_channels")
			d.FieldU32LE("sample_rate")
			d.FieldU32LE("byte_rate")
			d.FieldU16LE("block_align")
			d.FieldU16LE("bits_per_sample")

			if audioFormat == formatExtensible && d.BitsLeft() > 0 {
				d.FieldU16LE("extension_size")
				d.FieldU16LE("valid_bits_per_sample")
				d.FieldU32LE("channel_mask")
				d.FieldStringUUIDMapFn("sub_format", subFormatNames, "Unknown", func() []byte { return d.BytesLen(16) })
			}
		},
		"data": func(d *decode.D) {
			d.FieldBitBufLen("samples", d.BitsLeft())
		},
		"LIST": func(d *decode.D) {
			d.FieldUTF8("list_type", 4)
			decodeChunks(d, true)
		},
		"fact": func(d *decode.D) {
			d.FieldU32LE("sample_length")
		},
	}

	trimChunkID := d.FieldStrFn("id", func() (string, string) {
		return strings.TrimSpace(d.UTF8(4)), ""
	})
	if expectedChunkID != "" && trimChunkID != expectedChunkID {
		d.Invalid(fmt.Sprintf("expected chunk id %q found %q", expectedChunkID, trimChunkID))
	}
	const restOfFileLen = 0xffffffff
	chunkLen := int64(d.FieldUFn("size", func() (uint64, decode.DisplayFormat, string) {
		l := d.U32LE()
		if l == restOfFileLen {
			return l, decode.NumberHex, "rest of file"
		}
		return l, decode.NumberDecimal, ""
	}))

	if chunkLen == restOfFileLen {
		chunkLen = d.BitsLeft() / 8
	}

	if fn, ok := chunks[trimChunkID]; ok {
		d.DecodeLenFn(chunkLen*8, fn)
	} else {
		if stringData {
			d.FieldStrFn("data", func() (string, string) {
				return strings.Trim(d.UTF8(int(chunkLen)), " \x00"), ""
			})
		} else {
			d.FieldBitBufLen("data", chunkLen*8)
		}
	}

	if chunkLen%2 != 0 {
		d.FieldBitBufLen("align", 8)
	}

	return chunkLen + 8
}

func decodeChunks(d *decode.D, stringData bool) {
	d.FieldStructArrayLoopFn("chunks", "chunk", d.NotEnd, func(d *decode.D) {
		decodeChunk(d, "", stringData)
	})
}

func wavDecode(d *decode.D, in interface{}) interface{} {
	// there are wav files in the wild with id3v2 header id3v1 footer
	_, _, _ = d.FieldTryFormat("header", headerFormat)

	decodeChunk(d, "RIFF", false)

	_, _, _ = d.FieldTryFormat("footer", footerFormat)

	return nil
}
