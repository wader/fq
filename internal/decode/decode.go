package decode

import (
	"fmt"
	"fq/internal/bitbuf"
	"log"
	"math"
	"strconv"
	"strings"
)

type Error struct {
	Err   error
	Op    string
	Size  uint64
	Delta int64
	Pos   uint64
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: failed at bit position %d (size %d, delta %d): %s", e.Op, e.Pos, e.Size, e.Delta, e.Err)
}
func (e Error) Unwrap() error { return e.Err }

type Options struct {
	Probe bool
}

type Register struct {
	Name string
	MIME string
	New  func(common Common) Decoder
}

type Decoder interface {
	Decode(Options) bool
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

type Common struct {
	Current   *Field // TODO: need root field also?
	Parent    *Common
	Registers []*Register

	bitBuf *bitbuf.Buffer
}

func probe(bb *bitbuf.Buffer, registers []*Register, decoderNames []string) (*Register, Common, bool) {
	// TODO: order..
	var namesMap = map[string]bool{}
	for _, s := range decoderNames {
		namesMap[s] = true
	}

	for _, r := range registers {
		if decoderNames != nil {
			if _, ok := namesMap[r.Name]; !ok {
				continue
			}
		}

		// TODO: how to pass regsiters? do later? current field?
		c := Common{
			Current:   &Field{Name: r.Name},
			bitBuf:    bb.Copy(),
			Registers: registers,
		}
		d := r.New(c)
		if d.Decode(Options{}) {
			return r, c, true
		}
	}
	return nil, Common{}, false
}

func New(parent *Common, bb *bitbuf.Buffer, registers []*Register, decoderNames []string) (*Register, Common) {
	// TODO: add common,register to Decoder interface? rename register?
	r, c, ok := probe(bb, registers, decoderNames)
	if !ok {
		return nil, Common{}
	}

	return r, c
}

func (c *Common) fieldU(nBits uint64, name string, endian bitbuf.Endian) uint64 {
	start := c.bitBuf.Pos
	n, err := c.bitBuf.UE(nBits, endian)
	if err != nil {
		panic(Error{Err: err, Op: "FieldU" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
	}
	c.Current.Children = append(c.Current.Children, &Field{
		Name:  name,
		Range: Range{Start: start, Stop: start + uint64(nBits)},
		Value: Value{Type: TypeUInt, UInt: n},
	})
	return n
}

// TODO: return decooder?
func (c *Common) FieldDecode(name string, nBits uint64, decoderNames []string) bool {

	//start := c.Pos
	bb, err := c.bitBuf.BitBufRange(c.bitBuf.Pos, nBits)
	if err != nil {
		panic(Error{Err: err, Op: "FieldDecode", Size: nBits, Pos: c.bitBuf.Pos})
	}

	r, fieldC := New(c, bb, c.Registers, decoderNames)

	if r == nil {
		log.Printf("FieldDecode nope %#+v\n", decoderNames)
		return false
	}
	log.Printf("FieldDecode r: %#+v\n", r)

	// TODO: translate positions?
	// TODO: what out muxed stream?

	c.Current.Children = append(c.Current.Children, fieldC.Current)

	c.bitBuf.SeekRel(int64(fieldC.bitBuf.Pos))

	// TODO: what to return?
	return true
}

// TODO: return decooder?
func (c *Common) FieldDecodeRange(name string, start uint64, nBits uint64, decoderNames []string) bool {

	//start := c.Pos

	bb, err := c.bitBuf.BitBufRange(start, nBits)
	if err != nil {
		panic(Error{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: c.bitBuf.Pos})
	}

	r, fieldC := New(c, bb, c.Registers, decoderNames)

	// log.Printf("bb: %#+v\n", bb)

	if r == nil {
		log.Printf("FieldDecodeRange nope %#+v\n", decoderNames)
		return false
	}
	log.Printf("FieldDecodeRange r: %#+v\n", r)

	// TODO: translate positions?
	// TODO: what out muxed stream?

	c.Current.Children = append(c.Current.Children, fieldC.Current)

	// TODO: what to return?
	return true
}

func (c *Common) PeekBits(nBits uint64) uint64 {
	n, err := c.bitBuf.PeekBits(nBits)
	if err != nil {
		panic(Error{Err: err, Op: "PeekBits", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) BytesRange(firstBit uint64, nBytes uint64) []byte {
	bs, err := c.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(Error{Err: err, Op: "BytesRange", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *Common) BytesLen(nBytes uint64) []byte {
	bs, err := c.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(Error{Err: err, Op: "BytesLen", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *Common) Pos() uint64           { return c.bitBuf.Pos }
func (c *Common) Len() uint64           { return c.bitBuf.Len }
func (c *Common) End() bool             { return c.bitBuf.End() }
func (c *Common) BitsLeft() uint64      { return c.bitBuf.BitsLeft() }
func (c *Common) ByteAlignBits() uint64 { return c.bitBuf.ByteAlignBits() }
func (c *Common) BytePos() uint64       { return c.bitBuf.BytePos() }

func (c *Common) SeekRel(delta int64) uint64 {
	pos, err := c.bitBuf.SeekRel(delta)
	if err != nil {
		panic(Error{Err: err, Op: "SeekRel", Delta: delta, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *Common) SeekAbs(pos uint64) uint64 {
	pos, err := c.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(Error{Err: err, Op: "SeekAbs", Size: pos, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *Common) UE(nBits uint64, endian bitbuf.Endian) uint64 {
	n, err := c.bitBuf.UE(nBits, endian)
	if err != nil {
		panic(Error{Err: err, Op: "UE", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) Bool() bool {
	b, err := c.bitBuf.Bool()
	if err != nil {
		panic(Error{Err: err, Op: "Bool", Size: 1, Pos: c.bitBuf.Pos})
	}
	return b
}

func (c *Common) U(nBits uint64) uint64 { return c.UE(nBits, bitbuf.BigEndian) }
func (c *Common) U1() uint64            { return c.UE(1, bitbuf.BigEndian) }
func (c *Common) U2() uint64            { return c.UE(2, bitbuf.BigEndian) }
func (c *Common) U3() uint64            { return c.UE(3, bitbuf.BigEndian) }
func (c *Common) U4() uint64            { return c.UE(4, bitbuf.BigEndian) }
func (c *Common) U5() uint64            { return c.UE(5, bitbuf.BigEndian) }
func (c *Common) U6() uint64            { return c.UE(6, bitbuf.BigEndian) }
func (c *Common) U7() uint64            { return c.UE(7, bitbuf.BigEndian) }
func (c *Common) U8() uint64            { return c.UE(8, bitbuf.BigEndian) }
func (c *Common) U9() uint64            { return c.UE(9, bitbuf.BigEndian) }
func (c *Common) U10() uint64           { return c.UE(10, bitbuf.BigEndian) }
func (c *Common) U11() uint64           { return c.UE(11, bitbuf.BigEndian) }
func (c *Common) U12() uint64           { return c.UE(12, bitbuf.BigEndian) }
func (c *Common) U13() uint64           { return c.UE(13, bitbuf.BigEndian) }
func (c *Common) U14() uint64           { return c.UE(14, bitbuf.BigEndian) }
func (c *Common) U15() uint64           { return c.UE(15, bitbuf.BigEndian) }
func (c *Common) U16() uint64           { return c.UE(16, bitbuf.BigEndian) }
func (c *Common) U24() uint64           { return c.UE(24, bitbuf.BigEndian) }
func (c *Common) U32() uint64           { return c.UE(32, bitbuf.BigEndian) }
func (c *Common) U64() uint64           { return c.UE(64, bitbuf.BigEndian) }

func (c *Common) UBE(nBits uint64) uint64 { return c.UE(nBits, bitbuf.BigEndian) }
func (c *Common) U9BE() uint64            { return c.UE(9, bitbuf.BigEndian) }
func (c *Common) U10BE() uint64           { return c.UE(10, bitbuf.BigEndian) }
func (c *Common) U11BE() uint64           { return c.UE(11, bitbuf.BigEndian) }
func (c *Common) U12BE() uint64           { return c.UE(12, bitbuf.BigEndian) }
func (c *Common) U13BE() uint64           { return c.UE(13, bitbuf.BigEndian) }
func (c *Common) U14BE() uint64           { return c.UE(14, bitbuf.BigEndian) }
func (c *Common) U15BE() uint64           { return c.UE(15, bitbuf.BigEndian) }
func (c *Common) U16BE() uint64           { return c.UE(16, bitbuf.BigEndian) }
func (c *Common) U24BE() uint64           { return c.UE(24, bitbuf.BigEndian) }
func (c *Common) U32BE() uint64           { return c.UE(32, bitbuf.BigEndian) }
func (c *Common) U64BE() uint64           { return c.UE(64, bitbuf.BigEndian) }

func (c *Common) ULE(nBits uint64) uint64 { return c.UE(nBits, bitbuf.LittleEndian) }
func (c *Common) U9LE() uint64            { return c.UE(9, bitbuf.LittleEndian) }
func (c *Common) U10LE() uint64           { return c.UE(10, bitbuf.LittleEndian) }
func (c *Common) U11LE() uint64           { return c.UE(11, bitbuf.LittleEndian) }
func (c *Common) U12LE() uint64           { return c.UE(12, bitbuf.LittleEndian) }
func (c *Common) U13LE() uint64           { return c.UE(13, bitbuf.LittleEndian) }
func (c *Common) U14LE() uint64           { return c.UE(14, bitbuf.LittleEndian) }
func (c *Common) U15LE() uint64           { return c.UE(15, bitbuf.LittleEndian) }
func (c *Common) U16LE() uint64           { return c.UE(16, bitbuf.LittleEndian) }
func (c *Common) U24LE() uint64           { return c.UE(24, bitbuf.LittleEndian) }
func (c *Common) U32LE() uint64           { return c.UE(32, bitbuf.LittleEndian) }
func (c *Common) U64LE() uint64           { return c.UE(64, bitbuf.LittleEndian) }

func (c *Common) FieldBool(name string) bool {
	b, err := c.bitBuf.Bool()
	if err != nil {
		panic(Error{Err: err, Op: "FieldBool", Size: 1, Pos: c.bitBuf.Pos})
	}
	return b
}

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

func (c *Common) SE(nBits uint64, endian bitbuf.Endian) int64 {
	n, err := c.bitBuf.SE(nBits, endian)
	if err != nil {
		panic(Error{Err: err, Op: "SE", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) S(nBits uint64) int64 { return c.SE(nBits, bitbuf.BigEndian) }
func (c *Common) S1() int64            { return c.SE(1, bitbuf.BigEndian) }
func (c *Common) S2() int64            { return c.SE(2, bitbuf.BigEndian) }
func (c *Common) S3() int64            { return c.SE(3, bitbuf.BigEndian) }
func (c *Common) S4() int64            { return c.SE(4, bitbuf.BigEndian) }
func (c *Common) S5() int64            { return c.SE(5, bitbuf.BigEndian) }
func (c *Common) S6() int64            { return c.SE(6, bitbuf.BigEndian) }
func (c *Common) S7() int64            { return c.SE(7, bitbuf.BigEndian) }
func (c *Common) S8() int64            { return c.SE(8, bitbuf.BigEndian) }
func (c *Common) S9() int64            { return c.SE(9, bitbuf.BigEndian) }
func (c *Common) S10() int64           { return c.SE(10, bitbuf.BigEndian) }
func (c *Common) S11() int64           { return c.SE(11, bitbuf.BigEndian) }
func (c *Common) S12() int64           { return c.SE(12, bitbuf.BigEndian) }
func (c *Common) S13() int64           { return c.SE(13, bitbuf.BigEndian) }
func (c *Common) S14() int64           { return c.SE(14, bitbuf.BigEndian) }
func (c *Common) S15() int64           { return c.SE(15, bitbuf.BigEndian) }
func (c *Common) S16() int64           { return c.SE(16, bitbuf.BigEndian) }
func (c *Common) S24() int64           { return c.SE(24, bitbuf.BigEndian) }
func (c *Common) S32() int64           { return c.SE(32, bitbuf.BigEndian) }
func (c *Common) S64() int64           { return c.SE(64, bitbuf.BigEndian) }

func (c *Common) SBE(nBits uint64) int64 { return c.SE(nBits, bitbuf.BigEndian) }
func (c *Common) S9BE() int64            { return c.SE(9, bitbuf.BigEndian) }
func (c *Common) S10BE() int64           { return c.SE(10, bitbuf.BigEndian) }
func (c *Common) S11BE() int64           { return c.SE(11, bitbuf.BigEndian) }
func (c *Common) S12BE() int64           { return c.SE(12, bitbuf.BigEndian) }
func (c *Common) S13BE() int64           { return c.SE(13, bitbuf.BigEndian) }
func (c *Common) S14BE() int64           { return c.SE(14, bitbuf.BigEndian) }
func (c *Common) S15BE() int64           { return c.SE(15, bitbuf.BigEndian) }
func (c *Common) S16BE() int64           { return c.SE(16, bitbuf.BigEndian) }
func (c *Common) S24BE() int64           { return c.SE(24, bitbuf.BigEndian) }
func (c *Common) S32BE() int64           { return c.SE(32, bitbuf.BigEndian) }
func (c *Common) S64BE() int64           { return c.SE(64, bitbuf.BigEndian) }

func (c *Common) SLE(nBits uint64) int64 { return c.SE(nBits, bitbuf.LittleEndian) }
func (c *Common) S9LE() int64            { return c.SE(9, bitbuf.LittleEndian) }
func (c *Common) S10LE() int64           { return c.SE(10, bitbuf.LittleEndian) }
func (c *Common) S11LE() int64           { return c.SE(11, bitbuf.LittleEndian) }
func (c *Common) S12LE() int64           { return c.SE(12, bitbuf.LittleEndian) }
func (c *Common) S13LE() int64           { return c.SE(13, bitbuf.LittleEndian) }
func (c *Common) S14LE() int64           { return c.SE(14, bitbuf.LittleEndian) }
func (c *Common) S15LE() int64           { return c.SE(15, bitbuf.LittleEndian) }
func (c *Common) S16LE() int64           { return c.SE(16, bitbuf.LittleEndian) }
func (c *Common) S24LE() int64           { return c.SE(24, bitbuf.LittleEndian) }
func (c *Common) S32LE() int64           { return c.SE(32, bitbuf.LittleEndian) }
func (c *Common) S64LE() int64           { return c.SE(64, bitbuf.LittleEndian) }

func (c *Common) fieldS(name string, nBits uint64, endian bitbuf.Endian) int64 {
	start := c.bitBuf.Pos
	n, err := c.bitBuf.SE(nBits, endian)
	if err != nil {
		panic(Error{Err: err, Op: "FieldS" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
	}
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

func (c *Common) float32(fn func() (uint64, error)) float32 {
	n, err := fn()
	if err != nil {
		panic(Error{Err: err, Op: "Float32", Size: 32 * 8, Pos: c.bitBuf.Pos})
	}
	return math.Float32frombits(uint32(n))
}

func (c *Common) Float32(s uint) float32   { return c.float32(c.bitBuf.U32BE) }
func (c *Common) Float32BE(s uint) float32 { return c.float32(c.bitBuf.U32BE) }
func (c *Common) Float32LE(s uint) float32 { return c.float32(c.bitBuf.U32LE) }

func (c *Common) float64(fn func() (uint64, error)) float64 {
	n, err := fn()
	if err != nil {
		panic(Error{Err: err, Op: "Float64", Size: 64 * 8, Pos: c.bitBuf.Pos})
	}
	return math.Float64frombits(uint64(n))
}
func (c *Common) Float64(s uint) float64   { return c.float64(c.bitBuf.U64BE) }
func (c *Common) Float64BE(s uint) float64 { return c.float64(c.bitBuf.U64BE) }
func (c *Common) Float64LE(s uint) float64 { return c.float64(c.bitBuf.U64LE) }

func (c *Common) UTF8(nBytes uint64) string {
	s, err := c.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(Error{Err: err, Op: "UTF8", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return string(s)
}

func (c *Common) Unary(s uint64) uint {
	var n uint
	for {
		b, err := c.bitBuf.U1()
		if err != nil {
			panic(Error{Err: err, Op: "Unary", Size: 1, Pos: c.bitBuf.Pos})
		}
		if b != s {
			break
		}
		n++
	}
	return n
}

func (c *Common) FieldFn(name string, fn func() Value) Value {
	prev := c.Current

	f := &Field{Name: name}
	c.Current = f
	prev.Children = append(prev.Children, f)
	start := c.bitBuf.Pos
	f.Range.Start = start
	v := fn()
	f.Range.Stop = c.bitBuf.Pos
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
		bs, _ := c.bitBuf.BytesLen(nBytes)
		return bs, ""
	})
}

func (c *Common) FieldUTF8(name string, nBytes uint64) string {
	return c.FieldStrFn(name, func() (string, string) {
		str, _ := c.bitBuf.UTF8(nBytes)
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

func (c *Common) FieldVerifyString(name string, nBytes uint64, v string) bool {
	return c.FieldStrFn(name, func() (string, string) {
		str, _ := c.bitBuf.UTF8(nBytes)
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
