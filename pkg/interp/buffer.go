package interp

import (
	"bytes"
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
)

func init() {
	functionRegisterFns = append(functionRegisterFns, func(i *Interp) []Function {
		return []Function{
			{"_tobitsrange", 0, 2, i._toBitsRange, nil},
			{"open", 0, 0, i._open, nil},
		}
	})
}

type ToBuffer interface {
	ToBuffer() (Buffer, error)
}

func toBitBuf(v interface{}) (bitio.ReaderAtSeeker, error) {
	return toBitBufEx(v, false)
}

func toBitBufEx(v interface{}, inArray bool) (bitio.ReaderAtSeeker, error) {
	switch vv := v.(type) {
	case ToBuffer:
		bv, err := vv.ToBuffer()
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
				return nil, fmt.Errorf("buffer byte list must be bytes (0-255) got %v", bi)
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
	case []interface{}:
		rr := make([]bitio.ReadAtSeeker, 0, len(vv))
		// TODO: optimize byte array case, flatten into one slice
		for _, e := range vv {
			eBR, eErr := toBitBufEx(e, true)
			if eErr != nil {
				return nil, eErr
			}
			rr = append(rr, eBR)
		}

		mb, err := bitio.NewMultiBitReader(rr...)
		if err != nil {
			return nil, err
		}

		return mb, nil
	default:
		return nil, fmt.Errorf("value can't be a buffer")
	}
}

func toBuffer(v interface{}) (Buffer, error) {
	switch vv := v.(type) {
	case ToBuffer:
		return vv.ToBuffer()
	default:
		br, err := toBitBuf(v)
		if err != nil {
			return Buffer{}, err
		}
		return newBufferFromBuffer(br, 8)
	}
}

// note is used to implement tobytes*/0 also
func (i *Interp) _toBitsRange(c interface{}, a []interface{}) interface{} {
	var unit int
	var r bool
	var ok bool

	if len(a) >= 1 {
		unit, ok = gojqextra.ToInt(a[0])
		if !ok {
			return gojqextra.FuncTypeError{Name: "_tobitsrange", V: a[0]}
		}
	} else {
		unit = 1
	}

	if len(a) >= 2 {
		r, ok = gojqextra.ToBoolean(a[1])
		if !ok {
			return gojqextra.FuncTypeError{Name: "_tobitsrange", V: a[1]}
		}
	} else {
		r = true
	}

	// TODO: unit > 8?

	bv, err := toBuffer(c)
	if err != nil {
		return err
	}
	bv.unit = unit

	if !r {
		br, err := bv.toBuffer()
		if err != nil {
			return err
		}
		bb, err := newBufferFromBuffer(br, unit)
		if err != nil {
			return err
		}
		return bb
	}

	return bv
}

type openFile struct {
	Buffer
	filename   string
	progressFn progressreadseeker.ProgressFn
}

var _ Value = (*openFile)(nil)
var _ ToBuffer = (*openFile)(nil)

func (of *openFile) Display(w io.Writer, opts Options) error {
	_, err := fmt.Fprintf(w, "<openfile %q>\n", of.filename)
	return err
}

func (of *openFile) ToBuffer() (Buffer, error) {
	return newBufferFromBuffer(of.br, 8)
}

// def open: #:: string| => buffer
// opens a file for reading from filesystem
// TODO: when to close? when br loses all refs? need to use finalizer somehow?
func (i *Interp) _open(c interface{}, a []interface{}) interface{} {
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
			return fmt.Errorf("%s: %w", path, err)
		}
		f, err = i.os.FS().Open(path)
		if err != nil {
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
	// TODO: ctxreadseeker might leak if the underlaying call hangs forever

	// a regular file should be seekable but fallback below to read whole file if not
	if fFI.Mode().IsRegular() {
		if rs, ok := f.(io.ReadSeeker); ok {
			fRS = ctxreadseeker.New(i.evalContext.ctx, rs)
			bEnd = fFI.Size()
		}
	}

	if fRS == nil {
		buf, err := ioutil.ReadAll(ctxreadseeker.New(i.evalContext.ctx, &ioextra.ReadErrSeeker{Reader: f}))
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
	if err != nil {
		return err
	}

	return bbf
}

var _ Value = Buffer{}
var _ ToBuffer = Buffer{}

type Buffer struct {
	br   bitio.ReaderAtSeeker
	r    ranges.Range
	unit int
}

func newBufferFromBuffer(br bitio.ReaderAtSeeker, unit int) (Buffer, error) {
	l, err := bitioextra.Len(br)
	if err != nil {
		return Buffer{}, err
	}

	return Buffer{
		br:   br,
		r:    ranges.Range{Start: 0, Len: l},
		unit: unit,
	}, nil
}

func (b Buffer) toBytesBuffer(r ranges.Range) (*bytes.Buffer, error) {
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

func (Buffer) ExtType() string { return "buffer" }

func (Buffer) ExtKeys() []string {
	return []string{
		"size",
		"start",
		"stop",
		"bits",
		"bytes",
	}
}

func (b Buffer) ToBuffer() (Buffer, error) {
	return b, nil
}

func (b Buffer) JQValueLength() interface{} {
	return int(b.r.Len / int64(b.unit))
}
func (b Buffer) JQValueSliceLen() interface{} {
	return b.JQValueLength()
}

func (b Buffer) JQValueIndex(index int) interface{} {
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
func (b Buffer) JQValueSlice(start int, end int) interface{} {
	rStart := int64(start * b.unit)
	rLen := int64((end - start) * b.unit)

	return Buffer{
		br:   b.br,
		r:    ranges.Range{Start: b.r.Start + rStart, Len: rLen},
		unit: b.unit,
	}
}
func (b Buffer) JQValueKey(name string) interface{} {
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
		return Buffer{br: b.br, r: b.r, unit: 1}
	case "bytes":
		if b.unit == 8 {
			return b
		}
		return Buffer{br: b.br, r: b.r, unit: 8}
	}
	return nil
}
func (b Buffer) JQValueEach() interface{} {
	return nil
}
func (b Buffer) JQValueType() string {
	return "buffer"
}
func (b Buffer) JQValueKeys() interface{} {
	return gojqextra.FuncTypeNameError{Name: "keys", Typ: "buffer"}
}
func (b Buffer) JQValueHas(key interface{}) interface{} {
	return gojqextra.HasKeyTypeError{L: "buffer", R: fmt.Sprintf("%v", key)}
}
func (b Buffer) JQValueToNumber() interface{} {
	buf, err := b.toBytesBuffer(b.r)
	if err != nil {
		return err
	}
	extraBits := uint((8 - b.r.Len%8) % 8)
	return new(big.Int).Rsh(new(big.Int).SetBytes(buf.Bytes()), extraBits)
}
func (b Buffer) JQValueToString() interface{} {
	return b.JQValueToGoJQ()
}
func (b Buffer) JQValueToGoJQ() interface{} {
	buf, err := b.toBytesBuffer(b.r)
	if err != nil {
		return err
	}
	return buf.String()
}
func (b Buffer) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return gojqextra.NonUpdatableTypeError{Key: fmt.Sprintf("%v", key), Typ: "buffer"}
}

func (b Buffer) Display(w io.Writer, opts Options) error {
	if opts.RawOutput {
		br, err := b.toBuffer()
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

func (b Buffer) toBuffer() (bitio.ReaderAtSeeker, error) {
	return bitioextra.Range(b.br, b.r.Start, b.r.Len)
}
