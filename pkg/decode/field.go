package decode

import (
	"fmt"
	"fq/pkg/bitbuf"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// TODO: interface? Display(v interface{})
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
	V             interface{} // int64, uint64, float64, string, bool, []byte, Array, Struct
	Range         Range
	BitBuf        *bitbuf.Buffer
	Name          string
	MIME          string
	DisplayFormat DisplayFormat
	Symbol        string
	Desc          string
	Error         error
}

var lookupRe = regexp.MustCompile(`^([\w_]*)(?:\[(\d+)\])?$`)

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

	parts := strings.SplitN(path, ".", 2)
	first := parts[0]
	rest := ""
	if len(parts) > 1 {
		rest = parts[1]
	}

	index := -1
	firstSM := lookupRe.FindStringSubmatch(first)
	if firstSM == nil {
		return nil
	}
	name := firstSM[1]
	indexStr := firstSM[2]
	if indexStr != "" {
		index, _ = strconv.Atoi(indexStr)
	}

	switch v := v.V.(type) {
	case Struct:
		for _, f := range v {
			if f.Name != name {
				continue
			}

			if index != -1 {
				if vs, ok := f.V.(Array); ok {
					return vs[index]
				}
				return nil
			}

			return f.Lookup(rest)
		}
	}

	return nil
}

func (v *Value) Walk(fn func(v *Value, depth int) error) error {
	var walkFn func(v *Value, depth int) error
	walkFn = func(v *Value, depth int) error {
		if err := fn(v, depth); err != nil {
			return err
		}
		switch v := v.V.(type) {
		case Struct:
			for _, wv := range v {
				if err := walkFn(wv, depth+1); err != nil {
					return err
				}
			}
		case Array:
			for _, wv := range v {
				if err := walkFn(wv, depth+1); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return walkFn(v, 0)
}

func (v *Value) Errors() []error {
	var errs []error
	_ = v.Walk(func(v *Value, depth int) error {
		if v.Error != nil {
			errs = append(errs, v.Error)
		}
		return nil
	})
	return errs
}

func (v *Value) Sort() {
	vfs, _ := v.V.(Struct)
	if vfs == nil {
		return
	}

	sort.Slice(vfs, func(i, j int) bool {
		return vfs[i].Range.Start < vfs[j].Range.Start
	})

	for _, vf := range vfs {
		vf.Sort()
	}
}
