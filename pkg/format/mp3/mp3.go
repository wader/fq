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

var File = format.MustRegister(&decode.Format{
	Name:  "mp3",
	MIMEs: []string{"audio/mpeg"},
	New:   func() decode.Decoder { return &FileDecoder{} },
	Deps: []decode.Dep{
		{Names: []string{"id3v2"}, Formats: &headerTag},
		{Names: []string{"id3v1", "id3v11"}, Formats: &footerTags},
		{Names: []string{"apev2"}, Formats: &apeTag},
		{Names: []string{"mp3_frame"}, Formats: &mp3Frame},
	},
})

// FileDecoder is a MP3 decoder
type FileDecoder struct{ decode.Common }

// Decode decodes a MP3 stream
func (d *FileDecoder) Decode() {

	d.FieldTryDecode("header", headerTag)

	footerLen := int64(0)

	id3v1Len := int64(128 * 8)
	if d.BitsLeft() >= id3v1Len {
		if fd, _ := d.FieldTryDecodeRange(
			"footer", d.Pos()+d.BitsLeft()-id3v1Len, id3v1Len,
			footerTags); fd != nil {
			footerLen = id3v1Len
		}
	}

	validFrames := 0
	d.SubLenFn(d.BitsLeft()-footerLen, func() {
		d.MultiField("frame", func() {
			for !d.End() {
				if _, errs := d.FieldTryDecode("frame", mp3Frame); errs != nil {
					break
				}

				validFrames++
			}
		})

		d.FieldTryDecode("footer", apeTag)

		// TODO: truncated last frame?
		if d.BitsLeft() > 0 {
			d.FieldNoneFn("unknown", func() { d.SeekRel(int64(d.BitsLeft())) })
		}
	})

	if validFrames == 0 {
		d.Invalid("no frames found")
	}
}
