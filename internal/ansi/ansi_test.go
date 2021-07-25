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
	c := ansi.FromString("blue+underline")

	b := &bytes.Buffer{}

	_, _ = c.W(b).Write([]byte("bla"))

	log.Printf("b.String(): %#+v\n", b.String())
}

func TestLen(t *testing.T) {
	testCases := []struct {
		s string
		l int
	}{
		{"", 0},
		{"abc", 3},
		{ansi.Red.SetString + "a" + "bc" + ansi.Red.ResetString + "d", 4},
		{"a" + ansi.Red.SetString + "bc" + ansi.Red.ResetString + "d", 4},
		{"a" + ansi.Red.SetString + "bcd" + ansi.Red.ResetString, 4},
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
