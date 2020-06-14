package mp3

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3

// TODO: crc
// TODO: same sample decode?

import (
	"fq/internal/decode"
	"fq/internal/format/id3v2"
)

var Register = &decode.Register{
	Name: "mp3",
	MIME: "",
	New:  func() decode.Decoder { return &Decoder{} },
}

// Decoder is a mp3 decoder
type Decoder struct {
	decode.Common
}

// Decode MP3
func (d *Decoder) Decode(opts decode.Options) bool {
	p := id3v2.Decoder{Common: d.Common}
	p.Decode(opts)

	d.Common = p.Common

	mp3FramesLen := d.BitsLeft()
	if mp3FramesLen >= 128*8 {
		if d.Decode("x-fq/id3v1", d.BitBufRange(d.Len-(128*8), 128*8)) {
			// - return value?
			mp3FramesLen -= 128 * 8
		}
	}

	mp3frameBitBuf := d.BitBufLen(mp3FramesLen)

	for !mp3frameBitBuf.End() {
		d.FieldNoneFn("frame", func() {
		})
	}

	return true
}
