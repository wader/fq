//nolint:gosimple
package gojqextra

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/wader/fq/internal/colorjson"

	"github.com/wader/gojq"
)

func Typeof(v any) string {
	switch v := v.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case int, float64, *big.Int:
		return "number"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case gojq.JQValue:
		return v.JQValueType()
	default:
		panic(fmt.Sprintf("invalid value: %v", v))
	}
}

// TODO: preview errors

// array

var _ gojq.JQValue = Array{}

type Array []any

func (v Array) JQValueLength() any   { return len(v) }
func (v Array) JQValueSliceLen() any { return len(v) }
func (v Array) JQValueIndex(index int) any {
	if index < 0 {
		return nil
	}
	return v[index]
}
func (v Array) JQValueSlice(start int, end int) any { return v[start:end] }
func (v Array) JQValueKey(name string) any {
	return ExpectedObjectError{Typ: "array"}
}
func (v Array) JQValueEach() any {
	vs := make([]gojq.PathValue, len(v))
	for i, v := range v {
		vs[i] = gojq.PathValue{Path: i, Value: v}
	}
	return vs
}
func (v Array) JQValueKeys() any {
	vs := make([]any, len(v))
	for i := range v {
		vs[i] = i
	}
	return vs
}
func (v Array) JQValueHas(key any) any {
	intKey, ok := key.(int)
	if !ok {
		return HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
	}
	return intKey >= 0 && intKey < len(v)
}
func (v Array) JQValueType() string { return "array" }
func (v Array) JQValueToNumber() any {
	return FuncTypeNameError{Name: "tonumber", Typ: "array"}
}
func (v Array) JQValueToString() any {
	return FuncTypeNameError{Name: "tostring", Typ: "array"}
}
func (v Array) JQValueToGoJQ() any { return []any(v) }

// object

var _ gojq.JQValue = Object{}

type Object map[string]any

func (v Object) JQValueLength() any         { return len(v) }
func (v Object) JQValueSliceLen() any       { return ExpectedArrayError{Typ: "object"} }
func (v Object) JQValueIndex(index int) any { return ExpectedArrayError{Typ: "object"} }
func (v Object) JQValueSlice(start int, end int) any {
	return ExpectedArrayError{Typ: "object"}
}
func (v Object) JQValueKey(name string) any { return v[name] }
func (v Object) JQValueEach() any {
	vs := make([]gojq.PathValue, len(v))
	i := 0
	for k, v := range v {
		vs[i] = gojq.PathValue{Path: k, Value: v}
		i++
	}
	return vs
}
func (v Object) JQValueKeys() any {
	vs := make([]any, len(v))
	i := 0
	for k := range v {
		vs[i] = k
		i++
	}
	return vs
}
func (v Object) JQValueHas(key any) any {
	stringKey, ok := key.(string)
	if !ok {
		return HasKeyTypeError{L: "object", R: fmt.Sprintf("%v", key)}
	}
	_, ok = v[stringKey]
	return ok
}
func (v Object) JQValueType() string { return "object" }
func (v Object) JQValueToNumber() any {
	return FuncTypeNameError{Name: "tonumber", Typ: "object"}
}
func (v Object) JQValueToString() any {
	return FuncTypeNameError{Name: "tostring", Typ: "object"}
}
func (v Object) JQValueToGoJQ() any { return map[string]any(v) }

// number

var _ gojq.JQValue = Number{}

type Number struct {
	V any
}

func (v Number) JQValueLength() any         { return v.V }
func (v Number) JQValueSliceLen() any       { return ExpectedArrayError{Typ: "number"} }
func (v Number) JQValueIndex(index int) any { return ExpectedArrayError{Typ: "number"} }
func (v Number) JQValueSlice(start int, end int) any {
	return ExpectedArrayError{Typ: "number"}
}
func (v Number) JQValueKey(name string) any { return ExpectedObjectError{Typ: "number"} }
func (v Number) JQValueEach() any           { return IteratorError{Typ: "number"} }
func (v Number) JQValueKeys() any           { return FuncTypeNameError{Name: "keys", Typ: "number"} }
func (v Number) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: "number"}
}
func (v Number) JQValueType() string  { return "number" }
func (v Number) JQValueToNumber() any { return v.V }
func (v Number) JQValueToString() any {
	b := &bytes.Buffer{}
	// uses colorjson encode based on gojq encoder to support big.Int
	if err := colorjson.NewEncoder(false, false, 0, nil, colorjson.Colors{}).Marshal(v.V, b); err != nil {
		return err
	}
	return b.String()
}
func (v Number) JQValueToGoJQ() any { return v.V }

// string

var _ gojq.JQValue = String("")

type String []rune

func (v String) JQValueLength() any   { return len(v) }
func (v String) JQValueSliceLen() any { return len(v) }
func (v String) JQValueIndex(index int) any {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return ""
	}
	return fmt.Sprintf("%c", v[index])
}
func (v String) JQValueSlice(start int, end int) any { return string(v[start:end]) }
func (v String) JQValueKey(name string) any          { return ExpectedObjectError{Typ: "string"} }
func (v String) JQValueEach() any                    { return IteratorError{Typ: "string"} }
func (v String) JQValueKeys() any                    { return FuncTypeNameError{Name: "keys", Typ: "string"} }
func (v String) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: "string"}
}
func (v String) JQValueType() string  { return "string" }
func (v String) JQValueToNumber() any { return gojq.NormalizeNumbers(string(v)) }
func (v String) JQValueToString() any { return string(v) }
func (v String) JQValueToGoJQ() any   { return string(v) }

// boolean

var _ gojq.JQValue = Boolean(true)

type Boolean bool

func (v Boolean) JQValueLength() any {
	return FuncTypeNameError{Name: "length", Typ: "boolean"}
}
func (v Boolean) JQValueSliceLen() any       { return ExpectedArrayError{Typ: "boolean"} }
func (v Boolean) JQValueIndex(index int) any { return ExpectedArrayError{Typ: "boolean"} }
func (v Boolean) JQValueSlice(start int, end int) any {
	return ExpectedArrayError{Typ: "boolean"}
}
func (v Boolean) JQValueKey(name string) any { return ExpectedObjectError{Typ: "boolean"} }
func (v Boolean) JQValueEach() any           { return IteratorError{Typ: "boolean"} }
func (v Boolean) JQValueKeys() any           { return FuncTypeNameError{Name: "keys", Typ: "boolean"} }
func (v Boolean) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: "boolean"}
}
func (v Boolean) JQValueType() string { return "boolean" }
func (v Boolean) JQValueToNumber() any {
	return FuncTypeNameError{Name: "tonumber", Typ: "boolean"}
}
func (v Boolean) JQValueToString() any {
	if v {
		return "true"
	}
	return "false"
}
func (v Boolean) JQValueToGoJQ() any { return bool(v) }

// null

var _ gojq.JQValue = Null{}

type Null struct{}

func (v Null) JQValueLength() any                  { return 0 }
func (v Null) JQValueSliceLen() any                { return ExpectedArrayError{Typ: "null"} }
func (v Null) JQValueIndex(index int) any          { return ExpectedArrayError{Typ: "null"} }
func (v Null) JQValueSlice(start int, end int) any { return ExpectedArrayError{Typ: "null"} }
func (v Null) JQValueKey(name string) any          { return ExpectedObjectError{Typ: "null"} }

func (v Null) JQValueEach() any { return IteratorError{Typ: "null"} }
func (v Null) JQValueKeys() any { return FuncTypeNameError{Name: "keys", Typ: "null"} }
func (v Null) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: "null"}
}
func (v Null) JQValueType() string  { return "null" }
func (v Null) JQValueToNumber() any { return FuncTypeNameError{Name: "tonumber", Typ: "null"} }
func (v Null) JQValueToString() any { return "null" }
func (v Null) JQValueToGoJQ() any   { return nil }

// Base

var _ gojq.JQValue = Base{}

type Base struct {
	Typ string
}

func (v Base) JQValueLength() any   { return ExpectedArrayError{Typ: v.Typ} }
func (v Base) JQValueSliceLen() any { return ExpectedArrayError{Typ: v.Typ} }
func (v Base) JQValueIndex(index int) any {
	return ExpectedArrayWithIndexError{Typ: v.Typ, Index: index}
}
func (v Base) JQValueSlice(start int, end int) any { return ExpectedArrayError{Typ: v.Typ} }
func (v Base) JQValueKey(name string) any {
	return ExpectedObjectWithKeyError{Typ: v.Typ, Key: name}
}
func (v Base) JQValueEach() any { return IteratorError{Typ: v.Typ} }
func (v Base) JQValueKeys() any { return FuncTypeNameError{Name: "keys", Typ: v.Typ} }
func (v Base) JQValueHas(key any) any {
	return HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
}
func (v Base) JQValueType() string  { return v.Typ }
func (v Base) JQValueToNumber() any { return FuncTypeNameError{Name: "tonumber", Typ: v.Typ} }
func (v Base) JQValueToString() any { return FuncTypeNameError{Name: "tostring", Typ: v.Typ} }
func (v Base) JQValueToGoJQ() any   { return nil }

// lazy

var _ gojq.JQValue = &Lazy{}

type Lazy struct {
	Type     string
	IsScalar bool
	Fn       func() (gojq.JQValue, error)

	called bool
	err    error
	jv     gojq.JQValue
}

func (v *Lazy) v() (gojq.JQValue, error) {
	if !v.called {
		v.jv, v.err = v.Fn()
		v.called = true
	}
	return v.jv, v.err
}

func (v *Lazy) f(fn func(jv gojq.JQValue) any) any {
	jv, err := v.v()
	if err != nil {
		return err
	}
	return fn(jv)
}

func (v *Lazy) JQValueLength() any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueLength() })
}
func (v *Lazy) JQValueSliceLen() any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueSliceLen() })
}
func (v *Lazy) JQValueIndex(index int) any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueIndex(index) })
}
func (v *Lazy) JQValueSlice(start int, end int) any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueSlice(start, end) })
}
func (v *Lazy) JQValueKey(name string) any {
	if v.IsScalar {
		return ExpectedObjectWithKeyError{Typ: v.Type, Key: name}
	}
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueKey(name) })
}
func (v *Lazy) JQValueEach() any {
	if v.IsScalar {
		return IteratorError{Typ: v.Type}
	}
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueEach() })
}
func (v *Lazy) JQValueKeys() any {
	if v.IsScalar {
		return FuncTypeNameError{Name: "keys", Typ: "string"}
	}
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueKeys() })
}
func (v *Lazy) JQValueHas(key any) any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueHas(key) })
}
func (v *Lazy) JQValueType() string { return v.Type }
func (v *Lazy) JQValueToNumber() any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueToNumber() })
}
func (v *Lazy) JQValueToString() any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueToString() })
}
func (v *Lazy) JQValueToGoJQ() any {
	return v.f(func(jv gojq.JQValue) any { return jv.JQValueToGoJQ() })
}
