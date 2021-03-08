// Package ctxreadseeker wraps a io.ReadSeeker and optionally io.Closer to make it context aware
// Warning: this might leak a go routine and a reader if underlaying reader can block forever.
// Only use if it's not an issue, your going to exit soon anyway or there is some other mechism for
// cleaning up.
package ctxreadseeker

import (
	"context"
	"io"
)

type readCall struct {
	p []byte
}
type readReturn struct {
	n   int
	err error
}
type seekCall struct {
	offset int64
	whence int
}
type seekReturn struct {
	n   int64
	err error
}
type closeCall struct{}
type closeReturn struct {
	err error
}

type Reader struct {
	rs  io.ReadSeeker
	ctx context.Context
	ch  chan interface{}
}

func New(ctx context.Context, rs io.ReadSeeker) *Reader {
	ch := make(chan interface{})
	r := &Reader{
		rs:  rs,
		ctx: ctx,
		ch:  ch,
	}
	go r.loop()
	return r
}

func (r *Reader) loop() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case v, ok := <-r.ch:
			if !ok {
				return
			}
			switch v := v.(type) {
			case readCall:
				n, err := r.rs.Read(v.p)
				r.ch <- readReturn{n: n, err: err}
			case seekCall:
				n, err := r.rs.Seek(v.offset, v.whence)
				r.ch <- seekReturn{n: n, err: err}
			case closeCall:
				var err error
				if c, ok := r.rs.(io.Closer); ok {
					err = c.Close()
				}
				r.ch <- closeReturn{err: err}
				return
			default:
				panic("unreachable")
			}
		}
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case r.ch <- readCall{p: p}:
		select {
		case <-r.ctx.Done():
			return 0, r.ctx.Err()
		case v := <-r.ch:
			r := v.(readReturn)
			return r.n, r.err
		}
	}
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case r.ch <- seekCall{offset: offset, whence: whence}:
		select {
		case <-r.ctx.Done():
			return 0, r.ctx.Err()
		case v := <-r.ch:
			r := v.(seekReturn)
			return r.n, r.err
		}
	}
}

func (r *Reader) Close() error {
	select {
	case <-r.ctx.Done():
		return r.ctx.Err()
	case r.ch <- closeCall{}:
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		case v := <-r.ch:
			r := v.(closeReturn)
			return r.err
		}
	}
}
