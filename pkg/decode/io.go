package decode

import (
	"bytes"
	"io"

	"github.com/wader/fq/pkg/bitio"
)

func Copy(d *D, r io.Writer, w io.Reader) (int64, error) {
	// TODO: what size?
	buf := d.AllocReadBuf(64 * 1024)
	return io.CopyBuffer(r, w, buf)
}

func MustCopy(d *D, r io.Writer, w io.Reader) int64 {
	n, err := Copy(d, r, w)
	if err != nil {
		panic(IOError{Err: err, Op: "MustCopyBuffer"})
	}
	return n
}

func MustNewBitBufFromReader(d *D, r io.Reader) *bitio.Buffer {
	b := &bytes.Buffer{}
	MustCopy(d, b, r)
	return bitio.NewBufferFromBytes(b.Bytes(), -1)
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
