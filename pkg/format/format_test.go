package format_test

import (
	_ "fq/pkg/format/all"

	"fq/pkg/format"
	"fq/pkg/test"
	"testing"
)

func Test(t *testing.T) {
	test.TestPath(t, format.DefaultRegistry)
}
