package midi

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// MIDI meta-event status byte values.
const (
	SequenceNumber         uint64 = 0x00
	Text                   uint64 = 0x01
	Copyright              uint64 = 0x02
	TrackName              uint64 = 0x03
	InstrumentName         uint64 = 0x04
	Lyric                  uint64 = 0x05
	Marker                 uint64 = 0x06
	CuePoint               uint64 = 0x07
	ProgramName            uint64 = 0x08
	DeviceName             uint64 = 0x09
	MIDIChannelPrefix      uint64 = 0x20
	MIDIPort               uint64 = 0x21
	Tempo                  uint64 = 0x51
	SMPTEOffset            uint64 = 0x54
	TimeSignature          uint64 = 0x58
	KeySignature           uint64 = 0x59
	EndOfTrack             uint64 = 0x2f
	SequencerSpecificEvent uint64 = 0x7f
)

// Maps MIDI meta-events to a human readable name.
var metaevents = scalar.UintMapSymStr{
	SequenceNumber:         "sequence_number",
	Text:                   "text",
	Copyright:              "copyright",
	TrackName:              "track_name",
	InstrumentName:         "instrument_name",
	Lyric:                  "lyric",
	Marker:                 "marker",
	CuePoint:               "cue_point",
	ProgramName:            "program_name",
	DeviceName:             "device_name",
	MIDIChannelPrefix:      "midi_channel_prefix",
	MIDIPort:               "midi_port",
	Tempo:                  "tempo",
	SMPTEOffset:            "smpte_offset",
	TimeSignature:          "time_signature",
	KeySignature:           "key_signature",
	EndOfTrack:             "end_of_track",
	SequencerSpecificEvent: "sequencer_specific_event",
}

// Internal map of MIDI meta-events to the associated event parser.
var metafns = map[uint64]func(d *decode.D){
	SequenceNumber:         decodeSequenceNumber,
	Text:                   decodeText,
	Copyright:              decodeCopyright,
	TrackName:              decodeTrackName,
	InstrumentName:         decodeInstrumentName,
	Lyric:                  decodeLyric,
	Marker:                 decodeMarker,
	CuePoint:               decodeCuePoint,
	ProgramName:            decodeProgramName,
	DeviceName:             decodeDeviceName,
	MIDIChannelPrefix:      decodeMIDIChannelPrefix,
	MIDIPort:               decodeMIDIPort,
	Tempo:                  decodeTempo,
	SMPTEOffset:            decodeSMPTEOffset,
	TimeSignature:          decodeTimeSignature,
	KeySignature:           decodeKeySignature,
	EndOfTrack:             decodeEndOfTrack,
	SequencerSpecificEvent: decodeSequencerSpecificEvent,
}

// decodeMetaEvent extracts the meta-event delta time, event status and event detail.
func decodeMetaEvent(d *decode.D, event uint8, ctx *context) {
	ctx.running = 0x00
	ctx.casio = false

	delta := func(d *decode.D) {
		ctx.tick += d.FieldUintFn("delta", vlq)
		d.FieldValueUint("tick", ctx.tick)
	}

	if fn, ok := metafns[uint64(event)]; ok {
		d.FieldStruct("meta_event", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldU8("status")
			d.FieldU8("event", metaevents)
			fn(d)
		})
	} else {
		flush(d, "unknown meta event (%02x)", event)
	}
}

// decodeSequenceNumber parses a Sequence Number MIDI meta event to a struct comprising:
//   - sequence_number
func decodeSequenceNumber(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(2))
	d.FieldU16("sequence_number")
}

func decodeText(d *decode.D) {
	d.FieldStrFn("text", vlstring)
}

func decodeCopyright(d *decode.D) {
	d.FieldStrFn("copyright", vlstring)
}

func decodeTrackName(d *decode.D) {
	d.FieldStrFn("track_name", vlstring)
}

func decodeInstrumentName(d *decode.D) {
	d.FieldStrFn("instrument_name", vlstring)
}

func decodeLyric(d *decode.D) {
	d.FieldStrFn("lyric", vlstring)
}

func decodeMarker(d *decode.D) {
	d.FieldStrFn("marker", vlstring)
}

func decodeCuePoint(d *decode.D) {
	d.FieldStrFn("cue_point", vlstring)
}

func decodeProgramName(d *decode.D) {
	d.FieldStrFn("program_name", vlstring)
}

func decodeDeviceName(d *decode.D) {
	d.FieldStrFn("device_name", vlstring)
}

func decodeMIDIChannelPrefix(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(1))
	d.FieldU8("midi_channel_prefix")
}

func decodeMIDIPort(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(1))
	d.FieldU8("midi_port")
}

func decodeTempo(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(3))
	d.FieldU24("tempo")
}

func decodeSMPTEOffset(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(5))

	d.FieldStruct("smpte_offset", func(d *decode.D) {
		d.FieldU3("framerate", frameratesMap)
		d.FieldU5("hour")
		d.FieldU8("minute")
		d.FieldU8("second")
		d.FieldU8("frames")
		d.FieldU8("fractions")
	})
}

func decodeTimeSignature(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(4))

	d.FieldStruct("time_signature", func(d *decode.D) {
		d.FieldU8("numerator")

		d.FieldUintFn("denominator", func(d *decode.D) uint64 {
			denominator := uint64(1)
			v := d.U8()
			for i := uint64(0); i < v; i++ {
				denominator <<= 1
			}
			return denominator
		})

		d.FieldU8("ticks_per_click")
		d.FieldU8("thirty_seconds_per_quarter")
	})
}

func decodeKeySignature(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(2))
	d.FieldU16("key_signature", keys)
}

func decodeEndOfTrack(d *decode.D) {
	d.FieldUintFn("length", vlq, d.UintRequire(0))
}

func decodeSequencerSpecificEvent(d *decode.D) {
	length := d.FieldUintFn("length", vlq)

	d.FieldStruct("sequencer_specific_event", func(d *decode.D) {
		if length > 0 {
			b := d.PeekUintBits(8)

			if length > 2 && b == 0 {
				d.FieldU24("manufacturer", manufacturersExtendedMap)

				if length > 3 {
					d.FieldRawLen("data", 8*(int64(length)-3))
				}
			} else if length > 0 {
				d.FieldU8("manufacturer", manufacturersMap)

				if length > 1 {
					d.FieldRawLen("data", 8*(int64(length)-1))
				}
			}
		}
	})
}
