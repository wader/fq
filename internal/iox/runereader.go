package iox

import (
	"io"
	"math/bits"
	"unicode/utf8"
)

// ByteRuneReader reads each byte as a rune from a io.ReadSeeker
// ex: when used with regexp \u00ff code point will match byte 0xff and not the
// utf-8 encoded version of 0xff
type ByteRuneReader struct {
	RS io.ReadSeeker
}

func (brr ByteRuneReader) ReadRune() (r rune, size int, err error) {
	var b [1]byte
	_, err = io.ReadFull(brr.RS, b[:])
	if err != nil {
		return 0, 0, err
	}
	r = rune(b[0])
	return r, 1, nil
}

func (brr ByteRuneReader) Seek(offset int64, whence int) (int64, error) {
	return brr.RS.Seek(offset, whence)
}

type RuneReadSeeker struct {
	RS io.ReadSeeker
}

func utf8Bytes(b byte) int {
	c := bits.LeadingZeros8(^b)
	// 0b0xxxxxxx 1 byte
	// 0b110xxxxx 2 byte
	// 0b1110xxxx 3 byte
	// 0b11110xxx 4 byte
	switch c {
	case 0:
		return 1
	case 2, 3, 4:
		return c
	default:
		return -1
	}
}

// ReadRune reads rune from a io.ReadSeeker
func (brr RuneReadSeeker) ReadRune() (r rune, size int, err error) {
	var b [utf8.UTFMax]byte

	_, err = io.ReadFull(brr.RS, b[0:1])
	if err != nil {
		return 0, 0, err
	}

	c := b[0]
	if c < utf8.RuneSelf {
		return rune(c), 1, nil
	}

	ss := utf8Bytes(b[0])
	if ss < 0 {
		return utf8.RuneError, 1, nil
	}

	_, err = io.ReadFull(brr.RS, b[1:ss])
	if err != nil {
		return 0, 0, err
	}

	r, s := utf8.DecodeRune(b[0:ss])
	// possibly rewind if DecodeRune fails as there was a invalid multi byte code point
	// TODO: better way that don't require seek back? buffer? one at a time?
	d := ss - s
	if d > 0 {
		if _, err := brr.Seek(int64(-d), io.SeekCurrent); err != nil {
			return 0, 0, err
		}
	}

	return r, s, nil
}

func (brr RuneReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return brr.RS.Seek(offset, whence)
}
