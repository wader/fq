package decode

import (
	"encoding/hex"
	"fmt"
	"fq/pkg/bitbuf"
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

// TODO: base instead?

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
