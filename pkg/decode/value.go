package decode

// TODO: Encoding, u16le, varint etc, encode?
// TODO: Value/Compound interface? can have per type and save memory

import (
	"errors"
	"sort"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
)

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

type Compound struct {
	IsArray  bool
	Children *[]*Value

	Description string
	Format      *Format
	Err         error
}

type Scalar struct {
	Actual        interface{} // int, int64, uint64, float64, string, bool, []byte, *bitio.Buffer
	Sym           interface{}
	Description   string
	DisplayFormat DisplayFormat
	Unknown       bool
}

func (s Scalar) Value() interface{} {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}

type Value struct {
	Parent     *Value
	Name       string
	V          interface{} // Scalar, Array, Struct
	Index      int         // index in parent array/struct
	Range      ranges.Range
	RootBitBuf *bitio.Buffer
	IsRoot     bool // TODO: rework?
}

type WalkFn func(v *Value, rootV *Value, depth int, rootDepth int) error

var ErrWalkSkipChildren = errors.New("skip children")
var ErrWalkBreak = errors.New("break")
var ErrWalkStop = errors.New("stop")

type WalkOpts struct {
	PreOrder bool
	OneRoot  bool
	Fn       WalkFn
}

func (v *Value) Walk(opts WalkOpts) error {
	var walkFn WalkFn

	walkFn = func(wv *Value, rootV *Value, depth int, rootDepth int) error {
		if opts.OneRoot && wv != v && wv.IsRoot {
			return nil
		}

		rootDepthDelta := 0
		if wv.IsRoot {
			rootV = wv
			rootDepthDelta = 1
		}

		if opts.PreOrder {
			err := opts.Fn(wv, rootV, depth, rootDepth+rootDepthDelta)
			switch {
			case errors.Is(err, ErrWalkSkipChildren):
				return nil
			case errors.Is(err, ErrWalkStop):
				fallthrough
			default:
				if err != nil {
					return err
				}
			}
		}

		switch wvv := wv.V.(type) {
		case Compound:
			for _, wv := range *wvv.Children {
				if err := walkFn(wv, rootV, depth+1, rootDepth+rootDepthDelta); err != nil {
					if errors.Is(err, ErrWalkBreak) {
						break
					}
					return err
				}
			}
		}

		if !opts.PreOrder {
			err := opts.Fn(wv, rootV, depth, rootDepth+rootDepthDelta)
			switch {
			case errors.Is(err, ErrWalkSkipChildren):
				return errors.New("can't skip children in post-order")
			case errors.Is(err, ErrWalkStop):
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
	rootV := v.BufferRoot()

	err := walkFn(v, rootV, 0, 0)
	if errors.Is(err, ErrWalkStop) {
		err = nil
	}

	return err
}

func (v *Value) WalkPreOrder(fn WalkFn) error {
	return v.Walk(WalkOpts{
		PreOrder: true,
		Fn:       fn,
	})
}

func (v *Value) WalkPostOrder(fn WalkFn) error {
	return v.Walk(WalkOpts{
		PreOrder: false,
		Fn:       fn,
	})
}

func (v *Value) WalkRootPreOrder(fn WalkFn) error {
	return v.Walk(WalkOpts{
		PreOrder: true,
		OneRoot:  true,
		Fn:       fn,
	})
}

func (v *Value) WalkRootPostOrder(fn WalkFn) error {
	return v.Walk(WalkOpts{
		PreOrder: false,
		OneRoot:  true,
		Fn:       fn,
	})
}

func (v *Value) root(findSubRoot bool, findFormatRoot bool) *Value {
	rootV := v
	for rootV.Parent != nil {
		if findSubRoot && rootV.IsRoot {
			break
		}
		if findFormatRoot {
			if c, ok := rootV.V.(Compound); ok {
				if c.Format != nil {
					break
				}
			}
		}

		rootV = rootV.Parent
	}
	return rootV
}

func (v *Value) Root() *Value       { return v.root(false, false) }
func (v *Value) BufferRoot() *Value { return v.root(true, false) }
func (v *Value) FormatRoot() *Value { return v.root(true, true) }

func (v *Value) Errors() []error {
	var errs []error
	_ = v.WalkPreOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		switch vv := rootV.V.(type) {
		case Compound:
			if vv.Err != nil {
				errs = append(errs, vv.Err)
			}
		}
		return nil
	})
	return errs
}

func (v *Value) postProcess() {
	if err := v.WalkPostOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		switch vv := v.V.(type) {
		case Compound:
			first := true
			for _, f := range *vv.Children {
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

			// TODO: really sort array?
			sort.Slice(*vv.Children, func(i, j int) bool {
				return (*vv.Children)[i].Range.Start < (*vv.Children)[j].Range.Start
			})

			for i, f := range *vv.Children {
				f.Index = i
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func (v *Value) ScalarFn(sfns ...ScalarFn) error {
	var err error
	s, ok := v.V.(Scalar)
	if !ok {
		panic("not a scalar value")
	}
	for _, sfn := range sfns {
		s, err = sfn(s)
		if err != nil {
			break
		}
	}
	v.V = s
	return err
}
