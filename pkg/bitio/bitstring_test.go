package bitio_test

import (
	"fq/pkg/bitio"
	"testing"
)

func TestBitString(t *testing.T) {
	testCases := []string{
		"",
		"1",
		"0",
		"10",
		"01",
		"11",
		"00",
		"1000001",
		"0000000",
		"10000001",
		"00000000",
		"100000001",
		"000000000",
		"101010101",
		"111100000",
	}
	for _, tC := range testCases {
		t.Run(tC, func(t *testing.T) {
			bb, bbBits := bitio.BytesFromBitString(tC)
			actual := bitio.BitStringFromBytes(bb, bbBits)
			if tC != actual {
				t.Errorf("expected %s, got %s", tC, actual)
			}
		})
	}
}
