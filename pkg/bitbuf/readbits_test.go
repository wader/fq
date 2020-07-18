package bitbuf_test

import (
	"encoding/hex"
	"fmt"
	"fq/pkg/bitbuf"
	"testing"
)

func TestReadBits(t *testing.T) {
	testCases := []struct {
		buf      []byte
		bitPos   uint64
		bits     uint64
		expected uint64
	}{
		{buf: []byte{0xff}, bitPos: 0, bits: 8, expected: 0b11111111},
		{buf: []byte{0xff}, bitPos: 1, bits: 7, expected: 0b1111111},
		{buf: []byte{0xff}, bitPos: 2, bits: 6, expected: 0b111111},
		{buf: []byte{0xff}, bitPos: 3, bits: 5, expected: 0b11111},
		{buf: []byte{0xff}, bitPos: 4, bits: 4, expected: 0b1111},
		{buf: []byte{0xff}, bitPos: 5, bits: 3, expected: 0b111},
		{buf: []byte{0xff}, bitPos: 6, bits: 2, expected: 0b11},
		{buf: []byte{0xff}, bitPos: 7, bits: 1, expected: 0b1},
		{buf: []byte{0xff}, bitPos: 8, bits: 0, expected: 0},

		{buf: []byte{0xff, 0xff}, bitPos: 0, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 1, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 2, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 3, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 4, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 5, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 6, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 7, bits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, bitPos: 8, bits: 8, expected: 0xff},

		{buf: []byte{0x0f, 0x01}, bitPos: 6, bits: 10, expected: 0x301},

		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}, bitPos: 0, bits: 64, expected: 0x1122334455667788},
		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}, bitPos: 0, bits: 56, expected: 0x11223344556677},
		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}, bitPos: 0, bits: 48, expected: 0x112233445566},
		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55}, bitPos: 0, bits: 40, expected: 0x1122334455},
		{buf: []byte{0x11, 0x22, 0x33, 0x44}, bitPos: 0, bits: 32, expected: 0x11223344},
		{buf: []byte{0x11, 0x22, 0x33}, bitPos: 0, bits: 24, expected: 0x112233},
		{buf: []byte{0x11, 0x22}, bitPos: 0, bits: 16, expected: 0x1122},
		{buf: []byte{0x11}, bitPos: 0, bits: 8, expected: 0x11},

		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 12, expected: 0b111100001111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 1, bits: 12, expected: 0b111000011110},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 2, bits: 12, expected: 0b110000111100},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 3, bits: 12, expected: 0b100001111000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 4, bits: 12, expected: 0b000011110000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 5, bits: 12, expected: 0b000111100001},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 6, bits: 12, expected: 0b001111000011},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 7, bits: 12, expected: 0b011110000111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 8, bits: 12, expected: 0b111100001111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 9, bits: 12, expected: 0b111000011110},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 10, bits: 12, expected: 0b110000111100},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 11, bits: 12, expected: 0b100001111000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 12, bits: 12, expected: 0b000011110000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 13, bits: 12, expected: 0b000111100001},

		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 1, expected: 0b1},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 2, expected: 0b11},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 3, expected: 0b111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 4, expected: 0b1111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 5, expected: 0b11110},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 6, expected: 0b111100},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 7, expected: 0b1111000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 8, expected: 0b11110000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 9, expected: 0b111100001},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 10, expected: 0b1111000011},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 11, expected: 0b11110000111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 12, expected: 0b111100001111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, bitPos: 0, bits: 13, expected: 0b1111000011110},

		{buf: []byte{0xf8}, bitPos: 6, bits: 1, expected: 0},

		{buf: []byte{0x40}, bitPos: 1, bits: 6, expected: 0b100000},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%s %d %d", hex.EncodeToString(tC.buf), tC.bitPos, tC.bits), func(t *testing.T) {
			actual := bitbuf.ReadBits(tC.buf, tC.bitPos, tC.bits)
			if tC.expected != actual {
				t.Errorf("expected %x, got %x", tC.expected, actual)
			}
		})
	}
}

func TestReadBitsPanic(t *testing.T) {
	// TODO: check panic string
	defer func() { recover() }()
	bitbuf.ReadBits([]byte{}, 0, 65)
	t.Error("should panic")
}
