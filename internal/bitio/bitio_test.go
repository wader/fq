package bitio_test

import (
	"bytes"
	"io"
	"log"
	"testing"

	"fq/internal/bitio"
)

func TestReader(t *testing.T) {

	br, _ := bitio.NewFromReadSeeker(bytes.NewReader([]byte{0x37}))
	sbr := bitio.NewSectionBitReader(br, 4, 4)
	sbr2 := bitio.NewSectionBitReader(sbr, 0, 4)

	b := make([]byte, 1)
	var n int
	var err error

	n, err = sbr.ReadBits(b, 2)
	log.Printf("n: %d err: %v b=%v\n", n, err, b)
	n, err = sbr.ReadBits(b, 2)
	log.Printf("n: %d err: %v b=%v\n", n, err, b)
	n, err = sbr.ReadBits(b, 2)
	log.Printf("n: %d err: %v b=%v\n", n, err, b)

	log.Println("----")

	n, err = sbr2.ReadBits(b, 2)
	log.Printf("n: %d err: %v b=%v\n", n, err, b)
	n, err = sbr2.ReadBits(b, 2)
	log.Printf("n: %d err: %v b=%v\n", n, err, b)
	n, err = sbr2.ReadBits(b, 2)
	log.Printf("n: %d err: %v b=%v\n", n, err, b)
}

func TestCopy(t *testing.T) {

	br, _ := bitio.NewFromReadSeeker(bytes.NewReader([]byte{0xff, 0xff}))
	sbr := bitio.NewSectionBitReader(br, 1, 8)

	b := &bytes.Buffer{}

	n, err := io.Copy(b, sbr)
	log.Printf("n: %d err: %v b=%v\n", n, err, b.Bytes())

}
