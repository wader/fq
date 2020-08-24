package bitbuf

// TODO:
// inline for speed?
// F -> FLT?
// UTF16/UTF32

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

const cacheReadAheadSize = 1024 * 1024

//rangeContain does a contain b
func rangeContain(aStart, aEnd, bStart, bEnd int64) bool {
	return bStart >= aStart && bStart <= aEnd && bEnd <= aEnd
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Endian byte order
type Endian int

const (
	// BigEndian byte order
	BigEndian Endian = iota
	// LittleEndian byte order
	LittleEndian
)

var ErrEOF = errors.New("EOF")

<<<<<<< HEAD
type cacheReader struct {
	rs           io.ReadSeeker
	Len          int64
	cacheBytePos int64
	cacheByteLen int64
	cache        []byte
}

func (b *cacheReader) read(buf []byte, bitPos int64, nBits int64) (int64, int64, error) {
	// log.Printf("bitPos=%d nBits=%d", bitPos, nBits)

	readBitPos := bitPos
	readBytePos := int64(readBitPos / 8)
	readSkipBits := readBitPos % 8
=======
func (b *Buffer) read(buf []byte, bitPos int64, nBits int64) (int64, error) {
	//log.Printf("bitPos=%d nBits=%d", bitPos, nBits)

	readBytePos := int64(bitPos / 8)
	readSkipBits := bitPos % 8
>>>>>>> readseeker
	readBits := readSkipBits + nBits
	readBytes := readBits / 8
	if readBits%8 > 0 {
		readBytes++
	}

<<<<<<< HEAD
	readByteEnd := readBytePos + readBytes
	for {
		cacheByteEnd := b.cacheBytePos + b.cacheByteLen
		if rangeContain(b.cacheBytePos, cacheByteEnd, readBytePos, readByteEnd) {
			offset := readBytePos - b.cacheBytePos
			copy(buf[0:readBytes], b.cache[offset:offset+readBytes])
			// log.Printf("cached buf[0:%d], b.cache[%d:%d]",
			// 	readBytes, offset, offset+readBytes)
			return readBytes, readSkipBits, nil
		}

		// log.Println("NOPE")

		// log.Printf("rangeContain(b.cacheBytePos %d, cacheByteEnd %d, readBytePos %d, readByteEnd %d)",
		// 	b.cacheBytePos, cacheByteEnd, readBytePos, readByteEnd)

		if _, err := b.rs.Seek(readBytePos, io.SeekStart); err != nil {
			return 0, 0, err
		}

		if readBytes > cacheReadAheadSize {
			if _, err := io.ReadFull(b.rs, buf[0:readBytes]); err != nil {
				return 0, 0, err
			}
			// log.Printf("to big %d", readBytes)
			return readBytes, readSkipBits, nil
		}

		// var cacheReadAheadKeep int64

		// if readBytePos == cacheByteEnd-1 && b.cacheByteLen != 0 {
		// cacheReadAheadKeep = min(int64(cacheReadAheadSize), b.cacheByteLen)
		// readAheadBytes = min(cacheReadAheadSize, maxReadBytes-readBytePos)
		// } else {
		// log.Printf("maxReadByte-readBytePoss: %#+v\n", maxReadBytes-readBytePos)
		// log.Printf("b.Len=%d", b.Len)
		readAheadBytes := min(cacheReadAheadSize*2, b.Len-readBytePos)
		// }

		// if cacheReadAheadKeep > 0 {
		// 	log.Printf("keep b.cache[0:%d], b.cache[%d:]",
		// 		cacheReadAheadKeep, b.cacheByteLen-cacheReadAheadKeep)
		// 	copy(b.cache[0:cacheReadAheadKeep], b.cache[b.cacheByteLen-cacheReadAheadKeep:])
		// }

		if _, err := io.ReadFull(b.rs, b.cache[0:readAheadBytes]); err != nil {
			return 0, 0, err
		}

		// log.Printf("read b.cache[%d:%d]",
		// 	0, readAheadBytes)

		b.cacheByteLen = readAheadBytes
		b.cacheBytePos = readBytePos
	}
=======
	if readBytes > int64(len(b.buf)) {
		b.buf = make([]byte, readBytes)
	}

	if _, err := b.crs.Seek(readBytePos, io.SeekStart); err != nil {
		return 0, err
	}

	// TODO: nBits should be available
	_, err := io.ReadFull(b.crs, b.buf[0:readBytes])
	if err != nil {
		return 0, err
	}

	// log.Printf("  n: %#+v\n", n)

	if readSkipBits == 0 && nBits%8 == 0 {
		// log.Println("  aligned")
		copy(buf[0:readBytes], b.buf[0:readBytes])
		return readBytes, nil
	}

	nBytes := nBits / 8
	restBits := nBits % 8

	// TODO: copy smartness if many bytes
	for i := int64(0); i < nBytes; i++ {
		buf[i] = byte(ReadBits(b.buf, readSkipBits+i*8, 8))
	}
	if restBits != 0 {
		buf[nBytes] = byte(ReadBits(b.buf, readSkipBits+nBytes*8, restBits)) << (8 - restBits)
	}

	// log.Printf("  readBytes: %#+v\n", readBytes)

	return readBytes, nil
>>>>>>> readseeker
}

// Buffer is a bitbuf buffer
type Buffer struct {
	// Len is bit length of buffer
	Len int64
	// Pos is current bit position in buffer
	Pos int64

	firstBitOffset int64

<<<<<<< HEAD
	cr *cacheReader
=======
	buf []byte

	crs *CachingReadSeeker
>>>>>>> readseeker
}

// NewFromReadSeeker bitbuf.Buffer from io.ReadSeeker, start at firstBit with bit length lenBits
// buf is not copied.
func NewFromReadSeeker(rs io.ReadSeeker, firstBitOffset int64) (*Buffer, error) {
	len, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	if _, err := rs.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	if firstBitOffset > len*8 {
		return nil, ErrEOF
	}

	return &Buffer{
		Len:            len*8 - firstBitOffset,
		Pos:            0,
		firstBitOffset: firstBitOffset,
<<<<<<< HEAD
		cr: &cacheReader{
			rs:           rs,
			Len:          len,
			cacheBytePos: 0,
			cache:        make([]byte, cacheReadAheadSize*2),
		},
=======
		crs:            NewCachingReadSeeker(rs, cacheReadAheadSize),
>>>>>>> readseeker
	}, nil
}

// NewFromBytes bitbuf.Buffer from bytes
func NewFromBytes(buf []byte, firstBitOffset int64) (*Buffer, error) {
	return NewFromReadSeeker(bytes.NewReader(buf), firstBitOffset)
}

// NewFromBitBuf bitbuf.Buffer from other bitbuf.Buffer
// Will be a shallow copy with position reset to zero.
func NewFromBitBuf(b *Buffer, firstBitOffset int64) (*Buffer, error) {
	if firstBitOffset > b.Len {
		return nil, ErrEOF
	}

	return &Buffer{
		Len:            b.Len - firstBitOffset,
		Pos:            0,
		firstBitOffset: b.firstBitOffset + firstBitOffset,
<<<<<<< HEAD
		cr:             b.cr,
=======
		crs:            b.crs,
>>>>>>> readseeker
	}, nil
}

// NewFromBitString bitbuf.Buffer from bit string, ex: "0101"
func NewFromBitString(s string) (*Buffer, error) {
	r := len(s) % 8
	bufLen := len(s) / 8
	if r > 0 {
		bufLen++
	}
	firstBifOffset := int64((8 - r) % 8)
	buf := make([]byte, bufLen)

	for i := 0; i < len(s); i++ {
		b := 0
		switch s[i] {
		case '0':
			b = 0
		case '1':
			b = 1
		default:
			return nil, fmt.Errorf("invalid bit string %q at index %d %q", s, i, s[i])
		}

		bufBitOffset := firstBifOffset + int64(i)
		bufByteOffset := bufBitOffset / 8
		buf[bufByteOffset] |= byte(b << (7 - bufBitOffset%8))
	}

	return NewFromBytes(buf, firstBifOffset)
}

// BitBufRange reads nBits bits starting from start
// Does not update current position.
func (b *Buffer) BitBufRange(firstBitOffset int64, nBits int64) (*Buffer, error) {
	endPos := firstBitOffset + nBits
	if endPos > b.Len {
		return nil, ErrEOF
	}

	nb := &Buffer{
		Len:            nBits,
		Pos:            0,
		firstBitOffset: b.firstBitOffset + firstBitOffset,
		crs:            b.crs,
	}

	return nb, nil
}

// BitBufLen reads nBits
func (b *Buffer) BitBufLen(nBits int64) (*Buffer, error) {
	bb, err := b.BitBufRange(b.Pos, nBits)
	if err != nil {
		return nil, err
	}
	b.Pos += nBits
	return bb, nil
}

// Copy bitbuf
// TODO: rename? remove?
// TODO: no error?
func (b *Buffer) Copy() (*Buffer, error) {
	return NewFromBitBuf(b, 0)
}

// Bits reads nBits bits from buffer
func (b *Buffer) bits(nBits int64) (uint64, error) {
	if b.Pos+nBits > b.Len {
		return 0, ErrEOF
	}

	var bufArray [10]byte
	buf := bufArray[:]
	_, err := b.read(buf[:], b.firstBitOffset+b.Pos, nBits)
	if err != nil {
		return 0, err
	}

	n := ReadBits(buf[:], 0, nBits)

	return n, nil
}

// Bits reads nBits bits from buffer
func (b *Buffer) Bits(nBits int64) (uint64, error) {
	n, err := b.bits(nBits)
	if err != nil {
		return 0, err
	}
	b.Pos += nBits
	return n, nil
}

// PeekBits peek nBits bits from buffer
// TODO: share code?
func (b *Buffer) PeekBits(nBits int64) (uint64, error) {
	return b.bits(nBits)
}

// PeekBytes peek nBytes bytes from buffer
func (b *Buffer) PeekBytes(nBytes int64) ([]byte, error) {
	bs, err := b.BytesRange(b.Pos, nBytes)
	if err != nil {
		return bs, nil
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

func (b *Buffer) ReadBits(buf []byte, bitOffset int64, nBits int64) error {
	if bitOffset+nBits > b.Len {
		return ErrEOF
	}

	_, err := b.read(buf, b.firstBitOffset+bitOffset, nBits)
	return err
}

func (b *Buffer) BytesBitRange(firstBitOffset int64, nBits int64, pad uint8) ([]byte, error) {
	if firstBitOffset+nBits > b.Len {
		return nil, ErrEOF
	}

	nBytes := nBits / 8
	if nBits%8 != 0 {
		nBytes++
	}

	buf := make([]byte, nBytes)
	_, err := b.read(buf, b.firstBitOffset+firstBitOffset, nBits)

	return buf, err
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	bitsLeft := b.Len - b.Pos
	bytesLeft := bitsLeft / 8
	if bitsLeft%8 != 0 {
		bytesLeft = 1
	}

	readBytes := int64(len(p))
	if readBytes > bytesLeft {
		readBytes = bytesLeft
	}

	if _, err := b.read(p, b.firstBitOffset+b.Pos, readBytes*8); err != nil {
		return 0, nil
	}

	b.Pos += readBytes * 8

	return int(readBytes), nil
}

// BytesRange reads nBytes bytes starting bit position start
// Does not update current position.
func (b *Buffer) BytesRange(firstBit int64, nBytes int64) ([]byte, error) {
	return b.BytesBitRange(firstBit, nBytes*8, 0)
}

// BytesLen reads nBytes bytes
func (b *Buffer) BytesLen(nBytes int64) ([]byte, error) {
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
func (b *Buffer) BitsLeft() int64 {
	return b.Len - b.Pos
}

// ByteAlignBits number of bits to next byte align
func (b *Buffer) ByteAlignBits() int64 {
	return (8 - (b.Pos & 0x7)) & 0x7
}

// BytePos byte position of current bit position
func (b *Buffer) BytePos() int64 {
	return b.Pos & 0x7
}

// SeekRel seeks relative to current bit position
// TODO: better name?
func (b *Buffer) SeekRel(delta int64) (int64, error) {
	endPos := b.Pos + delta
	if endPos > b.Len {
		return b.Pos, ErrEOF
	}
	b.Pos = endPos

	return b.Pos, nil
}

// SeekAbs seeks to absolute position
func (b *Buffer) SeekAbs(pos int64) (int64, error) {
	if pos > b.Len {
		return b.Pos, ErrEOF
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

// TruncateRel length of buffer to current position plus n
func (b *Buffer) TruncateRel(nBits int64) error {
	endPos := b.Pos + nBits
	if endPos > b.Len {
		return ErrEOF
	}

	b.Len = endPos

	return nil
}

// UE reads a nBits bits unsigned integer with byte order endian
// MSB first
func (b *Buffer) UE(nBits int64, endian Endian) (uint64, error) {
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

func (b *Buffer) U(nBits int64) (uint64, error) { return b.UE(nBits, BigEndian) }
func (b *Buffer) U1() (uint64, error)           { return b.UE(1, BigEndian) }
func (b *Buffer) U2() (uint64, error)           { return b.UE(2, BigEndian) }
func (b *Buffer) U3() (uint64, error)           { return b.UE(3, BigEndian) }
func (b *Buffer) U4() (uint64, error)           { return b.UE(4, BigEndian) }
func (b *Buffer) U5() (uint64, error)           { return b.UE(5, BigEndian) }
func (b *Buffer) U6() (uint64, error)           { return b.UE(6, BigEndian) }
func (b *Buffer) U7() (uint64, error)           { return b.UE(7, BigEndian) }
func (b *Buffer) U8() (uint64, error)           { return b.UE(8, BigEndian) }
func (b *Buffer) U9() (uint64, error)           { return b.UE(9, BigEndian) }
func (b *Buffer) U10() (uint64, error)          { return b.UE(10, BigEndian) }
func (b *Buffer) U11() (uint64, error)          { return b.UE(11, BigEndian) }
func (b *Buffer) U12() (uint64, error)          { return b.UE(12, BigEndian) }
func (b *Buffer) U13() (uint64, error)          { return b.UE(13, BigEndian) }
func (b *Buffer) U14() (uint64, error)          { return b.UE(14, BigEndian) }
func (b *Buffer) U15() (uint64, error)          { return b.UE(15, BigEndian) }
func (b *Buffer) U16() (uint64, error)          { return b.UE(16, BigEndian) }
func (b *Buffer) U24() (uint64, error)          { return b.UE(24, BigEndian) }
func (b *Buffer) U32() (uint64, error)          { return b.UE(32, BigEndian) }
func (b *Buffer) U64() (uint64, error)          { return b.UE(64, BigEndian) }

func (b *Buffer) UBE(nBits int64) (uint64, error) { return b.UE(nBits, BigEndian) }
func (b *Buffer) U9BE() (uint64, error)           { return b.UE(9, BigEndian) }
func (b *Buffer) U10BE() (uint64, error)          { return b.UE(10, BigEndian) }
func (b *Buffer) U11BE() (uint64, error)          { return b.UE(11, BigEndian) }
func (b *Buffer) U12BE() (uint64, error)          { return b.UE(12, BigEndian) }
func (b *Buffer) U13BE() (uint64, error)          { return b.UE(13, BigEndian) }
func (b *Buffer) U14BE() (uint64, error)          { return b.UE(14, BigEndian) }
func (b *Buffer) U15BE() (uint64, error)          { return b.UE(15, BigEndian) }
func (b *Buffer) U16BE() (uint64, error)          { return b.UE(16, BigEndian) }
func (b *Buffer) U24BE() (uint64, error)          { return b.UE(24, BigEndian) }
func (b *Buffer) U32BE() (uint64, error)          { return b.UE(32, BigEndian) }
func (b *Buffer) U64BE() (uint64, error)          { return b.UE(64, BigEndian) }

func (b *Buffer) ULE(nBits int64) (uint64, error) { return b.UE(nBits, LittleEndian) }
func (b *Buffer) U9LE() (uint64, error)           { return b.UE(9, LittleEndian) }
func (b *Buffer) U10LE() (uint64, error)          { return b.UE(10, LittleEndian) }
func (b *Buffer) U11LE() (uint64, error)          { return b.UE(11, LittleEndian) }
func (b *Buffer) U12LE() (uint64, error)          { return b.UE(12, LittleEndian) }
func (b *Buffer) U13LE() (uint64, error)          { return b.UE(13, LittleEndian) }
func (b *Buffer) U14LE() (uint64, error)          { return b.UE(14, LittleEndian) }
func (b *Buffer) U15LE() (uint64, error)          { return b.UE(15, LittleEndian) }
func (b *Buffer) U16LE() (uint64, error)          { return b.UE(16, LittleEndian) }
func (b *Buffer) U24LE() (uint64, error)          { return b.UE(24, LittleEndian) }
func (b *Buffer) U32LE() (uint64, error)          { return b.UE(32, LittleEndian) }
func (b *Buffer) U64LE() (uint64, error)          { return b.UE(64, LittleEndian) }

// SE reads a nBits signed (two's-complement) integer with byte order endian
// MSB first
func (b *Buffer) SE(nBits int64, endian Endian) (int64, error) {
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

func (b *Buffer) S(nBits int64) (int64, error) { return b.SE(nBits, BigEndian) }
func (b *Buffer) S1() (int64, error)           { return b.SE(1, BigEndian) }
func (b *Buffer) S2() (int64, error)           { return b.SE(2, BigEndian) }
func (b *Buffer) S3() (int64, error)           { return b.SE(3, BigEndian) }
func (b *Buffer) S4() (int64, error)           { return b.SE(4, BigEndian) }
func (b *Buffer) S5() (int64, error)           { return b.SE(5, BigEndian) }
func (b *Buffer) S6() (int64, error)           { return b.SE(6, BigEndian) }
func (b *Buffer) S7() (int64, error)           { return b.SE(7, BigEndian) }
func (b *Buffer) S8() (int64, error)           { return b.SE(8, BigEndian) }
func (b *Buffer) S9() (int64, error)           { return b.SE(9, BigEndian) }
func (b *Buffer) S10() (int64, error)          { return b.SE(10, BigEndian) }
func (b *Buffer) S11() (int64, error)          { return b.SE(11, BigEndian) }
func (b *Buffer) S12() (int64, error)          { return b.SE(12, BigEndian) }
func (b *Buffer) S13() (int64, error)          { return b.SE(13, BigEndian) }
func (b *Buffer) S14() (int64, error)          { return b.SE(14, BigEndian) }
func (b *Buffer) S15() (int64, error)          { return b.SE(15, BigEndian) }
func (b *Buffer) S16() (int64, error)          { return b.SE(16, BigEndian) }
func (b *Buffer) S24() (int64, error)          { return b.SE(24, BigEndian) }
func (b *Buffer) S32() (int64, error)          { return b.SE(32, BigEndian) }
func (b *Buffer) S64() (int64, error)          { return b.SE(64, BigEndian) }

func (b *Buffer) SBE(nBits int64) (int64, error) { return b.SE(nBits, BigEndian) }
func (b *Buffer) S9BE() (int64, error)           { return b.SE(9, BigEndian) }
func (b *Buffer) S10BE() (int64, error)          { return b.SE(10, BigEndian) }
func (b *Buffer) S11BE() (int64, error)          { return b.SE(11, BigEndian) }
func (b *Buffer) S12BE() (int64, error)          { return b.SE(12, BigEndian) }
func (b *Buffer) S13BE() (int64, error)          { return b.SE(13, BigEndian) }
func (b *Buffer) S14BE() (int64, error)          { return b.SE(14, BigEndian) }
func (b *Buffer) S15BE() (int64, error)          { return b.SE(15, BigEndian) }
func (b *Buffer) S16BE() (int64, error)          { return b.SE(16, BigEndian) }
func (b *Buffer) S24BE() (int64, error)          { return b.SE(24, BigEndian) }
func (b *Buffer) S32BE() (int64, error)          { return b.SE(32, BigEndian) }
func (b *Buffer) S64BE() (int64, error)          { return b.SE(64, BigEndian) }

func (b *Buffer) SLE(nBits int64) (int64, error) { return b.SE(nBits, LittleEndian) }
func (b *Buffer) S9LE() (int64, error)           { return b.SE(9, LittleEndian) }
func (b *Buffer) S10LE() (int64, error)          { return b.SE(10, LittleEndian) }
func (b *Buffer) S11LE() (int64, error)          { return b.SE(11, LittleEndian) }
func (b *Buffer) S12LE() (int64, error)          { return b.SE(12, LittleEndian) }
func (b *Buffer) S13LE() (int64, error)          { return b.SE(13, LittleEndian) }
func (b *Buffer) S14LE() (int64, error)          { return b.SE(14, LittleEndian) }
func (b *Buffer) S15LE() (int64, error)          { return b.SE(15, LittleEndian) }
func (b *Buffer) S16LE() (int64, error)          { return b.SE(16, LittleEndian) }
func (b *Buffer) S24LE() (int64, error)          { return b.SE(24, LittleEndian) }
func (b *Buffer) S32LE() (int64, error)          { return b.SE(32, LittleEndian) }
func (b *Buffer) S64LE() (int64, error)          { return b.SE(64, LittleEndian) }

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

func (b *Buffer) UTF8(nBytes int64) (string, error) {
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
