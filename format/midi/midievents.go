package midi

import (
	"github.com/wader/fq/pkg/decode"
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

func decodeProgramChange(d *decode.D) {
	d.FieldUintFn("delta", vlq)
	d.FieldUintFn("channel", func(d *decode.D) uint64 {
		b := d.BytesLen(1)

		return uint64(b[0] & 0x0f)
	})

	d.FieldU8("program")
}
