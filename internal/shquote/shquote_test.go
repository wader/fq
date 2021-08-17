package shquote_test

import (
	"reflect"
	"testing"

	"github.com/wader/fq/internal/shquote"
)

func TestSplit(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{``, nil},
		{`abbc`, []string{`abbc`}},
		{` abbc `, []string{`abbc`}},
		{`a bb c`, []string{`a`, `bb`, `c`}},
		{`"b b"`, []string{`b b`}},
		{`"b ' b"`, []string{`b ' b`}},
		{`"b \"b"`, []string{`b "b`}},
		{`'b b'`, []string{`b b`}},
		{`'b " b'`, []string{`b " b`}},
		{`'b \"b'`, []string{`b \"b`}},
		{`a'b'"c"`, []string{`abc`}},
		{`a"b"c`, []string{`abc`}},
		{`a'b'c`, []string{`abc`}},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			actual := shquote.Split(tC.input)
			if !reflect.DeepEqual(tC.expected, actual) {
				t.Errorf("expected %#v, got %#v", tC.expected, actual)
			}
		})
	}
}
