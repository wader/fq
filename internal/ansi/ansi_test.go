package ansi_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/wader/fq/internal/ansi"
)

func TestF(t *testing.T) {
	c := ansi.FromString("blue")
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "%s", c.F("test"))
	actual := b.String()
	expected := "\x1b[34mtest\x1b[39m"
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestW(t *testing.T) {
	c := ansi.FromString("blue+underline")
	b := &bytes.Buffer{}
	_, _ = c.W(b).Write([]byte("test"))
	actual := b.String()
	expected := "\x1b[34;4mtest\x1b[39;24m"
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
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
