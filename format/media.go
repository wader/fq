package format

import (
	"github.com/wader/fq/pkg/scalar"
)

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

var MpegObjectTypeNames = scalar.UintMapSymStr{
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

const (
	MPEGAudioObjectTypeMain      = 1
	MPEGAudioObjectTypeLC        = 2
	MPEGAudioObjectTypeSSR       = 3
	MPEGAudioObjectTypeLTP       = 4
	MPEGAudioObjectTypeSBR       = 5
	MPEGAudioObjectTypeER_AAC_LD = 23
	MPEGAudioObjectTypePS        = 29
)

var MPEGAudioObjectTypeNames = scalar.UintMap{
	0:                            {Sym: "mull", Description: "Null"},
	MPEGAudioObjectTypeMain:      {Sym: "aac_main", Description: "AAC Main"},
	MPEGAudioObjectTypeLC:        {Sym: "aac_lc", Description: "AAC Low Complexity"},
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
	MPEGAudioObjectTypePS:        {Sym: "ps", Description: "Parametric Stereo"},
	30:                           {Sym: "mpeg_surround", Description: "MPEG Surround"},
	31:                           {Description: "(Escape value)"},
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

// based on ffmpeg/libavutil/pixfmt.h
var ISO_23091_2_ColourPrimariesMap = scalar.UintMap{
	1:  {Sym: "bt709", Description: "ITU-R BT1361 / IEC 61966-2-4 / SMPTE RP 177 Annex B"},
	2:  {Sym: "unspecified", Description: "Unspecified"},
	3:  {Sym: "reserved", Description: "Reserved"},
	4:  {Sym: "bt470m", Description: "FCC Title 47 Code of Federal Regulations 73.682 (a)(20)"},
	5:  {Sym: "bt470bg", Description: "ITU-R BT601-6 625 / ITU-R BT1358 625 / ITU-R BT1700 625 PAL & SECAM"},
	6:  {Sym: "smpte170m", Description: "ITU-R BT601-6 525 / ITU-R BT1358 525 / ITU-R BT1700 NTSC"},
	7:  {Sym: "smpte240m", Description: "ITU-R BT601-6 525 / ITU-R BT1358 525 / ITU-R BT1700 NTSC"},
	8:  {Sym: "film", Description: "Illuminant C"},
	9:  {Sym: "bt2020", Description: "ITU-R BT2020"},
	10: {Sym: "smpte428", Description: "SMPTE ST 428-1 (CIE 1931 XYZ)"},
	11: {Sym: "smpte431", Description: "SMPTE ST 431-2 (2011) / DCI P3"},
	12: {Sym: "smpte432", Description: "SMPTE ST 432-1 (2010) / P3 D65 / Display P3"},
	22: {Sym: "ebu3213", Description: "EBU Tech. 3213-E (nothing there) / one of JEDEC P22 group phosphors"},
}

var ISO_23091_2_TransferCharacteristicMap = scalar.UintMap{
	1:  {Sym: "bt709", Description: "ITU-R BT1361"},
	2:  {Sym: "unspecified", Description: "Unspecified"},
	3:  {Sym: "reserved", Description: "Reserved"},
	4:  {Sym: "gamma22", Description: "ITU-R BT470M / ITU-R BT1700 625 PAL & SECAM"},
	5:  {Sym: "gamma28", Description: "ITU-R BT470BG"},
	6:  {Sym: "smpte170m", Description: "ITU-R BT601-6 525 or 625 / ITU-R BT1358 525 or 625 / ITU-R BT1700 NTSC"},
	7:  {Sym: "smpte240m"},
	8:  {Sym: "linear", Description: "Linear transfer characteristics"},
	9:  {Sym: "log", Description: "Logarithmic transfer characteristic (100:1 range)"},
	10: {Sym: "log_sqrt", Description: "Logarithmic transfer characteristic (100 * Sqrt(10) : 1 range)"},
	11: {Sym: "iec61966_2_4", Description: "IEC 61966-2-4"},
	12: {Sym: "bt1361_ecg", Description: "ITU-R BT1361 Extended Colour Gamut"},
	13: {Sym: "iec61966_2_1", Description: "IEC 61966-2-1 (sRGB or sYCC)"},
	14: {Sym: "bt2020_10", Description: "ITU-R BT2020 for 10-bit system"},
	15: {Sym: "bt2020_12", Description: "ITU-R BT2020 for 12-bit system"},
	16: {Sym: "smpte2084", Description: "SMPTE ST 2084 for 10-, 12-, 14- and 16-bit systems"},
	17: {Sym: "smpte428", Description: "SMPTE ST 428-1"},
	18: {Sym: "arib_std_b67", Description: "ARIB STD-B67, known as Hybrid log-gamma"},
}

var ISO_23091_2_MatrixCoefficients = scalar.UintMap{
	0:  {Sym: "rgb", Description: "GBR, IEC 61966-2-1 (sRGB), YZX and ST 428-1"},
	1:  {Sym: "bt709", Description: "ITU-R BT1361 / IEC 61966-2-4 xvYCC709 / derived in SMPTE RP 177 Annex B"},
	2:  {Sym: "unspecified", Description: "Unspecified"},
	3:  {Sym: "reserved", Description: "Reserved"},
	4:  {Sym: "fcc", Description: "FCC Title 47 Code of Federal Regulations 73.682 (a)(20)"},
	5:  {Sym: "bt470bg", Description: "ITU-R BT601-6 625 / ITU-R BT1358 625 / ITU-R BT1700 625 PAL & SECAM / IEC 61966-2-4 xvYCC601"},
	6:  {Sym: "smpte170m", Description: "ITU-R BT601-6 525 / ITU-R BT1358 525 / ITU-R BT1700 NTSC"},
	7:  {Sym: "smpte240m", Description: "Derived from 170M primaries and D65 white point"},
	8:  {Sym: "ycgco", Description: "VC-2 and H.264 FRext"},
	9:  {Sym: "bt2020_ncl", Description: "ITU-R BT2020 non-constant luminance system"},
	10: {Sym: "bt2020_cl", Description: "ITU-R BT2020 constant luminance system"},
	11: {Sym: "smpte2085", Description: "SMPTE 2085, Y'D'zD'x"},
	12: {Sym: "chroma_derived_ncl", Description: "Chromaticity-derived non-constant luminance system"},
	13: {Sym: "chroma_derived_cl", Description: "Chromaticity-derived constant luminance system"},
	14: {Sym: "ictcp", Description: "ITU-R BT.2100-0, ICtCp"},
}

// based on ffmpeg libavformat/riff.c
const (
	WAVTagPCM             = 0x0001
	WAVTagADPCM_MS        = 0x0002
	WAVTagPCM_Float       = 0x0003
	WAVTagPCM_ALAW        = 0x0006
	WAVTagPCM_MULAW       = 0x0007
	WAVTagWMAVOICE        = 0x000a
	WAVTagADPCM_IMA_OKI   = 0x0010
	WAVTagADPCM_IMA_WAV   = 0x0011
	WAVTagADPCM_IMA_OKI_2 = 0x0017
	WAVTagADPCM_YAMAHA    = 0x0020
	WAVTagTRUESPEECH      = 0x0022
	WAVTagGSM_MS          = 0x0031
	WAVTagGSM_MS_2        = 0x0032
	WAVTagAMR_NB          = 0x0038
	WAVTagG723_1          = 0x0042
	WAVTagADPCM_G726      = 0x0045
	WAVTagADPCM_G726_2    = 0x0014
	WAVTagADPCM_G726_3    = 0x0040
	WAVTagMP2             = 0x0050
	WAVTagMP3             = 0x0055
	WAVTagAMR_NB_2        = 0x0057
	WAVTagAMR_WB          = 0x0058
	WAVTagADPCM_IMA_DK4   = 0x0061
	WAVTagADPCM_IMA_DK3   = 0x0062
	WAVTagADPCM_G726_4    = 0x0064
	WAVTagADPCM_IMA_WAV_2 = 0x0069
	WAVTagMETASOUND       = 0x0075
	WAVTagG729            = 0x0083
	WAVTagAAC             = 0x00ff
	WAVTagG723_1_2        = 0x0111
	WAVTagSIPR            = 0x0130
	WAVTagACELP_KELVIN    = 0x0135
	WAVTagWMAV1           = 0x0160
	WAVTagWMAV2           = 0x0161
	WAVTagWMAPRO          = 0x0162
	WAVTagWMALOSSLESS     = 0x0163
	WAVTagXMA1            = 0x0165
	WAVTagXMA2            = 0x0166
	WAVTagFTR             = 0x0180
	WAVTagADPCM_CT        = 0x0200
	WAVTagDVAUDIO         = 0x0215
	WAVTagDVAUDIO_2       = 0x0216
	WAVTagATRAC3          = 0x0270
	WAVTagMSNSIREN        = 0x028e
	WAVTagADPCM_G722      = 0x028f
	WAVTagMISC4           = 0x0350
	WAVTagIMC             = 0x0401
	WAVTagIAC             = 0x0402
	WAVTagON2AVC          = 0x0500
	WAVTagON2AVC_2        = 0x0501
	WAVTagGSM_MS_3        = 0x1500
	WAVTagTRUESPEECH_2    = 0x1501
	WAVTagAAC_2           = 0x1600
	WAVTagAAC_LATM        = 0x1602
	WAVTagAC3             = 0x2000
	WAVTagDTS             = 0x2001
	WAVTagSONIC           = 0x2048
	WAVTagG729_2          = 0x2222
	WAVTagPCM_MULAW_2     = 0x6c75
	WAVTagAAC_3           = 0x706d
	WAVTagAAC_4           = 0x4143
	WAVTagFTR_2           = 0x4180
	WAVTagXAN_DPCM        = 0x594a
	WAVTagG729_3          = 0x729a
	WAVTagFTR_3           = 0x8180
	WAVTagG723_1_3        = 0xa100
	WAVTagAAC_5           = 0xa106
	WAVTagSPEEX           = 0xa109
	WAVTagFLAC            = 0xf1ac
	WAVTagFORMATEX        = 0xfffe
	WAVTagADPCM_SWF       = 0x5356
	WAVTagVORBIS          = 0x566f
)

var WAVTagNames = scalar.UintMapSymStr{
	WAVTagPCM:             "pcm",
	WAVTagADPCM_MS:        "adpcm_ms",
	WAVTagPCM_Float:       "pcm_float",
	WAVTagPCM_ALAW:        "pcm_alaw",
	WAVTagPCM_MULAW:       "pcm_mulaw",
	WAVTagWMAVOICE:        "wmavoice",
	WAVTagADPCM_IMA_OKI:   "adpcm_ima_oki",
	WAVTagADPCM_IMA_WAV:   "adpcm_ima_wav",
	WAVTagADPCM_IMA_OKI_2: "adpcm_ima_oki_2",
	WAVTagADPCM_YAMAHA:    "adpcm_yamaha",
	WAVTagTRUESPEECH:      "truespeech",
	WAVTagGSM_MS:          "gsm_ms",
	WAVTagGSM_MS_2:        "gsm_ms_2",
	WAVTagAMR_NB:          "amr_nb",
	WAVTagG723_1:          "g723_1",
	WAVTagADPCM_G726:      "adpcm_g726",
	WAVTagADPCM_G726_2:    "adpcm_g726_2",
	WAVTagADPCM_G726_3:    "adpcm_g726_3",
	WAVTagMP2:             "mp2",
	WAVTagMP3:             "mp3",
	WAVTagAMR_NB_2:        "amr_nb_2",
	WAVTagAMR_WB:          "amr_wb",
	WAVTagADPCM_IMA_DK4:   "adpcm_ima_dk4",
	WAVTagADPCM_IMA_DK3:   "adpcm_ima_dk3",
	WAVTagADPCM_G726_4:    "adpcm_g726_4",
	WAVTagADPCM_IMA_WAV_2: "adpcm_ima_wav_2",
	WAVTagMETASOUND:       "metasound",
	WAVTagG729:            "g729",
	WAVTagAAC:             "aac",
	WAVTagG723_1_2:        "g723_1_2",
	WAVTagSIPR:            "sipr",
	WAVTagACELP_KELVIN:    "acelp_kelvin",
	WAVTagWMAV1:           "wmav1",
	WAVTagWMAV2:           "wmav2",
	WAVTagWMAPRO:          "wmapro",
	WAVTagWMALOSSLESS:     "wmalossless",
	WAVTagXMA1:            "xma1",
	WAVTagXMA2:            "xma2",
	WAVTagFTR:             "ftr",
	WAVTagADPCM_CT:        "adpcm_ct",
	WAVTagDVAUDIO:         "dvaudio",
	WAVTagDVAUDIO_2:       "dvaudio_2",
	WAVTagATRAC3:          "atrac3",
	WAVTagMSNSIREN:        "msnsiren",
	WAVTagADPCM_G722:      "adpcm_g722",
	WAVTagMISC4:           "misc4",
	WAVTagIMC:             "imc",
	WAVTagIAC:             "iac",
	WAVTagON2AVC:          "on2avc",
	WAVTagON2AVC_2:        "on2avc_2",
	WAVTagGSM_MS_3:        "gsm_ms_3",
	WAVTagTRUESPEECH_2:    "truespeech_2",
	WAVTagAAC_2:           "aac_2",
	WAVTagAAC_LATM:        "aac_latm",
	WAVTagAC3:             "ac3",
	WAVTagDTS:             "dts",
	WAVTagSONIC:           "sonic",
	WAVTagG729_2:          "g729_2",
	WAVTagPCM_MULAW_2:     "pcm_mulaw_2",
	WAVTagAAC_3:           "aac_3",
	WAVTagAAC_4:           "aac_4",
	WAVTagFTR_2:           "ftr_2",
	WAVTagXAN_DPCM:        "xan_dpcm",
	WAVTagG729_3:          "g729_3",
	WAVTagFTR_3:           "ftr_3",
	WAVTagG723_1_3:        "g723_1_3",
	WAVTagAAC_5:           "aac_5",
	WAVTagSPEEX:           "speex",
	WAVTagFLAC:            "flac",
	WAVTagFORMATEX:        "formatex",
	WAVTagADPCM_SWF:       "adpcm_swf",
	WAVTagVORBIS:          "vorbis",
}

// based on ffmpeg libavformat/riff.c
const (
	BMPTagH264                = "H264"
	BMPTagH264_h264           = "h264"
	BMPTagH264_X264           = "X264"
	BMPTagH264_x264           = "x264"
	BMPTagH264_avc1           = "avc1"
	BMPTagH264_DAVC           = "DAVC"
	BMPTagH264_SMV2           = "SMV2"
	BMPTagH264_VSSH           = "VSSH"
	BMPTagH264_Q264           = "Q264" // QNAP surveillance system
	BMPTagH264_V264           = "V264" // CCTV recordings
	BMPTagH264_GAVC           = "GAVC" // GeoVision camera
	BMPTagH264_UMSV           = "UMSV"
	BMPTagH264_tshd           = "tshd"
	BMPTagH264_INMC           = "INMC"
	BMPTagH263                = "H263"
	BMPTagH263_X263           = "X263"
	BMPTagH263_T263           = "T263"
	BMPTagH263_L263           = "L263"
	BMPTagH263_VX1K           = "VX1K"
	BMPTagH263_ZyGo           = "ZyGo"
	BMPTagH263_M263           = "M263"
	BMPTagH263_lsvm           = "lsvm"
	BMPTagH263P               = "H263"
	BMPTagH263I               = "I263" // Intel H.263
	BMPTagH261                = "H261"
	BMPTagH263_U263           = "U263"
	BMPTagH263_VSM4           = "VSM4" // needs -vf il=l=i:c=i
	BMPTagMPEG4               = "FMP4"
	BMPTagMPEG4_DIVX          = "DIVX"
	BMPTagMPEG4_DX50          = "DX50"
	BMPTagMPEG4_XVID          = "XVID"
	BMPTagMPEG4_MP4S          = "MP4S"
	BMPTagMPEG4_M4S2          = "M4S2"             // some broken AVIs use this
	BMPTagMPEG4_04000000      = "\x04\x00\x00\x00" // some broken AVIs use this
	BMPTagMPEG4_ZMP4          = "ZMP4"
	BMPTagMPEG4_DIV1          = "DIV1"
	BMPTagMPEG4_BLZ0          = "BLZ0"
	BMPTagMPEG4_mp4v          = "mp4v"
	BMPTagMPEG4_UMP4          = "UMP4"
	BMPTagMPEG4_WV1F          = "WV1F"
	BMPTagMPEG4_SEDG          = "SEDG"
	BMPTagMPEG4_RMP4          = "RMP4"
	BMPTagMPEG4_3IV2          = "3IV2" // WaWv MPEG-4 Video Codec
	BMPTagMPEG4_WAWV          = "WAWV"
	BMPTagMPEG4_FFDS          = "FFDS"
	BMPTagMPEG4_FVFW          = "FVFW"
	BMPTagMPEG4_DCOD          = "DCOD"
	BMPTagMPEG4_MVXM          = "MVXM"
	BMPTagMPEG4_PM4V          = "PM4V"
	BMPTagMPEG4_SMP4          = "SMP4"
	BMPTagMPEG4_DXGM          = "DXGM"
	BMPTagMPEG4_VIDM          = "VIDM"
	BMPTagMPEG4_M4T3          = "M4T3"
	BMPTagMPEG4_GEOX          = "GEOX" // flipped video
	BMPTagMPEG4_G264          = "G264" // flipped video
	BMPTagMPEG4_HDX4          = "HDX4"
	BMPTagMPEG4_DM4V          = "DM4V"
	BMPTagMPEG4_DMK2          = "DMK2"
	BMPTagMPEG4_DYM4          = "DYM4"
	BMPTagMPEG4_DIGI          = "DIGI" // Ephv MPEG-4
	BMPTagMPEG4_EPHV          = "EPHV"
	BMPTagMPEG4_EM4A          = "EM4A" // Divio MPEG-4
	BMPTagMPEG4_M4CC          = "M4CC"
	BMPTagMPEG4_SN40          = "SN40"
	BMPTagMPEG4_VSPX          = "VSPX"
	BMPTagMPEG4_ULDX          = "ULDX"
	BMPTagMPEG4_GEOV          = "GEOV" // Samsung SHR-6040
	BMPTagMPEG4_SIPP          = "SIPP"
	BMPTagMPEG4_SM4V          = "SM4V"
	BMPTagMPEG4_XVIX          = "XVIX"
	BMPTagMPEG4_DreX          = "DreX"
	BMPTagMPEG4_QMP4          = "QMP4" // QNAP Systems
	BMPTagMPEG4_PLV1          = "PLV1" // Pelco DVR MPEG-4
	BMPTagMPEG4_GLV4          = "GLV4"
	BMPTagMPEG4_GMP4          = "GMP4" // GeoVision camera
	BMPTagMPEG4_MNM4          = "MNM4" // March Networks DVR
	BMPTagMPEG4_GTM4          = "GTM4" // Telefactor
	BMPTagMSMPEG4V3           = "MP43"
	BMPTagMSMPEG4V3_DIV3      = "DIV3"
	BMPTagMSMPEG4V3_MPG3      = "MPG3"
	BMPTagMSMPEG4V3_DIV5      = "DIV5"
	BMPTagMSMPEG4V3_DIV6      = "DIV6"
	BMPTagMSMPEG4V3_DIV4      = "DIV4"
	BMPTagMSMPEG4V3_DVX3      = "DVX3"
	BMPTagMSMPEG4V3_AP41      = "AP41"
	BMPTagMSMPEG4V3_COL1      = "COL1"
	BMPTagMSMPEG4V3_COL0      = "COL0"
	BMPTagMSMPEG4V2           = "MP42"
	BMPTagMSMPEG4V2_DIV2      = "DIV2"
	BMPTagMSMPEG4V1           = "MPG4"
	BMPTagMSMPEG4V1_MP41      = "MP41"
	BMPTagWMV1                = "WMV1"
	BMPTagWMV2                = "WMV2"
	BMPTagWMV2_GXVE           = "GXVE"
	BMPTagDVVIDEO             = "dvsd"
	BMPTagDVVIDEO_dvhd        = "dvhd"
	BMPTagDVVIDEO_dvh1        = "dvh1"
	BMPTagDVVIDEO_dvsl        = "dvsl"
	BMPTagDVVIDEO_dv25        = "dv25"
	BMPTagDVVIDEO_dv50        = "dv50" // Canopus DV
	BMPTagDVVIDEO_cdvc        = "cdvc" // Canopus DV
	BMPTagDVVIDEO_CDVH        = "CDVH" // Canopus DV
	BMPTagDVVIDEO_CDV5        = "CDV5"
	BMPTagDVVIDEO_dvc         = "dvc "
	BMPTagDVVIDEO_dvcs        = "dvcs"
	BMPTagDVVIDEO_dvis        = "dvis"
	BMPTagDVVIDEO_pdvc        = "pdvc"
	BMPTagDVVIDEO_SL25        = "SL25"
	BMPTagDVVIDEO_SLDV        = "SLDV"
	BMPTagDVVIDEO_AVd1        = "AVd1"
	BMPTagMPEG1VIDEO          = "mpg1"
	BMPTagMPEG2VIDEO_mpg2     = "mpg2"
	BMPTagMPEG2VIDEO_MPEG     = "MPEG"
	BMPTagMPEG1VIDEO_PIM1     = "PIM1"
	BMPTagMPEG2VIDEO_PIM2     = "PIM2"
	BMPTagMPEG1VIDEO_VCR2     = "VCR2"
	BMPTagMPEG1VIDEO_01000016 = "\x01\x00\x00\x16"
	BMPTagMPEG2VIDEO_02000016 = "\x02\x00\x00\x16"
	BMPTagMPEG4_04000016      = "\x04\x00\x00\x16"
	BMPTagMPEG2VIDEO          = "DVR "
	BMPTagMPEG2VIDEOMMES      = "MMES" // Lead MPEG-2 in AVI
	BMPTagMPEG2VIDEOLMP2      = "LMP2"
	BMPTagMPEG2VIDEOslif      = "slif"
	BMPTagMPEG2VIDEOEM2V      = "EM2V" // Matrox MPEG-2 intra-only
	BMPTagMPEG2VIDEOM701      = "M701"
	BMPTagMPEG2VIDEOM702      = "M702"
	BMPTagMPEG2VIDEOM703      = "M703"
	BMPTagMPEG2VIDEOM704      = "M704"
	BMPTagMPEG2VIDEOM705      = "M705"
	BMPTagMPEG2VIDEOmpgv      = "mpgv"
	BMPTagMPEG1VIDEO_BW10     = "BW10"
	BMPTagMPEG1VIDEO_XMPG     = "XMPG" // Xing MPEG intra only
	BMPTagMJPEG               = "MJPG"
	BMPTagMJPEG_MSC2          = "MSC2" // Multiscope II
	BMPTagMJPEG_LJPG          = "LJPG"
	BMPTagMJPEG_dmb1          = "dmb1"
	BMPTagMJPEG_mjpa          = "mjpa"
	BMPTagMJPEG_JR24          = "JR24" // Quadrox Mjpeg
	BMPTagLJPEG               = "LJPG" // Pegasus lossless JPEG
	BMPTagMJPEG_JPGL          = "JPGL" // JPEG-LS custom FOURCC for AVI - encoder
	BMPTagJPEGLS              = "MJLS"
	BMPTagJPEGLS_MJPG         = "MJPG" // JPEG-LS custom FOURCC for AVI - decoder
	BMPTagMJPEG_MJLS          = "MJLS"
	BMPTagMJPEG_jpeg          = "jpeg"
	BMPTagMJPEG_IJPG          = "IJPG"
	BMPTagAVRN                = "AVRn"
	BMPTagMJPEG_ACDV          = "ACDV"
	BMPTagMJPEG_QIVG          = "QIVG" // SL M-JPEG
	BMPTagMJPEG_SLMJ          = "SLMJ" // Creative Webcam JPEG
	BMPTagMJPEG_CJPG          = "CJPG" // Intel JPEG Library Video Codec
	BMPTagMJPEG_IJLV          = "IJLV" // Midvid JPEG Video Codec
	BMPTagMJPEG_MVJP          = "MVJP"
	BMPTagMJPEG_AVI1          = "AVI1"
	BMPTagMJPEG_AVI2          = "AVI2"
	BMPTagMJPEG_MTSJ          = "MTSJ" // Paradigm Matrix M-JPEG Codec
	BMPTagMJPEG_ZJPG          = "ZJPG"
	BMPTagMJPEG_MMJP          = "MMJP"
	BMPTagHUFFYUV             = "HFYU"
	BMPTagFFVHUFF             = "FFVH"
	BMPTagCYUV                = "CYUV"
	BMPTagRAWVIDEO_00000000   = "\x00\x00\x00\x00"
	BMPTagRAWVIDEO_03000000   = "\x03\x00\x00\x00"
	BMPTagRAWVIDEO_I420       = "I420"
	BMPTagRAWVIDEO_YUY2       = "YUY2"
	BMPTagRAWVIDEO_Y422       = "Y422"
	BMPTagRAWVIDEO_V422       = "V422"
	BMPTagRAWVIDEO_YUNV       = "YUNV"
	BMPTagRAWVIDEO_UYNV       = "UYNV"
	BMPTagRAWVIDEO_UYNY       = "UYNY"
	BMPTagRAWVIDEO_uyv1       = "uyv1"
	BMPTagRAWVIDEO_2Vu1       = "2Vu1"
	BMPTagRAWVIDEO_2vuy       = "2vuy"
	BMPTagRAWVIDEO_yuvs       = "yuvs"
	BMPTagRAWVIDEO_yuv2       = "yuv2"
	BMPTagRAWVIDEO_P422       = "P422"
	BMPTagRAWVIDEO_YV12       = "YV12"
	BMPTagRAWVIDEO_YV16       = "YV16"
	BMPTagRAWVIDEO_YV24       = "YV24"
	BMPTagRAWVIDEO_UYVY       = "UYVY"
	BMPTagRAWVIDEO_VYUY       = "VYUY"
	BMPTagRAWVIDEO_IYUV       = "IYUV"
	BMPTagRAWVIDEO_AYUV       = "AYUV"
	BMPTagRAWVIDEO_Y800       = "Y800"
	BMPTagRAWVIDEO_Y8         = "Y8  "
	BMPTagRAWVIDEO_HDYC       = "HDYC"
	BMPTagRAWVIDEO_VDTZ       = "VDTZ"
	BMPTagRAWVIDEO_Y411       = "Y411"
	BMPTagRAWVIDEO_NV12       = "NV12"
	BMPTagRAWVIDEO_NV21       = "NV21"
	BMPTagRAWVIDEO_Y41B       = "Y41B"
	BMPTagRAWVIDEO_Y42B       = "Y42B"
	BMPTagRAWVIDEO_YUV9       = "YUV9"
	BMPTagRAWVIDEO_YVU9       = "YVU9"
	BMPTagRAWVIDEO_auv2       = "auv2"
	BMPTagRAWVIDEO_YVYU       = "YVYU"
	BMPTagRAWVIDEO_YUYV       = "YUYV"
	BMPTagRAWVIDEO_I410       = "I410"
	BMPTagRAWVIDEO_I411       = "I411"
	BMPTagRAWVIDEO_I422       = "I422"
	BMPTagRAWVIDEO_I440       = "I440"
	BMPTagRAWVIDEO_I444       = "I444"
	BMPTagRAWVIDEO_J420       = "J420"
	BMPTagRAWVIDEO_J422       = "J422"
	BMPTagRAWVIDEO_J440       = "J440"
	BMPTagRAWVIDEO_J444       = "J444"
	BMPTagRAWVIDEO_YUVA       = "YUVA"
	BMPTagRAWVIDEO_I40A       = "I40A"
	BMPTagRAWVIDEO_I42A       = "I42A"
	BMPTagRAWVIDEO_RGB2       = "RGB2"
	BMPTagRAWVIDEO_RV15       = "RV15"
	BMPTagRAWVIDEO_RV16       = "RV16"
	BMPTagRAWVIDEO_RV24       = "RV24"
	BMPTagRAWVIDEO_RV32       = "RV32"
	BMPTagRAWVIDEO_RGBA       = "RGBA"
	BMPTagRAWVIDEO_AV32       = "AV32"
	BMPTagRAWVIDEO_GREY       = "GREY"
	BMPTagRAWVIDEO_I09L       = "I09L"
	BMPTagRAWVIDEO_I09B       = "I09B"
	BMPTagRAWVIDEO_I29L       = "I29L"
	BMPTagRAWVIDEO_I29B       = "I29B"
	BMPTagRAWVIDEO_I49L       = "I49L"
	BMPTagRAWVIDEO_I49B       = "I49B"
	BMPTagRAWVIDEO_I0AL       = "I0AL"
	BMPTagRAWVIDEO_I0AB       = "I0AB"
	BMPTagRAWVIDEO_I2AL       = "I2AL"
	BMPTagRAWVIDEO_I2AB       = "I2AB"
	BMPTagRAWVIDEO_I4AL       = "I4AL"
	BMPTagRAWVIDEO_I4AB       = "I4AB"
	BMPTagRAWVIDEO_I4FL       = "I4FL"
	BMPTagRAWVIDEO_I4FB       = "I4FB"
	BMPTagRAWVIDEO_I0CL       = "I0CL"
	BMPTagRAWVIDEO_I0CB       = "I0CB"
	BMPTagRAWVIDEO_I2CL       = "I2CL"
	BMPTagRAWVIDEO_I2CB       = "I2CB"
	BMPTagRAWVIDEO_I4CL       = "I4CL"
	BMPTagRAWVIDEO_I4CB       = "I4CB"
	BMPTagRAWVIDEO_I0FL       = "I0FL"
	BMPTagRAWVIDEO_I0FB       = "I0FB"
	BMPTagFRWU                = "FRWU"
	BMPTagR10K                = "R10k"
	BMPTagR210                = "r210"
	BMPTagV210_v210           = "v210"
	BMPTagV210_C210           = "C210"
	BMPTagV308                = "v308"
	BMPTagV408                = "v408"
	BMPTagV410                = "v410"
	BMPTagYUV4                = "yuv4"
	BMPTagINDEO3_IV31         = "IV31"
	BMPTagINDEO3_IV32         = "IV32"
	BMPTagINDEO4              = "IV41"
	BMPTagINDEO5              = "IV50"
	BMPTagVP3_VP31            = "VP31"
	BMPTagVP3_VP30            = "VP30"
	BMPTagVP4                 = "VP40"
	BMPTagVP5                 = "VP50"
	BMPTagVP6_VP60            = "VP60"
	BMPTagVP6_VP61            = "VP61"
	BMPTagVP6_VP62            = "VP62"
	BMPTagVP6A                = "VP6A"
	BMPTagVP6F_VP6F           = "VP6F"
	BMPTagVP6F_FLV4           = "FLV4"
	BMPTagVP7_VP70            = "VP70"
	BMPTagVP7_VP71            = "VP71"
	BMPTagVP8                 = "VP80"
	BMPTagVP9                 = "VP90"
	BMPTagASV1                = "ASV1"
	BMPTagASV2                = "ASV2"
	BMPTagVCR1                = "VCR1"
	BMPTagFFV1                = "FFV1"
	BMPTagXAN_WC4             = "Xxan"
	BMPTagMIMIC               = "LM20"
	BMPTagMSRLE_mrle          = "mrle"
	BMPTagMSRLE_01000000      = "\x01\x00\x00\x00"
	BMPTagMSRLE_02000000      = "\x02\x00\x00\x00"
	BMPTagMSVIDEO1_MSVC       = "MSVC"
	BMPTagMSVIDEO1_msvc       = "msvc"
	BMPTagMSVIDEO1_CRAM       = "CRAM"
	BMPTagMSVIDEO1_cram       = "cram"
	BMPTagMSVIDEO1_WHAM       = "WHAM"
	BMPTagMSVIDEO1_wham       = "wham"
	BMPTagCINEPAK             = "cvid"
	BMPTagTRUEMOTION1_DUCK    = "DUCK"
	BMPTagTRUEMOTION1_PVEZ    = "PVEZ"
	BMPTagMSZH                = "MSZH"
	BMPTagZLIB                = "ZLIB"
	BMPTagSNOW                = "SNOW"
	BMPTag4XM                 = "4XMV"
	BMPTagFLV1                = "FLV1"
	BMPTagFLV1_S263           = "S263"
	BMPTagFLASHSV             = "FSV1"
	BMPTagSVQ1                = "svq1"
	BMPTagTSCC                = "tscc"
	BMPTagULTI                = "ULTI"
	BMPTagVIXL                = "VIXL"
	BMPTagQPEG                = "QPEG"
	BMPTagQPEG_Q1_0           = "Q1.0"
	BMPTagQPEG_Q1_1           = "Q1.1"
	BMPTagWMV3                = "WMV3"
	BMPTagWMV3IMAGE           = "WMVP"
	BMPTagVC1_WVC1            = "WVC1"
	BMPTagVC1_WMVA            = "WMVA"
	BMPTagVC1IMAGE            = "WVP2"
	BMPTagLOCO                = "LOCO"
	BMPTagWNV1_WNV1           = "WNV1"
	BMPTagWNV1_YUV8           = "YUV8"
	BMPTagAASC_AAS4           = "AAS4" // Autodesk 24 bit RLE compressor
	BMPTagAASC                = "AASC"
	BMPTagINDEO2              = "RT21"
	BMPTagFRAPS               = "FPS1"
	BMPTagTHEORA              = "theo"
	BMPTagTRUEMOTION2         = "TM20"
	BMPTagTRUEMOTION2RT       = "TR20"
	BMPTagCSCD                = "CSCD"
	BMPTagZMBV                = "ZMBV"
	BMPTagKMVC                = "KMVC"
	BMPTagCAVS                = "CAVS"
	BMPTagAVS2                = "AVS2"
	BMPTagJPEG2000_mjp2       = "mjp2"
	BMPTagJPEG2000_MJ2C       = "MJ2C"
	BMPTagJPEG2000_LJ2C       = "LJ2C"
	BMPTagJPEG2000_LJ2K       = "LJ2K"
	BMPTagJPEG2000_IPJ2       = "IPJ2"
	BMPTagJPEG2000_AVj2       = "AVj2" // Avid jpeg2000
	BMPTagVMNC                = "VMnc"
	BMPTagTARGA               = "tga "
	BMPTagPNG_MPNG            = "MPNG"
	BMPTagPNG_PNG1            = "PNG1"
	BMPTagPNG                 = "png " // ImageJ
	BMPTagCLJR                = "CLJR"
	BMPTagDIRAC               = "drac"
	BMPTagRPZA_azpr           = "azpr"
	BMPTagRPZA                = "RPZA"
	BMPTagRPZA_rpza           = "rpza"
	BMPTagSP5X                = "SP54"
	BMPTagAURA                = "AURA"
	BMPTagAURA2               = "AUR2"
	BMPTagDPX                 = "dpx "
	BMPTagKGV1                = "KGV1"
	BMPTagLAGARITH            = "LAGS"
	BMPTagAMV                 = "AMVF"
	BMPTagUTVIDEO_ULRA        = "ULRA"
	BMPTagUTVIDEO_ULRG        = "ULRG"
	BMPTagUTVIDEO_ULY0        = "ULY0"
	BMPTagUTVIDEO_ULY2        = "ULY2"
	BMPTagUTVIDEO_ULY4        = "ULY4" // Ut Video version 13.0.1 BT.709 codecs
	BMPTagUTVIDEO_ULH0        = "ULH0"
	BMPTagUTVIDEO_ULH2        = "ULH2"
	BMPTagUTVIDEO_ULH4        = "ULH4"
	BMPTagUTVIDEO_UQY0        = "UQY0"
	BMPTagUTVIDEO_UQY2        = "UQY2"
	BMPTagUTVIDEO_UQRA        = "UQRA"
	BMPTagUTVIDEO_UQRG        = "UQRG"
	BMPTagUTVIDEO_UMY2        = "UMY2"
	BMPTagUTVIDEO_UMH2        = "UMH2"
	BMPTagUTVIDEO_UMY4        = "UMY4"
	BMPTagUTVIDEO_UMH4        = "UMH4"
	BMPTagUTVIDEO_UMRA        = "UMRA"
	BMPTagUTVIDEO_UMRG        = "UMRG"
	BMPTagVBLE                = "VBLE"
	BMPTagESCAPE130           = "E130"
	BMPTagDXTORY              = "xtor"
	BMPTagZEROCODEC           = "ZECO"
	BMPTagY41P                = "Y41P"
	BMPTagFLIC                = "AFLC"
	BMPTagMSS1                = "MSS1"
	BMPTagMSA1                = "MSA1"
	BMPTagTSCC2               = "TSC2"
	BMPTagMTS2                = "MTS2"
	BMPTagCLLC                = "CLLC"
	BMPTagMSS2                = "MSS2"
	BMPTagSVQ3                = "SVQ3"
	BMPTag012V                = "012v"
	BMPTag012V_a12v           = "a12v"
	BMPTagG2M_G2M2            = "G2M2"
	BMPTagG2M_G2M3            = "G2M3"
	BMPTagG2M_G2M4            = "G2M4"
	BMPTagG2M_G2M5            = "G2M5"
	BMPTagFIC                 = "FICV"
	BMPTagHQX                 = "CHQX"
	BMPTagTDSC                = "TDSC"
	BMPTagHQ_HQA              = "CUVC"
	BMPTagRV40                = "RV40"
	BMPTagSCREENPRESSO        = "SPV1"
	BMPTagRSCC                = "RSCC"
	BMPTagRSCC_ISCC           = "ISCC"
	BMPTagCFHD                = "CFHD"
	BMPTagM101                = "M101"
	BMPTagM101_M102           = "M102"
	BMPTagMAGICYUV_MAGY       = "MAGY"
	BMPTagMAGICYUV_M8RG       = "M8RG"
	BMPTagMAGICYUV_M8RA       = "M8RA"
	BMPTagMAGICYUV_M8G0       = "M8G0"
	BMPTagMAGICYUV_M8Y0       = "M8Y0"
	BMPTagMAGICYUV_M8Y2       = "M8Y2"
	BMPTagMAGICYUV_M8Y4       = "M8Y4"
	BMPTagMAGICYUV_M8YA       = "M8YA"
	BMPTagMAGICYUV_M0RA       = "M0RA"
	BMPTagMAGICYUV_M0RG       = "M0RG"
	BMPTagMAGICYUV_M0G0       = "M0G0"
	BMPTagMAGICYUV_M0Y0       = "M0Y0"
	BMPTagMAGICYUV_M0Y2       = "M0Y2"
	BMPTagMAGICYUV_M0Y4       = "M0Y4"
	BMPTagMAGICYUV_M2RA       = "M2RA"
	BMPTagMAGICYUV_M2RG       = "M2RG"
	BMPTagYLC                 = "YLC0"
	BMPTagSPEEDHQ_SHQ0        = "SHQ0"
	BMPTagSPEEDHQ_SHQ1        = "SHQ1"
	BMPTagSPEEDHQ_SHQ2        = "SHQ2"
	BMPTagSPEEDHQ_SHQ3        = "SHQ3"
	BMPTagSPEEDHQ_SHQ4        = "SHQ4"
	BMPTagSPEEDHQ_SHQ5        = "SHQ5"
	BMPTagSPEEDHQ_SHQ7        = "SHQ7"
	BMPTagSPEEDHQ_SHQ9        = "SHQ9"
	BMPTagFMVC                = "FMVC"
	BMPTagSCPR                = "SCPR"
	BMPTagCLEARVIDEO          = "UCOD"
	BMPTagAV1                 = "AV01"
	BMPTagMSCC                = "MSCC"
	BMPTagSRGC                = "SRGC"
	BMPTagIMM4                = "IMM4"
	BMPTagPROSUMER            = "BT20"
	BMPTagMWSC                = "MWSC"
	BMPTagWCMV                = "WCMV"
	BMPTagRASC                = "RASC"
	BMPTagHYMT                = "HYMT"
	BMPTagARBC                = "ARBC"
	BMPTagAGM_AGM0            = "AGM0"
	BMPTagAGM_AGM1            = "AGM1"
	BMPTagAGM_AGM2            = "AGM2"
	BMPTagAGM_AGM3            = "AGM3"
	BMPTagAGM_AGM4            = "AGM4"
	BMPTagAGM_AGM5            = "AGM5"
	BMPTagAGM_AGM6            = "AGM6"
	BMPTagAGM_AGM7            = "AGM7"
	BMPTagLSCR                = "LSCR"
	BMPTagIMM5                = "IMM5"
	BMPTagMVDV                = "MVDV"
	BMPTagMVHA                = "MVHA"
	BMPTagMV30                = "MV30"
	BMPTagNOTCHLC             = "nlc1"
	BMPTagVQC_VQC1            = "VQC1"
	BMPTagVQC_VQC2            = "VQC2"
	// unofficial
	BMPTagHEVC      = "HEVC"
	BMPTagHEVC_H265 = "H265"
)
