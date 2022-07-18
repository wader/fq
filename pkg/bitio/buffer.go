package bitio

// TODO: NewBuffer with []byte arg to save on alloc

import (
	"io"
)

// Buffer is a bitio.Reader and bitio.Writer providing a bit buffer.
// Similar to bytes.Buffer.
type Buffer struct {
	buf     []byte
	bufBits int64
	bitsOff int64
}

func (b *Buffer) Len() int64 { return b.bufBits - b.bitsOff }

func (b *Buffer) Reset() {
	b.bufBits = 0
	b.bitsOff = 0
}

// Bits return unread bits in buffer
func (b *Buffer) Bits() ([]byte, int64) {
	l := b.Len()
	buf := make([]byte, BitsByteCount(l))
	copyBufBits(buf, 0, b.buf, b.bitsOff, l, true)
	return buf, b.bufBits
}

func (b *Buffer) WriteBits(p []byte, nBits int64) (n int64, err error) {
	tBytes := BitsByteCount(b.bufBits + nBits)

	if tBytes > int64(len(b.buf)) {
		if tBytes <= int64(cap(b.buf)) {
			b.buf = b.buf[:tBytes]
		} else {
			buf := make([]byte, tBytes, tBytes*2)
			copy(buf, b.buf)
			b.buf = buf
		}
	}

	copyBufBits(b.buf, b.bufBits, p, 0, nBits, true)
	b.bufBits += nBits

	return nBits, nil
}

func (b *Buffer) empty() bool { return b.bufBits <= b.bitsOff }

func (b *Buffer) ReadBits(p []byte, nBits int64) (n int64, err error) {
	if b.empty() {
		b.Reset()
		if nBits == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	c := nBits
	left := b.Len()
	if c > left {
		c = left
	}

	copyBufBits(p, 0, b.buf, b.bitsOff, c, true)
	b.bitsOff += c

	return c, nil
}
