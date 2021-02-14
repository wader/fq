package hexpairwriter_test

import (
	"bytes"
	"fq/internal/hexpairwriter"
	"log"
	"testing"
)

func TestWrite(t *testing.T) {
	b := &bytes.Buffer{}
	h := hexpairwriter.New(b, 4, 0, hexpairwriter.Pair)
	h.Write([]byte(""))
	h.Write([]byte("ab"))
	h.Write([]byte("c"))
	h.Write([]byte("d"))

	log.Printf("b.Bytes(): '%s'\n", b.Bytes())
}
