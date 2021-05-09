// Code below generated from ebml_matroska.xml
package mkv

import "fq/pkg/format/mkv/ebml"

const (
	EBMLMaxIDLengthID             = 0x42F2
	EBMLMaxSizeLengthID           = 0x42F3
	SegmentID                     = 0x18538067
	SeekHeadID                    = 0x114D9B74
	SeekID                        = 0x4DBB
	SeekIDID                      = 0x53AB
	SeekPositionID                = 0x53AC
	InfoID                        = 0x1549A966
	SegmentUIDID                  = 0x73A4
	SegmentFilenameID             = 0x7384
	PrevUIDID                     = 0x3CB923
	PrevFilenameID                = 0x3C83AB
	NextUIDID                     = 0x3EB923
	NextFilenameID                = 0x3E83BB
	SegmentFamilyID               = 0x4444
	ChapterTranslateID            = 0x6924
	ChapterTranslateEditionUIDID  = 0x69FC
	ChapterTranslateCodecID       = 0x69BF
	ChapterTranslateIDID          = 0x69A5
	TimestampScaleID              = 0x2AD7B1
	DurationID                    = 0x4489
	DateUTCID                     = 0x4461
	TitleID                       = 0x7BA9
	MuxingAppID                   = 0x4D80
	WritingAppID                  = 0x5741
	ClusterID                     = 0x1F43B675
	TimestampID                   = 0xE7
	SilentTracksID                = 0x5854
	SilentTrackNumberID           = 0x58D7
	PositionID                    = 0xA7
	PrevSizeID                    = 0xAB
	SimpleBlockID                 = 0xA3
	BlockGroupID                  = 0xA0
	BlockID                       = 0xA1
	BlockVirtualID                = 0xA2
	BlockAdditionsID              = 0x75A1
	BlockMoreID                   = 0xA6
	BlockAddIDID                  = 0xEE
	BlockAdditionalID             = 0xA5
	BlockDurationID               = 0x9B
	ReferencePriorityID           = 0xFA
	ReferenceBlockID              = 0xFB
	ReferenceVirtualID            = 0xFD
	CodecStateID                  = 0xA4
	DiscardPaddingID              = 0x75A2
	SlicesID                      = 0x8E
	TimeSliceID                   = 0xE8
	LaceNumberID                  = 0xCC
	FrameNumberID                 = 0xCD
	BlockAdditionIDID             = 0xCB
	DelayID                       = 0xCE
	SliceDurationID               = 0xCF
	ReferenceFrameID              = 0xC8
	ReferenceOffsetID             = 0xC9
	ReferenceTimestampID          = 0xCA
	EncryptedBlockID              = 0xAF
	TracksID                      = 0x1654AE6B
	TrackEntryID                  = 0xAE
	TrackNumberID                 = 0xD7
	TrackUIDID                    = 0x73C5
	TrackTypeID                   = 0x83
	FlagEnabledID                 = 0xB9
	FlagDefaultID                 = 0x88
	FlagForcedID                  = 0x55AA
	FlagLacingID                  = 0x9C
	MinCacheID                    = 0x6DE7
	MaxCacheID                    = 0x6DF8
	DefaultDurationID             = 0x23E383
	DefaultDecodedFieldDurationID = 0x234E7A
	TrackTimestampScaleID         = 0x23314F
	TrackOffsetID                 = 0x537F
	MaxBlockAdditionIDID          = 0x55EE
	BlockAdditionMappingID        = 0x41E4
	BlockAddIDValueID             = 0x41F0
	BlockAddIDNameID              = 0x41A4
	BlockAddIDTypeID              = 0x41E7
	BlockAddIDExtraDataID         = 0x41ED
	NameID                        = 0x536E
	LanguageID                    = 0x22B59C
	LanguageIETFID                = 0x22B59D
	CodecIDID                     = 0x86
	CodecPrivateID                = 0x63A2
	CodecNameID                   = 0x258688
	AttachmentLinkID              = 0x7446
	CodecSettingsID               = 0x3A9697
	CodecInfoURLID                = 0x3B4040
	CodecDownloadURLID            = 0x26B240
	CodecDecodeAllID              = 0xAA
	TrackOverlayID                = 0x6FAB
	CodecDelayID                  = 0x56AA
	SeekPreRollID                 = 0x56BB
	TrackTranslateID              = 0x6624
	TrackTranslateEditionUIDID    = 0x66FC
	TrackTranslateCodecID         = 0x66BF
	TrackTranslateTrackIDID       = 0x66A5
	VideoID                       = 0xE0
	FlagInterlacedID              = 0x9A
	FieldOrderID                  = 0x9D
	StereoModeID                  = 0x53B8
	AlphaModeID                   = 0x53C0
	OldStereoModeID               = 0x53B9
	PixelWidthID                  = 0xB0
	PixelHeightID                 = 0xBA
	PixelCropBottomID             = 0x54AA
	PixelCropTopID                = 0x54BB
	PixelCropLeftID               = 0x54CC
	PixelCropRightID              = 0x54DD
	DisplayWidthID                = 0x54B0
	DisplayHeightID               = 0x54BA
	DisplayUnitID                 = 0x54B2
	AspectRatioTypeID             = 0x54B3
	ColourSpaceID                 = 0x2EB524
	GammaValueID                  = 0x2FB523
	FrameRateID                   = 0x2383E3
	ColourID                      = 0x55B0
	MatrixCoefficientsID          = 0x55B1
	BitsPerChannelID              = 0x55B2
	ChromaSubsamplingHorzID       = 0x55B3
	ChromaSubsamplingVertID       = 0x55B4
	CbSubsamplingHorzID           = 0x55B5
	CbSubsamplingVertID           = 0x55B6
	ChromaSitingHorzID            = 0x55B7
	ChromaSitingVertID            = 0x55B8
	RangeID                       = 0x55B9
	TransferCharacteristicsID     = 0x55BA
	PrimariesID                   = 0x55BB
	MaxCLLID                      = 0x55BC
	MaxFALLID                     = 0x55BD
	MasteringMetadataID           = 0x55D0
	PrimaryRChromaticityXID       = 0x55D1
	PrimaryRChromaticityYID       = 0x55D2
	PrimaryGChromaticityXID       = 0x55D3
	PrimaryGChromaticityYID       = 0x55D4
	PrimaryBChromaticityXID       = 0x55D5
	PrimaryBChromaticityYID       = 0x55D6
	WhitePointChromaticityXID     = 0x55D7
	WhitePointChromaticityYID     = 0x55D8
	LuminanceMaxID                = 0x55D9
	LuminanceMinID                = 0x55DA
	ProjectionID                  = 0x7670
	ProjectionTypeID              = 0x7671
	ProjectionPrivateID           = 0x7672
	ProjectionPoseYawID           = 0x7673
	ProjectionPosePitchID         = 0x7674
	ProjectionPoseRollID          = 0x7675
	AudioID                       = 0xE1
	SamplingFrequencyID           = 0xB5
	OutputSamplingFrequencyID     = 0x78B5
	ChannelsID                    = 0x9F
	ChannelPositionsID            = 0x7D7B
	BitDepthID                    = 0x6264
	TrackOperationID              = 0xE2
	TrackCombinePlanesID          = 0xE3
	TrackPlaneID                  = 0xE4
	TrackPlaneUIDID               = 0xE5
	TrackPlaneTypeID              = 0xE6
	TrackJoinBlocksID             = 0xE9
	TrackJoinUIDID                = 0xED
	TrickTrackUIDID               = 0xC0
	TrickTrackSegmentUIDID        = 0xC1
	TrickTrackFlagID              = 0xC6
	TrickMasterTrackUIDID         = 0xC7
	TrickMasterTrackSegmentUIDID  = 0xC4
	ContentEncodingsID            = 0x6D80
	ContentEncodingID             = 0x6240
	ContentEncodingOrderID        = 0x5031
	ContentEncodingScopeID        = 0x5032
	ContentEncodingTypeID         = 0x5033
	ContentCompressionID          = 0x5034
	ContentCompAlgoID             = 0x4254
	ContentCompSettingsID         = 0x4255
	ContentEncryptionID           = 0x5035
	ContentEncAlgoID              = 0x47E1
	ContentEncKeyIDID             = 0x47E2
	ContentEncAESSettingsID       = 0x47E7
	AESSettingsCipherModeID       = 0x47E8
	ContentSignatureID            = 0x47E3
	ContentSigKeyIDID             = 0x47E4
	ContentSigAlgoID              = 0x47E5
	ContentSigHashAlgoID          = 0x47E6
	CuesID                        = 0x1C53BB6B
	CuePointID                    = 0xBB
	CueTimeID                     = 0xB3
	CueTrackPositionsID           = 0xB7
	CueTrackID                    = 0xF7
	CueClusterPositionID          = 0xF1
	CueRelativePositionID         = 0xF0
	CueDurationID                 = 0xB2
	CueBlockNumberID              = 0x5378
	CueCodecStateID               = 0xEA
	CueReferenceID                = 0xDB
	CueRefTimeID                  = 0x96
	CueRefClusterID               = 0x97
	CueRefNumberID                = 0x535F
	CueRefCodecStateID            = 0xEB
	AttachmentsID                 = 0x1941A469
	AttachedFileID                = 0x61A7
	FileDescriptionID             = 0x467E
	FileNameID                    = 0x466E
	FileMimeTypeID                = 0x4660
	FileDataID                    = 0x465C
	FileUIDID                     = 0x46AE
	FileReferralID                = 0x4675
	FileUsedStartTimeID           = 0x4661
	FileUsedEndTimeID             = 0x4662
	ChaptersID                    = 0x1043A770
	EditionEntryID                = 0x45B9
	EditionUIDID                  = 0x45BC
	EditionFlagHiddenID           = 0x45BD
	EditionFlagDefaultID          = 0x45DB
	EditionFlagOrderedID          = 0x45DD
	ChapterAtomID                 = 0xB6
	ChapterUIDID                  = 0x73C4
	ChapterStringUIDID            = 0x5654
	ChapterTimeStartID            = 0x91
	ChapterTimeEndID              = 0x92
	ChapterFlagHiddenID           = 0x98
	ChapterFlagEnabledID          = 0x4598
	ChapterSegmentUIDID           = 0x6E67
	ChapterSegmentEditionUIDID    = 0x6EBC
	ChapterPhysicalEquivID        = 0x63C3
	ChapterTrackID                = 0x8F
	ChapterTrackUIDID             = 0x89
	ChapterDisplayID              = 0x80
	ChapStringID                  = 0x85
	ChapLanguageID                = 0x437C
	ChapLanguageIETFID            = 0x437D
	ChapCountryID                 = 0x437E
	ChapProcessID                 = 0x6944
	ChapProcessCodecIDID          = 0x6955
	ChapProcessPrivateID          = 0x450D
	ChapProcessCommandID          = 0x6911
	ChapProcessTimeID             = 0x6922
	ChapProcessDataID             = 0x6933
	TagsID                        = 0x1254C367
	TagID                         = 0x7373
	TargetsID                     = 0x63C0
	TargetTypeValueID             = 0x68CA
	TargetTypeID                  = 0x63CA
	TagTrackUIDID                 = 0x63C5
	TagEditionUIDID               = 0x63C9
	TagChapterUIDID               = 0x63C4
	TagAttachmentUIDID            = 0x63C6
	SimpleTagID                   = 0x67C8
	TagNameID                     = 0x45A3
	TagLanguageID                 = 0x447A
	TagLanguageIETFID             = 0x447B
	TagDefaultID                  = 0x4484
	TagStringID                   = 0x4487
	TagBinaryID                   = 0x4485
)

var mkvSegment = ebml.Tag{
	SeekHeadID: {
		Name:       "SeekHead",
		Definition: "Contains the Segment Position of other Top-Level Elements.",
		Type:       ebml.Master, Tag: mkvSeekHead,
	},
	InfoID: {
		Name:       "Info",
		Definition: "Contains general information about the Segment.",
		Type:       ebml.Master, Tag: mkvInfo,
	},
	ClusterID: {
		Name:       "Cluster",
		Definition: "The Top-Level Element containing the (monolithic) Block structure.",
		Type:       ebml.Master, Tag: mkvCluster,
	},
	TracksID: {
		Name:       "Tracks",
		Definition: "A Top-Level Element of information with many tracks described.",
		Type:       ebml.Master, Tag: mkvTracks,
	},
	CuesID: {
		Name:       "Cues",
		Definition: "A Top-Level Element to speed seeking access. All entries are local to the Segment.",
		Type:       ebml.Master, Tag: mkvCues,
	},
	AttachmentsID: {
		Name:       "Attachments",
		Definition: "Contain attached files.",
		Type:       ebml.Master, Tag: mkvAttachments,
	},
	ChaptersID: {
		Name:       "Chapters",
		Definition: "A system to define basic menus and partition data. For more detailed information, look at the .",
		Type:       ebml.Master, Tag: mkvChapters,
	},
	TagsID: {
		Name:       "Tags",
		Definition: "Element containing metadata describing Tracks, Editions, Chapters, Attachments, or the Segment as a whole. A list of valid tags can be found",
		Type:       ebml.Master, Tag: mkvTags,
	},
}

var mkvSeekHead = ebml.Tag{
	SeekID: {
		Name:       "Seek",
		Definition: "Contains a single seek entry to an EBML Element.",
		Type:       ebml.Master, Tag: mkvSeek,
	},
}

var mkvSeek = ebml.Tag{
	SeekIDID: {
		Name:       "SeekID",
		Definition: "The binary ID corresponding to the Element name.",
		Type:       ebml.Binary,
	},
	SeekPositionID: {
		Name:       "SeekPosition",
		Definition: "The Segment Position of the Element.",
		Type:       ebml.Uinteger,
	},
}

var mkvInfo = ebml.Tag{
	SegmentUIDID: {
		Name:       "SegmentUID",
		Definition: "A randomly generated unique ID to identify the Segment amongst many others (128 bits).",
		Type:       ebml.Binary,
	},
	SegmentFilenameID: {
		Name:       "SegmentFilename",
		Definition: "A filename corresponding to this Segment.",
		Type:       ebml.UTF8,
	},
	PrevUIDID: {
		Name:       "PrevUID",
		Definition: "A unique ID to identify the previous Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	PrevFilenameID: {
		Name:       "PrevFilename",
		Definition: "A filename corresponding to the file of the previous Linked Segment.",
		Type:       ebml.UTF8,
	},
	NextUIDID: {
		Name:       "NextUID",
		Definition: "A unique ID to identify the next Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	NextFilenameID: {
		Name:       "NextFilename",
		Definition: "A filename corresponding to the file of the next Linked Segment.",
		Type:       ebml.UTF8,
	},
	SegmentFamilyID: {
		Name:       "SegmentFamily",
		Definition: "A randomly generated unique ID that all Segments of a Linked Segment MUST share (128 bits).",
		Type:       ebml.Binary,
	},
	ChapterTranslateID: {
		Name:       "ChapterTranslate",
		Definition: "A tuple of corresponding ID used by chapter codecs to represent this Segment.",
		Type:       ebml.Master, Tag: mkvChapterTranslate,
	},
	TimestampScaleID: {
		Name:       "TimestampScale",
		Definition: "Timestamp scale in nanoseconds (1.000.000 means all timestamps in the Segment are expressed in milliseconds).",
		Type:       ebml.Uinteger,
	},
	DurationID: {
		Name:       "Duration",
		Definition: "Duration of the Segment in nanoseconds based on TimestampScale.",
		Type:       ebml.Float,
	},
	DateUTCID: {
		Name:       "DateUTC",
		Definition: "The date and time that the Segment was created by the muxing application or library.",
		Type:       ebml.Date,
	},
	TitleID: {
		Name:       "Title",
		Definition: "General name of the Segment.",
		Type:       ebml.UTF8,
	},
	MuxingAppID: {
		Name:       "MuxingApp",
		Definition: "Muxing application or library (example: \"libmatroska-0.4.3\").",
		Type:       ebml.UTF8,
	},
	WritingAppID: {
		Name:       "WritingApp",
		Definition: "Writing application (example: \"mkvmerge-0.3.3\").",
		Type:       ebml.UTF8,
	},
}

var mkvChapterTranslate = ebml.Tag{
	ChapterTranslateEditionUIDID: {
		Name:       "ChapterTranslateEditionUID",
		Definition: "Specify an edition UID on which this correspondence applies. When not specified, it means for all editions found in the Segment.",
		Type:       ebml.Uinteger,
	},
	ChapterTranslateCodecID: {
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
	ChapterTranslateIDID: {
		Name:       "ChapterTranslateID",
		Definition: "The binary value used to represent this Segment in the chapter codec data. The format depends on the  used.",
		Type:       ebml.Binary,
	},
}

var mkvCluster = ebml.Tag{
	TimestampID: {
		Name:       "Timestamp",
		Definition: "Absolute timestamp of the cluster (based on TimestampScale).",
		Type:       ebml.Uinteger,
	},
	SilentTracksID: {
		Name:       "SilentTracks",
		Definition: "The list of tracks that are not used in that part of the stream. It is useful when using overlay tracks on seeking or to decide what track to use.",
		Type:       ebml.Master, Tag: mkvSilentTracks,
	},
	PositionID: {
		Name:       "Position",
		Definition: "The Segment Position of the Cluster in the Segment (0 in live streams). It might help to resynchronise offset on damaged streams.",
		Type:       ebml.Uinteger,
	},
	PrevSizeID: {
		Name:       "PrevSize",
		Definition: "Size of the previous Cluster, in octets. Can be useful for backward playing.",
		Type:       ebml.Uinteger,
	},
	SimpleBlockID: {
		Name:       "SimpleBlock",
		Definition: "Similar to  but without all the extra information, mostly used to reduced overhead when no extra feature is needed. (see )",
		Type:       ebml.Binary,
	},
	BlockGroupID: {
		Name:       "BlockGroup",
		Definition: "Basic container of information containing a single Block and information specific to that Block.",
		Type:       ebml.Master, Tag: mkvBlockGroup,
	},
	EncryptedBlockID: {
		Name:       "EncryptedBlock",
		Definition: "Similar to  but the data inside the Block are Transformed (encrypt and/or signed). (see )",
		Type:       ebml.Binary,
	},
}

var mkvSilentTracks = ebml.Tag{
	SilentTrackNumberID: {
		Name:       "SilentTrackNumber",
		Definition: "One of the track number that are not used from now on in the stream. It could change later if not specified as silent in a further Cluster.",
		Type:       ebml.Uinteger,
	},
}

var mkvBlockGroup = ebml.Tag{
	BlockID: {
		Name:       "Block",
		Definition: "Block containing the actual data to be rendered and a timestamp relative to the Cluster Timestamp. (see )",
		Type:       ebml.Binary,
	},
	BlockVirtualID: {
		Name:       "BlockVirtual",
		Definition: "A Block with no data. It MUST be stored in the stream at the place the real Block would be in display order. (see )",
		Type:       ebml.Binary,
	},
	BlockAdditionsID: {
		Name:       "BlockAdditions",
		Definition: "Contain additional blocks to complete the main one. An EBML parser that has no knowledge of the Block structure could still see and use/skip these data.",
		Type:       ebml.Master, Tag: mkvBlockAdditions,
	},
	BlockDurationID: {
		Name:       "BlockDuration",
		Definition: "The duration of the Block (based on TimestampScale). The BlockDuration Element can be useful at the end of a Track to define the duration of the last frame (as there is no subsequent Block available), or when there is a break in a track like for subtitle tracks.",
		Type:       ebml.Uinteger,
	},
	ReferencePriorityID: {
		Name:       "ReferencePriority",
		Definition: "This frame is referenced and has the specified cache priority. In cache only a frame of the same or higher priority can replace this frame. A value of 0 means the frame is not referenced.",
		Type:       ebml.Uinteger,
	},
	ReferenceBlockID: {
		Name:       "ReferenceBlock",
		Definition: "Timestamp of another frame used as a reference (ie: B or P frame). The timestamp is relative to the block it's attached to.",
		Type:       ebml.Integer,
	},
	ReferenceVirtualID: {
		Name:       "ReferenceVirtual",
		Definition: "The Segment Position of the data that would otherwise be in position of the virtual block.",
		Type:       ebml.Integer,
	},
	CodecStateID: {
		Name:       "CodecState",
		Definition: "The new codec state to use. Data interpretation is private to the codec. This information SHOULD always be referenced by a seek entry.",
		Type:       ebml.Binary,
	},
	DiscardPaddingID: {
		Name:       "DiscardPadding",
		Definition: "Duration in nanoseconds of the silent data added to the Block (padding at the end of the Block for positive value, at the beginning of the Block for negative value). The duration of DiscardPadding is not calculated in the duration of the TrackEntry and SHOULD be discarded during playback.",
		Type:       ebml.Integer,
	},
	SlicesID: {
		Name:       "Slices",
		Definition: "Contains slices description.",
		Type:       ebml.Master, Tag: mkvSlices,
	},
	ReferenceFrameID: {
		Name:       "ReferenceFrame",
		Definition: "",
		Type:       ebml.Master, Tag: mkvReferenceFrame,
	},
}

var mkvBlockAdditions = ebml.Tag{
	BlockMoreID: {
		Name:       "BlockMore",
		Definition: "Contain the BlockAdditional and some parameters.",
		Type:       ebml.Master, Tag: mkvBlockMore,
	},
}

var mkvBlockMore = ebml.Tag{
	BlockAddIDID: {
		Name:       "BlockAddID",
		Definition: "An ID to identify the BlockAdditional level. A value of 1 means the BlockAdditional data is interpreted as additional data passed to the codec with the Block data.",
		Type:       ebml.Uinteger,
	},
	BlockAdditionalID: {
		Name:       "BlockAdditional",
		Definition: "Interpreted by the codec as it wishes (using the BlockAddID).",
		Type:       ebml.Binary,
	},
}

var mkvSlices = ebml.Tag{
	TimeSliceID: {
		Name:       "TimeSlice",
		Definition: "Contains extra time information about the data contained in the Block. Being able to interpret this Element is not REQUIRED for playback.",
		Type:       ebml.Master, Tag: mkvTimeSlice,
	},
}

var mkvTimeSlice = ebml.Tag{
	LaceNumberID: {
		Name:       "LaceNumber",
		Definition: "The reverse number of the frame in the lace (0 is the last frame, 1 is the next to last, etc). Being able to interpret this Element is not REQUIRED for playback.",
		Type:       ebml.Uinteger,
	},
	FrameNumberID: {
		Name:       "FrameNumber",
		Definition: "The number of the frame to generate from this lace with this delay (allow you to generate many frames from the same Block/Frame).",
		Type:       ebml.Uinteger,
	},
	BlockAdditionIDID: {
		Name:       "BlockAdditionID",
		Definition: "The ID of the BlockAdditional Element (0 is the main Block).",
		Type:       ebml.Uinteger,
	},
	DelayID: {
		Name:       "Delay",
		Definition: "The (scaled) delay to apply to the Element.",
		Type:       ebml.Uinteger,
	},
	SliceDurationID: {
		Name:       "SliceDuration",
		Definition: "The (scaled) duration to apply to the Element.",
		Type:       ebml.Uinteger,
	},
}

var mkvReferenceFrame = ebml.Tag{
	ReferenceOffsetID: {
		Name:       "ReferenceOffset",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	ReferenceTimestampID: {
		Name:       "ReferenceTimestamp",
		Definition: "",
		Type:       ebml.Uinteger,
	},
}

var mkvTracks = ebml.Tag{
	TrackEntryID: {
		Name:       "TrackEntry",
		Definition: "Describes a track with all Elements.",
		Type:       ebml.Master, Tag: mkvTrackEntry,
	},
}

var mkvTrackEntry = ebml.Tag{
	TrackNumberID: {
		Name:       "TrackNumber",
		Definition: "The track number as used in the Block Header (using more than 127 tracks is not encouraged, though the design allows an unlimited number).",
		Type:       ebml.Uinteger,
	},
	TrackUIDID: {
		Name:       "TrackUID",
		Definition: "A unique ID to identify the Track. This SHOULD be kept the same when making a direct stream copy of the Track to another file.",
		Type:       ebml.Uinteger,
	},
	TrackTypeID: {
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
	FlagEnabledID: {
		Name:       "FlagEnabled",
		Definition: "Set if the track is usable. (1 bit)",
		Type:       ebml.Uinteger,
	},
	FlagDefaultID: {
		Name:       "FlagDefault",
		Definition: "Set if that track (audio, video or subs) SHOULD be active if no language found matches the user preference. (1 bit)",
		Type:       ebml.Uinteger,
	},
	FlagForcedID: {
		Name:       "FlagForced",
		Definition: "Set if that track MUST be active during playback. There can be many forced track for a kind (audio, video or subs), the player SHOULD select the one which language matches the user preference or the default + forced track. Overlay MAY happen between a forced and non-forced track of the same kind. (1 bit)",
		Type:       ebml.Uinteger,
	},
	FlagLacingID: {
		Name:       "FlagLacing",
		Definition: "Set if the track MAY contain blocks using lacing. (1 bit)",
		Type:       ebml.Uinteger,
	},
	MinCacheID: {
		Name:       "MinCache",
		Definition: "The minimum number of frames a player SHOULD be able to cache during playback. If set to 0, the reference pseudo-cache system is not used.",
		Type:       ebml.Uinteger,
	},
	MaxCacheID: {
		Name:       "MaxCache",
		Definition: "The maximum cache size necessary to store referenced frames in and the current frame. 0 means no cache is needed.",
		Type:       ebml.Uinteger,
	},
	DefaultDurationID: {
		Name:       "DefaultDuration",
		Definition: "Number of nanoseconds (not scaled via TimestampScale) per frame ('frame' in the Matroska sense -- one Element put into a (Simple)Block).",
		Type:       ebml.Uinteger,
	},
	DefaultDecodedFieldDurationID: {
		Name:       "DefaultDecodedFieldDuration",
		Definition: "The period in nanoseconds (not scaled by TimestampScale) between two successive fields at the output of the decoding process (see )",
		Type:       ebml.Uinteger,
	},
	TrackTimestampScaleID: {
		Name:       "TrackTimestampScale",
		Definition: "DEPRECATED, DO NOT USE. The scale to apply on this track to work at normal speed in relation with other tracks (mostly used to adjust video speed when the audio length differs).",
		Type:       ebml.Float,
	},
	TrackOffsetID: {
		Name:       "TrackOffset",
		Definition: "A value to add to the Block's Timestamp. This can be used to adjust the playback offset of a track.",
		Type:       ebml.Integer,
	},
	MaxBlockAdditionIDID: {
		Name:       "MaxBlockAdditionID",
		Definition: "The maximum value of . A value 0 means there is no  for this track.",
		Type:       ebml.Uinteger,
	},
	BlockAdditionMappingID: {
		Name:       "BlockAdditionMapping",
		Definition: "Contains elements that describe each value of  found in the Track.",
		Type:       ebml.Master, Tag: mkvBlockAdditionMapping,
	},
	NameID: {
		Name:       "Name",
		Definition: "A human-readable track name.",
		Type:       ebml.UTF8,
	},
	LanguageID: {
		Name:       "Language",
		Definition: "Specifies the language of the track in the . This Element MUST be ignored if the LanguageIETF Element is used in the same TrackEntry.",
		Type:       ebml.String,
	},
	LanguageIETFID: {
		Name:       "LanguageIETF",
		Definition: "Specifies the language of the track according to  and using the . If this Element is used, then any Language Elements used in the same TrackEntry MUST be ignored.",
		Type:       ebml.String,
	},
	CodecIDID: {
		Name:       "CodecID",
		Definition: "An ID corresponding to the codec, see the  for more info.",
		Type:       ebml.String,
	},
	CodecPrivateID: {
		Name:       "CodecPrivate",
		Definition: "Private data only known to the codec.",
		Type:       ebml.Binary,
	},
	CodecNameID: {
		Name:       "CodecName",
		Definition: "A human-readable string specifying the codec.",
		Type:       ebml.UTF8,
	},
	AttachmentLinkID: {
		Name:       "AttachmentLink",
		Definition: "The UID of an attachment that is used by this codec.",
		Type:       ebml.Uinteger,
	},
	CodecSettingsID: {
		Name:       "CodecSettings",
		Definition: "A string describing the encoding setting used.",
		Type:       ebml.UTF8,
	},
	CodecInfoURLID: {
		Name:       "CodecInfoURL",
		Definition: "A URL to find information about the codec used.",
		Type:       ebml.String,
	},
	CodecDownloadURLID: {
		Name:       "CodecDownloadURL",
		Definition: "A URL to download about the codec used.",
		Type:       ebml.String,
	},
	CodecDecodeAllID: {
		Name:       "CodecDecodeAll",
		Definition: "The codec can decode potentially damaged data (1 bit).",
		Type:       ebml.Uinteger,
	},
	TrackOverlayID: {
		Name:       "TrackOverlay",
		Definition: "Specify that this track is an overlay track for the Track specified (in the u-integer). That means when this track has a gap (see ) the overlay track SHOULD be used instead. The order of multiple TrackOverlay matters, the first one is the one that SHOULD be used. If not found it SHOULD be the second, etc.",
		Type:       ebml.Uinteger,
	},
	CodecDelayID: {
		Name:       "CodecDelay",
		Definition: "CodecDelay is The codec-built-in delay in nanoseconds. This value MUST be subtracted from each block timestamp in order to get the actual timestamp. The value SHOULD be small so the muxing of tracks with the same actual timestamp are in the same Cluster.",
		Type:       ebml.Uinteger,
	},
	SeekPreRollID: {
		Name:       "SeekPreRoll",
		Definition: "After a discontinuity, SeekPreRoll is the duration in nanoseconds of the data the decoder MUST decode before the decoded data is valid.",
		Type:       ebml.Uinteger,
	},
	TrackTranslateID: {
		Name:       "TrackTranslate",
		Definition: "The track identification for the given Chapter Codec.",
		Type:       ebml.Master, Tag: mkvTrackTranslate,
	},
	VideoID: {
		Name:       "Video",
		Definition: "Video settings.",
		Type:       ebml.Master, Tag: mkvVideo,
	},
	AudioID: {
		Name:       "Audio",
		Definition: "Audio settings.",
		Type:       ebml.Master, Tag: mkvAudio,
	},
	TrackOperationID: {
		Name:       "TrackOperation",
		Definition: "Operation that needs to be applied on tracks to create this virtual track. For more details  on the subject.",
		Type:       ebml.Master, Tag: mkvTrackOperation,
	},
	TrickTrackUIDID: {
		Name:       "TrickTrackUID",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	TrickTrackSegmentUIDID: {
		Name:       "TrickTrackSegmentUID",
		Definition: "",
		Type:       ebml.Binary,
	},
	TrickTrackFlagID: {
		Name:       "TrickTrackFlag",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	TrickMasterTrackUIDID: {
		Name:       "TrickMasterTrackUID",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	TrickMasterTrackSegmentUIDID: {
		Name:       "TrickMasterTrackSegmentUID",
		Definition: "",
		Type:       ebml.Binary,
	},
	ContentEncodingsID: {
		Name:       "ContentEncodings",
		Definition: "Settings for several content encoding mechanisms like compression or encryption.",
		Type:       ebml.Master, Tag: mkvContentEncodings,
	},
}

var mkvBlockAdditionMapping = ebml.Tag{
	BlockAddIDValueID: {
		Name:       "BlockAddIDValue",
		Definition: "The  value being described. To keep MaxBlockAdditionID as low as possible, small values SHOULD be used.",
		Type:       ebml.Uinteger,
	},
	BlockAddIDNameID: {
		Name:       "BlockAddIDName",
		Definition: "A human-friendly name describing the type of BlockAdditional data as defined by the associated Block Additional Mapping.",
		Type:       ebml.String,
	},
	BlockAddIDTypeID: {
		Name:       "BlockAddIDType",
		Definition: "Stores the registered identifer of the Block Additional Mapping to define how the BlockAdditional data should be handled.",
		Type:       ebml.Uinteger,
	},
	BlockAddIDExtraDataID: {
		Name:       "BlockAddIDExtraData",
		Definition: "Extra binary data that the BlockAddIDType can use to interpret the BlockAdditional data. The intepretation of the binary data depends on the BlockAddIDType value and the corresponding Block Additional Mapping.",
		Type:       ebml.Binary,
	},
}

var mkvTrackTranslate = ebml.Tag{
	TrackTranslateEditionUIDID: {
		Name:       "TrackTranslateEditionUID",
		Definition: "Specify an edition UID on which this translation applies. When not specified, it means for all editions found in the Segment.",
		Type:       ebml.Uinteger,
	},
	TrackTranslateCodecID: {
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
	TrackTranslateTrackIDID: {
		Name:       "TrackTranslateTrackID",
		Definition: "The binary value used to represent this track in the chapter codec data. The format depends on the  used.",
		Type:       ebml.Binary,
	},
}

var mkvVideo = ebml.Tag{
	FlagInterlacedID: {
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
	FieldOrderID: {
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
	StereoModeID: {
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
	AlphaModeID: {
		Name:       "AlphaMode",
		Definition: "Alpha Video Mode. Presence of this Element indicates that the BlockAdditional Element could contain Alpha data.",
		Type:       ebml.Uinteger,
	},
	OldStereoModeID: {
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
	PixelWidthID: {
		Name:       "PixelWidth",
		Definition: "Width of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	PixelHeightID: {
		Name:       "PixelHeight",
		Definition: "Height of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	PixelCropBottomID: {
		Name:       "PixelCropBottom",
		Definition: "The number of video pixels to remove at the bottom of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropTopID: {
		Name:       "PixelCropTop",
		Definition: "The number of video pixels to remove at the top of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropLeftID: {
		Name:       "PixelCropLeft",
		Definition: "The number of video pixels to remove on the left of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropRightID: {
		Name:       "PixelCropRight",
		Definition: "The number of video pixels to remove on the right of the image.",
		Type:       ebml.Uinteger,
	},
	DisplayWidthID: {
		Name:       "DisplayWidth",
		Definition: "Width of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	DisplayHeightID: {
		Name:       "DisplayHeight",
		Definition: "Height of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	DisplayUnitID: {
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
	AspectRatioTypeID: {
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
	ColourSpaceID: {
		Name:       "ColourSpace",
		Definition: "Specify the pixel format used for the Track's data as a FourCC. This value is similar in scope to the biCompression value of AVI's BITMAPINFOHEADER.",
		Type:       ebml.Binary,
	},
	GammaValueID: {
		Name:       "GammaValue",
		Definition: "Gamma Value.",
		Type:       ebml.Float,
	},
	FrameRateID: {
		Name:       "FrameRate",
		Definition: "Number of frames per second.  only.",
		Type:       ebml.Float,
	},
	ColourID: {
		Name:       "Colour",
		Definition: "Settings describing the colour format.",
		Type:       ebml.Master, Tag: mkvColour,
	},
	ProjectionID: {
		Name:       "Projection",
		Definition: "Describes the video projection details. Used to render spherical and VR videos.",
		Type:       ebml.Master, Tag: mkvProjection,
	},
}

var mkvColour = ebml.Tag{
	MatrixCoefficientsID: {
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
	BitsPerChannelID: {
		Name:       "BitsPerChannel",
		Definition: "Number of decoded bits per channel. A value of 0 indicates that the BitsPerChannel is unspecified.",
		Type:       ebml.Uinteger,
	},
	ChromaSubsamplingHorzID: {
		Name:       "ChromaSubsamplingHorz",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed horizontally. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingHorz SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	ChromaSubsamplingVertID: {
		Name:       "ChromaSubsamplingVert",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed vertically. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingVert SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	CbSubsamplingHorzID: {
		Name:       "CbSubsamplingHorz",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed horizontally. This is additive with ChromaSubsamplingHorz. Example: For video with 4:2:1 chroma subsampling, the ChromaSubsamplingHorz SHOULD be set to 1 and CbSubsamplingHorz SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	CbSubsamplingVertID: {
		Name:       "CbSubsamplingVert",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed vertically. This is additive with ChromaSubsamplingVert.",
		Type:       ebml.Uinteger,
	},
	ChromaSitingHorzID: {
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
	ChromaSitingVertID: {
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
	RangeID: {
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
	TransferCharacteristicsID: {
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
	PrimariesID: {
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
	MaxCLLID: {
		Name:       "MaxCLL",
		Definition: "Maximum brightness of a single pixel (Maximum Content Light Level) in candelas per square meter (cd/m²).",
		Type:       ebml.Uinteger,
	},
	MaxFALLID: {
		Name:       "MaxFALL",
		Definition: "Maximum brightness of a single full frame (Maximum Frame-Average Light Level) in candelas per square meter (cd/m²).",
		Type:       ebml.Uinteger,
	},
	MasteringMetadataID: {
		Name:       "MasteringMetadata",
		Definition: "SMPTE 2086 mastering data.",
		Type:       ebml.Master, Tag: mkvMasteringMetadata,
	},
}

var mkvMasteringMetadata = ebml.Tag{
	PrimaryRChromaticityXID: {
		Name:       "PrimaryRChromaticityX",
		Definition: "Red X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryRChromaticityYID: {
		Name:       "PrimaryRChromaticityY",
		Definition: "Red Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryGChromaticityXID: {
		Name:       "PrimaryGChromaticityX",
		Definition: "Green X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryGChromaticityYID: {
		Name:       "PrimaryGChromaticityY",
		Definition: "Green Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryBChromaticityXID: {
		Name:       "PrimaryBChromaticityX",
		Definition: "Blue X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryBChromaticityYID: {
		Name:       "PrimaryBChromaticityY",
		Definition: "Blue Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	WhitePointChromaticityXID: {
		Name:       "WhitePointChromaticityX",
		Definition: "White X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	WhitePointChromaticityYID: {
		Name:       "WhitePointChromaticityY",
		Definition: "White Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	LuminanceMaxID: {
		Name:       "LuminanceMax",
		Definition: "Maximum luminance. Represented in candelas per square meter (cd/m²).",
		Type:       ebml.Float,
	},
	LuminanceMinID: {
		Name:       "LuminanceMin",
		Definition: "Minimum luminance. Represented in candelas per square meter (cd/m²).",
		Type:       ebml.Float,
	},
}

var mkvProjection = ebml.Tag{
	ProjectionTypeID: {
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
	ProjectionPrivateID: {
		Name:       "ProjectionPrivate",
		Definition: "Private data that only applies to a specific projection.SemanticsIf ProjectionType equals 0 (Rectangular),\n     then this element must not be present.If ProjectionType equals 1 (Equirectangular), then this element must be present and contain the same binary data that would be stored inside\n      an ISOBMFF Equirectangular Projection Box ('equi').If ProjectionType equals 2 (Cubemap), then this element must be present and contain the same binary data that would be stored \n      inside an ISOBMFF Cubemap Projection Box ('cbmp').If ProjectionType equals 3 (Mesh), then this element must be present and contain the same binary data that would be stored inside\n       an ISOBMFF Mesh Projection Box ('mshp').Note: ISOBMFF box size and fourcc fields are not included in the binary data, but the FullBox version and flag fields are. This is to avoid \n       redundant framing information while preserving versioning and semantics between the two container formats.",
		Type:       ebml.Binary,
	},
	ProjectionPoseYawID: {
		Name:       "ProjectionPoseYaw",
		Definition: "Specifies a yaw rotation to the projection.SemanticsValue represents a clockwise rotation, in degrees, around the up vector. This rotation must be applied before any ProjectionPosePitch or ProjectionPoseRoll rotations. The value of this field should be in the -180 to 180 degree range.",
		Type:       ebml.Float,
	},
	ProjectionPosePitchID: {
		Name:       "ProjectionPosePitch",
		Definition: "Specifies a pitch rotation to the projection.SemanticsValue represents a counter-clockwise rotation, in degrees, around the right vector. This rotation must be applied after the ProjectionPoseYaw rotation and before the ProjectionPoseRoll rotation. The value of this field should be in the -90 to 90 degree range.",
		Type:       ebml.Float,
	},
	ProjectionPoseRollID: {
		Name:       "ProjectionPoseRoll",
		Definition: "Specifies a roll rotation to the projection.SemanticsValue represents a counter-clockwise rotation, in degrees, around the forward vector. This rotation must be applied after the ProjectionPoseYaw and ProjectionPosePitch rotations. The value of this field should be in the -180 to 180 degree range.",
		Type:       ebml.Float,
	},
}

var mkvAudio = ebml.Tag{
	SamplingFrequencyID: {
		Name:       "SamplingFrequency",
		Definition: "Sampling frequency in Hz.",
		Type:       ebml.Float,
	},
	OutputSamplingFrequencyID: {
		Name:       "OutputSamplingFrequency",
		Definition: "Real output sampling frequency in Hz (used for SBR techniques).",
		Type:       ebml.Float,
	},
	ChannelsID: {
		Name:       "Channels",
		Definition: "Numbers of channels in the track.",
		Type:       ebml.Uinteger,
	},
	ChannelPositionsID: {
		Name:       "ChannelPositions",
		Definition: "Table of horizontal angles for each successive channel, see .",
		Type:       ebml.Binary,
	},
	BitDepthID: {
		Name:       "BitDepth",
		Definition: "Bits per sample, mostly used for PCM.",
		Type:       ebml.Uinteger,
	},
}

var mkvTrackOperation = ebml.Tag{
	TrackCombinePlanesID: {
		Name:       "TrackCombinePlanes",
		Definition: "Contains the list of all video plane tracks that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: mkvTrackCombinePlanes,
	},
	TrackJoinBlocksID: {
		Name:       "TrackJoinBlocks",
		Definition: "Contains the list of all tracks whose Blocks need to be combined to create this virtual track",
		Type:       ebml.Master, Tag: mkvTrackJoinBlocks,
	},
}

var mkvTrackCombinePlanes = ebml.Tag{
	TrackPlaneID: {
		Name:       "TrackPlane",
		Definition: "Contains a video plane track that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: mkvTrackPlane,
	},
}

var mkvTrackPlane = ebml.Tag{
	TrackPlaneUIDID: {
		Name:       "TrackPlaneUID",
		Definition: "The trackUID number of the track representing the plane.",
		Type:       ebml.Uinteger,
	},
	TrackPlaneTypeID: {
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
	TrackJoinUIDID: {
		Name:       "TrackJoinUID",
		Definition: "The trackUID number of a track whose blocks are used to create this virtual track.",
		Type:       ebml.Uinteger,
	},
}

var mkvContentEncodings = ebml.Tag{
	ContentEncodingID: {
		Name:       "ContentEncoding",
		Definition: "Settings for one content encoding like compression or encryption.",
		Type:       ebml.Master, Tag: mkvContentEncoding,
	},
}

var mkvContentEncoding = ebml.Tag{
	ContentEncodingOrderID: {
		Name:       "ContentEncodingOrder",
		Definition: "Tells when this modification was used during encoding/muxing starting with 0 and counting upwards. The decoder/demuxer has to start with the highest order number it finds and work its way down. This value has to be unique over all ContentEncodingOrder Elements in the TrackEntry that contains this ContentEncodingOrder element.",
		Type:       ebml.Uinteger,
	},
	ContentEncodingScopeID: {
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
	ContentEncodingTypeID: {
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
	ContentCompressionID: {
		Name:       "ContentCompression",
		Definition: "Settings describing the compression used. This Element MUST be present if the value of ContentEncodingType is 0 and absent otherwise. Each block MUST be decompressable even if no previous block is available in order not to prevent seeking.",
		Type:       ebml.Master, Tag: mkvContentCompression,
	},
	ContentEncryptionID: {
		Name:       "ContentEncryption",
		Definition: "Settings describing the encryption used. This Element MUST be present if the value of `ContentEncodingType` is 1 (encryption) and MUST be ignored otherwise.",
		Type:       ebml.Master, Tag: mkvContentEncryption,
	},
}

var mkvContentCompression = ebml.Tag{
	ContentCompAlgoID: {
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
	ContentCompSettingsID: {
		Name:       "ContentCompSettings",
		Definition: "Settings that might be needed by the decompressor. For Header Stripping (`ContentCompAlgo`=3), the bytes that were removed from the beginning of each frames of the track.",
		Type:       ebml.Binary,
	},
}

var mkvContentEncryption = ebml.Tag{
	ContentEncAlgoID: {
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
	ContentEncKeyIDID: {
		Name:       "ContentEncKeyID",
		Definition: "For public key algorithms this is the ID of the public key the the data was encrypted with.",
		Type:       ebml.Binary,
	},
	ContentEncAESSettingsID: {
		Name:       "ContentEncAESSettings",
		Definition: "Settings describing the encryption algorithm used. If `ContentEncAlgo` != 5 this MUST be ignored.",
		Type:       ebml.Master, Tag: mkvContentEncAESSettings,
	},
	ContentSignatureID: {
		Name:       "ContentSignature",
		Definition: "A cryptographic signature of the contents.",
		Type:       ebml.Binary,
	},
	ContentSigKeyIDID: {
		Name:       "ContentSigKeyID",
		Definition: "This is the ID of the private key the data was signed with.",
		Type:       ebml.Binary,
	},
	ContentSigAlgoID: {
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
	ContentSigHashAlgoID: {
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
	AESSettingsCipherModeID: {
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
	CuePointID: {
		Name:       "CuePoint",
		Definition: "Contains all information relative to a seek point in the Segment.",
		Type:       ebml.Master, Tag: mkvCuePoint,
	},
}

var mkvCuePoint = ebml.Tag{
	CueTimeID: {
		Name:       "CueTime",
		Definition: "Absolute timestamp according to the Segment time base.",
		Type:       ebml.Uinteger,
	},
	CueTrackPositionsID: {
		Name:       "CueTrackPositions",
		Definition: "Contain positions for different tracks corresponding to the timestamp.",
		Type:       ebml.Master, Tag: mkvCueTrackPositions,
	},
}

var mkvCueTrackPositions = ebml.Tag{
	CueTrackID: {
		Name:       "CueTrack",
		Definition: "The track for which a position is given.",
		Type:       ebml.Uinteger,
	},
	CueClusterPositionID: {
		Name:       "CueClusterPosition",
		Definition: "The Segment Position of the Cluster containing the associated Block.",
		Type:       ebml.Uinteger,
	},
	CueRelativePositionID: {
		Name:       "CueRelativePosition",
		Definition: "The relative position inside the Cluster of the referenced SimpleBlock or BlockGroup with 0 being the first possible position for an Element inside that Cluster.",
		Type:       ebml.Uinteger,
	},
	CueDurationID: {
		Name:       "CueDuration",
		Definition: "The duration of the block according to the Segment time base. If missing the track's DefaultDuration does not apply and no duration information is available in terms of the cues.",
		Type:       ebml.Uinteger,
	},
	CueBlockNumberID: {
		Name:       "CueBlockNumber",
		Definition: "Number of the Block in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	CueCodecStateID: {
		Name:       "CueCodecState",
		Definition: "The Segment Position of the Codec State corresponding to this Cue Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
	CueReferenceID: {
		Name:       "CueReference",
		Definition: "The Clusters containing the referenced Blocks.",
		Type:       ebml.Master, Tag: mkvCueReference,
	},
}

var mkvCueReference = ebml.Tag{
	CueRefTimeID: {
		Name:       "CueRefTime",
		Definition: "Timestamp of the referenced Block.",
		Type:       ebml.Uinteger,
	},
	CueRefClusterID: {
		Name:       "CueRefCluster",
		Definition: "The Segment Position of the Cluster containing the referenced Block.",
		Type:       ebml.Uinteger,
	},
	CueRefNumberID: {
		Name:       "CueRefNumber",
		Definition: "Number of the referenced Block of Track X in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	CueRefCodecStateID: {
		Name:       "CueRefCodecState",
		Definition: "The Segment Position of the Codec State corresponding to this referenced Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
}

var mkvAttachments = ebml.Tag{
	AttachedFileID: {
		Name:       "AttachedFile",
		Definition: "An attached file.",
		Type:       ebml.Master, Tag: mkvAttachedFile,
	},
}

var mkvAttachedFile = ebml.Tag{
	FileDescriptionID: {
		Name:       "FileDescription",
		Definition: "A human-friendly name for the attached file.",
		Type:       ebml.UTF8,
	},
	FileNameID: {
		Name:       "FileName",
		Definition: "Filename of the attached file.",
		Type:       ebml.UTF8,
	},
	FileMimeTypeID: {
		Name:       "FileMimeType",
		Definition: "MIME type of the file.",
		Type:       ebml.String,
	},
	FileDataID: {
		Name:       "FileData",
		Definition: "The data of the file.",
		Type:       ebml.Binary,
	},
	FileUIDID: {
		Name:       "FileUID",
		Definition: "Unique ID representing the file, as random as possible.",
		Type:       ebml.Uinteger,
	},
	FileReferralID: {
		Name:       "FileReferral",
		Definition: "A binary value that a track/codec can refer to when the attachment is needed.",
		Type:       ebml.Binary,
	},
	FileUsedStartTimeID: {
		Name:       "FileUsedStartTime",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	FileUsedEndTimeID: {
		Name:       "FileUsedEndTime",
		Definition: "",
		Type:       ebml.Uinteger,
	},
}

var mkvChapters = ebml.Tag{
	EditionEntryID: {
		Name:       "EditionEntry",
		Definition: "Contains all information about a Segment edition.",
		Type:       ebml.Master, Tag: mkvEditionEntry,
	},
}

var mkvEditionEntry = ebml.Tag{
	EditionUIDID: {
		Name:       "EditionUID",
		Definition: "A unique ID to identify the edition. It's useful for tagging an edition.",
		Type:       ebml.Uinteger,
	},
	EditionFlagHiddenID: {
		Name:       "EditionFlagHidden",
		Definition: "If an edition is hidden (1), it SHOULD NOT be available to the user interface (but still to Control Tracks; see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	EditionFlagDefaultID: {
		Name:       "EditionFlagDefault",
		Definition: "If a flag is set (1) the edition SHOULD be used as the default one. (1 bit)",
		Type:       ebml.Uinteger,
	},
	EditionFlagOrderedID: {
		Name:       "EditionFlagOrdered",
		Definition: "Specify if the chapters can be defined multiple times and the order to play them is enforced. (1 bit)",
		Type:       ebml.Uinteger,
	},
	ChapterAtomID: {
		Name:       "ChapterAtom",
		Definition: "Contains the atom information to use as the chapter atom (apply to all tracks).",
		Type:       ebml.Master, Tag: mkvChapterAtom,
	},
}

var mkvChapterAtom = ebml.Tag{
	ChapterUIDID: {
		Name:       "ChapterUID",
		Definition: "A unique ID to identify the Chapter.",
		Type:       ebml.Uinteger,
	},
	ChapterStringUIDID: {
		Name:       "ChapterStringUID",
		Definition: "A unique string ID to identify the Chapter. Use for .",
		Type:       ebml.UTF8,
	},
	ChapterTimeStartID: {
		Name:       "ChapterTimeStart",
		Definition: "Timestamp of the start of Chapter (not scaled).",
		Type:       ebml.Uinteger,
	},
	ChapterTimeEndID: {
		Name:       "ChapterTimeEnd",
		Definition: "Timestamp of the end of Chapter (timestamp excluded, not scaled).",
		Type:       ebml.Uinteger,
	},
	ChapterFlagHiddenID: {
		Name:       "ChapterFlagHidden",
		Definition: "If a chapter is hidden (1), it SHOULD NOT be available to the user interface (but still to Control Tracks; see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	ChapterFlagEnabledID: {
		Name:       "ChapterFlagEnabled",
		Definition: "Specify whether the chapter is enabled. It can be enabled/disabled by a Control Track. When disabled, the movie SHOULD skip all the content between the TimeStart and TimeEnd of this chapter (see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	ChapterSegmentUIDID: {
		Name:       "ChapterSegmentUID",
		Definition: "The SegmentUID of another Segment to play during this chapter.",
		Type:       ebml.Binary,
	},
	ChapterSegmentEditionUIDID: {
		Name:       "ChapterSegmentEditionUID",
		Definition: "The EditionUID to play from the Segment linked in ChapterSegmentUID. If ChapterSegmentEditionUID is undeclared then no Edition of the linked Segment is used.",
		Type:       ebml.Uinteger,
	},
	ChapterPhysicalEquivID: {
		Name:       "ChapterPhysicalEquiv",
		Definition: "Specify the physical equivalent of this ChapterAtom like \"DVD\" (60) or \"SIDE\" (50), see .",
		Type:       ebml.Uinteger,
	},
	ChapterTrackID: {
		Name:       "ChapterTrack",
		Definition: "List of tracks on which the chapter applies. If this Element is not present, all tracks apply",
		Type:       ebml.Master, Tag: mkvChapterTrack,
	},
	ChapterDisplayID: {
		Name:       "ChapterDisplay",
		Definition: "Contains all possible strings to use for the chapter display.",
		Type:       ebml.Master, Tag: mkvChapterDisplay,
	},
	ChapProcessID: {
		Name:       "ChapProcess",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: mkvChapProcess,
	},
}

var mkvChapterTrack = ebml.Tag{
	ChapterTrackUIDID: {
		Name:       "ChapterTrackUID",
		Definition: "UID of the Track to apply this chapter too. In the absence of a control track, choosing this chapter will select the listed Tracks and deselect unlisted tracks. Absence of this Element indicates that the Chapter SHOULD be applied to any currently used Tracks.",
		Type:       ebml.Uinteger,
	},
}

var mkvChapterDisplay = ebml.Tag{
	ChapStringID: {
		Name:       "ChapString",
		Definition: "Contains the string to use as the chapter atom.",
		Type:       ebml.UTF8,
	},
	ChapLanguageID: {
		Name:       "ChapLanguage",
		Definition: "The languages corresponding to the string, in the . This Element MUST be ignored if the ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
	ChapLanguageIETFID: {
		Name:       "ChapLanguageIETF",
		Definition: "Specifies the language used in the ChapString according to  and using the . If this Element is used, then any ChapLanguage Elements used in the same ChapterDisplay MUST be ignored.",
		Type:       ebml.String,
	},
	ChapCountryID: {
		Name:       "ChapCountry",
		Definition: "The countries corresponding to the string, same 2 octets as in . This Element MUST be ignored if the ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
}

var mkvChapProcess = ebml.Tag{
	ChapProcessCodecIDID: {
		Name:       "ChapProcessCodecID",
		Definition: "Contains the type of the codec used for the processing. A value of 0 means native Matroska processing (to be defined), a value of 1 means the  command set is used. More codec IDs can be added later.",
		Type:       ebml.Uinteger,
	},
	ChapProcessPrivateID: {
		Name:       "ChapProcessPrivate",
		Definition: "Some optional data attached to the ChapProcessCodecID information. , it is the \"DVD level\" equivalent.",
		Type:       ebml.Binary,
	},
	ChapProcessCommandID: {
		Name:       "ChapProcessCommand",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: mkvChapProcessCommand,
	},
}

var mkvChapProcessCommand = ebml.Tag{
	ChapProcessTimeID: {
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
	ChapProcessDataID: {
		Name:       "ChapProcessData",
		Definition: "Contains the command information. The data SHOULD be interpreted depending on the ChapProcessCodecID value. , the data correspond to the binary DVD cell pre/post commands.",
		Type:       ebml.Binary,
	},
}

var mkvTags = ebml.Tag{
	TagID: {
		Name:       "Tag",
		Definition: "A single metadata descriptor.",
		Type:       ebml.Master, Tag: mkvTag,
	},
}

var mkvTag = ebml.Tag{
	TargetsID: {
		Name:       "Targets",
		Definition: "Specifies which other elements the metadata represented by the Tag applies to. If empty or not present, then the Tag describes everything in the Segment.",
		Type:       ebml.Master, Tag: mkvTargets,
	},
	SimpleTagID: {
		Name:       "SimpleTag",
		Definition: "Contains general information about the target.",
		Type:       ebml.Master, Tag: mkvSimpleTag,
	},
}

var mkvTargets = ebml.Tag{
	TargetTypeValueID: {
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
	TargetTypeID: {
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
	TagTrackUIDID: {
		Name:       "TagTrackUID",
		Definition: "A unique ID to identify the Track(s) the tags belong to. If the value is 0 at this level, the tags apply to all tracks in the Segment.",
		Type:       ebml.Uinteger,
	},
	TagEditionUIDID: {
		Name:       "TagEditionUID",
		Definition: "A unique ID to identify the EditionEntry(s) the tags belong to. If the value is 0 at this level, the tags apply to all editions in the Segment.",
		Type:       ebml.Uinteger,
	},
	TagChapterUIDID: {
		Name:       "TagChapterUID",
		Definition: "A unique ID to identify the Chapter(s) the tags belong to. If the value is 0 at this level, the tags apply to all chapters in the Segment.",
		Type:       ebml.Uinteger,
	},
	TagAttachmentUIDID: {
		Name:       "TagAttachmentUID",
		Definition: "A unique ID to identify the Attachment(s) the tags belong to. If the value is 0 at this level, the tags apply to all the attachments in the Segment.",
		Type:       ebml.Uinteger,
	},
}

var mkvSimpleTag = ebml.Tag{
	TagNameID: {
		Name:       "TagName",
		Definition: "The name of the Tag that is going to be stored.",
		Type:       ebml.UTF8,
	},
	TagLanguageID: {
		Name:       "TagLanguage",
		Definition: "Specifies the language of the tag specified, in the . This Element MUST be ignored if the TagLanguageIETF Element is used within the same SimpleTag Element.",
		Type:       ebml.String,
	},
	TagLanguageIETFID: {
		Name:       "TagLanguageIETF",
		Definition: "Specifies the language used in the TagString according to  and using the . If this Element is used, then any TagLanguage Elements used in the same SimpleTag MUST be ignored.",
		Type:       ebml.String,
	},
	TagDefaultID: {
		Name:       "TagDefault",
		Definition: "A boolean value to indicate if this is the default/original language to use for the given tag.",
		Type:       ebml.Uinteger,
	},
	TagStringID: {
		Name:       "TagString",
		Definition: "The value of the Tag.",
		Type:       ebml.UTF8,
	},
	TagBinaryID: {
		Name:       "TagBinary",
		Definition: "The values of the Tag if it is binary. Note that this cannot be used in the same SimpleTag as TagString.",
		Type:       ebml.Binary,
	},
}
