package bitiox

import (
	"io"

	"github.com/wader/fq/pkg/bitio"
)

type ZeroReadAtSeeker struct {
	pos   int64
	nBits int64
}

func NewZeroAtSeeker(nBits int64) *ZeroReadAtSeeker {
	return &ZeroReadAtSeeker{nBits: nBits}
}

func (z *ZeroReadAtSeeker) SeekBits(bitOffset int64, whence int) (int64, error) {
	p := z.pos
	switch whence {
	case io.SeekStart:
		p = bitOffset
	case io.SeekCurrent:
		p += bitOffset
	case io.SeekEnd:
		p = z.nBits + bitOffset
	default:
		panic("unknown whence")
	}

	if p < 0 || p > z.nBits {
		return z.pos, bitio.ErrOffset
	}
	z.pos = p

	return p, nil
}

func (z *ZeroReadAtSeeker) ReadBitsAt(p []byte, nBits int64, bitOff int64) (n int64, err error) {
	if bitOff < 0 || bitOff > z.nBits {
		return 0, bitio.ErrOffset
	}
	if bitOff == z.nBits {
		return 0, io.EOF
	}

	lBits := z.nBits - bitOff
	rBits := nBits
	if rBits > lBits {
		rBits = lBits
	}
	rBytes := bitio.BitsByteCount(rBits)
	for i := int64(0); i < rBytes; i++ {
		p[i] = 0
	}

	return rBits, nil
}

func (z *ZeroReadAtSeeker) CloneReadAtSeeker() (bitio.ReadAtSeeker, error) {
	return NewZeroAtSeeker(z.nBits), nil
}
