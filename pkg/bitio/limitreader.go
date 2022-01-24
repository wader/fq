package bitio

import (
	"io"
)

func NewLimitReader(r Reader, n int64) *LimitReader { return &LimitReader{r, n} }

// LimitReader is a bitio.Reader that reads a limited amount of bits from a bitio.Reader
// Similar to bytes.LimitedReader but for bits
type LimitReader struct {
	r Reader
	n int64
}

func (l *LimitReader) ReadBits(p []byte, nBits int64) (n int64, err error) {
	if l.n <= 0 {
		return 0, io.EOF
	}
	if nBits > l.n {
		nBits = l.n
	}
	n, err = l.r.ReadBits(p, nBits)
	l.n -= n
	return n, err
}
