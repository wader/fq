package bitio

import "fmt"

// ReverseBytes64 reverses the bytes part of the lowest nBits.
// Similar to bits.ReverseBytes64 but only rotates the lowest bytes and rest of bytes will be zero.
func ReverseBytes64(nBits int, n uint64) uint64 {
	switch {
	case nBits <= 8:
		return n
	case nBits <= 16:
		return n&0xff00>>8 | n&0xff<<8
	case nBits <= 24:
		return n&0xff<<16 | n&0xff00 | n&0xff0000>>16
	case nBits <= 32:
		return n&0xff<<24 | n&0xff00<<8 | n&0xff0000>>8 | n&0xff000000>>24
	case nBits <= 40:
		return n&0xff<<32 | n&0xff00<<16 | n&0xff0000 | n&0xff000000>>16 | n&0xff00000000>>32
	case nBits <= 48:
		return n&0xff<<40 | n&0xff00<<24 | n&0xff0000<<8 | n&0xff000000>>8 | n&0xff00000000>>24 | n&0xff0000000000>>40
	case nBits <= 56:
		return n&0xff<<48 | n&0xff00<<32 | n&0xff0000<<16 | n&0xff000000 | n&0xff00000000>>16 | n&0xff0000000000>>32 | n&0xff000000000000>>48
	case nBits <= 64:
		return n&0xff<<56 | n&0xff00<<40 | n&0xff0000<<24 | n&0xff000000<<8 | n&0xff00000000>>8 | n&0xff0000000000>>24 | n&0xff000000000000>>40 | n&0xff00000000000000>>56
	default:
		panic(fmt.Sprintf("unsupported bit length %d", nBits))
	}
}
