package binwriter

import (
	"io"

	"github.com/wader/fq/pkg/bitio"
)

type Writer struct {
	w               io.Writer
	width           int
	startLineOffset int
	offset          int
	fn              func(b byte) string
}

func New(w io.Writer, width int, startLineOffset int, fn func(b byte) string) *Writer {
	return &Writer{
		w:               w,
		width:           width,
		startLineOffset: startLineOffset,
		offset:          0,
		fn:              fn,
	}
}

func (w *Writer) WriteBits(p []byte, nBits int64) (n int64, err error) {
	for w.offset < w.startLineOffset {
		b := []byte(" ")
		if w.offset%w.width == w.width-1 {
			b = []byte(" \n")
		}
		if _, err := w.w.Write(b); err != nil {
			return 0, err
		}
		w.offset++
	}

	for i := int64(0); i < nBits; i++ {
		var v byte
		if bitio.Read64(p, i, 1) == 1 {
			v = 1
		}
		if _, err := w.w.Write([]byte(w.fn(v))); err != nil {
			return 0, err
		}

		w.offset++
		if w.offset%w.width == 0 {
			if _, err := w.w.Write([]byte("\n")); err != nil {
				return 0, err
			}
		}
	}

	return nBits, nil
}
