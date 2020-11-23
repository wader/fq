package bitio

import (
	"errors"
	"io"
)

func BitsByteCount(nBits int64) int64 {
	n := nBits / 8
	if nBits%8 != 0 {
		n++
	}
	return n
}

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

// Reader is a BitReadSeeker and BitReaderAt reading from a io.ReadSeeker
type Reader struct {
	bitPos int64
	rs     io.ReadSeeker
	buf    []byte
}

func NewReaderFromReadSeeker(rs io.ReadSeeker) *Reader {
	return &Reader{
		bitPos: 0,
		rs:     rs,
	}
}

func (r *Reader) ReadBitsAt(p []byte, nBits int, bitOffset int64) (int, error) {
	readBytePos := bitOffset / 8
	readSkipBits := int(bitOffset % 8)
	wantReadBits := readSkipBits + nBits
	wantReadBytes := int(BitsByteCount(int64(wantReadBits)))

	if wantReadBytes > len(r.buf) {
		// TODO: use append somehow?
		r.buf = make([]byte, wantReadBytes)
	}

	_, err := r.rs.Seek(readBytePos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	// TODO: nBits should be available
	readBytes, err := io.ReadFull(r.rs, r.buf[0:wantReadBytes])
	if err != nil && err != io.ErrUnexpectedEOF {
		return 0, err
	} else if err == io.ErrUnexpectedEOF {
		diffBytes := wantReadBytes - readBytes
		nBits = readSkipBits - 8
		if readSkipBits != 0 {
			diffBytes--
		}
		nBits += 8 * diffBytes
	}

	if readSkipBits == 0 && nBits%8 == 0 {
		copy(p[0:readBytes], r.buf[0:readBytes])
		return nBits, err
	}

	nBytes := int(nBits / 8)
	restBits := nBits % 8

	// TODO: copy smartness if many bytes
	for i := 0; i < nBytes; i++ {
		p[i] = byte(Uint64(r.buf, readSkipBits+i*8, 8))
	}
	if restBits != 0 {
		p[nBytes] = byte(Uint64(r.buf, readSkipBits+nBytes*8, restBits)) << (8 - restBits)
	}

	return nBits, err
}

func (r *Reader) ReadBits(p []byte, nBits int) (n int, err error) {
	rBits, err := r.ReadBitsAt(p, nBits, r.bitPos)
	r.bitPos += int64(rBits)
	return rBits, err
}

func (r *Reader) SeekBits(bitOff int64, whence int) (int64, error) {
	seekBytesPos, err := r.rs.Seek(bitOff/8, whence)
	if err != nil {
		return 0, err
	}
	seekBitPos := seekBytesPos*8 + bitOff%8
	r.bitPos = seekBitPos

	return seekBitPos, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.ReadBitsAt(p, len(p)*8, r.bitPos)
	r.bitPos += int64(n)
	if err != nil {
		return int(BitsByteCount(int64(n))), err
	}

	return int(BitsByteCount(int64(n))), nil
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	seekBytesPos, err := r.rs.Seek(offset, whence)
	if err != nil {
		return 0, err
	}
	r.bitPos = seekBytesPos * 8
	return seekBytesPos, nil
}

// SectionBitReader is a BitReadSeeker reading from a BitReaderAt
// modelled after io.SectionReader
type SectionBitReader struct {
	r        BitReaderAt
	bitBase  int64
	bitOff   int64
	bitLimit int64
}

func NewSectionBitReader(r BitReaderAt, bitOff int64, nBits int64) *SectionBitReader {
	return &SectionBitReader{
		r:        r,
		bitBase:  bitOff,
		bitOff:   bitOff,
		bitLimit: bitOff + nBits,
	}
}

func (r *SectionBitReader) ReadBitsAt(p []byte, nBits int, bitOff int64) (int, error) {
	if bitOff < 0 || bitOff >= r.bitLimit-r.bitBase {
		return 0, io.EOF
	}
	bitOff += r.bitBase
	if maxBits := int(r.bitLimit - bitOff); nBits > maxBits {
		nBits = maxBits
		rBits, err := r.r.ReadBitsAt(p, nBits, bitOff)
		if err == nil {
			err = io.EOF
		}
		return rBits, err
	}
	return r.r.ReadBitsAt(p, nBits, bitOff)
}

func (r *SectionBitReader) ReadBits(p []byte, nBits int) (n int, err error) {
	rBits, err := r.ReadBitsAt(p, nBits, r.bitOff-r.bitBase)
	r.bitOff += int64(rBits)
	return rBits, err
}

var errOffset = errors.New("invalid seek offset")

func (r *SectionBitReader) SeekBits(bitOff int64, whence int) (int64, error) {
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
		return 0, errOffset
	}
	r.bitOff = bitOff
	return bitOff - r.bitBase, nil
}

func (r *SectionBitReader) Read(p []byte) (n int, err error) {
	n, err = r.ReadBitsAt(p, len(p)*8, r.bitOff-r.bitBase)
	r.bitOff += int64(n)
	if err != nil {
		return int(BitsByteCount(int64(n))), err
	}

	return int(BitsByteCount(int64(n))), nil
}

func (r *SectionBitReader) Seek(offset int64, whence int) (int64, error) {
	seekBytePos, err := r.SeekBits(offset*8, whence)
	return seekBytePos * 8, err
}
