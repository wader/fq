package hexpairwriter_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/wader/fq/internal/hexpairwriter"
)

func TestWrite(t *testing.T) {
	b := &bytes.Buffer{}
	h := hexpairwriter.New(b, 4, 0, hexpairwriter.Pair)
	_, _ = h.Write([]byte(""))
	_, _ = h.Write([]byte("ab"))
	_, _ = h.Write([]byte("c"))
	_, _ = h.Write([]byte("d"))

	log.Printf("b.Bytes(): '%s'\n", b.Bytes())
}
