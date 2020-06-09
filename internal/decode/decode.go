package decode

import (
	"fmt"
	"strconv"
	"strings"
)

type Endian int

const (
	BigEndian Endian = iota
	LittleEndian
)

type Common struct {
	Current *Field
	BitPos  uint64
	Buf     []byte
	Endian  Endian
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

func (c *Common) u(nBits uint, endian Endian) uint64 {
	n := ReadBits(c.Buf, c.BitPos, nBits)
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}
	c.BitPos += uint64(nBits)
	return n
}

func (c *Common) U(nBits uint) uint64 { return c.u(nBits, c.Endian) }
func (c *Common) U1() uint64          { return c.u(1, c.Endian) }
func (c *Common) U2() uint64          { return c.u(2, c.Endian) }
func (c *Common) U3() uint64          { return c.u(3, c.Endian) }
func (c *Common) U4() uint64          { return c.u(4, c.Endian) }
func (c *Common) U5() uint64          { return c.u(5, c.Endian) }
func (c *Common) U6() uint64          { return c.u(6, c.Endian) }
func (c *Common) U7() uint64          { return c.u(7, c.Endian) }
func (c *Common) U8() uint64          { return c.u(8, c.Endian) }
func (c *Common) U9() uint64          { return c.u(9, c.Endian) }
func (c *Common) U10() uint64         { return c.u(10, c.Endian) }
func (c *Common) U11() uint64         { return c.u(11, c.Endian) }
func (c *Common) U12() uint64         { return c.u(12, c.Endian) }
func (c *Common) U13() uint64         { return c.u(13, c.Endian) }
func (c *Common) U14() uint64         { return c.u(14, c.Endian) }
func (c *Common) U15() uint64         { return c.u(15, c.Endian) }
func (c *Common) U16() uint64         { return c.u(16, c.Endian) }
func (c *Common) U24() uint64         { return c.u(24, c.Endian) }
func (c *Common) U32() uint64         { return c.u(32, c.Endian) }
func (c *Common) U64() uint64         { return c.u(64, c.Endian) }

func (c *Common) UBE(bits uint) uint64 { return c.u(bits, BigEndian) }
func (c *Common) U9BE() uint64         { return c.u(9, BigEndian) }
func (c *Common) U10BE() uint64        { return c.u(10, BigEndian) }
func (c *Common) U11BE() uint64        { return c.u(11, BigEndian) }
func (c *Common) U12BE() uint64        { return c.u(12, BigEndian) }
func (c *Common) U13BE() uint64        { return c.u(13, BigEndian) }
func (c *Common) U14BE() uint64        { return c.u(14, BigEndian) }
func (c *Common) U15BE() uint64        { return c.u(15, BigEndian) }
func (c *Common) U16BE() uint64        { return c.u(16, BigEndian) }
func (c *Common) U24BE() uint64        { return c.u(24, BigEndian) }
func (c *Common) U32BE() uint64        { return c.u(32, BigEndian) }
func (c *Common) U64BE() uint64        { return c.u(64, BigEndian) }

func (c *Common) ULE(bits uint) uint64 { return c.u(bits, LittleEndian) }
func (c *Common) U9LE() uint64         { return c.u(9, LittleEndian) }
func (c *Common) U10LE() uint64        { return c.u(10, LittleEndian) }
func (c *Common) U11LE() uint64        { return c.u(11, LittleEndian) }
func (c *Common) U12LE() uint64        { return c.u(12, LittleEndian) }
func (c *Common) U13LE() uint64        { return c.u(13, LittleEndian) }
func (c *Common) U14LE() uint64        { return c.u(14, LittleEndian) }
func (c *Common) U15LE() uint64        { return c.u(15, LittleEndian) }
func (c *Common) U16LE() uint64        { return c.u(16, LittleEndian) }
func (c *Common) U24LE() uint64        { return c.u(24, LittleEndian) }
func (c *Common) U32LE() uint64        { return c.u(32, LittleEndian) }
func (c *Common) U64LE() uint64        { return c.u(64, LittleEndian) }

func (c *Common) fieldU(bits uint, name string, endian Endian) uint64 {
	start := c.BitPos
	n := c.u(bits, endian)
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: start + uint64(bits)},
		Value: Value{Type: TypeUInt, UInt: n},
	})
	return n
}

func (c *Common) FieldU(nBits uint, name string) uint64 { return c.fieldU(nBits, name, c.Endian) }
func (c *Common) FieldU1(name string) uint64            { return c.fieldU(1, name, c.Endian) }
func (c *Common) FieldU2(name string) uint64            { return c.fieldU(2, name, c.Endian) }
func (c *Common) FieldU3(name string) uint64            { return c.fieldU(3, name, c.Endian) }
func (c *Common) FieldU4(name string) uint64            { return c.fieldU(4, name, c.Endian) }
func (c *Common) FieldU5(name string) uint64            { return c.fieldU(5, name, c.Endian) }
func (c *Common) FieldU6(name string) uint64            { return c.fieldU(6, name, c.Endian) }
func (c *Common) FieldU7(name string) uint64            { return c.fieldU(7, name, c.Endian) }
func (c *Common) FieldU8(name string) uint64            { return c.fieldU(8, name, c.Endian) }
func (c *Common) FieldU9(name string) uint64            { return c.fieldU(9, name, c.Endian) }
func (c *Common) FieldU10(name string) uint64           { return c.fieldU(10, name, c.Endian) }
func (c *Common) FieldU11(name string) uint64           { return c.fieldU(11, name, c.Endian) }
func (c *Common) FieldU12(name string) uint64           { return c.fieldU(12, name, c.Endian) }
func (c *Common) FieldU13(name string) uint64           { return c.fieldU(13, name, c.Endian) }
func (c *Common) FieldU14(name string) uint64           { return c.fieldU(14, name, c.Endian) }
func (c *Common) FieldU15(name string) uint64           { return c.fieldU(15, name, c.Endian) }
func (c *Common) FieldU16(name string) uint64           { return c.fieldU(16, name, c.Endian) }
func (c *Common) FieldU24(name string) uint64           { return c.fieldU(24, name, c.Endian) }
func (c *Common) FieldU32(name string) uint64           { return c.fieldU(32, name, c.Endian) }
func (c *Common) FieldU64(name string) uint64           { return c.fieldU(64, name, c.Endian) }

func (c *Common) FieldUBE(nBits uint, name string) uint64 { return c.fieldU(nBits, name, BigEndian) }
func (c *Common) FieldU9BE(name string) uint64            { return c.fieldU(9, name, BigEndian) }
func (c *Common) FieldU10BE(name string) uint64           { return c.fieldU(10, name, BigEndian) }
func (c *Common) FieldU11BE(name string) uint64           { return c.fieldU(11, name, BigEndian) }
func (c *Common) FieldU12BE(name string) uint64           { return c.fieldU(12, name, BigEndian) }
func (c *Common) FieldU13BE(name string) uint64           { return c.fieldU(13, name, BigEndian) }
func (c *Common) FieldU14BE(name string) uint64           { return c.fieldU(14, name, BigEndian) }
func (c *Common) FieldU15BE(name string) uint64           { return c.fieldU(15, name, BigEndian) }
func (c *Common) FieldU16BE(name string) uint64           { return c.fieldU(16, name, BigEndian) }
func (c *Common) FieldU24BE(name string) uint64           { return c.fieldU(24, name, BigEndian) }
func (c *Common) FieldU32BE(name string) uint64           { return c.fieldU(32, name, BigEndian) }
func (c *Common) FieldU64BE(name string) uint64           { return c.fieldU(64, name, BigEndian) }

func (c *Common) FieldULE(nBits uint, name string) uint64 { return c.fieldU(nBits, name, LittleEndian) }
func (c *Common) FieldU9LE(name string) uint64            { return c.fieldU(9, name, LittleEndian) }
func (c *Common) FieldU10LE(name string) uint64           { return c.fieldU(10, name, LittleEndian) }
func (c *Common) FieldU11LE(name string) uint64           { return c.fieldU(11, name, LittleEndian) }
func (c *Common) FieldU12LE(name string) uint64           { return c.fieldU(12, name, LittleEndian) }
func (c *Common) FieldU13LE(name string) uint64           { return c.fieldU(13, name, LittleEndian) }
func (c *Common) FieldU14LE(name string) uint64           { return c.fieldU(14, name, LittleEndian) }
func (c *Common) FieldU15LE(name string) uint64           { return c.fieldU(15, name, LittleEndian) }
func (c *Common) FieldU16LE(name string) uint64           { return c.fieldU(16, name, LittleEndian) }
func (c *Common) FieldU24LE(name string) uint64           { return c.fieldU(24, name, LittleEndian) }
func (c *Common) FieldU32LE(name string) uint64           { return c.fieldU(32, name, LittleEndian) }
func (c *Common) FieldU64LE(name string) uint64           { return c.fieldU(64, name, LittleEndian) }

func (c *Common) s(nBits uint, endian Endian) int64 {
	n := ReadBits(c.Buf, c.BitPos, nBits)
	if endian == LittleEndian {
		n = ReverseBytes(nBits, n)
	}
	var s int64
	if n&(1<<(nBits-1)) > 0 {
		s = -int64((^n & ((1 << nBits) - 1)) + 1)
	} else {
		s = int64(n)
	}
	c.BitPos += uint64(nBits)
	return s
}

func (c *Common) S(nBits uint) int64 { return c.s(nBits, c.Endian) }
func (c *Common) S1() int64          { return c.s(1, c.Endian) }
func (c *Common) S2() int64          { return c.s(2, c.Endian) }
func (c *Common) S3() int64          { return c.s(3, c.Endian) }
func (c *Common) S4() int64          { return c.s(4, c.Endian) }
func (c *Common) S5() int64          { return c.s(5, c.Endian) }
func (c *Common) S6() int64          { return c.s(6, c.Endian) }
func (c *Common) S7() int64          { return c.s(7, c.Endian) }
func (c *Common) S8() int64          { return c.s(8, c.Endian) }
func (c *Common) S9() int64          { return c.s(9, c.Endian) }
func (c *Common) S10() int64         { return c.s(10, c.Endian) }
func (c *Common) S11() int64         { return c.s(11, c.Endian) }
func (c *Common) S12() int64         { return c.s(12, c.Endian) }
func (c *Common) S13() int64         { return c.s(13, c.Endian) }
func (c *Common) S14() int64         { return c.s(14, c.Endian) }
func (c *Common) S15() int64         { return c.s(15, c.Endian) }
func (c *Common) S16() int64         { return c.s(16, c.Endian) }
func (c *Common) S24() int64         { return c.s(24, c.Endian) }
func (c *Common) S32() int64         { return c.s(32, c.Endian) }
func (c *Common) S64() int64         { return c.s(64, c.Endian) }

func (c *Common) SBE(nBits uint) int64 { return c.s(nBits, BigEndian) }
func (c *Common) S9BE() int64          { return c.s(9, BigEndian) }
func (c *Common) S10BE() int64         { return c.s(10, BigEndian) }
func (c *Common) S11BE() int64         { return c.s(11, BigEndian) }
func (c *Common) S12BE() int64         { return c.s(12, BigEndian) }
func (c *Common) S13BE() int64         { return c.s(13, BigEndian) }
func (c *Common) S14BE() int64         { return c.s(14, BigEndian) }
func (c *Common) S15BE() int64         { return c.s(15, BigEndian) }
func (c *Common) S16BE() int64         { return c.s(16, BigEndian) }
func (c *Common) S24BE() int64         { return c.s(24, BigEndian) }
func (c *Common) S32BE() int64         { return c.s(32, BigEndian) }
func (c *Common) S64BE() int64         { return c.s(64, BigEndian) }

func (c *Common) SLE(nBits uint) int64 { return c.s(nBits, LittleEndian) }
func (c *Common) S9LE() int64          { return c.s(9, LittleEndian) }
func (c *Common) S10LE() int64         { return c.s(10, LittleEndian) }
func (c *Common) S11LE() int64         { return c.s(11, LittleEndian) }
func (c *Common) S12LE() int64         { return c.s(12, LittleEndian) }
func (c *Common) S13LE() int64         { return c.s(13, LittleEndian) }
func (c *Common) S14LE() int64         { return c.s(14, LittleEndian) }
func (c *Common) S15LE() int64         { return c.s(15, LittleEndian) }
func (c *Common) S16LE() int64         { return c.s(16, LittleEndian) }
func (c *Common) S24LE() int64         { return c.s(24, LittleEndian) }
func (c *Common) S32LE() int64         { return c.s(32, LittleEndian) }
func (c *Common) S64LE() int64         { return c.s(64, LittleEndian) }

func (c *Common) fieldS(bits uint, name string, endian Endian) int64 {
	start := c.BitPos
	n := c.s(bits, endian)
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: start + uint64(bits)},
		Value: Value{Type: TypeSInt, SInt: n},
	})
	return n
}

func (c *Common) FieldS(nBits uint, name string) int64 { return c.fieldS(nBits, name, c.Endian) }
func (c *Common) FieldS1(name string) int64            { return c.fieldS(1, name, c.Endian) }
func (c *Common) FieldS2(name string) int64            { return c.fieldS(2, name, c.Endian) }
func (c *Common) FieldS3(name string) int64            { return c.fieldS(3, name, c.Endian) }
func (c *Common) FieldS4(name string) int64            { return c.fieldS(4, name, c.Endian) }
func (c *Common) FieldS5(name string) int64            { return c.fieldS(5, name, c.Endian) }
func (c *Common) FieldS6(name string) int64            { return c.fieldS(6, name, c.Endian) }
func (c *Common) FieldS7(name string) int64            { return c.fieldS(7, name, c.Endian) }
func (c *Common) FieldS8(name string) int64            { return c.fieldS(8, name, c.Endian) }
func (c *Common) FieldS9(name string) int64            { return c.fieldS(9, name, c.Endian) }
func (c *Common) FieldS10(name string) int64           { return c.fieldS(10, name, c.Endian) }
func (c *Common) FieldS11(name string) int64           { return c.fieldS(11, name, c.Endian) }
func (c *Common) FieldS12(name string) int64           { return c.fieldS(12, name, c.Endian) }
func (c *Common) FieldS13(name string) int64           { return c.fieldS(13, name, c.Endian) }
func (c *Common) FieldS14(name string) int64           { return c.fieldS(14, name, c.Endian) }
func (c *Common) FieldS15(name string) int64           { return c.fieldS(15, name, c.Endian) }
func (c *Common) FieldS16(name string) int64           { return c.fieldS(16, name, c.Endian) }
func (c *Common) FieldS24(name string) int64           { return c.fieldS(24, name, c.Endian) }
func (c *Common) FieldS32(name string) int64           { return c.fieldS(32, name, c.Endian) }
func (c *Common) FieldS64(name string) int64           { return c.fieldS(64, name, c.Endian) }

func (c *Common) FieldSBE(nBits uint, name string) int64 { return c.fieldS(nBits, name, BigEndian) }
func (c *Common) FieldS9BE(name string) int64            { return c.fieldS(9, name, BigEndian) }
func (c *Common) FieldS10BE(name string) int64           { return c.fieldS(10, name, BigEndian) }
func (c *Common) FieldS11BE(name string) int64           { return c.fieldS(11, name, BigEndian) }
func (c *Common) FieldS12BE(name string) int64           { return c.fieldS(12, name, BigEndian) }
func (c *Common) FieldS13BE(name string) int64           { return c.fieldS(13, name, BigEndian) }
func (c *Common) FieldS14BE(name string) int64           { return c.fieldS(14, name, BigEndian) }
func (c *Common) FieldS15BE(name string) int64           { return c.fieldS(15, name, BigEndian) }
func (c *Common) FieldS16BE(name string) int64           { return c.fieldS(16, name, BigEndian) }
func (c *Common) FieldS24BE(name string) int64           { return c.fieldS(24, name, BigEndian) }
func (c *Common) FieldS32BE(name string) int64           { return c.fieldS(32, name, BigEndian) }
func (c *Common) FieldS64BE(name string) int64           { return c.fieldS(64, name, BigEndian) }

func (c *Common) FieldSLE(nBits uint, name string) int64 { return c.fieldS(nBits, name, LittleEndian) }
func (c *Common) FieldS9LE(name string) int64            { return c.fieldS(9, name, LittleEndian) }
func (c *Common) FieldS10LE(name string) int64           { return c.fieldS(10, name, LittleEndian) }
func (c *Common) FieldS11LE(name string) int64           { return c.fieldS(11, name, LittleEndian) }
func (c *Common) FieldS12LE(name string) int64           { return c.fieldS(12, name, LittleEndian) }
func (c *Common) FieldS13LE(name string) int64           { return c.fieldS(13, name, LittleEndian) }
func (c *Common) FieldS14LE(name string) int64           { return c.fieldS(14, name, LittleEndian) }
func (c *Common) FieldS15LE(name string) int64           { return c.fieldS(15, name, LittleEndian) }
func (c *Common) FieldS16LE(name string) int64           { return c.fieldS(16, name, LittleEndian) }
func (c *Common) FieldS24LE(name string) int64           { return c.fieldS(24, name, LittleEndian) }
func (c *Common) FieldS32LE(name string) int64           { return c.fieldS(32, name, LittleEndian) }
func (c *Common) FieldS64LE(name string) int64           { return c.fieldS(64, name, LittleEndian) }

func (c *Common) Bytes(length uint) []byte {
	if c.BitPos%8 == 0 {
		bytePos := c.BitPos / 8
		bs := c.Buf[bytePos : bytePos+uint64(length)]
		c.BitPos += uint64(length * 8)
		return bs
	}

	var bs []byte
	for i := uint(0); i < length; i++ {
		bs = append(bs, byte(ReadBits(c.Buf, c.BitPos, 8)))
		c.BitPos += 8
	}
	return bs
}

func (c *Common) FieldBytes(length uint, name string) []byte {
	start := c.BitPos
	bs := c.Bytes(length)
	stop := c.BitPos
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: stop},
		Value: Value{Type: TypeBytes, Bytes: bs},
	})
	return bs
}

func (c *Common) UTF8(length uint) string {
	return string(c.Bytes(length))
}

func (c *Common) FieldUTF8(length uint, name string) string {
	start := c.BitPos
	s := c.UTF8(length)
	stop := c.BitPos
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: stop},
		Value: Value{Type: TypeStr, Str: s},
	})
	return s
}

func (c *Common) Unary(s uint) uint {
	var n uint
	for uint(c.U1()) == s {
		n++
	}
	return n
}

func (c *Common) FieldFn(name string, fn func() Value) Value {
	prev := c.Current

	f := &Field{Name: name}
	c.Current = f
	prev.Children = append(prev.Children, f)
	start := c.BitPos
	f.Range.Start = start
	v := fn()
	f.Range.Stop = c.BitPos
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

func (c *Common) FieldVerifyFn(name string, v uint64, fn func() uint64) uint64 {
	return c.FieldUFn(name, func() (uint64, Format, string) {
		n := fn()
		s := "Correct"
		if n != v {
			s = "Incorrect"
		}
		return n, FormatHex, s
	})
}

func (c *Common) EOF() bool {
	return c.BitPos>>3 >= uint64(len(c.Buf))
}

func (c *Common) ByteAlignBits() uint {
	return uint((8 - (c.BitPos & 0x7)) & 0x7)
}

func (c *Common) BytePos() uint64 {
	return c.BitPos & 0x7
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
