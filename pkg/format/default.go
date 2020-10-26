package format

import "fq/pkg/decode"

var DefaultRegistry = decode.NewRegistry()

func MustRegister(format *decode.Format) *decode.Format {
	return DefaultRegistry.MustRegister(format)
}
