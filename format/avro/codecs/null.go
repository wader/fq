package codecs

import (
	"github.com/wader/fq/format/avro/schema"
	"github.com/wader/fq/pkg/decode"
)

const intMask = byte(127)
const intFlag = byte(128)

// readLong reads a variable length zig zag long from the current position in decoder
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

type LongCodec struct {}

func (l LongCodec) Decode(d *decode.D) interface{}{
	return VarZigZag(d)
}

func (l LongCodec) Type() CodecType {
	return SCALAR
}

func BuildLongCodec(schema schema.SimplifiedSchema) (Codec, error) {
	return &LongCodec{}, nil
}
