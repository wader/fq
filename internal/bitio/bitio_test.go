package bitio_test

import (
	"bytes"
	"io"
	"log"
	"testing"

	"fq/internal/bitio"
)

func TestReader(t *testing.T) {

	bb, bbBits := bitio.BytesFromBitString("011011110")

	br, _ := bitio.NewFromReadSeeker(bytes.NewReader(bb))
	sbr := bitio.NewSectionBitReader(br, 0, int64(bbBits))
	sbr2 := bitio.NewSectionBitReader(sbr, 0, 4)

	ob := make([]byte, 2)

	obBits, _ := sbr2.ReadBits(ob, 4)

	obs := bitio.BitStringFromBytes(ob, obBits)

	log.Printf("obs: %#+v\n", obs)
}

func TestCopy(t *testing.T) {

	br, _ := bitio.NewFromReadSeeker(bytes.NewReader([]byte{0xff, 0xff}))
	sbr := bitio.NewSectionBitReader(br, 1, 8)

	b := &bytes.Buffer{}

	n, err := io.Copy(b, sbr)
	log.Printf("n: %d err: %v b=%v\n", n, err, b.Bytes())

}
