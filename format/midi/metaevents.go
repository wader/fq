package midi

import (
	"fmt"

	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const (
	TypeSequenceNumber         uint64 = 0x00
	TypeText                   uint64 = 0x01
	TypeCopyright              uint64 = 0x02
	TypeTrackName              uint64 = 0x03
	TypeInstrumentName         uint64 = 0x04
	TypeLyric                  uint64 = 0x05
	TypeMarker                 uint64 = 0x06
	TypeCuePoint               uint64 = 0x07
	TypeProgramName            uint64 = 0x08
	TypeDeviceName             uint64 = 0x09
	TypeMIDIChannelPrefix      uint64 = 0x20
	TypeMIDIPort               uint64 = 0x21
	TypeTempo                  uint64 = 0x51
	TypeSMPTEOffset            uint64 = 0x54
	TypeTimeSignature          uint64 = 0x58
	TypeKeySignature           uint64 = 0x59
	TypeEndOfTrack             uint64 = 0x2f
	TypeSequencerSpecificEvent uint64 = 0x7f
)

var metaevents = scalar.UintMapSymStr{
	TypeSequenceNumber:         "sequence_number",
	TypeText:                   "text",
	TypeCopyright:              "copyright",
	TypeTrackName:              "track_name",
	TypeInstrumentName:         "instrument_name",
	TypeLyric:                  "lyric",
	TypeMarker:                 "marker",
	TypeCuePoint:               "cue_point",
	TypeProgramName:            "program_name",
	TypeDeviceName:             "device_name",
	TypeMIDIChannelPrefix:      "midi_channel_prefix",
	TypeMIDIPort:               "midi_port",
	TypeTempo:                  "tempo",
	TypeSMPTEOffset:            "smpte_offset",
	TypeTimeSignature:          "time_signature",
	TypeKeySignature:           "key_signature",
	TypeEndOfTrack:             "end_of_track",
	TypeSequencerSpecificEvent: "sequencer_specific_event",
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
			d.FieldU8("status")
			d.FieldU8("event", metaevents)
			f(d)
		})
	}

	switch uint64(event) {
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
			d.Fatalf("%v", err)
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
			d.Fatalf("%v", err)
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
			d.Fatalf("%v", err)
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
			d.Fatalf("%v", err)
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
				d.Fatalf("%v", err)
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
		length := d.FieldUintFn("length", vlq)

		if length > 0 {
			d.FieldU8("numerator")
		}

		if length > 1 {
			d.FieldUintFn("denominator", func(d *decode.D) uint64 {
				denominator := uint64(1)
				v := d.U8()
				for i := uint64(0); i < v; i++ {
					denominator <<= 1
				}
				return denominator
			})
		}

		if length > 2 {
			d.FieldU8("ticks_per_click")
		}

		if length > 3 {
			d.FieldU8("thirty_seconds_per_quarter")
		}
	})
}

func decodeKeySignature(d *decode.D) {
	d.FieldUintFn("key", func(d *decode.D) uint64 {
		key := uint64(0)

		if data, err := vlf(d); err != nil {
			d.Fatalf("%v", err)
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
			d.Fatalf("%v", err)
		} else {
			length = len(data)
		}

		return uint64(length)
	})
}

func decodeSequencerSpecificEvent(d *decode.D) {
	d.FieldStruct("info", func(d *decode.D) {
		if length := d.FieldUintFn("length", vlq); length > 0 {
			b := d.PeekUintBits(8)

			if length > 2 && b == 0 {
				d.FieldStrFn("manufacturer", func(d *decode.D) string {
					manufacturer := d.BytesLen(3)

					return fmt.Sprintf("%02X%02X", manufacturer[1], manufacturer[2])
				}, manufacturers)

				if length > 3 {
					d.FieldRawLen("data", 8*(int64(length)-3))
				}
			} else if length > 0 {
				d.FieldStrFn("manufacturer", func(d *decode.D) string {
					return fmt.Sprintf("%02X", d.U8())
				}, manufacturers)

				if length > 1 {
					d.FieldRawLen("data", 8*(int64(length)-1))
				}
			}
		}
	})
}
