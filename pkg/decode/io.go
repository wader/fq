package decode

import (
	"bytes"
	"fq/pkg/bitio"
	"io"
)

func MustCopy(r io.Writer, w io.Reader) int64 {
	n, err := io.Copy(r, w)
	if err != nil {
		panic(err)
	}
	return n
}

func MustNewBitBufFromReader(r io.Reader) *bitio.Buffer {
	buf := &bytes.Buffer{}
	MustCopy(buf, r)
	return bitio.NewBufferFromBytes(buf.Bytes(), -1)
}
