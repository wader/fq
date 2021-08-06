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

type JQArray interface {
	JQArrayLength() interface{}
	// index -1 outside after string, -2 outside before string
	JQArrayIndex(index int) interface{}
	JQArraySlice(start int, end int) interface{}
	JQArrayEach() interface{}
	JQArrayKeys() interface{}
	JQArrayHasKey(key int) interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = ArrayBase{}

type ArrayBase struct {
	JQArray
}

func (v ArrayBase) JQValueLength() interface{} {
	return v.JQArrayLength()
}
func (v ArrayBase) JQValueSliceLen() interface{} {
	return v.JQArrayLength()
}
func (v ArrayBase) JQValueIndex(index int) interface{} {
	return v.JQArrayIndex(index)
}
func (v ArrayBase) JQValueSlice(start int, end int) interface{} {
	return v.JQArraySlice(start, end)
}
func (v ArrayBase) JQValueKey(name string) interface{} {
	return ExpectedObjectError{Typ: "array"}
}
func (v ArrayBase) JQValueEach() interface{} {
	return v.JQArrayEach()
}
func (v ArrayBase) JQValueKeys() interface{} {
	return v.JQArrayKeys()
}
func (v ArrayBase) JQValueHas(key interface{}) interface{} {
	intKey, ok := key.(int)
	if !ok {
		return HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
	}
	return v.JQArrayHasKey(intKey)
}
func (v ArrayBase) JQValueType() string { return "array" }
func (v ArrayBase) JQValueToNumber() interface{} {
	return FuncTypeError{Name: "tonumber", Typ: "array"}
}
func (v ArrayBase) JQValueToString() interface{} {
	return FuncTypeError{Name: "tostring", Typ: "array"}
}
func (v ArrayBase) JQValueToGoJQ() interface{} { return v.JQArray.JQValueToGoJQ() }

var _ JQArray = ArrayValue{}

type ArrayValue []interface{}

func (v ArrayValue) JQArrayLength() interface{}                  { return len(v) }
func (v ArrayValue) JQArrayIndex(index int) interface{}          { return v[index] }
func (v ArrayValue) JQArraySlice(start int, end int) interface{} { return v[start:end] }
func (v ArrayValue) JQArrayEach() interface{} {
	vs := make([]gojq.PathValue, len(v))
	for i, v := range v {
		vs[i] = gojq.PathValue{Path: i, Value: v}
	}
	return vs
}
func (v ArrayValue) JQArrayKeys() interface{} {
	vs := make([]interface{}, len(v))
	for i := range v {
		vs[i] = i
	}
	return vs
}
func (v ArrayValue) JQArrayHasKey(key int) interface{} {
	return key >= 0 && key < len(v)
}
func (v ArrayValue) JQValueToGoJQ() interface{} { return []interface{}(v) }

// object

type JQObject interface {
	JQObjectLength() interface{}
	JQObjectKey(name string) interface{}
	JQObjectEach() interface{}
	JQObjectKeys() interface{}
	JQObjectHasKey(key string) interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = ObjectBase{}

type ObjectBase struct {
	JQObject
}

func (v ObjectBase) JQValueLength() interface{} {
	return v.JQObjectLength()
}
func (v ObjectBase) JQValueSliceLen() interface{} {
	return ExpectedArrayError{Typ: "object"}
}
func (v ObjectBase) JQValueIndex(index int) interface{} {
	return ExpectedArrayError{Typ: "object"}
}
func (v ObjectBase) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "object"}
}
func (v ObjectBase) JQValueKey(name string) interface{} {
	return v.JQObjectKey(name)
}
func (v ObjectBase) JQValueEach() interface{} {
	return v.JQObjectEach()
}
func (v ObjectBase) JQValueKeys() interface{} {
	return v.JQObjectKeys()
}
func (v ObjectBase) JQValueHas(key interface{}) interface{} {
	stringKey, ok := key.(string)
	if !ok {
		return HasKeyTypeError{L: "object", R: fmt.Sprintf("%v", key)}
	}
	return v.JQObjectHasKey(stringKey)
}
func (v ObjectBase) JQValueType() string { return "object" }
func (v ObjectBase) JQValueToNumber() interface{} {
	return FuncTypeError{Name: "tonumber", Typ: "object"}
}
func (v ObjectBase) JQValueToString() interface{} {
	return FuncTypeError{Name: "tostring", Typ: "object"}
}
func (v ObjectBase) JQValueToGoJQ() interface{} { return v.JQObject.JQValueToGoJQ() }

var _ JQObject = ObjectValue{}

type ObjectValue map[string]interface{}

func (v ObjectValue) JQObjectLength() interface{} { return len(v) }
func (v ObjectValue) JQObjectKey(name string) interface{} {
	return v[name]
}
func (v ObjectValue) JQObjectEach() interface{} {
	vs := make([]gojq.PathValue, len(v))
	i := 0
	for k, v := range v {
		vs[i] = gojq.PathValue{Path: k, Value: v}
		i++
	}
	return vs
}
func (v ObjectValue) JQObjectKeys() interface{} {
	vs := make([]interface{}, len(v))
	i := 0
	for k := range v {
		vs[i] = k
		i++
	}
	return vs
}
func (v ObjectValue) JQObjectHasKey(key string) interface{} {
	_, ok := v[key]
	return ok
}
func (v ObjectValue) JQValueToGoJQ() interface{} { return map[string]interface{}(v) }

// number

type JQNumber interface {
	JQNumberLength() interface{}
	JQNumberToNumber() interface{}
	JQNumberToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = NumberBase{}

type NumberBase struct {
	JQNumber
}

func (v NumberBase) JQValueLength() interface{} {
	return v.JQNumberLength()
}
func (v NumberBase) JQValueSliceLen() interface{} {
	return ExpectedArrayError{Typ: "number"}
}
func (v NumberBase) JQValueIndex(index int) interface{} {
	return ExpectedArrayError{Typ: "number"}
}
func (v NumberBase) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "number"}
}
func (v NumberBase) JQValueKey(name string) interface{} {
	return ExpectedObjectError{Typ: "number"}
}
func (v NumberBase) JQValueEach() interface{} {
	return IteratorError{Typ: "number"}
}
func (v NumberBase) JQValueKeys() interface{} {
	return FuncTypeError{Name: "keys", Typ: "number"}
}
func (v NumberBase) JQValueHas(key interface{}) interface{} {
	return FuncTypeError{Name: "has", Typ: "number"}
}
func (v NumberBase) JQValueType() string { return "number" }
func (v NumberBase) JQValueToNumber() interface{} {
	return v.JQNumberToNumber()
}
func (v NumberBase) JQValueToString() interface{} {
	return v.JQNumberToString()
}
func (v NumberBase) JQValueToGoJQ() interface{} { return v.JQNumber.JQValueToGoJQ() }

var _ JQNumber = NumberValue{}

// TODO: per number type?
type NumberValue struct {
	V interface{}
}

func (v NumberValue) JQNumberLength() interface{}   { return v.V }
func (v NumberValue) JQNumberToNumber() interface{} { return v.V }
func (v NumberValue) JQNumberToString() interface{} {
	b := &bytes.Buffer{}
	// uses colorjson encode based on gojq encoder to support big.Int
	if err := colorjson.NewEncoder(false, false, 0, nil, colorjson.Colors{}).Marshal(v.V, b); err != nil {
		return err
	}
	return b.String()
}
func (v NumberValue) JQValueToGoJQ() interface{} { return v.V }

// string

type JQString interface {
	JQStringLength() interface{}
	// index -1 outside after string, -2 outside before string
	JQStringIndex(index int) interface{}
	JQStringSlice(start int, end int) interface{}
	JQStringToNumber() interface{}
	JQStringToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = StringBase{}

type StringBase struct {
	JQString
}

func (v StringBase) JQValueLength() interface{} {
	return v.JQStringLength()
}
func (v StringBase) JQValueSliceLen() interface{} {
	return v.JQStringLength()
}
func (v StringBase) JQValueIndex(index int) interface{} {
	return v.JQStringIndex(index)
}
func (v StringBase) JQValueSlice(start int, end int) interface{} {
	return v.JQStringSlice(start, end)
}
func (v StringBase) JQValueKey(name string) interface{} {
	return ExpectedObjectError{Typ: "string"}
}
func (v StringBase) JQValueEach() interface{} {
	return IteratorError{Typ: "string"}
}
func (v StringBase) JQValueKeys() interface{} {
	return FuncTypeError{Name: "keys", Typ: "string"}
}
func (v StringBase) JQValueHas(key interface{}) interface{} {
	return FuncTypeError{Name: "has", Typ: "string"}
}
func (v StringBase) JQValueType() string { return "string" }
func (v StringBase) JQValueToNumber() interface{} {
	return v.JQStringToNumber()
}
func (v StringBase) JQValueToString() interface{} {
	return v.JQStringToString()
}
func (v StringBase) JQValueToGoJQ() interface{} { return v.JQString.JQValueToGoJQ() }

var _ JQString = StringValue("")

type StringValue string

func (v StringValue) JQStringLength() interface{} {
	return len(v)
}
func (v StringValue) JQStringIndex(index int) interface{} {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return ""
	}
	return fmt.Sprintf("%c", v[index])
}
func (v StringValue) JQStringSlice(start int, end int) interface{} {
	return string(v[start:end])
}
func (v StringValue) JQStringToNumber() interface{} {
	return gojq.NormalizeNumbers(string(v))
}
func (v StringValue) JQStringToString() interface{} {
	return string(v)
}
func (v StringValue) JQValueToGoJQ() interface{} { return string(v) }

// boolean

type JQBoolean interface {
	JQBooleanToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = BooleanBase{}

type BooleanBase struct {
	JQBoolean
}

func (v BooleanBase) JQValueLength() interface{} {
	return FuncTypeError{Name: "length", Typ: "boolean"}
}
func (v BooleanBase) JQValueSliceLen() interface{} {
	return ExpectedArrayError{Typ: "boolean"}
}
func (v BooleanBase) JQValueIndex(index int) interface{} {
	return ExpectedArrayError{Typ: "boolean"}
}
func (v BooleanBase) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "boolean"}
}
func (v BooleanBase) JQValueKey(name string) interface{} {
	return ExpectedObjectError{Typ: "boolean"}
}
func (v BooleanBase) JQValueEach() interface{} {
	return IteratorError{Typ: "boolean"}
}
func (v BooleanBase) JQValueKeys() interface{} {
	return FuncTypeError{Name: "keys", Typ: "boolean"}
}
func (v BooleanBase) JQValueHas(key interface{}) interface{} {
	return FuncTypeError{Name: "has", Typ: "boolean"}
}
func (v BooleanBase) JQValueType() string { return "boolean" }
func (v BooleanBase) JQValueToNumber() interface{} {
	return FuncTypeError{Name: "tonumber", Typ: "boolean"}
}
func (v BooleanBase) JQValueToString() interface{} {
	return v.JQBooleanToString()
}
func (v BooleanBase) JQValueToGoJQ() interface{} { return v.JQBoolean.JQValueToGoJQ() }

var _ JQBoolean = BooleanValue(true)

type BooleanValue bool

func (v BooleanValue) JQBooleanToString() interface{} {
	if v {
		return "true"
	}
	return "false"
}
func (v BooleanValue) JQValueToGoJQ() interface{} { return bool(v) }

// null

type JQNull interface {
	JQNullLength() interface{}
	JQNullToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = NullBase{}

type NullBase struct {
	JQNull
}

func (v NullBase) JQValueLength() interface{} {
	return v.JQNullLength()
}
func (v NullBase) JQValueSliceLen() interface{} {
	return ExpectedArrayError{Typ: "null"}
}
func (v NullBase) JQValueIndex(index int) interface{} {
	return ExpectedArrayError{Typ: "null"}
}
func (v NullBase) JQValueSlice(start int, end int) interface{} {
	return ExpectedArrayError{Typ: "null"}
}
func (v NullBase) JQValueKey(name string) interface{} {
	return ExpectedObjectError{Typ: "null"}
}
func (v NullBase) JQValueEach() interface{} {
	return IteratorError{Typ: "null"}
}
func (v NullBase) JQValueKeys() interface{} {
	return FuncTypeError{Name: "keys", Typ: "null"}
}
func (v NullBase) JQValueHas(key interface{}) interface{} {
	return FuncTypeError{Name: "has", Typ: "null"}
}
func (v NullBase) JQValueType() string { return "null" }
func (v NullBase) JQValueToNumber() interface{} {
	return FuncTypeError{Name: "tonumber", Typ: "null"}
}
func (v NullBase) JQValueToString() interface{} {
	return v.JQNullToString()
}
func (v NullBase) JQValueToGoJQ() interface{} { return v.JQNull.JQValueToGoJQ() }

var _ JQNull = NullValue{}

type NullValue struct{}

func (v NullValue) JQNullLength() interface{}   { return 0 }
func (v NullValue) JQNullToString() interface{} { return "null" }
func (v NullValue) JQValueToGoJQ() interface{}  { return nil }
