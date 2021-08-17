package decode

//go:generate sh -c "cat decode_decoder_gen.go.tmpl | go run ../../dev/tmpl.go decode_decoder_gen.go.json | gofmt > decode_decoder_gen.go"

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/internal/recoverfn"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
)

type DecodeFormatsError struct {
	Errs []FormatError
}

func (de DecodeFormatsError) Error() string {
	var errs []string
	for _, err := range de.Errs {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, ", ")
}

type FormatError struct {
	Err        error
	Format     *Format
	Stacktrace recoverfn.Raw
}

func (fe FormatError) Error() string {
	// var fns []string
	// for _, f := range fe.Stacktrace.Frames() {
	// 	fns = append(fns, fmt.Sprintf("%s:%d:%s", f.File, f.Line, f.Function))
	// }

	return fe.Err.Error()
}

func (fe FormatError) Value() interface{} {
	var st []interface{}
	for _, f := range fe.Stacktrace.Frames() {
		st = append(st, f.Function)
	}

	return map[string]interface{}{
		"format":     fe.Format.Name,
		"error":      fe.Err.Error(),
		"stacktrace": st,
	}
}

type IOError struct {
	Err   error
	Name  string
	Op    string
	Size  int64
	Delta int64
	Pos   int64
}

func (e IOError) Error() string {
	var prefix string
	if e.Name != "" {
		prefix = e.Op + "(" + e.Name + ")"
	} else {
		prefix = e.Op
	}

	return fmt.Sprintf("%s: failed at position %s (size %s delta %s): %s",
		prefix, num.Bits(e.Pos).StringByteBits(10), num.Bits(e.Size).StringByteBits(10), num.Bits(e.Delta).StringByteBits(10), e.Err)
}
func (e IOError) Unwrap() error { return e.Err }

type ValidateError struct {
	Reason string
	Pos    int64
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("failed to validate at position %s: %s", num.Bits(e.Pos).StringByteBits(16), e.Reason)
}

type Endian int

const (
	// BigEndian byte order
	BigEndian = iota
	// LittleEndian byte order
	LittleEndian
)

// just to get some type safety
type Options interface{ decodeOptions() }

type DecodeOptions struct {
	IsRoot        bool
	StartOffset   int64
	FormatOptions map[string]interface{}
	FormatInArg   interface{}
}

func (DecodeOptions) decodeOptions() {}

type FormatOptions struct {
	InArg interface{}
}

func (FormatOptions) decodeOptions() {}

// Decode try decode formats and return first success and all other decoder errors
func Decode(name string, description string, bb *bitio.Buffer, formats []*Format, opts ...Options) (*Value, interface{}, error) {
	opts = append(opts, DecodeOptions{IsRoot: true})
	return decode(name, description, bb, formats, opts)
}

func decode(name string, description string, bb *bitio.Buffer, formats []*Format, opts []Options) (*Value, interface{}, error) {
	if formats == nil {
		panic("formats is nil, failed to register format?")
	}

	var forceOne = len(formats) == 1

	var decodeOpts DecodeOptions
	var formatOpts FormatOptions
	for _, opts := range opts {
		switch opt := opts.(type) {
		case DecodeOptions:
			decodeOpts = opt
		case FormatOptions:
			formatOpts = opt
		default:
			panic(fmt.Sprintf("unknown decode option %#+v", opts))
		}
	}

	decodeErr := DecodeFormatsError{}

	for _, f := range formats {
		d := NewDecoder(name, description, f, bb, decodeOpts)

		var decodeV interface{}

		r, rOk := recoverfn.Run(func() {
			decodeV = f.DecodeFn(d, formatOpts.InArg)
		})

		if !rOk {
			switch panicV := r.RecoverV.(type) {
			case IOError, ValidateError, DecodeFormatsError:
				panicErr, _ := panicV.(error)
				formatErr := FormatError{
					Err:        panicErr,
					Format:     f,
					Stacktrace: r,
				}
				decodeErr.Errs = append(decodeErr.Errs, formatErr)

				d.Value.Err = formatErr

				if !forceOne {
					continue
				}
			default:
				r.RePanic()
			}
		}

		var maxRange ranges.Range
		if err := d.Value.WalkPreOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
			if d.Value != v && v.IsRoot {
				return ErrWalkSkipChildren
			}

			maxRange = ranges.MinMax(maxRange, v.Range)
			v.Range.Start += decodeOpts.StartOffset
			v.RootBitBuf = d.Value.RootBitBuf
			return nil
		}); err != nil {
			return nil, nil, err
		}

		d.Value.Range = ranges.Range{Start: decodeOpts.StartOffset, Len: maxRange.Len}

		if decodeOpts.IsRoot {
			d.FillGaps("unknown")

			// sort and set ranges for struct and arrays
			d.Value.postProcess()
		}

		return d.Value, decodeV, decodeErr
	}

	return nil, nil, decodeErr
}

type D struct {
	Endian  Endian
	Value   *Value
	Options map[string]interface{}

	bitBuf *bitio.Buffer

	bitsBuf []byte
}

// TODO: new struct decoder?
func NewDecoder(name string, description string, format *Format, bb *bitio.Buffer, opts DecodeOptions) *D {
	cbb := bb.Copy()

	return &D{
		Endian: BigEndian,
		Value: &Value{
			Name:        name,
			Description: description,
			Format:      format,
			V:           Struct{},
			IsRoot:      opts.IsRoot,
			RootBitBuf:  cbb,
			Range:       ranges.Range{Start: 0, Len: 0},
		},
		Options: opts.FormatOptions,

		bitBuf: cbb,
	}
}

func (d *D) FillGaps(namePrefix string) {
	// TODO: d.Value is array?

	makeWalkFn := func(fn func(iv *Value)) func(iv *Value, rootV *Value, depth int, rootDepth int) error {
		return func(iv *Value, rootV *Value, depth int, rootDepth int) error {
			if iv.RootBitBuf != d.Value.RootBitBuf && iv.IsRoot {
				return ErrWalkSkipChildren
			}
			switch iv.V.(type) {
			case Struct, Array:
			default:
				fn(iv)
			}
			return nil
		}
	}

	// TODO: redo this, tries to get rid of slice glow
	// TODO: gaps things should be done in Framed* funcs
	// TODO: pre-sorted somehow?
	n := 0
	_ = d.Value.WalkPreOrder(makeWalkFn(func(iv *Value) { n++ }))
	valueRanges := make([]ranges.Range, n)
	i := 0
	_ = d.Value.WalkPreOrder(makeWalkFn(func(iv *Value) {
		valueRanges[i] = iv.Range
		i++
	}))

	gaps := ranges.Gaps(ranges.Range{Start: 0, Len: d.Len()}, valueRanges)
	for i, gap := range gaps {
		v := d.FieldValueBitBufRange(
			fmt.Sprintf("%s%d", namePrefix, i), gap.Start, gap.Len,
		)
		v.Unknown = true
	}
}

// Invalid stops decode with a reason
func (d *D) Invalid(reason string) {
	panic(ValidateError{Reason: reason, Pos: d.Pos()})
}

func (d *D) PeekBits(nBits int) uint64 {
	n, err := d.TryPeekBits(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekBits", Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) PeekBytes(nBytes int) []byte {
	bs, err := d.bitBuf.PeekBytes(nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekBytes", Size: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
}

func (d *D) PeekFind(nBits int, seekBits int64, fn func(v uint64) bool, maxLen int64) int64 {
	peekBits, err := d.TryPeekFind(nBits, seekBits, fn, maxLen)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekFind", Size: 0, Pos: d.Pos()})
	}
	if peekBits == -1 {
		panic(IOError{Err: fmt.Errorf("not found"), Op: "PeekFind", Size: 0, Pos: d.Pos()})
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
func (d *D) PeekFindByte(findV uint8, maxLen int64) int64 {
	peekBits, err := d.TryPeekFind(8, 8, func(v uint64) bool {
		return uint64(findV) == v
	}, maxLen*8)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekFindByte", Size: 0, Pos: d.Pos()})

	}
	return peekBits / 8
}

func (d *D) BytesRange(firstBit int64, nBytes int) []byte {
	bs, err := d.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "BytesRange", Size: int64(nBytes) * 8, Pos: firstBit})
	}
	return bs
}

func (d *D) BytesLen(nBytes int) []byte {
	bs, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "BytesLen", Size: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
}

// TODO: rename/remove BitBuf name?
func (d *D) BitBufRange(firstBit int64, nBits int64) *bitio.Buffer {
	bb, err := d.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "BitBufRange", Size: nBits, Pos: firstBit})
	}
	return bb
}

func (d *D) BitBufLen(nBits int64) *bitio.Buffer {
	bs, err := d.bitBuf.BitBufLen(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "BitBufLen", Size: nBits, Pos: d.Pos()})
	}
	return bs
}

func (d *D) Pos() int64 {
	bPos, err := d.bitBuf.Pos()
	if err != nil {
		panic(IOError{Err: err, Op: "Pos", Size: 0, Pos: bPos})
	}
	return bPos
}

func (d *D) Len() int64 {
	return d.bitBuf.Len()
}

func (d *D) End() bool {
	bEnd, err := d.bitBuf.End()
	if err != nil {
		panic(IOError{Err: err, Op: "Len", Size: 0, Pos: d.Pos()})
	}
	return bEnd
}

func (d *D) NotEnd() bool { return !d.End() }

func (d *D) BitsLeft() int64 {
	bBitsLeft, err := d.bitBuf.BitsLeft()
	if err != nil {
		panic(IOError{Err: err, Op: "BitsLeft", Size: 0, Pos: d.Pos()})
	}
	return bBitsLeft
}

func (d *D) ByteAlignBits() int {
	bByteAlignBits, err := d.bitBuf.ByteAlignBits()
	if err != nil {
		panic(IOError{Err: err, Op: "ByteAlignBits", Size: 0, Pos: d.Pos()})
	}
	return bByteAlignBits
}

func (d *D) BytePos() int64 {
	bBytePos, err := d.bitBuf.BytePos()
	if err != nil {
		panic(IOError{Err: err, Op: "BytePos", Size: 0, Pos: d.Pos()})
	}
	return bBytePos
}

func (d *D) SeekRel(deltaBits int64) int64 {
	pos, err := d.bitBuf.SeekRel(deltaBits)
	if err != nil {
		panic(IOError{Err: err, Op: "SeekRel", Delta: deltaBits, Pos: d.Pos()})
	}
	return pos
}

func (d *D) SeekAbs(pos int64) int64 {
	pos, err := d.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(IOError{Err: err, Op: "SeekAbs", Size: pos, Pos: d.Pos()})
	}
	return pos
}

func (d *D) AddChild(v *Value) {
	v.Parent = d.Value

	switch fv := d.Value.V.(type) {
	case Struct:
		for _, ff := range fv {
			if ff.Name == v.Name {
				d.Invalid(fmt.Sprintf("%s already exist in struct %s", v.Name, d.Value.Name))
			}
		}
		d.Value.V = append(fv, v)
		return
	case Array:
		d.Value.V = append(fv, v)
	}
}

func (d *D) FieldFormatr(name string, bitBuf *bitio.Buffer, v interface{}) *D {
	return &D{
		Endian: d.Endian,
		Value: &Value{
			Name:       name,
			V:          v,
			Range:      ranges.Range{Start: d.Pos(), Len: 0},
			RootBitBuf: d.bitBuf,
		},
		Options: d.Options,

		bitBuf: bitBuf,
	}
}

func (d *D) FieldRemove(name string) *Value {
	switch fv := d.Value.V.(type) {
	case Struct:
		for fi, ff := range fv {
			if ff.Name == name {
				d.Value.V = append(fv[0:fi], fv[fi+1:]...)
				return ff
			}
		}
		panic(fmt.Sprintf("%s not found in struct %s", name, d.Value.Name))
	default:
		panic(fmt.Sprintf("%s is not a struct", d.Value.Name))
	}
}

func (d *D) FieldMustRemove(name string) *Value {
	if v := d.FieldRemove(name); v != nil {
		return v
	}
	panic(fmt.Sprintf("%s not found in struct %s", name, d.Value.Name))
}

func (d *D) FieldGet(name string) *Value {
	switch fv := d.Value.V.(type) {
	case Struct:
		for _, ff := range fv {
			if ff.Name == name {
				return ff
			}
		}
	default:
		panic(fmt.Sprintf("%s is not a struct", d.Value.Name))
	}
	return nil
}

func (d *D) FieldMustGet(name string) *Value {
	if v := d.FieldGet(name); v != nil {
		return v
	}
	panic(fmt.Sprintf("%s not found in struct %s", name, d.Value.Name))
}

func (d *D) FieldArray(name string) *D {
	cd := d.FieldFormatr(name, d.bitBuf, Array{})
	d.AddChild(cd.Value)
	return cd
}

func (d *D) FieldArrayFn(name string, fn func(d *D)) *D {
	cd := d.FieldArray(name)
	fn(cd)
	return cd
}

func (d *D) FieldStruct(name string) *D {
	cd := d.FieldFormatr(name, d.bitBuf, Struct{})
	d.AddChild(cd.Value)
	return cd
}

func (d *D) FieldStructArrayLoopFn(name string, structName string, condFn func() bool, fn func(d *D)) *D {
	return d.FieldArrayFn(name, func(d *D) {
		for condFn() {
			d.FieldStructFn(structName, fn)
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

func (d *D) FieldRangeFn(name string, firstBit int64, nBits int64, fn func() *Value) *Value {
	v := fn()
	v.Name = name
	v.RootBitBuf = d.bitBuf
	v.Range = ranges.Range{Start: firstBit, Len: nBits}
	d.AddChild(v)

	return v
}

func (d *D) TryFieldFn(name string, fn func() (*Value, error)) (*Value, error) {
	start := d.Pos()
	v, err := fn()
	stop := d.Pos()
	v.Name = name
	v.RootBitBuf = d.bitBuf
	v.Range = ranges.Range{Start: start, Len: stop - start}
	d.AddChild(v)

	return v, err
}

func (d *D) FieldFn(name string, fn func() *Value) *Value {
	start := d.Pos()
	v := fn()
	stop := d.Pos()
	v.Name = name
	v.RootBitBuf = d.bitBuf
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

func (d *D) FieldStrFn(name string, fn func() (string, string)) string {
	return d.FieldFn(name, func() *Value {
		str, desc := fn()
		return &Value{V: str, Description: desc}
	}).V.(string)
}

func (d *D) FieldBytesFn(name string, fn func() ([]byte, string)) []byte {
	return d.FieldFn(name, func() *Value {
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

func (d *D) FieldValueBitBufFn(name string, firstBit int64, nBits int64, fn func() (*bitio.Buffer, string)) *Value {
	return d.FieldRangeFn(name, firstBit, nBits, func() *Value {
		bb, disp := fn()
		return &Value{V: bb, Symbol: disp}
	})
}

func (d *D) FieldBoolMapFn(name string, trueS string, falseS string, fn func() bool) (bool, bool) {
	var ok bool
	return d.FieldBoolFn(name, func() (bool, string) {
		n := fn()
		d := falseS
		if n {
			d = trueS
		}
		return n, d
	}), ok
}

func (d *D) FieldStringMapFn(name string, sm map[uint64]string, def string, fn func() uint64, df DisplayFormat) (uint64, bool) {
	var ok bool
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n := fn()
		var d string
		d, ok = sm[n]
		if !ok {
			d = def
		}
		return n, df, d
	}), ok
}

func (d *D) FieldStringRangeMapFn(name string, rm map[[2]uint64]string, def string, fn func() uint64, df DisplayFormat) (uint64, bool) {
	var ok bool
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n := fn()
		for r, s := range rm {
			if n >= r[0] && n <= r[1] {
				return n, NumberDecimal, s
			}
		}
		return n, df, def
	}), ok
}

func (d *D) FieldStringUUIDMapFn(name string, um map[[16]byte]string, def string, fn func() []byte) ([]byte, bool) {
	var ok bool
	return d.FieldBytesFn(name, func() ([]byte, string) {
		uuid := fn()
		for u, s := range um {
			if bytes.Equal(u[:], uuid[:]) {
				return uuid, s
			}
		}
		return uuid, def
	}), ok
}

func (d *D) FieldChecksumRange(name string, firstBit int64, nBits int64, calculated []byte, endian Endian) {
	nBytes := int(nBits / 8)
	d.FieldRangeFn(name, firstBit, nBits, func() *Value {
		expectedBB := d.BitBufRange(firstBit, nBits)
		expected, _ := expectedBB.BytesLen(nBytes)

		if endian == LittleEndian {
			bitio.ReverseBytes(expected)
			expectedBB = bitio.NewBufferFromBytes(expected, -1)
		}

		if bytes.Equal(expected, calculated) {
			return &Value{V: expectedBB.Copy(), Symbol: "Correct"}
		}

		return &Value{V: expectedBB.Copy(), Symbol: fmt.Sprintf("Incorrect (calculated %s)", hex.EncodeToString(calculated))}
	})
}

func (d *D) FieldChecksumLen(name string, nBits int64, calculated []byte, endian Endian) {
	d.FieldChecksumRange(name, d.Pos(), nBits, calculated, endian)
	d.SeekRel(nBits)
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
		str, err := d.TryUTF8(nBytes)
		if err != nil {
			panic(IOError{Err: err, Name: name, Op: "FieldValidateUTF8", Size: int64(nBytes) * 8, Pos: d.Pos()})
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

// TODO: rethink
func (d *D) FieldValueU(name string, v uint64, symbol string) {
	d.FieldUFn(name, func() (uint64, DisplayFormat, string) { return v, NumberDecimal, symbol })
}

func (d *D) FieldValueS(name string, v int64, symbol string) {
	d.FieldSFn(name, func() (int64, DisplayFormat, string) { return v, NumberDecimal, symbol })
}

func (d *D) FieldValueBool(name string, v bool, symbol string) {
	d.FieldBoolFn(name, func() (bool, string) { return v, symbol })
}

func (d *D) FieldValueFloat(name string, v float64, symbol string) {
	d.FieldFloatFn(name, func() (float64, string) { return v, symbol })
}

func (d *D) FieldValueStr(name string, v string, symbol string) {
	d.FieldStrFn(name, func() (string, string) { return v, symbol })
}

func (d *D) FieldValueBytes(name string, b []byte, symbol string) {
	d.FieldBytesFn(name, func() ([]byte, string) { return b, symbol })
}

// TODO: rename?
func (d *D) DecodeLenFn(nBits int64, fn func(d *D)) {
	d.DecodeRangeFn(d.Pos(), nBits, fn)
	d.SeekRel(nBits)
}

func (d *D) DecodeRangeFn(firstBit int64, nBits int64, fn func(d *D)) {
	var subV interface{}
	switch d.Value.V.(type) {
	case Struct:
		subV = Struct{}
	case Array:
		subV = Array{}
	default:
		panic("unreachable")
	}

	// TODO: do some kind of DecodeLimitedLen/RangeFn?
	bb := d.BitBufRange(0, firstBit+nBits)
	if _, err := bb.SeekAbs(firstBit); err != nil {
		panic(IOError{Err: err, Op: "SeekAbs", Pos: firstBit})
	}
	sd := d.FieldFormatr("", bb, subV)

	fn(sd)

	// TODO: refactor, similar to decode()
	if err := sd.Value.WalkPreOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		if v.IsRoot {
			return ErrWalkSkipChildren
		}

		//v.Range.Start += firstBit
		v.RootBitBuf = d.Value.RootBitBuf

		return nil
	}); err != nil {
		panic(err)
	}

	switch vv := sd.Value.V.(type) {
	case Struct:
		for _, f := range vv {
			d.AddChild(f)
		}
	case Array:
		for _, f := range vv {
			d.AddChild(f)
		}
	default:
		panic("unreachable")
	}
}

func (d *D) Format(formats []*Format, opts ...Options) interface{} {
	bb := d.BitBufRange(d.Pos(), d.BitsLeft())
	opts = append(opts, DecodeOptions{IsRoot: false, StartOffset: d.Pos()})
	dv, v, err := decode("", "", bb, formats, opts)
	if dv == nil || dv.Errors() != nil {
		panic(err)
	}

	switch vv := dv.V.(type) {
	case Struct:
		for _, f := range vv {
			d.AddChild(f)
		}
	case Array:
		for _, f := range vv {
			d.AddChild(f)
		}
	default:
		panic("unreachable")
	}

	if _, err := d.bitBuf.SeekRel(dv.Range.Len); err != nil {
		panic(err)
	}

	return v
}

func (d *D) FieldTryFormat(name string, formats []*Format, opts ...Options) (*Value, interface{}, error) {
	bb := d.BitBufRange(d.Pos(), d.BitsLeft())
	opts = append(opts, DecodeOptions{IsRoot: false, StartOffset: d.Pos()})
	dv, v, err := decode(name, "", bb, formats, opts)
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)
	if _, err := d.bitBuf.SeekRel(dv.Range.Len); err != nil {
		panic(err)
	}

	return dv, v, err
}

func (d *D) FieldFormat(name string, formats []*Format, opts ...Options) (*Value, interface{}) {
	dv, v, err := d.FieldTryFormat(name, formats, opts...)
	if dv == nil || dv.Errors() != nil {
		panic(err)
	}
	return dv, v
}

func (d *D) FieldTryFormatLen(name string, nBits int64, formats []*Format, opts ...Options) (*Value, interface{}, error) {
	bb := d.BitBufRange(d.Pos(), nBits)
	opts = append(opts, DecodeOptions{IsRoot: false, StartOffset: d.Pos()})
	dv, v, err := decode(name, "", bb, formats, opts)
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)
	if _, err := d.bitBuf.SeekRel(nBits); err != nil {
		panic(err)
	}

	return dv, v, err
}

func (d *D) FieldFormatLen(name string, nBits int64, formats []*Format, opts ...Options) (*Value, interface{}) {
	dv, v, err := d.FieldTryFormatLen(name, nBits, formats, opts...)
	if dv == nil || dv.Errors() != nil {
		panic(err)
	}
	return dv, v
}

// TODO: return decooder?
func (d *D) FieldTryFormatRange(name string, firstBit int64, nBits int64, formats []*Format, opts ...Options) (*Value, interface{}, error) {
	bb := d.BitBufRange(firstBit, nBits)
	opts = append(opts, DecodeOptions{IsRoot: false, StartOffset: firstBit})
	dv, v, err := decode(name, "", bb, formats, opts)
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)

	return dv, v, err
}

func (d *D) FieldFormatRange(name string, firstBit int64, nBits int64, formats []*Format, opts ...Options) (*Value, interface{}) {
	dv, v, err := d.FieldTryFormatRange(name, firstBit, nBits, formats, opts...)
	if dv == nil || dv.Errors() != nil {
		panic(err)
	}

	return dv, v
}

func (d *D) FieldTryFormatBitBuf(name string, bb *bitio.Buffer, formats []*Format, opts ...Options) (*Value, interface{}, error) {
	opts = append(opts, DecodeOptions{IsRoot: true})
	dv, v, err := decode(name, "", bb, formats, opts)
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)

	return dv, v, err
}

func (d *D) FieldFormatBitBuf(name string, bb *bitio.Buffer, formats []*Format, opts ...Options) (*Value, interface{}) {
	dv, v, err := d.FieldTryFormatBitBuf(name, bb, formats, opts...)
	if dv == nil || dv.Errors() != nil {
		panic(err)
	}
	return dv, v
}

// TODO: rethink this
func (d *D) FieldBitBuf(name string, bb *bitio.Buffer) *bitio.Buffer {
	return d.FieldBitBufFn(name, d.Pos(), 0, func() (*bitio.Buffer, string) {
		return bb, ""
	})
}

// TODO: rethink this
func (d *D) FieldRootBitBuf(name string, bb *bitio.Buffer) *Value {
	v := &Value{}
	v.V = bb
	v.Name = name
	v.IsRoot = true
	v.RootBitBuf = bb
	v.Range = ranges.Range{Start: 0, Len: bb.Len()}
	d.AddChild(v)

	return v
}

func (d *D) FieldValueBitBufRange(name string, firstBit int64, nBits int64) *Value {
	return d.FieldValueBitBufFn(name, firstBit, nBits, func() (*bitio.Buffer, string) {
		return d.BitBufRange(firstBit, nBits), ""
	})
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
func (d *D) FieldFormatZlibLen(name string, nBits int64, formats []*Format) (*Value, interface{}) {
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

	return d.FieldFormatBitBuf(name, zbb, formats)
}

func (d *D) FieldStrNullTerminated(name string) string {
	return d.FieldStrFn(name, func() (string, string) {
		return d.StrNullTerminated(), ""
	})
}

func (d *D) StrNullTerminated() string {
	c := d.PeekFindByte(0, -1) + 1
	s := d.UTF8(int(c))
	return s[:len(s)-1]
}

func (d *D) FieldStrNullTerminatedLen(name string, nBytes int) string {
	return d.FieldStrFn(name, func() (string, string) {
		return d.StrNullTerminatedLen(nBytes), ""
	})
}

func (d *D) StrNullTerminatedLen(nBytes int) string {
	s := d.UTF8(nBytes)
	nullIndex := strings.IndexByte(s, 0)
	if nullIndex == -1 {
		return s
	}
	return s[:nullIndex]
}
