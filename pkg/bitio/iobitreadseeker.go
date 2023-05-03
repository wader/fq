package bitio

import (
	"errors"
	"io"
)

// IOBitReadSeeker is a bitio.ReadAtSeeker reading from a io.ReadSeeker.
type IOBitReadSeeker struct {
	bitPos int64
	rs     io.ReadSeeker
	buf    []byte
}

// NewIOBitReadSeeker returns a new bitio.IOBitReadSeeker
func NewIOBitReadSeeker(rs io.ReadSeeker) *IOBitReadSeeker {
	return &IOBitReadSeeker{
		bitPos: 0,
		rs:     rs,
	}
}

func (r *IOBitReadSeeker) ReadBitsAt(p []byte, nBits int64, bitOffset int64) (int64, error) {
	if nBits < 0 {
		return 0, ErrNegativeNBits
	}

	readBytePos := bitOffset / 8
	readSkipBits := bitOffset % 8
	wantReadBits := readSkipBits + nBits
	wantReadBytes := int(BitsByteCount(wantReadBits))

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
		nBits = int64(readBytes) * 8
		err = io.EOF
	}

	if readSkipBits == 0 && nBits%8 == 0 {
		copy(p[0:readBytes], r.buf[0:readBytes])
		return nBits, err
	}

	nBytes := nBits / 8
	restBits := nBits % 8

	// TODO: copy smartness if many bytes
	for i := int64(0); i < nBytes; i++ {
		p[i] = byte(Read64(r.buf, readSkipBits+i*8, 8))
	}
	if restBits != 0 {
		p[nBytes] = byte(Read64(r.buf, readSkipBits+nBytes*8, restBits)) << (8 - restBits)
	}

	return nBits, err
}

func (r *IOBitReadSeeker) ReadBits(p []byte, nBits int64) (n int64, err error) {
	rBits, err := r.ReadBitsAt(p, nBits, r.bitPos)
	r.bitPos += rBits
	return rBits, err
}

func (r *IOBitReadSeeker) SeekBits(bitOff int64, whence int) (int64, error) {
	seekBytesPos, err := r.rs.Seek(bitOff/8, whence)
	if err != nil {
		return 0, err
	}
	seekBitPos := seekBytesPos*8 + bitOff%8
	r.bitPos = seekBitPos

	return seekBitPos, nil
}

func (r *IOBitReadSeeker) CloneReader() (Reader, error) {
	return r.CloneReaderAtSeeker()
}

func (r *IOBitReadSeeker) CloneReadSeeker() (ReadSeeker, error) {
	return r.CloneReaderAtSeeker()
}

func (r *IOBitReadSeeker) CloneReadAtSeeker() (ReadAtSeeker, error) {
	return r.CloneReaderAtSeeker()
}

func (r *IOBitReadSeeker) CloneReaderAtSeeker() (ReaderAtSeeker, error) {
	return NewIOBitReadSeeker(r.rs), nil
}

func (r *IOBitReadSeeker) Unwrap() any {
	return r.rs
}
