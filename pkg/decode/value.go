package decode

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"fq/internal/bitio"
	"fq/internal/ranges"
	"regexp"
	"sort"
	"strconv"
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

var ErrWalkSkip = errors.New("skip")
var ErrWalkStop = errors.New("stop")

func (v *Value) walk(preOrder bool, fn func(v *Value, depth int, rootDepth int) error) error {
	var walkFn func(v *Value, depth int, rootDepth int) error
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
	return walkFn(v, 0, 0)
}

func (v *Value) WalkPreOrder(fn func(v *Value, depth int, rootDepth int) error) error {
	return v.walk(true, fn)
}

func (v *Value) WalkPostOrder(fn func(v *Value, depth int, rootDepth int) error) error {
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
		vvLen, err := vv.Len()
		if err != nil {
			return err.Error()
		}
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
		vvLen, err := vv.Len()
		if err != nil {
			return err.Error()
		}
		bs, _ := v.BitBuf.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
		return string(bs)
	case nil:
		return ""
	default:
		panic("unreachable")
	}
}

func (v *Value) ToJQ() interface{} {
	obj := map[string]interface{}{
		"name":        v.Name,
		"field":       v,
		"description": v.Desc,
		"range":       []int64{v.Range.Start, v.Range.Len},
	}

	switch vv := v.V.(type) {
	case Struct:
		fields := map[string]interface{}{}
		for _, f := range vv {
			fields[f.Name] = f.ToJQ()
		}
		obj["value"] = fields
	case Array:
		fields := []interface{}{}
		for _, f := range vv {
			fields = append(fields, f.ToJQ())
		}
		obj["value"] = fields
	default:
		obj["value"] = v.V
	}

	return obj
}

func (v *Value) MarshalJSON() ([]byte, error) {
	return json.Marshal("test")
}
