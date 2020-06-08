package decode

import (
	"fmt"
	"fq/pkg/bitio"
	"math"
	"strconv"
)

//go:generate sh -c "cat numbers.go.tmpl | go run ../../dev/tmpl.go numbers.go.json | gofmt > numbers_gen.go"

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
		panic(IOError{Err: err, Op: "UE", Size: int64(nBits), Pos: d.Pos()})
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
		panic(IOError{Err: err, Op: "SE", Size: int64(nBits), Pos: d.Pos()})
	}
	return n
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
		panic(IOError{Err: err, Op: "Unary", Size: 1, Pos: d.Pos()})
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
		panic(IOError{Err: err, Op: "FPE" + (strconv.Itoa(nBits)), Size: int64(nBits), Pos: d.Pos()})
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
		return math.Float64frombits(n), nil
	default:
		return 0, fmt.Errorf("unsupported float size %d", nBits)
	}
}

func (d *D) FE(nBits int, endian Endian) float64 {
	n, err := d.TryFE(nBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FE" + (strconv.Itoa(nBits)), Size: int64(nBits), Pos: d.Pos()})
	}
	return n
}

func (d *D) TryFieldUFn(name string, fn func() (uint64, DisplayFormat, string)) (uint64, error) {
	if v, err := d.TryFieldFn(name, func() (*Value, error) {
		u, fmt, d := fn()
		return &Value{V: u, DisplayFormat: fmt, Symbol: d}, nil
	}); err != nil {
		return 0, err
	} else {
		return v.V.(uint64), err
	}
}

func (d *D) FieldUFn(name string, fn func() (uint64, DisplayFormat, string)) uint64 {
	return d.FieldFn(name, func() *Value {
		u, fmt, d := fn()
		return &Value{V: u, DisplayFormat: fmt, Symbol: d}
	}).V.(uint64)
}

func (d *D) FieldUDescFn(name string, fn func() (uint64, DisplayFormat, string, string)) uint64 {
	return d.FieldFn(name, func() *Value {
		u, fmt, s, d := fn()
		return &Value{V: u, DisplayFormat: fmt, Symbol: s, Description: d}
	}).V.(uint64)
}

func (d *D) FieldSFn(name string, fn func() (int64, DisplayFormat, string)) int64 {
	return d.FieldFn(name, func() *Value {
		s, fmt, d := fn()
		return &Value{V: s, DisplayFormat: fmt, Symbol: d}
	}).V.(int64)
}

func (d *D) FieldFloatFn(name string, fn func() (float64, string)) float64 {
	return d.FieldFn(name, func() *Value {
		f, d := fn()
		return &Value{V: f, Symbol: d}
	}).V.(float64)
}

func (d *D) TryFieldUE(name string, nBits int, endian Endian) (uint64, error) {
	return d.TryFieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.TryUE(nBits, endian)
		if err != nil {
			panic(IOError{Err: err, Name: name, Op: "FieldUE" + (strconv.Itoa(nBits)), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FieldUE(name string, nBits int, endian Endian) uint64 {
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.TryUE(nBits, endian)
		if err != nil {
			panic(IOError{Err: err, Name: name, Op: "FieldUE" + (strconv.Itoa(nBits)), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FieldSE(name string, nBits int, endian Endian) int64 {
	return d.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := d.TrySE(nBits, endian)
		if err != nil {
			panic(IOError{Err: err, Name: name, Op: "FieldSE" + (strconv.Itoa(nBits)), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FieldFE(name string, nBits int, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.TryFE(nBits, endian)
		if err != nil {
			panic(IOError{Err: err, Name: name, Op: "FieldFE" + (strconv.Itoa(nBits)), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, ""
	})
}

func (d *D) FieldFPE(name string, nBits int, fBits int64, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.TryFPE(nBits, fBits, endian)
		if err != nil {
			panic(IOError{Err: err, Name: name, Op: "FieldFPE" + (strconv.Itoa(nBits)), Size: int64(nBits), Pos: d.Pos()})
		}
		return n, ""
	})
}
