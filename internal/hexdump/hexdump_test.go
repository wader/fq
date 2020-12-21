package hexdump_test

import (
	"bytes"
	"fmt"
	"fq/internal/hexdump"
	"testing"
)

func Test(t *testing.T) {
	testCases := []struct {
		startOffset int
		writes      [][]byte
		expected    string
	}{
		{0, [][]byte{},
			"00|           |    |\n"},
		{0, [][]byte{{65}},
			"00|41         |A   |\n"},
		{0, [][]byte{{0, 0xff}},
			"00|00 ff      |..  |\n"},
		{0, [][]byte{{65, 66, 67}, {68}},
			"00|41 42 43 44|ABCD|\n"},
		{0, [][]byte{{65, 66, 67}, {68}, {69}}, "" +
			"00|41 42 43 44|ABCD|\n" +
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
			hd := hexdump.New(b, int64(tC.startOffset), 2, 4)
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
