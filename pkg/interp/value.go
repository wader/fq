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

type notUpdateableError struct {
	Typ string
	Key string
}

func (err notUpdateableError) Error() string {
	return fmt.Sprintf("cannot update key %s for %s", err.Key, err.Typ)
}

// TODO: rename
type valueIf interface {
	InterpValue
	ToBuffer
}

func valueUnderscoreKey(name string, a, b func(name string) interface{}) interface{} {
	if strings.HasPrefix(name, "_") {
		return a(name)
	}
	return b(name)
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

var _ valueIf = decodeValue{}

type decodeValue struct {
	gojq.JQValue
	decodeValueBase
}

func (v decodeValue) JQValueKey(name string) interface{} {
	return valueUnderscoreKey(name, v.decodeValueBase.JQValueKey, v.JQValue.JQValueKey)
}

// string (*bitio.Buffer)

var _ valueIf = bufferDecodeValue{}

type bufferDecodeValue struct {
	gojqextra.Base
	decodeValueBase
	*bitio.Buffer
}

func NewStringBufferValueObject(dv *decode.Value, bb *bitio.Buffer) bufferDecodeValue {
	return bufferDecodeValue{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "string"},
		Buffer:          bb,
	}
}

func (v bufferDecodeValue) JQValueKey(name string) interface{} {
	return valueUnderscoreKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueKey)
}

func (v bufferDecodeValue) JQValueLength() interface{} {
	return int(v.Buffer.Len()) / 8
}
func (v bufferDecodeValue) JQValueIndex(index int) interface{} {
	if index < 0 {
		return ""
	}
	// TODO: funcIndexSlice, string outside should return "" not null
	return v.JQValueSlice(index, index+1)
}
func (v bufferDecodeValue) JQValueSlice(start int, end int) interface{} {
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
func (v bufferDecodeValue) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "string"}
}
func (v bufferDecodeValue) JQValueToNumber() interface{} {
	s, ok := v.JQValueToString().(string)
	if ok {
		gojq.NormalizeNumbers(s)
	}
	return s
}
func (v bufferDecodeValue) JQValueToString() interface{} {
	return v.JQValueSlice(0, int(v.Buffer.Len())/8)
}
func (v bufferDecodeValue) JQValueToGoJQ() interface{} {
	return v.JQValueToString()
}
func (v bufferDecodeValue) JQValueToGoJQEx(opts Options) interface{} {
	s, err := opts.BitsFormatFn(v.Buffer.Copy())
	if err != nil {
		return err
	}
	return s
}

// decode value array

var _ valueIf = arrayDecodeValue{}

type arrayDecodeValue struct {
	gojqextra.Base
	decodeValueBase
	decode.Array
}

func NewArrayDecodeValue(dv *decode.Value, a decode.Array) arrayDecodeValue {
	return arrayDecodeValue{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "array"},
		Array:           a,
	}
}

func (v arrayDecodeValue) JQValueKey(name string) interface{} {
	return valueUnderscoreKey(name, v.decodeValueBase.JQValueKey, v.Base.JQValueKey)
}
func (v arrayDecodeValue) JQValueSliceLen() interface{} { return len(v.Array) }
func (v arrayDecodeValue) JQValueLength() interface{}   { return len(v.Array) }
func (v arrayDecodeValue) JQValueIndex(index int) interface{} {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return nil
	}
	return makeDecodeValue(v.Array[index])
}
func (v arrayDecodeValue) JQValueSlice(start int, end int) interface{} {
	vs := make([]interface{}, end-start)
	for i, e := range v.Array[start:end] {
		vs[i] = makeDecodeValue(e)
	}
	return vs
}
func (v arrayDecodeValue) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "array"}
}
func (v arrayDecodeValue) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(v.Array))
	for i, f := range v.Array {
		props[i] = gojq.PathValue{Path: i, Value: makeDecodeValue(f)}
	}
	return props
}
func (v arrayDecodeValue) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i := range v.Array {
		vs[i] = i
	}
	return vs
}
func (v arrayDecodeValue) JQValueHas(key interface{}) interface{} {
	intKey, ok := key.(int)
	if !ok {
		return gojqextra.HasKeyTypeError{L: "array", R: fmt.Sprintf("%v", key)}
	}
	return intKey >= 0 && intKey < len(v.Array)
}
func (v arrayDecodeValue) JQValueToGoJQ() interface{} {
	vs := make([]interface{}, len(v.Array))
	for i, f := range v.Array {
		vs[i] = makeDecodeValue(f)
	}
	return vs
}

// decode value struct

var _ valueIf = structDecodeValue{}

type structDecodeValue struct {
	gojqextra.Base
	decodeValueBase
	decode.Struct
}

func NewStructDecodeValue(dv *decode.Value, s decode.Struct) structDecodeValue {
	return structDecodeValue{
		decodeValueBase: decodeValueBase{dv},
		Base:            gojqextra.Base{Typ: "object"},
		Struct:          s,
	}
}

func (v structDecodeValue) JQValueLength() interface{}   { return len(v.Struct) }
func (v structDecodeValue) JQValueSliceLen() interface{} { return len(v.Struct) }
func (v structDecodeValue) JQValueKey(name string) interface{} {
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
func (v structDecodeValue) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "object"}
}
func (v structDecodeValue) JQValueEach() interface{} {
	props := make([]gojq.PathValue, len(v.Struct))
	for i, f := range v.Struct {
		props[i] = gojq.PathValue{Path: f.Name, Value: makeDecodeValue(f)}
	}
	return props
}
func (v structDecodeValue) JQValueKeys() interface{} {
	vs := make([]interface{}, len(v.Struct))
	for i, f := range v.Struct {
		vs[i] = f.Name
	}
	return vs
}
func (v structDecodeValue) JQValueHas(key interface{}) interface{} {
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
func (v structDecodeValue) JQValueToGoJQ() interface{} {
	vm := make(map[string]interface{}, len(v.Struct))
	for _, f := range v.Struct {
		vm[f.Name] = makeDecodeValue(f)
	}
	return vm
}
