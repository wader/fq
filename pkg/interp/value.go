package interp

import (
	"fmt"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/ranges"
	"io"
	"math/big"
	"sort"
	"strings"
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

func makeValueObject(dv *decode.Value) valueObjectIf {
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
		return baseValueObject{dv: dv, vFn: func() interface{} { return big.NewInt(int64(vv)) }, typ: "number"}
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

func (bv baseValueObject) ExtValueKeys() []string {
	kv := []string{
		"_type",
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
func (bv baseValueObject) JQValueProperty(name string) interface{} {
	if strings.HasPrefix(name, "_") {
		dv := bv.dv

		switch name {
		case "_type":
			return bv.typ
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
			switch err := dv.Err.(type) {
			case decode.FormatError:
				return formatError{err}
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
func (bv baseValueObject) JQValueType() string  { return bv.typ }
func (bv baseValueObject) JQValue() interface{} { return bv.vFn() }

// string

type stringValueObject struct {
	baseValueObject
	vv string
}

func (sv stringValueObject) ToBuffer() (*bitio.Buffer, error) {
	return bitio.NewBufferFromBytes([]byte(sv.vv), -1), nil
}
func (sv stringValueObject) ToBufferRange() (bufferRange, error) {
	bb := bitio.NewBufferFromBytes([]byte(sv.vv), -1)
	return bufferRange{bb: bb, r: ranges.Range{Start: 0, Len: bb.Len()}}, nil
}
func (sv stringValueObject) JQValueLength() interface{} { return len(sv.vv) }
func (sv stringValueObject) JQValueIndex(index int) interface{} {
	return fmt.Sprintf("%c", sv.vv[index])
}
func (sv stringValueObject) JQValueSlice(start int, end int) interface{} {
	return sv.vv[start:end]
}
func (sv stringValueObject) JQValue() interface{} {
	return sv.vv
}

// string (*bitio.Buffer)

type stringBufferValueObject struct {
	baseValueObject
	vv *bitio.Buffer
}

func (sv stringBufferValueObject) ToBuffer() (*bitio.Buffer, error) {
	return sv.vv.Copy(), nil
}
func (sv stringBufferValueObject) ToBufferRange() (bufferRange, error) {
	bb := sv.vv.Copy()
	return bufferRange{bb: sv.vv.Copy(), r: ranges.Range{Start: 0, Len: bb.Len()}}, nil
}
func (sv stringBufferValueObject) JQValueLength() interface{} {
	return int(sv.vv.Len()) / 8
}
func (sv stringBufferValueObject) JQValueIndex(index int) interface{} {
	bb := sv.vv.Copy()
	bb.SeekAbs(int64(index) * 8)
	s, err := bb.UTF8(1)
	if err != nil {
		return err
	}
	return s
}
func (sv stringBufferValueObject) JQValueSlice(start int, end int) interface{} {
	bb := sv.vv.Copy()
	bb.SeekAbs(int64(start) * 8)
	s, err := bb.UTF8(end - start)
	if err != nil {
		return err
	}
	return s
}
func (sv stringBufferValueObject) JQValue() interface{} {
	bb := sv.vv.Copy()
	s, err := bb.UTF8(int(bb.Len() / 8))
	if err != nil {
		return err
	}
	return s
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
	props := make([][2]interface{}, len(av.vv))
	for i, v := range av.vv {
		props[i] = [2]interface{}{i, makeValueObject(v)}
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
		vs[i] = makeValueObject(v).JQValue()
	}
	return vs
}

// struct

type structValueObject struct {
	baseValueObject
	vv decode.Struct
}

func (sv structValueObject) JQValueLength() interface{} { return len(sv.vv) }
func (sv structValueObject) JQValueProperty(name string) interface{} {
	if strings.HasPrefix(name, "_") {
		return sv.baseValueObject.JQValueProperty(name)
	}
	for _, v := range sv.vv {
		if v.Name == name {
			return makeValueObject(v)
		}
	}
	return nil
}
func (sv structValueObject) JQValueEach() interface{} {
	props := make([][2]interface{}, len(sv.vv))
	for i, v := range sv.vv {
		props[i] = [2]interface{}{v.Name, makeValueObject(v)}
	}
	sort.Slice(props, func(i, j int) bool {
		iString, _ := props[i][0].(string)
		jString, _ := props[j][0].(string)
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
		vm[v.Name] = makeValueObject(v).JQValue()
	}
	return vm
}
