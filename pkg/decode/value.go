package decode

import (
	"encoding/hex"
	"fmt"
	"fq/pkg/bitbuf"
	"regexp"
	"sort"
	"strconv"
)

type Bits uint64

func (b Bits) String() string {
	if b&0x7 != 0 {
		return strconv.FormatUint(uint64(b)>>3, 10) + "+" + strconv.FormatUint(uint64(b)&0x7, 10)
	}
	return strconv.FormatUint(uint64(b>>3), 10)
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

func (r Range) String() string {
	return fmt.Sprintf("%s-%s", Bits(r.Start), Bits(r.Stop))
}

func (r Range) Length() int64 {
	return r.Stop - r.Start
}

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

// TODO: base instead?

var lookupRe = regexp.MustCompile(`` +
	`^(?:` +
	`([\w_]+)` +
	`|` +
	`\[(\d+)\]` +
	`)` +
	`(?:\.?)`)

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

func (v Value) String() string {
	f := ""
	switch iv := v.V.(type) {
	case Array:
		f = "array"
	case Struct:
		f = "struct"
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
		f = strconv.FormatFloat(iv, 'f', -1, 64)
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
		bs, _ := v.BitBuf.BytesBitRange(0, 16*8, 0)
		if v.BitBuf.Len > 16 {
			f = hex.EncodeToString(bs) + "..."
		} else {
			f = hex.EncodeToString(bs)
		}
	case nil:
		f = "none"
		// TODO:
		//return hex.EncodeToString(v.Bytes)
	// case TypeDecoder:
	// 	c := v.Decoder
	// 	f = fmt.Sprintf("%s (%s) %s", c.Format().Name, c.MIME(), Bits(c.BitBuf().Len))
	// case TypeArray:
	// 	f = "array"
	default:
		panic("unreachable")
	}
	symbol := ""
	if v.Symbol != "" {
		symbol = fmt.Sprintf(" (%s)", v.Symbol)
	}
	desc := ""
	if v.Desc != "" {
		desc = fmt.Sprintf(" (%s)", v.Desc)
	}
	return fmt.Sprintf("%s%s%s", f, symbol, desc)
}

func (v Value) RawString() string {
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
