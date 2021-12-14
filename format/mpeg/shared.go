package mpeg

import (
	"github.com/wader/fq/pkg/decode"
)

func decodeEscapeValueFn(add int, b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return func(d *decode.D) uint64 {
		n1 := d.U(b1)
		n := n1
		if n1 == (1<<b1)-1 {
			n2 := d.U(b2)
			if add != -1 {
				n += n2 + uint64(add)
			} else {
				n = n2
			}
			if n2 == (1<<b2)-1 {
				n3 := d.U(b3)
				if add != -1 {
					n += n3 + uint64(add)
				} else {
					n = n3
				}
			}
		}
		return n
	}
}

// use last non-escaped value
func decodeEscapeValueAbsFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(-1, b1, b2, b3)
}

// add values and escaped values
//nolint: deadcode,unused
func decodeEscapeValueAddFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(0, b1, b2, b3)
}

// add values and escaped values+1
func decodeEscapeValueCarryFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(1, b1, b2, b3)
}
