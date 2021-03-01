package interp

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"fq"
	"fq/internal/ansi"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/ranges"
	"io"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"

	"github.com/itchyny/gojq"
)

const builtinPrefix = "@builtin"

//go:embed *.jq
var builtinFS embed.FS

//go:embed fq.jq
var fqJq []byte

type Output interface {
	io.Writer
	Size() (int, int)
	IsTerminal() bool
}

type OS interface {
	Stdin() io.Reader
	Stdout() Output
	Stderr() io.Writer
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
	if v, ok := m["repl"]; ok {
		d.REPL = toBoolZ(v)
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

type emptyIter struct{}

func (emptyIter) Next() (interface{}, bool) { return nil, false }

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

type Display interface {
	Display(w io.Writer, opts DisplayOptions) error
}

type ToBitBuf interface {
	ToBifBuf() *bitio.Buffer
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
	Names    []string
	MinArity int
	MaxArity int
	Fn       func(interface{}, []interface{}) interface{}
}

type RunMode int

const (
	ScriptMode RunMode = iota
	REPLMode
	CompletionMode
)

type evalContext struct {
	ctx      context.Context
	optsExpr map[string]interface{}
	opts     map[string]interface{}
	stdout   Output // TODO: rename?
	mode     RunMode
	inEval   bool
}

type Interp struct {
	variables map[string]interface{}
	registry  *decode.Registry
	os        OS

	builtinQueryCache map[string]*gojq.Query
	includeFqQuery    *gojq.Query

	evalContext *evalContext
}

func New(opts InterpOptions) (*Interp, error) {
	var err error

	i := &Interp{
		variables: opts.Variables,
		registry:  opts.Registry,
		os:        opts.OS,
	}

	// TODO: cleanup group names and panics

	i.builtinQueryCache = map[string]*gojq.Query{}
	i.evalContext = &evalContext{
		optsExpr: map[string]interface{}{},
		opts:     map[string]interface{}{},
	}
	i.includeFqQuery, err = gojq.Parse(string(fqJq))
	if err != nil {
		return nil, fmt.Errorf("%d: %w", queryErrorLine(err), err)
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

	iter, err := i.Eval(ctx, runMode, input, "main", i.os.Stdout(), nil)
	if err != nil {
		log.Printf("err: %#+v\n", err)
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

	// make copy of query
	ci := *i
	ni := &ci
	if optsExpr == nil {
		optsExpr = map[string]interface{}{}
	}
	ni.evalContext = &evalContext{
		ctx:      ctx,
		mode:     mode,
		optsExpr: optsExpr,
		opts:     i.evalContext.opts,
		inEval:   false,
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

	opts := buildDisplayOptions(i.evalContext.opts)
	cleanupFn := func() {}
	stdoutCtx := ctx

	if opts.REPL {
		i.evalContext.inEval = true
		interruptChan := make(chan os.Signal, 1)
		signal.Notify(interruptChan, os.Interrupt)
		interruptCtx, interruptCtxCancelFn := context.WithCancel(ctx)
		stdoutCtx = interruptCtx
		go func() {
			select {
			case <-interruptChan:
				if !ni.evalContext.inEval {
					interruptCtxCancelFn()
				}
			case <-interruptCtx.Done():
				// nop
			}
		}()
		cleanupFn = func() {
			signal.Stop(interruptChan)
			// stop interruptChan goroutine
			interruptCtxCancelFn()
			i.evalContext.inEval = false
		}
	}

	ni.evalContext.stdout = CtxOutput{Output: stdout, Ctx: stdoutCtx}

	iter := gc.RunWithContext(ctx, c, variableValues...)

	iterCtxWrapped := iterFn(func() (interface{}, bool) {
		v, ok := iter.Next()
		if v == context.Canceled {
			cleanupFn()
			return nil, false
		}
		if !ok {
			cleanupFn()
		}
		return v, ok
	})

	return iterCtxWrapped, nil
}

func (i *Interp) EvalFunc(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output, optsExpr map[string]interface{}) (gojq.Iter, error) {
	var argsJSON []string
	for _, arg := range args {
		b, err := json.Marshal(arg)
		if err != nil {
			return nil, err
		}
		argsJSON = append(argsJSON, string(b))
	}
	argsStr := ""
	if len(argsJSON) > 0 {
		argsStr = "(" + strings.Join(argsJSON, ";") + ")"
	}

	iter, err := i.Eval(ctx, mode, c, fmt.Sprintf("%s%s", name, argsStr), stdout, optsExpr)
	if err != nil {
		return nil, err
	}
	return iter, nil
}

func (i *Interp) EvalFuncValue(ctx context.Context, mode RunMode, c interface{}, name string, args []interface{}, stdout Output, optsExpr map[string]interface{}) interface{} {
	iter, err := i.EvalFunc(ctx, mode, c, name, args, stdout, optsExpr)
	if err != nil {
		return err
	}
	v, _ := iter.Next()
	return v
}
