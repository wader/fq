package mp3

// https://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header

import (
	"fq/internal/decode"
)

var XingHeader = &decode.Format{
	Name:      "xing_header",
	MIME:      "",
	New:       func() decode.Decoder { return &XingHeaderDecoder{} },
	SkipProbe: true,
}

// XingHeaderDecoder is a xing header decoder
type XingHeaderDecoder struct {
	decode.Common
}

// Decode decodes a xing header
func (d *XingHeaderDecoder) Decode() {
	// TODO: info has lame extension?
	hasLameExtension := false
	switch string(d.PeekBytes(4)) {
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
	d.FieldNoneFn("present_flags", func() {
		d.FieldU("unused", 28)
		qualityPresent = d.FieldBool("quality")
		tocPresent = d.FieldBool("toc")
		bytesPresent = d.FieldBool("bytes")
		framesPresent = d.FieldBool("frames")
	})

	if framesPresent {
		d.FieldU32("frames")
	}
	if bytesPresent {
		d.FieldU32("bytes")
	}
	if tocPresent {
		d.FieldBytesLen("toc", 100)
	}
	if qualityPresent {
		d.FieldU32("quality")
	}

	if hasLameExtension {
		d.FieldNoneFn("lame_extension", func() {
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
			d.FieldU12("encoder_delay_x") // TODO:
			d.FieldU12("encoder_delay_y") // TODO:
			d.FieldU8("misc")             // TODO:
			d.FieldU8("mp3_gain")         // TODO:
			d.FieldU16("preset")          // TODO:
			d.FieldU32("length")
			d.FieldU16("music_crc") // TODO:
			d.FieldU16("tag_crc")   // TODO:
		})
	}
}
