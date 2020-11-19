package decode

//go:generate sh -c "cat decode_readers_gen.go.tmpl | go run ../../_dev/tmpl.go | gofmt > decode_readers_gen.go"

import (
	"fq/pkg/bitbuf"
	"strconv"
)

func (d *D) ZeroPadding(nBits int) bool {
	isZero := true
	left := nBits
	for {
		// TODO: smart skip?
		rbits := left
		if rbits == 0 {
			break
		}
		if rbits > 64 {
			rbits = 64
		}
		n, err := d.bitBuf.Bits(int(rbits))
		if err != nil {
			panic(ReadError{Err: err, Op: "ZeroPadding", Size: int64(rbits), Pos: d.bitBuf.Pos()})
		}
		isZero = isZero && n == 0
		left -= rbits
	}
	return isZero
}

func (d *D) FieldValidateZeroPadding(name string, nBits int) {
	pos := d.bitBuf.Pos()
	var isZero bool
	d.FieldFn(name, func() *Value {
		isZero = d.ZeroPadding(nBits)
		s := "Correct"
		if !isZero {
			s = "Incorrect"
		}
		return &Value{Symbol: s, Desc: "zero padding"}
	})
	if !isZero {
		panic(ValidateError{Reason: "expected zero padding", Pos: pos})
	}
}

func (d *D) Bool() bool {
	b, err := d.bitBuf.Bool()
	if err != nil {
		panic(ReadError{Err: err, Op: "Bool", Size: 1, Pos: d.bitBuf.Pos()})
	}
	return b
}

func (d *D) FieldBool(name string) bool {
	return d.FieldBoolFn(name, func() (bool, string) {
		b, err := d.bitBuf.Bool()
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldBool", Size: 1, Pos: d.bitBuf.Pos()})
		}
		return b, ""
	})
}

func (d *D) UE(nBits int, endian Endian) uint64 {
	n, err := d.bitBuf.UE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "UE", Size: int64(nBits), Pos: d.bitBuf.Pos()})
	}
	return n
}

func (d *D) FieldUE(name string, nBits int, endian Endian) uint64 {
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.bitBuf.UE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.bitBuf.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) SE(nBits int, endian Endian) int64 {
	n, err := d.bitBuf.SE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "SE", Size: int64(nBits), Pos: d.bitBuf.Pos()})
	}
	return n
}

func (d *D) FieldSE(name string, nBits int, endian Endian) int64 {
	return d.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := d.bitBuf.SE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldSE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.bitBuf.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FE(nBits int, endian Endian) float64 {
	n, err := d.bitBuf.FE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "FE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.bitBuf.Pos()})
	}
	return n
}

func (d *D) FieldFE(name string, nBits int, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.bitBuf.FE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldFE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.bitBuf.Pos()})
		}
		return n, ""
	})
}

func (d *D) FPE(nBits int, dBits int64, endian Endian) float64 {
	n, err := d.bitBuf.FPE(nBits, dBits, bitbuf.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "FPE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.bitBuf.Pos()})
	}
	return n
}

func (d *D) FieldFPE(name string, nBits int, dBits int64, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.bitBuf.FPE(nBits, dBits, bitbuf.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldFPE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.bitBuf.Pos()})
		}
		return n, ""
	})
}

func (d *D) Unary(s uint64) uint64 {
	n, err := d.bitBuf.Unary(s)
	if err != nil {
		panic(ReadError{Err: err, Op: "Unary", Size: 1, Pos: d.bitBuf.Pos()})
	}
	return n
}

func (d *D) FieldBytesLen(name string, nBytes int) []byte {
	return d.FieldBytesFn(name, d.bitBuf.Pos(), int64(nBytes)*8, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldBytesLen", Size: int64(nBytes) * 8, Pos: d.bitBuf.Pos()})
		}
		return bs, ""
	})
}

func (d *D) FieldBytesRange(name string, firstBit int64, nBytes int) []byte {
	return d.FieldBytesFn(name, firstBit, int64(nBytes)*8, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesRange(firstBit, nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldBytesRange", Size: int64(nBytes) * 8, Pos: firstBit})
		}
		return bs, ""
	})
}

// UTF8 read nBytes utf8 string
func (d *D) UTF8(nBytes int) string {
	s, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		panic(ReadError{Err: err, Op: "UTF8", Size: int64(nBytes) * 8, Pos: d.bitBuf.Pos()})
	}
	return string(s)
}

// FieldUTF8 read nBytes utf8 string and add a field
func (d *D) FieldUTF8(name string, nBytes int) string {
	return d.FieldStrFn(name, func() (string, string) {
		str, err := d.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUTF8", Size: int64(nBytes) * 8, Pos: d.bitBuf.Pos()})
		}
		return str, ""
	})
}
