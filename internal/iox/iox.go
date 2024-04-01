package iox

import (
	"context"
	"errors"
	"io"
)

func SeekerEnd(s io.Seeker) (int64, error) {
	cPos, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	epos, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	if _, err := s.Seek(cPos, io.SeekStart); err != nil {
		return 0, err
	}

	return epos, nil
}

type ReadErrSeeker struct{ io.Reader }

func (r *ReadErrSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, errors.New("seek")
}

type CtxWriter struct {
	io.Writer
	Ctx context.Context
}

func (o CtxWriter) Write(p []byte) (n int, err error) {
	if o.Ctx != nil {
		if err := o.Ctx.Err(); err != nil {
			return 0, err
		}
	}
	return o.Writer.Write(p)
}

type DiscardCtxWriter struct {
	Ctx context.Context
}

func (o DiscardCtxWriter) Write(p []byte) (n int, err error) {
	if o.Ctx != nil {
		if err := o.Ctx.Err(); err != nil {
			return 0, err
		}
	}
	return n, nil
}

func Unwrap(r any) any {
	for {
		u, ok := r.(interface {
			Unwrap() any
		})
		if !ok {
			return r
		}
		r = u.Unwrap()
	}
}
