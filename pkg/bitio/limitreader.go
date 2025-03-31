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

func (r *LimitReader) ReadBits(p []byte, nBits int64) (n int64, err error) {
	if r.n <= 0 {
		return 0, io.EOF
	}
	if nBits > r.n {
		nBits = r.n
	}
	n, err = r.r.ReadBits(p, nBits)
	r.n -= n
	return n, err
}

func (r *LimitReader) CloneReader() (Reader, error) {
	rc, err := CloneReader(r.r)
	if err != nil {
		return nil, err
	}
	return &LimitReader{r: rc, n: r.n}, nil
}

func (r *LimitReader) Unwrap() any {
	return r.r
}
