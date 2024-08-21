package interp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/big"

	"github.com/wader/fq/internal/aheadreadseeker"
	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/internal/ctxreadseeker"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/internal/iox"
	"github.com/wader/fq/internal/progressreadseeker"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/gojq"
)

func init() {
	RegisterFunc1("_tobits", (*Interp)._toBits)
	RegisterFunc0("open", (*Interp)._open)
}

type ToBinary interface {
	ToBinary() (Binary, error)
}

func toBinary(v any) (Binary, error) {
	switch vv := v.(type) {
	case ToBinary:
		return vv.ToBinary()
	default:
		br, err := ToBitReader(v)
		if err != nil {
			return Binary{}, err
		}
		return NewBinaryFromBitReader(br, 8, 0)
	}
}

func ToBitReader(v any) (bitio.ReaderAtSeeker, error) {
	return toBitReaderEx(v, false)
}

type byteRangeError int

func (b byteRangeError) Error() string {
	return fmt.Sprintf("byte in binary list must be bytes (0-255) got %d", int(b))

}

func toBitReaderEx(v any, inArray bool) (bitio.ReaderAtSeeker, error) {
	switch vv := v.(type) {
	case ToBinary:
		bv, err := vv.ToBinary()
		if err != nil {
			return nil, err
		}
		return bitiox.Range(bv.br, bv.r.Start, bv.r.Len)
	case string:
		return bitio.NewBitReader([]byte(vv), -1), nil
	case int, float64, *big.Int:
		bi, err := toBigInt(v)
		if err != nil {
			return nil, err
		}

		if inArray {
			if bi.Cmp(big.NewInt(255)) > 0 || bi.Cmp(big.NewInt(0)) < 0 {
				return nil, byteRangeError(bi.Int64())
			}
			n := bi.Uint64()
			b := [1]byte{byte(n)}
			return bitio.NewBitReader(b[:], -1), nil
		}

		bitLen := int64(bi.BitLen())
		// bit.Int "The bit length of 0 is 0."
		if bitLen == 0 {
			var z [1]byte
			return bitio.NewBitReader(z[:], 1), nil
		}
		// TODO: how should this work? "0xf | tobytes" 4bits or 8bits? now 4
		padBefore := (8 - (bitLen % 8)) % 8
		// padBefore := 0
		br, err := bitiox.Range(bitio.NewBitReader(bi.Bytes(), -1), padBefore, bitLen)
		if err != nil {
			return nil, err
		}
		return br, nil
	case []any:
		rr := make([]bitio.ReadAtSeeker, 0, len(vv))

		// fast path for slice containing only 0-255 numbers and strings
		bs := &bytes.Buffer{}
		for _, e := range vv {
			if bs == nil {
				break
			}
			switch ev := e.(type) {
			case int:
				if ev >= 0 && ev <= 255 {
					bs.WriteByte(byte(ev))
					continue
				}
			case float64:
				b := int(ev)
				if b >= 0 && b <= 255 {
					bs.WriteByte(byte(ev))
					continue
				}
			case *big.Int:
				if ev.Cmp(big.NewInt(0)) >= 0 && ev.Cmp(big.NewInt(255)) <= 0 {
					bs.WriteByte(byte(ev.Uint64()))
					continue
				}
			case string:
				// TODO: maybe only if less then some length?
				bs.WriteString(ev)
				continue
			}
			bs = nil
		}
		if bs != nil {
			return bitio.NewBitReader(bs.Bytes(), -1), nil
		}

		// TODO: optimize byte array case, flatten into one slice
		for _, e := range vv {
			eBR, eErr := toBitReaderEx(e, true)
			if eErr != nil {
				return nil, eErr
			}
			rr = append(rr, eBR)
		}

		mb, err := bitio.NewMultiReader(rr...)
		if err != nil {
			return nil, err
		}

		return mb, nil
	default:
		return nil, fmt.Errorf("value can't be a binary")
	}
}

type toBitsOpts struct {
	Unit       int
	KeepRange  bool
	PadToUnits int
}

// note is used to implement tobytes* also
func (i *Interp) _toBits(c any, opts toBitsOpts) any {
	// TODO: unit > 8?

	bv, err := toBinary(c)
	if err != nil {
		return err
	}

	pad := int64(opts.Unit * opts.PadToUnits)
	if pad == 0 {
		pad = int64(opts.Unit)
	}

	bv.unit = opts.Unit
	bv.pad = (pad - bv.r.Len%pad) % pad

	if opts.KeepRange {
		return bv
	}

	br, err := bv.toReader()
	if err != nil {
		return err
	}
	bb, err := NewBinaryFromBitReader(br, bv.unit, 0)
	if err != nil {
		return err
	}
	return bb
}

type openFile struct {
	Binary
	filename   string
	progressFn progressreadseeker.ProgressFn
}

var _ Value = (*openFile)(nil)
var _ ToBinary = (*openFile)(nil)

func (of *openFile) Display(w io.Writer, opts *Options) error {
	_, err := fmt.Fprintf(w, "<openfile %q>\n", of.filename)
	return err
}

func (of *openFile) ToBinary() (Binary, error) {
	return NewBinaryFromBitReader(of.br, 8, 0)
}

// opens a file for reading from filesystem
// TODO: when to close? when br loses all refs? need to use finalizer somehow?
func (i *Interp) _open(c any) any {
	if i.EvalInstance.IsCompleting {
		// TODO: have dummy values for each type for completion?
		br, _ := NewBinaryFromBitReader(bitio.NewBitReader([]byte{}, -1), 8, 0)
		return br
	}

	var err error
	var f fs.File
	var path string

	switch c.(type) {
	case nil:
		path = "<stdin>"
		f = i.OS.Stdin()
	default:
		path, err = toString(c)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		f, err = i.OS.FS().Open(path)
		if err != nil {
			// path context added in jq error code
			var pe *fs.PathError
			if errors.As(err, &pe) {
				return pe.Err
			}
			return err
		}
	}

	var bEnd int64
	var fRS io.ReadSeeker

	fFI, err := f.Stat()
	if err != nil {
		f.Close()
		return err
	}

	// ctxreadseeker is used to make sure any io calls can be canceled
	// TODO: ctxreadseeker might leak if the underlying call hangs forever

	// a regular file should be seekable but fallback below to read whole file if not
	if fFI.Mode().IsRegular() {
		if rs, ok := f.(io.ReadSeeker); ok {
			fRS = ctxreadseeker.New(i.EvalInstance.Ctx, rs)
			bEnd = fFI.Size()
		}
	}

	if fRS == nil {
		buf, err := io.ReadAll(ctxreadseeker.New(i.EvalInstance.Ctx, &iox.ReadErrSeeker{Reader: f}))
		if err != nil {
			f.Close()
			return err
		}
		fRS = bytes.NewReader(buf)
		bEnd = int64(len(buf))
	}

	bbf := &openFile{
		filename: path,
	}

	const progressPrecision = 1024
	fRS = progressreadseeker.New(fRS, progressPrecision, bEnd,
		func(approxReadBytes int64, totalSize int64) {
			// progressFn is assign by decode etc
			if bbf.progressFn != nil {
				bbf.progressFn(approxReadBytes, totalSize)
			}
		},
	)

	const cacheReadAheadSize = 512 * 1024
	aheadRs := aheadreadseeker.New(fRS, cacheReadAheadSize)

	// bitio.Buffer -> (bitio.Reader) -> aheadreadseeker -> progressreadseeker -> ctxreadseeker -> readseeker

	bbf.br = bitio.NewIOBitReadSeeker(aheadRs)

	return bbf
}

var _ Value = Binary{}
var _ ToBinary = Binary{}

type Binary struct {
	br   bitio.ReaderAtSeeker
	r    ranges.Range
	unit int
	pad  int64
}

func NewBinaryFromBitReader(br bitio.ReaderAtSeeker, unit int, pad int64) (Binary, error) {
	l, err := bitiox.Len(br)
	if err != nil {
		return Binary{}, err
	}

	return Binary{
		br:   br,
		r:    ranges.Range{Start: 0, Len: l},
		unit: unit,
		pad:  pad,
	}, nil
}

func (b Binary) toBytesBuffer(r ranges.Range) (*bytes.Buffer, error) {
	br, err := bitiox.Range(b.br, r.Start, r.Len)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if _, err := bitiox.CopyBits(buf, br); err != nil {
		return nil, err
	}

	return buf, nil
}

func (Binary) ExtType() string { return "binary" }

func (Binary) ExtKeys() []string {
	return []string{
		"bits",
		"bytes",
		"name",
		"size",
		"start",
		"stop",
		"unit",
	}
}

func (b Binary) ToBinary() (Binary, error) {
	return b, nil
}

func (b Binary) JQValueLength() any {
	return int(b.r.Len / int64(b.unit))
}
func (b Binary) JQValueSliceLen() any {
	return b.JQValueLength()
}

func (b Binary) JQValueIndex(index int) any {
	if index < 0 {
		return nil
	}

	buf, err := b.toBytesBuffer(ranges.Range{Start: b.r.Start + int64(index*b.unit), Len: int64(b.unit)})
	if err != nil {
		return err
	}

	extraBits := uint((8 - b.unit%8) % 8)

	return new(big.Int).Rsh(new(big.Int).SetBytes(buf.Bytes()), extraBits)
}
func (b Binary) JQValueSlice(start int, end int) any {
	rStart := int64(start * b.unit)
	rLen := int64((end - start) * b.unit)

	return Binary{
		br:   b.br,
		r:    ranges.Range{Start: b.r.Start + rStart, Len: rLen},
		unit: b.unit,
	}
}
func (b Binary) JQValueKey(name string) any {
	switch name {
	case "bits":
		if b.unit == 1 {
			return b
		}
		return Binary{br: b.br, r: b.r, unit: 1}
	case "bytes":
		if b.unit == 8 {
			return b
		}
		return Binary{br: b.br, r: b.r, unit: 8}

	case "name":
		f := iox.Unwrap(b.br)
		// this exploits the fact that *os.File has Name()
		if n, ok := f.(interface{ Name() string }); ok {
			return n.Name()
		}
		return nil
	case "size":
		return new(big.Int).SetInt64(b.r.Len / int64(b.unit))
	case "start":
		return new(big.Int).SetInt64(b.r.Start / int64(b.unit))
	case "stop":
		stop := b.r.Stop()
		stopUnits := stop / int64(b.unit)
		if stop%int64(b.unit) != 0 {
			stopUnits++
		}
		return new(big.Int).SetInt64(stopUnits)
	case "unit":
		return b.unit
	}
	return nil
}
func (b Binary) JQValueEach() any {
	return nil
}
func (b Binary) JQValueType() string {
	return gojq.JQTypeString
}
func (b Binary) JQValueKeys() any {
	return gojqx.FuncTypeNameError{Name: "keys", Typ: gojq.JQTypeString}
}
func (b Binary) JQValueHas(key any) any {
	return gojqx.HasKeyTypeError{L: gojq.JQTypeString, R: fmt.Sprintf("%v", key)}
}
func (b Binary) JQValueToNumber() any {
	buf, err := b.toBytesBuffer(b.r)
	if err != nil {
		return err
	}
	extraBits := uint((8 - b.r.Len%8) % 8)
	return new(big.Int).Rsh(new(big.Int).SetBytes(buf.Bytes()), extraBits)
}
func (b Binary) JQValueToString() any {
	return b.JQValueToGoJQ()
}
func (b Binary) JQValueToGoJQEx(optsFn func() (*Options, error)) any {
	br, err := b.toReader()
	if err != nil {
		return err
	}

	brC, err := bitio.CloneReaderAtSeeker(br)
	if err != nil {
		return err
	}

	opts, err := optsFn()
	if err != nil {
		return err
	}

	s, err := opts.BitsFormatFn(brC)
	if err != nil {
		return err
	}

	return s
}

func (b Binary) JQValueToGoJQ() any {
	buf, err := b.toBytesBuffer(b.r)
	if err != nil {
		return err
	}
	return buf.String()
}

func (b Binary) Display(w io.Writer, opts *Options) error {
	if opts.RawOutput {
		br, err := b.toReader()
		if err != nil {
			return err
		}

		if _, err := bitiox.CopyBits(w, br); err != nil {
			return err
		}

		return nil
	}

	return hexdump(w, b, opts)
}

func (b Binary) toReader() (bitio.ReaderAtSeeker, error) {
	br, err := bitiox.Range(b.br, b.r.Start, b.r.Len)
	if err != nil {
		return nil, err
	}
	if b.pad == 0 {
		return br, nil
	}
	return bitio.NewMultiReader(bitiox.NewZeroAtSeeker(b.pad), br)
}
