package bitio

import (
	"fmt"
	"strings"
)

// BytesFromBitString []byte from bit string, ex: "0101" -> ([]byte{0x50}, 4)
func BytesFromBitString(s string) ([]byte, int) {
	r := len(s) % 8
	bufLen := len(s) / 8
	if r > 0 {
		bufLen++
	}
	buf := make([]byte, bufLen)

	for i := 0; i < len(s); i++ {
		d := s[i] - '0'
		if d != 0 && d != 1 {
			panic(fmt.Sprintf("invalid bit string %q at index %d %q", s, i, s[i]))
		}
		buf[i/8] |= d << (7 - i%8)
	}

	return buf, len(s)
}

// BitStringFromBytes string from []byte], ex: ([]byte{0x50}, 4) -> "0101"
func BitStringFromBytes(buf []byte, nBits int) string {
	sb := &strings.Builder{}
	for i := 0; i < nBits; i++ {
		if buf[i/8]&(1<<(7-i%8)) > 0 {
			sb.WriteString("1")
		} else {
			sb.WriteString("0")
		}
	}
	return sb.String()
}
