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

var metaevents = scalar.UintMapSymStr{
	0xff00: "sequence_number",
	0xff01: "text",
	0xff02: "copyright",
	0xff03: "track_name",
	0xff04: "instrument_name",
	0xff05: "lyric",
	0xff06: "marker",
	0xff07: "cue_point",
	0xff08: "program_name",
	0xff09: "device_name",
	0xff20: "midi_channel_prefix",
	0xff21: "midi_port",
	0xff51: "tempo",
	0xff54: "smpte_offset",
	0xff58: "time_signature",
	0xff59: "key_signature",
	0xff2f: "end_of_track",
	0xff7f: "sequencer_specific_event",
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

	metaevent := func(name string, f func(d *decode.D)) {
		d.FieldStruct(name, func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU16("event", metaevents)
			f(d)
		})
	}

	switch MetaEventType(event) {
	case TypeSequenceNumber:
		metaevent("sequence_number", decodeSequenceNumber)

	case TypeText:
		metaevent("text", decodeText)

	case TypeCopyright:
		metaevent("copyright", decodeCopyright)

	case TypeTrackName:
		metaevent("track_name", decodeTrackName)

	case TypeInstrumentName:
		metaevent("instrument_name", decodeInstrumentName)

	case TypeLyric:
		metaevent("lyric", decodeLyric)

	case TypeMarker:
		metaevent("marker", decodeMarker)

	case TypeCuePoint:
		metaevent("cue_point", decodeCuePoint)

	case TypeProgramName:
		metaevent("program_name", decodeProgramName)

	case TypeDeviceName:
		metaevent("device_name", decodeDeviceName)

	case TypeMIDIChannelPrefix:
		metaevent("midi_channel_prefix", decodeMIDIChannelPrefix)

	case TypeMIDIPort:
		metaevent("midi_port", decodeMIDIPort)

	case TypeTempo:
		metaevent("tempo", decodeTempo)

	case TypeSMPTEOffset:
		metaevent("smpte_offset", decodeSMPTEOffset)

	case TypeTimeSignature:
		metaevent("time_signature", decodeTimeSignature)

	case TypeKeySignature:
		metaevent("key_signature", decodeKeySignature)

	case TypeEndOfTrack:
		metaevent("end_of_track", decodeEndOfTrack)

	case TypeSequencerSpecificEvent:
		metaevent("sequencer_specific_event", decodeSequencerSpecificEvent)

	default:
		flush(d, "unknown meta event (%02x)", event)
	}
}

func decodeSequenceNumber(d *decode.D) {
	d.FieldUintFn("sequence_number", func(d *decode.D) uint64 {
		seqno := uint64(0)

		if data, err := vlf(d); err != nil {
			d.Errorf("%v", err)
		} else {
			if len(data) > 0 {
				seqno += uint64(data[0])
			}

			if len(data) > 1 {
				seqno <<= 8
				seqno += uint64(data[1])
			}
		}

		return seqno
	})
}

func decodeText(d *decode.D) {
	d.FieldStrFn("text", vlstring)
}

func decodeCopyright(d *decode.D) {
	d.FieldStrFn("copyright", vlstring)
}

func decodeTrackName(d *decode.D) {
	d.FieldStrFn("name", vlstring)
}

func decodeInstrumentName(d *decode.D) {
	d.FieldStrFn("instrument", vlstring)
}

func decodeLyric(d *decode.D) {
	d.FieldStrFn("lyric", vlstring)
}

func decodeMarker(d *decode.D) {
	d.FieldStrFn("marker", vlstring)
}

func decodeCuePoint(d *decode.D) {
	d.FieldStrFn("cue", vlstring)
}

func decodeProgramName(d *decode.D) {
	d.FieldStrFn("program", vlstring)
}

func decodeDeviceName(d *decode.D) {
	d.FieldStrFn("device", vlstring)
}

func decodeMIDIChannelPrefix(d *decode.D) {
	d.FieldUintFn("channel", func(d *decode.D) uint64 {
		channel := uint64(0)

		if data, err := vlf(d); err != nil {
			d.Errorf("%v", err)
		} else {
			for _, b := range data {
				channel <<= 8
				channel |= uint64(b & 0x00ff)
			}
		}

		return channel
	})
}

func decodeMIDIPort(d *decode.D) {
	d.FieldUintFn("port", func(d *decode.D) uint64 {
		port := uint64(0)

		if data, err := vlf(d); err != nil {
			d.Errorf("%v", err)
		} else {
			for _, b := range data {
				port <<= 8
				port |= uint64(b & 0x00ff)
			}
		}

		return port
	})
}

func decodeTempo(d *decode.D) {
	d.FieldUintFn("tempo", func(d *decode.D) uint64 {
		tempo := uint64(0)

		if data, err := vlf(d); err != nil {
			d.Errorf("%v", err)
		} else {
			for _, b := range data {
				tempo <<= 8
				tempo |= uint64(b & 0x00ff)
			}
		}

		return tempo
	})
}

func decodeSMPTEOffset(d *decode.D) {
	d.FieldStruct("offset", func(d *decode.D) {
		var data []uint8
		var err error

		d.FieldStrFn("bytes", func(d *decode.D) string {
			if data, err = vlf(d); err != nil {
				d.Errorf("%v", err)
			} else {
				return fmt.Sprintf("%v", data)
			}

			return "[]"
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
		var err error

		d.FieldStrFn("bytes", func(d *decode.D) string {
			if data, err = vlf(d); err != nil {
				d.Errorf("%v", err)
			} else {
				return fmt.Sprintf("%v", data)
			}

			return "[]"
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
			d.FieldValueUint("ticks_per_click", uint64(data[2]))
		}

		if len(data) > 3 {
			d.FieldValueUint("thirty_seconds_per_quarter", uint64(data[3]))
		}
	})
}

func decodeKeySignature(d *decode.D) {
	d.FieldUintFn("key", func(d *decode.D) uint64 {
		key := uint64(0)

		if data, err := vlf(d); err != nil {
			d.Errorf("%v", err)
		} else {
			if len(data) > 0 {
				key <<= 8
				key |= uint64(data[0]) & 0x00ff
			}

			if len(data) > 1 {
				key <<= 8
				key |= uint64(data[1]) & 0x00ff
			}
		}

		return key

	}, keys)
}

func decodeEndOfTrack(d *decode.D) {
	d.FieldUintFn("length", func(d *decode.D) uint64 {
		length := 0

		if data, err := vlf(d); err != nil {
			d.Errorf("%v", err)
		} else {
			length = len(data)
		}

		return uint64(length)
	})
}

func decodeSequencerSpecificEvent(d *decode.D) {
	d.FieldStruct("info", func(d *decode.D) {
		var data []uint8
		var err error

		d.FieldStrFn("bytes", func(d *decode.D) string {
			if data, err = vlf(d); err != nil {
				d.Errorf("%v", err)
			} else {
				return fmt.Sprintf("%v", data)
			}

			return "[]"
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
