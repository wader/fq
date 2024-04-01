// Copyright (C) 2014 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package mathx

import "unsafe"

// Float16 represents a 16-bit floating point number, containing a single sign bit, 5 exponent bits
// and 10 fractional bits. This corresponds to IEEE 754-2008 binary16 (or half precision float) type.
//
//	 MSB                                                                         LSB
//	╔════╦════╤════╤════╤════╤════╦════╤════╤════╤════╤════╤════╤════╤════╤════╤════╗
//	║Sign║ E₄ │ E₃ │ E₂ │ E₁ │ E₀ ║ F₉ │ F₈ │ F₇ │ F₆ │ F₅ │ F₄ │ F₃ │ F₂ │ F₁ │ F₀ ║
//	╚════╩════╧════╧════╧════╧════╩════╧════╧════╧════╧════╧════╧════╧════╧════╧════╝
//	Where E is the exponent bits and F is the fractional bits.
type Float16 uint16

const (
	float16ExpMask  Float16 = 0x7c00
	float16ExpBias  uint32  = 0xf
	float16ExpShift uint32  = 10
	float16FracMask Float16 = 0x03ff
	float16SignMask Float16 = 0x8000
	float32ExpMask  uint32  = 0x7f800000
	float32ExpBias  uint32  = 0x7f
	float32ExpShift uint32  = 23
	float32FracMask uint32  = 0x007fffff
)

// Float32 returns the Float16 value expanded to a float32. Infinities and NaNs are expanded as
// such.
func (f Float16) Float32() float32 {
	u32 := expandF16ToF32(f)
	ptr := unsafe.Pointer(&u32)
	f32 := *(*float32)(ptr)
	return f32
}

// IsNaN reports whether f is an “not-a-number” value.
func (f Float16) IsNaN() bool {
	return (f&float16ExpMask == float16ExpMask) && (f&float16FracMask != 0)
}

// IsInf reports whether f is an infinity, according to sign. If sign > 0, IsInf reports whether
// f is positive infinity. If sign < 0, IsInf reports whether f is negative infinity. If sign ==
// 0, IsInf reports whether f is either infinity.
func (f Float16) IsInf(sign int) bool {
	return ((f == float16ExpMask) && sign >= 0) ||
		(f == (float16SignMask|float16ExpMask) && sign <= 0)
}

// Float16NaN returns an “not-a-number” value.
func NewFloat16NaN() Float16 { return float16ExpMask | float16FracMask }

// Float16Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func NewFloat16Inf(sign int) Float16 {
	if sign >= 0 {
		return float16ExpMask
	} else {
		return float16SignMask | float16ExpMask
	}
}

// NewFloat16 returns a Float16 encoding of a 32-bit floating point number. Infinities and NaNs
// are encoded as such. Very large and very small numbers get rounded to infinity and zero
// respectively.
func NewFloat16(f32 float32) Float16 {
	ptr := unsafe.Pointer(&f32)
	u32 := *(*uint32)(ptr)
	sign := Float16(u32>>16) & float16SignMask
	exp := (u32 & float32ExpMask) >> float32ExpShift
	frac := u32 & 0x7fffff
	if exp == 0xff {
		// NaN or Infinity
		if frac != 0 { // NaN
			frac = 0x3f
		}
		return sign | float16ExpMask | Float16(frac)
	}
	if exp+float16ExpBias <= float32ExpBias {
		// Exponent is too small to represent in a Float16 (or a zero). We need to output
		// denormalized numbers (possibly rounding very small numbers to zero).
		denorm := float32ExpBias - exp - 1
		frac += 1 << float32ExpShift
		frac >>= denorm
		return sign | Float16(frac)
	}
	if exp > float32ExpBias+float16ExpBias {
		// Number too large to represent in a Float16 => round to Infinity.
		return sign | float16ExpMask
	}
	// General case.
	return sign | Float16(((exp+float16ExpBias-float32ExpBias)<<float16ExpShift)|(frac>>13))
}
func expandF16ToF32(in Float16) uint32 {
	sign := uint32(in&float16SignMask) << 16
	frac := uint32(in&float16FracMask) << 13
	exp := uint32(in&float16ExpMask) >> float16ExpShift
	if exp == 0x1f {
		// NaN of Infinity
		return sign | float32ExpMask | frac
	}
	if exp == 0 {
		if frac == 0 {
			// Zero
			return sign
		}
		// Denormalized number. In a float32 it must be stored in a normalized form, so
		// we normalize it.
		exp++
		for frac&float32ExpMask == 0 {
			frac <<= 1
			exp--
		}
		frac &= float32FracMask
	}
	exp += (float32ExpBias - float16ExpBias)
	return sign | (exp << float32ExpShift) | frac
}
