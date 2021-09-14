package mp3

// TODO: vbri
// TODO: mime audio/mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var headerFormat []*decode.Format
var footerFormat []*decode.Format
var mp3Frame []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.MP3,
		ProbeOrder:  10, // after most others (overlap some with other formats)
		Description: "MP3 file",
		Groups:      []string{format.PROBE},
		DecodeFn:    mp3Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ID3V2}, Formats: &headerFormat},
			{Names: []string{
				format.ID3V1,
				format.ID3V11,
				format.APEV2,
			}, Formats: &footerFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3Frame},
		},
	})
}

func mp3Decode(d *decode.D, in interface{}) interface{} {
	// there are mp3s files in the wild with multiple headers, two id3v2 tags etc
	d.FieldArrayFn("headers", func(d *decode.D) {
		for d.NotEnd() {
			if dv, _, _ := d.FieldTryFormat("header", headerFormat); dv == nil {
				return
			}
		}
	})

	lastValidEnd := int64(0)
	validFrames := 0
	decodeFailures := 0
	d.FieldArrayFn("frames", func(d *decode.D) {
		for d.NotEnd() {
			syncLen, _, err := d.TryPeekFind(16, 8, -1, func(v uint64) bool {
				return (v&0b1111_1111_1110_0000 == 0b1111_1111_1110_0000 && // sync header
					v&0b0000_0000_0001_1000 != 0b0000_0000_0000_1000 && // not reserved mpeg version
					v&0b0000_0000_0000_0110 == 0b0000_0000_0000_0010) // layer 3
			})
			if err != nil {
				break
			}
			if syncLen > 0 {
				d.SeekRel(syncLen)
			}

			if dv, _, _ := d.FieldTryFormat("frame", mp3Frame); dv == nil {
				decodeFailures++
				d.SeekRel(8)
				continue
			}
			lastValidEnd = d.Pos()
			validFrames++
		}
	})
	// TODO: better validate
	if validFrames == 0 || (validFrames < 2 && decodeFailures > 0) {
		d.Invalid("no frames found")
	}

	d.SeekAbs(lastValidEnd)

	d.FieldArrayFn("footers", func(d *decode.D) {
		for d.NotEnd() {
			if dv, _, _ := d.FieldTryFormat("footer", footerFormat); dv == nil {
				return
			}
		}
	})

	return nil
}
