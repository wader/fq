package deepequal_test

import (
	"fmt"
	"testing"

	"github.com/wader/bump/internal/deepequal"
)

type tfFn func(format string, args ...interface{})

func (fn tfFn) Errorf(format string, args ...interface{}) {
	fn(format, args...)
}

func (fn tfFn) Fatalf(format string, args ...interface{}) {
	fn(format, args...)
}

func TestError(t *testing.T) {
	deepequal.Error(
		tfFn(func(format string, args ...interface{}) {
			expected := `name
expected: "aaaaaaaaa"
  actual: "aaaaaabba"
    diff:        ^^  `
			actual := fmt.Sprintf(format, args...)
			if expected != actual {
				t.Errorf("expected %s, got %s", expected, actual)
			}
		}),
		"name",
		"aaaaaaaaa", "aaaaaabba",
	)
}
