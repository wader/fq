package decode

import (
	"fq/pkg/bitio"
	"io"
)

// TODO: FP64,unsigned/BE/LE? rename SFP32?

func (d *D) TryUTF8(nBytes int) (string, error) {
	s, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

// TryUTF8ShortString read pascal short string, max nBytes
func (d *D) TryUTF8ShortString(nBytes int) (string, error) {
	l, err := d.TryU8()
	if err != nil {
		return "", err
	}

	n := int(l)
	if nBytes != -1 {
		n = nBytes - 1
	}

	s, err := d.bitBuf.BytesLen(int(n))
	if err != nil {
		return "", err
	}

	return string(s[0:l]), nil
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

// FieldUTF8ShortString read nBytes utf8 pascal short string and add a field
func (d *D) FieldUTF8ShortString(name string, nBytes int) string {
	return d.FieldStrFn(name, func() (string, string) {
		str, err := d.TryUTF8ShortString(nBytes)
		if err != nil {
			panic(ReadError{Err: err, Name: name, Op: "FieldUTF8ShortString", Size: int64(nBytes) * 8, Pos: d.Pos()})
		}
		return str, ""
	})
}
