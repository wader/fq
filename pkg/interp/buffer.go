package interp

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
)

var _ Value = BufferView{}
var _ ToBufferView = BufferView{}

type BufferView struct {
	bb   *bitio.Buffer
	r    ranges.Range
	unit int
}

func bufferViewFromBuffer(bb *bitio.Buffer, unit int) BufferView {
	return BufferView{
		bb:   bb,
		r:    ranges.Range{Start: 0, Len: bb.Len()},
		unit: unit,
	}
}

func (bv BufferView) toBytesBuffer(r ranges.Range) (*bytes.Buffer, error) {
	bb, err := bv.bb.BitBufRange(r.Start, r.Len)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bb.Copy()); err != nil {
		return nil, err
	}
	return buf, nil
}

func (BufferView) DisplayName() string { return "buffer" }
func (BufferView) ExtKeys() []string {
	return []string{
		"size",
		"start",
		"stop",
		"bits",
		"bytes",
	}
}

func (bv BufferView) ToBufferView() (BufferView, error) {
	return bv, nil
}

func (bv BufferView) JQValueLength() interface{} {
	return int(bv.r.Len / int64(bv.unit))
}
func (bv BufferView) JQValueSliceLen() interface{} {
	return bv.JQValueLength()
}

func (bv BufferView) JQValueIndex(index int) interface{} {
	if index < 0 {
		return ""
	}
	buf, err := bv.toBytesBuffer(ranges.Range{Start: bv.r.Start + int64(index*bv.unit), Len: int64(bv.unit)})
	if err != nil {
		return err
	}
	s := buf.String()
	return s[0:1]
}
func (bv BufferView) JQValueSlice(start int, end int) interface{} {
	rStart := int64(start * bv.unit)
	rLen := int64((end - start) * bv.unit)

	return BufferView{
		bb:   bv.bb,
		r:    ranges.Range{Start: bv.r.Start + rStart, Len: rLen},
		unit: bv.unit,
	}
}
func (bv BufferView) JQValueKey(name string) interface{} {
	switch name {
	case "size":
		return new(big.Int).SetInt64(bv.r.Len / int64(bv.unit))
	case "start":
		return new(big.Int).SetInt64(bv.r.Start / int64(bv.unit))
	case "stop":
		stop := bv.r.Stop()
		stopUnits := stop / int64(bv.unit)
		if stop%int64(bv.unit) != 0 {
			stopUnits++
		}
		return new(big.Int).SetInt64(stopUnits)
	case "bits":
		if bv.unit == 1 {
			return bv
		}
		return BufferView{bb: bv.bb, r: bv.r, unit: 1}
	case "bytes":
		if bv.unit == 8 {
			return bv
		}
		return BufferView{bb: bv.bb, r: bv.r, unit: 8}
	}
	return nil
}
func (bv BufferView) JQValueEach() interface{} {
	return nil
}
func (bv BufferView) JQValueType() string {
	return "buffer"
}
func (bv BufferView) JQValueKeys() interface{} {
	return gojqextra.FuncTypeError{Name: "keys", Typ: "buffer"}
}
func (bv BufferView) JQValueHas(key interface{}) interface{} {
	return gojqextra.HasKeyTypeError{L: "buffer", R: fmt.Sprintf("%v", key)}
}
func (bv BufferView) JQValueToNumber() interface{} {
	buf, err := bv.toBytesBuffer(bv.r)
	if err != nil {
		return err
	}
	extraBits := uint((8 - bv.r.Len%8) % 8)
	return new(big.Int).Rsh(new(big.Int).SetBytes(buf.Bytes()), extraBits)
}
func (bv BufferView) JQValueToString() interface{} {
	return bv.JQValueToGoJQ()
}
func (bv BufferView) JQValueToGoJQ() interface{} {
	buf, err := bv.toBytesBuffer(bv.r)
	if err != nil {
		return err
	}
	return buf.String()
}
func (bv BufferView) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "buffer"}
}

func (bv BufferView) Display(w io.Writer, opts Options) error {
	if opts.RawOutput {
		bb, err := bv.toBuffer()
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, bb.Copy()); err != nil {
			return err
		}
		return nil
	}

	return hexdump(w, bv, opts)
}

func (bv BufferView) toBuffer() (*bitio.Buffer, error) {
	return bv.bb.BitBufRange(bv.r.Start, bv.r.Len)
}
