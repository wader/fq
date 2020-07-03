package bitbuf

// TODO:
// inline for speed?
// F -> FLT?
// UTF16/UTF32

import (
	"errors"
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

var ErrUnexpectedEOF = errors.New("unexpected EOF")

// Buffer is a bitbuf buffer
type Buffer struct {
	// Len is bit length of buffer
	Len uint64
	// Pos is current bit position in buffer
	Pos uint64

	buf         []byte
	bufFirstBit uint64
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
func NewFromBitString(s string) (*Buffer, error) {
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
			return nil, fmt.Errorf("invalid bit string %q at index %d %q", s, i, c)
		}

		p := 8 - (i % 8) - 1
		n |= byte(b) << p
		if (i > 0 && p == 0) || i == len(s)-1 {
			buf = append(buf, n)
			n = 0
		}
	}

	return New(buf, 0, uint64(len(s))), nil
}

// BitBufRange reads nBits bits starting from start
// Does not update current position.
func (b *Buffer) BitBufRange(firstBit uint64, nBits uint64) (*Buffer, error) {
	endPos := uint64(firstBit) + uint64(nBits)
	if endPos > b.Len {
		return nil, ErrUnexpectedEOF
	}

	nb := &Buffer{
		buf:         b.buf,
		bufFirstBit: b.bufFirstBit + firstBit,
		Len:         nBits,
		Pos:         0,
	}

	return nb, nil
}

// BitBufLen reads nBits
func (b *Buffer) BitBufLen(nBits uint64) (*Buffer, error) {
	bb, err := b.BitBufRange(b.Pos, nBits)
	if err != nil {
		return nil, err
	}
	b.Pos += nBits
	return bb, nil
}

// Copy bitbuf
// TODO: rename? remove?
func (b *Buffer) Copy() *Buffer {
	return NewFromBitBuf(b)
}

// Bits reads nBits bits from buffer
func (b *Buffer) Bits(nBits uint64) (uint64, error) {
	p := uint64(b.Pos) + uint64(nBits)
	if p > b.Len {
		return 0, ErrUnexpectedEOF
	}
	n := ReadBits(b.buf, b.bufFirstBit+b.Pos, nBits)
	b.Pos += nBits

	return n, nil
}

// PeekBits peek nBits bits from buffer
// TODO: share code?
func (b *Buffer) PeekBits(nBits uint64) (uint64, error) {
	p := uint64(b.Pos) + uint64(nBits)
	if p > b.Len {
		return 0, ErrUnexpectedEOF
	}
	n := ReadBits(b.buf, b.bufFirstBit+b.Pos, nBits)

	return n, nil
}

// PeekBytes peek nBytes bytes from buffer
func (b *Buffer) PeekBytes(nBytes uint64) ([]byte, error) {
	bs, err := b.BytesRange(b.Pos, nBytes)
	if err != nil {
		return bs, nil
	}
	return bs, nil
}

// BytesRange reads nBytes bytes starting bit position start
// Does not update current position.
func (b *Buffer) BytesRange(firstBit uint64, nBytes uint64) ([]byte, error) {
	endPos := firstBit + nBytes*8
	if endPos > b.Len {
		return nil, ErrUnexpectedEOF
	}

	bufFirstBit := b.bufFirstBit + firstBit
	if bufFirstBit%8 == 0 {
		bufFirstBytePos := bufFirstBit >> 3
		nb := b.buf[bufFirstBytePos : bufFirstBytePos+nBytes]
		return nb, nil
	}

	var buf []byte
	for i := uint64(0); i < nBytes; i++ {
		buf = append(buf, byte(ReadBits(b.buf, bufFirstBit+i, 8)))
	}

	return buf, nil
}

// BytesLen reads nBytes bytes
func (b *Buffer) BytesLen(nBytes uint64) ([]byte, error) {
	bb, err := b.BytesRange(b.Pos, nBytes)
	if err != nil {
		return nil, err
	}
	b.Pos += nBytes * 8
	return bb, nil
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

// SeekRel seeks relative to current bit position
// TODO: better name?
func (b *Buffer) SeekRel(delta int64) (uint64, error) {
	endPos := uint64(int64(b.Pos) + delta)
	if endPos > b.Len {
		return b.Pos, ErrUnexpectedEOF
	}
	b.Pos = endPos

	return b.Pos, nil
}

// SeekAbs seeks to absolute position
func (b *Buffer) SeekAbs(pos uint64) (uint64, error) {
	// TODO: panic? bitbuf should never panic?
	if pos > b.Len {
		return b.Pos, ErrUnexpectedEOF
	}
	b.Pos = pos
	return b.Pos, nil
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

// UE reads a nBits bits unsigned integer with byte order endian
// MSB first
func (b *Buffer) UE(nBits uint64, endian Endian) (uint64, error) {
	n, err := b.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}

	return n, nil
}

// Bool reads one bit as a boolean
func (b *Buffer) Bool() (bool, error) {
	n, err := b.UE(1, BigEndian)
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (b *Buffer) U(nBits uint64) (uint64, error) { return b.UE(nBits, BigEndian) }
func (b *Buffer) U1() (uint64, error)            { return b.UE(1, BigEndian) }
func (b *Buffer) U2() (uint64, error)            { return b.UE(2, BigEndian) }
func (b *Buffer) U3() (uint64, error)            { return b.UE(3, BigEndian) }
func (b *Buffer) U4() (uint64, error)            { return b.UE(4, BigEndian) }
func (b *Buffer) U5() (uint64, error)            { return b.UE(5, BigEndian) }
func (b *Buffer) U6() (uint64, error)            { return b.UE(6, BigEndian) }
func (b *Buffer) U7() (uint64, error)            { return b.UE(7, BigEndian) }
func (b *Buffer) U8() (uint64, error)            { return b.UE(8, BigEndian) }
func (b *Buffer) U9() (uint64, error)            { return b.UE(9, BigEndian) }
func (b *Buffer) U10() (uint64, error)           { return b.UE(10, BigEndian) }
func (b *Buffer) U11() (uint64, error)           { return b.UE(11, BigEndian) }
func (b *Buffer) U12() (uint64, error)           { return b.UE(12, BigEndian) }
func (b *Buffer) U13() (uint64, error)           { return b.UE(13, BigEndian) }
func (b *Buffer) U14() (uint64, error)           { return b.UE(14, BigEndian) }
func (b *Buffer) U15() (uint64, error)           { return b.UE(15, BigEndian) }
func (b *Buffer) U16() (uint64, error)           { return b.UE(16, BigEndian) }
func (b *Buffer) U24() (uint64, error)           { return b.UE(24, BigEndian) }
func (b *Buffer) U32() (uint64, error)           { return b.UE(32, BigEndian) }
func (b *Buffer) U64() (uint64, error)           { return b.UE(64, BigEndian) }

func (b *Buffer) UBE(nBits uint64) (uint64, error) { return b.UE(nBits, BigEndian) }
func (b *Buffer) U9BE() (uint64, error)            { return b.UE(9, BigEndian) }
func (b *Buffer) U10BE() (uint64, error)           { return b.UE(10, BigEndian) }
func (b *Buffer) U11BE() (uint64, error)           { return b.UE(11, BigEndian) }
func (b *Buffer) U12BE() (uint64, error)           { return b.UE(12, BigEndian) }
func (b *Buffer) U13BE() (uint64, error)           { return b.UE(13, BigEndian) }
func (b *Buffer) U14BE() (uint64, error)           { return b.UE(14, BigEndian) }
func (b *Buffer) U15BE() (uint64, error)           { return b.UE(15, BigEndian) }
func (b *Buffer) U16BE() (uint64, error)           { return b.UE(16, BigEndian) }
func (b *Buffer) U24BE() (uint64, error)           { return b.UE(24, BigEndian) }
func (b *Buffer) U32BE() (uint64, error)           { return b.UE(32, BigEndian) }
func (b *Buffer) U64BE() (uint64, error)           { return b.UE(64, BigEndian) }

func (b *Buffer) ULE(nBits uint64) (uint64, error) { return b.UE(nBits, LittleEndian) }
func (b *Buffer) U9LE() (uint64, error)            { return b.UE(9, LittleEndian) }
func (b *Buffer) U10LE() (uint64, error)           { return b.UE(10, LittleEndian) }
func (b *Buffer) U11LE() (uint64, error)           { return b.UE(11, LittleEndian) }
func (b *Buffer) U12LE() (uint64, error)           { return b.UE(12, LittleEndian) }
func (b *Buffer) U13LE() (uint64, error)           { return b.UE(13, LittleEndian) }
func (b *Buffer) U14LE() (uint64, error)           { return b.UE(14, LittleEndian) }
func (b *Buffer) U15LE() (uint64, error)           { return b.UE(15, LittleEndian) }
func (b *Buffer) U16LE() (uint64, error)           { return b.UE(16, LittleEndian) }
func (b *Buffer) U24LE() (uint64, error)           { return b.UE(24, LittleEndian) }
func (b *Buffer) U32LE() (uint64, error)           { return b.UE(32, LittleEndian) }
func (b *Buffer) U64LE() (uint64, error)           { return b.UE(64, LittleEndian) }

// SE reads a nBits signed (two's-complement) integer with byte order endian
// MSB first
func (b *Buffer) SE(nBits uint64, endian Endian) (int64, error) {
	n, err := b.Bits(nBits)
	if err != nil {
		return 0, err
	}
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

	return s, nil
}

func (b *Buffer) S(nBits uint64) (int64, error) { return b.SE(nBits, BigEndian) }
func (b *Buffer) S1() (int64, error)            { return b.SE(1, BigEndian) }
func (b *Buffer) S2() (int64, error)            { return b.SE(2, BigEndian) }
func (b *Buffer) S3() (int64, error)            { return b.SE(3, BigEndian) }
func (b *Buffer) S4() (int64, error)            { return b.SE(4, BigEndian) }
func (b *Buffer) S5() (int64, error)            { return b.SE(5, BigEndian) }
func (b *Buffer) S6() (int64, error)            { return b.SE(6, BigEndian) }
func (b *Buffer) S7() (int64, error)            { return b.SE(7, BigEndian) }
func (b *Buffer) S8() (int64, error)            { return b.SE(8, BigEndian) }
func (b *Buffer) S9() (int64, error)            { return b.SE(9, BigEndian) }
func (b *Buffer) S10() (int64, error)           { return b.SE(10, BigEndian) }
func (b *Buffer) S11() (int64, error)           { return b.SE(11, BigEndian) }
func (b *Buffer) S12() (int64, error)           { return b.SE(12, BigEndian) }
func (b *Buffer) S13() (int64, error)           { return b.SE(13, BigEndian) }
func (b *Buffer) S14() (int64, error)           { return b.SE(14, BigEndian) }
func (b *Buffer) S15() (int64, error)           { return b.SE(15, BigEndian) }
func (b *Buffer) S16() (int64, error)           { return b.SE(16, BigEndian) }
func (b *Buffer) S24() (int64, error)           { return b.SE(24, BigEndian) }
func (b *Buffer) S32() (int64, error)           { return b.SE(32, BigEndian) }
func (b *Buffer) S64() (int64, error)           { return b.SE(64, BigEndian) }

func (b *Buffer) SBE(nBits uint64) (int64, error) { return b.SE(nBits, BigEndian) }
func (b *Buffer) S9BE() (int64, error)            { return b.SE(9, BigEndian) }
func (b *Buffer) S10BE() (int64, error)           { return b.SE(10, BigEndian) }
func (b *Buffer) S11BE() (int64, error)           { return b.SE(11, BigEndian) }
func (b *Buffer) S12BE() (int64, error)           { return b.SE(12, BigEndian) }
func (b *Buffer) S13BE() (int64, error)           { return b.SE(13, BigEndian) }
func (b *Buffer) S14BE() (int64, error)           { return b.SE(14, BigEndian) }
func (b *Buffer) S15BE() (int64, error)           { return b.SE(15, BigEndian) }
func (b *Buffer) S16BE() (int64, error)           { return b.SE(16, BigEndian) }
func (b *Buffer) S24BE() (int64, error)           { return b.SE(24, BigEndian) }
func (b *Buffer) S32BE() (int64, error)           { return b.SE(32, BigEndian) }
func (b *Buffer) S64BE() (int64, error)           { return b.SE(64, BigEndian) }

func (b *Buffer) SLE(nBits uint64) (int64, error) { return b.SE(nBits, LittleEndian) }
func (b *Buffer) S9LE() (int64, error)            { return b.SE(9, LittleEndian) }
func (b *Buffer) S10LE() (int64, error)           { return b.SE(10, LittleEndian) }
func (b *Buffer) S11LE() (int64, error)           { return b.SE(11, LittleEndian) }
func (b *Buffer) S12LE() (int64, error)           { return b.SE(12, LittleEndian) }
func (b *Buffer) S13LE() (int64, error)           { return b.SE(13, LittleEndian) }
func (b *Buffer) S14LE() (int64, error)           { return b.SE(14, LittleEndian) }
func (b *Buffer) S15LE() (int64, error)           { return b.SE(15, LittleEndian) }
func (b *Buffer) S16LE() (int64, error)           { return b.SE(16, LittleEndian) }
func (b *Buffer) S24LE() (int64, error)           { return b.SE(24, LittleEndian) }
func (b *Buffer) S32LE() (int64, error)           { return b.SE(32, LittleEndian) }
func (b *Buffer) S64LE() (int64, error)           { return b.SE(64, LittleEndian) }

func (b *Buffer) F32E(endian Endian) (float32, error) {
	n, err := b.Bits(32)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = ReverseBytes(32, n)
	}
	return math.Float32frombits(uint32(n)), nil
}
func (b *Buffer) F32(s uint) (float32, error)   { return b.F32E(BigEndian) }
func (b *Buffer) F32BE(s uint) (float32, error) { return b.F32E(BigEndian) }
func (b *Buffer) F32LE(s uint) (float32, error) { return b.F32E(LittleEndian) }

func (b *Buffer) F64E(endian Endian) (float64, error) {
	n, err := b.Bits(64)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = ReverseBytes(64, n)
	}
	return math.Float64frombits(n), nil
}
func (b *Buffer) F64(s uint) (float64, error)   { return b.F64E(BigEndian) }
func (b *Buffer) F64BE(s uint) (float64, error) { return b.F64E(BigEndian) }
func (b *Buffer) F64LE(s uint) (float64, error) { return b.F64E(LittleEndian) }

// TODO: FP64,unsigned/BE/LE? rename SFP32?

// FP64 signed fixed point 1:31:32
func (b *Buffer) FP64() (float64, error) {
	n, err := b.S64()
	if err != nil {
		return 0, err
	}
	return float64(float64(n) / (1 << 32)), nil
}

// FP32 signed fixed point 1:15:16
func (b *Buffer) FP32() (float64, error) {
	n, err := b.S32()
	if err != nil {
		return 0, err
	}
	return float64(float64(n) / (1 << 16)), nil
}

// FP16 signed fixed point 1:7:8
func (b *Buffer) FP16() (float64, error) {
	n, err := b.S16()
	if err != nil {
		return 0, err
	}
	return float64(float64(n) / (1 << 8)), nil
}

func (b *Buffer) UTF8(nBytes uint64) (string, error) {
	s, err := b.BytesLen(nBytes)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func (b *Buffer) Unary(s uint64) (uint64, error) {
	var n uint64
	for {
		b, err := b.U1()
		if err != nil {
			return 0, err
		}
		if b != s {
			break
		}
		n++
	}
	return n, nil
}
