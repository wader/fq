package decode

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"reflect"
	"regexp"

	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/internal/iox"
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
	Name        string
	Description string
	Force       bool
	FillGaps    bool
	IsRoot      bool
	Range       ranges.Range // if zero use whole buffer
	InArg       any
	ParseOptsFn func(init any) any
	ReadBuf     *[]byte
}

// Decode try decode group and return first success and all other decoder errors
func Decode(ctx context.Context, br bitio.ReaderAtSeeker, group *Group, opts Options) (*Value, any, error) {
	return decode(ctx, br, group, opts)
}

func decode(ctx context.Context, br bitio.ReaderAtSeeker, group *Group, opts Options) (*Value, any, error) {
	brLen, err := bitiox.Len(br)
	if err != nil {
		return nil, nil, err
	}

	decodeRange := opts.Range
	if decodeRange.IsZero() {
		decodeRange = ranges.Range{Len: brLen}
	}

	if group == nil {
		panic("group is nil, failed to register format?")
	}

	hasGroupOpts := false
	groupArg := group.DefaultInArg
	if opts.ParseOptsFn != nil && groupArg != nil {
		if groupOptArg := opts.ParseOptsFn(groupArg); groupOptArg != nil {
			hasGroupOpts = true
			groupArg = groupOptArg
		}
	}

	formatsErr := FormatsError{}

	for _, f := range group.Formats {
		var inArgs []any

		// figure out if there are format specific arg passed as options
		hasFormatOpts := false
		formatArg := f.DefaultInArg
		if opts.ParseOptsFn != nil && formatArg != nil {
			if formatOptArg := opts.ParseOptsFn(formatArg); formatOptArg != nil {
				hasFormatOpts = true
				formatArg = formatOptArg
			}
		}

		// format options have priority
		if hasFormatOpts {
			inArgs = append(inArgs, formatArg)
		}
		if opts.InArg != nil {
			inArgs = append(inArgs, opts.InArg)
		}
		if !hasFormatOpts && f.DefaultInArg != nil {
			inArgs = append(inArgs, f.DefaultInArg)
		}

		if hasGroupOpts {
			inArgs = append(inArgs, groupArg)
		}

		cBR, err := bitiox.Range(br, decodeRange.Start, decodeRange.Len)
		if err != nil {
			return nil, nil, IOError{Err: err, Op: "BitBufRange", ReadSize: decodeRange.Len, Pos: decodeRange.Start}
		}

		d := newDecoder(ctx, f, cBR, opts)

		d.inArgs = inArgs

		var decodeV any
		r, rOk := recoverfn.Run(func() {
			decodeV = f.DecodeFn(d)
		})

		if ctx != nil && ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}

		if !rOk {
			var panicErr error
			if err, ok := r.RecoverV.(error); ok {
				panicErr = err
			} else {
				panicErr = fmt.Errorf("recoverable non-panic error :%v", r.RecoverV)
			}

			formatErr := FormatError{
				Err:        panicErr,
				Format:     f,
				Stacktrace: r,
			}
			formatsErr.Errs = append(formatsErr.Errs, formatErr)

			switch vv := d.Value.V.(type) {
			case *Compound:
				// TODO: hack, changes V
				d.Value.V = vv
				d.Value.Err = formatErr
			}

			if len(group.Formats) != 1 {
				continue
			}
		}

		// TODO: maybe move to Format* funcs?
		if opts.FillGaps {
			d.FillGaps(ranges.Range{Start: 0, Len: decodeRange.Len}, "gap")
		}

		var minMaxRange ranges.Range
		if err := d.Value.WalkRootPreOrder(func(v *Value, _ *Value, _ int, _ int) error {
			minMaxRange = ranges.MinMax(minMaxRange, v.Range)
			v.Range.Start += decodeRange.Start
			v.RootReader = br
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

	bitBuf bitio.ReaderAtSeeker

	readBuf *[]byte

	inArgs []any
}

// TODO: new struct decoder?
// note br is assumed to be a non-shared buffer
func newDecoder(ctx context.Context, format *Format, br bitio.ReaderAtSeeker, opts Options) *D {
	name := format.RootName
	if opts.Name != "" {
		name = opts.Name
	}
	rootV := &Compound{
		IsArray:     format.RootArray,
		Children:    nil,
		Description: opts.Description,
	}

	return &D{
		Ctx:    ctx,
		Endian: BigEndian,
		Value: &Value{
			Name:       name,
			V:          rootV,
			RootReader: br,
			Range:      ranges.Range{Start: 0, Len: 0},
			IsRoot:     opts.IsRoot,
			Format:     format,
		},
		Options: opts,

		bitBuf:  br,
		readBuf: opts.ReadBuf,
	}
}

// ArgAs looks for a decoder option with type of target, similar to errors.As
// return true if found
func (d *D) ArgAs(target any) bool {
	targeType := reflect.TypeOf(target)
	for _, in := range d.inArgs {
		// target will be &format.<Format>In
		// in format.<Format>In so have to get ptr type
		inType := reflect.PointerTo(reflect.TypeOf(in))
		if inType.AssignableTo(targeType) {
			targetVal := reflect.ValueOf(target)
			targetVal.Elem().Set(reflect.ValueOf(in))
			return true
		}
	}

	return false
}

func (d *D) fieldDecoder(name string, bitBuf bitio.ReaderAtSeeker, v any) *D {
	return &D{
		Ctx:    d.Ctx,
		Endian: d.Endian,
		Value: &Value{
			Name:       name,
			V:          v,
			Range:      ranges.Range{Start: d.Pos(), Len: 0},
			RootReader: bitBuf,
		},
		Options: d.Options,

		bitBuf:  bitBuf,
		readBuf: d.readBuf,
	}
}

func (d *D) TryCopyBits(w io.Writer, r bitio.Reader) (int64, error) {
	// TODO: what size? now same as io.Copy
	buf := d.SharedReadBuf(32 * 1024)
	return bitiox.CopyBitsBuffer(w, r, buf)
}

func (d *D) CopyBits(w io.Writer, r bitio.Reader) int64 {
	n, err := d.TryCopyBits(w, r)
	if err != nil {
		d.IOPanic(err, "", "CopyBits")
	}
	return n
}

func (d *D) TryCopy(w io.Writer, r io.Reader) (int64, error) {
	// TODO: what size? now same as io.Copy
	buf := d.SharedReadBuf(32 * 1024)
	return io.CopyBuffer(w, r, buf)
}

func (d *D) Copy(w io.Writer, r io.Reader) int64 {
	n, err := d.TryCopy(w, r)
	if err != nil {
		d.IOPanic(err, "", "Copy")
	}
	return n
}

func (d *D) CloneReadSeeker(br bitio.ReadSeeker) bitio.ReadSeeker {
	br, err := bitio.CloneReadSeeker(br)
	if err != nil {
		d.IOPanic(err, "", "CloneReadSeeker")
	}
	return br
}

func (d *D) NewBitBufFromReader(r io.Reader) bitio.ReaderAtSeeker {
	b := &bytes.Buffer{}
	d.Copy(b, r)
	return bitio.NewBitReader(b.Bytes(), -1)
}

func (d *D) TryReadAllBits(r bitio.Reader) ([]byte, error) {
	bb := &bytes.Buffer{}
	buf := d.SharedReadBuf(32 * 1024)
	if _, err := bitiox.CopyBitsBuffer(bb, r, buf); err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

func (d *D) ReadAllBits(r bitio.Reader) []byte {
	buf, err := d.TryReadAllBits(r)
	if err != nil {
		d.IOPanic(err, "", "ReadAllBytes")
	}
	return buf
}

func (d *D) SharedReadBuf(n int) []byte {
	if d.readBuf == nil {
		d.readBuf = new([]byte)
	}
	if len(*d.readBuf) < n {
		*d.readBuf = make([]byte, n)
	}
	return (*d.readBuf)[:n]
}

func (d *D) FillGaps(r ranges.Range, namePrefix string) {
	makeWalkFn := func(fn func(iv *Value)) func(iv *Value, rootV *Value, depth int, rootDepth int) error {
		return func(iv *Value, _ *Value, _ int, _ int) error {
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
	_ = d.Value.WalkRootPreOrder(makeWalkFn(func(_ *Value) { n++ }))
	valueRanges := make([]ranges.Range, n)
	i := 0
	_ = d.Value.WalkRootPreOrder(makeWalkFn(func(iv *Value) {
		valueRanges[i] = iv.Range
		i++
	}))

	gaps := ranges.Gaps(r, valueRanges)
	for i, gap := range gaps {
		br, err := bitiox.Range(d.bitBuf, gap.Start, gap.Len)
		if err != nil {
			d.IOPanic(err, namePrefix, "bitiox.Range")
		}

		v := &Value{
			Name: fmt.Sprintf("%s%d", namePrefix, i),
			V: &scalar.BitBuf{
				Actual: br,
				Flags:  scalar.FlagGap,
			},
			RootReader: d.bitBuf,
			Range:      gap,
		}

		// TODO: for arrays not great that we just append gap fields
		d.AddChild(v)
	}
}

// Errorf stops decode with a reason unless forced
func (d *D) Errorf(format string, a ...any) {
	if !d.Options.Force {
		panic(DecoderError{Reason: fmt.Sprintf(format, a...), Pos: d.Pos()})
	}
}

// Fatalf stops decode with a reason regardless of forced
func (d *D) Fatalf(format string, a ...any) {
	panic(DecoderError{Reason: fmt.Sprintf(format, a...), Pos: d.Pos()})
}

func (d *D) IOPanic(err error, name string, op string) {
	panic(IOError{Err: err, Name: name, Pos: d.Pos(), Op: op})
}

// TryBits reads nBits bits from buffer
func (d *D) TryBits(nBits int) ([]byte, error) {
	if nBits < 0 {
		return nil, fmt.Errorf("nBits must be >= 0 (%d)", nBits)
	}
	buf := d.SharedReadBuf(int(bitio.BitsByteCount(int64(nBits))))
	_, err := bitio.ReadFull(d.bitBuf, buf, int64(nBits))
	if err != nil {
		return nil, err
	}

	return buf[:], nil
}

// Bits reads nBits bits from buffer
func (d *D) Bits(nBits int) []byte {
	b, err := d.TryBits(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "Bits", ReadSize: int64(nBits), Pos: d.Pos()})
	}
	return b
}

// TryUintBits reads nBits bits as a uint64 from buffer
func (d *D) TryUintBits(nBits int) (uint64, error) {
	if nBits < 0 || nBits > 64 {
		return 0, fmt.Errorf("nBits must be 0-64 (%d)", nBits)
	}
	buf, err := d.TryBits(nBits)
	if err != nil {
		return 0, err
	}

	return bitio.Read64(buf[:], 0, int64(nBits)), nil // TODO: int64
}

// UintBits reads nBits bits as uint64 from buffer
func (d *D) UintBits(nBits int) uint64 {
	n, err := d.TryUintBits(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "UintBits", ReadSize: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) PeekUintBits(nBits int) uint64 {
	n, err := d.TryPeekBits(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekUintBits", ReadSize: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) TryPeekBytes(nBytes int) ([]byte, error) {
	start, err := d.bitBuf.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	bs, err := d.TryBytesLen(nBytes)
	if _, err := d.bitBuf.SeekBits(start, io.SeekStart); err != nil {
		return nil, err
	}
	return bs, err
}

func (d *D) PeekBytes(nBytes int) []byte {
	bs, err := d.TryPeekBytes(nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "PeekBytes", ReadSize: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
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
	n, err := d.TryUintBits(nBits)
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

func (d *D) PeekFind(nBits int, seekBits int64, maxLen int64, fn func(v uint64) bool) (int64, uint64) {
	peekBits, v, err := d.TryPeekFind(nBits, seekBits, maxLen, fn)
	if err != nil {
		d.IOPanic(err, "", "PeekFind")
	}
	if peekBits == -1 {
		d.Errorf("peek not found")
	}
	return peekBits, v
}

// BytesRange reads nBytes bytes starting bit position start
// Does not update current position.
// TODO: nBytes -1?
func (d *D) TryBytesRange(bitOffset int64, nBytes int) ([]byte, error) {
	if nBytes < 0 {
		return nil, fmt.Errorf("negative nBytes %d", nBytes)
	}
	buf := make([]byte, nBytes)
	n, err := bitio.ReadAtFull(d.bitBuf, buf, int64(nBytes)*8, bitOffset)
	if n == int64(nBytes)*8 {
		err = nil
	}
	return buf, err
}

func (d *D) BytesRange(firstBit int64, nBytes int) []byte {
	bs, err := d.TryBytesRange(firstBit, nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "BytesRange", ReadSize: int64(nBytes) * 8, Pos: firstBit})
	}
	return bs
}

func (d *D) TryBytesLen(nBytes int) ([]byte, error) {
	buf := make([]byte, nBytes)
	_, err := bitio.ReadFull(d.bitBuf, buf, int64(nBytes)*8)
	return buf, err
}

func (d *D) BytesLen(nBytes int) []byte {
	bs, err := d.TryBytesLen(nBytes)
	if err != nil {
		panic(IOError{Err: err, Op: "BytesLen", ReadSize: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return bs
}

// TODO: rename/remove BitBuf name?
func (d *D) TryBitBufRange(firstBit int64, nBits int64) (bitio.ReaderAtSeeker, error) {
	return bitiox.Range(d.bitBuf, firstBit, nBits)
}

func (d *D) BitBufRange(firstBit int64, nBits int64) bitio.ReaderAtSeeker {
	br, err := bitiox.Range(d.bitBuf, firstBit, nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "BitBufRange", ReadSize: nBits, Pos: firstBit})
	}
	return br
}

func (d *D) TryPos() (int64, error) {
	return d.bitBuf.SeekBits(0, io.SeekCurrent)
}

func (d *D) Pos() int64 {
	bPos, err := d.TryPos()
	if err != nil {
		panic(IOError{Err: err, Op: "Pos", ReadSize: 0, Pos: bPos})
	}
	return bPos
}

func (d *D) TryLen() (int64, error) {
	return bitiox.Len(d.bitBuf)
}

func (d *D) Len() int64 {
	l, err := d.TryLen()
	if err != nil {
		panic(IOError{Err: err, Op: "Len"})
	}
	return l
}

////

// BitBufLen reads nBits
func (d *D) TryBitBufLen(nBits int64) (bitio.ReaderAtSeeker, error) {
	bPos, err := d.TryPos()
	if err != nil {
		return nil, err
	}
	br, err := d.TryBitBufRange(bPos, nBits)
	if err != nil {
		return nil, err
	}
	if _, err := d.TrySeekRel(nBits); err != nil {
		return nil, err
	}

	return br, nil
}

// End is true if current position is at the end
func (d *D) TryEnd() (bool, error) {
	bPos, err := d.TryPos()
	if err != nil {
		return false, err
	}
	bLen, err := d.TryLen()
	if err != nil {
		return false, err
	}
	return bPos >= bLen, nil
}

func (d *D) End() bool {
	bEnd, err := d.TryEnd()
	if err != nil {
		panic(IOError{Err: err, Op: "End", ReadSize: 0, Pos: d.Pos()})
	}
	return bEnd
}

func (d *D) NotEnd() bool { return !d.End() }

// BitsLeft number of bits left until end
func (d *D) TryBitsLeft() (int64, error) {
	bPos, err := d.TryPos()
	if err != nil {
		return 0, err
	}
	bLen, err := d.TryLen()
	if err != nil {
		return 0, err
	}
	return bLen - bPos, nil
}

func (d *D) BitsLeft() int64 {
	bBitsLeft, err := d.TryBitsLeft()
	if err != nil {
		panic(IOError{Err: err, Op: "BitsLeft", ReadSize: 0, Pos: d.Pos()})
	}
	return bBitsLeft
}

// AlignBits number of bits to next nBits align
func (d *D) TryAlignBits(nBits int) (int, error) {
	bPos, err := d.TryPos()
	if err != nil {
		return 0, err
	}
	return int((int64(nBits) - (bPos % int64(nBits))) % int64(nBits)), nil
}

func (d *D) AlignBits(nBits int) int {
	bByteAlignBits, err := d.TryAlignBits(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "AlignBits", ReadSize: 0, Pos: d.Pos()})
	}
	return bByteAlignBits
}

// ByteAlignBits number of bits to next byte align
func (d *D) TryByteAlignBits() (int, error) {
	return d.TryAlignBits(8)
}

func (d *D) ByteAlignBits() int {
	bByteAlignBits, err := d.TryByteAlignBits()
	if err != nil {
		panic(IOError{Err: err, Op: "ByteAlignBits", ReadSize: 0, Pos: d.Pos()})
	}
	return bByteAlignBits
}

// BytePos byte position of current bit position
func (d *D) TryBytePos() (int64, error) {
	bPos, err := d.TryPos()
	if err != nil {
		return 0, err
	}
	return bPos & 0x7, nil
}

func (d *D) BytePos() int64 {
	bBytePos, err := d.TryBytePos()
	if err != nil {
		panic(IOError{Err: err, Op: "BytePos", ReadSize: 0, Pos: d.Pos()})
	}
	return bBytePos
}

func (d *D) trySeekAbs(pos int64, fns ...func(d *D)) (int64, error) {
	var oldPos int64
	if len(fns) > 0 {
		oldPos = d.Pos()
	}

	pos, err := d.bitBuf.SeekBits(pos, io.SeekStart)
	if err != nil {
		return 0, err
	}

	if len(fns) > 0 {
		for _, fn := range fns {
			fn(d)
		}
		_, err := d.bitBuf.SeekBits(oldPos, io.SeekStart)
		if err != nil {
			return 0, err
		}
	}

	return pos, nil
}

// TrySeekRel seeks relative to current bit position. If one or more functions is
// passed, the original seek position will be restored after decoding.
func (d *D) TrySeekRel(delta int64, fns ...func(d *D)) (int64, error) {
	return d.trySeekAbs(d.Pos()+delta, fns...)
}

// SeekRel seeks relative to current bit position. If one or more functions is
// passed, the original seek position will be restored after decoding. Panics if
// an error occurs while seeking.
func (d *D) SeekRel(delta int64, fns ...func(d *D)) int64 {
	n, err := d.trySeekAbs(d.Pos()+delta, fns...)
	if err != nil {
		d.IOPanic(err, "", "SeekRel")
	}
	return n
}

// TrySeekAbs seeks to absolute position. If one or more functions is passed,
// the original seek position will be restored after decoding.
func (d *D) TrySeekAbs(pos int64, fns ...func(d *D)) (int64, error) {
	return d.trySeekAbs(pos, fns...)
}

func (d *D) SeekAbs(pos int64, fns ...func(d *D)) int64 {
	n, err := d.trySeekAbs(pos, fns...)
	if err != nil {
		d.IOPanic(err, "", "SeekAbs")
	}
	return n
}

func (d *D) AddChild(v *Value) {
	v.Parent = d.Value

	switch fv := d.Value.V.(type) {
	case *Compound:
		if !fv.IsArray {
			if fv.ByName == nil {
				fv.ByName = make(map[string]*Value)
			}
			if _, ok := fv.ByName[v.Name]; ok {
				d.Fatalf("%q already exist in struct %s", v.Name, d.Value.Name)
			}
			fv.ByName[v.Name] = v
		}
		fv.Children = append(fv.Children, v)
	}
}

func (d *D) FieldGet(name string) *Value {
	switch fv := d.Value.V.(type) {
	case *Compound:
		if fv.IsArray {
			for _, ff := range fv.Children {
				if ff.Name == name {
					return ff
				}
			}
		} else {
			if fv.ByName != nil {
				if ff, ok := fv.ByName[name]; ok {
					return ff
				}
			}
			return nil
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

// FieldArray decode array of fields. Will not be range sorted.
func (d *D) FieldArray(name string, fn func(d *D)) *D {
	c := &Compound{IsArray: true}
	cd := d.fieldDecoder(name, d.bitBuf, c)
	d.AddChild(cd.Value)
	fn(cd)
	return cd
}

// FieldArrayValue decode array of fields. Will not be range sorted.
func (d *D) FieldArrayValue(name string) *D {
	return d.FieldArray(name, func(d *D) {})
}

// FieldStruct decode array of fields. Will be range sorted.
func (d *D) FieldStruct(name string, fn func(d *D)) *D {
	c := &Compound{IsArray: false}
	cd := d.fieldDecoder(name, d.bitBuf, c)
	d.AddChild(cd.Value)
	fn(cd)
	return cd
}

// FieldStructValue decode array of fields. Will be range sorted.
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

// FieldStructNArray decodes an array of elements with a known count.
func (d *D) FieldStructNArray(name string, structName string, count int64, fn func(d *D)) *D {
	return d.FieldArray(name, func(d *D) {
		for i := int64(0); i < count; i++ {
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
	v.RootReader = d.bitBuf
	v.Range = ranges.Range{Start: firstBit, Len: nBits}
	d.AddChild(v)

	return v
}

func (d *D) AssertPos(pos int64) {
	if d.Pos() != pos {
		panic(DecoderError{Reason: fmt.Sprintf("expected bits position %d", pos), Pos: d.Pos()})
	}
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

// FramedFn decode from current position nBits forward. When done position will be nBits forward.
func (d *D) FramedFn(nBits int64, fn func(d *D)) int64 {
	if nBits < 0 {
		d.Fatalf("%d nBits < 0", nBits)
	}
	decodeLen := d.RangeFn(d.Pos(), nBits, fn)
	d.SeekRel(nBits)
	return decodeLen
}

// LimitedFn decode from current position nBits forward. When done position will be after last bit decoded.
func (d *D) LimitedFn(nBits int64, fn func(d *D)) int64 {
	if nBits < 0 {
		d.Fatalf("%d nBits < 0", nBits)
	}
	decodeLen := d.RangeFn(d.Pos(), nBits, fn)
	d.SeekRel(decodeLen)
	return decodeLen
}

// RangeFn decode from firstBit position nBits forward. Position will not change.
func (d *D) RangeFn(firstBit int64, nBits int64, fn func(d *D)) int64 {
	startPos := d.Pos()

	// TODO: do some kind of DecodeLimitedLen/RangeFn?
	br := d.BitBufRange(0, firstBit+nBits)
	if _, err := br.SeekBits(firstBit, io.SeekStart); err != nil {
		d.IOPanic(err, "", "RangeFn: SeekAbs")
	}

	nd := *d
	nd.bitBuf = br

	fn(&nd)

	endPos := nd.Pos()

	return endPos - startPos
}

func (d *D) Format(group *Group, inArg any) any {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Force:       d.Options.Force,
		FillGaps:    false,
		IsRoot:      false,
		Range:       ranges.Range{Start: d.Pos(), Len: d.BitsLeft()},
		InArg:       inArg,
		ParseOptsFn: d.Options.ParseOptsFn,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, "", "Format: decode")
	}

	switch vv := dv.V.(type) {
	case *Compound:
		for _, f := range vv.Children {
			d.AddChild(f)
		}
	default:
		panic("unreachable")
	}

	if _, err := d.bitBuf.SeekBits(dv.Range.Len, io.SeekCurrent); err != nil {
		d.IOPanic(err, "", "Format: SeekRel")
	}

	return v
}

func (d *D) TryFieldFormat(name string, group *Group, inArg any) (*Value, any, error) {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    false,
		IsRoot:      false,
		Range:       ranges.Range{Start: d.Pos(), Len: d.BitsLeft()},
		InArg:       inArg,
		ParseOptsFn: d.Options.ParseOptsFn,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)
	if _, err := d.bitBuf.SeekBits(dv.Range.Len, io.SeekCurrent); err != nil {
		d.IOPanic(err, name, "TryFieldFormat: SeekRel")
	}

	return dv, v, err
}

func (d *D) FieldFormat(name string, group *Group, inArg any) (*Value, any) {
	dv, v, err := d.TryFieldFormat(name, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, name, "FieldFormat: TryFieldFormat")
	}
	return dv, v
}

func (d *D) FieldFormatOrRaw(name string, group *Group, inArg any) (*Value, any) {
	dv, v, _ := d.TryFieldFormat(name, group, inArg)
	if dv == nil {
		d.FieldRawLen(name, d.BitsLeft())
	}
	return dv, v
}

func (d *D) TryFieldFormatLen(name string, nBits int64, group *Group, inArg any) (*Value, any, error) {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    true,
		IsRoot:      false,
		Range:       ranges.Range{Start: d.Pos(), Len: nBits},
		InArg:       inArg,
		ParseOptsFn: d.Options.ParseOptsFn,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)
	if _, err := d.bitBuf.SeekBits(nBits, io.SeekCurrent); err != nil {
		d.IOPanic(err, name, "TryFieldFormatLen: SeekRel")
	}

	return dv, v, err
}

func (d *D) FieldFormatLen(name string, nBits int64, group *Group, inArg any) (*Value, any) {
	dv, v, err := d.TryFieldFormatLen(name, nBits, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, name, "FieldFormatLen: TryFieldFormatLen")
	}
	return dv, v
}

func (d *D) FieldFormatOrRawLen(name string, nBits int64, group *Group, inArg any) (*Value, any) {
	dv, v, _ := d.TryFieldFormatLen(name, nBits, group, inArg)
	if dv == nil {
		d.FieldRawLen(name, nBits)
	}
	return dv, v
}

// TODO: return decooder?
func (d *D) TryFieldFormatRange(name string, firstBit int64, nBits int64, group *Group, inArg any) (*Value, any, error) {
	dv, v, err := decode(d.Ctx, d.bitBuf, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    true,
		IsRoot:      false,
		Range:       ranges.Range{Start: firstBit, Len: nBits},
		InArg:       inArg,
		ParseOptsFn: d.Options.ParseOptsFn,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	d.AddChild(dv)

	return dv, v, err
}

func (d *D) FieldFormatRange(name string, firstBit int64, nBits int64, group *Group, inArg any) (*Value, any) {
	dv, v, err := d.TryFieldFormatRange(name, firstBit, nBits, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, name, "FieldFormatRange: TryFieldFormatRange")
	}

	return dv, v
}

func (d *D) TryFieldFormatBitBuf(name string, br bitio.ReaderAtSeeker, group *Group, inArg any) (*Value, any, error) {
	dv, v, err := decode(d.Ctx, br, group, Options{
		Name:        name,
		Force:       d.Options.Force,
		FillGaps:    true,
		IsRoot:      true,
		InArg:       inArg,
		ParseOptsFn: d.Options.ParseOptsFn,
		ReadBuf:     d.readBuf,
	})
	if dv == nil || dv.Errors() != nil {
		return nil, nil, err
	}

	dv.Range.Start = d.Pos()

	d.AddChild(dv)

	return dv, v, err
}

func (d *D) FieldFormatBitBuf(name string, br bitio.ReaderAtSeeker, group *Group, inArg any) (*Value, any) {
	dv, v, err := d.TryFieldFormatBitBuf(name, br, group, inArg)
	if dv == nil || dv.Errors() != nil {
		d.IOPanic(err, name, "FieldFormatBitBuf: TryFieldFormatBitBuf")
	}

	return dv, v
}

// TODO: rethink these

func (d *D) FieldRootBitBuf(name string, br bitio.ReaderAtSeeker, sms ...scalar.BitBufMapper) *Value {
	brLen, err := bitiox.Len(br)
	if err != nil {
		d.IOPanic(err, name, "br Len")
	}

	v := &Value{}
	v.V = &scalar.BitBuf{Actual: br}
	v.Name = name
	v.RootReader = br
	v.IsRoot = true
	v.Range = ranges.Range{Start: d.Pos(), Len: brLen}

	// if err := v.TryScalarFn(sms...); err != nil {
	// 	d.Fatalf("%v", err)
	// }

	d.AddChild(v)

	return v
}

func (d *D) FieldArrayRootBitBufFn(name string, br bitio.ReaderAtSeeker, fn func(d *D)) *Value {
	c := &Compound{IsArray: true}
	cd := d.fieldDecoder(name, br, c)
	cd.Value.IsRoot = true
	d.AddChild(cd.Value)
	fn(cd)

	cd.Value.postProcess()

	return cd.Value
}

func (d *D) FieldStructRootBitBufFn(name string, br bitio.ReaderAtSeeker, fn func(d *D)) *Value {
	c := &Compound{IsArray: false}
	cd := d.fieldDecoder(name, br, c)
	cd.Value.IsRoot = true
	d.AddChild(cd.Value)
	fn(cd)

	cd.Value.postProcess()

	return cd.Value
}

// TODO: range?
func (d *D) FieldFormatReaderLen(name string, nBits int64, fn func(r io.Reader) (io.ReadCloser, error), group *Group) (*Value, any) {
	br, err := d.TryBitBufLen(nBits)
	if err != nil {
		d.IOPanic(err, name, "FieldFormatReaderLen: BitBufLen")
	}

	bbBR := bitio.NewIOReader(br)
	r, err := fn(bbBR)
	if err != nil {
		d.IOPanic(err, name, "FieldFormatReaderLen: fn")
	}
	rBuf, err := io.ReadAll(r)
	if err != nil {
		d.IOPanic(err, name, "FieldFormatReaderLen: ReadAll")
	}
	rBR := bitio.NewBitReader(rBuf, -1)

	return d.FieldFormatBitBuf(name, rBR, group, nil)
}

// TODO: too many return values
func (d *D) TryFieldReaderRangeFormat(name string, startBit int64, nBits int64, fn func(r io.Reader) io.Reader, group *Group, inArg any) (int64, bitio.ReaderAtSeeker, *Value, any, error) {
	bitLen := nBits
	if bitLen == -1 {
		bitLen = d.BitsLeft()
	}
	br, err := d.TryBitBufRange(startBit, bitLen)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	r := bitio.NewIOReadSeeker(br)
	rb, err := io.ReadAll(fn(r))
	if err != nil {
		return 0, nil, nil, nil, err
	}
	cz, err := r.Seek(0, io.SeekCurrent)
	rbr := bitio.NewBitReader(rb, -1)
	if err != nil {
		return 0, nil, nil, nil, err
	}

	if group != nil {
		dv, v, err := d.TryFieldFormatBitBuf(name, rbr, group, inArg)
		return cz * 8, rbr, dv, v, err
	}

	d.FieldRootBitBuf(name, rbr)

	return cz * 8, rbr, nil, nil, nil
}

// TODO: refactor this
func (d *D) FieldReaderRange(name string, startBit int64, nBits int64, fn func(r io.Reader) io.Reader) (int64, bitio.ReaderAtSeeker, error) {
	sz, br, _, _, err := d.TryFieldReaderRangeFormat(name, startBit, nBits, fn, nil, nil)
	return sz, br, err
}

func (d *D) FieldReaderRangeFormat(name string, startBit int64, nBits int64, fn func(r io.Reader) io.Reader, group *Group, inArg any) (int64, bitio.ReaderAtSeeker, *Value, any) {
	cz, rBR, dv, v, err := d.TryFieldReaderRangeFormat(name, startBit, nBits, fn, group, inArg)
	if err != nil {
		d.IOPanic(err, name, "TryFieldReaderRangeFormat")
	}
	return cz, rBR, dv, v
}

func (d *D) TryFieldValue(name string, fn func() (*Value, error)) (*Value, error) {
	start := d.Pos()
	v, err := fn()
	stop := d.Pos()
	v.Name = name
	v.RootReader = d.bitBuf
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
		d.IOPanic(err, name, "FieldValue: TryFieldValue")
	}
	return v
}

func (d *D) RE(re *regexp.Regexp) []ranges.Range {
	startPos := d.Pos()

	rr := iox.ByteRuneReader{RS: bitio.NewIOReadSeeker(d.bitBuf)}
	locs := re.FindReaderSubmatchIndex(rr)
	if locs == nil {
		return nil
	}
	d.SeekAbs(startPos)

	var rs []ranges.Range
	l := len(locs) / 2
	for i := 0; i < l; i++ {
		loc := locs[i*2 : i*2+2]
		if loc[0] == -1 {
			rs = append(rs, ranges.Range{Start: -1})
		} else {
			rs = append(rs, ranges.Range{
				Start: startPos + int64(loc[0]*8),
				Len:   int64((loc[1] - loc[0]) * 8)},
			)
		}
	}

	return rs
}

func (d *D) FieldRE(re *regexp.Regexp, mRef *map[string]string, sms ...scalar.StrMapper) {
	subexpNames := re.SubexpNames()

	rs := d.RE(re)
	for i, r := range rs {
		if i == 0 || r.Start == -1 {
			continue
		}
		d.SeekAbs(r.Start)
		name := subexpNames[i]
		value := d.FieldUTF8(name, int(r.Len/8), sms...)
		if mRef != nil {
			(*mRef)[name] = value
		}
	}

	d.SeekAbs(rs[0].Stop())
}
