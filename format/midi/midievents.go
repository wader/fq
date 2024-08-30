package midi

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

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

	event := uint64(status & 0xf0)

	channel := func(d *decode.D) uint64 {
		b := d.PeekBytes(1)
		if b[0] >= 0x80 {
			d.U8()
		}

		return uint64(status & 0x0f)
	}

	midievent := func(name string, f func(d *decode.D)) {
		d.FieldStruct(name, func(d *decode.D) {
			d.FieldStruct("time", delta)
			d.FieldValueUint("event", event, midievents)
			d.FieldUintFn("channel", channel)

			f(d)
		})
	}

	switch event {
	case NoteOff:
		midievent("midievent", decodeNoteOff)

	case NoteOn:
		midievent("midievent", decodeNoteOn)

	case PolyphonicPressure:
		midievent("midievent", decodePolyphonicPressure)

	case Controller:
		midievent("midievent", decodeController)

	case ProgramChange:
		midievent("midievent", decodeProgramChange)

	case ChannelPressure:
		midievent("midievent", decodeChannelPressure)

	case PitchBend:
		midievent("midievent", decodePitchBend)

	default:
		flush(d, "unknown MIDI event (%02x)", status&0xf0)
	}
}

func decodeNoteOff(d *decode.D) {
	d.FieldStruct("note", func(d *decode.D) {
		d.FieldU8("note", notes)
		d.FieldUintFn("velocity", func(d *decode.D) uint64 {
			return d.U8() & 0x7f
		})
	})
}

func decodeNoteOn(d *decode.D) {
	d.FieldStruct("note", func(d *decode.D) {
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
		d.FieldU8("controller", controllers)
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
		bytes := d.BytesLen(2)

		bend := uint64(bytes[0])
		bend <<= 7
		bend |= uint64(bytes[1]) & 0x7f

		// ... centre value (0) is 0x2000 (81920)
		return int64(bend) - 8192
	})
}
