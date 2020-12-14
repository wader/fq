package bitio_test

import (
	"bytes"
	"fmt"
	"fq/pkg/bitio"
	"testing"
)

func TestReverseBytes(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{nil, nil},
		{[]byte{1}, []byte{1}},
		{[]byte{1, 2}, []byte{2, 1}},
		{[]byte{1, 2, 3}, []byte{3, 2, 1}},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v", tC.input), func(t *testing.T) {
			actual := append([]byte(nil), tC.input...)
			bitio.ReverseBytes(actual)
			if !bytes.Equal(tC.expected, actual) {
				t.Errorf("expected %v, got %v", tC.expected, actual)
			}
		})
	}
}
