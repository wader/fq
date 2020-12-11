package raw

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.RAW,
		Description: "Raw bits",
		DecodeFn:    func(d *decode.D) interface{} { return nil },
	})
}
