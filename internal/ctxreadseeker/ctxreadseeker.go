// Package ctxreadseeker wraps a io.ReadSeeker and optionally io.Closer to make it context aware
// Warning: this might leak a go routine and a reader if underlaying reader can block forever.
// Only use if it's not an issue, your going to exit soon anyway or there is some other mechism for
// cleaning up.
package ctxreadseeker

import (
	"context"
	"io"
)

type Reader struct {
	rs     io.ReadSeeker
	ctx    context.Context
	fnCh   chan func(rs io.ReadSeeker)
	waitCh chan struct{}
}

func New(ctx context.Context, rs io.ReadSeeker) *Reader {
	r := &Reader{
		rs:     rs,
		ctx:    ctx,
		fnCh:   make(chan func(rs io.ReadSeeker)),
		waitCh: make(chan struct{}),
	}
	go r.loop()
	return r
}

func (r *Reader) loop() {
	for {
		select {
		case <-r.ctx.Done():
			if c, ok := r.rs.(io.Closer); ok {
				c.Close()
			}
			return
		case fn, ok := <-r.fnCh:
			if !ok {
				panic("unreachable")
			}
			fn(r.rs)
			r.waitCh <- struct{}{}
		}
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case r.fnCh <- func(rs io.ReadSeeker) {
		n, err = rs.Read(p)
	}:
		<-r.waitCh
	}
	return n, err
}

func (r *Reader) Seek(offset int64, whence int) (n int64, err error) {
	select {
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	case r.fnCh <- func(rs io.ReadSeeker) {
		n, err = rs.Seek(offset, whence)
	}:
		<-r.waitCh
	}
	return n, err
}

func (r *Reader) Close() (err error) {
	select {
	case <-r.ctx.Done():
		return r.ctx.Err()
	case r.fnCh <- func(rs io.ReadSeeker) {
		if c, ok := r.rs.(io.Closer); ok {
			err = c.Close()
		}
	}:
		<-r.waitCh
	}
	return err
}
