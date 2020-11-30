package decode

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"fq/internal/ranges"
	"fq/pkg/bitio"
	"io/ioutil"
	"runtime"
)

type DecodeError struct {
	Err        error
	PanicStack string
}

func (de *DecodeError) Error() string { return de.Err.Error() }
func (de *DecodeError) Unwrap() error { return de.Err }

type ReadError struct {
	Err   error
	Name  string
	Op    string
	Size  int64
	Delta int64
	Pos   int64
}

func (e ReadError) Error() string {
	var prefix string
	if e.Name != "" {
		prefix = e.Name + ": " + e.Op
	} else {
		prefix = e.Op
	}

	return fmt.Sprintf("%s: failed at position %s (size %s delta %s): %s",
		prefix, Bits(e.Pos).StringByteBits(16), Bits(e.Size).StringByteBits(10), Bits(e.Delta).StringByteBits(10), e.Err)
}
func (e ReadError) Unwrap() error { return e.Err }

type ValidateError struct {
	Reason string
	Pos    int64
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("failed to validate at position %s: %s", Bits(e.Pos).StringByteBits(16), e.Reason)
}

type Endian bitio.Endian

var (
	// BigEndian byte order
	BigEndian Endian = Endian(bitio.BigEndian)
	// LittleEndian byte order
	LittleEndian Endian = Endian(bitio.LittleEndian)
)

type probeOptions struct {
	isRoot    bool
	relBitBuf *bitio.Buffer
	relStart  int64
}

// Probe probes all probeable formats and turns first found Decoder and all other decoder errors
func Probe(name string, bb *bitio.Buffer, formats []*Format) (*Value, interface{}, []error) {
	return probe(name, bb, formats, probeOptions{isRoot: true})
}

func probe(name string, bb *bitio.Buffer, formats []*Format, opts probeOptions) (*Value, interface{}, []error) {
	var forceOne = len(formats) == 1

	var errs []error
	for _, f := range formats {
		d := NewDecoder(name, f.Name, bb, opts.isRoot)

		decodeErr, dv := d.SafeDecodeFn(f.DecodeFn)
		if decodeErr != nil {
			d.Value.Error = decodeErr

			errs = append(errs, decodeErr)
			if !forceOne {
				continue
			}
		}

		var maxRange ranges.Range
		d.Value.WalkPreOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
			if v.IsRoot {
				return ErrWalkSkip
			}

			maxRange = ranges.MinMax(maxRange, v.Range)
			v.Range.Start += opts.relStart
			if opts.relBitBuf != nil {
				v.BitBuf = opts.relBitBuf
			}
			return nil
		})

		d.Value.Range = ranges.Range{Start: opts.relStart, Len: maxRange.Len}

		if opts.isRoot {
			d.FillGaps("unknown")

			// sort and set ranges for struct and arrays
			d.Value.postProcess()
		}

		return d.Value, dv, errs
	}

	return nil, nil, errs
}

type D struct {
	Endian Endian
	Value  *Value

	bitBuf   *bitio.Buffer
	registry *Registry
}

// TODO: new struct decoder?
func NewDecoder(name string, description string, bb *bitio.Buffer, isRoot bool) *D {
	cbb := bb.Copy()

	d := (&D{Endian: BigEndian, bitBuf: cbb}).FieldStructBitBuf(name, cbb)
	d.Value.Desc = description
	d.Value.BitBuf = cbb
	d.Value.IsRoot = isRoot

	return d
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
				case ReadError:
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

func (d *D) FillGaps(namePrefix string) {
	// TODO: d.Value is array?
	var valueRanges []ranges.Range
	d.Value.WalkPreOrder(func(iv *Value, rootV *Value, depth int, rootDepth int) error {
		if iv.BitBuf != d.Value.BitBuf && iv.IsRoot {
			return ErrWalkSkip
		}
		switch iv.V.(type) {
		case Struct, Array:
		default:
			valueRanges = append(valueRanges, iv.Range)
		}
		return nil
	})

	gaps := ranges.Gaps(ranges.Range{Start: 0, Len: d.Len()}, valueRanges)
	for i, gap := range gaps {
		d.FieldBitBufRange(
			fmt.Sprintf("%s%d", namePrefix, i), gap.Start, gap.Len,
		)
	}
}

// Invalid stops decode with a reason
func (d *D) Invalid(reason string) {
	panic(ValidateError{Reason: reason, Pos: d.Pos()})
}

func (d *D) PeekBits(nBits int) uint64 {
	n, err := d.bitBuf.PeekBits(nBits)
	if err != nil {
		panic(ReadError{Err: err, Op: "PeekBits", Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) PeekBytes(nBytes int) []byte {
	bs, err := d.bitBuf.PeekBytes(nBytes)
	if err != nil {
		panic(ReadError{Err: err, Op: "PeekBytes", Size: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
}

func (d *D) PeekFind(nBits int, v uint8, maxLen int64) int64 {
	peekBits, err := d.bitBuf.PeekFind(nBits, v, maxLen)
	if err != nil {
		panic(ReadError{Err: err, Op: "PeekFind", Size: 0, Pos: d.Pos()})
	}
	return peekBits
}

func (d *D) TryHasBytes(hb []byte) bool {
	lenHb := len(hb)
	if d.BitsLeft() < int64(lenHb*8) {
		return false
	}
	bs := d.PeekBytes(lenHb)
	return bytes.Equal(hb, bs)
}

// PeekFindByte number of bytes to next v
func (d *D) PeekFindByte(v uint8, maxLen int64) int64 {
	peekBits, err := d.bitBuf.PeekFind(8, v, maxLen)
	if err != nil {
		panic(ReadError{Err: err, Op: "PeekFindByte", Size: 0, Pos: d.Pos()})

	}
	return peekBits / 8
}

func (d *D) BytesRange(firstBit int64, nBytes int) []byte {
	bs, err := d.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(ReadError{Err: err, Op: "BytesRange", Size: int64(nBytes) * 8, Pos: firstBit})
	}
	return bs
}

func (d *D) BytesLen(nBytes int) []byte {
	bs, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(ReadError{Err: err, Op: "BytesLen", Size: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
}

// TODO: rename/remove BitBuf name?
func (d *D) BitBufRange(firstBit int64, nBits int64) *bitio.Buffer {
	bb, err := d.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(ReadError{Err: err, Op: "BitBufRange", Size: nBits, Pos: firstBit})
	}
	return bb
}

func (d *D) BitBufLen(nBits int64) *bitio.Buffer {
	bs, err := d.bitBuf.BitBufLen(nBits)
	if err != nil {
		panic(ReadError{Err: err, Op: "BitBufLen", Size: nBits, Pos: d.Pos()})
	}
	return bs
}

func (d *D) Pos() int64 {
	bPos, err := d.bitBuf.Pos()
	if err != nil {
		panic(ReadError{Err: err, Op: "Pos", Size: 0, Pos: bPos})
	}
	return bPos
}

func (d *D) Len() int64 {
	return d.bitBuf.Len()
}

func (d *D) End() bool {
	bEnd, err := d.bitBuf.End()
	if err != nil {
		panic(ReadError{Err: err, Op: "Len", Size: 0, Pos: d.Pos()})
	}
	return bEnd
}

func (d *D) NotEnd() bool { return !d.End() }

func (d *D) BitsLeft() int64 {
	bBitsLeft, err := d.bitBuf.BitsLeft()
	if err != nil {
		panic(ReadError{Err: err, Op: "BitsLeft", Size: 0, Pos: d.Pos()})
	}
	return bBitsLeft
}

func (d *D) ByteAlignBits() int {
	bByteAlignBits, err := d.bitBuf.ByteAlignBits()
	if err != nil {
		panic(ReadError{Err: err, Op: "ByteAlignBits", Size: 0, Pos: d.Pos()})
	}
	return bByteAlignBits
}

func (d *D) BytePos() int64 {
	bBytePos, err := d.bitBuf.BytePos()
	if err != nil {
		panic(ReadError{Err: err, Op: "BytePos", Size: 0, Pos: d.Pos()})
	}
	return bBytePos
}

func (d *D) SeekRel(deltaBits int64) int64 {
	pos, err := d.bitBuf.SeekRel(deltaBits)
	if err != nil {
		panic(ReadError{Err: err, Op: "SeekRel", Delta: deltaBits, Pos: d.Pos()})
	}
	return pos
}

func (d *D) SeekAbs(pos int64) int64 {
	pos, err := d.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(ReadError{Err: err, Op: "SeekAbs", Size: pos, Pos: d.Pos()})
	}
	return pos
}

func (d *D) AddChild(v *Value) {
	v.Parent = d.Value

	switch fv := d.Value.V.(type) {
	case Struct:
		for _, ff := range fv {
			if ff.Name == v.Name {
				panic(fmt.Sprintf("%s already exist in struct %s", v.Name, d.Value.Name))
			}
		}
		d.Value.V = append(fv, v)
		return
	case Array:
		d.Value.V = append(fv, v)
	}

}

func (d *D) fieldDecoder(name string, bitBuf *bitio.Buffer, v interface{}) *D {
	cd := &D{
		Endian: d.Endian,
		bitBuf: bitBuf,
		Value: &Value{
			Name:   name,
			V:      v,
			Range:  ranges.Range{Start: d.Pos(), Len: 0},
			BitBuf: d.bitBuf,
		},
		registry: d.registry,
	}

	// TODO: refactor
	if d.Value != nil {
		d.AddChild(cd.Value)
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

func (d *D) FieldStructArrayLoopFn(name string, condFn func() bool, fn func(d *D)) *D {
	return d.FieldArrayFn(name, func(d *D) {
		for condFn() {
			d.FieldStructFn(name, fn)
		}
	})
}

func (d *D) FieldArrayLoopFn(name string, condFn func() bool, fn func(d *D)) *D {
	return d.FieldArrayFn(name, func(d *D) {
		for condFn() {
			fn(d)
		}
	})
}

func (d *D) FieldStructFn(name string, fn func(d *D)) *D {
	cd := d.FieldStruct(name)
	fn(cd)
	return cd
}

func (d *D) FieldStructBitBuf(name string, bitBuf *bitio.Buffer) *D {
	return d.fieldDecoder(name, bitBuf, Struct{})
}

func (d *D) FieldStructBitBufFn(name string, bitBuf *bitio.Buffer, fn func(d *D)) *D {
	cd := d.FieldStructBitBuf(name, bitBuf)
	fn(cd)
	return cd
}

func (d *D) FieldRangeFn(name string, firstBit int64, nBits int64, fn func() *Value) *Value {
	v := fn()
	v.Name = name
	v.BitBuf = d.bitBuf
	v.Range = ranges.Range{Start: firstBit, Len: nBits}
	d.AddChild(v)

	return v
}

func (d *D) FieldFn(name string, fn func() *Value) *Value {
	start := d.Pos()
	v := fn()
	stop := d.Pos()
	v.Name = name
	v.BitBuf = d.bitBuf
	v.Range = ranges.Range{Start: start, Len: stop - start}
	d.AddChild(v)

	return v
}

func (d *D) FieldBoolFn(name string, fn func() (bool, string)) bool {
	return d.FieldFn(name, func() *Value {
		b, d := fn()
		return &Value{V: b, Symbol: d}
	}).V.(bool)
}

func (d *D) FieldUFn(name string, fn func() (uint64, DisplayFormat, string)) uint64 {
	return d.FieldFn(name, func() *Value {
		u, fmt, d := fn()
		return &Value{V: u, DisplayFormat: fmt, Symbol: d}
	}).V.(uint64)
}

func (d *D) FieldSFn(name string, fn func() (int64, DisplayFormat, string)) int64 {
	return d.FieldFn(name, func() *Value {
		s, fmt, d := fn()
		return &Value{V: s, DisplayFormat: fmt, Symbol: d}
	}).V.(int64)
}

func (d *D) FieldFloatFn(name string, fn func() (float64, string)) float64 {
	return d.FieldFn(name, func() *Value {
		f, d := fn()
		return &Value{V: f, Symbol: d}
	}).V.(float64)
}

func (d *D) FieldStrFn(name string, fn func() (string, string)) string {
	return d.FieldFn(name, func() *Value {
		str, disp := fn()
		return &Value{V: str, Symbol: disp}
	}).V.(string)
}

func (d *D) FieldBytesFn(name string, firstBit int64, nBits int64, fn func() ([]byte, string)) []byte {
	return d.FieldRangeFn(name, firstBit, nBits, func() *Value {
		bs, disp := fn()
		return &Value{V: bs, Symbol: disp}
	}).V.([]byte)
}

func (d *D) FieldBitBufFn(name string, firstBit int64, nBits int64, fn func() (*bitio.Buffer, string)) *bitio.Buffer {
	return d.FieldRangeFn(name, firstBit, nBits, func() *Value {
		bb, disp := fn()
		return &Value{V: bb, Symbol: disp}
	}).V.(*bitio.Buffer)
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
	pos := d.Pos()
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

func (d *D) FieldValidateUTF8Fn(name string, v string, fn func() string) {
	pos := d.Pos()
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

func (d *D) FieldValidateUTF8(name string, v string) {
	pos := d.Pos()
	s := d.FieldStrFn(name, func() (string, string) {
		nBytes := len(v)
		str, err := d.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldValidateUTF8", Size: int64(nBytes) * 8, Pos: d.Pos()})
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

func (d *D) ValidateAtLeastBitsLeft(nBits int64) {
	bl := d.BitsLeft()
	if bl < nBits {
		// TODO:
		panic(ValidateError{Reason: fmt.Sprintf("expected bits left %d, found %d", nBits, bl), Pos: d.Pos()})
	}
}

func (d *D) ValidateAtLeastBytesLeft(nBytes int64) {
	bl := d.BitsLeft()
	if bl < nBytes*8 {
		// TODO:
		panic(ValidateError{Reason: fmt.Sprintf("expected bytes left %d, found %d bits", nBytes, bl), Pos: d.Pos()})
	}
}

// TODO: rename?
func (d *D) SubLenFn(nBits int64, fn func(d *D)) {
	prevBb := d.bitBuf
	prevEndian := d.Endian
	endPos := d.Pos() + nBits

	bb := d.BitBufRange(0, d.Pos()+nBits)
	if _, err := bb.SeekAbs(d.Pos()); err != nil {
		panic(err)
	}
	d.bitBuf = bb

	fn(d)

	d.bitBuf = prevBb
	d.bitBuf.SeekAbs(endPos) // TODO: check err?
	d.Endian = prevEndian
}

func (d *D) SubRangeFn(firstBit int64, nBits int64, fn func(d *D)) {
	prevBb := d.bitBuf
	prevEndian := d.Endian

	bb := d.BitBufRange(0, firstBit+nBits)
	if _, err := bb.SeekAbs(firstBit); err != nil {
		panic(err)
	}
	d.bitBuf = bb

	fn(d)

	d.bitBuf = prevBb
	d.Endian = prevEndian
}

func (d *D) FieldTryDecode(name string, formats []*Format) (*Value, interface{}, []error) {
	bb := d.BitBufRange(d.Pos(), d.BitsLeft())
	v, dv, errs := probe(name, bb, formats, probeOptions{isRoot: false, relStart: d.Pos(), relBitBuf: d.bitBuf})
	if v == nil || v.Errors() != nil {
		return nil, nil, errs
	}

	d.AddChild(v)
	if _, err := d.bitBuf.SeekRel(int64(v.Range.Len)); err != nil {
		panic(err)
	}

	return v, dv, errs
}

func (d *D) FieldDecode(name string, formats []*Format) (*Value, interface{}, []error) {
	v, dv, errs := d.FieldTryDecode(name, formats)
	if v == nil || v.Errors() != nil {
		panic(errs)
	}
	return v, dv, errs
}

func (d *D) FieldTryDecodeLen(name string, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	bb := d.BitBufRange(d.Pos(), nBits)
	v, dv, errs := probe(name, bb, formats, probeOptions{isRoot: false, relStart: d.Pos(), relBitBuf: d.bitBuf})
	if v == nil || v.Errors() != nil {
		return nil, nil, errs
	}

	d.AddChild(v)
	if _, err := d.bitBuf.SeekRel(int64(nBits)); err != nil {
		panic(err)
	}

	return v, dv, errs
}

func (d *D) FieldDecodeLen(name string, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	v, dv, errs := d.FieldTryDecodeLen(name, nBits, formats)
	if v == nil || v.Errors() != nil {
		panic(errs)
	}
	return v, dv, errs
}

// TODO: return decooder?
func (d *D) FieldTryDecodeRange(name string, firstBit int64, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	bb := d.BitBufRange(firstBit, nBits)
	v, dv, errs := probe(name, bb, formats, probeOptions{isRoot: false, relStart: firstBit, relBitBuf: d.bitBuf})
	if v == nil || v.Errors() != nil {
		return nil, nil, errs
	}

	d.AddChild(v)

	return v, dv, errs
}

func (d *D) FieldDecodeRange(name string, firstBit int64, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	v, dv, errs := d.FieldTryDecodeRange(name, firstBit, nBits, formats)
	if v == nil || v.Errors() != nil {
		panic(errs)
	}

	return v, dv, errs
}

func (d *D) FieldTryDecodeBitBuf(name string, bb *bitio.Buffer, formats []*Format) (*Value, interface{}, []error) {
	v, dv, errs := probe(name, bb, formats, probeOptions{isRoot: true})
	if v == nil || v.Errors() != nil {
		return nil, nil, errs
	}

	d.AddChild(v)

	return v, dv, errs
}

func (d *D) FieldDecodeBitBuf(name string, bb *bitio.Buffer, formats []*Format) (*Value, interface{}, []error) {
	v, dv, errs := d.FieldTryDecodeBitBuf(name, bb, formats)
	if v == nil || v.Errors() != nil {
		panic(errs)
	}
	return v, dv, errs
}

func (d *D) FieldBitBufRange(name string, firstBit int64, nBits int64) *bitio.Buffer {
	return d.FieldBitBufFn(name, firstBit, nBits, func() (*bitio.Buffer, string) {
		return d.BitBufRange(firstBit, nBits), ""
	})
}

func (d *D) FieldBitBufLen(name string, nBits int64) *bitio.Buffer {
	return d.FieldBitBufFn(name, d.Pos(), nBits, func() (*bitio.Buffer, string) {
		return d.BitBufLen(nBits), ""
	})
}

// TODO: range?
func (d *D) FieldDecodeZlibLen(name string, nBits int64, formats []*Format) (*Value, interface{}, []error) {
	bb, err := d.bitBuf.BitBufLen(nBits)
	if err != nil {
		panic(err)
	}
	zr, err := zlib.NewReader(bb)
	if err != nil {
		panic(err)
	}
	zd, err := ioutil.ReadAll(zr)
	if err != nil {
		panic(err)
	}
	zbb := bitio.NewBufferFromBytes(zd, -1)

	return d.FieldDecodeBitBuf(name, zbb, formats)
}
