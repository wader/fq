package asciiwriter

import (
	"io"
)

type Writer struct {
	w           io.Writer
	width       int
	startOffset int
	offset      int
	buf         []byte
	bufOffset   int
}

func New(w io.Writer, width int, startOffset int) *Writer {
	return &Writer{
		w:           w,
		width:       width,
		startOffset: startOffset,
		offset:      0,
		buf:         make([]byte, width+2), // worst case " " or "\n" + width + "\n"
		bufOffset:   0,
	}
}

func (h *Writer) Write(p []byte) (n int, err error) {
	for h.offset < h.startOffset {
		b := []byte(" ")
		if h.offset%h.width == h.width-1 {
			b = []byte("\n")
		}
		if _, err := h.w.Write(b); err != nil {
			return 0, err
		}
		h.offset++
	}

	if h.offset > h.startOffset && h.offset%h.width == 0 {
		h.buf[0] = '\n'
		h.bufOffset = 1
	}

	for i := 0; i < len(p); i++ {
		lineOffset := h.offset % h.width
		v := p[i]
		if v < 32 || v > 126 {
			h.buf[h.bufOffset] = '.'
		} else {
			h.buf[h.bufOffset] = v
		}
		h.bufOffset++

		var b []byte
		switch {
		case i < len(p)-1 && lineOffset == h.width-1:
			h.buf[h.bufOffset] = '\n'
			h.bufOffset++
			b = h.buf[:h.bufOffset]
		case i == len(p)-1:
			b = h.buf[:h.bufOffset]
		}

		// log.Printf("i=%d h.bufOffset=%d lineOffset=%d h.width-1=%d b=%q\n", i, h.bufOffset, lineOffset, h.width-1, b)
		if b != nil {
			if _, err := h.w.Write(b); err != nil {
				return 0, err
			}
			h.bufOffset = 0
		}

		h.offset++
	}

	return len(p), nil
}
