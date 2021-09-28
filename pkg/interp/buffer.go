package interp

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/ranges"
)

var _ Value = (*BufferView)(nil)
var _ ToBuffer = (*BufferView)(nil)

type BufferView struct {
	bb   *bitio.Buffer
	r    ranges.Range
	unit int
}

// TODO: JQArray

func newBifBufObject(bb *bitio.Buffer, unit int) BufferView {
	return BufferView{
		bb:   bb,
		r:    ranges.Range{Start: 0, Len: bb.Len()},
		unit: unit,
	}
}

func (*BufferView) DisplayName() string { return "buffer" }
func (*BufferView) ExtKeys() []string {
	return []string{
		"size",
		"start",
		"stop",
		"bits",
		"bytes",
	}
}

func (bo BufferView) JQValueLength() interface{} {
	return int(bo.r.Len / int64(bo.unit))
}
func (bo BufferView) JQValueSliceLen() interface{} {
	return bo.JQValueLength()
}
func (bo BufferView) JQValueIndex(index int) interface{} {
	// TODO: use bitio
	/*
		pos, err := bo.bbr.bb.Pos()
		if err != nil {
			return err
		}
		if _, err := bo.bbr.bb.SeekAbs(int64(index) * int64(bo.unit)); err != nil {
			return err
		}
		v, err := bo.bbr.bb.U(bo.unit)
		if err != nil {
			return err
		}
		if _, err := bo.bbr.bb.SeekAbs(pos); err != nil {
			return err
		}
		return int(v)
	*/
	return nil
}
func (bo BufferView) JQValueSlice(start int, end int) interface{} {
	rStart := int64(start * bo.unit)
	rLen := int64((end - start) * bo.unit)

	return BufferView{
		bb:   bo.bb,
		r:    ranges.Range{Start: bo.r.Start + rStart, Len: rLen},
		unit: bo.unit,
	}
}
func (bo BufferView) JQValueKey(name string) interface{} {
	switch name {
	case "size":
		return new(big.Int).SetInt64(bo.r.Len / int64(bo.unit))
	case "start":
		return new(big.Int).SetInt64(bo.r.Start / int64(bo.unit))
	case "stop":
		stop := bo.r.Stop()
		stopUnits := stop / int64(bo.unit)
		if stop%int64(bo.unit) != 0 {
			stopUnits++
		}
		return new(big.Int).SetInt64(stopUnits)
	case "bits":
		if bo.unit == 1 {
			return bo
		}
		return BufferView{bb: bo.bb, r: bo.r, unit: 1}
	case "bytes":
		if bo.unit == 8 {
			return bo
		}
		return BufferView{bb: bo.bb, r: bo.r, unit: 8}
	}
	return nil
}
func (bo BufferView) JQValueEach() interface{} {
	return nil
}
func (bo BufferView) JQValueType() string {
	return "buffer"
}
func (bo BufferView) JQValueKeys() interface{} {
	return gojqextra.FuncTypeError{Name: "keys", Typ: "buffer"}
}
func (bo BufferView) JQValueHas(key interface{}) interface{} {
	return gojqextra.HasKeyTypeError{L: "buffer", R: fmt.Sprintf("%v", key)}
}
func (bo BufferView) JQValueToNumber() interface{} {
	bb, err := bo.ToBuffer()
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bb.Copy()); err != nil {
		return err
	}
	extraBits := uint((8 - bo.r.Len%8) % 8)
	return new(big.Int).Rsh(new(big.Int).SetBytes(buf.Bytes()), extraBits)
}
func (bo BufferView) JQValueToString() interface{} {
	return bo.JQValueToGoJQ()
}

func (bo BufferView) JQValueToGoJQ() interface{} {
	bb, err := bo.ToBuffer()
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bb.Copy()); err != nil {
		return err
	}
	return buf.String()
}

func (bo BufferView) JQValueUpdate(key interface{}, u interface{}, delpath bool) interface{} {
	return notUpdateableError{Key: fmt.Sprintf("%v", key), Typ: "buffer"}
}

func (bo BufferView) Display(w io.Writer, opts Options) error {
	if opts.RawOutput {
		bb, err := bo.ToBuffer()
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, bb.Copy()); err != nil {
			return err
		}
		return nil
	}

	unitNames := map[int]string{
		1: "bits",
		8: "bytes",
	}
	unitName := unitNames[bo.unit]
	if unitName == "" {
		unitName = "units"
	}

	// TODO: hack
	return dump(
		&decode.Value{
			Range:       bo.r,
			RootBitBuf:  bo.bb.Copy(),
			Description: fmt.Sprintf("%d %s", bo.r.Len/int64(bo.unit), unitName),
		},
		w,
		opts,
	)
}

func (bo BufferView) ToBuffer() (*bitio.Buffer, error) {
	return bo.bb.BitBufRange(bo.r.Start, bo.r.Len)
}
