package bits

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed bits.md
//go:embed bytes.md
var bitsFS embed.FS

func decodeBits(unit int) func(d *decode.D) any {
	return func(d *decode.D) any {
		var s scalar.Any
		b, _ := interp.NewBinaryFromBitReader(d.BitBufRange(0, d.Len()), unit, 0)
		s.Actual = b
		d.Value.V = &s
		d.Value.Range.Len = d.Len()
		return nil
	}
}

func init() {
	interp.RegisterFormat(
		format.Bits,
		&decode.Format{
			Description:        "Raw bits",
			DecodeFn:           decodeBits(1),
			SkipDecodeFunction: true, // skip add bits and frombits function
		})
	interp.RegisterFormat(
		format.Bytes,
		&decode.Format{
			Description:        "Raw bytes",
			DecodeFn:           decodeBits(8),
			SkipDecodeFunction: true, // skip add bytes and frombytes function
		})
	interp.RegisterFS(bitsFS)
}
