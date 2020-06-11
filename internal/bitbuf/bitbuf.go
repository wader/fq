package bitbuf

type Buffer struct {
	Buf   []byte
	Start uint64
	Len   uint64
	Pos   uint64
}

func New(buf []byte) *Buffer {
	return &Buffer{
		Buf:   buf,
		Start: 0,
		Len:   uint64(len(buf)) * 8,
		Pos:   0,
	}
}

func (b *Buffer) Bits(nBits uint64) (uint64, uint64) {
	p := uint64(b.Pos) + uint64(nBits)
	if p > b.Len {
		return 0, uint64(p) - b.Len
	}
	return ReadBits(b.Buf, b.Start+b.Pos, nBits), nBits
}

func (b *Buffer) BitBufRange(start uint64, nBits uint64) (*Buffer, uint64) {
	endPos := uint64(start) + uint64(nBits)
	if endPos > b.Len {
		return nil, uint64(endPos) - b.Len
	}
	return &Buffer{
		Buf:   b.Buf,
		Start: b.Start + start,
		Len:   nBits,
		Pos:   0,
	}, nBits
}

func (b *Buffer) BitBufLen(nBits uint64) (*Buffer, uint64) {
	return b.BitBufRange(b.Pos, nBits)
}

func (b *Buffer) BytesRange(start uint64, nBytes uint64) ([]byte, uint64) {
	endPos := uint64(start) + uint64(nBytes*8)
	if endPos > b.Len {
		return nil, uint64(endPos) - b.Len
	}

	bufStart := b.Start + b.Pos
	if bufStart%8 == 0 {
		return b.Buf[bufStart : bufStart+nBytes], nBytes * 8
	}

	var buf []byte
	for i := uint64(0); i < nBytes; i++ {
		buf = append(buf, byte(ReadBits(b.Buf, bufStart+i, 8)))
	}
	return buf, nBytes * 8
}

func (b *Buffer) BytesLen(nBytes uint64) ([]byte, uint64) {
	return b.BytesRange(b.Pos, nBytes)
}
