package all

import (
	"fq/format/registry"
	"fq/pkg/decode"
)

var Registry = registry.New()

func MustRegister(format *decode.Format) *decode.Format {
	return Registry.MustRegister(format)
}
