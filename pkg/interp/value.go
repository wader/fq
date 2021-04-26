package interp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"fq/pkg/ranges"
	"io"
	"math/big"
	"sort"
)

// assert that *Value implements InterpObject and ToBitBuf
var _ InterpObject = (*valueObject)(nil)
var _ ToBitBuf = (*valueObject)(nil)

type valueObject struct {
	v *decode.Value
}

// TODO: jq function somehow?
func (vo valueObject) Path() string {
	return valuePath(vo.v)
}

func (vo valueObject) MarshalJSON() ([]byte, error) {
	v := vo.v

	// TODO: range, bits etc?
	switch vv := v.V.(type) {
	case decode.Array:
		arr := []interface{}{}
		for _, f := range vv {
			arr = append(arr, valueObject{v: f})
		}
		return json.Marshal(arr)
	case decode.Struct:
		obj := map[string]interface{}{}
		for _, f := range vv {
			obj[f.Name] = valueObject{v: f}
		}
		return json.Marshal(obj)
	case bool, int64, uint64, float64, string, []byte, nil:
		return json.Marshal(vv)
	case *bitio.Buffer:
		bb := &bytes.Buffer{}
		if _, err := io.Copy(bb, vv.Copy()); err != nil {
			return nil, err
		}
		return json.Marshal(bb.Bytes())
	default:
		panic("unreachable")
	}
}

func (vo valueObject) JQValueLength() interface{} {
	v := vo.v
	switch vv := v.V.(type) {
	case decode.Struct:
		// log.Printf("JsonLength struct %d", len(vv)+5)

		return len(vv)
	case decode.Array:
		//log.Printf("JsonLength array %d", len(vv))

		return len(vv)
	default:
		// log.Printf("JsonLength value 0")

		return fmt.Errorf("%v has no length", v)
	}
}

func (vo valueObject) JQValueIndex(index int) interface{} {
	v := vo.v

	switch vv := v.V.(type) {
	case decode.Struct:
		// log.Printf("JsonIndex struct %d nil", index)

		return nil
	case decode.Array:
		// log.Printf("JsonIndex array %d %#+v", index, vv[index])

		return valueObject{v: vv[index]}
	default:
		// log.Printf("JsonIndex value %d nil", index)

		return nil
	}
}

func (vo valueObject) JQValueSlice(start int, end int) interface{} {
	v := vo.v

	switch vv := v.V.(type) {
	case decode.Struct:
		// log.Printf("JQValueSlice struct %d-%d nil", start, end)

		return nil
	case decode.Array:
		a := []interface{}{}
		for _, e := range vv[start:end] {
			a = append(a, valueObject{v: e})
		}

		// log.Printf("JQValueSlice array %d-%d %#+v", start, end, a)

		return a
	default:
		// log.Printf("JQValueSlice value %d-%d nil", start, end)

		return fmt.Errorf("%v can't be indexed", v)
	}
}

func (vo valueObject) SpecialPropNames() []string {
	return []string{
		"_type",
		"_name",
		"_value",
		"_symbol",
		"_description",
		"_size",
		"_path",
		"_bits",
		"_bytes",
		"_error",
	}
}

func (vo valueObject) DisplayName() string {
	v := vo.v
	if v.Description != "" {
		return vo.v.Description
	}
	switch v.V.(type) {
	case decode.Struct:
		return "{}"
	case decode.Array:
		return "[]"
	default:
		return "field"
	}
}

func (vo valueObject) JQValueProperty(name string) interface{} {
	v := vo.v

	// TODO: parent index useful?
	// TODO: mime, isRoot

	var r interface{}
	switch name {
	case "_type":
		switch v.V.(type) {
		case decode.Struct:
			return "struct"
		case decode.Array:
			return "array"
		default:
			return "field"
		}
	case "_name":
		r = v.Name
	case "_value":
		r = vo.JQValue()
	case "_symbol":
		r = v.Symbol
	case "_description":
		r = v.Description
	case "_size":
		r = big.NewInt(bitio.BitsByteCount(v.Range.Len))
	case "_path":
		r = valuePath(v)
	case "_error":
		switch err := v.Err.(type) {
		case decode.FormatError:
			return formatError{err}
		}

		return v.Err

	case "_bits":
		bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
		if err != nil {
			return err
		}
		r = &bitBufObject{bb: bb, unit: 1, r: v.Range}
	case "_bytes":
		bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
		if err != nil {
			return err
		}
		r = &bitBufObject{bb: bb, unit: 8, r: v.Range}
	}

	if r == nil {
		switch vv := v.V.(type) {
		case decode.Struct:
			for _, f := range vv {
				if f.Name == name {
					r = valueObject{v: f}
					break
				}
			}
		case decode.Array:
		default:
			//r = v
			//panic("unreachable")
		}
	}

	return r
}

func (vo valueObject) JQValueEach() interface{} {
	v := vo.v

	props := [][2]interface{}{}
	switch vv := v.V.(type) {
	case decode.Struct:
		for _, f := range vv {
			props = append(props, [2]interface{}{f.Name, valueObject{v: f}})
		}
	case decode.Array:
		for i, f := range vv {
			props = append(props, [2]interface{}{i, valueObject{v: f}})
		}
	}

	// for _, p := range v.specialPropNames() {
	// 	props = append(props, [2]interface{}{p, v.JsonProperty(p)})
	// }

	sort.Slice(props, func(i, j int) bool {
		iString, iIsString := props[i][0].(string)
		jString, jIsString := props[j][0].(string)
		iInt, iIsInt := props[i][0].(string)
		jInt, jIsInt := props[j][0].(string)
		if iIsString && jIsString {
			return iString < jString
		} else if iIsInt && jIsInt {
			return iInt < jInt
		} else if iIsInt {
			return true
		}

		return false
	})

	return props
}

func (vo valueObject) JQValueKeys() interface{} {
	var kvs []interface{}

	v := vo.v
	switch vv := v.V.(type) {
	case decode.Struct:
		for _, f := range vv {
			kvs = append(kvs, f.Name)
		}
	case decode.Array:
		for i := range vv {
			kvs = append(kvs, i)
		}
	default:
		return fmt.Errorf("can't get keys from %v", v.V)
	}

	return kvs
}

func (vo valueObject) JQValueHasKey(key interface{}) interface{} {
	v := vo.v
	switch vv := v.V.(type) {
	case decode.Struct:
		s, sOk := key.(string)
		if !sOk {
			return fmt.Errorf("can't check key for %#v", v.V)
		}
		for _, f := range vv {
			if f.Name == s {
				return true
			}
		}
		return false
	case decode.Array:
		// TODO: toInt? int64?
		i, iOk := key.(int)
		if !iOk {
			return fmt.Errorf("can't check key for %#v", v.V)
		}
		return i >= 0 && i < len(vv)
	default:
		return fmt.Errorf("can't check key for %#v", v.V)
	}
}

func (vo valueObject) JQValueType() string {
	v := vo.v
	switch v.V.(type) {
	case decode.Struct:
		return "object"
	case decode.Array:
		return "array"
	case int, float64, int64, uint64:
		return "number"
	case bool:
		return "boolean"
	case string, []byte:
		return "string"
	case nil:
		return "null"
	default:
		return "field"
	}
}

func (vo valueObject) JQValue() interface{} {
	v := vo.v
	switch vv := v.V.(type) {
	case decode.Array:
		arr := []interface{}{}
		for _, f := range vv {
			arr = append(arr, valueObject{v: f}.JQValue())
		}
		return arr
	case decode.Struct:
		obj := map[string]interface{}{}
		for _, f := range vv {
			obj[f.Name] = valueObject{v: f}.JQValue()
		}
		return obj
	case int, bool, float64:
		return vv
	case string:
		return vv
	case int64:
		return big.NewInt(vv)
	case uint64:
		return big.NewInt(int64(vv))
	case []byte:
		return string(vv)
	case *bitio.Buffer:
		return fmt.Sprintf("<%s bytes>", num.Bits(v.Range.Len).StringByteBits(10))
		// bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
		// if err != nil {
		// 	return err
		// }
		// buf := &bytes.Buffer{}
		// if _, err := io.Copy(buf, bb.Copy()); err != nil {
		// 	return err
		// }
		// return buf.String()
	case nil:
		return vv
	default:
		// TODO: error?
		return nil
	}
}

func (vo valueObject) Display(w io.Writer, opts Options) error {
	return dump(vo.v, w, opts)
}

func (vo valueObject) Preview(w io.Writer, opts Options) error {
	return preview(vo.v, w, opts)
}

func (vo valueObject) ToBitBuf() (*bitio.Buffer, ranges.Range) {
	v := vo.v

	switch vv := v.V.(type) {
	case []byte:
		bb := bitio.NewBufferFromBytes(vv, -1)
		return bb, ranges.Range{Start: 0, Len: bb.Len()}
	default:
		return v.RootBitBuf.Copy(), v.Range
	}

}
