package decode

import (
	"errors"
	"fmt"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/ranges"
	"sort"
	"strconv"
	"strings"
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
