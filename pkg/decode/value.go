package decode

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"fq/internal/ranges"
	"fq/pkg/bitio"
	"io"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Bits uint64

func (b Bits) StringByteBits(base int) string {
	if b&0x7 != 0 {
		return strconv.FormatUint(uint64(b)>>3, base) + "+" + strconv.FormatUint(uint64(b)&0x7, base)
	}
	return strconv.FormatUint(uint64(b>>3), base)
}

func (b Bits) StringBits(base int) string {
	return strconv.FormatUint(uint64(b), base)
}

type BitRange ranges.Range

func (r BitRange) StringByteBits(base int) string {
	if r.Len == 0 {
		return fmt.Sprintf("%s-NA", Bits(r.Start).StringByteBits(base))
	}
	return fmt.Sprintf("%s-%s", Bits(r.Start).StringByteBits(base), Bits(r.Start+r.Len-1).StringByteBits(base))
}

func (r BitRange) StringBits(base int) string {
	if r.Len == 0 {
		return fmt.Sprintf("%s-NA", Bits(r.Start).StringBits(base))
	}
	return fmt.Sprintf("%s-%s", Bits(r.Start).StringBits(base), Bits(r.Start+r.Len-1).StringBits(base))
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

// TODO: encoding? endian, string encoding, compression, etc?
type Value struct {
	Parent        *Value
	V             interface{} // int64, uint64, float64, string, bool, []byte, Array, Struct
	Index         int         // index in parent array/struct
	Range         ranges.Range
	BitBuf        *bitio.Buffer
	IsRoot        bool
	Name          string
	MIME          string
	DisplayFormat DisplayFormat
	Symbol        string
	Desc          string
	Error         error
}

// TODO: base instead?

var lookupRe = regexp.MustCompile(`` +
	`^(?:` +
	`([\w_]+)` + // .name
	`|` + // or
	`\[(\d+)\]` + // [123]
	`)` +
	`(?:\.?)`) // dot separator

func (v *Value) Eval(exp string) (*Value, error) {
	lf := v.Lookup(exp)
	if lf == nil {
		return lf, fmt.Errorf("not found")
	}

	return lf, nil
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

func (v *Value) Lookup(path string) *Value {
	if path == "" {
		return v
	}

	lookupSM := lookupRe.FindStringSubmatch(path)
	if lookupSM == nil {
		return nil
	}
	rest := path[len(lookupSM[0]):]

	switch {
	case lookupSM == nil:
		return nil
	case lookupSM[1] != "": // struct lookup
		name := lookupSM[1]
		if s, ok := v.V.(Struct); ok {
			for _, f := range s {
				if f.Name == name {
					return f.Lookup(rest)
				}
			}
			return nil
		} else {
			return nil
		}
	case lookupSM[2] != "": // array lookup
		indexStr := lookupSM[2]
		index, _ := strconv.Atoi(indexStr)
		if a, ok := v.V.(Array); ok {
			return a[index].Lookup(rest)
		} else {
			return nil
		}
	default:
		panic("unreachable")
	}
}

type WalkFn func(v *Value, depth int, rootDepth int) error

var ErrWalkSkip = errors.New("skip")
var ErrWalkStop = errors.New("stop")

func (v *Value) walk(preOrder bool, fn WalkFn) error {
	var walkFn WalkFn
	walkFn = func(v *Value, depth int, rootDepth int) error {
		rootDepthDelta := 0
		if v.IsRoot {
			rootDepthDelta = 1
		}

		if preOrder {
			err := fn(v, depth, rootDepth+rootDepthDelta)
			switch err {
			case ErrWalkSkip:
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
				if err := walkFn(wv, depth+1, rootDepth+rootDepthDelta); err != nil {
					return err
				}
			}
		case Array:
			for _, wv := range v {
				if err := walkFn(wv, depth+1, rootDepth+rootDepthDelta); err != nil {
					return err
				}
			}
		}
		if !preOrder {
			err := fn(v, depth, rootDepth+rootDepthDelta)
			switch err {
			case ErrWalkSkip:
				return errors.New("can't skip in post-order")
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

	err := walkFn(v, 0, 0)
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
	_ = v.WalkPreOrder(func(v *Value, depth int, rootDepth int) error {
		if v.Error != nil {
			errs = append(errs, v.Error)
		}
		return nil
	})
	return errs
}

func (v *Value) postProcess() {
	v.WalkPostOrder(func(v *Value, depth int, rootDepth int) error {
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
		f = fmt.Sprintf("array %s", v.Name)
	case Struct:
		f = fmt.Sprintf("struct %s", v.Name)
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

	desc := ""
	if v.Desc != "" {
		desc = fmt.Sprintf(" (%s)", v.Desc)
	}

	return fmt.Sprintf("%s%s", s, desc)
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
		bs, _ := v.BitBuf.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
		return string(bs)
	case nil:
		return ""
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
	case bool:
		if vv {
			return true
		} else {
			return false
		}
	case int64:
		return big.NewInt(vv)
	case uint64:
		return big.NewInt(int64(vv))
	case float64:
		return vv
	case string:
		return vv
	case []byte:
		return string(vv)
	case *bitio.Buffer:
		// TODO: RawString, switch to writer somehow?
		vvLen := vv.Len()
		bs, _ := v.BitBuf.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
		return string(bs)
	case nil:
		return nil
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

func (v *Value) JsonLength() int {

	switch vv := v.V.(type) {
	case Struct:
		// log.Printf("JsonLength struct %d", len(vv)+5)

		return len(vv)
	case Array:
		//log.Printf("JsonLength array %d", len(vv))

		return len(vv)
	default:
		// log.Printf("JsonLength value 0")

		return 0
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

func (v *Value) JsonRange(start int, end int) []interface{} {

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

		panic("unreachable")
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
		r = v.Desc
	case "_range":
		r = map[string]interface{}{
			"start":  big.NewInt(v.Range.Start),
			"stop":   big.NewInt(v.Range.Stop()),
			"length": big.NewInt(v.Range.Len),
		}
	case "_size":
		r = big.NewInt(v.Range.Len)
	case "_raw":
		bb, err := v.BitBuf.BitBufRange(v.Range.Start, v.Range.Len)
		if err != nil {
			return err
		}
		r = bb
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

	//log.Printf("JsonProperty %s %#+v\n", name, r)

	return r
}

func (v *Value) JsonEach() [][2]interface{} {

	switch vv := v.V.(type) {
	case Struct:
		props := [][2]interface{}{}
		for _, f := range vv {
			props = append(props, [2]interface{}{f.Name, f})
		}
		// log.Printf("JsonEach struct %#+v", props)
		sort.Slice(props, func(i, j int) bool {
			return props[i][0].(string) < props[j][0].(string)
		})

		return props
	case Array:
		props := [][2]interface{}{}
		for i, f := range vv {
			props = append(props, [2]interface{}{i, f})
		}

		// log.Printf("JsonEach array %#+v", props)

		return props
	default:
		// log.Printf("JsonEach value nil")
		//panic("unreachable")

		return nil
	}

}
