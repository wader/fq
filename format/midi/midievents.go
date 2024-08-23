package midi

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

type MidiEventType byte

const (
	TypeNoteOff            MidiEventType = 0x80
	TypeNoteOn             MidiEventType = 0x90
	TypePolyphonicPressure MidiEventType = 0xa0
	TypeController         MidiEventType = 0xb0
	TypeProgramChange      MidiEventType = 0xc0
	TypeChannelPressure    MidiEventType = 0xd0
	TypePitchBend          MidiEventType = 0xe0
)

var midievents = scalar.UintMapSymStr{
	0x80: "NoteOff",
	0x90: "NoteOn",
	0xa0: "PolyphonicPressure",
	0xb0: "Controller",
	0xc0: "ProgramChange",
	0xd0: "ChannelPressure",
	0xe0: "PitchBend",
}

func decodeMIDIEvent(d *decode.D, status uint8, ctx *context) {
	if status < 0x80 {
		status = ctx.running
	}

	ctx.running = status
	ctx.casio = false

	delta := func(d *decode.D) {
		dt := d.FieldUintFn("delta", vlq)
		d.FieldValueUint("tick", ctx.tick)

		ctx.tick += dt
	}

	channel := func(d *decode.D) uint64 {
		b := d.PeekBytes(1)
		if b[0] >= 0x80 {
			d.BytesLen(1)
		}

		return uint64(status & 0x0f)
	}

	event := uint64(status & 0xf0)

	switch MidiEventType(event) {
	case TypeNoteOff:
		d.FieldStruct("NoteOff", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			decodeNoteOff(d)
		})

	case TypeNoteOn:
		d.FieldStruct("NoteOn", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			decodeNoteOn(d)
		})

	case TypePolyphonicPressure:
		d.FieldStruct("PolyphonicPressure", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			decodePolyphonicPressure(d)
		})

	case TypeController:
		d.FieldStruct("Controller", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			decodeController(d)
		})

	case TypeProgramChange:
		d.FieldStruct("ProgramChange", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			decodeProgramChange(d)
		})

	case TypeChannelPressure:
		d.FieldStruct("ChannelPressure", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			decodeChannelPressure(d)
		})

	case TypePitchBend:
		d.FieldStruct("PitchBend", func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			decodePitchBend(d)
		})

	default:
		flush(d, "unknown MIDI event (%02x)", status&0xf0)
	}
}

func decodeNoteOff(d *decode.D) {
	d.FieldU8("note", notes)
	d.FieldU8("velocity")
}

func decodeNoteOn(d *decode.D) {
	d.FieldU8("note", notes)
	d.FieldU8("velocity")
}

func decodePolyphonicPressure(d *decode.D) {
	d.FieldU8("pressure")
}

func decodeController(d *decode.D) {
	d.FieldU8("controller", controllers)
	d.FieldU8("value")
}

func decodeProgramChange(d *decode.D) {
	d.FieldU8("program")
}

func decodeChannelPressure(d *decode.D) {
	d.FieldU8("pressure")
}

func decodePitchBend(d *decode.D) {
	d.FieldUintFn("bend", func(d *decode.D) uint64 {
		data := d.BytesLen(2)

		bend := uint64(data[0])
		bend <<= 7
		bend |= uint64(data[1]) & 0x7f

		return bend
	})
}
