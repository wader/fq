package wav

// http://soundfile.sapp.org/doc/WaveFormat/
// https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/wavdec.c
// https://tech.ebu.ch/docs/tech/tech3285.pdf
// http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/WAVE/WAVE.html
// TODO: audio/wav
// TODO: default little endian

import (
	"fmt"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var headerFormat decode.Group
var footerFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.WAV,
		ProbeOrder:  10, // after most others (overlap some with webp)
		Description: "WAV file",
		Groups:      []string{format.PROBE},
		DecodeFn:    wavDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ID3V2}, Group: &headerFormat},
			{Names: []string{format.ID3V1, format.ID3V11}, Group: &footerFormat},
		},
	})
}

const (
	formatExtensible = 0xfffe
)

// transformed from ffmpeg libavformat/riff.c
var audioFormatName = scalar.UToSymStr{
	0x0001: "pcm",
	0x0002: "adpcm_ms",
	0x0003: "pcm_float",
	/* must come after f32le in this list */
	0x0006: "pcm_alaw",
	0x0007: "pcm_mulaw",
	0x000a: "wmavoice",
	0x0010: "adpcm_ima_oki",
	0x0011: "adpcm_ima_wav",
	/* must come after adpcm_ima_wav in this list */
	0x0017: "adpcm_ima_oki",
	0x0020: "adpcm_yamaha",
	0x0022: "truespeech",
	0x0031: "gsm_ms",
	0x0032: "gsm_ms", /* msn audio */
	0x0038: "amr_nb", /* rogue format number */
	0x0042: "g723_1",
	0x0045: "adpcm_g726",
	0x0014: "adpcm_g726", /* g723 Antex */
	0x0040: "adpcm_g726", /* g721 Antex */
	0x0050: "mp2",
	0x0055: "mp3",
	0x0057: "amr_nb",
	0x0058: "amr_wb",
	/* rogue format number */
	0x0061: "adpcm_ima_dk4",
	/* rogue format number */
	0x0062:           "adpcm_ima_dk3",
	0x0064:           "adpcm_g726",
	0x0069:           "adpcm_ima_wav",
	0x0075:           "metasound",
	0x0083:           "g729",
	0x00ff:           "aac",
	0x0111:           "g723_1",
	0x0130:           "sipr",
	0x0135:           "acelp_kelvin",
	0x0160:           "wmav1",
	0x0161:           "wmav2",
	0x0162:           "wmapro",
	0x0163:           "wmalossless",
	0x0165:           "xma1",
	0x0166:           "xma2",
	0x0200:           "adpcm_ct",
	0x0215:           "dvaudio",
	0x0216:           "dvaudio",
	0x0270:           "atrac3",
	0x028f:           "adpcm_g722",
	0x0401:           "imc",
	0x0402:           "iac",
	0x0500:           "on2avc",
	0x0501:           "on2avc",
	0x1500:           "gsm_ms",
	0x1501:           "truespeech",
	0x1600:           "aac",
	0x1602:           "aac_latm",
	0x2000:           "ac3",
	0x2001:           "dts",
	0x2048:           "sonic",
	0x6c75:           "pcm_mulaw",
	0x706d:           "aac",
	0x4143:           "aac",
	0x594a:           "xan_dpcm",
	0x729a:           "g729",
	0xa100:           "g723_1", /* Comverse Infosys Ltd. G723 1 */
	0xa106:           "aac",
	0xa109:           "speex",
	0xf1ac:           "flac",
	('S' << 8) + 'F': "adpcm_swf",
	/* HACK/FIXME: Does Vorbis in WAV/AVI have an (in)official ID? */
	('V' << 8) + 'o': "vorbis",

	formatExtensible: "extensible",
}

var (
	subFormatPCMBytes  = [16]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}
	subFormatIEEEFloat = [16]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}
)

var subFormatNames = scalar.BytesToScalar{
	{Bytes: subFormatPCMBytes[:], Scalar: scalar.S{Sym: "pcm"}},
	{Bytes: subFormatIEEEFloat[:], Scalar: scalar.S{Sym: "ieee_float"}},
}

func decodeChunk(d *decode.D, expectedChunkID string, stringData bool) int64 {
	d.Endian = decode.LittleEndian

	chunks := map[string]func(d *decode.D){
		"RIFF": func(d *decode.D) {
			d.FieldUTF8("format", 4)
			decodeChunks(d, false)
		},
		"fmt": func(d *decode.D) {
			audioFormat := d.FieldU16("audio_format", audioFormatName)
			d.FieldU16("num_channels")
			d.FieldU32("sample_rate")
			d.FieldU32("byte_rate")
			d.FieldU16("block_align")
			d.FieldU16("bits_per_sample")

			if audioFormat == formatExtensible && d.BitsLeft() > 0 {
				d.FieldU16("extension_size")
				d.FieldU16("valid_bits_per_sample")
				d.FieldU32("channel_mask")
				d.FieldRawLen("sub_format", 16*8, subFormatNames)
			}
		},
		"data": func(d *decode.D) {
			d.FieldRawLen("samples", d.BitsLeft())
		},
		"LIST": func(d *decode.D) {
			d.FieldUTF8("list_type", 4)
			decodeChunks(d, true)
		},
		"fact": func(d *decode.D) {
			d.FieldU32("sample_length")
		},
	}

	trimChunkID := d.FieldStrFn("id", func(d *decode.D) string {
		return strings.TrimSpace(d.UTF8(4))
	})
	if expectedChunkID != "" && trimChunkID != expectedChunkID {
		d.Errorf(fmt.Sprintf("expected chunk id %q found %q", expectedChunkID, trimChunkID))
	}
	const restOfFileLen = 0xffffffff
	chunkLen := int64(d.FieldUScalarFn("size", func(d *decode.D) scalar.S {
		l := d.U32()
		if l == restOfFileLen {
			return scalar.S{Actual: l, ActualDisplay: scalar.NumberHex, Description: "Rest of file"}
		}
		return scalar.S{Actual: l, ActualDisplay: scalar.NumberDecimal}
	}))

	if chunkLen == restOfFileLen {
		chunkLen = d.BitsLeft() / 8
	}

	if fn, ok := chunks[trimChunkID]; ok {
		d.FramedFn(chunkLen*8, fn)
	} else {
		if stringData {
			d.FieldUTF8("data", int(chunkLen), scalar.ActualTrim(" \x00"))
		} else {
			d.FieldRawLen("data", chunkLen*8)
		}
	}

	if chunkLen%2 != 0 {
		d.FieldRawLen("align", 8)
	}

	return chunkLen + 8
}

func decodeChunks(d *decode.D, stringData bool) {
	d.FieldStructArrayLoop("chunks", "chunk", d.NotEnd, func(d *decode.D) {
		decodeChunk(d, "", stringData)
	})
}

func wavDecode(d *decode.D, in any) any {
	// there are wav files in the wild with id3v2 header id3v1 footer
	_, _, _ = d.TryFieldFormat("header", headerFormat, nil)

	decodeChunk(d, "RIFF", false)

	_, _, _ = d.TryFieldFormat("footer", footerFormat, nil)

	return nil
}
