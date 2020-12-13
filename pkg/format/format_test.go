package format_test

import (
	_ "fq/pkg/format/all"

	"fq/pkg/format"
	"fq/pkg/fqtest"
	"testing"
)

func Test(t *testing.T) {
	fqtest.TestPath(t, format.DefaultRegistry)
}
