package decode

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/fq/pkg/scalar"
)

type Compound struct {
	ByName      map[string]*Value
	Description string
	Children    []*Value
	IsArray     bool
}

// TODO: Encoding, u16le, varint etc, encode?
// TODO: Value/Compound interface? can have per type and save memory
// TODO: Make some fields optional somehow? map/slice?
type Value struct {
	V           any // scalar.S or Compound (array/struct)
	RootReader  bitio.ReaderAtSeeker
	Err         error
	Parent      *Value
	Format      *Format // TODO: rework
	Name        string
	Description string
	Range       ranges.Range
	Index       int  // index in parent array/struct
	IsRoot      bool // TODO: rework?
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
		// only count switching to a new root
		if wv.IsRoot && wv != rootV {
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
		case *Compound:
			for _, wv := range wvv.Children {
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
		if findFormatRoot && rootV.Format != nil {
			break
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
	_ = v.WalkPreOrder(func(v *Value, _ *Value, _ int, _ int) error {
		if v.Err != nil {
			errs = append(errs, v.Err)
		}
		return nil
	})
	return errs
}

func (v *Value) InnerRange() ranges.Range {
	if v.IsRoot {
		return ranges.Range{Start: 0, Len: v.Range.Len}
	}
	return v.Range
}

func (v *Value) postProcess() {
	if err := v.WalkRootPostOrder(func(v *Value, _ *Value, _ int, _ int) error {
		switch vv := v.V.(type) {
		case *Compound:
			first := true
			for _, f := range vv.Children {
				if f.IsRoot {
					continue
				}
				if s, ok := f.V.(scalar.Scalarable); ok && s.ScalarFlags().IsSynthetic() {
					continue
				}

				if first {
					v.Range = f.Range
					first = false
				} else {
					v.Range = ranges.MinMax(v.Range, f.Range)
				}
			}

			// sort struct fields and make sure to keep order if range is the same
			if !vv.IsArray {
				slices.SortStableFunc(vv.Children, func(a, b *Value) int {
					return cmp.Compare(a.Range.Start, b.Range.Start)
				})
			}

			v.Index = -1
			if vv.IsArray {
				for i, f := range vv.Children {
					f.Index = i
				}
			} else {
				for _, f := range vv.Children {
					f.Index = -1
				}
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

// TODO: rethink this
func (v *Value) TryUintScalarFn(sms ...scalar.UintMapper) error {
	var err error
	sr, ok := v.V.(*scalar.Uint)
	if !ok {
		panic("not a scalar value")
	}
	s := *sr
	for _, sm := range sms {
		s, err = sm.MapUint(s)
		if err != nil {
			break
		}
	}
	v.V = &s
	return err
}

func (v *Value) TryBitBufScalarFn(sms ...scalar.BitBufMapper) error {
	var err error
	sr, ok := v.V.(*scalar.BitBuf)
	if !ok {
		panic("not a scalar value")
	}
	s := *sr
	for _, sm := range sms {
		s, err = sm.MapBitBuf(s)
		if err != nil {
			break
		}
	}
	v.V = &s
	return err
}

func (v *Value) Remove() error {
	p := v.Parent
	if p == nil {
		return fmt.Errorf("d has no parent")
	}

	switch fv := p.V.(type) {
	case *Compound:
		if !fv.IsArray {
			if _, ok := fv.ByName[v.Name]; !ok {
				return fmt.Errorf("d not in parent ByName")
			}
			delete(fv.ByName, p.Name)
		}
		found := false
		var cs []*Value
		for _, c := range fv.Children {
			if c == v {
				found = true
				continue
			}
			cs = append(cs, c)
		}
		if !found {
			return fmt.Errorf("d not in parent children")
		}
		fv.Children = cs
	}
	return nil
}
