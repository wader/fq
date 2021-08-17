package deepequal_test

import (
	"fmt"
	"testing"

	"github.com/wader/fq/internal/deepequal"
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
			expected := `
name diff:
--- expected
+++ actual
@@ -1 +1 @@
-aaaaaaaaa
+aaaaaabba

`
			actual := fmt.Sprintf(format, args...)
			if expected != actual {
				t.Errorf("expected %q, got %q", expected, actual)
			}
		}),
		"name",
		"aaaaaaaaa", "aaaaaabba",
	)
}
