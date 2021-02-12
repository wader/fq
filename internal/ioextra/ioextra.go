package ioextra

import (
	"context"
	"io"
)

type ContextWriter struct {
	C context.Context
	W io.Writer
}

func (cw ContextWriter) Write(p []byte) (n int, err error) {
	if err := cw.C.Err(); err != nil {
		return 0, err
	}
	return cw.W.Write(p)
}

func MustCopy(r io.Writer, w io.Reader) int64 {
	n, err := io.Copy(r, w)
	if err != nil {
		panic(err)
	}
	return n
}
