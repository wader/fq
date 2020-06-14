package decode

import (
	"fmt"
	"fq/internal/bitbuf"
	"strconv"
	"strings"
)

type Options struct {
	Probe bool
}

type Register struct {
	Name string
	MIME string
	New  func() Decoder
}

type Decoder interface {
	Decode(Options) bool
}

type Common struct {
	*bitbuf.Buffer
	Current *Field
}

type Type int

const (
	TypeNone Type = iota
	TypeSInt
	TypeUInt
	TypeStr
	TypeBytes
)

type Format int

const (
	FormatDecimal Format = iota
	FormatBinary
	FormatOctal
	FormatHex
)

type Value struct {
	Type Type

	SInt  int64
	UInt  uint64
	Str   string
	Bytes []byte

	Format  Format
	Display string
	Mime    string
}

func (v Value) String() string {
	f := ""
	switch v.Type {
	case TypeNone:
		f = ""
	case TypeSInt:
		f = strconv.FormatInt(v.SInt, 10)
	case TypeUInt:
		f = strconv.FormatUint(v.UInt, 10)
	case TypeStr:
		f = v.Str
	case TypeBytes:
		f = fmt.Sprintf("%d bytes", len(v.Bytes))
		// TODO:
		//return hex.EncodeToString(v.Bytes)
	default:
		panic("unreachable")
	}
	if v.Display != "" {
		return fmt.Sprintf("%s (%s)", v.Display, f)
	}
	return f
}

type Range struct {
	Start uint64
	Stop  uint64
}

func (r Range) String() string {
	return fmt.Sprintf("%d-%d", r.Start, r.Stop)
}

type Field struct {
	Name     string
	Range    Range
	Value    Value
	Children []*Field
}

func (c *Common) fieldU(nBits uint64, name string, endian bitbuf.Endian) uint64 {
	start := c.Pos
	n := c.UE(nBits, endian)
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: start + uint64(nBits)},
		Value: Value{Type: TypeUInt, UInt: n},
	})
	return n
}

func (c *Common) FieldBool(name string) bool { return c.Bool() }

func (c *Common) FieldU(name string, nBits uint64) uint64 {
	return c.fieldU(nBits, name, bitbuf.BigEndian)
}
func (c *Common) FieldU1(name string) uint64  { return c.fieldU(1, name, bitbuf.BigEndian) }
func (c *Common) FieldU2(name string) uint64  { return c.fieldU(2, name, bitbuf.BigEndian) }
func (c *Common) FieldU3(name string) uint64  { return c.fieldU(3, name, bitbuf.BigEndian) }
func (c *Common) FieldU4(name string) uint64  { return c.fieldU(4, name, bitbuf.BigEndian) }
func (c *Common) FieldU5(name string) uint64  { return c.fieldU(5, name, bitbuf.BigEndian) }
func (c *Common) FieldU6(name string) uint64  { return c.fieldU(6, name, bitbuf.BigEndian) }
func (c *Common) FieldU7(name string) uint64  { return c.fieldU(7, name, bitbuf.BigEndian) }
func (c *Common) FieldU8(name string) uint64  { return c.fieldU(8, name, bitbuf.BigEndian) }
func (c *Common) FieldU9(name string) uint64  { return c.fieldU(9, name, bitbuf.BigEndian) }
func (c *Common) FieldU10(name string) uint64 { return c.fieldU(10, name, bitbuf.BigEndian) }
func (c *Common) FieldU11(name string) uint64 { return c.fieldU(11, name, bitbuf.BigEndian) }
func (c *Common) FieldU12(name string) uint64 { return c.fieldU(12, name, bitbuf.BigEndian) }
func (c *Common) FieldU13(name string) uint64 { return c.fieldU(13, name, bitbuf.BigEndian) }
func (c *Common) FieldU14(name string) uint64 { return c.fieldU(14, name, bitbuf.BigEndian) }
func (c *Common) FieldU15(name string) uint64 { return c.fieldU(15, name, bitbuf.BigEndian) }
func (c *Common) FieldU16(name string) uint64 { return c.fieldU(16, name, bitbuf.BigEndian) }
func (c *Common) FieldU24(name string) uint64 { return c.fieldU(24, name, bitbuf.BigEndian) }
func (c *Common) FieldU32(name string) uint64 { return c.fieldU(32, name, bitbuf.BigEndian) }
func (c *Common) FieldU64(name string) uint64 { return c.fieldU(64, name, bitbuf.BigEndian) }

func (c *Common) FieldUBE(nBits uint64, name string) uint64 {
	return c.fieldU(nBits, name, bitbuf.BigEndian)
}
func (c *Common) FieldU9BE(name string) uint64  { return c.fieldU(9, name, bitbuf.BigEndian) }
func (c *Common) FieldU10BE(name string) uint64 { return c.fieldU(10, name, bitbuf.BigEndian) }
func (c *Common) FieldU11BE(name string) uint64 { return c.fieldU(11, name, bitbuf.BigEndian) }
func (c *Common) FieldU12BE(name string) uint64 { return c.fieldU(12, name, bitbuf.BigEndian) }
func (c *Common) FieldU13BE(name string) uint64 { return c.fieldU(13, name, bitbuf.BigEndian) }
func (c *Common) FieldU14BE(name string) uint64 { return c.fieldU(14, name, bitbuf.BigEndian) }
func (c *Common) FieldU15BE(name string) uint64 { return c.fieldU(15, name, bitbuf.BigEndian) }
func (c *Common) FieldU16BE(name string) uint64 { return c.fieldU(16, name, bitbuf.BigEndian) }
func (c *Common) FieldU24BE(name string) uint64 { return c.fieldU(24, name, bitbuf.BigEndian) }
func (c *Common) FieldU32BE(name string) uint64 { return c.fieldU(32, name, bitbuf.BigEndian) }
func (c *Common) FieldU64BE(name string) uint64 { return c.fieldU(64, name, bitbuf.BigEndian) }

func (c *Common) FieldULE(nBits uint64, name string) uint64 {
	return c.fieldU(nBits, name, bitbuf.LittleEndian)
}
func (c *Common) FieldU9LE(name string) uint64  { return c.fieldU(9, name, bitbuf.LittleEndian) }
func (c *Common) FieldU10LE(name string) uint64 { return c.fieldU(10, name, bitbuf.LittleEndian) }
func (c *Common) FieldU11LE(name string) uint64 { return c.fieldU(11, name, bitbuf.LittleEndian) }
func (c *Common) FieldU12LE(name string) uint64 { return c.fieldU(12, name, bitbuf.LittleEndian) }
func (c *Common) FieldU13LE(name string) uint64 { return c.fieldU(13, name, bitbuf.LittleEndian) }
func (c *Common) FieldU14LE(name string) uint64 { return c.fieldU(14, name, bitbuf.LittleEndian) }
func (c *Common) FieldU15LE(name string) uint64 { return c.fieldU(15, name, bitbuf.LittleEndian) }
func (c *Common) FieldU16LE(name string) uint64 { return c.fieldU(16, name, bitbuf.LittleEndian) }
func (c *Common) FieldU24LE(name string) uint64 { return c.fieldU(24, name, bitbuf.LittleEndian) }
func (c *Common) FieldU32LE(name string) uint64 { return c.fieldU(32, name, bitbuf.LittleEndian) }
func (c *Common) FieldU64LE(name string) uint64 { return c.fieldU(64, name, bitbuf.LittleEndian) }

func (c *Common) fieldS(name string, nBits uint64, endian bitbuf.Endian) int64 {
	start := c.Pos
	n := c.SE(nBits, endian)
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: start + uint64(nBits)},
		Value: Value{Type: TypeSInt, SInt: n},
	})
	return n
}

func (c *Common) FieldS(name string, nBits uint64) int64 {
	return c.fieldS(name, nBits, bitbuf.BigEndian)
}
func (c *Common) FieldS1(name string) int64  { return c.fieldS(name, 1, bitbuf.BigEndian) }
func (c *Common) FieldS2(name string) int64  { return c.fieldS(name, 2, bitbuf.BigEndian) }
func (c *Common) FieldS3(name string) int64  { return c.fieldS(name, 3, bitbuf.BigEndian) }
func (c *Common) FieldS4(name string) int64  { return c.fieldS(name, 4, bitbuf.BigEndian) }
func (c *Common) FieldS5(name string) int64  { return c.fieldS(name, 5, bitbuf.BigEndian) }
func (c *Common) FieldS6(name string) int64  { return c.fieldS(name, 6, bitbuf.BigEndian) }
func (c *Common) FieldS7(name string) int64  { return c.fieldS(name, 7, bitbuf.BigEndian) }
func (c *Common) FieldS8(name string) int64  { return c.fieldS(name, 8, bitbuf.BigEndian) }
func (c *Common) FieldS9(name string) int64  { return c.fieldS(name, 9, bitbuf.BigEndian) }
func (c *Common) FieldS10(name string) int64 { return c.fieldS(name, 10, bitbuf.BigEndian) }
func (c *Common) FieldS11(name string) int64 { return c.fieldS(name, 11, bitbuf.BigEndian) }
func (c *Common) FieldS12(name string) int64 { return c.fieldS(name, 12, bitbuf.BigEndian) }
func (c *Common) FieldS13(name string) int64 { return c.fieldS(name, 13, bitbuf.BigEndian) }
func (c *Common) FieldS14(name string) int64 { return c.fieldS(name, 14, bitbuf.BigEndian) }
func (c *Common) FieldS15(name string) int64 { return c.fieldS(name, 15, bitbuf.BigEndian) }
func (c *Common) FieldS16(name string) int64 { return c.fieldS(name, 16, bitbuf.BigEndian) }
func (c *Common) FieldS24(name string) int64 { return c.fieldS(name, 24, bitbuf.BigEndian) }
func (c *Common) FieldS32(name string) int64 { return c.fieldS(name, 32, bitbuf.BigEndian) }
func (c *Common) FieldS64(name string) int64 { return c.fieldS(name, 64, bitbuf.BigEndian) }

func (c *Common) FieldSBE(name string, nBits uint64) int64 {
	return c.fieldS(name, nBits, bitbuf.BigEndian)
}
func (c *Common) FieldS9BE(name string) int64  { return c.fieldS(name, 9, bitbuf.BigEndian) }
func (c *Common) FieldS10BE(name string) int64 { return c.fieldS(name, 10, bitbuf.BigEndian) }
func (c *Common) FieldS11BE(name string) int64 { return c.fieldS(name, 11, bitbuf.BigEndian) }
func (c *Common) FieldS12BE(name string) int64 { return c.fieldS(name, 12, bitbuf.BigEndian) }
func (c *Common) FieldS13BE(name string) int64 { return c.fieldS(name, 13, bitbuf.BigEndian) }
func (c *Common) FieldS14BE(name string) int64 { return c.fieldS(name, 14, bitbuf.BigEndian) }
func (c *Common) FieldS15BE(name string) int64 { return c.fieldS(name, 15, bitbuf.BigEndian) }
func (c *Common) FieldS16BE(name string) int64 { return c.fieldS(name, 16, bitbuf.BigEndian) }
func (c *Common) FieldS24BE(name string) int64 { return c.fieldS(name, 24, bitbuf.BigEndian) }
func (c *Common) FieldS32BE(name string) int64 { return c.fieldS(name, 32, bitbuf.BigEndian) }
func (c *Common) FieldS64BE(name string) int64 { return c.fieldS(name, 64, bitbuf.BigEndian) }

func (c *Common) FieldSLE(nBits uint64, name string) int64 {
	return c.fieldS(name, nBits, bitbuf.LittleEndian)
}
func (c *Common) FieldS9LE(name string) int64  { return c.fieldS(name, 9, bitbuf.LittleEndian) }
func (c *Common) FieldS10LE(name string) int64 { return c.fieldS(name, 10, bitbuf.LittleEndian) }
func (c *Common) FieldS11LE(name string) int64 { return c.fieldS(name, 11, bitbuf.LittleEndian) }
func (c *Common) FieldS12LE(name string) int64 { return c.fieldS(name, 12, bitbuf.LittleEndian) }
func (c *Common) FieldS13LE(name string) int64 { return c.fieldS(name, 13, bitbuf.LittleEndian) }
func (c *Common) FieldS14LE(name string) int64 { return c.fieldS(name, 14, bitbuf.LittleEndian) }
func (c *Common) FieldS15LE(name string) int64 { return c.fieldS(name, 15, bitbuf.LittleEndian) }
func (c *Common) FieldS16LE(name string) int64 { return c.fieldS(name, 16, bitbuf.LittleEndian) }
func (c *Common) FieldS24LE(name string) int64 { return c.fieldS(name, 24, bitbuf.LittleEndian) }
func (c *Common) FieldS32LE(name string) int64 { return c.fieldS(name, 32, bitbuf.LittleEndian) }
func (c *Common) FieldS64LE(name string) int64 { return c.fieldS(name, 64, bitbuf.LittleEndian) }

func (c *Common) FieldFn(name string, fn func() Value) Value {
	prev := c.Current

	f := &Field{Name: name}
	c.Current = f
	prev.Children = append(prev.Children, f)
	start := c.Pos
	f.Range.Start = start
	v := fn()
	f.Range.Stop = c.Pos
	f.Value = v

	c.Current = prev

	return v
}

func (c *Common) FieldNoneFn(name string, fn func()) {
	c.FieldFn(name, func() Value {
		fn()
		return Value{}
	})
}

func (c *Common) FieldUFn(name string, fn func() (uint64, Format, string)) uint64 {
	return c.FieldFn(name, func() Value {
		u, fmt, d := fn()
		return Value{Type: TypeUInt, UInt: u, Format: fmt, Display: d}
	}).UInt
}

func (c *Common) FieldSFn(name string, fn func() (int64, Format, string)) int64 {
	return c.FieldFn(name, func() Value {
		s, fmt, d := fn()
		return Value{Type: TypeSInt, SInt: s, Format: fmt, Display: d}
	}).SInt
}

func (c *Common) FieldStrFn(name string, fn func() (string, string)) string {
	return c.FieldFn(name, func() Value {
		str, disp := fn()
		return Value{Type: TypeStr, Str: str, Display: disp}
	}).Str
}

func (c *Common) FieldBytesFn(name string, fn func() ([]byte, string)) []byte {
	return c.FieldFn(name, func() Value {
		bs, disp := fn()
		return Value{Type: TypeBytes, Bytes: bs, Display: disp}
	}).Bytes
}

func (c *Common) FieldStringMapFn(name string, sm map[uint64]string, def string, fn func() uint64) uint64 {
	return c.FieldUFn(name, func() (uint64, Format, string) {
		n := fn()
		var d string
		d, ok := sm[n]
		if !ok {
			d = def
		}
		return n, FormatDecimal, d
	})
}

func (c *Common) FieldVerifyUFn(name string, v uint64, fn func() uint64) bool {
	n := c.FieldUFn(name, func() (uint64, Format, string) {
		n := fn()
		s := "Correct"
		if n != v {
			s = "Incorrect"
		}
		return n, FormatHex, s
	})
	return n == v
}

// TODO: FieldBytesRange or?
func (c *Common) FieldBytes(name string, nBytes uint64) []byte {
	return c.FieldBytesFn(name, func() ([]byte, string) {
		bs, _ := c.BytesLen(nBytes)
		return bs, ""
	})
}

func (c *Common) FieldUTF8(name string, nBytes uint64) string {
	return c.FieldStrFn(name, func() (string, string) {
		str, _ := c.UTF8(nBytes)
		return str, ""
	})
}

func (c *Common) FieldVerifyStringFn(name string, v string, fn func() string) bool {
	return c.FieldStrFn(name, func() (string, string) {
		str := fn()
		s := "Correct"
		if str != v {
			s = "Incorrect"
		}
		return str, s
	}) == v
}

// --------------

func Dump(f *Field, depth int) {
	indent := strings.Repeat("  ", depth)
	if (len(f.Children)) != 0 {
		fmt.Printf("%s%s: %s %s {\n", indent, f.Name, f.Range, f.Value)
		for _, c := range f.Children {
			Dump(c, depth+1)
		}
		fmt.Printf("%s}\n", indent)
	} else {
		fmt.Printf("%s%s: %s %s\n", indent, f.Name, f.Range, f.Value)
	}
}
