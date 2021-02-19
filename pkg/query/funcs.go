package query

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/hexdump"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"fq/internal/progressreadseeker"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"io"
	"io/ioutil"
	"math/big"
	"net/url"
	"strings"
)

const builtinPrefix = "@builtin"

//go:embed *.jq
var builtinFS embed.FS

func buildDisplayOptions(ms ...map[string]interface{}) DisplayOptions {
	var opts DisplayOptions
	for _, m := range ms {
		if m != nil {
			mapSetDisplayOptions(&opts, m)
		}
	}
	opts.Decorator = decoratorFromDumpOptions(opts)

	return opts
}

func mapSetDisplayOptions(d *DisplayOptions, m map[string]interface{}) {
	if v, ok := m["maxdepth"]; ok {
		d.MaxDepth = num.MaxInt(0, toIntZ(v))
	}
	if v, ok := m["verbose"]; ok {
		d.Verbose = toBoolZ(v)
	}
	if v, ok := m["color"]; ok {
		d.Color = toBoolZ(v)
	}
	if v, ok := m["unicode"]; ok {
		d.Unicode = toBoolZ(v)
	}
	if v, ok := m["raw"]; ok {
		d.Raw = toBoolZ(v)
	}
	if v, ok := m["linebytes"]; ok {
		d.LineBytes = num.MaxInt(0, toIntZ(v))
	}
	if v, ok := m["displaybytes"]; ok {
		d.DisplayBytes = num.MaxInt64(0, toInt64Z(v))
	}
	if v, ok := m["addrbase"]; ok {
		d.AddrBase = num.ClampInt(2, 36, toIntZ(v))
	}
	if v, ok := m["sizebase"]; ok {
		d.SizeBase = num.ClampInt(2, 36, toIntZ(v))
	}
}

func decoratorFromDumpOptions(opts DisplayOptions) Decorator {
	colStr := "|"
	if opts.Unicode {
		colStr = "\xe2\x94\x82"
	}
	nameFn := func(s string) string { return s }
	valueFn := func(s string) string { return s }
	byteFn := func(b byte, s string) string { return s }
	column := colStr + "\n"
	if opts.Color {
		nameFn = func(s string) string { return ansi.FgBrightBlue + s + ansi.Reset }
		valueFn = func(s string) string { return ansi.FgBrightCyan + s + ansi.Reset }
		byteFn = func(b byte, s string) string {
			switch {
			case b == 0:
				return ansi.FgBrightBlack + s + ansi.Reset
			case b >= 32 && b <= 126, b == '\r', b == '\n', b == '\f', b == '\t', b == '\v':
				return ansi.FgWhite + s + ansi.Reset
			default:
				return ansi.FgBrightWhite + s + ansi.Reset
			}
		}
		column = ansi.FgWhite + colStr + ansi.Reset + "\n"
	}

	return Decorator{
		Name:   nameFn,
		Value:  valueFn,
		Byte:   byteFn,
		Column: column,
	}
}

type Decorators struct {
	Name   func(s string) string
	Value  func(s string) string
	Byte   func(b byte, s string) string
	Column string
}

// TODO: make it nicer somehow?
func (q *Query) makeFunctions(opts QueryOptions) []Function {
	fs := []Function{
		{[]string{"tty"}, 0, 0, q.tty},
		{[]string{"options"}, 0, 1, q.options},

		{[]string{"help"}, 0, 0, q.help},
		{[]string{"open"}, 0, 1, q.open},
		{[]string{"display", "d"}, 0, 1, q.makeDisplayFn(nil)},
		{[]string{"verbose", "v"}, 0, 1, q.makeDisplayFn(map[string]interface{}{"verbose": true})},
		{[]string{"hexdump", "hd", "h"}, 0, 1, q.hexdump},
		{[]string{"string"}, 0, 0, q.string_},
		{[]string{"decode"}, 0, 1, q.makeDecodeFn(opts.Registry, opts.Registry.MustGroup(format.PROBE))},
		{[]string{"u"}, 0, 1, q.u},
		{[]string{"push"}, 0, 0, q.push},
		{[]string{"pop"}, 0, 0, q.pop},
		{[]string{"_value_keys"}, 0, 0, q._valueKeys},
		{[]string{"formats"}, 0, 0, q.formats},
		{[]string{"preview", "p"}, 0, 0, q.preview},
		{[]string{"md5"}, 0, 0, q.md5},
		{[]string{"base64"}, 0, 0, q.base64},
		{[]string{"unbase64"}, 0, 0, q.unbase64},
		{[]string{"hex"}, 0, 0, q.hex},
		{[]string{"unhex"}, 0, 0, q.unhex},
		{[]string{"query_escape"}, 0, 0, q.queryEscape},
		{[]string{"query_unescape"}, 0, 0, q.queryUnescape},
		{[]string{"path_escape"}, 0, 0, q.pathEscape},
		{[]string{"path_unescape"}, 0, 0, q.pathUnescape},
		{[]string{"aes_ctr"}, 1, 2, q.aesCtr},

		{[]string{"json"}, 0, 0, q._json},
	}
	for name, f := range q.opts.Registry.Groups {
		fs = append(fs, Function{[]string{name}, 0, 0, q.makeDecodeFn(opts.Registry, f)})
	}

	return fs
}

func (q *Query) tty(c interface{}, a []interface{}) interface{} {
	w, h := q.runContext.stdout.Size()
	return map[string]interface{}{
		"is_terminal": q.runContext.stdout.IsTerminal(),
		"size":        []interface{}{w, h},
	}
}

func (q *Query) options(c interface{}, a []interface{}) interface{} {
	if len(a) > 0 {
		opts, ok := a[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("%v: value is not object", a[0])
		}
		q.runContext.opts = opts
	}
	return q.runContext.opts
}

func (q *Query) _json(c interface{}, a []interface{}) interface{} {
	bb, _, _, err := toBitBuf(c)
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

		var opts DisplayOptions
		if len(a) >= 1 {
			opts = buildDisplayOptions(q.runContext.opts, a[0].(map[string]interface{}))
		} else {
			opts = buildDisplayOptions(q.runContext.opts)
		}

		d := opts.Decorator
		hw := hexdump.New(
			stdout,
			(r.Start-bitsByteAlign)/8,
			num.DigitsInBase(bitio.BitsByteCount(r.Stop()+bitsByteAlign), true, opts.AddrBase),
			opts.AddrBase,
			opts.LineBytes,
			func(b byte) string { return d.Byte(b, hexpairwriter.Pair(b)) },
			func(b byte) string { return d.Byte(b, asciiwriter.SafeASCII(b)) },
			d.Column,
		)
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
	vo, ok := c.(valueObject)
	if !ok {
		return fmt.Errorf("%v: value is not a decode value", c)
	}
	return func(stdout io.Writer) error {
		if err := preview(vo.v, stdout); err != nil {
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
		fmt.Fprintf(stdout, "^D to exit\n")
		fmt.Fprintf(stdout, "^C to interrupt\n")
		return nil
	})
}

type bitBufFile struct {
	bb       *bitio.Buffer
	filename string

	decodeDoneFn func()
}

func (bbf *bitBufFile) Display(w io.Writer, opts DisplayOptions) error {
	_, err := fmt.Fprintf(w, "<%s>\n", bbf.filename)
	return err
}

func (bbf *bitBufFile) ToBifBuf() *bitio.Buffer {
	return bbf.bb.Copy()
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

	//TODO: how to know when decode is done?
	// TODO: refactor
	bPos, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	bEnd, err := rs.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	if _, err := rs.Seek(bPos, io.SeekStart); err != nil {
		return err
	}

	// TODO: make nicer
	// we don't want to print any progress things after decode is done
	var decodeDoneFn func()
	if q.runContext.mode == REPLMode {
		decodeDone := false
		decodeDoneFn = func() {
			fmt.Fprint(q.runContext.stdout, "\r")
			decodeDone = true
		}

		rs = progressreadseeker.New(rs, bEnd, func(readBytes int64, length int64) {
			if decodeDone {
				return
			}
			fmt.Fprintf(q.runContext.stdout, "\r%.1f%%", (float64(readBytes)/float64(length))*100)
		})
	}

	bb, err := bitio.NewBufferFromReadSeeker(rs)
	if err != nil {
		return err
	}

	return &bitBufFile{
		bb:           bb,
		filename:     filename,
		decodeDoneFn: decodeDoneFn,
	}
}

func (q *Query) makeDisplayFn(fnOpts map[string]interface{}) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		return func(stdout io.Writer) error {
			var opts DisplayOptions
			if len(a) >= 1 {
				opts = buildDisplayOptions(q.runContext.opts, fnOpts, a[0].(map[string]interface{}))
			} else {
				opts = buildDisplayOptions(q.runContext.opts, fnOpts)
			}

			switch v := c.(type) {
			case Display:
				if err := v.Display(stdout, opts); err != nil {
					return err
				}
			default:
				return fmt.Errorf("%v: not displayable", c)
			}
			return nil
		}
	}
}

func (q *Query) makeDecodeFn(registry *decode.Registry, decodeFormats []*decode.Format) func(c interface{}, a []interface{}) interface{} {
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
			decodeFormats, err = registry.Group(formatName)
			if err != nil {
				return fmt.Errorf("%s: %w", formatName, err)
			}
		}

		dv, _, errs := decode.Decode(name, bb, decodeFormats, decode.DecodeOptions{FormatOptions: opts})
		if dv == nil {
			return errs
		}

		return valueObject{v: dv}
	}
}

func (q *Query) _valueKeys(c interface{}, a []interface{}) interface{} {
	if v, ok := c.(valueObject); ok {
		var vs []interface{}
		for _, s := range v.SpecialPropNames() {
			vs = append(vs, s)
		}
		return vs
	}
	return nil
}

func (q *Query) string_(c interface{}, a []interface{}) interface{} {
	var bb *bitio.Buffer
	switch cc := c.(type) {
	case valueObject:
		var err error
		bb, err = cc.v.RootBitBuf.BitBufRange(cc.v.Range.Start, cc.v.Range.Len)
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
