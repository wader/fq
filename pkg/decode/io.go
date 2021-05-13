package decode

import (
	"bytes"
	"fq/pkg/bitio"
	"io"
)

func MustCopy(r io.Writer, w io.Reader) int64 {
	n, err := io.Copy(r, w)
	if err != nil {
		panic(ReadError{Err: err, Op: "MustCopy"})
	}
	return n
}

func MustNewBitBufFromReader(r io.Reader) *bitio.Buffer {
	buf := &bytes.Buffer{}
	MustCopy(buf, r)
	return bitio.NewBufferFromBytes(buf.Bytes(), -1)
}

// TODO: move?
// TODO: make generic replace reader? share with id3v2 unsync?
type NALUnescapeReader struct {
	io.Reader
	lastTwoZeros [2]bool
}

func (r NALUnescapeReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)

	ni := 0
	for i, b := range p[0:n] {
		if r.lastTwoZeros[0] && r.lastTwoZeros[1] && b == 0x03 {
			n--
			r.lastTwoZeros[0] = false
			r.lastTwoZeros[1] = false
			continue
		} else {
			r.lastTwoZeros[1] = r.lastTwoZeros[0]
			r.lastTwoZeros[0] = b == 0
		}
		p[ni] = p[i]
		ni++
	}

	return n, err
}
