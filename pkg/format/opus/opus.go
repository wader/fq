package opus

// TODO: maybe this file format should be ogg_opus?
// https://tools.ietf.org/html/rfc7845

import (
	"bytes"
	"fq/pkg/decode"
	"fq/pkg/format"
)

var vorbisComment []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.OPUS_PACKET,
		Description: "Opus packet",
		DecodeFn:    vorbisDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.VORBIS_COMMENT}, Formats: &vorbisComment},
		},
	})
}

func vorbisDecode(d *decode.D) interface{} {
	var prefix []byte
	if d.BitsLeft() >= 8*8 {
		prefix = d.PeekBytes(8)
	}
	switch {
	case bytes.Equal(prefix, []byte("OpusHead")):
		d.FieldValueStr("type", "head", "")
		d.FieldUTF8("prefix", 8)
		d.FieldU8("version")
		channelCount := d.FieldU8("channel_count")
		d.FieldU16("pre_skip")
		d.FieldU32LE("sample_rate")
		d.FieldU16LE("output_gain")
		mapFamily := d.FieldU8("map_family")
		if mapFamily != 0 {
			d.FieldU8("stream_count")
			d.FieldU8("coupled_count")
			i := uint64(0)
			d.FieldArrayLoopFn("channel_mapping", func() bool { return i < channelCount }, func(d *decode.D) {
				d.FieldU8("channel_mapping")
			})
		}
	case bytes.Equal(prefix, []byte("OpusTags")):
		d.FieldValueStr("type", "tags", "")
		d.FieldUTF8("prefix", 8)
		d.FieldDecode("comment", vorbisComment)
	default:
		d.FieldValueStr("type", "audio", "")
		d.FieldBitBufLen("data", d.BitsLeft())
	}

	return nil
}
