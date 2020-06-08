package raw

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.RAW,
		Description: "Raw bits",
		DecodeFn:    func(d *decode.D, in interface{}) interface{} { return nil },
	})
}
