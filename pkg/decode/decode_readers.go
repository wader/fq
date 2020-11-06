package decode

//go:generate sh -c "cat decode_readers_gen.go.tmpl | go run ../../_dev/tmpl.go | gofmt > decode_readers_gen.go"

import (
	"fq/pkg/bitbuf"
	"strconv"
)

func (d *D) Bool() bool {
	b, err := d.bitBuf.Bool()
	if err != nil {
		panic(BitBufError{Err: err, Op: "Bool", Size: 1, Pos: d.bitBuf.Pos})
	}
	return b
}

func (d *D) FieldBool(name string) bool {
	return d.FieldBoolFn(name, func() (bool, string) {
		b, err := d.bitBuf.Bool()
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBool", Size: 1, Pos: d.bitBuf.Pos})
		}
		return b, ""
	})
}

func (d *D) UE(nBits int64, endian Endian) uint64 {
	n, err := d.bitBuf.UE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(BitBufError{Err: err, Op: "UE", Size: nBits, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) FieldUE(name string, nBits int64, endian Endian) uint64 {
	return d.FieldUFn(name, func() (uint64, DisplayFormat, string) {
		n, err := d.bitBuf.UE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldU" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) SE(nBits int64, endian Endian) int64 {
	n, err := d.bitBuf.SE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(BitBufError{Err: err, Op: "SE", Size: nBits, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) FieldSE(name string, nBits int64, endian Endian) int64 {
	return d.FieldSFn(name, func() (int64, DisplayFormat, string) {
		n, err := d.bitBuf.SE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldS" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
		}
		return n, NumberDecimal, ""
	})
}

func (d *D) FE(nBits int64, endian Endian) float64 {
	n, err := d.bitBuf.FE(nBits, bitbuf.Endian(endian))
	if err != nil {
		panic(BitBufError{Err: err, Op: "FE" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) FieldFE(name string, nBits int64, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.bitBuf.FE(nBits, bitbuf.Endian(endian))
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldFE" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
		}
		return n, ""
	})
}

func (d *D) FPE(nBits int64, dBits int64, endian Endian) float64 {
	n, err := d.bitBuf.FPE(nBits, dBits, bitbuf.Endian(endian))
	if err != nil {
		panic(BitBufError{Err: err, Op: "FPE" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) FieldFPE(name string, nBits int64, dBits int64, endian Endian) float64 {
	return d.FieldFloatFn(name, func() (float64, string) {
		n, err := d.bitBuf.FPE(nBits, dBits, bitbuf.Endian(endian))
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldFPE" + (strconv.Itoa(int(nBits))), Size: nBits, Pos: d.bitBuf.Pos})
		}
		return n, ""
	})
}

func (d *D) Unary(s uint64) uint64 {
	n, err := d.bitBuf.Unary(s)
	if err != nil {
		panic(BitBufError{Err: err, Op: "Unary", Size: 1, Pos: d.bitBuf.Pos})
	}
	return n
}

func (d *D) FieldBytesLen(name string, nBytes int64) []byte {
	return d.FieldBytesFn(name, d.bitBuf.Pos, nBytes*8, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesLen(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesLen", Size: nBytes * 8, Pos: d.bitBuf.Pos})
		}
		return bs, ""
	})
}

func (d *D) FieldBytesRange(name string, firstBit int64, nBytes int64) []byte {
	return d.FieldBytesFn(name, firstBit, nBytes*8, func() ([]byte, string) {
		bs, err := d.bitBuf.BytesRange(firstBit, nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldBytesRange", Size: nBytes * 8, Pos: firstBit})
		}
		return bs, ""
	})
}

func (d *D) FieldUTF8(name string, nBytes int64) string {
	return d.FieldStrFn(name, func() (string, string) {
		str, err := d.bitBuf.UTF8(nBytes)
		if err != nil {
			panic(BitBufError{Err: err, Op: "FieldUTF8", Size: nBytes * 8, Pos: d.bitBuf.Pos})
		}
		return str, ""
	})
}
