package mp3

// TODO: vbri
// TOFO: resync on garbage? between id3v2 and first frame for example

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var headerFormat []*decode.Format
var footerFormat []*decode.Format
var mp3Frame []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MP3,
		Description: "MPEG audio layer 3 file",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"audio/mpeg"},
		DecodeFn:    mp3Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ID3_V2}, Formats: &headerFormat},
			{Names: []string{format.ID3_V1, format.ID3_V11, format.APEV2}, Formats: &footerFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3Frame},
		},
	})
}

func mp3Decode(d *decode.D, in interface{}) interface{} {
	// there are mp3s files in the wild with multiple headers, two id3v2 tags etc
	d.FieldArrayFn("headers", func(d *decode.D) {
		for d.NotEnd() {
			if _, _, err := d.FieldTryDecode("header", headerFormat); err != nil {
				return
			}
		}
	})

	validFrames := 0
	d.FieldArrayFn("frames", func(d *decode.D) {
		for d.NotEnd() {
			if _, _, errs := d.FieldTryDecode("frame", mp3Frame); errs != nil {
				break
			}
			validFrames++
		}
	})
	if validFrames == 0 {
		d.Invalid("no frames found")
	}

	d.FieldArrayFn("footers", func(d *decode.D) {
		for d.NotEnd() {
			if _, _, err := d.FieldTryDecode("footer", headerFormat); err != nil {
				return
			}
		}
	})

	return nil
}
