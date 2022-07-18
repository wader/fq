package bitio

import (
	"errors"
	"io"
)

// TODO: smarter, track index?

func endPos(rs Seeker) (int64, error) {
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

// MultiReader is a bitio.ReaderAtSeeker concatinating multiple bitio.ReadAtSeeker:s.
// Similar to io.MultiReader.
type MultiReader struct {
	pos        int64
	readers    []ReadAtSeeker
	readerEnds []int64
}

// NewMultiReader returns a new bitio.MultiReader.
func NewMultiReader(rs ...ReadAtSeeker) (*MultiReader, error) {
	readerEnds := make([]int64, len(rs))
	var esSum int64
	for i, r := range rs {
		e, err := endPos(r)
		if err != nil {
			return nil, err
		}
		esSum += e
		readerEnds[i] = esSum
	}
	return &MultiReader{readers: rs, readerEnds: readerEnds}, nil
}

func (m *MultiReader) ReadBitsAt(p []byte, nBits int64, bitOff int64) (n int64, err error) {
	var end int64
	if len(m.readers) > 0 {
		end = m.readerEnds[len(m.readers)-1]
	}
	if end <= bitOff {
		return 0, io.EOF
	}

	prevAtEnd := int64(0)
	readerAt := m.readers[0]
	for i, end := range m.readerEnds {
		if bitOff < end {
			readerAt = m.readers[i]
			break
		}
		prevAtEnd = end
	}

	rBits, err := readerAt.ReadBitsAt(p, nBits, bitOff-prevAtEnd)

	if errors.Is(err, io.EOF) {
		if bitOff+rBits < end {
			err = nil
		}
	}

	return rBits, err
}

func (m *MultiReader) ReadBits(p []byte, nBits int64) (n int64, err error) {
	n, err = m.ReadBitsAt(p, nBits, m.pos)
	m.pos += n
	return n, err
}

func (m *MultiReader) SeekBits(bitOff int64, whence int) (int64, error) {
	var p int64
	var end int64
	if len(m.readers) > 0 {
		end = m.readerEnds[len(m.readers)-1]
	}

	switch whence {
	case io.SeekStart:
		p = bitOff
	case io.SeekCurrent:
		p = m.pos + bitOff
	case io.SeekEnd:
		p = end + bitOff
	default:
		panic("unknown whence")
	}
	if p < 0 || p > end {
		return 0, ErrOffset
	}

	m.pos = p

	return p, nil
}

func (m *MultiReader) CloneReader() (Reader, error) {
	return m.CloneReaderAtSeeker()
}

func (m *MultiReader) CloneReadSeeker() (ReadSeeker, error) {
	return m.CloneReaderAtSeeker()
}

func (m *MultiReader) CloneReadAtSeeker() (ReadAtSeeker, error) {
	return m.CloneReaderAtSeeker()
}

func (m *MultiReader) CloneReaderAtSeeker() (ReaderAtSeeker, error) {
	return &MultiReader{
		pos:        0,
		readers:    m.readers,
		readerEnds: m.readerEnds,
	}, nil
}
