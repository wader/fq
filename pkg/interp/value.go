package interp

import (
	"bytes"
	"errors"
	"fmt"
	"fq/internal/gojqextra"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"io"
	"math/big"
	"strings"

	"github.com/itchyny/gojq"
)

type expectedExtkeyError struct {
	Key string
}

func (err expectedExtkeyError) Error() string {
	return "expected a extkey but got: " + err.Key
}

// TODO: rename
type valueObjectIf interface {
	InterpObject
	ToBuffer
}

func valueUnderscoreKey(name string, a, b func(name string) interface{}) interface{} {
	if strings.HasPrefix(name, "_") {
		return a(name)
	}
	return b(name)
}

func makeValueObject(dv *decode.Value) valueObjectIf {
	switch vv := dv.V.(type) {
	case decode.Array:
		return NewArrayValueObject(dv, vv)
	case decode.Struct:
		return NewStructValueObject(dv, vv)
	case *bitio.Buffer:
		return NewStringBufferValueObject(dv, vv)
	case bool:
		return valueObject{
			JQValue:         gojqextra.Boolean(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case int:
		return valueObject{
			JQValue:         gojqextra.Number{V: vv},
			decodeValueBase: decodeValueBase{dv},
		}
	case int64:
		return valueObject{
			JQValue:         gojqextra.Number{V: big.NewInt(vv)},
			decodeValueBase: decodeValueBase{dv},
		}
	case uint64:
		return valueObject{
			JQValue:         gojqextra.Number{V: new(big.Int).SetUint64(vv)},
			decodeValueBase: decodeValueBase{dv},
		}
	case float64:
		return valueObject{
			JQValue:         gojqextra.Number{V: vv},
			decodeValueBase: decodeValueBase{dv},
		}
	case string:
		return valueObject{
			JQValue:         gojqextra.String(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case []byte:
		return valueObject{
			JQValue:         gojqextra.String(string(vv)),
			decodeValueBase: decodeValueBase{dv},
		}
	case []interface{}:
		return valueObject{
			JQValue:         gojqextra.Array(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case map[string]interface{}:
		return valueObject{
			JQValue:         gojqextra.Object(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case nil:
		return valueObject{
			JQValue:         gojqextra.Null{},
			decodeValueBase: decodeValueBase{dv},
		}

	default:
		panic("unreachable")
	}
}

type decodeValueBase struct {
	dv *decode.Value
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

func (dvb decodeValueBase) JQValueKey(name string) interface{} {
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

	return expectedExtkeyError{Key: name}
}

// string (*bitio.Buffer)

var _ valueObjectIf = valueObject{}

type valueObject struct {
	gojq.JQValue
	decodeValueBase
}

func (v valueObject) JQValueKey(name string) interface{} {
	return valueUnderscoreKey(name, v.decodeValueBase.JQValueKey, v.JQValue.JQValueKey)
}

// string (*bitio.Buffer)

var _ valueObjectIf = stringBufferValueObject2{}

type stringBufferValueObject2 struct {
	gojqextra.Base
	decodeValueBase
	*bitio.Buffer
}

func NewStringBufferValueObject(dv *decode.Value, bb *bitio.Buffer) stringBufferValueObject2 {
	return stringBufferValueObject2{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "string"},
		Buffer:          bb,
	}
}

func (v stringBufferValueObject2) JQValueKey(name string) interface{} {
	return valueUnderscoreKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueKey)
}

func (v stringBufferValueObject2) JQValueLength() interface{} {
	return int(v.Buffer.Len()) / 8
}
func (v stringBufferValueObject2) JQValueIndex(index int) interface{} {
	if index < 0 {
		return ""
	}
	// TODO: funcIndexSlice, string outside should return "" not null
	return v.JQValueSlice(index, index+1)
}
func (v stringBufferValueObject2) JQValueSlice(start int, end int) interface{} {
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
func (v stringBufferValueObject2) JQValueToNumber() interface{} {
	s, ok := v.JQValueToString().(string)
	if ok {
		gojq.NormalizeNumbers(s)
	}
	return s
}
func (v stringBufferValueObject2) JQValueToString() interface{} {
	return v.JQValueSlice(0, int(v.Buffer.Len())/8)
}
func (v stringBufferValueObject2) JQValueToGoJQ() interface{} {
	return v.JQValueToString()
}
func (v stringBufferValueObject2) JQValueToGoJQEx(i *Interp) interface{} {
	return v.JQValueToGoJQ()
}

// decode value array

var _ valueObjectIf = arrayValueObject{}

type arrayValueObject struct {
	gojqextra.Base
	decodeValueBase
	decode.Array
}

func NewArrayValueObject(dv *decode.Value, a decode.Array) arrayValueObject {
	return arrayValueObject{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "array"},
		Array:           a,
	}
}

func (v arrayValueObject) JQValueKey(name string) interface{} {
	return valueUnderscoreKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueKey)
}
func (v arrayValueObject) JQValueSliceLen() interface{} { return len(v.Array) }
func (v arrayValueObject) JQValueLength() interface{}   { return len(v.Array) }
func (v arrayValueObject) JQValueIndex(index int) interface{} {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return nil
	}
	return makeValueObject(v.Array[index])
}
func (v arrayValueObject) JQValueSlice(start int, end int) interface{} {
	vs := make([]interface{}, end-start)
	for i, e := range v.Array[start:end] {
		vs[i] = makeValueObject(e)
	}
	return vs
}
func (v arrayValueObject) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(v.Array))
	for i, f := range v.Array {
		props[i] = gojq.PathValue{Path: i, Value: makeValueObject(f)}
	}
	return props
}
func (v arrayValueObject) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i := range v.Array {
		vs[i] = i
	}
	return vs
}
func (v arrayValueObject) JQValueHas(key interface{}) interface{} {
	intKey, ok := key.(int)
	if !ok {
		return gojqextra.HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
	}
	return intKey >= 0 && intKey < len(v.Array)
}
func (v arrayValueObject) JQValueToGoJQ() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i, f := range v.Array {
		vs[i] = makeValueObject(f)
	}
	return vs
}

// decode value struct

var _ valueObjectIf = structValueObject{}

type structValueObject struct {
	gojqextra.Base
	decodeValueBase
	decode.Struct
}

func NewStructValueObject(dv *decode.Value, s decode.Struct) structValueObject {
	return structValueObject{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "object"},
		Struct:          s,
	}
}

func (v structValueObject) JQValueLength() interface{}   { return len(v.Struct) }
func (v structValueObject) JQValueSliceLen() interface{} { return len(v.Struct) }
func (v structValueObject) JQValueKey(name string) interface{} {
	if strings.HasPrefix(name, "_") {
		return v.decodeValueBase.JQValueKey(name)
	}

	for _, f := range v.Struct {
		if f.Name == name {
			return makeValueObject(f)
		}
	}
	return nil
}
func (v structValueObject) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(v.Struct))
	for i, f := range v.Struct {
		props[i] = gojq.PathValue{Path: f.Name, Value: makeValueObject(f)}
	}
	return props
}
func (v structValueObject) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v.Struct))
	for i, f := range v.Struct {
		vs[i] = f.Name
	}
	return vs
}
func (v structValueObject) JQValueHas(key interface{}) interface{} {
	stringKey, ok := key.(string)
	if !ok {
		return gojqextra.HasKeyTypeError{L: "object", R: fmt.Sprintf("%v", key)}
	}
	for _, f := range v.Struct {
		if f.Name == stringKey {
			return true
		}
	}
	return false
}
func (v structValueObject) JQValueToGoJQ() interface{} {
	vm := make(map[string]interface{}, len(v.Struct))
	for _, f := range v.Struct {
		vm[f.Name] = makeValueObject(f)
	}
	return vm
}
