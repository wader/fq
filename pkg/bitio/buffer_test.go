package bitio_test

import (
	"fmt"
	"fq/pkg/bitio"
	"log"
	"math/rand"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	bb := bitio.NewBufferFromBytes([]byte{0xf5}, -1)
	e4, e4n := bb.Bits(3)

	log.Printf("e4: %#+v\n", e4)
	log.Printf("e4n: %#+v\n", e4n)

	// sbb, sbbn := bb.BitBufLen(5)
	// log.Printf("sbbn: %#+v\n", sbbn)
	// log.Printf("sbb.BitString(): %#+v\n", sbb.BitString())
}

func TestBufferBitString(t *testing.T) {
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
			bb := bitio.NewBufferFromBitString(tC)
			actual := bb.BitString()
			if tC != actual {
				t.Errorf("expected %s, got %s", tC, actual)
			}

			for i := int64(0); i < bb.Len; i++ {
				t.Run(fmt.Sprintf("%s_%d", tC, i), func(t *testing.T) {
					startBb, _ := bb.BitBufRange(i, bb.Len-i)
					startExpected := tC[i : i+bb.Len-i]
					startActual := startBb.BitString()
					if startExpected != startActual {
						t.Errorf("startBb expected %s, got %s", startExpected, startActual)
					}

					lenBb, _ := bb.BitBufRange(0, i)
					lenExpected := tC[0:i]
					lenActual := lenBb.BitString()
					if lenExpected != lenActual {
						t.Errorf("lenBb expected %s, got %s", lenExpected, lenActual)
					}
				})
			}
		})
	}
}

func TestBitStringRandom(t *testing.T) {
	r := rand.New(rand.NewSource(0))

	for i := 0; i < 10000; i++ {
		var ss []string
		for j := uint32(0); j < r.Uint32()%1000; j++ {
			switch r.Uint32() % 2 {
			case 0:
				ss = append(ss, "0")
			case 1:
				ss = append(ss, "1")
			}
		}
		expected := strings.Join(ss, "")
		bb := bitio.NewBufferFromBitString(expected)
		actual := bb.BitString()
		if expected != actual {
			t.Errorf("expected %s, got %s", expected, actual)
		}
	}
}

func TestInvalidBitString(t *testing.T) {
	// TODO: check panic?
	bitio.NewBufferFromBitString("01invalid")
}
