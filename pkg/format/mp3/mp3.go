package mp3

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3

// TODO: crc
// TODO: same sample decode?
// TODO: vbri
// TOFO: resync on garbage?

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var headerTag []*decode.Format
var id3v1Tags []*decode.Format
var apeTag []*decode.Format
var mp3Frame []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MP3,
		Description: "MPEG audio layer 3 file",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"audio/mpeg"},
		DecodeFn:    mp3Decode,
		Deps: []decode.Dep{
			{Names: []string{format.ID3_V2}, Formats: &headerTag},
			{Names: []string{format.ID3_V1, format.ID3_V11}, Formats: &id3v1Tags},
			{Names: []string{format.APEV2}, Formats: &apeTag},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3Frame},
		},
	})
}

func mp3Decode(d *decode.D) interface{} {
	d.FieldTryDecode("header", headerTag)

	validFrames := 0
	d.FieldArrayFn("frame", func(d *decode.D) {
		for !d.End() {
			if _, _, errs := d.FieldTryDecode("frame", mp3Frame); errs != nil {
				break
			}

			validFrames++
		}
	})

	if validFrames == 0 {
		d.Invalid("no frames found")
	}

	// only check for footer if there was some frames
	d.FieldArrayFn("footer", func(d *decode.D) {
		d.FieldTryDecode("footer", apeTag)
		if d.BitsLeft() >= 128*8 {
			d.SeekAbs(d.Len() - 128*8)
			d.FieldTryDecode("footer", id3v1Tags)

		}
	})

	return nil
}
