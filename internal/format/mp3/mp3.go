package mp3

// http://mpgedit.org/mpgedit/mpeg_format/MP3Format.html
// http://www.multiweb.cz/twoinches/MP3inside.htm
// https://wiki.hydrogenaud.io/index.php?title=MP3

// TODO: crc
// TODO: same sample decode?

import (
	"fq/internal/decode"
	"log"
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
func (d *Decoder) Decode(opts decode.Options) bool {
	return false

	d.FieldDecode("header", d.BitsLeft(), []string{"id3v2"})

	// TODO: recuseive.. stackverflow.. pass list of decoders?
	mp3FramesLen := d.BitsLeft()
	id3v1Len := uint64(128 * 8)
	log.Printf("mp3FramesLen: %#+v\n", mp3FramesLen)

	log.Printf("d.Pos1: %#+v\n", d.Pos)

	if mp3FramesLen >= id3v1Len {
		// TODO: added before? sort when presenting? probe? add later?
		if d.FieldDecodeRange("footer", d.Pos+mp3FramesLen-id3v1Len, id3v1Len, []string{"id3v1", "id3v11"}) {
			//mp3FramesLen -= id3v1Len
			d.Len -= id3v1Len
			log.Println("TASDSAD")
		}
	}

	log.Printf("d.Pos2: %#+v\n", d.Pos)

	// TODO: sub m3p frames thiny?
	//mp3frameBitBuf, _ := d.BitBufLen(mp3FramesLen)
	//d.Len = mp3FramesLen

	for !d.End() {
		if !d.FieldDecode("frame", d.BitsLeft(), []string{"mp3frame"}) {
			break
		}
	}

	return true
}
