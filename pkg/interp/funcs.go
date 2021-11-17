package interp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"

	"github.com/wader/gojq"
)

func init() {
	functionRegisterFns = append(functionRegisterFns, func(i *Interp) []Function {
		return []Function{
			{"_hexdump", 1, 1, nil, i._hexdump},

			{"hex", 0, 0, makeStringBitBufTransformFn(
				func(r io.Reader) (io.Reader, error) { return hex.NewDecoder(r), nil },
				func(r io.Writer) (io.Writer, error) { return hex.NewEncoder(r), nil },
			), nil},

			{"base64", 0, 0, makeStringBitBufTransformFn(
				func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.StdEncoding, r), nil },
				func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.StdEncoding, r), nil },
			), nil},
			{"rawbase64", 0, 0, makeStringBitBufTransformFn(
				func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.RawURLEncoding, r), nil },
				func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.RawURLEncoding, r), nil },
			), nil},

			{"urlbase64", 0, 0, makeStringBitBufTransformFn(
				func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.URLEncoding, r), nil },
				func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.URLEncoding, r), nil },
			), nil},

			{"nal_unescape", 0, 0, makeBitBufTransformFn(func(r io.Reader) (io.Reader, error) {
				return &decode.NALUnescapeReader{Reader: r}, nil
			}), nil},

			{"md5", 0, 0, makeHashFn(func() (hash.Hash, error) { return md5.New(), nil }), nil},

			{"query_escape", 0, 0, i.queryEscape, nil},
			{"query_unescape", 0, 0, i.queryUnescape, nil},
			{"path_escape", 0, 0, i.pathEscape, nil},
			{"path_unescape", 0, 0, i.pathUnescape, nil},
			{"aes_ctr", 1, 2, i.aesCtr, nil},
		}
	})
}

// transform byte string <-> buffer using fn:s
func makeStringBitBufTransformFn(
	decodeFn func(r io.Reader) (io.Reader, error),
	encodeFn func(w io.Writer) (io.Writer, error),
) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		switch c := c.(type) {
		case string:
			bb, err := toBuffer(c)
			if err != nil {
				return err
			}

			r, err := decodeFn(bb)
			if err != nil {
				return err
			}

			buf := &bytes.Buffer{}
			if _, err := io.Copy(buf, r); err != nil {
				return err
			}
			outBB := bitio.NewBufferFromBytes(buf.Bytes(), -1)

			return newBufferRangeFromBuffer(outBB, 8)
		default:
			bb, err := toBuffer(c)
			if err != nil {
				return err
			}

			buf := &bytes.Buffer{}
			w, err := encodeFn(buf)
			if err != nil {
				return err
			}

			if _, err := io.Copy(w, bb); err != nil {
				return err
			}

			if c, ok := w.(io.Closer); ok {
				c.Close()
			}

			return buf.String()
		}
	}
}

// transform to buffer using fn
func makeBitBufTransformFn(fn func(r io.Reader) (io.Reader, error)) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		inBB, err := toBuffer(c)
		if err != nil {
			return err
		}

		r, err := fn(inBB)
		if err != nil {
			return err
		}

		outBuf := &bytes.Buffer{}
		if _, err := io.Copy(outBuf, r); err != nil {
			return err
		}

		outBB := bitio.NewBufferFromBytes(outBuf.Bytes(), -1)

		return newBufferRangeFromBuffer(outBB, 8)
	}
}

// transform to buffer using fn
func makeHashFn(fn func() (hash.Hash, error)) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		inBB, err := toBuffer(c)
		if err != nil {
			return err
		}

		h, err := fn()
		if err != nil {
			return err
		}
		if _, err := io.Copy(h, inBB); err != nil {
			return err
		}

		outBB := bitio.NewBufferFromBytes(h.Sum(nil), -1)

		return newBufferRangeFromBuffer(outBB, 8)
	}
}

func (i *Interp) _registry(c interface{}, a []interface{}) interface{} {
	uniqueFormats := map[string]decode.Format{}

	groups := map[string]interface{}{}
	formats := map[string]interface{}{}

	for fsName := range i.registry.Groups {
		var group []interface{}

		for _, f := range i.registry.MustGroup(fsName) {
			group = append(group, f.Name)
			if _, ok := uniqueFormats[f.Name]; ok {
				continue
			}
			uniqueFormats[f.Name] = f
		}

		groups[fsName] = group
	}

	for _, f := range uniqueFormats {
		vf := map[string]interface{}{
			"name":        f.Name,
			"description": f.Description,
			"probe_order": f.ProbeOrder,
			"root_name":   f.RootName,
			"root_array":  f.RootArray,
		}

		var dependenciesVs []interface{}
		for _, d := range f.Dependencies {
			var dNamesVs []interface{}
			for _, n := range d.Names {
				dNamesVs = append(dNamesVs, n)
			}
			dependenciesVs = append(dependenciesVs, dNamesVs)
		}
		if len(dependenciesVs) > 0 {
			vf["dependencies"] = dependenciesVs
		}
		var groupsVs []interface{}
		for _, n := range f.Groups {
			groupsVs = append(groupsVs, n)
		}
		if len(groupsVs) > 0 {
			vf["groups"] = groupsVs
		}

		if f.Files != nil {
			files := map[string]interface{}{}

			entries, err := f.Files.ReadDir(".")
			if err != nil {
				return err
			}

			for _, e := range entries {
				f, err := f.Files.Open(e.Name())
				if err != nil {
					return err
				}
				b, err := ioutil.ReadAll(f)
				if err != nil {
					return err
				}
				files[e.Name()] = string(b)
			}

			vf["files"] = files
		}

		formats[f.Name] = vf
	}

	return map[string]interface{}{
		"groups":  groups,
		"formats": formats,
	}
}

func (i *Interp) queryEscape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	return url.QueryEscape(s)
}

func (i *Interp) queryUnescape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	u, err := url.QueryUnescape(s)
	if err != nil {
		return err
	}
	return u
}
func (i *Interp) pathEscape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	return url.PathEscape(s)
}

func (i *Interp) pathUnescape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	u, err := url.PathUnescape(s)
	if err != nil {
		return err
	}
	return u
}

func (i *Interp) aesCtr(c interface{}, a []interface{}) interface{} {
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

	bb, err := toBuffer(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	reader := &cipher.StreamReader{S: cipher.NewCTR(block, ivBytes), R: bb}
	if _, err := io.Copy(buf, reader); err != nil {
		return err
	}

	return buf.Bytes()
}

func (i *Interp) _hexdump(c interface{}, a []interface{}) gojq.Iter {
	opts := i.Options(a[0])
	bv, err := toBufferView(c)
	if err != nil {
		return gojq.NewIter(err)
	}
	if err := hexdump(i.evalContext.output, bv, opts); err != nil {
		return gojq.NewIter(err)
	}

	return gojq.NewIter()
}
