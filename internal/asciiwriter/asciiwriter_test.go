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
	_, _ = h.Write([]byte("\x00b"))
	_, _ = h.Write([]byte("c"))
	_, _ = h.Write([]byte("d"))
	_, _ = h.Write([]byte("e"))

	log.Printf("b.Bytes(): '%s'\n", b.Bytes())
}
