package bitiox

// bitio helpers that im not sure belong in bitio

import (
	"errors"
	"io"

	"github.com/wader/fq/pkg/bitio"
)

func CopyBitsBuffer(dst io.Writer, src bitio.Reader, buf []byte) (int64, error) {
	return io.CopyBuffer(dst, bitio.NewIOReader(src), buf)
}

func CopyBits(dst io.Writer, src bitio.Reader) (int64, error) {
	return CopyBitsBuffer(dst, src, nil)
}

func Len(br bitio.ReadAtSeeker) (int64, error) {
	bPos, err := br.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	bEnd, err := br.SeekBits(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	if _, err := br.SeekBits(bPos, io.SeekStart); err != nil {
		return 0, err
	}
	return bEnd, nil
}

func Range(br bitio.ReadAtSeeker, firstBitOffset int64, nBits int64) (bitio.ReaderAtSeeker, error) {
	l, err := Len(br)
	if err != nil {
		return nil, err
	}
	// TODO: move error check?
	if nBits < 0 {
		return nil, errors.New("negative nBits")
	}
	if firstBitOffset+nBits > l {
		return nil, errors.New("outside buffer")
	}
	return bitio.NewSectionReader(br, firstBitOffset, nBits), nil
}
