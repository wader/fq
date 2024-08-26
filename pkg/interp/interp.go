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
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/wader/fq/internal/ansi"
	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/internal/colorjson"
	"github.com/wader/fq/internal/ctxstack"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/internal/iox"
	"github.com/wader/fq/internal/mapstruct"
	"github.com/wader/fq/internal/mathx"
	"github.com/wader/fq/internal/pos"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"

	"github.com/wader/gojq"
)

//go:embed interp.jq
//go:embed internal.jq
//go:embed options.jq
//go:embed binary.jq
//go:embed decode.jq
//go:embed registry_include.jq
//go:embed format_decode.jq
//go:embed format_func.jq
//go:embed grep.jq
//go:embed args.jq
//go:embed eval.jq
//go:embed query.jq
//go:embed repl.jq
//go:embed help.jq
//go:embed funcs.jq
//go:embed ansi.jq
//go:embed init.jq
var builtinFS embed.FS

var initSource = `include "@builtin/init"; .`

func init() {
	RegisterIter1("_readline", (*Interp)._readline)
	RegisterIter2("_eval", (*Interp)._eval)

	RegisterIter2("_stdio_read", (*Interp)._stdioRead)
	RegisterIter1("_stdio_write", (*Interp)._stdioWrite)
	RegisterFunc1("_stdio_info", (*Interp)._stdioInfo)

	RegisterFunc0("_extkeys", (*Interp)._extKeys)
	RegisterFunc0("_exttype", (*Interp)._extType)

	RegisterFunc0("_global_state", func(i *Interp, c any) any { return *i.state })
	RegisterFunc1("_global_state", func(i *Interp, _ any, v any) any { *i.state = v; return v })

	RegisterFunc0("history", (*Interp).history)
	RegisterIter1("_display", (*Interp)._display)
	RegisterFunc0("_can_display", (*Interp)._canDisplay)
	RegisterIter1("_hexdump", (*Interp)._hexdump)
	RegisterIter1("_print_color_json", (*Interp)._printColorJSON)

	RegisterFunc0("_is_completing", (*Interp)._isCompleting)
}

type valueError struct {
	v any
}

func (v valueError) Error() string { return fmt.Sprintf("error: %v", v.v) }
func (v valueError) Value() any    { return v.v }

type compileError struct {
	err      error
	what     string
	filename string
	pos      pos.Pos
}

func (ce compileError) Value() any {
	return map[string]any{
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
		filename = "expr"
	}
	return fmt.Sprintf("%s:%d:%d: %s: %s", filename, ce.pos.Line, ce.pos.Column, ce.what, ce.err.Error())
}

var ErrEOF = io.EOF
var ErrInterrupt = errors.New("Interrupt")

// gojq errors can implement this to signal exit code
type Exiter interface {
	ExitCode() int
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

type Platform struct {
	OS        string
	Arch      string
	GoVersion string
}

type CompleteFn func(line string, pos int) (newLine []string, shared int)

type ReadlineOpts struct {
	Prompt     string
	CompleteFn CompleteFn
}

type OS interface {
	Platform() Platform
	Stdin() Input
	Stdout() Output
	Stderr() Output
	InterruptChan() chan struct{}
	Args() []string
	Environ() []string
	ConfigDir() (string, error)
	// FS.File returned by FS().Open() can optionally implement io.Seeker
	FS() fs.FS
	Readline(opts ReadlineOpts) (string, error)
	History() ([]string, error)
}

type FixedFileInfo struct {
	FName    string
	FSize    int64
	FMode    fs.FileMode
	FModTime time.Time
	FIsDir   bool
	FSys     any
}

func (ffi FixedFileInfo) Name() string       { return ffi.FName }
func (ffi FixedFileInfo) Size() int64        { return ffi.FSize }
func (ffi FixedFileInfo) Mode() fs.FileMode  { return ffi.FMode }
func (ffi FixedFileInfo) ModTime() time.Time { return ffi.FModTime }
func (ffi FixedFileInfo) IsDir() bool        { return ffi.FIsDir }
func (ffi FixedFileInfo) Sys() any           { return ffi.FSys }

type FileReader struct {
	R        io.Reader
	FileInfo FixedFileInfo
}

func (rf FileReader) Stat() (fs.FileInfo, error) { return rf.FileInfo, nil }
func (rf FileReader) Read(p []byte) (int, error) { return rf.R.Read(p) }
func (FileReader) Close() error                  { return nil }

type Value interface {
	gojq.JQValue

	ExtType() string
	ExtKeys() []string
}

type Display interface {
	Display(w io.Writer, opts *Options) error
}

type JQValueEx interface {
	gojq.JQValue
	JQValueToGoJQEx(optsFn func() (*Options, error)) any
}

func valuePath(v *decode.Value) []any {
	var parts []any

	for v.Parent != nil {
		switch vv := v.Parent.V.(type) {
		case *decode.Compound:
			if vv.IsArray {
				parts = append([]any{v.Index}, parts...)
			} else {
				parts = append([]any{v.Name}, parts...)
			}
		}
		v = v.Parent
	}

	return parts
}

func valuePathExprDecorated(v *decode.Value, d Decorator) string {
	parts := []string{"."}

	for i, p := range valuePath(v) {
		switch p := p.(type) {
		case string:
			if i > 0 {
				parts = append(parts, ".")
			}
			parts = append(parts, d.ObjectKey.Wrap(p))
		case int:
			indexStr := strconv.Itoa(p)
			parts = append(parts, fmt.Sprintf("%s%s%s", d.Index.F("["), d.Number.F(indexStr), d.Index.F("]")))
		}
	}

	return strings.Join(parts, "")
}

type iterFn func() (any, bool)

func (i iterFn) Next() (any, bool) { return i() }

type loadModule struct {
	init func() ([]*gojq.Query, error)
	load func(name string) (*gojq.Query, error)
}

func (l loadModule) LoadInitModules() ([]*gojq.Query, error)     { return l.init() }
func (l loadModule) LoadModule(name string) (*gojq.Query, error) { return l.load(name) }

func toString(v any) (string, error) {
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

func toBigInt(v any) (*big.Int, error) {
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

func toBytes(v any) ([]byte, error) {
	switch v := v.(type) {
	default:
		br, err := ToBitReader(v)
		if err != nil {
			return nil, fmt.Errorf("value is not bytes")
		}
		buf := &bytes.Buffer{}
		if _, err := bitiox.CopyBits(buf, br); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}
}

func queryErrorPosition(expr string, v error) pos.Pos {
	var offset int

	var e *gojq.ParseError
	if errors.As(v, &e) {
		offset = e.Offset
	}
	if offset >= 0 {
		return pos.NewFromOffset(expr, offset)
	}
	return pos.Pos{}
}

type Variable struct {
	Name  string
	Value any
}

type RunMode int

const (
	ScriptMode RunMode = iota
	REPLMode
	CompletionMode
)

type EvalInstance struct {
	Ctx          context.Context
	Output       io.Writer
	IsCompleting bool

	includeSeen map[string]struct{}
}

type Interp struct {
	Registry *Registry
	OS       OS

	initQuery      *gojq.Query
	includeCache   map[string]*gojq.Query
	interruptStack *ctxstack.Stack
	// global state, is ref as Interp is cloned per eval
	state *any

	// new for each eval, other values are copied by value
	EvalInstance EvalInstance
}

func New(os OS, registry *Registry) (*Interp, error) {
	var err error

	i := &Interp{
		OS:       os,
		Registry: registry,
	}

	i.includeCache = map[string]*gojq.Query{}
	i.initQuery, err = gojq.Parse(initSource)
	if err != nil {
		return nil, fmt.Errorf("init:%s: %w", queryErrorPosition(initSource, err), err)
	}
	// TODO: refactor ctxstack have a CancelTop and return c context to Stop?
	i.interruptStack = ctxstack.New(func(stopCh chan struct{}) {
		select {
		case <-stopCh:
			return
		case <-os.InterruptChan():
			return
		}
	})
	i.state = new(any)

	return i, nil
}

func (i *Interp) Stop() {
	// TODO: cancel all run instances?
	i.interruptStack.Stop()
}

func (i *Interp) Main(ctx context.Context, output Output, versionStr string) error {
	var args []any
	for _, a := range i.OS.Args() {
		args = append(args, a)
	}

	platform := i.OS.Platform()
	input := map[string]any{
		"args":       args,
		"version":    versionStr,
		"os":         platform.OS,
		"arch":       platform.Arch,
		"go_version": platform.GoVersion,
	}

	iter, err := i.EvalFunc(ctx, input, "_main", nil, EvalOpts{output: output})
	if err != nil {
		fmt.Fprintln(i.OS.Stderr(), err)
		return err
	}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		switch v := v.(type) {
		case error:
			var haltErr *gojq.HaltError
			if errors.As(v, &haltErr) {
				if haltErrV := haltErr.Value(); haltErrV != nil {
					if str, ok := haltErrV.(string); ok {
						if _, err := i.OS.Stderr().Write([]byte(str)); err != nil {
							return err
						}
					} else {
						bs, _ := gojq.Marshal(haltErrV)
						if _, err := i.OS.Stderr().Write(bs); err != nil {
							return err
						}
						if _, err := i.OS.Stderr().Write([]byte{'\n'}); err != nil {
							return err
						}
					}
				}

				return haltErr
			} else if errors.Is(v, context.Canceled) {
				// ignore context cancel here for now, which means user somehow interrupted the interpreter
				// TODO: handle this inside interp.jq instead but then we probably have to do nested
				// eval and or also use different contexts for the interpreter and reading/decoding
			} else {
				fmt.Fprintln(i.OS.Stderr(), v)
			}
			return v
		default:
			// TODO: can this happen?
			fmt.Fprintln(i.OS.Stderr(), v)
		}
	}

	return nil
}

type completionResult struct {
	Names  []string
	Prefix string
}

type readlineOpts struct {
	Prompt   string
	Complete string
	Timeout  float64
}

func (i *Interp) _readline(c any, opts readlineOpts) gojq.Iter {
	if i.EvalInstance.IsCompleting {
		return gojq.NewIter()
	}

	expr, err := i.OS.Readline(ReadlineOpts{
		Prompt: opts.Prompt,
		CompleteFn: func(line string, pos int) (newLine []string, shared int) {
			completeCtx := i.EvalInstance.Ctx
			if opts.Timeout > 0 {
				var completeCtxCancelFn context.CancelFunc
				completeCtx, completeCtxCancelFn = context.WithTimeout(i.EvalInstance.Ctx, time.Duration(opts.Timeout*float64(time.Second)))
				defer completeCtxCancelFn()
			}

			names, shared, err := func() (newLine []string, shared int, err error) {
				// c | opts.Complete(line; pos)
				vs, err := i.EvalFuncValues(
					completeCtx,
					c,
					opts.Complete,
					[]any{line, pos},
					EvalOpts{
						output:       iox.DiscardCtxWriter{Ctx: completeCtx},
						isCompleting: true,
					},
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
				r, ok := gojqx.CastFn[completionResult](v, mapstruct.ToStruct)
				if !ok {
					return nil, pos, fmt.Errorf("completion result not a map")
				}

				sharedLen := len(r.Prefix)

				return r.Names, sharedLen, nil
			}()

			// TODO: how to report err?
			_ = err

			return names, shared
		},
	})

	if errors.Is(err, ErrInterrupt) {
		return gojq.NewIter(valueError{"interrupt"})
	} else if errors.Is(err, ErrEOF) {
		return gojq.NewIter(valueError{"eof"})
	} else if err != nil {
		return gojq.NewIter(err)
	}

	return gojq.NewIter(expr)
}

type evalOpts struct {
	Filename string
}

func (i *Interp) _eval(c any, expr string, opts evalOpts) gojq.Iter {
	var err error

	iter, err := i.Eval(i.EvalInstance.Ctx, c, expr, EvalOpts{
		filename: opts.Filename,
		output:   i.EvalInstance.Output,
	})
	if err != nil {
		return gojq.NewIter(err)
	}

	return iter
}

func (i *Interp) _extKeys(c any) any {
	if v, ok := c.(Value); ok {
		var vs []any
		for _, s := range v.ExtKeys() {
			vs = append(vs, s)
		}
		return vs
	}
	return nil
}

func (i *Interp) _extType(c any) any {
	if v, ok := c.(Value); ok {
		return v.ExtType()
	}
	return gojq.TypeOf(c)
}

func (i *Interp) _stdioFdName(s string) (any, error) {
	switch s {
	case "stdin":
		return i.OS.Stdin(), nil
	case "stdout":
		return i.OS.Stdout(), nil
	case "stderr":
		return i.OS.Stderr(), nil
	default:
		return nil, fmt.Errorf("unknown fd %s", s)
	}
}

func (i *Interp) _stdioRead(c any, fdName string, l int) gojq.Iter {
	fd, err := i._stdioFdName(fdName)
	if err != nil {
		return gojq.NewIter(err)
	}
	r, ok := fd.(io.Reader)
	if !ok {
		return gojq.NewIter(fmt.Errorf("%s is not a writeable", fdName))
	}

	if i.EvalInstance.IsCompleting {
		return gojq.NewIter("")
	}

	buf := make([]byte, l)
	n, err := io.ReadFull(r, buf)
	s := string(buf[0:n])

	vs := []any{s}
	switch {
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
		vs = append(vs, valueError{"eof"})
	default:
		vs = append(vs, err)
	}

	return gojq.NewIter(vs...)
}

func (i *Interp) _stdioWrite(c any, fdName string) gojq.Iter {
	fd, err := i._stdioFdName(fdName)
	if err != nil {
		return gojq.NewIter(err)
	}
	w, ok := fd.(io.Writer)
	if !ok {
		return gojq.NewIter(fmt.Errorf("%s is not a writeable", fdName))
	}
	if i.EvalInstance.IsCompleting {
		return gojq.NewIter()
	}

	if _, err := fmt.Fprint(w, c); err != nil {
		return gojq.NewIter(err)
	}
	return gojq.NewIter()
}

func (i *Interp) _stdioInfo(c any, fdName string) any {
	fd, err := i._stdioFdName(fdName)
	if err != nil {
		return err
	}
	t, ok := fd.(Terminal)
	if !ok {
		return fmt.Errorf("%s is not a terminal", fdName)
	}

	w, h := t.Size()
	return map[string]any{
		"is_terminal": t.IsTerminal(),
		"width":       w,
		"height":      h,
	}
}

func (i *Interp) history(c any) any {
	hs, err := i.OS.History()
	if err != nil {
		return err
	}
	var vs []any
	for _, s := range hs {
		vs = append(vs, s)
	}
	return vs
}

func (i *Interp) _display(c any, v any) gojq.Iter {
	opts, err := OptionsFromValue(v)
	if err != nil {
		return gojq.NewIter(err)
	}

	switch v := c.(type) {
	case Display:
		if err := v.Display(i.EvalInstance.Output, opts); err != nil {
			return gojq.NewIter(err)
		}
		return gojq.NewIter()
	default:
		return gojq.NewIter(fmt.Errorf("%+#v: not displayable", c))
	}
}

func (i *Interp) _canDisplay(c any) any {
	_, ok := c.(Display)
	return ok
}

func (i *Interp) _hexdump(c any, v any) gojq.Iter {
	opts, err := OptionsFromValue(v)
	if err != nil {
		return gojq.NewIter(err)
	}

	bv, err := toBinary(c)
	if err != nil {
		return gojq.NewIter(err)
	}
	if err := hexdump(i.EvalInstance.Output, bv, opts); err != nil {
		return gojq.NewIter(err)
	}

	return gojq.NewIter()
}

func (i *Interp) _printColorJSON(c any, v any) gojq.Iter {
	opts, err := OptionsFromValue(v)
	if err != nil {
		return gojq.NewIter(err)
	}

	indent := 2
	if opts.Compact {
		indent = 0
	}

	cj := colorjson.NewEncoder(colorjson.Options{
		Color:  opts.Color,
		Tab:    false,
		Indent: indent,
		// uses a function to cache OptionsFromValue
		ValueFn: func(v any) (any, error) { return toValue(func() (*Options, error) { return opts, nil }, v) },
		Colors: colorjson.Colors{
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
	})
	if err := cj.Marshal(c, i.EvalInstance.Output); err != nil {
		return gojq.NewIter(err)
	}

	return gojq.NewIter()
}

func (i *Interp) _isCompleting(c any) any {
	return i.EvalInstance.IsCompleting
}

type pathResolver struct {
	prefix string
	open   func(filename string) (io.ReadCloser, string, error)
}

func (i *Interp) lookupPathResolver(filename string) (pathResolver, error) {
	configDir, err := i.OS.ConfigDir()
	if err != nil {
		return pathResolver{}, err
	}

	resolvePaths := []pathResolver{
		{
			"@builtin/",
			func(filename string) (io.ReadCloser, string, error) {
				f, err := builtinFS.Open(filename)
				return f, "@builtin/" + filename, err
			},
		},
		{
			"@config/", func(filename string) (io.ReadCloser, string, error) {
				p := path.Join(configDir, filename)
				f, err := i.OS.FS().Open(p)
				return f, p, err
			},
		},
		{
			"", func(filename string) (io.ReadCloser, string, error) {
				if path.IsAbs(filename) {
					f, err := i.OS.FS().Open(filename)
					return f, filename, err
				}

				// TODO: jq $ORIGIN
				for _, includePath := range append([]string{"./"}, i.includePaths()...) {
					p := path.Join(includePath, filename)
					if f, err := i.OS.FS().Open(path.Join(includePath, filename)); err == nil {
						return f, p, nil
					}
				}

				return nil, "", &fs.PathError{Op: "open", Path: filename, Err: fs.ErrNotExist}
			},
		},
	}
	for _, p := range resolvePaths {
		if strings.HasPrefix(filename, p.prefix) {
			return p, nil
		}
	}
	return pathResolver{}, fmt.Errorf("could not resolve path: %s", filename)
}

type EvalOpts struct {
	filename     string
	output       io.Writer
	isCompleting bool
}

func (i *Interp) Eval(ctx context.Context, c any, expr string, opts EvalOpts) (gojq.Iter, error) {
	gq, err := gojq.Parse(expr)
	if err != nil {
		p := queryErrorPosition(expr, err)
		return nil, compileError{
			err:      err,
			what:     "parse",
			filename: opts.filename,
			pos:      p,
		}
	}

	// make copy of interp and give it its own eval context
	ci := *i
	ni := &ci
	ni.EvalInstance = EvalInstance{
		includeSeen: map[string]struct{}{},
	}

	var variableNames []string
	var variableValues []any
	for k, v := range i.slurps() {
		variableNames = append(variableNames, "$"+k)
		variableValues = append(variableValues, v)
	}

	var funcCompilerOpts []gojq.CompilerOption

	for _, fn := range i.Registry.EnvFuncFns {
		f := fn(ni)
		if f.IterFn != nil {
			funcCompilerOpts = append(funcCompilerOpts,
				gojq.WithIterFunction(f.Name, f.MinArity, f.MaxArity, f.IterFn))
		} else {
			funcCompilerOpts = append(funcCompilerOpts,
				gojq.WithFunction(f.Name, f.MinArity, f.MaxArity, f.FuncFn))
		}
	}

	compilerOpts := append([]gojq.CompilerOption{}, funcCompilerOpts...)
	compilerOpts = append(compilerOpts, gojq.WithEnvironLoader(ni.OS.Environ))
	compilerOpts = append(compilerOpts, gojq.WithVariables(variableNames))
	compilerOpts = append(compilerOpts, gojq.WithModuleLoader(loadModule{
		init: func() ([]*gojq.Query, error) {
			return []*gojq.Query{i.initQuery}, nil
		},
		load: func(name string) (*gojq.Query, error) {
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

			pr, err := i.lookupPathResolver(filename)
			if err != nil {
				return nil, err
			}

			// skip if this eval instance has already included the file
			if _, ok := ni.EvalInstance.includeSeen[filename]; ok {
				return &gojq.Query{Term: &gojq.Term{Type: gojq.TermTypeIdentity}}, nil
			}
			ni.EvalInstance.includeSeen[filename] = struct{}{}

			// return cached version if file has already been parsed
			if q, ok := ni.includeCache[filename]; ok {
				return q, nil
			}

			filenamePart := strings.TrimPrefix(filename, pr.prefix)

			f, absPath, err := pr.open(filenamePart)
			if err != nil {
				if !isTry {
					return nil, err
				}
				f = io.NopCloser(&bytes.Buffer{})
			}
			defer f.Close()

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
					filename: absPath,
					pos:      p,
				}
			}
			// has some root expression, threat as dynamic include
			if q.Term != nil || q.Op != gojq.Operator(0) {
				gc, err := gojq.Compile(q, funcCompilerOpts...)
				if err != nil {
					return nil, err
				}
				iter := gc.RunWithContext(context.Background(), nil)
				var vs []any
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
						what:     "dynamic include parse",
						filename: filenamePart,
						pos:      p,
					}
				}
			}

			// TODO: some better way of handling relative includes that
			// works with @builtin etc
			basePath := path.Dir(name)
			for _, qi := range q.Imports {
				rewritePath := func(base, includePath string) string {
					if strings.HasPrefix(includePath, "@") || path.IsAbs(includePath) {
						return includePath
					}

					return path.Join(base, includePath)
				}
				if qi.IncludePath != "" {
					qi.IncludePath = rewritePath(basePath, qi.IncludePath)
				}
				if qi.ImportPath != "" {
					qi.ImportPath = rewritePath(basePath, qi.ImportPath)
				}
			}

			i.includeCache[filename] = q

			return q, nil
		},
	}))

	gc, err := gojq.Compile(gq, compilerOpts...)
	if err != nil {
		p := queryErrorPosition(expr, err)
		return nil, compileError{
			err:      err,
			what:     "compile",
			filename: opts.filename,
			pos:      p,
		}
	}

	output := opts.output
	if opts.output == nil {
		output = io.Discard
	}

	runCtx, runCtxCancelFn := i.interruptStack.Push(ctx)
	ni.EvalInstance.Ctx = runCtx
	ni.EvalInstance.Output = iox.CtxWriter{Writer: output, Ctx: runCtx}
	// inherit or maybe set
	ni.EvalInstance.IsCompleting = i.EvalInstance.IsCompleting || opts.isCompleting
	iter := gc.RunWithContext(runCtx, c, variableValues...)

	iterWrapper := iterFn(func() (any, bool) {
		v, ok := iter.Next()
		// gojq ctx cancel will not return ok=false, just cancelled error
		if !ok {
			runCtxCancelFn()
		} else if _, ok := v.(error); ok {
			runCtxCancelFn()
		}
		return v, ok
	})

	return iterWrapper, nil
}

func (i *Interp) EvalFunc(ctx context.Context, c any, name string, args []any, opts EvalOpts) (gojq.Iter, error) {
	var argsExpr []string
	for i := range args {
		argsExpr = append(argsExpr, fmt.Sprintf("$_args[%d]", i))
	}
	argExpr := ""
	if len(argsExpr) > 0 {
		argExpr = "(" + strings.Join(argsExpr, ";") + ")"
	}

	trampolineInput := map[string]any{
		"input": c,
		"args":  args,
	}
	// _args to mark variable as internal and hide it from completion
	// {input: ..., args: [...]} | .args as {args: $_args} | .input | name[($_args[0]; ...)]
	trampolineExpr := fmt.Sprintf(". as {args: $_args} | .input | %s%s", name, argExpr)
	iter, err := i.Eval(ctx, trampolineInput, trampolineExpr, opts)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (i *Interp) EvalFuncValues(ctx context.Context, c any, name string, args []any, opts EvalOpts) ([]any, error) {
	iter, err := i.EvalFunc(ctx, c, name, args, opts)
	if err != nil {
		return nil, err
	}

	var vs []any
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
	Depth          int
	ArrayTruncate  int
	StringTruncate int
	Verbose        bool
	Width          int
	DecodeProgress bool
	Color          bool
	Colors         map[string]string
	ByteColors     []struct {
		Ranges [][2]int
		Value  string
	}
	Unicode      bool
	RawOutput    bool
	REPL         bool
	RawString    bool
	JoinString   string
	Compact      bool
	BitsFormat   string
	LineBytes    int
	DisplayBytes int
	Addrbase     int
	Sizebase     int
	SkipGaps     bool

	Decorator    Decorator
	BitsFormatFn func(br bitio.ReaderAtSeeker) (any, error)
}

func OptionsFromValue(v any) (*Options, error) {
	var opts Options
	_ = mapstruct.ToStruct(v, &opts)
	opts.ArrayTruncate = max(0, opts.ArrayTruncate)
	opts.StringTruncate = max(0, opts.StringTruncate)
	opts.Depth = max(0, opts.Depth)
	opts.Addrbase = mathx.Clamp(2, 36, opts.Addrbase)
	opts.Sizebase = mathx.Clamp(2, 36, opts.Sizebase)
	opts.LineBytes = max(1, opts.LineBytes)
	opts.DisplayBytes = max(0, opts.DisplayBytes)
	opts.Decorator = decoratorFromOptions(opts)
	if fn, err := bitsFormatFnFromOptions(opts); err != nil {
		return nil, err
	} else {
		opts.BitsFormatFn = fn
	}

	return &opts, nil
}

func bitsFormatFnFromOptions(opts Options) (func(br bitio.ReaderAtSeeker) (any, error), error) {
	switch opts.BitsFormat {
	case "md5":
		return func(br bitio.ReaderAtSeeker) (any, error) {
			d := md5.New()
			if _, err := bitiox.CopyBits(d, br); err != nil {
				return "", err
			}
			return hex.EncodeToString(d.Sum(nil)), nil
		}, nil
	case "hex":
		return func(br bitio.ReaderAtSeeker) (any, error) {
			b := &bytes.Buffer{}
			e := hex.NewEncoder(b)
			if _, err := bitiox.CopyBits(e, br); err != nil {
				return "", err
			}
			return b.String(), nil
		}, nil
	case "base64":
		return func(br bitio.ReaderAtSeeker) (any, error) {
			b := &bytes.Buffer{}
			e := base64.NewEncoder(base64.StdEncoding, b)
			if _, err := bitiox.CopyBits(e, br); err != nil {
				return "", err
			}
			e.Close()
			return b.String(), nil
		}, nil
	case "truncate":
		// TODO: configure
		return func(br bitio.ReaderAtSeeker) (any, error) {
			b := &bytes.Buffer{}
			if _, err := bitiox.CopyBits(b, bitio.NewLimitReader(br, 1024*8)); err != nil {
				return "", err
			}
			return b.String(), nil
		}, nil
	case "string":
		return func(br bitio.ReaderAtSeeker) (any, error) {
			b := &bytes.Buffer{}
			if _, err := bitiox.CopyBits(b, br); err != nil {
				return "", err
			}
			return b.String(), nil
		}, nil
	case "snippet":
		return func(br bitio.ReaderAtSeeker) (any, error) {
			b := &bytes.Buffer{}
			e := base64.NewEncoder(base64.StdEncoding, b)
			if _, err := bitiox.CopyBits(e, bitio.NewLimitReader(br, 256*8)); err != nil {
				return "", err
			}
			e.Close()
			brLen, err := bitiox.Len(br)
			if err != nil {
				return nil, err
			}

			return fmt.Sprintf("<%s>%s", mathx.Bits(brLen).StringByteBits(opts.Sizebase), b.String()), nil
		}, nil
	case "byte_array":
		return func(br bitio.ReaderAtSeeker) (any, error) {
			b := &bytes.Buffer{}
			if _, err := bitiox.CopyBits(b, br); err != nil {
				return "", err
			}
			var v []any
			for _, bv := range b.Bytes() {
				v = append(v, int(bv))
			}
			return v, nil
		}, nil
	default:
		return nil, fmt.Errorf("invalid bits format %q", opts.BitsFormat)
	}
}

func (i *Interp) lookupState(key string) any {
	if i.state == nil {
		return nil
	}
	m, ok := (*i.state).(map[string]any)
	if !ok {
		return nil
	}
	return m[key]
}

func (i *Interp) includePaths() []string {
	pathsAny, _ := i.lookupState("include_paths").([]any)
	var paths []string
	for _, pathAny := range pathsAny {
		path, ok := pathAny.(string)
		if !ok {
			panic("path not string")
		}
		paths = append(paths, path)
	}
	return paths
}

func (i *Interp) slurps() map[string]any {
	slurpsAny, _ := i.lookupState("slurps").(map[string]any)
	return slurpsAny
}
