package interp

import (
	"bytes"
	"context"
	"crypto/md5"
	"embed"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/internal/ansi"
	"github.com/wader/fq/internal/colorjson"
	"github.com/wader/fq/internal/ctxstack"
	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/internal/pos"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"

	"github.com/wader/gojq"
)

//go:embed interp.jq
//go:embed internal.jq
//go:embed funcs.jq
//go:embed grep.jq
//go:embed options.jq
//go:embed args.jq
//go:embed query.jq
//go:embed repl.jq
//go:embed formats.jq
var builtinFS embed.FS

var initSource = `include "@builtin/interp";`

type valueError struct {
	v interface{}
}

func (v valueError) Error() string      { return fmt.Sprintf("error: %v", v.v) }
func (v valueError) Value() interface{} { return v.v }

type compileError struct {
	err      error
	what     string
	filename string
	pos      pos.Pos
}

func (ce compileError) Value() interface{} {
	return map[string]interface{}{
		"error":    ce.err.Error(),
		"what":     ce.what,
		"filename": ce.filename,
		"line":     ce.pos.Line,
		"column":   ce.pos.Column,
	}
}
func (ce compileError) Error() string {
	filename := ce.filename
	if filename == "" {
		filename = "src"
	}
	return fmt.Sprintf("%s:%d:%d: %s: %s", filename, ce.pos.Line, ce.pos.Column, ce.what, ce.err.Error())
}

var ErrEOF = io.EOF
var ErrInterrupt = errors.New("Interrupt")

// gojq errors can implement this to signal exit code
type Exiter interface {
	ExitCode() int
}

// gojq halt_error uses this
type IsEmptyErrorer interface {
	IsEmptyError() bool
}

type Terminal interface {
	Size() (int, int)
	IsTerminal() bool
}

type Input interface {
	fs.File
	Terminal
}

type Output interface {
	io.Writer
	Terminal
}

type OS interface {
	Stdin() Input
	Stdout() Output
	Stderr() Output
	Interrupt() chan struct{}
	Args() []string
	Environ() []string
	ConfigDir() (string, error)
	// returned Open() io.ReadSeeker can optionally implement io.Closer
	FS() fs.FS
	Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error)
	History() ([]string, error)
}

type FixedFileInfo struct {
	FName    string
	FSize    int64
	FMode    fs.FileMode
	FModTime time.Time
	FIsDir   bool
	FSys     interface{}
}

func (ffi FixedFileInfo) Name() string       { return ffi.FName }
func (ffi FixedFileInfo) Size() int64        { return ffi.FSize }
func (ffi FixedFileInfo) Mode() fs.FileMode  { return ffi.FMode }
func (ffi FixedFileInfo) ModTime() time.Time { return ffi.FModTime }
func (ffi FixedFileInfo) IsDir() bool        { return ffi.FIsDir }
func (ffi FixedFileInfo) Sys() interface{}   { return ffi.FSys }

type FileReader struct {
	R        io.Reader
	FileInfo FixedFileInfo
}

func (rf FileReader) Stat() (fs.FileInfo, error) { return rf.FileInfo, nil }
func (rf FileReader) Read(p []byte) (int, error) { return rf.R.Read(p) }
func (FileReader) Close() error                  { return nil }

// TODO: move
type DiscardCtxWriter struct {
	Ctx context.Context
}

func (o DiscardCtxWriter) Write(p []byte) (n int, err error) {
	if o.Ctx != nil {
		if err := o.Ctx.Err(); err != nil {
			return 0, err
		}
	}
	return n, nil
}

type CtxWriter struct {
	io.Writer
	Ctx context.Context
}

func (o CtxWriter) Write(p []byte) (n int, err error) {
	if o.Ctx != nil {
		if err := o.Ctx.Err(); err != nil {
			return 0, err
		}
	}
	return o.Writer.Write(p)
}

type Value interface {
	gojq.JQValue

	DisplayName() string
	ExtKeys() []string
}

type Display interface {
	Display(w io.Writer, opts Options) error
}

type ToBufferView interface {
	ToBufferView() (BufferView, error)
}

type JQValueEx interface {
	JQValueToGoJQEx(optsFn func() Options) interface{}
}

func valuePath(v *decode.Value) []interface{} {
	var parts []interface{}

	for v.Parent != nil {
		switch v.Parent.V.(type) {
		case decode.Struct:
			parts = append([]interface{}{v.Name}, parts...)
		case decode.Array:
			parts = append([]interface{}{v.Index}, parts...)
		}
		v = v.Parent
	}

	return parts
}

func valuePathDecorated(v *decode.Value, d Decorator) string {
	var parts []string

	for _, p := range valuePath(v) {
		switch p := p.(type) {
		case string:
			parts = append(parts, ".", d.ObjectKey.Wrap(p))
		case int:
			indexStr := strconv.Itoa(p)
			parts = append(parts, fmt.Sprintf("%s%s%s", d.Index.F("["), d.Number.F(indexStr), d.Index.F("]")))
		}
	}

	if len(parts) == 0 {
		return "."
	}

	return strings.Join(parts, "")
}

type iterFn func() (interface{}, bool)

func (i iterFn) Next() (interface{}, bool) { return i() }

type loadModule struct {
	init func() ([]*gojq.Query, error)
	load func(name string) (*gojq.Query, error)
}

func (l loadModule) LoadInitModules() ([]*gojq.Query, error)     { return l.init() }
func (l loadModule) LoadModule(name string) (*gojq.Query, error) { return l.load(name) }

func toString(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case gojq.JQValue:
		return toString(v.JQValueToGoJQ())
	default:
		b, err := toBytes(v)
		if err != nil {
			return "", fmt.Errorf("value can't be a string")
		}

		return string(b), nil
	}
}

func toBigInt(v interface{}) (*big.Int, error) {
	switch v := v.(type) {
	case int:
		return new(big.Int).SetInt64(int64(v)), nil
	case float64:
		return new(big.Int).SetInt64(int64(v)), nil
	case *big.Int:
		return v, nil
	default:
		return nil, fmt.Errorf("value is not a number")
	}
}

func toBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case []byte:
		return v, nil
	default:
		bb, err := toBuffer(v)
		if err != nil {
			return nil, fmt.Errorf("value is not bytes")
		}
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, bb); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
}

func toBuffer(v interface{}) (*bitio.Buffer, error) {
	return toBufferEx(v, false)
}

// TODO: refactor to return struct?
func toBufferEx(v interface{}, inArray bool) (*bitio.Buffer, error) {
	switch vv := v.(type) {
	case ToBufferView:
		bv, err := vv.ToBufferView()
		if err != nil {
			return nil, err
		}
		return bv.bb.BitBufRange(bv.r.Start, bv.r.Len)
	case string:
		return bitio.NewBufferFromBytes([]byte(vv), -1), nil
	case []byte:
		return bitio.NewBufferFromBytes(vv, -1), nil
	case int, float64, *big.Int:
		bi, err := toBigInt(v)
		if err != nil {
			return nil, err
		}

		if inArray {
			b := [1]byte{byte(bi.Uint64())}
			return bitio.NewBufferFromBytes(b[:], -1), nil
		}

		padBefore := (8 - (bi.BitLen() % 8)) % 8
		bb, err := bitio.NewBufferFromBytes(bi.Bytes(), -1).BitBufRange(int64(padBefore), int64(bi.BitLen()))
		if err != nil {
			return nil, err
		}
		return bb, nil
	case []interface{}:
		var rr []bitio.BitReadAtSeeker
		// TODO: optimize byte array case
		for _, e := range vv {
			eBB, eErr := toBufferEx(e, true)
			if eErr != nil {
				return nil, eErr
			}
			rr = append(rr, eBB)
		}

		mb, err := bitio.NewMultiBitReader(rr)
		if err != nil {
			return nil, err
		}

		bb, err := bitio.NewBufferFromBitReadSeeker(mb)
		if err != nil {
			return nil, err
		}

		return bb, nil
	default:
		return nil, fmt.Errorf("value can't be a buffer")
	}
}

func toBufferView(v interface{}) (BufferView, error) {
	switch vv := v.(type) {
	case ToBufferView:
		return vv.ToBufferView()
	default:
		bb, err := toBuffer(v)
		if err != nil {
			return BufferView{}, err
		}
		return bufferViewFromBuffer(bb, 8), nil
	}
}

// optsFn is a function as toValue is used by tovalue/0 so needs to be fast
func toValue(optsFn func() Options, v interface{}) (interface{}, bool) {
	switch v := v.(type) {
	case JQValueEx:
		return v.JQValueToGoJQEx(optsFn), true
	case gojq.JQValue:
		return v.JQValueToGoJQ(), true
	case nil, bool, float64, int, string, *big.Int, map[string]interface{}, []interface{}:
		return v, true
	default:
		return nil, false
	}
}

func queryErrorPosition(src string, v error) pos.Pos {
	var offset int

	if tokIf, ok := v.(interface{ Token() (string, int) }); ok { //nolint:errorlint
		_, offset = tokIf.Token()
	}
	if offset >= 0 {
		return pos.NewFromOffset(src, offset)
	}
	return pos.Pos{}
}

type Variable struct {
	Name  string
	Value interface{}
}

type Function struct {
	Names    []string
	MinArity int
	MaxArity int
	Fn       func(interface{}, []interface{}) interface{}
	IterFn   func(interface{}, []interface{}) gojq.Iter
}

type RunMode int

const (
	ScriptMode RunMode = iota
	REPLMode
	CompletionMode
)

type evalContext struct {
	// structcheck has problems with embedding https://gitlab.com/opennota/check#known-limitations
	ctx    context.Context
	output io.Writer
}

type Interp struct {
	registry       *registry.Registry
	os             OS
	initFqQuery    *gojq.Query
	includeCache   map[string]*gojq.Query
	interruptStack *ctxstack.Stack
	// global state, is ref as Interp i cloned per eval
	state *interface{}

	// new for each run, other values are copied by value
	evalContext evalContext
}

func New(os OS, registry *registry.Registry) (*Interp, error) {
	var err error

	i := &Interp{
		os:       os,
		registry: registry,
	}

	i.includeCache = map[string]*gojq.Query{}
	i.initFqQuery, err = gojq.Parse(initSource)
	if err != nil {
		return nil, fmt.Errorf("init:%s: %w", queryErrorPosition(initSource, err), err)
	}
	// TODO: refactor ctxstack have a CancelTop and return c context to Stop?
	i.interruptStack = ctxstack.New(func(stopCh chan struct{}) {
		select {
		case <-stopCh:
			return
		case <-os.Interrupt():
			return
		}
	})
	i.state = new(interface{})

	return i, nil
}

func (i *Interp) Stop() {
	// TODO: cancel all run instances?
	i.interruptStack.Stop()
}

func (i *Interp) Main(ctx context.Context, output Output, version string) error {
	var args []interface{}
	for _, a := range i.os.Args() {
		args = append(args, a)
	}

	input := map[string]interface{}{
		"args":    args,
		"version": version,
	}

	iter, err := i.EvalFunc(ctx, input, "_main", nil, output)
	if err != nil {
		fmt.Fprintln(i.os.Stderr(), err)
		return err
	}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		switch v := v.(type) {
		case error:
			if emptyErr, ok := v.(IsEmptyErrorer); ok && emptyErr.IsEmptyError() { //nolint:errorlint
				// no output
			} else {
				fmt.Fprintln(i.os.Stderr(), v)
			}
			return v
		case [2]interface{}:
			fmt.Fprintln(i.os.Stderr(), v[:]...)
		default:
			// TODO: can this happen?
			fmt.Fprintln(i.os.Stderr(), v)
		}
	}

	return nil
}

func (i *Interp) Eval(ctx context.Context, c interface{}, src string, srcFilename string, output io.Writer) (gojq.Iter, error) {
	gq, err := gojq.Parse(src)
	if err != nil {
		p := queryErrorPosition(src, err)
		return nil, compileError{
			err:      err,
			what:     "parse",
			filename: srcFilename,
			pos:      p,
		}
	}

	// make copy of interp
	ci := *i
	ni := &ci
	ni.evalContext = evalContext{}

	var variableNames []string
	var variableValues []interface{}
	for k, v := range i.variables() {
		variableNames = append(variableNames, "$"+k)
		variableValues = append(variableValues, v)
	}

	var funcCompilerOpts []gojq.CompilerOption
	for _, f := range ni.makeFunctions() {
		for _, n := range f.Names {
			if f.IterFn != nil {
				funcCompilerOpts = append(funcCompilerOpts,
					gojq.WithIterFunction(n, f.MinArity, f.MaxArity, f.IterFn))
			} else {
				funcCompilerOpts = append(funcCompilerOpts,
					gojq.WithFunction(n, f.MinArity, f.MaxArity, f.Fn))
			}
		}
	}

	compilerOpts := append([]gojq.CompilerOption{}, funcCompilerOpts...)
	compilerOpts = append(compilerOpts, gojq.WithEnvironLoader(ni.os.Environ))
	compilerOpts = append(compilerOpts, gojq.WithVariables(variableNames))
	compilerOpts = append(compilerOpts, gojq.WithModuleLoader(loadModule{
		init: func() ([]*gojq.Query, error) {
			return []*gojq.Query{i.initFqQuery}, nil
		},
		load: func(name string) (*gojq.Query, error) {
			// log.Printf("name: %#+v\n", name)
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			var filename string
			// support include "nonexisting?" to ignore include error
			var isTry bool
			if strings.HasSuffix(name, "?") {
				isTry = true
				filename = name[0 : len(name)-1]
			} else {
				filename = name
			}
			filename = filename + ".jq"

			pathPrefixes := []struct {
				prefix string
				cache  bool
				fn     func(filename string) (io.Reader, error)
			}{
				{
					"@builtin/", true, func(filename string) (io.Reader, error) {
						return builtinFS.Open(filename)
					},
				},
				{
					"@config/", false, func(filename string) (io.Reader, error) {
						configDir, err := i.os.ConfigDir()
						if err != nil {
							return nil, err
						}
						return i.os.FS().Open(filepath.Join(configDir, filename))
					},
				},
				{
					"", false, func(filename string) (io.Reader, error) {
						// TODO: jq $ORIGIN
						for _, path := range append([]string{"./"}, i.includePaths()...) {
							if f, err := i.os.FS().Open(filepath.Join(path, filename)); err == nil {
								return f, nil
							}
						}
						return nil, &fs.PathError{Op: "open", Path: filename, Err: fs.ErrNotExist}
					},
				},
			}

			for _, p := range pathPrefixes {
				if !strings.HasPrefix(filename, p.prefix) {
					continue
				}

				if p.cache {
					if q, ok := ni.includeCache[filename]; ok {
						return q, nil
					}
				}

				filenamePart := strings.TrimPrefix(filename, p.prefix)
				f, err := p.fn(filenamePart)
				if err != nil {
					if !isTry {
						return nil, err
					}
					f = &bytes.Buffer{}
				}

				b, err := io.ReadAll(f)
				if err != nil {
					return nil, err
				}
				s := string(b)
				q, err := gojq.Parse(s)
				if err != nil {
					p := queryErrorPosition(s, err)
					return nil, compileError{
						err:      err,
						what:     "parse",
						filename: filenamePart,
						pos:      p,
					}
				}

				// not identity body means it returns something, threat as dynamic include
				if q.Term.Type != gojq.TermTypeIdentity {
					gc, err := gojq.Compile(q, funcCompilerOpts...)
					if err != nil {
						return nil, err
					}
					iter := gc.RunWithContext(context.Background(), nil)
					var vs []interface{}
					for {
						v, ok := iter.Next()
						if !ok {
							break
						}
						if err, ok := v.(error); ok {
							return nil, err
						}
						vs = append(vs, v)
					}
					if len(vs) != 1 {
						return nil, fmt.Errorf("dynamic include: must output one string, got: %#v", vs)
					}
					s, sOk := vs[0].(string)
					if !sOk {
						return nil, fmt.Errorf("dynamic include: must be string, got %#v", s)
					}
					q, err = gojq.Parse(s)
					if err != nil {
						p := queryErrorPosition(s, err)
						return nil, compileError{
							err:      err,
							what:     "parse",
							filename: filenamePart,
							pos:      p,
						}
					}
				}

				// TODO: some better way of handling relative includes that
				// works with @builtin etc
				basePath := filepath.Dir(name)
				for _, i := range q.Imports {
					rewritePath := func(base, path string) string {
						if strings.HasPrefix(i.IncludePath, "@") {
							return path
						}
						if filepath.IsAbs(i.IncludePath) {
							return path
						}
						return filepath.Join(base, path)
					}
					i.IncludePath = rewritePath(basePath, i.IncludePath)
					i.ImportPath = rewritePath(basePath, i.ImportPath)
				}

				if p.cache {
					i.includeCache[filename] = q
				}

				return q, nil
			}

			panic("unreachable")
		},
	}))

	gc, err := gojq.Compile(gq, compilerOpts...)
	if err != nil {
		p := queryErrorPosition(src, err)
		return nil, compileError{
			err:      err,
			what:     "compile",
			filename: srcFilename,
			pos:      p,
		}
	}

	runCtx, runCtxCancelFn := i.interruptStack.Push(ctx)
	ni.evalContext.ctx = runCtx
	ni.evalContext.output = CtxWriter{Writer: output, Ctx: runCtx}
	iter := gc.RunWithContext(runCtx, c, variableValues...)

	iterWrapper := iterFn(func() (interface{}, bool) {
		v, ok := iter.Next()
		// gojq ctx cancel will not return ok=false, just cancelled error
		if !ok {
			runCtxCancelFn()
		}
		return v, ok
	})

	return iterWrapper, nil
}

func (i *Interp) EvalFunc(ctx context.Context, c interface{}, name string, args []interface{}, output io.Writer) (gojq.Iter, error) {
	var argsExpr []string
	for i := range args {
		argsExpr = append(argsExpr, fmt.Sprintf("$_args[%d]", i))
	}
	argExpr := ""
	if len(argsExpr) > 0 {
		argExpr = "(" + strings.Join(argsExpr, ";") + ")"
	}

	trampolineInput := map[string]interface{}{
		"input": c,
		"args":  args,
	}
	// _args to mark variable as internal and hide it from completion
	// {input: ..., args: [...]} | .args as {args: $_args} | .input | name[($_args[0]; ...)]
	trampolineExpr := fmt.Sprintf(". as {args: $_args} | .input | %s%s", name, argExpr)
	iter, err := i.Eval(ctx, trampolineInput, trampolineExpr, "", output)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (i *Interp) EvalFuncValues(ctx context.Context, c interface{}, name string, args []interface{}, output io.Writer) ([]interface{}, error) {
	iter, err := i.EvalFunc(ctx, c, name, args, output)
	if err != nil {
		return nil, err
	}

	var vs []interface{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		vs = append(vs, v)
	}

	return vs, nil
}

type Options struct {
	Depth          int    `mapstructure:"depth"`
	ArrayTruncate  int    `mapstructure:"array_truncate"`
	Verbose        bool   `mapstructure:"verbose"`
	DecodeProgress bool   `mapstructure:"decode_progress"`
	Color          bool   `mapstructure:"color"`
	Colors         string `mapstructure:"colors"`
	ByteColors     string `mapstructure:"byte_colors"`
	Unicode        bool   `mapstructure:"unicode"`
	RawOutput      bool   `mapstructure:"raw_output"`
	REPL           bool   `mapstructure:"repl"`
	RawString      bool   `mapstructure:"raw_string"`
	JoinString     string `mapstructure:"join_string"`
	Compact        bool   `mapstructure:"compact"`
	BitsFormat     string `mapstructure:"bits_format"`
	LineBytes      int    `mapstructure:"line_bytes"`
	DisplayBytes   int    `mapstructure:"display_bytes"`
	AddrBase       int    `mapstructure:"addrbase"`
	SizeBase       int    `mapstructure:"sizebase"`

	Decorator    Decorator
	BitsFormatFn func(bb *bitio.Buffer) (interface{}, error)
}

func bitsFormatFnFromOptions(opts Options) func(bb *bitio.Buffer) (interface{}, error) {
	switch opts.BitsFormat {
	case "md5":
		return func(bb *bitio.Buffer) (interface{}, error) {
			d := md5.New()
			if _, err := io.Copy(d, bb); err != nil {
				return "", err
			}
			return hex.EncodeToString(d.Sum(nil)), nil
		}
	case "base64":
		return func(bb *bitio.Buffer) (interface{}, error) {
			b := &bytes.Buffer{}
			e := base64.NewEncoder(base64.StdEncoding, b)
			if _, err := io.Copy(e, bb); err != nil {
				return "", err
			}
			e.Close()
			return b.String(), nil
		}
	case "truncate":
		// TODO: configure
		return func(bb *bitio.Buffer) (interface{}, error) {
			b := &bytes.Buffer{}
			if _, err := io.Copy(b, io.LimitReader(bb, 1024)); err != nil {
				return "", err
			}
			return b.String(), nil
		}
	case "string":
		return func(bb *bitio.Buffer) (interface{}, error) {
			b := &bytes.Buffer{}
			if _, err := io.Copy(b, bb); err != nil {
				return "", err
			}
			return b.String(), nil
		}
	case "snippet":
		fallthrough
	default:
		return func(bb *bitio.Buffer) (interface{}, error) {
			b := &bytes.Buffer{}
			e := base64.NewEncoder(base64.StdEncoding, b)
			if _, err := io.Copy(e, io.LimitReader(bb, 256)); err != nil {
				return "", err
			}
			e.Close()
			return fmt.Sprintf("<%s>%s", num.Bits(bb.Len()).StringByteBits(opts.SizeBase), b.String()), nil
		}
	}
}

func (i *Interp) lookupState(key string) interface{} {
	if i.state == nil {
		return nil
	}
	m, ok := (*i.state).(map[string]interface{})
	if !ok {
		return nil
	}
	return m[key]
}

func (i *Interp) includePaths() []string {
	pathsAny, _ := i.lookupState("include_paths").([]interface{})
	var paths []string
	for _, pathAny := range pathsAny {
		paths = append(paths, pathAny.(string))
	}
	return paths
}

func (i *Interp) variables() map[string]interface{} {
	variablesAny, _ := i.lookupState("variables").(map[string]interface{})
	return variablesAny
}

func (i *Interp) OptionsEx(fnOptsV ...interface{}) Options {
	opts := func() Options {
		vs, err := i.EvalFuncValues(i.evalContext.ctx, nil, "options", []interface{}{fnOptsV}, DiscardCtxWriter{Ctx: i.evalContext.ctx})
		if err != nil {
			return Options{}
		}
		if len(vs) < 1 {
			return Options{}
		}
		v := vs[0]
		if _, ok := v.(error); ok {
			return Options{}
		}
		m, ok := v.(map[string]interface{})
		if !ok {
			return Options{}
		}
		var opts Options
		_ = mapstructure.Decode(m, &opts)
		opts.Depth = num.MaxInt(0, opts.Depth)

		return opts
	}()

	opts.ArrayTruncate = num.MaxInt(0, opts.ArrayTruncate)
	opts.AddrBase = num.ClampInt(2, 36, opts.AddrBase)
	opts.SizeBase = num.ClampInt(2, 36, opts.SizeBase)
	opts.LineBytes = num.MaxInt(0, opts.LineBytes)
	opts.DisplayBytes = num.MaxInt(0, opts.DisplayBytes)
	opts.Decorator = decoratorFromOptions(opts)
	opts.BitsFormatFn = bitsFormatFnFromOptions(opts)

	return opts
}

func (i *Interp) Options() Options {
	return i.OptionsEx(nil)
}

func (i *Interp) NewColorJSON(opts Options) (*colorjson.Encoder, error) {
	indent := 2
	if opts.Compact {
		indent = 0
	}

	return colorjson.NewEncoder(
		opts.Color,
		false,
		indent,
		func(v interface{}) interface{} {
			if v, ok := toValue(func() Options { return opts }, v); ok {
				return v
			}
			panic(fmt.Sprintf("toValue not a JQValue value: %#v", v))
		},
		colorjson.Colors{
			Reset:     []byte(ansi.Reset.SetString),
			Null:      []byte(opts.Decorator.Null.SetString),
			False:     []byte(opts.Decorator.False.SetString),
			True:      []byte(opts.Decorator.True.SetString),
			Number:    []byte(opts.Decorator.Number.SetString),
			String:    []byte(opts.Decorator.String.SetString),
			ObjectKey: []byte(opts.Decorator.ObjectKey.SetString),
			Array:     []byte(opts.Decorator.Array.SetString),
			Object:    []byte(opts.Decorator.Object.SetString),
		},
	), nil
}
