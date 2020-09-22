package bitbuf

// TODO:
// inline for speed?
// F -> FLT?
// UTF16/UTF32

import (
	"bytes"
	"fmt"
	"fq/internal/aheadreadseeker"
	"fq/internal/bitio"
	"io"
	"math"
	"strings"
)

const cacheReadAheadSize = 256 * 1024

// Endian byte order
type Endian int

const (
	// BigEndian byte order
	BigEndian Endian = iota
	// LittleEndian byte order
	LittleEndian
)

// Buffer is a bitbuf buffer
type Buffer struct {
	br interface {
		io.ReadSeeker
		bitio.BitReadSeeker
		bitio.BitReaderAt
	}
}

// NewFromReadSeeker bitbuf.Buffer from io.ReadSeeker, start at firstBit with bit length lenBits
// buf is not copied.
func NewFromReadSeeker(rs io.ReadSeeker) *Buffer {
	return &Buffer{
		br: bitio.NewFromReadSeeker(aheadreadseeker.New(rs, cacheReadAheadSize)),
	}
}

// NewFromBytes bitbuf.Buffer from bytes
func NewFromBytes(buf []byte, nBits int64) *Buffer {
	return &Buffer{
		br: bitio.NewSectionBitReader(bitio.NewFromReadSeeker(bytes.NewReader(buf)), 0, nBits),
	}
}

// NewFromBitString bitbuf.Buffer from bit string, ex: "0101"
func NewFromBitString(s string) *Buffer {
	b, bBits := bitio.BytesFromBitString(s)
	return NewFromBytes(b, int64(bBits))
}

// BitBufRange reads nBits bits starting from start
// Does not update current position.
func (b *Buffer) BitBufRange(firstBitOffset int64, nBits int64) *Buffer {
	return &Buffer{
		br: bitio.NewSectionBitReader(b.br, firstBitOffset, nBits),
	}
}

func (b *Buffer) Pos() int64 {
	pos, err := b.br.SeekBits(0, io.SeekCurrent)
	if err != nil {
		panic("pos seek failed")
	}
	return pos
}

func (b *Buffer) Len() int64 {
	pos := b.Pos()
	end, err := b.br.SeekBits(0, io.SeekEnd)
	if err != nil {
		panic("end seek failed")
	}
	if _, err := b.br.SeekBits(pos, io.SeekStart); err != nil {
		panic("len restore seek failed")
	}
	return end
}

// BitBufLen reads nBits
func (b *Buffer) BitBufLen(nBits int64) (*Buffer, error) {
	bb := b.BitBufRange(b.Pos(), nBits)
	if _, err := b.br.SeekBits(nBits, io.SeekCurrent); err != nil {
		return nil, err
	}

	return bb, nil
}

// Bits reads nBits bits from buffer
func (b *Buffer) bits(nBits int) (uint64, error) {
	var bufArray [10]byte
	buf := bufArray[:]
	_, err := b.br.ReadBits(buf[:], nBits)
	if err != nil {
		return 0, err
	}

	return bitio.Uint64(buf[:], 0, nBits), nil
}

// Bits reads nBits bits from buffer
func (b *Buffer) Bits(nBits int) (uint64, error) {
	n, err := b.bits(nBits)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// PeekBits peek nBits bits from buffer
// TODO: share code?
func (b *Buffer) PeekBits(nBits int) (uint64, error) {
	n, err := b.bits(nBits)
	if err == nil || err == io.EOF {
		_, err = b.br.SeekBits(-int64(nBits), io.SeekCurrent)
	}
	return n, err
}

// PeekBytes peek nBytes bytes from buffer
func (b *Buffer) PeekBytes(nBytes int) ([]byte, error) {
	bs, err := b.BytesLen(nBytes)
	if err == nil || err == io.EOF {
		_, err = b.br.SeekBits(-int64(nBytes)*8, io.SeekCurrent)
	}
	return bs, nil
}

func (b *Buffer) PeekFind(nBits int64, v uint8, maxLen int64) (int64, error) {
	var count int64
	for {
		bv, err := b.U(nBits)
		if err != nil {
			return 0, err
		}
		count++
		if uint8(bv) == v || count == maxLen {
			break
		}
	}
	_, err := b.SeekRel(-count * int64(nBits))
	if err != nil {
		return 0, err
	}

	return count * nBits, nil
}

func (b *Buffer) ReadBits(buf []byte, bitOffset int64, nBits int) error {
	_, err := b.br.ReadBitsAt(buf, nBits, bitOffset)
	return err
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	return b.br.Read(p)
}

// BytesRange reads nBytes bytes starting bit position start
// Does not update current position.
// TODO: swap args
func (b *Buffer) BytesRange(bitOffset int64, nBytes int) ([]byte, error) {
	buf := make([]byte, nBytes)
	_, err := b.br.ReadBitsAt(buf, nBytes, bitOffset)
	return buf, err
}

// BytesLen reads nBytes bytes
func (b *Buffer) BytesLen(nBytes int) ([]byte, error) {
	buf := make([]byte, nBytes)
	_, err := io.ReadFull(b.BitBufRange(b.Pos(), int64(nBytes)*8), buf)
	return buf, err
}

// End is true if current position if at the end
func (b *Buffer) End() bool {
	return b.Pos() >= b.Len()
}

// BitsLeft number of bits left until end
func (b *Buffer) BitsLeft() int64 {
	return b.Len() - b.Pos()
}

// ByteAlignBits number of bits to next byte align
func (b *Buffer) ByteAlignBits() int64 {
	return (8 - (b.Pos() & 0x7)) & 0x7
}

// BytePos byte position of current bit position
func (b *Buffer) BytePos() int64 {
	return b.Pos() & 0x7
}

// SeekRel seeks relative to current bit position
// TODO: better name?
func (b *Buffer) SeekRel(delta int64) (int64, error) {
	return b.br.SeekBits(delta, io.SeekCurrent)
}

// SeekAbs seeks to absolute position
func (b *Buffer) SeekAbs(pos int64) (int64, error) {
	return b.br.SeekBits(pos, io.SeekStart)
}

func (b *Buffer) String() string {
	truncLen, truncS := b.Len(), ""
	if truncLen > 64 {
		truncLen, truncS = 64, "..."
	}
	truncBB := b.BitBufRange(0, truncLen)

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

// TruncateRel length of buffer to current position plus n
func (b *Buffer) TruncateRel(nBits int64) error {
	endPos := b.Pos + nBits
	if endPos > b.Len {
		return io.ErrUnexpectedEOF
	}

	b.Len = endPos

	return nil
}

// UE reads a nBits bits unsigned integer with byte order endian
// MSB first
func (b *Buffer) UE(nBits int, endian Endian) (uint64, error) {
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

func (b *Buffer) U(nBits int) (uint64, error) { return b.UE(nBits, BigEndian) }
func (b *Buffer) U1() (uint64, error)         { return b.UE(1, BigEndian) }
func (b *Buffer) U2() (uint64, error)         { return b.UE(2, BigEndian) }
func (b *Buffer) U3() (uint64, error)         { return b.UE(3, BigEndian) }
func (b *Buffer) U4() (uint64, error)         { return b.UE(4, BigEndian) }
func (b *Buffer) U5() (uint64, error)         { return b.UE(5, BigEndian) }
func (b *Buffer) U6() (uint64, error)         { return b.UE(6, BigEndian) }
func (b *Buffer) U7() (uint64, error)         { return b.UE(7, BigEndian) }
func (b *Buffer) U8() (uint64, error)         { return b.UE(8, BigEndian) }
func (b *Buffer) U9() (uint64, error)         { return b.UE(9, BigEndian) }
func (b *Buffer) U10() (uint64, error)        { return b.UE(10, BigEndian) }
func (b *Buffer) U11() (uint64, error)        { return b.UE(11, BigEndian) }
func (b *Buffer) U12() (uint64, error)        { return b.UE(12, BigEndian) }
func (b *Buffer) U13() (uint64, error)        { return b.UE(13, BigEndian) }
func (b *Buffer) U14() (uint64, error)        { return b.UE(14, BigEndian) }
func (b *Buffer) U15() (uint64, error)        { return b.UE(15, BigEndian) }
func (b *Buffer) U16() (uint64, error)        { return b.UE(16, BigEndian) }
func (b *Buffer) U24() (uint64, error)        { return b.UE(24, BigEndian) }
func (b *Buffer) U32() (uint64, error)        { return b.UE(32, BigEndian) }
func (b *Buffer) U64() (uint64, error)        { return b.UE(64, BigEndian) }

func (b *Buffer) UBE(nBits int) (uint64, error) { return b.UE(nBits, BigEndian) }
func (b *Buffer) U9BE() (uint64, error)         { return b.UE(9, BigEndian) }
func (b *Buffer) U10BE() (uint64, error)        { return b.UE(10, BigEndian) }
func (b *Buffer) U11BE() (uint64, error)        { return b.UE(11, BigEndian) }
func (b *Buffer) U12BE() (uint64, error)        { return b.UE(12, BigEndian) }
func (b *Buffer) U13BE() (uint64, error)        { return b.UE(13, BigEndian) }
func (b *Buffer) U14BE() (uint64, error)        { return b.UE(14, BigEndian) }
func (b *Buffer) U15BE() (uint64, error)        { return b.UE(15, BigEndian) }
func (b *Buffer) U16BE() (uint64, error)        { return b.UE(16, BigEndian) }
func (b *Buffer) U24BE() (uint64, error)        { return b.UE(24, BigEndian) }
func (b *Buffer) U32BE() (uint64, error)        { return b.UE(32, BigEndian) }
func (b *Buffer) U64BE() (uint64, error)        { return b.UE(64, BigEndian) }

func (b *Buffer) ULE(nBits int) (uint64, error) { return b.UE(nBits, LittleEndian) }
func (b *Buffer) U9LE() (uint64, error)         { return b.UE(9, LittleEndian) }
func (b *Buffer) U10LE() (uint64, error)        { return b.UE(10, LittleEndian) }
func (b *Buffer) U11LE() (uint64, error)        { return b.UE(11, LittleEndian) }
func (b *Buffer) U12LE() (uint64, error)        { return b.UE(12, LittleEndian) }
func (b *Buffer) U13LE() (uint64, error)        { return b.UE(13, LittleEndian) }
func (b *Buffer) U14LE() (uint64, error)        { return b.UE(14, LittleEndian) }
func (b *Buffer) U15LE() (uint64, error)        { return b.UE(15, LittleEndian) }
func (b *Buffer) U16LE() (uint64, error)        { return b.UE(16, LittleEndian) }
func (b *Buffer) U24LE() (uint64, error)        { return b.UE(24, LittleEndian) }
func (b *Buffer) U32LE() (uint64, error)        { return b.UE(32, LittleEndian) }
func (b *Buffer) U64LE() (uint64, error)        { return b.UE(64, LittleEndian) }

// SE reads a nBits signed (two's-complement) integer with byte order endian
// MSB first
func (b *Buffer) SE(nBits int, endian Endian) (int64, error) {
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

func (b *Buffer) S(nBits int) (int64, error) { return b.SE(nBits, BigEndian) }
func (b *Buffer) S1() (int64, error)         { return b.SE(1, BigEndian) }
func (b *Buffer) S2() (int64, error)         { return b.SE(2, BigEndian) }
func (b *Buffer) S3() (int64, error)         { return b.SE(3, BigEndian) }
func (b *Buffer) S4() (int64, error)         { return b.SE(4, BigEndian) }
func (b *Buffer) S5() (int64, error)         { return b.SE(5, BigEndian) }
func (b *Buffer) S6() (int64, error)         { return b.SE(6, BigEndian) }
func (b *Buffer) S7() (int64, error)         { return b.SE(7, BigEndian) }
func (b *Buffer) S8() (int64, error)         { return b.SE(8, BigEndian) }
func (b *Buffer) S9() (int64, error)         { return b.SE(9, BigEndian) }
func (b *Buffer) S10() (int64, error)        { return b.SE(10, BigEndian) }
func (b *Buffer) S11() (int64, error)        { return b.SE(11, BigEndian) }
func (b *Buffer) S12() (int64, error)        { return b.SE(12, BigEndian) }
func (b *Buffer) S13() (int64, error)        { return b.SE(13, BigEndian) }
func (b *Buffer) S14() (int64, error)        { return b.SE(14, BigEndian) }
func (b *Buffer) S15() (int64, error)        { return b.SE(15, BigEndian) }
func (b *Buffer) S16() (int64, error)        { return b.SE(16, BigEndian) }
func (b *Buffer) S24() (int64, error)        { return b.SE(24, BigEndian) }
func (b *Buffer) S32() (int64, error)        { return b.SE(32, BigEndian) }
func (b *Buffer) S64() (int64, error)        { return b.SE(64, BigEndian) }

func (b *Buffer) SBE(nBits int) (int64, error) { return b.SE(nBits, BigEndian) }
func (b *Buffer) S9BE() (int64, error)         { return b.SE(9, BigEndian) }
func (b *Buffer) S10BE() (int64, error)        { return b.SE(10, BigEndian) }
func (b *Buffer) S11BE() (int64, error)        { return b.SE(11, BigEndian) }
func (b *Buffer) S12BE() (int64, error)        { return b.SE(12, BigEndian) }
func (b *Buffer) S13BE() (int64, error)        { return b.SE(13, BigEndian) }
func (b *Buffer) S14BE() (int64, error)        { return b.SE(14, BigEndian) }
func (b *Buffer) S15BE() (int64, error)        { return b.SE(15, BigEndian) }
func (b *Buffer) S16BE() (int64, error)        { return b.SE(16, BigEndian) }
func (b *Buffer) S24BE() (int64, error)        { return b.SE(24, BigEndian) }
func (b *Buffer) S32BE() (int64, error)        { return b.SE(32, BigEndian) }
func (b *Buffer) S64BE() (int64, error)        { return b.SE(64, BigEndian) }

func (b *Buffer) SLE(nBits int) (int64, error) { return b.SE(nBits, LittleEndian) }
func (b *Buffer) S9LE() (int64, error)         { return b.SE(9, LittleEndian) }
func (b *Buffer) S10LE() (int64, error)        { return b.SE(10, LittleEndian) }
func (b *Buffer) S11LE() (int64, error)        { return b.SE(11, LittleEndian) }
func (b *Buffer) S12LE() (int64, error)        { return b.SE(12, LittleEndian) }
func (b *Buffer) S13LE() (int64, error)        { return b.SE(13, LittleEndian) }
func (b *Buffer) S14LE() (int64, error)        { return b.SE(14, LittleEndian) }
func (b *Buffer) S15LE() (int64, error)        { return b.SE(15, LittleEndian) }
func (b *Buffer) S16LE() (int64, error)        { return b.SE(16, LittleEndian) }
func (b *Buffer) S24LE() (int64, error)        { return b.SE(24, LittleEndian) }
func (b *Buffer) S32LE() (int64, error)        { return b.SE(32, LittleEndian) }
func (b *Buffer) S64LE() (int64, error)        { return b.SE(64, LittleEndian) }

func (b *Buffer) FE(nBits int64, endian Endian) (float64, error) {
	n, err := b.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}
	switch nBits {
	case 32:
		return math.Float64frombits(n), nil
	case 64:
		return float64(math.Float32frombits(uint32(n))), nil
	default:
		return 0, fmt.Errorf("unsupported float size %d", nBits)
	}
}

func (b *Buffer) F32E(endian Endian) (float64, error) { return b.FE(32, endian) }
func (b *Buffer) F32() (float64, error)               { return b.FE(32, BigEndian) }
func (b *Buffer) F32BE() (float64, error)             { return b.FE(32, BigEndian) }
func (b *Buffer) F32LE() (float64, error)             { return b.FE(32, LittleEndian) }

func (b *Buffer) F64E(endian Endian) (float64, error) { return b.FE(64, endian) }
func (b *Buffer) F64() (float64, error)               { return b.F64E(BigEndian) }
func (b *Buffer) F64BE() (float64, error)             { return b.F64E(BigEndian) }
func (b *Buffer) F64LE() (float64, error)             { return b.F64E(LittleEndian) }

// TODO: FP64,unsigned/BE/LE? rename SFP32?

func (b *Buffer) FPE(nBits int64, dBits int64, endian Endian) (float64, error) {
	n, err := b.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}
	return float64(n) / float64(uint64(1<<dBits)), nil
}

// FP64 signed fixed point 1:31:32
func (b *Buffer) FP64E(endian Endian) (float64, error) { return b.FPE(64, 32, endian) }
func (b *Buffer) FP64() (float64, error)               { return b.FPE(64, 32, BigEndian) }
func (b *Buffer) FP64BE() (float64, error)             { return b.FPE(64, 32, BigEndian) }
func (b *Buffer) FP64LE() (float64, error)             { return b.FPE(64, 32, LittleEndian) }

// FP32 signed fixed point 1:15:16
func (b *Buffer) FP32E(endian Endian) (float64, error) { return b.FPE(32, 16, endian) }
func (b *Buffer) FP32() (float64, error)               { return b.FPE(32, 16, BigEndian) }
func (b *Buffer) FP32BE() (float64, error)             { return b.FPE(32, 16, BigEndian) }
func (b *Buffer) FP32LE() (float64, error)             { return b.FPE(32, 16, LittleEndian) }

// FP16 signed fixed point 1:15:16
func (b *Buffer) FP16E(endian Endian) (float64, error) { return b.FPE(16, 8, endian) }
func (b *Buffer) FP16() (float64, error)               { return b.FPE(16, 8, BigEndian) }
func (b *Buffer) FP16BE() (float64, error)             { return b.FPE(16, 8, BigEndian) }
func (b *Buffer) FP16LE() (float64, error)             { return b.FPE(16, 8, LittleEndian) }

func (b *Buffer) UTF8(nBytes int) (string, error) {
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
