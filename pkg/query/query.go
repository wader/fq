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
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
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
func toBool(v interface{}) (bool, error) {
	switch v := v.(type) {
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
	Variables   map[string]interface{}
	Registry    *decode.Registry
	DumpOptions decode.DumpOptions
	OS          osenv.OS
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
	opts       QueryOptions
	inputStack [][]interface{}
	variables  map[string]interface{}
	functions  []Function
	runContext *runContext
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
	mode   RunMode
	pushVs []interface{}
	pops   int
	stdout io.Writer
}

type queryErrorFn func(stdout io.Writer) error

func (queryErrorFn) Error() string { return "" }

func NewQuery(opts QueryOptions) *Query {
	q := &Query{opts: opts}

	// TODO: cleanup group names and panics
	q.functions = q.makeFunctions(opts)
	// TODO: redo args handling in jq? a cli_entry function that reads args?
	q.variables = opts.Variables

	return q
}

func (q *Query) Run(ctx context.Context, mode RunMode, src string, stdout io.Writer) ([]interface{}, error) {
	var err error

	q.runContext = &runContext{
		mode:   mode,
		stdout: stdout,
	}

	if src != "" {
		src = `include "fq" ; inputs | ` + src
	} else {
		src = `include "fq" ; inputs`
	}

	query, err := gojq.Parse(src)
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
	compilerOpts = append(compilerOpts, gojq.WithInputIter(iterFn(func() (interface{}, bool) {
		if len(inputs) == 0 {
			return nil, false
		}
		var input interface{}
		input, inputs = inputs[0], inputs[1:]
		return input, true
	})))
	compilerOpts = append(compilerOpts, gojq.WithModuleLoader(loadModuleFn(func(name string) (*gojq.Query, error) {
		switch name {
		case "fq":
			return gojq.Parse(`
				# convert number to array of bytes
				def number_to_bytes($bits):
					def _number_to_bytes($d):
						if . > 0 then
							. % $d, (. div $d | _number_to_bytes($d))
						else
							empty
						end;
					if . == 0 then [0]
					else [_number_to_bytes(1 bsl $bits)] | reverse end;
				def number_to_bytes:
					number_to_bytes(8);

				# from https://rosettacode.org/wiki/Non-decimal_radices/Convert#jq
				# unknown author
				# Convert the input integer to a string in the specified base (2 to 36 inclusive)
				def _convert(base):
					def stream:
						recurse(if . > 0 then . div base else empty end) | . % base;
					if . == 0 then
						"0"
					else
						[stream] |
						reverse  |
						.[1:] |
						if base <  10 then
							map(tostring) | join("")
						elif base <= 36 then
							map(if . < 10 then 48 + . else . + 87 end) | implode
						else
							error("base too large")
						end
					end;

				# input string is converted from "base" to an integer, within limits
				# of the underlying arithmetic operations, and without error-checking:
				def _to_i(base):
					explode
					| reverse
					| map(if . > 96  then . - 87 else . - 48 end)  # "a" ~ 97 => 10 ~ 87
					| reduce .[] as $c
						# state: [power, ans]
						([1,0]; (.[0] * base) as $b | [$b, .[1] + (.[0] * $c)])
					| .[1];

				# like iprint
				def i:
					{
						bin: "0b\(_convert(2))",
						oct: "0o\(_convert(8))",
						dec: "\(_convert(10))",
						hex: "0x\(_convert(16))",
						str: ([.] | implode),
					};

				def _formats_dot:
					"# ... | dot -Tsvg -o formats.svg",
					"digraph formats {",
					"  nodesep=0.5",
					"  ranksep=0.5",
					"  node [shape=\"box\",style=\"rounded,filled\"]",
					"  edge [arrowsize=\"0.7\"]",
					(.[] | "  \(.name) -> {\(.dependencies | flatten? | join(" "))}"),
					(.[] | .name as $name | .groups[]? | "  \(.) -> \($name)"),
					(keys[] | "  \(.) [color=\"paleturquoise\"]"),
					([.[].groups[]?] | unique[] | "  \(.) [color=\"palegreen\"]"),
					"}";

				def field_inrange($p): ._type == "field" and ._range.start <= $p and $p < ._range.stop;

			`)
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
			opts := q.opts.DumpOptions
			opts.MaxDepth = 1
			if err := vv.Dump(stdout, opts); err != nil {
				return nil, err
			}
		case *decode.D:
			// TODO: remove?
			if err := vv.Value.Dump(stdout, q.opts.DumpOptions); err != nil {
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

// REPL read-eval-print-loop
func (q *Query) REPL(ctx context.Context) error {
	// TODO: refactor
	historyFile := ""
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	historyFile = filepath.Join(cacheDir, "fq/history")
	_ = os.MkdirAll(filepath.Dir(historyFile), 0700)

	// log := log.New(func() io.Writer { f, _ := os.Create("/tmp/log"); return f }(), "", 0)

	l, err := readline.NewEx(&readline.Config{
		Stdin:       ioutil.NopCloser(q.opts.OS.Stdin()),
		Stdout:      q.opts.OS.Stdout(),
		Stderr:      q.opts.OS.Stderr(),
		HistoryFile: historyFile,
		AutoComplete: autoCompleterFn(func(line []rune, pos int) (newLine [][]rune, length int) {
			completeCtx, completeCtxCancelFn := context.WithTimeout(ctx, 1*time.Second)
			defer completeCtxCancelFn()
			return autoComplete(completeCtx, q, line, pos)
		}),
		// InterruptPrompt: "^C",
		// EOFPrompt:       "exit",

		HistorySearchFold: true,
		// FuncFilterInputRune: filterInput,
	})
	if err != nil {
		return err
	}

	for {
		var v []interface{}
		if len(q.inputStack) > 0 {
			v = q.inputStack[len(q.inputStack)-1]
		}
		var inputSummary []string
		if len(v) > 0 {
			first := v[0]
			if vv, ok := first.(*decode.Value); ok {
				inputSummary = append(inputSummary, vv.Path())
			} else if t, ok := valueToTypeString(first); ok {
				inputSummary = append(inputSummary, t)
			} else {
				inputSummary = append(inputSummary, "?")
			}
		}
		if len(v) > 1 {
			inputSummary = append(inputSummary, "...")
		}
		prompt := fmt.Sprintf("inputs[%d] [%s]> ", len(q.inputStack), strings.Join(inputSummary, ","))

		l.SetPrompt(prompt)

		src, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(src) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if _, err := q.Run(ctx, REPLMode, src, q.opts.OS.Stdout()); err != nil {
			fmt.Fprintf(q.opts.OS.Stdout(), "error: %s\n", err)
		}
	}

	return nil
}
