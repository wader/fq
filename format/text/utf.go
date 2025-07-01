package text

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/text/unicode/runenames"
)

func init() {
	interp.RegisterFormat(
		format.UTF,
		&decode.Format{
			Description: "Unicode Transformation Format",
			DecodeFn:    utfDecode,
			DefaultInArg: format.UTF_In{
				Encoding: "utf8",
			},
			RootArray: true,
		})
}

func decodeUTF8Codepoint(d *decode.D) uint64 {
	b0 := d.PeekBytes(1)[0]
	s := 0
	switch {
	case (b0 & 0b1000_0000) == 0:
		s = 1
	case (b0 & 0b1100_0000) == 0b1100_0000:
		s = 2
	case (b0 & 0b1110_0000) == 0b1110_0000:
		s = 3
	case (b0 & 0b1111_0000) == 0b1111_0000:
		s = 4
	}
	return uint64([]rune(d.UTF8(s))[0])
}

func utfDecode(d *decode.D) any {
	for !d.End() {
		d.FieldUintFn("codepoint", decodeUTF8Codepoint, scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
			var r = rune(s.Actual)
			s.Sym = string([]rune{r})
			s.Description = fmt.Sprintf("U+%.4x %s", r, runenames.Name(r))
			return s, nil
		}))
	}

	return nil
}
