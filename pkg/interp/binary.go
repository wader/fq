package interp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math/big"

	"github.com/wader/fq/internal/aheadreadseeker"
	"github.com/wader/fq/internal/bitioextra"
	"github.com/wader/fq/internal/ctxreadseeker"
	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/internal/ioextra"
	"github.com/wader/fq/internal/progressreadseeker"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
	"github.com/wader/gojq"
)

func init() {
	functionRegisterFns = append(functionRegisterFns, func(i *Interp) []Function {
		return []Function{
			{"_tobits", 3, 3, i._toBits, nil},
			{"open", 0, 0, nil, i._open},
		}
	})
}

type ToBinary interface {
	ToBinary() (Binary, error)
}

func toBinary(v any) (Binary, error) {
	switch vv := v.(type) {
	case ToBinary:
		return vv.ToBinary()
	default:
		br, err := toBitReader(v)
		if err != nil {
			return Binary{}, err
		}
		return newBinaryFromBitReader(br, 8, 0)
	}
}

func toBitReader(v any) (bitio.ReaderAtSeeker, error) {
	return toBitReaderEx(v, false)
}

func toBitReaderEx(v any, inArray bool) (bitio.ReaderAtSeeker, error) {
	switch vv := v.(type) {
	case ToBinary:
		bv, err := vv.ToBinary()
		if err != nil {
			return nil, err
		}
		return bitioextra.Range(bv.br, bv.r.Start, bv.r.Len)
	case string:
		return bitio.NewBitReader([]byte(vv), -1), nil
	case int, float64, *big.Int:
		bi, err := toBigInt(v)
		if err != nil {
			return nil, err
		}

		if inArray {
			if bi.Cmp(big.NewInt(255)) > 0 || bi.Cmp(big.NewInt(0)) < 0 {
				return nil, fmt.Errorf("byte in binary list must be bytes (0-255) got %v", bi)
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
		br, err := bitioextra.Range(bitio.NewBitReader(bi.Bytes(), -1), padBefore, bitLen)
		if err != nil {
			return nil, err
		}
		return br, nil
	case []any:
		rr := make([]bitio.ReadAtSeeker, 0, len(vv))
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

// note is used to implement tobytes* also
func (i *Interp) _toBits(c any, a []any) any {
	unit, ok := gojqextra.ToInt(a[0])
	if !ok {
		return gojqextra.FuncTypeError{Name: "_tobits", V: a[0]}
	}
	keepRange, ok := gojqextra.ToBoolean(a[1])
	if !ok {
		return gojqextra.FuncTypeError{Name: "_tobits", V: a[1]}
	}
	padToUnits, ok := gojqextra.ToInt(a[2])
	if !ok {
		return gojqextra.FuncTypeError{Name: "_tobits", V: a[2]}
	}

	// TODO: unit > 8?

	bv, err := toBinary(c)
	if err != nil {
		return err
	}

	pad := int64(unit * padToUnits)
	if pad == 0 {
		pad = int64(unit)
	}

	bv.unit = unit
	bv.pad = (pad - bv.r.Len%pad) % pad

	if keepRange {
		return bv
	}

	br, err := bv.toReader()
	if err != nil {
		return err
	}
	bb, err := newBinaryFromBitReader(br, bv.unit, 0)
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

func (of *openFile) Display(w io.Writer, opts Options) error {
	_, err := fmt.Fprintf(w, "<openfile %q>\n", of.filename)
	return err
}

func (of *openFile) ToBinary() (Binary, error) {
	return newBinaryFromBitReader(of.br, 8, 0)
}

// opens a file for reading from filesystem
// TODO: when to close? when br loses all refs? need to use finalizer somehow?
func (i *Interp) _open(c any, a []any) gojq.Iter {
	if i.evalInstance.isCompleting {
		return gojq.NewIter()
	}

	var err error
	var f fs.File
	var path string

	switch c.(type) {
	case nil:
		path = "<stdin>"
		f = i.os.Stdin()
	default:
		path, err = toString(c)
		if err != nil {
			return gojq.NewIter(fmt.Errorf("%s: %w", path, err))
		}
		f, err = i.os.FS().Open(path)
		if err != nil {
			// path context added in jq error code
			var pe *fs.PathError
			if errors.As(err, &pe) {
				return gojq.NewIter(pe.Err)
			}
			return gojq.NewIter(err)
		}
	}

	var bEnd int64
	var fRS io.ReadSeeker

	fFI, err := f.Stat()
	if err != nil {
		f.Close()
		return gojq.NewIter(err)
	}

	// ctxreadseeker is used to make sure any io calls can be canceled
	// TODO: ctxreadseeker might leak if the underlaying call hangs forever

	// a regular file should be seekable but fallback below to read whole file if not
	if fFI.Mode().IsRegular() {
		if rs, ok := f.(io.ReadSeeker); ok {
			fRS = ctxreadseeker.New(i.evalInstance.ctx, rs)
			bEnd = fFI.Size()
		}
	}

	if fRS == nil {
		buf, err := ioutil.ReadAll(ctxreadseeker.New(i.evalInstance.ctx, &ioextra.ReadErrSeeker{Reader: f}))
		if err != nil {
			f.Close()
			return gojq.NewIter(err)
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
	if err != nil {
		return gojq.NewIter(err)
	}

	return gojq.NewIter(bbf)
}

var _ Value = Binary{}
var _ ToBinary = Binary{}

type Binary struct {
	br   bitio.ReaderAtSeeker
	r    ranges.Range
	unit int
	pad  int64
}

func newBinaryFromBitReader(br bitio.ReaderAtSeeker, unit int, pad int64) (Binary, error) {
	l, err := bitioextra.Len(br)
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
	br, err := bitioextra.Range(b.br, r.Start, r.Len)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if _, err := bitioextra.CopyBits(buf, br); err != nil {
		return nil, err
	}

	return buf, nil
}

func (Binary) ExtType() string { return "binary" }

func (Binary) ExtKeys() []string {
	return []string{
		"size",
		"start",
		"stop",
		"bits",
		"bytes",
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
	}
	return nil
}
func (b Binary) JQValueEach() any {
	return nil
}
func (b Binary) JQValueType() string {
	return "binary"
}
func (b Binary) JQValueKeys() any {
	return gojqextra.FuncTypeNameError{Name: "keys", Typ: "binary"}
}
func (b Binary) JQValueHas(key any) any {
	return gojqextra.HasKeyTypeError{L: "binary", R: fmt.Sprintf("%v", key)}
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
func (b Binary) JQValueToGoJQ() any {
	buf, err := b.toBytesBuffer(b.r)
	if err != nil {
		return err
	}
	return buf.String()
}

func (b Binary) Display(w io.Writer, opts Options) error {
	if opts.RawOutput {
		br, err := b.toReader()
		if err != nil {
			return err
		}

		if _, err := bitioextra.CopyBits(w, br); err != nil {
			return err
		}

		return nil
	}

	return hexdump(w, b, opts)
}

func (b Binary) toReader() (bitio.ReaderAtSeeker, error) {
	br, err := bitioextra.Range(b.br, b.r.Start, b.r.Len)
	if err != nil {
		return nil, err
	}
	if b.pad == 0 {
		return br, nil
	}
	return bitio.NewMultiReader(bitioextra.NewZeroAtSeeker(b.pad), br)
}
