package hexdump_test

import (
	"bytes"
	"fmt"
	"fq/internal/hexdump"
	"testing"
)

func Test(t *testing.T) {
	testCases := []struct {
		writes   [][]byte
		expected string
	}{
		{[][]byte{}, "00|           |    |\n"},
		{[][]byte{{65}}, "00|41         |A   |\n"},
		{[][]byte{{0, 0xff}}, "00|00 ff      |..  |\n"},
		{[][]byte{{65, 66, 67}, {68}}, "00|41 42 43 44|ABCD|\n"},
		{[][]byte{{65, 66, 67}, {68}, {69}}, "00|41 42 43 44|ABCD|\n04|45         |E   |\n"},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.writes), func(t *testing.T) {
			b := &bytes.Buffer{}
			hd := hexdump.New(b, 2, 4)
			for _, w := range tC.writes {
				hd.Write(w)
			}
			hd.Close()
			actual := b.String()
			if tC.expected != actual {
				t.Errorf("expected %q, got %q", tC.expected, actual)
			}
		})
	}
}
