package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"io"
	"log"
	"math/big"
	"sort"
)

// assert that *Value implements QueryObject
var _ QueryObject = &valueObject{}

type valueObject struct {
	v *decode.Value
}

// TODO: jq function somehow?
func (vo valueObject) Path() string {
	return valuePath(vo.v)
}

func (vo valueObject) ToJQ() interface{} {
	v := vo.v
	switch vv := v.V.(type) {
	case decode.Array:
		return vo
	case decode.Struct:
		return vo
	case int, bool, float64, string, nil:
		return vv
	case int64:
		return big.NewInt(vv)
	case uint64:
		return big.NewInt(int64(vv))
	case []byte:
		return string(vv)
	// TODO:
	// case *bitio.Buffer:
	// 	// TODO: RawString, switch to writer somehow?
	// 	bs, _ := v.RootBitBuf.BytesRange(v.Range.Start, int(bitio.BitsByteCount(v.Range.Len)))
	// 	return string(bs)
	default:
		panic("unreachable")
	}
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

func (vo valueObject) JsonLength() interface{} {
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

func (vo valueObject) JsonIndex(index int) interface{} {
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

func (vo valueObject) JsonRange(start int, end int) interface{} {
	v := vo.v

	switch vv := v.V.(type) {
	case decode.Struct:
		// log.Printf("JsonRange struct %d-%d nil", start, end)

		return nil
	case decode.Array:
		a := []interface{}{}
		for _, e := range vv[start:end] {
			a = append(a, valueObject{v: e})
		}

		// log.Printf("JsonRange array %d-%d %#+v", start, end, a)

		return a
	default:
		// log.Printf("JsonRange value %d-%d nil", start, end)

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
		"_range",
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
		return "struct"
	case decode.Array:
		return "array"
	default:
		return "field"
	}
}

func (vo valueObject) JsonProperty(name string) interface{} {
	v := vo.v

	// TODO: parent index useful?
	// TODO: mime, isroot

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
		r = vo.ToJQ()
	case "_symbol":
		r = v.Symbol
	case "_description":
		r = v.Description
	case "_range":
		r = map[string]interface{}{
			"start":  big.NewInt(v.Range.Start),
			"stop":   big.NewInt(v.Range.Stop()),
			"length": big.NewInt(v.Range.Len),
		}
	case "_size":
		r = big.NewInt(v.Range.Len)
	case "_path":
		r = valuePath(v)
	case "_error":
		if de, ok := v.Error.(*decode.DecodeError); ok {
			return &decodeError2{de}
		}
		return v.Error

	case "_bits":
		bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
		if err != nil {
			return err
		}
		r = &bitBufObject{bb: bb, unit: 1}
	case "_bytes":
		bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
		if err != nil {
			return err
		}
		r = &bitBufObject{bb: bb, unit: 8}
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

func (vo valueObject) JsonEach() interface{} {
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

func (vo valueObject) JsonType() string {
	v := vo.v
	switch v.V.(type) {
	case decode.Struct:
		return "object"
	case decode.Array:
		return "array"
	default:
		return "field"
	}
}

func (vo valueObject) JsonPrimitiveValue() interface{} {
	v := vo.v
	switch vv := v.V.(type) {
	case decode.Array:
		arr := []interface{}{}
		for _, f := range vv {
			arr = append(arr, valueObject{v: f})
		}
		return arr
	case decode.Struct:
		obj := map[string]interface{}{}
		for _, f := range vv {
			obj[f.Name] = valueObject{v: f}
		}
		return obj
	case int, bool, float64, string, nil:
		return vv
	case int64:
		return big.NewInt(vv)
	case uint64:
		return big.NewInt(int64(vv))
	case []byte:
		return string(vv)
	default:
		// TODO: error?
		return nil
	}
}

func (vo valueObject) Display(w io.Writer, opts DisplayOptions) error {
	return dump(vo.v, w, opts)
}

func (vo valueObject) ToBifBuf() *bitio.Buffer {
	v := vo.v
	bb, err := v.RootBitBuf.BitBufRange(v.Range.Start, v.Range.Len)
	if err != nil {
		return nil
	}
	return bb.Copy()
}

type decodeError2 struct {
	v *decode.DecodeError
}

func (de *decodeError2) JsonLength() interface{} {
	log.Printf("JsonLength: %#+v\n", de)
	return nil
}
func (de *decodeError2) JsonIndex(index int) interface{} {
	log.Printf("JsonIndex: %#+v\n", de)

	return nil
}
func (de *decodeError2) JsonRange(start int, end int) interface{} {
	log.Printf("JsonRange: %#+v\n", de)

	return nil
}
func (de *decodeError2) JsonProperty(name string) interface{} {
	log.Printf("JsonProperty: %#+v\n", de)

	switch name {
	case "errs":
		var errs []interface{}
		for _, e := range de.v.Errs {
			if de, ok := e.(*decode.DecodeError); ok {
				errs = append(errs, &decodeError2{de})
			} else {
				errs = append(errs, e)
			}
		}
		return errs
	}

	return nil
}
func (de *decodeError2) JsonEach() interface{} {
	log.Printf("JsonEach: %#+v\n", de)

	return nil
}
func (de *decodeError2) JsonType() string {
	log.Printf("JsonType: %#+v\n", de)

	return "object"
}
func (de *decodeError2) JsonPrimitiveValue() interface{} {
	log.Printf("JsonPrimitiveValue: %#+v\n", de)

	var errs []interface{}
	for _, e := range de.v.Errs {
		if de, ok := e.(*decode.DecodeError); ok {
			errs = append(errs, &decodeError2{de})
		} else {
			errs = append(errs, e)
		}
	}

	var err interface{} = de.v.Err
	if de, ok := err.(*decode.DecodeError); ok {
		err = &decodeError2{de}
	}

	return map[string]interface{}{

		"stack": de.v.PanicStack,
		"err":   de.v.Err,
		"errs":  errs,
	}
}

var _ QueryObject = &bitBufObject{}

type bitBufObject struct {
	bb   *bitio.Buffer
	unit int
}

func (*bitBufObject) DisplayName() string        { return "buffer" }
func (*bitBufObject) SpecialPropNames() []string { return nil }

func (bo *bitBufObject) JsonLength() interface{} {
	return int(bo.bb.Len()) / bo.unit
}
func (bo *bitBufObject) JsonIndex(index int) interface{} {
	pos, err := bo.bb.Pos()
	if err != nil {
		return err
	}
	if _, err := bo.bb.SeekAbs(int64(index) * int64(bo.unit)); err != nil {
		return err
	}
	v, err := bo.bb.U(bo.unit)
	if err != nil {
		return err
	}
	if _, err := bo.bb.SeekAbs(pos); err != nil {
		return err
	}
	return int(v)
}
func (bo *bitBufObject) JsonRange(start int, end int) interface{} {
	rbb, err := bo.bb.BitBufRange(int64(start*bo.unit), int64((end-start)*bo.unit))
	if err != nil {
		return err
	}
	return &bitBufObject{bb: rbb, unit: bo.unit}
}
func (bo *bitBufObject) JsonProperty(name string) interface{} {
	return nil
}
func (bo *bitBufObject) JsonEach() interface{} {
	return nil
}
func (bo *bitBufObject) JsonType() string {
	return "buffer"
}
func (bo *bitBufObject) JsonPrimitiveValue() interface{} {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bo.bb.Copy()); err != nil {
		return err
	}
	return buf.String()
}

func (bo *bitBufObject) Display(w io.Writer, opts DisplayOptions) error {
	if opts.Raw {
		if _, err := io.Copy(w, bo.bb.Copy()); err != nil {
			return err
		}
		return nil
	}

	unitNames := map[int]string{
		1: "bits",
		8: "bytes",
	}
	unitName := unitNames[bo.unit]
	if unitName == "" {
		unitName = "units"
	}
	_, err := fmt.Fprintf(w, "<%d %s>\n", bo.bb.Len()/int64(bo.unit), unitName)
	return err
}

func (bo *bitBufObject) ToBifBuf() *bitio.Buffer {
	return bo.bb.Copy()
}
