package interp

import (
	"bytes"
	"fmt"
	"fq/pkg/bitio"
	"fq/pkg/ranges"
	"io"
	"math/big"
)

var _ InterpObject = (*bitBufObject)(nil)
var _ ToBitBuf = (*bitBufObject)(nil)

type bitBufObject struct {
	bb   *bitio.Buffer
	unit int
	r    ranges.Range
}

func (*bitBufObject) DisplayName() string { return "buffer" }
func (*bitBufObject) SpecialPropNames() []string {
	return []string{
		"size",
		"start",
		"stop",
		"bits",
		"bytes",
	}
}

func (bo *bitBufObject) JsonLength() interface{} {
	return int(bo.bb.Len()) / bo.unit
}
func (bo *bitBufObject) JsonIndex(index int) interface{} {
	pos, err := bo.bb.Pos()
	if err != nil {
		return err
	}
	if _, err := bo.bb.SeekAbs(int64(index) * int64(bo.unit)); err != nil {
		return err
	}
	v, err := bo.bb.U(bo.unit)
	if err != nil {
		return err
	}
	if _, err := bo.bb.SeekAbs(pos); err != nil {
		return err
	}
	return int(v)
}
func (bo *bitBufObject) JsonRange(start int, end int) interface{} {
	rstart := int64(start * bo.unit)
	rlen := int64((end - start) * bo.unit)
	rbb, err := bo.bb.BitBufRange(rstart, rlen)
	if err != nil {
		return err
	}

	return &bitBufObject{
		bb:   rbb,
		unit: bo.unit,
		r:    ranges.Range{Start: bo.r.Start + rstart, Len: rlen},
	}
}
func (bo *bitBufObject) JsonProperty(name string) interface{} {
	switch name {
	case "size":
		return new(big.Int).SetInt64(bo.bb.Len() / int64(bo.unit))
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
		return &bitBufObject{bb: bo.bb, unit: 1, r: bo.r}
	case "bytes":
		if bo.unit == 8 {
			return bo
		}
		return &bitBufObject{bb: bo.bb, unit: 8, r: bo.r}
	}
	return nil
}
func (bo *bitBufObject) JsonEach() interface{} {
	return nil
}
func (bo *bitBufObject) JsonType() string {
	return "buffer"
}
func (bo *bitBufObject) JsonPrimitiveValue() interface{} {
	buf := &bytes.Buffer{}
	if _, err := io.Copy(buf, bo.bb.Copy()); err != nil {
		return err
	}
	return buf.String()
}

func (bo *bitBufObject) Display(w io.Writer, opts DisplayOptions) error {
	if opts.Raw {
		if _, err := io.Copy(w, bo.bb.Copy()); err != nil {
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
	_, err := fmt.Fprintf(w, "<%d %s>\n", bo.bb.Len()/int64(bo.unit), unitName)
	return err
}

func (bo *bitBufObject) ToBitBuf() (*bitio.Buffer, ranges.Range) {
	return bo.bb.Copy(), ranges.Range{Start: 0, Len: bo.bb.Len()}
}
