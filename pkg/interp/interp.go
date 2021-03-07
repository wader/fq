package interp

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"fq"
	"fq/internal/ansi"
	"fq/internal/ctxstack"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/ranges"
	"io"
	"log"
	"math/big"
	"strings"

	"github.com/itchyny/gojq"
)

func contextWithChan(ctx context.Context, c chan struct{}) (context.Context, func()) {
	chanCtx, cancelFn := context.WithCancel(ctx)
	chanCtxCancelChan := make(chan struct{})
	go func() {
		select {
		case <-chanCtxCancelChan:
			log.Println("chanctx cancel")
			return
		case <-c:
			log.Println("chanctx got c chancel")
			cancelFn()
		}
	}()

	return chanCtx, func() {
		close(chanCtxCancelChan)
	}
}

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

	if tokif, ok := v.(interface{ Token() (string, int) }); ok {
		_, offset = tokif.Token()
	}
	if qeif, ok := v.(interface {
		QueryParseError() (string, string, string, error)
	}); ok {
		_, _, content, _ = qeif.QueryParseError()
	}

	if offset > 0 && content != "" {
		return offsetToLine(content, offset)
	}
	return 0
}

type DisplayOptions struct {
	Depth     int
	Verbose   bool
	Color     bool
	Unicode   bool
	Raw       bool
	REPL      bool
	RawString bool

	LineBytes    int
	DisplayBytes int64
	AddrBase     int
	SizeBase     int

	Decorator Decorator
}

// TODO: rename, not only display things
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
	if v, ok := m["depth"]; ok {
		d.Depth = num.MaxInt(0, toIntZ(v))
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
}

func decoratorFromDumpOptions(opts DisplayOptions) Decorator {
	colStr := "|"
	if opts.Unicode {
		colStr = "\xe2\x94\x82"
	}
	nameFn := func(s string) string { return s }
	valueFn := func(s string) string { return s }
	frameFn := func(s string) string { return s }
	byteFn := func(b byte, s string) string { return s }
	column := colStr + "\n"
	if opts.Color {
		nameFn = func(s string) string { return ansi.FgBrightBlue + s + ansi.Reset }
		valueFn = func(s string) string { return ansi.FgBrightCyan + s + ansi.Reset }
		frameFn = func(s string) string { return ansi.FgYellow + s + ansi.Reset }
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
		Frame:  frameFn,
		Byte:   byteFn,
		Column: column,
	}
}

type Decorator struct {
	Name   func(s string) string
	Value  func(s string) string
	Frame  func(s string) string
	Byte   func(b byte, s string) string
	Column string
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
	gojq.JSONObject

	DisplayName() string
	SpecialPropNames() []string
}

type Display interface {
	Display(w io.Writer, opts DisplayOptions) error
}

type Preview interface {
	Preview(w io.Writer, opts DisplayOptions) error
}

type ToBitBuf interface {
	ToBifBuf() *bitio.Buffer
}

// TODO: jq function somehow? escape keys?
func valuePath(v *decode.Value) string {
	var parts []string

	for v.Parent != nil {
		switch v.Parent.V.(type) {
		case decode.Struct:
			parts = append([]string{".", v.Name}, parts...)
		case decode.Array:
			parts = append([]string{fmt.Sprintf("[%d]", v.Index)}, parts...)
		}
		v = v.Parent
	}

	if len(parts) == 0 {
		return "."
	}

	return strings.Join(parts, "")

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
		bb, _, _, err := toBitBuf(v)
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
func toBitBuf(v interface{}) (*bitio.Buffer, ranges.Range, string, error) {
	switch vv := v.(type) {
	case ToBitBuf:
		bb := vv.ToBifBuf()
		return bb, ranges.Range{Start: 0, Len: bb.Len()}, "", nil
	case string:
		bb := bitio.NewBufferFromBytes([]byte(vv), -1)
		return bb, ranges.Range{Start: 0, Len: bb.Len()}, "", nil
	case int, float64, *big.Int:
		bi, err := toBigInt(v)
		if err != nil {
			return nil, ranges.Range{}, "", err
		}
		bb := bitio.NewBufferFromBytes(bi.Bytes(), -1)
		return bb, ranges.Range{Start: 0, Len: bb.Len()}, "", nil
	default:
		return nil, ranges.Range{}, "", fmt.Errorf("value should be decode value, bit buffer, byte slice or string")
	}
}

func toValue(v interface{}) interface{} {
	switch v := v.(type) {
	case gojq.JSONObject:
		return v.JsonPrimitiveValue()
	case nil, bool, float64, int, string, *big.Int, map[string]interface{}, []interface{}:
		return v
	default:
		return nil
	}
}

type InterpOptions struct {
	Variables map[string]interface{}
	Registry  *decode.Registry
	OS        OS
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
	ctx      context.Context
	optsExpr map[string]interface{}
	opts     map[string]interface{}
	stdout   Output // TODO: rename?
	mode     RunMode
}

type Interp struct {
	variables map[string]interface{}
	registry  *decode.Registry
	os        OS

	builtinQueryCache map[string]*gojq.Query
	includeFqQuery    *gojq.Query
	interruptStack    *ctxstack.Stack

	// new for each run other values are copied
	runContext
}

func New(opts InterpOptions) (*Interp, error) {
	var err error

	i := &Interp{
		variables: opts.Variables,
		registry:  opts.Registry,
		os:        opts.OS,
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
		case <-opts.OS.Interrupt():
			return
		}
	})
	i.runContext = runContext{
		optsExpr: map[string]interface{}{},
		opts:     map[string]interface{}{},
	}

	return i, nil
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

	iter, err := i.EvalFunc(ctx, runMode, input, "main", nil, i.os.Stdout(), nil)
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
			fmt.Fprintf(i.os.Stderr(), "%s: %v\n", d[0], d[1])
		}
	}

	return nil
}

func (i *Interp) Eval(ctx context.Context, mode RunMode, c interface{}, src string, stdout Output, optsExpr map[string]interface{}) (gojq.Iter, error) {
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
	if optsExpr == nil {
		optsExpr = map[string]interface{}{}
	}
	ni.runContext = runContext{
		mode:     mode,
		optsExpr: optsExpr,
		opts:     i.opts,
	}

	var variableNames []string
	var variableValues []interface{}
	for k, v := range ni.variables {
		variableNames = append(variableNames, k)
		variableValues = append(variableValues, v)
	}

	var compilerOpts []gojq.CompilerOption
	for _, f := range ni.makeFunctions(ni.registry) {
		for _, n := range f.Names {
			compilerOpts = append(compilerOpts,
				gojq.WithFunction(n, f.MinArity, f.MaxArity, f.Fn))
		}
	}
	compilerOpts = append(compilerOpts, gojq.WithEnvironLoader(ni.os.Environ))
	compilerOpts = append(compilerOpts, gojq.WithVariables(variableNames))
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

	ni.stdout = CtxOutput{Output: stdout, Ctx: runCtx}
	ni.ctx = runCtx
	iter := gc.RunWithContext(ctx, c, variableValues...)

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

func (i *Interp) EvalFunc(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output, optsExpr map[string]interface{}) (gojq.Iter, error) {
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
	iter, err := i.Eval(ctx, mode, trampolineInput, trampolineExpr, stdout, optsExpr)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (i *Interp) EvalFuncValues(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output, optsExpr map[string]interface{}) ([]interface{}, error) {
	iter, err := i.EvalFunc(ctx, mode, c, name, args, stdout, optsExpr)
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
