package interp

import (
	"bytes"
	"errors"
	"fq/internal/gojqextra"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"io"
	"math/big"
	"sort"
	"strings"

	"github.com/itchyny/gojq"
)

// TODO: rename
type valueObjectIf interface {
	InterpObject
	ToBuffer
}

func makeValueObject(dv *decode.Value) decodeValueBase {
	switch vv := dv.V.(type) {
	case decode.Array:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.ArrayBase{JQArray: arrayValueObject{vv}},
		}
	case decode.Struct:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.ObjectBase{JQObject: structValueObject{vv}},
		}
	case bool:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.BooleanBase{JQBoolean: gojqextra.BooleanValue(vv)},
		}
	case int:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.NumberBase{JQNumber: gojqextra.NumberValue{V: vv}},
		}
	case int64:
		return decodeValueBase{
			dv: dv,
			// TODO: int() instead? on some cpus?
			JQValue: gojqextra.NumberBase{JQNumber: gojqextra.NumberValue{V: big.NewInt(vv)}},
		}
	case uint64:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.NumberBase{JQNumber: gojqextra.NumberValue{V: new(big.Int).SetUint64(vv)}},
		}
	case float64:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.NumberBase{JQNumber: gojqextra.NumberValue{V: vv}},
		}
	case string:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.StringBase{JQString: gojqextra.StringValue(vv)},
		}
	case []byte:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.StringBase{JQString: gojqextra.StringValue(string(vv))},
		}
	case *bitio.Buffer:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.StringBase{JQString: stringBufferValueObject{vv}},
		}
	case []interface{}:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.ArrayBase{JQArray: gojqextra.ArrayValue(vv)},
		}
	case map[string]interface{}:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.ObjectBase{JQObject: gojqextra.ObjectValue(vv)},
		}
	case nil:
		return decodeValueBase{
			dv:      dv,
			JQValue: gojqextra.NullBase{JQNull: gojqextra.NullValue{}},
		}
	default:
		panic("unreachable")
	}
}

var _ valueObjectIf = decodeValueBase{}

type decodeValueBase struct {
	dv *decode.Value
	gojq.JQValue
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
			return dvb.JQValue.JQValueToGoJQ()
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
	return dvb.JQValue.JQValueKey(name)
}

// string (*bitio.Buffer)

var _ gojqextra.JQString = stringBufferValueObject{}

type stringBufferValueObject struct {
	*bitio.Buffer
}

func (v stringBufferValueObject) JQStringLength() interface{} {
	return int(v.Buffer.Len()) / 8
}
func (v stringBufferValueObject) JQStringIndex(index int) interface{} {
	if index < 0 {
		return ""
	}
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

// decode value array

var _ gojqextra.JQArray = arrayValueObject{}

type arrayValueObject struct {
	decode.Array
}

func (v arrayValueObject) JQArrayLength() interface{} { return len(v.Array) }
func (v arrayValueObject) JQArrayIndex(index int) interface{} {
	// -1 outside after string, -2 outside before string
	if index < 0 {
		return nil
	}
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
	for i, f := range v.Array {
		props[i] = gojq.PathValue{Path: i, Value: makeValueObject(f)}
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
	for i, f := range v.Array {
		vs[i] = makeValueObject(f).JQValueToGoJQ()
	}
	return vs
}

// decode value struct

var _ gojqextra.JQObject = structValueObject{}

type structValueObject struct {
	decode.Struct
}

func (v structValueObject) JQObjectLength() interface{} { return len(v.Struct) }
func (v structValueObject) JQObjectKey(name string) interface{} {
	for _, f := range v.Struct {
		if f.Name == name {
			return makeValueObject(f)
		}
	}
	return nil
}
func (v structValueObject) JQObjectEach() interface{} {
	props := make([]gojq.PathValue, len(v.Struct))
	for i, f := range v.Struct {
		props[i] = gojq.PathValue{Path: f.Name, Value: makeValueObject(f)}
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
	for i, f := range v.Struct {
		vs[i] = f.Name
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
	for _, f := range v.Struct {
		vm[f.Name] = makeValueObject(f).JQValueToGoJQ()
	}
	return vm
}
