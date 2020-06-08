package registry

import (
	"fq/pkg/decode"
)

// Default global registry that all standard formats register with
var Default = New()

func MustRegister(format *decode.Format) *decode.Format {
	return Default.MustRegister(format)
}
