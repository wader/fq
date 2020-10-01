package wav

// http://soundfile.sapp.org/doc/WaveFormat/
// https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/wavdec.c

import (
	"fmt"
	"fq/pkg/decode"
	"strings"
)

var File = &decode.Format{
	Name:  "wav",
	MIMEs: []string{"audio/wav"},
	New:   func() decode.Decoder { return &FileDecoder{} },
}

// transformed from ffmpeg libavformat/riff.c
var audioFormatName = map[uint64]string{
	0x0001: "PCM",
	0x0002: "ADPCM_MS",
	0x0003: "PCM_FLOAT",
	/* must come after f32le in this list */
	0x0006: "PCM_ALAW",
	0x0007: "PCM_MULAW",
	0x000A: "WMAVOICE",
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
	0x028F:           "ADPCM_G722",
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
	0x729A:           "G729",
	0xA100:           "G723_1", /* Comverse Infosys Ltd. G723 1 */
	0xA106:           "AAC",
	0xA109:           "SPEEX",
	0xF1AC:           "FLAC",
	('S' << 8) + 'F': "ADPCM_SWF",
	/* HACK/FIXME: Does Vorbis in WAV/AVI have an (in)official ID? */
	('V' << 8) + 'o': "VORBIS",
}

// FileDecoder is a WAV decoder
type FileDecoder struct{ decode.Common }

func (d *FileDecoder) decodeChunk(expectedChunkId string) int64 {
	chunks := map[string]func(){
		"RIFF": func() {
			d.FieldUTF8("format", 4)
			d.decodeChunks()
		},
		"fmt ": func() {
			d.FieldStringMapFn("audio_format", audioFormatName, "Unknown", d.U16LE)
			d.FieldU16LE("num_channels")
			d.FieldU32LE("sample_rate")
			d.FieldU32LE("byte_rate")
			d.FieldU16LE("block_align")
			d.FieldU16LE("bits_per_sample")
		},
		"data": func() {
			d.FieldBitBufLen("samples", d.BitsLeft())
		},
	}

	var chunkID string
	var chunkLen int64
	chunkID = d.UTF8(4)
	if expectedChunkId != "" && chunkID != expectedChunkId {
		d.Invalid(fmt.Sprintf("expected chunk id %q found %q", expectedChunkId, chunkID))
	}
	d.SeekRel(-4 * 8)

	trimChunkID := strings.TrimSpace(chunkID)
	d.FieldStrFn(trimChunkID, func() (string, string) {
		d.FieldUTF8("chunk_id", 4)
		chunkLen = int64(d.FieldU32LE("chunk_size"))

		if fn, ok := chunks[chunkID]; ok {
			d.SubLenFn(chunkLen*8, fn)
		} else {
			d.FieldBitBufLen("data", chunkLen*8)
		}

		if chunkLen%2 != 0 {
			d.FieldBitBufLen("chunk_align", 8)
		}

		return chunkID, ""
	})

	return chunkLen + 8
}

func (d *FileDecoder) decodeChunks() {
	for !d.End() {
		d.decodeChunk("")
	}
}

// Decode decodes a WAV stream
func (d *FileDecoder) Decode() {
	d.decodeChunk("RIFF")
}
