package format_test

import (
	"flag"
	"testing"

	_ "github.com/wader/fq/format/all"
	"github.com/wader/fq/pkg/fqtest"
	"github.com/wader/fq/pkg/interp"
)

var update = flag.Bool("update", false, "Update tests")

func TestFormats(t *testing.T) {
	fqtest.TestPath(t, interp.DefaultRegistry, *update)
}
