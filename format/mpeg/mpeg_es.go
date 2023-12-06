package mpeg

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var mpegASCGroup decode.Group
var vorbisPacketGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.MPEG_ES,
		&decode.Format{
			Description: "MPEG Elementary Stream",
			DecodeFn:    esDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.MPEG_ASC}, Out: &mpegASCGroup},
				{Groups: []*decode.Group{format.Vorbis_Packet}, Out: &vorbisPacketGroup},
			},
		})
}

const (
	Forbidden0                          = 0x00
	ObjectDescrTag                      = 0x01
	InitialObjectDescrTag               = 0x02
	ES_DescrTag                         = 0x03
	DecoderConfigDescrTag               = 0x04
	DecSpecificInfoTag                  = 0x05
	SLConfigDescrTag                    = 0x06
	ContentIdentDescrTag                = 0x07
	SupplContentIdentDescrTag           = 0x08
	IPI_DescrPointerTag                 = 0x09
	IPMP_DescrPointerTag                = 0x0A
	IPMP_DescrTag                       = 0x0B
	QoS_DescrTag                        = 0x0C
	RegistrationDescrTag                = 0x0D
	ES_ID_IncTag                        = 0x0E
	ES_ID_RefTag                        = 0x0F
	MP4_IOD_Tag                         = 0x10
	MP4_OD_Tag                          = 0x11
	IPL_DescrPointerRefTag              = 0x12
	ExtensionProfileLevelDescrTag       = 0x13
	profileLevelIndicationIndexDescrTag = 0x14
	ContentClassificationDescrTag       = 0x40
	KeyWordDescrTag                     = 0x41
	RatingDescrTag                      = 0x42
	LanguageDescrTag                    = 0x43
	ShortTextualDescrTag                = 0x44
	ExpandedTextualDescrTag             = 0x45
	ContentCreatorNameDescrTag          = 0x46
	ContentCreationDateDescrTag         = 0x47
	OCICreatorNameDescrTag              = 0x48
	OCICreationDateDescrTag             = 0x49
	SmpteCameraPositionDescrTag         = 0x4A
	SegmentDescrTag                     = 0x4B
	MediaTimeDescrTag                   = 0x4C
	IPMP_ToolsListDescrTag              = 0x60
	IPMP_ToolTag                        = 0x61
	M4MuxTimingDescrTag                 = 0x62
	M4MuxCodeTableDescrTag              = 0x63
	ExtSLConfigDescrTag                 = 0x64
	M4MuxBufferSizeDescrTag             = 0x65
	M4MuxIdentDescrTag                  = 0x66
	DependencyPointerTag                = 0x67
	DependencyMarkerTag                 = 0x68
	M4MuxChannelDescrTag                = 0x69
	Forbidden1                          = 0xFF
)

var odTagNames = scalar.UintMapSymStr{
	Forbidden0:                          "Forbidden",
	ObjectDescrTag:                      "ObjectDescrTag",
	InitialObjectDescrTag:               "InitialObjectDescrTag",
	ES_DescrTag:                         "ES_DescrTag",
	DecoderConfigDescrTag:               "DecoderConfigDescrTag",
	DecSpecificInfoTag:                  "DecSpecificInfoTag",
	SLConfigDescrTag:                    "SLConfigDescrTag",
	ContentIdentDescrTag:                "ContentIdentDescrTag",
	SupplContentIdentDescrTag:           "SupplContentIdentDescrTag",
	IPI_DescrPointerTag:                 "IPI_DescrPointerTag",
	IPMP_DescrPointerTag:                "IPMP_DescrPointerTag",
	IPMP_DescrTag:                       "IPMP_DescrTag",
	QoS_DescrTag:                        "QoS_DescrTag",
	RegistrationDescrTag:                "RegistrationDescrTag",
	ES_ID_IncTag:                        "ES_ID_IncTag",
	ES_ID_RefTag:                        "ES_ID_RefTag",
	MP4_IOD_Tag:                         "MP4_IOD_Tag",
	MP4_OD_Tag:                          "MP4_OD_Tag",
	IPL_DescrPointerRefTag:              "IPL_DescrPointerRefTag",
	ExtensionProfileLevelDescrTag:       "ExtensionProfileLevelDescrTag",
	profileLevelIndicationIndexDescrTag: "profileLevelIndicationIndexDescrTag",
	ContentClassificationDescrTag:       "ContentClassificationDescrTag",
	KeyWordDescrTag:                     "KeyWordDescrTag",
	RatingDescrTag:                      "RatingDescrTag",
	LanguageDescrTag:                    "LanguageDescrTag",
	ShortTextualDescrTag:                "ShortTextualDescrTag",
	ExpandedTextualDescrTag:             "ExpandedTextualDescrTag",
	ContentCreatorNameDescrTag:          "ContentCreatorNameDescrTag",
	ContentCreationDateDescrTag:         "ContentCreationDateDescrTag",
	OCICreatorNameDescrTag:              "OCICreatorNameDescrTag",
	OCICreationDateDescrTag:             "OCICreationDateDescrTag",
	SmpteCameraPositionDescrTag:         "SmpteCameraPositionDescrTag",
	SegmentDescrTag:                     "SegmentDescrTag",
	MediaTimeDescrTag:                   "MediaTimeDescrTag",
	IPMP_ToolsListDescrTag:              "IPMP_ToolsListDescrTag",
	IPMP_ToolTag:                        "IPMP_ToolTag",
	M4MuxTimingDescrTag:                 "M4MuxTimingDescrTag",
	M4MuxCodeTableDescrTag:              "M4MuxCodeTableDescrTag",
	ExtSLConfigDescrTag:                 "ExtSLConfigDescrTag",
	M4MuxBufferSizeDescrTag:             "M4MuxBufferSizeDescrTag",
	M4MuxIdentDescrTag:                  "M4MuxIdentDescrTag",
	DependencyPointerTag:                "DependencyPointerTag",
	DependencyMarkerTag:                 "DependencyMarkerTag",
	M4MuxChannelDescrTag:                "M4MuxChannelDescrTag",
	Forbidden1:                          "Forbidden",
}

const (
	Forbidden               = 0x00
	ObjectDescriptorStream  = 0x01
	ClockReferenceStream    = 0x02
	SceneDescriptionStream  = 0x03
	VisualStream            = 0x04
	AudioStream             = 0x05
	MPEG7Stream             = 0x06
	IPMPStream              = 0x07
	ObjectContentInfoStream = 0x08
	MPEGJStream             = 0x09
	InteractionStream       = 0x0A
	IPMPToolStream          = 0x0B
)

var streamTypeNames = scalar.UintMapSymStr{
	Forbidden:               "Forbidden",
	ObjectDescriptorStream:  "ObjectDescriptorStream",
	ClockReferenceStream:    "ClockReferenceStream",
	SceneDescriptionStream:  "SceneDescriptionStream",
	VisualStream:            "VisualStream",
	AudioStream:             "AudioStream",
	MPEG7Stream:             "MPEG7Stream",
	IPMPStream:              "IPMPStream",
	ObjectContentInfoStream: "ObjectContentInfoStream",
	MPEGJStream:             "MPEGJStream",
	InteractionStream:       "InteractionStream",
	IPMPToolStream:          "IPMPToolStream",
}

func esLengthEncoding(d *decode.D) uint64 {
	v := uint64(0)
	nextByte := true
	for nextByte {
		nextByte = d.Bool()
		v = v<<7 | d.U7()
	}
	return v
}

//nolint:unparam
func fieldODDecodeTag(d *decode.D, edc *esDecodeContext, name string, expectedTagID int, fn func(d *decode.D)) {
	d.FieldStruct(name, func(d *decode.D) {
		odDecodeTag(d, edc, expectedTagID, fn)
	})
}

type esDecodeContext struct {
	currentDecoderConfig *format.MPEG_Decoder_Config
	decoderConfigs       []format.MPEG_Decoder_Config
}

func odDecodeTag(d *decode.D, edc *esDecodeContext, _ int, fn func(d *decode.D)) {
	odDecoders := map[uint64]func(d *decode.D){
		ES_DescrTag: func(d *decode.D) {
			d.FieldU16("es_id")
			streamDependencyFlag := d.FieldBool("stream_dependency_flag")
			urlFlag := d.FieldBool("url_flag")
			ocrStreamFlag := d.FieldBool("ocr_stream_flag")
			d.FieldU5("stream_priority")
			if streamDependencyFlag {
				d.FieldU16("dependency_on_es_id")
			}
			if urlFlag {
				urlLen := d.FieldU8("url_length")
				d.FieldUTF8("url", int(urlLen))
			}
			if ocrStreamFlag {
				d.FieldU16("ocr_es_id")
			}
			fieldODDecodeTag(d, edc, "dec_config_descr", -1, nil)
			fieldODDecodeTag(d, edc, "sl_config_descr", -1, nil)
		},
		DecoderConfigDescrTag: func(d *decode.D) {
			objectType := d.FieldU8("object_type_indication", format.MpegObjectTypeNames)
			edc.decoderConfigs = append(edc.decoderConfigs, format.MPEG_Decoder_Config{
				ObjectType: int(objectType),
			})
			edc.currentDecoderConfig = &edc.decoderConfigs[len(edc.decoderConfigs)-1]

			d.FieldU6("stream_type", streamTypeNames)
			d.FieldBool("upstream")
			d.FieldBool("specific_info_flag")
			d.FieldU24("buffer_size_db")
			d.FieldU32("max_bit_rate")
			d.FieldU32("avg_bit_rate")

			// TODO: seems to be optional
			if d.BitsLeft() > 0 {
				switch objectType {
				// TODO: vorbis is special case?
				case format.MPEGObjectTypeVORBIS:
					fieldODDecodeTag(d, edc, "decoder_specific_info", -1, func(d *decode.D) {
						numPackets := d.FieldU8("num_packets")
						// TODO: lacing
						packetLengths := []int64{}
						// Xiph-style lacing (similar to ogg) of n-1 packets, last is reset of block
						d.FieldArray("laces", func(d *decode.D) {
							for i := uint64(0); i < numPackets; i++ {
								l := d.FieldUintFn("lace", func(d *decode.D) uint64 {
									var l uint64
									for {
										n := d.U8()
										l += n
										if n < 255 {
											return l
										}
									}
								})
								packetLengths = append(packetLengths, int64(l))
							}
						})
						d.FieldArray("packets", func(d *decode.D) {
							for _, l := range packetLengths {
								d.FieldFormatLen("packet", l*8, &vorbisPacketGroup, nil)
							}
							d.FieldFormatLen("packet", d.BitsLeft(), &vorbisPacketGroup, nil)
						})
					})
				default:
					switch format.MpegObjectTypeStreamType[objectType] {
					case format.MPEGStreamTypeAudio:
						fieldODDecodeTag(d, edc, "decoder_specific_info", -1, func(d *decode.D) {
							_, v := d.FieldFormat("audio_specific_config", &mpegASCGroup, nil)
							mpegASCout, ok := v.(format.MPEG_ASC_Out)
							if !ok {
								panic(fmt.Sprintf("expected MPEGASCOut got %#+v", v))
							}
							if edc.currentDecoderConfig != nil {
								edc.currentDecoderConfig.ASCObjectType = mpegASCout.ObjectType
							}
						})
					case format.MPEGStreamTypeVideo:
						// TODO:
					case format.MPEGStreamTypeText:
						// TODO:
					}
				}
			}
		},
		// TODO: where is the spec?
		// https://xhelmboyx.tripod.com/formats/mp4-layout.txt
		// https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/movenc.c mov_write_iods_tag
		// TODO: IOD = initial object descriptor?
		MP4_IOD_Tag: func(d *decode.D) {
			d.FieldU16("od_id")
			d.FieldU8("od_profile_level")
			d.FieldU8("scene_profile_level")
			d.FieldU8("audio_profile_level")
			d.FieldU8("video_profile_level")
			d.FieldU8("graphics_profile_level")
			d.FieldArray("descriptors", func(d *decode.D) {
				for !d.End() {
					fieldODDecodeTag(d, edc, "desc", -1, nil)
				}
			})
		},
		ES_ID_IncTag: func(d *decode.D) {
			d.FieldU32("track_id")
		},
	}

	// TODO: expectedTagID

	tagID := d.FieldU8("tag_id", odTagNames)
	tagLen := d.FieldUintFn("length", esLengthEncoding)

	if fn != nil {
		d.FramedFn(int64(tagLen)*8, fn)
	} else if tagDecoder, ok := odDecoders[tagID]; ok {
		d.FramedFn(int64(tagLen)*8, tagDecoder)
	} else {
		d.FieldRawLen("data", d.BitsLeft())
	}
}

func esDecode(d *decode.D) any {
	var edc esDecodeContext
	odDecodeTag(d, &edc, -1, nil)
	return format.MPEG_ES_Out{DecoderConfigs: edc.decoderConfigs}
}
