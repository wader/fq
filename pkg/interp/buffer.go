package interp

import (
	"bytes"
	"fmt"
	"fq/pkg/bitio"
	"fq/pkg/ranges"
	"io"
	"math/big"
)

type bufferRange struct {
	bb *bitio.Buffer
	r  ranges.Range
}

var _ InterpObject = (*bufferObject)(nil)
var _ ToBuffer = (*bufferObject)(nil)

type bufferObject struct {
	bbr  bufferRange
	unit int
}

func newBifBufObject(bb *bitio.Buffer, unit int) *bufferObject {
	return &bufferObject{
		bbr:  bufferRange{bb: bb, r: ranges.Range{Start: 0, Len: bb.Len()}},
		unit: unit,
	}
}

func (*bufferObject) DisplayName() string { return "buffer" }
func (*bufferObject) ExtValueKeys() []string {
	return []string{
		"size",
		"start",
		"stop",
		"bits",
		"bytes",
	}
}

func (bo *bufferObject) JQValueLength() interface{} {
	return int(bo.bbr.r.Len / int64(bo.unit))
}
func (bo *bufferObject) JQValueIndex(index int) interface{} {
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
func (bo *bufferObject) JQValueSlice(start int, end int) interface{} {
	rStart := int64(start * bo.unit)
	rLen := int64((end - start) * bo.unit)

	rbb, err := bo.bbr.bb.BitBufRange(rStart, rLen)
	if err != nil {
		return err
	}

	return &bufferObject{
		bbr:  bufferRange{bb: rbb, r: ranges.Range{Start: 0, Len: rLen}},
		unit: bo.unit,
	}
}
func (bo *bufferObject) JQValueProperty(name string) interface{} {
	switch name {
	case "size":
		return new(big.Int).SetInt64(bo.bbr.r.Len / int64(bo.unit))
	case "start":
		return new(big.Int).SetInt64(bo.bbr.r.Start / int64(bo.unit))
	case "stop":
		stop := bo.bbr.r.Stop()
		stopUnits := stop / int64(bo.unit)
		if stop%int64(bo.unit) != 0 {
			stopUnits++
		}
		return new(big.Int).SetInt64(stopUnits)
	case "bits":
		if bo.unit == 1 {
			return bo
		}
		return &bufferObject{bbr: bo.bbr, unit: 1}
	case "bytes":
		if bo.unit == 8 {
			return bo
		}
		return &bufferObject{bbr: bo.bbr, unit: 8}
	}
	return nil
}
func (bo *bufferObject) JQValueEach() interface{} {
	return nil
}
func (bo *bufferObject) JQValueType() string {
	return "buffer"
}
func (bo *bufferObject) JQValueKeys() interface{} {
	return funcTypeError{name: "keys", typ: "buffer"}
}
func (bo *bufferObject) JQValueHasKey(key interface{}) interface{} {
	return hasKeyTypeError{l: "buffer", r: fmt.Sprintf("%v", key)}
}
func (bo *bufferObject) JQValueToNumber() interface{} {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bo.bbr.bb); err != nil {
		return err
	}
	extraBits := uint((8 - bo.bbr.r.Len%8) % 8)
	return new(big.Int).Rsh(new(big.Int).SetBytes(buf.Bytes()), extraBits)
}
func (bo *bufferObject) JQValueToString() interface{} {
	return bo.JQValue()
}

func (bo *bufferObject) JQValue() interface{} {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bo.bbr.bb.Copy()); err != nil {
		return err
	}
	return buf.String()
}

func (bo *bufferObject) Display(w io.Writer, opts Options) error {
	if opts.Raw {
		if _, err := io.Copy(w, bo.bbr.bb.Copy()); err != nil {
			return err
		}
		return nil
	}

	bbr := bo.bbr
	if bbr.r.Len/8 > int64(opts.DisplayBytes) {
		bbr.r.Len = int64(opts.DisplayBytes) * 8
		bb, err := bbr.bb.BitBufRange(bbr.r.Start, bbr.r.Len)
		if err != nil {
			return err
		}
		bbr.bb = bb
	}

	if err := hexdumpRange(bbr, w, opts); err != nil {
		return err
	}

	unitNames := map[int]string{
		1: "bits",
		8: "bytes",
	}
	unitName := unitNames[bo.unit]
	if unitName == "" {
		unitName = "units"
	}
	_, err := fmt.Fprintf(w, "<%d %s>\n", bo.bbr.r.Len/int64(bo.unit), unitName)
	return err
}

func (bo *bufferObject) ToBuffer() (*bitio.Buffer, error) {
	return bo.bbr.bb, nil
}
