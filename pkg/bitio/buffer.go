package bitio

// TODO:
// cache pos, len
// inline for speed?
// F -> FLT?
// UTF16/UTF32

import (
	"bytes"
	"errors"
	"fmt"
	"fq/internal/aheadreadseeker"
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

type progressReaderSeeker struct {
	RS         io.ReadSeeker
	Length     int64
	Pos        int64
	MaxPos     int64
	ProgressFn func(pos int64, length int64)
}

func (prs *progressReaderSeeker) Read(p []byte) (n int, err error) {
	n, err = prs.RS.Read(p)
	prs.Pos += int64(n)
	if prs.Pos > prs.MaxPos {
		prs.MaxPos = prs.Pos
		prs.ProgressFn(prs.MaxPos, prs.Length)
	}
	return n, err
}

func (prs *progressReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	pos, err := prs.RS.Seek(offset, whence)
	prs.Pos = pos
	return pos, err
}

// Buffer is a bit buffer
type Buffer struct {
	br interface {
		io.Reader // both Reader and SectionBitReader implement io.Reader
		BitReadSeeker
		BitReaderAt
	}
	bitLen int64 // mostly to cache len
}

// NewBufferFromReadSeeker new Buffer from io.ReadSeeker, start at firstBit with bit length lenBits
// buf is not copied.
func NewBufferFromReadSeeker(rs io.ReadSeeker) (*Buffer, error) {
	bPos, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	bEnd, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	if _, err := rs.Seek(bPos, io.SeekStart); err != nil {
		return nil, err
	}

	// // TODO: move
	// prs := &progressReaderSeeker{RS: rs, Length: bEnd, ProgressFn: func(pos, length int64) {
	// 	fmt.Fprintf(os.Stderr, "\r%.1f%%", float64(pos*100)/float64(length))
	// }}

	return &Buffer{
		br:     NewReaderFromReadSeeker(aheadreadseeker.New(rs, cacheReadAheadSize)),
		bitLen: bEnd * 8,
	}, nil
}

// NewBufferFromBytes new Buffer from bytes
// if nBits is < 0 nBits is all bits in buf
func NewBufferFromBytes(buf []byte, nBits int64) *Buffer {
	if nBits < 0 {
		nBits = int64(len(buf)) * 8
	}
	return &Buffer{
		br:     NewReaderFromReadSeeker(bytes.NewReader(buf)),
		bitLen: nBits,
	}
}

// NewBufferFromBitString new Buffer from bit string, ex: "0101"
func NewBufferFromBitString(s string) *Buffer {
	b, bBits := BytesFromBitString(s)
	return NewBufferFromBytes(b, int64(bBits))
}

// BitBufRange reads nBits bits starting from start
// Does not update current position.
// if nBits is < 0 nBits is all bits after firstBitOffset
func (b *Buffer) BitBufRange(firstBitOffset int64, nBits int64) (*Buffer, error) {
	// TODO: move error check?
	if firstBitOffset+nBits > b.bitLen {
		return nil, errors.New("outside buffer")
	}
	if nBits < 0 {
		nBits = b.bitLen - firstBitOffset
	}
	return &Buffer{
		br:     NewSectionBitReader(b.br, firstBitOffset, nBits),
		bitLen: nBits,
	}, nil
}

func (b *Buffer) Copy() *Buffer {
	return &Buffer{
		br:     NewSectionBitReader(b.br, 0, b.bitLen),
		bitLen: b.bitLen,
	}
}

func (b *Buffer) Pos() (int64, error) {
	bPos, err := b.br.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	return bPos, nil
}

func (b *Buffer) Len() int64 {
	return b.bitLen
}

// BitBufLen reads nBits
func (b *Buffer) BitBufLen(nBits int64) (*Buffer, error) {
	bPos, err := b.Pos()
	if err != nil {
		return nil, err
	}
	bb, err := b.BitBufRange(bPos, nBits)
	if err != nil {
		return nil, err
	}
	if _, err := b.br.SeekBits(nBits, io.SeekCurrent); err != nil {
		return nil, err
	}

	return bb, nil
}

// Bits reads nBits bits from buffer
func (b *Buffer) bits(nBits int) (uint64, error) {
	// 64 bits max, 9 byte worse case if not byte aligned
	var bufArray [9]byte
	buf := bufArray[:]
	_, err := b.br.ReadBits(buf[:], nBits)
	if err != nil {
		return 0, err
	}

	return Uint64(buf[:], 0, nBits), nil
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
	start, err := b.br.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	n, err := b.bits(nBits)
	if _, err := b.br.SeekBits(start, io.SeekStart); err != nil {
		return 0, err
	}
	return n, err
}

// PeekBytes peek nBytes bytes from buffer
func (b *Buffer) PeekBytes(nBytes int) ([]byte, error) {
	start, err := b.br.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	bs, err := b.BytesLen(nBytes)
	if _, err := b.br.SeekBits(start, io.SeekStart); err != nil {
		return nil, err
	}
	return bs, err
}

// TODO: will return maxLen*nBits if not found
func (b *Buffer) PeekFind(nBits int, v uint8, maxLen int64) (int64, error) {
	start, err := b.br.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	var count int64
	for {
		bv, err := b.U(nBits)
		if err != nil {
			if _, err := b.br.SeekBits(start, io.SeekStart); err != nil {
				return 0, err
			}
			return 0, err
		}
		count++
		if uint8(bv) == v || count == maxLen {
			break
		}
	}
	if _, err := b.br.SeekBits(start, io.SeekStart); err != nil {
		return 0, err
	}

	return count * int64(nBits), nil
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
// TODO: nBytes -1?
func (b *Buffer) BytesRange(bitOffset int64, nBytes int) ([]byte, error) {
	buf := make([]byte, nBytes)
	_, err := b.br.ReadBitsAt(buf, nBytes*8, bitOffset)
	return buf, err
}

// BytesLen reads nBytes bytes
func (b *Buffer) BytesLen(nBytes int) ([]byte, error) {
	buf := make([]byte, nBytes)
	_, err := io.ReadAtLeast(b, buf, nBytes)
	return buf, err
}

// End is true if current position is at the end
func (b *Buffer) End() (bool, error) {
	bPos, err := b.Pos()
	if err != nil {
		return false, err
	}
	return bPos >= b.bitLen, nil
}

// BitsLeft number of bits left until end
func (b *Buffer) BitsLeft() (int64, error) {
	bPos, err := b.Pos()
	if err != nil {
		return 0, err
	}
	return b.bitLen - bPos, nil
}

// ByteAlignBits number of bits to next byte align
func (b *Buffer) ByteAlignBits() (int, error) {
	bPos, err := b.Pos()
	if err != nil {
		return 0, err
	}
	return int((8 - (bPos & 0x7)) & 0x7), nil
}

// BytePos byte position of current bit position
func (b *Buffer) BytePos() (int64, error) {
	bPos, err := b.Pos()
	if err != nil {
		return 0, err
	}
	return bPos & 0x7, nil
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
	truncLen := b.bitLen
	truncS := ""
	if truncLen > 64 {
		truncLen, truncS = 64, "..."
	}
	truncBB, err := b.BitBufRange(0, truncLen)
	bitString, err := truncBB.BitString()
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("0b%s%s /* %d bits */", bitString, truncS, b.bitLen)
}

// BitString return bit string representation
func (b *Buffer) BitString() (string, error) {
	var ss []string
	for {
		n, err := b.Bits(1)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
		if n == 0 {
			ss = append(ss, "0")
		} else {
			ss = append(ss, "1")
		}
	}

	return strings.Join(ss, ""), nil
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

func (b *Buffer) FE(nBits int, endian Endian) (float64, error) {
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

func (b *Buffer) FPE(nBits int, fBits int64, endian Endian) (float64, error) {
	n, err := b.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}
	return float64(n) / float64(uint64(1<<fBits)), nil
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
