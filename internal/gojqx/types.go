//nolint:gosimple
package gojqx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"

	"github.com/wader/fq/internal/colorjson"

	"github.com/wader/gojq"
)

// Cast gojq value to go value
//
//nolint:forcetypeassert,unconvert
func CastFn[T any](v any, structFn func(input any, result any) error) (T, bool) {
	var t T
	switch any(t).(type) {
	case bool:
		switch v := v.(type) {
		case bool:
			return any(v).(T), true
		case gojq.JQValue:
			return CastFn[T](v.JQValueToGoJQ(), structFn)
		default:
			return t, false
		}
	case int:
		switch v := v.(type) {
		case int:
			return any(v).(T), true
		case *big.Int:
			if !v.IsInt64() {
				return t, false
			}
			vi := v.Int64()
			if math.MinInt <= vi && vi <= math.MaxInt {
				return any(int(vi)).(T), true
			}
			return t, false
		case float64:
			if math.MinInt <= v && v <= math.MaxInt {
				return any(int(v)).(T), true
			}
			if v > 0 {
				return any(int(math.MaxInt)).(T), true
			}
			return any(int(math.MinInt)).(T), true
		case gojq.JQValue:
			return CastFn[T](v.JQValueToGoJQ(), structFn)
		default:
			return t, false
		}
	case float64:
		switch v := v.(type) {
		case float64:
			return any(v).(T), true
		case int:
			return any(float64(v)).(T), true
		case *big.Int:
			if v.IsInt64() {
				return any(float64(v.Int64())).(T), true
			}
			// TODO: use *big.Float SetInt
			if f, err := strconv.ParseFloat(v.String(), 64); err == nil {
				return any(f).(T), true
			}
			return any(float64(math.Inf(v.Sign()))).(T), true
		case gojq.JQValue:
			return CastFn[T](v.JQValueToGoJQ(), structFn)
		default:
			return t, false
		}
	case *big.Int:
		switch v := v.(type) {
		case *big.Int:
			return any(v).(T), true
		case int:
			return any(new(big.Int).SetInt64(int64(v))).(T), true
		case float64:
			return any(new(big.Int).SetInt64(int64(v))).(T), true
		case gojq.JQValue:
			return CastFn[T](v.JQValueToGoJQ(), structFn)
		default:
			return t, false
		}
	case string:
		switch v := v.(type) {
		case string:
			return any(v).(T), true
		case gojq.JQValue:
			return CastFn[T](v.JQValueToGoJQ(), structFn)
		default:
			return t, false
		}
	case map[string]any:
		switch v := v.(type) {
		case map[string]any:
			return any(v).(T), true
		case nil:
			// return empty instantiated map, not nil map
			return any(map[string]any{}).(T), true
		case gojq.JQValue:
			return CastFn[T](v.JQValueToGoJQ(), structFn)
		default:
			return t, false
		}
	case []any:
		switch v := v.(type) {
		case []any:
			return any(v).(T), true
		case nil:
			return t, true
		case gojq.JQValue:
			return CastFn[T](v.JQValueToGoJQ(), structFn)
		default:
			return t, false
		}
	default:
		ft := reflect.TypeOf(&t)
		if ft.Elem().Kind() == reflect.Struct {
			// TODO: some way to allow decode value passthru?
			m, err := ToGoJQValue(v)
			if err != nil {
				return t, false
			}
			if structFn == nil {
				panic("structFn nil")
			}
			err = structFn(m, &t)
			if err != nil {
				return t, false
			}

			return t, true
		} else if ft.Elem().Kind() == reflect.Interface {
			// TODO: panic on non any interface?
			// ignore failed type assert as v can be nil
			cv, ok := any(v).(T)
			if !ok && v != nil {
				return cv, false
			}

			return cv, true
		}

		panic(fmt.Sprintf("unsupported type %s", ft.Elem().Kind()))
	}
}

func Cast[T any](v any) (T, bool) {
	return CastFn[T](v, nil)
}

// convert to gojq compatible values and map scalars with fn
func NormalizeFn(v any, fn func(v any) any) any {
	switch v := v.(type) {
	case map[string]any:
		for k, e := range v {
			v[k] = NormalizeFn(e, fn)
		}
		return v
	case map[any]any:
		// for gopkg.in/yaml.v2
		vm := map[string]any{}
		for k, e := range v {
			switch i := k.(type) {
			case string:
				vm[i] = NormalizeFn(e, fn)
			case int:
				vm[strconv.Itoa(i)] = NormalizeFn(e, fn)
			}
		}
		return vm
	case []map[string]any:
		var vs []any
		for _, e := range v {
			vs = append(vs, NormalizeFn(e, fn))
		}
		return vs
	case []any:
		for i, e := range v {
			v[i] = NormalizeFn(e, fn)
		}
		return v
	case gojq.JQValue:
		return NormalizeFn(v.JQValueToGoJQ(), fn)
	default:
		return fn(v)
	}
}

// NormalizeToStrings normalizes to strings
// strings as is
// null to empty string
// others to JSON representation
func NormalizeToStrings(v any) any {
	return NormalizeFn(v, func(v any) any {
		r, _ := ToGoJQValue(v)
		switch r := r.(type) {
		case string:
			return r
		case nil:
			return ""
		default:
			b, _ := gojq.Marshal(r)
			return string(b)
		}
	})
}

func Normalize(v any) any {
	return NormalizeFn(v, func(v any) any { r, _ := ToGoJQValue(v); return r })
}

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
	return ExpectedObjectError{Typ: gojq.JQTypeArray}
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
		return HasKeyTypeError{L: gojq.JQTypeArray, R: fmt.Sprintf("%v", key)}
	}
	return intKey >= 0 && intKey < len(v)
}
func (v Array) JQValueType() string { return gojq.JQTypeArray }
func (v Array) JQValueToNumber() any {
	return FuncTypeNameError{Name: "tonumber", Typ: gojq.JQTypeArray}
}
func (v Array) JQValueToString() any {
	return FuncTypeNameError{Name: "tostring", Typ: gojq.JQTypeArray}
}
func (v Array) JQValueToGoJQ() any { return []any(v) }

// object

var _ gojq.JQValue = Object{}

type Object map[string]any

func (v Object) JQValueLength() any         { return len(v) }
func (v Object) JQValueSliceLen() any       { return ExpectedArrayError{Typ: gojq.JQTypeObject} }
func (v Object) JQValueIndex(index int) any { return ExpectedArrayError{Typ: gojq.JQTypeObject} }
func (v Object) JQValueSlice(start int, end int) any {
	return ExpectedArrayError{Typ: gojq.JQTypeObject}
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
		return HasKeyTypeError{L: gojq.JQTypeObject, R: fmt.Sprintf("%v", key)}
	}
	_, ok = v[stringKey]
	return ok
}
func (v Object) JQValueType() string { return gojq.JQTypeObject }
func (v Object) JQValueToNumber() any {
	return FuncTypeNameError{Name: "tonumber", Typ: gojq.JQTypeObject}
}
func (v Object) JQValueToString() any {
	return FuncTypeNameError{Name: "tostring", Typ: gojq.JQTypeObject}
}
func (v Object) JQValueToGoJQ() any { return map[string]any(v) }

// number

var _ gojq.JQValue = Number{}

type Number struct {
	V any
}

func (v Number) JQValueLength() any         { return v.V }
func (v Number) JQValueSliceLen() any       { return ExpectedArrayError{Typ: gojq.JQTypeNumber} }
func (v Number) JQValueIndex(index int) any { return ExpectedArrayError{Typ: gojq.JQTypeNumber} }
func (v Number) JQValueSlice(start int, end int) any {
	return ExpectedArrayError{Typ: gojq.JQTypeNumber}
}
func (v Number) JQValueKey(name string) any { return ExpectedObjectError{Typ: gojq.JQTypeNumber} }
func (v Number) JQValueEach() any           { return IteratorError{Typ: gojq.JQTypeNumber} }
func (v Number) JQValueKeys() any           { return FuncTypeNameError{Name: "keys", Typ: gojq.JQTypeNumber} }
func (v Number) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: gojq.JQTypeNumber}
}
func (v Number) JQValueType() string  { return gojq.JQTypeNumber }
func (v Number) JQValueToNumber() any { return v.V }
func (v Number) JQValueToString() any {
	b := &bytes.Buffer{}
	// uses colorjson encode based on gojq encoder to support big.Int
	if err := colorjson.NewEncoder(colorjson.Options{}).Marshal(v.V, b); err != nil {
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
func (v String) JQValueKey(name string) any          { return ExpectedObjectError{Typ: gojq.JQTypeString} }
func (v String) JQValueEach() any                    { return IteratorError{Typ: gojq.JQTypeString} }
func (v String) JQValueKeys() any                    { return FuncTypeNameError{Name: "keys", Typ: gojq.JQTypeString} }
func (v String) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: gojq.JQTypeString}
}
func (v String) JQValueType() string { return gojq.JQTypeString }
func (v String) JQValueToNumber() any {
	if !gojq.ValidNumber(string(v)) {
		return fmt.Errorf("invalid number: %q", string(v))
	}
	return gojq.NormalizeNumber(json.Number(string(v)))
}
func (v String) JQValueToString() any { return string(v) }
func (v String) JQValueToGoJQ() any   { return string(v) }

// boolean

var _ gojq.JQValue = Boolean(true)

type Boolean bool

func (v Boolean) JQValueLength() any {
	return FuncTypeNameError{Name: "length", Typ: gojq.JQTypeBoolean}
}
func (v Boolean) JQValueSliceLen() any       { return ExpectedArrayError{Typ: gojq.JQTypeBoolean} }
func (v Boolean) JQValueIndex(index int) any { return ExpectedArrayError{Typ: gojq.JQTypeBoolean} }
func (v Boolean) JQValueSlice(start int, end int) any {
	return ExpectedArrayError{Typ: gojq.JQTypeBoolean}
}
func (v Boolean) JQValueKey(name string) any { return ExpectedObjectError{Typ: gojq.JQTypeBoolean} }
func (v Boolean) JQValueEach() any           { return IteratorError{Typ: gojq.JQTypeBoolean} }
func (v Boolean) JQValueKeys() any           { return FuncTypeNameError{Name: "keys", Typ: gojq.JQTypeBoolean} }
func (v Boolean) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: gojq.JQTypeBoolean}
}
func (v Boolean) JQValueType() string { return gojq.JQTypeBoolean }
func (v Boolean) JQValueToNumber() any {
	return FuncTypeNameError{Name: "tonumber", Typ: gojq.JQTypeBoolean}
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
func (v Null) JQValueSliceLen() any                { return ExpectedArrayError{Typ: gojq.JQTypeNull} }
func (v Null) JQValueIndex(index int) any          { return ExpectedArrayError{Typ: gojq.JQTypeNull} }
func (v Null) JQValueSlice(start int, end int) any { return ExpectedArrayError{Typ: gojq.JQTypeNull} }
func (v Null) JQValueKey(name string) any          { return ExpectedObjectError{Typ: gojq.JQTypeNull} }

func (v Null) JQValueEach() any { return IteratorError{Typ: gojq.JQTypeNull} }
func (v Null) JQValueKeys() any { return FuncTypeNameError{Name: "keys", Typ: gojq.JQTypeNull} }
func (v Null) JQValueHas(key any) any {
	return FuncTypeNameError{Name: "has", Typ: gojq.JQTypeNull}
}
func (v Null) JQValueType() string  { return gojq.JQTypeNull }
func (v Null) JQValueToNumber() any { return FuncTypeNameError{Name: "tonumber", Typ: gojq.JQTypeNull} }
func (v Null) JQValueToString() any { return gojq.JQTypeNull }
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
	return HasKeyTypeError{L: gojq.JQTypeArray, R: fmt.Sprintf("%v", key)}
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
		return FuncTypeNameError{Name: "keys", Typ: gojq.JQTypeString}
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
