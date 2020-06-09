package mp3

import (
	"fq/internal/decode"
)

// MP3 decoder
type MP3 struct {
	decode.Common
}

// Decode MP3
func (m *MP3) Decode() {
	m.FieldUFn("sync", func() (uint64, decode.Format, string) {
		n := m.U11()
		s := "correct"
		if n != 0b11111111111 {
			s = "incorrect"
		}
		return n, decode.FormatHex, s
	})
	m.FieldUFn("block_size", func() (uint64, decode.Format, string) {
		switch m.U2() {
		case 0b00:
			return 0b00, decode.FormatDecimal, "MPEG Version 2.5"
		case 0b01:
			return 0b01, decode.FormatDecimal, "reserved"
		case 0b10:
			return 0b10, decode.FormatDecimal, "MPEG Version 2"
		case 0b11:
			return 0b11, decode.FormatDecimal, "MPEG Version 1"
		default:
			panic("unreachable")
		}
	})
}
