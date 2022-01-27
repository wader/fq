package ksexpr

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/wader/fq/format/kaitai/ksexpr/strconvx"
)

// TODO: Value interface?

func ToInt(v any) (int, bool) {
	switch v := v.(type) {
	case Integer:
		return int(v), true
	case Float:
		return v.Int(), true
	case BigInt:
		return v.Int(), true
	default:
		return 0, false
	}
}

func ToInt64(v any) (int64, bool) {
	switch v := v.(type) {
	case Integer:
		return int64(v), true
	case Float:
		return v.Int64(), true
	case BigInt:
		return v.Int64(), true
	default:
		return 0, false
	}
}

func ToValue(v any) any {
	switch v := v.(type) {
	case bool:
		return Boolean(v)
	case byte:
		return Integer(v)
	case int:
		return Integer(v)
	case int16:
		return Integer(v)
	case int32:
		return Integer(v)
	case int64:
		if v < 0 || v > math.MaxInt {
			return BigInt{new(big.Int).SetInt64(v)}
		}
		return Integer(v)
	case uint16:
		return Integer(v)
	case uint32:
		return Integer(v)
	case uint64:
		if v > math.MaxInt {
			return BigInt{new(big.Int).SetUint64(v)}
		}
		return Integer(v)
	case *big.Int:
		return BigInt{v}
	case float32:
		return Float(v)
	case float64:
		return Float(v)
	case string:
		return String(v)
	case []any:
		// TODO: ToValue elements?
		return Array(v)
	case map[string]any:
		// TODO: ToValue elements?
		return Object(v)
	}

	return v
}

func floatToNumber(v float64) any {
	// TODO: ok way to do lossless int case?
	i := int(v)
	if float64(i) == v {
		return Integer(i)
	}
	return Float(v)
}

func bigIntIsZero(b *big.Int) bool {
	return b.Sign() == 0
}

// TODO: bigint
// uses fork of strconv.Parse* to support trailing invalid characters and
// optional base prefix even when base is provided.
func strToInteger(s string, base int, strict bool) (any, int, error) {
	if i64, l, err := strconvx.ParseInt(s, base, 64, strict); err == nil {
		if i64 >= math.MinInt && i64 <= math.MaxInt {
			return int(i64), l, nil
		}
		return i64, l, nil
	} else if b, ok := new(big.Int).SetString(s, 0); ok {
		return b, len(s), nil
	}
	return nil, 0, &strconv.NumError{Func: "strToInteger", Num: s, Err: strconv.ErrSyntax}
}

func strReverse(s string) string {
	rs := []rune(s)
	l := len(rs)
	for i := 0; i < l/2; i++ {
		rs[i], rs[l-i-1] = rs[l-i-1], rs[i]
	}
	return string(rs)
}

type noSuchMethodError struct {
	input any
	name  string
}

func (a noSuchMethodError) Error() string {
	return fmt.Sprintf("%s.%s no such method", a.input, a.name)
}

type noSuchKeyError struct {
	input any
	name  string
}

func (a noSuchKeyError) Error() string {
	return fmt.Sprintf("%s[%s] no such key ", a.input, a.name)
}

type noSuchIndexError struct {
	input any
	index int
	max   int
}

func (a noSuchIndexError) Error() string {
	return fmt.Sprintf("%s[%d] no such index (0-%d)", a.input, a.index, a.max)
}

type invalidIndexError struct {
	input any
	index any
}

func (a invalidIndexError) Error() string {
	return fmt.Sprintf("%s[%s] invalid index", a.input, a.index)
}

type notIndexableError struct {
	input any
	index int
}

func (a notIndexableError) Error() string {
	return fmt.Sprintf("%s[%d] not indexable", a.input, a.index)
}

type noArgsError struct {
	input any
	name  string
}

func (a noArgsError) Error() string {
	return fmt.Sprintf("%s.%s takes no arguments", a.input, a.name)
}

type argsError struct {
	input any
	name  string
	args  string
}

func (a argsError) Error() string {
	return fmt.Sprintf("%s.%s %s", a.input, a.name, a.args)
}

type Boolean bool

func (b Boolean) String() string { return fmt.Sprintf("%t (boolean)", bool(b)) }
func (b Boolean) KSExprCall(ns []string, name string, args []any) (any, error) {
	switch name {
	case "to_i":
		if len(args) > 0 {
			return nil, noArgsError{b, name}
		}
		if bool(b) {
			return Integer(1), nil
		}
		return Integer(0), nil
	}
	return nil, noSuchMethodError{b, name}
}

type Integer int

func (i Integer) String() string { return fmt.Sprintf("%d (integer)", i) }
func (i Integer) KSExprCall(ns []string, name string, args []any) (any, error) {
	switch name {
	case "to_s":
		if len(args) > 0 {
			return nil, noArgsError{i, name}
		}
		return String(fmt.Sprintf("%d", i)), nil
	}
	return nil, noSuchMethodError{i, name}
}

type BigInt struct{ V *big.Int }

func NewBigIntFromInteger(i Integer) BigInt {
	return BigInt{new(big.Int).SetInt64(int64(i))}
}

func (i BigInt) String() string { return fmt.Sprintf("%s (integer)", i.V) }
func (i BigInt) KSExprCall(ns []string, name string, args []any) (any, error) {
	switch name {
	case "to_s":
		if len(args) > 0 {
			return nil, noArgsError{i, name}
		}
		return String(i.String()), nil
	}
	return nil, noSuchMethodError{i, name}
}
func (i BigInt) Int() int {
	if i.V.IsInt64() {
		n := i.V.Int64()
		if n > math.MaxInt {
			return math.MaxInt
		} else if n < math.MinInt {
			return math.MinInt
		}
		return int(n)
	}
	if i.V.Sign() > 0 {
		return math.MaxInt
	}
	return math.MinInt
}
func (i BigInt) Int64() int64 {
	if i.V.IsInt64() {
		return i.V.Int64()
	}
	if i.V.Sign() > 0 {
		return math.MaxInt64
	}
	return math.MinInt64
}
func (i BigInt) Float() float64 {
	if i.V.IsInt64() {
		return float64(i.V.Int64())
	}
	if f, err := strconv.ParseFloat(i.V.String(), 64); err == nil {
		return f
	}
	return math.Inf(i.V.Sign())
}

type Float float64

func NewFloatFromBigInt(b BigInt) Float {
	if b.V.IsInt64() {
		return Float(b.V.Int64())
	}
	// TODO: better?
	return Float(b.V.Int64())
}

func (f Float) String() string { return fmt.Sprintf("%f (float)", f) }
func (f Float) KSExprCall(ns []string, name string, args []any) (any, error) {
	switch name {
	case "to_i":
		if len(args) > 0 {
			return nil, noArgsError{f, name}
		}
		// TODO: better?
		return Integer(f.Int()), nil
	}
	return nil, noSuchMethodError{f, name}
}
func (f Float) Int() int {
	if f > math.MaxInt {
		return math.MaxInt
	} else if f < math.MinInt {
		return math.MinInt
	}
	return int(f)
}
func (f Float) Int64() int64 {
	if f > math.MaxInt64 {
		return math.MaxInt64
	} else if f < math.MinInt64 {
		return math.MinInt64
	}
	return int64(f)
}

type String string

func (s String) String() string { return fmt.Sprintf("%q (string)", string(s)) }
func (s String) KSExprCall(ns []string, name string, args []any) (any, error) {
	switch name {
	case "substring":
		if len(args) != 2 {
			return nil, argsError{s, name, "takes 2 arguments start and stop index"}
		}
		start, ok := ToInt(args[0])
		if !ok {
			return nil, argsError{s, name, "start has to be an integer"}
		}
		stop, ok := ToInt(args[1])
		if !ok {
			return nil, argsError{s, name, "stop has to be an integer"}
		}
		if start < 0 || start > len(s) {
			return nil, argsError{s, name, fmt.Sprintf("start %d out of range %d-%d", start, 0, len(s))}
		}
		if stop < 0 || stop > len(s) {
			return nil, argsError{s, name, fmt.Sprintf("substring stop %d out of range %d-%d", stop, 0, len(s))}
		}
		// seems to be how js kaitai works
		// if start > stop {
		// 	stop, start = start, stop
		// }
		if start > stop {
			start = stop
		}
		return String(s[start:stop]), nil
	case "length":
		if len(args) > 0 {
			return nil, noArgsError{s, name}
		}
		return Integer(len([]rune(s))), nil
	case "reverse":
		if len(args) > 0 {
			return nil, noArgsError{s, name}
		}
		// TODO: store String as []rune?
		return String(strReverse(string(s))), nil
	case "to_i":
		base := 0
		if len(args) == 1 {
			var ok bool
			base, ok = ToInt(args[0])
			if !ok {
				return nil, argsError{s, name, "base argument must be an integer"}
			}
		} else if len(args) > 1 {
			return nil, argsError{s, name, "takes no or one base argument"}
		}

		// TODO: bigint
		n, _, err := strToInteger(string(s), base, false)
		return n, err
	}

	return nil, noSuchMethodError{s, name}
}

type Enum struct {
	Name string
	V    any
}

func (e Enum) String() string { return fmt.Sprintf("%q=%v (enum)", e.Name, e.V) }
func (e Enum) KSExprCall(ns []string, name string, args []any) (any, error) {
	switch name {
	case "to_i":
		return e.V, nil
	}

	return nil, noSuchMethodError{e, name}
}

// TODO: byte array?
type Array []any

func (a Array) IsByteArray() bool {
	for _, e := range a {
		i, ok := ToInt(e)
		if !ok || i <= 0 || i >= 256 {
			return false
		}
	}
	return true
}

func (a Array) String() string {
	if a.IsByteArray() {
		sb := strings.Builder{}
		for _, e := range a {
			i, ok := ToInt(e)
			if !ok || i <= 0 || i >= 256 {
				panic("unreachable")
			}
			sb.WriteString(fmt.Sprintf("%.2x", i))
		}
		return sb.String()
	}

	if a == nil {
		return "[]"
	}
	// TODO:
	bs, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
func (a Array) KSExprCall(ns []string, name string, args []any) (any, error) {
	switch name {
	case "first", "last", "min", "max":
		if len(args) > 0 {
			return nil, noArgsError{a, name}
		}
		// TODO: empty array not allowed?
		if len(a) < 1 {
			return nil, fmt.Errorf("empty array")
		}

		switch name {
		case "first":
			return a[0], nil
		case "last":
			return a[len(a)-1], nil
		case "min", "max":
			cmp := InfixGT
			if name == "max" {
				cmp = InfixLT
			}
			m := a[0]
			for _, e := range a[1:] {
				v := cmp(m, e)
				if err, ok := v.(error); ok {
					return nil, fmt.Errorf("%s: %w", name, err)
				} else if b, ok := v.(Boolean); ok {
					if bool(b) {
						m = e
					}
					continue
				}
				panic("min/max cmp not a bool")
			}
			return m, nil
		default:
			panic("unreachable")
		}

	case "length":
		// TODO: byte elements
		if len(args) > 0 {
			return nil, noArgsError{a, name}
		}
		return Integer(len(a)), nil
	case "size":
		if len(args) > 0 {
			return nil, noArgsError{a, name}
		}
		return Integer(len(a)), nil
	}
	return nil, noSuchMethodError{a, name}
}
func (a Array) KSExprIndex(index int) (any, error) {
	if index < 0 || index >= len(a) {
		return nil, noSuchIndexError{a, index, len(a) - 1}
	}
	return a[index], nil
}

type Object map[string]any

func (o Object) String() string {
	if o == nil {
		return "{}"
	}

	// TODO:
	bs, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return string(bs)
}
func (o Object) KSExprCall(ns []string, name string, args []any) (any, error) {
	v, ok := o[name]
	if !ok {
		return nil, noSuchKeyError{o, name}
	}
	return v, nil
}
