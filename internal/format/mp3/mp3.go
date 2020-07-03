package mp3

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3

// TODO: crc
// TODO: same sample decode?

import (
	"fq/internal/decode"
)

var Register = &decode.Register{
	Name: "mp3",
	MIME: "",
	New:  func(common decode.Common) decode.Decoder { return &Decoder{Common: common} },
}

// Decoder is a mp3 decoder
type Decoder struct {
	decode.Common
}

// Decode MP3
func (d *Decoder) Decode(opts decode.Options) {
	d.FieldDecode("header", []string{"id3v2"})

	footerLen := uint64(0)
	id3v1Len := uint64(128 * 8)
	if d.BitsLeft() >= id3v1Len {
		// TODO: added before? sort when presenting? probe? add later?
		if d.FieldDecodeRange("footer", d.Pos()+d.BitsLeft()-id3v1Len, id3v1Len, []string{"id3v1", "id3v11"}) {
			footerLen = id3v1Len
		}
	}

	validFrames := 0
	d.SubLen(d.BitsLeft()-footerLen, func() {
		for !d.End() {
			if !d.FieldDecode("frame", []string{"mp3frame"}) {
				break
			}
			validFrames++
		}
	})

	if validFrames == 0 {
		d.Invalid("no frames found")
	}
}
