package interp_test

import (
	"testing"

	_ "github.com/wader/fq/format/all"
	"github.com/wader/fq/pkg/fqtest"
	"github.com/wader/fq/pkg/interp"
)

func TestFQTests(t *testing.T) {
	fqtest.TestPath(t, interp.DefaultRegistry)
}
