package mpeg

import (
	"bytes"
	"io"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFunc0("nal_unescape", makeBinaryTransformFn(func(r io.Reader) (io.Reader, error) {
		return &nalUnescapeReader{Reader: r}, nil
	}))
}

// transform to binary using fn
func makeBinaryTransformFn(fn func(r io.Reader) (io.Reader, error)) func(_ *interp.Interp, c any) any {
	return func(_ *interp.Interp, c any) any {
		inBR, err := interp.ToBitReader(c)
		if err != nil {
			return err
		}

		r, err := fn(bitio.NewIOReader(inBR))
		if err != nil {
			return err
		}

		outBuf := &bytes.Buffer{}
		if _, err := io.Copy(outBuf, r); err != nil {
			return err
		}

		outBR := bitio.NewBitReader(outBuf.Bytes(), -1)

		bb, err := interp.NewBinaryFromBitReader(outBR, 8, 0)
		if err != nil {
			return err
		}
		return bb
	}
}

func decodeEscapeValueFn(add int, b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return func(d *decode.D) uint64 {
		n1 := d.U(b1)
		n := n1
		if n1 == (1<<b1)-1 {
			n2 := d.U(b2)
			if add != -1 {
				n += n2 + uint64(add)
			} else {
				n = n2
			}
			if n2 == (1<<b2)-1 {
				n3 := d.U(b3)
				if add != -1 {
					n += n3 + uint64(add)
				} else {
					n = n3
				}
			}
		}
		return n
	}
}

// use last non-escaped value
func decodeEscapeValueAbsFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(-1, b1, b2, b3)
}

// add values and escaped values
//
//nolint:deadcode,unused
func decodeEscapeValueAddFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(0, b1, b2, b3)
}

// add values and escaped values+1
func decodeEscapeValueCarryFn(b1 int, b2 int, b3 int) func(d *decode.D) uint64 {
	return decodeEscapeValueFn(1, b1, b2, b3)
}

// TODO: move?
// TODO: make generic replace reader? share with id3v2 unsync?
type nalUnescapeReader struct {
	io.Reader
	lastTwoZeros [2]bool
}

func (r nalUnescapeReader) Read(p []byte) (n int, err error) {
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
