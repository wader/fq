package interp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"

	"github.com/wader/gojq"
)

// TODO: move things in funcs.go/jq elsewere

func init() {
	functionRegisterFns = append(functionRegisterFns, func(i *Interp) []Function {
		return []Function{
			{"_hexdump", 1, 1, nil, i._hexdump},

			{"nal_unescape", 0, 0, makeBinaryTransformFn(func(r io.Reader) (io.Reader, error) {
				return &decode.NALUnescapeReader{Reader: r}, nil
			}), nil},

			{"aes_ctr", 1, 2, i.aesCtr, nil},
		}
	})
}

// transform to binary using fn
func makeBinaryTransformFn(fn func(r io.Reader) (io.Reader, error)) func(c any, a []any) any {
	return func(c any, a []any) any {
		inBR, err := toBitReader(c)
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

		bb, err := newBinaryFromBitReader(outBR, 8, 0)
		if err != nil {
			return err
		}
		return bb
	}
}

func (i *Interp) aesCtr(c any, a []any) any {
	keyBytes, err := toBytes(a[0])
	if err != nil {
		return err
	}

	switch len(keyBytes) {
	case 16, 24, 32:
	default:
		return fmt.Errorf("key length should be 16, 24 or 32 bytes, is %d bytes", len(keyBytes))
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return err
	}

	var ivBytes []byte
	if len(a) >= 2 {
		var err error
		ivBytes, err = toBytes(a[1])
		if err != nil {
			return err
		}
		if len(ivBytes) != block.BlockSize() {
			return fmt.Errorf("iv length should be %d bytes, is %d bytes", block.BlockSize(), len(ivBytes))
		}
	} else {
		ivBytes = make([]byte, block.BlockSize())
	}

	br, err := toBitReader(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	reader := &cipher.StreamReader{S: cipher.NewCTR(block, ivBytes), R: bitio.NewIOReader(br)}
	if _, err := io.Copy(buf, reader); err != nil {
		return err
	}

	bb, err := newBinaryFromBitReader(bitio.NewBitReader(buf.Bytes(), -1), 8, 0)
	if err != nil {
		return err
	}
	return bb
}

func (i *Interp) _hexdump(c any, a []any) gojq.Iter {
	opts := i.Options(a[0])
	bv, err := toBinary(c)
	if err != nil {
		return gojq.NewIter(err)
	}
	if err := hexdump(i.evalInstance.output, bv, opts); err != nil {
		return gojq.NewIter(err)
	}

	return gojq.NewIter()
}
