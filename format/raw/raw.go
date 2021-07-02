package raw

import (
	"fq/format"
	"fq/format/all/all"
	"fq/pkg/decode"
)

func init() {
	all.MustRegister(&decode.Format{
		Name:        format.RAW,
		Description: "Raw bits",
		DecodeFn:    func(d *decode.D, in interface{}) interface{} { return nil },
	})
}
