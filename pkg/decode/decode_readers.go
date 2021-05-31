package decode

//go:generate sh -c "cat decode_readers_gen.go.tmpl | go run ../../_dev/tmpl.go decode_readers_gen.go.json | gofmt > decode_readers_gen.go"

import (
	"fmt"
	"fq/pkg/bitio"
	"io"
	"math"
	"strconv"
)

// TODO: FP64,unsigned/BE/LE? rename SFP32?

func (d *D) TryUTF8(nBytes int) (string, error) {
	s, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func (d *D) TryUnary(s uint64) (uint64, error) {
	var n uint64
	for {
		b, err := d.TryU1()
		if err != nil {
			return 0, err
		}
		if b != s {
			break
		}
		n++
	}
	return n, nil
}

func (d *D) Unary(s uint64) uint64 {
	n, err := d.TryUnary(s)
	if err != nil {
		panic(ReadError{Err: err, Op: "Unary", Size: 1, Pos: d.Pos()})
	}
	return n
}

func (d *D) TryFPE(nBits int, fBits int64, endian Endian) (float64, error) {
	n, err := d.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}
	return float64(n) / float64(uint64(1<<fBits)), nil
}

func (d *D) FPE(nBits int, fBits int64, endian Endian) float64 {
	n, err := d.TryFPE(nBits, fBits, endian)
	if err != nil {
		panic(ReadError{Err: err, Op: "FPE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) TryFE(nBits int, endian Endian) (float64, error) {
	n, err := d.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}
	switch nBits {
	case 32:
		return float64(math.Float32frombits(uint32(n))), nil
	case 64:
		return float64(math.Float64frombits(uint64(n))), nil
	default:
		return 0, fmt.Errorf("unsupported float size %d", nBits)
	}
}

func (d *D) FE(nBits int, endian Endian) float64 {
	n, err := d.TryFE(nBits, endian)
	if err != nil {
		panic(ReadError{Err: err, Op: "FE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

// PeekBits peek nBits bits from buffer
// TODO: share code?
func (d *D) TryPeekBits(nBits int) (uint64, error) {
	start, err := d.bitBuf.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	n, err := d.bits(nBits)
	if _, err := d.bitBuf.SeekBits(start, io.SeekStart); err != nil {
		return 0, err
	}
	return n, err
}

// Bits reads nBits bits from buffer
func (d *D) bits(nBits int) (uint64, error) {
	// 64 bits max, 9 byte worse case if not byte aligned
	buf := d.bitsBuf
	if buf == nil {
		d.bitsBuf = make([]byte, 9)
		buf = d.bitsBuf
	}

	_, err := bitio.ReadFull(d.bitBuf, buf, nBits)
	if err != nil {
		return 0, err
	}

	return bitio.Read64(buf[:], 0, nBits), nil
}

// Bits reads nBits bits from buffer
func (d *D) Bits(nBits int) (uint64, error) {
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// UE reads a nBits bits unsigned integer with byte order endian
// MSB first
func (d *D) TryUE(nBits int, endian Endian) (uint64, error) {
	n, err := d.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}

	return n, nil
}

func (d *D) UE(nBits int, endian Endian) uint64 {
	n, err := d.TryUE(nBits, endian)
	if err != nil {
		panic(ReadError{Err: err, Op: "UE", Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) TrySE(nBits int, endian Endian) (int64, error) {
	n, err := d.Bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}
	var s int64
	if n&(1<<(nBits-1)) > 0 {
		// two's complement
		s = -int64((^n & ((1 << nBits) - 1)) + 1)
	} else {
		s = int64(n)
	}

	return s, nil
}

func (d *D) SE(nBits int, endian Endian) int64 {
	n, err := d.TrySE(nBits, endian)
	if err != nil {
		panic(ReadError{Err: err, Op: "SE", Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) TryPeekFind(nBits int, seekBits int64, fn func(v uint64) bool, maxLen int64) (int64, error) {
	start, err := d.bitBuf.SeekBits(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	found := false
	var count int64
	for {
		if maxLen >= 0 && count >= maxLen {
			break
		}
		v, err := d.TryU(nBits)
		if err != nil {
			if _, err := d.bitBuf.SeekBits(start, io.SeekStart); err != nil {
				return 0, err
			}
			return 0, err
		}
		if fn(v) {
			found = true
			break
		}
		count += seekBits
		if _, err := d.bitBuf.SeekBits(start+count, io.SeekStart); err != nil {
			return 0, err
		}
	}
	if _, err := d.bitBuf.SeekBits(start, io.SeekStart); err != nil {
		return 0, err
	}

	if !found {
		return -1, nil
	}

	return count, nil
}

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
		n, err := d.Bits(int(rBits))
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

// Bool reads one bit as a boolean
func (d *D) TryBool() (bool, error) {
	n, err := d.TryU(1)
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

func (d *D) Bool() bool {
	b, err := d.TryBool()
	if err != nil {
		panic(ReadError{Err: err, Op: "Bool", Size: 1, Pos: d.Pos()})
	}
	return b
}

func (d *D) FieldBool(name string) bool {
	return d.FieldBoolFn(name, func() (bool, string) {
		b, err := d.TryBool()
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldBool", Size: 1, Pos: d.Pos()})
		}
		return b, ""
	})
}

func (d *D) TryFieldUE(name string, nBits int, endian Endian) (uint64, error) {
	return d.TryFieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.TryUE(nBits, endian)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FieldUE(name string, nBits int, endian Endian) uint64 {
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.TryUE(nBits, endian)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FieldSE(name string, nBits int, endian Endian) int64 {
	return d.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := d.TrySE(nBits, endian)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldSE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FieldFE(name string, nBits int, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.TryFE(nBits, endian)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldFE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, ""
	})
}

func (d *D) FieldFPE(name string, nBits int, fBits int64, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.TryFPE(nBits, fBits, endian)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldFPE" + (strconv.Itoa(int(nBits))), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, ""
	})
}

func (d *D) FieldBytesLen(name string, nBytes int) []byte {
	return d.FieldBytesFn(name, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldBytesLen", Size: int64(nBytes) * 8, Pos: d.Pos()})
		}
		return bs, ""
	})
}

// UTF8 read nBytes utf8 string
func (d *D) UTF8(nBytes int) string {
	s, err := d.TryUTF8(nBytes)
	if err != nil {
		panic(ReadError{Err: err, Op: "UTF8", Size: int64(nBytes) * 8, Pos: d.Pos()})
	}
	return string(s)
}

// FieldUTF8 read nBytes utf8 string and add a field
func (d *D) FieldUTF8(name string, nBytes int) string {
	return d.FieldStrFn(name, func() (string, string) {
		str, err := d.TryUTF8(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUTF8", Size: int64(nBytes) * 8, Pos: d.Pos()})
		}
		return str, ""
	})
}
