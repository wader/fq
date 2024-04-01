package mathx_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/wader/fq/internal/mathx"
)

func TestBigIntSetBytesSigned(t *testing.T) {
	testCases := []struct {
		buf []byte
		s   string
	}{
		{[]byte{1}, "1"},
		{[]byte{0b1111_1111}, "-1"},
		{[]byte{0b1000_0000}, "-128"},
		{[]byte{0b0111_1111}, "127"},
		{[]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "9223372036854775807"},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "-1"},
		{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "-9223372036854775808"},
		{[]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "2361183241434822606847"},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, "-1"},
		{[]byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, "-2361183241434822606848"},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v %s", tC.buf, tC.s), func(t *testing.T) {
			var n big.Int
			expected := tC.s
			actual := mathx.BigIntSetBytesSigned(&n, tC.buf).String()
			if expected != actual {
				t.Errorf("expected %s, got %s", expected, actual)
			}
		})
	}
}
