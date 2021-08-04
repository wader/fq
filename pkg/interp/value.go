package interp

import (
	"bytes"
	"errors"
	"fmt"
	"fq/internal/colorjson"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"io"
	"math/big"
	"sort"
	"strings"

	"github.com/itchyny/gojq"
)

// TODO: refactor to use errors from gojq?
// TODO: preview errors

type funcTypeError struct {
	name string
	typ  string
}

func (err funcTypeError) Error() string { return err.name + " cannot be applied to: " + err.typ }

type expectedObjectError struct {
	typ string
}

func (err expectedObjectError) Error() string {
	return "expected an object but got: " + err.typ
}

type expectedArrayError struct {
	typ string
}

func (err expectedArrayError) Error() string {
	return "expected an array but got: " + err.typ
}

type iteratorError struct {
	typ string
}

func (err iteratorError) Error() string {
	return "cannot iterate over: " + err.typ
}

type hasKeyTypeError struct {
	l, r string
}

func (err hasKeyTypeError) Error() string {
	return "cannot check whether " + err.l + " has a key: " + err.r
}

type valueObjectIf interface {
	InterpObject
	ToBuffer
}

func makeValueObject(dv *decode.Value) interface{} {
	switch vv := dv.V.(type) {
	case decode.Array:
		av := arrayValueObject{baseValueObject: baseValueObject{dv: dv, typ: "array"}, vv: vv}
		av.baseValueObject.vFn = av.JQValue
		return av
	case decode.Struct:
		sv := structValueObject{baseValueObject: baseValueObject{dv: dv, typ: "object"}, vv: vv}
		sv.baseValueObject.vFn = sv.JQValue
		return sv
	case bool:
		return baseValueObject{dv: dv, vFn: func() interface{} { return vv }, typ: "boolean"}
	case int, float64:
		return baseValueObject{dv: dv, vFn: func() interface{} { return vv }, typ: "number"}
	case string:
		sv := stringValueObject{baseValueObject: baseValueObject{dv: dv, typ: "string"}, vv: vv}
		sv.baseValueObject.vFn = sv.JQValue
		return sv
	case int64:
		return baseValueObject{dv: dv, vFn: func() interface{} { return big.NewInt(vv) }, typ: "number"}
	case uint64:
		return baseValueObject{dv: dv, vFn: func() interface{} { return new(big.Int).SetUint64(vv) }, typ: "number"}
	case []byte:
		sv := stringValueObject{baseValueObject: baseValueObject{dv: dv, typ: "string"}, vv: string(vv)}
		sv.baseValueObject.vFn = sv.JQValue
		return sv
	case *bitio.Buffer:
		sv := stringBufferValueObject{baseValueObject: baseValueObject{dv: dv, typ: "string"}, vv: vv}
		sv.baseValueObject.vFn = sv.JQValue
		return sv
	case decode.JSON:
		return vv.V
	case nil:
		return baseValueObject{dv: dv, vFn: func() interface{} { return nil }, typ: "null"}
	default:
		// TODO: error?
		panic("unreachable")
	}
}

var _ valueObjectIf = baseValueObject{}

type baseValueObject struct {
	dv  *decode.Value
	vFn func() interface{}
	typ string
}

func (bv baseValueObject) DisplayName() string {
	if bv.dv.Format != nil {
		return bv.dv.Format.Name
	}
	if bv.dv.Description != "" {
		return bv.dv.Description
	}
	return bv.typ
}
func (bv baseValueObject) Display(w io.Writer, opts Options) error { return dump(bv.dv, w, opts) }
func (bv baseValueObject) Preview(w io.Writer, opts Options) error { return preview(bv.dv, w, opts) }
func (bv baseValueObject) ToBuffer() (*bitio.Buffer, error) {
	return bv.dv.RootBitBuf.Copy().BitBufRange(bv.dv.Range.Start, bv.dv.Range.Len)
}
func (bv baseValueObject) ToBufferRange() (bufferRange, error) {
	return bufferRange{bb: bv.dv.RootBitBuf.Copy(), r: bv.dv.Range}, nil
}

func (bv baseValueObject) ExtKeys() []string {
	kv := []string{
		"_start",
		"_stop",
		"_len",
		"_name",
		"_value",
		"_symbol",
		"_description",
		"_path",
		"_bits",
		"_bytes",
		"_error",
		"_unknown",
	}

	if bv.dv.Format != nil {
		kv = append(kv, "_format")
	}

	return kv
}

func (bv baseValueObject) JQValueLength() interface{} {
	return funcTypeError{name: "length", typ: bv.typ}
}
func (bv baseValueObject) JQValueIndex(index int) interface{} {
	return expectedArrayError{typ: bv.typ}
}
func (bv baseValueObject) JQValueSlice(start int, end int) interface{} {
	return expectedArrayError{typ: bv.typ}
}
func (bv baseValueObject) JQValueKey(name string) interface{} {
	if strings.HasPrefix(name, "_") {
		dv := bv.dv

		switch name {
		case "_start":
			return big.NewInt(dv.Range.Start)
		case "_stop":
			return big.NewInt(dv.Range.Stop())
		case "_len":
			return big.NewInt(dv.Range.Len)
		case "_name":
			return dv.Name
		case "_value":
			return bv.vFn()
		case "_symbol":
			return dv.Symbol
		case "_description":
			return dv.Description
		case "_path":
			return valuePath(dv)
		case "_error":
			var formatErr decode.FormatError
			if errors.As(dv.Err, &formatErr) {
				return formatErr.Value()

			}

			return dv.Err
		case "_bits":
			bb, err := dv.RootBitBuf.BitBufRange(dv.Range.Start, dv.Range.Len)
			if err != nil {
				return err
			}
			return newBifBufObject(bb, 1)
		case "_bytes":
			bb, err := dv.RootBitBuf.BitBufRange(dv.Range.Start, dv.Range.Len)
			if err != nil {
				return err
			}
			return newBifBufObject(bb, 8)
		case "_format":
			if bv.dv.Format == nil {
				return nil
			}
			return bv.dv.Format.Name
		case "_unknown":
			return bv.dv.Unknown
		}

		// TODO: error?
		return nil
	}
	return expectedObjectError{typ: bv.typ}
}
func (bv baseValueObject) JQValueEach() interface{} {
	return iteratorError{typ: bv.typ}
}
func (bv baseValueObject) JQValueKeys() interface{} {
	return funcTypeError{name: "keys", typ: bv.typ}
}
func (bv baseValueObject) JQValueHasKey(key interface{}) interface{} {
	return hasKeyTypeError{l: bv.typ, r: fmt.Sprintf("%v", key)}
}
func (bv baseValueObject) JQValueType() string { return bv.typ }
func (bv baseValueObject) JQValueToNumber() interface{} {
	v := bv.vFn()
	switch bv.typ {
	case "number":
		return v
	case "string":
		return gojq.NormalizeNumbers(v.(string))
	default:
		return funcTypeError{name: "tonumber", typ: bv.typ}
	}
}
func (bv baseValueObject) JQValueToString() interface{} {
	v := bv.vFn()
	switch bv.typ {
	case "number":
		b := &bytes.Buffer{}
		if err := colorjson.NewEncoder(false, false, 0, nil, colorjson.Colors{}).Marshal(v, b); err != nil {
			return err
		}
		return b.String()
	case "string":
		return v
	default:
		return funcTypeError{name: "tostring", typ: bv.typ}
	}
}
func (bv baseValueObject) JQValue() interface{} { return bv.vFn() }

// string

type stringValueObject struct {
	baseValueObject
	vv string
}

func (sv stringValueObject) ToBuffer() (*bitio.Buffer, error) {
	return bitio.NewBufferFromBytes([]byte(sv.vv), -1), nil
}
func (sv stringValueObject) JQValueLength() interface{} { return len(sv.vv) }
func (sv stringValueObject) JQValueIndex(index int) interface{} {
	return fmt.Sprintf("%c", sv.vv[index])
}
func (sv stringValueObject) JQValueSlice(start int, end int) interface{} {
	return sv.vv[start:end]
}
func (sv stringValueObject) JQValueToString() interface{} {
	return sv.JQValue()
}
func (sv stringValueObject) JQValue() interface{} {
	return sv.vv
}

// string (*bitio.Buffer)

type stringBufferValueObject struct {
	baseValueObject
	vv *bitio.Buffer
}

func (sbv stringBufferValueObject) ToBuffer() (*bitio.Buffer, error) {
	return sbv.vv.Copy(), nil
}
func (sbv stringBufferValueObject) JQValueLength() interface{} {
	return int(sbv.vv.Len()) / 8
}
func (sbv stringBufferValueObject) JQValueIndex(index int) interface{} {
	return sbv.JQValueSlice(index, index+1)
}
func (sbv stringBufferValueObject) JQValueSlice(start int, end int) interface{} {
	bb := sbv.vv.Copy()
	if start != 0 {
		if _, err := bb.SeekAbs(int64(start) * 8); err != nil {
			return err
		}
	}
	b := &bytes.Buffer{}
	if _, err := io.CopyN(b, bb, int64(end-start)); err != nil {
		return err
	}
	return b.String()

}
func (sbv stringBufferValueObject) JQValueToString() interface{} {
	return sbv.JQValue()
}
func (sbv stringBufferValueObject) JQValue() interface{} {
	return sbv.JQValueSlice(0, int(sbv.vv.Len())/8)
}

// array

type arrayValueObject struct {
	baseValueObject
	vv decode.Array
}

func (av arrayValueObject) JQValueLength() interface{} { return len(av.vv) }
func (av arrayValueObject) JQValueIndex(index int) interface{} {
	return makeValueObject(av.vv[index])
}
func (av arrayValueObject) JQValueSlice(start int, end int) interface{} {
	vs := make([]interface{}, end-start)
	for i, e := range av.vv[start:end] {
		vs[i] = makeValueObject(e)
	}
	return vs
}
func (av arrayValueObject) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(av.vv))
	for i, v := range av.vv {
		props[i] = gojq.PathValue{Path: i, Value: makeValueObject(v)}
	}
	return props
}
func (av arrayValueObject) JQValueKeys() interface{} {
	vs := make([]interface{}, len(av.vv))
	for i := range av.vv {
		vs[i] = i
	}
	return vs
}
func (av arrayValueObject) JQValueHasKey(key interface{}) interface{} {
	// TODO: toInt? int64?
	i, iOk := key.(int)
	if !iOk {
		return hasKeyTypeError{l: av.typ, r: fmt.Sprintf("%v", key)}
	}
	return i >= 0 && i < len(av.vv)
}
func (av arrayValueObject) JQValue() interface{} {
	vs := make([]interface{}, len(av.vv))
	for i, v := range av.vv {
		vs[i] = makeValueObject(v)
	}
	return vs
}

// struct

type structValueObject struct {
	baseValueObject
	vv decode.Struct
}

func (sv structValueObject) JQValueLength() interface{} { return len(sv.vv) }
func (sv structValueObject) JQValueKey(name string) interface{} {
	if strings.HasPrefix(name, "_") {
		return sv.baseValueObject.JQValueKey(name)
	}
	for _, v := range sv.vv {
		if v.Name == name {
			return makeValueObject(v)
		}
	}
	return nil
}
func (sv structValueObject) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(sv.vv))
	for i, v := range sv.vv {
		props[i] = gojq.PathValue{Path: v.Name, Value: makeValueObject(v)}
	}
	sort.Slice(props, func(i, j int) bool {
		iString, _ := props[i].Path.(string)
		jString, _ := props[j].Path.(string)
		return iString < jString
	})
	return props
}
func (sv structValueObject) JQValueKeys() interface{} {
	vs := make([]interface{}, len(sv.vv))
	for i, v := range sv.vv {
		vs[i] = v.Name
	}
	return vs
}
func (sv structValueObject) JQValueHasKey(key interface{}) interface{} {
	s, sOk := key.(string)
	if !sOk {
		return hasKeyTypeError{l: sv.typ, r: fmt.Sprintf("%v", key)}
	}
	for _, f := range sv.vv {
		if f.Name == s {
			return true
		}
	}
	return false
}
func (sv structValueObject) JQValue() interface{} {
	vm := make(map[string]interface{}, len(sv.vv))
	for _, v := range sv.vv {
		vm[v.Name] = makeValueObject(v)
	}
	return vm
}

// array

type JQArray interface {
	JQArrayLength() interface{}
	JQArrayIndex(index int) interface{}
	JQArraySlice(start int, end int) interface{}
	JQArrayEach() interface{}
	JQArrayKeys() interface{}
	JQArrayHasKey(key int) interface{}
	JQValue() interface{}
}

var _ gojq.JQValue = arrayBase{}

type arrayBase struct {
	v JQArray
}

func (v arrayBase) JQValueLength() interface{} {
	return v.v.JQArrayLength()
}
func (v arrayBase) JQValueIndex(index int) interface{} {
	return v.v.JQArrayIndex(index)
}
func (v arrayBase) JQValueSlice(start int, end int) interface{} {
	return v.v.JQArraySlice(start, end)
}
func (v arrayBase) JQValueKey(name string) interface{} {
	return expectedObjectError{typ: "array"}
}
func (v arrayBase) JQValueEach() interface{} {
	return v.v.JQArrayEach()
}
func (v arrayBase) JQValueKeys() interface{} {
	return v.v.JQArrayKeys()
}
func (v arrayBase) JQValueHasKey(key interface{}) interface{} {
	intKey, ok := key.(int)
	if !ok {
		return hasKeyTypeError{l: "array", r: fmt.Sprintf("%v", key)}
	}
	return v.v.JQArrayHasKey(intKey)
}
func (v arrayBase) JQValueType() string { return "array" }
func (v arrayBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "array"}
}
func (v arrayBase) JQValueToString() interface{} {
	return funcTypeError{name: "tostring", typ: "array"}
}
func (v arrayBase) JQValue() interface{} { return v.v.JQValue() }

var _ JQArray = arrayValue{}

type arrayValue struct {
	v []interface{}
}

func (v arrayValue) JQArrayLength() interface{}                  { return len(v.v) }
func (v arrayValue) JQArrayIndex(index int) interface{}          { return v.v[index] }
func (v arrayValue) JQArraySlice(start int, end int) interface{} { return v.v[start:end] }
func (v arrayValue) JQArrayEach() interface{} {
	vs := make([]gojq.PathValue, len(v.v))
	for i, v := range v.v {
		vs[i] = gojq.PathValue{Path: i, Value: v}
	}
	return vs
}
func (v arrayValue) JQArrayKeys() interface{} {
	vs := make([]interface{}, len(v.v))
	for i := range v.v {
		vs[i] = i
	}
	return vs
}
func (v arrayValue) JQArrayHasKey(key int) interface{} {
	return key >= 0 && key < len(v.v)
}
func (v arrayValue) JQValue() interface{} { return v.v }

// object

type JQObject interface {
	JQObjectLength() interface{}
	JQObjectKey(name string) interface{}
	JQObjectEach() interface{}
	JQObjectKeys() interface{}
	JQObjectHasKey(key string) interface{}
	JQValue() interface{}
}

var _ gojq.JQValue = objectBase{}

type objectBase struct {
	v JQObject
}

func (v objectBase) JQValueLength() interface{} {
	return v.v.JQObjectLength()
}
func (v objectBase) JQValueIndex(index int) interface{} {
	return expectedArrayError{typ: "object"}
}
func (v objectBase) JQValueSlice(start int, end int) interface{} {
	return expectedArrayError{typ: "object"}
}
func (v objectBase) JQValueKey(name string) interface{} {
	return v.v.JQObjectKey(name)
}
func (v objectBase) JQValueEach() interface{} {
	return v.v.JQObjectEach()
}
func (v objectBase) JQValueKeys() interface{} {
	return v.v.JQObjectKeys()
}
func (v objectBase) JQValueHasKey(key interface{}) interface{} {
	stringKey, ok := key.(string)
	if !ok {
		return hasKeyTypeError{l: "object", r: fmt.Sprintf("%v", key)}
	}
	return v.v.JQObjectHasKey(stringKey)
}
func (v objectBase) JQValueType() string { return "object" }
func (v objectBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "object"}
}
func (v objectBase) JQValueToString() interface{} {
	return funcTypeError{name: "tostring", typ: "object"}
}
func (v objectBase) JQValue() interface{} { return v.v.JQValue() }

var _ JQObject = objectValue{}

type objectValue struct {
	v map[string]interface{}
}

func (v objectValue) JQObjectLength() interface{} { return len(v.v) }
func (v objectValue) JQObjectKey(name string) interface{} {
	return v.v[name]
}
func (v objectValue) JQObjectEach() interface{} {
	vs := make([]gojq.PathValue, len(v.v))
	i := 0
	for k, v := range v.v {
		vs[i] = gojq.PathValue{Path: k, Value: v}
		i++
	}
	return vs
}
func (v objectValue) JQObjectKeys() interface{} {
	vs := make([]interface{}, len(v.v))
	i := 0
	for k := range v.v {
		vs[i] = k
		i++
	}
	return vs
}
func (v objectValue) JQObjectHasKey(key string) interface{} {
	_, ok := v.v[key]
	return ok
}
func (v objectValue) JQValue() interface{} { return v.v }

// number

type JQNumber interface {
	JQNumberLength() interface{}
	JQNumberToNumber() interface{}
	JQNumberToString() interface{}
	JQValue() interface{}
}

var _ gojq.JQValue = numberBase{}

type numberBase struct {
	v JQNumber
}

func (v numberBase) JQValueLength() interface{} {
	return v.v.JQNumberLength()
}
func (v numberBase) JQValueIndex(index int) interface{} {
	return expectedArrayError{typ: "number"}
}
func (v numberBase) JQValueSlice(start int, end int) interface{} {
	return expectedArrayError{typ: "number"}
}
func (v numberBase) JQValueKey(name string) interface{} {
	return expectedObjectError{typ: "number"}
}
func (v numberBase) JQValueEach() interface{} {
	return iteratorError{typ: "number"}
}
func (v numberBase) JQValueKeys() interface{} {
	return funcTypeError{name: "keys", typ: "number"}
}
func (v numberBase) JQValueHasKey(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "number"}
}
func (v numberBase) JQValueType() string { return "number" }
func (v numberBase) JQValueToNumber() interface{} {
	return v.v.JQNumberToNumber()
}
func (v numberBase) JQValueToString() interface{} {
	return v.v.JQNumberToString()
}
func (v numberBase) JQValue() interface{} { return v.v.JQValue() }

var _ JQNumber = numberValue{}

type numberValue struct {
	v interface{}
}

func (v numberValue) JQNumberLength() interface{}   { return v.v }
func (v numberValue) JQNumberToNumber() interface{} { return v.v }
func (v numberValue) JQNumberToString() interface{} {
	b := &bytes.Buffer{}
	// uses colorjson encode based on gojq encoder to support big.Int
	if err := colorjson.NewEncoder(false, false, 0, nil, colorjson.Colors{}).Marshal(v.v, b); err != nil {
		return err
	}
	return b.String()
}
func (v numberValue) JQValue() interface{} { return v.v }

// string

type JQString interface {
	JQStringLength() interface{}
	JQStringIndex(index int) interface{}
	JQStringSlice(start int, end int) interface{}
	JQStringToNumber() interface{}
	JQStringToString() interface{}
	JQValue() interface{}
}

var _ gojq.JQValue = stringBase{}

type stringBase struct {
	v JQString
}

func (v stringBase) JQValueLength() interface{} {
	return v.v.JQStringLength()
}
func (v stringBase) JQValueIndex(index int) interface{} {
	return v.v.JQStringIndex(index)
}
func (v stringBase) JQValueSlice(start int, end int) interface{} {
	return v.v.JQStringSlice(start, end)
}
func (v stringBase) JQValueKey(name string) interface{} {
	return expectedObjectError{typ: "string"}
}
func (v stringBase) JQValueEach() interface{} {
	return iteratorError{typ: "string"}
}
func (v stringBase) JQValueKeys() interface{} {
	return funcTypeError{name: "keys", typ: "string"}
}
func (v stringBase) JQValueHasKey(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "string"}
}
func (v stringBase) JQValueType() string { return "string" }
func (v stringBase) JQValueToNumber() interface{} {
	return v.v.JQStringToNumber()
}
func (v stringBase) JQValueToString() interface{} {
	return v.v.JQStringToString()
}
func (v stringBase) JQValue() interface{} { return v.v.JQValue() }

var _ JQString = stringValue("")

type stringValue string

func (v stringValue) JQStringLength() interface{} {
	return len(v)
}
func (v stringValue) JQStringIndex(index int) interface{} {
	// TODO: funcIndexSlice, string outside should return "" not null
	return fmt.Sprintf("%c", v[index])
}
func (v stringValue) JQStringSlice(start int, end int) interface{} {
	return v[start:end]
}
func (v stringValue) JQStringToNumber() interface{} {
	return gojq.NormalizeNumbers(string(v))
}
func (v stringValue) JQStringToString() interface{} {
	return v
}
func (v stringValue) JQValue() interface{} { return v }

// boolean

type JQBoolean interface {
	JQBooleanToString() interface{}
	JQValue() interface{}
}

var _ gojq.JQValue = booleanBase{}

type booleanBase struct {
	v JQBoolean
}

func (v booleanBase) JQValueLength() interface{} {
	return funcTypeError{name: "length", typ: "boolean"}
}
func (v booleanBase) JQValueIndex(index int) interface{} {
	return expectedArrayError{typ: "boolean"}
}
func (v booleanBase) JQValueSlice(start int, end int) interface{} {
	return expectedArrayError{typ: "boolean"}
}
func (v booleanBase) JQValueKey(name string) interface{} {
	return expectedObjectError{typ: "boolean"}
}
func (v booleanBase) JQValueEach() interface{} {
	return iteratorError{typ: "boolean"}
}
func (v booleanBase) JQValueKeys() interface{} {
	return funcTypeError{name: "keys", typ: "boolean"}
}
func (v booleanBase) JQValueHasKey(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "boolean"}
}
func (v booleanBase) JQValueType() string { return "boolean" }
func (v booleanBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "boolean"}
}
func (v booleanBase) JQValueToString() interface{} {
	return v.v.JQBooleanToString()
}
func (v booleanBase) JQValue() interface{} { return v.v.JQValue() }

var _ JQBoolean = booleanValue(true)

type booleanValue bool

func (v booleanValue) JQBooleanToString() interface{} {
	if v {
		return "true"
	}
	return "false"
}
func (v booleanValue) JQValue() interface{} { return v }

// null

type JQNull interface {
	JQNullLength() interface{}
	JQNullToString() interface{}
	JQValue() interface{}
}

var _ gojq.JQValue = nullBase{}

type nullBase struct {
	v JQNull
}

func (v nullBase) JQValueLength() interface{} {
	return v.v.JQNullLength()
}
func (v nullBase) JQValueIndex(index int) interface{} {
	return expectedArrayError{typ: "boolean"}
}
func (v nullBase) JQValueSlice(start int, end int) interface{} {
	return expectedArrayError{typ: "boolean"}
}
func (v nullBase) JQValueKey(name string) interface{} {
	return expectedObjectError{typ: "boolean"}
}
func (v nullBase) JQValueEach() interface{} {
	return iteratorError{typ: "boolean"}
}
func (v nullBase) JQValueKeys() interface{} {
	return funcTypeError{name: "keys", typ: "boolean"}
}
func (v nullBase) JQValueHasKey(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "boolean"}
}
func (v nullBase) JQValueType() string { return "boolean" }
func (v nullBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "boolean"}
}
func (v nullBase) JQValueToString() interface{} {
	return v.v.JQNullToString()
}
func (v nullBase) JQValue() interface{} { return v.v.JQValue() }

var _ JQNull = nullValue{}

type nullValue struct{}

func (v nullValue) JQNullLength() interface{}   { return 0 }
func (v nullValue) JQNullToString() interface{} { return "null" }
func (v nullValue) JQValue() interface{}        { return nil }
