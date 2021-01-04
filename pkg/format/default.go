package format

import (
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/ranges"
)

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
	AAC_FRAME          = "aac_frame"
	AAC_STREAM         = "aac_stream"
	APEV2              = "apev2"
	BZIP2              = "bzip2"
	ELF                = "elf"
	FLAC               = "flac"
	FLAC_FRAME         = "flac_frame"
	FLAC_METADATABLOCK = "flac_metadatablock"
	FLAC_PICTURE       = "flac_picture"
	FLV                = "flv" // TODO:
	GZIP               = "gzip"
	ICC_PROFILE        = "icc_profile"
	ID3_V1             = "id3_v1"
	ID3_V11            = "id3_v11"
	ID3_V2             = "id3_v2"
	JPEG               = "jpeg"
	MKV                = "mkv"
	MP3                = "mp3"
	MP3_FRAME          = "mp3_frame"
	MP3_XING           = "mp3_xing"
	MP4                = "mp4"
	MPEG_ASC           = "mpeg_asc"
	MPEG_ES            = "mpeg_es"
	MPEG_PES           = "mpeg_pes"
	MPEG_PES_PACKET    = "mpeg_pes_packet"
	MPEG_SPU           = "mpeg_spu"
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
	MD5Range     ranges.Range
}

type FlacMetadatablockOut struct {
	IsLastBlock   bool
	HasStreamInfo bool
	StreamInfo    FlacMetadatablockStreamInfo
}

type FlacFrameIn struct {
	StreamInfo FlacMetadatablockStreamInfo
}

type FlacFrameOut struct {
	SamplesBuf []byte
}

type OggPageOut struct {
	IsLastPage         bool
	IsFirstPage        bool
	IsContinuedPacket  bool
	StreamSerialNumber uint32
	SequenceNo         uint32
	Segments           []*bitio.Buffer
}
