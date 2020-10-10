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

// TODO: encoding? endian, string encoding, compression, etc?
type Value struct {
	V             interface{} // int64, uint64, float64, string, bool, []byte, error etc
	Range         Range
	BitBuf        *bitbuf.Buffer
	MIME          string
	DisplayFormat DisplayFormat
	Symbol        string
	Desc          string
}

type Field struct {
	Name  string
	Value interface{} // Value, []Value (array) or []*Field (struct)
	Error error
}

var lookupRe = regexp.MustCompile(`^([\w_]*)(?:\[(\d+)\])?$`)

func (f *Field) Eval(exp string) (interface{}, error) {
	lf := f.Lookup(exp)
	if lf == nil {
		return lf, fmt.Errorf("not found")
	}

	return lf, nil
}

func (f *Field) Lookup(path string) interface{} {
	if path == "" {
		return f
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

	switch v := f.Value.(type) {
	case []*Field:
		for _, f := range v {
			if f.Name != name {
				continue
			}

			if index != -1 {
				if vs, ok := f.Value.([]Value); ok {
					return vs[index]
				}
				return nil
			}

			return f.Lookup(rest)
		}
	}

	return nil
}

func (f *Field) Walk(fn func(f *Field)) {
	var walkFn func(f *Field)
	walkFn = func(f *Field) {
		fn(f)
		switch v := f.Value.(type) {
		case []*Field:
			for _, wf := range v {
				walkFn(wf)
			}
		case []Value:
			for _, wv := range v {
				if vwf, ok := wv.V.(*Field); ok {
					walkFn(vwf)
				}
			}
		}
	}
}

func (f *Field) WalkValues(fn func(v Value)) {
	f.Walk(func(f *Field) {
		if v, ok := f.Value.(Value); ok {
			fn(v)
		}
	})
}

func (f *Field) Errors() []error {
	var errs []error
	f.Walk(func(f *Field) {
		if f.Error != nil {
			errs = append(errs, f.Error)
		}
	})
	return errs
}

func (f *Field) Start() int64 {
	switch v := f.Value.(type) {
	case []*Field:
		if len(v) > 0 {
			return v[0].Start()
		}
		// TODO: hmm
		return 0
	case []Value:
		if len(v) > 0 {
			return v[0].Range.Start
		}
		// TODO: hmm
		return 0
	case Value:
		return v.Range.Start
	case nil:
		// TODO: hmm
		return 0
	default:
		panic("field not sortable")
	}
}

func (f *Field) Sort() {
	vfs, _ := f.Value.([]*Field)
	if vfs == nil {
		return
	}

	sort.Slice(vfs, func(i, j int) bool {
		return vfs[i].Start() < vfs[j].Start()
	})

	for _, vf := range vfs {
		vf.Sort()
	}
}
