package decode

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/ranges"
	"io"
	"log"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

type Bits uint64

func (b Bits) StringByteBits(base int) string {
	if b&0x7 != 0 {
		return num.BasePrefixMap[base] + strconv.FormatUint(uint64(b)>>3, base) + "+" + strconv.FormatUint(uint64(b)&0x7, base)
	}
	return num.BasePrefixMap[base] + strconv.FormatUint(uint64(b>>3), base)
}

type BitRange ranges.Range

func (r BitRange) StringByteBits(base int) string {
	if r.Len == 0 {
		return fmt.Sprintf("%s-NA", Bits(r.Start).StringByteBits(base))
	}
	return fmt.Sprintf("%s-%s", Bits(r.Start).StringByteBits(base), Bits(r.Start+r.Len-1).StringByteBits(base))
}

type DisplayFormat int

const (
	NumberDecimal DisplayFormat = iota
	NumberBinary
	NumberOctal
	NumberHex
)

func DisplayFormatToBase(fmt DisplayFormat) int {
	switch fmt {
	case NumberDecimal:
		return 10
	case NumberBinary:
		return 2
	case NumberOctal:
		return 8
	case NumberHex:
		return 16
	default:
		return 0
	}
}

type Struct []*Value

type Array []*Value

// assert that *Value implements JSONObject
var _ gojq.JSONObject = &Value{}

// TODO: encoding? endian, string encoding, compression, etc?
type Value struct {
	Parent        *Value
	V             interface{} // int64, uint64, float64, string, bool, []byte, Array, Struct
	Index         int         // index in parent array/struct
	Range         ranges.Range
	RootBitBuf    *bitio.Buffer
	IsRoot        bool
	Name          string
	MIME          string
	DisplayFormat DisplayFormat
	Symbol        string
	Description   string
	Error         error
}

func (v *Value) Path() string {
	var parts []string

	for v.Parent != nil {
		switch v.Parent.V.(type) {
		case Struct:
			parts = append([]string{".", v.Name}, parts...)
		case Array:
			parts = append([]string{fmt.Sprintf("[%d]", v.Index)}, parts...)
		}
		v = v.Parent
	}

	if len(parts) == 0 {
		return "."
	}

	return strings.Join(parts, "")

}

type WalkFn func(v *Value, rootV *Value, depth int, rootDepth int) error

var ErrWalkSkipChildren = errors.New("skip children")
var ErrWalkStop = errors.New("stop")

func (v *Value) walk(preOrder bool, fn WalkFn) error {
	var walkFn WalkFn
	walkFn = func(v *Value, rootV *Value, depth int, rootDepth int) error {
		rootDepthDelta := 0
		if v.IsRoot {
			rootV = v
			rootDepthDelta = 1
		}

		if preOrder {
			err := fn(v, rootV, depth, rootDepth+rootDepthDelta)
			switch err {
			case ErrWalkSkipChildren:
				return nil
			case ErrWalkStop:
				fallthrough
			default:
				if err != nil {
					return err
				}
			}
		}

		switch v := v.V.(type) {
		case Struct:
			for _, wv := range v {
				if err := walkFn(wv, rootV, depth+1, rootDepth+rootDepthDelta); err != nil {
					return err
				}
			}
		case Array:
			for _, wv := range v {
				if err := walkFn(wv, rootV, depth+1, rootDepth+rootDepthDelta); err != nil {
					return err
				}
			}
		}
		if !preOrder {
			err := fn(v, rootV, depth, rootDepth+rootDepthDelta)
			switch err {
			case ErrWalkSkipChildren:
				return errors.New("can't skip children in post-order")
			case ErrWalkStop:
				fallthrough
			default:
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	// figure out root value for v as it might not be a root itself
	rootV := v
	for rootV != nil && !rootV.IsRoot {
		rootV = rootV.Parent
	}

	err := walkFn(v, rootV, 0, 0)
	if err == ErrWalkStop {
		err = nil
	}

	return err
}

func (v *Value) WalkPreOrder(fn WalkFn) error {
	return v.walk(true, fn)
}

func (v *Value) WalkPostOrder(fn WalkFn) error {
	return v.walk(false, fn)
}

func (v *Value) Errors() []error {
	var errs []error
	_ = v.WalkPreOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		if v.Error != nil {
			errs = append(errs, v.Error)
		}
		return nil
	})
	return errs
}

func (v *Value) postProcess() {
	v.WalkPostOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		switch vv := v.V.(type) {
		case Struct:
			first := true
			for _, f := range vv {
				if f.IsRoot {
					continue
				}

				if first {
					v.Range = f.Range
					first = false
				} else {
					v.Range = ranges.MinMax(v.Range, f.Range)
				}
			}

			sort.Slice(vv, func(i, j int) bool {
				return vv[i].Range.Start < vv[j].Range.Start
			})

			for i, f := range vv {
				f.Index = i
			}
		case Array:
			first := true
			for _, f := range vv {
				if f.IsRoot {
					continue
				}

				if first {
					v.Range = f.Range
					first = false
				} else {
					v.Range = ranges.MinMax(v.Range, f.Range)
				}
			}

			for i, f := range vv {
				f.Index = i
			}

			// TODO: also sort?
		}
		return nil
	})
}

func (v *Value) String() string {
	f := ""
	switch vv := v.V.(type) {
	case Array:
		f = "array"
	case Struct:
		f = "struct"
	case bool:
		f = "false"
		if vv {
			f = "true"
		}
	case int64:
		// TODO: DisplayFormat is weird
		f = strconv.FormatInt(vv, DisplayFormatToBase(v.DisplayFormat))
	case uint64:
		f = strconv.FormatUint(vv, DisplayFormatToBase(v.DisplayFormat))
	case float64:
		// TODO: float32? better truncated to significant digits?
		f = strconv.FormatFloat(vv, 'g', -1, 64)
	case string:
		f = vv
		if len(f) > 50 {
			f = fmt.Sprintf("%q", f[0:50]) + "..."
		} else {
			f = fmt.Sprintf("%q", vv)
		}
	case []byte:
		if len(vv) > 16 {
			f = hex.EncodeToString(vv[0:16]) + "..."

		} else {
			f = hex.EncodeToString(vv)
		}
	case *bitio.Buffer:
		vvLen := vv.Len()
		if vvLen > 16*8 {
			bs, _ := vv.BytesRange(0, 16)
			f = hex.EncodeToString(bs) + "..."
		} else {
			bs, _ := vv.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
			f = hex.EncodeToString(bs)
		}
	case nil:
		f = "none"
	default:
		panic("unreachable")
	}
	s := ""
	if v.Symbol != "" {
		s = fmt.Sprintf("%s (%s)", v.Symbol, f)
	} else {
		s = f
	}

	return s
}

func (v *Value) RawString() string {
	switch vv := v.V.(type) {
	case Array:
		return "array"
	case Struct:
		return "struct"
	case bool:
		if vv {
			return "1"
		} else {
			return "0"
		}
	case int64:
		// TODO: DisplayFormat is weird
		return strconv.FormatInt(vv, int(v.DisplayFormat))
	case uint64:
		return strconv.FormatUint(vv, int(v.DisplayFormat))
	case float64:
		return strconv.FormatFloat(vv, 'f', -1, 64)
	case string:
		return vv
	case []byte:
		return string(vv)
	case *bitio.Buffer:
		// TODO: RawString, switch to writer somehow?
		vvLen := vv.Len()
		bs, _ := v.RootBitBuf.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
		return string(bs)
	case nil:
		return ""
	default:
		panic("unreachable")
	}
}

func (v *Value) PreviewString() string {
	switch vv := v.V.(type) {
	case Array:
		return "[]"
	case Struct:
		return v.Description
	case bool:
		if vv {
			return "true"
		} else {
			return "false"
		}
	case int64:
		// TODO: DisplayFormat is weird
		return strconv.FormatInt(vv, DisplayFormatToBase(v.DisplayFormat))
	case uint64:
		return strconv.FormatUint(vv, DisplayFormatToBase(v.DisplayFormat))
	case float64:
		// TODO: float32? better truncated to significant digits?
		return strconv.FormatFloat(vv, 'g', -1, 64)
	case string:
		if len(vv) > 10 {
			return fmt.Sprintf("%s...", vv[0:10])
		} else {
			return vv
		}
	case []byte:
		if len(vv) > 16 {
			return hex.EncodeToString(vv[0:16]) + "..."
		} else {
			return hex.EncodeToString(vv)
		}
	case *bitio.Buffer:
		vvLen := vv.Len()
		if vvLen > 16*8 {
			bs, _ := vv.BytesRange(0, 16)
			return hex.EncodeToString(bs) + "..."
		} else {
			bs, _ := vv.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
			return hex.EncodeToString(bs)
		}
	case nil:
		return "nil"
	default:
		panic("unreachable")
	}
}

func (v *Value) ToJQ() interface{} {
	switch vv := v.V.(type) {
	case Array:
		return v
	case Struct:
		return v
	case int, bool, float64, string, nil:
		return vv
	case int64:
		return big.NewInt(vv)
	case uint64:
		return big.NewInt(int64(vv))
	case []byte:
		return string(vv)
	case *bitio.Buffer:
		// TODO: RawString, switch to writer somehow?
		bs, _ := v.RootBitBuf.BytesRange(v.Range.Start, int(bitio.BitsByteCount(v.Range.Len)))
		return string(bs)
	default:
		panic("unreachable")
	}
}

func (v *Value) MarshalJSON() ([]byte, error) {
	// TODO: range, bits etc?
	switch vv := v.V.(type) {
	case Array:
		arr := []interface{}{}
		for _, f := range vv {
			arr = append(arr, f)
		}
		return json.Marshal(arr)
	case Struct:
		obj := map[string]interface{}{}
		for _, f := range vv {
			obj[f.Name] = f.V
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

func (v *Value) JsonLength() interface{} {

	switch vv := v.V.(type) {
	case Struct:
		// log.Printf("JsonLength struct %d", len(vv)+5)

		return len(vv)
	case Array:
		//log.Printf("JsonLength array %d", len(vv))

		return len(vv)
	default:
		// log.Printf("JsonLength value 0")

		return fmt.Errorf("%v has no length", v)
	}
}

func (v *Value) JsonIndex(index int) interface{} {

	switch vv := v.V.(type) {
	case Struct:
		// log.Printf("JsonIndex struct %d nil", index)

		return nil
	case Array:
		// log.Printf("JsonIndex array %d %#+v", index, vv[index])

		return vv[index]
	default:
		// log.Printf("JsonIndex value %d nil", index)

		return nil
	}
}

func (v *Value) JsonRange(start int, end int) interface{} {
	switch vv := v.V.(type) {
	case Struct:
		// log.Printf("JsonRange struct %d-%d nil", start, end)

		return nil
	case Array:
		a := []interface{}{}
		for _, e := range vv[start:end] {
			a = append(a, e)
		}

		// log.Printf("JsonRange array %d-%d %#+v", start, end, a)

		return a
	default:
		// log.Printf("JsonRange value %d-%d nil", start, end)

		return fmt.Errorf("%v can't be indexed", v)
	}
}

func (v *Value) SpecialPropNames() []string {
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

func (v *Value) JsonProperty(name string) interface{} {

	// TODO: parent index useful?
	// TODO: mime, isroot

	var r interface{}
	switch name {
	case "_type":
		switch v.V.(type) {
		case Struct:
			return "struct"
		case Array:
			return "array"
		default:
			return "field"
		}
	case "_name":
		r = v.Name
	case "_value":
		r = v.ToJQ()
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
		r = v.Path()
	case "_error":
		if de, ok := v.Error.(*DecodeError); ok {
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
		case Struct:
			for _, f := range vv {
				if f.Name == name {
					r = f
					break
				}
			}
		case Array:
		default:
			//r = v
			//panic("unreachable")
		}
	}

	return r
}

func (v *Value) JsonEach() interface{} {
	props := [][2]interface{}{}
	switch vv := v.V.(type) {
	case Struct:
		for _, f := range vv {
			props = append(props, [2]interface{}{f.Name, f})
		}
	case Array:
		for i, f := range vv {
			props = append(props, [2]interface{}{i, f})
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

func (v *Value) JsonType() string {
	switch v.V.(type) {
	case Struct:
		return "object"
	case Array:
		return "array"
	default:
		return "field"
	}
}

func (v *Value) JsonPrimitiveValue() interface{} {
	switch vv := v.V.(type) {
	case Array:
		return v
	case Struct:
		return v
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

type decodeError2 struct {
	v *DecodeError
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
			if de, ok := e.(*DecodeError); ok {
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
		if de, ok := e.(*DecodeError); ok {
			errs = append(errs, &decodeError2{de})
		} else {
			errs = append(errs, e)
		}
	}

	var err interface{} = de.v.Err
	if de, ok := err.(*DecodeError); ok {
		err = &decodeError2{de}
	}

	return map[string]interface{}{

		"stack": de.v.PanicStack,
		"err":   de.v.Err,
		"errs":  errs,
	}
}

type bitBufObject struct {
	bb   *bitio.Buffer
	unit int
}

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
func (bo *bitBufObject) ToBifBuf() *bitio.Buffer {
	return bo.bb.Copy()
}
