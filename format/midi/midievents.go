package midi

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// MIDI event status byte values. A MIDI event status byte is a composite byte
// composed of the event type in the high order nibble and the event channel
// (0 to 15) in the low order nibble.
const (
	NoteOff            uint64 = 0x80
	NoteOn             uint64 = 0x90
	PolyphonicPressure uint64 = 0xa0
	Controller         uint64 = 0xb0
	ProgramChange      uint64 = 0xc0
	ChannelPressure    uint64 = 0xd0
	PitchBend          uint64 = 0xe0
)

var midievents = scalar.UintMapSymStr{
	NoteOff:            "note_off",
	NoteOn:             "note_on",
	PolyphonicPressure: "polyphonic_pressure",
	Controller:         "controller",
	ProgramChange:      "program_change",
	ChannelPressure:    "channel_pressure",
	PitchBend:          "pitch_bend",
}

var midifns = map[uint64]func(d *decode.D){
	NoteOff:            decodeNoteOff,
	NoteOn:             decodeNoteOn,
	PolyphonicPressure: decodePolyphonicPressure,
	Controller:         decodeController,
	ProgramChange:      decodeProgramChange,
	ChannelPressure:    decodeChannelPressure,
	PitchBend:          decodePitchBend,
}

func decodeMIDIEvent(d *decode.D, status uint8, ctx *context) {
	if status < 0x80 {
		status = ctx.running
	}

	ctx.running = status
	ctx.casio = false

	delta := func(d *decode.D) {
		ctx.tick += d.FieldUintFn("delta", vlq)
		d.FieldValueUint("tick", ctx.tick)
	}

	if fn, ok := midifns[uint64(status&0x00f0)]; ok {
		d.FieldStruct("midi_event", func(d *decode.D) {
			d.FieldStruct("time", delta)

			b := d.PeekBytes(1)
			if b[0] >= 0x80 {
				d.FieldUintFn("event", func(d *decode.D) uint64 {
					return d.U4() << 4
				}, midievents)
				d.FieldU4("channel")
			} else {
				d.FieldValueUint("event", uint64(status&0x00f0), midievents)
				d.FieldValueUint("channel", uint64(status&0x000f))
			}

			fn(d)
		})
	} else {
		flush(d, "unknown MIDI event (%02x)", status&0xf0)
	}
}

func decodeNoteOff(d *decode.D) {
	d.FieldStruct("note_off", func(d *decode.D) {
		d.FieldU8("note", notes)
		d.FieldUintFn("velocity", func(d *decode.D) uint64 {
			return d.U8() & 0x7f
		})
	})
}

func decodeNoteOn(d *decode.D) {
	d.FieldStruct("note_on", func(d *decode.D) {
		d.FieldU8("note", notes)
		d.FieldUintFn("velocity", func(d *decode.D) uint64 {
			return d.U8() & 0x7f
		})
	})
}

func decodePolyphonicPressure(d *decode.D) {
	d.FieldU8("polyphonic_pressure")
}

func decodeController(d *decode.D) {
	d.FieldStruct("controller", func(d *decode.D) {
		d.FieldU8("controller", controllersMap)
		d.FieldU8("value")
	})
}

func decodeProgramChange(d *decode.D) {
	d.FieldU8("program_change")
}

func decodeChannelPressure(d *decode.D) {
	d.FieldU8("channel_pressure")
}

func decodePitchBend(d *decode.D) {
	d.FieldSintFn("pitch_bend", func(d *decode.D) int64 {
		// ... 14 bit range i.e. [0..16383]
		bytes := d.BytesLen(2)

		bend := uint64(bytes[0]) & 0x7f
		bend <<= 7
		bend |= uint64(bytes[1]) & 0x7f

		// ... centre value (0) is 0x2000 (81920)
		return int64(bend) - 8192
	})
}
