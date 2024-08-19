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
		N := int(d.FieldUintFn("length", vlq))
		data := d.PeekBytes(N)

		if len(data) > 0 {
			d.FieldUintFn("framerate", func(d *decode.D) uint64 {
				d.BytesLen(1)
				rr := (data[0] >> 6) & 0x03

				switch rr {
				case 0:
					return 24

				case 1:
					return 25

				case 2:
					return 29

				case 3:
					return 30

				default:
					return 0
				}
			})
			d.FieldValueUint("hour", uint64(data[0]&0x01f))
		}

		if len(data) > 1 {
			d.FieldUintFn("minute", func(d *decode.D) uint64 {
				d.BytesLen(1)
				return uint64(data[1])
			})
		}

		if len(data) > 2 {
			d.FieldUintFn("second", func(d *decode.D) uint64 {
				d.BytesLen(1)
				return uint64(data[2])
			})
		}

		if len(data) > 3 {
			d.FieldUintFn("frames", func(d *decode.D) uint64 {
				d.BytesLen(1)
				return uint64(data[3])
			})
		}

		if len(data) > 4 {
			d.FieldUintFn("fractions", func(d *decode.D) uint64 {
				d.BytesLen(1)
				return uint64(data[4])
			})
		}
	})
}

func decodeTimeSignature(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldStruct("signature", func(d *decode.D) {
		N := int(d.FieldUintFn("length", vlq))
		data := d.PeekBytes(N)

		d.FieldUintFn("numerator", func(d *decode.D) uint64 {
			d.BytesLen(1)
			return uint64(data[0])
		})

		d.FieldUintFn("denominator", func(d *decode.D) uint64 {
			d.BytesLen(1)
			denominator := uint16(1)
			for i := uint8(0); i < data[1]; i++ {
				denominator *= 2
			}

			return uint64(denominator)
		})

		d.FieldUintFn("ticksPerClick", func(d *decode.D) uint64 {
			d.BytesLen(1)
			return uint64(data[2])
		})

		d.FieldUintFn("thirtySecondsPerQuarter", func(d *decode.D) uint64 {
			d.BytesLen(1)
			return uint64(data[3])
		})
	})
}

func decodeKeySignature(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	d.FieldUintFn("key", func(d *decode.D) uint64 {
		data := vlf(d)
		key := uint64(data[0]) & 0x00ff
		key <<= 8
		key |= uint64(data[1]) & 0x00ff

		return key

	}, keys)
}

func decodeEndOfTrack(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")
	d.FieldUintFn("length", func(d *decode.D) uint64 {
		return uint64(len(vlf(d)))
	})
}

func decodeSequencerSpecificEvent(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldU8("status")
	d.FieldU8("event")

	d.FieldStruct("message", func(d *decode.D) {
		N := int(d.FieldUintFn("length", vlq))
		data := d.PeekBytes(N)

		if len(data) > 2 && data[0] == 0x00 {
			d.FieldStrFn("manufacturer", func(d *decode.D) string {
				d.BytesLen(3)
				return fmt.Sprintf("%02X%02X", data[1], data[2])
			}, manufacturers)

			if len(data) > 3 {
				d.FieldStrFn("data", func(d *decode.D) string {
					d.BytesLen(N - 3)
					return fmt.Sprintf("%v", data[3:])
				})
			}

		} else if len(data) > 0 {
			d.FieldStrFn("manufacturer", func(d *decode.D) string {
				d.BytesLen(1)
				return fmt.Sprintf("%02x", data[0])
			}, manufacturers)
			if len(data) > 1 {
				d.BytesLen(N - 1)
				d.FieldStrFn("data", func(d *decode.D) string {
					return fmt.Sprintf("%v", data[1:])
				})
			}
		}
	})
}
