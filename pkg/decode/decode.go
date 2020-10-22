package decode

//go:generate sh -c "cat decode.gen.go.tmpl | go run ../../_dev/tmpl.go > decode.gen.go"

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

type Decoder interface {
	Decode()
	GetCommon() *Common // TODO: rename
}

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

type Common struct {
	bitBuf *bitbuf.Buffer

	current *Value // TODO: need root field also?

	registry *Registry
}

func (c *Common) Decode() {}

func (c *Common) SafeDecodeFn(fn func()) error {
	decodeErr := func() (err error) {
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

		fn()

		return nil
	}()

	return decodeErr
}

func (c *Common) SafeDecodeFn2(fn func(d *Common)) error {
	decodeErr := func() (err error) {
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

		fn(c)

		return nil
	}()

	return decodeErr
}

func (c *Common) GetCommon() *Common {
	return c
}

func (c *Common) PeekBits(nBits int64) uint64 {
	n, err := c.bitBuf.PeekBits(nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBits", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) PeekBytes(nBytes int64) []byte {
	bs, err := c.bitBuf.PeekBytes(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBytes", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *Common) PeekFind(nBits int64, v uint8, maxLen int64) int64 {
	peekBits, err := c.bitBuf.PeekFind(nBits, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFind", Size: 0, Pos: c.bitBuf.Pos})
	}
	return peekBits
}

func (c *Common) TryHasBytes(hb []byte) bool {
	lenHb := int64(len(hb))
	if c.BitsLeft() < lenHb*8 {
		return false
	}
	bs := c.PeekBytes(lenHb)
	return bytes.Equal(hb, bs)
}

// PeekFindByte number of bytes to next v
func (c *Common) PeekFindByte(v uint8, maxLen int64) int64 {
	peekBits, err := c.bitBuf.PeekFind(8, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFindByte", Size: 0, Pos: c.bitBuf.Pos})

	}
	return peekBits / 8
}

func (c *Common) BytesRange(firstBit int64, nBytes int64) []byte {
	bs, err := c.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesRange", Size: nBytes * 8, Pos: firstBit})
	}
	return bs
}

func (c *Common) BytesLen(nBytes int64) []byte {
	bs, err := c.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesLen", Size: nBytes * 8, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *Common) BitBufRange(firstBit int64, nBits int64) *bitbuf.Buffer {
	bs, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BitBufRange", Size: nBits, Pos: firstBit})
	}
	return bs
}

func (c *Common) BitBufLen(nBits int64) *bitbuf.Buffer {
	bs, err := c.bitBuf.BitBufLen(nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BitBufLen", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return bs
}

func (c *Common) Pos() int64           { return c.bitBuf.Pos }
func (c *Common) Len() int64           { return c.bitBuf.Len }
func (c *Common) End() bool            { return c.bitBuf.End() }
func (c *Common) BitsLeft() int64      { return c.bitBuf.BitsLeft() }
func (c *Common) ByteAlignBits() int64 { return c.bitBuf.ByteAlignBits() }
func (c *Common) BytePos() int64       { return c.bitBuf.BytePos() }

func (c *Common) SeekRel(deltaBits int64) int64 {
	pos, err := c.bitBuf.SeekRel(deltaBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SeekRel", Delta: deltaBits, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *Common) SeekAbs(pos int64) int64 {
	pos, err := c.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SeekAbs", Size: pos, Pos: c.bitBuf.Pos})
	}
	return pos
}

func (c *Common) UE(nBits int64, endian bitbuf.Endian) uint64 {
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

func (c *Common) FieldUE(name string, nBits int64, endian bitbuf.Endian) uint64 {
	return c.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := c.bitBuf.UE(nBits, endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldU" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (c *Common) SE(nBits int64, endian bitbuf.Endian) int64 {
	n, err := c.bitBuf.SE(nBits, endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SE", Size: nBits, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) FieldSE(name string, nBits int64, endian bitbuf.Endian) int64 {
	return c.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := c.bitBuf.SE(nBits, endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldS" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: c.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

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

func (c *Common) UTF8(nBytes int64) string {
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

func (c *Common) UFP64() float64 {
	f, err := c.bitBuf.UFP64()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP64", Size: 8, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *Common) FieldUFP64(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.UFP64(), ""
	})
}

func (c *Common) UFP32() float64 {
	f, err := c.bitBuf.UFP32()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP32", Size: 4, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *Common) FieldUFP32(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.UFP32(), ""
	})
}

func (c *Common) UFP16() float64 {
	f, err := c.bitBuf.UFP16()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP16", Size: 2, Pos: c.bitBuf.Pos})
	}
	return f
}

func (c *Common) FieldUFP16(name string) float64 {
	return c.FieldFloatFn(name, func() (float64, string) {
		return c.UFP16(), ""
	})
}

func (c *Common) Unary(s uint64) uint64 {
	n, err := c.bitBuf.Unary(s)
	if err != nil {
		panic(BitBufError{Err: err, Op: "Unary", Size: 1, Pos: c.bitBuf.Pos})
	}
	return n
}

func (c *Common) ZeroPadding(nBits int64) bool {
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

func (c *Common) AddChild(v *Value) {
	switch fv := c.current.V.(type) {
	case Struct:
		for _, ff := range fv {
			if ff.Name == v.Name {
				panic(fmt.Sprintf("%s already exist", v.Name))
			}
		}
		c.current.V = append(fv, v)
		return
	case Array:
		c.current.V = append(fv, v)
	}

}

func (c *Common) FieldArrayFn(name string, fn func()) {
	prev := c.current

	v := &Value{Name: name, V: Array{}}
	c.AddChild(v)
	c.current = v

	fn()

	var minMax Range
	for _, vf := range v.V.(Array) {
		minMax = RangeMinMax(minMax, vf.Range)
	}

	v.BitBuf = c.BitBufRange(minMax.Start, minMax.Stop-minMax.Start)
	v.Range = Range{Start: minMax.Start, Stop: minMax.Stop}

	c.current = prev
}

func (c *Common) FieldStructFn(name string, fn func()) {
	prev := c.current

	v := &Value{Name: name, V: Struct{}}
	c.AddChild(v)
	c.current = v

	fn()

	var minMax Range
	for _, vf := range v.V.(Struct) {
		minMax = RangeMinMax(minMax, vf.Range)
	}

	// TODO: find start/stop from Ranges instead? what if seekaround? concat bitbufs but want gaps? sort here, crash?
	v.BitBuf = c.BitBufRange(minMax.Start, minMax.Stop-minMax.Start)
	v.Range = Range{Start: minMax.Start, Stop: minMax.Stop}

	c.current = prev
}

func (c *Common) fieldDecoder(name string, bitBuf *bitbuf.Buffer, v interface{}) *Common {
	d := &Common{
		bitBuf: bitBuf,
		// TODO: rename current to value?
		current: &Value{
			Name: name,
			V:    v,
		},
		registry: c.registry,
	}
	// TODO: refactor
	if c.current != nil {
		c.AddChild(d.current)
	}
	return d
}

func (c *Common) FieldArray2(name string) *Common {
	return c.fieldDecoder(name, c.bitBuf, Array{})
}

func (c *Common) FieldArrayFn2(name string, fn func(d *Common)) *Common {
	d := c.FieldArray2(name)
	fn(d)
	return d
}

func (c *Common) FieldStruct2(name string) *Common {
	return c.fieldDecoder(name, c.bitBuf, Struct{})
}

func (c *Common) FieldStructFn2(name string, fn func(d *Common)) *Common {
	d := c.FieldStruct2(name)
	fn(d)
	return d
}

func (c *Common) FieldStructBitBuf(name string, bitBuf *bitbuf.Buffer) *Common {
	return c.fieldDecoder(name, bitBuf, Struct{})
}

func (c *Common) FieldStructBitBufFn(name string, bitBuf *bitbuf.Buffer, fn func(d *Common)) *Common {
	d := c.FieldStructBitBuf(name, bitBuf)
	fn(d)
	return d
}

func (c *Common) FieldRangeFn(name string, firstBit int64, nBits int64, fn func() Value) Value {
	v := fn()
	v.Name = name
	v.BitBuf = c.BitBufRange(firstBit, nBits)
	v.Range = Range{Start: firstBit, Stop: firstBit + nBits}
	c.AddChild(&v)

	return v
}

func (c *Common) FieldFn(name string, fn func() Value) Value {
	start := c.bitBuf.Pos
	v := fn()
	stop := c.bitBuf.Pos
	v.Name = name
	v.BitBuf = c.BitBufRange(start, stop-start)
	v.Range = Range{Start: start, Stop: stop}
	c.AddChild(&v)

	return v
}

// TODO: remove
func (c *Common) FieldNoneFn(name string, fn func()) {
	c.FieldStructFn(name, func() {
		fn()
	})
}

func (c *Common) FieldBoolFn(name string, fn func() (bool, string)) bool {
	return c.FieldFn(name, func() Value {
		b, d := fn()
		return Value{V: b, Symbol: d}
	}).V.(bool)
}

func (c *Common) FieldUFn(name string, fn func() (uint64, DisplayFormat, string)) uint64 {
	return c.FieldFn(name, func() Value {
		u, fmt, d := fn()
		return Value{V: u, DisplayFormat: fmt, Symbol: d}
	}).V.(uint64)
}

func (c *Common) FieldSFn(name string, fn func() (int64, DisplayFormat, string)) int64 {
	return c.FieldFn(name, func() Value {
		s, fmt, d := fn()
		return Value{V: s, DisplayFormat: fmt, Symbol: d}
	}).V.(int64)
}

func (c *Common) FieldFloatFn(name string, fn func() (float64, string)) float64 {
	return c.FieldFn(name, func() Value {
		f, d := fn()
		return Value{V: f, Symbol: d}
	}).V.(float64)
}

func (c *Common) FieldStrFn(name string, fn func() (string, string)) string {
	return c.FieldFn(name, func() Value {
		str, disp := fn()
		return Value{V: str, Symbol: disp}
	}).V.(string)
}

func (c *Common) FieldBytesFn(name string, firstBit int64, nBits int64, fn func() ([]byte, string)) []byte {
	return c.FieldRangeFn(name, firstBit, nBits, func() Value {
		bs, disp := fn()
		return Value{V: bs, Symbol: disp}
	}).V.([]byte)
}

func (c *Common) FieldBitBufFn(name string, firstBit int64, nBits int64, fn func() (*bitbuf.Buffer, string)) *bitbuf.Buffer {
	return c.FieldRangeFn(name, firstBit, nBits, func() Value {
		bb, disp := fn()
		return Value{V: bb, Symbol: disp}
	}).BitBuf
}

func (c *Common) FieldStringMapFn(name string, sm map[uint64]string, def string, fn func() uint64) (uint64, bool) {
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

func (c *Common) FieldValidateUFn(name string, v uint64, fn func() uint64) {
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
func (c *Common) FieldBytesLen(name string, nBytes int64) []byte {
	return c.FieldBytesFn(name, c.bitBuf.Pos, nBytes*8, func() ([]byte, string) {
		bs, err := c.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesLen", Size: nBytes * 8, Pos: c.bitBuf.Pos})
		}
		return bs, ""
	})
}

func (c *Common) FieldBytesRange(name string, firstBit int64, nBytes int64) []byte {
	return c.FieldBytesFn(name, firstBit, nBytes*8, func() ([]byte, string) {
		bs, err := c.bitBuf.BytesRange(firstBit, nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesRange", Size: nBytes * 8, Pos: firstBit})
		}
		return bs, ""
	})
}

func (c *Common) FieldUTF8(name string, nBytes int64) string {
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

func (c *Common) FieldValidateZeroPadding(name string, nBits int64) {
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

func (c *Common) ValidateAtLeastBitsLeft(nBits int64) {
	bl := c.bitBuf.BitsLeft()
	if bl < nBits {
		// TODO:
		panic(ValidateError{Reason: "not enough bits left", Pos: c.bitBuf.Pos})
	}
}

func (c *Common) ValidateAtLeastBytesLeft(nBytes int64) {
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
func (c *Common) SubLenFn(nBits int64, fn func()) {
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

func (c *Common) SubRangeFn(firstBit int64, nBits int64, fn func()) {
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
func (c *Common) FieldTryDecode(name string, forceFormats []*Format) (*Value, Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(c.bitBuf.Pos, c.BitsLeft())
	if err != nil {
		// TODO: can't happen?
		panic(BitBufError{Err: err, Op: "FieldDecode", Size: c.BitsLeft(), Pos: c.bitBuf.Pos})
	}

	v, fLen, d, errs := c.registry.Probe(c, name, Range{Start: c.bitBuf.Pos, Stop: c.bitBuf.Pos}, bb, forceFormats)
	if v == nil || v.Errors() != nil {
		return nil, nil, errs
	}

	// TODO: bitbuf len shorten!
	c.AddChild(v)
	_, err = c.bitBuf.SeekRel(int64(fLen))
	if err != nil {
		panic(err)
	}

	return v, d, errs
}

// TODO: FieldTryDecode? just TryDecode?
func (c *Common) FieldDecodeLen(name string, nBits int64, forceFormats []*Format) (*Value, Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(c.bitBuf.Pos, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeLen", Size: nBits, Pos: c.bitBuf.Pos})
	}

	v, _, d, errs := c.registry.Probe(c, name, Range{Start: c.bitBuf.Pos, Stop: c.bitBuf.Pos + nBits}, bb, forceFormats)
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

	return v, d, errs
}

// TODO: return decooder?
func (c *Common) FieldTryDecodeRange(name string, firstBit int64, nBits int64, forceFormats []*Format) (*Value, Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: c.bitBuf.Pos})
	}

	v, _, d, errs := c.registry.Probe(c, name, Range{Start: firstBit, Stop: firstBit + nBits}, bb, forceFormats)
	if v != nil {
		c.AddChild(v)
	}

	return v, d, errs
}

// TODO: return decooder?
func (c *Common) FieldDecodeRange(name string, firstBit int64, nBits int64, forceFormats []*Format) (*Value, Decoder, []error) {
	bb, err := c.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: c.bitBuf.Pos})
	}

	v, _, d, errs := c.registry.Probe(c, name, Range{Start: firstBit, Stop: firstBit + nBits}, bb, forceFormats)
	if v != nil {
		c.AddChild(v)
	} else {
		c.FieldRangeFn(name, firstBit, nBits, func() Value { return Value{} })
	}

	return v, d, errs
}

// TODO: list of ranges?
func (c *Common) FieldDecodeBitBuf(name string, firstBit int64, nBits int64, bb *bitbuf.Buffer, forceFormats []*Format) (*Value, Decoder, []error) {
	f, _, d, errs := c.registry.Probe(c, name, Range{Start: firstBit, Stop: nBits}, bb, forceFormats)
	if f != nil {
		c.AddChild(f)
	} else {
		c.FieldBitBufFn(name, firstBit, nBits, func() (*bitbuf.Buffer, string) {
			return bb, ""
		})
	}

	return f, d, errs
}

func (c *Common) FieldBitBufRange(name string, firstBit int64, nBits int64) *bitbuf.Buffer {
	return c.FieldBitBufFn(name, firstBit, nBits, func() (*bitbuf.Buffer, string) {
		return c.BitBufRange(firstBit, nBits), ""
	})
}

func (c *Common) FieldBitBufLen(name string, nBits int64) *bitbuf.Buffer {
	return c.FieldBitBufFn(name, c.bitBuf.Pos, nBits, func() (*bitbuf.Buffer, string) {
		return c.BitBufLen(nBits), ""
	})
}

func (c *Common) FieldZlib(name string, firsBit int64, nBits int64, b []byte, formats []*Format) (*Value, Decoder, []error) {
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
func (c *Common) FieldZlibLen(name string, nBytes int64, formats []*Format) (*Value, Decoder, []error) {
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
