package flac

// TODO: 24 bit picture length truncate warning

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"fq/pkg/ranges"
)

var flacPicture []*decode.Format
var vorbisCommentFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.FLAC_METADATABLOCK,
		Description: "FLAC metadatablock",
		DecodeFn:    metadatablockDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_PICTURE}, Formats: &flacPicture},
			{Names: []string{format.VORBIS_COMMENT}, Formats: &vorbisCommentFormat},
		},
	})
}

const (
	MetadataBlockStreaminfo    = 0
	MetadataBlockPadding       = 1
	MetadataBlockApplication   = 2
	MetadataBlockSeektable     = 3
	MetadataBlockVorbisComment = 4
	MetadataBlockCuesheet      = 5
	MetadataBlockPicture       = 6
)

var metadataBlockNames = map[uint]string{
	MetadataBlockStreaminfo:    "Streaminfo",
	MetadataBlockPadding:       "Padding",
	MetadataBlockApplication:   "Application",
	MetadataBlockSeektable:     "Seektable",
	MetadataBlockVorbisComment: "Vorbis comment",
	MetadataBlockCuesheet:      "Cuesheet",
	MetadataBlockPicture:       "Picture",
}

func metadatablockDecode(d *decode.D) interface{} {
	mb := &format.FlacMetadatablockOut{}

	mb.LastBlock = d.FieldBool("last_block")
	typ := d.FieldUFn("type", func() (uint64, decode.DisplayFormat, string) {
		t := d.U7()
		name := "Unknown"
		if s, ok := metadataBlockNames[uint(t)]; ok {
			name = s
		}
		return t, decode.NumberDecimal, name
	})
	length := d.FieldU24("length")

	switch typ {
	case MetadataBlockStreaminfo:
		d.FieldU16("minimum_block_size")
		d.FieldU16("maximum_block_size")
		d.FieldU24("minimum_frame_size")
		d.FieldU24("maximum_frame_size")
		sampleRate := d.FieldU("sample_rate", 20)
		// <3> (number of channels)-1. FLAC supports from 1 to 8 channels
		d.FieldUFn("channels", func() (uint64, decode.DisplayFormat, string) { return d.U3() + 1, decode.NumberDecimal, "" })
		// <5> (bits per sample)-1. FLAC supports from 4 to 32 bits per sample. Currently the reference encoder and decoders only support up to 24 bits per sample.
		bitPerSample := d.FieldUFn("bits_per_sample", func() (uint64, decode.DisplayFormat, string) {
			return d.U5() + 1, decode.NumberDecimal, ""
		})
		d.FieldU("total_samples_in_steam", 36)
		md5Range := ranges.Range{Start: d.Pos(), Len: 16 * 8}
		d.FieldBitBufLen("md5", 16*8)

		mb.StreamInfo = &format.FlacMetadatablockStreamInfo{
			SampleRate:   sampleRate,
			BitPerSample: bitPerSample,
			MD5Range:     md5Range,
		}
	case MetadataBlockVorbisComment:
		d.FieldDecodeLen("comment", int64(length*8), vorbisCommentFormat)
	case MetadataBlockPicture:
		d.FieldDecodeLen("picture", int64(length*8), flacPicture)
	case MetadataBlockSeektable:
		seektableCount := length / 18
		d.FieldArrayFn("seekpoint", func(d *decode.D) {
			for i := uint64(0); i < seektableCount; i++ {
				d.FieldStructFn("seekpoint", func(d *decode.D) {
					d.FieldUFn("sample_number", func() (uint64, decode.DisplayFormat, string) {
						n := d.U64()
						d := ""
						if n == 0xffffffffffffffff {
							d = "Placeholder"
						}
						return n, decode.NumberDecimal, d
					})
					d.FieldU64("offset")
					d.FieldU16("number_of_samples")
				})
			}
		})
	default:
		d.FieldBitBufLen("data", int64(length*8))
	}

	return mb
}
