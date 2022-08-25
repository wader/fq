// Code below generated from ebml_matroska.xml
//
//nolint:revive
package ebml_matroska

import (
	"github.com/wader/fq/format/matroska/ebml"
	"github.com/wader/fq/pkg/scalar"
)

var Root = ebml.Tag{
	ebml.HeaderID: {Name: "ebml", Type: ebml.Master, Tag: ebml.Header},
	SegmentID:     {Name: "segment", Type: ebml.Master, Tag: Segment},
}

const (
	EBMLMaxIDLengthID             = 0x42f2
	EBMLMaxSizeLengthID           = 0x42f3
	SegmentID                     = 0x18538067
	SeekHeadID                    = 0x114d9b74
	SeekID                        = 0x4dbb
	SeekIDID                      = 0x53ab
	SeekPositionID                = 0x53ac
	InfoID                        = 0x1549a966
	SegmentUIDID                  = 0x73a4
	SegmentFilenameID             = 0x7384
	PrevUIDID                     = 0x3cb923
	PrevFilenameID                = 0x3c83ab
	NextUIDID                     = 0x3eb923
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
	BlockAddIDID                  = 0xee
	BlockAdditionalID             = 0xa5
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
	LanguageIETFID                = 0x22b59d
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
	FileMimeTypeID                = 0x4660
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
	ChapterAtomID                 = 0xb6
	ChapterUIDID                  = 0x73c4
	ChapterStringUIDID            = 0x5654
	ChapterTimeStartID            = 0x91
	ChapterTimeEndID              = 0x92
	ChapterFlagHiddenID           = 0x98
	ChapterFlagEnabledID          = 0x4598
	ChapterSegmentUIDID           = 0x6e67
	ChapterSegmentEditionUIDID    = 0x6ebc
	ChapterPhysicalEquivID        = 0x63c3
	ChapterTrackID                = 0x8f
	ChapterTrackUIDID             = 0x89
	ChapterDisplayID              = 0x80
	ChapStringID                  = 0x85
	ChapLanguageID                = 0x437c
	ChapLanguageIETFID            = 0x437d
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
	TagLanguageIETFID             = 0x447b
	TagDefaultID                  = 0x4484
	TagDefaultBogusID             = 0x44b4
	TagStringID                   = 0x4487
	TagBinaryID                   = 0x4485
)

var Segment = ebml.Tag{
	SeekHeadID: {
		Name:       "seek_head",
		Definition: "Contains the Segment Position of other Top-Level Elements.",
		Type:       ebml.Master, Tag: SeekHead,
	},
	InfoID: {
		Name:       "info",
		Definition: "Contains general information about the Segment.",
		Type:       ebml.Master, Tag: Info,
	},
	ClusterID: {
		Name:       "cluster",
		Definition: "The Top-Level Element containing the (monolithic) Block structure.",
		Type:       ebml.Master, Tag: Cluster,
	},
	TracksID: {
		Name:       "tracks",
		Definition: "A Top-Level Element of information with many tracks described.",
		Type:       ebml.Master, Tag: Tracks,
	},
	CuesID: {
		Name:       "cues",
		Definition: "A Top-Level Element to speed seeking access. All entries are local to the Segment.",
		Type:       ebml.Master, Tag: Cues,
	},
	AttachmentsID: {
		Name:       "attachments",
		Definition: "Contain attached files.",
		Type:       ebml.Master, Tag: Attachments,
	},
	ChaptersID: {
		Name:       "chapters",
		Definition: "A system to define basic menus and partition data. For more detailed information, look at the Chapters explanation in (#chapters).",
		Type:       ebml.Master, Tag: Chapters,
	},
	TagsID: {
		Name:       "tags",
		Definition: "Element containing metadata describing Tracks, Editions, Chapters, Attachments, or the Segment as a whole. A list of valid tags can be found in [@!MatroskaTags].",
		Type:       ebml.Master, Tag: Tags,
	},
}

var SeekHead = ebml.Tag{
	SeekID: {
		Name:       "seek",
		Definition: "Contains a single seek entry to an EBML Element.",
		Type:       ebml.Master, Tag: Seek,
	},
}

var Seek = ebml.Tag{
	SeekIDID: {
		Name:       "seek_id",
		Definition: "The binary ID corresponding to the Element name.",
		Type:       ebml.Binary,
	},
	SeekPositionID: {
		Name:       "seek_position",
		Definition: "The Segment Position of the Element.",
		Type:       ebml.Uinteger,
	},
}

var Info = ebml.Tag{
	SegmentUIDID: {
		Name:       "segment_uid",
		Definition: "A randomly generated unique ID to identify the Segment amongst many others (128 bits).",
		Type:       ebml.Binary,
	},
	SegmentFilenameID: {
		Name:       "segment_filename",
		Definition: "A filename corresponding to this Segment.",
		Type:       ebml.UTF8,
	},
	PrevUIDID: {
		Name:       "prev_uid",
		Definition: "A unique ID to identify the previous Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	PrevFilenameID: {
		Name:       "prev_filename",
		Definition: "A filename corresponding to the file of the previous Linked Segment.",
		Type:       ebml.UTF8,
	},
	NextUIDID: {
		Name:       "next_uid",
		Definition: "A unique ID to identify the next Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	NextFilenameID: {
		Name:       "next_filename",
		Definition: "A filename corresponding to the file of the next Linked Segment.",
		Type:       ebml.UTF8,
	},
	SegmentFamilyID: {
		Name:       "segment_family",
		Definition: "A randomly generated unique ID that all Segments of a Linked Segment **MUST** share (128 bits).",
		Type:       ebml.Binary,
	},
	ChapterTranslateID: {
		Name:       "chapter_translate",
		Definition: "The mapping between this `Segment` and a segment value in the given Chapter Codec.",
		Type:       ebml.Master, Tag: ChapterTranslate,
	},
	TimestampScaleID: {
		Name:       "timestamp_scale",
		Definition: "Base unit for Segment Ticks and Track Ticks, in nanoseconds. A TimestampScale value of 1.000.000 means scaled timestamps in the Segment are expressed in milliseconds; see (#timestamps) on how to interpret timestamps.",
		Type:       ebml.Uinteger,
	},
	DurationID: {
		Name:       "duration",
		Definition: "Duration of the Segment, expressed in Segment Ticks which is based on TimestampScale; see (#timestamp-ticks).",
		Type:       ebml.Float,
	},
	DateUTCID: {
		Name:       "date_utc",
		Definition: "The date and time that the Segment was created by the muxing application or library.",
		Type:       ebml.Date,
	},
	TitleID: {
		Name:       "title",
		Definition: "General name of the Segment.",
		Type:       ebml.UTF8,
	},
	MuxingAppID: {
		Name:       "muxing_app",
		Definition: "Muxing application or library (example: \"libmatroska-0.4.3\").",
		Type:       ebml.UTF8,
	},
	WritingAppID: {
		Name:       "writing_app",
		Definition: "Writing application (example: \"mkvmerge-0.3.3\").",
		Type:       ebml.UTF8,
	},
}

var ChapterTranslate = ebml.Tag{
	ChapterTranslateIDID: {
		Name:       "chapter_translate_id",
		Definition: "The binary value used to represent this Segment in the chapter codec data. The format depends on the ChapProcessCodecID used; see (#chapprocesscodecid-element).",
		Type:       ebml.Binary,
	},
	ChapterTranslateCodecID: {
		Name:       "chapter_translate_codec",
		Definition: "This `ChapterTranslate` applies to this chapter codec of the given chapter edition(s); see (#chapprocesscodecid-element).",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym:         "matroska_script",
				Description: "Chapter commands using the Matroska Script codec.",
			},
			1: {
				Sym:         "dvd_menu",
				Description: "Chapter commands using the DVD-like codec.",
			},
		},
	},
	ChapterTranslateEditionUIDID: {
		Name:       "chapter_translate_edition_uid",
		Definition: "Specify a chapter edition UID on which this `ChapterTranslate` applies.",
		Type:       ebml.Uinteger,
	},
}

var Cluster = ebml.Tag{
	TimestampID: {
		Name:       "timestamp",
		Definition: "Absolute timestamp of the cluster, expressed in Segment Ticks which is based on TimestampScale; see (#timestamp-ticks).",
		Type:       ebml.Uinteger,
	},
	SilentTracksID: {
		Name:       "silent_tracks",
		Definition: "The list of tracks that are not used in that part of the stream. It is useful when using overlay tracks on seeking or to decide what track to use.",
		Type:       ebml.Master, Tag: SilentTracks,
	},
	PositionID: {
		Name:       "position",
		Definition: "The Segment Position of the Cluster in the Segment (0 in live streams). It might help to resynchronise offset on damaged streams.",
		Type:       ebml.Uinteger,
	},
	PrevSizeID: {
		Name:       "prev_size",
		Definition: "Size of the previous Cluster, in octets. Can be useful for backward playing.",
		Type:       ebml.Uinteger,
	},
	SimpleBlockID: {
		Name:       "simple_block",
		Definition: "Similar to Block, see (#block-structure), but without all the extra information, mostly used to reduced overhead when no extra feature is needed; see (#simpleblock-structure) on SimpleBlock Structure.",
		Type:       ebml.Binary,
	},
	BlockGroupID: {
		Name:       "block_group",
		Definition: "Basic container of information containing a single Block and information specific to that Block.",
		Type:       ebml.Master, Tag: BlockGroup,
	},
	EncryptedBlockID: {
		Name:       "encrypted_block",
		Definition: "Similar to SimpleBlock, see (#simpleblock-structure), but the data inside the Block are Transformed (encrypt and/or signed).",
		Type:       ebml.Binary,
	},
}

var SilentTracks = ebml.Tag{
	SilentTrackNumberID: {
		Name:       "silent_track_number",
		Definition: "One of the track number that are not used from now on in the stream. It could change later if not specified as silent in a further Cluster.",
		Type:       ebml.Uinteger,
	},
}

var BlockGroup = ebml.Tag{
	BlockID: {
		Name:       "block",
		Definition: "Block containing the actual data to be rendered and a timestamp relative to the Cluster Timestamp; see (#block-structure) on Block Structure.",
		Type:       ebml.Binary,
	},
	BlockVirtualID: {
		Name:       "block_virtual",
		Definition: "A Block with no data. It **MUST** be stored in the stream at the place the real Block would be in display order.",
		Type:       ebml.Binary,
	},
	BlockAdditionsID: {
		Name:       "block_additions",
		Definition: "Contain additional blocks to complete the main one. An EBML parser that has no knowledge of the Block structure could still see and use/skip these data.",
		Type:       ebml.Master, Tag: BlockAdditions,
	},
	BlockDurationID: {
		Name:       "block_duration",
		Definition: "The duration of the Block, expressed in Track Ticks; see (#timestamp-ticks). The BlockDuration Element can be useful at the end of a Track to define the duration of the last frame (as there is no subsequent Block available), or when there is a break in a track like for subtitle tracks.",
		Type:       ebml.Uinteger,
	},
	ReferencePriorityID: {
		Name:       "reference_priority",
		Definition: "This frame is referenced and has the specified cache priority. In cache only a frame of the same or higher priority can replace this frame. A value of 0 means the frame is not referenced.",
		Type:       ebml.Uinteger,
	},
	ReferenceBlockID: {
		Name:       "reference_block",
		Definition: "A timestamp value, relative to the timestamp of the Block in this BlockGroup, expressed in Track Ticks; see (#timestamp-ticks). This is used to reference other frames necessary to decode this frame. The relative value **SHOULD** correspond to a valid `Block` this `Block` depends on. Historically Matroska Writer didn't write the actual `Block(s)` this `Block` depends on, but *some* `Block` in the past.  The value \"0\" **MAY** also be used to signify this `Block` cannot be decoded on its own, but without knownledge of which `Block` is necessary. In this case, other `ReferenceBlock` **MUST NOT** be found in the same `BlockGroup`.  If the `BlockGroup` doesn't have any `ReferenceBlock` element, then the `Block` it contains can be decoded without using any other `Block` data.",
		Type:       ebml.Integer,
	},
	ReferenceVirtualID: {
		Name:       "reference_virtual",
		Definition: "The Segment Position of the data that would otherwise be in position of the virtual block.",
		Type:       ebml.Integer,
	},
	CodecStateID: {
		Name:       "codec_state",
		Definition: "The new codec state to use. Data interpretation is private to the codec. This information **SHOULD** always be referenced by a seek entry.",
		Type:       ebml.Binary,
	},
	DiscardPaddingID: {
		Name:       "discard_padding",
		Definition: "Duration of the silent data added to the Block, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks) (padding at the end of the Block for positive value, at the beginning of the Block for negative value). The duration of DiscardPadding is not calculated in the duration of the TrackEntry and **SHOULD** be discarded during playback.",
		Type:       ebml.Integer,
	},
	SlicesID: {
		Name:       "slices",
		Definition: "Contains slices description.",
		Type:       ebml.Master, Tag: Slices,
	},
	ReferenceFrameID: {
		Name:       "reference_frame",
		Definition: "Contains information about the last reference frame. See [@?DivXTrickTrack].",
		Type:       ebml.Master, Tag: ReferenceFrame,
	},
}

var BlockAdditions = ebml.Tag{
	BlockMoreID: {
		Name:       "block_more",
		Definition: "Contain the BlockAdditional and some parameters.",
		Type:       ebml.Master, Tag: BlockMore,
	},
}

var BlockMore = ebml.Tag{
	BlockAddIDID: {
		Name:       "block_add_id",
		Definition: "An ID to identify the BlockAdditional level. If BlockAddIDType of the corresponding block is 0, this value is also the value of BlockAddIDType for the meaning of the content of BlockAdditional.",
		Type:       ebml.Uinteger,
	},
	BlockAdditionalID: {
		Name:       "block_additional",
		Definition: "Interpreted by the codec as it wishes (using the BlockAddID).",
		Type:       ebml.Binary,
	},
}

var Slices = ebml.Tag{
	TimeSliceID: {
		Name:       "time_slice",
		Definition: "Contains extra time information about the data contained in the Block. Being able to interpret this Element is not **REQUIRED** for playback.",
		Type:       ebml.Master, Tag: TimeSlice,
	},
}

var TimeSlice = ebml.Tag{
	LaceNumberID: {
		Name:       "lace_number",
		Definition: "The reverse number of the frame in the lace (0 is the last frame, 1 is the next to last, etc). Being able to interpret this Element is not **REQUIRED** for playback.",
		Type:       ebml.Uinteger,
	},
	FrameNumberID: {
		Name:       "frame_number",
		Definition: "The number of the frame to generate from this lace with this delay (allow you to generate many frames from the same Block/Frame).",
		Type:       ebml.Uinteger,
	},
	BlockAdditionIDID: {
		Name:       "block_addition_id",
		Definition: "The ID of the BlockAdditional Element (0 is the main Block).",
		Type:       ebml.Uinteger,
	},
	DelayID: {
		Name:       "delay",
		Definition: "The delay to apply to the Element, expressed in Track Ticks; see (#timestamp-ticks).",
		Type:       ebml.Uinteger,
	},
	SliceDurationID: {
		Name:       "slice_duration",
		Definition: "The duration to apply to the Element, expressed in Track Ticks; see (#timestamp-ticks).",
		Type:       ebml.Uinteger,
	},
}

var ReferenceFrame = ebml.Tag{
	ReferenceOffsetID: {
		Name:       "reference_offset",
		Definition: "The relative offset, in bytes, from the previous BlockGroup element for this Smooth FF/RW video track to the containing BlockGroup element. See [@?DivXTrickTrack].",
		Type:       ebml.Uinteger,
	},
	ReferenceTimestampID: {
		Name:       "reference_timestamp",
		Definition: "The timestamp of the BlockGroup pointed to by ReferenceOffset, expressed in Track Ticks; see (#timestamp-ticks). See [@?DivXTrickTrack].",
		Type:       ebml.Uinteger,
	},
}

var Tracks = ebml.Tag{
	TrackEntryID: {
		Name:       "track_entry",
		Definition: "Describes a track with all Elements.",
		Type:       ebml.Master, Tag: TrackEntry,
	},
}

var TrackEntry = ebml.Tag{
	TrackNumberID: {
		Name:       "track_number",
		Definition: "The track number as used in the Block Header (using more than 127 tracks is not encouraged, though the design allows an unlimited number).",
		Type:       ebml.Uinteger,
	},
	TrackUIDID: {
		Name:       "track_uid",
		Definition: "A unique ID to identify the Track.",
		Type:       ebml.Uinteger,
	},
	TrackTypeID: {
		Name:       "track_type",
		Definition: "The `TrackType` defines the type of each frame found in the Track. The value **SHOULD** be stored on 1 octet.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			1: {
				Sym:         "video",
				Description: "An image.",
			},
			2: {
				Sym:         "audio",
				Description: "Audio samples.",
			},
			3: {
				Sym:         "complex",
				Description: "A mix of different other TrackType. The codec needs to define how the `Matroska Player` should interpret such data.",
			},
			16: {
				Sym:         "logo",
				Description: "An image to be rendered over the video track(s).",
			},
			17: {
				Sym:         "subtitle",
				Description: "Subtitle or closed caption data to be rendered over the video track(s).",
			},
			18: {
				Sym:         "buttons",
				Description: "Interactive button(s) to be rendered over the video track(s).",
			},
			32: {
				Sym:         "control",
				Description: "Metadata used to control the player of the `Matroska Player`.",
			},
			33: {
				Sym:         "metadata",
				Description: "Timed metadata that can be passed on to the `Matroska Player`.",
			},
		},
	},
	FlagEnabledID: {
		Name:       "flag_enabled",
		Definition: "Set to 1 if the track is usable. It is possible to turn a not usable track into a usable track using chapter codecs or control tracks.",
		Type:       ebml.Uinteger,
	},
	FlagDefaultID: {
		Name:       "flag_default",
		Definition: "Set if that track (audio, video or subs) **SHOULD** be eligible for automatic selection by the player; see (#default-track-selection) for more details.",
		Type:       ebml.Uinteger,
	},
	FlagForcedID: {
		Name:       "flag_forced",
		Definition: "Applies only to subtitles. Set if that track **SHOULD** be eligible for automatic selection by the player if it matches the user's language preference, even if the user's preferences would normally not enable subtitles with the selected audio track; this can be used for tracks containing only translations of foreign-language audio or onscreen text. See (#default-track-selection) for more details.",
		Type:       ebml.Uinteger,
	},
	FlagHearingImpairedID: {
		Name:       "flag_hearing_impaired",
		Definition: "Set to 1 if that track is suitable for users with hearing impairments, set to 0 if it is unsuitable for users with hearing impairments.",
		Type:       ebml.Uinteger,
	},
	FlagVisualImpairedID: {
		Name:       "flag_visual_impaired",
		Definition: "Set to 1 if that track is suitable for users with visual impairments, set to 0 if it is unsuitable for users with visual impairments.",
		Type:       ebml.Uinteger,
	},
	FlagTextDescriptionsID: {
		Name:       "flag_text_descriptions",
		Definition: "Set to 1 if that track contains textual descriptions of video content, set to 0 if that track does not contain textual descriptions of video content.",
		Type:       ebml.Uinteger,
	},
	FlagOriginalID: {
		Name:       "flag_original",
		Definition: "Set to 1 if that track is in the content's original language, set to 0 if it is a translation.",
		Type:       ebml.Uinteger,
	},
	FlagCommentaryID: {
		Name:       "flag_commentary",
		Definition: "Set to 1 if that track contains commentary, set to 0 if it does not contain commentary.",
		Type:       ebml.Uinteger,
	},
	FlagLacingID: {
		Name:       "flag_lacing",
		Definition: "Set to 1 if the track **MAY** contain blocks using lacing. When set to 0 all blocks **MUST** have their lacing flags set to No lacing; see (#block-lacing) on Block Lacing.",
		Type:       ebml.Uinteger,
	},
	MinCacheID: {
		Name:       "min_cache",
		Definition: "The minimum number of frames a player **SHOULD** be able to cache during playback. If set to 0, the reference pseudo-cache system is not used.",
		Type:       ebml.Uinteger,
	},
	MaxCacheID: {
		Name:       "max_cache",
		Definition: "The maximum cache size necessary to store referenced frames in and the current frame. 0 means no cache is needed.",
		Type:       ebml.Uinteger,
	},
	DefaultDurationID: {
		Name:       "default_duration",
		Definition: "Number of nanoseconds per frame, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks) (frame in the Matroska sense -- one Element put into a (Simple)Block).",
		Type:       ebml.Uinteger,
	},
	DefaultDecodedFieldDurationID: {
		Name:       "default_decoded_field_duration",
		Definition: "The period between two successive fields at the output of the decoding process, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks). see (#defaultdecodedfieldduration) for more information",
		Type:       ebml.Uinteger,
	},
	TrackTimestampScaleID: {
		Name:       "track_timestamp_scale",
		Definition: "DEPRECATED, DO NOT USE. The scale to apply on this track to work at normal speed in relation with other tracks (mostly used to adjust video speed when the audio length differs).",
		Type:       ebml.Float,
	},
	TrackOffsetID: {
		Name:       "track_offset",
		Definition: "A value to add to the Block's Timestamp, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks). This can be used to adjust the playback offset of a track.",
		Type:       ebml.Integer,
	},
	MaxBlockAdditionIDID: {
		Name:       "max_block_addition_id",
		Definition: "The maximum value of BlockAddID ((#blockaddid-element)). A value 0 means there is no BlockAdditions ((#blockadditions-element)) for this track.",
		Type:       ebml.Uinteger,
	},
	BlockAdditionMappingID: {
		Name:       "block_addition_mapping",
		Definition: "Contains elements that extend the track format, by adding content either to each frame, with BlockAddID ((#blockaddid-element)), or to the track as a whole with BlockAddIDExtraData.",
		Type:       ebml.Master, Tag: BlockAdditionMapping,
	},
	NameID: {
		Name:       "name",
		Definition: "A human-readable track name.",
		Type:       ebml.UTF8,
	},
	LanguageID: {
		Name:       "language",
		Definition: "Specifies the language of the track in the Matroska languages form; see (#language-codes) on language codes. This Element **MUST** be ignored if the LanguageIETF Element is used in the same TrackEntry.",
		Type:       ebml.String,
	},
	LanguageIETFID: {
		Name:       "language_ietf",
		Definition: "Specifies the language of the track according to [@!BCP47] and using the IANA Language Subtag Registry [@!IANALangRegistry]. If this Element is used, then any Language Elements used in the same TrackEntry **MUST** be ignored.",
		Type:       ebml.String,
	},
	CodecIDID: {
		Name:       "codec_id",
		Definition: "An ID corresponding to the codec, see [@!MatroskaCodec] for more info.",
		Type:       ebml.String,
	},
	CodecPrivateID: {
		Name:       "codec_private",
		Definition: "Private data only known to the codec.",
		Type:       ebml.Binary,
	},
	CodecNameID: {
		Name:       "codec_name",
		Definition: "A human-readable string specifying the codec.",
		Type:       ebml.UTF8,
	},
	AttachmentLinkID: {
		Name:       "attachment_link",
		Definition: "The UID of an attachment that is used by this codec.",
		Type:       ebml.Uinteger,
	},
	CodecSettingsID: {
		Name:       "codec_settings",
		Definition: "A string describing the encoding setting used.",
		Type:       ebml.UTF8,
	},
	CodecInfoURLID: {
		Name:       "codec_info_url",
		Definition: "A URL to find information about the codec used.",
		Type:       ebml.String,
	},
	CodecDownloadURLID: {
		Name:       "codec_download_url",
		Definition: "A URL to download about the codec used.",
		Type:       ebml.String,
	},
	CodecDecodeAllID: {
		Name:       "codec_decode_all",
		Definition: "Set to 1 if the codec can decode potentially damaged data.",
		Type:       ebml.Uinteger,
	},
	TrackOverlayID: {
		Name:       "track_overlay",
		Definition: "Specify that this track is an overlay track for the Track specified (in the u-integer). That means when this track has a gap, see (#silenttracks-element) on SilentTracks, the overlay track **SHOULD** be used instead. The order of multiple TrackOverlay matters, the first one is the one that **SHOULD** be used. If not found it **SHOULD** be the second, etc.",
		Type:       ebml.Uinteger,
	},
	CodecDelayID: {
		Name:       "codec_delay",
		Definition: "CodecDelay is The codec-built-in delay, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks). It represents the amount of codec samples that will be discarded by the decoder during playback. This timestamp value **MUST** be subtracted from each frame timestamp in order to get the timestamp that will be actually played. The value **SHOULD** be small so the muxing of tracks with the same actual timestamp are in the same Cluster.",
		Type:       ebml.Uinteger,
	},
	SeekPreRollID: {
		Name:       "seek_pre_roll",
		Definition: "After a discontinuity, SeekPreRoll is the duration of the data the decoder **MUST** decode before the decoded data is valid, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks).",
		Type:       ebml.Uinteger,
	},
	TrackTranslateID: {
		Name:       "track_translate",
		Definition: "The mapping between this `TrackEntry` and a track value in the given Chapter Codec.",
		Type:       ebml.Master, Tag: TrackTranslate,
	},
	VideoID: {
		Name:       "video",
		Definition: "Video settings.",
		Type:       ebml.Master, Tag: Video,
	},
	AudioID: {
		Name:       "audio",
		Definition: "Audio settings.",
		Type:       ebml.Master, Tag: Audio,
	},
	TrackOperationID: {
		Name:       "track_operation",
		Definition: "Operation that needs to be applied on tracks to create this virtual track. For more details look at (#track-operation).",
		Type:       ebml.Master, Tag: TrackOperation,
	},
	TrickTrackUIDID: {
		Name:       "trick_track_uid",
		Definition: "The TrackUID of the Smooth FF/RW video in the paired EBML structure corresponding to this video track. See [@?DivXTrickTrack].",
		Type:       ebml.Uinteger,
	},
	TrickTrackSegmentUIDID: {
		Name:       "trick_track_segment_uid",
		Definition: "The SegmentUID of the Segment containing the track identified by TrickTrackUID. See [@?DivXTrickTrack].",
		Type:       ebml.Binary,
	},
	TrickTrackFlagID: {
		Name:       "trick_track_flag",
		Definition: "Set to 1 if this video track is a Smooth FF/RW track. If set to 1, MasterTrackUID and MasterTrackSegUID should must be present and BlockGroups for this track must contain ReferenceFrame structures. Otherwise, TrickTrackUID and TrickTrackSegUID must be present if this track has a corresponding Smooth FF/RW track. See [@?DivXTrickTrack].",
		Type:       ebml.Uinteger,
	},
	TrickMasterTrackUIDID: {
		Name:       "trick_master_track_uid",
		Definition: "The TrackUID of the video track in the paired EBML structure that corresponds to this Smooth FF/RW track. See [@?DivXTrickTrack].",
		Type:       ebml.Uinteger,
	},
	TrickMasterTrackSegmentUIDID: {
		Name:       "trick_master_track_segment_uid",
		Definition: "The SegmentUID of the Segment containing the track identified by MasterTrackUID. See [@?DivXTrickTrack].",
		Type:       ebml.Binary,
	},
	ContentEncodingsID: {
		Name:       "content_encodings",
		Definition: "Settings for several content encoding mechanisms like compression or encryption.",
		Type:       ebml.Master, Tag: ContentEncodings,
	},
}

var BlockAdditionMapping = ebml.Tag{
	BlockAddIDValueID: {
		Name:       "block_add_idvalue",
		Definition: "If the track format extension needs content beside frames, the value refers to the BlockAddID ((#blockaddid-element)), value being described. To keep MaxBlockAdditionID as low as possible, small values **SHOULD** be used.",
		Type:       ebml.Uinteger,
	},
	BlockAddIDNameID: {
		Name:       "block_add_idname",
		Definition: "A human-friendly name describing the type of BlockAdditional data, as defined by the associated Block Additional Mapping.",
		Type:       ebml.String,
	},
	BlockAddIDTypeID: {
		Name:       "block_add_idtype",
		Definition: "Stores the registered identifier of the Block Additional Mapping to define how the BlockAdditional data should be handled.",
		Type:       ebml.Uinteger,
	},
	BlockAddIDExtraDataID: {
		Name:       "block_add_idextra_data",
		Definition: "Extra binary data that the BlockAddIDType can use to interpret the BlockAdditional data. The interpretation of the binary data depends on the BlockAddIDType value and the corresponding Block Additional Mapping.",
		Type:       ebml.Binary,
	},
}

var TrackTranslate = ebml.Tag{
	TrackTranslateTrackIDID: {
		Name:       "track_translate_track_id",
		Definition: "The binary value used to represent this `TrackEntry` in the chapter codec data. The format depends on the `ChapProcessCodecID` used; see (#chapprocesscodecid-element).",
		Type:       ebml.Binary,
	},
	TrackTranslateCodecID: {
		Name:       "track_translate_codec",
		Definition: "This `TrackTranslate` applies to this chapter codec of the given chapter edition(s); see (#chapprocesscodecid-element).",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym:         "matroska_script",
				Description: "Chapter commands using the Matroska Script codec.",
			},
			1: {
				Sym:         "dvd_menu",
				Description: "Chapter commands using the DVD-like codec.",
			},
		},
	},
	TrackTranslateEditionUIDID: {
		Name:       "track_translate_edition_uid",
		Definition: "Specify a chapter edition UID on which this `TrackTranslate` applies.",
		Type:       ebml.Uinteger,
	},
}

var Video = ebml.Tag{
	FlagInterlacedID: {
		Name:       "flag_interlaced",
		Definition: "Specify whether the video frames in this track are interlaced or not.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym:         "undetermined",
				Description: "Unknown status.",
			},
			1: {
				Sym:         "interlaced",
				Description: "Interlaced frames.",
			},
			2: {
				Sym:         "progressive",
				Description: "No interlacing.",
			},
		},
	},
	FieldOrderID: {
		Name:       "field_order",
		Definition: "Specify the field ordering of video frames in this track.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym:         "progressive",
				Description: "Interlaced frames.",
			},
			1: {
				Sym:         "tff",
				Description: "Top field displayed first. Top field stored first.",
			},
			2: {
				Sym:         "undetermined",
				Description: "Unknown field order.",
			},
			6: {
				Sym:         "bff",
				Description: "Bottom field displayed first. Bottom field stored first.",
			},
			9: {
				Description: "Top field displayed first. Fields are interleaved in storage with the top line of the top field stored first.",
			},
			14: {
				Description: "Bottom field displayed first. Fields are interleaved in storage with the top line of the top field stored first.",
			},
		},
	},
	StereoModeID: {
		Name:       "stereo_mode",
		Definition: "Stereo-3D video mode. There are some more details in (#multi-planar-and-3d-videos).",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "mono",
			},
			1: {
				Description: "side by side (left eye first)",
			},
			2: {
				Description: "top - bottom (right eye is first)",
			},
			3: {
				Description: "top - bottom (left eye is first)",
			},
			4: {
				Description: "checkboard (right eye is first)",
			},
			5: {
				Description: "checkboard (left eye is first)",
			},
			6: {
				Description: "row interleaved (right eye is first)",
			},
			7: {
				Description: "row interleaved (left eye is first)",
			},
			8: {
				Description: "column interleaved (right eye is first)",
			},
			9: {
				Description: "column interleaved (left eye is first)",
			},
			10: {
				Description: "anaglyph (cyan/red)",
			},
			11: {
				Description: "side by side (right eye first)",
			},
			12: {
				Description: "anaglyph (green/magenta)",
			},
			13: {
				Description: "both eyes laced in one Block (left eye is first)",
			},
			14: {
				Description: "both eyes laced in one Block (right eye is first)",
			},
		},
	},
	AlphaModeID: {
		Name:       "alpha_mode",
		Definition: "Indicate whether the BlockAdditional Element with BlockAddID of \"1\" contains Alpha data, as defined by to the Codec Mapping for the `CodecID`. Undefined values **SHOULD NOT** be used as the behavior of known implementations is different (considered either as 0 or 1).",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym:         "none",
				Description: "The BlockAdditional Element with BlockAddID of \"1\" does not exist or **SHOULD NOT** be considered as containing such data.",
			},
			1: {
				Sym:         "present",
				Description: "The BlockAdditional Element with BlockAddID of \"1\" contains alpha channel data.",
			},
		},
	},
	OldStereoModeID: {
		Name:       "old_stereo_mode",
		Definition: "DEPRECATED, DO NOT USE. Bogus StereoMode value used in old versions of libmatroska.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "mono",
			},
			1: {
				Sym: "right_eye",
			},
			2: {
				Sym: "left_eye",
			},
			3: {
				Sym: "both_eyes",
			},
		},
	},
	PixelWidthID: {
		Name:       "pixel_width",
		Definition: "Width of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	PixelHeightID: {
		Name:       "pixel_height",
		Definition: "Height of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	PixelCropBottomID: {
		Name:       "pixel_crop_bottom",
		Definition: "The number of video pixels to remove at the bottom of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropTopID: {
		Name:       "pixel_crop_top",
		Definition: "The number of video pixels to remove at the top of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropLeftID: {
		Name:       "pixel_crop_left",
		Definition: "The number of video pixels to remove on the left of the image.",
		Type:       ebml.Uinteger,
	},
	PixelCropRightID: {
		Name:       "pixel_crop_right",
		Definition: "The number of video pixels to remove on the right of the image.",
		Type:       ebml.Uinteger,
	},
	DisplayWidthID: {
		Name:       "display_width",
		Definition: "Width of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	DisplayHeightID: {
		Name:       "display_height",
		Definition: "Height of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	DisplayUnitID: {
		Name:       "display_unit",
		Definition: "How DisplayWidth & DisplayHeight are interpreted.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "pixels",
			},
			1: {
				Sym: "centimeters",
			},
			2: {
				Sym: "inches",
			},
			3: {
				Sym: "display_aspect_ratio",
			},
			4: {
				Sym: "unknown",
			},
		},
	},
	AspectRatioTypeID: {
		Name:       "aspect_ratio_type",
		Definition: "Specify the possible modifications to the aspect ratio.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "free_resizing",
			},
			1: {
				Sym: "keep_aspect_ratio",
			},
			2: {
				Sym: "fixed",
			},
		},
	},
	UncompressedFourCCID: {
		Name:       "uncompressed_four_cc",
		Definition: "Specify the uncompressed pixel format used for the Track's data as a FourCC. This value is similar in scope to the biCompression value of AVI's `BITMAPINFO` [@?AVIFormat]. See the YUV video formats [@?FourCC-YUV] and RGB video formats [@?FourCC-RGB] for common values.",
		Type:       ebml.Binary,
	},
	GammaValueID: {
		Name:       "gamma_value",
		Definition: "Gamma Value.",
		Type:       ebml.Float,
	},
	FrameRateID: {
		Name:       "frame_rate",
		Definition: "Number of frames per second. This value is Informational only. It is intended for constant frame rate streams, and **SHOULD NOT** be used for a variable frame rate TrackEntry.",
		Type:       ebml.Float,
	},
	ColourID: {
		Name:       "colour",
		Definition: "Settings describing the colour format.",
		Type:       ebml.Master, Tag: Colour,
	},
	ProjectionID: {
		Name:       "projection",
		Definition: "Describes the video projection details. Used to render spherical, VR videos or flipping videos horizontally/vertically.",
		Type:       ebml.Master, Tag: Projection,
	},
}

var Colour = ebml.Tag{
	MatrixCoefficientsID: {
		Name:       "matrix_coefficients",
		Definition: "The Matrix Coefficients of the video used to derive luma and chroma values from red, green, and blue color primaries. For clarity, the value and meanings for MatrixCoefficients are adopted from Table 4 of ISO/IEC 23001-8:2016 or ITU-T H.273.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "identity",
			},
			1: {
				Sym: "itu_r_bt_709",
			},
			2: {
				Sym: "unspecified",
			},
			3: {
				Sym: "reserved",
			},
			4: {
				Sym: "us_fcc_73_682",
			},
			5: {
				Sym: "itu_r_bt_470bg",
			},
			6: {
				Sym: "smpte_170m",
			},
			7: {
				Sym: "smpte_240m",
			},
			8: {
				Sym: "ycocg",
			},
			9: {
				Sym: "bt2020_non_constant_luminance",
			},
			10: {
				Sym: "bt2020_constant_luminance",
			},
			11: {
				Sym: "smpte_st_2085",
			},
			12: {
				Sym: "chroma_derived_non_constant_luminance",
			},
			13: {
				Sym: "chroma_derived_constant_luminance",
			},
			14: {
				Sym: "itu_r_bt_2100_0",
			},
		},
	},
	BitsPerChannelID: {
		Name:       "bits_per_channel",
		Definition: "Number of decoded bits per channel. A value of 0 indicates that the BitsPerChannel is unspecified.",
		Type:       ebml.Uinteger,
	},
	ChromaSubsamplingHorzID: {
		Name:       "chroma_subsampling_horz",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed horizontally. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingHorz **SHOULD** be set to 1.",
		Type:       ebml.Uinteger,
	},
	ChromaSubsamplingVertID: {
		Name:       "chroma_subsampling_vert",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed vertically. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingVert **SHOULD** be set to 1.",
		Type:       ebml.Uinteger,
	},
	CbSubsamplingHorzID: {
		Name:       "cb_subsampling_horz",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed horizontally. This is additive with ChromaSubsamplingHorz. Example: For video with 4:2:1 chroma subsampling, the ChromaSubsamplingHorz **SHOULD** be set to 1 and CbSubsamplingHorz **SHOULD** be set to 1.",
		Type:       ebml.Uinteger,
	},
	CbSubsamplingVertID: {
		Name:       "cb_subsampling_vert",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed vertically. This is additive with ChromaSubsamplingVert.",
		Type:       ebml.Uinteger,
	},
	ChromaSitingHorzID: {
		Name:       "chroma_siting_horz",
		Definition: "How chroma is subsampled horizontally.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "unspecified",
			},
			1: {
				Sym: "left_collocated",
			},
			2: {
				Sym: "half",
			},
		},
	},
	ChromaSitingVertID: {
		Name:       "chroma_siting_vert",
		Definition: "How chroma is subsampled vertically.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "unspecified",
			},
			1: {
				Sym: "top_collocated",
			},
			2: {
				Sym: "half",
			},
		},
	},
	RangeID: {
		Name:       "range",
		Definition: "Clipping of the color ranges.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "unspecified",
			},
			1: {
				Sym: "broadcast_range",
			},
			2: {
				Description: "full range (no clipping)",
			},
			3: {
				Sym: "defined_by_matrixcoefficients_transfercharacteristics",
			},
		},
	},
	TransferCharacteristicsID: {
		Name:       "transfer_characteristics",
		Definition: "The transfer characteristics of the video. For clarity, the value and meanings for TransferCharacteristics are adopted from Table 3 of ISO/IEC 23091-4 or ITU-T H.273.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "reserved",
			},
			1: {
				Sym: "itu_r_bt_709",
			},
			2: {
				Sym: "unspecified",
			},
			3: {
				Sym: "reserved",
			},
			4: {
				Sym: "gamma_2_2_curve_bt_470m",
			},
			5: {
				Sym: "gamma_2_8_curve_bt_470bg",
			},
			6: {
				Sym: "smpte_170m",
			},
			7: {
				Sym: "smpte_240m",
			},
			8: {
				Sym: "linear",
			},
			9: {
				Sym: "log",
			},
			10: {
				Sym: "log_sqrt",
			},
			11: {
				Sym: "iec_61966_2_4",
			},
			12: {
				Sym: "itu_r_bt_1361_extended_colour_gamut",
			},
			13: {
				Sym: "iec_61966_2_1",
			},
			14: {
				Sym: "itu_r_bt_2020_10_bit",
			},
			15: {
				Sym: "itu_r_bt_2020_12_bit",
			},
			16: {
				Sym: "itu_r_bt_2100_perceptual_quantization",
			},
			17: {
				Sym: "smpte_st_428_1",
			},
			18: {
				Description: "ARIB STD-B67 (HLG)",
			},
		},
	},
	PrimariesID: {
		Name:       "primaries",
		Definition: "The colour primaries of the video. For clarity, the value and meanings for Primaries are adopted from Table 2 of ISO/IEC 23091-4 or ITU-T H.273.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "reserved",
			},
			1: {
				Sym: "itu_r_bt_709",
			},
			2: {
				Sym: "unspecified",
			},
			3: {
				Sym: "reserved",
			},
			4: {
				Sym: "itu_r_bt_470m",
			},
			5: {
				Sym: "itu_r_bt_470bg_bt_601_625",
			},
			6: {
				Sym: "itu_r_bt_601_525_smpte_170m",
			},
			7: {
				Sym: "smpte_240m",
			},
			8: {
				Sym: "film",
			},
			9: {
				Sym: "itu_r_bt_2020",
			},
			10: {
				Sym: "smpte_st_428_1",
			},
			11: {
				Sym: "smpte_rp_432_2",
			},
			12: {
				Sym: "smpte_eg_432_2",
			},
			22: {
				Sym: "ebu_tech_3213_e_jedec_p22_phosphors",
			},
		},
	},
	MaxCLLID: {
		Name:       "max_cll",
		Definition: "Maximum brightness of a single pixel (Maximum Content Light Level) in candelas per square meter (cd/m^2^).",
		Type:       ebml.Uinteger,
	},
	MaxFALLID: {
		Name:       "max_fall",
		Definition: "Maximum brightness of a single full frame (Maximum Frame-Average Light Level) in candelas per square meter (cd/m^2^).",
		Type:       ebml.Uinteger,
	},
	MasteringMetadataID: {
		Name:       "mastering_metadata",
		Definition: "SMPTE 2086 mastering data.",
		Type:       ebml.Master, Tag: MasteringMetadata,
	},
}

var MasteringMetadata = ebml.Tag{
	PrimaryRChromaticityXID: {
		Name:       "primary_rchromaticity_x",
		Definition: "Red X chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryRChromaticityYID: {
		Name:       "primary_rchromaticity_y",
		Definition: "Red Y chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryGChromaticityXID: {
		Name:       "primary_gchromaticity_x",
		Definition: "Green X chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryGChromaticityYID: {
		Name:       "primary_gchromaticity_y",
		Definition: "Green Y chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryBChromaticityXID: {
		Name:       "primary_bchromaticity_x",
		Definition: "Blue X chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	PrimaryBChromaticityYID: {
		Name:       "primary_bchromaticity_y",
		Definition: "Blue Y chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	WhitePointChromaticityXID: {
		Name:       "white_point_chromaticity_x",
		Definition: "White X chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	WhitePointChromaticityYID: {
		Name:       "white_point_chromaticity_y",
		Definition: "White Y chromaticity coordinate, as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	LuminanceMaxID: {
		Name:       "luminance_max",
		Definition: "Maximum luminance. Represented in candelas per square meter (cd/m^2^).",
		Type:       ebml.Float,
	},
	LuminanceMinID: {
		Name:       "luminance_min",
		Definition: "Minimum luminance. Represented in candelas per square meter (cd/m^2^).",
		Type:       ebml.Float,
	},
}

var Projection = ebml.Tag{
	ProjectionTypeID: {
		Name:       "projection_type",
		Definition: "Describes the projection used for this video track.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "rectangular",
			},
			1: {
				Sym: "equirectangular",
			},
			2: {
				Sym: "cubemap",
			},
			3: {
				Sym: "mesh",
			},
		},
	},
	ProjectionPrivateID: {
		Name:       "projection_private",
		Definition: "Private data that only applies to a specific projection.  *  If `ProjectionType` equals 0 (Rectangular),      then this element must not be present. *  If `ProjectionType` equals 1 (Equirectangular), then this element must be present and contain the same binary data that would be stored inside       an ISOBMFF Equirectangular Projection Box ('equi'). *  If `ProjectionType` equals 2 (Cubemap), then this element must be present and contain the same binary data that would be stored       inside an ISOBMFF Cubemap Projection Box ('cbmp'). *  If `ProjectionType` equals 3 (Mesh), then this element must be present and contain the same binary data that would be stored inside        an ISOBMFF Mesh Projection Box ('mshp').",
		Type:       ebml.Binary,
	},
	ProjectionPoseYawID: {
		Name:       "projection_pose_yaw",
		Definition: "Specifies a yaw rotation to the projection.  Value represents a clockwise rotation, in degrees, around the up vector. This rotation must be applied before any `ProjectionPosePitch` or `ProjectionPoseRoll` rotations. The value of this element **MUST** be in the -180 to 180 degree range, both included.  Setting `ProjectionPoseYaw` to 180 or -180 degrees, with the `ProjectionPoseRoll` and `ProjectionPosePitch` set to 0 degrees flips the image horizontally.",
		Type:       ebml.Float,
	},
	ProjectionPosePitchID: {
		Name:       "projection_pose_pitch",
		Definition: "Specifies a pitch rotation to the projection.  Value represents a counter-clockwise rotation, in degrees, around the right vector. This rotation must be applied after the `ProjectionPoseYaw` rotation and before the `ProjectionPoseRoll` rotation. The value of this element **MUST** be in the -90 to 90 degree range, both included.",
		Type:       ebml.Float,
	},
	ProjectionPoseRollID: {
		Name:       "projection_pose_roll",
		Definition: "Specifies a roll rotation to the projection.  Value represents a counter-clockwise rotation, in degrees, around the forward vector. This rotation must be applied after the `ProjectionPoseYaw` and `ProjectionPosePitch` rotations. The value of this element **MUST** be in the -180 to 180 degree range, both included.  Setting `ProjectionPoseRoll` to 180 or -180 degrees, the `ProjectionPoseYaw` to 180 or -180 degrees with `ProjectionPosePitch` set to 0 degrees flips the image vertically.  Setting `ProjectionPoseRoll` to 180 or -180 degrees, with the `ProjectionPoseYaw` and `ProjectionPosePitch` set to 0 degrees flips the image horizontally and vertically.",
		Type:       ebml.Float,
	},
}

var Audio = ebml.Tag{
	SamplingFrequencyID: {
		Name:       "sampling_frequency",
		Definition: "Sampling frequency in Hz.",
		Type:       ebml.Float,
	},
	OutputSamplingFrequencyID: {
		Name:       "output_sampling_frequency",
		Definition: "Real output sampling frequency in Hz (used for SBR techniques).",
		Type:       ebml.Float,
	},
	ChannelsID: {
		Name:       "channels",
		Definition: "Numbers of channels in the track.",
		Type:       ebml.Uinteger,
	},
	ChannelPositionsID: {
		Name:       "channel_positions",
		Definition: "Table of horizontal angles for each successive channel.",
		Type:       ebml.Binary,
	},
	BitDepthID: {
		Name:       "bit_depth",
		Definition: "Bits per sample, mostly used for PCM.",
		Type:       ebml.Uinteger,
	},
}

var TrackOperation = ebml.Tag{
	TrackCombinePlanesID: {
		Name:       "track_combine_planes",
		Definition: "Contains the list of all video plane tracks that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: TrackCombinePlanes,
	},
	TrackJoinBlocksID: {
		Name:       "track_join_blocks",
		Definition: "Contains the list of all tracks whose Blocks need to be combined to create this virtual track",
		Type:       ebml.Master, Tag: TrackJoinBlocks,
	},
}

var TrackCombinePlanes = ebml.Tag{
	TrackPlaneID: {
		Name:       "track_plane",
		Definition: "Contains a video plane track that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: TrackPlane,
	},
}

var TrackPlane = ebml.Tag{
	TrackPlaneUIDID: {
		Name:       "track_plane_uid",
		Definition: "The trackUID number of the track representing the plane.",
		Type:       ebml.Uinteger,
	},
	TrackPlaneTypeID: {
		Name:       "track_plane_type",
		Definition: "The kind of plane this track corresponds to.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "left_eye",
			},
			1: {
				Sym: "right_eye",
			},
			2: {
				Sym: "background",
			},
		},
	},
}

var TrackJoinBlocks = ebml.Tag{
	TrackJoinUIDID: {
		Name:       "track_join_uid",
		Definition: "The trackUID number of a track whose blocks are used to create this virtual track.",
		Type:       ebml.Uinteger,
	},
}

var ContentEncodings = ebml.Tag{
	ContentEncodingID: {
		Name:       "content_encoding",
		Definition: "Settings for one content encoding like compression or encryption.",
		Type:       ebml.Master, Tag: ContentEncoding,
	},
}

var ContentEncoding = ebml.Tag{
	ContentEncodingOrderID: {
		Name:       "content_encoding_order",
		Definition: "Tell in which order to apply each `ContentEncoding` of the `ContentEncodings`. The decoder/demuxer **MUST** start with the `ContentEncoding` with the highest `ContentEncodingOrder` and work its way down to the `ContentEncoding` with the lowest `ContentEncodingOrder`. This value **MUST** be unique over for each `ContentEncoding` found in the `ContentEncodings` of this `TrackEntry`.",
		Type:       ebml.Uinteger,
	},
	ContentEncodingScopeID: {
		Name:       "content_encoding_scope",
		Definition: "A bit field that describes which Elements have been modified in this way. Values (big-endian) can be OR'ed.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			1: {
				Sym:         "block",
				Description: "All frame contents, excluding lacing data.",
			},
			2: {
				Sym:         "private",
				Description: "The track's private data.",
			},
			4: {
				Sym:         "next",
				Description: "The next ContentEncoding (next `ContentEncodingOrder`. Either the data inside `ContentCompression` and/or `ContentEncryption`).",
			},
		},
	},
	ContentEncodingTypeID: {
		Name:       "content_encoding_type",
		Definition: "A value describing what kind of transformation is applied.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "compression",
			},
			1: {
				Sym: "encryption",
			},
		},
	},
	ContentCompressionID: {
		Name:       "content_compression",
		Definition: "Settings describing the compression used. This Element **MUST** be present if the value of ContentEncodingType is 0 and absent otherwise. Each block **MUST** be decompressable even if no previous block is available in order not to prevent seeking.",
		Type:       ebml.Master, Tag: ContentCompression,
	},
	ContentEncryptionID: {
		Name:       "content_encryption",
		Definition: "Settings describing the encryption used. This Element **MUST** be present if the value of `ContentEncodingType` is 1 (encryption) and **MUST** be ignored otherwise.",
		Type:       ebml.Master, Tag: ContentEncryption,
	},
}

var ContentCompression = ebml.Tag{
	ContentCompAlgoID: {
		Name:       "content_comp_algo",
		Definition: "The compression algorithm used.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym:         "zlib",
				Description: "zlib compression [@!RFC1950].",
			},
			1: {
				Sym:         "bzlib",
				Description: "bzip2 compression [@!BZIP2], **SHOULD NOT** be used; see usage notes.",
			},
			2: {
				Sym:         "lzo1x",
				Description: "Lempel-Ziv-Oberhumer compression [@!LZO], **SHOULD NOT** be used; see usage notes.",
			},
			3: {
				Sym:         "header_stripping",
				Description: "Octets in `ContentCompSettings` ((#contentcompsettings-element)) have been stripped from each frame.",
			},
		},
	},
	ContentCompSettingsID: {
		Name:       "content_comp_settings",
		Definition: "Settings that might be needed by the decompressor. For Header Stripping (`ContentCompAlgo`=3), the bytes that were removed from the beginning of each frames of the track.",
		Type:       ebml.Binary,
	},
}

var ContentEncryption = ebml.Tag{
	ContentEncAlgoID: {
		Name:       "content_enc_algo",
		Definition: "The encryption algorithm used. The value \"0\" means that the contents have not been encrypted.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "not_encrypted",
			},
			1: {
				Sym:         "des",
				Description: "Data Encryption Standard (DES) [@!FIPS.46-3].",
			},
			2: {
				Sym:         "3des",
				Description: "Triple Data Encryption Algorithm [@!SP.800-67].",
			},
			3: {
				Sym:         "twofish",
				Description: "Twofish Encryption Algorithm [@!Twofish].",
			},
			4: {
				Sym:         "blowfish",
				Description: "Blowfish Encryption Algorithm [@!Blowfish].",
			},
			5: {
				Sym:         "aes",
				Description: "Advanced Encryption Standard (AES) [@!FIPS.197].",
			},
		},
	},
	ContentEncKeyIDID: {
		Name:       "content_enc_key_id",
		Definition: "For public key algorithms this is the ID of the public key the the data was encrypted with.",
		Type:       ebml.Binary,
	},
	ContentEncAESSettingsID: {
		Name:       "content_enc_aessettings",
		Definition: "Settings describing the encryption algorithm used. It **MUST** be ignored if `ContentEncAlgo` is not AES (5).",
		Type:       ebml.Master, Tag: ContentEncAESSettings,
	},
	ContentSignatureID: {
		Name:       "content_signature",
		Definition: "A cryptographic signature of the contents.",
		Type:       ebml.Binary,
	},
	ContentSigKeyIDID: {
		Name:       "content_sig_key_id",
		Definition: "This is the ID of the private key the data was signed with.",
		Type:       ebml.Binary,
	},
	ContentSigAlgoID: {
		Name:       "content_sig_algo",
		Definition: "The algorithm used for the signature.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "not_signed",
			},
			1: {
				Sym: "rsa",
			},
		},
	},
	ContentSigHashAlgoID: {
		Name:       "content_sig_hash_algo",
		Definition: "The hash algorithm used for the signature.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "not_signed",
			},
			1: {
				Sym: "sha1_160",
			},
			2: {
				Sym: "md5",
			},
		},
	},
}

var ContentEncAESSettings = ebml.Tag{
	AESSettingsCipherModeID: {
		Name:       "aessettings_cipher_mode",
		Definition: "The AES cipher mode used in the encryption. It **MUST** be ignored if `ContentEncAlgo` is not AES (5).",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			1: {
				Sym:         "aes_ctr",
				Description: "Counter [@!SP.800-38A].",
			},
			2: {
				Sym:         "aes_cbc",
				Description: "Cipher Block Chaining [@!SP.800-38A].",
			},
		},
	},
}

var Cues = ebml.Tag{
	CuePointID: {
		Name:       "cue_point",
		Definition: "Contains all information relative to a seek point in the Segment.",
		Type:       ebml.Master, Tag: CuePoint,
	},
}

var CuePoint = ebml.Tag{
	CueTimeID: {
		Name:       "cue_time",
		Definition: "Absolute timestamp of the seek point, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks).",
		Type:       ebml.Uinteger,
	},
	CueTrackPositionsID: {
		Name:       "cue_track_positions",
		Definition: "Contain positions for different tracks corresponding to the timestamp.",
		Type:       ebml.Master, Tag: CueTrackPositions,
	},
}

var CueTrackPositions = ebml.Tag{
	CueTrackID: {
		Name:       "cue_track",
		Definition: "The track for which a position is given.",
		Type:       ebml.Uinteger,
	},
	CueClusterPositionID: {
		Name:       "cue_cluster_position",
		Definition: "The Segment Position of the Cluster containing the associated Block.",
		Type:       ebml.Uinteger,
	},
	CueRelativePositionID: {
		Name:       "cue_relative_position",
		Definition: "The relative position inside the Cluster of the referenced SimpleBlock or BlockGroup with 0 being the first possible position for an Element inside that Cluster.",
		Type:       ebml.Uinteger,
	},
	CueDurationID: {
		Name:       "cue_duration",
		Definition: "The duration of the block, expressed in Segment Ticks which is based on TimestampScale; see (#timestamp-ticks). If missing, the track's DefaultDuration does not apply and no duration information is available in terms of the cues.",
		Type:       ebml.Uinteger,
	},
	CueBlockNumberID: {
		Name:       "cue_block_number",
		Definition: "Number of the Block in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	CueCodecStateID: {
		Name:       "cue_codec_state",
		Definition: "The Segment Position of the Codec State corresponding to this Cue Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
	CueReferenceID: {
		Name:       "cue_reference",
		Definition: "The Clusters containing the referenced Blocks.",
		Type:       ebml.Master, Tag: CueReference,
	},
}

var CueReference = ebml.Tag{
	CueRefTimeID: {
		Name:       "cue_ref_time",
		Definition: "Timestamp of the referenced Block, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks).",
		Type:       ebml.Uinteger,
	},
	CueRefClusterID: {
		Name:       "cue_ref_cluster",
		Definition: "The Segment Position of the Cluster containing the referenced Block.",
		Type:       ebml.Uinteger,
	},
	CueRefNumberID: {
		Name:       "cue_ref_number",
		Definition: "Number of the referenced Block of Track X in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	CueRefCodecStateID: {
		Name:       "cue_ref_codec_state",
		Definition: "The Segment Position of the Codec State corresponding to this referenced Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
}

var Attachments = ebml.Tag{
	AttachedFileID: {
		Name:       "attached_file",
		Definition: "An attached file.",
		Type:       ebml.Master, Tag: AttachedFile,
	},
}

var AttachedFile = ebml.Tag{
	FileDescriptionID: {
		Name:       "file_description",
		Definition: "A human-friendly name for the attached file.",
		Type:       ebml.UTF8,
	},
	FileNameID: {
		Name:       "file_name",
		Definition: "Filename of the attached file.",
		Type:       ebml.UTF8,
	},
	FileMimeTypeID: {
		Name:       "file_mime_type",
		Definition: "MIME type of the file.",
		Type:       ebml.String,
	},
	FileDataID: {
		Name:       "file_data",
		Definition: "The data of the file.",
		Type:       ebml.Binary,
	},
	FileUIDID: {
		Name:       "file_uid",
		Definition: "Unique ID representing the file, as random as possible.",
		Type:       ebml.Uinteger,
	},
	FileReferralID: {
		Name:       "file_referral",
		Definition: "A binary value that a track/codec can refer to when the attachment is needed.",
		Type:       ebml.Binary,
	},
	FileUsedStartTimeID: {
		Name:       "file_used_start_time",
		Definition: "The timestamp at which this optimized font attachment comes into context, expressed in Segment Ticks which is based on TimestampScale. See [@?DivXWorldFonts].",
		Type:       ebml.Uinteger,
	},
	FileUsedEndTimeID: {
		Name:       "file_used_end_time",
		Definition: "The timestamp at which this optimized font attachment goes out of context, expressed in Segment Ticks which is based on TimestampScale. See [@?DivXWorldFonts].",
		Type:       ebml.Uinteger,
	},
}

var Chapters = ebml.Tag{
	EditionEntryID: {
		Name:       "edition_entry",
		Definition: "Contains all information about a Segment edition.",
		Type:       ebml.Master, Tag: EditionEntry,
	},
}

var EditionEntry = ebml.Tag{
	EditionUIDID: {
		Name:       "edition_uid",
		Definition: "A unique ID to identify the edition. It's useful for tagging an edition.",
		Type:       ebml.Uinteger,
	},
	EditionFlagHiddenID: {
		Name:       "edition_flag_hidden",
		Definition: "Set to 1 if an edition is hidden. Hidden editions **SHOULD NOT** be available to the user interface (but still to Control Tracks; see (#chapter-flags) on Chapter flags).",
		Type:       ebml.Uinteger,
	},
	EditionFlagDefaultID: {
		Name:       "edition_flag_default",
		Definition: "Set to 1 if the edition **SHOULD** be used as the default one.",
		Type:       ebml.Uinteger,
	},
	EditionFlagOrderedID: {
		Name:       "edition_flag_ordered",
		Definition: "Set to 1 if the chapters can be defined multiple times and the order to play them is enforced; see (#editionflagordered).",
		Type:       ebml.Uinteger,
	},
	ChapterAtomID: {
		Name:       "chapter_atom",
		Definition: "Contains the atom information to use as the chapter atom (apply to all tracks).",
		Type:       ebml.Master, Tag: ChapterAtom,
	},
}

var ChapterAtom = ebml.Tag{
	ChapterUIDID: {
		Name:       "chapter_uid",
		Definition: "A unique ID to identify the Chapter.",
		Type:       ebml.Uinteger,
	},
	ChapterStringUIDID: {
		Name:       "chapter_string_uid",
		Definition: "A unique string ID to identify the Chapter. Use for WebVTT cue identifier storage [@!WebVTT].",
		Type:       ebml.UTF8,
	},
	ChapterTimeStartID: {
		Name:       "chapter_time_start",
		Definition: "Timestamp of the start of Chapter, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks).",
		Type:       ebml.Uinteger,
	},
	ChapterTimeEndID: {
		Name:       "chapter_time_end",
		Definition: "Timestamp of the end of Chapter timestamp excluded, expressed in Matroska Ticks -- ie in nanoseconds; see (#timestamp-ticks). The value **MUST** be greater than or equal to the `ChapterTimeStart` of the same `ChapterAtom`.",
		Type:       ebml.Uinteger,
	},
	ChapterFlagHiddenID: {
		Name:       "chapter_flag_hidden",
		Definition: "Set to 1 if a chapter is hidden. Hidden chapters **SHOULD NOT** be available to the user interface (but still to Control Tracks; see (#chapterflaghidden) on Chapter flags).",
		Type:       ebml.Uinteger,
	},
	ChapterFlagEnabledID: {
		Name:       "chapter_flag_enabled",
		Definition: "Set to 1 if the chapter is enabled. It can be enabled/disabled by a Control Track. When disabled, the movie **SHOULD** skip all the content between the TimeStart and TimeEnd of this chapter; see (#chapter-flags) on Chapter flags.",
		Type:       ebml.Uinteger,
	},
	ChapterSegmentUIDID: {
		Name:       "chapter_segment_uid",
		Definition: "The SegmentUID of another Segment to play during this chapter.",
		Type:       ebml.Binary,
	},
	ChapterSegmentEditionUIDID: {
		Name:       "chapter_segment_edition_uid",
		Definition: "The EditionUID to play from the Segment linked in ChapterSegmentUID. If ChapterSegmentEditionUID is undeclared, then no Edition of the linked Segment is used; see (#medium-linking) on medium-linking Segments.",
		Type:       ebml.Uinteger,
	},
	ChapterPhysicalEquivID: {
		Name:       "chapter_physical_equiv",
		Definition: "Specify the physical equivalent of this ChapterAtom like \"DVD\" (60) or \"SIDE\" (50); see (#physical-types) for a complete list of values.",
		Type:       ebml.Uinteger,
	},
	ChapterTrackID: {
		Name:       "chapter_track",
		Definition: "List of tracks on which the chapter applies. If this Element is not present, all tracks apply",
		Type:       ebml.Master, Tag: ChapterTrack,
	},
	ChapterDisplayID: {
		Name:       "chapter_display",
		Definition: "Contains all possible strings to use for the chapter display.",
		Type:       ebml.Master, Tag: ChapterDisplay,
	},
	ChapProcessID: {
		Name:       "chap_process",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: ChapProcess,
	},
}

var ChapterTrack = ebml.Tag{
	ChapterTrackUIDID: {
		Name:       "chapter_track_uid",
		Definition: "UID of the Track to apply this chapter to. In the absence of a control track, choosing this chapter will select the listed Tracks and deselect unlisted tracks. Absence of this Element indicates that the Chapter **SHOULD** be applied to any currently used Tracks.",
		Type:       ebml.Uinteger,
	},
}

var ChapterDisplay = ebml.Tag{
	ChapStringID: {
		Name:       "chap_string",
		Definition: "Contains the string to use as the chapter atom.",
		Type:       ebml.UTF8,
	},
	ChapLanguageID: {
		Name:       "chap_language",
		Definition: "A language corresponding to the string, in the bibliographic ISO-639-2 form [@!ISO639-2]. This Element **MUST** be ignored if a ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
	ChapLanguageIETFID: {
		Name:       "chap_language_ietf",
		Definition: "Specifies a language corresponding to the ChapString in the format defined in [@!BCP47] and using the IANA Language Subtag Registry [@!IANALangRegistry]. If a ChapLanguageIETF Element is used, then any ChapLanguage and ChapCountry Elements used in the same ChapterDisplay **MUST** be ignored.",
		Type:       ebml.String,
	},
	ChapCountryID: {
		Name:       "chap_country",
		Definition: "A country corresponding to the string, using the same 2 octets country-codes as in Internet domains [@!IANADomains] based on [@!ISO3166-1] alpha-2 codes. This Element **MUST** be ignored if a ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
}

var ChapProcess = ebml.Tag{
	ChapProcessCodecIDID: {
		Name:       "chap_process_codec_id",
		Definition: "Contains the type of the codec used for the processing. A value of 0 means native Matroska processing (to be defined), a value of 1 means the DVD command set is used; see (#menu-features) on DVD menus. More codec IDs can be added later.",
		Type:       ebml.Uinteger,
	},
	ChapProcessPrivateID: {
		Name:       "chap_process_private",
		Definition: "Some optional data attached to the ChapProcessCodecID information.     For ChapProcessCodecID = 1, it is the \"DVD level\" equivalent; see (#menu-features) on DVD menus.",
		Type:       ebml.Binary,
	},
	ChapProcessCommandID: {
		Name:       "chap_process_command",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: ChapProcessCommand,
	},
}

var ChapProcessCommand = ebml.Tag{
	ChapProcessTimeID: {
		Name:       "chap_process_time",
		Definition: "Defines when the process command **SHOULD** be handled",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			0: {
				Sym: "during_the_whole_chapter",
			},
			1: {
				Sym: "before_starting_playback",
			},
			2: {
				Sym: "after_playback_of_the_chapter",
			},
		},
	},
	ChapProcessDataID: {
		Name:       "chap_process_data",
		Definition: "Contains the command information. The data **SHOULD** be interpreted depending on the ChapProcessCodecID value. For ChapProcessCodecID = 1, the data correspond to the binary DVD cell pre/post commands; see (#menu-features) on DVD menus.",
		Type:       ebml.Binary,
	},
}

var Tags = ebml.Tag{
	TagID: {
		Name:       "tag",
		Definition: "A single metadata descriptor.",
		Type:       ebml.Master, Tag: Tag,
	},
}

var Tag = ebml.Tag{
	TargetsID: {
		Name:       "targets",
		Definition: "Specifies which other elements the metadata represented by the Tag applies to. If empty or not present, then the Tag describes everything in the Segment.",
		Type:       ebml.Master, Tag: Targets,
	},
	SimpleTagID: {
		Name:       "simple_tag",
		Definition: "Contains general information about the target.",
		Type:       ebml.Master, Tag: SimpleTag,
	},
}

var Targets = ebml.Tag{
	TargetTypeValueID: {
		Name:       "target_type_value",
		Definition: "A number to indicate the logical level of the target.",
		Type:       ebml.Uinteger,
		UintegerEnums: scalar.UToScalar{
			70: {
				Sym:         "collection",
				Description: "The highest hierarchical level that tags can describe.",
			},
			60: {
				Sym:         "edition_issue_volume_opus_season_sequel",
				Description: "A list of lower levels grouped together.",
			},
			50: {
				Sym:         "album_opera_concert_movie_episode",
				Description: "The most common grouping level of music and video (equals to an episode for TV series).",
			},
			40: {
				Sym:         "part_session",
				Description: "When an album or episode has different logical parts.",
			},
			30: {
				Sym:         "track_song_chapter",
				Description: "The common parts of an album or movie.",
			},
			20: {
				Sym:         "subtrack_part_movement_scene",
				Description: "Corresponds to parts of a track for audio (like a movement).",
			},
			10: {
				Sym:         "shot",
				Description: "The lowest hierarchy found in music or movies.",
			},
		},
	},
	TargetTypeID: {
		Name:       "target_type",
		Definition: "An informational string that can be used to display the logical level of the target like \"ALBUM\", \"TRACK\", \"MOVIE\", \"CHAPTER\", etc ; see Section 6.4 of [@!MatroskaTags].",
		Type:       ebml.String,
		StringEnums: scalar.StrToScalar{
			"COLLECTION": {
				Sym: "collection",
			},
			"EDITION": {
				Sym: "edition",
			},
			"ISSUE": {
				Sym: "issue",
			},
			"VOLUME": {
				Sym: "volume",
			},
			"OPUS": {
				Sym: "opus",
			},
			"SEASON": {
				Sym: "season",
			},
			"SEQUEL": {
				Sym: "sequel",
			},
			"ALBUM": {
				Sym: "album",
			},
			"OPERA": {
				Sym: "opera",
			},
			"CONCERT": {
				Sym: "concert",
			},
			"MOVIE": {
				Sym: "movie",
			},
			"EPISODE": {
				Sym: "episode",
			},
			"PART": {
				Sym: "part",
			},
			"SESSION": {
				Sym: "session",
			},
			"TRACK": {
				Sym: "track",
			},
			"SONG": {
				Sym: "song",
			},
			"CHAPTER": {
				Sym: "chapter",
			},
			"SUBTRACK": {
				Sym: "subtrack",
			},
			"MOVEMENT": {
				Sym: "movement",
			},
			"SCENE": {
				Sym: "scene",
			},
			"SHOT": {
				Sym: "shot",
			},
		},
	},
	TagTrackUIDID: {
		Name:       "tag_track_uid",
		Definition: "A unique ID to identify the Track(s) the tags belong to.",
		Type:       ebml.Uinteger,
	},
	TagEditionUIDID: {
		Name:       "tag_edition_uid",
		Definition: "A unique ID to identify the EditionEntry(s) the tags belong to.",
		Type:       ebml.Uinteger,
	},
	TagChapterUIDID: {
		Name:       "tag_chapter_uid",
		Definition: "A unique ID to identify the Chapter(s) the tags belong to.",
		Type:       ebml.Uinteger,
	},
	TagAttachmentUIDID: {
		Name:       "tag_attachment_uid",
		Definition: "A unique ID to identify the Attachment(s) the tags belong to.",
		Type:       ebml.Uinteger,
	},
}

var SimpleTag = ebml.Tag{
	TagNameID: {
		Name:       "tag_name",
		Definition: "The name of the Tag that is going to be stored.",
		Type:       ebml.UTF8,
	},
	TagLanguageID: {
		Name:       "tag_language",
		Definition: "Specifies the language of the tag specified, in the Matroska languages form; see (#language-codes) on language codes. This Element **MUST** be ignored if the TagLanguageIETF Element is used within the same SimpleTag Element.",
		Type:       ebml.String,
	},
	TagLanguageIETFID: {
		Name:       "tag_language_ietf",
		Definition: "Specifies the language used in the TagString according to [@!BCP47] and using the IANA Language Subtag Registry [@!IANALangRegistry]. If this Element is used, then any TagLanguage Elements used in the same SimpleTag **MUST** be ignored.",
		Type:       ebml.String,
	},
	TagDefaultID: {
		Name:       "tag_default",
		Definition: "A boolean value to indicate if this is the default/original language to use for the given tag.",
		Type:       ebml.Uinteger,
	},
	TagDefaultBogusID: {
		Name:       "tag_default_bogus",
		Definition: "A variant of the TagDefault element with a bogus Element ID; see (#tagdefault-element).",
		Type:       ebml.Uinteger,
	},
	TagStringID: {
		Name:       "tag_string",
		Definition: "The value of the Tag.",
		Type:       ebml.UTF8,
	},
	TagBinaryID: {
		Name:       "tag_binary",
		Definition: "The values of the Tag, if it is binary. Note that this cannot be used in the same SimpleTag as TagString.",
		Type:       ebml.Binary,
	},
}
