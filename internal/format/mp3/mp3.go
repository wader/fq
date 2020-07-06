package mp3

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3

// TODO: crc
// TODO: same sample decode?

import (
	"fq/internal/decode"
	"fq/internal/format/id3v1"
	"fq/internal/format/id3v11"
	"fq/internal/format/id3v2"
)

var File = &decode.Format{
	Name: "mp3",
	MIME: "",
	New:  func() decode.Decoder { return &FileDecoder{} },
}

// FileDecoder is a MP3 decoder
type FileDecoder struct {
	decode.Common
}

// Decode decodes a MP3 stream
func (d *FileDecoder) Decode() {
	d.FieldDecode("header", id3v2.Tag)

	footerLen := uint64(0)
	id3v1Len := uint64(128 * 8)
	if d.BitsLeft() >= id3v1Len {
		// TODO: added before? sort when presenting? probe? add later?
		if d.FieldDecodeRange("footer", d.Pos()+d.BitsLeft()-id3v1Len, id3v1Len, id3v1.Tag, id3v11.Tag) {
			footerLen = id3v1Len
		}
	}

	validFrames := 0
	d.SubLen(d.BitsLeft()-footerLen, func() {
		for !d.End() {
			if !d.FieldDecode("frame", Frame) {
				break
			}
			validFrames++
		}
	})

	if validFrames == 0 {
		d.Invalid("no frames found")
	}
}
