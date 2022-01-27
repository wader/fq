// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This is a fork of std/strconv ParseInt/ParseUint with a strict
// option to allow trailing invalid characters.

package strconvx

import (
	"errors"
	"strconv"
)

// lower(c) is a lower-case letter if and only if
// c is either that lower-case letter or the equivalent upper-case letter.
// Instead of writing c == 'x' || c == 'X' one can write lower(c) == 'x'.
// Note that lower of non-letters can produce other non-letters.
func lower(c byte) byte {
	return c | ('x' - 'X')
}

// ErrRange indicates that a value is out of range for the target type.
var ErrRange = errors.New("value out of range")

// ErrSyntax indicates that a value does not have the right syntax for the target type.
var ErrSyntax = errors.New("invalid syntax")

// A NumError records a failed conversion.
type NumError struct {
	Func string // the failing function (ParseBool, ParseInt, ParseUint, ParseFloat, ParseComplex)
	Num  string // the input
	Err  error  // the reason the conversion failed (e.g. ErrRange, ErrSyntax, etc.)
}

func (e *NumError) Error() string {
	return "strconv." + e.Func + ": " + "parsing " + strconv.Quote(e.Num) + ": " + e.Err.Error()
}

func (e *NumError) Unwrap() error { return e.Err }

// cloneString returns a string copy of x.
//
// All ParseXXX functions allow the input string to escape to the error value.
// This hurts strconv.ParseXXX(string(b)) calls where b is []byte since
// the conversion from []byte must allocate a string on the heap.
// If we assume errors are infrequent, then we can avoid escaping the input
// back to the output by copying it first. This allows the compiler to call
// strconv.ParseXXX without a heap allocation for most []byte to string
// conversions, since it can now prove that the string cannot escape Parse.
//
// TODO: Use strings.Clone instead? However, we cannot depend on "strings"
// since it incurs a transitive dependency on "unicode".
// Either move strings.Clone to an internal/bytealg or make the
// "strings" to "unicode" dependency lighter (see https://go.dev/issue/54098).
func cloneString(x string) string { return string([]byte(x)) }

func syntaxError(fn, str string) *NumError {
	return &NumError{fn, cloneString(str), ErrSyntax}
}

func rangeError(fn, str string) *NumError {
	return &NumError{fn, cloneString(str), ErrRange}
}

func baseError(fn, str string, base int) *NumError {
	return &NumError{fn, cloneString(str), errors.New("invalid base " + strconv.Itoa(base))}
}

func bitSizeError(fn, str string, bitSize int) *NumError {
	return &NumError{fn, cloneString(str), errors.New("invalid bit size " + strconv.Itoa(bitSize))}
}

const intSize = 32 << (^uint(0) >> 63)

// IntSize is the size in bits of an int or uint value.
const IntSize = intSize

const maxUint64 = 1<<64 - 1

// ParseUint is like ParseInt but for unsigned numbers.
//
// A sign prefix is not permitted.
func ParseUint(s string, base int, bitSize int, strict bool) (uint64, int, error) {
	const fnParseUint = "ParseUint"

	if s == "" {
		return 0, 0, syntaxError(fnParseUint, s)
	}

	base0 := base == 0

	s0 := s
	baseCut := 0
	switch {
	case 2 <= base && base <= 36:
		// valid base; strip optional prefix
		if s[0] == '0' && len(s) >= 2 {
			switch {
			case base == 2 && lower(s[1]) == 'b':
				baseCut = 2
			case base == 8 && lower(s[1]) == 'o':
				baseCut = 2
			case base == 16 && lower(s[1]) == 'x':
				baseCut = 2
			}
		}

	case base == 0:
		// Look for octal, hex prefix.
		base = 10
		if s[0] == '0' {
			switch {
			case len(s) >= 3 && lower(s[1]) == 'b':
				base = 2
				baseCut = 2
			case len(s) >= 3 && lower(s[1]) == 'o':
				base = 8
				baseCut = 2
			case len(s) >= 3 && lower(s[1]) == 'x':
				base = 16
				baseCut = 2
			default:
				base = 8
				baseCut = 1
			}
		}

	default:
		return 0, 0, baseError(fnParseUint, s0, base)
	}

	s = s[baseCut:]

	if bitSize == 0 {
		bitSize = IntSize
	} else if bitSize < 0 || bitSize > 64 {
		return 0, 0, bitSizeError(fnParseUint, s0, bitSize)
	}

	// Cutoff is the smallest number such that cutoff*base > maxUint64.
	// Use compile-time constants for common cases.
	var cutoff uint64
	switch base {
	case 10:
		cutoff = maxUint64/10 + 1
	case 16:
		cutoff = maxUint64/16 + 1
	default:
		cutoff = maxUint64/uint64(base) + 1
	}

	maxVal := uint64(1)<<uint(bitSize) - 1

	underscores := false
	var n uint64
	l := baseCut
out:
	for i, c := range []byte(s) {
		var d byte
		switch {
		case c == '_' && base0:
			underscores = true
			continue
		case '0' <= c && c <= '9':
			d = c - '0'
		case 'a' <= lower(c) && lower(c) <= 'z':
			d = lower(c) - 'a' + 10
		default:
			if i > 0 && !strict {
				s0 = s[0 : baseCut+i]
				break out
			}
			return 0, l, syntaxError(fnParseUint, s0)
		}

		if d >= byte(base) {
			if i > 0 && !strict {
				s0 = s[0 : baseCut+i]
				break out
			}
			return 0, l, syntaxError(fnParseUint, s0)
		}

		if n >= cutoff {
			// n*base overflows
			return maxVal, l, rangeError(fnParseUint, s0)
		}
		n *= uint64(base)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxVal {
			// n+d overflows
			return maxVal, l, rangeError(fnParseUint, s0)
		}
		n = n1
		l++
	}

	if underscores && !underscoreOK(s0) {
		return 0, l, syntaxError(fnParseUint, s0)
	}

	return n, l, nil
}

// ParseInt interprets a string s in the given base (0, 2 to 36) and
// bit size (0 to 64) and returns the corresponding value i.
//
// The string may begin with a leading sign: "+" or "-".
//
// If the base argument is 0, the true base is implied by the string's
// prefix following the sign (if present): 2 for "0b", 8 for "0" or "0o",
// 16 for "0x", and 10 otherwise. Also, for argument base 0 only,
// underscore characters are permitted as defined by the Go syntax for
// [integer literals].
//
// The bitSize argument specifies the integer type
// that the result must fit into. Bit sizes 0, 8, 16, 32, and 64
// correspond to int, int8, int16, int32, and int64.
// If bitSize is below 0 or above 64, an error is returned.
//
// The errors that ParseInt returns have concrete type *NumError
// and include err.Num = s. If s is empty or contains invalid
// digits, err.Err = ErrSyntax and the returned value is 0;
// if the value corresponding to s cannot be represented by a
// signed integer of the given size, err.Err = ErrRange and the
// returned value is the maximum magnitude integer of the
// appropriate bitSize and sign.
//
// [integer literals]: https://go.dev/ref/spec#Integer_literals
func ParseInt(s string, base int, bitSize int, strict bool) (i int64, l int, err error) {
	const fnParseInt = "ParseInt"

	if s == "" {
		return 0, 0, syntaxError(fnParseInt, s)
	}

	// Pick off leading sign.
	s0 := s
	neg := false
	negCut := 0
	if s[0] == '+' {
		negCut = 1
	} else if s[0] == '-' {
		neg = true
		negCut = 1
	}
	s = s[negCut:]

	// Convert unsigned and check range.
	var un uint64
	un, l, err = ParseUint(s, base, bitSize, strict)
	l += negCut
	if err != nil && err.(*NumError).Err != ErrRange {
		err.(*NumError).Func = fnParseInt
		err.(*NumError).Num = cloneString(s0)
		return 0, l, err
	}

	if bitSize == 0 {
		bitSize = IntSize
	}

	cutoff := uint64(1 << uint(bitSize-1))

	if !neg && un >= cutoff {
		return int64(cutoff - 1), l, rangeError(fnParseInt, s0)
	}
	if neg && un > cutoff {
		return -int64(cutoff), l, rangeError(fnParseInt, s0)
	}
	n := int64(un)
	if neg {
		n = -n
	}
	return n, l, nil
}

// underscoreOK reports whether the underscores in s are allowed.
// Checking them in this one function lets all the parsers skip over them simply.
// Underscore must appear only between digits or between a base prefix and a digit.
func underscoreOK(s string) bool {
	// saw tracks the last character (class) we saw:
	// ^ for beginning of number,
	// 0 for a digit or base prefix,
	// _ for an underscore,
	// ! for none of the above.
	saw := '^'
	i := 0

	// Optional sign.
	if len(s) >= 1 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}

	// Optional base prefix.
	hex := false
	if len(s) >= 2 && s[0] == '0' && (lower(s[1]) == 'b' || lower(s[1]) == 'o' || lower(s[1]) == 'x') {
		i = 2
		saw = '0' // base prefix counts as a digit for "underscore as digit separator"
		hex = lower(s[1]) == 'x'
	}

	// Number proper.
	for ; i < len(s); i++ {
		// Digits are always okay.
		if '0' <= s[i] && s[i] <= '9' || hex && 'a' <= lower(s[i]) && lower(s[i]) <= 'f' {
			saw = '0'
			continue
		}
		// Underscore must follow digit.
		if s[i] == '_' {
			if saw != '0' {
				return false
			}
			saw = '_'
			continue
		}
		// Underscore must also be followed by digit.
		if saw == '_' {
			return false
		}
		// Saw non-digit, non-underscore.
		saw = '!'
	}
	return saw != '_'
}
