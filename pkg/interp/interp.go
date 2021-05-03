package interp

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"fq"
	"fq/internal/ctxstack"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/ranges"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

const builtinPrefix = "@builtin"

//go:embed *.jq
var builtinFS embed.FS

//go:embed fq.jq
var fqJq []byte

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
	// returned io.ReadSeeker can optionally implement io.Closer
	Open(name string) (io.ReadSeeker, error)
	Readline(prompt string, complete func(line string, pos int) (newLine []string, shared int)) (string, error)
}

// TODO: would be nice if gojq had something for this? maybe missing something?
func offsetToLine(s string, offset int) int {
	co := 0
	line := 1
	for {
		no := strings.Index(s[co:], "\n")
		if no == -1 || co+no >= offset {
			return line
		}
		co += no + 1
		line++
	}
}

func queryErrorLine(v error) int {
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
		return offsetToLine(content, offset)
	}
	return 0
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

func valuePathDecorated(v *decode.Value, d Decorator) string {
	var parts []string

	for v.Parent != nil {
		switch v.Parent.V.(type) {
		case decode.Struct:
			parts = append([]string{".", d.ObjectKey.Wrap(v.Name)}, parts...)
		case decode.Array:
			indexStr := strconv.Itoa(v.Index)
			parts = append([]string{fmt.Sprintf("%s%s%s", d.Index.F("["), d.Number.F(indexStr), d.Index.F("]"))}, parts...)
		}
		v = v.Parent
	}

	if len(parts) == 0 {
		return "."
	}

	return strings.Join(parts, "")
}

// TODO: jq function somehow? escape keys?
func valuePath(v *decode.Value) string {
	return valuePathDecorated(v, PlainDecorator)
}

type EmptyError interface {
	IsEmptyError() bool
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

func toInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case *big.Int:
		return v.Int64(), nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("value is not a number")
	}
}

func toInt64Z(v interface{}) int64 {
	n, _ := toInt64(v)
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

// TODO: refactor to return struct?
func toBuffer(v interface{}) (*bitio.Buffer, error) {
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
		return bitio.NewBufferFromBytes(bi.Bytes(), -1), nil
	case []interface{}:
		var rr []bitio.BitReadAtSeeker
		for _, e := range vv {
			eBB, eErr := toBuffer(e)
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
		switch vv := v.(type) {
		case ToBuffer:
			bb, err := vv.ToBuffer()
			if err != nil {
				return bufferRange{}, err
			}
			return bufferRange{bb: bb, r: ranges.Range{Len: bb.Len()}}, nil
		default:
			return bufferRange{}, fmt.Errorf("value can't be buffer")
		}
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

type Variable struct {
	Name  string
	Value interface{}
}

type Function struct {
	Names     []string
	MinArity  int
	MaxArity  int
	Fn        func(interface{}, []interface{}) interface{}
	Generator bool
}

type RunMode int

const (
	ScriptMode RunMode = iota
	REPLMode
	CompletionMode
)

type runContext struct {
	ctx    context.Context
	stdout Output // TODO: rename?
	mode   RunMode
	state  map[string]interface{}
}

type Interp struct {
	// variables map[string]interface{}
	registry *decode.Registry
	os       OS

	builtinQueryCache map[string]*gojq.Query
	includeFqQuery    *gojq.Query
	interruptStack    *ctxstack.Stack

	// new for each run, other values are copied by ref
	runContext
}

func New(os OS, registry *decode.Registry) (*Interp, error) {
	var err error

	i := &Interp{
		os:       os,
		registry: registry,
	}

	i.builtinQueryCache = map[string]*gojq.Query{}
	i.includeFqQuery, err = gojq.Parse(string(fqJq))
	if err != nil {
		return nil, fmt.Errorf("%d: %w", queryErrorLine(err), err)
	}
	i.interruptStack = ctxstack.New(func(closeCh chan struct{}) {
		select {
		case <-closeCh:
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

func (i *Interp) Main(ctx context.Context, stdout io.Writer) error {
	runMode := ScriptMode

	var args []interface{}
	for _, a := range i.os.Args() {
		args = append(args, a)
	}

	input := map[string]interface{}{
		"args":    args,
		"version": fq.Version,
	}

	iter, err := i.EvalFunc(ctx, runMode, input, "main", nil, i.os.Stdout())
	if err != nil {
		fmt.Fprintln(i.os.Stderr(), err)
		return err
	}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		} else if err, ok := v.(error); ok {
			fmt.Fprintln(i.os.Stderr(), err)
			return err
		} else if d, ok := v.([2]interface{}); ok {
			fmt.Fprintln(i.os.Stderr(), d[:]...)
		}
	}

	return nil
}

func (i *Interp) Eval(ctx context.Context, mode RunMode, c interface{}, src string, stdout Output) (gojq.Iter, error) {
	var err error

	// TODO: did not work
	// nq := &(*q)

	gq, err := gojq.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("%d: %w", queryErrorLine(err), err)
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
		state: newState,
		mode:  mode,
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
			if f.Generator {
				compilerOpts = append(compilerOpts,
					gojq.WithIterator(n, f.MinArity, f.MaxArity, f.Fn))
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
			return []*gojq.Query{i.includeFqQuery}, nil
		},
		load: func(name string) (*gojq.Query, error) {
			parts := strings.Split(name, "/")

			if len(parts) > 0 && parts[0] == builtinPrefix {
				name = strings.Join(parts[1:], "/")
				if q, ok := ni.builtinQueryCache[name]; ok {
					return q, nil
				}
				b, err := builtinFS.ReadFile(name)
				if err != nil {
					return nil, err
				}
				mq, err := gojq.Parse(string(b))
				if err != nil {
					return nil, err
				}
				ni.builtinQueryCache[name] = mq
				return mq, nil
			}

			return nil, fmt.Errorf("module not found: %q", name)
		},
	}))

	gc, err := gojq.Compile(gq, compilerOpts...)
	if err != nil {
		return nil, fmt.Errorf("%d: %w", queryErrorLine(err), err)
	}

	runCtx, runCtxCancelFn := i.interruptStack.Push(ctx)
	ni.ctx = runCtx
	ni.stdout = CtxOutput{Output: stdout, Ctx: runCtx}

	iter := gc.RunWithContext(runCtx, c)
	// iter := gc.RunWithContext(ctx, c, variableValues...)

	iterWrapper := iterFn(func() (interface{}, bool) {
		v, ok := iter.Next()
		_, isErr := v.(error)
		if !ok || isErr {
			runCtxCancelFn()
		}

		return v, ok
	})

	return iterWrapper, nil
}

func (i *Interp) EvalFunc(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output) (gojq.Iter, error) {
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
	iter, err := i.Eval(ctx, mode, trampolineInput, trampolineExpr, stdout)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (i *Interp) EvalFuncValues(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output) ([]interface{}, error) {
	iter, err := i.EvalFunc(ctx, mode, c, name, args, stdout)
	if err != nil {
		return nil, err
	}

	var vs []interface{}
	for {
		v, ok := iter.Next()
		_, isErr := v.(error)
		vs = append(vs, v)
		if !ok || isErr {
			break
		}
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

	LineBytes    int   `json:"linebytes"`
	DisplayBytes int64 `json:"displaybytes"`
	AddrBase     int   `json:"addrbase"`
	SizeBase     int   `json:"sizebase"`

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
