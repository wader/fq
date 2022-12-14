package decoders

import (
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

const intMask = 0b0111_1111
const intFlag = 0b1000_0000

// VarZigZag reads a variable length zigzag long from the current position in decoder
func VarZigZag(d *decode.D) int64 {
	var value uint64
	var shift uint
	size := 0
	for d.NotEnd() && size < 8 {
		size++
		b := byte(d.U8())
		value |= uint64(b&intMask) << shift
		if b&intFlag == 0 {
			return int64(value>>1) ^ -int64(value&1)
		}
		shift += 7
	}
	if size >= 8 {
		d.Fatalf("long exceeds 8 bytes")
	}
	d.Fatalf("unexpected end of data")
	return 0
}

func decodeLongFn(sms ...scalar.SintMapper) (DecodeFn, error) {
	// Int and long values are written using variable-length zig-zag coding.
	return func(name string, d *decode.D) any {
		return d.FieldSintFn(name, VarZigZag, sms...)
	}, nil
}
