package format

import "fq/pkg/decode"

var DefaultRegistry = decode.NewRegistry()

func MustRegister(format *decode.Format) *decode.Format {
	return DefaultRegistry.MustRegister(format)
}

const (
	PROBEABLE = "probeable"
	IMAGE     = "image"

	ADTS           = "adts"
	AAC_FRAME      = "aac_frame"
	AAC_STREAM     = "aac_stream"
	APEV2          = "apev2"
	ELF            = "elf"
	FLAC           = "flac"
	FLAC_PICTURE   = "flac_picture"
	FLV            = "flv" // TODO:
	ICC            = "icc"
	ID3V1          = "id3v1"
	ID3V11         = "id3v11"
	ID3V2          = "id3v2"
	JPEG           = "jpeg"
	MP3_FRAME      = "mp3_frame"
	MP3            = "mp3"
	XING_HEADER    = "xing_header"
	MP4            = "mp4"
	OGG            = "ogg"
	OGG_PAGE       = "ogg_page"
	PNG            = "png"
	TAR            = "tar"
	TIFF           = "tiff"
	VORBIS_COMMENT = "vorbis_comment"
	VORBIS         = "vorbis"
	WAV            = "wav"
)
