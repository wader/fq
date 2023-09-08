package mp3

// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header
// http://gabriel.mp3-tech.org/mp3infotag.html
// https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/mp3dec.c

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(
		format.MP3_Frame_XING,
		&decode.Format{
			Description: "MP3 frame Xing/Info tag",
			Groups:      []*decode.Group{format.MP3_Frame_Tags},
			DecodeFn:    mp3FrameTagXingDecode,
		})
}

func mp3FrameTagXingDecode(d *decode.D) any {
	d.FieldUTF8("header", 4, d.StrAssert("Xing", "Info"))
	lamePresent := false
	qualityPresent := false
	tocPresent := false
	bytesPresent := false
	framesPresent := false
	d.FieldStruct("present_flags", func(d *decode.D) {
		d.FieldU("unused", 27)
		lamePresent = d.FieldBool("lame")
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

	// this is mix of what ffmpeg and mediainfo does to detect lame extensions
	peekLame, _ := d.TryPeekBytes(4)
	peekLaneStr := string(peekLame)
	hasLameHeader := (peekLaneStr == "LAME" ||
		peekLaneStr == "Lavf" ||
		peekLaneStr == "Lavc" ||
		peekLaneStr == "GOGO" ||
		peekLaneStr == "L3.9")

	if lamePresent || hasLameHeader {
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
	}

	return nil
}
