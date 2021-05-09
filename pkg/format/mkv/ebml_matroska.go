// Code below generated from ebml_matroska.xml
package mkv

import "fq/pkg/format/mkv/ebml"

const (
	EBMLMaxIDLength             = 0x42F2
	EBMLMaxSizeLength           = 0x42F3
	Segment                     = 0x18538067
	SeekHead                    = 0x114D9B74
	Seek                        = 0x4DBB
	SeekID                      = 0x53AB
	SeekPosition                = 0x53AC
	Info                        = 0x1549A966
	SegmentUID                  = 0x73A4
	SegmentFilename             = 0x7384
	PrevUID                     = 0x3CB923
	PrevFilename                = 0x3C83AB
	NextUID                     = 0x3EB923
	NextFilename                = 0x3E83BB
	SegmentFamily               = 0x4444
	ChapterTranslate            = 0x6924
	ChapterTranslateEditionUID  = 0x69FC
	ChapterTranslateCodec       = 0x69BF
	ChapterTranslateID          = 0x69A5
	TimestampScale              = 0x2AD7B1
	Duration                    = 0x4489
	DateUTC                     = 0x4461
	Title                       = 0x7BA9
	MuxingApp                   = 0x4D80
	WritingApp                  = 0x5741
	Cluster                     = 0x1F43B675
	Timestamp                   = 0xE7
	SilentTracks                = 0x5854
	SilentTrackNumber           = 0x58D7
	Position                    = 0xA7
	PrevSize                    = 0xAB
	SimpleBlock                 = 0xA3
	BlockGroup                  = 0xA0
	Block                       = 0xA1
	BlockVirtual                = 0xA2
	BlockAdditions              = 0x75A1
	BlockMore                   = 0xA6
	BlockAddID                  = 0xEE
	BlockAdditional             = 0xA5
	BlockDuration               = 0x9B
	ReferencePriority           = 0xFA
	ReferenceBlock              = 0xFB
	ReferenceVirtual            = 0xFD
	CodecState                  = 0xA4
	DiscardPadding              = 0x75A2
	Slices                      = 0x8E
	TimeSlice                   = 0xE8
	LaceNumber                  = 0xCC
	FrameNumber                 = 0xCD
	BlockAdditionID             = 0xCB
	Delay                       = 0xCE
	SliceDuration               = 0xCF
	ReferenceFrame              = 0xC8
	ReferenceOffset             = 0xC9
	ReferenceTimestamp          = 0xCA
	EncryptedBlock              = 0xAF
	Tracks                      = 0x1654AE6B
	TrackEntry                  = 0xAE
	TrackNumber                 = 0xD7
	TrackUID                    = 0x73C5
	TrackType                   = 0x83
	FlagEnabled                 = 0xB9
	FlagDefault                 = 0x88
	FlagForced                  = 0x55AA
	FlagLacing                  = 0x9C
	MinCache                    = 0x6DE7
	MaxCache                    = 0x6DF8
	DefaultDuration             = 0x23E383
	DefaultDecodedFieldDuration = 0x234E7A
	TrackTimestampScale         = 0x23314F
	TrackOffset                 = 0x537F
	MaxBlockAdditionID          = 0x55EE
	BlockAdditionMapping        = 0x41E4
	BlockAddIDValue             = 0x41F0
	BlockAddIDName              = 0x41A4
	BlockAddIDType              = 0x41E7
	BlockAddIDExtraData         = 0x41ED
	Name                        = 0x536E
	Language                    = 0x22B59C
	LanguageIETF                = 0x22B59D
	CodecID                     = 0x86
	CodecPrivate                = 0x63A2
	CodecName                   = 0x258688
	AttachmentLink              = 0x7446
	CodecSettings               = 0x3A9697
	CodecInfoURL                = 0x3B4040
	CodecDownloadURL            = 0x26B240
	CodecDecodeAll              = 0xAA
	TrackOverlay                = 0x6FAB
	CodecDelay                  = 0x56AA
	SeekPreRoll                 = 0x56BB
	TrackTranslate              = 0x6624
	TrackTranslateEditionUID    = 0x66FC
	TrackTranslateCodec         = 0x66BF
	TrackTranslateTrackID       = 0x66A5
	Video                       = 0xE0
	FlagInterlaced              = 0x9A
	FieldOrder                  = 0x9D
	StereoMode                  = 0x53B8
	AlphaMode                   = 0x53C0
	OldStereoMode               = 0x53B9
	PixelWidth                  = 0xB0
	PixelHeight                 = 0xBA
	PixelCropBottom             = 0x54AA
	PixelCropTop                = 0x54BB
	PixelCropLeft               = 0x54CC
	PixelCropRight              = 0x54DD
	DisplayWidth                = 0x54B0
	DisplayHeight               = 0x54BA
	DisplayUnit                 = 0x54B2
	AspectRatioType             = 0x54B3
	ColourSpace                 = 0x2EB524
	GammaValue                  = 0x2FB523
	FrameRate                   = 0x2383E3
	Colour                      = 0x55B0
	MatrixCoefficients          = 0x55B1
	BitsPerChannel              = 0x55B2
	ChromaSubsamplingHorz       = 0x55B3
	ChromaSubsamplingVert       = 0x55B4
	CbSubsamplingHorz           = 0x55B5
	CbSubsamplingVert           = 0x55B6
	ChromaSitingHorz            = 0x55B7
	ChromaSitingVert            = 0x55B8
	Range                       = 0x55B9
	TransferCharacteristics     = 0x55BA
	Primaries                   = 0x55BB
	MaxCLL                      = 0x55BC
	MaxFALL                     = 0x55BD
	MasteringMetadata           = 0x55D0
	PrimaryRChromaticityX       = 0x55D1
	PrimaryRChromaticityY       = 0x55D2
	PrimaryGChromaticityX       = 0x55D3
	PrimaryGChromaticityY       = 0x55D4
	PrimaryBChromaticityX       = 0x55D5
	PrimaryBChromaticityY       = 0x55D6
	WhitePointChromaticityX     = 0x55D7
	WhitePointChromaticityY     = 0x55D8
	LuminanceMax                = 0x55D9
	LuminanceMin                = 0x55DA
	Projection                  = 0x7670
	ProjectionType              = 0x7671
	ProjectionPrivate           = 0x7672
	ProjectionPoseYaw           = 0x7673
	ProjectionPosePitch         = 0x7674
	ProjectionPoseRoll          = 0x7675
	Audio                       = 0xE1
	SamplingFrequency           = 0xB5
	OutputSamplingFrequency     = 0x78B5
	Channels                    = 0x9F
	ChannelPositions            = 0x7D7B
	BitDepth                    = 0x6264
	TrackOperation              = 0xE2
	TrackCombinePlanes          = 0xE3
	TrackPlane                  = 0xE4
	TrackPlaneUID               = 0xE5
	TrackPlaneType              = 0xE6
	TrackJoinBlocks             = 0xE9
	TrackJoinUID                = 0xED
	TrickTrackUID               = 0xC0
	TrickTrackSegmentUID        = 0xC1
	TrickTrackFlag              = 0xC6
	TrickMasterTrackUID         = 0xC7
	TrickMasterTrackSegmentUID  = 0xC4
	ContentEncodings            = 0x6D80
	ContentEncoding             = 0x6240
	ContentEncodingOrder        = 0x5031
	ContentEncodingScope        = 0x5032
	ContentEncodingType         = 0x5033
	ContentCompression          = 0x5034
	ContentCompAlgo             = 0x4254
	ContentCompSettings         = 0x4255
	ContentEncryption           = 0x5035
	ContentEncAlgo              = 0x47E1
	ContentEncKeyID             = 0x47E2
	ContentEncAESSettings       = 0x47E7
	AESSettingsCipherMode       = 0x47E8
	ContentSignature            = 0x47E3
	ContentSigKeyID             = 0x47E4
	ContentSigAlgo              = 0x47E5
	ContentSigHashAlgo          = 0x47E6
	Cues                        = 0x1C53BB6B
	CuePoint                    = 0xBB
	CueTime                     = 0xB3
	CueTrackPositions           = 0xB7
	CueTrack                    = 0xF7
	CueClusterPosition          = 0xF1
	CueRelativePosition         = 0xF0
	CueDuration                 = 0xB2
	CueBlockNumber              = 0x5378
	CueCodecState               = 0xEA
	CueReference                = 0xDB
	CueRefTime                  = 0x96
	CueRefCluster               = 0x97
	CueRefNumber                = 0x535F
	CueRefCodecState            = 0xEB
	Attachments                 = 0x1941A469
	AttachedFile                = 0x61A7
	FileDescription             = 0x467E
	FileName                    = 0x466E
	FileMimeType                = 0x4660
	FileData                    = 0x465C
	FileUID                     = 0x46AE
	FileReferral                = 0x4675
	FileUsedStartTime           = 0x4661
	FileUsedEndTime             = 0x4662
	Chapters                    = 0x1043A770
	EditionEntry                = 0x45B9
	EditionUID                  = 0x45BC
	EditionFlagHidden           = 0x45BD
	EditionFlagDefault          = 0x45DB
	EditionFlagOrdered          = 0x45DD
	ChapterAtom                 = 0xB6
	ChapterUID                  = 0x73C4
	ChapterStringUID            = 0x5654
	ChapterTimeStart            = 0x91
	ChapterTimeEnd              = 0x92
	ChapterFlagHidden           = 0x98
	ChapterFlagEnabled          = 0x4598
	ChapterSegmentUID           = 0x6E67
	ChapterSegmentEditionUID    = 0x6EBC
	ChapterPhysicalEquiv        = 0x63C3
	ChapterTrack                = 0x8F
	ChapterTrackUID             = 0x89
	ChapterDisplay              = 0x80
	ChapString                  = 0x85
	ChapLanguage                = 0x437C
	ChapLanguageIETF            = 0x437D
	ChapCountry                 = 0x437E
	ChapProcess                 = 0x6944
	ChapProcessCodecID          = 0x6955
	ChapProcessPrivate          = 0x450D
	ChapProcessCommand          = 0x6911
	ChapProcessTime             = 0x6922
	ChapProcessData             = 0x6933
	Tags                        = 0x1254C367
	Tag                         = 0x7373
	Targets                     = 0x63C0
	TargetTypeValue             = 0x68CA
	TargetType                  = 0x63CA
	TagTrackUID                 = 0x63C5
	TagEditionUID               = 0x63C9
	TagChapterUID               = 0x63C4
	TagAttachmentUID            = 0x63C6
	SimpleTag                   = 0x67C8
	TagName                     = 0x45A3
	TagLanguage                 = 0x447A
	TagLanguageIETF             = 0x447B
	TagDefault                  = 0x4484
	TagString                   = 0x4487
	TagBinary                   = 0x4485
)

var mkvSegment = ebml.Tag{
	SeekHead: {
		Name:       "SeekHead",
		Definition: "Contains the Segment Position of other Top-Level Elements.",
		Type:       ebml.Master, Tag: mkvSeekHead,
	},
	Info: {
		Name:       "Info",
		Definition: "Contains general information about the Segment.",
		Type:       ebml.Master, Tag: mkvInfo,
	},
	Cluster: {
		Name:       "Cluster",
		Definition: "The Top-Level Element containing the (monolithic) Block structure.",
		Type:       ebml.Master, Tag: mkvCluster,
	},
	Tracks: {
		Name:       "Tracks",
		Definition: "A Top-Level Element of information with many tracks described.",
		Type:       ebml.Master, Tag: mkvTracks,
	},
	Cues: {
		Name:       "Cues",
		Definition: "A Top-Level Element to speed seeking access. All entries are local to the Segment.",
		Type:       ebml.Master, Tag: mkvCues,
	},
	Attachments: {
		Name:       "Attachments",
		Definition: "Contain attached files.",
		Type:       ebml.Master, Tag: mkvAttachments,
	},
	Chapters: {
		Name:       "Chapters",
		Definition: "A system to define basic menus and partition data. For more detailed information, look at the .",
		Type:       ebml.Master, Tag: mkvChapters,
	},
	Tags: {
		Name:       "Tags",
		Definition: "Element containing metadata describing Tracks, Editions, Chapters, Attachments, or the Segment as a whole. A list of valid tags can be found",
		Type:       ebml.Master, Tag: mkvTags,
	},
}

var mkvSeekHead = ebml.Tag{
	Seek: {
		Name:       "Seek",
		Definition: "Contains a single seek entry to an EBML Element.",
		Type:       ebml.Master, Tag: mkvSeek,
	},
}

var mkvSeek = ebml.Tag{
	SeekID: {
		Name:       "SeekID",
		Definition: "The binary ID corresponding to the Element name.",
		Type:       ebml.Binary,
	},
	SeekPosition: {
		Name:       "SeekPosition",
		Definition: "The Segment Position of the Element.",
		Type:       ebml.Uinteger,
	},
}

var mkvInfo = ebml.Tag{
	SegmentUID: {
		Name:       "SegmentUID",
		Definition: "A randomly generated unique ID to identify the Segment amongst many others (128 bits).",
		Type:       ebml.Binary,
	},
	SegmentFilename: {
		Name:       "SegmentFilename",
		Definition: "A filename corresponding to this Segment.",
		Type:       ebml.UTF8,
	},
	PrevUID: {
		Name:       "PrevUID",
		Definition: "A unique ID to identify the previous Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	PrevFilename: {
		Name:       "PrevFilename",
		Definition: "A filename corresponding to the file of the previous Linked Segment.",
		Type:       ebml.UTF8,
	},
	NextUID: {
		Name:       "NextUID",
		Definition: "A unique ID to identify the next Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	NextFilename: {
		Name:       "NextFilename",
		Definition: "A filename corresponding to the file of the next Linked Segment.",
		Type:       ebml.UTF8,
	},
	SegmentFamily: {
		Name:       "SegmentFamily",
		Definition: "A randomly generated unique ID that all Segments of a Linked Segment MUST share (128 bits).",
		Type:       ebml.Binary,
	},
	ChapterTranslate: {
		Name:       "ChapterTranslate",
		Definition: "A tuple of corresponding ID used by chapter codecs to represent this Segment.",
		Type:       ebml.Master, Tag: mkvChapterTranslate,
	},
	TimestampScale: {
		Name:       "TimestampScale",
		Definition: "Timestamp scale in nanoseconds (1.000.000 means all timestamps in the Segment are expressed in milliseconds).",
		Type:       ebml.Uinteger,
	},
	Duration: {
		Name:       "Duration",
		Definition: "Duration of the Segment in nanoseconds based on TimestampScale.",
		Type:       ebml.Float,
	},
	DateUTC: {
		Name:       "DateUTC",
		Definition: "The date and time that the Segment was created by the muxing application or library.",
		Type:       ebml.Date,
	},
	Title: {
		Name:       "Title",
		Definition: "General name of the Segment.",
		Type:       ebml.UTF8,
	},
	MuxingApp: {
		Name:       "MuxingApp",
		Definition: "Muxing application or library (example: \"libmatroska-0.4.3\").",
		Type:       ebml.UTF8,
	},
	WritingApp: {
		Name:       "WritingApp",
		Definition: "Writing application (example: \"mkvmerge-0.3.3\").",
		Type:       ebml.UTF8,
	},
}

var mkvChapterTranslate = ebml.Tag{
	ChapterTranslateEditionUID: {
		Name:       "ChapterTranslateEditionUID",
		Definition: "Specify an edition UID on which this correspondence applies. When not specified, it means for all editions found in the Segment.",
		Type:       ebml.Uinteger,
	},
	ChapterTranslateCodec: {
		Name:       "ChapterTranslateCodec",
		Definition: "The",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "Matroska Script",
			},
			1: {
				Value: "1",
				Label: "DVD-menu",
			},
		},
	},
	ChapterTranslateID: {
		Name:       "ChapterTranslateID",
		Definition: "The binary value used to represent this Segment in the chapter codec data. The format depends on the  used.",
		Type:       ebml.Binary,
	},
}

var mkvCluster = ebml.Tag{
	Timestamp: {
		Name:       "Timestamp",
		Definition: "Absolute timestamp of the cluster (based on TimestampScale).",
		Type:       ebml.Uinteger,
	},
	SilentTracks: {
		Name:       "SilentTracks",
		Definition: "The list of tracks that are not used in that part of the stream. It is useful when using overlay tracks on seeking or to decide what track to use.",
		Type:       ebml.Master, Tag: mkvSilentTracks,
	},
	Position: {
		Name:       "Position",
		Definition: "The Segment Position of the Cluster in the Segment (0 in live streams). It might help to resynchronise offset on damaged streams.",
		Type:       ebml.Uinteger,
	},
	PrevSize: {
		Name:       "PrevSize",
		Definition: "Size of the previous Cluster, in octets. Can be useful for backward playing.",
		Type:       ebml.Uinteger,
	},
	SimpleBlock: {
		Name:       "SimpleBlock",
		Definition: "Similar to  but without all the extra information, mostly used to reduced overhead when no extra feature is needed. (see )",
		Type:       ebml.Binary,
	},
	BlockGroup: {
		Name:       "BlockGroup",
		Definition: "Basic container of information containing a single Block and information specific to that Block.",
		Type:       ebml.Master, Tag: mkvBlockGroup,
	},
	EncryptedBlock: {
		Name:       "EncryptedBlock",
		Definition: "Similar to  but the data inside the Block are Transformed (encrypt and/or signed). (see )",
		Type:       ebml.Binary,
	},
}

var mkvSilentTracks = ebml.Tag{
	SilentTrackNumber: {
		Name:       "SilentTrackNumber",
		Definition: "One of the track number that are not used from now on in the stream. It could change later if not specified as silent in a further Cluster.",
		Type:       ebml.Uinteger,
	},
}

var mkvBlockGroup = ebml.Tag{
	Block: {
		Name:       "Block",
		Definition: "Block containing the actual data to be rendered and a timestamp relative to the Cluster Timestamp. (see )",
		Type:       ebml.Binary,
	},
	BlockVirtual: {
		Name:       "BlockVirtual",
		Definition: "A Block with no data. It MUST be stored in the stream at the place the real Block would be in display order. (see )",
		Type:       ebml.Binary,
	},
	BlockAdditions: {
		Name:       "BlockAdditions",
		Definition: "Contain additional blocks to complete the main one. An EBML parser that has no knowledge of the Block structure could still see and use/skip these data.",
		Type:       ebml.Master, Tag: mkvBlockAdditions,
	},
	BlockDuration: {
		Name:       "BlockDuration",
		Definition: "The duration of the Block (based on TimestampScale). The BlockDuration Element can be useful at the end of a Track to define the duration of the last frame (as there is no subsequent Block available), or when there is a break in a track like for subtitle tracks.",
		Type:       ebml.Uinteger,
	},
	ReferencePriority: {
		Name:       "ReferencePriority",
		Definition: "This frame is referenced and has the specified cache priority. In cache only a frame of the same or higher priority can replace this frame. A value of 0 means the frame is not referenced.",
		Type:       ebml.Uinteger,
	},
	ReferenceBlock: {
		Name:       "ReferenceBlock",
		Definition: "Timestamp of another frame used as a reference (ie: B or P frame). The timestamp is relative to the block it's attached to.",
		Type:       ebml.Integer,
	},
	ReferenceVirtual: {
		Name:       "ReferenceVirtual",
		Definition: "The Segment Position of the data that would otherwise be in position of the virtual block.",
		Type:       ebml.Integer,
	},
	CodecState: {
		Name:       "CodecState",
		Definition: "The new codec state to use. Data interpretation is private to the codec. This information SHOULD always be referenced by a seek entry.",
		Type:       ebml.Binary,
	},
	DiscardPadding: {
		Name:       "DiscardPadding",
		Definition: "Duration in nanoseconds of the silent data added to the Block (padding at the end of the Block for positive value, at the beginning of the Block for negative value). The duration of DiscardPadding is not calculated in the duration of the TrackEntry and SHOULD be discarded during playback.",
		Type:       ebml.Integer,
	},
	Slices: {
		Name:       "Slices",
		Definition: "Contains slices description.",
		Type:       ebml.Master, Tag: mkvSlices,
	},
	ReferenceFrame: {
		Name:       "ReferenceFrame",
		Definition: "",
		Type:       ebml.Master, Tag: mkvReferenceFrame,
	},
}

var mkvBlockAdditions = ebml.Tag{
	BlockMore: {
		Name:       "BlockMore",
		Definition: "Contain the BlockAdditional and some parameters.",
		Type:       ebml.Master, Tag: mkvBlockMore,
	},
}

var mkvBlockMore = ebml.Tag{
	BlockAddID: {
		Name:       "BlockAddID",
		Definition: "An ID to identify the BlockAdditional level. A value of 1 means the BlockAdditional data is interpreted as additional data passed to the codec with the Block data.",
		Type:       ebml.Uinteger,
	},
	BlockAdditional: {
		Name:       "BlockAdditional",
		Definition: "Interpreted by the codec as it wishes (using the BlockAddID).",
		Type:       ebml.Binary,
	},
}

var mkvSlices = ebml.Tag{
	TimeSlice: {
		Name:       "TimeSlice",
		Definition: "Contains extra time information about the data contained in the Block. Being able to interpret this Element is not REQUIRED for playback.",
		Type:       ebml.Master, Tag: mkvTimeSlice,
	},
}

var mkvTimeSlice = ebml.Tag{
	LaceNumber: {
		Name:       "LaceNumber",
		Definition: "The reverse number of the frame in the lace (0 is the last frame, 1 is the next to last, etc). Being able to interpret this Element is not REQUIRED for playback.",
		Type:       ebml.Uinteger,
	},
	FrameNumber: {
		Name:       "FrameNumber",
		Definition: "The number of the frame to generate from this lace with this delay (allow you to generate many frames from the same Block/Frame).",
		Type:       ebml.Uinteger,
	},
	BlockAdditionID: {
		Name:       "BlockAdditionID",
		Definition: "The ID of the BlockAdditional Element (0 is the main Block).",
		Type:       ebml.Uinteger,
	},
	Delay: {
		Name:       "Delay",
		Definition: "The (scaled) delay to apply to the Element.",
		Type:       ebml.Uinteger,
	},
	SliceDuration: {
		Name:       "SliceDuration",
		Definition: "The (scaled) duration to apply to the Element.",
		Type:       ebml.Uinteger,
	},
}

var mkvReferenceFrame = ebml.Tag{
	ReferenceOffset: {
		Name:       "ReferenceOffset",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	ReferenceTimestamp: {
		Name:       "ReferenceTimestamp",
		Definition: "",
		Type:       ebml.Uinteger,
	},
}

var mkvTracks = ebml.Tag{
	TrackEntry: {
		Name:       "TrackEntry",
		Definition: "Describes a track with all Elements.",
		Type:       ebml.Master, Tag: mkvTrackEntry,
	},
}

var mkvTrackEntry = ebml.Tag{
	TrackNumber: {
		Name:       "TrackNumber",
		Definition: "The track number as used in the Block Header (using more than 127 tracks is not encouraged, though the design allows an unlimited number).",
		Type:       ebml.Uinteger,
	},
	TrackUID: {
		Name:       "TrackUID",
		Definition: "A unique ID to identify the Track. This SHOULD be kept the same when making a direct stream copy of the Track to another file.",
		Type:       ebml.Uinteger,
	},
	TrackType: {
		Name:       "TrackType",
		Definition: "A set of track types coded on 8 bits.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			1: {
				Value: "1",
				Label: "video",
			},
			2: {
				Value: "2",
				Label: "audio",
			},
			3: {
				Value: "3",
				Label: "complex",
			},
			16: {
				Value: "16",
				Label: "logo",
			},
			17: {
				Value: "17",
				Label: "subtitle",
			},
			18: {
				Value: "18",
				Label: "buttons",
			},
			32: {
				Value: "32",
				Label: "control",
			},
			33: {
				Value: "33",
				Label: "metadata",
			},
		},
	},
	FlagEnabled: {
		Name:       "FlagEnabled",
		Definition: "Set if the track is usable. (1 bit)",
		Type:       ebml.Uinteger,
	},
	FlagDefault: {
		Name:       "FlagDefault",
		Definition: "Set if that track (audio, video or subs) SHOULD be active if no language found matches the user preference. (1 bit)",
		Type:       ebml.Uinteger,
	},
	FlagForced: {
		Name:       "FlagForced",
		Definition: "Set if that track MUST be active during playback. There can be many forced track for a kind (audio, video or subs), the player SHOULD select the one which language matches the user preference or the default + forced track. Overlay MAY happen between a forced and non-forced track of the same kind. (1 bit)",
		Type:       ebml.Uinteger,
	},
	FlagLacing: {
		Name:       "FlagLacing",
		Definition: "Set if the track MAY contain blocks using lacing. (1 bit)",
		Type:       ebml.Uinteger,
	},
	MinCache: {
		Name:       "MinCache",
		Definition: "The minimum number of frames a player SHOULD be able to cache during playback. If set to 0, the reference pseudo-cache system is not used.",
		Type:       ebml.Uinteger,
	},
	MaxCache: {
		Name:       "MaxCache",
		Definition: "The maximum cache size necessary to store referenced frames in and the current frame. 0 means no cache is needed.",
		Type:       ebml.Uinteger,
	},
	DefaultDuration: {
		Name:       "DefaultDuration",
		Definition: "Number of nanoseconds (not scaled via TimestampScale) per frame ('frame' in the Matroska sense -- one Element put into a (Simple)Block).",
		Type:       ebml.Uinteger,
	},
	DefaultDecodedFieldDuration: {
		Name:       "DefaultDecodedFieldDuration",
		Definition: "The period in nanoseconds (not scaled by TimestampScale) between two successive fields at the output of the decoding process (see )",
		Type:       ebml.Uinteger,
	},
	TrackTimestampScale: {
		Name:       "TrackTimestampScale",
		Definition: "DEPRECATED, DO NOT USE. The scale to apply on this track to work at normal speed in relation with other tracks (mostly used to adjust video speed when the audio length differs).",
		Type:       ebml.Float,
	},
	TrackOffset: {
		Name:       "TrackOffset",
		Definition: "A value to add to the Block's Timestamp. This can be used to adjust the playback offset of a track.",
		Type:       ebml.Integer,
	},
	MaxBlockAdditionID: {
		Name:       "MaxBlockAdditionID",
		Definition: "The maximum value of . A value 0 means there is no  for this track.",
		Type:       ebml.Uinteger,
	},
	BlockAdditionMapping: {
		Name:       "BlockAdditionMapping",
		Definition: "Contains elements that describe each value of  found in the Track.",
		Type:       ebml.Master, Tag: mkvBlockAdditionMapping,
	},
	Name: {
		Name:       "Name",
		Definition: "A human-readable track name.",
		Type:       ebml.UTF8,
	},
	Language: {
		Name:       "Language",
		Definition: "Specifies the language of the track in the . This Element MUST be ignored if the LanguageIETF Element is used in the same TrackEntry.",
		Type:       ebml.String,
	},
	LanguageIETF: {
		Name:       "LanguageIETF",
		Definition: "Specifies the language of the track according to  and using the . If this Element is used, then any Language Elements used in the same TrackEntry MUST be ignored.",
		Type:       ebml.String,
	},
	CodecID: {
		Name:       "CodecID",
		Definition: "An ID corresponding to the codec, see the  for more info.",
		Type:       ebml.String,
	},
	CodecPrivate: {
		Name:       "CodecPrivate",
		Definition: "Private data only known to the codec.",
		Type:       ebml.Binary,
	},
	CodecName: {
		Name:       "CodecName",
		Definition: "A human-readable string specifying the codec.",
		Type:       ebml.UTF8,
	},
	AttachmentLink: {
		Name:       "AttachmentLink",
		Definition: "The UID of an attachment that is used by this codec.",
		Type:       ebml.Uinteger,
	},
	CodecSettings: {
		Name:       "CodecSettings",
		Definition: "A string describing the encoding setting used.",
		Type:       ebml.UTF8,
	},
	CodecInfoURL: {
		Name:       "CodecInfoURL",
		Definition: "A URL to find information about the codec used.",
		Type:       ebml.String,
	},
	CodecDownloadURL: {
		Name:       "CodecDownloadURL",
		Definition: "A URL to download about the codec used.",
		Type:       ebml.String,
	},
	CodecDecodeAll: {
		Name:       "CodecDecodeAll",
		Definition: "The codec can decode potentially damaged data (1 bit).",
		Type:       ebml.Uinteger,
	},
	TrackOverlay: {
		Name:       "TrackOverlay",
		Definition: "Specify that this track is an overlay track for the Track specified (in the u-integer). That means when this track has a gap (see ) the overlay track SHOULD be used instead. The order of multiple TrackOverlay matters, the first one is the one that SHOULD be used. If not found it SHOULD be the second, etc.",
		Type:       ebml.Uinteger,
	},
	CodecDelay: {
		Name:       "CodecDelay",
		Definition: "CodecDelay is The codec-built-in delay in nanoseconds. This value MUST be subtracted from each block timestamp in order to get the actual timestamp. The value SHOULD be small so the muxing of tracks with the same actual timestamp are in the same Cluster.",
		Type:       ebml.Uinteger,
	},
	SeekPreRoll: {
		Name:       "SeekPreRoll",
		Definition: "After a discontinuity, SeekPreRoll is the duration in nanoseconds of the data the decoder MUST decode before the decoded data is valid.",
		Type:       ebml.Uinteger,
	},
	TrackTranslate: {
		Name:       "TrackTranslate",
		Definition: "The track identification for the given Chapter Codec.",
		Type:       ebml.Master, Tag: mkvTrackTranslate,
	},
	Video: {
		Name:       "Video",
		Definition: "Video settings.",
		Type:       ebml.Master, Tag: mkvVideo,
	},
	Audio: {
		Name:       "Audio",
		Definition: "Audio settings.",
		Type:       ebml.Master, Tag: mkvAudio,
	},
	TrackOperation: {
		Name:       "TrackOperation",
		Definition: "Operation that needs to be applied on tracks to create this virtual track. For more details  on the subject.",
		Type:       ebml.Master, Tag: mkvTrackOperation,
	},
	TrickTrackUID: {
		Name:       "TrickTrackUID",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	TrickTrackSegmentUID: {
		Name:       "TrickTrackSegmentUID",
		Definition: "",
		Type:       ebml.Binary,
	},
	TrickTrackFlag: {
		Name:       "TrickTrackFlag",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	TrickMasterTrackUID: {
		Name:       "TrickMasterTrackUID",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	TrickMasterTrackSegmentUID: {
		Name:       "TrickMasterTrackSegmentUID",
		Definition: "",
		Type:       ebml.Binary,
	},
	ContentEncodings: {
		Name:       "ContentEncodings",
		Definition: "Settings for several content encoding mechanisms like compression or encryption.",
		Type:       ebml.Master, Tag: mkvContentEncodings,
	},
}

var mkvBlockAdditionMapping = ebml.Tag{
	BlockAddIDValue: {
		Name:       "BlockAddIDValue",
		Definition: "The  value being described. To keep MaxBlockAdditionID as low as possible, small values SHOULD be used.",
		Type:       ebml.Uinteger,
	},
	BlockAddIDName: {
		Name:       "BlockAddIDName",
		Definition: "A human-friendly name describing the type of BlockAdditional data as defined by the associated Block Additional Mapping.",
		Type:       ebml.String,
	},
	BlockAddIDType: {
		Name:       "BlockAddIDType",
		Definition: "Stores the registered identifer of the Block Additional Mapping to define how the BlockAdditional data should be handled.",
		Type:       ebml.Uinteger,
	},
	BlockAddIDExtraData: {
		Name:       "BlockAddIDExtraData",
		Definition: "Extra binary data that the BlockAddIDType can use to interpret the BlockAdditional data. The intepretation of the binary data depends on the BlockAddIDType value and the corresponding Block Additional Mapping.",
		Type:       ebml.Binary,
	},
}

var mkvTrackTranslate = ebml.Tag{
	TrackTranslateEditionUID: {
		Name:       "TrackTranslateEditionUID",
		Definition: "Specify an edition UID on which this translation applies. When not specified, it means for all editions found in the Segment.",
		Type:       ebml.Uinteger,
	},
	TrackTranslateCodec: {
		Name:       "TrackTranslateCodec",
		Definition: "The .",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "Matroska Script",
			},
			1: {
				Value: "1",
				Label: "DVD-menu",
			},
		},
	},
	TrackTranslateTrackID: {
		Name:       "TrackTranslateTrackID",
		Definition: "The binary value used to represent this track in the chapter codec data. The format depends on the  used.",
		Type:       ebml.Binary,
	},
}

var mkvVideo = ebml.Tag{
	FlagInterlaced: {
		Name:       "FlagInterlaced",
		Definition: "A flag to declare if the video is known to be progressive or interlaced and if applicable to declare details about the interlacement.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "undetermined",
			},
			1: {
				Value: "1",
				Label: "interlaced",
			},
			2: {
				Value: "2",
				Label: "progressive",
			},
		},
	},
	FieldOrder: {
		Name:       "FieldOrder",
		Definition: "Declare the field ordering of the video. If FlagInterlaced is not set to 1, this Element MUST be ignored.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "progressive",
			},
			1: {
				Value:      "1",
				Label:      "tff",
				Definition: "Top field displayed first. Top field stored first.",
			},
			2: {
				Value: "2",
				Label: "undetermined",
			},
			6: {
				Value:      "6",
				Label:      "bff",
				Definition: "Bottom field displayed first. Bottom field stored first.",
			},
			9: {
				Value:      "9",
				Label:      "bff(swapped)",
				Definition: "Top field displayed first. Fields are interleaved in storage with the top line of the top field stored first.",
			},
			14: {
				Value:      "14",
				Label:      "tff(swapped)",
				Definition: "Bottom field displayed first. Fields are interleaved in storage with the top line of the top field stored first.",
			},
		},
	},
	StereoMode: {
		Name:       "StereoMode",
		Definition: "Stereo-3D video mode. There are some more details on .",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "mono",
			},
			1: {
				Value: "1",
				Label: "side by side (left eye first)",
			},
			2: {
				Value: "2",
				Label: "top - bottom (right eye is first)",
			},
			3: {
				Value: "3",
				Label: "top - bottom (left eye is first)",
			},
			4: {
				Value: "4",
				Label: "checkboard (right eye is first)",
			},
			5: {
				Value: "5",
				Label: "checkboard (left eye is first)",
			},
			6: {
				Value: "6",
				Label: "row interleaved (right eye is first)",
			},
			7: {
				Value: "7",
				Label: "row interleaved (left eye is first)",
			},
			8: {
				Value: "8",
				Label: "column interleaved (right eye is first)",
			},
			9: {
				Value: "9",
				Label: "column interleaved (left eye is first)",
			},
			10: {
				Value: "10",
				Label: "anaglyph (cyan/red)",
			},
			11: {
				Value: "11",
				Label: "side by side (right eye first)",
			},
			12: {
				Value: "12",
				Label: "anaglyph (green/magenta)",
			},
			13: {
				Value: "13",
				Label: "both eyes laced in one Block (left eye is first)",
			},
			14: {
				Value: "14",
				Label: "both eyes laced in one Block (right eye is first)",
			},
		},
	},
	AlphaMode: {
		Name:       "AlphaMode",
		Definition: "Alpha Video Mode. Presence of this Element indicates that the BlockAdditional Element could contain Alpha data.",
		Type:       ebml.Uinteger,
	},
	OldStereoMode: {
		Name:       "OldStereoMode",
		Definition: "DEPRECATED, DO NOT USE. Bogus StereoMode value used in old versions of libmatroska.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "mono",
			},
			1: {
				Value: "1",
				Label: "right eye",
			},
			2: {
				Value: "2",
				Label: "left eye",
			},
			3: {
				Value: "3",
				Label: "both eyes",
			},
		},
	},
	PixelWidth: {
		Name:       "PixelWidth",
		Definition: "Width of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	PixelHeight: {
		Name:       "PixelHeight",
		Definition: "Height of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	PixelCropBottom: {
		Name:       "PixelCropBottom",
		Definition: "The number of video pixels to remove at the bottom of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropTop: {
		Name:       "PixelCropTop",
		Definition: "The number of video pixels to remove at the top of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropLeft: {
		Name:       "PixelCropLeft",
		Definition: "The number of video pixels to remove on the left of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropRight: {
		Name:       "PixelCropRight",
		Definition: "The number of video pixels to remove on the right of the image.",
		Type:       ebml.Uinteger,
	},
	DisplayWidth: {
		Name:       "DisplayWidth",
		Definition: "Width of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	DisplayHeight: {
		Name:       "DisplayHeight",
		Definition: "Height of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	DisplayUnit: {
		Name:       "DisplayUnit",
		Definition: "How DisplayWidth & DisplayHeight are interpreted.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "pixels",
			},
			1: {
				Value: "1",
				Label: "centimeters",
			},
			2: {
				Value: "2",
				Label: "inches",
			},
			3: {
				Value: "3",
				Label: "display aspect ratio",
			},
			4: {
				Value: "4",
				Label: "unknown",
			},
		},
	},
	AspectRatioType: {
		Name:       "AspectRatioType",
		Definition: "Specify the possible modifications to the aspect ratio.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "free resizing",
			},
			1: {
				Value: "1",
				Label: "keep aspect ratio",
			},
			2: {
				Value: "2",
				Label: "fixed",
			},
		},
	},
	ColourSpace: {
		Name:       "ColourSpace",
		Definition: "Specify the pixel format used for the Track's data as a FourCC. This value is similar in scope to the biCompression value of AVI's BITMAPINFOHEADER.",
		Type:       ebml.Binary,
	},
	GammaValue: {
		Name:       "GammaValue",
		Definition: "Gamma Value.",
		Type:       ebml.Float,
	},
	FrameRate: {
		Name:       "FrameRate",
		Definition: "Number of frames per second.  only.",
		Type:       ebml.Float,
	},
	Colour: {
		Name:       "Colour",
		Definition: "Settings describing the colour format.",
		Type:       ebml.Master, Tag: mkvColour,
	},
	Projection: {
		Name:       "Projection",
		Definition: "Describes the video projection details. Used to render spherical and VR videos.",
		Type:       ebml.Master, Tag: mkvProjection,
	},
}

var mkvColour = ebml.Tag{
	MatrixCoefficients: {
		Name:       "MatrixCoefficients",
		Definition: "The Matrix Coefficients of the video used to derive luma and chroma values from red, green, and blue color primaries. For clarity, the value and meanings for MatrixCoefficients are adopted from Table 4 of ISO/IEC 23001-8:2016 or ITU-T H.273.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "Identity",
			},
			1: {
				Value: "1",
				Label: "ITU-R BT.709",
			},
			2: {
				Value: "2",
				Label: "unspecified",
			},
			3: {
				Value: "3",
				Label: "reserved",
			},
			4: {
				Value: "4",
				Label: "US FCC 73.682",
			},
			5: {
				Value: "5",
				Label: "ITU-R BT.470BG",
			},
			6: {
				Value: "6",
				Label: "SMPTE 170M",
			},
			7: {
				Value: "7",
				Label: "SMPTE 240M",
			},
			8: {
				Value: "8",
				Label: "YCoCg",
			},
			9: {
				Value: "9",
				Label: "BT2020 Non-constant Luminance",
			},
			10: {
				Value: "10",
				Label: "BT2020 Constant Luminance",
			},
			11: {
				Value: "11",
				Label: "SMPTE ST 2085",
			},
			12: {
				Value: "12",
				Label: "Chroma-derived Non-constant Luminance",
			},
			13: {
				Value: "13",
				Label: "Chroma-derived Constant Luminance",
			},
			14: {
				Value: "14",
				Label: "ITU-R BT.2100-0",
			},
		},
	},
	BitsPerChannel: {
		Name:       "BitsPerChannel",
		Definition: "Number of decoded bits per channel. A value of 0 indicates that the BitsPerChannel is unspecified.",
		Type:       ebml.Uinteger,
	},
	ChromaSubsamplingHorz: {
		Name:       "ChromaSubsamplingHorz",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed horizontally. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingHorz SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	ChromaSubsamplingVert: {
		Name:       "ChromaSubsamplingVert",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed vertically. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingVert SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	CbSubsamplingHorz: {
		Name:       "CbSubsamplingHorz",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed horizontally. This is additive with ChromaSubsamplingHorz. Example: For video with 4:2:1 chroma subsampling, the ChromaSubsamplingHorz SHOULD be set to 1 and CbSubsamplingHorz SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	CbSubsamplingVert: {
		Name:       "CbSubsamplingVert",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed vertically. This is additive with ChromaSubsamplingVert.",
		Type:       ebml.Uinteger,
	},
	ChromaSitingHorz: {
		Name:       "ChromaSitingHorz",
		Definition: "How chroma is subsampled horizontally.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "unspecified",
			},
			1: {
				Value: "1",
				Label: "left collocated",
			},
			2: {
				Value: "2",
				Label: "half",
			},
		},
	},
	ChromaSitingVert: {
		Name:       "ChromaSitingVert",
		Definition: "How chroma is subsampled vertically.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "unspecified",
			},
			1: {
				Value: "1",
				Label: "top collocated",
			},
			2: {
				Value: "2",
				Label: "half",
			},
		},
	},
	Range: {
		Name:       "Range",
		Definition: "Clipping of the color ranges.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "unspecified",
			},
			1: {
				Value: "1",
				Label: "broadcast range",
			},
			2: {
				Value: "2",
				Label: "full range (no clipping)",
			},
			3: {
				Value: "3",
				Label: "defined by MatrixCoefficients / TransferCharacteristics",
			},
		},
	},
	TransferCharacteristics: {
		Name:       "TransferCharacteristics",
		Definition: "The transfer characteristics of the video. For clarity, the value and meanings for TransferCharacteristics are adopted from Table 3 of  ISO/IEC 23091-4 or ITU-T H.273.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "reserved",
			},
			1: {
				Value: "1",
				Label: "ITU-R BT.709",
			},
			2: {
				Value: "2",
				Label: "unspecified",
			},
			3: {
				Value: "3",
				Label: "reserved",
			},
			4: {
				Value: "4",
				Label: "Gamma 2.2 curve - BT.470M",
			},
			5: {
				Value: "5",
				Label: "Gamma 2.8 curve - BT.470BG",
			},
			6: {
				Value: "6",
				Label: "SMPTE 170M",
			},
			7: {
				Value: "7",
				Label: "SMPTE 240M",
			},
			8: {
				Value: "8",
				Label: "Linear",
			},
			9: {
				Value: "9",
				Label: "Log",
			},
			10: {
				Value: "10",
				Label: "Log Sqrt",
			},
			11: {
				Value: "11",
				Label: "IEC 61966-2-4",
			},
			12: {
				Value: "12",
				Label: "ITU-R BT.1361 Extended Colour Gamut",
			},
			13: {
				Value: "13",
				Label: "IEC 61966-2-1",
			},
			14: {
				Value: "14",
				Label: "ITU-R BT.2020 10 bit",
			},
			15: {
				Value: "15",
				Label: "ITU-R BT.2020 12 bit",
			},
			16: {
				Value: "16",
				Label: "ITU-R BT.2100 Perceptual Quantization",
			},
			17: {
				Value: "17",
				Label: "SMPTE ST 428-1",
			},
			18: {
				Value: "18",
				Label: "ARIB STD-B67 (HLG)",
			},
		},
	},
	Primaries: {
		Name:       "Primaries",
		Definition: "The colour primaries of the video. For clarity, the value and meanings for Primaries are adopted from Table 2 of ISO/IEC 23091-4 or ITU-T H.273.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "reserved",
			},
			1: {
				Value: "1",
				Label: "ITU-R BT.709",
			},
			2: {
				Value: "2",
				Label: "unspecified",
			},
			3: {
				Value: "3",
				Label: "reserved",
			},
			4: {
				Value: "4",
				Label: "ITU-R BT.470M",
			},
			5: {
				Value: "5",
				Label: "ITU-R BT.470BG - BT.601 625",
			},
			6: {
				Value: "6",
				Label: "ITU-R BT.601 525 - SMPTE 170M",
			},
			7: {
				Value: "7",
				Label: "SMPTE 240M",
			},
			8: {
				Value: "8",
				Label: "FILM",
			},
			9: {
				Value: "9",
				Label: "ITU-R BT.2020",
			},
			10: {
				Value: "10",
				Label: "SMPTE ST 428-1",
			},
			11: {
				Value: "11",
				Label: "SMPTE RP 432-2",
			},
			12: {
				Value: "12",
				Label: "SMPTE EG 432-2",
			},
			22: {
				Value: "22",
				Label: "EBU Tech. 3213-E - JEDEC P22 phosphors",
			},
		},
	},
	MaxCLL: {
		Name:       "MaxCLL",
		Definition: "Maximum brightness of a single pixel (Maximum Content Light Level) in candelas per square meter (cd/m²).",
		Type:       ebml.Uinteger,
	},
	MaxFALL: {
		Name:       "MaxFALL",
		Definition: "Maximum brightness of a single full frame (Maximum Frame-Average Light Level) in candelas per square meter (cd/m²).",
		Type:       ebml.Uinteger,
	},
	MasteringMetadata: {
		Name:       "MasteringMetadata",
		Definition: "SMPTE 2086 mastering data.",
		Type:       ebml.Master, Tag: mkvMasteringMetadata,
	},
}

var mkvMasteringMetadata = ebml.Tag{
	PrimaryRChromaticityX: {
		Name:       "PrimaryRChromaticityX",
		Definition: "Red X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryRChromaticityY: {
		Name:       "PrimaryRChromaticityY",
		Definition: "Red Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryGChromaticityX: {
		Name:       "PrimaryGChromaticityX",
		Definition: "Green X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryGChromaticityY: {
		Name:       "PrimaryGChromaticityY",
		Definition: "Green Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryBChromaticityX: {
		Name:       "PrimaryBChromaticityX",
		Definition: "Blue X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryBChromaticityY: {
		Name:       "PrimaryBChromaticityY",
		Definition: "Blue Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	WhitePointChromaticityX: {
		Name:       "WhitePointChromaticityX",
		Definition: "White X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	WhitePointChromaticityY: {
		Name:       "WhitePointChromaticityY",
		Definition: "White Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	LuminanceMax: {
		Name:       "LuminanceMax",
		Definition: "Maximum luminance. Represented in candelas per square meter (cd/m²).",
		Type:       ebml.Float,
	},
	LuminanceMin: {
		Name:       "LuminanceMin",
		Definition: "Minimum luminance. Represented in candelas per square meter (cd/m²).",
		Type:       ebml.Float,
	},
}

var mkvProjection = ebml.Tag{
	ProjectionType: {
		Name:       "ProjectionType",
		Definition: "Describes the projection used for this video track.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "rectangular",
			},
			1: {
				Value: "1",
				Label: "equirectangular",
			},
			2: {
				Value: "2",
				Label: "cubemap",
			},
			3: {
				Value: "3",
				Label: "mesh",
			},
		},
	},
	ProjectionPrivate: {
		Name:       "ProjectionPrivate",
		Definition: "Private data that only applies to a specific projection.SemanticsIf ProjectionType equals 0 (Rectangular),\n     then this element must not be present.If ProjectionType equals 1 (Equirectangular), then this element must be present and contain the same binary data that would be stored inside\n      an ISOBMFF Equirectangular Projection Box ('equi').If ProjectionType equals 2 (Cubemap), then this element must be present and contain the same binary data that would be stored \n      inside an ISOBMFF Cubemap Projection Box ('cbmp').If ProjectionType equals 3 (Mesh), then this element must be present and contain the same binary data that would be stored inside\n       an ISOBMFF Mesh Projection Box ('mshp').Note: ISOBMFF box size and fourcc fields are not included in the binary data, but the FullBox version and flag fields are. This is to avoid \n       redundant framing information while preserving versioning and semantics between the two container formats.",
		Type:       ebml.Binary,
	},
	ProjectionPoseYaw: {
		Name:       "ProjectionPoseYaw",
		Definition: "Specifies a yaw rotation to the projection.SemanticsValue represents a clockwise rotation, in degrees, around the up vector. This rotation must be applied before any ProjectionPosePitch or ProjectionPoseRoll rotations. The value of this field should be in the -180 to 180 degree range.",
		Type:       ebml.Float,
	},
	ProjectionPosePitch: {
		Name:       "ProjectionPosePitch",
		Definition: "Specifies a pitch rotation to the projection.SemanticsValue represents a counter-clockwise rotation, in degrees, around the right vector. This rotation must be applied after the ProjectionPoseYaw rotation and before the ProjectionPoseRoll rotation. The value of this field should be in the -90 to 90 degree range.",
		Type:       ebml.Float,
	},
	ProjectionPoseRoll: {
		Name:       "ProjectionPoseRoll",
		Definition: "Specifies a roll rotation to the projection.SemanticsValue represents a counter-clockwise rotation, in degrees, around the forward vector. This rotation must be applied after the ProjectionPoseYaw and ProjectionPosePitch rotations. The value of this field should be in the -180 to 180 degree range.",
		Type:       ebml.Float,
	},
}

var mkvAudio = ebml.Tag{
	SamplingFrequency: {
		Name:       "SamplingFrequency",
		Definition: "Sampling frequency in Hz.",
		Type:       ebml.Float,
	},
	OutputSamplingFrequency: {
		Name:       "OutputSamplingFrequency",
		Definition: "Real output sampling frequency in Hz (used for SBR techniques).",
		Type:       ebml.Float,
	},
	Channels: {
		Name:       "Channels",
		Definition: "Numbers of channels in the track.",
		Type:       ebml.Uinteger,
	},
	ChannelPositions: {
		Name:       "ChannelPositions",
		Definition: "Table of horizontal angles for each successive channel, see .",
		Type:       ebml.Binary,
	},
	BitDepth: {
		Name:       "BitDepth",
		Definition: "Bits per sample, mostly used for PCM.",
		Type:       ebml.Uinteger,
	},
}

var mkvTrackOperation = ebml.Tag{
	TrackCombinePlanes: {
		Name:       "TrackCombinePlanes",
		Definition: "Contains the list of all video plane tracks that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: mkvTrackCombinePlanes,
	},
	TrackJoinBlocks: {
		Name:       "TrackJoinBlocks",
		Definition: "Contains the list of all tracks whose Blocks need to be combined to create this virtual track",
		Type:       ebml.Master, Tag: mkvTrackJoinBlocks,
	},
}

var mkvTrackCombinePlanes = ebml.Tag{
	TrackPlane: {
		Name:       "TrackPlane",
		Definition: "Contains a video plane track that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: mkvTrackPlane,
	},
}

var mkvTrackPlane = ebml.Tag{
	TrackPlaneUID: {
		Name:       "TrackPlaneUID",
		Definition: "The trackUID number of the track representing the plane.",
		Type:       ebml.Uinteger,
	},
	TrackPlaneType: {
		Name:       "TrackPlaneType",
		Definition: "The kind of plane this track corresponds to.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "left eye",
			},
			1: {
				Value: "1",
				Label: "right eye",
			},
			2: {
				Value: "2",
				Label: "background",
			},
		},
	},
}

var mkvTrackJoinBlocks = ebml.Tag{
	TrackJoinUID: {
		Name:       "TrackJoinUID",
		Definition: "The trackUID number of a track whose blocks are used to create this virtual track.",
		Type:       ebml.Uinteger,
	},
}

var mkvContentEncodings = ebml.Tag{
	ContentEncoding: {
		Name:       "ContentEncoding",
		Definition: "Settings for one content encoding like compression or encryption.",
		Type:       ebml.Master, Tag: mkvContentEncoding,
	},
}

var mkvContentEncoding = ebml.Tag{
	ContentEncodingOrder: {
		Name:       "ContentEncodingOrder",
		Definition: "Tells when this modification was used during encoding/muxing starting with 0 and counting upwards. The decoder/demuxer has to start with the highest order number it finds and work its way down. This value has to be unique over all ContentEncodingOrder Elements in the TrackEntry that contains this ContentEncodingOrder element.",
		Type:       ebml.Uinteger,
	},
	ContentEncodingScope: {
		Name:       "ContentEncodingScope",
		Definition: "A bit field that describes which Elements have been modified in this way. Values (big endian) can be OR'ed.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			1: {
				Value: "1",
				Label: "All frame contents, excluding lacing data",
			},
			2: {
				Value: "2",
				Label: "The track's private data",
			},
			4: {
				Value: "4",
				Label: "The next ContentEncoding (next `ContentEncodingOrder`. Either the data inside `ContentCompression` and/or `ContentEncryption`)",
			},
		},
	},
	ContentEncodingType: {
		Name:       "ContentEncodingType",
		Definition: "A value describing what kind of transformation is applied.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "Compression",
			},
			1: {
				Value: "1",
				Label: "Encryption",
			},
		},
	},
	ContentCompression: {
		Name:       "ContentCompression",
		Definition: "Settings describing the compression used. This Element MUST be present if the value of ContentEncodingType is 0 and absent otherwise. Each block MUST be decompressable even if no previous block is available in order not to prevent seeking.",
		Type:       ebml.Master, Tag: mkvContentCompression,
	},
	ContentEncryption: {
		Name:       "ContentEncryption",
		Definition: "Settings describing the encryption used. This Element MUST be present if the value of `ContentEncodingType` is 1 (encryption) and MUST be ignored otherwise.",
		Type:       ebml.Master, Tag: mkvContentEncryption,
	},
}

var mkvContentCompression = ebml.Tag{
	ContentCompAlgo: {
		Name:       "ContentCompAlgo",
		Definition: "The compression algorithm used.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "zlib",
			},
			1: {
				Value: "1",
				Label: "bzlib",
			},
			2: {
				Value: "2",
				Label: "lzo1x",
			},
			3: {
				Value: "3",
				Label: "Header Stripping",
			},
		},
	},
	ContentCompSettings: {
		Name:       "ContentCompSettings",
		Definition: "Settings that might be needed by the decompressor. For Header Stripping (`ContentCompAlgo`=3), the bytes that were removed from the beginning of each frames of the track.",
		Type:       ebml.Binary,
	},
}

var mkvContentEncryption = ebml.Tag{
	ContentEncAlgo: {
		Name:       "ContentEncAlgo",
		Definition: "The encryption algorithm used. The value '0' means that the contents have not been encrypted but only signed.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "Not encrypted",
			},
			1: {
				Value: "1",
				Label: "DES - FIPS 46-3",
			},
			2: {
				Value: "2",
				Label: "Triple DES - RFC 1851",
			},
			3: {
				Value: "3",
				Label: "Twofish",
			},
			4: {
				Value: "4",
				Label: "Blowfish",
			},
			5: {
				Value: "5",
				Label: "AES - FIPS 187",
			},
		},
	},
	ContentEncKeyID: {
		Name:       "ContentEncKeyID",
		Definition: "For public key algorithms this is the ID of the public key the the data was encrypted with.",
		Type:       ebml.Binary,
	},
	ContentEncAESSettings: {
		Name:       "ContentEncAESSettings",
		Definition: "Settings describing the encryption algorithm used. If `ContentEncAlgo` != 5 this MUST be ignored.",
		Type:       ebml.Master, Tag: mkvContentEncAESSettings,
	},
	ContentSignature: {
		Name:       "ContentSignature",
		Definition: "A cryptographic signature of the contents.",
		Type:       ebml.Binary,
	},
	ContentSigKeyID: {
		Name:       "ContentSigKeyID",
		Definition: "This is the ID of the private key the data was signed with.",
		Type:       ebml.Binary,
	},
	ContentSigAlgo: {
		Name:       "ContentSigAlgo",
		Definition: "The algorithm used for the signature.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "Not signed",
			},
			1: {
				Value: "1",
				Label: "RSA",
			},
		},
	},
	ContentSigHashAlgo: {
		Name:       "ContentSigHashAlgo",
		Definition: "The hash algorithm used for the signature.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "Not signed",
			},
			1: {
				Value: "1",
				Label: "SHA1-160",
			},
			2: {
				Value: "2",
				Label: "MD5",
			},
		},
	},
}

var mkvContentEncAESSettings = ebml.Tag{
	AESSettingsCipherMode: {
		Name:       "AESSettingsCipherMode",
		Definition: "The AES cipher mode used in the encryption.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			1: {
				Value: "1",
				Label: "AES-CTR / Counter, NIST SP 800-38A",
			},
			2: {
				Value: "2",
				Label: "AES-CBC / Cipher Block Chaining, NIST SP 800-38A",
			},
		},
	},
}

var mkvCues = ebml.Tag{
	CuePoint: {
		Name:       "CuePoint",
		Definition: "Contains all information relative to a seek point in the Segment.",
		Type:       ebml.Master, Tag: mkvCuePoint,
	},
}

var mkvCuePoint = ebml.Tag{
	CueTime: {
		Name:       "CueTime",
		Definition: "Absolute timestamp according to the Segment time base.",
		Type:       ebml.Uinteger,
	},
	CueTrackPositions: {
		Name:       "CueTrackPositions",
		Definition: "Contain positions for different tracks corresponding to the timestamp.",
		Type:       ebml.Master, Tag: mkvCueTrackPositions,
	},
}

var mkvCueTrackPositions = ebml.Tag{
	CueTrack: {
		Name:       "CueTrack",
		Definition: "The track for which a position is given.",
		Type:       ebml.Uinteger,
	},
	CueClusterPosition: {
		Name:       "CueClusterPosition",
		Definition: "The Segment Position of the Cluster containing the associated Block.",
		Type:       ebml.Uinteger,
	},
	CueRelativePosition: {
		Name:       "CueRelativePosition",
		Definition: "The relative position inside the Cluster of the referenced SimpleBlock or BlockGroup with 0 being the first possible position for an Element inside that Cluster.",
		Type:       ebml.Uinteger,
	},
	CueDuration: {
		Name:       "CueDuration",
		Definition: "The duration of the block according to the Segment time base. If missing the track's DefaultDuration does not apply and no duration information is available in terms of the cues.",
		Type:       ebml.Uinteger,
	},
	CueBlockNumber: {
		Name:       "CueBlockNumber",
		Definition: "Number of the Block in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	CueCodecState: {
		Name:       "CueCodecState",
		Definition: "The Segment Position of the Codec State corresponding to this Cue Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
	CueReference: {
		Name:       "CueReference",
		Definition: "The Clusters containing the referenced Blocks.",
		Type:       ebml.Master, Tag: mkvCueReference,
	},
}

var mkvCueReference = ebml.Tag{
	CueRefTime: {
		Name:       "CueRefTime",
		Definition: "Timestamp of the referenced Block.",
		Type:       ebml.Uinteger,
	},
	CueRefCluster: {
		Name:       "CueRefCluster",
		Definition: "The Segment Position of the Cluster containing the referenced Block.",
		Type:       ebml.Uinteger,
	},
	CueRefNumber: {
		Name:       "CueRefNumber",
		Definition: "Number of the referenced Block of Track X in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	CueRefCodecState: {
		Name:       "CueRefCodecState",
		Definition: "The Segment Position of the Codec State corresponding to this referenced Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
}

var mkvAttachments = ebml.Tag{
	AttachedFile: {
		Name:       "AttachedFile",
		Definition: "An attached file.",
		Type:       ebml.Master, Tag: mkvAttachedFile,
	},
}

var mkvAttachedFile = ebml.Tag{
	FileDescription: {
		Name:       "FileDescription",
		Definition: "A human-friendly name for the attached file.",
		Type:       ebml.UTF8,
	},
	FileName: {
		Name:       "FileName",
		Definition: "Filename of the attached file.",
		Type:       ebml.UTF8,
	},
	FileMimeType: {
		Name:       "FileMimeType",
		Definition: "MIME type of the file.",
		Type:       ebml.String,
	},
	FileData: {
		Name:       "FileData",
		Definition: "The data of the file.",
		Type:       ebml.Binary,
	},
	FileUID: {
		Name:       "FileUID",
		Definition: "Unique ID representing the file, as random as possible.",
		Type:       ebml.Uinteger,
	},
	FileReferral: {
		Name:       "FileReferral",
		Definition: "A binary value that a track/codec can refer to when the attachment is needed.",
		Type:       ebml.Binary,
	},
	FileUsedStartTime: {
		Name:       "FileUsedStartTime",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	FileUsedEndTime: {
		Name:       "FileUsedEndTime",
		Definition: "",
		Type:       ebml.Uinteger,
	},
}

var mkvChapters = ebml.Tag{
	EditionEntry: {
		Name:       "EditionEntry",
		Definition: "Contains all information about a Segment edition.",
		Type:       ebml.Master, Tag: mkvEditionEntry,
	},
}

var mkvEditionEntry = ebml.Tag{
	EditionUID: {
		Name:       "EditionUID",
		Definition: "A unique ID to identify the edition. It's useful for tagging an edition.",
		Type:       ebml.Uinteger,
	},
	EditionFlagHidden: {
		Name:       "EditionFlagHidden",
		Definition: "If an edition is hidden (1), it SHOULD NOT be available to the user interface (but still to Control Tracks; see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	EditionFlagDefault: {
		Name:       "EditionFlagDefault",
		Definition: "If a flag is set (1) the edition SHOULD be used as the default one. (1 bit)",
		Type:       ebml.Uinteger,
	},
	EditionFlagOrdered: {
		Name:       "EditionFlagOrdered",
		Definition: "Specify if the chapters can be defined multiple times and the order to play them is enforced. (1 bit)",
		Type:       ebml.Uinteger,
	},
	ChapterAtom: {
		Name:       "ChapterAtom",
		Definition: "Contains the atom information to use as the chapter atom (apply to all tracks).",
		Type:       ebml.Master, Tag: mkvChapterAtom,
	},
}

var mkvChapterAtom = ebml.Tag{
	ChapterUID: {
		Name:       "ChapterUID",
		Definition: "A unique ID to identify the Chapter.",
		Type:       ebml.Uinteger,
	},
	ChapterStringUID: {
		Name:       "ChapterStringUID",
		Definition: "A unique string ID to identify the Chapter. Use for .",
		Type:       ebml.UTF8,
	},
	ChapterTimeStart: {
		Name:       "ChapterTimeStart",
		Definition: "Timestamp of the start of Chapter (not scaled).",
		Type:       ebml.Uinteger,
	},
	ChapterTimeEnd: {
		Name:       "ChapterTimeEnd",
		Definition: "Timestamp of the end of Chapter (timestamp excluded, not scaled).",
		Type:       ebml.Uinteger,
	},
	ChapterFlagHidden: {
		Name:       "ChapterFlagHidden",
		Definition: "If a chapter is hidden (1), it SHOULD NOT be available to the user interface (but still to Control Tracks; see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	ChapterFlagEnabled: {
		Name:       "ChapterFlagEnabled",
		Definition: "Specify whether the chapter is enabled. It can be enabled/disabled by a Control Track. When disabled, the movie SHOULD skip all the content between the TimeStart and TimeEnd of this chapter (see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	ChapterSegmentUID: {
		Name:       "ChapterSegmentUID",
		Definition: "The SegmentUID of another Segment to play during this chapter.",
		Type:       ebml.Binary,
	},
	ChapterSegmentEditionUID: {
		Name:       "ChapterSegmentEditionUID",
		Definition: "The EditionUID to play from the Segment linked in ChapterSegmentUID. If ChapterSegmentEditionUID is undeclared then no Edition of the linked Segment is used.",
		Type:       ebml.Uinteger,
	},
	ChapterPhysicalEquiv: {
		Name:       "ChapterPhysicalEquiv",
		Definition: "Specify the physical equivalent of this ChapterAtom like \"DVD\" (60) or \"SIDE\" (50), see .",
		Type:       ebml.Uinteger,
	},
	ChapterTrack: {
		Name:       "ChapterTrack",
		Definition: "List of tracks on which the chapter applies. If this Element is not present, all tracks apply",
		Type:       ebml.Master, Tag: mkvChapterTrack,
	},
	ChapterDisplay: {
		Name:       "ChapterDisplay",
		Definition: "Contains all possible strings to use for the chapter display.",
		Type:       ebml.Master, Tag: mkvChapterDisplay,
	},
	ChapProcess: {
		Name:       "ChapProcess",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: mkvChapProcess,
	},
}

var mkvChapterTrack = ebml.Tag{
	ChapterTrackUID: {
		Name:       "ChapterTrackUID",
		Definition: "UID of the Track to apply this chapter too. In the absence of a control track, choosing this chapter will select the listed Tracks and deselect unlisted tracks. Absence of this Element indicates that the Chapter SHOULD be applied to any currently used Tracks.",
		Type:       ebml.Uinteger,
	},
}

var mkvChapterDisplay = ebml.Tag{
	ChapString: {
		Name:       "ChapString",
		Definition: "Contains the string to use as the chapter atom.",
		Type:       ebml.UTF8,
	},
	ChapLanguage: {
		Name:       "ChapLanguage",
		Definition: "The languages corresponding to the string, in the . This Element MUST be ignored if the ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
	ChapLanguageIETF: {
		Name:       "ChapLanguageIETF",
		Definition: "Specifies the language used in the ChapString according to  and using the . If this Element is used, then any ChapLanguage Elements used in the same ChapterDisplay MUST be ignored.",
		Type:       ebml.String,
	},
	ChapCountry: {
		Name:       "ChapCountry",
		Definition: "The countries corresponding to the string, same 2 octets as in . This Element MUST be ignored if the ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
}

var mkvChapProcess = ebml.Tag{
	ChapProcessCodecID: {
		Name:       "ChapProcessCodecID",
		Definition: "Contains the type of the codec used for the processing. A value of 0 means native Matroska processing (to be defined), a value of 1 means the  command set is used. More codec IDs can be added later.",
		Type:       ebml.Uinteger,
	},
	ChapProcessPrivate: {
		Name:       "ChapProcessPrivate",
		Definition: "Some optional data attached to the ChapProcessCodecID information. , it is the \"DVD level\" equivalent.",
		Type:       ebml.Binary,
	},
	ChapProcessCommand: {
		Name:       "ChapProcessCommand",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: mkvChapProcessCommand,
	},
}

var mkvChapProcessCommand = ebml.Tag{
	ChapProcessTime: {
		Name:       "ChapProcessTime",
		Definition: "Defines when the process command SHOULD be handled",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			0: {
				Value: "0",
				Label: "during the whole chapter",
			},
			1: {
				Value: "1",
				Label: "before starting playback",
			},
			2: {
				Value: "2",
				Label: "after playback of the chapter",
			},
		},
	},
	ChapProcessData: {
		Name:       "ChapProcessData",
		Definition: "Contains the command information. The data SHOULD be interpreted depending on the ChapProcessCodecID value. , the data correspond to the binary DVD cell pre/post commands.",
		Type:       ebml.Binary,
	},
}

var mkvTags = ebml.Tag{
	Tag: {
		Name:       "Tag",
		Definition: "A single metadata descriptor.",
		Type:       ebml.Master, Tag: mkvTag,
	},
}

var mkvTag = ebml.Tag{
	Targets: {
		Name:       "Targets",
		Definition: "Specifies which other elements the metadata represented by the Tag applies to. If empty or not present, then the Tag describes everything in the Segment.",
		Type:       ebml.Master, Tag: mkvTargets,
	},
	SimpleTag: {
		Name:       "SimpleTag",
		Definition: "Contains general information about the target.",
		Type:       ebml.Master, Tag: mkvSimpleTag,
	},
}

var mkvTargets = ebml.Tag{
	TargetTypeValue: {
		Name:       "TargetTypeValue",
		Definition: "A number to indicate the logical level of the target.",
		Type:       ebml.Uinteger,
		UintegerEnums: map[uint64]ebml.Enum{
			70: {
				Value:      "70",
				Label:      "COLLECTION",
				Definition: "The highest hierarchical level that tags can describe.",
			},
			60: {
				Value:      "60",
				Label:      "EDITION / ISSUE / VOLUME / OPUS / SEASON / SEQUEL",
				Definition: "A list of lower levels grouped together.",
			},
			50: {
				Value:      "50",
				Label:      "ALBUM / OPERA / CONCERT / MOVIE / EPISODE / CONCERT",
				Definition: "The most common grouping level of music and video (equals to an episode for TV series).",
			},
			40: {
				Value:      "40",
				Label:      "PART / SESSION",
				Definition: "When an album or episode has different logical parts.",
			},
			30: {
				Value:      "30",
				Label:      "TRACK / SONG / CHAPTER",
				Definition: "The common parts of an album or movie.",
			},
			20: {
				Value:      "20",
				Label:      "SUBTRACK / PART / MOVEMENT / SCENE",
				Definition: "Corresponds to parts of a track for audio (like a movement).",
			},
			10: {
				Value:      "10",
				Label:      "SHOT",
				Definition: "The lowest hierarchy found in music or movies.",
			},
		},
	},
	TargetType: {
		Name:       "TargetType",
		Definition: "An informational string that can be used to display the logical level of the target like \"ALBUM\", \"TRACK\", \"MOVIE\", \"CHAPTER\", etc (see ).",
		Type:       ebml.String,
		StringEnums: map[string]ebml.Enum{
			"COLLECTION": {
				Value: "COLLECTION",
				Label: "COLLECTION",
			},
			"EDITION": {
				Value: "EDITION",
				Label: "EDITION",
			},
			"ISSUE": {
				Value: "ISSUE",
				Label: "ISSUE",
			},
			"VOLUME": {
				Value: "VOLUME",
				Label: "VOLUME",
			},
			"OPUS": {
				Value: "OPUS",
				Label: "OPUS",
			},
			"SEASON": {
				Value: "SEASON",
				Label: "SEASON",
			},
			"SEQUEL": {
				Value: "SEQUEL",
				Label: "SEQUEL",
			},
			"ALBUM": {
				Value: "ALBUM",
				Label: "ALBUM",
			},
			"OPERA": {
				Value: "OPERA",
				Label: "OPERA",
			},
			"CONCERT": {
				Value: "CONCERT",
				Label: "CONCERT",
			},
			"MOVIE": {
				Value: "MOVIE",
				Label: "MOVIE",
			},
			"EPISODE": {
				Value: "EPISODE",
				Label: "EPISODE",
			},
			"PART": {
				Value: "PART",
				Label: "PART",
			},
			"SESSION": {
				Value: "SESSION",
				Label: "SESSION",
			},
			"TRACK": {
				Value: "TRACK",
				Label: "TRACK",
			},
			"SONG": {
				Value: "SONG",
				Label: "SONG",
			},
			"CHAPTER": {
				Value: "CHAPTER",
				Label: "CHAPTER",
			},
			"SUBTRACK": {
				Value: "SUBTRACK",
				Label: "SUBTRACK",
			},
			"MOVEMENT": {
				Value: "MOVEMENT",
				Label: "MOVEMENT",
			},
			"SCENE": {
				Value: "SCENE",
				Label: "SCENE",
			},
			"SHOT": {
				Value: "SHOT",
				Label: "SHOT",
			},
		},
	},
	TagTrackUID: {
		Name:       "TagTrackUID",
		Definition: "A unique ID to identify the Track(s) the tags belong to. If the value is 0 at this level, the tags apply to all tracks in the Segment.",
		Type:       ebml.Uinteger,
	},
	TagEditionUID: {
		Name:       "TagEditionUID",
		Definition: "A unique ID to identify the EditionEntry(s) the tags belong to. If the value is 0 at this level, the tags apply to all editions in the Segment.",
		Type:       ebml.Uinteger,
	},
	TagChapterUID: {
		Name:       "TagChapterUID",
		Definition: "A unique ID to identify the Chapter(s) the tags belong to. If the value is 0 at this level, the tags apply to all chapters in the Segment.",
		Type:       ebml.Uinteger,
	},
	TagAttachmentUID: {
		Name:       "TagAttachmentUID",
		Definition: "A unique ID to identify the Attachment(s) the tags belong to. If the value is 0 at this level, the tags apply to all the attachments in the Segment.",
		Type:       ebml.Uinteger,
	},
}

var mkvSimpleTag = ebml.Tag{
	TagName: {
		Name:       "TagName",
		Definition: "The name of the Tag that is going to be stored.",
		Type:       ebml.UTF8,
	},
	TagLanguage: {
		Name:       "TagLanguage",
		Definition: "Specifies the language of the tag specified, in the . This Element MUST be ignored if the TagLanguageIETF Element is used within the same SimpleTag Element.",
		Type:       ebml.String,
	},
	TagLanguageIETF: {
		Name:       "TagLanguageIETF",
		Definition: "Specifies the language used in the TagString according to  and using the . If this Element is used, then any TagLanguage Elements used in the same SimpleTag MUST be ignored.",
		Type:       ebml.String,
	},
	TagDefault: {
		Name:       "TagDefault",
		Definition: "A boolean value to indicate if this is the default/original language to use for the given tag.",
		Type:       ebml.Uinteger,
	},
	TagString: {
		Name:       "TagString",
		Definition: "The value of the Tag.",
		Type:       ebml.UTF8,
	},
	TagBinary: {
		Name:       "TagBinary",
		Definition: "The values of the Tag if it is binary. Note that this cannot be used in the same SimpleTag as TagString.",
		Type:       ebml.Binary,
	},
}
