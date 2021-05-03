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

var _ InterpObject = (*bitBufObject)(nil)
var _ ToBuffer = (*bitBufObject)(nil)

type bitBufObject struct {
	bbr  bufferRange
	unit int
}

func newBifBufObject(bb *bitio.Buffer, unit int) *bitBufObject {
	return &bitBufObject{
		bbr:  bufferRange{bb: bb, r: ranges.Range{Start: 0, Len: bb.Len()}},
		unit: unit,
	}
}

func (*bitBufObject) DisplayName() string { return "buffer" }
func (*bitBufObject) ExtValueKeys() []string {
	return []string{
		"size",
		"start",
		"stop",
		"bits",
		"bytes",
	}
}

func (bo *bitBufObject) JQValueLength() interface{} {
	return int(bo.bbr.r.Len / int64(bo.unit))
}
func (bo *bitBufObject) JQValueIndex(index int) interface{} {
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
}
func (bo *bitBufObject) JQValueSlice(start int, end int) interface{} {
	rStart := int64(start * bo.unit)
	rLen := int64((end - start) * bo.unit)
	rbb, err := bo.bbr.bb.BitBufRange(rStart, rLen)
	if err != nil {
		return err
	}

	return &bitBufObject{
		bbr:  bufferRange{bb: rbb, r: ranges.Range{Start: bo.bbr.r.Start + rStart, Len: rLen}},
		unit: bo.unit,
	}
}
func (bo *bitBufObject) JQValueProperty(name string) interface{} {
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
		return &bitBufObject{bbr: bo.bbr, unit: 1}
	case "bytes":
		if bo.unit == 8 {
			return bo
		}
		return &bitBufObject{bbr: bo.bbr, unit: 8}
	}
	return nil
}
func (bo *bitBufObject) JQValueEach() interface{} {
	return nil
}
func (bo *bitBufObject) JQValueType() string {
	return "buffer"
}
func (bo *bitBufObject) JQValueKeys() interface{} {
	return fmt.Errorf("can't get keys from bitbuf")
}
func (bo *bitBufObject) JQValueHasKey(key interface{}) interface{} {
	return fmt.Errorf("can't get keys from bitbuf")
}

func (bo *bitBufObject) JQValue() interface{} {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bo.bbr.bb.Copy()); err != nil {
		return err
	}
	return buf.String()
}

func (bo *bitBufObject) Display(w io.Writer, opts Options) error {
	if opts.Raw {
		if _, err := io.Copy(w, bo.bbr.bb.Copy()); err != nil {
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
	_, err := fmt.Fprintf(w, "<%d %s>\n", bo.bbr.r.Len/int64(bo.unit), unitName)
	return err
}

func (bo *bitBufObject) ToBuffer() (*bitio.Buffer, error) {
	return bo.bbr.bb, nil
}
