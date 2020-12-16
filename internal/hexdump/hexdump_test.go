package hexdump_test

import (
	"bytes"
	"fq/internal/hexdump"
	"testing"
)

func TestHexdumper(t *testing.T) {
	b := &bytes.Buffer{}
	hd := hexdump.New(b, 2, 4)
	_, _ = hd.Write([]byte{65, 66, 67})
	_, _ = hd.Write([]byte{68})
	hd.Close()
	expected := "01|41 42 43 44|ABCD|\n"
	actual := b.String()
	if expected != actual {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
