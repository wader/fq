package decode

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"fq/internal/bitbuf"
	"io/ioutil"
	"strconv"
)

type Decoder interface {
	Decode()

	Prepare(common Common)
	Finish(err error)

	Root() *Field
	BitBuf() *bitbuf.Buffer
	Format() *Format
	Error() error
}

type BitBufError struct {
	Err   error
	Op    string
	Size  uint64
	Delta int64
	Pos   uint64
}

func (e BitBufError) Error() string {
	return fmt.Sprintf("%s: failed at bit position %d (%d+%d) (size %d (%d+%d), delta %d (%d+%d)): %s",
		e.Op, e.Pos, e.Pos>>3, e.Pos&0x7, e.Size, e.Size>>3, e.Size&0x7, e.Delta, e.Delta>>3, e.Delta&0x7, e.Err)
}
func (e BitBufError) Unwrap() error { return e.Err }

type ValidateError struct {
	Reason string
	Pos    uint64
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("failed to validate at position %d (%d+%d): %s", e.Pos, e.Pos>>3, e.Pos&0x7, e.Reason)
}

type Format struct {
	Name      string
	MIME      string
	New       func() Decoder
	SkipProbe bool
}

func FormatFn(d func(c *Common)) *Format {
	return &Format{
		New: func() Decoder { return &DecoderFn{decode: d} },
	}
}

type DecoderFn struct {
	Common
	decode func(c *Common)
}

func (d *DecoderFn) Decode() {
	d.decode(&d.Common)
}

type Common struct {
	Parent   Decoder
	format   *Format
	registry *Registry
	bitBuf   *bitbuf.Buffer
	root     *Field
	current  *Field // TODO: need root field also?
	err      error
}

func (c *Common) Decode() {}

func (c *Common) Prepare(common Common) {
	*c = common
}

func (c *Common) Finish(err error) {
	c.err = err
	c.root.Sort()
}

func (c *Common) Root() *Field           { return c.root }
func (c *Common) BitBuf() *bitbuf.Buffer { return c.bitBuf }
func (c *Common) Format() *Format        { return c.format }
func (c *Common) Error() error           { return c.err }

func (c *Common) PeekBits(nBits uint64) uint64 {
	n, err := c.bitBuf.PeekBits(nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBits", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) PeekBytes(nBytes uint64) []byte {
	bs, err := c.bitBuf.PeekBytes(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBytes", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *Common) PeekFind(bits uint64, v uint8, maxLen int64) uint64 {
	peekBits, err := c.bitBuf.PeekFind(bits, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFind", Size: 0, Pos: c.bitBuf.Pos})
	}
	return peekBits
}

func (c *Common) PeekFindByte(v uint8, maxLen int64) uint64 {
	peekBits, err := c.bitBuf.PeekFind(8, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFindByte", Size: 0, Pos: c.bitBuf.Pos})

	}
	return peekBits / 8
}

func (c *Common) BytesRange(firstBit uint64, nBytes uint64) []byte {
	bs, err := c.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesRange", Size: nBytes * 8, Pos: firstBit})
	}
	return bs
}

func (c *Common) BytesLen(nBytes uint64) []byte {
	bs, err := c.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesLen", Size: nBytes * 8, Pos: c.bitBuf.Pos})
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
		panic(BitBufError{Err: err, Op: "SeekRel", Delta: delta, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *Common) SeekAbs(pos uint64) uint64 {
	pos, err := c.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SeekAbs", Size: pos, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *Common) UE(nBits uint64, endian bitbuf.Endian) uint64 {
	n, err := c.bitBuf.UE(nBits, endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "UE", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) Bool() bool {
	b, err := c.bitBuf.Bool()
	if err != nil {
		panic(BitBufError{Err: err, Op: "Bool", Size: 1, Pos: c.bitBuf.Pos})
	}
	return b
}

func (c *Common) FieldBool(name string) bool {
	return c.FieldBoolFn(name, func() (bool, string) {
		b, err := c.bitBuf.Bool()
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBool", Size: 1, Pos: c.bitBuf.Pos})
		}
		return b, ""
	})
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

func (c *Common) FieldUE(name string, nBits uint64, endian bitbuf.Endian) uint64 {
	return c.FieldUFn(name, func() (uint64, NumberFormat, string) {
		n, err := c.bitBuf.UE(nBits, endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldU" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (c *Common) FieldU(name string, nBits uint64) uint64 {
	return c.FieldUE(name, nBits, bitbuf.BigEndian)
}
func (c *Common) FieldU1(name string) uint64  { return c.FieldUE(name, 1, bitbuf.BigEndian) }
func (c *Common) FieldU2(name string) uint64  { return c.FieldUE(name, 2, bitbuf.BigEndian) }
func (c *Common) FieldU3(name string) uint64  { return c.FieldUE(name, 3, bitbuf.BigEndian) }
func (c *Common) FieldU4(name string) uint64  { return c.FieldUE(name, 4, bitbuf.BigEndian) }
func (c *Common) FieldU5(name string) uint64  { return c.FieldUE(name, 5, bitbuf.BigEndian) }
func (c *Common) FieldU6(name string) uint64  { return c.FieldUE(name, 6, bitbuf.BigEndian) }
func (c *Common) FieldU7(name string) uint64  { return c.FieldUE(name, 7, bitbuf.BigEndian) }
func (c *Common) FieldU8(name string) uint64  { return c.FieldUE(name, 8, bitbuf.BigEndian) }
func (c *Common) FieldU9(name string) uint64  { return c.FieldUE(name, 9, bitbuf.BigEndian) }
func (c *Common) FieldU10(name string) uint64 { return c.FieldUE(name, 10, bitbuf.BigEndian) }
func (c *Common) FieldU11(name string) uint64 { return c.FieldUE(name, 11, bitbuf.BigEndian) }
func (c *Common) FieldU12(name string) uint64 { return c.FieldUE(name, 12, bitbuf.BigEndian) }
func (c *Common) FieldU13(name string) uint64 { return c.FieldUE(name, 13, bitbuf.BigEndian) }
func (c *Common) FieldU14(name string) uint64 { return c.FieldUE(name, 14, bitbuf.BigEndian) }
func (c *Common) FieldU15(name string) uint64 { return c.FieldUE(name, 15, bitbuf.BigEndian) }
func (c *Common) FieldU16(name string) uint64 { return c.FieldUE(name, 16, bitbuf.BigEndian) }
func (c *Common) FieldU24(name string) uint64 { return c.FieldUE(name, 24, bitbuf.BigEndian) }
func (c *Common) FieldU32(name string) uint64 { return c.FieldUE(name, 32, bitbuf.BigEndian) }
func (c *Common) FieldU64(name string) uint64 { return c.FieldUE(name, 64, bitbuf.BigEndian) }

func (c *Common) FieldUBE(nBits uint64, name string) uint64 {
	return c.FieldUE(name, nBits, bitbuf.BigEndian)
}
func (c *Common) FieldU9BE(name string) uint64  { return c.FieldUE(name, 9, bitbuf.BigEndian) }
func (c *Common) FieldU10BE(name string) uint64 { return c.FieldUE(name, 10, bitbuf.BigEndian) }
func (c *Common) FieldU11BE(name string) uint64 { return c.FieldUE(name, 11, bitbuf.BigEndian) }
func (c *Common) FieldU12BE(name string) uint64 { return c.FieldUE(name, 12, bitbuf.BigEndian) }
func (c *Common) FieldU13BE(name string) uint64 { return c.FieldUE(name, 13, bitbuf.BigEndian) }
func (c *Common) FieldU14BE(name string) uint64 { return c.FieldUE(name, 14, bitbuf.BigEndian) }
func (c *Common) FieldU15BE(name string) uint64 { return c.FieldUE(name, 15, bitbuf.BigEndian) }
func (c *Common) FieldU16BE(name string) uint64 { return c.FieldUE(name, 16, bitbuf.BigEndian) }
func (c *Common) FieldU24BE(name string) uint64 { return c.FieldUE(name, 24, bitbuf.BigEndian) }
func (c *Common) FieldU32BE(name string) uint64 { return c.FieldUE(name, 32, bitbuf.BigEndian) }
func (c *Common) FieldU64BE(name string) uint64 { return c.FieldUE(name, 64, bitbuf.BigEndian) }

func (c *Common) FieldULE(nBits uint64, name string) uint64 {
	return c.FieldUE(name, nBits, bitbuf.LittleEndian)
}
func (c *Common) FieldU9LE(name string) uint64  { return c.FieldUE(name, 9, bitbuf.LittleEndian) }
func (c *Common) FieldU10LE(name string) uint64 { return c.FieldUE(name, 10, bitbuf.LittleEndian) }
func (c *Common) FieldU11LE(name string) uint64 { return c.FieldUE(name, 11, bitbuf.LittleEndian) }
func (c *Common) FieldU12LE(name string) uint64 { return c.FieldUE(name, 12, bitbuf.LittleEndian) }
func (c *Common) FieldU13LE(name string) uint64 { return c.FieldUE(name, 13, bitbuf.LittleEndian) }
func (c *Common) FieldU14LE(name string) uint64 { return c.FieldUE(name, 14, bitbuf.LittleEndian) }
func (c *Common) FieldU15LE(name string) uint64 { return c.FieldUE(name, 15, bitbuf.LittleEndian) }
func (c *Common) FieldU16LE(name string) uint64 { return c.FieldUE(name, 16, bitbuf.LittleEndian) }
func (c *Common) FieldU24LE(name string) uint64 { return c.FieldUE(name, 24, bitbuf.LittleEndian) }
func (c *Common) FieldU32LE(name string) uint64 { return c.FieldUE(name, 32, bitbuf.LittleEndian) }
func (c *Common) FieldU64LE(name string) uint64 { return c.FieldUE(name, 64, bitbuf.LittleEndian) }

func (c *Common) SE(nBits uint64, endian bitbuf.Endian) int64 {
	n, err := c.bitBuf.SE(nBits, endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SE", Size: nBits, Pos: c.bitBuf.Pos})
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

func (c *Common) FieldSE(name string, nBits uint64, endian bitbuf.Endian) int64 {
	return c.FieldSFn(name, func() (int64, NumberFormat, string) {
		n, err := c.bitBuf.SE(nBits, endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldS" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}
func (c *Common) FieldS(name string, nBits uint64) int64 {
	return c.FieldSE(name, nBits, bitbuf.BigEndian)
}
func (c *Common) FieldS1(name string) int64  { return c.FieldSE(name, 1, bitbuf.BigEndian) }
func (c *Common) FieldS2(name string) int64  { return c.FieldSE(name, 2, bitbuf.BigEndian) }
func (c *Common) FieldS3(name string) int64  { return c.FieldSE(name, 3, bitbuf.BigEndian) }
func (c *Common) FieldS4(name string) int64  { return c.FieldSE(name, 4, bitbuf.BigEndian) }
func (c *Common) FieldS5(name string) int64  { return c.FieldSE(name, 5, bitbuf.BigEndian) }
func (c *Common) FieldS6(name string) int64  { return c.FieldSE(name, 6, bitbuf.BigEndian) }
func (c *Common) FieldS7(name string) int64  { return c.FieldSE(name, 7, bitbuf.BigEndian) }
func (c *Common) FieldS8(name string) int64  { return c.FieldSE(name, 8, bitbuf.BigEndian) }
func (c *Common) FieldS9(name string) int64  { return c.FieldSE(name, 9, bitbuf.BigEndian) }
func (c *Common) FieldS10(name string) int64 { return c.FieldSE(name, 10, bitbuf.BigEndian) }
func (c *Common) FieldS11(name string) int64 { return c.FieldSE(name, 11, bitbuf.BigEndian) }
func (c *Common) FieldS12(name string) int64 { return c.FieldSE(name, 12, bitbuf.BigEndian) }
func (c *Common) FieldS13(name string) int64 { return c.FieldSE(name, 13, bitbuf.BigEndian) }
func (c *Common) FieldS14(name string) int64 { return c.FieldSE(name, 14, bitbuf.BigEndian) }
func (c *Common) FieldS15(name string) int64 { return c.FieldSE(name, 15, bitbuf.BigEndian) }
func (c *Common) FieldS16(name string) int64 { return c.FieldSE(name, 16, bitbuf.BigEndian) }
func (c *Common) FieldS24(name string) int64 { return c.FieldSE(name, 24, bitbuf.BigEndian) }
func (c *Common) FieldS32(name string) int64 { return c.FieldSE(name, 32, bitbuf.BigEndian) }
func (c *Common) FieldS64(name string) int64 { return c.FieldSE(name, 64, bitbuf.BigEndian) }

func (c *Common) FieldSBE(name string, nBits uint64) int64 {
	return c.FieldSE(name, nBits, bitbuf.BigEndian)
}
func (c *Common) FieldS9BE(name string) int64  { return c.FieldSE(name, 9, bitbuf.BigEndian) }
func (c *Common) FieldS10BE(name string) int64 { return c.FieldSE(name, 10, bitbuf.BigEndian) }
func (c *Common) FieldS11BE(name string) int64 { return c.FieldSE(name, 11, bitbuf.BigEndian) }
func (c *Common) FieldS12BE(name string) int64 { return c.FieldSE(name, 12, bitbuf.BigEndian) }
func (c *Common) FieldS13BE(name string) int64 { return c.FieldSE(name, 13, bitbuf.BigEndian) }
func (c *Common) FieldS14BE(name string) int64 { return c.FieldSE(name, 14, bitbuf.BigEndian) }
func (c *Common) FieldS15BE(name string) int64 { return c.FieldSE(name, 15, bitbuf.BigEndian) }
func (c *Common) FieldS16BE(name string) int64 { return c.FieldSE(name, 16, bitbuf.BigEndian) }
func (c *Common) FieldS24BE(name string) int64 { return c.FieldSE(name, 24, bitbuf.BigEndian) }
func (c *Common) FieldS32BE(name string) int64 { return c.FieldSE(name, 32, bitbuf.BigEndian) }
func (c *Common) FieldS64BE(name string) int64 { return c.FieldSE(name, 64, bitbuf.BigEndian) }

func (c *Common) FieldSLE(nBits uint64, name string) int64 {
	return c.FieldSE(name, nBits, bitbuf.LittleEndian)
}
func (c *Common) FieldS9LE(name string) int64  { return c.FieldSE(name, 9, bitbuf.LittleEndian) }
func (c *Common) FieldS10LE(name string) int64 { return c.FieldSE(name, 10, bitbuf.LittleEndian) }
func (c *Common) FieldS11LE(name string) int64 { return c.FieldSE(name, 11, bitbuf.LittleEndian) }
func (c *Common) FieldS12LE(name string) int64 { return c.FieldSE(name, 12, bitbuf.LittleEndian) }
func (c *Common) FieldS13LE(name string) int64 { return c.FieldSE(name, 13, bitbuf.LittleEndian) }
func (c *Common) FieldS14LE(name string) int64 { return c.FieldSE(name, 14, bitbuf.LittleEndian) }
func (c *Common) FieldS15LE(name string) int64 { return c.FieldSE(name, 15, bitbuf.LittleEndian) }
func (c *Common) FieldS16LE(name string) int64 { return c.FieldSE(name, 16, bitbuf.LittleEndian) }
func (c *Common) FieldS24LE(name string) int64 { return c.FieldSE(name, 24, bitbuf.LittleEndian) }
func (c *Common) FieldS32LE(name string) int64 { return c.FieldSE(name, 32, bitbuf.LittleEndian) }
func (c *Common) FieldS64LE(name string) int64 { return c.FieldSE(name, 64, bitbuf.LittleEndian) }

func (c *Common) F32E(endian bitbuf.Endian) float64 {
	f, err := c.bitBuf.F32E(endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "F32", Size: 32, Pos: c.bitBuf.Pos})
	}
	return float64(f)
}

func (c *Common) F32() float64   { return c.F32E(bitbuf.BigEndian) }
func (c *Common) F32BE() float64 { return c.F32E(bitbuf.BigEndian) }
func (c *Common) F32LE() float64 { return c.F32E(bitbuf.LittleEndian) }

func (c *Common) FieldF32E(name string, endian bitbuf.Endian) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		f, err := c.bitBuf.F32E(endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "F32", Size: 32, Pos: c.bitBuf.Pos})
		}
		return float64(f), ""
	})
}

func (c *Common) FieldF32(name string) float64   { return c.FieldF32E(name, bitbuf.BigEndian) }
func (c *Common) FieldF32BE(name string) float64 { return c.FieldF32E(name, bitbuf.BigEndian) }
func (c *Common) FieldF32LE(name string) float64 { return c.FieldF32E(name, bitbuf.LittleEndian) }

func (c *Common) F64E(endian bitbuf.Endian) float64 {
	f, err := c.bitBuf.F64E(endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "F64", Size: 64, Pos: c.bitBuf.Pos})
	}
	return float64(f)
}

func (c *Common) F64() float64   { return c.F64E(bitbuf.BigEndian) }
func (c *Common) F64BE() float64 { return c.F64E(bitbuf.BigEndian) }
func (c *Common) F64LE() float64 { return c.F64E(bitbuf.LittleEndian) }

func (c *Common) FieldF64E(name string, endian bitbuf.Endian) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		f, err := c.bitBuf.F64E(endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "F64", Size: 64, Pos: c.bitBuf.Pos})
		}
		return float64(f), ""
	})
}

func (c *Common) FieldF64(name string) float64   { return c.FieldF64E(name, bitbuf.BigEndian) }
func (c *Common) FieldF64BE(name string) float64 { return c.FieldF64E(name, bitbuf.BigEndian) }
func (c *Common) FieldF64LE(name string) float64 { return c.FieldF64E(name, bitbuf.LittleEndian) }

func (c *Common) UTF8(nBytes uint64) string {
	s, err := c.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "UTF8", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return string(s)
}

func (c *Common) FP64() float64 {
	f, err := c.bitBuf.FP64()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP64", Size: 8, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *Common) FieldFP64(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.FP64(), ""
	})
}

func (c *Common) FP32() float64 {
	f, err := c.bitBuf.FP32()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP32", Size: 4, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *Common) FieldFP32(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.FP32(), ""
	})
}

func (c *Common) FP16() float64 {
	f, err := c.bitBuf.FP16()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP16", Size: 2, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *Common) FieldFP16(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.FP16(), ""
	})
}

func (c *Common) Unary(s uint64) uint64 {
	n, err := c.bitBuf.Unary(s)
	if err != nil {
		panic(BitBufError{Err: err, Op: "Unary", Size: 1, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) ZeroPadding(nBits uint64) bool {
	isZero := true
	left := nBits
	for {
		// TODO: smart skip?
		rbits := left
		if rbits == 0 {
			break
		}
		if rbits > 64 {
			rbits = 64
		}
		n, err := c.bitBuf.Bits(rbits)
		if err != nil {
			panic(BitBufError{Err: err, Op: "ZeroPadding", Size: rbits, Pos: c.bitBuf.Pos})
		}
		isZero = isZero && n == 0
		left -= rbits
	}
	return isZero
}

func (c *Common) FieldRangeFn(name string, firstBit uint64, nBits uint64, fn func() Value) Value {
	f := &Field{
		Name:  name,
		Range: Range{Start: firstBit, Stop: firstBit + nBits},
	}
	c.current.Children = append(c.current.Children, f)
	f.Value = fn()

	return f.Value
}

func (c *Common) FieldFn(name string, fn func() Value) Value {
	prev := c.current

	f := &Field{Name: name}
	c.current = f
	prev.Children = append(prev.Children, f)
	start := c.bitBuf.Pos
	f.Range.Start = start
	v := fn()
	f.Range.Stop = c.bitBuf.Pos
	f.Value = v

	c.current = prev

	return v
}

func (c *Common) FieldNoneFn(name string, fn func()) {
	c.FieldFn(name, func() Value {
		fn()
		return Value{}
	})
}

func (c *Common) FieldBoolFn(name string, fn func() (bool, string)) bool {
	return c.FieldFn(name, func() Value {
		b, d := fn()
		return Value{Type: TypeBool, Bool: b, Display: d}
	}).Bool
}

func (c *Common) FieldUFn(name string, fn func() (uint64, NumberFormat, string)) uint64 {
	return c.FieldFn(name, func() Value {
		u, fmt, d := fn()
		return Value{Type: TypeUInt, UInt: u, Format: fmt, Display: d}
	}).UInt
}

func (c *Common) FieldSFn(name string, fn func() (int64, NumberFormat, string)) int64 {
	return c.FieldFn(name, func() Value {
		s, fmt, d := fn()
		return Value{Type: TypeSInt, SInt: s, Format: fmt, Display: d}
	}).SInt
}

func (c *Common) FieldFloatFn(name string, fn func() (float64, string)) float64 {
	return c.FieldFn(name, func() Value {
		f, d := fn()
		return Value{Type: TypeFloat, Float: f, Display: d}
	}).Float
}

func (c *Common) FieldStrFn(name string, fn func() (string, string)) string {
	return c.FieldFn(name, func() Value {
		str, disp := fn()
		return Value{Type: TypeStr, Str: str, Display: disp}
	}).Str
}

func (c *Common) FieldBytesFn(name string, firstBit uint64, nBits uint64, fn func() ([]byte, string)) []byte {
	return c.FieldRangeFn(name, firstBit, nBits, func() Value {
		bs, disp := fn()
		return Value{Type: TypeBytes, Bytes: bs, Display: disp}
	}).Bytes
}

func (c *Common) FieldStringMapFn(name string, sm map[uint64]string, def string, fn func() uint64) uint64 {
	return c.FieldUFn(name, func() (uint64, NumberFormat, string) {
		n := fn()
		var d string
		d, ok := sm[n]
		if !ok {
			d = def
		}
		return n, NumberDecimal, d
	})
}

func (c *Common) FieldValidateUFn(name string, v uint64, fn func() uint64) {
	pos := c.bitBuf.Pos
	n := c.FieldUFn(name, func() (uint64, NumberFormat, string) {
		n := fn()
		s := "Correct"
		if n != v {
			s = "Incorrect"
		}
		return n, NumberHex, s
	})
	if n != v {
		panic(ValidateError{Reason: fmt.Sprintf("expected %d found %d", v, n), Pos: pos})
	}
}

// TODO: FieldBytesRange or?
func (c *Common) FieldBytesLen(name string, nBytes uint64) []byte {
	return c.FieldBytesFn(name, c.bitBuf.Pos, nBytes*8, func() ([]byte, string) {
		bs, err := c.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesLen", Size: nBytes * 8, Pos: c.bitBuf.Pos})
		}
		return bs, ""
	})
}

func (c *Common) FieldBytesRange(name string, firstBit uint64, nBytes uint64) []byte {
	return c.FieldBytesFn(name, firstBit, nBytes*8, func() ([]byte, string) {
		bs, err := c.bitBuf.BytesRange(firstBit, nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesRange", Size: nBytes * 8, Pos: firstBit})
		}
		return bs, ""
	})
}

func (c *Common) FieldUTF8(name string, nBytes uint64) string {
	return c.FieldStrFn(name, func() (string, string) {
		str, err := c.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldUTF8", Size: nBytes * 8, Pos: c.bitBuf.Pos})
		}
		return str, ""
	})
}

func (c *Common) FieldValidateStringFn(name string, v string, fn func() string) {
	pos := c.bitBuf.Pos
	s := c.FieldStrFn(name, func() (string, string) {
		str := fn()
		s := "Correct"
		if str != v {
			s = "Incorrect"
		}
		return str, s
	})
	if s != v {
		panic(ValidateError{Pos: pos})
	}
}

func (c *Common) FieldValidateString(name string, v string) {
	pos := c.bitBuf.Pos
	s := c.FieldStrFn(name, func() (string, string) {
		nBytes := uint64(len(v))
		str, err := c.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldValidateString", Size: nBytes * 8, Pos: c.bitBuf.Pos})
		}
		s := "Correct"
		if str != v {
			s = "Incorrect"
		}
		return str, s
	})
	if s != v {
		panic(ValidateError{Reason: fmt.Sprintf("expected %s found %s", v, s), Pos: pos})
	}
}

func (c *Common) FieldValidateZeroPadding(name string, nBits uint64) {
	pos := c.bitBuf.Pos
	var isZero bool
	c.FieldFn(name, func() Value {
		isZero = c.ZeroPadding(nBits)
		s := "Correct"
		if !isZero {
			s = "Incorrect"
		}
		return Value{Type: TypePadding, Display: s}
	})
	if !isZero {
		panic(ValidateError{Reason: fmt.Sprintf("expected zero padding"), Pos: pos})
	}
}

func (c *Common) ValidateAtLeastBitsLeft(nBits uint64) {
	bl := c.bitBuf.BitsLeft()
	if bl < nBits {
		// TODO:
		panic(ValidateError{Reason: "not enough bits left", Pos: c.bitBuf.Pos})
	}
}

func (c *Common) ValidateAtLeastBytesLeft(nBytes uint64) {
	bl := c.bitBuf.BitsLeft()
	if bl < nBytes*8 {
		// TODO:
		panic(ValidateError{Reason: "not enough bytes left", Pos: c.bitBuf.Pos})
	}
}

// Invalid stops decode with a reason
func (c *Common) Invalid(reason string) {
	panic(ValidateError{Reason: reason, Pos: c.bitBuf.Pos})
}

// TODO: rename?
func (c *Common) SubLenFn(nBits uint64, fn func()) {
	prevBb := c.bitBuf

	bb, err := c.bitBuf.BitBufRange(0, c.bitBuf.Pos+nBits)
	bb.SeekAbs(c.bitBuf.Pos)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SubLen", Size: nBits, Pos: c.bitBuf.Pos})

	}
	c.bitBuf = bb

	fn()

	bitsLeft := nBits - (c.bitBuf.Pos - prevBb.Pos)
	c.SeekRel(int64(bitsLeft))

	prevBb.Pos = c.bitBuf.Pos
	c.bitBuf = prevBb
}

// TODO: TryDecode?
func (c *Common) FieldTryDecode(name string, forceFormats ...*Format) (Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(c.bitBuf.Pos, c.BitsLeft())
	if err != nil {
		// TODO: can't happen?
		panic(BitBufError{Err: err, Op: "FieldDecode", Size: c.BitsLeft(), Pos: c.bitBuf.Pos})
	}

	d, errs := c.registry.Probe(c, name, Range{Start: c.bitBuf.Pos, Stop: c.bitBuf.Pos}, bb, forceFormats)
	if d == nil || d.Error() != nil {
		return nil, errs
	}

	// TODO: bitbuf len shorten!

	dbb := d.BitBuf()
	dbb.TruncateRel(0)
	df := d.Root()
	df.Range.Stop += dbb.Pos
	c.current.Children = append(c.current.Children, df)

	c.bitBuf.SeekRel(int64(d.BitBuf().Pos))

	return d, errs
}

// TODO: FieldTryDecode? just TryDecode?
func (c *Common) FieldDecodeLen(name string, nBits uint64, forceFormats ...*Format) (Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(c.bitBuf.Pos, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeLen", Size: nBits, Pos: c.bitBuf.Pos})
	}

	d, errs := c.registry.Probe(c, name, Range{Start: c.bitBuf.Pos, Stop: c.bitBuf.Pos + nBits}, bb, forceFormats)
	if d != nil {
		c.current.Children = append(c.current.Children, d.Root())
	} else {
		// TODO: decoder unknown
		c.FieldRangeFn(name, c.bitBuf.Pos, nBits, func() Value { return Value{} })
	}

	c.bitBuf.SeekRel(int64(nBits))

	return d, errs
}

// TODO: return decooder?
func (c *Common) FieldTryDecodeRange(name string, firstBit uint64, nBits uint64, forceFormats ...*Format) (Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: c.bitBuf.Pos})
	}

	d, errs := c.registry.Probe(c, name, Range{Start: firstBit, Stop: firstBit + nBits}, bb, forceFormats)
	if d != nil {
		c.current.Children = append(c.current.Children, d.Root())
	}

	return d, errs
}

// TODO: return decooder?
func (c *Common) FieldDecodeRange(name string, firstBit uint64, nBits uint64, forceFormats ...*Format) (Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: c.bitBuf.Pos})
	}

	d, errs := c.registry.Probe(c, name, Range{Start: firstBit, Stop: firstBit + nBits}, bb, forceFormats)
	if d != nil {
		c.current.Children = append(c.current.Children, d.Root())
	} else {
		c.FieldRangeFn(name, firstBit, nBits, func() Value { return Value{} })
	}

	return d, errs
}

// TODO: list of ranges?
func (c *Common) FieldDecodeBitBuf(name string, firstBit uint64, nBits uint64, bb *bitbuf.Buffer, forceFormats ...*Format) (Decoder, []error) {
	d, errs := c.registry.Probe(c, name, Range{Start: firstBit, Stop: nBits}, bb, forceFormats)
	if d != nil {
		c.current.Children = append(c.current.Children, d.Root())
	} else {
		c.FieldRangeFn(name, firstBit, nBits, func() Value { return Value{} })
	}

	return d, errs
}

func (c *Common) FieldZlib(name string, firsBit uint64, nBits uint64, b []byte, forceFormats ...*Format) (Decoder, []error) {
	zr, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	zd, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	return c.FieldDecodeBitBuf(name, firsBit, nBits, bitbuf.NewFromBytes(zd), forceFormats...)
}

// TODO: range?
func (c *Common) FieldZlibLen(name string, nBytes uint64, forceFormats ...*Format) (Decoder, []error) {
	firstBit := c.bitBuf.Pos
	zr, err := zlib.NewReader(bytes.NewReader(c.BytesLen(nBytes)))
	if err != nil {
		panic(err)
	}
	zd, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	return c.FieldDecodeBitBuf(name, firstBit, firstBit+nBytes*8, bitbuf.NewFromBytes(zd), forceFormats...)
}
