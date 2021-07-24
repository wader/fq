// TODO: fix tests
//+build ignore

package hexdump_test

import (
	"bytes"
	"fmt"
	"fq/internal/asciiwriter"
	"fq/internal/hexdump"
	"fq/internal/hexpairwriter"
	"testing"
)

func TestHexdump(t *testing.T) {
	testCases := []struct {
		startOffset int
		writes      [][]byte
		expected    string
	}{
		{0, [][]byte{}, `
   |00 01 02 03|    |
0x0|           |    |
`[1:]},
		{0, [][]byte{{65}}, `
   |00 01 02 03|    |
0x0|41         |A   |
`[1:]},
		{0, [][]byte{{0, 0xff}},
			"0x0|00 ff      |..  |\n"},
		{0, [][]byte{{65, 66, 67}, {68}},
			"0x0|41 42 43 44|ABCD|\n"},
		{0, [][]byte{{65, 66, 67}, {68}, {69}}, "" +
			"0x0|41 42 43 44|ABCD|\n" +
			"04|45         |E   |\n"},

		{4, [][]byte{{65}},
			"04|41         |A   |\n"},
		{5, [][]byte{{65}},
			"04|   41      | A  |\n"},
		{6, [][]byte{{65}},
			"04|      41   |  A |\n"},
		{7, [][]byte{{65}},
			"04|         41|   A|\n"},
		{8, [][]byte{{65}},
			"08|41         |A   |\n"},

		{3, [][]byte{{65, 66, 67}, {68}, {69}}, "" +
			"00|         41|   A|\n" +
			"04|42 43 44 45|BCDE|\n"},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.writes), func(t *testing.T) {
			b := &bytes.Buffer{}
			hd := hexdump.New(b, int64(tC.startOffset), 2, 16, 4,
				func(b byte) string { return hexpairwriter.Pair(b) },
				func(b byte) string { return asciiwriter.SafeASCII(b) },
				func(s string) string { return s },
				func(s string) string { return s },
				"|",
			)
			for _, w := range tC.writes {
				if _, err := hd.Write(w); err != nil {
					t.Fatal(err)
				}
			}
			hd.Close()
			actual := b.String()
			if tC.expected != actual {
				t.Errorf("expected %q, got %q", tC.expected, actual)
			}
		})
	}
}
