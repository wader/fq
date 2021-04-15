// Code below generated from ebml_matroska.xml
package mkv

import "fq/pkg/format/mkv/ebml"

var mkvSegment = ebml.Tag{
	0x114d9b74: {
		Name:       "SeekHead",
		Definition: "Contains the Segment Position of other Top-Level Elements.",
		Type:       ebml.Master, Tag: mkvSeekHead,
	},
	0x1549a966: {
		Name:       "Info",
		Definition: "Contains general information about the Segment.",
		Type:       ebml.Master, Tag: mkvInfo,
	},
	0x1f43b675: {
		Name:       "Cluster",
		Definition: "The Top-Level Element containing the (monolithic) Block structure.",
		Type:       ebml.Master, Tag: mkvCluster,
	},
	0x1654ae6b: {
		Name:       "Tracks",
		Definition: "A Top-Level Element of information with many tracks described.",
		Type:       ebml.Master, Tag: mkvTracks,
	},
	0x1c53bb6b: {
		Name:       "Cues",
		Definition: "A Top-Level Element to speed seeking access. All entries are local to the Segment.",
		Type:       ebml.Master, Tag: mkvCues,
	},
	0x1941a469: {
		Name:       "Attachments",
		Definition: "Contain attached files.",
		Type:       ebml.Master, Tag: mkvAttachments,
	},
	0x1043a770: {
		Name:       "Chapters",
		Definition: "A system to define basic menus and partition data. For more detailed information, look at the .",
		Type:       ebml.Master, Tag: mkvChapters,
	},
	0x1254c367: {
		Name:       "Tags",
		Definition: "Element containing metadata describing Tracks, Editions, Chapters, Attachments, or the Segment as a whole. A list of valid tags can be found",
		Type:       ebml.Master, Tag: mkvTags,
	},
}

var mkvSeekHead = ebml.Tag{
	0x4dbb: {
		Name:       "Seek",
		Definition: "Contains a single seek entry to an EBML Element.",
		Type:       ebml.Master, Tag: mkvSeek,
	},
}

var mkvSeek = ebml.Tag{
	0x53ab: {
		Name:       "SeekID",
		Definition: "The binary ID corresponding to the Element name.",
		Type:       ebml.Binary,
	},
	0x53ac: {
		Name:       "SeekPosition",
		Definition: "The Segment Position of the Element.",
		Type:       ebml.Uinteger,
	},
}

var mkvInfo = ebml.Tag{
	0x73a4: {
		Name:       "SegmentUID",
		Definition: "A randomly generated unique ID to identify the Segment amongst many others (128 bits).",
		Type:       ebml.Binary,
	},
	0x7384: {
		Name:       "SegmentFilename",
		Definition: "A filename corresponding to this Segment.",
		Type:       ebml.UTF8,
	},
	0x3cb923: {
		Name:       "PrevUID",
		Definition: "A unique ID to identify the previous Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	0x3c83ab: {
		Name:       "PrevFilename",
		Definition: "A filename corresponding to the file of the previous Linked Segment.",
		Type:       ebml.UTF8,
	},
	0x3eb923: {
		Name:       "NextUID",
		Definition: "A unique ID to identify the next Segment of a Linked Segment (128 bits).",
		Type:       ebml.Binary,
	},
	0x3e83bb: {
		Name:       "NextFilename",
		Definition: "A filename corresponding to the file of the next Linked Segment.",
		Type:       ebml.UTF8,
	},
	0x4444: {
		Name:       "SegmentFamily",
		Definition: "A randomly generated unique ID that all Segments of a Linked Segment MUST share (128 bits).",
		Type:       ebml.Binary,
	},
	0x6924: {
		Name:       "ChapterTranslate",
		Definition: "A tuple of corresponding ID used by chapter codecs to represent this Segment.",
		Type:       ebml.Master, Tag: mkvChapterTranslate,
	},
	0x2ad7b1: {
		Name:       "TimestampScale",
		Definition: "Timestamp scale in nanoseconds (1.000.000 means all timestamps in the Segment are expressed in milliseconds).",
		Type:       ebml.Uinteger,
	},
	0x4489: {
		Name:       "Duration",
		Definition: "Duration of the Segment in nanoseconds based on TimestampScale.",
		Type:       ebml.Float,
	},
	0x4461: {
		Name:       "DateUTC",
		Definition: "The date and time that the Segment was created by the muxing application or library.",
		Type:       ebml.Date,
	},
	0x7ba9: {
		Name:       "Title",
		Definition: "General name of the Segment.",
		Type:       ebml.UTF8,
	},
	0x4d80: {
		Name:       "MuxingApp",
		Definition: "Muxing application or library (example: \"libmatroska-0.4.3\").",
		Type:       ebml.UTF8,
	},
	0x5741: {
		Name:       "WritingApp",
		Definition: "Writing application (example: \"mkvmerge-0.3.3\").",
		Type:       ebml.UTF8,
	},
}

var mkvChapterTranslate = ebml.Tag{
	0x69fc: {
		Name:       "ChapterTranslateEditionUID",
		Definition: "Specify an edition UID on which this correspondence applies. When not specified, it means for all editions found in the Segment.",
		Type:       ebml.Uinteger,
	},
	0x69bf: {
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
	0x69a5: {
		Name:       "ChapterTranslateID",
		Definition: "The binary value used to represent this Segment in the chapter codec data. The format depends on the  used.",
		Type:       ebml.Binary,
	},
}

var mkvCluster = ebml.Tag{
	0xe7: {
		Name:       "Timestamp",
		Definition: "Absolute timestamp of the cluster (based on TimestampScale).",
		Type:       ebml.Uinteger,
	},
	0x5854: {
		Name:       "SilentTracks",
		Definition: "The list of tracks that are not used in that part of the stream. It is useful when using overlay tracks on seeking or to decide what track to use.",
		Type:       ebml.Master, Tag: mkvSilentTracks,
	},
	0xa7: {
		Name:       "Position",
		Definition: "The Segment Position of the Cluster in the Segment (0 in live streams). It might help to resynchronise offset on damaged streams.",
		Type:       ebml.Uinteger,
	},
	0xab: {
		Name:       "PrevSize",
		Definition: "Size of the previous Cluster, in octets. Can be useful for backward playing.",
		Type:       ebml.Uinteger,
	},
	0xa3: {
		Name:       "SimpleBlock",
		Definition: "Similar to  but without all the extra information, mostly used to reduced overhead when no extra feature is needed. (see )",
		Type:       ebml.Binary,
	},
	0xa0: {
		Name:       "BlockGroup",
		Definition: "Basic container of information containing a single Block and information specific to that Block.",
		Type:       ebml.Master, Tag: mkvBlockGroup,
	},
	0xaf: {
		Name:       "EncryptedBlock",
		Definition: "Similar to  but the data inside the Block are Transformed (encrypt and/or signed). (see )",
		Type:       ebml.Binary,
	},
}

var mkvSilentTracks = ebml.Tag{
	0x58d7: {
		Name:       "SilentTrackNumber",
		Definition: "One of the track number that are not used from now on in the stream. It could change later if not specified as silent in a further Cluster.",
		Type:       ebml.Uinteger,
	},
}

var mkvBlockGroup = ebml.Tag{
	0xa1: {
		Name:       "Block",
		Definition: "Block containing the actual data to be rendered and a timestamp relative to the Cluster Timestamp. (see )",
		Type:       ebml.Binary,
	},
	0xa2: {
		Name:       "BlockVirtual",
		Definition: "A Block with no data. It MUST be stored in the stream at the place the real Block would be in display order. (see )",
		Type:       ebml.Binary,
	},
	0x75a1: {
		Name:       "BlockAdditions",
		Definition: "Contain additional blocks to complete the main one. An EBML parser that has no knowledge of the Block structure could still see and use/skip these data.",
		Type:       ebml.Master, Tag: mkvBlockAdditions,
	},
	0x9b: {
		Name:       "BlockDuration",
		Definition: "The duration of the Block (based on TimestampScale). The BlockDuration Element can be useful at the end of a Track to define the duration of the last frame (as there is no subsequent Block available), or when there is a break in a track like for subtitle tracks.",
		Type:       ebml.Uinteger,
	},
	0xfa: {
		Name:       "ReferencePriority",
		Definition: "This frame is referenced and has the specified cache priority. In cache only a frame of the same or higher priority can replace this frame. A value of 0 means the frame is not referenced.",
		Type:       ebml.Uinteger,
	},
	0xfb: {
		Name:       "ReferenceBlock",
		Definition: "Timestamp of another frame used as a reference (ie: B or P frame). The timestamp is relative to the block it's attached to.",
		Type:       ebml.Integer,
	},
	0xfd: {
		Name:       "ReferenceVirtual",
		Definition: "The Segment Position of the data that would otherwise be in position of the virtual block.",
		Type:       ebml.Integer,
	},
	0xa4: {
		Name:       "CodecState",
		Definition: "The new codec state to use. Data interpretation is private to the codec. This information SHOULD always be referenced by a seek entry.",
		Type:       ebml.Binary,
	},
	0x75a2: {
		Name:       "DiscardPadding",
		Definition: "Duration in nanoseconds of the silent data added to the Block (padding at the end of the Block for positive value, at the beginning of the Block for negative value). The duration of DiscardPadding is not calculated in the duration of the TrackEntry and SHOULD be discarded during playback.",
		Type:       ebml.Integer,
	},
	0x8e: {
		Name:       "Slices",
		Definition: "Contains slices description.",
		Type:       ebml.Master, Tag: mkvSlices,
	},
	0xc8: {
		Name:       "ReferenceFrame",
		Definition: "",
		Type:       ebml.Master, Tag: mkvReferenceFrame,
	},
}

var mkvBlockAdditions = ebml.Tag{
	0xa6: {
		Name:       "BlockMore",
		Definition: "Contain the BlockAdditional and some parameters.",
		Type:       ebml.Master, Tag: mkvBlockMore,
	},
}

var mkvBlockMore = ebml.Tag{
	0xee: {
		Name:       "BlockAddID",
		Definition: "An ID to identify the BlockAdditional level. A value of 1 means the BlockAdditional data is interpreted as additional data passed to the codec with the Block data.",
		Type:       ebml.Uinteger,
	},
	0xa5: {
		Name:       "BlockAdditional",
		Definition: "Interpreted by the codec as it wishes (using the BlockAddID).",
		Type:       ebml.Binary,
	},
}

var mkvSlices = ebml.Tag{
	0xe8: {
		Name:       "TimeSlice",
		Definition: "Contains extra time information about the data contained in the Block. Being able to interpret this Element is not REQUIRED for playback.",
		Type:       ebml.Master, Tag: mkvTimeSlice,
	},
}

var mkvTimeSlice = ebml.Tag{
	0xcc: {
		Name:       "LaceNumber",
		Definition: "The reverse number of the frame in the lace (0 is the last frame, 1 is the next to last, etc). Being able to interpret this Element is not REQUIRED for playback.",
		Type:       ebml.Uinteger,
	},
	0xcd: {
		Name:       "FrameNumber",
		Definition: "The number of the frame to generate from this lace with this delay (allow you to generate many frames from the same Block/Frame).",
		Type:       ebml.Uinteger,
	},
	0xcb: {
		Name:       "BlockAdditionID",
		Definition: "The ID of the BlockAdditional Element (0 is the main Block).",
		Type:       ebml.Uinteger,
	},
	0xce: {
		Name:       "Delay",
		Definition: "The (scaled) delay to apply to the Element.",
		Type:       ebml.Uinteger,
	},
	0xcf: {
		Name:       "SliceDuration",
		Definition: "The (scaled) duration to apply to the Element.",
		Type:       ebml.Uinteger,
	},
}

var mkvReferenceFrame = ebml.Tag{
	0xc9: {
		Name:       "ReferenceOffset",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	0xca: {
		Name:       "ReferenceTimestamp",
		Definition: "",
		Type:       ebml.Uinteger,
	},
}

var mkvTracks = ebml.Tag{
	0xae: {
		Name:       "TrackEntry",
		Definition: "Describes a track with all Elements.",
		Type:       ebml.Master, Tag: mkvTrackEntry,
	},
}

var mkvTrackEntry = ebml.Tag{
	0xd7: {
		Name:       "TrackNumber",
		Definition: "The track number as used in the Block Header (using more than 127 tracks is not encouraged, though the design allows an unlimited number).",
		Type:       ebml.Uinteger,
	},
	0x73c5: {
		Name:       "TrackUID",
		Definition: "A unique ID to identify the Track. This SHOULD be kept the same when making a direct stream copy of the Track to another file.",
		Type:       ebml.Uinteger,
	},
	0x83: {
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
	0xb9: {
		Name:       "FlagEnabled",
		Definition: "Set if the track is usable. (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x88: {
		Name:       "FlagDefault",
		Definition: "Set if that track (audio, video or subs) SHOULD be active if no language found matches the user preference. (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x55aa: {
		Name:       "FlagForced",
		Definition: "Set if that track MUST be active during playback. There can be many forced track for a kind (audio, video or subs), the player SHOULD select the one which language matches the user preference or the default + forced track. Overlay MAY happen between a forced and non-forced track of the same kind. (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x9c: {
		Name:       "FlagLacing",
		Definition: "Set if the track MAY contain blocks using lacing. (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x6de7: {
		Name:       "MinCache",
		Definition: "The minimum number of frames a player SHOULD be able to cache during playback. If set to 0, the reference pseudo-cache system is not used.",
		Type:       ebml.Uinteger,
	},
	0x6df8: {
		Name:       "MaxCache",
		Definition: "The maximum cache size necessary to store referenced frames in and the current frame. 0 means no cache is needed.",
		Type:       ebml.Uinteger,
	},
	0x23e383: {
		Name:       "DefaultDuration",
		Definition: "Number of nanoseconds (not scaled via TimestampScale) per frame ('frame' in the Matroska sense -- one Element put into a (Simple)Block).",
		Type:       ebml.Uinteger,
	},
	0x234e7a: {
		Name:       "DefaultDecodedFieldDuration",
		Definition: "The period in nanoseconds (not scaled by TimestampScale) between two successive fields at the output of the decoding process (see )",
		Type:       ebml.Uinteger,
	},
	0x23314f: {
		Name:       "TrackTimestampScale",
		Definition: "DEPRECATED, DO NOT USE. The scale to apply on this track to work at normal speed in relation with other tracks (mostly used to adjust video speed when the audio length differs).",
		Type:       ebml.Float,
	},
	0x537f: {
		Name:       "TrackOffset",
		Definition: "A value to add to the Block's Timestamp. This can be used to adjust the playback offset of a track.",
		Type:       ebml.Integer,
	},
	0x55ee: {
		Name:       "MaxBlockAdditionID",
		Definition: "The maximum value of . A value 0 means there is no  for this track.",
		Type:       ebml.Uinteger,
	},
	0x41e4: {
		Name:       "BlockAdditionMapping",
		Definition: "Contains elements that describe each value of  found in the Track.",
		Type:       ebml.Master, Tag: mkvBlockAdditionMapping,
	},
	0x536e: {
		Name:       "Name",
		Definition: "A human-readable track name.",
		Type:       ebml.UTF8,
	},
	0x22b59c: {
		Name:       "Language",
		Definition: "Specifies the language of the track in the . This Element MUST be ignored if the LanguageIETF Element is used in the same TrackEntry.",
		Type:       ebml.String,
	},
	0x22b59d: {
		Name:       "LanguageIETF",
		Definition: "Specifies the language of the track according to  and using the . If this Element is used, then any Language Elements used in the same TrackEntry MUST be ignored.",
		Type:       ebml.String,
	},
	0x86: {
		Name:       "CodecID",
		Definition: "An ID corresponding to the codec, see the  for more info.",
		Type:       ebml.String,
	},
	0x63a2: {
		Name:       "CodecPrivate",
		Definition: "Private data only known to the codec.",
		Type:       ebml.Binary,
	},
	0x258688: {
		Name:       "CodecName",
		Definition: "A human-readable string specifying the codec.",
		Type:       ebml.UTF8,
	},
	0x7446: {
		Name:       "AttachmentLink",
		Definition: "The UID of an attachment that is used by this codec.",
		Type:       ebml.Uinteger,
	},
	0x3a9697: {
		Name:       "CodecSettings",
		Definition: "A string describing the encoding setting used.",
		Type:       ebml.UTF8,
	},
	0x3b4040: {
		Name:       "CodecInfoURL",
		Definition: "A URL to find information about the codec used.",
		Type:       ebml.String,
	},
	0x26b240: {
		Name:       "CodecDownloadURL",
		Definition: "A URL to download about the codec used.",
		Type:       ebml.String,
	},
	0xaa: {
		Name:       "CodecDecodeAll",
		Definition: "The codec can decode potentially damaged data (1 bit).",
		Type:       ebml.Uinteger,
	},
	0x6fab: {
		Name:       "TrackOverlay",
		Definition: "Specify that this track is an overlay track for the Track specified (in the u-integer). That means when this track has a gap (see ) the overlay track SHOULD be used instead. The order of multiple TrackOverlay matters, the first one is the one that SHOULD be used. If not found it SHOULD be the second, etc.",
		Type:       ebml.Uinteger,
	},
	0x56aa: {
		Name:       "CodecDelay",
		Definition: "CodecDelay is The codec-built-in delay in nanoseconds. This value MUST be subtracted from each block timestamp in order to get the actual timestamp. The value SHOULD be small so the muxing of tracks with the same actual timestamp are in the same Cluster.",
		Type:       ebml.Uinteger,
	},
	0x56bb: {
		Name:       "SeekPreRoll",
		Definition: "After a discontinuity, SeekPreRoll is the duration in nanoseconds of the data the decoder MUST decode before the decoded data is valid.",
		Type:       ebml.Uinteger,
	},
	0x6624: {
		Name:       "TrackTranslate",
		Definition: "The track identification for the given Chapter Codec.",
		Type:       ebml.Master, Tag: mkvTrackTranslate,
	},
	0xe0: {
		Name:       "Video",
		Definition: "Video settings.",
		Type:       ebml.Master, Tag: mkvVideo,
	},
	0xe1: {
		Name:       "Audio",
		Definition: "Audio settings.",
		Type:       ebml.Master, Tag: mkvAudio,
	},
	0xe2: {
		Name:       "TrackOperation",
		Definition: "Operation that needs to be applied on tracks to create this virtual track. For more details  on the subject.",
		Type:       ebml.Master, Tag: mkvTrackOperation,
	},
	0xc0: {
		Name:       "TrickTrackUID",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	0xc1: {
		Name:       "TrickTrackSegmentUID",
		Definition: "",
		Type:       ebml.Binary,
	},
	0xc6: {
		Name:       "TrickTrackFlag",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	0xc7: {
		Name:       "TrickMasterTrackUID",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	0xc4: {
		Name:       "TrickMasterTrackSegmentUID",
		Definition: "",
		Type:       ebml.Binary,
	},
	0x6d80: {
		Name:       "ContentEncodings",
		Definition: "Settings for several content encoding mechanisms like compression or encryption.",
		Type:       ebml.Master, Tag: mkvContentEncodings,
	},
}

var mkvBlockAdditionMapping = ebml.Tag{
	0x41f0: {
		Name:       "BlockAddIDValue",
		Definition: "The  value being described. To keep MaxBlockAdditionID as low as possible, small values SHOULD be used.",
		Type:       ebml.Uinteger,
	},
	0x41a4: {
		Name:       "BlockAddIDName",
		Definition: "A human-friendly name describing the type of BlockAdditional data as defined by the associated Block Additional Mapping.",
		Type:       ebml.String,
	},
	0x41e7: {
		Name:       "BlockAddIDType",
		Definition: "Stores the registered identifer of the Block Additional Mapping to define how the BlockAdditional data should be handled.",
		Type:       ebml.Uinteger,
	},
	0x41ed: {
		Name:       "BlockAddIDExtraData",
		Definition: "Extra binary data that the BlockAddIDType can use to interpret the BlockAdditional data. The intepretation of the binary data depends on the BlockAddIDType value and the corresponding Block Additional Mapping.",
		Type:       ebml.Binary,
	},
}

var mkvTrackTranslate = ebml.Tag{
	0x66fc: {
		Name:       "TrackTranslateEditionUID",
		Definition: "Specify an edition UID on which this translation applies. When not specified, it means for all editions found in the Segment.",
		Type:       ebml.Uinteger,
	},
	0x66bf: {
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
	0x66a5: {
		Name:       "TrackTranslateTrackID",
		Definition: "The binary value used to represent this track in the chapter codec data. The format depends on the  used.",
		Type:       ebml.Binary,
	},
}

var mkvVideo = ebml.Tag{
	0x9a: {
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
	0x9d: {
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
	0x53b8: {
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
	0x53c0: {
		Name:       "AlphaMode",
		Definition: "Alpha Video Mode. Presence of this Element indicates that the BlockAdditional Element could contain Alpha data.",
		Type:       ebml.Uinteger,
	},
	0x53b9: {
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
	0xb0: {
		Name:       "PixelWidth",
		Definition: "Width of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	0xba: {
		Name:       "PixelHeight",
		Definition: "Height of the encoded video frames in pixels.",
		Type:       ebml.Uinteger,
	},
	0x54aa: {
		Name:       "PixelCropBottom",
		Definition: "The number of video pixels to remove at the bottom of the image.",
		Type:       ebml.Uinteger,
	},
	0x54bb: {
		Name:       "PixelCropTop",
		Definition: "The number of video pixels to remove at the top of the image.",
		Type:       ebml.Uinteger,
	},
	0x54cc: {
		Name:       "PixelCropLeft",
		Definition: "The number of video pixels to remove on the left of the image.",
		Type:       ebml.Uinteger,
	},
	0x54dd: {
		Name:       "PixelCropRight",
		Definition: "The number of video pixels to remove on the right of the image.",
		Type:       ebml.Uinteger,
	},
	0x54b0: {
		Name:       "DisplayWidth",
		Definition: "Width of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	0x54ba: {
		Name:       "DisplayHeight",
		Definition: "Height of the video frames to display. Applies to the video frame after cropping (PixelCrop* Elements).",
		Type:       ebml.Uinteger,
	},
	0x54b2: {
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
	0x54b3: {
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
	0x2eb524: {
		Name:       "ColourSpace",
		Definition: "Specify the pixel format used for the Track's data as a FourCC. This value is similar in scope to the biCompression value of AVI's BITMAPINFOHEADER.",
		Type:       ebml.Binary,
	},
	0x2fb523: {
		Name:       "GammaValue",
		Definition: "Gamma Value.",
		Type:       ebml.Float,
	},
	0x2383e3: {
		Name:       "FrameRate",
		Definition: "Number of frames per second.  only.",
		Type:       ebml.Float,
	},
	0x55b0: {
		Name:       "Colour",
		Definition: "Settings describing the colour format.",
		Type:       ebml.Master, Tag: mkvColour,
	},
	0x7670: {
		Name:       "Projection",
		Definition: "Describes the video projection details. Used to render spherical and VR videos.",
		Type:       ebml.Master, Tag: mkvProjection,
	},
}

var mkvColour = ebml.Tag{
	0x55b1: {
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
	0x55b2: {
		Name:       "BitsPerChannel",
		Definition: "Number of decoded bits per channel. A value of 0 indicates that the BitsPerChannel is unspecified.",
		Type:       ebml.Uinteger,
	},
	0x55b3: {
		Name:       "ChromaSubsamplingHorz",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed horizontally. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingHorz SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	0x55b4: {
		Name:       "ChromaSubsamplingVert",
		Definition: "The amount of pixels to remove in the Cr and Cb channels for every pixel not removed vertically. Example: For video with 4:2:0 chroma subsampling, the ChromaSubsamplingVert SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	0x55b5: {
		Name:       "CbSubsamplingHorz",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed horizontally. This is additive with ChromaSubsamplingHorz. Example: For video with 4:2:1 chroma subsampling, the ChromaSubsamplingHorz SHOULD be set to 1 and CbSubsamplingHorz SHOULD be set to 1.",
		Type:       ebml.Uinteger,
	},
	0x55b6: {
		Name:       "CbSubsamplingVert",
		Definition: "The amount of pixels to remove in the Cb channel for every pixel not removed vertically. This is additive with ChromaSubsamplingVert.",
		Type:       ebml.Uinteger,
	},
	0x55b7: {
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
	0x55b8: {
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
	0x55b9: {
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
	0x55ba: {
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
	0x55bb: {
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
	0x55bc: {
		Name:       "MaxCLL",
		Definition: "Maximum brightness of a single pixel (Maximum Content Light Level) in candelas per square meter (cd/m²).",
		Type:       ebml.Uinteger,
	},
	0x55bd: {
		Name:       "MaxFALL",
		Definition: "Maximum brightness of a single full frame (Maximum Frame-Average Light Level) in candelas per square meter (cd/m²).",
		Type:       ebml.Uinteger,
	},
	0x55d0: {
		Name:       "MasteringMetadata",
		Definition: "SMPTE 2086 mastering data.",
		Type:       ebml.Master, Tag: mkvMasteringMetadata,
	},
}

var mkvMasteringMetadata = ebml.Tag{
	0x55d1: {
		Name:       "PrimaryRChromaticityX",
		Definition: "Red X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d2: {
		Name:       "PrimaryRChromaticityY",
		Definition: "Red Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d3: {
		Name:       "PrimaryGChromaticityX",
		Definition: "Green X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d4: {
		Name:       "PrimaryGChromaticityY",
		Definition: "Green Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d5: {
		Name:       "PrimaryBChromaticityX",
		Definition: "Blue X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d6: {
		Name:       "PrimaryBChromaticityY",
		Definition: "Blue Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d7: {
		Name:       "WhitePointChromaticityX",
		Definition: "White X chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d8: {
		Name:       "WhitePointChromaticityY",
		Definition: "White Y chromaticity coordinate as defined by CIE 1931.",
		Type:       ebml.Float,
	},
	0x55d9: {
		Name:       "LuminanceMax",
		Definition: "Maximum luminance. Represented in candelas per square meter (cd/m²).",
		Type:       ebml.Float,
	},
	0x55da: {
		Name:       "LuminanceMin",
		Definition: "Minimum luminance. Represented in candelas per square meter (cd/m²).",
		Type:       ebml.Float,
	},
}

var mkvProjection = ebml.Tag{
	0x7671: {
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
	0x7672: {
		Name:       "ProjectionPrivate",
		Definition: "Private data that only applies to a specific projection.SemanticsIf ProjectionType equals 0 (Rectangular),\n     then this element must not be present.If ProjectionType equals 1 (Equirectangular), then this element must be present and contain the same binary data that would be stored inside\n      an ISOBMFF Equirectangular Projection Box ('equi').If ProjectionType equals 2 (Cubemap), then this element must be present and contain the same binary data that would be stored \n      inside an ISOBMFF Cubemap Projection Box ('cbmp').If ProjectionType equals 3 (Mesh), then this element must be present and contain the same binary data that would be stored inside\n       an ISOBMFF Mesh Projection Box ('mshp').Note: ISOBMFF box size and fourcc fields are not included in the binary data, but the FullBox version and flag fields are. This is to avoid \n       redundant framing information while preserving versioning and semantics between the two container formats.",
		Type:       ebml.Binary,
	},
	0x7673: {
		Name:       "ProjectionPoseYaw",
		Definition: "Specifies a yaw rotation to the projection.SemanticsValue represents a clockwise rotation, in degrees, around the up vector. This rotation must be applied before any ProjectionPosePitch or ProjectionPoseRoll rotations. The value of this field should be in the -180 to 180 degree range.",
		Type:       ebml.Float,
	},
	0x7674: {
		Name:       "ProjectionPosePitch",
		Definition: "Specifies a pitch rotation to the projection.SemanticsValue represents a counter-clockwise rotation, in degrees, around the right vector. This rotation must be applied after the ProjectionPoseYaw rotation and before the ProjectionPoseRoll rotation. The value of this field should be in the -90 to 90 degree range.",
		Type:       ebml.Float,
	},
	0x7675: {
		Name:       "ProjectionPoseRoll",
		Definition: "Specifies a roll rotation to the projection.SemanticsValue represents a counter-clockwise rotation, in degrees, around the forward vector. This rotation must be applied after the ProjectionPoseYaw and ProjectionPosePitch rotations. The value of this field should be in the -180 to 180 degree range.",
		Type:       ebml.Float,
	},
}

var mkvAudio = ebml.Tag{
	0xb5: {
		Name:       "SamplingFrequency",
		Definition: "Sampling frequency in Hz.",
		Type:       ebml.Float,
	},
	0x78b5: {
		Name:       "OutputSamplingFrequency",
		Definition: "Real output sampling frequency in Hz (used for SBR techniques).",
		Type:       ebml.Float,
	},
	0x9f: {
		Name:       "Channels",
		Definition: "Numbers of channels in the track.",
		Type:       ebml.Uinteger,
	},
	0x7d7b: {
		Name:       "ChannelPositions",
		Definition: "Table of horizontal angles for each successive channel, see .",
		Type:       ebml.Binary,
	},
	0x6264: {
		Name:       "BitDepth",
		Definition: "Bits per sample, mostly used for PCM.",
		Type:       ebml.Uinteger,
	},
}

var mkvTrackOperation = ebml.Tag{
	0xe3: {
		Name:       "TrackCombinePlanes",
		Definition: "Contains the list of all video plane tracks that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: mkvTrackCombinePlanes,
	},
	0xe9: {
		Name:       "TrackJoinBlocks",
		Definition: "Contains the list of all tracks whose Blocks need to be combined to create this virtual track",
		Type:       ebml.Master, Tag: mkvTrackJoinBlocks,
	},
}

var mkvTrackCombinePlanes = ebml.Tag{
	0xe4: {
		Name:       "TrackPlane",
		Definition: "Contains a video plane track that need to be combined to create this 3D track",
		Type:       ebml.Master, Tag: mkvTrackPlane,
	},
}

var mkvTrackPlane = ebml.Tag{
	0xe5: {
		Name:       "TrackPlaneUID",
		Definition: "The trackUID number of the track representing the plane.",
		Type:       ebml.Uinteger,
	},
	0xe6: {
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
	0xed: {
		Name:       "TrackJoinUID",
		Definition: "The trackUID number of a track whose blocks are used to create this virtual track.",
		Type:       ebml.Uinteger,
	},
}

var mkvContentEncodings = ebml.Tag{
	0x6240: {
		Name:       "ContentEncoding",
		Definition: "Settings for one content encoding like compression or encryption.",
		Type:       ebml.Master, Tag: mkvContentEncoding,
	},
}

var mkvContentEncoding = ebml.Tag{
	0x5031: {
		Name:       "ContentEncodingOrder",
		Definition: "Tells when this modification was used during encoding/muxing starting with 0 and counting upwards. The decoder/demuxer has to start with the highest order number it finds and work its way down. This value has to be unique over all ContentEncodingOrder Elements in the TrackEntry that contains this ContentEncodingOrder element.",
		Type:       ebml.Uinteger,
	},
	0x5032: {
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
	0x5033: {
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
	0x5034: {
		Name:       "ContentCompression",
		Definition: "Settings describing the compression used. This Element MUST be present if the value of ContentEncodingType is 0 and absent otherwise. Each block MUST be decompressable even if no previous block is available in order not to prevent seeking.",
		Type:       ebml.Master, Tag: mkvContentCompression,
	},
	0x5035: {
		Name:       "ContentEncryption",
		Definition: "Settings describing the encryption used. This Element MUST be present if the value of `ContentEncodingType` is 1 (encryption) and MUST be ignored otherwise.",
		Type:       ebml.Master, Tag: mkvContentEncryption,
	},
}

var mkvContentCompression = ebml.Tag{
	0x4254: {
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
	0x4255: {
		Name:       "ContentCompSettings",
		Definition: "Settings that might be needed by the decompressor. For Header Stripping (`ContentCompAlgo`=3), the bytes that were removed from the beginning of each frames of the track.",
		Type:       ebml.Binary,
	},
}

var mkvContentEncryption = ebml.Tag{
	0x47e1: {
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
	0x47e2: {
		Name:       "ContentEncKeyID",
		Definition: "For public key algorithms this is the ID of the public key the the data was encrypted with.",
		Type:       ebml.Binary,
	},
	0x47e7: {
		Name:       "ContentEncAESSettings",
		Definition: "Settings describing the encryption algorithm used. If `ContentEncAlgo` != 5 this MUST be ignored.",
		Type:       ebml.Master, Tag: mkvContentEncAESSettings,
	},
	0x47e3: {
		Name:       "ContentSignature",
		Definition: "A cryptographic signature of the contents.",
		Type:       ebml.Binary,
	},
	0x47e4: {
		Name:       "ContentSigKeyID",
		Definition: "This is the ID of the private key the data was signed with.",
		Type:       ebml.Binary,
	},
	0x47e5: {
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
	0x47e6: {
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
	0x47e8: {
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
	0xbb: {
		Name:       "CuePoint",
		Definition: "Contains all information relative to a seek point in the Segment.",
		Type:       ebml.Master, Tag: mkvCuePoint,
	},
}

var mkvCuePoint = ebml.Tag{
	0xb3: {
		Name:       "CueTime",
		Definition: "Absolute timestamp according to the Segment time base.",
		Type:       ebml.Uinteger,
	},
	0xb7: {
		Name:       "CueTrackPositions",
		Definition: "Contain positions for different tracks corresponding to the timestamp.",
		Type:       ebml.Master, Tag: mkvCueTrackPositions,
	},
}

var mkvCueTrackPositions = ebml.Tag{
	0xf7: {
		Name:       "CueTrack",
		Definition: "The track for which a position is given.",
		Type:       ebml.Uinteger,
	},
	0xf1: {
		Name:       "CueClusterPosition",
		Definition: "The Segment Position of the Cluster containing the associated Block.",
		Type:       ebml.Uinteger,
	},
	0xf0: {
		Name:       "CueRelativePosition",
		Definition: "The relative position inside the Cluster of the referenced SimpleBlock or BlockGroup with 0 being the first possible position for an Element inside that Cluster.",
		Type:       ebml.Uinteger,
	},
	0xb2: {
		Name:       "CueDuration",
		Definition: "The duration of the block according to the Segment time base. If missing the track's DefaultDuration does not apply and no duration information is available in terms of the cues.",
		Type:       ebml.Uinteger,
	},
	0x5378: {
		Name:       "CueBlockNumber",
		Definition: "Number of the Block in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	0xea: {
		Name:       "CueCodecState",
		Definition: "The Segment Position of the Codec State corresponding to this Cue Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
	0xdb: {
		Name:       "CueReference",
		Definition: "The Clusters containing the referenced Blocks.",
		Type:       ebml.Master, Tag: mkvCueReference,
	},
}

var mkvCueReference = ebml.Tag{
	0x96: {
		Name:       "CueRefTime",
		Definition: "Timestamp of the referenced Block.",
		Type:       ebml.Uinteger,
	},
	0x97: {
		Name:       "CueRefCluster",
		Definition: "The Segment Position of the Cluster containing the referenced Block.",
		Type:       ebml.Uinteger,
	},
	0x535f: {
		Name:       "CueRefNumber",
		Definition: "Number of the referenced Block of Track X in the specified Cluster.",
		Type:       ebml.Uinteger,
	},
	0xeb: {
		Name:       "CueRefCodecState",
		Definition: "The Segment Position of the Codec State corresponding to this referenced Element. 0 means that the data is taken from the initial Track Entry.",
		Type:       ebml.Uinteger,
	},
}

var mkvAttachments = ebml.Tag{
	0x61a7: {
		Name:       "AttachedFile",
		Definition: "An attached file.",
		Type:       ebml.Master, Tag: mkvAttachedFile,
	},
}

var mkvAttachedFile = ebml.Tag{
	0x467e: {
		Name:       "FileDescription",
		Definition: "A human-friendly name for the attached file.",
		Type:       ebml.UTF8,
	},
	0x466e: {
		Name:       "FileName",
		Definition: "Filename of the attached file.",
		Type:       ebml.UTF8,
	},
	0x4660: {
		Name:       "FileMimeType",
		Definition: "MIME type of the file.",
		Type:       ebml.String,
	},
	0x465c: {
		Name:       "FileData",
		Definition: "The data of the file.",
		Type:       ebml.Binary,
	},
	0x46ae: {
		Name:       "FileUID",
		Definition: "Unique ID representing the file, as random as possible.",
		Type:       ebml.Uinteger,
	},
	0x4675: {
		Name:       "FileReferral",
		Definition: "A binary value that a track/codec can refer to when the attachment is needed.",
		Type:       ebml.Binary,
	},
	0x4661: {
		Name:       "FileUsedStartTime",
		Definition: "",
		Type:       ebml.Uinteger,
	},
	0x4662: {
		Name:       "FileUsedEndTime",
		Definition: "",
		Type:       ebml.Uinteger,
	},
}

var mkvChapters = ebml.Tag{
	0x45b9: {
		Name:       "EditionEntry",
		Definition: "Contains all information about a Segment edition.",
		Type:       ebml.Master, Tag: mkvEditionEntry,
	},
}

var mkvEditionEntry = ebml.Tag{
	0x45bc: {
		Name:       "EditionUID",
		Definition: "A unique ID to identify the edition. It's useful for tagging an edition.",
		Type:       ebml.Uinteger,
	},
	0x45bd: {
		Name:       "EditionFlagHidden",
		Definition: "If an edition is hidden (1), it SHOULD NOT be available to the user interface (but still to Control Tracks; see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x45db: {
		Name:       "EditionFlagDefault",
		Definition: "If a flag is set (1) the edition SHOULD be used as the default one. (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x45dd: {
		Name:       "EditionFlagOrdered",
		Definition: "Specify if the chapters can be defined multiple times and the order to play them is enforced. (1 bit)",
		Type:       ebml.Uinteger,
	},
	0xb6: {
		Name:       "ChapterAtom",
		Definition: "Contains the atom information to use as the chapter atom (apply to all tracks).",
		Type:       ebml.Master, Tag: mkvChapterAtom,
	},
}

var mkvChapterAtom = ebml.Tag{
	0x73c4: {
		Name:       "ChapterUID",
		Definition: "A unique ID to identify the Chapter.",
		Type:       ebml.Uinteger,
	},
	0x5654: {
		Name:       "ChapterStringUID",
		Definition: "A unique string ID to identify the Chapter. Use for .",
		Type:       ebml.UTF8,
	},
	0x91: {
		Name:       "ChapterTimeStart",
		Definition: "Timestamp of the start of Chapter (not scaled).",
		Type:       ebml.Uinteger,
	},
	0x92: {
		Name:       "ChapterTimeEnd",
		Definition: "Timestamp of the end of Chapter (timestamp excluded, not scaled).",
		Type:       ebml.Uinteger,
	},
	0x98: {
		Name:       "ChapterFlagHidden",
		Definition: "If a chapter is hidden (1), it SHOULD NOT be available to the user interface (but still to Control Tracks; see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x4598: {
		Name:       "ChapterFlagEnabled",
		Definition: "Specify whether the chapter is enabled. It can be enabled/disabled by a Control Track. When disabled, the movie SHOULD skip all the content between the TimeStart and TimeEnd of this chapter (see ). (1 bit)",
		Type:       ebml.Uinteger,
	},
	0x6e67: {
		Name:       "ChapterSegmentUID",
		Definition: "The SegmentUID of another Segment to play during this chapter.",
		Type:       ebml.Binary,
	},
	0x6ebc: {
		Name:       "ChapterSegmentEditionUID",
		Definition: "The EditionUID to play from the Segment linked in ChapterSegmentUID. If ChapterSegmentEditionUID is undeclared then no Edition of the linked Segment is used.",
		Type:       ebml.Uinteger,
	},
	0x63c3: {
		Name:       "ChapterPhysicalEquiv",
		Definition: "Specify the physical equivalent of this ChapterAtom like \"DVD\" (60) or \"SIDE\" (50), see .",
		Type:       ebml.Uinteger,
	},
	0x8f: {
		Name:       "ChapterTrack",
		Definition: "List of tracks on which the chapter applies. If this Element is not present, all tracks apply",
		Type:       ebml.Master, Tag: mkvChapterTrack,
	},
	0x80: {
		Name:       "ChapterDisplay",
		Definition: "Contains all possible strings to use for the chapter display.",
		Type:       ebml.Master, Tag: mkvChapterDisplay,
	},
	0x6944: {
		Name:       "ChapProcess",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: mkvChapProcess,
	},
}

var mkvChapterTrack = ebml.Tag{
	0x89: {
		Name:       "ChapterTrackUID",
		Definition: "UID of the Track to apply this chapter too. In the absence of a control track, choosing this chapter will select the listed Tracks and deselect unlisted tracks. Absence of this Element indicates that the Chapter SHOULD be applied to any currently used Tracks.",
		Type:       ebml.Uinteger,
	},
}

var mkvChapterDisplay = ebml.Tag{
	0x85: {
		Name:       "ChapString",
		Definition: "Contains the string to use as the chapter atom.",
		Type:       ebml.UTF8,
	},
	0x437c: {
		Name:       "ChapLanguage",
		Definition: "The languages corresponding to the string, in the . This Element MUST be ignored if the ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
	0x437d: {
		Name:       "ChapLanguageIETF",
		Definition: "Specifies the language used in the ChapString according to  and using the . If this Element is used, then any ChapLanguage Elements used in the same ChapterDisplay MUST be ignored.",
		Type:       ebml.String,
	},
	0x437e: {
		Name:       "ChapCountry",
		Definition: "The countries corresponding to the string, same 2 octets as in . This Element MUST be ignored if the ChapLanguageIETF Element is used within the same ChapterDisplay Element.",
		Type:       ebml.String,
	},
}

var mkvChapProcess = ebml.Tag{
	0x6955: {
		Name:       "ChapProcessCodecID",
		Definition: "Contains the type of the codec used for the processing. A value of 0 means native Matroska processing (to be defined), a value of 1 means the  command set is used. More codec IDs can be added later.",
		Type:       ebml.Uinteger,
	},
	0x450d: {
		Name:       "ChapProcessPrivate",
		Definition: "Some optional data attached to the ChapProcessCodecID information. , it is the \"DVD level\" equivalent.",
		Type:       ebml.Binary,
	},
	0x6911: {
		Name:       "ChapProcessCommand",
		Definition: "Contains all the commands associated to the Atom.",
		Type:       ebml.Master, Tag: mkvChapProcessCommand,
	},
}

var mkvChapProcessCommand = ebml.Tag{
	0x6922: {
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
	0x6933: {
		Name:       "ChapProcessData",
		Definition: "Contains the command information. The data SHOULD be interpreted depending on the ChapProcessCodecID value. , the data correspond to the binary DVD cell pre/post commands.",
		Type:       ebml.Binary,
	},
}

var mkvTags = ebml.Tag{
	0x7373: {
		Name:       "Tag",
		Definition: "A single metadata descriptor.",
		Type:       ebml.Master, Tag: mkvTag,
	},
}

var mkvTag = ebml.Tag{
	0x63c0: {
		Name:       "Targets",
		Definition: "Specifies which other elements the metadata represented by the Tag applies to. If empty or not present, then the Tag describes everything in the Segment.",
		Type:       ebml.Master, Tag: mkvTargets,
	},
	0x67c8: {
		Name:       "SimpleTag",
		Definition: "Contains general information about the target.",
		Type:       ebml.Master, Tag: mkvSimpleTag,
	},
}

var mkvTargets = ebml.Tag{
	0x68ca: {
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
	0x63ca: {
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
	0x63c5: {
		Name:       "TagTrackUID",
		Definition: "A unique ID to identify the Track(s) the tags belong to. If the value is 0 at this level, the tags apply to all tracks in the Segment.",
		Type:       ebml.Uinteger,
	},
	0x63c9: {
		Name:       "TagEditionUID",
		Definition: "A unique ID to identify the EditionEntry(s) the tags belong to. If the value is 0 at this level, the tags apply to all editions in the Segment.",
		Type:       ebml.Uinteger,
	},
	0x63c4: {
		Name:       "TagChapterUID",
		Definition: "A unique ID to identify the Chapter(s) the tags belong to. If the value is 0 at this level, the tags apply to all chapters in the Segment.",
		Type:       ebml.Uinteger,
	},
	0x63c6: {
		Name:       "TagAttachmentUID",
		Definition: "A unique ID to identify the Attachment(s) the tags belong to. If the value is 0 at this level, the tags apply to all the attachments in the Segment.",
		Type:       ebml.Uinteger,
	},
}

var mkvSimpleTag = ebml.Tag{
	0x45a3: {
		Name:       "TagName",
		Definition: "The name of the Tag that is going to be stored.",
		Type:       ebml.UTF8,
	},
	0x447a: {
		Name:       "TagLanguage",
		Definition: "Specifies the language of the tag specified, in the . This Element MUST be ignored if the TagLanguageIETF Element is used within the same SimpleTag Element.",
		Type:       ebml.String,
	},
	0x447b: {
		Name:       "TagLanguageIETF",
		Definition: "Specifies the language used in the TagString according to  and using the . If this Element is used, then any TagLanguage Elements used in the same SimpleTag MUST be ignored.",
		Type:       ebml.String,
	},
	0x4484: {
		Name:       "TagDefault",
		Definition: "A boolean value to indicate if this is the default/original language to use for the given tag.",
		Type:       ebml.Uinteger,
	},
	0x4487: {
		Name:       "TagString",
		Definition: "The value of the Tag.",
		Type:       ebml.UTF8,
	},
	0x4485: {
		Name:       "TagBinary",
		Definition: "The values of the Tag if it is binary. Note that this cannot be used in the same SimpleTag as TagString.",
		Type:       ebml.Binary,
	},
}
