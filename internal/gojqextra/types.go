//nolint:gosimple
package gojqextra

import (
	"bytes"
	"fmt"

	"github.com/wader/fq/internal/colorjson"

	"github.com/wader/gojq"
)

// TODO: preview errors

func expectedArrayOrObject(key interface{}, typ string) error {
	switch v := key.(type) {
	case string:
		return ExpectedObjectWithKeyError{Typ: typ, Key: v}
	case int:
		return ExpectedArrayWithIndexError{Typ: typ, Index: v}
	default:
		panic("unreachable")
	}
}

// array

var _ gojq.JQValue = Array{}

type Array []interface{}

func (v Array) JQValueLength() interface{}   { return len(v) }
func (v Array) JQValueSliceLen() interface{} { return len(v) }
func (v Array) JQValueIndex(index int) interface{} {
	if index < 0 {
		return nil
	}
	return v[index]
}
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
func (v Array) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	// TODO: handle {start:, end: }
	// TODO: maybe should use gojq implementation as it's quite complex
	intKey, ok := key.(int)
	if !ok {
		return HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
	}

	if intKey > 0x3ffffff {
		return ArrayIndexTooLargeError{V: intKey}
	}

	l := len(v)
	if intKey >= l {
		if delpath {
			return v
		}
		l = intKey + 1
	} else if intKey < -l {
		if delpath {
			return v
		}
		// TODO: wrong error?
		return FuncTypeNameError{Name: "setpath", Typ: "number"}
	} else if intKey < 0 {
		intKey += len(v)
	}

	var uu []interface{}
	if delpath {
		uu = append(v[0:intKey], v[intKey+1:]...)
	} else {
		uu = make([]interface{}, l)
		copy(uu, v)
		uu[intKey] = u
	}

	return uu
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
	return FuncTypeNameError{Name: "tonumber", Typ: "array"}
}
func (v Array) JQValueToString() interface{} {
	return FuncTypeNameError{Name: "tostring", Typ: "array"}
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
func (v Object) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	stringKey, ok := key.(string)
	if !ok {
		return HasKeyTypeError{L: "object", R: fmt.Sprintf("%v", key)}
	}

	uu := make(map[string]interface{}, len(v))
	for kv, vv := range v {
		uu[kv] = vv
	}
	if delpath {
		delete(uu, stringKey)
	} else {
		uu[stringKey] = u
	}

	return uu
}
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
	return FuncTypeNameError{Name: "tonumber", Typ: "object"}
}
func (v Object) JQValueToString() interface{} {
	return FuncTypeNameError{Name: "tostring", Typ: "object"}
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
func (v Number) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return expectedArrayOrObject(key, "number")
}
func (v Number) JQValueEach() interface{} { return IteratorError{Typ: "number"} }
func (v Number) JQValueKeys() interface{} { return FuncTypeNameError{Name: "keys", Typ: "number"} }
func (v Number) JQValueHas(key interface{}) interface{} {
	return FuncTypeNameError{Name: "has", Typ: "number"}
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

type String []rune

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
func (v String) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return expectedArrayOrObject(key, "string")
}
func (v String) JQValueEach() interface{} { return IteratorError{Typ: "string"} }
func (v String) JQValueKeys() interface{} { return FuncTypeNameError{Name: "keys", Typ: "string"} }
func (v String) JQValueHas(key interface{}) interface{} {
	return FuncTypeNameError{Name: "has", Typ: "string"}
}
func (v String) JQValueType() string          { return "string" }
func (v String) JQValueToNumber() interface{} { return gojq.NormalizeNumbers(string(v)) }
func (v String) JQValueToString() interface{} { return string(v) }
func (v String) JQValueToGoJQ() interface{}   { return string(v) }

// boolean

var _ gojq.JQValue = Boolean(true)

type Boolean bool

func (v Boolean) JQValueLength() interface{} {
	return FuncTypeNameError{Name: "length", Typ: "boolean"}
}
func (v Boolean) JQValueSliceLen() interface{}       { return ExpectedArrayError{Typ: "boolean"} }
func (v Boolean) JQValueIndex(index int) interface{} { return ExpectedArrayError{Typ: "boolean"} }
func (v Boolean) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "boolean"}
}
func (v Boolean) JQValueKey(name string) interface{} { return ExpectedObjectError{Typ: "boolean"} }
func (v Boolean) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return expectedArrayOrObject(key, "boolean")
}
func (v Boolean) JQValueEach() interface{} { return IteratorError{Typ: "boolean"} }
func (v Boolean) JQValueKeys() interface{} { return FuncTypeNameError{Name: "keys", Typ: "boolean"} }
func (v Boolean) JQValueHas(key interface{}) interface{} {
	return FuncTypeNameError{Name: "has", Typ: "boolean"}
}
func (v Boolean) JQValueType() string { return "boolean" }
func (v Boolean) JQValueToNumber() interface{} {
	return FuncTypeNameError{Name: "tonumber", Typ: "boolean"}
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
func (v Null) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return expectedArrayOrObject(key, "null")
}
func (v Null) JQValueEach() interface{} { return IteratorError{Typ: "null"} }
func (v Null) JQValueKeys() interface{} { return FuncTypeNameError{Name: "keys", Typ: "null"} }
func (v Null) JQValueHas(key interface{}) interface{} {
	return FuncTypeNameError{Name: "has", Typ: "null"}
}
func (v Null) JQValueType() string          { return "null" }
func (v Null) JQValueToNumber() interface{} { return FuncTypeNameError{Name: "tonumber", Typ: "null"} }
func (v Null) JQValueToString() interface{} { return "null" }
func (v Null) JQValueToGoJQ() interface{}   { return nil }

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
func (v Base) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return expectedArrayOrObject(key, v.Typ)
}
func (v Base) JQValueEach() interface{} { return IteratorError{Typ: v.Typ} }
func (v Base) JQValueKeys() interface{} { return FuncTypeNameError{Name: "keys", Typ: v.Typ} }
func (v Base) JQValueHas(key interface{}) interface{} {
	return HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
}
func (v Base) JQValueType() string          { return v.Typ }
func (v Base) JQValueToNumber() interface{} { return FuncTypeNameError{Name: "tonumber", Typ: v.Typ} }
func (v Base) JQValueToString() interface{} { return FuncTypeNameError{Name: "tostring", Typ: v.Typ} }
func (v Base) JQValueToGoJQ() interface{}   { return nil }

// lazy

var _ gojq.JQValue = &Lazy{}

type Lazy struct {
	Fn     func() (gojq.JQValue, error)
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

func (v *Lazy) f(fn func(jv gojq.JQValue) interface{}) interface{} {
	jv, err := v.v()
	if err != nil {
		return err
	}
	return fn(jv)
}

func (v *Lazy) JQValueLength() interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueLength() })
}
func (v *Lazy) JQValueSliceLen() interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueSliceLen() })
}
func (v *Lazy) JQValueIndex(index int) interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueIndex(index) })
}
func (v *Lazy) JQValueSlice(start int, end int) interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueSlice(start, end) })
}
func (v *Lazy) JQValueKey(name string) interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueKey(name) })
}
func (v *Lazy) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueUpdate(key, u, delpath) })
}
func (v *Lazy) JQValueEach() interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueEach() })
}
func (v *Lazy) JQValueKeys() interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueKeys() })
}
func (v *Lazy) JQValueHas(key interface{}) interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueHas(key) })
}
func (v *Lazy) JQValueType() string {
	jv, err := v.v()
	if err != nil {
		return "error"
	}
	return jv.JQValueType()
}
func (v *Lazy) JQValueToNumber() interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueToNumber() })
}
func (v *Lazy) JQValueToString() interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueToString() })
}
func (v *Lazy) JQValueToGoJQ() interface{} {
	return v.f(func(jv gojq.JQValue) interface{} { return jv.JQValueToGoJQ() })
}
