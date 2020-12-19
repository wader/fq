package query

// TODO: rename to context etc? env?

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"fq/internal/hexdump"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/format"
	"fq/pkg/osenv"
	"io"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

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

func toBB(v interface{}) (*bitio.Buffer, string, error) {
	var bb *bitio.Buffer
	switch vv := v.(type) {
	case *queryBB:
		return vv.bb, vv.filename, nil
	case *decode.Value:
		var err error
		bb, err = vv.RootBitBuf.BitBufRange(vv.Range.Start, vv.Range.Len)
		if err != nil {
			return nil, "", err
		}
	case *bitio.Buffer:
		bb = vv
	case []byte:
		bb = bitio.NewBufferFromBytes(vv, -1)
	case string:
		bb = bitio.NewBufferFromBytes([]byte(vv), -1)
	default:
		return nil, "", fmt.Errorf("value should be decode value, bit buffer, byte slice or string")
	}

	return bb, "", nil
}

type QueryOptions struct {
	Variables   []Variable
	FormatName  string
	Filename    string
	Registry    *decode.Registry
	DumpOptions decode.DumpOptions
	OS          osenv.OS
}

type Variable struct {
	Name  string
	Value interface{}
}

type Function struct {
	Name     string
	MinArity int
	MaxArity int
	Fn       func(interface{}, []interface{}) interface{}
}

type Query struct {
	opts         QueryOptions
	allFormats   []*decode.Format
	probeFormats []*decode.Format
	dotValue     interface{}
	variables    []Variable
	functions    []Function
	last         interface{}
	outCount     int
}

func NewQuery(opts QueryOptions) *Query {
	q := &Query{opts: opts}

	// TODO: cleanup group names and panics
	q.allFormats = opts.Registry.MustAll()
	q.probeFormats = opts.Registry.MustGroup(format.PROBE)
	q.functions = []Function{
		{"help", 0, 0, q.help},
		{"bits", 0, 2, q.bits},
		{"string", 0, 0, q.string_},
		{"probe", 0, 1, q.makeProbeFn(q.probeFormats)},
		{"hexdump", 0, 0, q.hexdump},
		{"dump", 0, 1, q.dump},
		{"open", 0, 1, q.open},
		{"u", 1, 1, q.u},
		{"dot", 0, 0, q.dot},
	}
	for _, f := range q.allFormats {
		q.functions = append(q.functions, Function{f.Name, 0, 0, q.makeProbeFn([]*decode.Format{f})})
	}
	q.variables = []Variable{
		{Name: "FORMAT", Value: opts.FormatName},
		{Name: "FILENAME", Value: opts.Filename},
	}

	return q
}

type queryBB struct {
	bb       *bitio.Buffer
	filename string
}

type queryDump struct {
	maxDepth int
	v        *decode.Value
}

type queryHexDump struct {
	bb *bitio.Buffer
}

type queryHelp struct{}

func (q *Query) help(c interface{}, a []interface{}) interface{} {
	return &queryHelp{}
}

func (q *Query) makeProbeFn(formats []*decode.Format) func(c interface{}, a []interface{}) interface{} {
	return func(c interface{}, a []interface{}) interface{} {
		bb, filename, err := toBB(c)
		if err != nil {
			return err
		}

		opts := map[string]interface{}{}

		name := "unnamed"
		if filename != "" {
			name = filename
		}

		dv, _, errs := decode.Probe(name, bb, formats, decode.ProbeOptions{FormatOptions: opts})
		if dv == nil {
			return errs
		}

		return dv
	}
}

func (q *Query) bits(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBB(c)
	if err != nil {
		return err
	}

	startArg := int64(0)
	endArg := int64(-1)
	toAbs := func(v int64, l int64) int64 {
		if v < 0 {
			return l + v + 1
		}
		return v
	}

	if len(a) >= 1 {
		startArg, err = toInt64(a[0])
		if err != nil {
			return err
		}
	}
	if len(a) >= 2 {
		endArg, err = toInt64(a[1])
		if err != nil {
			return err
		}
	}

	startArg = toAbs(startArg, bb.Len())
	endArg = toAbs(endArg, bb.Len())

	bb, err = bb.BitBufRange(startArg, endArg-startArg)
	if err != nil {
		return err
	}

	return bb
}

func (q *Query) string_(c interface{}, a []interface{}) interface{} {
	var bb *bitio.Buffer
	switch cc := c.(type) {
	case *decode.Value:
		var err error
		bb, err = cc.RootBitBuf.BitBufRange(cc.Range.Start, cc.Range.Len)
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

func (q *Query) hexdump(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBB(c)
	if err != nil {
		return err
	}

	return &queryHexDump{
		bb: bb,
	}
}

func (q *Query) dump(c interface{}, a []interface{}) interface{} {
	var v *decode.Value
	switch cc := c.(type) {
	case *decode.Value:
		v = cc
	case *decode.D:
		// TODO: remove?
		v = cc.Value
	default:
		return fmt.Errorf("%v: value is not a decode value", c)
	}

	maxDepth := 0
	if len(a) == 1 {
		var ok bool
		maxDepth, ok = a[0].(int)
		if !ok {
			return fmt.Errorf("max depth is not a int")
		}
		if maxDepth < 0 {
			return fmt.Errorf("max depth can't be negative")
		}
	}

	return &queryDump{
		maxDepth: maxDepth,
		v:        v,
	}
}

func (q *Query) open(c interface{}, a []interface{}) interface{} {
	var rs io.ReadSeeker

	var filename string
	if len(a) == 1 {
		var filenameOk bool
		filename, filenameOk = a[0].(string)
		if !filenameOk {
			return fmt.Errorf("filename must be a string")
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
		// TODO: query Close method that cleanups?
		// if c, ok := f.(io.Closer); ok {
		// 	defer c.Close()
		// }
		rs = f
	}

	bb, err := bitio.NewBufferFromReadSeeker(rs)
	if err != nil {
		return err
	}

	return &queryBB{
		bb:       bb,
		filename: filename,
	}
}

func (q *Query) u(c interface{}, a []interface{}) interface{} {
	bb, _, err := toBB(c)
	if err != nil {
		return err
	}

	nBits, err := toInt64(a[0])
	if err != nil {
		return err
	}
	n, err := bb.U(int(nBits))
	if err != nil {
		return err
	}

	return new(big.Int).SetUint64(n)
}

func (q *Query) dot(c interface{}, a []interface{}) interface{} {
	q.dotValue = c
	return c
}

func (q *Query) Run(src string) ([]interface{}, error) {
	var err error

	query, err := gojq.Parse(src)
	if err != nil {
		return nil, err
	}

	var variableNames []string
	var variableValues []interface{}
	variableNames = append(variableNames, "$last")
	variableValues = append(variableValues, q.last)
	for _, v := range q.variables {
		variableNames = append(variableNames, "$"+v.Name)
		variableValues = append(variableValues, v.Value)
	}

	var compilerOpts []gojq.CompilerOption
	for _, f := range q.functions {
		compilerOpts = append(compilerOpts,
			gojq.WithFunction(f.Name, f.MinArity, f.MaxArity, f.Fn))
	}
	compilerOpts = append(compilerOpts, gojq.WithVariables(variableNames))
	code, err := gojq.Compile(query, compilerOpts...)
	if err != nil {
		return nil, err
	}
	iter := code.Run(q.dotValue, variableValues...)

	var vs []interface{}
	for {
		var ok bool
		var v interface{}

		if v, ok = iter.Next(); !ok {
			break
		}
		if err, ok = v.(error); ok {
			break
		}

		switch vv := v.(type) {
		case *queryHelp:
			for _, f := range q.functions {
				for i := f.MinArity; i <= f.MaxArity; i++ {
					fmt.Fprintf(q.opts.OS.Stdout(), "%s/%d", f.Name, i)
					if i != f.MaxArity {
						fmt.Fprintf(q.opts.OS.Stdout(), ", ")
					}
				}
				fmt.Fprintf(q.opts.OS.Stdout(), "\n")
			}
		case *queryDump:
			opts := q.opts.DumpOptions
			opts.MaxDepth = vv.maxDepth
			if err := vv.v.Dump(q.opts.OS.Stdout(), opts); err != nil {
				return nil, err
			}
		case *decode.Value:
			opts := q.opts.DumpOptions
			opts.MaxDepth = 1
			if err := vv.Dump(q.opts.OS.Stdout(), opts); err != nil {
				return nil, err
			}
		case *decode.D:
			// TODO: remove?
			if err := vv.Value.Dump(q.opts.OS.Stdout(), q.opts.DumpOptions); err != nil {
				return nil, err
			}
		case *queryHexDump:
			hw := hexdump.New(
				q.opts.OS.Stdout(),
				num.DigitsInBase(bitio.BitsByteCount(vv.bb.Len()), 16),
				q.opts.DumpOptions.LineBytes)
			defer hw.Close()
			if _, err := io.Copy(hw, vv.bb); err != nil {
				return nil, err
			}
		case *bitio.Buffer:
			io.Copy(q.opts.OS.Stdout(), vv.Copy())
		case string, int, int32, int64, uint, uint32, uint64:
			fmt.Fprintln(q.opts.OS.Stdout(), vv)
		case float32:
			fmt.Fprintln(q.opts.OS.Stdout(), strconv.FormatFloat(float64(vv), 'f', -1, 32))
		case float64:
			fmt.Fprintln(q.opts.OS.Stdout(), strconv.FormatFloat(vv, 'f', -1, 64))
		default:
			e := json.NewEncoder(q.opts.OS.Stdout())
			e.SetIndent("", "  ")
			if err := e.Encode(v); err != nil {
				return nil, err
			}
		}

		vs = append(vs, v)
	}

	return vs, err
}

func (q *Query) REPL() error {
	scanner := bufio.NewScanner(q.opts.OS.Stdin())

	for {
		prompt := "> "
		if q.dotValue != nil {
			if v, ok := q.dotValue.(*decode.Value); ok {
				prompt = v.Path() + "> "
			}
		}

		fmt.Fprint(q.opts.OS.Stdout(), prompt)
		if !scanner.Scan() {
			return scanner.Err()
		}
		src := scanner.Text()

		vs, err := q.Run(src)
		if err != nil {
			fmt.Fprintf(q.opts.OS.Stdout(), "error: %s\n", err)
		}
		varName := fmt.Sprintf("out%d", q.outCount)
		q.variables = append(q.variables, Variable{Name: varName, Value: vs})
		q.outCount++

		if len(vs) > 0 {
			q.last = vs[0]
		}

		//fmt.Fprintf(q.opts.OS.Stdout(), "%s\n", varName)
	}
}
