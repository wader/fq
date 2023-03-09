// Float80 type from https://github.com/mewspring/mewmew-l
// modified as bit to read bytes instead of hex string
//
// This is free and unencumbered software released into the public domain.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that recognize copyright laws, the author or authors
// of this software dedicate any and all copyright interest in the
// software to the public domain. We make this dedication for the benefit
// of the public at large and to the detriment of our heirs and
// successors. We intend this dedication to be an overt act of
// relinquishment in perpetuity of all present and future rights to this
// software under copyright law.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// For more information, please refer to <http://unlicense.org>
package mathex

import (
	"fmt"
	"log"
	"math"
	"math/big"
)

// Float80 represents an 80-bit IEEE 754 extended precision floating-point
// value, in x86 extended precision format.
//
// References:
//
//	https://en.wikipedia.org/wiki/Extended_precision#x86_extended_precision_format
type Float80 struct {
	// Sign and exponent.
	//
	//    1 bit:   sign
	//    15 bits: exponent
	se uint16
	// Integer part and fraction.
	//
	//    1 bit:   integer part
	//    63 bits: fraction
	m uint64
}

// Bits returns the IEEE 754 binary representation of f, with the sign and
// exponent in se and the mantissa in m.
func (f Float80) Bits() (se uint16, m uint64) {
	return f.se, f.m
}

// Bytes returns the x86 extended precision binary representation of f as a byte
// slice.
func (f Float80) Bytes() []byte {
	return []byte(f.String())
}

// String returns the IEEE 754 binary representation of f as a string,
// containing 10 bytes in hexadecimal format.
func (f Float80) String() string {
	return fmt.Sprintf("%04X%016X", f.se, f.m)
}

// Float64 returns the float64 representation of f.
func (f Float80) Float64() float64 {
	se := uint64(f.se)
	m := f.m
	// 1 bit: sign
	sign := se >> 15
	// 15 bits: exponent
	exp := se & 0x7FFF
	// Adjust for exponent bias.
	//
	// === [ binary64 ] =========================================================
	//
	// Exponent bias 1023.
	//
	//    +===========================+=======================+
	//    | Exponent (in binary)      | Notes                 |
	//    +===========================+=======================+
	//    | 00000000000               | zero/subnormal number |
	//    +---------------------------+-----------------------+
	//    | 00000000001 - 11111111110 | normalized value      |
	//    +---------------------------+-----------------------+
	//    | 11111111111               | infinity/NaN          |
	//    +---------------------------+-----------------------+
	//
	// References:
	//    https://en.wikipedia.org/wiki/Double-precision_floating-point_format#Exponent_encoding
	exp64 := int64(exp) - 16383 + 1023
	switch {
	case exp == 0:
		// exponent is all zeroes.
		exp64 = 0
	case exp == 0x7FFF:
		// exponent is all ones.
		exp64 = 0x7FF
	default:
	}
	// 63 bits: fraction
	frac := m & 0x7FFFFFFFFFFFFFFF
	// Sign, exponent and fraction of binary64.
	//
	//    1 bit:   sign
	//    11 bits: exponent
	//    52 bits: fraction
	//
	// References:
	//    https://en.wikipedia.org/wiki/Double-precision_floating-point_format#IEEE_754_double-precision_binary_floating-point_format:_binary64
	bits := sign<<63 | uint64(exp64)<<52 | frac>>11
	return math.Float64frombits(bits)
}

// BigFloat returns the *big.Float representation of f.
func (f Float80) BigFloat() *big.Float {
	x := &big.Float{}
	sign := (f.se & 0x8000) != 0
	e := f.se & 0x7FFF
	s := fmt.Sprintf("0x.%Xp%d", f.m, e-16383+1)
	if sign {
		s = "-" + s
	}
	x.SetPrec(52)
	_, _, err := x.Parse(s, 0)
	if err != nil {
		log.Printf("big.Float.Parse: error %v", err)
	}
	return x
}

// NewFloat80FromFloat64 returns the nearest 80-bit floating-point value for x.
func NewFloat80FromFloat64(x float64) Float80 {
	// Sign, exponent and fraction of binary64.
	//
	//    1 bit:   sign
	//    11 bits: exponent
	//    52 bits: fraction
	bits := math.Float64bits(x)
	// 1 bit: sign
	sign := uint16(bits >> 63)
	// 11 bits: exponent
	exp := bits >> 52 & 0x7FF
	// 52 bits: fraction
	frac := bits & 0xFFFFFFFFFFFFF

	if exp == 0 && frac == 0 {
		// zero value.
		return Float80{}
	}

	// Sign, exponent and fraction of binary80.
	//
	//    1 bit:   sign
	//    15 bits: exponent
	//    1 bit:   integer part
	//    63 bits: fraction

	// 15 bits: exponent.
	//
	// Exponent bias 1023  (binary64)
	// Exponent bias 16383 (binary80)
	exp80 := int64(exp) - 1023 + 16383
	// 63 bits: fraction.
	//
	frac80 := frac << 11
	switch {
	case exp == 0:
		exp80 = 0
	case exp == 0x7FF:
		exp80 = 0x7FFF
	}
	se := sign<<15 | uint16(exp80)
	// Integer part set to specify normalized value.
	m := 0x8000000000000000 | frac80
	return NewFloat80FromBits(se, m)
}

// NewFloat80FromBytes returns a new 80-bit floating-point value based on b,
func NewFloat80FromBytes(b []byte) Float80 {
	var f Float80
	if len(b) != 10 {
		panic(fmt.Errorf("invalid length of float80 representation, expected 10, got %d", len(b)))
	}
	f.se = uint16(int64(b[0])<<8 | int64(b[1]<<0))
	f.m = uint64(0 |
		int64(b[2])<<56 |
		int64(b[3])<<48 |
		int64(b[4])<<40 |
		int64(b[5])<<32 |
		int64(b[6])<<24 |
		int64(b[7])<<16 |
		int64(b[8])<<8 |
		int64(b[9])<<0,
	)
	return f
}

// NewFloat80FromBits returns a new 80-bit floating-point value based on the
// sign, exponent and mantissa bits.
func NewFloat80FromBits(se uint16, m uint64) Float80 {
	return Float80{
		se: se,
		m:  m,
	}
}
