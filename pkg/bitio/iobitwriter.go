package bitio

import (
	"io"
)

// IOBitWriter is a bitio.Writer that writes to a io.Writer.
// Use Flush to write possible unaligned byte zero bit padded.
type IOBitWriter struct {
	w io.Writer
	b Buffer
}

// NewIOBitWriter returns a new bitio.IOBitWriter.
func NewIOBitWriter(w io.Writer) *IOBitWriter {
	return &IOBitWriter{w: w}
}

func (w *IOBitWriter) WriteBits(p []byte, nBits int64) (n int64, err error) {
	if n, err = w.b.WriteBits(p, nBits); err != nil {
		return n, err
	}

	var sn int64

	for {
		l := w.b.Len()
		if l < 8 {
			break
		}

		var buf [32 * 1024]byte
		n, err := w.b.ReadBits(buf[:], l-(l%8))
		if err != nil {
			return sn, err
		}

		rn, err := w.w.Write(buf[:n/8])
		sn += int64(rn)
		if err != nil {
			return sn, err
		}

		sn += n
	}

	return nBits, nil
}

// Flush write possible unaligned byte zero bit padded.
func (w *IOBitWriter) Flush() error {
	if w.b.Len() == 0 {
		return nil
	}

	var buf [1]byte
	_, err := w.b.ReadBits(buf[:], w.b.Len())
	if err != nil {
		return err
	}
	_, err = w.w.Write(buf[:])

	return err
}

func (w *IOBitWriter) Unwrap() any {
	return w.w
}
