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

var controllers = scalar.UintMapSymStr{
	// High resolution continuous controllers (MSB)
	0:  "Bank Select (MSB)",
	1:  "Modulation Wheel (MSB)",
	2:  "Breath Controller (MSB)",
	4:  "Foot Controller (MSB)",
	5:  "Portamento Time (MSB)",
	6:  "Data Entry (MSB)",
	7:  "Channel Volume (MSB)",
	8:  "Balance (MSB)",
	10: "Pan (MSB)",
	11: "Expression Controller (MSB)",
	12: "Effect Control 1 (MSB)",
	13: "Effect Control 2 (MSB)",
	16: "General Purpose Controller 1 (MSB)",
	17: "General Purpose Controller 2 (MSB)",
	18: "General Purpose Controller 3 (MSB)",
	19: "General Purpose Controller 4 (MSB)",

	// High resolution continuous controllers (LSB)
	32: "Bank Select (LSB)",
	33: "Modulation Wheel (LSB)",
	34: "Breath Controller (LSB)",
	36: "Foot Controller (LSB)",
	37: "Portamento Time (LSB)",
	38: "Data Entry (LSB)",
	39: "Channel Volume (LSB)",
	40: "Balance (LSB)",
	42: "Pan (LSB)",
	43: "Expression Controller (LSB)",
	44: "Effect Control 1 (LSB)",
	45: "Effect Control 2 (LSB)",
	48: "General Purpose Controller 1 (LSB)",
	49: "General Purpose Controller 2 (LSB)",
	50: "General Purpose Controller 3 (LSB)",
	51: "General Purpose Controller 4 (LSB)",

	// Switches
	64: "Sustain On/Off",
	65: "Portamento On/Off",
	66: "Sostenuto On/Off",
	67: "Soft Pedal On/Off",
	68: "Legato On/Off",
	69: "Hold 2 On/Off",

	// Low resolution continuous controllers
	70: "Sound Controller 1  (TG: Sound Variation;  FX: Exciter On/Off)",
	71: "Sound Controller 2  (TG: Harmonic Content; FX: Compressor On/Off)",
	72: "Sound Controller 3  (TG: Release Time;     FX: Distortion On/Off)",
	73: "Sound Controller 4  (TG: Attack Time;      FX: EQ On/Off)",
	74: "Sound Controller 5  (TG: Brightness;       FX: Expander On/Off)",
	75: "Sound Controller 6  (TG: Decay Time;       FX: Reverb On/Off)",
	76: "Sound Controller 7  (TG: Vibrato Rate;     FX: Delay On/Off)",
	77: "Sound Controller 8  (TG: Vibrato Depth;    FX: Pitch Transpose On/Off)",
	78: "Sound Controller 9  (TG: Vibrato Delay;    FX: Flange/Chorus On/Off)",
	79: "Sound Controller 10 (TG: Undefined;        FX: Special Effects On/Off)",
	80: "General Purpose Controller 5",
	81: "General Purpose Controller 6",
	82: "General Purpose Controller 7",
	83: "General Purpose Controller 8",
	84: "Portamento Control",
	88: "High Resolution Velocity Prefix",
	91: "Effects 1 Depth (Reverb Send Level)",
	92: "Effects 2 Depth (Tremelo Depth)",
	93: "Effects 3 Depth (Chorus Send Level)",
	94: "Effects 4 Depth (Celeste Depth)",
	95: "Effects 5 Depth (Phaser Depth)",

	// RPNs / NRPNs
	96:  "Data Increment",
	97:  "Data Decrement",
	98:  "Non-Registered Parameter Number (LSB)",
	99:  "Non-Registered Parameter Number (MSB)",
	100: "Registered Parameter Number (LSB)",
	101: "Registered Parameter Number (MSB)",

	// Channel Mode messages
	120: "All Sound Off",
	121: "Reset All Controllers",
	122: "Local Control On/Off",
	123: "All Notes Off",
	124: "Omni Mode Off",
	125: "Omni Mode On ",
	126: "Mono Mode On",
	127: "Poly Mode On",
}

func decodeMIDIEvent(d *decode.D, status uint8) {
	event := status & 0xf0

	channel := func(d *decode.D) uint64 {
		b := d.PeekBytes(1)
		if b[0] >= 0x80 {
			d.BytesLen(1)
		}

		return uint64(status & 0x0f)
	}

	switch MidiEventType(event) {
	case TypeNoteOff:
		d.FieldStruct("NoteOff", func(d *decode.D) {
			d.FieldUintFn("delta", vlq)
			d.FieldUintFn("channel", channel)

			decodeNoteOff(d)
		})
		return

	case TypeNoteOn:
		d.FieldStruct("NoteOn", func(d *decode.D) {
			d.FieldUintFn("delta", vlq)
			d.FieldUintFn("channel", channel)

			decodeNoteOn(d)
		})
		return

	case TypePolyphonicPressure:
		d.FieldStruct("PolyphonicPressure", func(d *decode.D) {
			d.FieldUintFn("delta", vlq)
			d.FieldUintFn("channel", channel)

			decodePolyphonicPressure(d)
		})
		return

	case TypeController:
		d.FieldStruct("Controller", func(d *decode.D) {
			d.FieldUintFn("delta", vlq)
			d.FieldUintFn("channel", channel)

			decodeController(d)
		})
		return

	case TypeProgramChange:
		d.FieldStruct("ProgramChange", func(d *decode.D) {
			d.FieldUintFn("delta", vlq)
			d.FieldUintFn("channel", channel)

			decodeProgramChange(d)
		})
		return

	case TypeChannelPressure:
		d.FieldStruct("ChannelPressure", func(d *decode.D) {
			d.FieldUintFn("delta", vlq)
			d.FieldUintFn("channel", channel)

			decodeChannelPressure(d)
		})
		return

	case TypePitchBend:
		d.FieldStruct("PitchBend", func(d *decode.D) {
			d.FieldUintFn("delta", vlq)
			d.FieldUintFn("channel", channel)

			decodePitchBend(d)
		})
		return
	}

	// ... unknown event - flush remaining data
	d.Errorf("unknown MIDI event (%02x)", event)

	var N int = int(d.BitsLeft())

	d.Bits(N)
}

func decodeNoteOff(d *decode.D) {
	d.FieldU8("note")
	d.FieldU8("velocity")
}

func decodeNoteOn(d *decode.D) {
	d.FieldU8("note")
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
