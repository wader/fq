package format

import "github.com/wader/fq/pkg/decode"

// based on ffmpeg libavformat/isom.c ff_mp4_obj_type
//nolint:revive
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

var MpegObjectTypeNames = decode.UToStr{
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

var ProtoBufTypeNames = decode.UToStr{
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

//nolint:revive
const (
	MPEGAudioObjectTypeMain      = 1
	MPEGAudioObjectTypeLC        = 2
	MPEGAudioObjectTypeSSR       = 3
	MPEGAudioObjectTypeLTP       = 4
	MPEGAudioObjectTypeSBR       = 5
	MPEGAudioObjectTypeER_AAC_LD = 23
)

var MPEGAudioObjectTypeNames = decode.UToScalar{
	0:                            {Sym: "mull", Description: "Null"},
	MPEGAudioObjectTypeMain:      {Sym: "aac_main", Description: "AAC Main"},
	MPEGAudioObjectTypeLC:        {Sym: "aac_lc", Description: "AAC Low Complexity)"},
	MPEGAudioObjectTypeSSR:       {Sym: "aac_ssr", Description: "AAC Scalable Sample Rate"},
	MPEGAudioObjectTypeLTP:       {Sym: "aac_ltp", Description: "AAC Long Term Prediction"},
	MPEGAudioObjectTypeSBR:       {Sym: "sbr", Description: "Spectral Band Replication"},
	6:                            {Sym: "aac_scalable", Description: "AAC Scalable"},
	7:                            {Sym: "twinvq", Description: "TwinVQ"},
	8:                            {Sym: "celp", Description: "Code Excited Linear Prediction"},
	9:                            {Sym: "hxvc", Description: "Harmonic Vector eXcitation Coding"},
	10:                           {Sym: "reserved", Description: "Reserved"},
	11:                           {Sym: "reserved", Description: "Reserved"},
	12:                           {Sym: "ttsi", Description: "TTSI (Text-To-Speech Interface)"},
	13:                           {Sym: "main_synthesis", Description: "Main Synthesis"},
	14:                           {Sym: "wavetable_synthesis", Description: "Wavetable Synthesis"},
	15:                           {Sym: "general_midi", Description: "General MIDI"},
	16:                           {Sym: "algorithmic", Description: "Algorithmic Synthesis and Audio Effects"},
	17:                           {Sym: "er_aac_lc", Description: "Error Resilient AAC LC"},
	18:                           {Sym: "reserved", Description: "Reserved"},
	19:                           {Sym: "er_aac_ltp", Description: "ER AAC LTP"},
	20:                           {Sym: "er_aac_Scalable", Description: "ER AAC Scalable"},
	21:                           {Sym: "er_twinvq", Description: "ER TwinVQ"},
	22:                           {Sym: "er_bsac", Description: "ER BSAC Bit-Sliced Arithmetic Coding"},
	MPEGAudioObjectTypeER_AAC_LD: {Sym: "er_aac_ld", Description: "ER AAC LD Low Delay"},
	24:                           {Sym: "er_celp", Description: "ER CELP"},
	25:                           {Sym: "er_hvxc", Description: "ER HVXC"},
	26:                           {Sym: "er_hiln", Description: "ER HILN Harmonic and Individual Lines plus Noise"},
	27:                           {Sym: "er_parametric", Description: "ER Parametric"},
	28:                           {Sym: "ssc", Description: "SinuSoidal Coding"},
	29:                           {Sym: "ps", Description: "Parametric Stereo"},
	30:                           {Sym: "mpeg_surround", Description: "MPEG Surround"},
	31:                           {Sym: "(escape value)", Description: "(Escape value)"},
	32:                           {Sym: "layer_1", Description: "MPEG Layer-1"},
	33:                           {Sym: "layer_2", Description: "MPEG Layer-2"},
	34:                           {Sym: "layer_3", Description: "MPEG Layer-3"},
	35:                           {Sym: "dst", Description: "Direct Stream Transfer"},
	36:                           {Sym: "als", Description: "Audio Lossless"},
	37:                           {Sym: "sls", Description: "Scalable Lossless"},
	38:                           {Sym: "sls_non_core", Description: "SLS non-core"},
	39:                           {Sym: "er_aac_eld", Description: "ER AAC ELD Enhanced Low Delay"},
	40:                           {Sym: "smr_simple", Description: "Symbolic Music Representation Simple"},
	41:                           {Sym: "smr_main", Description: "Symbolic Music Representation Main"},
	42:                           {Sym: "usac_no_sbr", Description: "Unified Speech and Audio Coding (no SBR)"},
	43:                           {Sym: "saoc", Description: "Spatial Audio Object Coding"},
	44:                           {Sym: "ld_mpeg_surround", Description: "LD MPEG Surround"},
	45:                           {Sym: "usac", Description: "USAC"},
}
