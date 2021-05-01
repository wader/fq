package interp

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"fq/internal/aheadreadseeker"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/colorjson"
	"fq/internal/ctxreadseeker"
	"fq/internal/hexdump"
	"fq/internal/hexpairwriter"
	"fq/internal/ioextra"
	"fq/internal/num"
	"fq/internal/progressreadseeker"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"fq/pkg/ranges"
	"io"
	"io/ioutil"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/itchyny/gojq"
)

// TODO: make it nicer somehow? generate generators? remove from struct?
func (i *Interp) makeFunctions(registry *decode.Registry) []Function {
	fs := []Function{
		{[]string{"tty"}, 0, 0, i.tty, false},

		{[]string{"read"}, 0, 2, i.read, false},
		{[]string{"eval"}, 1, 1, i.eval, true},
		{[]string{"print"}, 0, 0, i.print, true},

		{[]string{"complete_query"}, 0, 0, i.completeQuery, false},
		{[]string{"display_name"}, 0, 0, i.displayName, false},
		{[]string{"_value_keys"}, 0, 0, i._valueKeys, false},
		{[]string{"formats"}, 0, 0, i.formats, false},

		{[]string{"open"}, 0, 1, i._open, false},
		{[]string{"decode"}, 0, 1, i.makeDecodeFn(registry, registry.MustGroup(format.PROBE)), false},

		{[]string{"display", "d"}, 0, 1, i.makeDisplayFn(nil), true},
		{[]string{"verbose", "v"}, 0, 1, i.makeDisplayFn(map[string]interface{}{"verbose": true}), true},
		{[]string{"preview", "p"}, 0, 1, i.preview, true},
		{[]string{"hexdump", "hd", "h"}, 0, 1, i.hexdump, true},

		{[]string{"string"}, 0, 0, i.string_, false},
		{[]string{"bytes"}, 0, 0, i.bytes, false},
		{[]string{"bits"}, 0, 0, i.bits, false},
		{[]string{"tovalue"}, 0, 0, i.tovalue, false},

		{[]string{"u"}, 0, 1, i.u, false},

		{[]string{"md5"}, 0, 0, i.md5, false},

		{[]string{"hex"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return hex.NewDecoder(r), nil },
			func(r io.Writer) (io.Writer, error) { return hex.NewEncoder(r), nil },
		), false},

		{[]string{"base64"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.StdEncoding, r), nil },
			func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.StdEncoding, r), nil },
		), false},
		{[]string{"rawbase64"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.RawURLEncoding, r), nil },
			func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.RawURLEncoding, r), nil },
		), false},

		{[]string{"urlbase64"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.URLEncoding, r), nil },
			func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.URLEncoding, r), nil },
		), false},

		{[]string{"nal_unescape"}, 0, 0, i.nalUnescape(), false},

		{[]string{"query_escape"}, 0, 0, i.queryEscape, false},
		{[]string{"query_unescape"}, 0, 0, i.queryUnescape, false},
		{[]string{"path_escape"}, 0, 0, i.pathEscape, false},
		{[]string{"path_unescape"}, 0, 0, i.pathUnescape, false},
		{[]string{"aes_ctr"}, 1, 2, i.aesCtr, false},
		{[]string{"json"}, 0, 0, i._json, false},

		{[]string{"_state"}, 1, 2, i._state, false},
		{[]string{"options"}, 0, 0, i.options, false},
	}
	for name, f := range i.registry.Groups {
		fs = append(fs, Function{[]string{name}, 0, 0, i.makeDecodeFn(registry, f), false})
	}

	return fs
}

func (i *Interp) tty(c interface{}, a []interface{}) interface{} {
	w, h := i.stdout.Size()
	return map[string]interface{}{
		"is_terminal": i.stdout.IsTerminal(),
		"size":        []interface{}{w, h},
	}
}

func makeStringBitBufTransformFn(
	decodeFn func(r io.Reader) (io.Reader, error),
	encodeFn func(w io.Writer) (io.Writer, error),
) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		switch c := c.(type) {
		case string:
			bb, _, err := toBitBuf(c)
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

			return &bitBufObject{bb: outBB, unit: 8, r: ranges.Range{Len: outBB.Len()}}
		default:
			bb, _, err := toBitBuf(c)
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

func (i *Interp) read(c interface{}, a []interface{}) interface{} {
	var ok bool
	completeFn := ""
	prompt := ""

	if len(a) > 0 {
		prompt, ok = a[0].(string)
		if !ok {
			return fmt.Errorf("%v: prompt is not a string", a[1])
		}
	}
	if len(a) > 1 {
		completeFn, ok = a[1].(string)
		if !ok {
			return fmt.Errorf("%v: complete function name is not a string", a[0])
		}
	}

	src, err := i.os.Readline(prompt, func(line string, pos int) (newLine []string, shared int) {
		completeCtx, completeCtxCancelFn := context.WithTimeout(i.ctx, 1*time.Second)
		defer completeCtxCancelFn()
		// TODO: err
		names, shared, _ := completeTrampoline(completeCtx, completeFn, c, i, string(line), pos)
		return names, shared
	})

	if err == ErrInterrupt {
		return valueErr{"interrupt"}
	} else if err == ErrEOF {
		return valueErr{"eof"}
	} else if err != nil {
		return err
	}

	return src
}

func (i *Interp) eval(c interface{}, a []interface{}) interface{} {
	src, ok := a[0].(string)
	if !ok {
		return fmt.Errorf("%v: src is not a string", a[0])
	}
	iter, err := i.Eval(i.ctx, ScriptMode, c, src, i.stdout)
	if err != nil {
		return err
	}

	return iter
}

func (i *Interp) print(c interface{}, a []interface{}) interface{} {
	if _, err := fmt.Fprintln(i.stdout, c); err != nil {
		return err
	}
	return gojq.EmptyIter{}
}

func (i *Interp) completeQuery(c interface{}, a []interface{}) interface{} {
	s, ok := c.(string)
	if !ok {
		return fmt.Errorf("%v: value is not a string", c)
	}

	gq, typ, prefix := BuildCompletionQuery(s)
	queryStr := ""
	if gq != nil {
		queryStr = gq.String()
	}

	return map[string]interface{}{
		"query":  queryStr,
		"type":   string(typ),
		"prefix": prefix,
	}
}

func (i *Interp) displayName(c interface{}, a []interface{}) interface{} {
	qo, ok := c.(InterpObject)
	if !ok {
		return fmt.Errorf("%v: value is not query object", c)
	}
	return qo.DisplayName()
}

func (i *Interp) _valueKeys(c interface{}, a []interface{}) interface{} {
	if v, ok := c.(InterpObject); ok {
		var vs []interface{}
		for _, s := range v.SpecialPropNames() {
			vs = append(vs, s)
		}
		return vs
	}
	return nil
}

func (i *Interp) formats(c interface{}, a []interface{}) interface{} {

	allFormats := map[string]*decode.Format{}

	for _, fs := range i.registry.Groups {
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

type bitBufFile struct {
	bb       *bitio.Buffer
	filename string

	decodeDoneFn func()
}

var _ ToBitBuf = (*bitBufFile)(nil)

func (bbf *bitBufFile) Display(w io.Writer, opts Options) error {
	_, err := fmt.Fprintf(w, "<%s>\n", bbf.filename)
	return err
}

func (bbf *bitBufFile) ToBitBuf() (*bitio.Buffer, ranges.Range) {
	return bbf.bb.Copy(), ranges.Range{Start: 0, Len: bbf.bb.Len()}
}

// def open($path): #:: [a]|(string) => buffer
// open read file from filesystem
func (i *Interp) _open(c interface{}, a []interface{}) interface{} {
	var rs io.ReadSeeker
	opts, err := i.Options()
	if err != nil {
		return err
	}

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
		buf, err := ioutil.ReadAll(ctxreadseeker.New(i.ctx, &ioextra.NopSeeker{Reader: i.os.Stdin()}))
		if err != nil {
			return err
		}
		rs = bytes.NewReader(buf)
	} else {
		f, err := i.os.Open(filename)
		if err != nil {
			return err
		}

		rs = f
	}

	ctxRs := ctxreadseeker.New(i.ctx, rs)
	// TODO: cleanup? bitbuf have optional close method etc?
	// TODO: can call ctxRs directory of need to forward close thru aheadreadseeker etc?
	// if c, ok := rs.(io.Closer); ok {
	// 	ctxRs.Close()
	// }

	bEnd, err := ioextra.SeekerEnd(ctxRs)
	if err != nil {
		return err
	}

	var progressRs io.ReadSeeker = ctxRs

	// TODO: make nicer
	// we don't want to print any progress things after decode is done
	var decodeDoneFn func()
	if opts.REPL {
		decodeDone := false
		progressFn := func(r, l int64) {
			if decodeDone {
				return
			}
			fmt.Fprintf(i.os.Stderr(), "\r%.1f%%", (float64(r)/float64(l))*100)
		}
		decodeDoneFn = func() {
			decodeDone = true
			// cleanup when done
			fmt.Fprint(i.os.Stderr(), "\r      \r")
		}
		const progressPrecision = 1024
		progressRs = progressreadseeker.New(ctxRs, progressPrecision, bEnd, progressFn)
	}

	const cacheReadAheadSize = 512 * 1024
	aheadRs := aheadreadseeker.New(progressRs, cacheReadAheadSize)

	// bb -> aheadreadseeker -> progressreadseeker -> ctxreadseeker -> readerseeker

	bb, err := bitio.NewBufferFromReadSeeker(aheadRs)
	if err != nil {
		return err
	}

	return &bitBufFile{
		bb:           bb,
		filename:     filename,
		decodeDoneFn: decodeDoneFn,
	}
}

func (i *Interp) makeDecodeFn(registry *decode.Registry, decodeFormats []*decode.Format) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		// TODO: progress hack
		// would be nice to move progress code into decode but it might be
		// tricky to keep track of absolute positions in the underlaying readers
		// when it uses BitBuf slices, maybe only in Pos()?
		if bbf, ok := c.(*bitBufFile); ok {
			if bbf.decodeDoneFn != nil {
				defer bbf.decodeDoneFn()
			}
		}

		bb, r, err := toBitBuf(c)
		if err != nil {
			return err
		}
		bb, err = bb.BitBufRange(r.Start, r.Len)
		if err != nil {
			return err
		}

		opts := map[string]interface{}{}

		// TODO:
		name := "unnamed"

		if len(a) >= 1 {
			formatName, err := toString(a[0])
			if err != nil {
				return fmt.Errorf("%s: %w", formatName, err)
			}
			decodeFormats, err = registry.Group(formatName)
			if err != nil {
				return fmt.Errorf("%s: %w", formatName, err)
			}
		}

		dv, _, err := decode.Decode(i.ctx, name, bb, decodeFormats, decode.DecodeOptions{FormatOptions: opts})
		if dv == nil {
			switch err := err.(type) {
			case decode.DecodeFormatsError:
				return decodeError{err}
			}
			return valueErr{err}
		}

		return valueObject{v: dv}
	}
}

func (i *Interp) makeDisplayFn(fnOpts map[string]interface{}) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		opts, err := i.Options(append([]interface{}{fnOpts}, a...)...)
		if err != nil {
			return err
		}

		switch v := c.(type) {
		case Display:
			if err := v.Display(i.stdout, opts); err != nil {
				return err
			}
			return []interface{}{}
		case nil, bool, float64, int, string, *big.Int, map[string]interface{}, []interface{}, gojq.JQValue:
			if s, ok := v.(string); ok && opts.RawString {
				fmt.Fprintln(i.stdout, s)
				return []interface{}{}
			}

			if err := colorjson.NewEncoder(
				opts.Color, false, 2,
				func(v interface{}) interface{} {
					if o, ok := v.(gojq.JQValue); ok {
						return o.JQValue()
					}
					return nil
				},
				colorjson.Colors{
					Reset:     []byte(ansi.Reset),
					Null:      []byte(opts.Decorator.Null),
					False:     []byte(opts.Decorator.False),
					True:      []byte(opts.Decorator.True),
					Number:    []byte(opts.Decorator.Number),
					String:    []byte(opts.Decorator.String),
					ObjectKey: []byte(opts.Decorator.ObjectKey),
					Array:     []byte(opts.Decorator.Array),
					Object:    []byte(opts.Decorator.Object),
				},
			).Marshal(v, i.stdout); err != nil {
				return err
			}
			fmt.Fprintln(i.stdout)

			return gojq.EmptyIter{}
		case error:
			return v
		default:
			return fmt.Errorf("%+#v: not displayable", c)
		}
	}
}

// TODO: opts and colors?
func (i *Interp) preview(c interface{}, a []interface{}) interface{} {
	opts, err := i.Options(a...)
	if err != nil {
		return err
	}

	switch v := c.(type) {
	case Preview:
		if err := v.Preview(i.stdout, opts); err != nil {
			return err
		}
		return gojq.EmptyIter{}
	default:
		return fmt.Errorf("%v: not previewable", c)
	}
}

func (i *Interp) hexdump(c interface{}, a []interface{}) interface{} {
	bb, r, err := toBitBuf(c)
	if err != nil {
		return err
	}

	bitsByteAlign := r.Start % 8
	bb, err = bb.BitBufRange(r.Start-bitsByteAlign, r.Len+bitsByteAlign)
	if err != nil {
		return err
	}

	opts, err := i.Options(a...)
	if err != nil {
		return err
	}

	d := opts.Decorator
	hw := hexdump.New(
		i.stdout,
		(r.Start-bitsByteAlign)/8,
		num.DigitsInBase(bitio.BitsByteCount(r.Stop()+bitsByteAlign), true, opts.AddrBase),
		opts.AddrBase,
		opts.LineBytes,
		func(b byte) string { return d.ByteColor(b).Wrap(hexpairwriter.Pair(b)) },
		func(b byte) string { return d.ByteColor(b).Wrap(asciiwriter.SafeASCII(b)) },
		func(s string) string { return d.Frame.Wrap(s) },
		d.Column,
	)
	defer hw.Close()
	if _, err = io.Copy(hw, bb); err != nil {
		return err
	}

	return gojq.EmptyIter{}
}

func (i *Interp) string_(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	sb := &strings.Builder{}
	if _, err := io.Copy(sb, bb); err != nil {
		return err
	}

	return string(sb.String())
}

func (i *Interp) bytes(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	return &bitBufObject{bb: bb, unit: 8, r: ranges.Range{Len: bb.Len()}}
}

func (i *Interp) bits(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	return &bitBufObject{bb: bb, unit: 1, r: ranges.Range{Len: bb.Len()}}
}

func (i *Interp) tovalue(c interface{}, a []interface{}) interface{} {
	return toValue(c)
}

func (i *Interp) u(c interface{}, a []interface{}) interface{} {
	bb, r, err := toBitBuf(c)
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

func (i *Interp) md5(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	md5 := md5.New()
	if _, err := io.Copy(md5, bb); err != nil {
		return err
	}

	return md5.Sum(nil)
}

func makeBitBufTransformFn(fn func(r io.Reader) (io.Reader, error)) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		inBB, _, err := toBitBuf(c)
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

		return &bitBufObject{bb: outBB, unit: 8, r: ranges.Range{Len: outBB.Len()}}
	}
}

func (i *Interp) nalUnescape() func(c interface{}, a []interface{}) interface{} {
	return makeBitBufTransformFn(func(r io.Reader) (io.Reader, error) {
		return &decode.NALUnescapeReader{Reader: r}, nil
	})
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

	bb, _, err := toBitBuf(c)
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

func (i *Interp) _json(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBitBuf(c)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bb); err != nil {
		return err
	}

	var vv interface{}
	if err := json.Unmarshal(buf.Bytes(), &vv); err != nil {
		return err
	}

	return vv
}

func (i *Interp) _state(c interface{}, a []interface{}) interface{} {
	name, ok := a[0].(string)
	if !ok {
		return fmt.Errorf("%v: value is not a string", c)
	}

	s := i.state[name]
	if len(a) > 1 {
		i.state[name] = a[1]
	}

	return s
}

func (i *Interp) options(c interface{}, a []interface{}) interface{} {
	opts, err := i.Options()
	if err != nil {
		return err
	}

	// TODO: ugly
	var v map[string]interface{}
	b := &bytes.Buffer{}
	if err := json.NewEncoder(b).Encode(&opts); err != nil {
		return err
	}
	if err := json.NewDecoder(b).Decode(&v); err != nil {
		return err
	}

	return v
}
