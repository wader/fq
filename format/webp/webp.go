package webp

// https://developers.google.com/speed/webp/docs/riff_container

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var vp8Frame decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.WEBP,
		Description: "WebP image",
		Groups:      []string{format.PROBE, format.IMAGE},
		DecodeFn:    webpDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.VP8_FRAME}, Group: &vp8Frame},
		},
	})
}

func decodeChunk(d *decode.D, expectedChunkID string, fn func(d *decode.D)) bool { //nolint:unparam
	trimChunkID := d.FieldUTF8("id", 4, d.TrimSpace)
	if expectedChunkID != "" && trimChunkID != expectedChunkID {
		return false
	}
	chunkLen := int64(d.FieldU32LE("size"))

	if fn != nil {
		d.LenFn(chunkLen*8, fn)
	} else {
		d.FieldRawLen("data", chunkLen*8)
	}

	return true
}

func webpDecode(d *decode.D, in interface{}) interface{} {
	d.FieldUTF8("riff_id", 4, d.AssertStr("RIFF"))
	riffLength := d.FieldU32LE("riff_length")
	d.FieldUTF8("webp_id", 4, d.AssertStr("WEBP"))

	d.LenFn(int64(riffLength-4)*8, func(d *decode.D) {
		p := d.PeekBytes(4)

		// TODO: VP8X

		switch {
		case bytes.Equal(p, []byte("VP8 ")):
			d.FieldStruct("image", func(d *decode.D) {
				decodeChunk(d, "VP8", func(d *decode.D) {
					d.Format(vp8Frame, nil)
				})
			})
		case bytes.Equal(p, []byte("VP8L")):
			d.FieldStruct("image", func(d *decode.D) {
				decodeChunk(d, "VP8L", func(d *decode.D) {
					// TODO
				})
			})
		default:
			d.Fatalf("could not find VP8 or VP8L chunk")
		}
	})

	return nil
}
