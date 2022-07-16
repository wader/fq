package raw

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.RAW,
		Description: "Raw bits",
		DecodeFn:    func(d *decode.D, in any) any { return nil },
	})
}
