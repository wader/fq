package bitio

// not concurrency safe as bitsBuf is reused

// TODO:
// cache pos, len
// inline for speed?
// F -> FLT?
// UTF16/UTF32

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Buffer is a bit buffer
type Buffer struct {
	br interface {
		io.Reader // both Reader and SectionBitReader implement io.Reader
		BitReadSeeker
		BitReader
		BitReaderAt
	}

	bitLen int64 // mostly to cache len
	ctx    context.Context
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

	return &Buffer{
		br:     NewReaderFromReadSeeker(rs),
		bitLen: bEnd * 8,
	}, nil
}

func NewBufferFromBitReadSeeker(br interface {
	io.Reader
	BitReadSeeker
	BitReaderAt
}) (*Buffer, error) {
	bPos, err := br.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	bEnd, err := br.SeekBits(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	if _, err := br.SeekBits(bPos, io.SeekStart); err != nil {
		return nil, err
	}

	return &Buffer{
		br:     br,
		bitLen: bEnd,
	}, nil
}

// NewBufferFromBytes new Buffer from bytes
// if nBits is < 0 nBits is all bits in buf
func NewBufferFromBytes(buf []byte, nBits int64) *Buffer {
	if nBits < 0 {
		nBits = int64(len(buf)) * 8
	}
	return &Buffer{
		br:     NewSectionBitReader(NewReaderFromReadSeeker(bytes.NewReader(buf)), 0, nBits),
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
	if nBits < 0 {
		return nil, errors.New("negative nBits")
	}
	if firstBitOffset+nBits > b.bitLen {
		return nil, errors.New("outside buffer")
	}
	if nBits < 0 {
		nBits = b.bitLen - firstBitOffset
	}
	return &Buffer{
		br:     NewSectionBitReader(b.br, firstBitOffset, nBits),
		bitLen: nBits,
		ctx:    b.ctx,
	}, nil
}

func (b *Buffer) Copy() *Buffer {
	return b.CopyWithContext(b.ctx)
}

func (b *Buffer) CopyWithContext(ctx context.Context) *Buffer {
	return &Buffer{
		br:     NewSectionBitReader(b.br, 0, b.bitLen),
		bitLen: b.bitLen,
		ctx:    ctx,
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

func (b *Buffer) ReadBits(p []byte, nBits int) (n int, err error) {
	return b.br.ReadBits(p, nBits)
}

func (b *Buffer) ReadBitsAt(p []byte, nBits int, bitOff int64) (n int, err error) {
	return b.br.ReadBitsAt(p, nBits, bitOff)
}

func (b *Buffer) SeekBits(bitOffset int64, whence int) (int64, error) {
	return b.br.SeekBits(bitOffset, whence)
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
	n, err := ReadAtFull(b.br, buf, nBytes*8, bitOffset)
	if n == nBytes*8 {
		err = nil
	}
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
	if err != nil {
		return err.Error()
	}
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
		var buf [1]byte
		_, err := ReadFull(b, buf[:], 1)
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
		if buf[0] != 0b1000_0000 {
			ss = append(ss, "0")
		} else {
			ss = append(ss, "1")
		}
	}

	return strings.Join(ss, ""), nil
}
