package format_test

import (
	_ "fq/format/all"

	"fq/format"
	"fq/pkg/fqtest"
	"testing"
)

func TestFQTests(t *testing.T) {
	fqtest.TestPath(t, format.DefaultRegistry)
}
