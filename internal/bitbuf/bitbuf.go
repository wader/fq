package bitbuf

import (
	"fmt"
	"strings"
)

type Buffer struct {
	Buf         []byte
	BufFirstBit uint64
	Len         uint64
	Pos         uint64
}

func New(firstBit uint64, buf []byte, lenBits uint64) *Buffer {
	return &Buffer{
		Buf:         buf,
		BufFirstBit: firstBit,
		Len:         lenBits,
		Pos:         0,
	}
}

func NewFromBytes(buf []byte) *Buffer {
	return New(0, buf, uint64(len(buf)*8))
}

func NewFromBitString(s string) *Buffer {
	var buf []byte
	i := 0
	var n byte
	for ; i < len(s); i++ {
		c := s[i]
		b := 0
		if c == '0' {
			b = 0
		} else if c == '1' {
			b = 1
		} else {
			panic(fmt.Sprintf("invalid bit string %q at index %d %q", s, i, c))
		}

		p := 8 - (i % 8) - 1
		n |= byte(b) << p
		if (i > 0 && p == 0) || i == len(s)-1 {
			buf = append(buf, n)
			n = 0
		}
	}

	return New(0, buf, uint64(len(s)))
}

func (b *Buffer) Bits(nBits uint64) (uint64, uint64) {
	p := uint64(b.Pos) + uint64(nBits)
	if p > b.Len {
		return 0, uint64(p) - b.Len
	}

	n := ReadBits(b.Buf, b.BufFirstBit+b.Pos, nBits)
	b.Pos += nBits

	return n, nBits
}

func (b *Buffer) BitBufRange(start uint64, nBits uint64) (*Buffer, uint64) {
	endPos := uint64(start) + uint64(nBits)
	if endPos > b.Len {
		return nil, endPos - b.Len
	}

	nb := &Buffer{
		Buf:         b.Buf,
		BufFirstBit: b.BufFirstBit + start,
		Len:         nBits,
		Pos:         0,
	}
	b.Pos += nBits

	return nb, nBits
}

func (b *Buffer) BitBufLen(nBits uint64) (*Buffer, uint64) {
	return b.BitBufRange(b.Pos, nBits)
}

func (b *Buffer) BytesRange(firstBit uint64, nBytes uint64) ([]byte, uint64) {
	endPos := firstBit + nBytes*8
	if endPos > b.Len {
		return nil, endPos - b.Len
	}

	bufFirstBit := b.BufFirstBit + firstBit
	if bufFirstBit%8 == 0 {
		nb := b.Buf[bufFirstBit : bufFirstBit+nBytes]
		b.Pos += nBytes * 8

		return nb, nBytes * 8
	}

	var buf []byte
	for i := uint64(0); i < nBytes; i++ {
		buf = append(buf, byte(ReadBits(b.Buf, bufFirstBit+i, 8)))
	}
	b.Pos += nBytes * 8

	return buf, nBytes * 8
}

func (b *Buffer) BytesLen(nBytes uint64) ([]byte, uint64) {
	return b.BytesRange(b.Pos, nBytes)
}

func (b *Buffer) End() bool {
	return b.Pos >= b.Len
}

func (b *Buffer) ByteAlignBits() uint {
	return uint((8 - (b.Pos & 0x7)) & 0x7)
}

func (b *Buffer) BytePos() uint64 {
	return b.Pos & 0x7
}

func (b *Buffer) String() string {
	truncLen, truncS := b.Len, ""
	if truncLen > 64 {
		truncLen, truncS = 64, "..."
	}
	truncBB, _ := b.BitBufLen(truncLen)

	return fmt.Sprintf("0b%s%s /* %d bits */", truncBB.BitString(), truncS, b.Len)
}

func (b *Buffer) BitString() string {
	var ss []string
	for !b.End() {
		if n, _ := b.Bits(1); n == 0 {
			ss = append(ss, "0")
		} else {
			ss = append(ss, "1")
		}
	}

	return strings.Join(ss, "")
}
