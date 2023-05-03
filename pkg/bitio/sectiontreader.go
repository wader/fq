package bitio

import (
	"io"
)

// SectionReader is a bitio.BitReaderAtSeeker reading a section of a bitio.ReaderAt.
// Similar to io.SectionReader.
type SectionReader struct {
	r        ReaderAt
	bitBase  int64
	bitOff   int64
	bitLimit int64
}

// NewSectionReader returns a new bitio.SectionReader.
func NewSectionReader(r ReaderAt, bitOff int64, nBits int64) *SectionReader {
	return &SectionReader{
		r:        r,
		bitBase:  bitOff,
		bitOff:   bitOff,
		bitLimit: bitOff + nBits,
	}
}

func (r *SectionReader) ReadBitsAt(p []byte, nBits int64, bitOff int64) (int64, error) {
	if bitOff < 0 || bitOff >= r.bitLimit-r.bitBase {
		return 0, io.EOF
	}
	bitOff += r.bitBase
	if maxBits := r.bitLimit - bitOff; nBits > maxBits {
		nBits = maxBits
		rBits, err := r.r.ReadBitsAt(p, nBits, bitOff)
		return rBits, err
	}
	return r.r.ReadBitsAt(p, nBits, bitOff)
}

func (r *SectionReader) ReadBits(p []byte, nBits int64) (n int64, err error) {
	rBits, err := r.ReadBitsAt(p, nBits, r.bitOff-r.bitBase)
	r.bitOff += rBits
	return rBits, err
}

func (r *SectionReader) SeekBits(bitOff int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		bitOff += r.bitBase
	case io.SeekCurrent:
		bitOff += r.bitOff
	case io.SeekEnd:
		bitOff += r.bitLimit
	default:
		panic("unknown whence")
	}
	if bitOff < r.bitBase {
		return 0, ErrOffset
	}
	r.bitOff = bitOff
	return bitOff - r.bitBase, nil
}

func (r *SectionReader) CloneReader() (Reader, error) {
	return r.CloneReaderAtSeeker()
}

func (r *SectionReader) CloneReadSeeker() (ReadSeeker, error) {
	return r.CloneReaderAtSeeker()
}

func (r *SectionReader) CloneReaderSeeker() (ReadAtSeeker, error) {
	return r.CloneReaderAtSeeker()
}

func (r *SectionReader) CloneReaderAtSeeker() (ReaderAtSeeker, error) {
	return &SectionReader{
		r:        r.r,
		bitBase:  r.bitBase,
		bitOff:   r.bitBase,
		bitLimit: r.bitLimit,
	}, nil
}

func (r *SectionReader) Unwrap() any {
	return r.r
}
