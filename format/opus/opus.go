package opus

// https://tools.ietf.org/html/rfc7845

import (
	"bytes"
	"fq/format"
	"fq/format/all/all"
	"fq/pkg/decode"
)

var vorbisComment []*decode.Format

func init() {
	all.MustRegister(&decode.Format{
		Name:        format.OPUS_PACKET,
		Description: "Opus packet",
		DecodeFn:    opusDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.VORBIS_COMMENT}, Formats: &vorbisComment},
		},
	})
}

func opusDecode(d *decode.D, in interface{}) interface{} {
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
			d.FieldArrayLoopFn("channel_mappings", func() bool { return i < channelCount }, func(d *decode.D) {
				d.FieldU8("channel_mapping")
			})
		}
	case bytes.Equal(prefix, []byte("OpusTags")):
		d.FieldValueStr("type", "tags", "")
		d.FieldUTF8("prefix", 8)
		d.FieldDecode("comment", vorbisComment)
	default:
		d.FieldValueStr("type", "audio", "")
		d.FieldStructFn("toc", func(d *decode.D) {
			d.FieldStructFn("config", func(d *decode.D) {
				configurations := map[uint64]struct {
					mode      string
					bandwidth string
					frameSize float64
				}{
					0:  {"SILK-only", "NB", 10},
					1:  {"SILK-only", "NB", 20},
					2:  {"SILK-only", "NB", 40},
					3:  {"SILK-only", "NB", 60},
					4:  {"SILK-only", "MB", 10},
					5:  {"SILK-only", "MB", 20},
					6:  {"SILK-only", "MB", 40},
					7:  {"SILK-only", "MB", 60},
					8:  {"SILK-only", "WB", 10},
					9:  {"SILK-only", "WB", 20},
					10: {"SILK-only", "WB", 40},
					11: {"SILK-only", "WB", 60},
					12: {"Hybrid", "SWB", 10},
					13: {"Hybrid", "SWB", 20},
					14: {"Hybrid", "FB", 10},
					15: {"Hybrid", "FB", 20},
					16: {"CELT-only", "NB", 2.5},
					17: {"CELT-only", "NB", 5},
					18: {"CELT-only", "NB", 10},
					19: {"CELT-only", "NB", 20},
					20: {"CELT-only", "WB", 2.5},
					21: {"CELT-only", "WB", 5},
					22: {"CELT-only", "WB", 10},
					23: {"CELT-only", "WB", 20},
					24: {"CELT-only", "SWB", 2.5},
					25: {"CELT-only", "SWB", 5},
					26: {"CELT-only", "SWB", 10},
					27: {"CELT-only", "SWB", 20},
					28: {"CELT-only", "FB", 2.5},
					29: {"CELT-only", "FB", 5},
					30: {"CELT-only", "FB", 10},
					31: {"CELT-only", "FB", 20},
				}
				n := d.FieldU5("config")
				config := configurations[n]
				d.FieldValueStr("mode", config.mode, "")
				d.FieldValueStr("bandwidth", config.bandwidth, "")
				d.FieldValueFloat("frame_size", config.frameSize, "")
			})
			d.FieldBool("stereo")
			d.FieldStructFn("frames_per_packet", func(d *decode.D) {
				framesPerPacketConfigs := map[uint64]struct {
					frames uint64
					mode   string
				}{
					0: {1, "1 frame"},
					1: {2, "2 frames, equal size"},
					2: {2, "2 frames, different size"},
					3: {0, "arbitrary number of frames"},
				}
				n := d.FieldU2("config")
				config := framesPerPacketConfigs[n]
				d.FieldValueU("frames", config.frames, "")
				d.FieldValueStr("mode", config.mode, "")
			})
			d.FieldBitBufLen("data", d.BitsLeft())
		})
	}

	return nil
}
