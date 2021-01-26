package query

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"fq/internal/hexdump"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"io"
	"io/ioutil"
	"math/big"
	"net/url"
	"strings"
)

func (q *Query) hexdump(c interface{}, a []interface{}) interface{} {
	bb, r, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	return func(stdout io.Writer) error {
		bitsByteAlign := r.Start % 8
		bb, err := bb.BitBufRange(r.Start-bitsByteAlign, r.Len+bitsByteAlign)
		if err != nil {
			return err
		}
		hw := hexdump.New(
			stdout,
			(r.Start-bitsByteAlign)/8,
			num.DigitsInBase(bitio.BitsByteCount(r.Stop()+bitsByteAlign), 16),
			q.opts.DumpOptions.LineBytes)
		if _, err := io.Copy(hw, bb); err != nil {
			return err
		}
		hw.Close()
		return nil
	}
}

func (q *Query) formats(c interface{}, a []interface{}) interface{} {

	allFormats := map[string]*decode.Format{}

	for _, fs := range q.opts.Registry.Groups {
		for _, f := range fs {
			if _, ok := allFormats[f.Name]; ok {
				continue
			}
			allFormats[f.Name] = f
		}
	}

	vs := map[string]interface{}{}
	for _, f := range allFormats {
		vf := map[string]interface{}{
			"name":        f.Name,
			"description": f.Description,
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

		vs[f.Name] = vf
	}

	return vs
}

func (q *Query) preview(c interface{}, a []interface{}) interface{} {
	v, err := toValue(c)
	if err != nil {
		return fmt.Errorf("%v: value is not a decode value", c)
	}
	return func(stdout io.Writer) error {
		if err := v.Preview(stdout); err != nil {
			return err
		}
		return nil
	}
}

func (q *Query) help(c interface{}, a []interface{}) interface{} {
	return queryErrorFn(func(stdout io.Writer) error {
		for _, f := range q.functions {
			var names []string
			for _, n := range f.Names {
				for j := f.MinArity; j <= f.MaxArity; j++ {
					names = append(names, fmt.Sprintf("%s/%d", n, j))
				}
			}
			fmt.Fprintf(stdout, "%s\n", strings.Join(names, ", "))
		}
		return nil
	})
}

func (q *Query) open(c interface{}, a []interface{}) interface{} {
	var rs io.ReadSeeker

	var filename string
	if len(a) == 1 {
		var err error
		filename, err = toString(a[0])
		if err != nil {
			return fmt.Errorf("%s: %w", filename, err)
		}
	}

	if filename == "" || filename == "-" {
		filename = "stdin"
		buf, err := ioutil.ReadAll(q.opts.OS.Stdin())
		if err != nil {
			return err
		}
		rs = bytes.NewReader(buf)
	} else {
		f, err := q.opts.OS.Open(filename)
		if err != nil {
			return err
		}

		// TODO: cleanup? bitbuf have optional close method etc?
		// if c, ok := f.(io.Closer); ok {
		// 	c.Close()
		// }

		rs = f
	}

	bb, err := bitio.NewBufferFromReadSeeker(rs)
	if err != nil {
		return err
	}

	return &bitBufFile{
		bb:       bb,
		filename: filename,
	}
}

func (q *Query) makeDumpFn(fnOpts decode.DumpOptions) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		v, err := toValue(c)
		if err != nil {
			return fmt.Errorf("%v: value is not a decode value", c)
		}
		opts := q.opts.DumpOptions
		opts.MaxDepth = fnOpts.MaxDepth
		opts.Verbose = fnOpts.Verbose

		return func(stdout io.Writer) error {
			if err := v.Dump(stdout, opts); err != nil {
				return err
			}
			return nil
		}
	}
}

func (q *Query) makeProbeFn(registry *decode.Registry, probeFormats []*decode.Format) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		bb, r, filename, err := toBitBuf(c)
		if err != nil {
			return err
		}
		bb, err = bb.BitBufRange(r.Start, r.Len)
		if err != nil {
			return err
		}

		opts := map[string]interface{}{}

		name := "unnamed"
		if filename != "" {
			name = filename
		}

		if len(a) >= 1 {
			formatName, err := toString(a[0])
			if err != nil {
				return fmt.Errorf("%s: %w", formatName, err)
			}
			probeFormats, err = registry.Group(formatName)
			if err != nil {
				return fmt.Errorf("%s: %w", formatName, err)
			}
		}

		dv, _, errs := decode.Probe(name, bb, probeFormats, decode.ProbeOptions{FormatOptions: opts})
		if dv == nil {
			return errs
		}

		return dv
	}
}

func (q *Query) _valueKeys(c interface{}, a []interface{}) interface{} {
	if v, ok := c.(*decode.Value); ok {
		var vs []interface{}
		for _, s := range v.SpecialPropNames() {
			vs = append(vs, s)
		}
		return vs
	}
	return nil
}

func (q *Query) bits(c interface{}, a []interface{}) interface{} {
	bb, r, _, err := toBitBuf(c)
	if err != nil {
		return err
	}
	bb, err = bb.BitBufRange(r.Start, r.Len)
	if err != nil {
		return err
	}

	startArg := int64(0)
	endArg := int64(-1)
	toAbs := func(v int64, l int64) int64 {
		if v < 0 {
			return l + v + 1
		}
		return v
	}

	if len(a) >= 1 {
		startArg, err = toInt64(a[0])
		if err != nil {
			return err
		}
	}
	if len(a) >= 2 {
		endArg, err = toInt64(a[1])
		if err != nil {
			return err
		}
	}

	startArg = toAbs(startArg, bb.Len())
	endArg = toAbs(endArg, bb.Len())

	bb, err = bb.BitBufRange(startArg, endArg-startArg)
	if err != nil {
		return err
	}

	return bb
}

func (q *Query) string_(c interface{}, a []interface{}) interface{} {
	var bb *bitio.Buffer
	switch cc := c.(type) {
	case *decode.Value:
		var err error
		bb, err = cc.RootBitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
		if err != nil {
			return err
		}
	case *bitio.Buffer:
		bb = cc
	default:
		return fmt.Errorf("value is not a decode value or bit buffer")
	}

	sb := &strings.Builder{}
	if _, err := io.Copy(sb, bb); err != nil {
		return err
	}

	return string(sb.String())
}

func (q *Query) u(c interface{}, a []interface{}) interface{} {
	bb, r, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	nBits := r.Len
	if len(a) == 1 {
		n, err := toInt64(a[0])
		if err != nil {
			return err
		}
		nBits = n
	}

	bb, err = bb.BitBufRange(r.Start, nBits)
	if err != nil {
		return err
	}

	// TODO: smart and maybe use int if bits can fit?
	bi := new(big.Int)
	for i := bb.Len() - 1; i >= 0; i-- {
		v, err := bb.Bool()
		if err != nil {
			return err
		}
		if v {
			bi.SetBit(bi, int(i), 1)
		}
	}

	return bi
}

func (q *Query) push(c interface{}, a []interface{}) interface{} {
	if _, ok := c.(error); !ok {
		q.runContext.pushVs = append(q.runContext.pushVs, c)
	}
	return func(stdout io.Writer) error {
		// nop
		return nil
	}

}

func (q *Query) pop(c interface{}, a []interface{}) interface{} {
	q.runContext.pops++
	return func(stdout io.Writer) error {
		// nop
		return nil
	}
}

func (q *Query) md5(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	md5 := md5.New()
	if _, err := io.Copy(md5, bb); err != nil {
		return err
	}

	return md5.Sum(nil)
}

func (q *Query) base64(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	b64Buf := &bytes.Buffer{}
	b64 := base64.NewEncoder(base64.StdEncoding, b64Buf)
	if _, err := io.Copy(b64Buf, bb); err != nil {
		return err
	}
	b64.Close()

	return b64Buf.Bytes()
}

func (q *Query) unbase64(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, base64.NewDecoder(base64.StdEncoding, bb)); err != nil {
		return err
	}

	return buf.Bytes()
}

func (q *Query) hex(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(hex.NewEncoder(buf), bb); err != nil {
		return err
	}

	return buf.String()
}

func (q *Query) unhex(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	b64Buf := &bytes.Buffer{}
	if _, err := io.Copy(b64Buf, hex.NewDecoder(bb)); err != nil {
		return err
	}

	return b64Buf.Bytes()
}

func (q *Query) queryEscape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	return url.QueryEscape(s)
}

func (q *Query) queryUnescape(c interface{}, a []interface{}) interface{} {
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
func (q *Query) pathEscape(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	return url.PathEscape(s)
}

func (q *Query) pathUnescape(c interface{}, a []interface{}) interface{} {
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

func (q *Query) aesCtr(c interface{}, a []interface{}) interface{} {
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

	bb, _, _, err := toBitBuf(c)
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
