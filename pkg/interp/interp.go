package interp

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"fq/format/registry"
	"fq/internal/ansi"
	"fq/internal/colorjson"
	"fq/internal/ctxstack"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/ranges"
	"io"
	"math/big"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

//go:embed *.jq
var builtinFS embed.FS

var fqInitSource = `include "@builtin/fq";`

type valueErr struct {
	v interface{}
}

func (v valueErr) Error() string      { return fmt.Sprintf("error: %v", v.v) }
func (v valueErr) Value() interface{} { return v.v }

var ErrEOF = io.EOF
var ErrInterrupt = errors.New("Interrupt")

type Output interface {
	io.Writer
	Size() (int, int)
	IsTerminal() bool
}

type OS interface {
	Stdin() io.Reader
	Stdout() Output
	Stderr() io.Writer
	Interrupt() chan struct{}
	Args() []string
	Environ() []string
	ConfigDir() (string, error)
	// returned io.ReadSeeker can optionally implement io.Closer
	Open(name string) (io.ReadSeeker, error)
	Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error)
	History() ([]string, error)
}

// TODO: would be nice if gojq had something for this? maybe missing something?
func offsetToLineColumn(s string, offset int) (int, int) {
	co := 0
	line := 1
	for {
		no := strings.Index(s[co:], "\n")
		if no == -1 || co+no >= offset {
			return line, offset - co
		}
		co += no + 1
		line++
	}
}

// TODO: move
type DiscardOutput struct {
	Output
	Ctx context.Context
}

func (o DiscardOutput) Write(p []byte) (n int, err error) {
	if o.Ctx != nil {
		if err := o.Ctx.Err(); err != nil {
			return 0, err
		}
	}
	return n, nil
}

type CtxOutput struct {
	Output
	Ctx context.Context
}

func (o CtxOutput) Write(p []byte) (n int, err error) {
	if o.Ctx != nil {
		if err := o.Ctx.Err(); err != nil {
			return 0, err
		}
	}
	return o.Output.Write(p)
}

type InterpObject interface {
	gojq.JQValue

	DisplayName() string
	ExtValueKeys() []string
}

type Display interface {
	Display(w io.Writer, opts Options) error
}

type Preview interface {
	Preview(w io.Writer, opts Options) error
}

type ToBuffer interface {
	ToBuffer() (*bitio.Buffer, error)
}

type ToBufferRange interface {
	ToBufferRange() (bufferRange, error)
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

func (l loadModule) LoadModule(name string) (*gojq.Query, error) { return l.load(name) }
func (l loadModule) LoadInitModules() ([]*gojq.Query, error)     { return l.init() }

func toBool(v interface{}) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil
	case *big.Int:
		return v.Int64() != 0, nil
	case int:
		return v != 0, nil
	case float64:
		return v != 0, nil
	default:
		return false, fmt.Errorf("value is not a number")
	}
}

func toBoolZ(v interface{}) bool {
	b, _ := toBool(v)
	return b
}

func toInt(v interface{}) (int, error) {
	switch v := v.(type) {
	case *big.Int:
		return int(v.Int64()), nil
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, fmt.Errorf("value is not a number")
	}
}

func toIntZ(v interface{}) int {
	n, _ := toInt(v)
	return n
}

func toString(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	default:
		b, err := toBytes(v)
		if err != nil {
			return "", fmt.Errorf("value can't be a string")
		}

		return string(b), nil
	}
}

func toStringZ(v interface{}) string {
	s, _ := toString(v)
	return s
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
	case ToBuffer:
		return vv.ToBuffer()
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
		} else {
			padBefore := (8 - (bi.BitLen() % 8)) % 8
			bb, err := bitio.NewBufferFromBytes(bi.Bytes(), -1).BitBufRange(int64(padBefore), int64(bi.BitLen()))
			if err != nil {
				return nil, err
			}
			return bb, nil
		}
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
		return nil, fmt.Errorf("value can't be buffer")
	}
}

func toBufferRange(v interface{}) (bufferRange, error) {
	switch vv := v.(type) {
	case ToBufferRange:
		return vv.ToBufferRange()
	default:
		bb, err := toBuffer(v)
		if err != nil {
			return bufferRange{}, err
		}
		return bufferRange{bb: bb, r: ranges.Range{Len: bb.Len()}}, nil
	}
}

func toValue(v interface{}) interface{} {
	switch v := v.(type) {
	case gojq.JQValue:
		return v.JQValue()
	case nil, bool, float64, int, string, *big.Int, map[string]interface{}, []interface{}:
		return v
	default:
		return nil
	}
}

func queryErrorPosition(v error) string {
	var offset int
	var content string

	if tokIf, ok := v.(interface{ Token() (string, int) }); ok {
		_, offset = tokIf.Token()
	}
	if qeIf, ok := v.(interface {
		QueryParseError() (string, string, string, error)
	}); ok {
		_, _, content, _ = qeIf.QueryParseError()
	}

	if offset > 0 && content != "" {
		l, c := offsetToLineColumn(content, offset)
		return fmt.Sprintf(":%d:%d", l, c)
	}
	return ""
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

type runContext struct {
	ctx          context.Context
	stdout       Output // TODO: rename?
	mode         RunMode
	state        map[string]interface{}
	debugFn      string
	includeStack []string
}

type Interp struct {
	// variables map[string]interface{}
	registry *registry.Registry
	os       OS

	initFqQuery *gojq.Query

	includeCache map[string]*gojq.Query

	interruptStack *ctxstack.Stack

	// new for each run, other values are copied by ref
	runContext
}

func New(os OS, registry *registry.Registry) (*Interp, error) {
	var err error

	i := &Interp{
		os:       os,
		registry: registry,
	}

	i.includeCache = map[string]*gojq.Query{}
	i.initFqQuery, err = gojq.Parse(fqInitSource)
	if err != nil {
		return nil, fmt.Errorf("init%s %w", queryErrorPosition(err), err)
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

	return i, nil
}

func (i *Interp) Stop() {
	// TODO: cancel all run instances?
	i.interruptStack.Stop()
}

func (i *Interp) Main(ctx context.Context, stdout io.Writer, version string) error {
	runMode := ScriptMode

	var args []interface{}
	for _, a := range i.os.Args() {
		args = append(args, a)
	}

	input := map[string]interface{}{
		"args":    args,
		"version": version,
	}

	iter, err := i.EvalFunc(ctx, runMode, input, "main", nil, i.os.Stdout(), "")
	if err != nil {
		return err
	}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		switch v := v.(type) {
		case error:
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

func (i *Interp) Eval(ctx context.Context, mode RunMode, c interface{}, src string, stdout Output, debugFn string) (gojq.Iter, error) {
	var err error

	// TODO: did not work
	// nq := &(*q)

	gq, err := gojq.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("eval%s: %w", queryErrorPosition(err), err)
	}

	// make copy of interp
	ci := *i
	ni := &ci

	newState := map[string]interface{}{}
	if i.runContext.state != nil {
		for k, v := range i.runContext.state {
			newState[k] = v
		}
	}

	ni.runContext = runContext{
		state:        newState,
		mode:         mode,
		debugFn:      debugFn,
		includeStack: []string{"eval"},
	}

	// var variableNames []string
	// var variableValues []interface{}
	// for k, v := range ni.variables {
	// 	variableNames = append(variableNames, k)
	// 	variableValues = append(variableValues, v)
	// }

	var compilerOpts []gojq.CompilerOption
	for _, f := range ni.makeFunctions(ni.registry) {
		for _, n := range f.Names {
			if f.IterFn != nil {
				compilerOpts = append(compilerOpts,
					gojq.WithIterFunction(n, f.MinArity, f.MaxArity, f.IterFn))
			} else {
				compilerOpts = append(compilerOpts,
					gojq.WithFunction(n, f.MinArity, f.MaxArity, f.Fn))
			}
		}
	}
	compilerOpts = append(compilerOpts, gojq.WithEnvironLoader(ni.os.Environ))
	// compilerOpts = append(compilerOpts, gojq.WithVariables(variableNames))
	compilerOpts = append(compilerOpts, gojq.WithModuleLoader(loadModule{
		init: func() ([]*gojq.Query, error) {
			return []*gojq.Query{i.initFqQuery}, nil
		},
		load: func(name string) (*gojq.Query, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			var filename string
			// suport include "nonexisting?" to ignore include error
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
						return i.os.Open(filepath.Join(configDir, filename))
					},
				},
				{
					"", false, func(filename string) (io.Reader, error) {
						return i.os.Open(filename)
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
					err = nil
					f = &bytes.Buffer{}
				}

				b, err := io.ReadAll(f)
				if err != nil {
					return nil, err
				}
				q, err := gojq.Parse(string(b))
				if err != nil {
					return nil, fmt.Errorf("%s%s: %w", name, queryErrorPosition(err), err)
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
		return nil, fmt.Errorf("eval%s: %w", queryErrorPosition(err), err)
	}

	runCtx, runCtxCancelFn := i.interruptStack.Push(ctx)
	ni.ctx = runCtx
	ni.stdout = CtxOutput{Output: stdout, Ctx: runCtx}
	iter := gc.RunWithContext(runCtx, c)

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

func (i *Interp) EvalFunc(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output, debugFn string) (gojq.Iter, error) {
	var argsExpr []string
	for i := range args {
		argsExpr = append(argsExpr, fmt.Sprintf("$a[%d]", i))
	}
	argExpr := ""
	if len(argsExpr) > 0 {
		argExpr = "(" + strings.Join(argsExpr, ";") + ")"
	}

	trampolineInput := map[string]interface{}{
		"input": c,
		"args":  args,
	}
	// {input: ..., args: [...]} | .args as $a | .input | name[($a[0]; ...)]
	trampolineExpr := fmt.Sprintf(".args as $a | .input | %s%s", name, argExpr)
	iter, err := i.Eval(ctx, mode, trampolineInput, trampolineExpr, stdout, debugFn)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (i *Interp) EvalFuncValues(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output, debugFn string) ([]interface{}, error) {
	iter, err := i.EvalFunc(ctx, mode, c, name, args, stdout, debugFn)
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
	Depth      int    `json:"depth"`
	Verbose    bool   `json:"verbose"`
	Color      bool   `json:"color"`
	Colors     string `json:"colors"`
	ByteColors string `json:"bytecolors"`
	Unicode    bool   `json:"unicode"`
	Raw        bool   `json:"raw"`
	REPL       bool   `json:"repl"`
	RawString  bool   `json:"rawstring"`
	Compact    bool   `json:"compact"`

	LineBytes    int `json:"linebytes"`
	DisplayBytes int `json:"displaybytes"`
	AddrBase     int `json:"addrbase"`
	SizeBase     int `json:"sizebase"`

	REPLLevel int `json:"repllevel"`

	Decorator Decorator `json:"-"`
}

func mapSetOptions(d *Options, m map[string]interface{}) {
	if v, ok := m["depth"]; ok {
		d.Depth = num.MaxInt(0, toIntZ(v))
	}
	if v, ok := m["verbose"]; ok {
		d.Verbose = toBoolZ(v)
	}
	if v, ok := m["color"]; ok {
		d.Color = toBoolZ(v)
	}
	if v, ok := m["colors"]; ok {
		d.Colors = toStringZ(v)
	}
	if v, ok := m["bytecolors"]; ok {
		d.ByteColors = toStringZ(v)
	}
	if v, ok := m["unicode"]; ok {
		d.Unicode = toBoolZ(v)
	}
	if v, ok := m["raw"]; ok {
		d.Raw = toBoolZ(v)
	}
	if v, ok := m["repl"]; ok {
		d.REPL = toBoolZ(v)
	}
	if v, ok := m["rawstring"]; ok {
		d.RawString = toBoolZ(v)
	}
	if v, ok := m["compact"]; ok {
		d.Compact = toBoolZ(v)
	}

	if v, ok := m["linebytes"]; ok {
		d.LineBytes = num.MaxInt(0, toIntZ(v))
	}
	if v, ok := m["displaybytes"]; ok {
		d.DisplayBytes = num.MaxInt(0, toIntZ(v))
	}
	if v, ok := m["addrbase"]; ok {
		d.AddrBase = num.ClampInt(2, 36, toIntZ(v))
	}
	if v, ok := m["sizebase"]; ok {
		d.SizeBase = num.ClampInt(2, 36, toIntZ(v))
	}

	if v, ok := m["repllevel"]; ok {
		d.REPLLevel = toIntZ(v)
	}
}

func (i *Interp) Options(fnOptsV ...interface{}) (Options, error) {
	var opts Options

	defaultOptsV := i.state["default_options"]
	if defaultOptsV == nil {
		return Options{}, fmt.Errorf("default_options state not set")
	}
	defaultOpts, ok := defaultOptsV.(map[string]interface{})
	if !ok {
		return Options{}, fmt.Errorf("default_options not an object")
	}
	mapSetOptions(&opts, defaultOpts)

	optsStackV := i.state["options_stack"]
	if optsStackV == nil {
		return Options{}, fmt.Errorf("options_stack state not set")
	}
	optsStack, ok := optsStackV.([]interface{})
	if !ok {
		return Options{}, fmt.Errorf("options_stack is not an array")
	}
	for i := len(optsStack) - 1; i >= 0; i-- {
		ov := optsStack[i]
		o, ok := ov.(map[string]interface{})
		if !ok {
			return Options{}, fmt.Errorf("optsStack[%d] not an object: %v", i, ov)
		}
		mapSetOptions(&opts, o)
	}

	for _, fnOptsV := range fnOptsV {
		fnOpts, ok := fnOptsV.(map[string]interface{})
		if !ok {
			return Options{}, fmt.Errorf("options not an object")
		}
		mapSetOptions(&opts, fnOpts)
	}

	opts.Decorator = decoratorFromOptions(opts)

	return opts, nil
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
			if o, ok := v.(gojq.JQValue); ok {
				return o.JQValue()
			}
			return nil
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
