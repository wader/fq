package format_test

import (
	_ "fq/format/all"
	"fq/format/registry"
	"fq/pkg/fqtest"
	"testing"
)

func TestFQTests(t *testing.T) {
	fqtest.TestPath(t, registry.Default)
}
