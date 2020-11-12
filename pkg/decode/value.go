package decode

import (
	"encoding/hex"
	"errors"
	"fmt"
	"fq/pkg/bitbuf"
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

type Range struct {
	Start int64
	Stop  int64
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func RangeMinMax(a, b Range) Range {
	return Range{Start: min(a.Start, b.Start), Stop: max(a.Stop, b.Stop)}
}

func (r Range) StringByteBits(base int) string {
	return fmt.Sprintf("%s-%s", Bits(r.Start).StringByteBits(base), Bits(r.Stop).StringByteBits(base))
}

func (r Range) StringBits(base int) string {
	return fmt.Sprintf("%s-%s", Bits(r.Start).StringBits(base), Bits(r.Stop).StringBits(base))
}

func (r Range) Length() int64 {
	return r.Stop - r.Start
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
	Range         Range
	BitBuf        *bitbuf.Buffer
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
	// TODO: find start/stop from Ranges instead? what if seekaround? concat bitbufs but want gaps? sort here, crash?
	// TDOO: if bitbuf set?

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
					v.Range = RangeMinMax(v.Range, f.Range)
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
					v.Range = RangeMinMax(v.Range, f.Range)
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
	switch iv := v.V.(type) {
	case Array:
		f = fmt.Sprintf("array %s", v.Name)
	case Struct:
		f = fmt.Sprintf("struct %s", v.Name)
	case bool:
		f = "false"
		if iv {
			f = "true"
		}
	case int64:
		// TODO: DisplayFormat is weird
		f = strconv.FormatInt(iv, DisplayFormatToBase(v.DisplayFormat))
	case uint64:
		f = strconv.FormatUint(iv, DisplayFormatToBase(v.DisplayFormat))
	case float64:
		// TODO: float32? better truncated to significant digits?
		f = strconv.FormatFloat(iv, 'g', -1, 64)
	case string:
		f = iv
		if len(f) > 50 {
			f = fmt.Sprintf("%q", f[0:50]) + "..."
		} else {
			f = fmt.Sprintf("%q", iv)
		}
	case []byte:
		if len(iv) > 16 {
			f = hex.EncodeToString(iv[0:16]) + "..."

		} else {
			f = hex.EncodeToString(iv)
		}
	case *bitbuf.Buffer:
		if iv.Len > 16*8 {
			bs, _ := iv.BytesBitRange(0, 16*8, 0)
			f = hex.EncodeToString(bs) + "..."
		} else {
			bs, _ := iv.BytesBitRange(0, iv.Len, 0)
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
		s = fmt.Sprintf("%s", f)
	}

	desc := ""
	if v.Desc != "" {
		desc = fmt.Sprintf(" (%s)", v.Desc)
	}

	return fmt.Sprintf("%s%s", s, desc)
}

func (v *Value) RawString() string {
	switch iv := v.V.(type) {
	case Array:
		return "array"
	case Struct:
		return "struct"
	case bool:
		if iv {
			return "1"
		} else {
			return "0"
		}
	case int64:
		// TODO: DisplayFormat is weird
		return strconv.FormatInt(iv, int(v.DisplayFormat))
	case uint64:
		return strconv.FormatUint(iv, int(v.DisplayFormat))
	case float64:
		return strconv.FormatFloat(iv, 'f', -1, 64)
	case string:
		return iv
	case []byte:
		return string(iv)
	case *bitbuf.Buffer:
		bs, _ := v.BitBuf.BytesBitRange(0, 16*8, 0)
		return string(bs)
	case nil:
		return ""
	default:
		panic("unreachable")
	}
}
