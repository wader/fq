package format

import (
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
)

const (
	ALL = "all"

	PROBE = "probe"
	RAW   = "raw"

	// TODO: rename PROBE_* something?
	IMAGE = "image"

	JSON = "json"

	AAC_FRAME           = "aac_frame"
	ADTS                = "adts"
	ADTS_FRAME          = "adts_frame"
	APEV2               = "apev2"
	AV1_CCR             = "av1_ccr"
	AV1_FRAME           = "av1_frame"
	AV1_OBU             = "av1_obu"
	BZIP2               = "bzip2"
	DNS                 = "dns"
	ELF                 = "elf"
	EXIF                = "exif"
	FLAC                = "flac"
	FLAC_FRAME          = "flac_frame"
	FLAC_METADATABLOCKS = "flac_metadatablocks"
	FLAC_PICTURE        = "flac_picture"
	FLV                 = "flv" // TODO:
	GIF                 = "gif"
	GZIP                = "gzip"
	ICC_PROFILE         = "icc_profile"
	ID3V1               = "id3v1"
	ID3V11              = "id3v11"
	ID3V2               = "id3v2"
	JPEG                = "jpeg"
	MATROSKA            = "matroska"
	MP3                 = "mp3"
	MP3_FRAME           = "mp3_frame"
	XING                = "xing"
	MP4                 = "mp4"
	MPEG_ASC            = "mpeg_asc"
	AVC_ANNEXB          = "avc_annexb"
	AVC_DCR             = "avc_dcr"
	AVC_SPS             = "avc_sps"
	AVC_PPS             = "avc_pps"
	AVC_SEI             = "avc_sei"
	AVC_NALU            = "avc_nalu"
	AVC_AU              = "avc_au"
	HEVC_ANNEXB         = "hevc_annexb"
	HEVC_AU             = "hevc_au"
	HEVC_NALU           = "hevc_nalu"
	HEVC_DCR            = "hevc_dcr"
	MPEG_ES             = "mpeg_es"
	MPEG_PES            = "mpeg_pes"
	MPEG_PES_PACKET     = "mpeg_pes_packet"
	MPEG_SPU            = "mpeg_spu"
	MPEG_TS             = "mpeg_ts"
	OGG                 = "ogg"
	OGG_PAGE            = "ogg_page"
	OPUS_PACKET         = "opus_packet"
	PNG                 = "png"
	PROTOBUF            = "protobuf"
	PROTOBUF_WIDEVINE   = "protobuf_widevine"
	PSSH_PLAYREADY      = "pssh_playready"
	TAR                 = "tar"
	TIFF                = "tiff"
	VORBIS_COMMENT      = "vorbis_comment"
	VORBIS_PACKET       = "vorbis_packet"
	VP8_FRAME           = "vp8_frame"
	VP9_FRAME           = "vp9_frame"
	VP9_CFM             = "vp9_cfm"
	VPX_CCR             = "vpx_ccr"
	WAV                 = "wav"
	WEBP                = "webp"
)

type FlacMetadatablockStreamInfo struct {
	SampleRate           uint64
	BitPerSample         uint64
	TotalSamplesInStream uint64
	MD5Range             ranges.Range
}

type FlacMetadatablocksOut struct {
	HasStreamInfo bool
	StreamInfo    FlacMetadatablockStreamInfo
}

type FlacFrameIn struct {
	SamplesBuf []byte
	StreamInfo FlacMetadatablockStreamInfo
}

type FlacFrameOut struct {
	SamplesBuf    []byte
	Samples       uint64
	Channels      int
	BitsPerSample int
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
	MPEGObjectTypeMPEG1VIDEO        = 0x6a /* 11172-2 */
	MPEGObjectTypeMP3               = 0x6b /* 11172-3 */
	MPEGObjectTypeMJPEG             = 0x6c /* 10918-1 */
	MPEGObjectTypePNG               = 0x6d
	MPEGObjectTypeJPEG2000          = 0x6e /* 15444-1 */
	MPEGObjectTypeVC1               = 0xa3
	MPEGObjectTypeDIRAC             = 0xa4
	MPEGObjectTypeAC3               = 0xa5
	MPEGObjectTypeEAC3              = 0xa6
	MPEGObjectTypeDTS               = 0xa9 /* mp4ra.org */
	MPEGObjectTypeOPUS              = 0xad /* mp4ra.org */
	MPEGObjectTypeVP9               = 0xb1 /* mp4ra.org */
	MPEGObjectTypeFLAC              = 0xc1 /* nonstandard, update when there is a standard value */
	MPEGObjectTypeTSCC2             = 0xd0 /* nonstandard, camtasia uses it */
	MPEGObjectTypeEVRC              = 0xd1 /* nonstandard, pvAuthor uses it */
	MPEGObjectTypeVORBIS            = 0xdd /* nonstandard, gpac uses it */
	MPEGObjectTypeDVDSubtitle       = 0xe0 /* nonstandard, see unsupported-embedded-subs-2.mp4 */
	MPEGObjectTypeQCELP             = 0xe1
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

const (
	MPEGStreamTypeUnknown = iota
	MPEGStreamTypeVideo
	MPEGStreamTypeAudio
	MPEGStreamTypeText
)

var MpegObjectTypeStreamType = map[uint64]int{
	MPEGObjectTypeMOV_TEXT:          MPEGStreamTypeText,
	MPEGObjectTypeMPEG4:             MPEGStreamTypeVideo,
	MPEGObjectTypeH264:              MPEGStreamTypeVideo,
	MPEGObjectTypeHEVC:              MPEGStreamTypeVideo,
	MPEGObjectTypeAAC:               MPEGStreamTypeAudio,
	MPEGObjectTypeMPEG2VideoMain:    MPEGStreamTypeVideo,
	MPEGObjectTypeMPEG2VideoSimple:  MPEGStreamTypeVideo,
	MPEGObjectTypeMPEG2VideoSNR:     MPEGStreamTypeVideo,
	MPEGObjectTypeMPEG2VideoSpatial: MPEGStreamTypeVideo,
	MPEGObjectTypeMPEG2VideoHigh:    MPEGStreamTypeVideo,
	MPEGObjectTypeMPEG2Video422:     MPEGStreamTypeVideo,
	MPEGObjectTypeAACMain:           MPEGStreamTypeAudio,
	MPEGObjectTypeAACLow:            MPEGStreamTypeAudio,
	MPEGObjectTypeAACSSR:            MPEGStreamTypeAudio,
	MPEGObjectTypeMP32MP3:           MPEGStreamTypeAudio,
	MPEGObjectTypeMPEG1VIDEO:        MPEGStreamTypeVideo,
	MPEGObjectTypeMP3:               MPEGStreamTypeAudio,
	MPEGObjectTypeMJPEG:             MPEGStreamTypeVideo,
	MPEGObjectTypePNG:               MPEGStreamTypeVideo,
	MPEGObjectTypeJPEG2000:          MPEGStreamTypeVideo,
	MPEGObjectTypeVC1:               MPEGStreamTypeVideo,
	MPEGObjectTypeDIRAC:             MPEGStreamTypeVideo,
	MPEGObjectTypeAC3:               MPEGStreamTypeAudio,
	MPEGObjectTypeEAC3:              MPEGStreamTypeAudio,
	MPEGObjectTypeDTS:               MPEGStreamTypeAudio,
	MPEGObjectTypeOPUS:              MPEGStreamTypeAudio,
	MPEGObjectTypeVP9:               MPEGStreamTypeAudio,
	MPEGObjectTypeFLAC:              MPEGStreamTypeAudio,
	MPEGObjectTypeTSCC2:             MPEGStreamTypeVideo,
	MPEGObjectTypeEVRC:              MPEGStreamTypeAudio,
	MPEGObjectTypeVORBIS:            MPEGStreamTypeAudio,
	MPEGObjectTypeDVDSubtitle:       MPEGStreamTypeText,
	MPEGObjectTypeQCELP:             MPEGStreamTypeAudio,
	MPEGObjectTypeMPEG4SYSTEMS1:     MPEGStreamTypeUnknown,
	MPEGObjectTypeMPEG4SYSTEMS2:     MPEGStreamTypeUnknown,
	MPEGObjectTypeNONE:              MPEGStreamTypeUnknown,
}

type MpegDecoderConfig struct {
	ObjectType    int
	ASCObjectType int
}

type MpegEsOut struct {
	DecoderConfigs []MpegDecoderConfig
}

type ProtoBufType int

const (
	ProtoBufTypeInt32 = iota
	ProtoBufTypeInt64
	ProtoBufTypeUInt32
	ProtoBufTypeUInt64
	ProtoBufTypeSInt32
	ProtoBufTypeSInt64
	ProtoBufTypeBool
	ProtoBufTypeEnum
	ProtoBufTypeFixed64
	ProtoBufTypeSFixed64
	ProtoBufTypeDouble
	ProtoBufTypeString
	ProtoBufTypeBytes
	ProtoBufTypeMessage
	ProtoBufTypePackedRepeated
	ProtoBufTypeFixed32
	ProtoBufTypeSFixed32
	ProtoBufTypeFloat
)

var ProtoBufTypeNames = map[uint64]string{
	ProtoBufTypeInt32:          "Int32",
	ProtoBufTypeInt64:          "Int64",
	ProtoBufTypeUInt32:         "UInt32",
	ProtoBufTypeUInt64:         "UInt64",
	ProtoBufTypeSInt32:         "SInt32",
	ProtoBufTypeSInt64:         "SInt64",
	ProtoBufTypeBool:           "Bool",
	ProtoBufTypeEnum:           "Enum",
	ProtoBufTypeFixed64:        "Fixed64",
	ProtoBufTypeSFixed64:       "SFixed64",
	ProtoBufTypeDouble:         "Double",
	ProtoBufTypeString:         "String",
	ProtoBufTypeBytes:          "Bytes",
	ProtoBufTypeMessage:        "Message",
	ProtoBufTypePackedRepeated: "PackedRepeated",
	ProtoBufTypeFixed32:        "Fixed32",
	ProtoBufTypeSFixed32:       "SFixed32",
	ProtoBufTypeFloat:          "Float",
}

type ProtoBufField struct {
	Type    int
	Name    string
	Message ProtoBufMessage
	Enums   map[uint64]string
}

type ProtoBufMessage map[int]ProtoBufField

type ProtoBufIn struct {
	Message ProtoBufMessage
}

const (
	MPEGAudioObjectTypeMain      = 1
	MPEGAudioObjectTypeLC        = 2
	MPEGAudioObjectTypeSSR       = 3
	MPEGAudioObjectTypeLTP       = 4
	MPEGAudioObjectTypeSBR       = 5
	MPEGAudioObjectTypeER_AAC_LD = 23
)

var MPEGAudioObjectTypeNames = map[uint64]string{
	0:                            "Null",
	MPEGAudioObjectTypeMain:      "AAC Main",
	MPEGAudioObjectTypeLC:        "AAC LC (Low Complexity)",
	MPEGAudioObjectTypeSSR:       "AAC SSR (Scalable Sample Rate)",
	MPEGAudioObjectTypeLTP:       "AAC LTP (Long Term Prediction)",
	MPEGAudioObjectTypeSBR:       "SBR (Spectral Band Replication)",
	6:                            "AAC Scalable",
	7:                            "TwinVQ",
	8:                            "CELP (Code Excited Linear Prediction)",
	9:                            "HXVC (Harmonic Vector eXcitation Coding)",
	10:                           "Reserved",
	11:                           "Reserved",
	12:                           "TTSI (Text-To-Speech Interface)",
	13:                           "Main Synthesis",
	14:                           "Wavetable Synthesis",
	15:                           "General MIDI",
	16:                           "Algorithmic Synthesis and Audio Effects",
	17:                           "ER (Error Resilient) AAC LC",
	18:                           "Reserved",
	19:                           "ER AAC LTP",
	20:                           "ER AAC Scalable",
	21:                           "ER TwinVQ",
	22:                           "ER BSAC (Bit-Sliced Arithmetic Coding)",
	MPEGAudioObjectTypeER_AAC_LD: "ER AAC LD (Low Delay)",
	24:                           "ER CELP",
	25:                           "ER HVXC",
	26:                           "ER HILN (Harmonic and Individual Lines plus Noise)",
	27:                           "ER Parametric",
	28:                           "SSC (SinuSoidal Coding)",
	29:                           "PS (Parametric Stereo)",
	30:                           "MPEG Surround",
	31:                           "(Escape value)",
	32:                           "Layer-1",
	33:                           "Layer-2",
	34:                           "Layer-3",
	35:                           "DST (Direct Stream Transfer)",
	36:                           "ALS (Audio Lossless)",
	37:                           "SLS (Scalable LosslesS)",
	38:                           "SLS non-core",
	39:                           "ER AAC ELD (Enhanced Low Delay)",
	40:                           "SMR (Symbolic Music Representation) Simple",
	41:                           "SMR Main",
	42:                           "USAC (Unified Speech and Audio Coding) (no SBR)",
	43:                           "SAOC (Spatial Audio Object Coding)",
	44:                           "LD MPEG Surround",
	45:                           "USAC",
}

type MPEGASCOut struct {
	ObjectType int
}

type AACFrameIn struct {
	ObjectType int
}
