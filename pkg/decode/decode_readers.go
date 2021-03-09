package decode

//go:generate sh -c "cat decode_readers_gen.go.tmpl | go run ../../_dev/tmpl.go decode_readers_gen.go.json | gofmt > decode_readers_gen.go"

import (
	"fq/pkg/bitio"
	"strconv"
)

func (d *D) ZeroPadding(nBits int) bool {
	isZero := true
	left := nBits
	for {
		// TODO: smart skip?
		rBits := left
		if rBits == 0 {
			break
		}
		if rBits > 64 {
			rBits = 64
		}
		n, err := d.bitBuf.Bits(int(rBits))
		if err != nil {
			panic(ReadError{Err: err, Op: "ZeroPadding", Size: int64(rBits), Pos: d.Pos()})
		}
		isZero = isZero && n == 0
		left -= rBits
	}
	return isZero
}

func (d *D) FieldOptionalFillFn(name string, fn func(d *D)) int64 {
	start := d.Pos()
	fn(d)
	fillLen := d.Pos() - start
	if fillLen > 0 {
		d.FieldBitBufRange(name, start, fillLen)
	}

	return fillLen
}

func (d *D) FieldOptionalZeroBytes(name string) int64 {
	return d.FieldOptionalFillFn(name, func(d *D) {
		for d.BitsLeft() >= 8 && d.PeekBits(8) == 0 {
			d.SeekRel(8)
		}
	})
}

func (d *D) FieldValidateZeroPadding(name string, nBits int) {
	pos := d.Pos()
	var isZero bool
	d.FieldFn(name, func() *Value {
		isZero = d.ZeroPadding(nBits)
		s := "Correct"
		if !isZero {
			s = "Incorrect"
		}
		return &Value{Symbol: s, Description: "zero padding"}
	})
	if !isZero {
		panic(ValidateError{Reason: "expected zero padding", Pos: pos})
	}
}

func (d *D) Bool() bool {
	b, err := d.bitBuf.Bool()
	if err != nil {
		panic(ReadError{Err: err, Op: "Bool", Size: 1, Pos: d.Pos()})
	}
	return b
}

func (d *D) FieldBool(name string) bool {
	return d.FieldBoolFn(name, func() (bool, string) {
		b, err := d.bitBuf.Bool()
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldBool", Size: 1, Pos: d.Pos()})
		}
		return b, ""
	})
}

func (d *D) UE(nBits int, endian Endian) uint64 {
	n, err := d.bitBuf.UE(nBits, bitio.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "UE", Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) FieldUE(name string, nBits int, endian Endian) uint64 {
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.bitBuf.UE(nBits, bitio.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) SE(nBits int, endian Endian) int64 {
	n, err := d.bitBuf.SE(nBits, bitio.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "SE", Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) FieldSE(name string, nBits int, endian Endian) int64 {
	return d.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := d.bitBuf.SE(nBits, bitio.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldSE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FE(nBits int, endian Endian) float64 {
	n, err := d.bitBuf.FE(nBits, bitio.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "FE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) FieldFE(name string, nBits int, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.bitBuf.FE(nBits, bitio.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldFE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, ""
	})
}

func (d *D) FPE(nBits int, fBits int64, endian Endian) float64 {
	n, err := d.bitBuf.FPE(nBits, fBits, bitio.Endian(endian))
	if err != nil {
		panic(ReadError{Err: err, Op: "FPE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) FieldFPE(name string, nBits int, fBits int64, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.bitBuf.FPE(nBits, fBits, bitio.Endian(endian))
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldFPE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, ""
	})
}

func (d *D) Unary(s uint64) uint64 {
	n, err := d.bitBuf.Unary(s)
	if err != nil {
		panic(ReadError{Err: err, Op: "Unary", Size: 1, Pos: d.Pos()})
	}
	return n
}

func (d *D) FieldBytesLen(name string, nBytes int) []byte {
	return d.FieldBytesFn(name, d.Pos(), int64(nBytes)*8, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldBytesLen", Size: int64(nBytes) * 8, Pos: d.Pos()})
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
		panic(ReadError{Err: err, Op: "UTF8", Size: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return string(s)
}

// FieldUTF8 read nBytes utf8 string and add a field
func (d *D) FieldUTF8(name string, nBytes int) string {
	return d.FieldStrFn(name, func() (string, string) {
		str, err := d.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUTF8", Size: int64(nBytes) * 8, Pos: d.Pos()})
		}
		return str, ""
	})
}
