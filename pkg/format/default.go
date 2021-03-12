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

	APEV2              = "apev2"
	AV1_CCR            = "av1_ccr"
	AV1_FRAME          = "av1_frame"
	AV1_OBU            = "av1_obu"
	BZIP2              = "bzip2"
	DNS                = "dns"
	ELF                = "elf"
	EXIF               = "exif"
	FLAC               = "flac"
	FLAC_FRAME         = "flac_frame"
	FLAC_METADATABLOCK = "flac_metadatablock"
	FLAC_PICTURE       = "flac_picture"
	FLV                = "flv" // TODO:
	GIF                = "gif"
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
	MPEG_AAC_FRAME     = "mpeg_aac_frame"
	MPEG_AAC_STREAM    = "mpeg_aac_stream"
	MPEG_ADTS          = "mpeg_adts"
	MPEG_ASC           = "mpeg_asc"
	MPEG_AVC           = "mpeg_avc"
	MPEG_AVC_DCR       = "mpeg_avc_dcr"
	MPEG_ES            = "mpeg_es"
	MPEG_HEVC          = "mpeg_hevc"
	MPEG_HEVC_DCR      = "mpeg_hevc_dcr"
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
	VP8_FRAME          = "vp8_frame"
	VP9_FRAME          = "vp9_frame"
	WAV                = "wav"
	WEBP               = "webp"
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
	Segments           []*bitio.Buffer // TODO: bitio.Reader (bitio.MultiReader internally?)
}

type AvcIn struct {
	LengthSize uint64
}

type AvcDcrOut struct {
	LengthSize uint64
}

type HevcIn struct {
	LengthSize uint64
}

type HevcDcrOut struct {
	LengthSize uint64
}

// based on ffmpeg libavformat/isom.c ff_mp4_obj_type
const (
	MPEGObjectTypeMOV_TEXT          = 0x08
	MPEGObjectTypeMPEG4             = 0x20
	MPEGObjectTypeH264              = 0x21
	MPEGObjectTypeHEVC              = 0x23
	MPEGObjectTypeAAC               = 0x40
	MPEGObjectTypeMPEG2VideoMain    = 0x61 /* MPEG-2 Main */
	MPEGObjectTypeMPEG2VideoSimple  = 0x60 /* MPEG-2 Simple */
	MPEGObjectTypeMPEG2VideoSNR     = 0x62 /* MPEG-2 SNR */
	MPEGObjectTypeMPEG2VideoSpatial = 0x63 /* MPEG-2 Spatial */
	MPEGObjectTypeMPEG2VideoHigh    = 0x64 /* MPEG-2 High */
	MPEGObjectTypeMPEG2Video422     = 0x65 /* MPEG-2 422 */
	MPEGObjectTypeAACMain           = 0x66 /* MPEG-2 AAC Main */
	MPEGObjectTypeAACLow            = 0x67 /* MPEG-2 AAC Low */
	MPEGObjectTypeAACSSR            = 0x68 /* MPEG-2 AAC SSR */
	MPEGObjectTypeMP32MP3           = 0x69 /* 13818-3 */
	MPEGObjectTypeMPEG1VIDEO        = 0x6A /* 11172-2 */
	MPEGObjectTypeMP3               = 0x6B /* 11172-3 */
	MPEGObjectTypeMJPEG             = 0x6C /* 10918-1 */
	MPEGObjectTypePNG               = 0x6D
	MPEGObjectTypeJPEG2000          = 0x6E /* 15444-1 */
	MPEGObjectTypeVC1               = 0xA3
	MPEGObjectTypeDIRAC             = 0xA4
	MPEGObjectTypeAC3               = 0xA5
	MPEGObjectTypeEAC3              = 0xA6
	MPEGObjectTypeDTS               = 0xA9 /* mp4ra.org */
	MPEGObjectTypeOPUS              = 0xAD /* mp4ra.org */
	MPEGObjectTypeVP9               = 0xB1 /* mp4ra.org */
	MPEGObjectTypeFLAC              = 0xC1 /* nonstandard, update when there is a standard value */
	MPEGObjectTypeTSCC2             = 0xD0 /* nonstandard, camtasia uses it */
	MPEGObjectTypeEVRC              = 0xD1 /* nonstandard, pvAuthor uses it */
	MPEGObjectTypeVORBIS            = 0xDD /* nonstandard, gpac uses it */
	MPEGObjectTypeDVDSubtitle       = 0xE0 /* nonstandard, see unsupported-embedded-subs-2.mp4 */
	MPEGObjectTypeQCELP             = 0xE1
	MPEGObjectTypeMPEG4SYSTEMS1     = 0x01
	MPEGObjectTypeMPEG4SYSTEMS2     = 0x02
	MPEGObjectTypeNONE              = 0
)

var MpegObjectTypeNames = map[uint64]string{
	MPEGObjectTypeMOV_TEXT:          "MPEGObjectTypeMOV_TEXT",
	MPEGObjectTypeMPEG4:             "MPEGObjectTypeMPEG4",
	MPEGObjectTypeH264:              "MPEGObjectTypeH264",
	MPEGObjectTypeHEVC:              "MPEGObjectTypeHEVC",
	MPEGObjectTypeAAC:               "MPEGObjectTypeAAC",
	MPEGObjectTypeMPEG2VideoMain:    "MPEGObjectTypeMPEG2VideoMain",
	MPEGObjectTypeMPEG2VideoSimple:  "MPEGObjectTypeMPEG2VideoSimple",
	MPEGObjectTypeMPEG2VideoSNR:     "MPEGObjectTypeMPEG2VideoSNR",
	MPEGObjectTypeMPEG2VideoSpatial: "MPEGObjectTypeMPEG2VideoSpatial",
	MPEGObjectTypeMPEG2VideoHigh:    "MPEGObjectTypeMPEG2VideoHigh",
	MPEGObjectTypeMPEG2Video422:     "MPEGObjectTypeMPEG2Video422",
	MPEGObjectTypeAACMain:           "MPEGObjectTypeAACMain",
	MPEGObjectTypeAACLow:            "MPEGObjectTypeAACLow",
	MPEGObjectTypeAACSSR:            "MPEGObjectTypeAACSSR",
	MPEGObjectTypeMP32MP3:           "MPEGObjectTypeMP32MP3",
	MPEGObjectTypeMPEG1VIDEO:        "MPEGObjectTypeMPEG1VIDEO",
	MPEGObjectTypeMP3:               "MPEGObjectTypeMP3",
	MPEGObjectTypeMJPEG:             "MPEGObjectTypeMJPEG",
	MPEGObjectTypePNG:               "MPEGObjectTypePNG",
	MPEGObjectTypeJPEG2000:          "MPEGObjectTypeJPEG2000",
	MPEGObjectTypeVC1:               "MPEGObjectTypeVC1",
	MPEGObjectTypeDIRAC:             "MPEGObjectTypeDIRAC",
	MPEGObjectTypeAC3:               "MPEGObjectTypeAC3",
	MPEGObjectTypeEAC3:              "MPEGObjectTypeEAC3",
	MPEGObjectTypeDTS:               "MPEGObjectTypeDTS",
	MPEGObjectTypeOPUS:              "MPEGObjectTypeOPUS",
	MPEGObjectTypeVP9:               "MPEGObjectTypeVP9",
	MPEGObjectTypeFLAC:              "MPEGObjectTypeFLAC",
	MPEGObjectTypeTSCC2:             "MPEGObjectTypeTSCC2",
	MPEGObjectTypeEVRC:              "MPEGObjectTypeEVRC",
	MPEGObjectTypeVORBIS:            "MPEGObjectTypeVORBIS",
	MPEGObjectTypeDVDSubtitle:       "MPEGObjectTypeDVDSubtitle",
	MPEGObjectTypeQCELP:             "MPEGObjectTypeQCELP",
	MPEGObjectTypeMPEG4SYSTEMS1:     "MPEGObjectTypeMPEG4SYSTEMS1",
	MPEGObjectTypeMPEG4SYSTEMS2:     "MPEGObjectTypeMPEG4SYSTEMS2",
	MPEGObjectTypeNONE:              "MPEGObjectTypeNONE",
}

type MpegDecoderConfig struct {
	ObjectType int
}

type MpegEsOut struct {
	DecoderConfigs []MpegDecoderConfig
}
