package decode

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/wader/fq/internal/recoverfn"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/fq/pkg/scalar"
)

//go:generate sh -c "cat decode_gen.go.tmpl | go run ../../dev/tmpl.go types.json | gofmt > decode_gen.go"

type Endian int

const (
	// BigEndian byte order
	BigEndian = iota
	// LittleEndian byte order
	LittleEndian
)

type Options struct {
	Name          string
	Description   string
	Force         bool
	FillGaps      bool
	IsRoot        bool
	Range         ranges.Range // if zero use whole buffer
	FormatOptions map[string]interface{}
	FormatInArg   interface{}
	ReadBuf       *[]byte
}

// Decode try decode group and return first success and all other decoder errors
func Decode(ctx context.Context, bb *bitio.Buffer, group Group, opts Options) (*Value, interface{}, error) {
	return decode(ctx, bb, group, opts)
}

func decode(ctx context.Context, bb *bitio.Buffer, group Group, opts Options) (*Value, interface{}, error) {
	decodeRange := opts.Range
	if decodeRange.IsZero() {
		decodeRange = ranges.Range{Len: bb.Len()}
	}

	if group == nil {
		panic("group is nil, failed to register format?")
	}

	formatsErr := FormatsError{}

	for _, g := range group {
		cbb, err := bb.BitBufRange(decodeRange.Start, decodeRange.Len)
		if err != nil {
			return nil, nil, IOError{Err: err, Op: "BitBufRange", ReadSize: decodeRange.Len, Pos: decodeRange.Start}
		}

		d := newDecoder(ctx, g, cbb, opts)

		var decodeV interface{}
		r, rOk := recoverfn.Run(func() {
			decodeV = g.DecodeFn(d, opts.FormatInArg)
		})

		if ctx != nil && ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}

		if !rOk {
			if re, ok := r.RecoverV.(RecoverableErrorer); ok && re.IsRecoverableError() {
				panicErr, _ := re.(error)
				formatErr := FormatError{
					Err:        panicErr,
					Format:     g,
					Stacktrace: r,
				}
				formatsErr.Errs = append(formatsErr.Errs, formatErr)

				switch vv := d.Value.V.(type) {
				case *Compound:
					// TODO: hack, changes V
					vv.Err = formatErr
					d.Value.V = vv
				}

				if len(group) != 1 {
					continue
				}
			} else {
				r.RePanic()
			}
		}

		// TODO: maybe move to Format* funcs?
		if opts.FillGaps {
			d.FillGaps(ranges.Range{Start: 0, Len: decodeRange.Len}, "unknown")
		}

		var minMaxRange ranges.Range
		if err := d.Value.WalkRootPreOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
			minMaxRange = ranges.MinMax(minMaxRange, v.Range)
			v.Range.Start += decodeRange.Start
			v.RootBitBuf = bb
			return nil
		}); err != nil {
			return nil, nil, err
		}

		d.Value.Range = ranges.Range{Start: decodeRange.Start, Len: minMaxRange.Len}

		if opts.IsRoot {
			d.Value.postProcess()
		}

		if len(formatsErr.Errs) > 0 {
			return d.Value, decodeV, formatsErr
		}

		return d.Value, decodeV, nil
	}

	return nil, nil, formatsErr
}

type D struct {
	Ctx     context.Context
	Endian  Endian
	Value   *Value
	Options Options

	bitBuf *bitio.Buffer

	readBuf *[]byte
}

// TODO: new struct decoder?
// note bb is assumed to be a non-shared buffer
func newDecoder(ctx context.Context, format Format, bb *bitio.Buffer, opts Options) *D {
	name := format.RootName
	if opts.Name != "" {
		name = opts.Name
	}
	rootV := &Compound{
		IsArray:     format.RootArray,
		Children:    nil,
		Description: opts.Description,
		Format:      &format,
	}

	return &D{
		Ctx:    ctx,
		Endian: BigEndian,
		Value: &Value{
			Name:       name,
			V:          rootV,
			RootBitBuf: bb,
			Range:      ranges.Range{Start: 0, Len: 0},
			IsRoot:     opts.IsRoot,
		},
		Options: opts,

		bitBuf:  bb,
		readBuf: opts.ReadBuf,
	}
}

func (d *D) FieldDecoder(name string, bitBuf *bitio.Buffer, v interface{}) *D {
	return &D{
		Ctx:    d.Ctx,
		Endian: d.Endian,
		Value: &Value{
			Name:       name,
			V:          v,
			Range:      ranges.Range{Start: d.Pos(), Len: 0},
			RootBitBuf: bitBuf,
		},
		Options: d.Options,

		bitBuf:  bitBuf,
		readBuf: d.readBuf,
	}
}

func (d *D) Copy(r io.Writer, w io.Reader) (int64, error) {
	// TODO: what size? now same as io.Copy
	buf := d.SharedReadBuf(32 * 1024)
	return io.CopyBuffer(r, w, buf)
}

func (d *D) MustCopy(r io.Writer, w io.Reader) int64 {
	n, err := d.Copy(r, w)
	if err != nil {
		d.IOPanic(err, "MustCopy: Copy")
	}
	return n
}

func (d *D) MustNewBitBufFromReader(r io.Reader) *bitio.Buffer {
	b := &bytes.Buffer{}
	d.MustCopy(b, r)
	return bitio.NewBufferFromBytes(b.Bytes(), -1)
}

func (d *D) SharedReadBuf(n int) []byte {
	if d.readBuf == nil {
		d.readBuf = new([]byte)
	}
	if len(*d.readBuf) < n {
		*d.readBuf = make([]byte, n)
	}
	return *d.readBuf
}

func (d *D) FillGaps(r ranges.Range, namePrefix string) {
	makeWalkFn := func(fn func(iv *Value)) func(iv *Value, rootV *Value, depth int, rootDepth int) error {
		return func(iv *Value, rootV *Value, depth int, rootDepth int) error {
			switch iv.V.(type) {
			case *Compound:
			default:
				fn(iv)
			}
			return nil
		}
	}

	// TODO: redo this, tries to get rid of slice grow
	// TODO: pre-sorted somehow?
	n := 0
	_ = d.Value.WalkRootPreOrder(makeWalkFn(func(iv *Value) { n++ }))
	valueRanges := make([]ranges.Range, n)
	i := 0
	_ = d.Value.WalkRootPreOrder(makeWalkFn(func(iv *Value) {
		valueRanges[i] = iv.Range
		i++
	}))

	gaps := ranges.Gaps(r, valueRanges)
	for i, gap := range gaps {
		bb, err := d.bitBuf.BitBufRange(gap.Start, gap.Len)
		if err != nil {
			d.IOPanic(err, "FillGaps: BitBufRange")
		}

		v := &Value{
			Name: fmt.Sprintf("%s%d", namePrefix, i),
			V: &scalar.S{
				Actual:  bb,
				Unknown: true,
			},
			RootBitBuf: d.bitBuf,
			Range:      gap,
		}

		d.AddChild(v)
	}
}

// Errorf stops decode with a reason unless forced
func (d *D) Errorf(format string, a ...interface{}) {
	if !d.Options.Force {
		panic(DecoderError{Reason: fmt.Sprintf(format, a...), Pos: d.Pos()})
	}
}

// Fatalf stops decode with a reason regardless of forced
func (d *D) Fatalf(format string, a ...interface{}) {
	panic(DecoderError{Reason: fmt.Sprintf(format, a...), Pos: d.Pos()})
}

func (d *D) IOPanic(err error, op string) {
	panic(IOError{Err: err, Pos: d.Pos(), Op: op})
}

// Bits reads nBits bits from buffer
func (d *D) bits(nBits int) (uint64, error) {
	if nBits < 0 || nBits > 64 {
		return 0, fmt.Errorf("nBits must be 0-64 (%d)", nBits)
	}
	// 64 bits max, 9 byte worse case if not byte aligned
	buf := d.SharedReadBuf(9)
	_, err := bitio.ReadFull(d.bitBuf, buf, nBits)
	if err != nil {
		return 0, err
	}

	return bitio.Read64(buf[:], 0, nBits), nil
}

// Bits reads nBits bits from buffer
func (d *D) Bits(nBits int) (uint64, error) {
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (d *D) PeekBits(nBits int) uint64 {
	n, err := d.TryPeekBits(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekBits", ReadSize: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) PeekBytes(nBytes int) []byte {
	bs, err := d.bitBuf.PeekBytes(nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekBytes", ReadSize: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
}

func (d *D) PeekFind(nBits int, seekBits int64, fn func(v uint64) bool, maxLen int64) (int64, uint64) {
	peekBits, v, err := d.TryPeekFind(nBits, seekBits, maxLen, fn)
	if err != nil {
		d.IOPanic(err, "PeekFind: TryPeekFind")
	}
	if peekBits == -1 {
		d.Errorf("peek not found")
	}
	return peekBits, v
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
	peekBits, _, err := d.TryPeekFind(8, 8, maxLen*8, func(v uint64) bool {
		return uint64(findV) == v
	})
	if err != nil {
		panic(IOError{Err: err, Op: "PeekFindByte", ReadSize: 0, Pos: d.Pos()})

	}
	return peekBits / 8
}

// PeekBits peek nBits bits from buffer
// TODO: share code?
func (d *D) TryPeekBits(nBits int) (uint64, error) {
	start, err := d.bitBuf.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	n, err := d.bits(nBits)
	if _, err := d.bitBuf.SeekBits(start, io.SeekStart); err != nil {
		return 0, err
	}
	return n, err
}

func (d *D) TryPeekFind(nBits int, seekBits int64, maxLen int64, fn func(v uint64) bool) (int64, uint64, error) {
	start, err := d.bitBuf.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}

	var count int64

	if seekBits < 0 {
		count = int64(-nBits)
		if _, err := d.bitBuf.SeekBits(start+count, io.SeekStart); err != nil {
			return 0, 0, err
		}
	}

	found := false
	var v uint64
	for {
		if (seekBits > 0 && maxLen > 0 && count >= maxLen) || (seekBits < 0 && maxLen > 0 && count < -maxLen) {
			break
		}
		v, err = d.TryU(nBits)
		if err != nil {
			if _, err := d.bitBuf.SeekBits(start, io.SeekStart); err != nil {
				return 0, 0, err
			}
			return 0, 0, err
		}
		if fn(v) {
			found = true
			break
		}
		count += seekBits
		if _, err := d.bitBuf.SeekBits(start+count, io.SeekStart); err != nil {
			return 0, 0, err
		}
	}
	if _, err := d.bitBuf.SeekBits(start, io.SeekStart); err != nil {
		return 0, 0, err
	}

	if !found {
		return -1, 0, nil
	}

	return count, v, nil
}

func (d *D) BytesRange(firstBit int64, nBytes int) []byte {
	bs, err := d.bitBuf.BytesRange(firstBit, nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "BytesRange", ReadSize: int64(nBytes) * 8, Pos: firstBit})
	}
	return bs
}

func (d *D) BytesLen(nBytes int) []byte {
	bs, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "BytesLen", ReadSize: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
}

// TODO: rename/remove BitBuf name?
func (d *D) BitBufRange(firstBit int64, nBits int64) *bitio.Buffer {
	bb, err := d.bitBuf.BitBufRange(firstBit, nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "BitBufRange", ReadSize: nBits, Pos: firstBit})
	}
	return bb
}

func (d *D) Pos() int64 {
	bPos, err := d.bitBuf.Pos()
	if err != nil {
		panic(IOError{Err: err, Op: "Pos", ReadSize: 0, Pos: bPos})
	}
	return bPos
}

func (d *D) Len() int64 {
	return d.bitBuf.Len()
}

func (d *D) End() bool {
	bEnd, err := d.bitBuf.End()
	if err != nil {
		panic(IOError{Err: err, Op: "Len", ReadSize: 0, Pos: d.Pos()})
	}
	return bEnd
}

func (d *D) NotEnd() bool { return !d.End() }

func (d *D) BitsLeft() int64 {
	bBitsLeft, err := d.bitBuf.BitsLeft()
	if err != nil {
		panic(IOError{Err: err, Op: "BitsLeft", ReadSize: 0, Pos: d.Pos()})
	}
	return bBitsLeft
}

func (d *D) AlignBits(nBits int) int {
	bByteAlignBits, err := d.bitBuf.AlignBits(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "AlignBits", ReadSize: 0, Pos: d.Pos()})
	}
	return bByteAlignBits
}

func (d *D) ByteAlignBits() int {
	bByteAlignBits, err := d.bitBuf.ByteAlignBits()
	if err != nil {
		panic(IOError{Err: err, Op: "ByteAlignBits", ReadSize: 0, Pos: d.Pos()})
	}
	return bByteAlignBits
}

func (d *D) BytePos() int64 {
	bBytePos, err := d.bitBuf.BytePos()
	if err != nil {
		panic(IOError{Err: err, Op: "BytePos", ReadSize: 0, Pos: d.Pos()})
	}
	return bBytePos
}

func (d *D) seekAbs(pos int64, name string, fns ...func(d *D)) int64 {
	var oldPos int64
	if len(fns) > 0 {
		oldPos = d.Pos()
	}

	pos, err := d.bitBuf.SeekAbs(pos)
	if err != nil {
		panic(IOError{Err: err, Op: name, SeekPos: pos, Pos: d.Pos()})
	}

	if len(fns) > 0 {
		for _, fn := range fns {
			fn(d)
		}
		_, err := d.bitBuf.SeekAbs(oldPos)
		if err != nil {
			panic(IOError{Err: err, Op: name, SeekPos: pos, Pos: d.Pos()})
		}
	}

	return pos
}

func (d *D) SeekRel(deltaPos int64, fns ...func(d *D)) int64 {
	return d.seekAbs(d.Pos()+deltaPos, "SeekRel", fns...)
}

func (d *D) SeekAbs(pos int64, fns ...func(d *D)) int64 {
	return d.seekAbs(pos, "SeekAbs", fns...)
}

func (d *D) AddChild(v *Value) {
	v.Parent = d.Value

	switch fv := d.Value.V.(type) {
	case *Compound:
		if !fv.IsArray {
			for _, ff := range fv.Children {
				if ff.Name == v.Name {
					d.Fatalf("%q already exist in struct %s", v.Name, d.Value.Name)
				}
			}
		}
		fv.Children = append(fv.Children, v)
	}
}

func (d *D) FieldGet(name string) *Value {
	switch fv := d.Value.V.(type) {
	case *Compound:
		for _, ff := range fv.Children {
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

func (d *D) FieldArray(name string, fn func(d *D), sms ...scalar.Mapper) *D {
	cd := d.FieldDecoder(name, d.bitBuf, &Compound{IsArray: true})
	d.AddChild(cd.Value)
	fn(cd)
	return cd
}

func (d *D) FieldArrayValue(name string) *D {
	return d.FieldArray(name, func(d *D) {})
}

func (d *D) FieldStruct(name string, fn func(d *D)) *D {
	cd := d.FieldDecoder(name, d.bitBuf, &Compound{})
	d.AddChild(cd.Value)
	fn(cd)
	return cd
}

func (d *D) FieldStructValue(name string) *D {
	return d.FieldStruct(name, func(d *D) {})
}

func (d *D) FieldStructArrayLoop(name string, structName string, condFn func() bool, fn func(d *D)) *D {
	return d.FieldArray(name, func(d *D) {
		for condFn() {
			d.FieldStruct(structName, fn)
		}
	})
}

func (d *D) FieldArrayLoop(name string, condFn func() bool, fn func(d *D)) *D {
	return d.FieldArray(name, func(d *D) {
		for condFn() {
			fn(d)
		}
	})
}

func (d *D) FieldRangeFn(name string, firstBit int64, nBits int64, fn func() *Value) *Value {
	v := fn()
	v.Name = name
	v.RootBitBuf = d.bitBuf
	v.Range = ranges.Range{Start: firstBit, Len: nBits}
	d.AddChild(v)

	return v
}

func (d *D) AssertAtLeastBitsLeft(nBits int64) {
	if d.Options.Force {
		return
	}
	bl := d.BitsLeft()
	if bl < nBits {
		// TODO:
		panic(DecoderError{Reason: fmt.Sprintf("expected bits left %d, found %d", nBits, bl), Pos: d.Pos()})
	}
}

func (d *D) AssertLeastBytesLeft(nBytes int64) {
	if d.Options.Force {
		return
	}
	bl := d.BitsLeft()
	if bl < nBytes*8 {
		// TODO:
		panic(DecoderError{Reason: fmt.Sprintf("expected bytes left %d, found %d bits", nBytes, bl), Pos: d.Pos()})
	}
}

// TODO: rethink
func (d *D) FieldValueU(name string, a uint64, sms ...scalar.Mapper) {
	d.FieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return scalar.S{Actual: a}, nil }, sms...)
}

func (d *D) FieldValueS(name string, a int64, sms ...scalar.Mapper) {
	d.FieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return scalar.S{Actual: a}, nil }, sms...)
}

func (d *D) FieldValueBool(name string, a bool, sms ...scalar.Mapper) {
	d.FieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return scalar.S{Actual: a}, nil }, sms...)
}

func (d *D) FieldValueFloat(name string, a float64, sms ...scalar.Mapper) {
	d.FieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return scalar.S{Actual: a}, nil }, sms...)
}

func (d *D) FieldValueStr(name string, a string, sms ...scalar.Mapper) {
	d.FieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return scalar.S{Actual: a}, nil }, sms...)
}

func (d *D) FieldValueRaw(name string, a []byte, sms ...scalar.Mapper) {
	d.FieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		return scalar.S{Actual: bitio.NewBufferFromBytes(a, -1)}, nil
	}, sms...)
}

func (d *D) LenFn(nBits int64, fn func(d *D)) {
	d.RangeFn(d.Pos(), nBits, fn)
	d.SeekRel(nBits)
}

func (d *D) RangeFn(firstBit int64, nBits int64, fn func(d *D)) {
	var subV interface{}
	switch vv := d.Value.V.(type) {
	case *Compound:
		subV = &Compound{IsArray: vv.IsArray}
	default:
		panic("unreachable")
	}

	if nBits < 0 {
		nBits = d.Len() - firstBit
	}

	// TODO: do some kind of DecodeLimitedLen/RangeFn?
	bb := d.BitBufRange(0, firstBit+nBits)
	if _, err := bb.SeekAbs(firstBit); err != nil {
		d.IOPanic(err, "RangeFn: SeekAbs")
	}
	sd := d.FieldDecoder("", bb, subV)

	fn(sd)

	// TODO: refactor, similar to decode()
	if err := sd.Value.WalkRootPreOrder(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		//v.Range.Start += firstBit
		v.RootBitBuf = d.Value.RootBitBuf

		return nil
	}); err != nil {
		panic(err)
	}

	switch vv := sd.Value.V.(type) {
	case *Compound:
		for _, f := range vv.Children {
			d.AddChild(f)
		}
	default:
		panic("unreachable")
	}
}

func (d *D) Format(group Group, inArg interface{}) interface{} {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Force:       d.Options.Force,
		FillGaps:    false,
		IsRoot:      false,
		Range:       ranges.Range{Start: d.Pos(), Len: d.BitsLeft()},
		FormatInArg: inArg,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, "Format: decode")
	}

	switch vv := dv.V.(type) {
	case *Compound:
		for _, f := range vv.Children {
			d.AddChild(f)
		}
	default:
		panic("unreachable")
	}

	if _, err := d.bitBuf.SeekRel(dv.Range.Len); err != nil {
		d.IOPanic(err, "Format: SeekRel")
	}

	return v
}

func (d *D) TryFieldFormat(name string, group Group, inArg interface{}) (*Value, interface{}, error) {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    false,
		IsRoot:      false,
		Range:       ranges.Range{Start: d.Pos(), Len: d.BitsLeft()},
		FormatInArg: inArg,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)
	if _, err := d.bitBuf.SeekRel(dv.Range.Len); err != nil {
		d.IOPanic(err, "TryFieldFormat: SeekRel")
	}

	return dv, v, err
}

func (d *D) FieldFormat(name string, group Group, inArg interface{}) (*Value, interface{}) {
	dv, v, err := d.TryFieldFormat(name, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, "FieldFormat: TryFieldFormat")
	}
	return dv, v
}

func (d *D) TryFieldFormatLen(name string, nBits int64, group Group, inArg interface{}) (*Value, interface{}, error) {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    true,
		IsRoot:      false,
		Range:       ranges.Range{Start: d.Pos(), Len: nBits},
		FormatInArg: inArg,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)
	if _, err := d.bitBuf.SeekRel(nBits); err != nil {
		d.IOPanic(err, "TryFieldFormatLen: SeekRel")
	}

	return dv, v, err
}

func (d *D) FieldFormatLen(name string, nBits int64, group Group, inArg interface{}) (*Value, interface{}) {
	dv, v, err := d.TryFieldFormatLen(name, nBits, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, "FieldFormatLen: TryFieldFormatLen")
	}
	return dv, v
}

// TODO: return decooder?
func (d *D) TryFieldFormatRange(name string, firstBit int64, nBits int64, group Group, inArg interface{}) (*Value, interface{}, error) {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    true,
		IsRoot:      false,
		Range:       ranges.Range{Start: firstBit, Len: nBits},
		FormatInArg: inArg,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)

	return dv, v, err
}

func (d *D) FieldFormatRange(name string, firstBit int64, nBits int64, group Group, inArg interface{}) (*Value, interface{}) {
	dv, v, err := d.TryFieldFormatRange(name, firstBit, nBits, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, "FieldFormatRange: TryFieldFormatRange")
	}

	return dv, v
}

func (d *D) TryFieldFormatBitBuf(name string, bb *bitio.Buffer, group Group, inArg interface{}) (*Value, interface{}, error) {
	dv, v, err := decode(d.Ctx, bb, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    true,
		IsRoot:      true,
		FormatInArg: inArg,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	dv.Range.Start = d.Pos()

	d.AddChild(dv)

	return dv, v, err
}

func (d *D) FieldFormatBitBuf(name string, bb *bitio.Buffer, group Group, inArg interface{}) (*Value, interface{}) {
	dv, v, err := d.TryFieldFormatBitBuf(name, bb, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, "FieldFormatBitBuf: TryFieldFormatBitBuf")
	}

	return dv, v
}

// TODO: rethink this
func (d *D) FieldRootBitBuf(name string, bb *bitio.Buffer) *Value {
	v := &Value{}
	v.V = &scalar.S{Actual: bb}
	v.Name = name
	v.RootBitBuf = bb
	v.IsRoot = true
	v.Range = ranges.Range{Start: d.Pos(), Len: bb.Len()}
	d.AddChild(v)

	return v
}

func (d *D) FieldStructRootBitBufFn(name string, bb *bitio.Buffer, fn func(d *D)) *Value {
	cd := d.FieldDecoder(name, bb, &Compound{})
	cd.Value.IsRoot = true
	d.AddChild(cd.Value)
	fn(cd)

	cd.Value.postProcess()

	return cd.Value
}

// TODO: range?
func (d *D) FieldFormatReaderLen(name string, nBits int64, fn func(r io.Reader) (io.ReadCloser, error), group Group) (*Value, interface{}) {
	bb, err := d.bitBuf.BitBufLen(nBits)
	if err != nil {
		d.IOPanic(err, "FieldFormatReaderLen: BitBufLen")
	}
	zr, err := fn(bb)
	if err != nil {
		d.IOPanic(err, "FieldFormatReaderLen: fn")
	}
	zd, err := ioutil.ReadAll(zr)
	if err != nil {
		d.IOPanic(err, "FieldFormatReaderLen: ReadAll")
	}
	zbb := bitio.NewBufferFromBytes(zd, -1)

	return d.FieldFormatBitBuf(name, zbb, group, nil)
}

// TODO: too mant return values
func (d *D) TryFieldReaderRangeFormat(name string, startBit int64, nBits int64, fn func(r io.Reader) io.Reader, group Group, inArg interface{}) (int64, *bitio.Buffer, *Value, interface{}, error) {
	bitLen := nBits
	if bitLen == -1 {
		bitLen = d.BitsLeft()
	}
	bb, err := d.bitBuf.BitBufRange(startBit, bitLen)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	r := fn(bb)
	// TODO: check if io.Closer?
	rb, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	cz, err := bb.Pos()
	rbb := bitio.NewBufferFromBytes(rb, -1)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	dv, v, err := d.TryFieldFormatBitBuf(name, rbb, group, inArg)

	return cz, rbb, dv, v, err
}

func (d *D) FieldReaderRangeFormat(name string, startBit int64, nBits int64, fn func(r io.Reader) io.Reader, group Group, inArg interface{}) (int64, *bitio.Buffer, *Value, interface{}) {
	cz, rbb, dv, v, err := d.TryFieldReaderRangeFormat(name, startBit, nBits, fn, group, inArg)
	if err != nil {
		d.IOPanic(err, "TryFieldReaderRangeFormat")
	}
	return cz, rbb, dv, v
}

func (d *D) TryFieldValue(name string, fn func() (*Value, error)) (*Value, error) {
	start := d.Pos()
	v, err := fn()
	stop := d.Pos()
	v.Name = name
	v.RootBitBuf = d.bitBuf
	v.Range = ranges.Range{Start: start, Len: stop - start}
	if err != nil {
		return nil, err
	}
	d.AddChild(v)

	return v, err
}

func (d *D) FieldValue(name string, fn func() *Value) *Value {
	v, err := d.TryFieldValue(name, func() (*Value, error) { return fn(), nil })
	if err != nil {
		d.IOPanic(err, "FieldValue: TryFieldValue")
	}
	return v
}

// looks a bit weird to force at least one ScalarFn arg
func (d *D) TryFieldScalarFn(name string, sfn scalar.Fn, sms ...scalar.Mapper) (*scalar.S, error) {
	v, err := d.TryFieldValue(name, func() (*Value, error) {
		s, err := sfn(scalar.S{})
		if err != nil {
			return &Value{V: &s}, err
		}
		for _, sm := range sms {
			s, err = sm.MapScalar(s)
			if err != nil {
				return &Value{V: &s}, err
			}
		}
		return &Value{V: &s}, nil
	})
	if err != nil {
		return &scalar.S{}, err
	}
	return v.V.(*scalar.S), nil
}

func (d *D) FieldScalarFn(name string, sfn scalar.Fn, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarFn(name, sfn, sms...)
	if err != nil {
		d.IOPanic(err, "FieldScalarFn: TryFieldScalarFn")
	}
	return v
}

func (v *Value) TryScalarFn(sms ...scalar.Mapper) error {
	var err error
	sr, ok := v.V.(*scalar.S)
	if !ok {
		panic("not a scalar value")
	}
	s := *sr
	for _, sm := range sms {
		s, err = sm.MapScalar(s)
		if err != nil {
			break
		}
	}
	v.V = &s
	return err
}
