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
	"errors"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"math/big"
	"net/url"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/wader/fq/internal/aheadreadseeker"
	"github.com/wader/fq/internal/ctxreadseeker"
	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/internal/ioextra"
	"github.com/wader/fq/internal/progressreadseeker"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"

	"github.com/wader/gojq"
)

// TODO: make it nicer somehow? generate generators? remove from struct?
func (i *Interp) makeFunctions() []Function {
	fs := []Function{
		{[]string{"_readline"}, 0, 2, i.readline, nil},
		{[]string{"eval"}, 1, 2, nil, i.eval},
		{[]string{"stdin"}, 0, 0, nil, i.makeStdioFn(i.os.Stdin())},
		{[]string{"stdout"}, 0, 0, nil, i.makeStdioFn(i.os.Stdout())},
		{[]string{"stderr"}, 0, 0, nil, i.makeStdioFn(i.os.Stderr())},

		{[]string{"_query_fromstring"}, 0, 0, i.queryFromString, nil},
		{[]string{"_query_tostring"}, 0, 0, i.queryToString, nil},

		{[]string{"_display_name"}, 0, 0, i._displayName, nil},
		{[]string{"_extkeys"}, 0, 0, i._extKeys, nil},
		{[]string{"_global_state"}, 0, 1, i.makeStateFn(i.state), nil},

		{[]string{"_registry"}, 0, 0, i._registry, nil},
		{[]string{"history"}, 0, 0, i.history, nil},

		{[]string{"open"}, 0, 0, i._open, nil},
		{[]string{"_decode"}, 2, 2, i._decode, nil},

		{[]string{"format"}, 0, 0, i.format, nil},
		{[]string{"_display"}, 1, 1, nil, i._display},
		{[]string{"preview", "p"}, 0, 1, nil, i.preview},
		{[]string{"hexdump", "hd", "h"}, 0, 1, nil, i.hexdump},

		{[]string{"tobytes"}, 0, 0, i.bytes, nil},
		{[]string{"tobits"}, 0, 0, i.bits, nil},
		{[]string{"tovalue"}, 0, 1, i.tovalue, nil},

		{[]string{"hex"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return hex.NewDecoder(r), nil },
			func(r io.Writer) (io.Writer, error) { return hex.NewEncoder(r), nil },
		), nil},

		{[]string{"base64"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.StdEncoding, r), nil },
			func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.StdEncoding, r), nil },
		), nil},
		{[]string{"rawbase64"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.RawURLEncoding, r), nil },
			func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.RawURLEncoding, r), nil },
		), nil},

		{[]string{"urlbase64"}, 0, 0, makeStringBitBufTransformFn(
			func(r io.Reader) (io.Reader, error) { return base64.NewDecoder(base64.URLEncoding, r), nil },
			func(r io.Writer) (io.Writer, error) { return base64.NewEncoder(base64.URLEncoding, r), nil },
		), nil},

		{[]string{"nal_unescape"}, 0, 0, makeBitBufTransformFn(func(r io.Reader) (io.Reader, error) {
			return &decode.NALUnescapeReader{Reader: r}, nil
		}), nil},

		{[]string{"md5"}, 0, 0, makeHashFn(func() (hash.Hash, error) { return md5.New(), nil }), nil},

		{[]string{"query_escape"}, 0, 0, i.queryEscape, nil},
		{[]string{"query_unescape"}, 0, 0, i.queryUnescape, nil},
		{[]string{"path_escape"}, 0, 0, i.pathEscape, nil},
		{[]string{"path_unescape"}, 0, 0, i.pathUnescape, nil},
		{[]string{"aes_ctr"}, 1, 2, i.aesCtr, nil},

		{[]string{"find"}, 1, 1, nil, i.find},
	}

	return fs
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

			return newBifBufObject(outBB, 8)
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

		return newBifBufObject(outBB, 8)
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

		return newBifBufObject(outBB, 8)
	}
}

func (i *Interp) readline(c interface{}, a []interface{}) interface{} {
	var opts struct {
		Complete string  `mapstructure:"complete"`
		Timeout  float64 `mapstructure:"timeout"`
	}

	var err error
	prompt := ""

	if len(a) > 0 {
		prompt, err = toString(a[0])
		if err != nil {
			return fmt.Errorf("prompt: %w", err)
		}
	}
	if len(a) > 1 {
		_ = mapstructure.Decode(a[1], &opts)
	}

	src, err := i.os.Readline(
		prompt,
		func(line string, pos int) (newLine []string, shared int) {
			completeCtx := i.evalContext.ctx
			if opts.Timeout > 0 {
				var completeCtxCancelFn context.CancelFunc
				completeCtx, completeCtxCancelFn = context.WithTimeout(i.evalContext.ctx, time.Duration(opts.Timeout*float64(time.Second)))
				defer completeCtxCancelFn()
			}

			names, shared, err := func() (newLine []string, shared int, err error) {
				vs, err := i.EvalFuncValues(
					completeCtx, c, opts.Complete, []interface{}{line, pos}, DiscardCtxWriter{Ctx: completeCtx},
				)
				if err != nil {
					return nil, pos, err
				}
				if len(vs) < 1 {
					return nil, pos, fmt.Errorf("no values")
				}
				v := vs[0]
				if vErr, ok := v.(error); ok {
					return nil, pos, vErr
				}

				// {abc: 123, abd: 123} | complete(".ab"; 3) will return {prefix: "ab", names: ["abc", "abd"]}

				var result struct {
					Names  []string `mapstructure:"names"`
					Prefix string   `mapstructure:"prefix"`
				}

				_ = mapstructure.Decode(v, &result)
				if len(result.Names) == 0 {
					return nil, pos, nil
				}

				sharedLen := len(result.Prefix)

				return result.Names, sharedLen, nil
			}()

			// TODO: how to report err?
			_ = err

			return names, shared
		},
	)

	if errors.Is(err, ErrInterrupt) {
		return valueError{"interrupt"}
	} else if errors.Is(err, ErrEOF) {
		return valueError{"eof"}
	} else if err != nil {
		return err
	}

	return src
}

func (i *Interp) eval(c interface{}, a []interface{}) gojq.Iter {
	var err error
	src, err := toString(a[0])
	if err != nil {
		return gojq.NewIter(fmt.Errorf("src: %w", err))
	}
	var filenameHint string
	if len(a) >= 2 {
		filenameHint, err = toString(a[1])
		if err != nil {
			return gojq.NewIter(fmt.Errorf("filename hint: %w", err))
		}
	}

	iter, err := i.Eval(i.evalContext.ctx, c, src, filenameHint, i.evalContext.output)
	if err != nil {
		return gojq.NewIter(err)
	}

	return iter
}

func (i *Interp) makeStdioFn(t Terminal) func(c interface{}, a []interface{}) gojq.Iter {
	return func(c interface{}, a []interface{}) gojq.Iter {
		if c == nil {
			w, h := t.Size()
			return gojq.NewIter(map[string]interface{}{
				"is_terminal": t.IsTerminal(),
				"width":       w,
				"height":      h,
			})
		}

		if w, ok := t.(io.Writer); ok {
			if _, err := fmt.Fprint(w, c); err != nil {
				return gojq.NewIter(err)
			}
			return gojq.NewIter()
		}

		return gojq.NewIter(fmt.Errorf("%v: it not writeable", c))
	}
}

func (i *Interp) queryFromString(c interface{}, a []interface{}) interface{} {
	s, err := toString(c)
	if err != nil {
		return err
	}
	q, err := gojq.Parse(s)
	if err != nil {
		p := queryErrorPosition(s, err)
		return compileError{
			err:  err,
			what: "parse",
			pos:  p,
		}
	}

	// TODO: use mapstruct?
	b, err := json.Marshal(q)
	if err != nil {
		return err
	}
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	return v

}

func (i *Interp) queryToString(c interface{}, a []interface{}) interface{} {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	var q gojq.Query
	if err := json.Unmarshal(b, &q); err != nil {
		return err
	}

	return q.String()
}

func (i *Interp) _displayName(c interface{}, a []interface{}) interface{} {
	qo, ok := c.(Value)
	if !ok {
		return fmt.Errorf("%v: value is not query object", c)
	}
	return qo.DisplayName()
}

func (i *Interp) _extKeys(c interface{}, a []interface{}) interface{} {
	if v, ok := c.(Value); ok {
		var vs []interface{}
		for _, s := range v.ExtKeys() {
			vs = append(vs, s)
		}
		return vs
	}
	return nil
}

func (i *Interp) makeStateFn(state *interface{}) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		if len(a) > 0 {
			*state = a[0]
		}
		return *state
	}
}

func (i *Interp) _registry(c interface{}, a []interface{}) interface{} {
	uniqueFormats := map[string]*decode.Format{}

	groups := map[string]interface{}{}
	formats := map[string]interface{}{}

	for fsName, fs := range i.registry.Groups {
		var group []interface{}

		for _, f := range fs {
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

func (i *Interp) history(c interface{}, a []interface{}) interface{} {
	hs, err := i.os.History()
	if err != nil {
		return err
	}
	var vs []interface{}
	for _, s := range hs {
		vs = append(vs, s)
	}
	return vs
}

type bitBufFile struct {
	gojq.JQValue

	bb       *bitio.Buffer
	filename string

	progressFn progressreadseeker.ProgressFn
}

var _ ToBuffer = (*bitBufFile)(nil)

func (bbf *bitBufFile) Display(w io.Writer, opts Options) error {
	_, err := fmt.Fprintln(w, bbf.JQValue.JQValueToString())
	return err
}

func (bbf *bitBufFile) ToBuffer() (*bitio.Buffer, error) {
	return bbf.bb.Copy(), nil
}

// def open: #:: string| => buffer
// opens a file for reading from filesystem
// TODO: when to close? when bb loses all refs? need to use finalizer somehow?
func (i *Interp) _open(c interface{}, a []interface{}) interface{} {
	var err error

	var path string
	path, err = toString(c)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	var bEnd int64
	var f fs.File
	if path == "" || path == "-" {
		f = i.os.Stdin()
	} else {
		f, err = i.os.FS().Open(path)
		if err != nil {
			return err
		}
	}

	var fRS io.ReadSeeker
	fFI, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	// TODO: ctxreadseeker might leak
	if fFI.Mode().IsRegular() {
		if rs, ok := f.(io.ReadSeeker); ok {
			fRS = ctxreadseeker.New(i.evalContext.ctx, rs)
			bEnd = fFI.Size()
		}
	}

	if fRS == nil {
		buf, err := ioutil.ReadAll(ctxreadseeker.New(i.evalContext.ctx, &ioextra.ReadErrSeeker{Reader: f}))
		if err != nil {
			f.Close()
			return err
		}
		fRS = bytes.NewReader(buf)
		bEnd = int64(len(buf))
	}

	bbf := &bitBufFile{
		JQValue:  gojqextra.String(fmt.Sprintf("<bitBufFile %s>", path)),
		filename: path,
	}

	const progressPrecision = 1024
	fRS = progressreadseeker.New(fRS, progressPrecision, bEnd,
		func(approxReadBytes int64, totalSize int64) {
			// progressFn is assign by decode etc
			if bbf.progressFn != nil {
				bbf.progressFn(approxReadBytes, totalSize)
			}
		},
	)

	const cacheReadAheadSize = 512 * 1024
	aheadRs := aheadreadseeker.New(fRS, cacheReadAheadSize)

	// bb -> aheadreadseeker -> progressreadseeker -> ctxreadseeker -> readerseeker

	bbf.bb, err = bitio.NewBufferFromReadSeeker(aheadRs)
	if err != nil {
		return err
	}

	return bbf
}

func (i *Interp) _decode(c interface{}, a []interface{}) interface{} {
	var opts struct {
		Filename string                 `mapstructure:"filename"`
		Progress string                 `mapstructure:"_progress"`
		Remain   map[string]interface{} `mapstructure:",remain"`
	}
	_ = mapstructure.Decode(a[1], &opts)

	// TODO: progress hack
	// would be nice to move all progress code into decode but it might be
	// tricky to keep track of absolute positions in the underlaying readers
	// when it uses BitBuf slices, maybe only in Pos()?
	if bbf, ok := c.(*bitBufFile); ok {
		opts.Filename = bbf.filename

		if opts.Progress != "" {
			evalProgress := func(c interface{}) {
				_, _ = i.EvalFuncValues(
					i.evalContext.ctx,
					c,
					opts.Progress,
					nil,
					DiscardCtxWriter{Ctx: i.evalContext.ctx},
				)
			}
			bbf.progressFn = func(approxReadBytes, totalSize int64) {
				evalProgress(
					map[string]interface{}{
						"approx_read_bytes": approxReadBytes,
						"total_size":        totalSize,
					},
				)
			}
			defer evalProgress(nil)
		}
	}

	bb, err := toBuffer(c)
	if err != nil {
		return err
	}

	formatName, err := toString(a[0])
	if err != nil {
		return fmt.Errorf("%s: %w", formatName, err)
	}
	decodeFormats, err := i.registry.Group(formatName)
	if err != nil {
		return fmt.Errorf("%s: %w", formatName, err)
	}

	dv, _, err := decode.Decode(i.evalContext.ctx, bb, decodeFormats,
		decode.Options{
			Description:   opts.Filename,
			FormatOptions: opts.Remain,
		},
	)
	if dv == nil {
		var decodeFormatsErr decode.FormatsError
		if errors.As(err, &decodeFormatsErr) {
			var vs []interface{}
			for _, fe := range decodeFormatsErr.Errs {
				vs = append(vs, fe.Value())
			}

			return valueError{vs}
		}

		return valueError{err}
	}

	return makeDecodeValue(dv)
}

func (i *Interp) format(c interface{}, a []interface{}) interface{} {
	cj, ok := c.(gojq.JQValue)
	if !ok {
		return nil
	}
	f, ok := cj.JQValueKey("_format").(string)
	if !ok {
		return nil
	}
	return f
}

func (i *Interp) _display(c interface{}, a []interface{}) gojq.Iter {
	opts, err := i.Options(a...)
	if err != nil {
		return gojq.NewIter(err)
	}

	switch v := c.(type) {
	case Display:
		if err := v.Display(i.evalContext.output, opts); err != nil {
			return gojq.NewIter(err)
		}
		return gojq.NewIter()
	case nil, bool, float64, int, string, *big.Int, map[string]interface{}, []interface{}, gojq.JQValue:
		if s, ok := v.(string); ok && opts.RawString {
			fmt.Fprint(i.evalContext.output, s)
		} else {
			cj, err := i.NewColorJSON(opts)
			if err != nil {
				return gojq.NewIter(err)
			}
			if err := cj.Marshal(v, i.evalContext.output); err != nil {
				return gojq.NewIter(err)
			}
		}
		fmt.Fprint(i.evalContext.output, opts.JoinString)

		return gojq.NewIter()
	case error:
		return gojq.NewIter(v)
	default:
		return gojq.NewIter(fmt.Errorf("%+#v: not displayable", c))
	}
}

// TODO: opts and colors?
func (i *Interp) preview(c interface{}, a []interface{}) gojq.Iter {
	opts, err := i.Options(a...)
	if err != nil {
		return gojq.NewIter(err)
	}

	switch v := c.(type) {
	case Preview:
		if err := v.Preview(i.evalContext.output, opts); err != nil {
			return gojq.NewIter(err)
		}
		return gojq.NewIter()
	default:
		return gojq.NewIter(fmt.Errorf("%v: not previewable", c))
	}
}

func (i *Interp) hexdump(c interface{}, a []interface{}) gojq.Iter {
	bbr, err := toBufferRange(c)
	if err != nil {
		return gojq.NewIter(err)
	}

	opts, err := i.Options(a...)
	if err != nil {
		return gojq.NewIter(err)
	}

	if err := hexdumpRange(bbr, i.evalContext.output, opts); err != nil {
		return gojq.NewIter(err)
	}

	return gojq.NewIter()
}

func (i *Interp) bytes(c interface{}, a []interface{}) interface{} {
	bb, err := toBuffer(c)
	if err != nil {
		return err
	}
	return newBifBufObject(bb, 8)
}

func (i *Interp) bits(c interface{}, a []interface{}) interface{} {
	bb, err := toBuffer(c)
	if err != nil {
		return err
	}
	return newBifBufObject(bb, 1)
}

func (i *Interp) tovalue(c interface{}, a []interface{}) interface{} {
	opts, err := i.Options(append([]interface{}{}, a...)...)
	if err != nil {
		return err
	}
	v, _ := toValue(opts, c)
	return v
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

func (i *Interp) find(c interface{}, a []interface{}) gojq.Iter {
	bb, err := toBuffer(c)
	if err != nil {
		return gojq.NewIter(err)
	}

	sbb, err := toBuffer(a[0])
	if err != nil {
		return gojq.NewIter(err)
	}

	log.Printf("sbb: %#+v\n", sbb)

	// TODO: error, bitio.Copy?

	bbBytes := &bytes.Buffer{}
	_, _ = io.Copy(bbBytes, bb)

	sbbBytes := &bytes.Buffer{}
	_, _ = io.Copy(sbbBytes, sbb)

	// log.Printf("bbBytes.Bytes(): %#+v\n", bbBytes.Bytes())
	// log.Printf("sbbBytes.Bytes(): %#+v\n", sbbBytes.Bytes())

	idx := bytes.Index(bbBytes.Bytes(), sbbBytes.Bytes())
	if idx == -1 {
		return gojq.NewIter()
	}

	bbo := newBifBufObject(bb, 8)
	// log.Printf("bbo: %#+v\n", bbo)

	return gojq.NewIter(bbo)
}
