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
	TypeInt
	TypeUint
	TypeStr
	TypeBytes
)

type Value struct {
	Type  Type
	Int   int64
	Uint  uint64
	Str   string
	Bytes []byte
	Mime  string
}

func (v Value) String() string {
	switch v.Type {
	case TypeNone:
		return "None"
	case TypeInt:
		return strconv.FormatInt(v.Int, 10)
	case TypeUint:
		return strconv.FormatUint(v.Uint, 10)
	case TypeStr:
		return v.Str
	case TypeBytes:
		return fmt.Sprintf("%d bytes", len(v.Bytes))
		// TODO:
		//return hex.EncodeToString(v.Bytes)
	default:
		panic("unreachable")
	}
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
	Display  string
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

func (c *Common) U(bits uint) uint64 { return c.u(bits, c.Endian) }
func (c *Common) U1() uint64         { return c.u(1, c.Endian) }
func (c *Common) U2() uint64         { return c.u(2, c.Endian) }
func (c *Common) U3() uint64         { return c.u(3, c.Endian) }
func (c *Common) U4() uint64         { return c.u(4, c.Endian) }
func (c *Common) U5() uint64         { return c.u(5, c.Endian) }
func (c *Common) U6() uint64         { return c.u(6, c.Endian) }
func (c *Common) U7() uint64         { return c.u(7, c.Endian) }
func (c *Common) U8() uint64         { return c.u(8, c.Endian) }
func (c *Common) U9() uint64         { return c.u(9, c.Endian) }
func (c *Common) U10() uint64        { return c.u(10, c.Endian) }
func (c *Common) U11() uint64        { return c.u(11, c.Endian) }
func (c *Common) U12() uint64        { return c.u(12, c.Endian) }
func (c *Common) U13() uint64        { return c.u(13, c.Endian) }
func (c *Common) U14() uint64        { return c.u(14, c.Endian) }
func (c *Common) U15() uint64        { return c.u(15, c.Endian) }
func (c *Common) U16() uint64        { return c.u(16, c.Endian) }
func (c *Common) U24() uint64        { return c.u(24, c.Endian) }
func (c *Common) U32() uint64        { return c.u(32, c.Endian) }
func (c *Common) U64() uint64        { return c.u(64, c.Endian) }

func (c *Common) ULE(bits uint) uint64 { return c.u(bits, LittleEndian) }
func (c *Common) U1LE() uint64         { return c.u(1, LittleEndian) }
func (c *Common) U2LE() uint64         { return c.u(2, LittleEndian) }
func (c *Common) U3LE() uint64         { return c.u(3, LittleEndian) }
func (c *Common) U4LE() uint64         { return c.u(4, LittleEndian) }
func (c *Common) U5LE() uint64         { return c.u(5, LittleEndian) }
func (c *Common) U6LE() uint64         { return c.u(6, LittleEndian) }
func (c *Common) U7LE() uint64         { return c.u(7, LittleEndian) }
func (c *Common) U8LE() uint64         { return c.u(8, LittleEndian) }
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

func (c *Common) UBE(bits uint) uint64 { return c.u(bits, BigEndian) }
func (c *Common) U1BE() uint64         { return c.u(1, BigEndian) }
func (c *Common) U2BE() uint64         { return c.u(2, BigEndian) }
func (c *Common) U3BE() uint64         { return c.u(3, BigEndian) }
func (c *Common) U4BE() uint64         { return c.u(4, BigEndian) }
func (c *Common) U5BE() uint64         { return c.u(5, BigEndian) }
func (c *Common) U6BE() uint64         { return c.u(6, BigEndian) }
func (c *Common) U7BE() uint64         { return c.u(7, BigEndian) }
func (c *Common) U8BE() uint64         { return c.u(8, BigEndian) }
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

func (c *Common) SLE(nBits uint) int64 { return c.s(nBits, LittleEndian) }
func (c *Common) S1LE() int64          { return c.s(1, LittleEndian) }
func (c *Common) S2LE() int64          { return c.s(2, LittleEndian) }
func (c *Common) S3LE() int64          { return c.s(3, LittleEndian) }
func (c *Common) S4LE() int64          { return c.s(4, LittleEndian) }
func (c *Common) S5LE() int64          { return c.s(5, LittleEndian) }
func (c *Common) S6LE() int64          { return c.s(6, LittleEndian) }
func (c *Common) S7LE() int64          { return c.s(7, LittleEndian) }
func (c *Common) S8LE() int64          { return c.s(8, LittleEndian) }
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

func (c *Common) SBE(nBits uint) int64 { return c.s(nBits, BigEndian) }
func (c *Common) S1BE() int64          { return c.s(1, BigEndian) }
func (c *Common) S2BE() int64          { return c.s(2, BigEndian) }
func (c *Common) S3BE() int64          { return c.s(3, BigEndian) }
func (c *Common) S4BE() int64          { return c.s(4, BigEndian) }
func (c *Common) S5BE() int64          { return c.s(5, BigEndian) }
func (c *Common) S6BE() int64          { return c.s(6, BigEndian) }
func (c *Common) S7BE() int64          { return c.s(7, BigEndian) }
func (c *Common) S8BE() int64          { return c.s(8, BigEndian) }
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

func (c *Common) Bytes(length uint) Value {
	if c.BitPos%8 == 0 {
		bytePos := c.BitPos / 8
		bs := c.Buf[bytePos : bytePos+uint64(length)]
		c.BitPos += uint64(length * 8)
		return Value{Type: TypeBytes, Bytes: bs}
	} else {
		var bs []byte
		for i := uint(0); i < length; i++ {
			bs = append(bs, byte(ReadBits(c.Buf, c.BitPos, 8)))
			c.BitPos += 8
		}
		return Value{Type: TypeBytes, Bytes: bs}
	}
}

func (c *Common) UTF8(length uint) Value {
	v := c.Bytes(length)
	return Value{Type: TypeStr, Str: string(v.Bytes)}
}

func (c *Common) Field(name string, fn func() (Value, string)) (Value, string) {
	prev := c.Current

	f := &Field{Name: name}
	c.Current = f
	prev.Children = append(prev.Children, f)
	start := c.BitPos
	f.Range.Start = start
	v, d := fn()
	f.Range.Stop = c.BitPos
	f.Value = v
	f.Display = d

	c.Current = prev

	return v, d
}

func (c *Common) FieldU(bits uint, name string) uint64 {
	start := c.BitPos
	n := c.U(bits)
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: start + uint64(bits)},
		Value: Value{Type: TypeUint, Uint: n},
	})
	return n
}

func (c *Common) FieldU1(name string) uint64  { return c.FieldU(1, name) }
func (c *Common) FieldU2(name string) uint64  { return c.FieldU(2, name) }
func (c *Common) FieldU3(name string) uint64  { return c.FieldU(3, name) }
func (c *Common) FieldU4(name string) uint64  { return c.FieldU(4, name) }
func (c *Common) FieldU5(name string) uint64  { return c.FieldU(5, name) }
func (c *Common) FieldU6(name string) uint64  { return c.FieldU(6, name) }
func (c *Common) FieldU7(name string) uint64  { return c.FieldU(7, name) }
func (c *Common) FieldU8(name string) uint64  { return c.FieldU(8, name) }
func (c *Common) FieldU9(name string) uint64  { return c.FieldU(9, name) }
func (c *Common) FieldU10(name string) uint64 { return c.FieldU(10, name) }
func (c *Common) FieldU11(name string) uint64 { return c.FieldU(11, name) }
func (c *Common) FieldU12(name string) uint64 { return c.FieldU(12, name) }
func (c *Common) FieldU13(name string) uint64 { return c.FieldU(13, name) }
func (c *Common) FieldU14(name string) uint64 { return c.FieldU(14, name) }
func (c *Common) FieldU15(name string) uint64 { return c.FieldU(15, name) }
func (c *Common) FieldU16(name string) uint64 { return c.FieldU(16, name) }
func (c *Common) FieldU24(name string) uint64 { return c.FieldU(24, name) }
func (c *Common) FieldU32(name string) uint64 { return c.FieldU(32, name) }
func (c *Common) FieldU64(name string) uint64 { return c.FieldU(64, name) }

func (c *Common) FieldS(bits uint, name string) int64 {
	start := c.BitPos
	n := c.S(bits)
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: start + uint64(bits)},
		Value: Value{Type: TypeInt, Int: n},
	})
	return n
}

func (c *Common) FieldS1(name string) int64  { return c.FieldS(1, name) }
func (c *Common) FieldS2(name string) int64  { return c.FieldS(2, name) }
func (c *Common) FieldS3(name string) int64  { return c.FieldS(3, name) }
func (c *Common) FieldS4(name string) int64  { return c.FieldS(4, name) }
func (c *Common) FieldS5(name string) int64  { return c.FieldS(5, name) }
func (c *Common) FieldS6(name string) int64  { return c.FieldS(6, name) }
func (c *Common) FieldS7(name string) int64  { return c.FieldS(7, name) }
func (c *Common) FieldS8(name string) int64  { return c.FieldS(8, name) }
func (c *Common) FieldS9(name string) int64  { return c.FieldS(9, name) }
func (c *Common) FieldS10(name string) int64 { return c.FieldS(10, name) }
func (c *Common) FieldS11(name string) int64 { return c.FieldS(11, name) }
func (c *Common) FieldS12(name string) int64 { return c.FieldS(12, name) }
func (c *Common) FieldS13(name string) int64 { return c.FieldS(13, name) }
func (c *Common) FieldS14(name string) int64 { return c.FieldS(14, name) }
func (c *Common) FieldS15(name string) int64 { return c.FieldS(15, name) }
func (c *Common) FieldS16(name string) int64 { return c.FieldS(16, name) }
func (c *Common) FieldS24(name string) int64 { return c.FieldS(24, name) }
func (c *Common) FieldS32(name string) int64 { return c.FieldS(32, name) }
func (c *Common) FieldS64(name string) int64 { return c.FieldS(64, name) }

func (c *Common) FieldBytes(length uint, name string) Value {
	start := c.BitPos
	v := c.Bytes(length)
	stop := c.BitPos
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: stop},
		Value: v,
	})
	return v
}

func (c *Common) FieldUTF8(length uint, name string) Value {
	start := c.BitPos
	v := c.UTF8(length)
	stop := c.BitPos
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: stop},
		Value: v,
	})
	return v
}

func (c *Common) EOF() bool {
	return c.BitPos>>3 >= uint64(len(c.Buf))
}

func (c *Common) ByteAlignBits() uint {
	return uint((8 - (c.BitPos & 0x7)) & 0x7)
}

// --------------

func Dump(f *Field, depth int) {
	indent := strings.Repeat("  ", depth)
	if (len(f.Children)) != 0 {
		fmt.Printf("%s%s: %s %s %s {\n", indent, f.Name, f.Range, f.Value, f.Display)
		for _, c := range f.Children {
			Dump(c, depth+1)
		}
		fmt.Printf("%s}\n", indent)
	} else {
		fmt.Printf("%s%s: %s %s %s\n", indent, f.Name, f.Range, f.Value, f.Display)
	}
}
