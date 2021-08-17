package interp_test

import (
	"testing"

	_ "github.com/wader/fq/format/all"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/fqtest"
)

func TestFQTests(t *testing.T) {
	fqtest.TestPath(t, registry.Default)
}
