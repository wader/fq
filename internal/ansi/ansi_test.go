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

func TestLen(t *testing.T) {
	testCases := []struct {
		s string
		l int
	}{
		{"", 0},
		{"abc", 3},
		{ansi.FgRed + "a" + "bc" + ansi.Reset + "d", 4},
		{"a" + ansi.FgRed + "bc" + ansi.Reset + "d", 4},
		{"a" + ansi.FgRed + "bcd" + ansi.Reset, 4},
		{"a│b", 3},
	}
	for _, tC := range testCases {
		t.Run(tC.s, func(t *testing.T) {
			actualL := ansi.Len(tC.s)
			if tC.l != actualL {
				t.Errorf("expected %d, got %d", tC.l, actualL)
			}
		})
	}
}
