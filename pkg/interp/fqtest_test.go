package interp_test

import (
	"fq/format/all"
	"fq/pkg/fqtest"
	"testing"
)

func TestFQTests(t *testing.T) {
	fqtest.TestPath(t, all.Registry)
}
