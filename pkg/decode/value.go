package decode

import (
	"encoding/hex"
	"fmt"
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

func (r Range) String() string {
	return fmt.Sprintf("%s-%s", Bits(r.Start), Bits(r.Stop))
}

func (r Range) Length() int64 {
	return r.Stop - r.Start
}

type Type int

const (
	TypeNone Type = iota
	TypeBool
	TypeSInt
	TypeUInt
	TypeFloat
	TypeStr
	TypeBytes
	TypePadding
	TypeDecoder
)

// TODO: base instead?
type NumberFormat int

const (
	NumberDecimal = 10
	NumberBinary  = 2
	NumberOctal   = 8
	NumberHex     = 16
)

type Value struct {
	Type Type

	Bool    bool
	SInt    int64
	UInt    uint64
	Float   float64
	Str     string
	Bytes   []byte
	Decoder Decoder

	Format  NumberFormat
	Display string
	Mime    string
}

func (v Value) String() string {
	f := ""
	switch v.Type {
	case TypeNone:
		f = "none"
	case TypeBool:
		f = "false"
		if v.Bool {
			f = "true"
		}
	case TypeSInt:
		f = strconv.FormatInt(v.SInt, int(v.Format))
	case TypeUInt:
		f = strconv.FormatUint(v.UInt, int(v.Format))
	case TypeFloat:
		f = strconv.FormatFloat(v.Float, 'f', -1, 64)
	case TypeStr:
		f = v.Str
		if len(f) > 50 {
			f = fmt.Sprintf("%q", f[0:50]) + "..."
		} else {
			f = fmt.Sprintf("%q", v.Str)
		}
	case TypeBytes:
		if len(v.Bytes) > 50 {
			f = hex.EncodeToString(v.Bytes[0:25]) + "..."

		} else {
			f = hex.EncodeToString(v.Bytes)
		}
	case TypePadding:
		f = "padding"
		// TODO:
		//return hex.EncodeToString(v.Bytes)
	case TypeDecoder:
		c := v.Decoder
		f = fmt.Sprintf("%s (%s) %s", c.Format().Name, c.MIME(), Bits(c.BitBuf().Len))
	default:
		panic("unreachable")
	}
	if v.Display != "" {
		return fmt.Sprintf("%s (%s)", v.Display, f)
	}
	return f
}

func (v Value) RawString() string {
	f := ""
	switch v.Type {
	case TypeNone:
		f = "none"
	case TypeBool:
		f = "false"
		if v.Bool {
			f = "true"
		}
	case TypeSInt:
		f = strconv.FormatInt(v.SInt, 10)
	case TypeUInt:
		f = strconv.FormatUint(v.UInt, 10)
	case TypeFloat:
		f = strconv.FormatFloat(v.Float, 'f', -1, 64)
	case TypeStr:
		f = v.Str
	case TypeBytes:
		f = string(v.Bytes)
	case TypePadding:
		f = string(v.Bytes)
	case TypeDecoder:
		c := v.Decoder
		b, _ := c.BitBuf().BytesRange(0, c.BitBuf().Len/8)
		f = string(b)
	default:
		panic("unreachable")
	}
	return f
}
