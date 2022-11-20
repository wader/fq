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

func decodeBits(unit int) func(d *decode.D, _ any) any {
	return func(d *decode.D, _ any) any {
		var s scalar.S
		b, _ := interp.NewBinaryFromBitReader(d.BitBufRange(0, d.Len()), unit, 0)
		s.Actual = b
		d.Value.V = &s
		d.Value.Range.Len = d.Len()
		return nil
	}
}

func init() {
	interp.RegisterFormat(decode.Format{
		Name:               format.BITS,
		Description:        "Raw bits",
		DecodeFn:           decodeBits(1),
		SkipDecodeFunction: true, // skip add bits and frombits function
	})
	interp.RegisterFormat(decode.Format{
		Name:               format.BYTES,
		Description:        "Raw bytes",
		DecodeFn:           decodeBits(8),
		SkipDecodeFunction: true, // skip add bytes and frombytes function
	})
	interp.RegisterFS(bitsFS)
}
