package registry

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/decode/registry"
)

// Default global registry that all builtin formats register with
var Default = registry.New()

func MustRegister(format decode.Format) {
	Default.MustRegister(format)
}
