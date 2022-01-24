package bitio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrOffset = errors.New("invalid seek offset")
var ErrNegativeNBits = errors.New("negative number of bits")

type ReaderAt interface {
	ReadBitsAt(p []byte, nBits int64, bitOff int64) (n int64, err error)
}

type Reader interface {
	ReadBits(p []byte, nBits int64) (n int64, err error)
}

type Seeker interface {
	SeekBits(bitOffset int64, whence int) (int64, error)
}

type ReadSeeker interface {
	Reader
	Seeker
}

type ReadAtSeeker interface {
	ReaderAt
	Seeker
}

type ReaderAtSeeker interface {
	Reader
	ReaderAt
	Seeker
}

type Writer interface {
	WriteBits(p []byte, nBits int64) (n int64, err error)
}

// NewBitReader reading nBits bits from a []byte
// If nBits is -1 all bits will be used.
// Similar to bytes.NewReader
func NewBitReader(buf []byte, nBits int64) *SectionReader {
	if nBits < 0 {
		nBits = int64(len(buf)) * 8
	}
	return NewSectionReader(
		NewIOBitReadSeeker(bytes.NewReader(buf)),
		0,
		nBits,
	)
}

// BitsByteCount returns smallest amount of bytes to fit nBits bits
func BitsByteCount(nBits int64) int64 {
	n := nBits / 8
	if nBits%8 != 0 {
		n++
	}
	return n
}

// BytesFromBitString []byte from bit string, ex: "0101" -> ([]byte{0x50}, 4)
func BytesFromBitString(s string) ([]byte, int64) {
	r := len(s) % 8
	bufLen := len(s) / 8
	if r > 0 {
		bufLen++
	}
	buf := make([]byte, bufLen)

	for i := 0; i < len(s); i++ {
		d := s[i] - '0'
		if d != 0 && d != 1 {
			panic(fmt.Sprintf("invalid bit string %q at index %d %q", s, i, s[i]))
		}
		buf[i/8] |= d << (7 - i%8)
	}

	return buf, int64(len(s))
}

// BitStringFromBytes string from []byte], ex: ([]byte{0x50}, 4) -> "0101"
func BitStringFromBytes(buf []byte, nBits int64) string {
	sb := &strings.Builder{}
	for i := int64(0); i < nBits; i++ {
		if buf[i/8]&(1<<(7-i%8)) > 0 {
			sb.WriteString("1")
		} else {
			sb.WriteString("0")
		}
	}
	return sb.String()
}

// CopyBuffer bits from src to dst using provided buffer
// Similar to io.CopyBuffer
func CopyBuffer(dst Writer, src Reader, buf []byte) (n int64, err error) {
	// same default size as io.Copy
	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	var written int64

	for {
		rBits, rErr := src.ReadBits(buf, int64(len(buf))*8)
		if rBits > 0 {
			wBits, wErr := dst.WriteBits(buf, rBits)
			written += wBits
			if wErr != nil {
				err = wErr
				break
			}
			if rBits != wBits {
				err = io.ErrShortWrite
				break
			}
		}
		if rErr != nil {
			if !errors.Is(rErr, io.EOF) {
				err = rErr
			}
			break
		}
	}

	return written, err
}

// Copy bits from src to dst
// Similar to io.Copy
func Copy(dst Writer, src Reader) (n int64, err error) {
	return CopyBuffer(dst, src, nil)
}

// TODO: make faster, align and use copy()
func copyBits(dst []byte, dstStart int64, src []byte, srcStart int64, n int64, zero bool) {
	l := n
	off := int64(0)
	for l > 0 {
		c := int64(64)
		if l < c {
			c = l
		}
		u := Read64(src, srcStart+off, c)
		Write64(u, c, dst, dstStart+off)
		off += c
		l -= c
	}

	// zero fill last bits if not aligned
	e := dstStart + n
	if zero && e%8 != 0 {
		Write64(0, 8-(e%8), dst, e)
	}
}

// TODO: redo?
func readFull(p []byte, nBits int64, bitOff int64, fn func(p []byte, nBits int64, bitOff int64) (int64, error)) (int64, error) {
	if nBits < 0 {
		return 0, ErrNegativeNBits
	}

	readBitOffset := int64(0)
	for readBitOffset < nBits {
		byteOffset := readBitOffset / 8
		byteBitsOffset := readBitOffset % 8
		partialByteBitsLeft := (8 - byteBitsOffset) % 8
		leftBits := nBits - readBitOffset

		if partialByteBitsLeft != 0 || leftBits < 8 {
			readBits := partialByteBitsLeft
			if partialByteBitsLeft == 0 || leftBits < readBits {
				readBits = leftBits
			}

			var pb [1]byte
			rBits, err := fn(pb[:], readBits, bitOff+readBitOffset)
			Write64(uint64(pb[0]>>(8-rBits)), rBits, p, readBitOffset)
			readBitOffset += rBits

			if err != nil {
				return nBits - readBitOffset, err
			}

			continue
		}

		rBits, err := fn(p[byteOffset:], nBits-readBitOffset, bitOff+readBitOffset)

		readBitOffset += rBits
		if err != nil {
			return nBits - readBitOffset, err
		}
	}

	return nBits, nil
}

func ReadAtFull(r ReaderAt, p []byte, nBits int64, bitOff int64) (int64, error) {
	return readFull(p, nBits, bitOff, func(p []byte, nBits int64, bitOff int64) (int64, error) {
		return r.ReadBitsAt(p, nBits, bitOff)
	})
}

func ReadFull(r Reader, p []byte, nBits int64) (int64, error) {
	return readFull(p, nBits, 0, func(p []byte, nBits int64, bitOff int64) (int64, error) {
		return r.ReadBits(p, nBits)
	})
}
