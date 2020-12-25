package format

import "fq/pkg/decode"

var DefaultRegistry = decode.NewRegistry()

func MustRegister(format *decode.Format) *decode.Format {
	return DefaultRegistry.MustRegister(format)
}

const (
	ALL = "all"

	PROBE = "probe"
	RAW   = "raw"

	// TODO: rename PROBE_* something?
	IMAGE = "image"

	AAC_ADTS           = "aac_adts"
	MPEG_ASC           = "mpeg_asc"
	MPEG_ES            = "mpeg_es"
	MPEG_SPU           = "mpeg_spu"
	MPEG_PES           = "mpeg_pes"
	MPEG_PES_PACKET    = "mpeg_pes_packet"
	AAC_FRAME          = "aac_frame"
	AAC_STREAM         = "aac_stream"
	APEV2              = "apev2"
	ELF                = "elf"
	FLAC               = "flac"
	FLAC_PICTURE       = "flac_picture"
	FLAC_METADATABLOCK = "flac_metadatablock"
	FLAC_FRAME         = "flac_frame"
	FLV                = "flv" // TODO:
	ICC_PROFILE        = "icc_profile"
	ID3_V1             = "id3_v1"
	ID3_V11            = "id3_v11"
	ID3_V2             = "id3_v2"
	JPEG               = "jpeg"
	MKV                = "mkv"
	MP3_FRAME          = "mp3_frame"
	MP3                = "mp3"
	XING_HEADER        = "xing_header"
	MP4                = "mp4"
	OGG                = "ogg"
	OGG_PAGE           = "ogg_page"
	OPUS_PACKET        = "opus_packet"
	PNG                = "png"
	TAR                = "tar"
	TIFF               = "tiff"
	VORBIS_COMMENT     = "vorbis_comment"
	VORBIS_PACKET      = "vorbis_packet"
	VP9_FRAME          = "vp9_frame"
	WAV                = "wav"
)

type FlacMetadatablockStreamInfo struct {
	SampleRate   uint64
	BitPerSample uint64
}

type FlacMetadatablockOut struct {
	StreamInfo FlacMetadatablockStreamInfo
	//            *decode.D
}

type FlacFrameIn struct {
	StreamInfo FlacMetadatablockStreamInfo
}

type FlacFrameOut struct {
	SamplesBuf []byte
}
