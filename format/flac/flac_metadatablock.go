package flac

// TODO: 24 bit picture length truncate warning
// TODO: Cuesheet

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var flacStreaminfoFormat []*decode.Format
var flacPicture []*decode.Format
var vorbisCommentFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.FLAC_METADATABLOCK,
		Description: "FLAC metadatablock",
		DecodeFn:    metadatablockDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_STREAMINFO}, Formats: &flacStreaminfoFormat},
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

func metadatablockDecode(d *decode.D, in interface{}) interface{} {
	var hasStreamInfo bool
	var streamInfo format.FlacStreamInfo

	isLastBlock := d.FieldBool("last_block")
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
		flacStreaminfoOut, ok := d.Format(flacStreaminfoFormat, nil).(format.FlacStreaminfoOut)
		if !ok {
			d.Invalid(fmt.Sprintf("expected FlacStreaminfoOut, got %#+v", flacStreaminfoOut))
		}
		hasStreamInfo = true
		streamInfo = flacStreaminfoOut.StreamInfo
	case MetadataBlockVorbisComment:
		d.FieldFormatLen("comment", int64(length*8), vorbisCommentFormat, nil)
	case MetadataBlockPicture:
		d.FieldFormatLen("picture", int64(length*8), flacPicture, nil)
	case MetadataBlockSeektable:
		seektableCount := length / 18
		d.FieldArrayFn("seekpoints", func(d *decode.D) {
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
	case MetadataBlockApplication:
		d.FieldUTF8("id", 4)
		d.FieldBitBufLen("data", int64((length-4)*8))
	default:
		d.FieldBitBufLen("data", int64(length*8))
	}

	return format.FlacMetadatablockOut{
		IsLastBlock:   isLastBlock,
		HasStreamInfo: hasStreamInfo,
		StreamInfo:    streamInfo,
	}
}
