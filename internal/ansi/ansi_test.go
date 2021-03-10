package ansi_test

import (
	"bytes"
	"fmt"
	"fq/internal/ansi"
	"log"
	"testing"
)

func Test(t *testing.T) {

	c := ansi.FromString("blue")

	b := &bytes.Buffer{}

	fmt.Fprintf(b, "%s", c.F("bla"))

	log.Printf("b.String(): %#+v\n", b.String())

}

func Test2(t *testing.T) {
	c := ansi.FromString("blue")

	b := &bytes.Buffer{}

	c.W(b).Write([]byte("bla"))

	log.Printf("b.String(): %#+v\n", b.String())
}
