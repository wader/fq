package bitbuf

import (
	"fmt"
	"strings"
)

type Endian int

const (
	BigEndian Endian = iota
	LittleEndian
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

func (b *Buffer) ByteAlignBits() uint64 {
	return (8 - (b.Pos & 0x7)) & 0x7
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

func (b *Buffer) UE(nBits uint64, endian Endian) uint64 {
	n, _ := b.Bits(nBits)
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}

	return n
}

func (b *Buffer) Bool() bool { return b.UE(1, BigEndian) == 1 }

func (b *Buffer) U(nBits uint64) uint64 { return b.UE(nBits, BigEndian) }
func (b *Buffer) U1() uint64            { return b.UE(1, BigEndian) }
func (b *Buffer) U2() uint64            { return b.UE(2, BigEndian) }
func (b *Buffer) U3() uint64            { return b.UE(3, BigEndian) }
func (b *Buffer) U4() uint64            { return b.UE(4, BigEndian) }
func (b *Buffer) U5() uint64            { return b.UE(5, BigEndian) }
func (b *Buffer) U6() uint64            { return b.UE(6, BigEndian) }
func (b *Buffer) U7() uint64            { return b.UE(7, BigEndian) }
func (b *Buffer) U8() uint64            { return b.UE(8, BigEndian) }
func (b *Buffer) U9() uint64            { return b.UE(9, BigEndian) }
func (b *Buffer) U10() uint64           { return b.UE(10, BigEndian) }
func (b *Buffer) U11() uint64           { return b.UE(11, BigEndian) }
func (b *Buffer) U12() uint64           { return b.UE(12, BigEndian) }
func (b *Buffer) U13() uint64           { return b.UE(13, BigEndian) }
func (b *Buffer) U14() uint64           { return b.UE(14, BigEndian) }
func (b *Buffer) U15() uint64           { return b.UE(15, BigEndian) }
func (b *Buffer) U16() uint64           { return b.UE(16, BigEndian) }
func (b *Buffer) U24() uint64           { return b.UE(24, BigEndian) }
func (b *Buffer) U32() uint64           { return b.UE(32, BigEndian) }
func (b *Buffer) U64() uint64           { return b.UE(64, BigEndian) }

func (b *Buffer) UBE(nBits uint64) uint64 { return b.UE(nBits, BigEndian) }
func (b *Buffer) U9BE() uint64            { return b.UE(9, BigEndian) }
func (b *Buffer) U10BE() uint64           { return b.UE(10, BigEndian) }
func (b *Buffer) U11BE() uint64           { return b.UE(11, BigEndian) }
func (b *Buffer) U12BE() uint64           { return b.UE(12, BigEndian) }
func (b *Buffer) U13BE() uint64           { return b.UE(13, BigEndian) }
func (b *Buffer) U14BE() uint64           { return b.UE(14, BigEndian) }
func (b *Buffer) U15BE() uint64           { return b.UE(15, BigEndian) }
func (b *Buffer) U16BE() uint64           { return b.UE(16, BigEndian) }
func (b *Buffer) U24BE() uint64           { return b.UE(24, BigEndian) }
func (b *Buffer) U32BE() uint64           { return b.UE(32, BigEndian) }
func (b *Buffer) U64BE() uint64           { return b.UE(64, BigEndian) }

func (b *Buffer) ULE(nBits uint64) uint64 { return b.UE(nBits, LittleEndian) }
func (b *Buffer) U9LE() uint64            { return b.UE(9, LittleEndian) }
func (b *Buffer) U10LE() uint64           { return b.UE(10, LittleEndian) }
func (b *Buffer) U11LE() uint64           { return b.UE(11, LittleEndian) }
func (b *Buffer) U12LE() uint64           { return b.UE(12, LittleEndian) }
func (b *Buffer) U13LE() uint64           { return b.UE(13, LittleEndian) }
func (b *Buffer) U14LE() uint64           { return b.UE(14, LittleEndian) }
func (b *Buffer) U15LE() uint64           { return b.UE(15, LittleEndian) }
func (b *Buffer) U16LE() uint64           { return b.UE(16, LittleEndian) }
func (b *Buffer) U24LE() uint64           { return b.UE(24, LittleEndian) }
func (b *Buffer) U32LE() uint64           { return b.UE(32, LittleEndian) }
func (b *Buffer) U64LE() uint64           { return b.UE(64, LittleEndian) }

func (b *Buffer) SE(nBits uint64, endian Endian) int64 {
	n, _ := b.Bits(nBits)
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}
	var s int64
	if n&(1<<(nBits-1)) > 0 {
		// two's complement
		s = -int64((^n & ((1 << nBits) - 1)) + 1)
	} else {
		s = int64(n)
	}

	return s
}

func (b *Buffer) S(nBits uint64) int64 { return b.SE(nBits, BigEndian) }
func (b *Buffer) S1() int64            { return b.SE(1, BigEndian) }
func (b *Buffer) S2() int64            { return b.SE(2, BigEndian) }
func (b *Buffer) S3() int64            { return b.SE(3, BigEndian) }
func (b *Buffer) S4() int64            { return b.SE(4, BigEndian) }
func (b *Buffer) S5() int64            { return b.SE(5, BigEndian) }
func (b *Buffer) S6() int64            { return b.SE(6, BigEndian) }
func (b *Buffer) S7() int64            { return b.SE(7, BigEndian) }
func (b *Buffer) S8() int64            { return b.SE(8, BigEndian) }
func (b *Buffer) S9() int64            { return b.SE(9, BigEndian) }
func (b *Buffer) S10() int64           { return b.SE(10, BigEndian) }
func (b *Buffer) S11() int64           { return b.SE(11, BigEndian) }
func (b *Buffer) S12() int64           { return b.SE(12, BigEndian) }
func (b *Buffer) S13() int64           { return b.SE(13, BigEndian) }
func (b *Buffer) S14() int64           { return b.SE(14, BigEndian) }
func (b *Buffer) S15() int64           { return b.SE(15, BigEndian) }
func (b *Buffer) S16() int64           { return b.SE(16, BigEndian) }
func (b *Buffer) S24() int64           { return b.SE(24, BigEndian) }
func (b *Buffer) S32() int64           { return b.SE(32, BigEndian) }
func (b *Buffer) S64() int64           { return b.SE(64, BigEndian) }

func (b *Buffer) SBE(nBits uint64) int64 { return b.SE(nBits, BigEndian) }
func (b *Buffer) S9BE() int64            { return b.SE(9, BigEndian) }
func (b *Buffer) S10BE() int64           { return b.SE(10, BigEndian) }
func (b *Buffer) S11BE() int64           { return b.SE(11, BigEndian) }
func (b *Buffer) S12BE() int64           { return b.SE(12, BigEndian) }
func (b *Buffer) S13BE() int64           { return b.SE(13, BigEndian) }
func (b *Buffer) S14BE() int64           { return b.SE(14, BigEndian) }
func (b *Buffer) S15BE() int64           { return b.SE(15, BigEndian) }
func (b *Buffer) S16BE() int64           { return b.SE(16, BigEndian) }
func (b *Buffer) S24BE() int64           { return b.SE(24, BigEndian) }
func (b *Buffer) S32BE() int64           { return b.SE(32, BigEndian) }
func (b *Buffer) S64BE() int64           { return b.SE(64, BigEndian) }

func (b *Buffer) SLE(nBits uint64) int64 { return b.SE(nBits, LittleEndian) }
func (b *Buffer) S9LE() int64            { return b.SE(9, LittleEndian) }
func (b *Buffer) S10LE() int64           { return b.SE(10, LittleEndian) }
func (b *Buffer) S11LE() int64           { return b.SE(11, LittleEndian) }
func (b *Buffer) S12LE() int64           { return b.SE(12, LittleEndian) }
func (b *Buffer) S13LE() int64           { return b.SE(13, LittleEndian) }
func (b *Buffer) S14LE() int64           { return b.SE(14, LittleEndian) }
func (b *Buffer) S15LE() int64           { return b.SE(15, LittleEndian) }
func (b *Buffer) S16LE() int64           { return b.SE(16, LittleEndian) }
func (b *Buffer) S24LE() int64           { return b.SE(24, LittleEndian) }
func (b *Buffer) S32LE() int64           { return b.SE(32, LittleEndian) }
func (b *Buffer) S64LE() int64           { return b.SE(64, LittleEndian) }

func (b *Buffer) UTF8(nBytes uint64) string {
	// TODO: panic?
	s, _ := b.BytesLen(nBytes)
	return string(s)
}

func (b *Buffer) Unary(s uint) uint {
	var n uint
	for uint(b.U1()) == s {
		n++
	}
	return n
}
