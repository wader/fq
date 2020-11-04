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
var footerTags []*decode.Format
var apeTag []*decode.Format
var mp3Frame []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:     format.MP3,
		Groups:   []string{format.PROBE},
		MIMEs:    []string{"audio/mpeg"},
		DecodeFn: mp3Decode,
		Deps: []decode.Dep{
			{Names: []string{format.ID3V2}, Formats: &headerTag},
			{Names: []string{format.ID3V1, "id3v11"}, Formats: &footerTags},
			{Names: []string{format.APEV2}, Formats: &apeTag},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3Frame},
		},
	})
}

func mp3Decode(d *decode.D) interface{} {
	d.FieldTryDecode("header", headerTag)

	footerLen := int64(0)

	id3v1Len := int64(128 * 8)
	if d.BitsLeft() >= id3v1Len {
		if fd, _, _ := d.FieldTryDecodeRange(
			"footer", d.Pos()+d.BitsLeft()-id3v1Len, id3v1Len,
			footerTags); fd != nil {
			footerLen = id3v1Len
		}
	}

	validFrames := 0
	d.SubLenFn(d.BitsLeft()-footerLen, func() {
		d.FieldArrayFn("frame", func(d *decode.D) {
			for !d.End() {
				if _, _, errs := d.FieldTryDecode("frame", mp3Frame); errs != nil {
					break
				}

				validFrames++
			}
		})

		d.FieldTryDecode("footer", apeTag)

		// TODO: truncated last frame?
		if d.BitsLeft() > 0 {
			// TODO: some better unknown/garbage handling? generic gap filling?
			d.FieldBitBufLen("unknown", d.BitsLeft())
		}
	})

	if validFrames == 0 {
		d.Invalid("no frames found")
	}

	return nil
}
