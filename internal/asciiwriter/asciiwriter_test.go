package asciiwriter_test

import (
	"bytes"
	"fq/internal/asciiwriter"
	"log"
	"testing"
)

func TestWrite(t *testing.T) {
	b := &bytes.Buffer{}
	h := asciiwriter.New(b, 4, 0, asciiwriter.SafeASCII)
	h.Write([]byte("\x00b"))
	h.Write([]byte("c"))
	h.Write([]byte("d"))
	h.Write([]byte("e"))

	log.Printf("b.Bytes(): '%s'\n", b.Bytes())
}
