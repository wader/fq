// Package ctxreadseeker wraps a io.ReadSeeker and optionally io.Closer to make it context aware
// Warning: this might leak a go routine and a reader if underlying reader can block forever.
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
	fnCh   chan func()
	waitCh chan struct{}
}

func New(ctx context.Context, rs io.ReadSeeker) *Reader {
	r := &Reader{
		rs:     rs,
		ctx:    ctx,
		fnCh:   make(chan func()),
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
			fn()
			r.waitCh <- struct{}{}
		}
	}
}

func (r *Reader) callWait(fn func()) error {
	select {
	case <-r.ctx.Done():
		return r.ctx.Err()
	case r.fnCh <- fn:
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		case <-r.waitCh:
		}
	}
	return nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if err := r.callWait(func() {
		n, err = r.rs.Read(p)
	}); err != nil {
		return 0, err
	}
	return n, err
}

func (r *Reader) Seek(offset int64, whence int) (n int64, err error) {
	if err := r.callWait(func() {
		n, err = r.rs.Seek(offset, whence)
	}); err != nil {
		return 0, err
	}
	return n, err
}

func (r *Reader) Close() (err error) {
	if err := r.callWait(func() {
		if c, ok := r.rs.(io.Closer); ok {
			err = c.Close()
		}
	}); err != nil {
		return err
	}
	return err
}

func (r *Reader) Unwrap() any {
	return r.rs
}
