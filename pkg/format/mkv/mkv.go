package mkv

// https://tools.ietf.org/html/draft-ietf-cellar-ebml-00
// https://matroska.org/technical/specs/index.html

import (
	"fq/pkg/decode"
	"fq/pkg/format"
	"log"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MKV,
		Description: "Matroska EBML",
		Groups:      []string{format.PROBE},
		DecodeFn:    mkvDecode,
	})
}

/*
proc uint {size {name ""}} {
    set n 0
    for {set i 0} {$i < $size} {incr i} {
        set n [expr $n<<8 | [uint8]]
    }
    if {$name != ""} {
	    entry $name $n $size [expr [pos]-$size]
    }
    return $n
}
*/
/*
func decodeUint(d *decode.D, nBytes int) uint64 {
	var n uint64
	for i:=0; i < nBytes; i++ {
		n |=
	}


}
*/

/*
proc int {size {name ""}} {
    set n 0
    for {set i 0} {$i < $size} {incr i} {
        set n [expr $n<<8 | [uint8]]
    }
    if {$n & (1 << ($size*8-1))} {
        # 2-complement
        set n [expr -((~$n & (1<<($size*8))-1)+1)]
    }

    if {$name != ""} {
	    entry $name $n $size [expr [pos]-$size]
    }
    return $n
}

proc vint {} {
    set n [uint8]

    set width 1
    for {set i 0} {($n & (1<<(7-$i))) == 0} {incr i} {
        incr width
    }
    for {set i 1} {$i < $width} {incr i} {
        set n [expr ($n<<8) | [uint8]]
    }

    # return byte-width raw-n n
    return [list $width $n [expr ((1<<(($width-1)*8+(8-$width)))-1) & $n]]
}
*/

// TODO: smarter?
func decodeRawVintWidth(d *decode.D) (uint64, int) {
	n := d.U8()
	w := 1
	for i := 0; (n & (1 << (7 - i))) == 0; i++ {
		w++
	}
	for i := 1; i < w; i++ {
		n = n<<8 | d.U8()
	}
	return n, w
}

func decodeRawVint(d *decode.D) uint64 {
	n, _ := decodeRawVintWidth(d)
	return n
}

func decodeVint(d *decode.D) uint64 {
	n, w := decodeRawVintWidth(d)
	m := (uint64(1<<((w-1)*8+(8-w))) - 1)
	return n & m
}

func fieldDecodeRawVint(d *decode.D, name string, displayFormat decode.DisplayFormat) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return decodeRawVint(d), displayFormat, ""
	})
}
func fieldDecodeVint(d *decode.D, name string, displayFormat decode.DisplayFormat) uint64 {
	return d.FieldUFn(name, func() (uint64, decode.DisplayFormat, string) {
		return decodeVint(d), displayFormat, ""
	})
}

/*
set ebml_root [dict create \
    1a45dfa3 {EBML master Header} \
    18538067 {Segment master Segment} \
]

set ebml_Header [dict create \
    4286 {EBMLVersion uinteger {}} \
    42f7 {EBMLReadVersion uinteger {}} \
    42f2 {EBMLMaxIDLength uinteger {}} \
    42f3 {EBMLMaxSizeLength uinteger {}} \
    4282 {DocType string {}} \
    4287 {DocTypeVersion uinteger {}} \
    4285 {DocTypeReadVersion uinteger {}} \
]
*/

type ebmlType int

const (
	ebmlInteger ebmlType = iota
	ebmlUinteger
	ebmlFloat
	ebmlString
	ebmlUTF8
	ebmlDate
	ebmlBinary
	ebmlMaster
)

var ebmlTypeNames = map[ebmlType]string{
	ebmlInteger:  "integer",
	ebmlUinteger: "uinteger",
	ebmlFloat:    "float",
	ebmlString:   "string",
	ebmlUTF8:     "UTF8",
	ebmlDate:     "data",
	ebmlBinary:   "binary",
	ebmlMaster:   "master",
}

type ebmlAttribute struct {
	name string
	typ  ebmlType
	tag  ebmlTag
}

type ebmlTag map[uint64]ebmlAttribute

var ebmlGlobal = ebmlTag{
	0xbf: {name: "CRC-32", typ: ebmlBinary},
	0xec: {name: "Void", typ: ebmlBinary},
}

var ebmlHeader = ebmlTag{
	0x4286: {name: "EBMLVersion", typ: ebmlUinteger},
	0x42f7: {name: "EBMLReadVersion", typ: ebmlUinteger},
	0x42f2: {name: "EBMLMaxIDLength", typ: ebmlUinteger},
	0x42f3: {name: "EBMLMaxSizeLength", typ: ebmlUinteger},
	0x4282: {name: "DocType", typ: ebmlString},
	0x4287: {name: "DocTypeVersion", typ: ebmlUinteger},
	0x4285: {name: "DocTypeReadVersion", typ: ebmlUinteger},
}

var ebmlRoot = ebmlTag{
	0x1a45dfa3: {name: "EBML", typ: ebmlMaster, tag: ebmlHeader},
	0x18538067: {name: "Segment", typ: ebmlMaster, tag: mkvSegment},
}

var mkvSegment = ebmlTag{
	0x114d9b74: {name: "SeekHead", typ: ebmlMaster, tag: mkvSeekHead},
	0x1549a966: {name: "Info", typ: ebmlMaster, tag: mkvInfo},
	0x1f43b675: {name: "Cluster", typ: ebmlMaster, tag: mkvCluster},
	0x1654ae6b: {name: "Tracks", typ: ebmlMaster, tag: mkvTracks},
	0x1c53bb6b: {name: "Cues", typ: ebmlMaster, tag: mkvCues},
	0x1941a469: {name: "Attachments", typ: ebmlMaster, tag: mkvAttachments},
	0x1043a770: {name: "Chapters", typ: ebmlMaster, tag: mkvChapters},
	0x1254c367: {name: "Tags", typ: ebmlMaster, tag: mkvTags},
}

var mkvSeekHead = ebmlTag{
	0x4dbb: {name: "Seek", typ: ebmlMaster, tag: mkvSeek},
}

var mkvSeek = ebmlTag{
	0x53ab: {name: "SeekID", typ: ebmlBinary},
	0x53ac: {name: "SeekPosition", typ: ebmlUinteger},
}

var mkvInfo = ebmlTag{
	0x73a4:   {name: "SegmentUID", typ: ebmlBinary},
	0x7384:   {name: "SegmentFilename", typ: ebmlUTF8},
	0x3cb923: {name: "PrevUID", typ: ebmlBinary},
	0x3c83ab: {name: "PrevFilename", typ: ebmlUTF8},
	0x3eb923: {name: "NextUID", typ: ebmlBinary},
	0x3e83bb: {name: "NextFilename", typ: ebmlUTF8},
	0x4444:   {name: "SegmentFamily", typ: ebmlBinary},
	0x6924:   {name: "ChapterTranslate", typ: ebmlMaster, tag: mkvChapterTranslate},
	0x2ad7b1: {name: "TimestampScale", typ: ebmlUinteger},
	0x4489:   {name: "Duration", typ: ebmlFloat},
	0x4461:   {name: "DateUTC", typ: ebmlDate},
	0x7ba9:   {name: "Title", typ: ebmlUTF8},
	0x4d80:   {name: "MuxingApp", typ: ebmlUTF8},
	0x5741:   {name: "WritingApp", typ: ebmlUTF8},
}

var mkvChapterTranslate = ebmlTag{
	0x69fc: {name: "ChapterTranslateEditionUID", typ: ebmlUinteger},
	0x69bf: {name: "ChapterTranslateCodec", typ: ebmlUinteger},
	0x69a5: {name: "ChapterTranslateID", typ: ebmlBinary},
}

var mkvCluster = ebmlTag{
	0xe7:   {name: "Timestamp", typ: ebmlUinteger},
	0x5854: {name: "SilentTracks", typ: ebmlMaster, tag: mkvSilentTracks},
	0xa7:   {name: "Position", typ: ebmlUinteger},
	0xab:   {name: "PrevSize", typ: ebmlUinteger},
	0xa3:   {name: "SimpleBlock", typ: ebmlBinary},
	0xa0:   {name: "BlockGroup", typ: ebmlMaster, tag: mkvBlockGroup},
	0xaf:   {name: "EncryptedBlock", typ: ebmlBinary},
}

var mkvSilentTracks = ebmlTag{
	0x58d7: {name: "SilentTrackNumber", typ: ebmlUinteger},
}

var mkvBlockGroup = ebmlTag{
	0xa1:   {name: "Block", typ: ebmlBinary},
	0xa2:   {name: "BlockVirtual", typ: ebmlBinary},
	0x75a1: {name: "BlockAdditions", typ: ebmlMaster, tag: mkvBlockAdditions},
	0x9b:   {name: "BlockDuration", typ: ebmlUinteger},
	0xfa:   {name: "ReferencePriority", typ: ebmlUinteger},
	0xfb:   {name: "ReferenceBlock", typ: ebmlInteger},
	0xfd:   {name: "ReferenceVirtual", typ: ebmlInteger},
	0xa4:   {name: "CodecState", typ: ebmlBinary},
	0x75a2: {name: "DiscardPadding", typ: ebmlInteger},
	0x8e:   {name: "Slices", typ: ebmlMaster, tag: mkvSlices},
	0xc8:   {name: "ReferenceFrame", typ: ebmlMaster, tag: mkvReferenceFrame},
}

var mkvBlockAdditions = ebmlTag{
	0xa6: {name: "BlockMore", typ: ebmlMaster, tag: mkvBlockMore},
}

var mkvBlockMore = ebmlTag{
	0xee: {name: "BlockAddID", typ: ebmlUinteger},
	0xa5: {name: "BlockAdditional", typ: ebmlBinary},
}

var mkvSlices = ebmlTag{
	0xe8: {name: "TimeSlice", typ: ebmlMaster, tag: mkvTimeSlice},
}

var mkvTimeSlice = ebmlTag{
	0xcc: {name: "LaceNumber", typ: ebmlUinteger},
	0xcd: {name: "FrameNumber", typ: ebmlUinteger},
	0xcb: {name: "BlockAdditionID", typ: ebmlUinteger},
	0xce: {name: "Delay", typ: ebmlUinteger},
	0xcf: {name: "SliceDuration", typ: ebmlUinteger},
}

var mkvReferenceFrame = ebmlTag{
	0xc9: {name: "ReferenceOffset", typ: ebmlUinteger},
	0xca: {name: "ReferenceTimestamp", typ: ebmlUinteger},
}

var mkvTracks = ebmlTag{
	0xae: {name: "TrackEntry", typ: ebmlMaster, tag: mkvTrackEntry},
}

var mkvTrackEntry = ebmlTag{
	0xd7:     {name: "TrackNumber", typ: ebmlUinteger},
	0x73c5:   {name: "TrackUID", typ: ebmlUinteger},
	0x83:     {name: "TrackType", typ: ebmlUinteger},
	0xb9:     {name: "FlagEnabled", typ: ebmlUinteger},
	0x88:     {name: "FlagDefault", typ: ebmlUinteger},
	0x55aa:   {name: "FlagForced", typ: ebmlUinteger},
	0x9c:     {name: "FlagLacing", typ: ebmlUinteger},
	0x6de7:   {name: "MinCache", typ: ebmlUinteger},
	0x6df8:   {name: "MaxCache", typ: ebmlUinteger},
	0x23e383: {name: "DefaultDuration", typ: ebmlUinteger},
	0x234e7a: {name: "DefaultDecodedFieldDuration", typ: ebmlUinteger},
	0x23314f: {name: "TrackTimestampScale", typ: ebmlFloat},
	0x537f:   {name: "TrackOffset", typ: ebmlInteger},
	0x55ee:   {name: "MaxBlockAdditionID", typ: ebmlUinteger},
	0x41e4:   {name: "BlockAdditionMapping", typ: ebmlMaster, tag: mkvBlockAdditionMapping},
	0x536e:   {name: "Name", typ: ebmlUTF8},
	0x22b59c: {name: "Language", typ: ebmlString},
	0x22b59d: {name: "LanguageIETF", typ: ebmlString},
	0x86:     {name: "CodecID", typ: ebmlString},
	0x63a2:   {name: "CodecPrivate", typ: ebmlBinary},
	0x258688: {name: "CodecName", typ: ebmlUTF8},
	0x7446:   {name: "AttachmentLink", typ: ebmlUinteger},
	0x3a9697: {name: "CodecSettings", typ: ebmlUTF8},
	0x3b4040: {name: "CodecInfoURL", typ: ebmlString},
	0x26b240: {name: "CodecDownloadURL", typ: ebmlString},
	0xaa:     {name: "CodecDecodeAll", typ: ebmlUinteger},
	0x6fab:   {name: "TrackOverlay", typ: ebmlUinteger},
	0x56aa:   {name: "CodecDelay", typ: ebmlUinteger},
	0x56bb:   {name: "SeekPreRoll", typ: ebmlUinteger},
	0x6624:   {name: "TrackTranslate", typ: ebmlMaster, tag: mkvTrackTranslate},
	0xe0:     {name: "Video", typ: ebmlMaster, tag: mkvVideo},
	0xe1:     {name: "Audio", typ: ebmlMaster, tag: mkvAudio},
	0xe2:     {name: "TrackOperation", typ: ebmlMaster, tag: mkvTrackOperation},
	0xc0:     {name: "TrickTrackUID", typ: ebmlUinteger},
	0xc1:     {name: "TrickTrackSegmentUID", typ: ebmlBinary},
	0xc6:     {name: "TrickTrackFlag", typ: ebmlUinteger},
	0xc7:     {name: "TrickMasterTrackUID", typ: ebmlUinteger},
	0xc4:     {name: "TrickMasterTrackSegmentUID", typ: ebmlBinary},
	0x6d80:   {name: "ContentEncodings", typ: ebmlMaster, tag: mkvContentEncodings},
}

var mkvBlockAdditionMapping = ebmlTag{
	0x41f0: {name: "BlockAddIDValue", typ: ebmlUinteger},
	0x41a4: {name: "BlockAddIDName", typ: ebmlString},
	0x41e7: {name: "BlockAddIDType", typ: ebmlUinteger},
	0x41ed: {name: "BlockAddIDExtraData", typ: ebmlBinary},
}

var mkvTrackTranslate = ebmlTag{
	0x66fc: {name: "TrackTranslateEditionUID", typ: ebmlUinteger},
	0x66bf: {name: "TrackTranslateCodec", typ: ebmlUinteger},
	0x66a5: {name: "TrackTranslateTrackID", typ: ebmlBinary},
}

var mkvVideo = ebmlTag{
	0x9a:     {name: "FlagInterlaced", typ: ebmlUinteger},
	0x9d:     {name: "FieldOrder", typ: ebmlUinteger},
	0x53b8:   {name: "StereoMode", typ: ebmlUinteger},
	0x53c0:   {name: "AlphaMode", typ: ebmlUinteger},
	0x53b9:   {name: "OldStereoMode", typ: ebmlUinteger},
	0xb0:     {name: "PixelWidth", typ: ebmlUinteger},
	0xba:     {name: "PixelHeight", typ: ebmlUinteger},
	0x54aa:   {name: "PixelCropBottom", typ: ebmlUinteger},
	0x54bb:   {name: "PixelCropTop", typ: ebmlUinteger},
	0x54cc:   {name: "PixelCropLeft", typ: ebmlUinteger},
	0x54dd:   {name: "PixelCropRight", typ: ebmlUinteger},
	0x54b0:   {name: "DisplayWidth", typ: ebmlUinteger},
	0x54ba:   {name: "DisplayHeight", typ: ebmlUinteger},
	0x54b2:   {name: "DisplayUnit", typ: ebmlUinteger},
	0x54b3:   {name: "AspectRatioType", typ: ebmlUinteger},
	0x2eb524: {name: "ColourSpace", typ: ebmlBinary},
	0x2fb523: {name: "GammaValue", typ: ebmlFloat},
	0x2383e3: {name: "FrameRate", typ: ebmlFloat},
	0x55b0:   {name: "Colour", typ: ebmlMaster, tag: mkvColour},
	0x7670:   {name: "Projection", typ: ebmlMaster, tag: mkvProjection},
}

var mkvColour = ebmlTag{
	0x55b1: {name: "MatrixCoefficients", typ: ebmlUinteger},
	0x55b2: {name: "BitsPerChannel", typ: ebmlUinteger},
	0x55b3: {name: "ChromaSubsamplingHorz", typ: ebmlUinteger},
	0x55b4: {name: "ChromaSubsamplingVert", typ: ebmlUinteger},
	0x55b5: {name: "CbSubsamplingHorz", typ: ebmlUinteger},
	0x55b6: {name: "CbSubsamplingVert", typ: ebmlUinteger},
	0x55b7: {name: "ChromaSitingHorz", typ: ebmlUinteger},
	0x55b8: {name: "ChromaSitingVert", typ: ebmlUinteger},
	0x55b9: {name: "Range", typ: ebmlUinteger},
	0x55ba: {name: "TransferCharacteristics", typ: ebmlUinteger},
	0x55bb: {name: "Primaries", typ: ebmlUinteger},
	0x55bc: {name: "MaxCLL", typ: ebmlUinteger},
	0x55bd: {name: "MaxFALL", typ: ebmlUinteger},
	0x55d0: {name: "MasteringMetadata", typ: ebmlMaster, tag: mkvMasteringMetadata},
}

var mkvMasteringMetadata = ebmlTag{
	0x55d1: {name: "PrimaryRChromaticityX", typ: ebmlFloat},
	0x55d2: {name: "PrimaryRChromaticityY", typ: ebmlFloat},
	0x55d3: {name: "PrimaryGChromaticityX", typ: ebmlFloat},
	0x55d4: {name: "PrimaryGChromaticityY", typ: ebmlFloat},
	0x55d5: {name: "PrimaryBChromaticityX", typ: ebmlFloat},
	0x55d6: {name: "PrimaryBChromaticityY", typ: ebmlFloat},
	0x55d7: {name: "WhitePointChromaticityX", typ: ebmlFloat},
	0x55d8: {name: "WhitePointChromaticityY", typ: ebmlFloat},
	0x55d9: {name: "LuminanceMax", typ: ebmlFloat},
	0x55da: {name: "LuminanceMin", typ: ebmlFloat},
}

var mkvProjection = ebmlTag{
	0x7671: {name: "ProjectionType", typ: ebmlUinteger},
	0x7672: {name: "ProjectionPrivate", typ: ebmlBinary},
	0x7673: {name: "ProjectionPoseYaw", typ: ebmlFloat},
	0x7674: {name: "ProjectionPosePitch", typ: ebmlFloat},
	0x7675: {name: "ProjectionPoseRoll", typ: ebmlFloat},
}

var mkvAudio = ebmlTag{
	0xb5:   {name: "SamplingFrequency", typ: ebmlFloat},
	0x78b5: {name: "OutputSamplingFrequency", typ: ebmlFloat},
	0x9f:   {name: "Channels", typ: ebmlUinteger},
	0x7d7b: {name: "ChannelPositions", typ: ebmlBinary},
	0x6264: {name: "BitDepth", typ: ebmlUinteger},
}

var mkvTrackOperation = ebmlTag{
	0xe3: {name: "TrackCombinePlanes", typ: ebmlMaster, tag: mkvTrackCombinePlanes},
	0xe9: {name: "TrackJoinBlocks", typ: ebmlMaster, tag: mkvTrackJoinBlocks},
}

var mkvTrackCombinePlanes = ebmlTag{
	0xe4: {name: "TrackPlane", typ: ebmlMaster, tag: mkvTrackPlane},
}

var mkvTrackPlane = ebmlTag{
	0xe5: {name: "TrackPlaneUID", typ: ebmlUinteger},
	0xe6: {name: "TrackPlaneType", typ: ebmlUinteger},
}

var mkvTrackJoinBlocks = ebmlTag{
	0xed: {name: "TrackJoinUID", typ: ebmlUinteger},
}

var mkvContentEncodings = ebmlTag{
	0x6240: {name: "ContentEncoding", typ: ebmlMaster, tag: mkvContentEncoding},
}

var mkvContentEncoding = ebmlTag{
	0x5031: {name: "ContentEncodingOrder", typ: ebmlUinteger},
	0x5032: {name: "ContentEncodingScope", typ: ebmlUinteger},
	0x5033: {name: "ContentEncodingType", typ: ebmlUinteger},
	0x5034: {name: "ContentCompression", typ: ebmlMaster, tag: mkvContentCompression},
	0x5035: {name: "ContentEncryption", typ: ebmlMaster, tag: mkvContentEncryption},
}

var mkvContentCompression = ebmlTag{
	0x4254: {name: "ContentCompAlgo", typ: ebmlUinteger},
	0x4255: {name: "ContentCompSettings", typ: ebmlBinary},
}

var mkvContentEncryption = ebmlTag{
	0x47e1: {name: "ContentEncAlgo", typ: ebmlUinteger},
	0x47e2: {name: "ContentEncKeyID", typ: ebmlBinary},
	0x47e7: {name: "ContentEncAESSettings", typ: ebmlMaster, tag: mkvContentEncAESSettings},
	0x47e3: {name: "ContentSignature", typ: ebmlBinary},
	0x47e4: {name: "ContentSigKeyID", typ: ebmlBinary},
	0x47e5: {name: "ContentSigAlgo", typ: ebmlUinteger},
	0x47e6: {name: "ContentSigHashAlgo", typ: ebmlUinteger},
}

var mkvContentEncAESSettings = ebmlTag{
	0x47e8: {name: "AESSettingsCipherMode", typ: ebmlUinteger},
}

var mkvCues = ebmlTag{
	0xbb: {name: "CuePoint", typ: ebmlMaster, tag: mkvCuePoint},
}

var mkvCuePoint = ebmlTag{
	0xb3: {name: "CueTime", typ: ebmlUinteger},
	0xb7: {name: "CueTrackPositions", typ: ebmlMaster, tag: mkvCueTrackPositions},
}

var mkvCueTrackPositions = ebmlTag{
	0xf7:   {name: "CueTrack", typ: ebmlUinteger},
	0xf1:   {name: "CueClusterPosition", typ: ebmlUinteger},
	0xf0:   {name: "CueRelativePosition", typ: ebmlUinteger},
	0xb2:   {name: "CueDuration", typ: ebmlUinteger},
	0x5378: {name: "CueBlockNumber", typ: ebmlUinteger},
	0xea:   {name: "CueCodecState", typ: ebmlUinteger},
	0xdb:   {name: "CueReference", typ: ebmlMaster, tag: mkvCueReference},
}

var mkvCueReference = ebmlTag{
	0x96:   {name: "CueRefTime", typ: ebmlUinteger},
	0x97:   {name: "CueRefCluster", typ: ebmlUinteger},
	0x535f: {name: "CueRefNumber", typ: ebmlUinteger},
	0xeb:   {name: "CueRefCodecState", typ: ebmlUinteger},
}

var mkvAttachments = ebmlTag{
	0x61a7: {name: "AttachedFile", typ: ebmlMaster, tag: mkvAttachedFile},
}

var mkvAttachedFile = ebmlTag{
	0x467e: {name: "FileDescription", typ: ebmlUTF8},
	0x466e: {name: "FileName", typ: ebmlUTF8},
	0x4660: {name: "FileMimeType", typ: ebmlString},
	0x465c: {name: "FileData", typ: ebmlBinary},
	0x46ae: {name: "FileUID", typ: ebmlUinteger},
	0x4675: {name: "FileReferral", typ: ebmlBinary},
	0x4661: {name: "FileUsedStartTime", typ: ebmlUinteger},
	0x4662: {name: "FileUsedEndTime", typ: ebmlUinteger},
}

var mkvChapters = ebmlTag{
	0x45b9: {name: "EditionEntry", typ: ebmlMaster, tag: mkvEditionEntry},
}

var mkvEditionEntry = ebmlTag{
	0x45bc: {name: "EditionUID", typ: ebmlUinteger},
	0x45bd: {name: "EditionFlagHidden", typ: ebmlUinteger},
	0x45db: {name: "EditionFlagDefault", typ: ebmlUinteger},
	0x45dd: {name: "EditionFlagOrdered", typ: ebmlUinteger},
	0xb6:   {name: "ChapterAtom", typ: ebmlMaster, tag: mkvChapterAtom},
}

var mkvChapterAtom = ebmlTag{
	0x73c4: {name: "ChapterUID", typ: ebmlUinteger},
	0x5654: {name: "ChapterStringUID", typ: ebmlUTF8},
	0x91:   {name: "ChapterTimeStart", typ: ebmlUinteger},
	0x92:   {name: "ChapterTimeEnd", typ: ebmlUinteger},
	0x98:   {name: "ChapterFlagHidden", typ: ebmlUinteger},
	0x4598: {name: "ChapterFlagEnabled", typ: ebmlUinteger},
	0x6e67: {name: "ChapterSegmentUID", typ: ebmlBinary},
	0x6ebc: {name: "ChapterSegmentEditionUID", typ: ebmlUinteger},
	0x63c3: {name: "ChapterPhysicalEquiv", typ: ebmlUinteger},
	0x8f:   {name: "ChapterTrack", typ: ebmlMaster, tag: mkvChapterTrack},
	0x80:   {name: "ChapterDisplay", typ: ebmlMaster, tag: mkvChapterDisplay},
	0x6944: {name: "ChapProcess", typ: ebmlMaster, tag: mkvChapProcess},
}

var mkvChapterTrack = ebmlTag{
	0x89: {name: "ChapterTrackUID", typ: ebmlUinteger},
}

var mkvChapterDisplay = ebmlTag{
	0x85:   {name: "ChapString", typ: ebmlUTF8},
	0x437c: {name: "ChapLanguage", typ: ebmlString},
	0x437d: {name: "ChapLanguageIETF", typ: ebmlString},
	0x437e: {name: "ChapCountry", typ: ebmlString},
}

var mkvChapProcess = ebmlTag{
	0x6955: {name: "ChapProcessCodecID", typ: ebmlUinteger},
	0x450d: {name: "ChapProcessPrivate", typ: ebmlBinary},
	0x6911: {name: "ChapProcessCommand", typ: ebmlMaster, tag: mkvChapProcessCommand},
}

var mkvChapProcessCommand = ebmlTag{
	0x6922: {name: "ChapProcessTime", typ: ebmlUinteger},
	0x6933: {name: "ChapProcessData", typ: ebmlBinary},
}

var mkvTags = ebmlTag{
	0x7373: {name: "Tag", typ: ebmlMaster, tag: mkvTag},
}

var mkvTag = ebmlTag{
	0x63c0: {name: "Targets", typ: ebmlMaster, tag: mkvTargets},
	0x67c8: {name: "SimpleTag", typ: ebmlMaster, tag: mkvSimpleTag},
}

var mkvTargets = ebmlTag{
	0x68ca: {name: "TargetTypeValue", typ: ebmlUinteger},
	0x63ca: {name: "TargetType", typ: ebmlString},
	0x63c5: {name: "TagTrackUID", typ: ebmlUinteger},
	0x63c9: {name: "TagEditionUID", typ: ebmlUinteger},
	0x63c4: {name: "TagChapterUID", typ: ebmlUinteger},
	0x63c6: {name: "TagAttachmentUID", typ: ebmlUinteger},
}

var mkvSimpleTag = ebmlTag{
	0x45a3: {name: "TagName", typ: ebmlUTF8},
	0x447a: {name: "TagLanguage", typ: ebmlString},
	0x447b: {name: "TagLanguageIETF", typ: ebmlString},
	0x4484: {name: "TagDefault", typ: ebmlUinteger},
	0x4487: {name: "TagString", typ: ebmlUTF8},
	0x4485: {name: "TagBinary", typ: ebmlBinary},
}

/*
proc type_master {size _label extra} {
    upvar #0 "ebml_$extra" tags
    global ebml_Global
    set garbage_size 0

    # TODO: unknown-size might not be correct handled
    while {![end] && ($size > 0 || $size == -1)} {
        lassign [vint] tag_id_width tag_idnr
        set tag_id [format "%x" $tag_idnr]
        lassign [vint] tag_size_width tag_size_raw tag_size

        set tag_name "Unknown"
        set tag_type "binary"
        set tag_extra {}
        set tag_desc ""
        if {[dict exists $tags $tag_id]} {
            lassign [dict get $tags $tag_id] tag_name tag_type tag_extra tag_desc
        } elseif {[dict exists $ebml_Global $tag_id]} {
            lassign [dict get $ebml_Global $tag_id] tag_name tag_type tag_extra tag_desc
        } elseif {$size == -1} {
            incr garbage_size
            move [expr -($tag_id_width+$tag_size_width-1)]
            continue
        }

        if {$garbage_size != 0} {
            entry "Garbage" {} $garbage_size [expr [pos]-$garbage_size-$tag_id_width-$tag_size_width]
            set garbage_size 0
        }

        set type_fn "type_$tag_type"

        section "$tag_name ($tag_type)" {
            entry "ID" $tag_id $tag_id_width [expr [pos]-$tag_id_width-$tag_size_width]
            set tag_size_str $tag_size
            if {$tag_size_raw == 0xff} {
                append tag_size_str " (unknown)"
                set tag_size -1
            }
            entry "Size" "$tag_size_str" $tag_size_width [expr [pos]-$tag_size_width]
            $type_fn $tag_size $tag_name $tag_extra
        }

        if {$size == -1} {
            continue
        }
        incr size [expr -($tag_id_width+$tag_size_width+$tag_size)]
    }
}
*/

func decodeMaster(d *decode.D, nBytes uint64, tag ebmlTag) {

	d.FieldArrayFn("bla", func(d *decode.D) {

		for d.NotEnd() {

			startPos := d.Pos()

			tagID := decodeRawVint(d)
			d.SeekAbs(startPos)

			a, ok := tag[tagID]
			if !ok {
				a, ok = ebmlGlobal[tagID]
				if !ok {
					//return
					panic("asdsad")
				}
			}

			d.FieldStructFn(a.name, func(d *decode.D) {
				fieldDecodeRawVint(d, "id", decode.NumberHex)
				tagSize := fieldDecodeVint(d, "size", decode.NumberDecimal)

				switch a.typ {
				case ebmlInteger:
					d.FieldS("value", int(tagSize)*8)
				case ebmlUinteger:
					d.FieldU("value", int(tagSize)*8)
				case ebmlFloat:
					d.FieldF("value", int(tagSize)*8)
				case ebmlString:
					d.FieldUTF8("value", int(tagSize))
				case ebmlUTF8:
					d.FieldUTF8("value", int(tagSize))
				case ebmlDate:
					// TODO:
					d.FieldBitBufLen("value", int64(tagSize)*8)
				case ebmlBinary:
					d.FieldBitBufLen("value", int64(tagSize)*8)
				case ebmlMaster:
					d.DecodeLenFn(int64(tagSize)*8, func(d *decode.D) {
						decodeMaster(d, 0, a.tag)
					})
				}

				log.Println("bla")
			})
		}
	})

}

func mkvDecode(d *decode.D) interface{} {
	decodeMaster(d, 0, ebmlRoot)
	return nil
}

/*
# EBML, Matroska and webm binary template
#
# Specification can be found at:

# Copyright (c) 2020 Mattias Wadman
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
# of the Software, and to permit persons to whom the Software is furnished to do
# so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

# TODO:
# enums
# default value when zero length

requires 0 "1a 45 df a3"
big_endian

proc ascii_maybe_empty {size {name ""}} {
    if {$size > 0} {
        if {$name != ""} {
            return [ascii $size $name]
        } else {
            return [ascii $size]
        }
    } else {
        if {$name != ""} {
            entry $name ""
        }
        return ""
    }
}

proc bytes_maybe_empty {size {name ""}} {
    if {$size > 0} {
        if {$name != ""} {
            return [bytes $size $name]
        } else {
            return [bytes $size]
        }
    } else {
        if {$name != ""} {
            entry $name
        }
        return ""
    }
}

proc utf8 {size label} {
    set s ""
    if {$size > 0} {
        set bytes [bytes $size]
        set s [encoding convertfrom utf-8 [string trimright $bytes "\x00"]]
    }
    entry $label $s $size [expr [pos]-$size]
}

proc uint {size {name ""}} {
    set n 0
    for {set i 0} {$i < $size} {incr i} {
        set n [expr $n<<8 | [uint8]]
    }
    if {$name != ""} {
	    entry $name $n $size [expr [pos]-$size]
    }
    return $n
}

proc int {size {name ""}} {
    set n 0
    for {set i 0} {$i < $size} {incr i} {
        set n [expr $n<<8 | [uint8]]
    }
    if {$n & (1 << ($size*8-1))} {
        # 2-complement
        set n [expr -((~$n & (1<<($size*8))-1)+1)]
    }

    if {$name != ""} {
	    entry $name $n $size [expr [pos]-$size]
    }
    return $n
}

proc vint {} {
    set n [uint8]

    set width 1
    for {set i 0} {($n & (1<<(7-$i))) == 0} {incr i} {
        incr width
    }
    for {set i 1} {$i < $width} {incr i} {
        set n [expr ($n<<8) | [uint8]]
    }

    # return byte-width raw-n n
    return [list $width $n [expr ((1<<(($width-1)*8+(8-$width)))-1) & $n]]
}

proc type_string {size label _extra} {
    ascii_maybe_empty $size $label
}

proc type_binary {size label _extra} {
    bytes_maybe_empty $size $label
}

proc type_utf-8 {size label _extra} {
    utf8 $size $label
}

proc type_uinteger {size label _extra} {
    switch $size {
        0 {entry $label 0}
        1 {uint8 $label}
        2 {uint16 $label}
        4 {uint32 $label}
        8 {uint64 $label}
        3 -
        5 -
        6 -
        7 {uint $size $label}
        default {bytes $size $label}
    }
}

proc type_integer {size label _extra} {
    switch $size {
        0 {entry $label 0}
        1 {int8 $label}
        2 {int16 $label}
        4 {int32 $label}
        8 {int64 $label}
        3 -
        5 -
        6 -
        7 {int $size $label}
        default {bytes $size $label}
    }
}

proc type_float {size label _extr} {
    switch $size {
        0 {entry $label 0}
        4 {float $label}
        8 {double $label}
        default {bytes $size $label}
    }
}

proc type_date {size label _extra} {
    set s [clock scan {2001-01-01 00:00:00}]
    set frac 0
    switch $size {
        0 {}
        8 {
            set nano [int64]
            set s [clock add $s [expr $nano/1000000000] seconds]
            set frac [expr ($nano%1000000000)/1000000000.0]
        }
        default {
            bytes $size $label
            return
        }
    }

    entry $label "[clock format $s] ${frac}s" $size [expr [pos]-$size]
}

proc type_master {size _label extra} {
    upvar #0 "ebml_$extra" tags
    global ebml_Global
    set garbage_size 0

    # TODO: unknown-size might not be correct handled
    while {![end] && ($size > 0 || $size == -1)} {
        lassign [vint] tag_id_width tag_idnr
        set tag_id [format "%x" $tag_idnr]
        lassign [vint] tag_size_width tag_size_raw tag_size

        set tag_name "Unknown"
        set tag_type "binary"
        set tag_extra {}
        set tag_desc ""
        if {[dict exists $tags $tag_id]} {
            lassign [dict get $tags $tag_id] tag_name tag_type tag_extra tag_desc
        } elseif {[dict exists $ebml_Global $tag_id]} {
            lassign [dict get $ebml_Global $tag_id] tag_name tag_type tag_extra tag_desc
        } elseif {$size == -1} {
            incr garbage_size
            move [expr -($tag_id_width+$tag_size_width-1)]
            continue
        }

        if {$garbage_size != 0} {
            entry "Garbage" {} $garbage_size [expr [pos]-$garbage_size-$tag_id_width-$tag_size_width]
            set garbage_size 0
        }

        set type_fn "type_$tag_type"

        section "$tag_name ($tag_type)" {
            entry "ID" $tag_id $tag_id_width [expr [pos]-$tag_id_width-$tag_size_width]
            set tag_size_str $tag_size
            if {$tag_size_raw == 0xff} {
                append tag_size_str " (unknown)"
                set tag_size -1
            }
            entry "Size" "$tag_size_str" $tag_size_width [expr [pos]-$tag_size_width]
            $type_fn $tag_size $tag_name $tag_extra
        }

        if {$size == -1} {
            continue
        }
        incr size [expr -($tag_id_width+$tag_size_width+$tag_size)]
    }
}

# generated from https://raw.githubusercontent.com/cellar-wg/matroska-specification/aa2144a58b661baf54b99bab41113d66b0f5ff62/ebml_matroska.xml
# using https://gist.github.com/wader/e15b0966dc464db5d70c2a155537ba1f
set ebml_Global [dict create \
    bf {CRC-32 binary {}} \
    ec {Void binary {}} \
]

set ebml_root [dict create \
    1a45dfa3 {EBML master Header} \
    18538067 {Segment master Segment} \
]

set ebml_Header [dict create \
    4286 {EBMLVersion uinteger {}} \
    42f7 {EBMLReadVersion uinteger {}} \
    42f2 {EBMLMaxIDLength uinteger {}} \
    42f3 {EBMLMaxSizeLength uinteger {}} \
    4282 {DocType string {}} \
    4287 {DocTypeVersion uinteger {}} \
    4285 {DocTypeReadVersion uinteger {}} \
]

set ebml_Segment [dict create \
    114d9b74 {SeekHead master SeekHead} \
    1549a966 {Info master Info} \
    1f43b675 {Cluster master Cluster} \
    1654ae6b {Tracks master Tracks} \
    1c53bb6b {Cues master Cues} \
    1941a469 {Attachments master Attachments} \
    1043a770 {Chapters master Chapters} \
    1254c367 {Tags master Tags} \
]

set ebml_SeekHead [dict create \
    4dbb {Seek master Seek} \
]

set ebml_Seek [dict create \
    53ab {SeekID binary {}} \
    53ac {SeekPosition uinteger {}} \
]

set ebml_Info [dict create \
    73a4 {SegmentUID binary {}} \
    7384 {SegmentFilename utf-8 {}} \
    3cb923 {PrevUID binary {}} \
    3c83ab {PrevFilename utf-8 {}} \
    3eb923 {NextUID binary {}} \
    3e83bb {NextFilename utf-8 {}} \
    4444 {SegmentFamily binary {}} \
    6924 {ChapterTranslate master ChapterTranslate} \
    2ad7b1 {TimestampScale uinteger {}} \
    4489 {Duration float {}} \
    4461 {DateUTC date {}} \
    7ba9 {Title utf-8 {}} \
    4d80 {MuxingApp utf-8 {}} \
    5741 {WritingApp utf-8 {}} \
]

set ebml_ChapterTranslate [dict create \
    69fc {ChapterTranslateEditionUID uinteger {}} \
    69bf {ChapterTranslateCodec uinteger {}} \
    69a5 {ChapterTranslateID binary {}} \
]

set ebml_Cluster [dict create \
    e7 {Timestamp uinteger {}} \
    5854 {SilentTracks master SilentTracks} \
    a7 {Position uinteger {}} \
    ab {PrevSize uinteger {}} \
    a3 {SimpleBlock binary {}} \
    a0 {BlockGroup master BlockGroup} \
    af {EncryptedBlock binary {}} \
]

set ebml_SilentTracks [dict create \
    58d7 {SilentTrackNumber uinteger {}} \
]

set ebml_BlockGroup [dict create \
    a1 {Block binary {}} \
    a2 {BlockVirtual binary {}} \
    75a1 {BlockAdditions master BlockAdditions} \
    9b {BlockDuration uinteger {}} \
    fa {ReferencePriority uinteger {}} \
    fb {ReferenceBlock integer {}} \
    fd {ReferenceVirtual integer {}} \
    a4 {CodecState binary {}} \
    75a2 {DiscardPadding integer {}} \
    8e {Slices master Slices} \
    c8 {ReferenceFrame master ReferenceFrame} \
]

set ebml_BlockAdditions [dict create \
    a6 {BlockMore master BlockMore} \
]

set ebml_BlockMore [dict create \
    ee {BlockAddID uinteger {}} \
    a5 {BlockAdditional binary {}} \
]

set ebml_Slices [dict create \
    e8 {TimeSlice master TimeSlice} \
]

set ebml_TimeSlice [dict create \
    cc {LaceNumber uinteger {}} \
    cd {FrameNumber uinteger {}} \
    cb {BlockAdditionID uinteger {}} \
    ce {Delay uinteger {}} \
    cf {SliceDuration uinteger {}} \
]

set ebml_ReferenceFrame [dict create \
    c9 {ReferenceOffset uinteger {}} \
    ca {ReferenceTimestamp uinteger {}} \
]

set ebml_Tracks [dict create \
    ae {TrackEntry master TrackEntry} \
]

set ebml_TrackEntry [dict create \
    d7 {TrackNumber uinteger {}} \
    73c5 {TrackUID uinteger {}} \
    83 {TrackType uinteger {}} \
    b9 {FlagEnabled uinteger {}} \
    88 {FlagDefault uinteger {}} \
    55aa {FlagForced uinteger {}} \
    9c {FlagLacing uinteger {}} \
    6de7 {MinCache uinteger {}} \
    6df8 {MaxCache uinteger {}} \
    23e383 {DefaultDuration uinteger {}} \
    234e7a {DefaultDecodedFieldDuration uinteger {}} \
    23314f {TrackTimestampScale float {}} \
    537f {TrackOffset integer {}} \
    55ee {MaxBlockAdditionID uinteger {}} \
    41e4 {BlockAdditionMapping master BlockAdditionMapping} \
    536e {Name utf-8 {}} \
    22b59c {Language string {}} \
    22b59d {LanguageIETF string {}} \
    86 {CodecID string {}} \
    63a2 {CodecPrivate binary {}} \
    258688 {CodecName utf-8 {}} \
    7446 {AttachmentLink uinteger {}} \
    3a9697 {CodecSettings utf-8 {}} \
    3b4040 {CodecInfoURL string {}} \
    26b240 {CodecDownloadURL string {}} \
    aa {CodecDecodeAll uinteger {}} \
    6fab {TrackOverlay uinteger {}} \
    56aa {CodecDelay uinteger {}} \
    56bb {SeekPreRoll uinteger {}} \
    6624 {TrackTranslate master TrackTranslate} \
    e0 {Video master Video} \
    e1 {Audio master Audio} \
    e2 {TrackOperation master TrackOperation} \
    c0 {TrickTrackUID uinteger {}} \
    c1 {TrickTrackSegmentUID binary {}} \
    c6 {TrickTrackFlag uinteger {}} \
    c7 {TrickMasterTrackUID uinteger {}} \
    c4 {TrickMasterTrackSegmentUID binary {}} \
    6d80 {ContentEncodings master ContentEncodings} \
]

set ebml_BlockAdditionMapping [dict create \
    41f0 {BlockAddIDValue uinteger {}} \
    41a4 {BlockAddIDName string {}} \
    41e7 {BlockAddIDType uinteger {}} \
    41ed {BlockAddIDExtraData binary {}} \
]

set ebml_TrackTranslate [dict create \
    66fc {TrackTranslateEditionUID uinteger {}} \
    66bf {TrackTranslateCodec uinteger {}} \
    66a5 {TrackTranslateTrackID binary {}} \
]

set ebml_Video [dict create \
    9a {FlagInterlaced uinteger {}} \
    9d {FieldOrder uinteger {}} \
    53b8 {StereoMode uinteger {}} \
    53c0 {AlphaMode uinteger {}} \
    53b9 {OldStereoMode uinteger {}} \
    b0 {PixelWidth uinteger {}} \
    ba {PixelHeight uinteger {}} \
    54aa {PixelCropBottom uinteger {}} \
    54bb {PixelCropTop uinteger {}} \
    54cc {PixelCropLeft uinteger {}} \
    54dd {PixelCropRight uinteger {}} \
    54b0 {DisplayWidth uinteger {}} \
    54ba {DisplayHeight uinteger {}} \
    54b2 {DisplayUnit uinteger {}} \
    54b3 {AspectRatioType uinteger {}} \
    2eb524 {ColourSpace binary {}} \
    2fb523 {GammaValue float {}} \
    2383e3 {FrameRate float {}} \
    55b0 {Colour master Colour} \
    7670 {Projection master Projection} \
]

set ebml_Colour [dict create \
    55b1 {MatrixCoefficients uinteger {}} \
    55b2 {BitsPerChannel uinteger {}} \
    55b3 {ChromaSubsamplingHorz uinteger {}} \
    55b4 {ChromaSubsamplingVert uinteger {}} \
    55b5 {CbSubsamplingHorz uinteger {}} \
    55b6 {CbSubsamplingVert uinteger {}} \
    55b7 {ChromaSitingHorz uinteger {}} \
    55b8 {ChromaSitingVert uinteger {}} \
    55b9 {Range uinteger {}} \
    55ba {TransferCharacteristics uinteger {}} \
    55bb {Primaries uinteger {}} \
    55bc {MaxCLL uinteger {}} \
    55bd {MaxFALL uinteger {}} \
    55d0 {MasteringMetadata master MasteringMetadata} \
]

set ebml_MasteringMetadata [dict create \
    55d1 {PrimaryRChromaticityX float {}} \
    55d2 {PrimaryRChromaticityY float {}} \
    55d3 {PrimaryGChromaticityX float {}} \
    55d4 {PrimaryGChromaticityY float {}} \
    55d5 {PrimaryBChromaticityX float {}} \
    55d6 {PrimaryBChromaticityY float {}} \
    55d7 {WhitePointChromaticityX float {}} \
    55d8 {WhitePointChromaticityY float {}} \
    55d9 {LuminanceMax float {}} \
    55da {LuminanceMin float {}} \
]

set ebml_Projection [dict create \
    7671 {ProjectionType uinteger {}} \
    7672 {ProjectionPrivate binary {}} \
    7673 {ProjectionPoseYaw float {}} \
    7674 {ProjectionPosePitch float {}} \
    7675 {ProjectionPoseRoll float {}} \
]

set ebml_Audio [dict create \
    b5 {SamplingFrequency float {}} \
    78b5 {OutputSamplingFrequency float {}} \
    9f {Channels uinteger {}} \
    7d7b {ChannelPositions binary {}} \
    6264 {BitDepth uinteger {}} \
]

set ebml_TrackOperation [dict create \
    e3 {TrackCombinePlanes master TrackCombinePlanes} \
    e9 {TrackJoinBlocks master TrackJoinBlocks} \
]

set ebml_TrackCombinePlanes [dict create \
    e4 {TrackPlane master TrackPlane} \
]

set ebml_TrackPlane [dict create \
    e5 {TrackPlaneUID uinteger {}} \
    e6 {TrackPlaneType uinteger {}} \
]

set ebml_TrackJoinBlocks [dict create \
    ed {TrackJoinUID uinteger {}} \
]

set ebml_ContentEncodings [dict create \
    6240 {ContentEncoding master ContentEncoding} \
]

set ebml_ContentEncoding [dict create \
    5031 {ContentEncodingOrder uinteger {}} \
    5032 {ContentEncodingScope uinteger {}} \
    5033 {ContentEncodingType uinteger {}} \
    5034 {ContentCompression master ContentCompression} \
    5035 {ContentEncryption master ContentEncryption} \
]

set ebml_ContentCompression [dict create \
    4254 {ContentCompAlgo uinteger {}} \
    4255 {ContentCompSettings binary {}} \
]

set ebml_ContentEncryption [dict create \
    47e1 {ContentEncAlgo uinteger {}} \
    47e2 {ContentEncKeyID binary {}} \
    47e7 {ContentEncAESSettings master ContentEncAESSettings} \
    47e3 {ContentSignature binary {}} \
    47e4 {ContentSigKeyID binary {}} \
    47e5 {ContentSigAlgo uinteger {}} \
    47e6 {ContentSigHashAlgo uinteger {}} \
]

set ebml_ContentEncAESSettings [dict create \
    47e8 {AESSettingsCipherMode uinteger {}} \
]

set ebml_Cues [dict create \
    bb {CuePoint master CuePoint} \
]

set ebml_CuePoint [dict create \
    b3 {CueTime uinteger {}} \
    b7 {CueTrackPositions master CueTrackPositions} \
]

set ebml_CueTrackPositions [dict create \
    f7 {CueTrack uinteger {}} \
    f1 {CueClusterPosition uinteger {}} \
    f0 {CueRelativePosition uinteger {}} \
    b2 {CueDuration uinteger {}} \
    5378 {CueBlockNumber uinteger {}} \
    ea {CueCodecState uinteger {}} \
    db {CueReference master CueReference} \
]

set ebml_CueReference [dict create \
    96 {CueRefTime uinteger {}} \
    97 {CueRefCluster uinteger {}} \
    535f {CueRefNumber uinteger {}} \
    eb {CueRefCodecState uinteger {}} \
]

set ebml_Attachments [dict create \
    61a7 {AttachedFile master AttachedFile} \
]

set ebml_AttachedFile [dict create \
    467e {FileDescription utf-8 {}} \
    466e {FileName utf-8 {}} \
    4660 {FileMimeType string {}} \
    465c {FileData binary {}} \
    46ae {FileUID uinteger {}} \
    4675 {FileReferral binary {}} \
    4661 {FileUsedStartTime uinteger {}} \
    4662 {FileUsedEndTime uinteger {}} \
]

set ebml_Chapters [dict create \
    45b9 {EditionEntry master EditionEntry} \
]

set ebml_EditionEntry [dict create \
    45bc {EditionUID uinteger {}} \
    45bd {EditionFlagHidden uinteger {}} \
    45db {EditionFlagDefault uinteger {}} \
    45dd {EditionFlagOrdered uinteger {}} \
    b6 {ChapterAtom master ChapterAtom} \
]

set ebml_ChapterAtom [dict create \
    73c4 {ChapterUID uinteger {}} \
    5654 {ChapterStringUID utf-8 {}} \
    91 {ChapterTimeStart uinteger {}} \
    92 {ChapterTimeEnd uinteger {}} \
    98 {ChapterFlagHidden uinteger {}} \
    4598 {ChapterFlagEnabled uinteger {}} \
    6e67 {ChapterSegmentUID binary {}} \
    6ebc {ChapterSegmentEditionUID uinteger {}} \
    63c3 {ChapterPhysicalEquiv uinteger {}} \
    8f {ChapterTrack master ChapterTrack} \
    80 {ChapterDisplay master ChapterDisplay} \
    6944 {ChapProcess master ChapProcess} \
]

set ebml_ChapterTrack [dict create \
    89 {ChapterTrackUID uinteger {}} \
]

set ebml_ChapterDisplay [dict create \
    85 {ChapString utf-8 {}} \
    437c {ChapLanguage string {}} \
    437d {ChapLanguageIETF string {}} \
    437e {ChapCountry string {}} \
]

set ebml_ChapProcess [dict create \
    6955 {ChapProcessCodecID uinteger {}} \
    450d {ChapProcessPrivate binary {}} \
    6911 {ChapProcessCommand master ChapProcessCommand} \
]

set ebml_ChapProcessCommand [dict create \
    6922 {ChapProcessTime uinteger {}} \
    6933 {ChapProcessData binary {}} \
]

set ebml_Tags [dict create \
    7373 {Tag master Tag} \
]

set ebml_Tag [dict create \
    63c0 {Targets master Targets} \
    67c8 {SimpleTag master SimpleTag} \
]

set ebml_Targets [dict create \
    68ca {TargetTypeValue uinteger {}} \
    63ca {TargetType string {}} \
    63c5 {TagTrackUID uinteger {}} \
    63c9 {TagEditionUID uinteger {}} \
    63c4 {TagChapterUID uinteger {}} \
    63c6 {TagAttachmentUID uinteger {}} \
]

set ebml_SimpleTag [dict create \
    45a3 {TagName utf-8 {}} \
    447a {TagLanguage string {}} \
    447b {TagLanguageIETF string {}} \
    4484 {TagDefault uinteger {}} \
    4487 {TagString utf-8 {}} \
    4485 {TagBinary binary {}} \
]

type_master [len] "" root

*/
