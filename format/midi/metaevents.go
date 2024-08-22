package midi

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
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

var metaevents = scalar.UintMapSymUint{
	0xff00: 0x00,
	0xff01: 0x01,
	0xff02: 0x02,
	0xff03: 0x03,
	0xff04: 0x04,
	0xff05: 0x05,
	0xff06: 0x06,
	0xff07: 0x07,
	0xff08: 0x08,
	0xff09: 0x09,
	0xff20: 0x20,
	0xff21: 0x21,
	0xff51: 0x51,
	0xff54: 0x54,
	0xff58: 0x58,
	0xff59: 0x59,
	0xff2f: 0x2f,
	0xff7f: 0x7f,
}

var framerates = scalar.UintMapSymUint{
	0: 24,
	1: 25,
	2: 29,
	3: 30,
}

func decodeMetaEvent(d *decode.D, event uint8, ctx *context) {
	ctx.running = 0x00
	ctx.casio = false

	delta := func(d *decode.D) {
		dt := d.FieldUintFn("delta", vlq)
		d.FieldValueUint("tick", ctx.tick)

		ctx.tick += dt
	}

	switch MetaEventType(event) {
	case TypeSequenceNumber:
		d.FieldStruct("SequenceNumber", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeSequenceNumber(d)
		})

	case TypeText:
		d.FieldStruct("Text", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeText(d)
		})

	case TypeCopyright:
		d.FieldStruct("Copyright", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeCopyright(d)
		})

	case TypeTrackName:
		d.FieldStruct("TrackName", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeTrackName(d)
		})

	case TypeInstrumentName:
		d.FieldStruct("InstrumentName", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeInstrumentName(d)
		})

	case TypeLyric:
		d.FieldStruct("Lyric", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeLyric(d)
		})

	case TypeMarker:
		d.FieldStruct("Marker", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeMarker(d)
		})

	case TypeCuePoint:
		d.FieldStruct("CuePoint", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeCuePoint(d)
		})

	case TypeProgramName:
		d.FieldStruct("ProgramName", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeProgramName(d)
		})

	case TypeDeviceName:
		d.FieldStruct("DeviceName", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeDeviceName(d)
		})

	case TypeMIDIChannelPrefix:
		d.FieldStruct("TypeMIDIChannelPrefix", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeMIDIChannelPrefix(d)
		})

	case TypeMIDIPort:
		d.FieldStruct("TypeMIDIPort", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeMIDIPort(d)
		})

	case TypeTempo:
		d.FieldStruct("Tempo", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeTempo(d)
		})

	case TypeSMPTEOffset:
		d.FieldStruct("SMPTEOffset", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeSMPTEOffset(d)
		})

	case TypeTimeSignature:
		d.FieldStruct("TimeSignature", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeTimeSignature(d)
		})

	case TypeKeySignature:
		d.FieldStruct("KeySignature", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeKeySignature(d)
		})

	case TypeEndOfTrack:
		d.FieldStruct("EndOfTrack", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeEndOfTrack(d)
		})

	case TypeSequencerSpecificEvent:
		d.FieldStruct("SequencerSpecific", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("metaevent", metaevents)
			decodeSequencerSpecificEvent(d)
		})

	default:
		flush(d, "unknown meta event (%02x)", event)
	}
}

func decodeSequenceNumber(d *decode.D) {
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
	d.FieldStrFn("text", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeCopyright(d *decode.D) {
	d.FieldStrFn("copyright", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeTrackName(d *decode.D) {
	d.FieldStrFn("name", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeInstrumentName(d *decode.D) {
	d.FieldStrFn("instrument", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeLyric(d *decode.D) {
	d.FieldStrFn("lyric", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeMarker(d *decode.D) {
	d.FieldStrFn("marker", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeCuePoint(d *decode.D) {
	d.FieldStrFn("cue", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeProgramName(d *decode.D) {
	d.FieldStrFn("program", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeDeviceName(d *decode.D) {
	d.FieldStrFn("device", func(d *decode.D) string {
		return string(vlf(d))
	})
}

func decodeMIDIChannelPrefix(d *decode.D) {
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
	d.FieldStruct("offset", func(d *decode.D) {
		var data []uint8
		d.FieldStrFn("bytes", func(d *decode.D) string {
			data = vlf(d)

			return fmt.Sprintf("%v", data)
		})

		if len(data) > 0 {
			d.FieldUintFn("framerate", func(d *decode.D) uint64 {
				return uint64((data[0] >> 6) & 0x03)
			}, framerates)

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
	d.FieldStruct("signature", func(d *decode.D) {
		var data []uint8
		d.FieldStrFn("bytes", func(d *decode.D) string {
			data = vlf(d)

			return fmt.Sprintf("%v", data)
		})

		if len(data) > 0 {
			d.FieldValueUint("numerator", uint64(data[0]))
		}

		if len(data) > 1 {
			denominator := uint64(1)
			for i := uint8(0); i < data[1]; i++ {
				denominator <<= 1
			}

			d.FieldValueUint("denominator", denominator)
		}

		if len(data) > 2 {
			d.FieldValueUint("ticksPerClick", uint64(data[2]))
		}

		if len(data) > 3 {
			d.FieldValueUint("thirtySecondsPerQuarter", uint64(data[3]))
		}
	})
}

func decodeKeySignature(d *decode.D) {
	d.FieldUintFn("key", func(d *decode.D) uint64 {
		data := vlf(d)
		key := uint64(data[0]) & 0x00ff
		key <<= 8
		key |= uint64(data[1]) & 0x00ff

		return key

	}, keys)
}

func decodeEndOfTrack(d *decode.D) {
	d.FieldUintFn("length", func(d *decode.D) uint64 {
		return uint64(len(vlf(d)))
	})
}

func decodeSequencerSpecificEvent(d *decode.D) {
	d.FieldStruct("info", func(d *decode.D) {
		var data []uint8
		d.FieldStrFn("bytes", func(d *decode.D) string {
			data = vlf(d)

			return fmt.Sprintf("%v", data)
		})

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
	})
}
