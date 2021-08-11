//nolint:gosimple
package gojqextra

import (
	"bytes"
	"fmt"
	"fq/internal/colorjson"

	"github.com/itchyny/gojq"
)

// TODO: refactor to use errors from gojq?
// TODO: preview errors

type FuncTypeError struct {
	Name string
	Typ  string
}

func (err FuncTypeError) Error() string { return err.Name + " cannot be applied to: " + err.Typ }

type ExpectedObjectError struct {
	Typ string
}

func (err ExpectedObjectError) Error() string {
	return "expected an object but got: " + err.Typ
}

type ExpectedArrayError struct {
	Typ string
}

func (err ExpectedArrayError) Error() string {
	return "expected an array but got: " + err.Typ
}

type ExpectedObjectWithKeyError struct {
	Typ string
	Key string
}

func (err ExpectedObjectWithKeyError) Error() string {
	return fmt.Sprintf("expected an object with key %q but got: %s", err.Key, err.Typ)
}

type ExpectedArrayWithIndexError struct {
	Typ   string
	Index int
}

func (err ExpectedArrayWithIndexError) Error() string {
	return fmt.Sprintf("expected an array with index %d but got: %s", err.Index, err.Typ)
}

type IteratorError struct {
	Typ string
}

func (err IteratorError) Error() string {
	return "cannot iterate over: " + err.Typ
}

type HasKeyTypeError struct {
	L, R string
}

func (err HasKeyTypeError) Error() string {
	return "cannot check whether " + err.L + " has a key: " + err.R
}

// array

var _ gojq.JQValue = Array{}

type Array []interface{}

func (v Array) JQValueLength() interface{}                  { return len(v) }
func (v Array) JQValueSliceLen() interface{}                { return len(v) }
func (v Array) JQValueIndex(index int) interface{}          { return v[index] }
func (v Array) JQValueSlice(start int, end int) interface{} { return v[start:end] }
func (v Array) JQValueKey(name string) interface{} {
	return ExpectedObjectError{Typ: "array"}
}
func (v Array) JQValueEach() interface{} {
	vs := make([]gojq.PathValue, len(v))
	for i, v := range v {
		vs[i] = gojq.PathValue{Path: i, Value: v}
	}
	return vs
}
func (v Array) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v))
	for i := range v {
		vs[i] = i
	}
	return vs
}
func (v Array) JQValueHas(key interface{}) interface{} {
	intKey, ok := key.(int)
	if !ok {
		return HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
	}
	return intKey >= 0 && intKey < len(v)
}
func (v Array) JQValueType() string { return "array" }
func (v Array) JQValueToNumber() interface{} {
	return FuncTypeError{Name: "tonumber", Typ: "array"}
}
func (v Array) JQValueToString() interface{} {
	return FuncTypeError{Name: "tostring", Typ: "array"}
}
func (v Array) JQValueToGoJQ() interface{} { return []interface{}(v) }

// object

var _ gojq.JQValue = Object{}

type Object map[string]interface{}

func (v Object) JQValueLength() interface{}         { return len(v) }
func (v Object) JQValueSliceLen() interface{}       { return ExpectedArrayError{Typ: "object"} }
func (v Object) JQValueIndex(index int) interface{} { return ExpectedArrayError{Typ: "object"} }
func (v Object) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "object"}
}
func (v Object) JQValueKey(name string) interface{} { return v[name] }
func (v Object) JQValueEach() interface{} {
	vs := make([]gojq.PathValue, len(v))
	i := 0
	for k, v := range v {
		vs[i] = gojq.PathValue{Path: k, Value: v}
		i++
	}
	return vs
}
func (v Object) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v))
	i := 0
	for k := range v {
		vs[i] = k
		i++
	}
	return vs
}
func (v Object) JQValueHas(key interface{}) interface{} {
	stringKey, ok := key.(string)
	if !ok {
		return HasKeyTypeError{L: "object", R: fmt.Sprintf("%v", key)}
	}
	_, ok = v[stringKey]
	return ok
}
func (v Object) JQValueType() string { return "object" }
func (v Object) JQValueToNumber() interface{} {
	return FuncTypeError{Name: "tonumber", Typ: "object"}
}
func (v Object) JQValueToString() interface{} {
	return FuncTypeError{Name: "tostring", Typ: "object"}
}
func (v Object) JQValueToGoJQ() interface{} { return map[string]interface{}(v) }

// number

var _ gojq.JQValue = Number{}

type Number struct {
	V interface{}
}

func (v Number) JQValueLength() interface{}         { return v.V }
func (v Number) JQValueSliceLen() interface{}       { return ExpectedArrayError{Typ: "number"} }
func (v Number) JQValueIndex(index int) interface{} { return ExpectedArrayError{Typ: "number"} }
func (v Number) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "number"}
}
func (v Number) JQValueKey(name string) interface{} { return ExpectedObjectError{Typ: "number"} }
func (v Number) JQValueEach() interface{}           { return IteratorError{Typ: "number"} }
func (v Number) JQValueKeys() interface{}           { return FuncTypeError{Name: "keys", Typ: "number"} }
func (v Number) JQValueHas(key interface{}) interface{} {
	return FuncTypeError{Name: "has", Typ: "number"}
}
func (v Number) JQValueType() string          { return "number" }
func (v Number) JQValueToNumber() interface{} { return v.V }
func (v Number) JQValueToString() interface{} {
	b := &bytes.Buffer{}
	// uses colorjson encode based on gojq encoder to support big.Int
	if err := colorjson.NewEncoder(false, false, 0, nil, colorjson.Colors{}).Marshal(v.V, b); err != nil {
		return err
	}
	return b.String()
}
func (v Number) JQValueToGoJQ() interface{} { return v.V }

// string

var _ gojq.JQValue = String("")

type String string

func (v String) JQValueLength() interface{}   { return len(v) }
func (v String) JQValueSliceLen() interface{} { return len(v) }
func (v String) JQValueIndex(index int) interface{} {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return ""
	}
	return fmt.Sprintf("%c", v[index])
}
func (v String) JQValueSlice(start int, end int) interface{} { return string(v[start:end]) }
func (v String) JQValueKey(name string) interface{}          { return ExpectedObjectError{Typ: "string"} }
func (v String) JQValueEach() interface{}                    { return IteratorError{Typ: "string"} }
func (v String) JQValueKeys() interface{}                    { return FuncTypeError{Name: "keys", Typ: "string"} }
func (v String) JQValueHas(key interface{}) interface{} {
	return FuncTypeError{Name: "has", Typ: "string"}
}
func (v String) JQValueType() string          { return "string" }
func (v String) JQValueToNumber() interface{} { return gojq.NormalizeNumbers(string(v)) }
func (v String) JQValueToString() interface{} { return string(v) }
func (v String) JQValueToGoJQ() interface{}   { return string(v) }

// boolean

var _ gojq.JQValue = Boolean(true)

type Boolean bool

func (v Boolean) JQValueLength() interface{}         { return FuncTypeError{Name: "length", Typ: "boolean"} }
func (v Boolean) JQValueSliceLen() interface{}       { return ExpectedArrayError{Typ: "boolean"} }
func (v Boolean) JQValueIndex(index int) interface{} { return ExpectedArrayError{Typ: "boolean"} }
func (v Boolean) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "boolean"}
}
func (v Boolean) JQValueKey(name string) interface{} { return ExpectedObjectError{Typ: "boolean"} }
func (v Boolean) JQValueEach() interface{}           { return IteratorError{Typ: "boolean"} }
func (v Boolean) JQValueKeys() interface{}           { return FuncTypeError{Name: "keys", Typ: "boolean"} }
func (v Boolean) JQValueHas(key interface{}) interface{} {
	return FuncTypeError{Name: "has", Typ: "boolean"}
}
func (v Boolean) JQValueType() string { return "boolean" }
func (v Boolean) JQValueToNumber() interface{} {
	return FuncTypeError{Name: "tonumber", Typ: "boolean"}
}
func (v Boolean) JQValueToString() interface{} {
	if v {
		return "true"
	}
	return "false"
}
func (v Boolean) JQValueToGoJQ() interface{} { return bool(v) }

// null

var _ gojq.JQValue = Null{}

type Null struct{}

func (v Null) JQValueLength() interface{}                  { return 0 }
func (v Null) JQValueSliceLen() interface{}                { return ExpectedArrayError{Typ: "null"} }
func (v Null) JQValueIndex(index int) interface{}          { return ExpectedArrayError{Typ: "null"} }
func (v Null) JQValueSlice(start int, end int) interface{} { return ExpectedArrayError{Typ: "null"} }
func (v Null) JQValueKey(name string) interface{}          { return ExpectedObjectError{Typ: "null"} }
func (v Null) JQValueEach() interface{}                    { return IteratorError{Typ: "null"} }
func (v Null) JQValueKeys() interface{}                    { return FuncTypeError{Name: "keys", Typ: "null"} }
func (v Null) JQValueHas(key interface{}) interface{}      { return FuncTypeError{Name: "has", Typ: "null"} }
func (v Null) JQValueType() string                         { return "null" }
func (v Null) JQValueToNumber() interface{}                { return FuncTypeError{Name: "tonumber", Typ: "null"} }
func (v Null) JQValueToString() interface{}                { return "null" }
func (v Null) JQValueToGoJQ() interface{}                  { return nil }

// Base

var _ gojq.JQValue = Base{}

type Base struct {
	Typ string
}

func (v Base) JQValueLength() interface{}   { return ExpectedArrayError{Typ: v.Typ} }
func (v Base) JQValueSliceLen() interface{} { return ExpectedArrayError{Typ: v.Typ} }
func (v Base) JQValueIndex(index int) interface{} {
	return ExpectedArrayWithIndexError{Typ: v.Typ, Index: index}
}
func (v Base) JQValueSlice(start int, end int) interface{} { return ExpectedArrayError{Typ: v.Typ} }
func (v Base) JQValueKey(name string) interface{} {
	return ExpectedObjectWithKeyError{Typ: v.Typ, Key: name}
}
func (v Base) JQValueEach() interface{} { return IteratorError{Typ: v.Typ} }
func (v Base) JQValueKeys() interface{} { return FuncTypeError{Name: "keys", Typ: v.Typ} }
func (v Base) JQValueHas(key interface{}) interface{} {
	return HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
}
func (v Base) JQValueType() string          { return v.Typ }
func (v Base) JQValueToNumber() interface{} { return FuncTypeError{Name: "tonumber", Typ: v.Typ} }
func (v Base) JQValueToString() interface{} { return FuncTypeError{Name: "tostring", Typ: v.Typ} }
func (v Base) JQValueToGoJQ() interface{}   { return nil }
