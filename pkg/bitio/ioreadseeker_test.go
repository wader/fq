package bitio_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/wader/fq/pkg/bitio"
)

func sb(s string) *bitio.SectionReader {
	buf, nBits := bitio.BytesFromBitString(s)
	return bitio.NewBitReader(buf, nBits)
}

func bs(br bitio.Reader) string {
	bib := &bitio.Buffer{}
	_, err := bitio.Copy(bib, br)
	if err != nil {
		panic(err)
	}
	buf, nBits := bib.Bits()
	return bitio.BitStringFromBytes(buf, nBits)
}

func Test(t *testing.T) {
	testCases := []struct {
		bs string
	}{
		{""},
		{"|"},
		{"0"},
		{"1"},
		{"1|"},
		{"|1|"},
		{"|1|"},
		{"0|"},
		{"|0|"},
		{"0|"},
		{"1"},
		{"10"},
		{"101"},
		{"1011"},
		{"10110"},
		{"101101"},
		{"1011011"},
		{"10110110"},

		{"101101101"},
		{"1011011010"},
		{"10110110101"},
		{"101101101011"},
		{"1011011010110"},
		{"10110110101101"},
		{"101101101011011"},
		{"1011011010110110"},

		{"10110110|1"},
		{"10110110|10"},
		{"10110110|101"},
		{"10110110|1011"},
		{"10110110|10110"},
		{"10110110|101101"},
		{"10110110|1011011"},
		{"10110110|10110110"},

		{"1|10110110|1"},
		{"10|10110110|10"},
		{"101|10110110|101"},
		{"1011|10110110|1011"},
		{"10110|10110110|10110"},
		{"101101|10110110|101101"},
		{"1011011|10110110|1011011"},
		{"10110110|10110110|10110110"},

		{"1|1|0110110|1"},
		{"10|10|110110|10"},
		{"101|101|10110|101"},
		{"1011|1011|0110|1011"},
		{"10110|10110|110|10110"},
		{"101101|101101|10|101101"},
		{"1011011|1011011|0|1011011"},

		{"1|10110110|10110110|1"},
		{"10|10110110|10110110|10"},
		{"101|10110110|10110110|101"},
		{"1011|10110110|10110110|1011"},
		{"10110|10110110|10110110|10110"},
		{"101101|10110110|10110110|101101"},
		{"1011011|10110110|10110110|1011011"},
		{"10110110|10110110|10110110|10110110"},

		{"1|101101101011011010110110101101101011011010110110|1"},
		{"10|101101101011011010110110101101101011011010110110|10"},
		{"101|101101101011011010110110101101101011011010110110|101"},
		{"1011|101101101011011010110110101101101011011010110110|1011"},
		{"10110|101101101011011010110110101101101011011010110110|10110"},
		{"101101|101101101011011010110110101101101011011010110110|101101"},
		{"1011011|101101101011011010110110101101101011011010110110|1011011"},
		{"10110110|101101101011011010110110101101101011011010110110|10110110"},
	}
	for _, tC := range testCases {

		bsParts := strings.Split(tC.bs, "|")
		var bsBRs []bitio.ReadAtSeeker
		for _, p := range bsParts {
			bsBRs = append(bsBRs, sb(p))
		}
		bsBR, err := bitio.NewMultiBitReader(bsBRs...)
		if err != nil {
			panic(err)
		}

		bsBitString := strings.ReplaceAll(tC.bs, "|", "")

		for i := 0; i < len(bsBitString); i++ {
			t.Run(fmt.Sprintf("%s_%d", tC.bs, i), func(t *testing.T) {
				_, err = bsBR.SeekBits(int64(i), io.SeekStart)
				if err != nil {
					t.Fatal(err)
				}

				expectedBitString := bsBitString[i:]
				actualBitString := bs(bsBR)

				if expectedBitString != actualBitString {
					t.Errorf("expected bits %q, got %q", expectedBitString, actualBitString)
				}

				_, err = bsBR.SeekBits(int64(i), io.SeekStart)
				if err != nil {
					t.Fatal(err)
				}

				r := bitio.NewIOReader(bsBR)
				bb := &bytes.Buffer{}
				if _, err := io.Copy(bb, r); err != nil {
					t.Fatal(err)
				}

				expecetdByteBitString := expectedBitString + strings.Repeat("0", (8-(len(expectedBitString)%8))%8)
				actualByteBitString := bitio.BitStringFromBytes(bb.Bytes(), int64(bb.Len()*8))

				if expecetdByteBitString != actualByteBitString {
					t.Errorf("expected bytes %q, got %q", expecetdByteBitString, actualByteBitString)
				}
			})
		}
	}
}
