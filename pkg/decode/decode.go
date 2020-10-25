package decode

//go:generate sh -c "cat decode.gen.go.tmpl | go run ../../_dev/tmpl.go | gofmt > decode.gen.go"

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"fq/pkg/bitbuf"
	"io/ioutil"
	"runtime"
	"strconv"
)

type DecodeError struct {
	Err        error
	PanicStack string
}

func (de *DecodeError) Error() string { return de.Err.Error() }
func (de *DecodeError) Unwrap() error { return de.Err }

type BitBufError struct {
	Err   error
	Op    string
	Size  int64
	Delta int64
	Pos   int64
}

func (e BitBufError) Error() string {
	return fmt.Sprintf("%s: failed at position %s (size %s delta %s): %s",
		e.Op, Bits(e.Pos), Bits(e.Size), Bits(e.Delta), e.Err)
}
func (e BitBufError) Unwrap() error { return e.Err }

type ValidateError struct {
	Reason string
	Pos    int64
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("failed to validate at position %s: %s", Bits(e.Pos), e.Reason)
}

type D struct {
	bitBuf   *bitbuf.Buffer
	value    *Value
	registry *Registry
}

func (c *D) SafeDecodeFn(fn func(d *D) interface{}) (error, interface{}) {
	decodeErr, dv := func() (err error, dv interface{}) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				// https://github.com/golang/go/blob/master/src/net/http/server.go#L1770
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]

				pe := &DecodeError{
					PanicStack: string(buf),
				}
				switch panicErr := recoverErr.(type) {
				case BitBufError:
					pe.Err = panicErr
				case ValidateError:
					pe.Err = panicErr
				default:
					pe.Err = fmt.Errorf("%s", panicErr)
				}

				err = pe
			}
		}()

		return nil, fn(c)
	}()

	return decodeErr, dv
}

func (c *D) GetCommon() *D {
	return c
}

func (c *D) PeekBits(nBits int64) uint64 {
	n, err := c.bitBuf.PeekBits(nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBits", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *D) PeekBytes(nBytes int64) []byte {
	bs, err := c.bitBuf.PeekBytes(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBytes", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *D) PeekFind(nBits int64, v uint8, maxLen int64) int64 {
	peekBits, err := c.bitBuf.PeekFind(nBits, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFind", Size: 0, Pos: c.bitBuf.Pos})
	}
	return peekBits
}

func (c *D) TryHasBytes(hb []byte) bool {
	lenHb := int64(len(hb))
	if c.BitsLeft() < lenHb*8 {
		return false
	}
	bs := c.PeekBytes(lenHb)
	return bytes.Equal(hb, bs)
}

// PeekFindByte number of bytes to next v
func (c *D) PeekFindByte(v uint8, maxLen int64) int64 {
	peekBits, err := c.bitBuf.PeekFind(8, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFindByte", Size: 0, Pos: c.bitBuf.Pos})

	}
	return peekBits / 8
}

func (c *D) BytesRange(firstBit int64, nBytes int64) []byte {
	bs, err := c.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesRange", Size: nBytes * 8, Pos: firstBit})
	}
	return bs
}

func (c *D) BytesLen(nBytes int64) []byte {
	bs, err := c.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesLen", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *D) BitBufRange(firstBit int64, nBits int64) *bitbuf.Buffer {
	bs, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BitBufRange", Size: nBits, Pos: firstBit})
	}
	return bs
}

func (c *D) BitBufLen(nBits int64) *bitbuf.Buffer {
	bs, err := c.bitBuf.BitBufLen(nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BitBufLen", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *D) Pos() int64           { return c.bitBuf.Pos }
func (c *D) Len() int64           { return c.bitBuf.Len }
func (c *D) End() bool            { return c.bitBuf.End() }
func (c *D) BitsLeft() int64      { return c.bitBuf.BitsLeft() }
func (c *D) ByteAlignBits() int64 { return c.bitBuf.ByteAlignBits() }
func (c *D) BytePos() int64       { return c.bitBuf.BytePos() }

func (c *D) SeekRel(deltaBits int64) int64 {
	pos, err := c.bitBuf.SeekRel(deltaBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SeekRel", Delta: deltaBits, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *D) SeekAbs(pos int64) int64 {
	pos, err := c.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SeekAbs", Size: pos, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *D) UE(nBits int64, endian bitbuf.Endian) uint64 {
	n, err := c.bitBuf.UE(nBits, endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "UE", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *D) Bool() bool {
	b, err := c.bitBuf.Bool()
	if err != nil {
		panic(BitBufError{Err: err, Op: "Bool", Size: 1, Pos: c.bitBuf.Pos})
	}
	return b
}

func (c *D) FieldBool(name string) bool {
	return c.FieldBoolFn(name, func() (bool, string) {
		b, err := c.bitBuf.Bool()
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBool", Size: 1, Pos: c.bitBuf.Pos})
		}
		return b, ""
	})
}

func (c *D) FieldUE(name string, nBits int64, endian bitbuf.Endian) uint64 {
	return c.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := c.bitBuf.UE(nBits, endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldU" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (c *D) SE(nBits int64, endian bitbuf.Endian) int64 {
	n, err := c.bitBuf.SE(nBits, endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SE", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *D) FieldSE(name string, nBits int64, endian bitbuf.Endian) int64 {
	return c.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := c.bitBuf.SE(nBits, endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldS" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (c *D) F32E(endian bitbuf.Endian) float64 {
	f, err := c.bitBuf.F32E(endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "F32", Size: 32, Pos: c.bitBuf.Pos})
	}
	return float64(f)
}

func (c *D) F32() float64   { return c.F32E(bitbuf.BigEndian) }
func (c *D) F32BE() float64 { return c.F32E(bitbuf.BigEndian) }
func (c *D) F32LE() float64 { return c.F32E(bitbuf.LittleEndian) }

func (c *D) FieldF32E(name string, endian bitbuf.Endian) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		f, err := c.bitBuf.F32E(endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "F32", Size: 32, Pos: c.bitBuf.Pos})
		}
		return float64(f), ""
	})
}

func (c *D) FieldF32(name string) float64   { return c.FieldF32E(name, bitbuf.BigEndian) }
func (c *D) FieldF32BE(name string) float64 { return c.FieldF32E(name, bitbuf.BigEndian) }
func (c *D) FieldF32LE(name string) float64 { return c.FieldF32E(name, bitbuf.LittleEndian) }

func (c *D) F64E(endian bitbuf.Endian) float64 {
	f, err := c.bitBuf.F64E(endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "F64", Size: 64, Pos: c.bitBuf.Pos})
	}
	return float64(f)
}

func (c *D) F64() float64   { return c.F64E(bitbuf.BigEndian) }
func (c *D) F64BE() float64 { return c.F64E(bitbuf.BigEndian) }
func (c *D) F64LE() float64 { return c.F64E(bitbuf.LittleEndian) }

func (c *D) FieldF64E(name string, endian bitbuf.Endian) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		f, err := c.bitBuf.F64E(endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "F64", Size: 64, Pos: c.bitBuf.Pos})
		}
		return float64(f), ""
	})
}

func (c *D) FieldF64(name string) float64   { return c.FieldF64E(name, bitbuf.BigEndian) }
func (c *D) FieldF64BE(name string) float64 { return c.FieldF64E(name, bitbuf.BigEndian) }
func (c *D) FieldF64LE(name string) float64 { return c.FieldF64E(name, bitbuf.LittleEndian) }

func (c *D) UTF8(nBytes int64) string {
	s, err := c.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "UTF8", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return string(s)
}

func (c *D) FP64() float64 {
	f, err := c.bitBuf.FP64()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP64", Size: 8, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *D) FieldFP64(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.FP64(), ""
	})
}

func (c *D) FP32() float64 {
	f, err := c.bitBuf.FP32()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP32", Size: 4, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *D) FieldFP32(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.FP32(), ""
	})
}

func (c *D) FP16() float64 {
	f, err := c.bitBuf.FP16()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP16", Size: 2, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *D) FieldFP16(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.FP16(), ""
	})
}

func (c *D) UFP64() float64 {
	f, err := c.bitBuf.UFP64()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP64", Size: 8, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *D) FieldUFP64(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.UFP64(), ""
	})
}

func (c *D) UFP32() float64 {
	f, err := c.bitBuf.UFP32()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP32", Size: 4, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *D) FieldUFP32(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.UFP32(), ""
	})
}

func (c *D) UFP16() float64 {
	f, err := c.bitBuf.UFP16()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP16", Size: 2, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *D) FieldUFP16(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.UFP16(), ""
	})
}

func (c *D) Unary(s uint64) uint64 {
	n, err := c.bitBuf.Unary(s)
	if err != nil {
		panic(BitBufError{Err: err, Op: "Unary", Size: 1, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *D) ZeroPadding(nBits int64) bool {
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

func (c *D) AddChild(v *Value) {
	switch fv := c.value.V.(type) {
	case Struct:
		for _, ff := range fv {
			if ff.Name == v.Name {
				panic(fmt.Sprintf("%s already exist", v.Name))
			}
		}
		c.value.V = append(fv, v)
		return
	case Array:
		c.value.V = append(fv, v)
	}

}

func (c *D) fieldDecoder(name string, bitBuf *bitbuf.Buffer, v interface{}) *D {
	d := &D{
		bitBuf: bitBuf,
		// TODO: rename current to value?
		value: &Value{
			Name: name,
			V:    v,
		},
		registry: c.registry,
	}

	// TODO: find start/stop from Ranges instead? what if seekaround? concat bitbufs but want gaps? sort here, crash?

	// TODO: refactor
	if c.value != nil {
		c.AddChild(d.value)
	}
	return d
}

func (c *D) FieldArray(name string) *D {
	return c.fieldDecoder(name, c.bitBuf, Array{})
}

func (c *D) FieldArrayFn(name string, fn func(d *D)) *D {
	d := c.FieldArray(name)
	fn(d)
	return d
}

func (c *D) FieldStruct(name string) *D {
	return c.fieldDecoder(name, c.bitBuf, Struct{})
}

func (c *D) FieldStructFn(name string, fn func(d *D)) *D {
	d := c.FieldStruct(name)
	fn(d)
	return d
}

func (c *D) FieldStructBitBuf(name string, bitBuf *bitbuf.Buffer) *D {
	return c.fieldDecoder(name, bitBuf, Struct{})
}

func (c *D) FieldStructBitBufFn(name string, bitBuf *bitbuf.Buffer, fn func(d *D)) *D {
	d := c.FieldStructBitBuf(name, bitBuf)
	fn(d)
	return d
}

func (c *D) FieldRangeFn(name string, firstBit int64, nBits int64, fn func() Value) Value {
	v := fn()
	v.Name = name
	v.BitBuf = c.BitBufRange(firstBit, nBits)
	v.Range = Range{Start: firstBit, Stop: firstBit + nBits}
	c.AddChild(&v)

	return v
}

func (c *D) FieldFn(name string, fn func() Value) Value {
	start := c.bitBuf.Pos
	v := fn()
	stop := c.bitBuf.Pos
	v.Name = name
	v.BitBuf = c.BitBufRange(start, stop-start)
	v.Range = Range{Start: start, Stop: stop}
	c.AddChild(&v)

	return v
}

func (c *D) FieldBoolFn(name string, fn func() (bool, string)) bool {
	return c.FieldFn(name, func() Value {
		b, d := fn()
		return Value{V: b, Symbol: d}
	}).V.(bool)
}

func (c *D) FieldUFn(name string, fn func() (uint64, DisplayFormat, string)) uint64 {
	return c.FieldFn(name, func() Value {
		u, fmt, d := fn()
		return Value{V: u, DisplayFormat: fmt, Symbol: d}
	}).V.(uint64)
}

func (c *D) FieldSFn(name string, fn func() (int64, DisplayFormat, string)) int64 {
	return c.FieldFn(name, func() Value {
		s, fmt, d := fn()
		return Value{V: s, DisplayFormat: fmt, Symbol: d}
	}).V.(int64)
}

func (c *D) FieldFloatFn(name string, fn func() (float64, string)) float64 {
	return c.FieldFn(name, func() Value {
		f, d := fn()
		return Value{V: f, Symbol: d}
	}).V.(float64)
}

func (c *D) FieldStrFn(name string, fn func() (string, string)) string {
	return c.FieldFn(name, func() Value {
		str, disp := fn()
		return Value{V: str, Symbol: disp}
	}).V.(string)
}

func (c *D) FieldBytesFn(name string, firstBit int64, nBits int64, fn func() ([]byte, string)) []byte {
	return c.FieldRangeFn(name, firstBit, nBits, func() Value {
		bs, disp := fn()
		return Value{V: bs, Symbol: disp}
	}).V.([]byte)
}

func (c *D) FieldBitBufFn(name string, firstBit int64, nBits int64, fn func() (*bitbuf.Buffer, string)) *bitbuf.Buffer {
	return c.FieldRangeFn(name, firstBit, nBits, func() Value {
		bb, disp := fn()
		return Value{V: bb, Symbol: disp}
	}).BitBuf
}

func (c *D) FieldStringMapFn(name string, sm map[uint64]string, def string, fn func() uint64) (uint64, bool) {
	var ok bool
	return c.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n := fn()
		var d string
		d, ok = sm[n]
		if !ok {
			d = def
		}
		return n, NumberDecimal, d
	}), ok
}

func (c *D) FieldValidateUFn(name string, v uint64, fn func() uint64) {
	pos := c.bitBuf.Pos
	n := c.FieldUFn(name, func() (uint64, DisplayFormat, string) {
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
func (c *D) FieldBytesLen(name string, nBytes int64) []byte {
	return c.FieldBytesFn(name, c.bitBuf.Pos, nBytes*8, func() ([]byte, string) {
		bs, err := c.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesLen", Size: nBytes * 8, Pos: c.bitBuf.Pos})
		}
		return bs, ""
	})
}

func (c *D) FieldBytesRange(name string, firstBit int64, nBytes int64) []byte {
	return c.FieldBytesFn(name, firstBit, nBytes*8, func() ([]byte, string) {
		bs, err := c.bitBuf.BytesRange(firstBit, nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesRange", Size: nBytes * 8, Pos: firstBit})
		}
		return bs, ""
	})
}

func (c *D) FieldUTF8(name string, nBytes int64) string {
	return c.FieldStrFn(name, func() (string, string) {
		str, err := c.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldUTF8", Size: nBytes * 8, Pos: c.bitBuf.Pos})
		}
		return str, ""
	})
}

func (c *D) FieldValidateStringFn(name string, v string, fn func() string) {
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

func (c *D) FieldValidateString(name string, v string) {
	pos := c.bitBuf.Pos
	s := c.FieldStrFn(name, func() (string, string) {
		nBytes := int64(len(v))
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

func (c *D) FieldValidateZeroPadding(name string, nBits int64) {
	pos := c.bitBuf.Pos
	var isZero bool
	c.FieldFn(name, func() Value {
		isZero = c.ZeroPadding(nBits)
		s := "Correct"
		if !isZero {
			s = "Incorrect"
		}
		return Value{Symbol: s, Desc: "zero padding"}
	})
	if !isZero {
		panic(ValidateError{Reason: "expected zero padding", Pos: pos})
	}
}

func (c *D) ValidateAtLeastBitsLeft(nBits int64) {
	bl := c.bitBuf.BitsLeft()
	if bl < nBits {
		// TODO:
		panic(ValidateError{Reason: fmt.Sprintf("expected bits left %d, found %d", nBits, bl), Pos: c.bitBuf.Pos})
	}
}

func (c *D) ValidateAtLeastBytesLeft(nBytes int64) {
	bl := c.bitBuf.BitsLeft()
	if bl < nBytes*8 {
		// TODO:
		panic(ValidateError{Reason: fmt.Sprintf("expected bytes left %d, found %d bits", nBytes, bl), Pos: c.bitBuf.Pos})
	}
}

// Invalid stops decode with a reason
func (c *D) Invalid(reason string) {
	panic(ValidateError{Reason: reason, Pos: c.bitBuf.Pos})
}

// TODO: rename?
func (c *D) SubLenFn(nBits int64, fn func()) {
	prevBb := c.bitBuf

	bb, err := c.bitBuf.BitBufRange(0, c.bitBuf.Pos+nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SubLen", Size: nBits, Pos: c.bitBuf.Pos})
	}
	_, err = bb.SeekAbs(c.bitBuf.Pos)
	if err != nil {
		panic(err)
	}
	c.bitBuf = bb

	fn()

	bitsLeft := nBits - (c.bitBuf.Pos - prevBb.Pos)
	c.SeekRel(int64(bitsLeft))

	prevBb.Pos = c.bitBuf.Pos
	c.bitBuf = prevBb
}

func (c *D) SubRangeFn(firstBit int64, nBits int64, fn func()) {
	prevBb := c.bitBuf

	bb, err := c.bitBuf.BitBufRange(0, firstBit+nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SubRangeFn", Size: nBits, Pos: firstBit})
	}
	_, err = bb.SeekAbs(firstBit)
	if err != nil {
		panic(err)
	}
	c.bitBuf = bb

	fn()

	c.bitBuf = prevBb
}

// TODO: TryDecode?
func (c *D) FieldTryDecode(name string, forceFormats []*Format) (*Value, interface{}, []error) {
	bb, err := c.bitBuf.BitBufRange(c.bitBuf.Pos, c.BitsLeft())
	if err != nil {
		// TODO: can't happen?
		panic(BitBufError{Err: err, Op: "FieldDecode", Size: c.BitsLeft(), Pos: c.bitBuf.Pos})
	}

	v, fLen, dv, errs := c.registry.Probe(name, Range{Start: c.bitBuf.Pos, Stop: c.bitBuf.Pos}, bb, forceFormats)
	if v == nil || v.Errors() != nil {
		return nil, nil, errs
	}

	// TODO: bitbuf len shorten!
	c.AddChild(v)
	_, err = c.bitBuf.SeekRel(int64(fLen))
	if err != nil {
		panic(err)
	}

	return v, dv, errs
}

// TODO: FieldTryDecode? just TryDecode?
func (c *D) FieldDecodeLen(name string, nBits int64, forceFormats []*Format) (*Value, interface{}, []error) {
	bb, err := c.bitBuf.BitBufRange(c.bitBuf.Pos, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeLen", Size: nBits, Pos: c.bitBuf.Pos})
	}

	v, _, dv, errs := c.registry.Probe(name, Range{Start: c.bitBuf.Pos, Stop: c.bitBuf.Pos + nBits}, bb, forceFormats)
	if v != nil {
		c.AddChild(v)
	} else {
		// TODO: decoder unknown
		c.FieldRangeFn(name, c.bitBuf.Pos, nBits, func() Value { return Value{} })
	}

	// TODO: nBits - fLen gap?

	_, err = c.bitBuf.SeekRel(int64(nBits))
	if err != nil {
		panic(err)
	}

	return v, dv, errs
}

// TODO: return decooder?
func (c *D) FieldTryDecodeRange(name string, firstBit int64, nBits int64, forceFormats []*Format) (*Value, interface{}, []error) {
	bb, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: c.bitBuf.Pos})
	}

	v, _, dv, errs := c.registry.Probe(name, Range{Start: firstBit, Stop: firstBit + nBits}, bb, forceFormats)
	if v != nil {
		c.AddChild(v)
	}

	return v, dv, errs
}

// TODO: return decooder?
func (c *D) FieldDecodeRange(name string, firstBit int64, nBits int64, forceFormats []*Format) (*Value, interface{}, []error) {
	bb, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: c.bitBuf.Pos})
	}

	v, _, dv, errs := c.registry.Probe(name, Range{Start: firstBit, Stop: firstBit + nBits}, bb, forceFormats)
	if v != nil {
		c.AddChild(v)
	} else {
		c.FieldRangeFn(name, firstBit, nBits, func() Value { return Value{} })
	}

	return v, dv, errs
}

// TODO: list of ranges?
func (c *D) FieldDecodeBitBuf(name string, firstBit int64, nBits int64, bb *bitbuf.Buffer, forceFormats []*Format) (*Value, interface{}, []error) {
	f, _, dv, errs := c.registry.Probe(name, Range{Start: firstBit, Stop: nBits}, bb, forceFormats)
	if f != nil {
		c.AddChild(f)
	} else {
		c.FieldBitBufFn(name, firstBit, nBits, func() (*bitbuf.Buffer, string) {
			return bb, ""
		})
	}

	return f, dv, errs
}

func (c *D) FieldBitBufRange(name string, firstBit int64, nBits int64) *bitbuf.Buffer {
	return c.FieldBitBufFn(name, firstBit, nBits, func() (*bitbuf.Buffer, string) {
		return c.BitBufRange(firstBit, nBits), ""
	})
}

func (c *D) FieldBitBufLen(name string, nBits int64) *bitbuf.Buffer {
	return c.FieldBitBufFn(name, c.bitBuf.Pos, nBits, func() (*bitbuf.Buffer, string) {
		return c.BitBufLen(nBits), ""
	})
}

func (c *D) FieldZlib(name string, firsBit int64, nBits int64, b []byte, formats []*Format) (*Value, interface{}, []error) {
	zr, err := zlib.NewReader(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	zd, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	zbb, err := bitbuf.NewFromBytes(zd, 0)
	if err != nil {
		return nil, nil, []error{err}
	}

	return c.FieldDecodeBitBuf(name, firsBit, nBits, zbb, formats)
}

// TODO: range?
func (c *D) FieldZlibLen(name string, nBytes int64, formats []*Format) (*Value, interface{}, []error) {
	firstBit := c.bitBuf.Pos
	zr, err := zlib.NewReader(bytes.NewReader(c.BytesLen(nBytes)))
	if err != nil {
		panic(err)
	}
	zd, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	zbb, err := bitbuf.NewFromBytes(zd, 0)
	if err != nil {
		return nil, nil, []error{err}
	}

	return c.FieldDecodeBitBuf(name, firstBit, firstBit+nBytes*8, zbb, formats)
}
