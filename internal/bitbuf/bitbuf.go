package bitbuf

import (
	"fmt"
	"math"
	"strings"
)

// Endian byte order
type Endian int

const (
	// BigEndian byte order
	BigEndian Endian = iota
	// LittleEndian byte order
	LittleEndian
)

// Buffer is a bitbuf buffer
// TODO: make Buf/BufFirstBit private?
type Buffer struct {
	buf         []byte
	bufFirstBit uint64
	Len         uint64
	Pos         uint64
}

// New bitbuf.Buffer from byte buffer buf, start at firstBit with bit length lenBits
// buf is not copied.
func New(buf []byte, firstBit uint64, lenBits uint64) *Buffer {
	return &Buffer{
		buf:         buf,
		bufFirstBit: firstBit,
		Len:         lenBits,
		Pos:         0,
	}
}

// NewFromBytes bitbuf.Buffer from bytes
func NewFromBytes(buf []byte) *Buffer {
	return New(buf, 0, uint64(len(buf)*8))
}

// NewFromBitBuf bitbuf.Buffer from other bitbuf.Buffer
// Will be a shallow copy with position reset to zero.
func NewFromBitBuf(b *Buffer) *Buffer {
	return New(b.buf, b.bufFirstBit, b.Len)
}

// NewFromBitString bitbuf.Buffer from bit string, ex: "0101"
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

	return New(buf, 0, uint64(len(s)))
}

// Copy bitbuf
// TODO: rename? remove?
func (b *Buffer) Copy() *Buffer {
	return NewFromBitBuf(b)
}

// Bits reads nBits bits from buffer
func (b *Buffer) Bits(nBits uint64) (uint64, uint64) {
	p := uint64(b.Pos) + uint64(nBits)
	if p > b.Len {
		return 0, uint64(p) - b.Len
	}
	n := ReadBits(b.buf, b.bufFirstBit+b.Pos, nBits)
	b.Pos += nBits

	return n, nBits
}

// PeekBits peek nBits bits from buffer
// TODO: share code?
func (b *Buffer) PeekBits(nBits uint64) (uint64, uint64) {
	p := uint64(b.Pos) + uint64(nBits)
	if p > b.Len {
		return 0, uint64(p) - b.Len
	}
	n := ReadBits(b.buf, b.bufFirstBit+b.Pos, nBits)

	return n, nBits
}

// BitBufRange reads nBits bits starting from start
// Does not update current position.
func (b *Buffer) BitBufRange(start uint64, nBits uint64) (*Buffer, uint64) {
	endPos := uint64(start) + uint64(nBits)
	if endPos > b.Len {
		return nil, endPos - b.Len
	}

	nb := &Buffer{
		buf:         b.buf,
		bufFirstBit: b.bufFirstBit + start,
		Len:         nBits,
		Pos:         0,
	}

	return nb, nBits
}

// BitBufLen reads nBits
func (b *Buffer) BitBufLen(nBits uint64) (*Buffer, uint64) {
	bb, rBits := b.BitBufRange(b.Pos, nBits)
	b.Pos += rBits
	return bb, rBits
}

// BytesRange reads nBytes bytes starting bit position start
// Does not update current position.
func (b *Buffer) BytesRange(firstBit uint64, nBytes uint64) ([]byte, uint64) {
	endPos := firstBit + nBytes*8
	if endPos > b.Len {
		return nil, endPos - b.Len
	}

	bufFirstBit := b.bufFirstBit + firstBit
	if bufFirstBit%8 == 0 {
		bufFirstBytePos := bufFirstBit >> 3
		nb := b.buf[bufFirstBytePos : bufFirstBytePos+nBytes]
		return nb, nBytes * 8
	}

	var buf []byte
	for i := uint64(0); i < nBytes; i++ {
		buf = append(buf, byte(ReadBits(b.buf, bufFirstBit+i, 8)))
	}

	return buf, nBytes * 8
}

// BytesLen reads nBytes bytes
func (b *Buffer) BytesLen(nBytes uint64) ([]byte, uint64) {
	bb, rBits := b.BytesRange(b.Pos, nBytes)
	b.Pos += rBits
	return bb, rBits
}

// End is true if current position if at the end
func (b *Buffer) End() bool {
	return b.Pos >= b.Len
}

// BitsLeft number of bits left until end
func (b *Buffer) BitsLeft() uint64 {
	return b.Len - b.Pos
}

// ByteAlignBits number of bits to next byte align
func (b *Buffer) ByteAlignBits() uint64 {
	return (8 - (b.Pos & 0x7)) & 0x7
}

// BytePos byte position of current bit position
func (b *Buffer) BytePos() uint64 {
	return b.Pos & 0x7
}

// SeekRel relative to current bit position
// TODO: better name?
func (b *Buffer) SeekRel(delta int64) uint64 {
	// TODO: panic? bitbuf should never panic? return error, set error flag ignore rest?
	b.Pos = uint64(int64(b.Pos) + delta)
	return b.Pos
}

// SeekAbs to absolute position
func (b *Buffer) SeekAbs(pos uint64) uint64 {
	// TODO: panic? bitbuf should never panic?
	b.Pos = pos
	return b.Pos
}

func (b *Buffer) String() string {
	truncLen, truncS := b.Len, ""
	if truncLen > 64 {
		truncLen, truncS = 64, "..."
	}
	truncBB, _ := b.BitBufLen(truncLen)

	return fmt.Sprintf("0b%s%s /* %d bits */", truncBB.BitString(), truncS, b.Len)
}

// BitString return bit string representation
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

// UE reads nBits unsigned integer with byte order endian
func (b *Buffer) UE(nBits uint64, endian Endian) uint64 {
	n, _ := b.Bits(nBits)
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}

	return n
}

// Bool reads one bit as a boolean
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

func (b *Buffer) Float32(s uint) float32   { return math.Float32frombits(uint32(b.U32())) }
func (b *Buffer) Float32BE(s uint) float32 { return math.Float32frombits(uint32(b.U32BE())) }
func (b *Buffer) Float32LE(s uint) float32 { return math.Float32frombits(uint32(b.U32LE())) }

func (b *Buffer) Float64(s uint) float64   { return math.Float64frombits(uint64(b.U64())) }
func (b *Buffer) Float64BE(s uint) float64 { return math.Float64frombits(uint64(b.U64BE())) }
func (b *Buffer) Float64LE(s uint) float64 { return math.Float64frombits(uint64(b.U64LE())) }

func (b *Buffer) UTF8(nBytes uint64) (string, uint64) {
	// TODO: panic?
	s, rBits := b.BytesLen(nBytes)
	return string(s), rBits
}

func (b *Buffer) Unary(s uint) uint {
	var n uint
	for uint(b.U1()) == s {
		n++
	}
	return n
}
