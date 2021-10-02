package interp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"

	"github.com/wader/gojq"
)

type expectedExtkeyError struct {
	Key string
}

func (err expectedExtkeyError) Error() string {
	return "expected a extkey but got: " + err.Key
}

type notUpdateableError struct {
	Typ string
	Key string
}

func (err notUpdateableError) Error() string {
	return fmt.Sprintf("cannot update key %s for %s", err.Key, err.Typ)
}

// TODO: rename
type valueIf interface {
	Value
	ToBuffer
}

func valueKey(name string, a, b func(name string) interface{}) interface{} {
	if strings.HasPrefix(name, "_") {
		return a(name)
	}
	return b(name)
}
func valueHas(key interface{}, a func(name string) interface{}, b func(key interface{}) interface{}) interface{} {
	stringKey, ok := key.(string)
	if ok && strings.HasPrefix(stringKey, "_") {
		if err, ok := a(stringKey).(error); ok {
			return err
		}
		return true
	}
	return b(key)
}

func makeDecodeValue(dv *decode.Value) valueIf {
	switch vv := dv.V.(type) {
	case decode.Array:
		return NewArrayDecodeValue(dv, vv)
	case decode.Struct:
		return NewStructDecodeValue(dv, vv)
	case *bitio.Buffer:
		return NewStringBufferValueObject(dv, vv)
	case bool:
		return decodeValue{
			JQValue:         gojqextra.Boolean(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case int:
		return decodeValue{
			JQValue:         gojqextra.Number{V: vv},
			decodeValueBase: decodeValueBase{dv},
		}
	case int64:
		return decodeValue{
			JQValue:         gojqextra.Number{V: big.NewInt(vv)},
			decodeValueBase: decodeValueBase{dv},
		}
	case uint64:
		return decodeValue{
			JQValue:         gojqextra.Number{V: new(big.Int).SetUint64(vv)},
			decodeValueBase: decodeValueBase{dv},
		}
	case float64:
		return decodeValue{
			JQValue:         gojqextra.Number{V: vv},
			decodeValueBase: decodeValueBase{dv},
		}
	case string:
		return decodeValue{
			JQValue:         gojqextra.String(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case []byte:
		return decodeValue{
			JQValue:         gojqextra.String(string(vv)),
			decodeValueBase: decodeValueBase{dv},
		}
	case []interface{}:
		return decodeValue{
			JQValue:         gojqextra.Array(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case map[string]interface{}:
		return decodeValue{
			JQValue:         gojqextra.Object(vv),
			decodeValueBase: decodeValueBase{dv},
		}
	case nil:
		return decodeValue{
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
func (dvb decodeValueBase) ToBuffer() (*bitio.Buffer, error) {
	return dvb.dv.RootBitBuf.Copy().BitBufRange(dvb.dv.Range.Start, dvb.dv.Range.Len)
}
func (dvb decodeValueBase) ToBufferView() (BufferView, error) {
	return BufferView{bb: dvb.dv.RootBitBuf.Copy(), r: dvb.dv.Range, unit: 8}, nil
}
func (dvb decodeValueBase) ExtKeys() []string {
	kv := []string{
		"_start",
		"_stop",
		"_len",
		"_name",
		"_root",
		"_buffer_root",
		"_format_root",
		"_parent",
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
	case "_root":
		return makeDecodeValue(dv.Root())
	case "_buffer_root":
		// TODO: rename?
		return makeDecodeValue(dv.BufferRoot())
	case "_format_root":
		// TODO: rename?
		return makeDecodeValue(dv.FormatRoot())
	case "_parent":
		if dv.Parent == nil {
			return nil
		}
		return makeDecodeValue(dv.Parent)
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

var _ valueIf = decodeValue{}

type decodeValue struct {
	gojq.JQValue
	decodeValueBase
}

func (v decodeValue) JQValueKey(name string) interface{} {
	return valueKey(name, v.decodeValueBase.JQValueKey, v.JQValue.JQValueKey)
}
func (v decodeValue) JQValueHas(key interface{}) interface{} {
	return valueHas(key, v.decodeValueBase.JQValueKey, v.JQValue.JQValueHas)
}

// string (*bitio.Buffer)

var _ valueIf = BufferDecodeValue{}

type BufferDecodeValue struct {
	gojqextra.Base
	decodeValueBase
	*bitio.Buffer
}

func NewStringBufferValueObject(dv *decode.Value, bb *bitio.Buffer) BufferDecodeValue {
	return BufferDecodeValue{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "string"},
		Buffer:          bb,
	}
}

func (v BufferDecodeValue) JQValueKey(name string) interface{} {
	return valueKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueKey)
}
func (v BufferDecodeValue) JQValueHas(key interface{}) interface{} {
	return valueHas(key, v.decodeValueBase.JQValueKey, v.Base.JQValueHas)
}
func (v BufferDecodeValue) JQValueLength() interface{} {
	return int(v.Buffer.Len()) / 8
}
func (v BufferDecodeValue) JQValueIndex(index int) interface{} {
	if index < 0 {
		return ""
	}
	// TODO: funcIndexSlice, string outside should return "" not null
	return v.JQValueSlice(index, index+1)
}
func (v BufferDecodeValue) JQValueSlice(start int, end int) interface{} {
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
func (v BufferDecodeValue) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "string"}
}
func (v BufferDecodeValue) JQValueToNumber() interface{} {
	s, ok := v.JQValueToString().(string)
	if ok {
		gojq.NormalizeNumbers(s)
	}
	return s
}
func (v BufferDecodeValue) JQValueToString() interface{} {
	return v.JQValueSlice(0, int(v.Buffer.Len())/8)
}
func (v BufferDecodeValue) JQValueToGoJQ() interface{} {
	return v.JQValueToString()
}
func (v BufferDecodeValue) JQValueToGoJQEx(opts Options) interface{} {
	s, err := opts.BitsFormatFn(v.Buffer.Copy())
	if err != nil {
		return err
	}
	return s
}

// decode value array

var _ valueIf = ArrayDecodeValue{}

type ArrayDecodeValue struct {
	gojqextra.Base
	decodeValueBase
	decode.Array
}

func NewArrayDecodeValue(dv *decode.Value, a decode.Array) ArrayDecodeValue {
	return ArrayDecodeValue{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "array"},
		Array:           a,
	}
}

func (v ArrayDecodeValue) JQValueKey(name string) interface{} {
	return valueKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueKey)
}
func (v ArrayDecodeValue) JQValueSliceLen() interface{} { return len(v.Array) }
func (v ArrayDecodeValue) JQValueLength() interface{}   { return len(v.Array) }
func (v ArrayDecodeValue) JQValueIndex(index int) interface{} {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return nil
	}
	return makeDecodeValue(v.Array[index])
}
func (v ArrayDecodeValue) JQValueSlice(start int, end int) interface{} {
	vs := make([]interface{}, end-start)
	for i, e := range v.Array[start:end] {
		vs[i] = makeDecodeValue(e)
	}
	return vs
}
func (v ArrayDecodeValue) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "array"}
}
func (v ArrayDecodeValue) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(v.Array))
	for i, f := range v.Array {
		props[i] = gojq.PathValue{Path: i, Value: makeDecodeValue(f)}
	}
	return props
}
func (v ArrayDecodeValue) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i := range v.Array {
		vs[i] = i
	}
	return vs
}
func (v ArrayDecodeValue) JQValueHas(key interface{}) interface{} {
	return valueHas(
		key,
		v.decodeValueBase.JQValueKey,
		func(key interface{}) interface{} {
			intKey, ok := key.(int)
			if !ok {
				return gojqextra.HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
			}
			return intKey >= 0 && intKey < len(v.Array)
		})
}
func (v ArrayDecodeValue) JQValueToGoJQ() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i, f := range v.Array {
		vs[i] = makeDecodeValue(f)
	}
	return vs
}

// decode value struct

var _ valueIf = StructDecodeValue{}

type StructDecodeValue struct {
	gojqextra.Base
	decodeValueBase
	decode.Struct
}

func NewStructDecodeValue(dv *decode.Value, s decode.Struct) StructDecodeValue {
	return StructDecodeValue{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "object"},
		Struct:          s,
	}
}

func (v StructDecodeValue) JQValueLength() interface{}   { return len(v.Struct) }
func (v StructDecodeValue) JQValueSliceLen() interface{} { return len(v.Struct) }
func (v StructDecodeValue) JQValueKey(name string) interface{} {
	if strings.HasPrefix(name, "_") {
		return v.decodeValueBase.JQValueKey(name)
	}

	for _, f := range v.Struct {
		if f.Name == name {
			return makeDecodeValue(f)
		}
	}
	return nil
}
func (v StructDecodeValue) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "object"}
}
func (v StructDecodeValue) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(v.Struct))
	for i, f := range v.Struct {
		props[i] = gojq.PathValue{Path: f.Name, Value: makeDecodeValue(f)}
	}
	return props
}
func (v StructDecodeValue) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v.Struct))
	for i, f := range v.Struct {
		vs[i] = f.Name
	}
	return vs
}
func (v StructDecodeValue) JQValueHas(key interface{}) interface{} {
	return valueHas(
		key,
		v.decodeValueBase.JQValueKey,
		func(key interface{}) interface{} {
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
		},
	)
}
func (v StructDecodeValue) JQValueToGoJQ() interface{} {
	vm := make(map[string]interface{}, len(v.Struct))
	for _, f := range v.Struct {
		vm[f.Name] = makeDecodeValue(f)
	}
	return vm
}
