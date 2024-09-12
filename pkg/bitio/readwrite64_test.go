package bitio_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/wader/fq/pkg/bitio"
)

func TestRead64(t *testing.T) {
	testCases := []struct {
		buf      []byte
		firstBit int64
		nBits    int64
		expected uint64
	}{
		{buf: []byte{0xff}, firstBit: 0, nBits: 8, expected: 0b11111111},
		{buf: []byte{0xff}, firstBit: 1, nBits: 7, expected: 0b1111111},
		{buf: []byte{0xff}, firstBit: 2, nBits: 6, expected: 0b111111},
		{buf: []byte{0xff}, firstBit: 3, nBits: 5, expected: 0b11111},
		{buf: []byte{0xff}, firstBit: 4, nBits: 4, expected: 0b1111},
		{buf: []byte{0xff}, firstBit: 5, nBits: 3, expected: 0b111},
		{buf: []byte{0xff}, firstBit: 6, nBits: 2, expected: 0b11},
		{buf: []byte{0xff}, firstBit: 7, nBits: 1, expected: 0b1},
		{buf: []byte{0xff}, firstBit: 8, nBits: 0, expected: 0},

		{buf: []byte{0xff, 0xff}, firstBit: 0, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 1, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 2, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 3, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 4, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 5, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 6, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 7, nBits: 8, expected: 0xff},
		{buf: []byte{0xff, 0xff}, firstBit: 8, nBits: 8, expected: 0xff},

		{buf: []byte{0x0f, 0x01}, firstBit: 6, nBits: 10, expected: 0x301},

		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}, firstBit: 0, nBits: 64, expected: 0x1122334455667788},
		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}, firstBit: 0, nBits: 56, expected: 0x11223344556677},
		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}, firstBit: 0, nBits: 48, expected: 0x112233445566},
		{buf: []byte{0x11, 0x22, 0x33, 0x44, 0x55}, firstBit: 0, nBits: 40, expected: 0x1122334455},
		{buf: []byte{0x11, 0x22, 0x33, 0x44}, firstBit: 0, nBits: 32, expected: 0x11223344},
		{buf: []byte{0x11, 0x22, 0x33}, firstBit: 0, nBits: 24, expected: 0x112233},
		{buf: []byte{0x11, 0x22}, firstBit: 0, nBits: 16, expected: 0x1122},
		{buf: []byte{0x11}, firstBit: 0, nBits: 8, expected: 0x11},

		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 12, expected: 0b111100001111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 1, nBits: 12, expected: 0b111000011110},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 2, nBits: 12, expected: 0b110000111100},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 3, nBits: 12, expected: 0b100001111000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 4, nBits: 12, expected: 0b000011110000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 5, nBits: 12, expected: 0b000111100001},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 6, nBits: 12, expected: 0b001111000011},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 7, nBits: 12, expected: 0b011110000111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 8, nBits: 12, expected: 0b111100001111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 9, nBits: 12, expected: 0b111000011110},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 10, nBits: 12, expected: 0b110000111100},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 11, nBits: 12, expected: 0b100001111000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 12, nBits: 12, expected: 0b000011110000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 13, nBits: 12, expected: 0b000111100001},

		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 1, expected: 0b1},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 2, expected: 0b11},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 3, expected: 0b111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 4, expected: 0b1111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 5, expected: 0b11110},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 6, expected: 0b111100},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 7, expected: 0b1111000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 8, expected: 0b11110000},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 9, expected: 0b111100001},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 10, expected: 0b1111000011},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 11, expected: 0b11110000111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 12, expected: 0b111100001111},
		{buf: []byte{0xf0, 0xf0, 0xf0, 0xf0}, firstBit: 0, nBits: 13, expected: 0b1111000011110},

		{buf: []byte{0xf8}, firstBit: 6, nBits: 1, expected: 0},

		{buf: []byte{0x40}, firstBit: 1, nBits: 6, expected: 0b100000},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%s %d %d", hex.EncodeToString(tC.buf), tC.firstBit, tC.nBits), func(t *testing.T) {
			actual := bitio.Read64(tC.buf, tC.firstBit, tC.nBits)
			if tC.expected != actual {
				t.Errorf("expected %x, got %x", tC.expected, actual)
			}
		})
	}
}

func TestWrite64(t *testing.T) {
	testCases := []struct {
		v           uint64
		nBits       int64
		buf         []byte
		firstBit    int64
		expectedBuf []byte
	}{
		{0x0123456789abcdef, 8, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0xef, 0x00, 0x00, 0x0, 0x0, 0x00, 0x00, 0x00}},
		{0x0123456789abcdef, 16, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0xcd, 0xef, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{0x0123456789abcdef, 24, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{0x0123456789abcdef, 32, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0x89, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00, 0x00}},
		{0x0123456789abcdef, 40, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0x67, 0x89, 0xab, 0xcd, 0xef, 0x00, 0x00, 0x00}},
		{0x0123456789abcdef, 48, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x00, 0x00}},
		{0x0123456789abcdef, 56, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x00}},
		{0x0123456789abcdef, 64, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0, []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}},

		{0b1111, 4, []byte{0b0000_1001, 0}, 0, []byte{0b1111_1001, 00}},

		{0b1111, 4, []byte{0b0000_1001, 0}, 1, []byte{0b0111_1001, 0b0000_0000}},

		{0b1111_0000_1111, 12, []byte{0b0000_0000, 0b0000_1001}, 0, []byte{0b1111_0000, 0b1111_1001}},

		{0b1111_0000_1111, 12, []byte{0b1000_0000, 0b0000_1001}, 1, []byte{0b1111_1000, 0b0111_1001}},

		{0xf, 4, []byte{0x0e, 0}, 0, []byte{0xfe, 00}},

		{0b0, 1, []byte{0b1111_1111}, 0, []byte{0b0111_1111}},
		{0b0, 1, []byte{0b1111_1111}, 1, []byte{0b1011_1111}},
		{0b0, 1, []byte{0b1111_1111}, 2, []byte{0b1101_1111}},
		{0b0, 1, []byte{0b1111_1111}, 3, []byte{0b1110_1111}},
		{0b0, 1, []byte{0b1111_1111}, 4, []byte{0b1111_0111}},
		{0b0, 1, []byte{0b1111_1111}, 5, []byte{0b1111_1011}},
		{0b0, 1, []byte{0b1111_1111}, 6, []byte{0b1111_1101}},
		{0b0, 1, []byte{0b1111_1111}, 7, []byte{0b1111_1110}},

		{0b1, 1, []byte{0b0000_0000}, 0, []byte{0b1000_0000}},
		{0b1, 1, []byte{0b0000_0000}, 1, []byte{0b0100_0000}},
		{0b1, 1, []byte{0b0000_0000}, 2, []byte{0b0010_0000}},
		{0b1, 1, []byte{0b0000_0000}, 3, []byte{0b0001_0000}},
		{0b1, 1, []byte{0b0000_0000}, 4, []byte{0b0000_1000}},
		{0b1, 1, []byte{0b0000_0000}, 5, []byte{0b0000_0100}},
		{0b1, 1, []byte{0b0000_0000}, 6, []byte{0b0000_0010}},
		{0b1, 1, []byte{0b0000_0000}, 7, []byte{0b0000_0001}},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%d %d %s", tC.v, tC.nBits, hex.EncodeToString(tC.buf)), func(t *testing.T) {
			bitio.Write64(tC.v, tC.nBits, tC.buf, tC.firstBit)
			if !bytes.Equal(tC.expectedBuf, tC.buf) {
				t.Errorf("expected %s, got %s", hex.EncodeToString(tC.expectedBuf), hex.EncodeToString(tC.buf))
			}
		})
	}
}
