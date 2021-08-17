package hexpairwriter

// TODO: generalize and rename? make buffer more flexible

import (
	"io"

	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/bitio"
)

type Writer struct {
	w               io.Writer
	width           int
	startLineOffset int
	fn              func(v byte) string
	offset          int
	buf             []byte
	bufOffset       int

	bitsBuf  []byte
	bitsBufN int
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
		// TODO: ansi length? nicer reusable buffer?
		buf:       make([]byte, width*200+1), // worst case " " or "\n" + width*(XX "+ansi) + "\n"
		bufOffset: 0,
		bitsBuf:   make([]byte, 1),
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

func (h *Writer) WriteBits(p []byte, nBits int) (n int, err error) {
	pos := 0
	rBits := nBits
	if h.bitsBufN > 0 {
		r := num.MinInt(8-h.bitsBufN, nBits)
		v := bitio.Read64(p, 0, r)
		bitio.Write64(v, r, h.bitsBuf, h.bitsBufN)

		h.bitsBufN += r

		if h.bitsBufN < 8 {
			return nBits, nil
		}
		if n, err := h.Write(h.bitsBuf); err != nil {
			return n * 8, err
		}
		pos = r
		rBits -= r
	}

	for rBits >= 8 {
		b := [1]byte{0}

		b[0] = byte(bitio.Read64(p, pos, 8))
		if n, err := h.Write(b[:]); err != nil {
			return n * 8, err
		}

		pos += 8
		rBits -= 8

	}

	if rBits > 0 {
		h.bitsBuf[0] = byte(bitio.Read64(p, pos, rBits)) << (8 - rBits)
		h.bitsBufN = rBits
	} else {
		h.bitsBufN = 0
	}

	return nBits, nil
}
