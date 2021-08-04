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
		return decodeValueBase{
			dv:  dv,
			jqv: arrayBase{arrayValueObject{vv}},
		}
	case decode.Struct:
		return decodeValueBase{
			dv:  dv,
			jqv: objectBase{structValueObject{vv}},
		}
	case bool:
		return decodeValueBase{
			dv:  dv,
			jqv: booleanBase{booleanValue(vv)},
		}
	case int:
		return decodeValueBase{
			dv:  dv,
			jqv: numberBase{numberValue{vv}},
		}
	case int64:
		return decodeValueBase{
			dv: dv,
			// TODO: int() instead? on some cpus?
			jqv: numberBase{numberValue{big.NewInt(vv)}},
		}
	case uint64:
		return decodeValueBase{
			dv:  dv,
			jqv: numberBase{numberValue{v: new(big.Int).SetUint64(vv)}},
		}
	case float64:
		return decodeValueBase{
			dv:  dv,
			jqv: numberBase{numberValue{vv}},
		}
	case string:
		return decodeValueBase{
			dv:  dv,
			jqv: stringBase{stringValue(vv)},
		}
	case []byte:
		return decodeValueBase{
			dv:  dv,
			jqv: stringBase{stringValue(string(vv))},
		}
	case *bitio.Buffer:
		return decodeValueBase{
			dv:  dv,
			jqv: stringBase{stringBufferValueObject{vv}},
		}
	case []interface{}:
		return decodeValueBase{
			dv:  dv,
			jqv: arrayBase{arrayValue(vv)},
		}
	case map[string]interface{}:
		return decodeValueBase{
			dv:  dv,
			jqv: objectBase{objectValue(vv)},
		}
	case nil:
		return decodeValueBase{
			dv:  dv,
			jqv: nullBase{nullValue{}},
		}
	default:
		// TODO: error?
		panic("unreachable")
	}
}

var _ valueObjectIf = decodeValueBase{}

type decodeValueBase struct {
	dv  *decode.Value
	jqv gojq.JQValue // TODO: embed? need to rename interface JQValueToGoJQ() method and interface name collides
}

func (dvb decodeValueBase) DisplayName() string {
	if dvb.dv.Format != nil {
		return dvb.dv.Format.Name
	}
	if dvb.dv.Description != "" {
		return dvb.dv.Description
	}
	return ""
}

func (dvb decodeValueBase) Display(w io.Writer, opts Options) error { return dump(dvb.dv, w, opts) }
func (dvb decodeValueBase) Preview(w io.Writer, opts Options) error { return preview(dvb.dv, w, opts) }
func (dvb decodeValueBase) ToBuffer() (*bitio.Buffer, error) {
	return dvb.dv.RootBitBuf.Copy().BitBufRange(dvb.dv.Range.Start, dvb.dv.Range.Len)
}
func (dvb decodeValueBase) ToBufferRange() (bufferRange, error) {
	return bufferRange{bb: dvb.dv.RootBitBuf.Copy(), r: dvb.dv.Range}, nil
}
func (dvb decodeValueBase) ExtKeys() []string {
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

	if dvb.dv.Format != nil {
		kv = append(kv, "_format")
	}

	return kv
}
func (dvb decodeValueBase) JQValueLength() interface{} {
	return dvb.jqv.JQValueLength()
}
func (dvb decodeValueBase) JQValueIndex(index int) interface{} {
	return dvb.jqv.JQValueIndex(index)
}
func (dvb decodeValueBase) JQValueSlice(start int, end int) interface{} {
	return dvb.jqv.JQValueSlice(start, end)
}
func (dvb decodeValueBase) JQValueKey(name string) interface{} {
	if strings.HasPrefix(name, "_") {
		dv := dvb.dv

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
			return dvb.jqv.JQValueToGoJQ()
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
			if dvb.dv.Format == nil {
				return nil
			}
			return dvb.dv.Format.Name
		case "_unknown":
			return dvb.dv.Unknown
		}

		// TODO: error?
		return nil
	}
	return dvb.jqv.JQValueKey(name)
}
func (dvb decodeValueBase) JQValueEach() interface{} {
	return dvb.jqv.JQValueEach()
}
func (dvb decodeValueBase) JQValueKeys() interface{} {
	return dvb.jqv.JQValueKeys()
}
func (dvb decodeValueBase) JQValueHas(key interface{}) interface{} {
	return dvb.jqv.JQValueHas(key)
}
func (dvb decodeValueBase) JQValueType() string {
	return dvb.jqv.JQValueType()
}
func (dvb decodeValueBase) JQValueToNumber() interface{} {
	return dvb.jqv.JQValueToNumber()
}
func (dvb decodeValueBase) JQValueToString() interface{} {
	return dvb.jqv.JQValueToString()
}
func (dvb decodeValueBase) JQValueToGoJQ() interface{} { return dvb.jqv.JQValueToGoJQ() }

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
func (bv baseValueObject) JQValueHas(key interface{}) interface{} {
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
func (bv baseValueObject) JQValueToGoJQ() interface{} { return bv.vFn() }

// string

var _ JQString = stringValueObject{}

type stringValueObject struct {
	stringValue
}

// string (*bitio.Buffer)

var _ JQString = stringBufferValueObject{}

type stringBufferValueObject struct {
	*bitio.Buffer
}

func (v stringBufferValueObject) JQStringLength() interface{} {
	return int(v.Buffer.Len()) / 8
}
func (v stringBufferValueObject) JQStringIndex(index int) interface{} {
	// TODO: funcIndexSlice, string outside should return "" not null
	return v.JQStringSlice(index, index+1)
}
func (v stringBufferValueObject) JQStringSlice(start int, end int) interface{} {
	bb := v.Buffer.Copy()
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
func (v stringBufferValueObject) JQStringToNumber() interface{} {
	s, ok := v.JQStringToString().(string)
	if ok {
		gojq.NormalizeNumbers(s)
	}
	return s
}
func (v stringBufferValueObject) JQStringToString() interface{} {
	return v.JQStringSlice(0, int(v.Buffer.Len())/8)
}
func (v stringBufferValueObject) JQValueToGoJQ() interface{} {
	return v.JQStringToString()
}

// array

var _ JQArray = arrayValueObject{}

type arrayValueObject struct {
	decode.Array
}

func (v arrayValueObject) JQArrayLength() interface{} { return len(v.Array) }
func (v arrayValueObject) JQArrayIndex(index int) interface{} {
	return makeValueObject(v.Array[index])
}
func (v arrayValueObject) JQArraySlice(start int, end int) interface{} {
	vs := make([]interface{}, end-start)
	for i, e := range v.Array[start:end] {
		vs[i] = makeValueObject(e)
	}
	return vs
}
func (v arrayValueObject) JQArrayEach() interface{} {
	props := make([]gojq.PathValue, len(v.Array))
	for i, v := range v.Array {
		props[i] = gojq.PathValue{Path: i, Value: makeValueObject(v)}
	}
	return props
}
func (v arrayValueObject) JQArrayKeys() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i := range v.Array {
		vs[i] = i
	}
	return vs
}
func (v arrayValueObject) JQArrayHasKey(index int) interface{} {
	return index >= 0 && index < len(v.Array)
}
func (v arrayValueObject) JQValueToGoJQ() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i, v := range v.Array {
		vs[i] = makeValueObject(v)
	}
	return vs
}

// struct

var _ JQObject = structValueObject{}

type structValueObject struct {
	decode.Struct
}

func (v structValueObject) JQObjectLength() interface{} { return len(v.Struct) }
func (v structValueObject) JQObjectKey(name string) interface{} {
	for _, v := range v.Struct {
		if v.Name == name {
			return makeValueObject(v)
		}
	}
	return nil
}
func (v structValueObject) JQObjectEach() interface{} {
	props := make([]gojq.PathValue, len(v.Struct))
	for i, v := range v.Struct {
		props[i] = gojq.PathValue{Path: v.Name, Value: makeValueObject(v)}
	}
	sort.Slice(props, func(i, j int) bool {
		iString, _ := props[i].Path.(string)
		jString, _ := props[j].Path.(string)
		return iString < jString
	})
	return props
}
func (v structValueObject) JQObjectKeys() interface{} {
	vs := make([]interface{}, len(v.Struct))
	for i, v := range v.Struct {
		vs[i] = v.Name
	}
	return vs
}
func (v structValueObject) JQObjectHasKey(key string) interface{} {
	for _, f := range v.Struct {
		if f.Name == key {
			return true
		}
	}
	return false
}
func (v structValueObject) JQValueToGoJQ() interface{} {
	vm := make(map[string]interface{}, len(v.Struct))
	for _, v := range v.Struct {
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
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = arrayBase{}

type arrayBase struct {
	JQArray
}

func (v arrayBase) JQValueLength() interface{} {
	return v.JQArrayLength()
}
func (v arrayBase) JQValueIndex(index int) interface{} {
	return v.JQArrayIndex(index)
}
func (v arrayBase) JQValueSlice(start int, end int) interface{} {
	return v.JQArraySlice(start, end)
}
func (v arrayBase) JQValueKey(name string) interface{} {
	return expectedObjectError{typ: "array"}
}
func (v arrayBase) JQValueEach() interface{} {
	return v.JQArrayEach()
}
func (v arrayBase) JQValueKeys() interface{} {
	return v.JQArrayKeys()
}
func (v arrayBase) JQValueHas(key interface{}) interface{} {
	intKey, ok := key.(int)
	if !ok {
		return hasKeyTypeError{l: "array", r: fmt.Sprintf("%v", key)}
	}
	return v.JQArrayHasKey(intKey)
}
func (v arrayBase) JQValueType() string { return "array" }
func (v arrayBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "array"}
}
func (v arrayBase) JQValueToString() interface{} {
	return funcTypeError{name: "tostring", typ: "array"}
}
func (v arrayBase) JQValueToGoJQ() interface{} { return v.JQArray.JQValueToGoJQ() }

var _ JQArray = arrayValue{}

type arrayValue []interface{}

func (v arrayValue) JQArrayLength() interface{}                  { return len(v) }
func (v arrayValue) JQArrayIndex(index int) interface{}          { return v[index] }
func (v arrayValue) JQArraySlice(start int, end int) interface{} { return v[start:end] }
func (v arrayValue) JQArrayEach() interface{} {
	vs := make([]gojq.PathValue, len(v))
	for i, v := range v {
		vs[i] = gojq.PathValue{Path: i, Value: v}
	}
	return vs
}
func (v arrayValue) JQArrayKeys() interface{} {
	vs := make([]interface{}, len(v))
	for i := range v {
		vs[i] = i
	}
	return vs
}
func (v arrayValue) JQArrayHasKey(key int) interface{} {
	return key >= 0 && key < len(v)
}
func (v arrayValue) JQValueToGoJQ() interface{} { return []interface{}(v) }

// object

type JQObject interface {
	JQObjectLength() interface{}
	JQObjectKey(name string) interface{}
	JQObjectEach() interface{}
	JQObjectKeys() interface{}
	JQObjectHasKey(key string) interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = objectBase{}

type objectBase struct {
	JQObject
}

func (v objectBase) JQValueLength() interface{} {
	return v.JQObjectLength()
}
func (v objectBase) JQValueIndex(index int) interface{} {
	return expectedArrayError{typ: "object"}
}
func (v objectBase) JQValueSlice(start int, end int) interface{} {
	return expectedArrayError{typ: "object"}
}
func (v objectBase) JQValueKey(name string) interface{} {
	return v.JQObjectKey(name)
}
func (v objectBase) JQValueEach() interface{} {
	return v.JQObjectEach()
}
func (v objectBase) JQValueKeys() interface{} {
	return v.JQObjectKeys()
}
func (v objectBase) JQValueHas(key interface{}) interface{} {
	stringKey, ok := key.(string)
	if !ok {
		return hasKeyTypeError{l: "object", r: fmt.Sprintf("%v", key)}
	}
	return v.JQObjectHasKey(stringKey)
}
func (v objectBase) JQValueType() string { return "object" }
func (v objectBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "object"}
}
func (v objectBase) JQValueToString() interface{} {
	return funcTypeError{name: "tostring", typ: "object"}
}
func (v objectBase) JQValueToGoJQ() interface{} { return v.JQObject.JQValueToGoJQ() }

var _ JQObject = objectValue{}

type objectValue map[string]interface{}

func (v objectValue) JQObjectLength() interface{} { return len(v) }
func (v objectValue) JQObjectKey(name string) interface{} {
	return v[name]
}
func (v objectValue) JQObjectEach() interface{} {
	vs := make([]gojq.PathValue, len(v))
	i := 0
	for k, v := range v {
		vs[i] = gojq.PathValue{Path: k, Value: v}
		i++
	}
	return vs
}
func (v objectValue) JQObjectKeys() interface{} {
	vs := make([]interface{}, len(v))
	i := 0
	for k := range v {
		vs[i] = k
		i++
	}
	return vs
}
func (v objectValue) JQObjectHasKey(key string) interface{} {
	_, ok := v[key]
	return ok
}
func (v objectValue) JQValueToGoJQ() interface{} { return map[string]interface{}(v) }

// number

type JQNumber interface {
	JQNumberLength() interface{}
	JQNumberToNumber() interface{}
	JQNumberToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = numberBase{}

type numberBase struct {
	JQNumber
}

func (v numberBase) JQValueLength() interface{} {
	return v.JQNumberLength()
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
func (v numberBase) JQValueHas(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "number"}
}
func (v numberBase) JQValueType() string { return "number" }
func (v numberBase) JQValueToNumber() interface{} {
	return v.JQNumberToNumber()
}
func (v numberBase) JQValueToString() interface{} {
	return v.JQNumberToString()
}
func (v numberBase) JQValueToGoJQ() interface{} { return v.JQNumber.JQValueToGoJQ() }

var _ JQNumber = numberValue{}

// TODO: per number type?
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
func (v numberValue) JQValueToGoJQ() interface{} { return v.v }

// string

type JQString interface {
	JQStringLength() interface{}
	JQStringIndex(index int) interface{}
	JQStringSlice(start int, end int) interface{}
	JQStringToNumber() interface{}
	JQStringToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = stringBase{}

type stringBase struct {
	JQString
}

func (v stringBase) JQValueLength() interface{} {
	return v.JQStringLength()
}
func (v stringBase) JQValueIndex(index int) interface{} {
	return v.JQStringIndex(index)
}
func (v stringBase) JQValueSlice(start int, end int) interface{} {
	return v.JQStringSlice(start, end)
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
func (v stringBase) JQValueHas(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "string"}
}
func (v stringBase) JQValueType() string { return "string" }
func (v stringBase) JQValueToNumber() interface{} {
	return v.JQStringToNumber()
}
func (v stringBase) JQValueToString() interface{} {
	return v.JQStringToString()
}
func (v stringBase) JQValueToGoJQ() interface{} { return v.JQString.JQValueToGoJQ() }

var _ JQString = stringValue("")

type stringValue string

func (v stringValue) JQStringLength() interface{} {
	return len(v)
}
func (v stringValue) JQStringIndex(index int) interface{} {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return ""
	}
	return fmt.Sprintf("%c", v[index])
}
func (v stringValue) JQStringSlice(start int, end int) interface{} {
	return string(v[start:end])
}
func (v stringValue) JQStringToNumber() interface{} {
	return gojq.NormalizeNumbers(string(v))
}
func (v stringValue) JQStringToString() interface{} {
	return string(v)
}
func (v stringValue) JQValueToGoJQ() interface{} { return string(v) }

// boolean

type JQBoolean interface {
	JQBooleanToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = booleanBase{}

type booleanBase struct {
	JQBoolean
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
func (v booleanBase) JQValueHas(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "boolean"}
}
func (v booleanBase) JQValueType() string { return "boolean" }
func (v booleanBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "boolean"}
}
func (v booleanBase) JQValueToString() interface{} {
	return v.JQBooleanToString()
}
func (v booleanBase) JQValueToGoJQ() interface{} { return v.JQBoolean.JQValueToGoJQ() }

var _ JQBoolean = booleanValue(true)

type booleanValue bool

func (v booleanValue) JQBooleanToString() interface{} {
	if v {
		return "true"
	}
	return "false"
}
func (v booleanValue) JQValueToGoJQ() interface{} { return bool(v) }

// null

type JQNull interface {
	JQNullLength() interface{}
	JQNullToString() interface{}
	JQValueToGoJQ() interface{}
}

var _ gojq.JQValue = nullBase{}

type nullBase struct {
	JQNull
}

func (v nullBase) JQValueLength() interface{} {
	return v.JQNullLength()
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
func (v nullBase) JQValueHas(key interface{}) interface{} {
	return funcTypeError{name: "has", typ: "boolean"}
}
func (v nullBase) JQValueType() string { return "boolean" }
func (v nullBase) JQValueToNumber() interface{} {
	return funcTypeError{name: "tonumber", typ: "boolean"}
}
func (v nullBase) JQValueToString() interface{} {
	return v.JQNullToString()
}
func (v nullBase) JQValueToGoJQ() interface{} { return v.JQNull.JQValueToGoJQ() }

var _ JQNull = nullValue{}

type nullValue struct{}

func (v nullValue) JQNullLength() interface{}   { return 0 }
func (v nullValue) JQNullToString() interface{} { return "null" }
func (v nullValue) JQValueToGoJQ() interface{}  { return nil }
