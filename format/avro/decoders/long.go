package decoders

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

const intMask = byte(127)
const intFlag = byte(128)

// VarZigZag reads a variable length zigzag long from the current position in decoder
func VarZigZag(d *decode.D) int64 {
	var value uint64
	var shift uint
	for d.NotEnd() {
		b := byte(d.U8())
		value |= uint64(b&intMask) << shift
		if b&intFlag == 0 {
			return int64(value>>1) ^ -int64(value&1)
		}
		shift += 7
	}
	panic("unexpected end of data")
}

func decodeLongFn(schema schema.SimplifiedSchema) (func(string, *decode.D), error) {
	// int and long values are written using variable-length zig-zag coding.
	return func(name string, d *decode.D) {
		d.FieldSFn(name, VarZigZag)
	}, nil
}
