package bits

import (
	"embed"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqex"
	"github.com/wader/fq/internal/mathex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed bits.md
//go:embed bits.jq
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
	interp.RegisterFunc2("_from_float", func(_ *interp.Interp, c any, nBits int, isLE bool) any {
		switch nBits {
		case 16, 32, 64:
		default:
			return fmt.Errorf("unsupported bit size %d, must be 16, 32 or 64", nBits)
		}

		br, err := interp.ToBitReader(c)
		if err != nil {
			return err
		}
		var b [8]byte
		bs := b[:][0 : nBits/8]
		_, err = br.ReadBits(bs[:], int64(nBits))
		if err != nil {
			return err
		}
		if isLE {
			decode.ReverseBytes(bs[:])
		}

		switch nBits {
		case 64:
			return math.Float64frombits(binary.BigEndian.Uint64(bs[:]))
		case 32:
			return float64(math.Float32frombits(binary.BigEndian.Uint32(bs[:])))
		case 16:
			return float64(mathex.Float16(binary.BigEndian.Uint16(bs[:])).Float32())
		default:
			panic("unreachable")
		}
	})
	interp.RegisterFunc2("_to_float", func(_ *interp.Interp, c any, nBits int, isLE bool) any {
		switch nBits {
		case 16, 32, 64:
		default:
			return fmt.Errorf("unsupported bit size %d, must be 16, 32 or 64", nBits)
		}

		v, ok := gojqex.Cast[float64](c)
		if !ok {
			return gojqex.FuncTypeError{Name: "_to_float", V: v}
		}

		var b [8]byte
		bs := b[:][0 : nBits/8]
		switch nBits {
		case 64:
			binary.BigEndian.PutUint64(bs, math.Float64bits(v))
		case 32:
			binary.BigEndian.PutUint32(bs, math.Float32bits(float32(v)))
		case 16:
			binary.BigEndian.PutUint16(bs, uint16(mathex.NewFloat16(float32(v))))
		default:
			panic("unreachable")
		}
		if isLE {
			decode.ReverseBytes(bs[:])
		}

		br, err := interp.NewBinaryFromBitReader(bitio.NewBitReader(bs, -1), 8, 0)
		if err != nil {
			return err
		}

		return br

	})

	interp.RegisterFS(bitsFS)
}
