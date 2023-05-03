package bitio

import (
	"io"
)

// LimitReader is a bitio.Reader that reads a limited amount of bits from a bitio.Reader.
// Similar to bytes.LimitedReader.
type LimitReader struct {
	r Reader
	n int64
}

// NewLimitReader returns a new bitio.LimitReader.
func NewLimitReader(r Reader, n int64) *LimitReader { return &LimitReader{r, n} }

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

func (l *LimitReader) CloneReader() (Reader, error) {
	rc, err := CloneReader(l.r)
	if err != nil {
		return nil, err
	}
	return &LimitReader{r: rc, n: l.n}, nil
}

func (r *LimitReader) Unwrap() any {
	return r.r
}
