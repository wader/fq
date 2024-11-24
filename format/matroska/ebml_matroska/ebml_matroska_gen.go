// Code below generated from ebml_matroska.xml
package ebml_matroska

import (
	"github.com/wader/fq/format/matroska/ebml"
)

var RootElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:       RootID,
		ParentID: -1,
		Name:     "",
	},
	Master: map[ebml.ID]ebml.Element{
		ebml.HeaderID: ebml.Header,
		SegmentID:     SegmentElement,
	},
}

const (
	RootID                        = ebml.RootID
	EBMLMaxIDLengthID             = 0x42f2
	EBMLMaxSizeLengthID           = 0x42f3
	SegmentID                     = 0x18538067
	SeekHeadID                    = 0x114d9b74
	SeekID                        = 0x4dbb
	SeekIDID                      = 0x53ab
	SeekPositionID                = 0x53ac
	InfoID                        = 0x1549a966
	SegmentUUIDID                 = 0x73a4
	SegmentFilenameID             = 0x7384
	PrevUUIDID                    = 0x3cb923
	PrevFilenameID                = 0x3c83ab
	NextUUIDID                    = 0x3eb923
	NextFilenameID                = 0x3e83bb
	SegmentFamilyID               = 0x4444
	ChapterTranslateID            = 0x6924
	ChapterTranslateIDID          = 0x69a5
	ChapterTranslateCodecID       = 0x69bf
	ChapterTranslateEditionUIDID  = 0x69fc
	TimestampScaleID              = 0x2ad7b1
	DurationID                    = 0x4489
	DateUTCID                     = 0x4461
	TitleID                       = 0x7ba9
	MuxingAppID                   = 0x4d80
	WritingAppID                  = 0x5741
	ClusterID                     = 0x1f43b675
	TimestampID                   = 0xe7
	SilentTracksID                = 0x5854
	SilentTrackNumberID           = 0x58d7
	PositionID                    = 0xa7
	PrevSizeID                    = 0xab
	SimpleBlockID                 = 0xa3
	BlockGroupID                  = 0xa0
	BlockID                       = 0xa1
	BlockVirtualID                = 0xa2
	BlockAdditionsID              = 0x75a1
	BlockMoreID                   = 0xa6
	BlockAdditionalID             = 0xa5
	BlockAddIDID                  = 0xee
	BlockDurationID               = 0x9b
	ReferencePriorityID           = 0xfa
	ReferenceBlockID              = 0xfb
	ReferenceVirtualID            = 0xfd
	CodecStateID                  = 0xa4
	DiscardPaddingID              = 0x75a2
	SlicesID                      = 0x8e
	TimeSliceID                   = 0xe8
	LaceNumberID                  = 0xcc
	FrameNumberID                 = 0xcd
	BlockAdditionIDID             = 0xcb
	DelayID                       = 0xce
	SliceDurationID               = 0xcf
	ReferenceFrameID              = 0xc8
	ReferenceOffsetID             = 0xc9
	ReferenceTimestampID          = 0xca
	EncryptedBlockID              = 0xaf
	TracksID                      = 0x1654ae6b
	TrackEntryID                  = 0xae
	TrackNumberID                 = 0xd7
	TrackUIDID                    = 0x73c5
	TrackTypeID                   = 0x83
	FlagEnabledID                 = 0xb9
	FlagDefaultID                 = 0x88
	FlagForcedID                  = 0x55aa
	FlagHearingImpairedID         = 0x55ab
	FlagVisualImpairedID          = 0x55ac
	FlagTextDescriptionsID        = 0x55ad
	FlagOriginalID                = 0x55ae
	FlagCommentaryID              = 0x55af
	FlagLacingID                  = 0x9c
	MinCacheID                    = 0x6de7
	MaxCacheID                    = 0x6df8
	DefaultDurationID             = 0x23e383
	DefaultDecodedFieldDurationID = 0x234e7a
	TrackTimestampScaleID         = 0x23314f
	TrackOffsetID                 = 0x537f
	MaxBlockAdditionIDID          = 0x55ee
	BlockAdditionMappingID        = 0x41e4
	BlockAddIDValueID             = 0x41f0
	BlockAddIDNameID              = 0x41a4
	BlockAddIDTypeID              = 0x41e7
	BlockAddIDExtraDataID         = 0x41ed
	NameID                        = 0x536e
	LanguageID                    = 0x22b59c
	LanguageBCP47ID               = 0x22b59d
	CodecIDID                     = 0x86
	CodecPrivateID                = 0x63a2
	CodecNameID                   = 0x258688
	AttachmentLinkID              = 0x7446
	CodecSettingsID               = 0x3a9697
	CodecInfoURLID                = 0x3b4040
	CodecDownloadURLID            = 0x26b240
	CodecDecodeAllID              = 0xaa
	TrackOverlayID                = 0x6fab
	CodecDelayID                  = 0x56aa
	SeekPreRollID                 = 0x56bb
	TrackTranslateID              = 0x6624
	TrackTranslateTrackIDID       = 0x66a5
	TrackTranslateCodecID         = 0x66bf
	TrackTranslateEditionUIDID    = 0x66fc
	VideoID                       = 0xe0
	FlagInterlacedID              = 0x9a
	FieldOrderID                  = 0x9d
	StereoModeID                  = 0x53b8
	AlphaModeID                   = 0x53c0
	OldStereoModeID               = 0x53b9
	PixelWidthID                  = 0xb0
	PixelHeightID                 = 0xba
	PixelCropBottomID             = 0x54aa
	PixelCropTopID                = 0x54bb
	PixelCropLeftID               = 0x54cc
	PixelCropRightID              = 0x54dd
	DisplayWidthID                = 0x54b0
	DisplayHeightID               = 0x54ba
	DisplayUnitID                 = 0x54b2
	AspectRatioTypeID             = 0x54b3
	UncompressedFourCCID          = 0x2eb524
	GammaValueID                  = 0x2fb523
	FrameRateID                   = 0x2383e3
	ColourID                      = 0x55b0
	MatrixCoefficientsID          = 0x55b1
	BitsPerChannelID              = 0x55b2
	ChromaSubsamplingHorzID       = 0x55b3
	ChromaSubsamplingVertID       = 0x55b4
	CbSubsamplingHorzID           = 0x55b5
	CbSubsamplingVertID           = 0x55b6
	ChromaSitingHorzID            = 0x55b7
	ChromaSitingVertID            = 0x55b8
	RangeID                       = 0x55b9
	TransferCharacteristicsID     = 0x55ba
	PrimariesID                   = 0x55bb
	MaxCLLID                      = 0x55bc
	MaxFALLID                     = 0x55bd
	MasteringMetadataID           = 0x55d0
	PrimaryRChromaticityXID       = 0x55d1
	PrimaryRChromaticityYID       = 0x55d2
	PrimaryGChromaticityXID       = 0x55d3
	PrimaryGChromaticityYID       = 0x55d4
	PrimaryBChromaticityXID       = 0x55d5
	PrimaryBChromaticityYID       = 0x55d6
	WhitePointChromaticityXID     = 0x55d7
	WhitePointChromaticityYID     = 0x55d8
	LuminanceMaxID                = 0x55d9
	LuminanceMinID                = 0x55da
	ProjectionID                  = 0x7670
	ProjectionTypeID              = 0x7671
	ProjectionPrivateID           = 0x7672
	ProjectionPoseYawID           = 0x7673
	ProjectionPosePitchID         = 0x7674
	ProjectionPoseRollID          = 0x7675
	AudioID                       = 0xe1
	SamplingFrequencyID           = 0xb5
	OutputSamplingFrequencyID     = 0x78b5
	ChannelsID                    = 0x9f
	ChannelPositionsID            = 0x7d7b
	BitDepthID                    = 0x6264
	EmphasisID                    = 0x52f1
	TrackOperationID              = 0xe2
	TrackCombinePlanesID          = 0xe3
	TrackPlaneID                  = 0xe4
	TrackPlaneUIDID               = 0xe5
	TrackPlaneTypeID              = 0xe6
	TrackJoinBlocksID             = 0xe9
	TrackJoinUIDID                = 0xed
	TrickTrackUIDID               = 0xc0
	TrickTrackSegmentUIDID        = 0xc1
	TrickTrackFlagID              = 0xc6
	TrickMasterTrackUIDID         = 0xc7
	TrickMasterTrackSegmentUIDID  = 0xc4
	ContentEncodingsID            = 0x6d80
	ContentEncodingID             = 0x6240
	ContentEncodingOrderID        = 0x5031
	ContentEncodingScopeID        = 0x5032
	ContentEncodingTypeID         = 0x5033
	ContentCompressionID          = 0x5034
	ContentCompAlgoID             = 0x4254
	ContentCompSettingsID         = 0x4255
	ContentEncryptionID           = 0x5035
	ContentEncAlgoID              = 0x47e1
	ContentEncKeyIDID             = 0x47e2
	ContentEncAESSettingsID       = 0x47e7
	AESSettingsCipherModeID       = 0x47e8
	ContentSignatureID            = 0x47e3
	ContentSigKeyIDID             = 0x47e4
	ContentSigAlgoID              = 0x47e5
	ContentSigHashAlgoID          = 0x47e6
	CuesID                        = 0x1c53bb6b
	CuePointID                    = 0xbb
	CueTimeID                     = 0xb3
	CueTrackPositionsID           = 0xb7
	CueTrackID                    = 0xf7
	CueClusterPositionID          = 0xf1
	CueRelativePositionID         = 0xf0
	CueDurationID                 = 0xb2
	CueBlockNumberID              = 0x5378
	CueCodecStateID               = 0xea
	CueReferenceID                = 0xdb
	CueRefTimeID                  = 0x96
	CueRefClusterID               = 0x97
	CueRefNumberID                = 0x535f
	CueRefCodecStateID            = 0xeb
	AttachmentsID                 = 0x1941a469
	AttachedFileID                = 0x61a7
	FileDescriptionID             = 0x467e
	FileNameID                    = 0x466e
	FileMediaTypeID               = 0x4660
	FileDataID                    = 0x465c
	FileUIDID                     = 0x46ae
	FileReferralID                = 0x4675
	FileUsedStartTimeID           = 0x4661
	FileUsedEndTimeID             = 0x4662
	ChaptersID                    = 0x1043a770
	EditionEntryID                = 0x45b9
	EditionUIDID                  = 0x45bc
	EditionFlagHiddenID           = 0x45bd
	EditionFlagDefaultID          = 0x45db
	EditionFlagOrderedID          = 0x45dd
	EditionDisplayID              = 0x4520
	EditionStringID               = 0x4521
	EditionLanguageIETFID         = 0x45e4
	ChapterAtomID                 = 0xb6
	ChapterUIDID                  = 0x73c4
	ChapterStringUIDID            = 0x5654
	ChapterTimeStartID            = 0x91
	ChapterTimeEndID              = 0x92
	ChapterFlagHiddenID           = 0x98
	ChapterFlagEnabledID          = 0x4598
	ChapterSegmentUUIDID          = 0x6e67
	ChapterSkipTypeID             = 0x4588
	ChapterSegmentEditionUIDID    = 0x6ebc
	ChapterPhysicalEquivID        = 0x63c3
	ChapterTrackID                = 0x8f
	ChapterTrackUIDID             = 0x89
	ChapterDisplayID              = 0x80
	ChapStringID                  = 0x85
	ChapLanguageID                = 0x437c
	ChapLanguageBCP47ID           = 0x437d
	ChapCountryID                 = 0x437e
	ChapProcessID                 = 0x6944
	ChapProcessCodecIDID          = 0x6955
	ChapProcessPrivateID          = 0x450d
	ChapProcessCommandID          = 0x6911
	ChapProcessTimeID             = 0x6922
	ChapProcessDataID             = 0x6933
	TagsID                        = 0x1254c367
	TagID                         = 0x7373
	TargetsID                     = 0x63c0
	TargetTypeValueID             = 0x68ca
	TargetTypeID                  = 0x63ca
	TagTrackUIDID                 = 0x63c5
	TagEditionUIDID               = 0x63c9
	TagChapterUIDID               = 0x63c4
	TagAttachmentUIDID            = 0x63c6
	SimpleTagID                   = 0x67c8
	TagNameID                     = 0x45a3
	TagLanguageID                 = 0x447a
	TagLanguageBCP47ID            = 0x447b
	TagDefaultID                  = 0x4484
	TagDefaultBogusID             = 0x44b4
	TagStringID                   = 0x4487
	TagBinaryID                   = 0x4485
)

var SegmentElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         SegmentID,
		ParentID:   RootID,
		Name:       "segment",
		Definition: "The Root Element that contains all other Top-Level Elements",
	},
	Master: map[ebml.ID]ebml.Element{
		SeekHeadID:    SeekHeadElement,
		InfoID:        InfoElement,
		ClusterID:     ClusterElement,
		TracksID:      TracksElement,
		CuesID:        CuesElement,
		AttachmentsID: AttachmentsElement,
		ChaptersID:    ChaptersElement,
		TagsID:        TagsElement,
	},
}

var SeekHeadElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         SeekHeadID,
		ParentID:   SegmentID,
		Name:       "seek_head",
		Definition: "Contains seeking information of Top-Level Elements",
	},
	Master: map[ebml.ID]ebml.Element{
		SeekID: SeekElement,
	},
}

var SeekElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         SeekID,
		ParentID:   SeekHeadID,
		Name:       "seek",
		Definition: "Contains a single seek entry to an EBML Element",
	},
	Master: map[ebml.ID]ebml.Element{
		SeekIDID:       SeekIDElement,
		SeekPositionID: SeekPositionElement,
	},
}
var SeekIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         SeekIDID,
		ParentID:   SeekID,
		Name:       "seek_id",
		Definition: "The binary EBML ID of a Top-Level Element",
	},
}
var SeekPositionElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         SeekPositionID,
		ParentID:   SeekID,
		Name:       "seek_position",
		Definition: "The Segment Position of a Top-Level Element",
	},
}

var InfoElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         InfoID,
		ParentID:   SegmentID,
		Name:       "info",
		Definition: "Contains general information about the Segment",
	},
	Master: map[ebml.ID]ebml.Element{
		SegmentUUIDID:      SegmentUUIDElement,
		SegmentFilenameID:  SegmentFilenameElement,
		PrevUUIDID:         PrevUUIDElement,
		PrevFilenameID:     PrevFilenameElement,
		NextUUIDID:         NextUUIDElement,
		NextFilenameID:     NextFilenameElement,
		SegmentFamilyID:    SegmentFamilyElement,
		ChapterTranslateID: ChapterTranslateElement,
		TimestampScaleID:   TimestampScaleElement,
		DurationID:         DurationElement,
		DateUTCID:          DateUTCElement,
		TitleID:            TitleElement,
		MuxingAppID:        MuxingAppElement,
		WritingAppID:       WritingAppElement,
	},
}
var SegmentUUIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         SegmentUUIDID,
		ParentID:   InfoID,
		Name:       "segment_uuid",
		Definition: "A randomly generated UID that identifies the Segment amongst many others v4 RFC9562 with all bits randomly (or pseudorandomly) chosen",
	},
}
var SegmentFilenameElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         SegmentFilenameID,
		ParentID:   InfoID,
		Name:       "segment_filename",
		Definition: "A filename corresponding to this Segment",
	},
}
var PrevUUIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         PrevUUIDID,
		ParentID:   InfoID,
		Name:       "prev_uuid",
		Definition: "An ID that identifies the previous Segment of a Linked Segment",
	},
}
var PrevFilenameElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         PrevFilenameID,
		ParentID:   InfoID,
		Name:       "prev_filename",
		Definition: "A filename corresponding to the file of the previous Linked Segment",
	},
}
var NextUUIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         NextUUIDID,
		ParentID:   InfoID,
		Name:       "next_uuid",
		Definition: "An ID that identifies the next Segment of a Linked Segment",
	},
}
var NextFilenameElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         NextFilenameID,
		ParentID:   InfoID,
		Name:       "next_filename",
		Definition: "A filename corresponding to the file of the next Linked Segment",
	},
}
var SegmentFamilyElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         SegmentFamilyID,
		ParentID:   InfoID,
		Name:       "segment_family",
		Definition: "A UID that all Segments of a Linked Segment **MUST** share chosen",
	},
}
var TimestampScaleElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TimestampScaleID,
		ParentID:   InfoID,
		Name:       "timestamp_scale",
		Definition: "Base unit for Segment Ticks and Track Ticks",
	},
}
var DurationElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         DurationID,
		ParentID:   InfoID,
		Name:       "duration",
		Definition: "Duration of the Segment",
	},
}
var DateUTCElement = &ebml.Date{
	ElementType: ebml.ElementType{
		ID:         DateUTCID,
		ParentID:   InfoID,
		Name:       "date_utc",
		Definition: "The date and time that the Segment was created by the muxing application or library",
	},
}
var TitleElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         TitleID,
		ParentID:   InfoID,
		Name:       "title",
		Definition: "General name of the Segment",
	},
}
var MuxingAppElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         MuxingAppID,
		ParentID:   InfoID,
		Name:       "muxing_app",
		Definition: "Muxing application or library",
	},
}
var WritingAppElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         WritingAppID,
		ParentID:   InfoID,
		Name:       "writing_app",
		Definition: "Writing application",
	},
}

var ChapterTranslateElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ChapterTranslateID,
		ParentID:   InfoID,
		Name:       "chapter_translate",
		Definition: "The mapping between this Segment and a segment value in the given Chapter Codec",
	},
	Master: map[ebml.ID]ebml.Element{
		ChapterTranslateIDID:         ChapterTranslateIDElement,
		ChapterTranslateCodecID:      ChapterTranslateCodecElement,
		ChapterTranslateEditionUIDID: ChapterTranslateEditionUIDElement,
	},
}
var ChapterTranslateIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ChapterTranslateIDID,
		ParentID:   ChapterTranslateID,
		Name:       "chapter_translate_id",
		Definition: "The binary value used to represent this Segment in the chapter codec data",
	},
}
var ChapterTranslateCodecElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterTranslateCodecID,
		ParentID:   ChapterTranslateID,
		Name:       "chapter_translate_codec",
		Definition: "Applies to the chapter codec of the given chapter edition",
	},
}
var ChapterTranslateEditionUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterTranslateEditionUIDID,
		ParentID:   ChapterTranslateID,
		Name:       "chapter_translate_edition_uid",
		Definition: "Specifies a chapter edition UID to which this ChapterTranslate applies",
	},
}

var ClusterElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ClusterID,
		ParentID:   SegmentID,
		Name:       "cluster",
		Definition: "The Top-Level Element containing the (monolithic) Block structure",
	},
	Master: map[ebml.ID]ebml.Element{
		TimestampID:      TimestampElement,
		SilentTracksID:   SilentTracksElement,
		PositionID:       PositionElement,
		PrevSizeID:       PrevSizeElement,
		SimpleBlockID:    SimpleBlockElement,
		BlockGroupID:     BlockGroupElement,
		EncryptedBlockID: EncryptedBlockElement,
	},
}
var TimestampElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TimestampID,
		ParentID:   ClusterID,
		Name:       "timestamp",
		Definition: "Absolute timestamp of the cluster",
	},
}
var PositionElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PositionID,
		ParentID:   ClusterID,
		Name:       "position",
		Definition: "The Segment Position of the Cluster in the Segment (0 in live streams)",
	},
}
var PrevSizeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PrevSizeID,
		ParentID:   ClusterID,
		Name:       "prev_size",
		Definition: "Size of the previous Cluster",
	},
}
var SimpleBlockElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         SimpleBlockID,
		ParentID:   ClusterID,
		Name:       "simple_block",
		Definition: "Similar to Block ) but without all the extra information",
	},
}
var EncryptedBlockElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         EncryptedBlockID,
		ParentID:   ClusterID,
		Name:       "encrypted_block",
		Definition: "Similar to SimpleBlock )",
	},
}

var SilentTracksElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         SilentTracksID,
		ParentID:   ClusterID,
		Name:       "silent_tracks",
		Definition: "The list of tracks that are not used in that part of the stream",
	},
	Master: map[ebml.ID]ebml.Element{
		SilentTrackNumberID: SilentTrackNumberElement,
	},
}
var SilentTrackNumberElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         SilentTrackNumberID,
		ParentID:   SilentTracksID,
		Name:       "silent_track_number",
		Definition: "One of the track numbers that is not used from now on in the stream",
	},
}

var BlockGroupElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         BlockGroupID,
		ParentID:   ClusterID,
		Name:       "block_group",
		Definition: "Basic container of information containing a single Block and information specific to that Block",
	},
	Master: map[ebml.ID]ebml.Element{
		BlockID:             BlockElement,
		BlockVirtualID:      BlockVirtualElement,
		BlockAdditionsID:    BlockAdditionsElement,
		BlockDurationID:     BlockDurationElement,
		ReferencePriorityID: ReferencePriorityElement,
		ReferenceBlockID:    ReferenceBlockElement,
		ReferenceVirtualID:  ReferenceVirtualElement,
		CodecStateID:        CodecStateElement,
		DiscardPaddingID:    DiscardPaddingElement,
		SlicesID:            SlicesElement,
		ReferenceFrameID:    ReferenceFrameElement,
	},
}
var BlockElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         BlockID,
		ParentID:   BlockGroupID,
		Name:       "block",
		Definition: "Block containing the actual data to be rendered and a timestamp relative to the Cluster Timestamp",
	},
}
var BlockVirtualElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         BlockVirtualID,
		ParentID:   BlockGroupID,
		Name:       "block_virtual",
		Definition: "A Block with no data",
	},
}
var BlockDurationElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         BlockDurationID,
		ParentID:   BlockGroupID,
		Name:       "block_duration",
		Definition: "The duration of the Block",
	},
}
var ReferencePriorityElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ReferencePriorityID,
		ParentID:   BlockGroupID,
		Name:       "reference_priority",
		Definition: "This frame is referenced and has the specified cache priority",
	},
}
var ReferenceBlockElement = &ebml.Integer{
	ElementType: ebml.ElementType{
		ID:         ReferenceBlockID,
		ParentID:   BlockGroupID,
		Name:       "reference_block",
		Definition: "A timestamp value",
	},
}
var ReferenceVirtualElement = &ebml.Integer{
	ElementType: ebml.ElementType{
		ID:         ReferenceVirtualID,
		ParentID:   BlockGroupID,
		Name:       "reference_virtual",
		Definition: "The Segment Position of the data that would otherwise be in position of the virtual block",
	},
}
var CodecStateElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         CodecStateID,
		ParentID:   BlockGroupID,
		Name:       "codec_state",
		Definition: "The new codec state to use",
	},
}
var DiscardPaddingElement = &ebml.Integer{
	ElementType: ebml.ElementType{
		ID:         DiscardPaddingID,
		ParentID:   BlockGroupID,
		Name:       "discard_padding",
		Definition: "Duration of the silent data added to the Block",
	},
}

var BlockAdditionsElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         BlockAdditionsID,
		ParentID:   BlockGroupID,
		Name:       "block_additions",
		Definition: "Contains additional binary data to complete the Block element",
	},
	Master: map[ebml.ID]ebml.Element{
		BlockMoreID: BlockMoreElement,
	},
}

var BlockMoreElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         BlockMoreID,
		ParentID:   BlockAdditionsID,
		Name:       "block_more",
		Definition: "Contains the BlockAdditional and some parameters",
	},
	Master: map[ebml.ID]ebml.Element{
		BlockAdditionalID: BlockAdditionalElement,
		BlockAddIDID:      BlockAddIDElement,
	},
}
var BlockAdditionalElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         BlockAdditionalID,
		ParentID:   BlockMoreID,
		Name:       "block_additional",
		Definition: "Interpreted by the codec as it wishes",
	},
}
var BlockAddIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         BlockAddIDID,
		ParentID:   BlockMoreID,
		Name:       "block_add_id",
		Definition: "An ID that identifies how to interpret the BlockAdditional data",
	},
}

var SlicesElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         SlicesID,
		ParentID:   BlockGroupID,
		Name:       "slices",
		Definition: "Contains slices description",
	},
	Master: map[ebml.ID]ebml.Element{
		TimeSliceID: TimeSliceElement,
	},
}

var TimeSliceElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TimeSliceID,
		ParentID:   SlicesID,
		Name:       "time_slice",
		Definition: "Contains extra time information about the data contained in the Block",
	},
	Master: map[ebml.ID]ebml.Element{
		LaceNumberID:      LaceNumberElement,
		FrameNumberID:     FrameNumberElement,
		BlockAdditionIDID: BlockAdditionIDElement,
		DelayID:           DelayElement,
		SliceDurationID:   SliceDurationElement,
	},
}
var LaceNumberElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         LaceNumberID,
		ParentID:   TimeSliceID,
		Name:       "lace_number",
		Definition: "The reverse number of the frame in the lace ",
	},
}
var FrameNumberElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FrameNumberID,
		ParentID:   TimeSliceID,
		Name:       "frame_number",
		Definition: "The number of the frame to generate from this lace with this delay",
	},
}
var BlockAdditionIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         BlockAdditionIDID,
		ParentID:   TimeSliceID,
		Name:       "block_addition_id",
		Definition: "The ID of the BlockAdditional element",
	},
}
var DelayElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         DelayID,
		ParentID:   TimeSliceID,
		Name:       "delay",
		Definition: "The delay to apply to the element",
	},
}
var SliceDurationElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         SliceDurationID,
		ParentID:   TimeSliceID,
		Name:       "slice_duration",
		Definition: "The duration to apply to the element",
	},
}

var ReferenceFrameElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ReferenceFrameID,
		ParentID:   BlockGroupID,
		Name:       "reference_frame",
		Definition: "Contains information about the last reference frame",
	},
	Master: map[ebml.ID]ebml.Element{
		ReferenceOffsetID:    ReferenceOffsetElement,
		ReferenceTimestampID: ReferenceTimestampElement,
	},
}
var ReferenceOffsetElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ReferenceOffsetID,
		ParentID:   ReferenceFrameID,
		Name:       "reference_offset",
		Definition: "The relative offset",
	},
}
var ReferenceTimestampElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ReferenceTimestampID,
		ParentID:   ReferenceFrameID,
		Name:       "reference_timestamp",
		Definition: "The timestamp of the BlockGroup pointed to by ReferenceOffset",
	},
}

var TracksElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TracksID,
		ParentID:   SegmentID,
		Name:       "tracks",
		Definition: "A Top-Level Element of information with many tracks described",
	},
	Master: map[ebml.ID]ebml.Element{
		TrackEntryID: TrackEntryElement,
	},
}

var TrackEntryElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TrackEntryID,
		ParentID:   TracksID,
		Name:       "track_entry",
		Definition: "Describes a track with all elements",
	},
	Master: map[ebml.ID]ebml.Element{
		TrackNumberID:                 TrackNumberElement,
		TrackUIDID:                    TrackUIDElement,
		TrackTypeID:                   TrackTypeElement,
		FlagEnabledID:                 FlagEnabledElement,
		FlagDefaultID:                 FlagDefaultElement,
		FlagForcedID:                  FlagForcedElement,
		FlagHearingImpairedID:         FlagHearingImpairedElement,
		FlagVisualImpairedID:          FlagVisualImpairedElement,
		FlagTextDescriptionsID:        FlagTextDescriptionsElement,
		FlagOriginalID:                FlagOriginalElement,
		FlagCommentaryID:              FlagCommentaryElement,
		FlagLacingID:                  FlagLacingElement,
		MinCacheID:                    MinCacheElement,
		MaxCacheID:                    MaxCacheElement,
		DefaultDurationID:             DefaultDurationElement,
		DefaultDecodedFieldDurationID: DefaultDecodedFieldDurationElement,
		TrackTimestampScaleID:         TrackTimestampScaleElement,
		TrackOffsetID:                 TrackOffsetElement,
		MaxBlockAdditionIDID:          MaxBlockAdditionIDElement,
		BlockAdditionMappingID:        BlockAdditionMappingElement,
		NameID:                        NameElement,
		LanguageID:                    LanguageElement,
		LanguageBCP47ID:               LanguageBCP47Element,
		CodecIDID:                     CodecIDElement,
		CodecPrivateID:                CodecPrivateElement,
		CodecNameID:                   CodecNameElement,
		AttachmentLinkID:              AttachmentLinkElement,
		CodecSettingsID:               CodecSettingsElement,
		CodecInfoURLID:                CodecInfoURLElement,
		CodecDownloadURLID:            CodecDownloadURLElement,
		CodecDecodeAllID:              CodecDecodeAllElement,
		TrackOverlayID:                TrackOverlayElement,
		CodecDelayID:                  CodecDelayElement,
		SeekPreRollID:                 SeekPreRollElement,
		TrackTranslateID:              TrackTranslateElement,
		VideoID:                       VideoElement,
		AudioID:                       AudioElement,
		TrackOperationID:              TrackOperationElement,
		TrickTrackUIDID:               TrickTrackUIDElement,
		TrickTrackSegmentUIDID:        TrickTrackSegmentUIDElement,
		TrickTrackFlagID:              TrickTrackFlagElement,
		TrickMasterTrackUIDID:         TrickMasterTrackUIDElement,
		TrickMasterTrackSegmentUIDID:  TrickMasterTrackSegmentUIDElement,
		ContentEncodingsID:            ContentEncodingsElement,
	},
}
var TrackNumberElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackNumberID,
		ParentID:   TrackEntryID,
		Name:       "track_number",
		Definition: "The track number as used in the Block Header",
	},
}
var TrackUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackUIDID,
		ParentID:   TrackEntryID,
		Name:       "track_uid",
		Definition: "A UID that identifies the Track",
	},
}
var TrackTypeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackTypeID,
		ParentID:   TrackEntryID,
		Name:       "track_type",
		Definition: "The TrackType defines the type of each frame found in the Track",
	},
	Enums: map[uint64]ebml.Enum{
		1:  {Name: "video", Description: "An image"},
		2:  {Name: "audio", Description: "Audio samples"},
		3:  {Name: "complex", Description: "A mix of different other TrackType"},
		16: {Name: "logo", Description: "An image to be rendered over the video track(s)"},
		17: {Name: "subtitle", Description: "Subtitle or closed caption data to be rendered over the video track(s)"},
		18: {Name: "buttons", Description: "Interactive button"},
		32: {Name: "control", Description: "Metadata used to control the player of the Matroska Player"},
		33: {Name: "metadata", Description: "Timed metadata that can be passed on to the Matroska Player"},
	},
}
var FlagEnabledElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagEnabledID,
		ParentID:   TrackEntryID,
		Name:       "flag_enabled",
		Definition: "Set to 1 if the track is usable",
	},
}
var FlagDefaultElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagDefaultID,
		ParentID:   TrackEntryID,
		Name:       "flag_default",
		Definition: "Set to 1 if the track is eligible for automatic selection by the player",
	},
}
var FlagForcedElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagForcedID,
		ParentID:   TrackEntryID,
		Name:       "flag_forced",
		Definition: "Applies only to subtitles",
	},
}
var FlagHearingImpairedElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagHearingImpairedID,
		ParentID:   TrackEntryID,
		Name:       "flag_hearing_impaired",
		Definition: "Set to 1 if and only if the track is suitable for users with hearing impairments",
	},
}
var FlagVisualImpairedElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagVisualImpairedID,
		ParentID:   TrackEntryID,
		Name:       "flag_visual_impaired",
		Definition: "Set to 1 if and only if the track is suitable for users with visual impairments",
	},
}
var FlagTextDescriptionsElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagTextDescriptionsID,
		ParentID:   TrackEntryID,
		Name:       "flag_text_descriptions",
		Definition: "Set to 1 if and only if the track contains textual descriptions of video content",
	},
}
var FlagOriginalElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagOriginalID,
		ParentID:   TrackEntryID,
		Name:       "flag_original",
		Definition: "Set to 1 if and only if the track is in the content's original language",
	},
}
var FlagCommentaryElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagCommentaryID,
		ParentID:   TrackEntryID,
		Name:       "flag_commentary",
		Definition: "Set to 1 if and only if the track contains commentary",
	},
}
var FlagLacingElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagLacingID,
		ParentID:   TrackEntryID,
		Name:       "flag_lacing",
		Definition: "Set to 1 if the track **MAY** contain blocks that use lacing",
	},
}
var MinCacheElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         MinCacheID,
		ParentID:   TrackEntryID,
		Name:       "min_cache",
		Definition: "The minimum number of frames a player should be able to cache during playback",
	},
}
var MaxCacheElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         MaxCacheID,
		ParentID:   TrackEntryID,
		Name:       "max_cache",
		Definition: "The maximum cache size necessary to store referenced frames in and the current frame",
	},
}
var DefaultDurationElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         DefaultDurationID,
		ParentID:   TrackEntryID,
		Name:       "default_duration",
		Definition: "Number of nanoseconds per frame",
	},
}
var DefaultDecodedFieldDurationElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         DefaultDecodedFieldDurationID,
		ParentID:   TrackEntryID,
		Name:       "default_decoded_field_duration",
		Definition: "The period between two successive fields at the output of the decoding process",
	},
}
var TrackTimestampScaleElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         TrackTimestampScaleID,
		ParentID:   TrackEntryID,
		Name:       "track_timestamp_scale",
		Definition: "The scale to apply on this track to work at normal speed in relation with other tracks",
	},
}
var TrackOffsetElement = &ebml.Integer{
	ElementType: ebml.ElementType{
		ID:         TrackOffsetID,
		ParentID:   TrackEntryID,
		Name:       "track_offset",
		Definition: "A value to add to the Block's Timestamp",
	},
}
var MaxBlockAdditionIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         MaxBlockAdditionIDID,
		ParentID:   TrackEntryID,
		Name:       "max_block_addition_id",
		Definition: "The maximum value of BlockAddID ",
	},
}
var NameElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         NameID,
		ParentID:   TrackEntryID,
		Name:       "name",
		Definition: "A human-readable track name",
	},
}
var LanguageElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         LanguageID,
		ParentID:   TrackEntryID,
		Name:       "language",
		Definition: "The language of the track",
	},
}
var LanguageBCP47Element = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         LanguageBCP47ID,
		ParentID:   TrackEntryID,
		Name:       "language_bcp47",
		Definition: "The language of the track",
	},
}
var CodecIDElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         CodecIDID,
		ParentID:   TrackEntryID,
		Name:       "codec_id",
		Definition: "An ID corresponding to the codec",
	},
}
var CodecPrivateElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         CodecPrivateID,
		ParentID:   TrackEntryID,
		Name:       "codec_private",
		Definition: "Private data only known to the codec",
	},
}
var CodecNameElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         CodecNameID,
		ParentID:   TrackEntryID,
		Name:       "codec_name",
		Definition: "A human-readable string specifying the codec",
	},
}
var AttachmentLinkElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         AttachmentLinkID,
		ParentID:   TrackEntryID,
		Name:       "attachment_link",
		Definition: "The UID of an attachment that is used by this codec",
	},
}
var CodecSettingsElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         CodecSettingsID,
		ParentID:   TrackEntryID,
		Name:       "codec_settings",
		Definition: "A string describing the encoding setting used",
	},
}
var CodecInfoURLElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         CodecInfoURLID,
		ParentID:   TrackEntryID,
		Name:       "codec_info_url",
		Definition: "A URL to find information about the codec used",
	},
}
var CodecDownloadURLElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         CodecDownloadURLID,
		ParentID:   TrackEntryID,
		Name:       "codec_download_url",
		Definition: "A URL to download information about the codec used",
	},
}
var CodecDecodeAllElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CodecDecodeAllID,
		ParentID:   TrackEntryID,
		Name:       "codec_decode_all",
		Definition: "Set to 1 if the codec can decode potentially damaged data",
	},
}
var TrackOverlayElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackOverlayID,
		ParentID:   TrackEntryID,
		Name:       "track_overlay",
		Definition: "Specify that this track is an overlay track for the Track specified (in the u-integer)",
	},
}
var CodecDelayElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CodecDelayID,
		ParentID:   TrackEntryID,
		Name:       "codec_delay",
		Definition: "The built-in delay for the codec",
	},
}
var SeekPreRollElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         SeekPreRollID,
		ParentID:   TrackEntryID,
		Name:       "seek_pre_roll",
		Definition: "After a discontinuity",
	},
}
var TrickTrackUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrickTrackUIDID,
		ParentID:   TrackEntryID,
		Name:       "trick_track_uid",
		Definition: "The TrackUID of the Smooth FF/RW video in the paired EBML structure corresponding to this video track",
	},
}
var TrickTrackSegmentUIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         TrickTrackSegmentUIDID,
		ParentID:   TrackEntryID,
		Name:       "trick_track_segment_uid",
		Definition: "The SegmentUUID of the Segment containing the track identified by TrickTrackUID",
	},
}
var TrickTrackFlagElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrickTrackFlagID,
		ParentID:   TrackEntryID,
		Name:       "trick_track_flag",
		Definition: "Set to 1 if this video track is a Smooth FF/RW track",
	},
}
var TrickMasterTrackUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrickMasterTrackUIDID,
		ParentID:   TrackEntryID,
		Name:       "trick_master_track_uid",
		Definition: "The TrackUID of the video track in the paired EBML structure that corresponds to this Smooth FF/RW track",
	},
}
var TrickMasterTrackSegmentUIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         TrickMasterTrackSegmentUIDID,
		ParentID:   TrackEntryID,
		Name:       "trick_master_track_segment_uid",
		Definition: "The SegmentUUID of the Segment containing the track identified by MasterTrackUID",
	},
}

var BlockAdditionMappingElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         BlockAdditionMappingID,
		ParentID:   TrackEntryID,
		Name:       "block_addition_mapping",
		Definition: "Contains elements that extend the track format by adding content either to each frame",
	},
	Master: map[ebml.ID]ebml.Element{
		BlockAddIDValueID:     BlockAddIDValueElement,
		BlockAddIDNameID:      BlockAddIDNameElement,
		BlockAddIDTypeID:      BlockAddIDTypeElement,
		BlockAddIDExtraDataID: BlockAddIDExtraDataElement,
	},
}
var BlockAddIDValueElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         BlockAddIDValueID,
		ParentID:   BlockAdditionMappingID,
		Name:       "block_add_idvalue",
		Definition: "If the track format extension needs content beside frames",
	},
}
var BlockAddIDNameElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         BlockAddIDNameID,
		ParentID:   BlockAdditionMappingID,
		Name:       "block_add_idname",
		Definition: "A human-friendly name describing the type of BlockAdditional data",
	},
}
var BlockAddIDTypeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         BlockAddIDTypeID,
		ParentID:   BlockAdditionMappingID,
		Name:       "block_add_idtype",
		Definition: "Stores the registered identifier of the Block Additional Mapping to define how the BlockAdditional data should be handled",
	},
}
var BlockAddIDExtraDataElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         BlockAddIDExtraDataID,
		ParentID:   BlockAdditionMappingID,
		Name:       "block_add_idextra_data",
		Definition: "Extra binary data that the BlockAddIDType can use to interpret the BlockAdditional data",
	},
}

var TrackTranslateElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TrackTranslateID,
		ParentID:   TrackEntryID,
		Name:       "track_translate",
		Definition: "The mapping between this TrackEntry and a track value in the given Chapter Codec",
	},
	Master: map[ebml.ID]ebml.Element{
		TrackTranslateTrackIDID:    TrackTranslateTrackIDElement,
		TrackTranslateCodecID:      TrackTranslateCodecElement,
		TrackTranslateEditionUIDID: TrackTranslateEditionUIDElement,
	},
}
var TrackTranslateTrackIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         TrackTranslateTrackIDID,
		ParentID:   TrackTranslateID,
		Name:       "track_translate_track_id",
		Definition: "The binary value used to represent this TrackEntry in the chapter codec data",
	},
}
var TrackTranslateCodecElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackTranslateCodecID,
		ParentID:   TrackTranslateID,
		Name:       "track_translate_codec",
		Definition: "Applies to the chapter codec of the given chapter edition",
	},
}
var TrackTranslateEditionUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackTranslateEditionUIDID,
		ParentID:   TrackTranslateID,
		Name:       "track_translate_edition_uid",
		Definition: "Specifies a chapter edition UID to which this TrackTranslate applies",
	},
}

var VideoElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         VideoID,
		ParentID:   TrackEntryID,
		Name:       "video",
		Definition: "Video settings",
	},
	Master: map[ebml.ID]ebml.Element{
		FlagInterlacedID:     FlagInterlacedElement,
		FieldOrderID:         FieldOrderElement,
		StereoModeID:         StereoModeElement,
		AlphaModeID:          AlphaModeElement,
		OldStereoModeID:      OldStereoModeElement,
		PixelWidthID:         PixelWidthElement,
		PixelHeightID:        PixelHeightElement,
		PixelCropBottomID:    PixelCropBottomElement,
		PixelCropTopID:       PixelCropTopElement,
		PixelCropLeftID:      PixelCropLeftElement,
		PixelCropRightID:     PixelCropRightElement,
		DisplayWidthID:       DisplayWidthElement,
		DisplayHeightID:      DisplayHeightElement,
		DisplayUnitID:        DisplayUnitElement,
		AspectRatioTypeID:    AspectRatioTypeElement,
		UncompressedFourCCID: UncompressedFourCCElement,
		GammaValueID:         GammaValueElement,
		FrameRateID:          FrameRateElement,
		ColourID:             ColourElement,
		ProjectionID:         ProjectionElement,
	},
}
var FlagInterlacedElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FlagInterlacedID,
		ParentID:   VideoID,
		Name:       "flag_interlaced",
		Definition: "Specifies whether the video frames in this track are interlaced",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "undetermined", Description: "Unknown status"},
		1: {Name: "interlaced", Description: "Interlaced frames"},
		2: {Name: "progressive", Description: "No interlacing"},
	},
}
var FieldOrderElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FieldOrderID,
		ParentID:   VideoID,
		Name:       "field_order",
		Definition: "Specifies the field ordering of video frames in this track",
	},
	Enums: map[uint64]ebml.Enum{
		0:  {Name: "progressive", Description: "Interlaced frames"},
		1:  {Name: "tff", Description: "Top field displayed first"},
		2:  {Name: "undetermined", Description: "Unknown field order"},
		6:  {Name: "bff", Description: "Bottom field displayed first"},
		9:  {Description: "Top field displayed first"},
		14: {Description: "Bottom field displayed first"},
	},
}
var StereoModeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         StereoModeID,
		ParentID:   VideoID,
		Name:       "stereo_mode",
		Definition: "Stereo-3D video mode",
	},
	Enums: map[uint64]ebml.Enum{
		0:  {Name: "mono"},
		1:  {Description: "side by side (left eye first)"},
		2:  {Description: "top - bottom (right eye is first)"},
		3:  {Description: "top - bottom (left eye is first)"},
		4:  {Description: "checkboard (right eye is first)"},
		5:  {Description: "checkboard (left eye is first)"},
		6:  {Description: "row interleaved (right eye is first)"},
		7:  {Description: "row interleaved (left eye is first)"},
		8:  {Description: "column interleaved (right eye is first)"},
		9:  {Description: "column interleaved (left eye is first)"},
		10: {Description: "anaglyph (cyan/red)"},
		11: {Description: "side by side (right eye first)"},
		12: {Description: "anaglyph (green/magenta)"},
		13: {Description: "both eyes laced in one Block (left eye is first)"},
		14: {Description: "both eyes laced in one Block (right eye is first)"},
	},
}
var AlphaModeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         AlphaModeID,
		ParentID:   VideoID,
		Name:       "alpha_mode",
		Definition: "Indicates whether the BlockAdditional element with BlockAddID of \"1\" contains Alpha data as defined by the Codec Mapping for the CodecID",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "none", Description: "The BlockAdditional element with BlockAddID of \"1\" does not exist or **SHOULD NOT** be considered as containing such data"},
		1: {Name: "present", Description: "The BlockAdditional element with BlockAddID of \"1\" contains alpha channel data"},
	},
}
var OldStereoModeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         OldStereoModeID,
		ParentID:   VideoID,
		Name:       "old_stereo_mode",
		Definition: "Bogus StereoMode value used in old versions of libmatroska",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "mono"},
		1: {Name: "right_eye"},
		2: {Name: "left_eye"},
		3: {Name: "both_eyes"},
	},
}
var PixelWidthElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PixelWidthID,
		ParentID:   VideoID,
		Name:       "pixel_width",
		Definition: "Width of the encoded video frames in pixels",
	},
}
var PixelHeightElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PixelHeightID,
		ParentID:   VideoID,
		Name:       "pixel_height",
		Definition: "Height of the encoded video frames in pixels",
	},
}
var PixelCropBottomElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PixelCropBottomID,
		ParentID:   VideoID,
		Name:       "pixel_crop_bottom",
		Definition: "The number of video pixels to remove at the bottom of the image",
	},
}
var PixelCropTopElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PixelCropTopID,
		ParentID:   VideoID,
		Name:       "pixel_crop_top",
		Definition: "The number of video pixels to remove at the top of the image",
	},
}
var PixelCropLeftElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PixelCropLeftID,
		ParentID:   VideoID,
		Name:       "pixel_crop_left",
		Definition: "The number of video pixels to remove on the left of the image",
	},
}
var PixelCropRightElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PixelCropRightID,
		ParentID:   VideoID,
		Name:       "pixel_crop_right",
		Definition: "The number of video pixels to remove on the right of the image",
	},
}
var DisplayWidthElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         DisplayWidthID,
		ParentID:   VideoID,
		Name:       "display_width",
		Definition: "Width of the video frames to display",
	},
}
var DisplayHeightElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         DisplayHeightID,
		ParentID:   VideoID,
		Name:       "display_height",
		Definition: "Height of the video frames to display",
	},
}
var DisplayUnitElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         DisplayUnitID,
		ParentID:   VideoID,
		Name:       "display_unit",
		Definition: "How DisplayWidth and DisplayHeight are interpreted",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "pixels"},
		1: {Name: "centimeters"},
		2: {Name: "inches"},
		3: {Name: "display_aspect_ratio"},
		4: {Name: "unknown"},
	},
}
var AspectRatioTypeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         AspectRatioTypeID,
		ParentID:   VideoID,
		Name:       "aspect_ratio_type",
		Definition: "Specifies the possible modifications to the aspect ratio",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "free_resizing"},
		1: {Name: "keep_aspect_ratio"},
		2: {Name: "fixed"},
	},
}
var UncompressedFourCCElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         UncompressedFourCCID,
		ParentID:   VideoID,
		Name:       "uncompressed_four_cc",
		Definition: "Specifies the uncompressed pixel format used for the Track's data as a FourCC",
	},
}
var GammaValueElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         GammaValueID,
		ParentID:   VideoID,
		Name:       "gamma_value",
		Definition: "Gamma value",
	},
}
var FrameRateElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         FrameRateID,
		ParentID:   VideoID,
		Name:       "frame_rate",
		Definition: "Number of frames per second",
	},
}

var ColourElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ColourID,
		ParentID:   VideoID,
		Name:       "colour",
		Definition: "Settings describing the color format",
	},
	Master: map[ebml.ID]ebml.Element{
		MatrixCoefficientsID:      MatrixCoefficientsElement,
		BitsPerChannelID:          BitsPerChannelElement,
		ChromaSubsamplingHorzID:   ChromaSubsamplingHorzElement,
		ChromaSubsamplingVertID:   ChromaSubsamplingVertElement,
		CbSubsamplingHorzID:       CbSubsamplingHorzElement,
		CbSubsamplingVertID:       CbSubsamplingVertElement,
		ChromaSitingHorzID:        ChromaSitingHorzElement,
		ChromaSitingVertID:        ChromaSitingVertElement,
		RangeID:                   RangeElement,
		TransferCharacteristicsID: TransferCharacteristicsElement,
		PrimariesID:               PrimariesElement,
		MaxCLLID:                  MaxCLLElement,
		MaxFALLID:                 MaxFALLElement,
		MasteringMetadataID:       MasteringMetadataElement,
	},
}
var MatrixCoefficientsElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         MatrixCoefficientsID,
		ParentID:   ColourID,
		Name:       "matrix_coefficients",
		Definition: "The Matrix Coefficients of the video used to derive luma and chroma values from red",
	},
	Enums: map[uint64]ebml.Enum{
		0:  {Name: "identity"},
		1:  {Name: "itu_r_bt_709"},
		2:  {Name: "unspecified"},
		3:  {Name: "reserved"},
		4:  {Name: "us_fcc_73_682"},
		5:  {Name: "itu_r_bt_470bg"},
		6:  {Name: "smpte_170m"},
		7:  {Name: "smpte_240m"},
		8:  {Name: "ycocg"},
		9:  {Name: "bt2020_non_constant_luminance"},
		10: {Name: "bt2020_constant_luminance"},
		11: {Name: "smpte_st_2085"},
		12: {Name: "chroma_derived_non_constant_luminance"},
		13: {Name: "chroma_derived_constant_luminance"},
		14: {Name: "itu_r_bt_2100_0"},
	},
}
var BitsPerChannelElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         BitsPerChannelID,
		ParentID:   ColourID,
		Name:       "bits_per_channel",
		Definition: "Number of decoded bits per channel",
	},
}
var ChromaSubsamplingHorzElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChromaSubsamplingHorzID,
		ParentID:   ColourID,
		Name:       "chroma_subsampling_horz",
		Definition: "The number of pixels to remove in the Cr and Cb channels for every pixel not removed horizontally",
	},
}
var ChromaSubsamplingVertElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChromaSubsamplingVertID,
		ParentID:   ColourID,
		Name:       "chroma_subsampling_vert",
		Definition: "The number of pixels to remove in the Cr and Cb channels for every pixel not removed vertically",
	},
}
var CbSubsamplingHorzElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CbSubsamplingHorzID,
		ParentID:   ColourID,
		Name:       "cb_subsampling_horz",
		Definition: "The number of pixels to remove in the Cb channel for every pixel not removed horizontally",
	},
}
var CbSubsamplingVertElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CbSubsamplingVertID,
		ParentID:   ColourID,
		Name:       "cb_subsampling_vert",
		Definition: "The number of pixels to remove in the Cb channel for every pixel not removed vertically",
	},
}
var ChromaSitingHorzElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChromaSitingHorzID,
		ParentID:   ColourID,
		Name:       "chroma_siting_horz",
		Definition: "How chroma is subsampled horizontally",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "unspecified"},
		1: {Name: "left_collocated"},
		2: {Name: "half"},
	},
}
var ChromaSitingVertElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChromaSitingVertID,
		ParentID:   ColourID,
		Name:       "chroma_siting_vert",
		Definition: "How chroma is subsampled vertically",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "unspecified"},
		1: {Name: "top_collocated"},
		2: {Name: "half"},
	},
}
var RangeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         RangeID,
		ParentID:   ColourID,
		Name:       "range",
		Definition: "Clipping of the color ranges",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "unspecified"},
		1: {Name: "broadcast_range"},
		2: {Description: "full range (no clipping)"},
		3: {Name: "defined_by_matrixcoefficients_transfercharacteristics"},
	},
}
var TransferCharacteristicsElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TransferCharacteristicsID,
		ParentID:   ColourID,
		Name:       "transfer_characteristics",
		Definition: "The transfer characteristics of the video",
	},
	Enums: map[uint64]ebml.Enum{
		0:  {Name: "reserved"},
		1:  {Name: "itu_r_bt_709"},
		2:  {Name: "unspecified"},
		3:  {Name: "reserved2"},
		4:  {Name: "gamma_2_2_curve_bt_470m"},
		5:  {Name: "gamma_2_8_curve_bt_470bg"},
		6:  {Name: "smpte_170m"},
		7:  {Name: "smpte_240m"},
		8:  {Name: "linear"},
		9:  {Name: "log"},
		10: {Name: "log_sqrt"},
		11: {Name: "iec_61966_2_4"},
		12: {Name: "itu_r_bt_1361_extended_colour_gamut"},
		13: {Name: "iec_61966_2_1"},
		14: {Name: "itu_r_bt_2020_10_bit"},
		15: {Name: "itu_r_bt_2020_12_bit"},
		16: {Name: "itu_r_bt_2100_perceptual_quantization"},
		17: {Name: "smpte_st_428_1"},
		18: {Description: "ARIB STD-B67 (HLG)"},
	},
}
var PrimariesElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         PrimariesID,
		ParentID:   ColourID,
		Name:       "primaries",
		Definition: "The color primaries of the video",
	},
	Enums: map[uint64]ebml.Enum{
		0:  {Name: "reserved"},
		1:  {Name: "itu_r_bt_709"},
		2:  {Name: "unspecified"},
		3:  {Name: "reserved2"},
		4:  {Name: "itu_r_bt_470m"},
		5:  {Name: "itu_r_bt_470bg_bt_601_625"},
		6:  {Name: "itu_r_bt_601_525_smpte_170m"},
		7:  {Name: "smpte_240m"},
		8:  {Name: "film"},
		9:  {Name: "itu_r_bt_2020"},
		10: {Name: "smpte_st_428_1"},
		11: {Name: "smpte_rp_432_2"},
		12: {Name: "smpte_eg_432_2"},
		22: {Name: "ebu_tech_3213_e_jedec_p22_phosphors"},
	},
}
var MaxCLLElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         MaxCLLID,
		ParentID:   ColourID,
		Name:       "max_cll",
		Definition: "Maximum brightness of a single pixel in candelas per square meter (cd/m^2^)",
	},
}
var MaxFALLElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         MaxFALLID,
		ParentID:   ColourID,
		Name:       "max_fall",
		Definition: "Maximum brightness of a single full frame in candelas per square meter (cd/m^2^)",
	},
}

var MasteringMetadataElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         MasteringMetadataID,
		ParentID:   ColourID,
		Name:       "mastering_metadata",
		Definition: "SMPTE 2086 mastering data",
	},
	Master: map[ebml.ID]ebml.Element{
		PrimaryRChromaticityXID:   PrimaryRChromaticityXElement,
		PrimaryRChromaticityYID:   PrimaryRChromaticityYElement,
		PrimaryGChromaticityXID:   PrimaryGChromaticityXElement,
		PrimaryGChromaticityYID:   PrimaryGChromaticityYElement,
		PrimaryBChromaticityXID:   PrimaryBChromaticityXElement,
		PrimaryBChromaticityYID:   PrimaryBChromaticityYElement,
		WhitePointChromaticityXID: WhitePointChromaticityXElement,
		WhitePointChromaticityYID: WhitePointChromaticityYElement,
		LuminanceMaxID:            LuminanceMaxElement,
		LuminanceMinID:            LuminanceMinElement,
	},
}
var PrimaryRChromaticityXElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         PrimaryRChromaticityXID,
		ParentID:   MasteringMetadataID,
		Name:       "primary_rchromaticity_x",
		Definition: "Red X chromaticity coordinate",
	},
}
var PrimaryRChromaticityYElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         PrimaryRChromaticityYID,
		ParentID:   MasteringMetadataID,
		Name:       "primary_rchromaticity_y",
		Definition: "Red Y chromaticity coordinate",
	},
}
var PrimaryGChromaticityXElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         PrimaryGChromaticityXID,
		ParentID:   MasteringMetadataID,
		Name:       "primary_gchromaticity_x",
		Definition: "Green X chromaticity coordinate",
	},
}
var PrimaryGChromaticityYElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         PrimaryGChromaticityYID,
		ParentID:   MasteringMetadataID,
		Name:       "primary_gchromaticity_y",
		Definition: "Green Y chromaticity coordinate",
	},
}
var PrimaryBChromaticityXElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         PrimaryBChromaticityXID,
		ParentID:   MasteringMetadataID,
		Name:       "primary_bchromaticity_x",
		Definition: "Blue X chromaticity coordinate",
	},
}
var PrimaryBChromaticityYElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         PrimaryBChromaticityYID,
		ParentID:   MasteringMetadataID,
		Name:       "primary_bchromaticity_y",
		Definition: "Blue Y chromaticity coordinate",
	},
}
var WhitePointChromaticityXElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         WhitePointChromaticityXID,
		ParentID:   MasteringMetadataID,
		Name:       "white_point_chromaticity_x",
		Definition: "White X chromaticity coordinate",
	},
}
var WhitePointChromaticityYElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         WhitePointChromaticityYID,
		ParentID:   MasteringMetadataID,
		Name:       "white_point_chromaticity_y",
		Definition: "White Y chromaticity coordinate",
	},
}
var LuminanceMaxElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         LuminanceMaxID,
		ParentID:   MasteringMetadataID,
		Name:       "luminance_max",
		Definition: "Maximum luminance",
	},
}
var LuminanceMinElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         LuminanceMinID,
		ParentID:   MasteringMetadataID,
		Name:       "luminance_min",
		Definition: "Minimum luminance",
	},
}

var ProjectionElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ProjectionID,
		ParentID:   VideoID,
		Name:       "projection",
		Definition: "Describes the video projection details",
	},
	Master: map[ebml.ID]ebml.Element{
		ProjectionTypeID:      ProjectionTypeElement,
		ProjectionPrivateID:   ProjectionPrivateElement,
		ProjectionPoseYawID:   ProjectionPoseYawElement,
		ProjectionPosePitchID: ProjectionPosePitchElement,
		ProjectionPoseRollID:  ProjectionPoseRollElement,
	},
}
var ProjectionTypeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ProjectionTypeID,
		ParentID:   ProjectionID,
		Name:       "projection_type",
		Definition: "Describes the projection used for this video track",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "rectangular"},
		1: {Name: "equirectangular"},
		2: {Name: "cubemap"},
		3: {Name: "mesh"},
	},
}
var ProjectionPrivateElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ProjectionPrivateID,
		ParentID:   ProjectionID,
		Name:       "projection_private",
		Definition: "Private data that only applies to a specific projection",
	},
}
var ProjectionPoseYawElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         ProjectionPoseYawID,
		ParentID:   ProjectionID,
		Name:       "projection_pose_yaw",
		Definition: "Specifies a yaw rotation to the projection",
	},
}
var ProjectionPosePitchElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         ProjectionPosePitchID,
		ParentID:   ProjectionID,
		Name:       "projection_pose_pitch",
		Definition: "Specifies a pitch rotation to the projection",
	},
}
var ProjectionPoseRollElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         ProjectionPoseRollID,
		ParentID:   ProjectionID,
		Name:       "projection_pose_roll",
		Definition: "Specifies a roll rotation to the projection",
	},
}

var AudioElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         AudioID,
		ParentID:   TrackEntryID,
		Name:       "audio",
		Definition: "Audio settings",
	},
	Master: map[ebml.ID]ebml.Element{
		SamplingFrequencyID:       SamplingFrequencyElement,
		OutputSamplingFrequencyID: OutputSamplingFrequencyElement,
		ChannelsID:                ChannelsElement,
		ChannelPositionsID:        ChannelPositionsElement,
		BitDepthID:                BitDepthElement,
		EmphasisID:                EmphasisElement,
	},
}
var SamplingFrequencyElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         SamplingFrequencyID,
		ParentID:   AudioID,
		Name:       "sampling_frequency",
		Definition: "Sampling frequency in Hz",
	},
}
var OutputSamplingFrequencyElement = &ebml.Float{
	ElementType: ebml.ElementType{
		ID:         OutputSamplingFrequencyID,
		ParentID:   AudioID,
		Name:       "output_sampling_frequency",
		Definition: "Real output sampling frequency in Hz that is used for Spectral Band Replication (SBR) techniques",
	},
}
var ChannelsElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChannelsID,
		ParentID:   AudioID,
		Name:       "channels",
		Definition: "Numbers of channels in the track",
	},
}
var ChannelPositionsElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ChannelPositionsID,
		ParentID:   AudioID,
		Name:       "channel_positions",
		Definition: "Table of horizontal angles for each successive channel",
	},
}
var BitDepthElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         BitDepthID,
		ParentID:   AudioID,
		Name:       "bit_depth",
		Definition: "Bits per sample",
	},
}
var EmphasisElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         EmphasisID,
		ParentID:   AudioID,
		Name:       "emphasis",
		Definition: "Audio emphasis applied on audio samples",
	},
	Enums: map[uint64]ebml.Enum{
		0:  {Name: "no_emphasis"},
		1:  {Name: "cd_audio", Description: "First order filter with zero point at 50 microseconds and a pole at 15 microseconds"},
		2:  {Name: "reserved"},
		3:  {Name: "ccit_j_17", Description: "Defined in ITU-J"},
		4:  {Name: "fm_50", Description: "FM Radio in Europe"},
		5:  {Name: "fm_75", Description: "FM Radio in the USA"},
		10: {Name: "phono_riaa", Description: "Phono filter with time constants of t1=3180"},
		11: {Name: "phono_iec_n78", Description: "Phono filter with time constants of t1=3180"},
		12: {Name: "phono_teldec", Description: "Phono filter with time constants of t1=3180"},
		13: {Name: "phono_emi", Description: "Phono filter with time constants of t1=2500"},
		14: {Name: "phono_columbia_lp", Description: "Phono filter with time constants of t1=1590"},
		15: {Name: "phono_london", Description: "Phono filter with time constants of t1=1590"},
		16: {Name: "phono_nartb", Description: "Phono filter with time constants of t1=3180"},
	},
}

var TrackOperationElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TrackOperationID,
		ParentID:   TrackEntryID,
		Name:       "track_operation",
		Definition: "Operation that needs to be applied on tracks to create this virtual track",
	},
	Master: map[ebml.ID]ebml.Element{
		TrackCombinePlanesID: TrackCombinePlanesElement,
		TrackJoinBlocksID:    TrackJoinBlocksElement,
	},
}

var TrackCombinePlanesElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TrackCombinePlanesID,
		ParentID:   TrackOperationID,
		Name:       "track_combine_planes",
		Definition: "Contains the list of all video plane tracks that need to be combined to create this 3D track",
	},
	Master: map[ebml.ID]ebml.Element{
		TrackPlaneID: TrackPlaneElement,
	},
}

var TrackPlaneElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TrackPlaneID,
		ParentID:   TrackCombinePlanesID,
		Name:       "track_plane",
		Definition: "Contains a video plane track that needs to be combined to create this 3D track",
	},
	Master: map[ebml.ID]ebml.Element{
		TrackPlaneUIDID:  TrackPlaneUIDElement,
		TrackPlaneTypeID: TrackPlaneTypeElement,
	},
}
var TrackPlaneUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackPlaneUIDID,
		ParentID:   TrackPlaneID,
		Name:       "track_plane_uid",
		Definition: "The TrackUID number of the track representing the plane",
	},
}
var TrackPlaneTypeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackPlaneTypeID,
		ParentID:   TrackPlaneID,
		Name:       "track_plane_type",
		Definition: "The kind of plane this track corresponds to",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "left_eye"},
		1: {Name: "right_eye"},
		2: {Name: "background"},
	},
}

var TrackJoinBlocksElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TrackJoinBlocksID,
		ParentID:   TrackOperationID,
		Name:       "track_join_blocks",
		Definition: "Contains the list of all tracks whose Blocks need to be combined to create this virtual track",
	},
	Master: map[ebml.ID]ebml.Element{
		TrackJoinUIDID: TrackJoinUIDElement,
	},
}
var TrackJoinUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TrackJoinUIDID,
		ParentID:   TrackJoinBlocksID,
		Name:       "track_join_uid",
		Definition: "The TrackUID number of a track whose blocks are used to create this virtual track",
	},
}

var ContentEncodingsElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ContentEncodingsID,
		ParentID:   TrackEntryID,
		Name:       "content_encodings",
		Definition: "Settings for several content encoding mechanisms like compression or encryption",
	},
	Master: map[ebml.ID]ebml.Element{
		ContentEncodingID: ContentEncodingElement,
	},
}

var ContentEncodingElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ContentEncodingID,
		ParentID:   ContentEncodingsID,
		Name:       "content_encoding",
		Definition: "Settings for one content encoding like compression or encryption",
	},
	Master: map[ebml.ID]ebml.Element{
		ContentEncodingOrderID: ContentEncodingOrderElement,
		ContentEncodingScopeID: ContentEncodingScopeElement,
		ContentEncodingTypeID:  ContentEncodingTypeElement,
		ContentCompressionID:   ContentCompressionElement,
		ContentEncryptionID:    ContentEncryptionElement,
	},
}
var ContentEncodingOrderElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ContentEncodingOrderID,
		ParentID:   ContentEncodingID,
		Name:       "content_encoding_order",
		Definition: "Defines the order to apply each ContentEncoding of the ContentEncodings",
	},
}
var ContentEncodingScopeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ContentEncodingScopeID,
		ParentID:   ContentEncodingID,
		Name:       "content_encoding_scope",
		Definition: "A bit field that describes which elements have been modified in this way",
	},
	Enums: map[uint64]ebml.Enum{
		1: {Name: "block", Description: "All frame contents"},
		2: {Name: "private", Description: "The track's CodecPrivate data"},
		4: {Name: "next", Description: "The next ContentEncoding"},
	},
}
var ContentEncodingTypeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ContentEncodingTypeID,
		ParentID:   ContentEncodingID,
		Name:       "content_encoding_type",
		Definition: "A value describing the kind of transformation that is applied",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "compression"},
		1: {Name: "encryption"},
	},
}

var ContentCompressionElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ContentCompressionID,
		ParentID:   ContentEncodingID,
		Name:       "content_compression",
		Definition: "Settings describing the compression used",
	},
	Master: map[ebml.ID]ebml.Element{
		ContentCompAlgoID:     ContentCompAlgoElement,
		ContentCompSettingsID: ContentCompSettingsElement,
	},
}
var ContentCompAlgoElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ContentCompAlgoID,
		ParentID:   ContentCompressionID,
		Name:       "content_comp_algo",
		Definition: "The compression algorithm used",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "zlib", Description: "zlib compression RFC1950"},
		1: {Name: "bzlib", Description: "bzip2 compression BZIP2 **SHOULD NOT** be used"},
		2: {Name: "lzo1x", Description: "Lempel-Ziv-Oberhumer compression LZO **SHOULD NOT** be used"},
		3: {Name: "header_stripping", Description: "Octets in ContentCompSettings have been stripped from each frame"},
	},
}
var ContentCompSettingsElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ContentCompSettingsID,
		ParentID:   ContentCompressionID,
		Name:       "content_comp_settings",
		Definition: "Settings that might be needed by the decompressor",
	},
}

var ContentEncryptionElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ContentEncryptionID,
		ParentID:   ContentEncodingID,
		Name:       "content_encryption",
		Definition: "Settings describing the encryption used",
	},
	Master: map[ebml.ID]ebml.Element{
		ContentEncAlgoID:        ContentEncAlgoElement,
		ContentEncKeyIDID:       ContentEncKeyIDElement,
		ContentEncAESSettingsID: ContentEncAESSettingsElement,
		ContentSignatureID:      ContentSignatureElement,
		ContentSigKeyIDID:       ContentSigKeyIDElement,
		ContentSigAlgoID:        ContentSigAlgoElement,
		ContentSigHashAlgoID:    ContentSigHashAlgoElement,
	},
}
var ContentEncAlgoElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ContentEncAlgoID,
		ParentID:   ContentEncryptionID,
		Name:       "content_enc_algo",
		Definition: "The encryption algorithm used",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "not_encrypted", Description: "The data are not encrypted"},
		1: {Name: "des", Description: "Data Encryption Standard (DES) FIPS46-3"},
		2: {Name: "3des", Description: "Triple Data Encryption Algorithm SP800-67"},
		3: {Name: "twofish", Description: "Twofish Encryption Algorithm Twofish"},
		4: {Name: "blowfish", Description: "Blowfish Encryption Algorithm Blowfish"},
		5: {Name: "aes", Description: "Advanced Encryption Standard (AES) FIPS197"},
	},
}
var ContentEncKeyIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ContentEncKeyIDID,
		ParentID:   ContentEncryptionID,
		Name:       "content_enc_key_id",
		Definition: "For public key algorithms",
	},
}
var ContentSignatureElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ContentSignatureID,
		ParentID:   ContentEncryptionID,
		Name:       "content_signature",
		Definition: "A cryptographic signature of the contents",
	},
}
var ContentSigKeyIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ContentSigKeyIDID,
		ParentID:   ContentEncryptionID,
		Name:       "content_sig_key_id",
		Definition: "This is the ID of the private key that the data was signed with",
	},
}
var ContentSigAlgoElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ContentSigAlgoID,
		ParentID:   ContentEncryptionID,
		Name:       "content_sig_algo",
		Definition: "The algorithm used for the signature",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "not_signed"},
		1: {Name: "rsa"},
	},
}
var ContentSigHashAlgoElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ContentSigHashAlgoID,
		ParentID:   ContentEncryptionID,
		Name:       "content_sig_hash_algo",
		Definition: "The hash algorithm used for the signature",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "not_signed"},
		1: {Name: "sha1_160"},
		2: {Name: "md5"},
	},
}

var ContentEncAESSettingsElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ContentEncAESSettingsID,
		ParentID:   ContentEncryptionID,
		Name:       "content_enc_aessettings",
		Definition: "Settings describing the encryption algorithm used",
	},
	Master: map[ebml.ID]ebml.Element{
		AESSettingsCipherModeID: AESSettingsCipherModeElement,
	},
}
var AESSettingsCipherModeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         AESSettingsCipherModeID,
		ParentID:   ContentEncAESSettingsID,
		Name:       "aessettings_cipher_mode",
		Definition: "The AES cipher mode used in the encryption",
	},
	Enums: map[uint64]ebml.Enum{
		1: {Name: "aes_ctr", Description: "Counter SP800-38A"},
		2: {Name: "aes_cbc", Description: "Cipher Block Chaining SP800-38A"},
	},
}

var CuesElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         CuesID,
		ParentID:   SegmentID,
		Name:       "cues",
		Definition: "A Top-Level Element to speed seeking access",
	},
	Master: map[ebml.ID]ebml.Element{
		CuePointID: CuePointElement,
	},
}

var CuePointElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         CuePointID,
		ParentID:   CuesID,
		Name:       "cue_point",
		Definition: "Contains all information relative to a seek point in the Segment",
	},
	Master: map[ebml.ID]ebml.Element{
		CueTimeID:           CueTimeElement,
		CueTrackPositionsID: CueTrackPositionsElement,
	},
}
var CueTimeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueTimeID,
		ParentID:   CuePointID,
		Name:       "cue_time",
		Definition: "Absolute timestamp of the seek point",
	},
}

var CueTrackPositionsElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         CueTrackPositionsID,
		ParentID:   CuePointID,
		Name:       "cue_track_positions",
		Definition: "Contains positions for different tracks corresponding to the timestamp",
	},
	Master: map[ebml.ID]ebml.Element{
		CueTrackID:            CueTrackElement,
		CueClusterPositionID:  CueClusterPositionElement,
		CueRelativePositionID: CueRelativePositionElement,
		CueDurationID:         CueDurationElement,
		CueBlockNumberID:      CueBlockNumberElement,
		CueCodecStateID:       CueCodecStateElement,
		CueReferenceID:        CueReferenceElement,
	},
}
var CueTrackElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueTrackID,
		ParentID:   CueTrackPositionsID,
		Name:       "cue_track",
		Definition: "The track for which a position is given",
	},
}
var CueClusterPositionElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueClusterPositionID,
		ParentID:   CueTrackPositionsID,
		Name:       "cue_cluster_position",
		Definition: "The Segment Position of the Cluster containing the associated Block",
	},
}
var CueRelativePositionElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueRelativePositionID,
		ParentID:   CueTrackPositionsID,
		Name:       "cue_relative_position",
		Definition: "The relative position inside the Cluster of the referenced SimpleBlock or BlockGroup with 0 being the first possible position for an element inside that Cluster",
	},
}
var CueDurationElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueDurationID,
		ParentID:   CueTrackPositionsID,
		Name:       "cue_duration",
		Definition: "The duration of the block",
	},
}
var CueBlockNumberElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueBlockNumberID,
		ParentID:   CueTrackPositionsID,
		Name:       "cue_block_number",
		Definition: "Number of the Block in the specified Cluster",
	},
}
var CueCodecStateElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueCodecStateID,
		ParentID:   CueTrackPositionsID,
		Name:       "cue_codec_state",
		Definition: "The Segment Position of the Codec State corresponding to this Cues element",
	},
}

var CueReferenceElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         CueReferenceID,
		ParentID:   CueTrackPositionsID,
		Name:       "cue_reference",
		Definition: "The Clusters containing the referenced Blocks",
	},
	Master: map[ebml.ID]ebml.Element{
		CueRefTimeID:       CueRefTimeElement,
		CueRefClusterID:    CueRefClusterElement,
		CueRefNumberID:     CueRefNumberElement,
		CueRefCodecStateID: CueRefCodecStateElement,
	},
}
var CueRefTimeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueRefTimeID,
		ParentID:   CueReferenceID,
		Name:       "cue_ref_time",
		Definition: "Timestamp of the referenced Block",
	},
}
var CueRefClusterElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueRefClusterID,
		ParentID:   CueReferenceID,
		Name:       "cue_ref_cluster",
		Definition: "The Segment Position of the Cluster containing the referenced Block",
	},
}
var CueRefNumberElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueRefNumberID,
		ParentID:   CueReferenceID,
		Name:       "cue_ref_number",
		Definition: "Number of the referenced Block of Track X in the specified Cluster",
	},
}
var CueRefCodecStateElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         CueRefCodecStateID,
		ParentID:   CueReferenceID,
		Name:       "cue_ref_codec_state",
		Definition: "The Segment Position of the Codec State corresponding to this referenced element",
	},
}

var AttachmentsElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         AttachmentsID,
		ParentID:   SegmentID,
		Name:       "attachments",
		Definition: "Contains attached files",
	},
	Master: map[ebml.ID]ebml.Element{
		AttachedFileID: AttachedFileElement,
	},
}

var AttachedFileElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         AttachedFileID,
		ParentID:   AttachmentsID,
		Name:       "attached_file",
		Definition: "An attached file",
	},
	Master: map[ebml.ID]ebml.Element{
		FileDescriptionID:   FileDescriptionElement,
		FileNameID:          FileNameElement,
		FileMediaTypeID:     FileMediaTypeElement,
		FileDataID:          FileDataElement,
		FileUIDID:           FileUIDElement,
		FileReferralID:      FileReferralElement,
		FileUsedStartTimeID: FileUsedStartTimeElement,
		FileUsedEndTimeID:   FileUsedEndTimeElement,
	},
}
var FileDescriptionElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         FileDescriptionID,
		ParentID:   AttachedFileID,
		Name:       "file_description",
		Definition: "A human-friendly name for the attached file",
	},
}
var FileNameElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         FileNameID,
		ParentID:   AttachedFileID,
		Name:       "file_name",
		Definition: "Filename of the attached file",
	},
}
var FileMediaTypeElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         FileMediaTypeID,
		ParentID:   AttachedFileID,
		Name:       "file_media_type",
		Definition: "Media type of the file following the format described in RFC6838",
	},
}
var FileDataElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         FileDataID,
		ParentID:   AttachedFileID,
		Name:       "file_data",
		Definition: "The data of the file",
	},
}
var FileUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FileUIDID,
		ParentID:   AttachedFileID,
		Name:       "file_uid",
		Definition: "UID representing the file",
	},
}
var FileReferralElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         FileReferralID,
		ParentID:   AttachedFileID,
		Name:       "file_referral",
		Definition: "A binary value that a track/codec can refer to when the attachment is needed",
	},
}
var FileUsedStartTimeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FileUsedStartTimeID,
		ParentID:   AttachedFileID,
		Name:       "file_used_start_time",
		Definition: "The timestamp at which this optimized font attachment comes into context",
	},
}
var FileUsedEndTimeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         FileUsedEndTimeID,
		ParentID:   AttachedFileID,
		Name:       "file_used_end_time",
		Definition: "The timestamp at which this optimized font attachment goes out of context",
	},
}

var ChaptersElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ChaptersID,
		ParentID:   SegmentID,
		Name:       "chapters",
		Definition: "A system to define basic menus and partition data",
	},
	Master: map[ebml.ID]ebml.Element{
		EditionEntryID: EditionEntryElement,
	},
}

var EditionEntryElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         EditionEntryID,
		ParentID:   ChaptersID,
		Name:       "edition_entry",
		Definition: "Contains all information about a Segment edition",
	},
	Master: map[ebml.ID]ebml.Element{
		EditionUIDID:         EditionUIDElement,
		EditionFlagHiddenID:  EditionFlagHiddenElement,
		EditionFlagDefaultID: EditionFlagDefaultElement,
		EditionFlagOrderedID: EditionFlagOrderedElement,
		EditionDisplayID:     EditionDisplayElement,
		ChapterAtomID:        ChapterAtomElement,
	},
}
var EditionUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         EditionUIDID,
		ParentID:   EditionEntryID,
		Name:       "edition_uid",
		Definition: "A UID that identifies the edition",
	},
}
var EditionFlagHiddenElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         EditionFlagHiddenID,
		ParentID:   EditionEntryID,
		Name:       "edition_flag_hidden",
		Definition: "Set to 1 if an edition is hidden",
	},
}
var EditionFlagDefaultElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         EditionFlagDefaultID,
		ParentID:   EditionEntryID,
		Name:       "edition_flag_default",
		Definition: "Set to 1 if the edition **SHOULD** be used as the default one",
	},
}
var EditionFlagOrderedElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         EditionFlagOrderedID,
		ParentID:   EditionEntryID,
		Name:       "edition_flag_ordered",
		Definition: "Set to 1 if the chapters can be defined multiple times and the order to play them is enforced",
	},
}

var EditionDisplayElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         EditionDisplayID,
		ParentID:   EditionEntryID,
		Name:       "edition_display",
		Definition: "Contains a possible string to use for the edition display for the given languages",
	},
	Master: map[ebml.ID]ebml.Element{
		EditionStringID:       EditionStringElement,
		EditionLanguageIETFID: EditionLanguageIETFElement,
	},
}
var EditionStringElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         EditionStringID,
		ParentID:   EditionDisplayID,
		Name:       "edition_string",
		Definition: "Contains the string to use as the edition name",
	},
}
var EditionLanguageIETFElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         EditionLanguageIETFID,
		ParentID:   EditionDisplayID,
		Name:       "edition_language_ietf",
		Definition: "One language corresponding to the EditionString",
	},
}

var ChapterAtomElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ChapterAtomID,
		ParentID:   EditionEntryID,
		Name:       "chapter_atom",
		Definition: "Contains the atom information to use as the chapter atom",
	},
	Master: map[ebml.ID]ebml.Element{
		ChapterUIDID:               ChapterUIDElement,
		ChapterStringUIDID:         ChapterStringUIDElement,
		ChapterTimeStartID:         ChapterTimeStartElement,
		ChapterTimeEndID:           ChapterTimeEndElement,
		ChapterFlagHiddenID:        ChapterFlagHiddenElement,
		ChapterFlagEnabledID:       ChapterFlagEnabledElement,
		ChapterSegmentUUIDID:       ChapterSegmentUUIDElement,
		ChapterSkipTypeID:          ChapterSkipTypeElement,
		ChapterSegmentEditionUIDID: ChapterSegmentEditionUIDElement,
		ChapterPhysicalEquivID:     ChapterPhysicalEquivElement,
		ChapterTrackID:             ChapterTrackElement,
		ChapterDisplayID:           ChapterDisplayElement,
		ChapProcessID:              ChapProcessElement,
	},
}
var ChapterUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterUIDID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_uid",
		Definition: "A UID that identifies the Chapter",
	},
}
var ChapterStringUIDElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         ChapterStringUIDID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_string_uid",
		Definition: "A unique string ID that identifies the Chapter",
	},
}
var ChapterTimeStartElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterTimeStartID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_time_start",
		Definition: "Timestamp of the start of Chapter",
	},
}
var ChapterTimeEndElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterTimeEndID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_time_end",
		Definition: "Timestamp of the end of Chapter ",
	},
}
var ChapterFlagHiddenElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterFlagHiddenID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_flag_hidden",
		Definition: "Set to 1 if a chapter is hidden",
	},
}
var ChapterFlagEnabledElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterFlagEnabledID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_flag_enabled",
		Definition: "Set to 1 if the chapter is enabled",
	},
}
var ChapterSegmentUUIDElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ChapterSegmentUUIDID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_segment_uuid",
		Definition: "The SegmentUUID of another Segment to play during this chapter",
	},
}
var ChapterSkipTypeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterSkipTypeID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_skip_type",
		Definition: "Indicates what type of content the ChapterAtom contains and might be skipped",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "no_skipping", Description: "Content which should not be skipped"},
		1: {Name: "opening_credits", Description: "Credits usually found at the beginning of the content"},
		2: {Name: "end_credits", Description: "Credits usually found at the end of the content"},
		3: {Name: "recap", Description: "Recap of previous episodes of the content"},
		4: {Name: "next_preview", Description: "Preview of the next episode of the content"},
		5: {Name: "preview", Description: "Preview of the current episode of the content"},
		6: {Name: "advertisement", Description: "Advertisement within the content"},
	},
}
var ChapterSegmentEditionUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterSegmentEditionUIDID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_segment_edition_uid",
		Definition: "The EditionUID to play from the Segment linked in ChapterSegmentUUID",
	},
}
var ChapterPhysicalEquivElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterPhysicalEquivID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_physical_equiv",
		Definition: "Specifies the physical equivalent of this ChapterAtom",
	},
}

var ChapterTrackElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ChapterTrackID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_track",
		Definition: "List of tracks on which the chapter applies",
	},
	Master: map[ebml.ID]ebml.Element{
		ChapterTrackUIDID: ChapterTrackUIDElement,
	},
}
var ChapterTrackUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapterTrackUIDID,
		ParentID:   ChapterTrackID,
		Name:       "chapter_track_uid",
		Definition: "UID of the Track to apply this chapter to",
	},
}

var ChapterDisplayElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ChapterDisplayID,
		ParentID:   ChapterAtomID,
		Name:       "chapter_display",
		Definition: "Contains all possible strings to use for the chapter display",
	},
	Master: map[ebml.ID]ebml.Element{
		ChapStringID:        ChapStringElement,
		ChapLanguageID:      ChapLanguageElement,
		ChapLanguageBCP47ID: ChapLanguageBCP47Element,
		ChapCountryID:       ChapCountryElement,
	},
}
var ChapStringElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         ChapStringID,
		ParentID:   ChapterDisplayID,
		Name:       "chap_string",
		Definition: "Contains the string to use as the chapter atom",
	},
}
var ChapLanguageElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         ChapLanguageID,
		ParentID:   ChapterDisplayID,
		Name:       "chap_language",
		Definition: "A language corresponding to the string",
	},
}
var ChapLanguageBCP47Element = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         ChapLanguageBCP47ID,
		ParentID:   ChapterDisplayID,
		Name:       "chap_language_bcp47",
		Definition: "A language corresponding to the ChapString",
	},
}
var ChapCountryElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         ChapCountryID,
		ParentID:   ChapterDisplayID,
		Name:       "chap_country",
		Definition: "A country corresponding to the string",
	},
}

var ChapProcessElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ChapProcessID,
		ParentID:   ChapterAtomID,
		Name:       "chap_process",
		Definition: "Contains all the commands associated with the Atom",
	},
	Master: map[ebml.ID]ebml.Element{
		ChapProcessCodecIDID: ChapProcessCodecIDElement,
		ChapProcessPrivateID: ChapProcessPrivateElement,
		ChapProcessCommandID: ChapProcessCommandElement,
	},
}
var ChapProcessCodecIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapProcessCodecIDID,
		ParentID:   ChapProcessID,
		Name:       "chap_process_codec_id",
		Definition: "Contains the type of the codec used for processing",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "matroska_script", Description: "Chapter commands using the Matroska Script codec"},
		1: {Name: "dvd_menu", Description: "Chapter commands using the DVD-like codec"},
	},
}
var ChapProcessPrivateElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ChapProcessPrivateID,
		ParentID:   ChapProcessID,
		Name:       "chap_process_private",
		Definition: "Optional data attached to the ChapProcessCodecID information",
	},
}

var ChapProcessCommandElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         ChapProcessCommandID,
		ParentID:   ChapProcessID,
		Name:       "chap_process_command",
		Definition: "Contains all the commands associated with the Atom",
	},
	Master: map[ebml.ID]ebml.Element{
		ChapProcessTimeID: ChapProcessTimeElement,
		ChapProcessDataID: ChapProcessDataElement,
	},
}
var ChapProcessTimeElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         ChapProcessTimeID,
		ParentID:   ChapProcessCommandID,
		Name:       "chap_process_time",
		Definition: "Defines when the process command **SHOULD** be handled",
	},
	Enums: map[uint64]ebml.Enum{
		0: {Name: "during_the_whole_chapter"},
		1: {Name: "before_starting_playback"},
		2: {Name: "after_playback_of_the_chapter"},
	},
}
var ChapProcessDataElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         ChapProcessDataID,
		ParentID:   ChapProcessCommandID,
		Name:       "chap_process_data",
		Definition: "Contains the command information",
	},
}

var TagsElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TagsID,
		ParentID:   SegmentID,
		Name:       "tags",
		Definition: "Element containing metadata describing Tracks",
	},
	Master: map[ebml.ID]ebml.Element{
		TagID: TagElement,
	},
}

var TagElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TagID,
		ParentID:   TagsID,
		Name:       "tag",
		Definition: "A single metadata descriptor",
	},
	Master: map[ebml.ID]ebml.Element{
		TargetsID:   TargetsElement,
		SimpleTagID: SimpleTagElement,
	},
}

var TargetsElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         TargetsID,
		ParentID:   TagID,
		Name:       "targets",
		Definition: "Specifies which other elements the metadata represented by the tag value applies to",
	},
	Master: map[ebml.ID]ebml.Element{
		TargetTypeValueID:  TargetTypeValueElement,
		TargetTypeID:       TargetTypeElement,
		TagTrackUIDID:      TagTrackUIDElement,
		TagEditionUIDID:    TagEditionUIDElement,
		TagChapterUIDID:    TagChapterUIDElement,
		TagAttachmentUIDID: TagAttachmentUIDElement,
	},
}
var TargetTypeValueElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TargetTypeValueID,
		ParentID:   TargetsID,
		Name:       "target_type_value",
		Definition: "A number to indicate the logical level of the target",
	},
	Enums: map[uint64]ebml.Enum{
		10: {Name: "shot", Description: "The lowest hierarchy found in music or movies"},
		20: {Name: "subtrack_movement_scene", Description: "Corresponds to parts of a track for audio"},
		30: {Name: "track_song_chapter", Description: "The common parts of an album or movie"},
		40: {Name: "part_session", Description: "When an album or episode has different logical parts"},
		50: {Name: "album_opera_concert_movie_episode", Description: "The most common grouping level of music and video"},
		60: {Name: "edition_issue_volume_opus_season_sequel", Description: "A list of lower levels grouped together"},
		70: {Name: "collection", Description: "The highest hierarchical level that tags can describe"},
	},
}
var TargetTypeElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         TargetTypeID,
		ParentID:   TargetsID,
		Name:       "target_type",
		Definition: "An informational string that can be used to display the logical level of the target",
	},
	Enums: map[string]ebml.Enum{
		"COLLECTION": {Name: "targettypevalue_70"},
		"EDITION":    {Name: "targettypevalue_60"},
		"ISSUE":      {Name: "targettypevalue_60"},
		"VOLUME":     {Name: "targettypevalue_60"},
		"OPUS":       {Name: "targettypevalue_60"},
		"SEASON":     {Name: "targettypevalue_60"},
		"SEQUEL":     {Name: "targettypevalue_60"},
		"ALBUM":      {Name: "targettypevalue_50"},
		"OPERA":      {Name: "targettypevalue_50"},
		"CONCERT":    {Name: "targettypevalue_50"},
		"MOVIE":      {Name: "targettypevalue_50"},
		"EPISODE":    {Name: "targettypevalue_50"},
		"PART":       {Name: "targettypevalue_40"},
		"SESSION":    {Name: "targettypevalue_40"},
		"TRACK":      {Name: "targettypevalue_30"},
		"SONG":       {Name: "targettypevalue_30"},
		"CHAPTER":    {Name: "targettypevalue_30"},
		"SUBTRACK":   {Name: "targettypevalue_20"},
		"MOVEMENT":   {Name: "targettypevalue_20"},
		"SCENE":      {Name: "targettypevalue_20"},
		"SHOT":       {Name: "targettypevalue_10"},
	},
}
var TagTrackUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TagTrackUIDID,
		ParentID:   TargetsID,
		Name:       "tag_track_uid",
		Definition: "A UID that identifies the Track(s) that the tags belong to",
	},
}
var TagEditionUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TagEditionUIDID,
		ParentID:   TargetsID,
		Name:       "tag_edition_uid",
		Definition: "A UID that identifies the EditionEntry(s) that the tags belong to",
	},
}
var TagChapterUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TagChapterUIDID,
		ParentID:   TargetsID,
		Name:       "tag_chapter_uid",
		Definition: "A UID that identifies the Chapter(s) that the tags belong to",
	},
}
var TagAttachmentUIDElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TagAttachmentUIDID,
		ParentID:   TargetsID,
		Name:       "tag_attachment_uid",
		Definition: "A UID that identifies the Attachment(s) that the tags belong to",
	},
}

var SimpleTagElement = &ebml.Master{
	ElementType: ebml.ElementType{
		ID:         SimpleTagID,
		ParentID:   TagID,
		Name:       "simple_tag",
		Definition: "Contains general information about the target",
	},
	Master: map[ebml.ID]ebml.Element{
		TagNameID:          TagNameElement,
		TagLanguageID:      TagLanguageElement,
		TagLanguageBCP47ID: TagLanguageBCP47Element,
		TagDefaultID:       TagDefaultElement,
		TagDefaultBogusID:  TagDefaultBogusElement,
		TagStringID:        TagStringElement,
		TagBinaryID:        TagBinaryElement,
	},
}
var TagNameElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         TagNameID,
		ParentID:   SimpleTagID,
		Name:       "tag_name",
		Definition: "The name of the tag value that is going to be stored",
	},
}
var TagLanguageElement = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         TagLanguageID,
		ParentID:   SimpleTagID,
		Name:       "tag_language",
		Definition: "Specifies the language of the specified tag in the Matroska languages form",
	},
}
var TagLanguageBCP47Element = &ebml.String{
	ElementType: ebml.ElementType{
		ID:         TagLanguageBCP47ID,
		ParentID:   SimpleTagID,
		Name:       "tag_language_bcp47",
		Definition: "The language used in the TagString",
	},
}
var TagDefaultElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TagDefaultID,
		ParentID:   SimpleTagID,
		Name:       "tag_default",
		Definition: "A boolean value to indicate if this is the default/original language to use for the given tag",
	},
}
var TagDefaultBogusElement = &ebml.Uinteger{
	ElementType: ebml.ElementType{
		ID:         TagDefaultBogusID,
		ParentID:   SimpleTagID,
		Name:       "tag_default_bogus",
		Definition: "A variant of the TagDefault element with a bogus element ID",
	},
}
var TagStringElement = &ebml.UTF8{
	ElementType: ebml.ElementType{
		ID:         TagStringID,
		ParentID:   SimpleTagID,
		Name:       "tag_string",
		Definition: "The tag value",
	},
}
var TagBinaryElement = &ebml.Binary{
	ElementType: ebml.ElementType{
		ID:         TagBinaryID,
		ParentID:   SimpleTagID,
		Name:       "tag_binary",
		Definition: "The tag value if it is binary",
	},
}

var IDToElement = map[ebml.ID]ebml.Element{
	RootID:                        RootElement,
	SegmentID:                     SegmentElement,
	SeekHeadID:                    SeekHeadElement,
	SeekID:                        SeekElement,
	SeekIDID:                      SeekIDElement,
	SeekPositionID:                SeekPositionElement,
	InfoID:                        InfoElement,
	SegmentUUIDID:                 SegmentUUIDElement,
	SegmentFilenameID:             SegmentFilenameElement,
	PrevUUIDID:                    PrevUUIDElement,
	PrevFilenameID:                PrevFilenameElement,
	NextUUIDID:                    NextUUIDElement,
	NextFilenameID:                NextFilenameElement,
	SegmentFamilyID:               SegmentFamilyElement,
	TimestampScaleID:              TimestampScaleElement,
	DurationID:                    DurationElement,
	DateUTCID:                     DateUTCElement,
	TitleID:                       TitleElement,
	MuxingAppID:                   MuxingAppElement,
	WritingAppID:                  WritingAppElement,
	ChapterTranslateID:            ChapterTranslateElement,
	ChapterTranslateIDID:          ChapterTranslateIDElement,
	ChapterTranslateCodecID:       ChapterTranslateCodecElement,
	ChapterTranslateEditionUIDID:  ChapterTranslateEditionUIDElement,
	ClusterID:                     ClusterElement,
	TimestampID:                   TimestampElement,
	PositionID:                    PositionElement,
	PrevSizeID:                    PrevSizeElement,
	SimpleBlockID:                 SimpleBlockElement,
	EncryptedBlockID:              EncryptedBlockElement,
	SilentTracksID:                SilentTracksElement,
	SilentTrackNumberID:           SilentTrackNumberElement,
	BlockGroupID:                  BlockGroupElement,
	BlockID:                       BlockElement,
	BlockVirtualID:                BlockVirtualElement,
	BlockDurationID:               BlockDurationElement,
	ReferencePriorityID:           ReferencePriorityElement,
	ReferenceBlockID:              ReferenceBlockElement,
	ReferenceVirtualID:            ReferenceVirtualElement,
	CodecStateID:                  CodecStateElement,
	DiscardPaddingID:              DiscardPaddingElement,
	BlockAdditionsID:              BlockAdditionsElement,
	BlockMoreID:                   BlockMoreElement,
	BlockAdditionalID:             BlockAdditionalElement,
	BlockAddIDID:                  BlockAddIDElement,
	SlicesID:                      SlicesElement,
	TimeSliceID:                   TimeSliceElement,
	LaceNumberID:                  LaceNumberElement,
	FrameNumberID:                 FrameNumberElement,
	BlockAdditionIDID:             BlockAdditionIDElement,
	DelayID:                       DelayElement,
	SliceDurationID:               SliceDurationElement,
	ReferenceFrameID:              ReferenceFrameElement,
	ReferenceOffsetID:             ReferenceOffsetElement,
	ReferenceTimestampID:          ReferenceTimestampElement,
	TracksID:                      TracksElement,
	TrackEntryID:                  TrackEntryElement,
	TrackNumberID:                 TrackNumberElement,
	TrackUIDID:                    TrackUIDElement,
	TrackTypeID:                   TrackTypeElement,
	FlagEnabledID:                 FlagEnabledElement,
	FlagDefaultID:                 FlagDefaultElement,
	FlagForcedID:                  FlagForcedElement,
	FlagHearingImpairedID:         FlagHearingImpairedElement,
	FlagVisualImpairedID:          FlagVisualImpairedElement,
	FlagTextDescriptionsID:        FlagTextDescriptionsElement,
	FlagOriginalID:                FlagOriginalElement,
	FlagCommentaryID:              FlagCommentaryElement,
	FlagLacingID:                  FlagLacingElement,
	MinCacheID:                    MinCacheElement,
	MaxCacheID:                    MaxCacheElement,
	DefaultDurationID:             DefaultDurationElement,
	DefaultDecodedFieldDurationID: DefaultDecodedFieldDurationElement,
	TrackTimestampScaleID:         TrackTimestampScaleElement,
	TrackOffsetID:                 TrackOffsetElement,
	MaxBlockAdditionIDID:          MaxBlockAdditionIDElement,
	NameID:                        NameElement,
	LanguageID:                    LanguageElement,
	LanguageBCP47ID:               LanguageBCP47Element,
	CodecIDID:                     CodecIDElement,
	CodecPrivateID:                CodecPrivateElement,
	CodecNameID:                   CodecNameElement,
	AttachmentLinkID:              AttachmentLinkElement,
	CodecSettingsID:               CodecSettingsElement,
	CodecInfoURLID:                CodecInfoURLElement,
	CodecDownloadURLID:            CodecDownloadURLElement,
	CodecDecodeAllID:              CodecDecodeAllElement,
	TrackOverlayID:                TrackOverlayElement,
	CodecDelayID:                  CodecDelayElement,
	SeekPreRollID:                 SeekPreRollElement,
	TrickTrackUIDID:               TrickTrackUIDElement,
	TrickTrackSegmentUIDID:        TrickTrackSegmentUIDElement,
	TrickTrackFlagID:              TrickTrackFlagElement,
	TrickMasterTrackUIDID:         TrickMasterTrackUIDElement,
	TrickMasterTrackSegmentUIDID:  TrickMasterTrackSegmentUIDElement,
	BlockAdditionMappingID:        BlockAdditionMappingElement,
	BlockAddIDValueID:             BlockAddIDValueElement,
	BlockAddIDNameID:              BlockAddIDNameElement,
	BlockAddIDTypeID:              BlockAddIDTypeElement,
	BlockAddIDExtraDataID:         BlockAddIDExtraDataElement,
	TrackTranslateID:              TrackTranslateElement,
	TrackTranslateTrackIDID:       TrackTranslateTrackIDElement,
	TrackTranslateCodecID:         TrackTranslateCodecElement,
	TrackTranslateEditionUIDID:    TrackTranslateEditionUIDElement,
	VideoID:                       VideoElement,
	FlagInterlacedID:              FlagInterlacedElement,
	FieldOrderID:                  FieldOrderElement,
	StereoModeID:                  StereoModeElement,
	AlphaModeID:                   AlphaModeElement,
	OldStereoModeID:               OldStereoModeElement,
	PixelWidthID:                  PixelWidthElement,
	PixelHeightID:                 PixelHeightElement,
	PixelCropBottomID:             PixelCropBottomElement,
	PixelCropTopID:                PixelCropTopElement,
	PixelCropLeftID:               PixelCropLeftElement,
	PixelCropRightID:              PixelCropRightElement,
	DisplayWidthID:                DisplayWidthElement,
	DisplayHeightID:               DisplayHeightElement,
	DisplayUnitID:                 DisplayUnitElement,
	AspectRatioTypeID:             AspectRatioTypeElement,
	UncompressedFourCCID:          UncompressedFourCCElement,
	GammaValueID:                  GammaValueElement,
	FrameRateID:                   FrameRateElement,
	ColourID:                      ColourElement,
	MatrixCoefficientsID:          MatrixCoefficientsElement,
	BitsPerChannelID:              BitsPerChannelElement,
	ChromaSubsamplingHorzID:       ChromaSubsamplingHorzElement,
	ChromaSubsamplingVertID:       ChromaSubsamplingVertElement,
	CbSubsamplingHorzID:           CbSubsamplingHorzElement,
	CbSubsamplingVertID:           CbSubsamplingVertElement,
	ChromaSitingHorzID:            ChromaSitingHorzElement,
	ChromaSitingVertID:            ChromaSitingVertElement,
	RangeID:                       RangeElement,
	TransferCharacteristicsID:     TransferCharacteristicsElement,
	PrimariesID:                   PrimariesElement,
	MaxCLLID:                      MaxCLLElement,
	MaxFALLID:                     MaxFALLElement,
	MasteringMetadataID:           MasteringMetadataElement,
	PrimaryRChromaticityXID:       PrimaryRChromaticityXElement,
	PrimaryRChromaticityYID:       PrimaryRChromaticityYElement,
	PrimaryGChromaticityXID:       PrimaryGChromaticityXElement,
	PrimaryGChromaticityYID:       PrimaryGChromaticityYElement,
	PrimaryBChromaticityXID:       PrimaryBChromaticityXElement,
	PrimaryBChromaticityYID:       PrimaryBChromaticityYElement,
	WhitePointChromaticityXID:     WhitePointChromaticityXElement,
	WhitePointChromaticityYID:     WhitePointChromaticityYElement,
	LuminanceMaxID:                LuminanceMaxElement,
	LuminanceMinID:                LuminanceMinElement,
	ProjectionID:                  ProjectionElement,
	ProjectionTypeID:              ProjectionTypeElement,
	ProjectionPrivateID:           ProjectionPrivateElement,
	ProjectionPoseYawID:           ProjectionPoseYawElement,
	ProjectionPosePitchID:         ProjectionPosePitchElement,
	ProjectionPoseRollID:          ProjectionPoseRollElement,
	AudioID:                       AudioElement,
	SamplingFrequencyID:           SamplingFrequencyElement,
	OutputSamplingFrequencyID:     OutputSamplingFrequencyElement,
	ChannelsID:                    ChannelsElement,
	ChannelPositionsID:            ChannelPositionsElement,
	BitDepthID:                    BitDepthElement,
	EmphasisID:                    EmphasisElement,
	TrackOperationID:              TrackOperationElement,
	TrackCombinePlanesID:          TrackCombinePlanesElement,
	TrackPlaneID:                  TrackPlaneElement,
	TrackPlaneUIDID:               TrackPlaneUIDElement,
	TrackPlaneTypeID:              TrackPlaneTypeElement,
	TrackJoinBlocksID:             TrackJoinBlocksElement,
	TrackJoinUIDID:                TrackJoinUIDElement,
	ContentEncodingsID:            ContentEncodingsElement,
	ContentEncodingID:             ContentEncodingElement,
	ContentEncodingOrderID:        ContentEncodingOrderElement,
	ContentEncodingScopeID:        ContentEncodingScopeElement,
	ContentEncodingTypeID:         ContentEncodingTypeElement,
	ContentCompressionID:          ContentCompressionElement,
	ContentCompAlgoID:             ContentCompAlgoElement,
	ContentCompSettingsID:         ContentCompSettingsElement,
	ContentEncryptionID:           ContentEncryptionElement,
	ContentEncAlgoID:              ContentEncAlgoElement,
	ContentEncKeyIDID:             ContentEncKeyIDElement,
	ContentSignatureID:            ContentSignatureElement,
	ContentSigKeyIDID:             ContentSigKeyIDElement,
	ContentSigAlgoID:              ContentSigAlgoElement,
	ContentSigHashAlgoID:          ContentSigHashAlgoElement,
	ContentEncAESSettingsID:       ContentEncAESSettingsElement,
	AESSettingsCipherModeID:       AESSettingsCipherModeElement,
	CuesID:                        CuesElement,
	CuePointID:                    CuePointElement,
	CueTimeID:                     CueTimeElement,
	CueTrackPositionsID:           CueTrackPositionsElement,
	CueTrackID:                    CueTrackElement,
	CueClusterPositionID:          CueClusterPositionElement,
	CueRelativePositionID:         CueRelativePositionElement,
	CueDurationID:                 CueDurationElement,
	CueBlockNumberID:              CueBlockNumberElement,
	CueCodecStateID:               CueCodecStateElement,
	CueReferenceID:                CueReferenceElement,
	CueRefTimeID:                  CueRefTimeElement,
	CueRefClusterID:               CueRefClusterElement,
	CueRefNumberID:                CueRefNumberElement,
	CueRefCodecStateID:            CueRefCodecStateElement,
	AttachmentsID:                 AttachmentsElement,
	AttachedFileID:                AttachedFileElement,
	FileDescriptionID:             FileDescriptionElement,
	FileNameID:                    FileNameElement,
	FileMediaTypeID:               FileMediaTypeElement,
	FileDataID:                    FileDataElement,
	FileUIDID:                     FileUIDElement,
	FileReferralID:                FileReferralElement,
	FileUsedStartTimeID:           FileUsedStartTimeElement,
	FileUsedEndTimeID:             FileUsedEndTimeElement,
	ChaptersID:                    ChaptersElement,
	EditionEntryID:                EditionEntryElement,
	EditionUIDID:                  EditionUIDElement,
	EditionFlagHiddenID:           EditionFlagHiddenElement,
	EditionFlagDefaultID:          EditionFlagDefaultElement,
	EditionFlagOrderedID:          EditionFlagOrderedElement,
	EditionDisplayID:              EditionDisplayElement,
	EditionStringID:               EditionStringElement,
	EditionLanguageIETFID:         EditionLanguageIETFElement,
	ChapterAtomID:                 ChapterAtomElement,
	ChapterUIDID:                  ChapterUIDElement,
	ChapterStringUIDID:            ChapterStringUIDElement,
	ChapterTimeStartID:            ChapterTimeStartElement,
	ChapterTimeEndID:              ChapterTimeEndElement,
	ChapterFlagHiddenID:           ChapterFlagHiddenElement,
	ChapterFlagEnabledID:          ChapterFlagEnabledElement,
	ChapterSegmentUUIDID:          ChapterSegmentUUIDElement,
	ChapterSkipTypeID:             ChapterSkipTypeElement,
	ChapterSegmentEditionUIDID:    ChapterSegmentEditionUIDElement,
	ChapterPhysicalEquivID:        ChapterPhysicalEquivElement,
	ChapterTrackID:                ChapterTrackElement,
	ChapterTrackUIDID:             ChapterTrackUIDElement,
	ChapterDisplayID:              ChapterDisplayElement,
	ChapStringID:                  ChapStringElement,
	ChapLanguageID:                ChapLanguageElement,
	ChapLanguageBCP47ID:           ChapLanguageBCP47Element,
	ChapCountryID:                 ChapCountryElement,
	ChapProcessID:                 ChapProcessElement,
	ChapProcessCodecIDID:          ChapProcessCodecIDElement,
	ChapProcessPrivateID:          ChapProcessPrivateElement,
	ChapProcessCommandID:          ChapProcessCommandElement,
	ChapProcessTimeID:             ChapProcessTimeElement,
	ChapProcessDataID:             ChapProcessDataElement,
	TagsID:                        TagsElement,
	TagID:                         TagElement,
	TargetsID:                     TargetsElement,
	TargetTypeValueID:             TargetTypeValueElement,
	TargetTypeID:                  TargetTypeElement,
	TagTrackUIDID:                 TagTrackUIDElement,
	TagEditionUIDID:               TagEditionUIDElement,
	TagChapterUIDID:               TagChapterUIDElement,
	TagAttachmentUIDID:            TagAttachmentUIDElement,
	SimpleTagID:                   SimpleTagElement,
	TagNameID:                     TagNameElement,
	TagLanguageID:                 TagLanguageElement,
	TagLanguageBCP47ID:            TagLanguageBCP47Element,
	TagDefaultID:                  TagDefaultElement,
	TagDefaultBogusID:             TagDefaultBogusElement,
	TagStringID:                   TagStringElement,
	TagBinaryID:                   TagBinaryElement,
}
