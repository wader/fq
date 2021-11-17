package flac

// TODO: 24 bit picture length truncate warning
// TODO: Cuesheet

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var flacStreaminfoFormat decode.Group
var flacPicture decode.Group
var vorbisCommentFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.FLAC_METADATABLOCK,
		Description: "FLAC metadatablock",
		DecodeFn:    metadatablockDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_STREAMINFO}, Group: &flacStreaminfoFormat},
			{Names: []string{format.FLAC_PICTURE}, Group: &flacPicture},
			{Names: []string{format.VORBIS_COMMENT}, Group: &vorbisCommentFormat},
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

var metadataBlockNames = decode.UToStr{
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
	typ := d.FieldU7("type", d.MapUToStrSym(metadataBlockNames))
	length := d.FieldU24("length")

	switch typ {
	case MetadataBlockStreaminfo:
		flacStreaminfoOut, ok := d.Format(flacStreaminfoFormat, nil).(format.FlacStreaminfoOut)
		if !ok {
			panic(fmt.Sprintf("expected FlacStreaminfoOut, got %#+v", flacStreaminfoOut))
		}
		hasStreamInfo = true
		streamInfo = flacStreaminfoOut.StreamInfo
	case MetadataBlockVorbisComment:
		d.FieldFormatLen("comment", int64(length*8), vorbisCommentFormat, nil)
	case MetadataBlockPicture:
		d.FieldFormatLen("picture", int64(length*8), flacPicture, nil)
	case MetadataBlockSeektable:
		seektableCount := length / 18
		d.FieldArray("seekpoints", func(d *decode.D) {
			for i := uint64(0); i < seektableCount; i++ {
				d.FieldStruct("seekpoint", func(d *decode.D) {
					d.FieldU64("sample_number", d.MapUToScalar(decode.UToScalar{
						0xffff_ffff_ffff_ffff: {Description: "Placeholder"},
					}))
					d.FieldU64("offset")
					d.FieldU16("number_of_samples")
				})
			}
		})
	case MetadataBlockApplication:
		d.FieldUTF8("id", 4)
		d.FieldRawLen("data", int64((length-4)*8))
	default:
		d.FieldRawLen("data", int64(length*8))
	}

	return format.FlacMetadatablockOut{
		IsLastBlock:   isLastBlock,
		HasStreamInfo: hasStreamInfo,
		StreamInfo:    streamInfo,
	}
}
