package rtmp

// https://rtmp.veriskope.com/docs/spec/

// TODO: audio/video message, coded header?
// TODO: split to rtmp/rtmp_message?
// TODO: format options, set default chunk size?
// TODO: support to skip handshake?

import (
	"bytes"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

var amf0Group decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.RTMP,
		Description: "Real-Time Messaging Protocol",
		Groups: []string{
			format.TCP_STREAM,
		},
		DecodeFn: rtmpDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AMF0}, Group: &amf0Group},
		},
	})
}

// from spec
const defaultChunkSize = 128

// names from RTMP spec
const (
	messageTypeSetChunkSize                = 1
	messageTypeAbortMessage                = 2
	messageTypeAcknowledgment              = 3
	messageTypeUserControlMessage          = 4
	messageTypeWindowAcknowledgementSize   = 5
	messageTypeSetPeerBandwidth            = 6
	messageTypeVirtualControl              = 7 // TODO: not in spec?
	messageTypeAudioMessage                = 8
	messageTypeVideoMessage                = 9
	messageTypeDataMessageExtended         = 15
	messageTypeSharedObjectMessageExtended = 16
	messageTypeCommandMessageExtended      = 17
	messageTypeDataMessage                 = 18
	messageTypeSharedObjectMessage         = 19
	messageTypeCommandMessage              = 20
	messageTypeUDP                         = 21 // TODO: not in spec?
	messageTypeAggregateMessage            = 22
	messageTypePresent                     = 23 // TODO: not in spec?
)

var rtmpMessageTypeIDNames = scalar.UToSymStr{
	messageTypeSetChunkSize:                "set_chunk_size",
	messageTypeAbortMessage:                "abort_message",
	messageTypeAcknowledgment:              "acknowledgment",
	messageTypeUserControlMessage:          "user_control_message",
	messageTypeWindowAcknowledgementSize:   "window_acknowledgement_size",
	messageTypeSetPeerBandwidth:            "set_peer_bandwidth",
	messageTypeVirtualControl:              "virtual_control",
	messageTypeAudioMessage:                "audio_message",
	messageTypeVideoMessage:                "video_message",
	messageTypeDataMessageExtended:         "data_message_extended",
	messageTypeSharedObjectMessageExtended: "shared_object_message_extended",
	messageTypeCommandMessageExtended:      "command_message_extended",
	messageTypeDataMessage:                 "data_message",
	messageTypeSharedObjectMessage:         "shared_object_message",
	messageTypeCommandMessage:              "command_message",
	messageTypeUDP:                         "udp",
	messageTypeAggregateMessage:            "aggregate_message",
	messageTypePresent:                     "present",
}

const (
	userControlEvenTypeStreamBegin      = 0
	userControlEvenTypeStreamEOF        = 1
	userControlEvenTypeStreamDry        = 2
	userControlEvenTypeSetBufferLength  = 3
	userControlEvenTypeStreamIsRecorded = 4
	userControlEvenTypePingRequest      = 6
	userControlEvenTypePingResponse     = 7
)

var userControlEvenTypNames = scalar.UToSymStr{
	userControlEvenTypeStreamBegin:      "stream_begin",
	userControlEvenTypeStreamEOF:        "stream_eof",
	userControlEvenTypeStreamDry:        "stream_dry",
	userControlEvenTypeSetBufferLength:  "set_buffer_length",
	userControlEvenTypeStreamIsRecorded: "stream_is_recorded",
	userControlEvenTypePingRequest:      "ping_request",
	userControlEvenTypePingResponse:     "ping_response",
}

var setPeerBandwidthLimitTypeName = scalar.UToSymStr{
	0: "hard",
	1: "soft",
	2: "dynamic",
}

const timestampExtended = 0xff_ff_ff

var timestampDescription = scalar.UToScalar{
	timestampExtended: scalar.S{Description: "extended"},
}

// based on https://github.com/wireshark/wireshark/blob/master/epan/dissectors/packet-rtmpt.c
// which in turn is based on rtmp and swf specifications and FLV v10.1 section E.4.3.1
var audioMessageCodecNames = scalar.UToSymStr{
	0:  "uncompressed",
	1:  "adpcm",
	2:  "mp3",
	3:  "uncompressed_le",
	4:  "nellymoser_16khz",
	5:  "nellymoser_8khz",
	6:  "nellymoser",
	7:  "g711a",
	8:  "g711u",
	9:  "nellymoser_16khz",
	10: "he-aac",
	11: "speex",
}

var audioMessageRateNames = scalar.UToSymU{
	0: 5500,
	1: 11000,
	2: 22000, // TODO: 22050?
	3: 44000, // TODO: 44100?
}

var audioMessageSampleSize = scalar.UToSymU{
	0: 8,
	1: 16,
}

var audioMessageChannels = scalar.UToSymU{
	0: 1,
	1: 2,
}

var videoMessageTypeNames = scalar.UToSymStr{
	1: "keyframe",
	2: "inter_frame",
	3: "disposable_inter_frame",
	4: "generated_key_frame",
	5: "video_info_or_command_frame",
}

var videoMessageCodecNames = scalar.UToSymStr{
	2: "sorensen_h263",
	3: "screen_video",
	4: "on2_vp6",
	5: "on2_vp6_alpha",
	6: "screen_video_version_2",
	7: "h264",
}

// TODO: invalid warning that timestampDelta is unused
//nolint: structcheck,unused
type messageHeader struct {
	timestamp       uint64
	timestampDelta  uint64
	messageStreamID uint64
	messageLength   uint64
	messageTypeID   uint64
}

func rtmpDecodeMessageType(d *decode.D, typ int, chunkSize *int) {
	switch typ {
	case messageTypeSetChunkSize:
		// TODO: zero bit, verify size? message size is 24 bit
		*chunkSize = int(d.FieldU32("size"))
	case messageTypeAbortMessage:
		d.FieldU32("chunk_stream_id")
	case messageTypeAcknowledgment:
		d.FieldU32("sequence_number")
	case messageTypeUserControlMessage:
		d.FieldU16("type", userControlEvenTypNames)
		d.FieldRawLen("data", d.BitsLeft())
	case messageTypeWindowAcknowledgementSize:
		d.FieldU32("window_size")
	case messageTypeSetPeerBandwidth:
		d.FieldU32("chunk_size")
		d.FieldU8("limit_type", setPeerBandwidthLimitTypeName)
	case messageTypeDataMessage:
		d.FieldFormat("message", amf0Group, nil)
	case messageTypeCommandMessage:
		d.FieldFormat("command_name", amf0Group, nil)
		d.FieldFormat("transaction_id", amf0Group, nil)
		d.FieldFormat("command_object", amf0Group, nil)
		d.FieldArray("arguments", func(d *decode.D) {
			for !d.End() {
				d.FieldFormat("argument", amf0Group, nil)
			}
		})
	case messageTypeAggregateMessage:
		d.FieldArray("messages", func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("message", func(d *decode.D) {
					var h messageHeader
					h.messageTypeID = d.FieldU8("message_type_id", rtmpMessageTypeIDNames)
					h.messageLength = d.FieldU24("message_length")
					h.timestamp = d.FieldU32("timestamp", timestampDescription)
					h.messageStreamID = d.FieldU24("message_stream_id")
					// TODO: possible to set chunk size in aggregated message?
					d.FramedFn(int64(h.messageLength*8), func(d *decode.D) {
						rtmpDecodeMessageType(d, int(h.messageTypeID), chunkSize)
					})
					d.FieldU32("back_pointer")
				})
			}
		})
	case messageTypeAudioMessage:
		if d.BitsLeft() == 0 {
			return
		}
		d.FieldU4("codec", audioMessageCodecNames)
		d.FieldU2("sample_rate", audioMessageRateNames)
		d.FieldU1("sample_size", audioMessageSampleSize)
		d.FieldU1("channels", audioMessageChannels)
		d.FieldRawLen("data", d.BitsLeft())
	case messageTypeVideoMessage:
		if d.BitsLeft() == 0 {
			return
		}
		d.FieldU4("type", videoMessageTypeNames)
		d.FieldU4("codec", videoMessageCodecNames)
		// TODO: flv header + h263 format?
		// TODO: ffmpeg rtmp proto seems to recrate a flv stream and demux it
		d.FieldRawLen("data", d.BitsLeft())
	default:
		d.FieldRawLen("data", d.BitsLeft())
	}
}

func rtmpDecode(d *decode.D, in interface{}) interface{} {
	var isClient bool

	if tsi, ok := in.(format.TCPStreamIn); ok {
		if tsi.DestinationPort != format.TCPPortRTMP {
			d.Fatalf("wrong port")
		}
		isClient = tsi.IsClient
	}

	// chunk size is global for one direction
	chunkSize := defaultChunkSize

	name := "s"
	if isClient {
		name = "c"
	}
	// TODO: 1536 byte blob instead?
	d.FieldStruct("handshake", func(d *decode.D) {
		d.FieldStruct(name+"0", func(d *decode.D) {
			d.FieldU8("version")
		})
		d.FieldStruct(name+"1", func(d *decode.D) {
			d.FieldU32("time")
			d.FieldU32("zero")
			d.FieldRawLen("random", 1528*8)
		})
		d.FieldStruct(name+"2", func(d *decode.D) {
			d.FieldU32("time")
			d.FieldU32("time2")
			d.FieldRawLen("random", 1528*8)
		})
	})

	type messageHeader struct {
		timestamp       uint64
		timestampDelta  uint64
		messageStreamID uint64
		messageLength   uint64
		messageTypeID   uint64
	}

	type message struct {
		l   uint64
		b   bytes.Buffer
		typ uint64
	}
	type chunkStream struct {
		messageSteams map[uint64]*message
		prevHeader    messageHeader
	}

	chunkStreams := map[uint64]*chunkStream{}

	messages := d.FieldArrayValue("messages")

	d.FieldArray("chunks", func(d *decode.D) {
		for !d.End() {
			d.FieldStruct("chunk", func(d *decode.D) {
				var chunkSteamID uint64

				fmt := d.FieldU2("fmt")
				switch d.PeekBits(6) {
				case 0:
					// 64-319: 1 byte
					chunkSteamID = d.FieldU14("chunk_stream_id", scalar.UAdd(64))
				case 1:
					// 64-65599: 1 byte
					chunkSteamID = d.FieldU30("chunk_stream_id", scalar.UAdd(64))
				default:
					// 2-63: 1 byte
					chunkSteamID = d.FieldU6("chunk_stream_id")
				}

				cs, ok := chunkStreams[chunkSteamID]
				if !ok {
					cs = &chunkStream{
						messageSteams: map[uint64]*message{},
					}
					chunkStreams[chunkSteamID] = cs
				}

				var h messageHeader

				switch fmt {
				case 0:
					h.timestamp = d.FieldU24("timestamp", timestampDescription)
					h.messageLength = d.FieldU24("message_length")
					h.messageTypeID = d.FieldU8("message_type_id", rtmpMessageTypeIDNames)
					h.messageStreamID = d.FieldU32LE("message_stream_id")
					if h.timestamp == timestampExtended {
						h.timestamp = d.FieldU32("extended_timestamp")
					}
				case 1:
					h.timestampDelta = d.FieldU24("timestamp_delta", timestampDescription)
					h.messageLength = d.FieldU24("message_length")
					h.messageTypeID = d.FieldU8("message_type_id", rtmpMessageTypeIDNames)
					if h.timestamp == timestampExtended {
						h.timestampDelta = d.FieldU32("extended_timestamp")
					}
					h.timestamp = cs.prevHeader.timestamp
					h.messageStreamID = cs.prevHeader.messageStreamID
					d.FieldValueU("message_stream_id", h.messageStreamID, scalar.Description("previous"))

				case 2:
					h.timestampDelta = d.FieldU24("timestamp_delta", timestampDescription)
					if h.timestamp == timestampExtended {
						h.timestampDelta = d.FieldU32("extended_timestamp")
					}
					h.timestamp = cs.prevHeader.timestamp
					h.messageLength = cs.prevHeader.messageLength
					h.messageStreamID = cs.prevHeader.messageStreamID
					h.messageTypeID = cs.prevHeader.messageTypeID
					d.FieldValueU("message_length", h.messageLength, scalar.Description("previous"))
					d.FieldValueU("message_type_id", h.messageTypeID, scalar.Description("previous"))
					d.FieldValueU("message_stream_id", h.messageStreamID, scalar.Description("previous"))
				case 3:
					h.timestamp = cs.prevHeader.timestamp
					h.timestampDelta = cs.prevHeader.timestampDelta
					h.messageLength = cs.prevHeader.messageLength
					h.messageStreamID = cs.prevHeader.messageStreamID
					h.messageTypeID = cs.prevHeader.messageTypeID
					d.FieldValueU("message_length", h.messageLength, scalar.Description("previous"))
					d.FieldValueU("message_type_id", h.messageTypeID, scalar.Description("previous"))
					d.FieldValueU("message_stream_id", h.messageStreamID, scalar.Description("previous"))
				}

				h.timestamp += h.timestampDelta

				d.FieldValueU("calculated_timestamp", h.timestamp)

				m, ok := cs.messageSteams[h.messageStreamID]
				if !ok {
					m = &message{
						l:   h.messageLength,
						typ: h.messageTypeID,
					}
					cs.messageSteams[h.messageStreamID] = m
				}

				payloadLength := chunkSize
				left := int(m.l) - m.b.Len()
				if left < payloadLength {
					payloadLength = left
				}

				if payloadLength > 0 {
					d.MustCopyBits(&m.b, d.FieldRawLen("data", int64(payloadLength)*8))
				}

				if m.l == uint64(m.b.Len()) {
					messageBR := bitio.NewBitReader(m.b.Bytes(), -1)
					messages.FieldStructRootBitBufFn("message", messageBR, func(d *decode.D) {
						d.FieldValueU("message_stream_id", h.messageStreamID)
						d.FieldValueU("message_type_id", m.typ, rtmpMessageTypeIDNames)
						rtmpDecodeMessageType(d, int(m.typ), &chunkSize)
					})

					// delete so that we create a new message{} with a new bytes.Buffer to
					// not share byte slice
					delete(cs.messageSteams, h.messageStreamID)
				}

				cs.prevHeader = h
			})
		}
	})

	return nil
}
