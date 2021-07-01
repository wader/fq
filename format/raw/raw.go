package raw

import (
	"fq/format"
	"fq/pkg/decode"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.RAW,
		Description: "Raw bits",
		DecodeFn:    func(d *decode.D, in interface{}) interface{} { return nil },
	})
}
