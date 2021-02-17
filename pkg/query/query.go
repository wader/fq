package query

// TODO: rename to context etc? env?
// TODO: per run context?

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/osenv"
	"fq/pkg/ranges"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

func jsonEscape(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func valueToTypeString(v interface{}) (string, bool) {
	switch v.(type) {
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64, complex64, complex128, uintptr, *big.Int:
		return "number", true
	case bool:
		return "boolean", true
	case string:
		return "string", true
	}
	return "?", false
}

type EmptyError interface {
	IsEmptyError() bool
}

type iterFn func() (interface{}, bool)

func (i iterFn) Next() (interface{}, bool) { return i() }

type autoCompleterFn func(line []rune, pos int) (newLine [][]rune, length int)

func (a autoCompleterFn) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return a(line, pos)
}

type loadModuleFn func(name string) (*gojq.Query, error)

func (l loadModuleFn) LoadModule(name string) (*gojq.Query, error) {
	return l(name)
}

type listenerFn func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool)

func (lf listenerFn) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	return lf(line, pos, key)
}

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

type ToBitBuf interface {
	ToBifBuf() *bitio.Buffer
}

// TODO: refactor to return struct?
func toBitBuf(v interface{}) (*bitio.Buffer, ranges.Range, string, error) {
	switch vv := v.(type) {
	case *bitBufFile:
		return vv.bb, ranges.Range{Start: 0, Len: vv.bb.Len()}, vv.filename, nil
	case *decode.Value:
		return vv.RootBitBuf, vv.Range, "", nil
	case *bitio.Buffer:
		return vv, ranges.Range{Start: 0, Len: vv.Len()}, "", nil
	case []byte:
		bb := bitio.NewBufferFromBytes(vv, -1)
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
	case ToBitBuf:
		bb := vv.ToBifBuf()
		return bb, ranges.Range{Start: 0, Len: bb.Len()}, "", nil
	default:
		return nil, ranges.Range{}, "", fmt.Errorf("value should be decode value, bit buffer, byte slice or string")
	}
}

func toValue(v interface{}) (*decode.Value, error) {
	switch v := v.(type) {
	case *decode.Value:
		return v, nil
	case *decode.D:
		// TODO: remove decode.D?
		return v.Value, nil
	default:
		// TODO: remove decode.D?
		return nil, fmt.Errorf("%v: value is not a decode value", v)
	}
}

type QueryOptions struct {
	Variables map[string]interface{}
	Registry  *decode.Registry
	Options   map[string]string
	OS        osenv.OS
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

type Query struct {
	opts              QueryOptions
	inputStack        [][]interface{}
	variables         map[string]interface{}
	functions         []Function
	runContext        *runContext
	builtinQueryCache map[string]*gojq.Query
}

type bitBufFile struct {
	bb       *bitio.Buffer
	filename string

	decodeDoneFn func()
}

type RunMode int

const (
	ScriptMode RunMode = iota
	REPLMode
	CompletionMode
)

type runContext struct {
	ctx    context.Context
	mode   RunMode
	stdout Output // TODO: rename?
	opts   map[string]interface{}

	pushVs []interface{}
	pops   int
}

type queryErrorFn func(stdout io.Writer) error

func (queryErrorFn) Error() string { return "" }

func NewQuery(opts QueryOptions) *Query {
	q := &Query{opts: opts}

	// TODO: cleanup group names and panics
	q.functions = q.makeFunctions(opts)
	// TODO: redo args handling in jq? a cli_entry function that reads args?
	q.variables = opts.Variables
	q.builtinQueryCache = map[string]*gojq.Query{}

	return q
}

func (q *Query) Run(ctx context.Context, mode RunMode, src string, stdout Output) ([]interface{}, error) {
	var err error

	q.runContext = &runContext{
		ctx:    ctx,
		mode:   mode,
		stdout: stdout,
		opts:   map[string]interface{}{},
	}

	optsExpr := "{"
	for k, v := range q.opts.Options {
		optsExpr += fmt.Sprintf(`"%s": (%s),`, k, v)
	}
	optsExpr += "}"

	runQuery := fmt.Sprintf(`include "%s/fq.jq"; options(%s) | inputs`, builtinPrefix, optsExpr)
	if src != "" {
		runQuery += `| ` + src
	}

	query, err := gojq.Parse(runQuery)
	if err != nil {
		return nil, err
	}

	var variableNames []string
	var variableValues []interface{}
	for k, v := range q.variables {
		variableNames = append(variableNames, k)
		variableValues = append(variableValues, v)
	}

	var compilerOpts []gojq.CompilerOption
	for _, f := range q.functions {
		for _, n := range f.Names {
			compilerOpts = append(compilerOpts,
				gojq.WithFunction(n, f.MinArity, f.MaxArity, f.Fn))
		}
	}
	compilerOpts = append(compilerOpts, gojq.WithVariables(variableNames))
	var inputs []interface{}
	if len(q.inputStack) > 0 {
		inputs = q.inputStack[len(q.inputStack)-1]
	} else {
		// TODO: hmm
		inputs = []interface{}{nil}
	}
	compilerOpts = append(compilerOpts, gojq.WithEnvironLoader(q.opts.OS.Environ))
	compilerOpts = append(compilerOpts, gojq.WithInputIter(iterFn(func() (interface{}, bool) {
		if len(inputs) == 0 {
			return nil, false
		}
		var input interface{}
		input, inputs = inputs[0], inputs[1:]
		return input, true
	})))
	compilerOpts = append(compilerOpts, gojq.WithModuleLoader(loadModuleFn(func(name string) (*gojq.Query, error) {
		parts := strings.Split(name, "/")

		if len(parts) > 0 && parts[0] == builtinPrefix {
			name = strings.Join(parts[1:], "/")
			if q, ok := q.builtinQueryCache[name]; ok {
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
			q.builtinQueryCache[name] = mq
			return mq, nil
		}

		return nil, fmt.Errorf("module not found: %q", name)
	})))

	code, err := gojq.Compile(query, compilerOpts...)
	if err != nil {
		return nil, err
	}

	iter := code.RunWithContext(ctx, nil, variableValues...)

	var vs []interface{}
	for {
		var ok bool
		var v interface{}

		if v, ok = iter.Next(); !ok {
			break
		}
		if err, ok = v.(error); ok {
			switch ee := err.(type) {
			case EmptyError:
				if ee.IsEmptyError() {
					err = nil
					continue
				}
			case queryErrorFn:
				return nil, ee(stdout)
			}
			break
		}

		vs = append(vs, v)

		switch vv := v.(type) {
		case func(stdout io.Writer) error:
			if err := vv(stdout); err != nil {
				return nil, err
			}
		case *bitBufFile:
			fmt.Fprintf(stdout, "<file %s>\n", vv.filename)
		case *decode.Value:
			if err := vv.Dump(stdout, buildDumpOptions(q.runContext.opts, map[string]interface{}{
				"maxdepth": 1,
			})); err != nil {
				return nil, err
			}
		case *decode.D:
			// TODO: remove?
			if err := vv.Value.Dump(stdout, buildDumpOptions(q.runContext.opts)); err != nil {
				return nil, err
			}

		case ToBitBuf:
			bb := vv.ToBifBuf()
			if _, err := io.Copy(stdout, bb.Copy()); err != nil {
				return nil, err
			}
		case *bitio.Buffer:
			if _, err := io.Copy(stdout, vv.Copy()); err != nil {
				return nil, err
			}
		case string, int, int32, int64, uint, uint32, uint64:
			fmt.Fprintln(stdout, vv)
		case float32:
			// TODO: should not happen?
			fmt.Fprintln(stdout, strconv.FormatFloat(float64(vv), 'f', -1, 32))
		case float64:
			fmt.Fprintln(stdout, strconv.FormatFloat(vv, 'f', -1, 64))
		case []byte:
			if _, err := io.Copy(stdout, bytes.NewBuffer(vv)); err != nil {
				return nil, err
			}
		default:
			e := json.NewEncoder(stdout)
			e.SetIndent("", "  ")
			if err := e.Encode(v); err != nil {
				return nil, err
			}
		}

	}

	if q.runContext.pops > 0 && len(q.inputStack) > 0 {
		q.inputStack = q.inputStack[0 : len(q.inputStack)-1]
	}

	if q.runContext.pushVs != nil {
		// TODO: use vs?
		q.inputStack = append(q.inputStack, q.runContext.pushVs)
	}

	return vs, err
}
