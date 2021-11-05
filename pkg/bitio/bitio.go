package bitio

import (
	"errors"
	"io"
)

var ErrOffset = errors.New("invalid seek offset")
var ErrNegativeNBits = errors.New("negative number of bits")

type BitReaderAt interface {
	ReadBitsAt(p []byte, nBits int, bitOff int64) (n int, err error)
}

type BitReader interface {
	ReadBits(p []byte, nBits int) (n int, err error)
}

type BitSeeker interface {
	SeekBits(bitOffset int64, whence int) (int64, error)
}

type BitReadSeeker interface {
	BitReader
	BitSeeker
}

type BitReadAtSeeker interface {
	BitReaderAt
	BitSeeker
}

type BitWriter interface {
	WriteBits(p []byte, nBits int) (n int, err error)
}

func CopyBuffer(dst BitWriter, src BitReader, buf []byte) (n int64, err error) {
	// same default size as io.Copy
	if buf == nil {
		buf = make([]byte, 32*1024)
	}
	var written int64

	for {
		rBits, rErr := src.ReadBits(buf, len(buf)*8)
		if rBits > 0 {
			wBits, wErr := dst.WriteBits(buf, rBits)
			written += int64(wBits)
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

func Copy(dst BitWriter, src BitReader) (n int64, err error) {
	return CopyBuffer(dst, src, nil)
}

// BitsByteCount returns smallest amount of bytes to fit nBits bits
func BitsByteCount(nBits int64) int64 {
	n := nBits / 8
	if nBits%8 != 0 {
		n++
	}
	return n
}

func readFull(p []byte, nBits int, bitOff int64, fn func(p []byte, nBits int, bitOff int64) (int, error)) (int, error) {
	if nBits < 0 {
		return 0, ErrNegativeNBits
	}

	readBitOffset := 0
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
			rBits, err := fn(pb[:], readBits, bitOff+int64(readBitOffset))
			Write64(uint64(pb[0]>>(8-rBits)), rBits, p, readBitOffset)
			readBitOffset += rBits

			if err != nil {
				return nBits - readBitOffset, err
			}

			continue
		}

		rBits, err := fn(p[byteOffset:], nBits-readBitOffset, bitOff+int64(readBitOffset))

		readBitOffset += rBits
		if err != nil {
			return nBits - readBitOffset, err
		}
	}

	return nBits, nil
}

func ReadAtFull(r BitReaderAt, p []byte, nBits int, bitOff int64) (int, error) {
	return readFull(p, nBits, bitOff, func(p []byte, nBits int, bitOff int64) (int, error) {
		return r.ReadBitsAt(p, nBits, bitOff)
	})
}

func ReadFull(r BitReader, p []byte, nBits int) (int, error) {
	return readFull(p, nBits, 0, func(p []byte, nBits int, bitOff int64) (int, error) {
		return r.ReadBits(p, nBits)
	})
}

// TODO: move?
func EndPos(rs BitSeeker) (int64, error) {
	c, err := rs.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	e, err := rs.SeekBits(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	_, err = rs.SeekBits(c, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return e, nil
}
