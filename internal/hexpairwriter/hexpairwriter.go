package hexpairwriter

// TODO: generalize and rename? make buffer more flexible

import (
	"io"
)

type Writer struct {
	w               io.Writer
	width           int
	startLineOffset int
	fn              func(v byte) string
	offset          int
	buf             []byte
	bufOffset       int
}

// for i in $(seq 0 255) ; do printf "%02x" $i ; done | while read -n 32 s ; do echo "\"$s\""+; done
const hexstring = "" +
	"000102030405060708090a0b0c0d0e0f" +
	"101112131415161718191a1b1c1d1e1f" +
	"202122232425262728292a2b2c2d2e2f" +
	"303132333435363738393a3b3c3d3e3f" +
	"404142434445464748494a4b4c4d4e4f" +
	"505152535455565758595a5b5c5d5e5f" +
	"606162636465666768696a6b6c6d6e6f" +
	"707172737475767778797a7b7c7d7e7f" +
	"808182838485868788898a8b8c8d8e8f" +
	"909192939495969798999a9b9c9d9e9f" +
	"a0a1a2a3a4a5a6a7a8a9aaabacadaeaf" +
	"b0b1b2b3b4b5b6b7b8b9babbbcbdbebf" +
	"c0c1c2c3c4c5c6c7c8c9cacbcccdcecf" +
	"d0d1d2d3d4d5d6d7d8d9dadbdcdddedf" +
	"e0e1e2e3e4e5e6e7e8e9eaebecedeeef" +
	"f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff"

func Pair(c byte) string {
	return hexstring[int(c)*2 : int(c)*2+2]
}

func New(w io.Writer, width int, startLineOffset int, fn func(b byte) string) *Writer {
	return &Writer{
		w:               w,
		width:           width,
		startLineOffset: startLineOffset,
		fn:              fn,
		offset:          0,
		buf:             make([]byte, width*12+1), // worst case " " or "\n" + width*(XX "+ansi) + "\n"
		bufOffset:       0,
	}
}

func (h *Writer) Write(p []byte) (n int, err error) {
	for h.offset < h.startLineOffset {
		b := []byte("   ")
		if h.offset%h.width == h.width-1 {
			b = []byte("  \n")
		}
		if _, err := h.w.Write(b); err != nil {
			return 0, err
		}
		h.offset++
	}

	if h.offset > h.startLineOffset {
		if h.offset%h.width == 0 {
			h.buf[0] = '\n'
		} else {
			h.buf[0] = ' '
		}
		h.bufOffset = 1
	}

	for i := 0; i < len(p); i++ {
		lineOffset := h.offset % h.width
		s := []byte(h.fn(p[i]))
		copy(h.buf[h.bufOffset:], s)
		h.bufOffset += len(s)
		h.buf[h.bufOffset] = ' '
		h.bufOffset++

		var b []byte
		switch {
		case i < len(p)-1 && lineOffset == h.width-1:
			h.buf[h.bufOffset-1] = '\n'
			b = h.buf[:h.bufOffset]
		case i == len(p)-1:
			b = h.buf[:h.bufOffset-1]
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
