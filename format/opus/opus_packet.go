package opus

// https://tools.ietf.org/html/rfc7845

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var vorbisComment decode.Group

func init() {
	interp.RegisterFormat(
		format.Opus_Packet,
		&decode.Format{
			Description: "Opus packet",
			DecodeFn:    opusDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Vorbis_Comment}, Out: &vorbisComment},
			},
		})
}

func opusDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	var prefix []byte
	if d.BitsLeft() >= 8*8 {
		prefix = d.PeekBytes(8)
	}
	switch {
	case bytes.Equal(prefix, []byte("OpusHead")):
		d.FieldValueStr("type", "head")
		d.FieldUTF8("prefix", 8)
		d.FieldU8("version")
		channelCount := d.FieldU8("channel_count")
		d.FieldU16("pre_skip")
		d.FieldU32("sample_rate")
		d.FieldU16("output_gain")
		mapFamily := d.FieldU8("map_family")
		if mapFamily != 0 {
			d.FieldU8("stream_count")
			d.FieldU8("coupled_count")
			i := uint64(0)
			d.FieldArrayLoop("channel_mappings", func() bool { return i < channelCount }, func(d *decode.D) {
				d.FieldU8("channel_mapping")
			})
		}
	case bytes.Equal(prefix, []byte("OpusTags")):
		d.FieldValueStr("type", "tags")
		d.FieldUTF8("prefix", 8)
		d.FieldFormat("comment", &vorbisComment, nil)
	default:
		d.FieldValueStr("type", "audio")
		d.FieldStruct("toc", func(d *decode.D) {
			d.FieldStruct("config", func(d *decode.D) {
				configurations := map[uint64]struct {
					mode      string
					bandwidth string
					frameSize float64
				}{
					0:  {"silk_only", "nb", 10},
					1:  {"silk_only", "nb", 20},
					2:  {"silk_only", "nb", 40},
					3:  {"silk_only", "nb", 60},
					4:  {"silk_only", "mb", 10},
					5:  {"silk_only", "mb", 20},
					6:  {"silk_only", "mb", 40},
					7:  {"silk_only", "mb", 60},
					8:  {"silk_only", "wb", 10},
					9:  {"silk_only", "wb", 20},
					10: {"silk_only", "wb", 40},
					11: {"silk_only", "wb", 60},
					12: {"hybrid", "swb", 10},
					13: {"hybrid", "swb", 20},
					14: {"hybrid", "fb", 10},
					15: {"hybrid", "fb", 20},
					16: {"celt_only", "nb", 2.5},
					17: {"celt_only", "nb", 5},
					18: {"celt_only", "nb", 10},
					19: {"celt_only", "nb", 20},
					20: {"celt_only", "wb", 2.5},
					21: {"celt_only", "wb", 5},
					22: {"celt_only", "wb", 10},
					23: {"celt_only", "wb", 20},
					24: {"celt_only", "swb", 2.5},
					25: {"celt_only", "swb", 5},
					26: {"celt_only", "swb", 10},
					27: {"celt_only", "swb", 20},
					28: {"celt_only", "fb", 2.5},
					29: {"celt_only", "fb", 5},
					30: {"celt_only", "fb", 10},
					31: {"celt_only", "fb", 20},
				}
				n := d.FieldU5("config")
				config := configurations[n]
				d.FieldValueStr("mode", config.mode)
				d.FieldValueStr("bandwidth", config.bandwidth)
				d.FieldValueFlt("frame_size", config.frameSize)
			})
			d.FieldBool("stereo")
			d.FieldStruct("frames_per_packet", func(d *decode.D) {
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
				d.FieldValueUint("frames", config.frames)
				d.FieldValueStr("mode", config.mode)
			})
			d.FieldRawLen("data", d.BitsLeft())
		})
	}

	return nil
}
