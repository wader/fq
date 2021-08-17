package mp3

// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.XING,
		Description: "Xing header",
		DecodeFn:    xingDecode,
	})
}

func xingDecode(d *decode.D, in interface{}) interface{} {
	// TODO: info has lame extension?
	hasLameExtension := false
	switch d.FieldUTF8("header", 4) {
	case "Xing":
	case "Info":
		hasLameExtension = true
	default:
		d.Invalid("no xing header found")
	}

	qualityPresent := false
	tocPresent := false
	bytesPresent := false
	framesPresent := false
	d.FieldStructFn("present_flags", func(d *decode.D) {
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
		d.FieldArrayFn("toc", func(d *decode.D) {
			for i := 0; i < 100; i++ {
				d.FieldU8("entry")
			}
		})
	}
	if qualityPresent {
		d.FieldU32BE("quality")
	}

	if hasLameExtension {
		d.FieldStructFn("lame_extension", func(d *decode.D) {
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
		})
	}

	return nil
}
