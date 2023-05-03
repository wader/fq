package rtmp

// https://rtmp.veriskope.com/docs/spec/
// https://rtmp.veriskope.com/pdf/video_file_format_spec_v10.pdf

// TODO: split to rtmp/rtmp_message?
// TODO: support to skip handshake?
// TODO: keep track of message stream format, decode aac etc

import (
	"bytes"
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var rtmpAmf0Group decode.Group
var rtmpMpegASCGroup decode.Group

//go:embed rtmp.md
var rtmpFS embed.FS

func init() {
	interp.RegisterFormat(
		format.RTMP,
		&decode.Format{
			Description: "Real-Time Messaging Protocol",
			Groups: []*decode.Group{
				format.TCP_Stream,
			},
			DecodeFn: rtmpDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.AMF0}, Out: &rtmpAmf0Group},
				{Groups: []*decode.Group{format.MPEG_ASC}, Out: &rtmpMpegASCGroup},
			},
		})
	interp.RegisterFS(rtmpFS)
}

// from RTMP spec
const defaultChunkSize = 128

// names from RTMP spec
const (
	messageTypeSetChunkSize                = 1
	messageTypeAbortMessage                = 2
	messageTypeAcknowledgment              = 3
	messageTypeUserControlMessage          = 4
	messageTypeWindowAcknowledgementSize   = 5
	messageTypeSetPeerBandwidth            = 6
	messageTypeVirtualControl              = 7 // TODO: not in spec but in wikipedia article
	messageTypeAudioMessage                = 8
	messageTypeVideoMessage                = 9
	messageTypeDataMessageExtended         = 15
	messageTypeSharedObjectMessageExtended = 16
	messageTypeCommandMessageExtended      = 17
	messageTypeDataMessage                 = 18
	messageTypeSharedObjectMessage         = 19
	messageTypeCommandMessage              = 20
	messageTypeUDP                         = 21 // TODO: not in spec but in wikipedia article
	messageTypeAggregateMessage            = 22
	messageTypePresent                     = 23 // TODO: not in spec but in wikipedia article
)

var rtmpMessageTypeIDNames = scalar.UintMapSymStr{
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

var userControlEvenTypNames = scalar.UintMapSymStr{
	userControlEvenTypeStreamBegin:      "stream_begin",
	userControlEvenTypeStreamEOF:        "stream_eof",
	userControlEvenTypeStreamDry:        "stream_dry",
	userControlEvenTypeSetBufferLength:  "set_buffer_length",
	userControlEvenTypeStreamIsRecorded: "stream_is_recorded",
	userControlEvenTypePingRequest:      "ping_request",
	userControlEvenTypePingResponse:     "ping_response",
}

var setPeerBandwidthLimitTypeName = scalar.UintMapSymStr{
	0: "hard",
	1: "soft",
	2: "dynamic",
}

const timestampExtended = 0xff_ff_ff

var timestampUintDescription = scalar.UintMapDescription{
	timestampExtended: "extended",
}

const (
	audioMessageCodecAAC = 10
)

// based on https://github.com/wireshark/wireshark/blob/master/epan/dissectors/packet-rtmpt.c
// which in turn is based on rtmp and swf specifications and FLV v10.1 section E.4.3.1
var audioMessageCodecNames = scalar.UintMapSymStr{
	0:                    "uncompressed",
	1:                    "adpcm",
	2:                    "mp3",
	3:                    "uncompressed_le",
	4:                    "nellymoser_16khz",
	5:                    "nellymoser_8khz",
	6:                    "nellymoser",
	7:                    "g711a",
	8:                    "g711u",
	9:                    "nellymoser_16khz",
	audioMessageCodecAAC: "aac",
	11:                   "speex",
}

const (
	audioMessageAACPacketTypeASC = 0
	audioMessageAACPacketTypeRaw = 1
)

var audioMessageAACPacketTypeNames = scalar.UintMapSymStr{
	audioMessageAACPacketTypeASC: "asc",
	audioMessageAACPacketTypeRaw: "raw",
}

var audioMessageRateNames = scalar.UintMapSymUint{
	0: 5500,
	1: 11025,
	2: 22050,
	3: 44100,
}

var audioMessageSampleSize = scalar.UintMapSymUint{
	0: 8,
	1: 16,
}

var audioMessageChannels = scalar.UintMapSymUint{
	0: 1,
	1: 2,
}

var videoMessageTypeNames = scalar.UintMapSymStr{
	1: "keyframe",
	2: "inter_frame",
	3: "disposable_inter_frame",
	4: "generated_key_frame",
	5: "video_info_or_command_frame",
}

const (
	videoMessageCodecH264 = 7
)

var videoMessageCodecNames = scalar.UintMapSymStr{
	2:                     "h263",
	3:                     "screen_video",
	4:                     "vp6",
	5:                     "vp6_alpha",
	6:                     "screen_video_v2",
	videoMessageCodecH264: "h264",
}

var videoMessageH264PacketTypeNames = scalar.UintMapSymStr{
	0: "dcr",
	1: "au", // TODO: is access unit?
	2: "empty",
}

// TODO: invalid warning that timestampDelta is unused

//nolint:unused
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
		typ := d.FieldU16("type", userControlEvenTypNames)
		switch typ {
		case userControlEvenTypeStreamBegin:
			d.FieldU32("stream_id")
		case userControlEvenTypeStreamEOF:
			d.FieldU32("stream_id")
		case userControlEvenTypeStreamDry:
			d.FieldU32("stream_id")
		case userControlEvenTypeSetBufferLength:
			d.FieldU32("stream_id")
			d.FieldU32("length")
		case userControlEvenTypeStreamIsRecorded:
			d.FieldU32("stream_id")
		case userControlEvenTypePingRequest:
			d.FieldU32("timestamp")
		case userControlEvenTypePingResponse:
			d.FieldU32("timestamp")
		default:
			d.FieldRawLen("data", d.BitsLeft())
		}
	case messageTypeWindowAcknowledgementSize:
		d.FieldU32("window_size")
	case messageTypeSetPeerBandwidth:
		d.FieldU32("chunk_size")
		d.FieldU8("limit_type", setPeerBandwidthLimitTypeName)
	case messageTypeDataMessage:
		d.FieldArray("messages", func(d *decode.D) {
			for !d.End() {
				d.FieldFormat("message", &rtmpAmf0Group, nil)
			}
		})
	case messageTypeCommandMessage:
		d.FieldFormat("command_name", &rtmpAmf0Group, nil)
		d.FieldFormat("transaction_id", &rtmpAmf0Group, nil)
		d.FieldFormat("command_object", &rtmpAmf0Group, nil)
		d.FieldArray("arguments", func(d *decode.D) {
			for !d.End() {
				d.FieldFormat("argument", &rtmpAmf0Group, nil)
			}
		})
	case messageTypeAggregateMessage:
		d.FieldArray("messages", func(d *decode.D) {
			for !d.End() {
				d.FieldStruct("message", func(d *decode.D) {
					var h messageHeader
					h.messageTypeID = d.FieldU8("message_type_id", rtmpMessageTypeIDNames)
					h.messageLength = d.FieldU24("message_length")
					h.timestamp = d.FieldU32("timestamp", timestampUintDescription)
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
		codec := d.FieldU4("codec", audioMessageCodecNames)
		d.FieldU2("sample_rate", audioMessageRateNames)
		d.FieldU1("sample_size", audioMessageSampleSize)
		d.FieldU1("channels", audioMessageChannels)
		if codec == audioMessageCodecAAC {
			switch d.FieldU8("type", audioMessageAACPacketTypeNames) {
			case audioMessageAACPacketTypeASC:
				d.FieldFormat("data", &rtmpMpegASCGroup, nil)
			default:
				d.FieldRawLen("data", d.BitsLeft())
			}
		} else {
			d.FieldRawLen("data", d.BitsLeft())
		}
	case messageTypeVideoMessage:
		if d.BitsLeft() == 0 {
			return
		}
		d.FieldU4("type", videoMessageTypeNames)
		codec := d.FieldU4("codec", videoMessageCodecNames)
		// TODO: flv header + h263 format?
		// TODO: ffmpeg rtmp proto seems to recrate a flv stream and demux it
		if codec == videoMessageCodecH264 {
			d.FieldU8("type", videoMessageH264PacketTypeNames)
		}

		d.FieldRawLen("data", d.BitsLeft())
	default:
		d.FieldRawLen("data", d.BitsLeft())
	}
}

func rtmpDecode(d *decode.D) any {
	var isClient bool

	var tsi format.TCP_Stream_In
	if d.ArgAs(&tsi) {
		tsi.MustIsPort(d.Fatalf, format.TCPPortRTMP)
		isClient = tsi.IsClient
	}

	// chunk size is global for one direction
	chunkSize := defaultChunkSize

	name := "s"
	if isClient {
		name = "c"
	}
	// TODO: 1536 byte blobs instead?
	d.FieldStruct("handshake", func(d *decode.D) {
		d.FieldStruct(name+"0", func(d *decode.D) {
			d.FieldU8("version")
		})
		d.FieldStruct(name+"1", func(d *decode.D) {
			d.FieldU32("time")
			d.FieldU32("zero") // TODO: does not seems to be zero sometimes?
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
				switch d.PeekUintBits(6) {
				case 0:
					// 64-319: 2 byte
					d.FieldU6("chunk_stream_id_prefix")
					chunkSteamID = d.FieldU8("chunk_stream_id", scalar.UintActualAdd(64))
				case 1:
					// 64-65599: 3 byte
					d.FieldU6("chunk_stream_id_prefix")
					chunkSteamID = d.FieldU16("chunk_stream_id", scalar.UintActualAdd(64))
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
					h.timestamp = d.FieldU24("timestamp", timestampUintDescription)
					h.messageLength = d.FieldU24("message_length")
					h.messageTypeID = d.FieldU8("message_type_id", rtmpMessageTypeIDNames)
					h.messageStreamID = d.FieldU32LE("message_stream_id")
					if h.timestamp == timestampExtended {
						h.timestamp = d.FieldU32("extended_timestamp")
					}
				case 1:
					h.timestampDelta = d.FieldU24("timestamp_delta", timestampUintDescription)
					h.messageLength = d.FieldU24("message_length")
					h.messageTypeID = d.FieldU8("message_type_id", rtmpMessageTypeIDNames)
					if h.timestamp == timestampExtended {
						h.timestampDelta = d.FieldU32("extended_timestamp")
					}
					h.timestamp = cs.prevHeader.timestamp
					h.messageStreamID = cs.prevHeader.messageStreamID
					d.FieldValueUint("message_stream_id", h.messageStreamID, scalar.UintDescription("previous"))

				case 2:
					h.timestampDelta = d.FieldU24("timestamp_delta", timestampUintDescription)
					if h.timestamp == timestampExtended {
						h.timestampDelta = d.FieldU32("extended_timestamp")
					}
					h.timestamp = cs.prevHeader.timestamp
					h.messageLength = cs.prevHeader.messageLength
					h.messageStreamID = cs.prevHeader.messageStreamID
					h.messageTypeID = cs.prevHeader.messageTypeID
					d.FieldValueUint("message_length", h.messageLength, scalar.UintDescription("previous"))
					d.FieldValueUint("message_type_id", h.messageTypeID, scalar.UintDescription("previous"))
					d.FieldValueUint("message_stream_id", h.messageStreamID, scalar.UintDescription("previous"))
				case 3:
					h.timestamp = cs.prevHeader.timestamp
					h.timestampDelta = cs.prevHeader.timestampDelta
					h.messageLength = cs.prevHeader.messageLength
					h.messageStreamID = cs.prevHeader.messageStreamID
					h.messageTypeID = cs.prevHeader.messageTypeID
					d.FieldValueUint("message_length", h.messageLength, scalar.UintDescription("previous"))
					d.FieldValueUint("message_type_id", h.messageTypeID, scalar.UintDescription("previous"))
					d.FieldValueUint("message_stream_id", h.messageStreamID, scalar.UintDescription("previous"))
				}

				h.timestamp += h.timestampDelta

				d.FieldValueUint("calculated_timestamp", h.timestamp)

				m, ok := cs.messageSteams[h.messageStreamID]
				if !ok {
					m = &message{
						l:   h.messageLength,
						typ: h.messageTypeID,
					}
					cs.messageSteams[h.messageStreamID] = m
				}

				payloadLength := int64(chunkSize)
				messageLeft := int64(m.l) - int64(m.b.Len())
				if messageLeft < payloadLength {
					payloadLength = messageLeft
				}
				// support decoding interrupted rtmp stream
				// TODO: throw away message buffer? currently only do tcp so no point?
				payloadLength *= 8
				if payloadLength > d.BitsLeft() {
					payloadLength = d.BitsLeft()
				}

				if payloadLength > 0 {
					d.CopyBits(&m.b, d.FieldRawLen("data", payloadLength))
				}

				if m.l == uint64(m.b.Len()) {
					messageBR := bitio.NewBitReader(m.b.Bytes(), -1)
					messages.FieldStructRootBitBufFn("message", messageBR, func(d *decode.D) {
						d.FieldValueUint("message_stream_id", h.messageStreamID)
						d.FieldValueUint("message_type_id", m.typ, rtmpMessageTypeIDNames)
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
