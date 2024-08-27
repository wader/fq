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
	0x80: "note_off",
	0x90: "note_on",
	0xa0: "polyphonic_pressure",
	0xb0: "controller",
	0xc0: "program_change",
	0xd0: "channel_pressure",
	0xe0: "pitch_bend",
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

	switch MidiEventType(event) {
	case TypeNoteOff:
		midievent("note_off", decodeNoteOff)

	case TypeNoteOn:
		midievent("note_on", decodeNoteOn)

	case TypePolyphonicPressure:
		midievent("polyphonic_pressure", decodePolyphonicPressure)

	case TypeController:
		midievent("controller", decodeController)

	case TypeProgramChange:
		midievent("program_change", decodeProgramChange)

	case TypeChannelPressure:
		midievent("channel_pressure", decodeChannelPressure)

	case TypePitchBend:
		midievent("pitch_bend", decodePitchBend)

	default:
		flush(d, "unknown MIDI event (%02x)", status&0xf0)
	}
}

func decodeNoteOff(d *decode.D) {
	d.AssertLeastBytesLeft(2)

	d.FieldU8("note", notes)
	d.FieldU8("velocity")
}

func decodeNoteOn(d *decode.D) {
	d.AssertLeastBytesLeft(2)

	d.FieldU8("note", notes)
	d.FieldU8("velocity")
}

func decodePolyphonicPressure(d *decode.D) {
	d.AssertLeastBytesLeft(1)

	d.FieldU8("pressure")
}

func decodeController(d *decode.D) {
	d.AssertLeastBytesLeft(2)

	d.FieldU8("controller", controllers)
	d.FieldU8("value")
}

func decodeProgramChange(d *decode.D) {
	d.AssertLeastBytesLeft(1)

	d.FieldU8("program")
}

func decodeChannelPressure(d *decode.D) {
	d.AssertLeastBytesLeft(1)

	d.FieldU8("pressure")
}

func decodePitchBend(d *decode.D) {
	d.AssertLeastBytesLeft(2)

	d.FieldUintFn("bend", func(d *decode.D) uint64 {
		data := d.BytesLen(2)

		bend := uint64(data[0])
		bend <<= 7
		bend |= uint64(data[1]) & 0x7f

		return bend
	})
}
