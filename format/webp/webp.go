package webp

// https://developers.google.com/speed/webp/docs/riff_container

import (
	"bytes"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var vp8Frame []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.WEBP,
		Description: "WebP image",
		Groups:      []string{format.PROBE, format.IMAGE},
		DecodeFn:    webpDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.VP8_FRAME}, Formats: &vp8Frame},
		},
	})
}

func decodeChunk(d *decode.D, expectedChunkID string, fn func(d *decode.D)) bool { //nolint:unparam
	trimChunkID := d.FieldStrFn("id", func() (string, string) {
		return strings.TrimSpace(d.UTF8(4)), ""
	})
	if expectedChunkID != "" && trimChunkID != expectedChunkID {
		return false
	}
	chunkLen := int64(d.FieldU32LE("size"))

	if fn != nil {
		d.DecodeLenFn(chunkLen*8, fn)
	} else {
		d.FieldBitBufLen("data", chunkLen*8)
	}

	return true
}

func webpDecode(d *decode.D, in interface{}) interface{} {
	d.FieldValidateUTF8("riff_id", "RIFF")
	riffLength := d.FieldU32LE("riff_length")
	d.FieldValidateUTF8("webp_id", "WEBP")

	d.DecodeLenFn(int64(riffLength-4)*8, func(d *decode.D) {
		p := d.PeekBytes(4)

		// TODO: VP8X

		switch {
		case bytes.Equal(p, []byte("VP8 ")):
			d.FieldStructFn("image", func(d *decode.D) {
				decodeChunk(d, "VP8", func(d *decode.D) {
					d.Format(vp8Frame, nil)
				})
			})
		case bytes.Equal(p, []byte("VP8L")):
			d.FieldStructFn("image", func(d *decode.D) {
				decodeChunk(d, "VP8L", func(d *decode.D) {
					// TODO
				})
			})
		default:
			d.Invalid("could not find VP8 or VP8L chunk")
		}
	})

	return nil
}
