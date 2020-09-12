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

// FileDecoder is a WAV decoder
type FileDecoder struct{ decode.Common }

func (d *FileDecoder) decodeChunk(expectedChunkId string) int64 {
	chunks := map[string]func(){
		"RIFF": func() {
			d.FieldU32LE("format")
			d.decodeChunks()
		},
		"fmt ": func() {
			d.FieldU16LE("audio_format")
			d.FieldU16LE("num_channels")
			d.FieldU32LE("sample_rate")
			d.FieldU32LE("byte_rate")
			d.FieldU16LE("block_align")
			d.FieldU16LE("bits_per_sample")
		},
		"data": func() {
			d.FieldNoneFn("samples", func() {})
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
			d.FieldBytesLen("data", chunkLen)
		}

		if chunkLen%2 != 0 {
			d.FieldBytesLen("chunk_align", 1)
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
