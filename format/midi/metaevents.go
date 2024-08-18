package midi

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
)

type MetaEventType uint8

const (
	TypeSequenceNumber         MetaEventType = 0x00
	TypeText                   MetaEventType = 0x01
	TypeCopyright              MetaEventType = 0x02
	TypeTrackName              MetaEventType = 0x03
	TypeInstrumentName         MetaEventType = 0x04
	TypeLyric                  MetaEventType = 0x05
	TypeMarker                 MetaEventType = 0x06
	TypeCuePoint               MetaEventType = 0x07
	TypeProgramName            MetaEventType = 0x08
	TypeDeviceName             MetaEventType = 0x09
	TypeMIDIChannelPrefix      MetaEventType = 0x20
	TypeMIDIPort               MetaEventType = 0x21
	TypeTempo                  MetaEventType = 0x51
	TypeSMPTEOffset            MetaEventType = 0x54
	TypeTimeSignature          MetaEventType = 0x58
	TypeKeySignature           MetaEventType = 0x59
	TypeEndOfTrack             MetaEventType = 0x2f
	TypeSequencerSpecificEvent MetaEventType = 0x7f
)

func decodeMetaEvent(d *decode.D, event uint8, ctx *context) {
	ctx.running = 0x00
	ctx.casio = false

	switch MetaEventType(event) {
	case TypeSequenceNumber:
		d.FieldStruct("SequenceNumber", decodeSequenceNumber)

	case TypeText:
		d.FieldStruct("Text", decodeText)

	case TypeCopyright:
		d.FieldStruct("Copyright", decodeCopyright)

	case TypeTrackName:
		d.FieldStruct("TrackName", decodeTrackName)

	case TypeInstrumentName:
		d.FieldStruct("InstrumentName", decodeInstrumentName)

	case TypeLyric:
		d.FieldStruct("Lyric", decodeLyric)

	case TypeMarker:
		d.FieldStruct("Marker", decodeMarker)

	case TypeCuePoint:
		d.FieldStruct("CuePoint", decodeCuePoint)

	case TypeProgramName:
		d.FieldStruct("ProgramName", decodeProgramName)

	case TypeDeviceName:
		d.FieldStruct("DeviceName", decodeDeviceName)

	case TypeMIDIChannelPrefix:
		d.FieldStruct("TypeMIDIChannelPrefix", decodeMIDIChannelPrefix)

	case TypeMIDIPort:
		d.FieldStruct("TypeMIDIPort", decodeMIDIPort)

	case TypeTempo:
		d.FieldStruct("Tempo", decodeTempo)

	case TypeSMPTEOffset:
		d.FieldStruct("SMPTEOffset", decodeSMPTEOffset)

	case TypeTimeSignature:
		d.FieldStruct("TimeSignature", decodeTimeSignature)

	case TypeKeySignature:
		d.FieldStruct("KeySignature", decodeKeySignature)

	case TypeEndOfTrack:
		d.FieldStruct("EndOfTrack", decodeEndOfTrack)

	case TypeSequencerSpecificEvent:
		d.FieldStruct("SequencerSpecific", decodeSequencerSpecificEvent)

	default:
		flush(d, "unknown meta event (%02x)", event)
	}
}

func decodeSequenceNumber(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldUintFn("sequenceNumber", func(d *decode.D) uint64 {
		data := vlf(d)
		seqno := uint64(0)

		if len(data) > 0 {
			seqno += uint64(data[0])
		}

		if len(data) > 1 {
			seqno <<= 8
			seqno += uint64(data[1])
		}

		return seqno
	})
}

func decodeText(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("text", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeCopyright(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("copyright", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeTrackName(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("name", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeInstrumentName(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("instrument", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeLyric(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("lyric", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeMarker(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("marker", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeCuePoint(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("cue", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeProgramName(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("program", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeDeviceName(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStrFn("device", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeMIDIChannelPrefix(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldUintFn("channel", func(d *decode.D) uint64 {
		channel := uint64(0)
		data := vlf(d)

		for _, b := range data {
			channel <<= 8
			channel |= uint64(b & 0x00ff)
		}

		return channel
	})
}

func decodeMIDIPort(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldUintFn("port", func(d *decode.D) uint64 {
		channel := uint64(0)
		data := vlf(d)

		for _, b := range data {
			channel <<= 8
			channel |= uint64(b & 0x00ff)
		}

		return channel
	})
}

func decodeTempo(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	d.FieldUintFn("tempo", func(d *decode.D) uint64 {
		tempo := uint64(0)
		data := vlf(d)

		for _, b := range data {
			tempo <<= 8
			tempo |= uint64(b & 0x00ff)
		}

		return tempo
	})
}

func decodeSMPTEOffset(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	d.FieldStruct("offset", func(d *decode.D) {
		data := vlf(d)

		if len(data) > 0 {
			rr := (data[0] >> 6) & 0x03

			switch rr {
			case 0:
				d.FieldValueUint("framerate", 24)

			case 1:
				d.FieldValueUint("framerate", 25)

			case 2:
				d.FieldValueUint("framerate", 29)

			case 3:
				d.FieldValueUint("framerate", 30)

			default:
				d.FieldValueUint("framerate", 0)
			}

			d.FieldValueUint("hour", uint64(data[0]&0x01f))
		}

		if len(data) > 1 {
			d.FieldValueUint("minute", uint64(data[1]))
		}

		if len(data) > 2 {
			d.FieldValueUint("second", uint64(data[2]))
		}

		if len(data) > 3 {
			d.FieldValueUint("frames", uint64(data[3]))
		}

		if len(data) > 4 {
			d.FieldValueUint("fractions", uint64(data[4]))
		}
	})
}

func decodeTimeSignature(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStruct("signature", func(d *decode.D) {
		bytes := vlf(d)

		if len(bytes) > 0 {
			d.FieldValueUint("numerator", uint64(bytes[0]))
		}

		if len(bytes) > 1 {
			denominator := uint16(1)
			for i := uint8(0); i < bytes[1]; i++ {
				denominator *= 2
			}

			d.FieldValueUint("denominator", uint64(denominator))
		}

		if len(bytes) > 2 {
			d.FieldValueUint("ticksPerClick", uint64(bytes[2]))
		}

		if len(bytes) > 3 {
			d.FieldValueUint("thirtySecondsPerQuarter", uint64(bytes[3]))
		}
	})
}

func decodeKeySignature(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	data := vlf(d)
	if len(data) > 1 {
		key := (uint64(data[0]) << 8) & 0xff00
		key |= (uint64(data[1]) << 0) & 0x00ff

		d.FieldValueUint("key", key, keys)
	}
}

func decodeEndOfTrack(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	vlf(d)
}

func decodeSequencerSpecificEvent(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	data := vlf(d)
	if len(data) > 2 && data[0] == 0x00 {
		d.FieldValueStr("manufacturer", fmt.Sprintf("%02X%02X", data[1], data[2]), manufacturers)

		if len(data) > 3 {
			d.FieldValueStr("data", fmt.Sprintf("%v", data[3:]))
		}

	} else if len(data) > 0 {
		d.FieldValueStr("manufacturer", fmt.Sprintf("%02x", data[0]), manufacturers)
		if len(data) > 1 {
			d.FieldValueStr("data", fmt.Sprintf("%v", data[1:]))
		}
	}
}
