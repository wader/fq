package mp3

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3

// TODO: crc
// TODO: same sample decode?
// TODO: vbri

import (
	"fq/pkg/decode"
	"fq/pkg/format/id3v1"
	"fq/pkg/format/id3v11"
	"fq/pkg/format/id3v2"
)

var File = &decode.Format{
	Name:  "mp3",
	MIMEs: []string{"audio/mpeg"},
	New:   func() decode.Decoder { return &FileDecoder{} },
}

// FileDecoder is a MP3 decoder
type FileDecoder struct{ decode.Common }

// Decode decodes a MP3 stream
func (d *FileDecoder) Decode() {
	d.FieldTryDecode("header", id3v2.Tag)

	footerLen := uint64(0)
	id3v1Len := uint64(128 * 8)
	if d.BitsLeft() >= id3v1Len {
		if fd, _ := d.FieldTryDecodeRange(
			"footer", d.Pos()+d.BitsLeft()-id3v1Len, id3v1Len,
			id3v1.Tag, id3v11.Tag); fd != nil {
			footerLen = id3v1Len
		}
	}

	validFrames := 0
	d.SubLenFn(d.BitsLeft()-footerLen, func() {
		for !d.End() {
			if _, errs := d.FieldTryDecode("frame", Frame); errs != nil {
				// TODO: truncated last frame?
				if d.BitsLeft() >= 0 {
					d.FieldNoneFn("unknown", func() { d.SeekRel(int64(d.BitsLeft())) })
				}
				break
			}

			validFrames++
		}
	})

	if validFrames == 0 {
		d.Invalid("no frames found")
	}
}
