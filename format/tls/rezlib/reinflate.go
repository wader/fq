// Package rezlib wraps a zlib reader and makes it possible to read
// until the last current input flush boundary by reading until EOF.
//
// Inspiration from https://github.com/golang/go/issues/48877
//
// TODO: force deflate only?
//
// This is used by TLS deflate (seems be zlib?) where a flush is done
// between each record
package rezlib

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
)

// max deflate dictionary size
const maxDictSize = 1 << 15 // 32KiB

type Reader struct {
	r                io.Reader
	fr               io.Reader
	res              zlib.Resetter
	currentDict      bytes.Buffer
	zlibHeader       [2]byte
	hadUnexpectedEOF bool
}

func NewReader(r io.Reader) (*Reader, error) {
	var zlibHeader [2]byte
	if _, err := io.ReadFull(r, zlibHeader[:]); err != nil {
		return nil, err
	}

	zr, zrErr := zlib.NewReader(io.MultiReader(
		bytes.NewBuffer(zlibHeader[:]),
		r,
	))
	if zrErr != nil {
		return nil, zrErr
	}

	res, resOk := zr.(zlib.Resetter)
	if !resOk {
		panic("zlib reader not a Resetter")
	}

	return &Reader{
		r:          r,
		fr:         zr,
		res:        res,
		zlibHeader: zlibHeader,
	}, nil
}

func (r *Reader) Read(b []byte) (int, error) {
	if r.hadUnexpectedEOF {
		r.hadUnexpectedEOF = false
		// inject zlib header again
		mr := io.MultiReader(
			bytes.NewReader(r.zlibHeader[:]),
			r.r,
		)
		if err := r.res.Reset(mr, r.currentDict.Bytes()); err != nil {
			panic(fmt.Sprintf("zlib reader could not reset %s", err))
		}
	}

	n, err := r.fr.Read(b)

	r.currentDict.Write(b[:n])
	l := r.currentDict.Len()
	if l > maxDictSize {
		r.currentDict.Next(l - maxDictSize)
	}

	// TODO: currently we assume ErrUnexpectedEOF means we have read
	// until a flush boundary. To get read of the error reset with
	// current dictionary
	// some better way?
	if errors.Is(err, io.ErrUnexpectedEOF) {
		r.hadUnexpectedEOF = true
		err = io.EOF
	}

	return n, err
}
