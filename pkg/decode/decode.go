package decode

//go:generate sh -c "cat decode_gen.go.tmpl | go run ../../_dev/tmpl.go | gofmt > decode_gen.go"

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

type Endian bitbuf.Endian

var (
	// BigEndian byte order
	BigEndian Endian = Endian(bitbuf.BigEndian)
	// LittleEndian byte order
	LittleEndian Endian = Endian(bitbuf.LittleEndian)
)

type D struct {
	Endian Endian

	bitBuf   *bitbuf.Buffer
	value    *Value
	registry *Registry
}

// Probe probes all probeable formats and turns first found Decoder and all other decoder errors
func Probe(name string, bb *bitbuf.Buffer, formats []*Format) (*Value, interface{}, []error) {
	var forceOne = len(formats) == 1

	// TODO: order..

	startPos := bb.Pos

	var errs []error
	for _, f := range formats {
		cbb := bb.Copy()

		// TODO: how to pass regsiters? do later? current field?

		d := (&D{Endian: BigEndian, bitBuf: cbb}).FieldStructBitBuf(name, cbb)
		d.value.Desc = f.Name
		d.value.BitBuf = cbb
		decodeErr, dv := d.SafeDecodeFn(f.DecodeFn)
		if decodeErr != nil {
			d.value.Error = decodeErr

			errs = append(errs, decodeErr)
			if !forceOne {
				continue
			}
		}

		// TODO: nicer
		d.value.Range = Range{Start: startPos, Stop: cbb.Pos}

		if d.value.Parent == nil {
			d.value.Sort()
		}

		// TODO: wrong keep track of largest?
		_ = cbb.TruncateRel(0)

		return d.value, dv, errs
	}

	return nil, nil, errs
}

func (d *D) SafeDecodeFn(fn func(d *D) interface{}) (error, interface{}) {
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

		return nil, fn(d)
	}()

	return decodeErr, dv
}

func (d *D) PeekBits(nBits int64) uint64 {
	n, err := d.bitBuf.PeekBits(nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBits", Size: nBits, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) PeekBytes(nBytes int64) []byte {
	bs, err := d.bitBuf.PeekBytes(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekBytes", Size: nBytes * 8, Pos: d.bitBuf.Pos})
	}
	return bs
}

func (d *D) PeekFind(nBits int64, v uint8, maxLen int64) int64 {
	peekBits, err := d.bitBuf.PeekFind(nBits, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFind", Size: 0, Pos: d.bitBuf.Pos})
	}
	return peekBits
}

func (d *D) TryHasBytes(hb []byte) bool {
	lenHb := int64(len(hb))
	if d.BitsLeft() < lenHb*8 {
		return false
	}
	bs := d.PeekBytes(lenHb)
	return bytes.Equal(hb, bs)
}

// PeekFindByte number of bytes to next v
func (d *D) PeekFindByte(v uint8, maxLen int64) int64 {
	peekBits, err := d.bitBuf.PeekFind(8, v, maxLen)
	if err != nil {
		panic(BitBufError{Err: err, Op: "PeekFindByte", Size: 0, Pos: d.bitBuf.Pos})

	}
	return peekBits / 8
}

func (d *D) BytesRange(firstBit int64, nBytes int64) []byte {
	bs, err := d.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesRange", Size: nBytes * 8, Pos: firstBit})
	}
	return bs
}

func (d *D) BytesLen(nBytes int64) []byte {
	bs, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BytesLen", Size: nBytes * 8, Pos: d.bitBuf.Pos})
	}
	return bs
}

func (d *D) BitBufRange(firstBit int64, nBits int64) *bitbuf.Buffer {
	bs, err := d.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BitBufRange", Size: nBits, Pos: firstBit})
	}
	return bs
}

func (d *D) BitBufLen(nBits int64) *bitbuf.Buffer {
	bs, err := d.bitBuf.BitBufLen(nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "BitBufLen", Size: nBits, Pos: d.bitBuf.Pos})
	}
	return bs
}

func (d *D) Pos() int64           { return d.bitBuf.Pos }
func (d *D) Len() int64           { return d.bitBuf.Len }
func (d *D) End() bool            { return d.bitBuf.End() }
func (d *D) BitsLeft() int64      { return d.bitBuf.BitsLeft() }
func (d *D) ByteAlignBits() int64 { return d.bitBuf.ByteAlignBits() }
func (d *D) BytePos() int64       { return d.bitBuf.BytePos() }

func (d *D) SeekRel(deltaBits int64) int64 {
	pos, err := d.bitBuf.SeekRel(deltaBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SeekRel", Delta: deltaBits, Pos: d.bitBuf.Pos})
	}
	return pos
}

func (d *D) SeekAbs(pos int64) int64 {
	pos, err := d.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SeekAbs", Size: pos, Pos: d.bitBuf.Pos})
	}
	return pos
}

func (d *D) UE(nBits int64, endian Endian) uint64 {
	n, err := d.bitBuf.UE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(BitBufError{Err: err, Op: "UE", Size: nBits, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) Bool() bool {
	b, err := d.bitBuf.Bool()
	if err != nil {
		panic(BitBufError{Err: err, Op: "Bool", Size: 1, Pos: d.bitBuf.Pos})
	}
	return b
}

func (d *D) FieldBool(name string) bool {
	return d.FieldBoolFn(name, func() (bool, string) {
		b, err := d.bitBuf.Bool()
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBool", Size: 1, Pos: d.bitBuf.Pos})
		}
		return b, ""
	})
}

func (d *D) FieldUE(name string, nBits int64, endian Endian) uint64 {
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.bitBuf.UE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldU" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) SE(nBits int64, endian Endian) int64 {
	n, err := d.bitBuf.SE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(BitBufError{Err: err, Op: "SE", Size: nBits, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) FieldSE(name string, nBits int64, endian Endian) int64 {
	return d.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := d.bitBuf.SE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldS" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) F32E(endian bitbuf.Endian) float64 {
	f, err := d.bitBuf.F32E(endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "F32", Size: 32, Pos: d.bitBuf.Pos})
	}
	return float64(f)
}

func (d *D) F32() float64   { return d.F32E(bitbuf.BigEndian) }
func (d *D) F32BE() float64 { return d.F32E(bitbuf.BigEndian) }
func (d *D) F32LE() float64 { return d.F32E(bitbuf.LittleEndian) }

func (d *D) FieldF32E(name string, endian bitbuf.Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		f, err := d.bitBuf.F32E(endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "F32", Size: 32, Pos: d.bitBuf.Pos})
		}
		return float64(f), ""
	})
}

func (d *D) FieldF32(name string) float64   { return d.FieldF32E(name, bitbuf.BigEndian) }
func (d *D) FieldF32BE(name string) float64 { return d.FieldF32E(name, bitbuf.BigEndian) }
func (d *D) FieldF32LE(name string) float64 { return d.FieldF32E(name, bitbuf.LittleEndian) }

func (d *D) F64E(endian bitbuf.Endian) float64 {
	f, err := d.bitBuf.F64E(endian)
	if err != nil {
		panic(BitBufError{Err: err, Op: "F64", Size: 64, Pos: d.bitBuf.Pos})
	}
	return float64(f)
}

func (d *D) F64() float64   { return d.F64E(bitbuf.BigEndian) }
func (d *D) F64BE() float64 { return d.F64E(bitbuf.BigEndian) }
func (d *D) F64LE() float64 { return d.F64E(bitbuf.LittleEndian) }

func (d *D) FieldF64E(name string, endian bitbuf.Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		f, err := d.bitBuf.F64E(endian)
		if err != nil {
			panic(BitBufError{Err: err, Op: "F64", Size: 64, Pos: d.bitBuf.Pos})
		}
		return float64(f), ""
	})
}

func (d *D) FieldF64(name string) float64   { return d.FieldF64E(name, bitbuf.BigEndian) }
func (d *D) FieldF64BE(name string) float64 { return d.FieldF64E(name, bitbuf.BigEndian) }
func (d *D) FieldF64LE(name string) float64 { return d.FieldF64E(name, bitbuf.LittleEndian) }

func (d *D) UTF8(nBytes int64) string {
	s, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(BitBufError{Err: err, Op: "UTF8", Size: nBytes * 8, Pos: d.bitBuf.Pos})
	}
	return string(s)
}

func (d *D) FP64() float64 {
	f, err := d.bitBuf.FP64()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP64", Size: 8, Pos: d.bitBuf.Pos})
	}
	return f
}

func (d *D) FieldFP64(name string) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		return d.FP64(), ""
	})
}

func (d *D) FP32() float64 {
	f, err := d.bitBuf.FP32()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP32", Size: 4, Pos: d.bitBuf.Pos})
	}
	return f
}

func (d *D) FieldFP32(name string) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		return d.FP32(), ""
	})
}

func (d *D) FP16() float64 {
	f, err := d.bitBuf.FP16()
	if err != nil {
		panic(BitBufError{Err: err, Op: "FP16", Size: 2, Pos: d.bitBuf.Pos})
	}
	return f
}

func (d *D) FieldFP16(name string) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		return d.FP16(), ""
	})
}

func (d *D) UFP64() float64 {
	f, err := d.bitBuf.UFP64()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP64", Size: 8, Pos: d.bitBuf.Pos})
	}
	return f
}

func (d *D) FieldUFP64(name string) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		return d.UFP64(), ""
	})
}

func (d *D) UFP32() float64 {
	f, err := d.bitBuf.UFP32()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP32", Size: 4, Pos: d.bitBuf.Pos})
	}
	return f
}

func (d *D) FieldUFP32(name string) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		return d.UFP32(), ""
	})
}

func (d *D) UFP16() float64 {
	f, err := d.bitBuf.UFP16()
	if err != nil {
		panic(BitBufError{Err: err, Op: "UFP16", Size: 2, Pos: d.bitBuf.Pos})
	}
	return f
}

func (d *D) FieldUFP16(name string) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		return d.UFP16(), ""
	})
}

func (d *D) Unary(s uint64) uint64 {
	n, err := d.bitBuf.Unary(s)
	if err != nil {
		panic(BitBufError{Err: err, Op: "Unary", Size: 1, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) ZeroPadding(nBits int64) bool {
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
		n, err := d.bitBuf.Bits(rbits)
		if err != nil {
			panic(BitBufError{Err: err, Op: "ZeroPadding", Size: rbits, Pos: d.bitBuf.Pos})
		}
		isZero = isZero && n == 0
		left -= rbits
	}
	return isZero
}

func (d *D) AddChild(v *Value) {
	v.Parent = d.value

	switch fv := d.value.V.(type) {
	case Struct:
		for _, ff := range fv {
			if ff.Name == v.Name {
				panic(fmt.Sprintf("%s already exist in struct %s", v.Name, d.value.Name))
			}
		}
		d.value.V = append(fv, v)
		return
	case Array:
		d.value.V = append(fv, v)
	}

}

func (d *D) fieldDecoder(name string, bitBuf *bitbuf.Buffer, v interface{}) *D {
	r := Range{}
	if d.bitBuf != nil {
		r = Range{Start: d.bitBuf.Pos, Stop: d.bitBuf.Pos}
	}

	cd := &D{
		Endian: d.Endian,

		bitBuf: bitBuf,
		// TODO: rename current to value?
		value: &Value{
			Name:  name,
			V:     v,
			Range: r,
		},
		registry: d.registry,
	}

	// TODO: find start/stop from Ranges instead? what if seekaround? concat bitbufs but want gaps? sort here, crash?

	// TODO: refactor
	if d.value != nil {
		d.AddChild(cd.value)
	}
	return cd
}

func (d *D) FieldArray(name string) *D {
	return d.fieldDecoder(name, d.bitBuf, Array{})
}

func (d *D) FieldArrayFn(name string, fn func(d *D)) *D {
	cd := d.FieldArray(name)
	fn(cd)
	return cd
}

func (d *D) FieldStruct(name string) *D {
	return d.fieldDecoder(name, d.bitBuf, Struct{})
}

func (d *D) FieldStructFn(name string, fn func(d *D)) *D {
	cd := d.FieldStruct(name)
	fn(cd)
	return cd
}

func (d *D) FieldStructBitBuf(name string, bitBuf *bitbuf.Buffer) *D {
	return d.fieldDecoder(name, bitBuf, Struct{})
}

func (d *D) FieldStructBitBufFn(name string, bitBuf *bitbuf.Buffer, fn func(d *D)) *D {
	cd := d.FieldStructBitBuf(name, bitBuf)
	fn(cd)
	return cd
}

func (d *D) FieldRangeFn(name string, firstBit int64, nBits int64, fn func() Value) Value {
	v := fn()
	v.Name = name
	//v.BitBuf = d.BitBufRange(firstBit, nBits)
	v.Range = Range{Start: firstBit, Stop: firstBit + nBits}
	d.AddChild(&v)

	return v
}

func (d *D) FieldFn(name string, fn func() Value) Value {
	start := d.bitBuf.Pos
	v := fn()
	stop := d.bitBuf.Pos
	v.Name = name
	//v.BitBuf = d.BitBufRange(start, stop-start)
	v.Range = Range{Start: start, Stop: stop}
	d.AddChild(&v)

	return v
}

func (d *D) FieldBoolFn(name string, fn func() (bool, string)) bool {
	return d.FieldFn(name, func() Value {
		b, d := fn()
		return Value{V: b, Symbol: d}
	}).V.(bool)
}

func (d *D) FieldUFn(name string, fn func() (uint64, DisplayFormat, string)) uint64 {
	return d.FieldFn(name, func() Value {
		u, fmt, d := fn()
		return Value{V: u, DisplayFormat: fmt, Symbol: d}
	}).V.(uint64)
}

func (d *D) FieldSFn(name string, fn func() (int64, DisplayFormat, string)) int64 {
	return d.FieldFn(name, func() Value {
		s, fmt, d := fn()
		return Value{V: s, DisplayFormat: fmt, Symbol: d}
	}).V.(int64)
}

func (d *D) FieldFloatFn(name string, fn func() (float64, string)) float64 {
	return d.FieldFn(name, func() Value {
		f, d := fn()
		return Value{V: f, Symbol: d}
	}).V.(float64)
}

func (d *D) FieldStrFn(name string, fn func() (string, string)) string {
	return d.FieldFn(name, func() Value {
		str, disp := fn()
		return Value{V: str, Symbol: disp}
	}).V.(string)
}

func (d *D) FieldBytesFn(name string, firstBit int64, nBits int64, fn func() ([]byte, string)) []byte {
	return d.FieldRangeFn(name, firstBit, nBits, func() Value {
		bs, disp := fn()
		return Value{V: bs, Symbol: disp}
	}).V.([]byte)
}

func (d *D) FieldBitBufFn(name string, firstBit int64, nBits int64, fn func() (*bitbuf.Buffer, string)) *bitbuf.Buffer {
	return d.FieldRangeFn(name, firstBit, nBits, func() Value {
		bb, disp := fn()
		return Value{V: bb, Symbol: disp}
	}).BitBuf
}

func (d *D) FieldStringMapFn(name string, sm map[uint64]string, def string, fn func() uint64) (uint64, bool) {
	var ok bool
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n := fn()
		var d string
		d, ok = sm[n]
		if !ok {
			d = def
		}
		return n, NumberDecimal, d
	}), ok
}

func (d *D) FieldValidateUFn(name string, v uint64, fn func() uint64) {
	pos := d.bitBuf.Pos
	n := d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
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
func (d *D) FieldBytesLen(name string, nBytes int64) []byte {
	return d.FieldBytesFn(name, d.bitBuf.Pos, nBytes*8, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesLen", Size: nBytes * 8, Pos: d.bitBuf.Pos})
		}
		return bs, ""
	})
}

func (d *D) FieldBytesRange(name string, firstBit int64, nBytes int64) []byte {
	return d.FieldBytesFn(name, firstBit, nBytes*8, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesRange(firstBit, nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesRange", Size: nBytes * 8, Pos: firstBit})
		}
		return bs, ""
	})
}

func (d *D) FieldUTF8(name string, nBytes int64) string {
	return d.FieldStrFn(name, func() (string, string) {
		str, err := d.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldUTF8", Size: nBytes * 8, Pos: d.bitBuf.Pos})
		}
		return str, ""
	})
}

func (d *D) FieldValidateStringFn(name string, v string, fn func() string) {
	pos := d.bitBuf.Pos
	s := d.FieldStrFn(name, func() (string, string) {
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

func (d *D) FieldValidateString(name string, v string) {
	pos := d.bitBuf.Pos
	s := d.FieldStrFn(name, func() (string, string) {
		nBytes := int64(len(v))
		str, err := d.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldValidateString", Size: nBytes * 8, Pos: d.bitBuf.Pos})
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

func (d *D) FieldValidateZeroPadding(name string, nBits int64) {
	pos := d.bitBuf.Pos
	var isZero bool
	d.FieldFn(name, func() Value {
		isZero = d.ZeroPadding(nBits)
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

func (d *D) ValidateAtLeastBitsLeft(nBits int64) {
	bl := d.bitBuf.BitsLeft()
	if bl < nBits {
		// TODO:
		panic(ValidateError{Reason: fmt.Sprintf("expected bits left %d, found %d", nBits, bl), Pos: d.bitBuf.Pos})
	}
}

func (d *D) ValidateAtLeastBytesLeft(nBytes int64) {
	bl := d.bitBuf.BitsLeft()
	if bl < nBytes*8 {
		// TODO:
		panic(ValidateError{Reason: fmt.Sprintf("expected bytes left %d, found %d bits", nBytes, bl), Pos: d.bitBuf.Pos})
	}
}

// Invalid stops decode with a reason
func (d *D) Invalid(reason string) {
	panic(ValidateError{Reason: reason, Pos: d.bitBuf.Pos})
}

// TODO: rename?
func (d *D) SubLenFn(nBits int64, fn func()) {
	prevBb := d.bitBuf

	bb, err := d.bitBuf.BitBufRange(0, d.bitBuf.Pos+nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SubLen", Size: nBits, Pos: d.bitBuf.Pos})
	}
	_, err = bb.SeekAbs(d.bitBuf.Pos)
	if err != nil {
		panic(err)
	}
	d.bitBuf = bb

	fn()

	bitsLeft := nBits - (d.bitBuf.Pos - prevBb.Pos)
	d.SeekRel(int64(bitsLeft))

	prevBb.Pos = d.bitBuf.Pos
	d.bitBuf = prevBb
}

func (d *D) SubRangeFn(firstBit int64, nBits int64, fn func()) {
	prevBb := d.bitBuf

	bb, err := d.bitBuf.BitBufRange(0, firstBit+nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "SubRangeFn", Size: nBits, Pos: firstBit})
	}
	_, err = bb.SeekAbs(firstBit)
	if err != nil {
		panic(err)
	}
	d.bitBuf = bb

	fn()

	d.bitBuf = prevBb
}

// TODO: TryDecode?
func (d *D) FieldTryDecode(name string, formats []*Format) (*Value, interface{}, []error) {
	bb, err := d.bitBuf.BitBufRange(d.bitBuf.Pos, d.BitsLeft())
	if err != nil {
		// TODO: can't happen?
		panic(BitBufError{Err: err, Op: "FieldTryDecode", Size: d.BitsLeft(), Pos: d.bitBuf.Pos})
	}

	v, dv, errs := Probe(name, bb, formats)
	if v == nil || v.Errors() != nil {
		return nil, nil, errs
	}

	v.BitBuf = nil
	//v.Range = Range{Start: v.Range.Start + d.Pos(), Stop: v.Range.Stop + d.Pos()}

	// TODO: bitbuf len shorten!
	d.AddChild(v)
	_, err = d.bitBuf.SeekRel(int64(v.Range.Length()))
	if err != nil {
		panic(err)
	}

	return v, dv, errs
}

// TODO: FieldTryDecode? just TryDecode?
func (d *D) FieldDecodeLen(name string, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	bb, err := d.bitBuf.BitBufRange(d.bitBuf.Pos, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeLen", Size: nBits, Pos: d.bitBuf.Pos})
	}

	v, dv, errs := Probe(name, bb, formats)
	if v != nil {
		v.BitBuf = nil
		d.AddChild(v)
	} else {
		// TODO: decoder unknown
		d.FieldRangeFn(name, d.bitBuf.Pos, nBits, func() Value { return Value{} })
	}

	// TODO: nBits - fLen gap?

	_, err = d.bitBuf.SeekRel(int64(nBits))
	if err != nil {
		panic(err)
	}

	return v, dv, errs
}

// TODO: return decooder?
func (d *D) FieldTryDecodeRange(name string, firstBit int64, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	bb, err := d.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: d.bitBuf.Pos})
	}

	v, dv, errs := Probe(name, bb, formats)
	if v != nil {
		v.BitBuf = nil
		//v.Range = Range{Start: v.Range.Start + firstBit, Stop: v.Range.Stop + firstBit}
		d.AddChild(v)
	}

	return v, dv, errs
}

// TODO: return decooder?
func (d *D) FieldDecodeRange(name string, firstBit int64, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	bb, err := d.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(BitBufError{Err: err, Op: "FieldDecodeRange", Size: nBits, Pos: d.bitBuf.Pos})
	}

	v, dv, errs := Probe(name, bb, formats)
	if v != nil {
		v.BitBuf = nil
		d.AddChild(v)
	} else {
		d.FieldRangeFn(name, firstBit, nBits, func() Value { return Value{} })
	}

	return v, dv, errs
}

// TODO: list of ranges?
func (d *D) FieldDecodeBitBuf(name string, firstBit int64, nBits int64, bb *bitbuf.Buffer, formats []*Format) (*Value, interface{}, []error) {
	f, dv, errs := Probe(name, bb, formats)
	if f != nil {
		d.AddChild(f)
	} else {
		d.FieldBitBufFn(name, firstBit, nBits, func() (*bitbuf.Buffer, string) {
			return bb, ""
		})
	}

	return f, dv, errs
}

func (d *D) FieldBitBufRange(name string, firstBit int64, nBits int64) *bitbuf.Buffer {
	return d.FieldBitBufFn(name, firstBit, nBits, func() (*bitbuf.Buffer, string) {
		return d.BitBufRange(firstBit, nBits), ""
	})
}

func (d *D) FieldBitBufLen(name string, nBits int64) *bitbuf.Buffer {
	return d.FieldBitBufFn(name, d.bitBuf.Pos, nBits, func() (*bitbuf.Buffer, string) {
		return d.BitBufLen(nBits), ""
	})
}

func (d *D) FieldZlib(name string, firsBit int64, nBits int64, b []byte, formats []*Format) (*Value, interface{}, []error) {
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

	return d.FieldDecodeBitBuf(name, firsBit, nBits, zbb, formats)
}

// TODO: range?
func (d *D) FieldZlibLen(name string, nBytes int64, formats []*Format) (*Value, interface{}, []error) {
	firstBit := d.bitBuf.Pos
	zr, err := zlib.NewReader(bytes.NewReader(d.BytesLen(nBytes)))
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

	return d.FieldDecodeBitBuf(name, firstBit, firstBit+nBytes*8, zbb, formats)
}
