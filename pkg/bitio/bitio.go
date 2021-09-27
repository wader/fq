package bitio

// TODO: should return int64?
// TODO: document len(p)/nBits, should be +1 for when not aligned

import (
	"errors"
	"io"
)

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

type AlignBitWriter struct {
	W BitWriter
	N int
	c int64
}

func (a *AlignBitWriter) WriteBits(p []byte, nBits int) (n int, err error) {
	n, err = a.W.WriteBits(p, nBits)
	a.c += int64(n)
	return n, err
}

func (a *AlignBitWriter) Close() error {
	n := int64(a.N)
	r := int((n - a.c%n) % n)
	if r == 0 {
		return nil
	}
	b := make([]byte, a.N/8+1)
	_, err := a.W.WriteBits(b, r)
	return err
}

type AlignBitReader struct {
	R BitReaderAt
	N int
	c int64
}

func (a *AlignBitReader) ReadBitsAt(p []byte, nBits int, bitOff int64) (n int, err error) {
	n, err = a.R.ReadBitsAt(p, nBits, bitOff)
	a.c += int64(n)
	return n, err
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

func BitsByteCount(nBits int64) int64 {
	n := nBits / 8
	if nBits%8 != 0 {
		n++
	}
	return n
}

func readFull(p []byte, nBits int, bitOff int64, fn func(p []byte, nBits int, bitOff int64) (int, error)) (int, error) {
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

// Reader is a BitReadSeeker and BitReaderAt reading from a io.ReadSeeker
// TODO: private?
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
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return 0, err
	} else if errors.Is(err, io.ErrUnexpectedEOF) {
		nBits = readBytes * 8
		err = io.EOF
	}

	if readSkipBits == 0 && nBits%8 == 0 {
		copy(p[0:readBytes], r.buf[0:readBytes])
		return nBits, err
	}

	nBytes := nBits / 8
	restBits := nBits % 8

	// TODO: copy smartness if many bytes
	for i := 0; i < nBytes; i++ {
		p[i] = byte(Read64(r.buf, readSkipBits+i*8, 8))
	}
	if restBits != 0 {
		p[nBytes] = byte(Read64(r.buf, readSkipBits+nBytes*8, restBits)) << (8 - restBits)
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
	return int(BitsByteCount(int64(n))), err
}

func (r *SectionBitReader) Seek(offset int64, whence int) (int64, error) {
	seekBytePos, err := r.SeekBits(offset*8, whence)
	return seekBytePos * 8, err
}

// TODO: smart, track index?
type MultiBitReader struct {
	pos        int64
	readers    []BitReadAtSeeker
	readerEnds []int64
}

func NewMultiBitReader(rs []BitReadAtSeeker) (*MultiBitReader, error) {
	readerEnds := make([]int64, len(rs))
	var esSum int64
	for i, r := range rs {
		e, err := EndPos(r)
		if err != nil {
			return nil, err
		}
		esSum += e
		readerEnds[i] = esSum
	}
	return &MultiBitReader{readers: rs, readerEnds: readerEnds}, nil
}

func (m *MultiBitReader) ReadBitsAt(p []byte, nBits int, bitOff int64) (n int, err error) {
	end := m.readerEnds[len(m.readers)-1]
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
		if bitOff+int64(rBits) < end {
			err = nil
		}
	}

	return rBits, err
}

func (m *MultiBitReader) ReadBits(p []byte, nBits int) (n int, err error) {
	n, err = m.ReadBitsAt(p, nBits, m.pos)
	m.pos += int64(n)
	return n, err
}

func (m *MultiBitReader) SeekBits(bitOff int64, whence int) (int64, error) {
	var p int64
	end := m.readerEnds[len(m.readerEnds)-1]

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
		return 0, errOffset
	}

	m.pos = p

	return p, nil
}

func (m *MultiBitReader) Read(p []byte) (n int, err error) {
	n, err = m.ReadBitsAt(p, len(p)*8, m.pos)
	m.pos += int64(n)

	// log.Printf("n: %#+v\n", n)
	// log.Printf("err: %#+v\n", err)

	if err != nil {
		return int(BitsByteCount(int64(n))), err
	}

	return int(BitsByteCount(int64(n))), nil
}
