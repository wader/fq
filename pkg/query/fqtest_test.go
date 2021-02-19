package query_test

import (
	_ "fq/pkg/format/all"

	"fq/pkg/format"
	"fq/pkg/fqtest"
	"testing"
)

func TestFQTests(t *testing.T) {
	fqtest.TestPath(t, format.DefaultRegistry)
}
