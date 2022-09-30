package mp3

// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header
// http://gabriel.mp3-tech.org/mp3infotag.html
// https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/mp3dec.c

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.MP3_FRAME_TAGS,
		Description: "MP3 frame info/xing tags",
		DecodeFn:    mp3FrameTagsDecode,
	})
}

var headerUintDescription = scalar.StrMapDescription{
	"Xing": "XING variable bitrate",
	"Info": "XING variable bitrate",
	"VBRI": "Fraunhofer Encoder variable bitrate info",
}

func mp3FrameTagsDecode(d *decode.D, _ any) any {
	switch d.FieldUTF8("header", 4, headerUintDescription) {
	case "Xing",
		"Info":
		qualityPresent := false
		tocPresent := false
		bytesPresent := false
		framesPresent := false
		d.FieldStruct("present_flags", func(d *decode.D) {
			d.FieldU("unused", 28)
			qualityPresent = d.FieldBool("quality")
			tocPresent = d.FieldBool("toc")
			bytesPresent = d.FieldBool("bytes")
			framesPresent = d.FieldBool("frames")
		})

		if framesPresent {
			d.FieldU32BE("frames")
		}
		if bytesPresent {
			d.FieldU32BE("bytes")
		}
		if tocPresent {
			d.FieldArray("toc", func(d *decode.D) {
				for i := 0; i < 100; i++ {
					d.FieldU8("entry")
				}
			})
		}
		if qualityPresent {
			d.FieldU32BE("quality")
		}

		d.FieldUTF8("encoder", 9)
		d.FieldU4("tag_revision")
		d.FieldU4("vbr_method")
		d.FieldU8("lowpass_filter") // TODO: /100
		d.FieldU32("replay_gain_peak")
		d.FieldU16("radio_replay_gain")
		d.FieldU16("audiophile_replay_gain")
		d.FieldU4("lame_flags")
		d.FieldU4("lame_ath_type")
		d.FieldU8("abr_vbr")          // TODO:
		d.FieldU12("encoder_delay")   // TODO:
		d.FieldU12("encoder_padding") // TODO:
		d.FieldU8("misc")             // TODO:
		d.FieldU8("mp3_gain")         // TODO:
		d.FieldU16("preset")          // TODO:
		d.FieldU32("length")
		d.FieldU16("music_crc") // TODO:
		d.FieldU16("tag_crc")   // TODO:
	case "VBRI":
		d.FieldU16("version_id")
		d.FieldU16("delay")
		d.FieldU16("quality")
		d.FieldU32("length", scalar.UintDescription("Number of bytes"))
		d.FieldU32("frames", scalar.UintDescription("Number of frames"))
		tocEntries := d.FieldU16("toc_entries", scalar.UintDescription("Number of entries within TOC table"))
		d.FieldU16("scale_factor", scalar.UintDescription("Scale factor of TOC table entries"))
		tocEntrySize := d.FieldU16("toc_entry_size", d.UintAssert(1, 2, 3, 4), scalar.UintDescription("Size per table entry"))
		d.FieldU16("frame_per_entry", scalar.UintDescription("Frames per table entry"))
		d.FieldArray("toc", func(d *decode.D) {
			for i := 0; i < int(tocEntries); i++ {
				d.FieldU("entry", int(tocEntrySize)*8)
			}
		})
	default:
		d.Errorf("no vbri, info or xing header found")
	}

	return nil
}
