package bitio_test

import (
	"bytes"
	"io"
	"log"
	"testing"

	"fq/pkg/bitio"
)

type shortBitReader struct {
	bitio.BitReaderAt
}

func (b shortBitReader) ReadBitsAt(p []byte, nBits int, bitOff int64) (n int, err error) {
	if nBits > 3 {
		nBits = 3
	}
	return b.BitReaderAt.ReadBitsAt(p, nBits, bitOff)
}

func TestReader(t *testing.T) {

	bb, bbBits := bitio.BytesFromBitString("011011110")

	br := bitio.NewReaderFromReadSeeker(bytes.NewReader(bb))
	sbr := bitio.NewSectionBitReader(br, 0, int64(bbBits))
	sbr2 := bitio.NewSectionBitReader(sbr, 0, 4)

	ob := make([]byte, 2)

	obBits, _ := sbr2.ReadBits(ob, 4)

	obs := bitio.BitStringFromBytes(ob, obBits)

	log.Printf("obs: %#+v\n", obs)
}

func TestCopy(t *testing.T) {

	br := bitio.NewReaderFromReadSeeker(bytes.NewReader([]byte{0xf0, 0xff, 0xff}))
	sbr := bitio.NewSectionBitReader(br, 0, 8*3-1)

	b := &bytes.Buffer{}

	n, err := io.Copy(b, sbr)
	log.Printf("n: %d err: %v b=%v\n", n, err, b.Bytes())

}

func TestReadFull(t *testing.T) {

	bb, bbBits := bitio.BytesFromBitString("011011110")

	br := bitio.NewReaderFromReadSeeker(bytes.NewReader(bb))

	ob := make([]byte, 2)

	obBits, _ := bitio.ReadAtFull(shortBitReader{br}, ob, bbBits, 0)

	obs := bitio.BitStringFromBytes(ob, obBits)

	log.Printf("obs: %#+v\n", obs)

}

func TestMultiBitReader(t *testing.T) {

	bb1, bb1Bits := bitio.BytesFromBitString("101")
	br1 := bitio.NewSectionBitReader(bitio.NewReaderFromReadSeeker(bytes.NewReader(bb1)), 0, int64(bb1Bits))
	bb2, bb2Bits := bitio.BytesFromBitString("0001")
	br2 := bitio.NewSectionBitReader(bitio.NewReaderFromReadSeeker(bytes.NewReader(bb2)), 0, int64(bb2Bits))

	mb, _ := bitio.NewMultiBitReader([]bitio.BitReadAtSeeker{br1, br2})

	ob := make([]byte, 2)

	obBits, _ := bitio.ReadAtFull(mb, ob, 7, 0)

	obs := bitio.BitStringFromBytes(ob, obBits)

	log.Printf("obs: %#+v\n", obs)

}
