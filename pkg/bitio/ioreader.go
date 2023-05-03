package bitio

import (
	"errors"
	"io"
)

// IOReader is a io.Reader and io.ByteReader that reads from a bitio.Reader.
// Unaligned byte at EOF will be zero bit padded.
type IOReader struct {
	r    Reader
	rErr error
	b    Buffer
}

// NewIOReader returns a new bitio.IOReader.
func NewIOReader(r Reader) *IOReader {
	return &IOReader{r: r}
}

func (r *IOReader) Read(p []byte) (n int, err error) {
	var ns int64

	for {
		// uses p even if returning nothing, io.Reader docs says:
		// "it may use all of p as scratch space during the call"
		if r.rErr == nil {
			var rn int64
			rn, err = r.r.ReadBits(p, int64(len(p))*8)
			r.rErr = err
			ns += rn
			_, err = r.b.WriteBits(p, rn)
			if err != nil {
				return 0, err
			}
		}

		if r.b.Len() >= 8 {
			// read whole bytes
			rBits := int64(len(p)) * 8
			bBits := r.b.Len()
			aBits := bBits - bBits%8
			if rBits > aBits {
				rBits = aBits
			}

			rn, rErr := r.b.ReadBits(p, rBits)
			if rErr != nil {
				return int(rn / 8), rErr
			}
			return int(rn / 8), nil
		} else if r.rErr != nil {
			if errors.Is(r.rErr, io.EOF) && r.b.Len() > 0 {
				// TODO: hmm io.Buffer does this
				if len(p) == 0 {
					return 0, nil
				}

				_, err := r.b.ReadBits(p, r.b.Len())
				if err != nil {
					return 0, err
				}
				return 1, r.rErr
			}
			return 0, r.rErr
		}
	}
}

// required to make some readers like deflate not do their own buffering
func (r *IOReader) ReadByte() (byte, error) {
	var rb [1]byte
	_, err := r.Read(rb[:])
	return rb[0], err
}

func (r *IOReader) Unwrap() any {
	return r.r
}
