// Code below generated from decode_gen.go.tmpl
package decode

import (
	"fmt"
	"math/big"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/text/encoding"
)

// Type BigInt

// TryFieldBigIntScalarFn tries to add a field, calls scalar functions and returns actual value as a BigInt
func (d *D) TryFieldBigIntScalarFn(name string, fn func(d *D) (scalar.S, error), sms ...scalar.Mapper) (*big.Int, error) {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d) }, sms...)
	if err != nil {
		return nil, err
	}
	return v.ActualBigInt(), err
}

// FieldBigIntScalarFn adds a field, calls scalar functions and returns actual value as a BigInt
func (d *D) FieldBigIntScalarFn(name string, fn func(d *D) scalar.S, sms ...scalar.Mapper) *big.Int {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "BigInt", Pos: d.Pos()})
	}
	return v.ActualBigInt()
}

// FieldBigIntFn adds a field, calls *big.Int decode function and returns actual value as a BigInt
func (d *D) FieldBigIntFn(name string, fn func(d *D) *big.Int, sms ...scalar.Mapper) *big.Int {
	return d.FieldBigIntScalarFn(name, func(d *D) scalar.S { return scalar.S{Actual: fn(d)} }, sms...)
}

// TryFieldBigIntFn tries to add a field, calls *big.Int decode function and returns actual value as a BigInt
func (d *D) TryFieldBigIntFn(name string, fn func(d *D) (*big.Int, error), sms ...scalar.Mapper) (*big.Int, error) {
	return d.TryFieldBigIntScalarFn(name, func(d *D) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// TryFieldScalarBigIntFn tries to add a field, calls *big.Int decode function and returns scalar
func (d *D) TryFieldScalarBigIntFn(name string, fn func(d *D) (*big.Int, error), sms ...scalar.Mapper) (*scalar.S, error) {
	return d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// FieldScalarBigIntFn tries to add a field, calls *big.Int decode function and returns scalar
func (d *D) FieldScalarBigIntFn(name string, fn func(d *D) *big.Int, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarBigIntFn(name, func(d *D) (*big.Int, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "BigInt", Pos: d.Pos()})
	}
	return v
}

// Type BitBuf

// TryFieldBitBufScalarFn tries to add a field, calls scalar functions and returns actual value as a BitBuf
func (d *D) TryFieldBitBufScalarFn(name string, fn func(d *D) (scalar.S, error), sms ...scalar.Mapper) (bitio.ReaderAtSeeker, error) {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d) }, sms...)
	if err != nil {
		return nil, err
	}
	return v.ActualBitBuf(), err
}

// FieldBitBufScalarFn adds a field, calls scalar functions and returns actual value as a BitBuf
func (d *D) FieldBitBufScalarFn(name string, fn func(d *D) scalar.S, sms ...scalar.Mapper) bitio.ReaderAtSeeker {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "BitBuf", Pos: d.Pos()})
	}
	return v.ActualBitBuf()
}

// FieldBitBufFn adds a field, calls bitio.ReaderAtSeeker decode function and returns actual value as a BitBuf
func (d *D) FieldBitBufFn(name string, fn func(d *D) bitio.ReaderAtSeeker, sms ...scalar.Mapper) bitio.ReaderAtSeeker {
	return d.FieldBitBufScalarFn(name, func(d *D) scalar.S { return scalar.S{Actual: fn(d)} }, sms...)
}

// TryFieldBitBufFn tries to add a field, calls bitio.ReaderAtSeeker decode function and returns actual value as a BitBuf
func (d *D) TryFieldBitBufFn(name string, fn func(d *D) (bitio.ReaderAtSeeker, error), sms ...scalar.Mapper) (bitio.ReaderAtSeeker, error) {
	return d.TryFieldBitBufScalarFn(name, func(d *D) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// TryFieldScalarBitBufFn tries to add a field, calls bitio.ReaderAtSeeker decode function and returns scalar
func (d *D) TryFieldScalarBitBufFn(name string, fn func(d *D) (bitio.ReaderAtSeeker, error), sms ...scalar.Mapper) (*scalar.S, error) {
	return d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// FieldScalarBitBufFn tries to add a field, calls bitio.ReaderAtSeeker decode function and returns scalar
func (d *D) FieldScalarBitBufFn(name string, fn func(d *D) bitio.ReaderAtSeeker, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarBitBufFn(name, func(d *D) (bitio.ReaderAtSeeker, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "BitBuf", Pos: d.Pos()})
	}
	return v
}

// Type Bool

// TryFieldBoolScalarFn tries to add a field, calls scalar functions and returns actual value as a Bool
func (d *D) TryFieldBoolScalarFn(name string, fn func(d *D) (scalar.S, error), sms ...scalar.Mapper) (bool, error) {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d) }, sms...)
	if err != nil {
		return false, err
	}
	return v.ActualBool(), err
}

// FieldBoolScalarFn adds a field, calls scalar functions and returns actual value as a Bool
func (d *D) FieldBoolScalarFn(name string, fn func(d *D) scalar.S, sms ...scalar.Mapper) bool {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Bool", Pos: d.Pos()})
	}
	return v.ActualBool()
}

// FieldBoolFn adds a field, calls bool decode function and returns actual value as a Bool
func (d *D) FieldBoolFn(name string, fn func(d *D) bool, sms ...scalar.Mapper) bool {
	return d.FieldBoolScalarFn(name, func(d *D) scalar.S { return scalar.S{Actual: fn(d)} }, sms...)
}

// TryFieldBoolFn tries to add a field, calls bool decode function and returns actual value as a Bool
func (d *D) TryFieldBoolFn(name string, fn func(d *D) (bool, error), sms ...scalar.Mapper) (bool, error) {
	return d.TryFieldBoolScalarFn(name, func(d *D) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// TryFieldScalarBoolFn tries to add a field, calls bool decode function and returns scalar
func (d *D) TryFieldScalarBoolFn(name string, fn func(d *D) (bool, error), sms ...scalar.Mapper) (*scalar.S, error) {
	return d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// FieldScalarBoolFn tries to add a field, calls bool decode function and returns scalar
func (d *D) FieldScalarBoolFn(name string, fn func(d *D) bool, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarBoolFn(name, func(d *D) (bool, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Bool", Pos: d.Pos()})
	}
	return v
}

// Type F

// TryFieldFScalarFn tries to add a field, calls scalar functions and returns actual value as a F
func (d *D) TryFieldFScalarFn(name string, fn func(d *D) (scalar.S, error), sms ...scalar.Mapper) (float64, error) {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d) }, sms...)
	if err != nil {
		return 0, err
	}
	return v.ActualF(), err
}

// FieldFScalarFn adds a field, calls scalar functions and returns actual value as a F
func (d *D) FieldFScalarFn(name string, fn func(d *D) scalar.S, sms ...scalar.Mapper) float64 {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F", Pos: d.Pos()})
	}
	return v.ActualF()
}

// FieldFFn adds a field, calls float64 decode function and returns actual value as a F
func (d *D) FieldFFn(name string, fn func(d *D) float64, sms ...scalar.Mapper) float64 {
	return d.FieldFScalarFn(name, func(d *D) scalar.S { return scalar.S{Actual: fn(d)} }, sms...)
}

// TryFieldFFn tries to add a field, calls float64 decode function and returns actual value as a F
func (d *D) TryFieldFFn(name string, fn func(d *D) (float64, error), sms ...scalar.Mapper) (float64, error) {
	return d.TryFieldFScalarFn(name, func(d *D) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// TryFieldScalarFFn tries to add a field, calls float64 decode function and returns scalar
func (d *D) TryFieldScalarFFn(name string, fn func(d *D) (float64, error), sms ...scalar.Mapper) (*scalar.S, error) {
	return d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// FieldScalarFFn tries to add a field, calls float64 decode function and returns scalar
func (d *D) FieldScalarFFn(name string, fn func(d *D) float64, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarFFn(name, func(d *D) (float64, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F", Pos: d.Pos()})
	}
	return v
}

// Type S

// TryFieldSScalarFn tries to add a field, calls scalar functions and returns actual value as a S
func (d *D) TryFieldSScalarFn(name string, fn func(d *D) (scalar.S, error), sms ...scalar.Mapper) (int64, error) {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d) }, sms...)
	if err != nil {
		return 0, err
	}
	return v.ActualS(), err
}

// FieldSScalarFn adds a field, calls scalar functions and returns actual value as a S
func (d *D) FieldSScalarFn(name string, fn func(d *D) scalar.S, sms ...scalar.Mapper) int64 {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S", Pos: d.Pos()})
	}
	return v.ActualS()
}

// FieldSFn adds a field, calls int64 decode function and returns actual value as a S
func (d *D) FieldSFn(name string, fn func(d *D) int64, sms ...scalar.Mapper) int64 {
	return d.FieldSScalarFn(name, func(d *D) scalar.S { return scalar.S{Actual: fn(d)} }, sms...)
}

// TryFieldSFn tries to add a field, calls int64 decode function and returns actual value as a S
func (d *D) TryFieldSFn(name string, fn func(d *D) (int64, error), sms ...scalar.Mapper) (int64, error) {
	return d.TryFieldSScalarFn(name, func(d *D) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// TryFieldScalarSFn tries to add a field, calls int64 decode function and returns scalar
func (d *D) TryFieldScalarSFn(name string, fn func(d *D) (int64, error), sms ...scalar.Mapper) (*scalar.S, error) {
	return d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// FieldScalarSFn tries to add a field, calls int64 decode function and returns scalar
func (d *D) FieldScalarSFn(name string, fn func(d *D) int64, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarSFn(name, func(d *D) (int64, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S", Pos: d.Pos()})
	}
	return v
}

// Type Str

// TryFieldStrScalarFn tries to add a field, calls scalar functions and returns actual value as a Str
func (d *D) TryFieldStrScalarFn(name string, fn func(d *D) (scalar.S, error), sms ...scalar.Mapper) (string, error) {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d) }, sms...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}

// FieldStrScalarFn adds a field, calls scalar functions and returns actual value as a Str
func (d *D) FieldStrScalarFn(name string, fn func(d *D) scalar.S, sms ...scalar.Mapper) string {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Str", Pos: d.Pos()})
	}
	return v.ActualStr()
}

// FieldStrFn adds a field, calls string decode function and returns actual value as a Str
func (d *D) FieldStrFn(name string, fn func(d *D) string, sms ...scalar.Mapper) string {
	return d.FieldStrScalarFn(name, func(d *D) scalar.S { return scalar.S{Actual: fn(d)} }, sms...)
}

// TryFieldStrFn tries to add a field, calls string decode function and returns actual value as a Str
func (d *D) TryFieldStrFn(name string, fn func(d *D) (string, error), sms ...scalar.Mapper) (string, error) {
	return d.TryFieldStrScalarFn(name, func(d *D) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// TryFieldScalarStrFn tries to add a field, calls string decode function and returns scalar
func (d *D) TryFieldScalarStrFn(name string, fn func(d *D) (string, error), sms ...scalar.Mapper) (*scalar.S, error) {
	return d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// FieldScalarStrFn tries to add a field, calls string decode function and returns scalar
func (d *D) FieldScalarStrFn(name string, fn func(d *D) string, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarStrFn(name, func(d *D) (string, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Str", Pos: d.Pos()})
	}
	return v
}

// Type U

// TryFieldUScalarFn tries to add a field, calls scalar functions and returns actual value as a U
func (d *D) TryFieldUScalarFn(name string, fn func(d *D) (scalar.S, error), sms ...scalar.Mapper) (uint64, error) {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d) }, sms...)
	if err != nil {
		return 0, err
	}
	return v.ActualU(), err
}

// FieldUScalarFn adds a field, calls scalar functions and returns actual value as a U
func (d *D) FieldUScalarFn(name string, fn func(d *D) scalar.S, sms ...scalar.Mapper) uint64 {
	v, err := d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U", Pos: d.Pos()})
	}
	return v.ActualU()
}

// FieldUFn adds a field, calls uint64 decode function and returns actual value as a U
func (d *D) FieldUFn(name string, fn func(d *D) uint64, sms ...scalar.Mapper) uint64 {
	return d.FieldUScalarFn(name, func(d *D) scalar.S { return scalar.S{Actual: fn(d)} }, sms...)
}

// TryFieldUFn tries to add a field, calls uint64 decode function and returns actual value as a U
func (d *D) TryFieldUFn(name string, fn func(d *D) (uint64, error), sms ...scalar.Mapper) (uint64, error) {
	return d.TryFieldUScalarFn(name, func(d *D) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// TryFieldScalarUFn tries to add a field, calls uint64 decode function and returns scalar
func (d *D) TryFieldScalarUFn(name string, fn func(d *D) (uint64, error), sms ...scalar.Mapper) (*scalar.S, error) {
	return d.TryFieldScalarFn(name, func(_ scalar.S) (scalar.S, error) {
		v, err := fn(d)
		return scalar.S{Actual: v}, err
	}, sms...)
}

// FieldScalarUFn tries to add a field, calls uint64 decode function and returns scalar
func (d *D) FieldScalarUFn(name string, fn func(d *D) uint64, sms ...scalar.Mapper) *scalar.S {
	v, err := d.TryFieldScalarUFn(name, func(d *D) (uint64, error) { return fn(d), nil }, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U", Pos: d.Pos()})
	}
	return v
}

// Require/Assert/Validate BigInt

func requireBigInt(name string, s scalar.S, desc bool, fail bool, vs ...*big.Int) (scalar.S, error) {
	a := s.ActualBigInt()
	for _, b := range vs {
		if a.Cmp(b) == 0 {
			if desc {
				s.Description = "valid"
			}
			return s, nil
		}
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s BigInt", name)
	}
	return s, nil
}

// RequireBigInt that actual value is one of given *big.Int values
func (d *D) RequireBigInt(vs ...*big.Int) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireBigInt("require", s, false, true, vs...) })
}

// AssertBigInt validate and asserts that actual value is one of given *big.Int values
func (d *D) AssertBigInt(vs ...*big.Int) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireBigInt("assert", s, true, !d.Options.Force, vs...) })
}

// ValidateBigInt validates that actual value is one of given *big.Int values
func (d *D) ValidateBigInt(vs ...*big.Int) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireBigInt("validate", s, true, false, vs...) })
}

// Require/Assert/ValidatRange BigInt

func requireRangeBigInt(name string, s scalar.S, desc bool, fail bool, start, end *big.Int) (scalar.S, error) {
	a := s.ActualBigInt()
	if a.Cmp(start) >= 0 && a.Cmp(end) <= 0 {
		if desc {
			s.Description = "valid"
		}
		return s, nil
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s BigInt range %v-%v", name, start, end)
	}
	return s, nil
}

// RequireBigIntRange require that actual value is in range
func (d *D) RequireBigIntRange(start, end *big.Int) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeBigInt("require", s, false, true, start, end) })
}

// AssertBigIntRange asserts that actual value is in range
func (d *D) AssertBigIntRange(start, end *big.Int) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return requireRangeBigInt("assert", s, true, !d.Options.Force, start, end)
	})
}

// ValidateBigIntRange validates that actual value is in range
func (d *D) ValidateBigIntRange(start, end *big.Int) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeBigInt("validate", s, true, false, start, end) })
}

// Require/Assert/Validate Bool

func requireBool(name string, s scalar.S, desc bool, fail bool, vs ...bool) (scalar.S, error) {
	a := s.ActualBool()
	for _, b := range vs {
		if a == b {
			if desc {
				s.Description = "valid"
			}
			return s, nil
		}
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s Bool", name)
	}
	return s, nil
}

// RequireBool that actual value is one of given bool values
func (d *D) RequireBool(vs ...bool) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireBool("require", s, false, true, vs...) })
}

// AssertBool validate and asserts that actual value is one of given bool values
func (d *D) AssertBool(vs ...bool) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireBool("assert", s, true, !d.Options.Force, vs...) })
}

// ValidateBool validates that actual value is one of given bool values
func (d *D) ValidateBool(vs ...bool) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireBool("validate", s, true, false, vs...) })
}

// Require/Assert/Validate F

func requireF(name string, s scalar.S, desc bool, fail bool, vs ...float64) (scalar.S, error) {
	a := s.ActualF()
	for _, b := range vs {
		if a == b {
			if desc {
				s.Description = "valid"
			}
			return s, nil
		}
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s F", name)
	}
	return s, nil
}

// RequireF that actual value is one of given float64 values
func (d *D) RequireF(vs ...float64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireF("require", s, false, true, vs...) })
}

// AssertF validate and asserts that actual value is one of given float64 values
func (d *D) AssertF(vs ...float64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireF("assert", s, true, !d.Options.Force, vs...) })
}

// ValidateF validates that actual value is one of given float64 values
func (d *D) ValidateF(vs ...float64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireF("validate", s, true, false, vs...) })
}

// Require/Assert/ValidatRange F

func requireRangeF(name string, s scalar.S, desc bool, fail bool, start, end float64) (scalar.S, error) {
	a := s.ActualF()
	if a >= start && a <= end {
		if desc {
			s.Description = "valid"
		}
		return s, nil
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s F range %v-%v", name, start, end)
	}
	return s, nil
}

// RequireFRange require that actual value is in range
func (d *D) RequireFRange(start, end float64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeF("require", s, false, true, start, end) })
}

// AssertFRange asserts that actual value is in range
func (d *D) AssertFRange(start, end float64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return requireRangeF("assert", s, true, !d.Options.Force, start, end)
	})
}

// ValidateFRange validates that actual value is in range
func (d *D) ValidateFRange(start, end float64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeF("validate", s, true, false, start, end) })
}

// Require/Assert/Validate S

func requireS(name string, s scalar.S, desc bool, fail bool, vs ...int64) (scalar.S, error) {
	a := s.ActualS()
	for _, b := range vs {
		if a == b {
			if desc {
				s.Description = "valid"
			}
			return s, nil
		}
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s S", name)
	}
	return s, nil
}

// RequireS that actual value is one of given int64 values
func (d *D) RequireS(vs ...int64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireS("require", s, false, true, vs...) })
}

// AssertS validate and asserts that actual value is one of given int64 values
func (d *D) AssertS(vs ...int64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireS("assert", s, true, !d.Options.Force, vs...) })
}

// ValidateS validates that actual value is one of given int64 values
func (d *D) ValidateS(vs ...int64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireS("validate", s, true, false, vs...) })
}

// Require/Assert/ValidatRange S

func requireRangeS(name string, s scalar.S, desc bool, fail bool, start, end int64) (scalar.S, error) {
	a := s.ActualS()
	if a >= start && a <= end {
		if desc {
			s.Description = "valid"
		}
		return s, nil
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s S range %v-%v", name, start, end)
	}
	return s, nil
}

// RequireSRange require that actual value is in range
func (d *D) RequireSRange(start, end int64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeS("require", s, false, true, start, end) })
}

// AssertSRange asserts that actual value is in range
func (d *D) AssertSRange(start, end int64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return requireRangeS("assert", s, true, !d.Options.Force, start, end)
	})
}

// ValidateSRange validates that actual value is in range
func (d *D) ValidateSRange(start, end int64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeS("validate", s, true, false, start, end) })
}

// Require/Assert/Validate Str

func requireStr(name string, s scalar.S, desc bool, fail bool, vs ...string) (scalar.S, error) {
	a := s.ActualStr()
	for _, b := range vs {
		if a == b {
			if desc {
				s.Description = "valid"
			}
			return s, nil
		}
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s Str", name)
	}
	return s, nil
}

// RequireStr that actual value is one of given string values
func (d *D) RequireStr(vs ...string) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireStr("require", s, false, true, vs...) })
}

// AssertStr validate and asserts that actual value is one of given string values
func (d *D) AssertStr(vs ...string) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireStr("assert", s, true, !d.Options.Force, vs...) })
}

// ValidateStr validates that actual value is one of given string values
func (d *D) ValidateStr(vs ...string) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireStr("validate", s, true, false, vs...) })
}

// Require/Assert/ValidatRange Str

func requireRangeStr(name string, s scalar.S, desc bool, fail bool, start, end string) (scalar.S, error) {
	a := s.ActualStr()
	if a >= start && a <= end {
		if desc {
			s.Description = "valid"
		}
		return s, nil
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s Str range %v-%v", name, start, end)
	}
	return s, nil
}

// RequireStrRange require that actual value is in range
func (d *D) RequireStrRange(start, end string) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeStr("require", s, false, true, start, end) })
}

// AssertStrRange asserts that actual value is in range
func (d *D) AssertStrRange(start, end string) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return requireRangeStr("assert", s, true, !d.Options.Force, start, end)
	})
}

// ValidateStrRange validates that actual value is in range
func (d *D) ValidateStrRange(start, end string) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeStr("validate", s, true, false, start, end) })
}

// Require/Assert/Validate U

func requireU(name string, s scalar.S, desc bool, fail bool, vs ...uint64) (scalar.S, error) {
	a := s.ActualU()
	for _, b := range vs {
		if a == b {
			if desc {
				s.Description = "valid"
			}
			return s, nil
		}
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s U", name)
	}
	return s, nil
}

// RequireU that actual value is one of given uint64 values
func (d *D) RequireU(vs ...uint64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireU("require", s, false, true, vs...) })
}

// AssertU validate and asserts that actual value is one of given uint64 values
func (d *D) AssertU(vs ...uint64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireU("assert", s, true, !d.Options.Force, vs...) })
}

// ValidateU validates that actual value is one of given uint64 values
func (d *D) ValidateU(vs ...uint64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireU("validate", s, true, false, vs...) })
}

// Require/Assert/ValidatRange U

func requireRangeU(name string, s scalar.S, desc bool, fail bool, start, end uint64) (scalar.S, error) {
	a := s.ActualU()
	if a >= start && a <= end {
		if desc {
			s.Description = "valid"
		}
		return s, nil
	}
	if desc {
		s.Description = "invalid"
	}
	if fail {
		return s, fmt.Errorf("failed to %s U range %v-%v", name, start, end)
	}
	return s, nil
}

// RequireURange require that actual value is in range
func (d *D) RequireURange(start, end uint64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeU("require", s, false, true, start, end) })
}

// AssertURange asserts that actual value is in range
func (d *D) AssertURange(start, end uint64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return requireRangeU("assert", s, true, !d.Options.Force, start, end)
	})
}

// ValidateURange validates that actual value is in range
func (d *D) ValidateURange(start, end uint64) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) { return requireRangeU("validate", s, true, false, start, end) })
}

// Reader RawLen

// TryRawLen tries to read nBits raw bits
func (d *D) TryRawLen(nBits int64) (bitio.ReaderAtSeeker, error) { return d.tryBitBuf(nBits) }

// RawLen reads nBits raw bits
func (d *D) RawLen(nBits int64) bitio.ReaderAtSeeker {
	v, err := d.tryBitBuf(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "RawLen", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarRawLen tries to add a field and read nBits raw bits
func (d *D) TryFieldScalarRawLen(name string, nBits int64, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBitBuf(nBits)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarRawLen adds a field and reads nBits raw bits
func (d *D) FieldScalarRawLen(name string, nBits int64, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarRawLen(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "RawLen", Pos: d.Pos()})
	}
	return s
}

// TryFieldRawLen tries to add a field and read nBits raw bits
func (d *D) TryFieldRawLen(name string, nBits int64, sms ...scalar.Mapper) (bitio.ReaderAtSeeker, error) {
	s, err := d.TryFieldScalarRawLen(name, nBits, sms...)
	return s.ActualBitBuf(), err
}

// FieldRawLen adds a field and reads nBits raw bits
func (d *D) FieldRawLen(name string, nBits int64, sms ...scalar.Mapper) bitio.ReaderAtSeeker {
	return d.FieldScalarRawLen(name, nBits, sms...).ActualBitBuf()
}

// Reader Bool

// TryBool tries to read 1 bit boolean
func (d *D) TryBool() (bool, error) { return d.tryBool() }

// Bool reads 1 bit boolean
func (d *D) Bool() bool {
	v, err := d.tryBool()
	if err != nil {
		panic(IOError{Err: err, Op: "Bool", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarBool tries to add a field and read 1 bit boolean
func (d *D) TryFieldScalarBool(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBool()
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarBool adds a field and reads 1 bit boolean
func (d *D) FieldScalarBool(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarBool(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Bool", Pos: d.Pos()})
	}
	return s
}

// TryFieldBool tries to add a field and read 1 bit boolean
func (d *D) TryFieldBool(name string, sms ...scalar.Mapper) (bool, error) {
	s, err := d.TryFieldScalarBool(name, sms...)
	return s.ActualBool(), err
}

// FieldBool adds a field and reads 1 bit boolean
func (d *D) FieldBool(name string, sms ...scalar.Mapper) bool {
	return d.FieldScalarBool(name, sms...).ActualBool()
}

// Reader U

// TryU tries to read nBits bits unsigned integer in current endian
func (d *D) TryU(nBits int) (uint64, error) { return d.tryUEndian(nBits, d.Endian) }

// U reads nBits bits unsigned integer in current endian
func (d *D) U(nBits int) uint64 {
	v, err := d.tryUEndian(nBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU tries to add a field and read nBits bits unsigned integer in current endian
func (d *D) TryFieldScalarU(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(nBits, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU adds a field and reads nBits bits unsigned integer in current endian
func (d *D) FieldScalarU(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U", Pos: d.Pos()})
	}
	return s
}

// TryFieldU tries to add a field and read nBits bits unsigned integer in current endian
func (d *D) TryFieldU(name string, nBits int, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU(name, nBits, sms...)
	return s.ActualU(), err
}

// FieldU adds a field and reads nBits bits unsigned integer in current endian
func (d *D) FieldU(name string, nBits int, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU(name, nBits, sms...).ActualU()
}

// Reader UE

// TryUE tries to read nBits unsigned integer in specified endian
func (d *D) TryUE(nBits int, endian Endian) (uint64, error) { return d.tryUEndian(nBits, endian) }

// UE reads nBits unsigned integer in specified endian
func (d *D) UE(nBits int, endian Endian) uint64 {
	v, err := d.tryUEndian(nBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "UE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUE tries to add a field and read nBits unsigned integer in specified endian
func (d *D) TryFieldScalarUE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(nBits, endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUE adds a field and reads nBits unsigned integer in specified endian
func (d *D) FieldScalarUE(name string, nBits int, endian Endian, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUE(name, nBits, endian, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UE", Pos: d.Pos()})
	}
	return s
}

// TryFieldUE tries to add a field and read nBits unsigned integer in specified endian
func (d *D) TryFieldUE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarUE(name, nBits, endian, sms...)
	return s.ActualU(), err
}

// FieldUE adds a field and reads nBits unsigned integer in specified endian
func (d *D) FieldUE(name string, nBits int, endian Endian, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarUE(name, nBits, endian, sms...).ActualU()
}

// Reader U1

// TryU1 tries to read 1 bit unsigned integer in current endian
func (d *D) TryU1() (uint64, error) { return d.tryUEndian(1, d.Endian) }

// U1 reads 1 bit unsigned integer in current endian
func (d *D) U1() uint64 {
	v, err := d.tryUEndian(1, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U1", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU1 tries to add a field and read 1 bit unsigned integer in current endian
func (d *D) TryFieldScalarU1(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(1, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU1 adds a field and reads 1 bit unsigned integer in current endian
func (d *D) FieldScalarU1(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU1(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U1", Pos: d.Pos()})
	}
	return s
}

// TryFieldU1 tries to add a field and read 1 bit unsigned integer in current endian
func (d *D) TryFieldU1(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU1(name, sms...)
	return s.ActualU(), err
}

// FieldU1 adds a field and reads 1 bit unsigned integer in current endian
func (d *D) FieldU1(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU1(name, sms...).ActualU()
}

// Reader U2

// TryU2 tries to read 2 bit unsigned integer in current endian
func (d *D) TryU2() (uint64, error) { return d.tryUEndian(2, d.Endian) }

// U2 reads 2 bit unsigned integer in current endian
func (d *D) U2() uint64 {
	v, err := d.tryUEndian(2, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U2", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU2 tries to add a field and read 2 bit unsigned integer in current endian
func (d *D) TryFieldScalarU2(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(2, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU2 adds a field and reads 2 bit unsigned integer in current endian
func (d *D) FieldScalarU2(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU2(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U2", Pos: d.Pos()})
	}
	return s
}

// TryFieldU2 tries to add a field and read 2 bit unsigned integer in current endian
func (d *D) TryFieldU2(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU2(name, sms...)
	return s.ActualU(), err
}

// FieldU2 adds a field and reads 2 bit unsigned integer in current endian
func (d *D) FieldU2(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU2(name, sms...).ActualU()
}

// Reader U3

// TryU3 tries to read 3 bit unsigned integer in current endian
func (d *D) TryU3() (uint64, error) { return d.tryUEndian(3, d.Endian) }

// U3 reads 3 bit unsigned integer in current endian
func (d *D) U3() uint64 {
	v, err := d.tryUEndian(3, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U3", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU3 tries to add a field and read 3 bit unsigned integer in current endian
func (d *D) TryFieldScalarU3(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(3, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU3 adds a field and reads 3 bit unsigned integer in current endian
func (d *D) FieldScalarU3(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU3(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U3", Pos: d.Pos()})
	}
	return s
}

// TryFieldU3 tries to add a field and read 3 bit unsigned integer in current endian
func (d *D) TryFieldU3(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU3(name, sms...)
	return s.ActualU(), err
}

// FieldU3 adds a field and reads 3 bit unsigned integer in current endian
func (d *D) FieldU3(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU3(name, sms...).ActualU()
}

// Reader U4

// TryU4 tries to read 4 bit unsigned integer in current endian
func (d *D) TryU4() (uint64, error) { return d.tryUEndian(4, d.Endian) }

// U4 reads 4 bit unsigned integer in current endian
func (d *D) U4() uint64 {
	v, err := d.tryUEndian(4, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U4", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU4 tries to add a field and read 4 bit unsigned integer in current endian
func (d *D) TryFieldScalarU4(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(4, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU4 adds a field and reads 4 bit unsigned integer in current endian
func (d *D) FieldScalarU4(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU4(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U4", Pos: d.Pos()})
	}
	return s
}

// TryFieldU4 tries to add a field and read 4 bit unsigned integer in current endian
func (d *D) TryFieldU4(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU4(name, sms...)
	return s.ActualU(), err
}

// FieldU4 adds a field and reads 4 bit unsigned integer in current endian
func (d *D) FieldU4(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU4(name, sms...).ActualU()
}

// Reader U5

// TryU5 tries to read 5 bit unsigned integer in current endian
func (d *D) TryU5() (uint64, error) { return d.tryUEndian(5, d.Endian) }

// U5 reads 5 bit unsigned integer in current endian
func (d *D) U5() uint64 {
	v, err := d.tryUEndian(5, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U5", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU5 tries to add a field and read 5 bit unsigned integer in current endian
func (d *D) TryFieldScalarU5(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(5, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU5 adds a field and reads 5 bit unsigned integer in current endian
func (d *D) FieldScalarU5(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU5(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U5", Pos: d.Pos()})
	}
	return s
}

// TryFieldU5 tries to add a field and read 5 bit unsigned integer in current endian
func (d *D) TryFieldU5(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU5(name, sms...)
	return s.ActualU(), err
}

// FieldU5 adds a field and reads 5 bit unsigned integer in current endian
func (d *D) FieldU5(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU5(name, sms...).ActualU()
}

// Reader U6

// TryU6 tries to read 6 bit unsigned integer in current endian
func (d *D) TryU6() (uint64, error) { return d.tryUEndian(6, d.Endian) }

// U6 reads 6 bit unsigned integer in current endian
func (d *D) U6() uint64 {
	v, err := d.tryUEndian(6, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U6", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU6 tries to add a field and read 6 bit unsigned integer in current endian
func (d *D) TryFieldScalarU6(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(6, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU6 adds a field and reads 6 bit unsigned integer in current endian
func (d *D) FieldScalarU6(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU6(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U6", Pos: d.Pos()})
	}
	return s
}

// TryFieldU6 tries to add a field and read 6 bit unsigned integer in current endian
func (d *D) TryFieldU6(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU6(name, sms...)
	return s.ActualU(), err
}

// FieldU6 adds a field and reads 6 bit unsigned integer in current endian
func (d *D) FieldU6(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU6(name, sms...).ActualU()
}

// Reader U7

// TryU7 tries to read 7 bit unsigned integer in current endian
func (d *D) TryU7() (uint64, error) { return d.tryUEndian(7, d.Endian) }

// U7 reads 7 bit unsigned integer in current endian
func (d *D) U7() uint64 {
	v, err := d.tryUEndian(7, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U7", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU7 tries to add a field and read 7 bit unsigned integer in current endian
func (d *D) TryFieldScalarU7(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(7, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU7 adds a field and reads 7 bit unsigned integer in current endian
func (d *D) FieldScalarU7(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU7(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U7", Pos: d.Pos()})
	}
	return s
}

// TryFieldU7 tries to add a field and read 7 bit unsigned integer in current endian
func (d *D) TryFieldU7(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU7(name, sms...)
	return s.ActualU(), err
}

// FieldU7 adds a field and reads 7 bit unsigned integer in current endian
func (d *D) FieldU7(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU7(name, sms...).ActualU()
}

// Reader U8

// TryU8 tries to read 8 bit unsigned integer in current endian
func (d *D) TryU8() (uint64, error) { return d.tryUEndian(8, d.Endian) }

// U8 reads 8 bit unsigned integer in current endian
func (d *D) U8() uint64 {
	v, err := d.tryUEndian(8, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U8", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU8 tries to add a field and read 8 bit unsigned integer in current endian
func (d *D) TryFieldScalarU8(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(8, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU8 adds a field and reads 8 bit unsigned integer in current endian
func (d *D) FieldScalarU8(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU8(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U8", Pos: d.Pos()})
	}
	return s
}

// TryFieldU8 tries to add a field and read 8 bit unsigned integer in current endian
func (d *D) TryFieldU8(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU8(name, sms...)
	return s.ActualU(), err
}

// FieldU8 adds a field and reads 8 bit unsigned integer in current endian
func (d *D) FieldU8(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU8(name, sms...).ActualU()
}

// Reader U9

// TryU9 tries to read 9 bit unsigned integer in current endian
func (d *D) TryU9() (uint64, error) { return d.tryUEndian(9, d.Endian) }

// U9 reads 9 bit unsigned integer in current endian
func (d *D) U9() uint64 {
	v, err := d.tryUEndian(9, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U9", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU9 tries to add a field and read 9 bit unsigned integer in current endian
func (d *D) TryFieldScalarU9(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(9, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU9 adds a field and reads 9 bit unsigned integer in current endian
func (d *D) FieldScalarU9(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU9(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U9", Pos: d.Pos()})
	}
	return s
}

// TryFieldU9 tries to add a field and read 9 bit unsigned integer in current endian
func (d *D) TryFieldU9(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU9(name, sms...)
	return s.ActualU(), err
}

// FieldU9 adds a field and reads 9 bit unsigned integer in current endian
func (d *D) FieldU9(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU9(name, sms...).ActualU()
}

// Reader U10

// TryU10 tries to read 10 bit unsigned integer in current endian
func (d *D) TryU10() (uint64, error) { return d.tryUEndian(10, d.Endian) }

// U10 reads 10 bit unsigned integer in current endian
func (d *D) U10() uint64 {
	v, err := d.tryUEndian(10, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U10", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU10 tries to add a field and read 10 bit unsigned integer in current endian
func (d *D) TryFieldScalarU10(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(10, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU10 adds a field and reads 10 bit unsigned integer in current endian
func (d *D) FieldScalarU10(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU10(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U10", Pos: d.Pos()})
	}
	return s
}

// TryFieldU10 tries to add a field and read 10 bit unsigned integer in current endian
func (d *D) TryFieldU10(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU10(name, sms...)
	return s.ActualU(), err
}

// FieldU10 adds a field and reads 10 bit unsigned integer in current endian
func (d *D) FieldU10(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU10(name, sms...).ActualU()
}

// Reader U11

// TryU11 tries to read 11 bit unsigned integer in current endian
func (d *D) TryU11() (uint64, error) { return d.tryUEndian(11, d.Endian) }

// U11 reads 11 bit unsigned integer in current endian
func (d *D) U11() uint64 {
	v, err := d.tryUEndian(11, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U11", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU11 tries to add a field and read 11 bit unsigned integer in current endian
func (d *D) TryFieldScalarU11(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(11, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU11 adds a field and reads 11 bit unsigned integer in current endian
func (d *D) FieldScalarU11(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU11(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U11", Pos: d.Pos()})
	}
	return s
}

// TryFieldU11 tries to add a field and read 11 bit unsigned integer in current endian
func (d *D) TryFieldU11(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU11(name, sms...)
	return s.ActualU(), err
}

// FieldU11 adds a field and reads 11 bit unsigned integer in current endian
func (d *D) FieldU11(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU11(name, sms...).ActualU()
}

// Reader U12

// TryU12 tries to read 12 bit unsigned integer in current endian
func (d *D) TryU12() (uint64, error) { return d.tryUEndian(12, d.Endian) }

// U12 reads 12 bit unsigned integer in current endian
func (d *D) U12() uint64 {
	v, err := d.tryUEndian(12, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U12", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU12 tries to add a field and read 12 bit unsigned integer in current endian
func (d *D) TryFieldScalarU12(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(12, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU12 adds a field and reads 12 bit unsigned integer in current endian
func (d *D) FieldScalarU12(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU12(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U12", Pos: d.Pos()})
	}
	return s
}

// TryFieldU12 tries to add a field and read 12 bit unsigned integer in current endian
func (d *D) TryFieldU12(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU12(name, sms...)
	return s.ActualU(), err
}

// FieldU12 adds a field and reads 12 bit unsigned integer in current endian
func (d *D) FieldU12(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU12(name, sms...).ActualU()
}

// Reader U13

// TryU13 tries to read 13 bit unsigned integer in current endian
func (d *D) TryU13() (uint64, error) { return d.tryUEndian(13, d.Endian) }

// U13 reads 13 bit unsigned integer in current endian
func (d *D) U13() uint64 {
	v, err := d.tryUEndian(13, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U13", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU13 tries to add a field and read 13 bit unsigned integer in current endian
func (d *D) TryFieldScalarU13(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(13, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU13 adds a field and reads 13 bit unsigned integer in current endian
func (d *D) FieldScalarU13(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU13(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U13", Pos: d.Pos()})
	}
	return s
}

// TryFieldU13 tries to add a field and read 13 bit unsigned integer in current endian
func (d *D) TryFieldU13(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU13(name, sms...)
	return s.ActualU(), err
}

// FieldU13 adds a field and reads 13 bit unsigned integer in current endian
func (d *D) FieldU13(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU13(name, sms...).ActualU()
}

// Reader U14

// TryU14 tries to read 14 bit unsigned integer in current endian
func (d *D) TryU14() (uint64, error) { return d.tryUEndian(14, d.Endian) }

// U14 reads 14 bit unsigned integer in current endian
func (d *D) U14() uint64 {
	v, err := d.tryUEndian(14, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U14", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU14 tries to add a field and read 14 bit unsigned integer in current endian
func (d *D) TryFieldScalarU14(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(14, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU14 adds a field and reads 14 bit unsigned integer in current endian
func (d *D) FieldScalarU14(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU14(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U14", Pos: d.Pos()})
	}
	return s
}

// TryFieldU14 tries to add a field and read 14 bit unsigned integer in current endian
func (d *D) TryFieldU14(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU14(name, sms...)
	return s.ActualU(), err
}

// FieldU14 adds a field and reads 14 bit unsigned integer in current endian
func (d *D) FieldU14(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU14(name, sms...).ActualU()
}

// Reader U15

// TryU15 tries to read 15 bit unsigned integer in current endian
func (d *D) TryU15() (uint64, error) { return d.tryUEndian(15, d.Endian) }

// U15 reads 15 bit unsigned integer in current endian
func (d *D) U15() uint64 {
	v, err := d.tryUEndian(15, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U15", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU15 tries to add a field and read 15 bit unsigned integer in current endian
func (d *D) TryFieldScalarU15(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(15, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU15 adds a field and reads 15 bit unsigned integer in current endian
func (d *D) FieldScalarU15(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU15(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U15", Pos: d.Pos()})
	}
	return s
}

// TryFieldU15 tries to add a field and read 15 bit unsigned integer in current endian
func (d *D) TryFieldU15(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU15(name, sms...)
	return s.ActualU(), err
}

// FieldU15 adds a field and reads 15 bit unsigned integer in current endian
func (d *D) FieldU15(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU15(name, sms...).ActualU()
}

// Reader U16

// TryU16 tries to read 16 bit unsigned integer in current endian
func (d *D) TryU16() (uint64, error) { return d.tryUEndian(16, d.Endian) }

// U16 reads 16 bit unsigned integer in current endian
func (d *D) U16() uint64 {
	v, err := d.tryUEndian(16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U16", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU16 tries to add a field and read 16 bit unsigned integer in current endian
func (d *D) TryFieldScalarU16(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(16, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU16 adds a field and reads 16 bit unsigned integer in current endian
func (d *D) FieldScalarU16(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU16(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U16", Pos: d.Pos()})
	}
	return s
}

// TryFieldU16 tries to add a field and read 16 bit unsigned integer in current endian
func (d *D) TryFieldU16(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU16(name, sms...)
	return s.ActualU(), err
}

// FieldU16 adds a field and reads 16 bit unsigned integer in current endian
func (d *D) FieldU16(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU16(name, sms...).ActualU()
}

// Reader U17

// TryU17 tries to read 17 bit unsigned integer in current endian
func (d *D) TryU17() (uint64, error) { return d.tryUEndian(17, d.Endian) }

// U17 reads 17 bit unsigned integer in current endian
func (d *D) U17() uint64 {
	v, err := d.tryUEndian(17, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U17", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU17 tries to add a field and read 17 bit unsigned integer in current endian
func (d *D) TryFieldScalarU17(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(17, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU17 adds a field and reads 17 bit unsigned integer in current endian
func (d *D) FieldScalarU17(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU17(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U17", Pos: d.Pos()})
	}
	return s
}

// TryFieldU17 tries to add a field and read 17 bit unsigned integer in current endian
func (d *D) TryFieldU17(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU17(name, sms...)
	return s.ActualU(), err
}

// FieldU17 adds a field and reads 17 bit unsigned integer in current endian
func (d *D) FieldU17(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU17(name, sms...).ActualU()
}

// Reader U18

// TryU18 tries to read 18 bit unsigned integer in current endian
func (d *D) TryU18() (uint64, error) { return d.tryUEndian(18, d.Endian) }

// U18 reads 18 bit unsigned integer in current endian
func (d *D) U18() uint64 {
	v, err := d.tryUEndian(18, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U18", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU18 tries to add a field and read 18 bit unsigned integer in current endian
func (d *D) TryFieldScalarU18(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(18, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU18 adds a field and reads 18 bit unsigned integer in current endian
func (d *D) FieldScalarU18(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU18(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U18", Pos: d.Pos()})
	}
	return s
}

// TryFieldU18 tries to add a field and read 18 bit unsigned integer in current endian
func (d *D) TryFieldU18(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU18(name, sms...)
	return s.ActualU(), err
}

// FieldU18 adds a field and reads 18 bit unsigned integer in current endian
func (d *D) FieldU18(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU18(name, sms...).ActualU()
}

// Reader U19

// TryU19 tries to read 19 bit unsigned integer in current endian
func (d *D) TryU19() (uint64, error) { return d.tryUEndian(19, d.Endian) }

// U19 reads 19 bit unsigned integer in current endian
func (d *D) U19() uint64 {
	v, err := d.tryUEndian(19, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U19", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU19 tries to add a field and read 19 bit unsigned integer in current endian
func (d *D) TryFieldScalarU19(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(19, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU19 adds a field and reads 19 bit unsigned integer in current endian
func (d *D) FieldScalarU19(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU19(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U19", Pos: d.Pos()})
	}
	return s
}

// TryFieldU19 tries to add a field and read 19 bit unsigned integer in current endian
func (d *D) TryFieldU19(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU19(name, sms...)
	return s.ActualU(), err
}

// FieldU19 adds a field and reads 19 bit unsigned integer in current endian
func (d *D) FieldU19(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU19(name, sms...).ActualU()
}

// Reader U20

// TryU20 tries to read 20 bit unsigned integer in current endian
func (d *D) TryU20() (uint64, error) { return d.tryUEndian(20, d.Endian) }

// U20 reads 20 bit unsigned integer in current endian
func (d *D) U20() uint64 {
	v, err := d.tryUEndian(20, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U20", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU20 tries to add a field and read 20 bit unsigned integer in current endian
func (d *D) TryFieldScalarU20(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(20, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU20 adds a field and reads 20 bit unsigned integer in current endian
func (d *D) FieldScalarU20(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU20(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U20", Pos: d.Pos()})
	}
	return s
}

// TryFieldU20 tries to add a field and read 20 bit unsigned integer in current endian
func (d *D) TryFieldU20(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU20(name, sms...)
	return s.ActualU(), err
}

// FieldU20 adds a field and reads 20 bit unsigned integer in current endian
func (d *D) FieldU20(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU20(name, sms...).ActualU()
}

// Reader U21

// TryU21 tries to read 21 bit unsigned integer in current endian
func (d *D) TryU21() (uint64, error) { return d.tryUEndian(21, d.Endian) }

// U21 reads 21 bit unsigned integer in current endian
func (d *D) U21() uint64 {
	v, err := d.tryUEndian(21, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U21", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU21 tries to add a field and read 21 bit unsigned integer in current endian
func (d *D) TryFieldScalarU21(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(21, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU21 adds a field and reads 21 bit unsigned integer in current endian
func (d *D) FieldScalarU21(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU21(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U21", Pos: d.Pos()})
	}
	return s
}

// TryFieldU21 tries to add a field and read 21 bit unsigned integer in current endian
func (d *D) TryFieldU21(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU21(name, sms...)
	return s.ActualU(), err
}

// FieldU21 adds a field and reads 21 bit unsigned integer in current endian
func (d *D) FieldU21(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU21(name, sms...).ActualU()
}

// Reader U22

// TryU22 tries to read 22 bit unsigned integer in current endian
func (d *D) TryU22() (uint64, error) { return d.tryUEndian(22, d.Endian) }

// U22 reads 22 bit unsigned integer in current endian
func (d *D) U22() uint64 {
	v, err := d.tryUEndian(22, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U22", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU22 tries to add a field and read 22 bit unsigned integer in current endian
func (d *D) TryFieldScalarU22(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(22, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU22 adds a field and reads 22 bit unsigned integer in current endian
func (d *D) FieldScalarU22(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU22(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U22", Pos: d.Pos()})
	}
	return s
}

// TryFieldU22 tries to add a field and read 22 bit unsigned integer in current endian
func (d *D) TryFieldU22(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU22(name, sms...)
	return s.ActualU(), err
}

// FieldU22 adds a field and reads 22 bit unsigned integer in current endian
func (d *D) FieldU22(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU22(name, sms...).ActualU()
}

// Reader U23

// TryU23 tries to read 23 bit unsigned integer in current endian
func (d *D) TryU23() (uint64, error) { return d.tryUEndian(23, d.Endian) }

// U23 reads 23 bit unsigned integer in current endian
func (d *D) U23() uint64 {
	v, err := d.tryUEndian(23, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U23", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU23 tries to add a field and read 23 bit unsigned integer in current endian
func (d *D) TryFieldScalarU23(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(23, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU23 adds a field and reads 23 bit unsigned integer in current endian
func (d *D) FieldScalarU23(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU23(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U23", Pos: d.Pos()})
	}
	return s
}

// TryFieldU23 tries to add a field and read 23 bit unsigned integer in current endian
func (d *D) TryFieldU23(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU23(name, sms...)
	return s.ActualU(), err
}

// FieldU23 adds a field and reads 23 bit unsigned integer in current endian
func (d *D) FieldU23(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU23(name, sms...).ActualU()
}

// Reader U24

// TryU24 tries to read 24 bit unsigned integer in current endian
func (d *D) TryU24() (uint64, error) { return d.tryUEndian(24, d.Endian) }

// U24 reads 24 bit unsigned integer in current endian
func (d *D) U24() uint64 {
	v, err := d.tryUEndian(24, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U24", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU24 tries to add a field and read 24 bit unsigned integer in current endian
func (d *D) TryFieldScalarU24(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(24, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU24 adds a field and reads 24 bit unsigned integer in current endian
func (d *D) FieldScalarU24(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU24(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U24", Pos: d.Pos()})
	}
	return s
}

// TryFieldU24 tries to add a field and read 24 bit unsigned integer in current endian
func (d *D) TryFieldU24(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU24(name, sms...)
	return s.ActualU(), err
}

// FieldU24 adds a field and reads 24 bit unsigned integer in current endian
func (d *D) FieldU24(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU24(name, sms...).ActualU()
}

// Reader U25

// TryU25 tries to read 25 bit unsigned integer in current endian
func (d *D) TryU25() (uint64, error) { return d.tryUEndian(25, d.Endian) }

// U25 reads 25 bit unsigned integer in current endian
func (d *D) U25() uint64 {
	v, err := d.tryUEndian(25, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U25", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU25 tries to add a field and read 25 bit unsigned integer in current endian
func (d *D) TryFieldScalarU25(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(25, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU25 adds a field and reads 25 bit unsigned integer in current endian
func (d *D) FieldScalarU25(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU25(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U25", Pos: d.Pos()})
	}
	return s
}

// TryFieldU25 tries to add a field and read 25 bit unsigned integer in current endian
func (d *D) TryFieldU25(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU25(name, sms...)
	return s.ActualU(), err
}

// FieldU25 adds a field and reads 25 bit unsigned integer in current endian
func (d *D) FieldU25(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU25(name, sms...).ActualU()
}

// Reader U26

// TryU26 tries to read 26 bit unsigned integer in current endian
func (d *D) TryU26() (uint64, error) { return d.tryUEndian(26, d.Endian) }

// U26 reads 26 bit unsigned integer in current endian
func (d *D) U26() uint64 {
	v, err := d.tryUEndian(26, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U26", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU26 tries to add a field and read 26 bit unsigned integer in current endian
func (d *D) TryFieldScalarU26(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(26, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU26 adds a field and reads 26 bit unsigned integer in current endian
func (d *D) FieldScalarU26(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU26(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U26", Pos: d.Pos()})
	}
	return s
}

// TryFieldU26 tries to add a field and read 26 bit unsigned integer in current endian
func (d *D) TryFieldU26(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU26(name, sms...)
	return s.ActualU(), err
}

// FieldU26 adds a field and reads 26 bit unsigned integer in current endian
func (d *D) FieldU26(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU26(name, sms...).ActualU()
}

// Reader U27

// TryU27 tries to read 27 bit unsigned integer in current endian
func (d *D) TryU27() (uint64, error) { return d.tryUEndian(27, d.Endian) }

// U27 reads 27 bit unsigned integer in current endian
func (d *D) U27() uint64 {
	v, err := d.tryUEndian(27, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U27", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU27 tries to add a field and read 27 bit unsigned integer in current endian
func (d *D) TryFieldScalarU27(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(27, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU27 adds a field and reads 27 bit unsigned integer in current endian
func (d *D) FieldScalarU27(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU27(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U27", Pos: d.Pos()})
	}
	return s
}

// TryFieldU27 tries to add a field and read 27 bit unsigned integer in current endian
func (d *D) TryFieldU27(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU27(name, sms...)
	return s.ActualU(), err
}

// FieldU27 adds a field and reads 27 bit unsigned integer in current endian
func (d *D) FieldU27(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU27(name, sms...).ActualU()
}

// Reader U28

// TryU28 tries to read 28 bit unsigned integer in current endian
func (d *D) TryU28() (uint64, error) { return d.tryUEndian(28, d.Endian) }

// U28 reads 28 bit unsigned integer in current endian
func (d *D) U28() uint64 {
	v, err := d.tryUEndian(28, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U28", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU28 tries to add a field and read 28 bit unsigned integer in current endian
func (d *D) TryFieldScalarU28(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(28, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU28 adds a field and reads 28 bit unsigned integer in current endian
func (d *D) FieldScalarU28(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU28(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U28", Pos: d.Pos()})
	}
	return s
}

// TryFieldU28 tries to add a field and read 28 bit unsigned integer in current endian
func (d *D) TryFieldU28(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU28(name, sms...)
	return s.ActualU(), err
}

// FieldU28 adds a field and reads 28 bit unsigned integer in current endian
func (d *D) FieldU28(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU28(name, sms...).ActualU()
}

// Reader U29

// TryU29 tries to read 29 bit unsigned integer in current endian
func (d *D) TryU29() (uint64, error) { return d.tryUEndian(29, d.Endian) }

// U29 reads 29 bit unsigned integer in current endian
func (d *D) U29() uint64 {
	v, err := d.tryUEndian(29, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U29", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU29 tries to add a field and read 29 bit unsigned integer in current endian
func (d *D) TryFieldScalarU29(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(29, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU29 adds a field and reads 29 bit unsigned integer in current endian
func (d *D) FieldScalarU29(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU29(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U29", Pos: d.Pos()})
	}
	return s
}

// TryFieldU29 tries to add a field and read 29 bit unsigned integer in current endian
func (d *D) TryFieldU29(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU29(name, sms...)
	return s.ActualU(), err
}

// FieldU29 adds a field and reads 29 bit unsigned integer in current endian
func (d *D) FieldU29(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU29(name, sms...).ActualU()
}

// Reader U30

// TryU30 tries to read 30 bit unsigned integer in current endian
func (d *D) TryU30() (uint64, error) { return d.tryUEndian(30, d.Endian) }

// U30 reads 30 bit unsigned integer in current endian
func (d *D) U30() uint64 {
	v, err := d.tryUEndian(30, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U30", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU30 tries to add a field and read 30 bit unsigned integer in current endian
func (d *D) TryFieldScalarU30(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(30, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU30 adds a field and reads 30 bit unsigned integer in current endian
func (d *D) FieldScalarU30(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU30(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U30", Pos: d.Pos()})
	}
	return s
}

// TryFieldU30 tries to add a field and read 30 bit unsigned integer in current endian
func (d *D) TryFieldU30(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU30(name, sms...)
	return s.ActualU(), err
}

// FieldU30 adds a field and reads 30 bit unsigned integer in current endian
func (d *D) FieldU30(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU30(name, sms...).ActualU()
}

// Reader U31

// TryU31 tries to read 31 bit unsigned integer in current endian
func (d *D) TryU31() (uint64, error) { return d.tryUEndian(31, d.Endian) }

// U31 reads 31 bit unsigned integer in current endian
func (d *D) U31() uint64 {
	v, err := d.tryUEndian(31, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U31", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU31 tries to add a field and read 31 bit unsigned integer in current endian
func (d *D) TryFieldScalarU31(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(31, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU31 adds a field and reads 31 bit unsigned integer in current endian
func (d *D) FieldScalarU31(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU31(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U31", Pos: d.Pos()})
	}
	return s
}

// TryFieldU31 tries to add a field and read 31 bit unsigned integer in current endian
func (d *D) TryFieldU31(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU31(name, sms...)
	return s.ActualU(), err
}

// FieldU31 adds a field and reads 31 bit unsigned integer in current endian
func (d *D) FieldU31(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU31(name, sms...).ActualU()
}

// Reader U32

// TryU32 tries to read 32 bit unsigned integer in current endian
func (d *D) TryU32() (uint64, error) { return d.tryUEndian(32, d.Endian) }

// U32 reads 32 bit unsigned integer in current endian
func (d *D) U32() uint64 {
	v, err := d.tryUEndian(32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U32", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU32 tries to add a field and read 32 bit unsigned integer in current endian
func (d *D) TryFieldScalarU32(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(32, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU32 adds a field and reads 32 bit unsigned integer in current endian
func (d *D) FieldScalarU32(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU32(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U32", Pos: d.Pos()})
	}
	return s
}

// TryFieldU32 tries to add a field and read 32 bit unsigned integer in current endian
func (d *D) TryFieldU32(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU32(name, sms...)
	return s.ActualU(), err
}

// FieldU32 adds a field and reads 32 bit unsigned integer in current endian
func (d *D) FieldU32(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU32(name, sms...).ActualU()
}

// Reader U33

// TryU33 tries to read 33 bit unsigned integer in current endian
func (d *D) TryU33() (uint64, error) { return d.tryUEndian(33, d.Endian) }

// U33 reads 33 bit unsigned integer in current endian
func (d *D) U33() uint64 {
	v, err := d.tryUEndian(33, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U33", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU33 tries to add a field and read 33 bit unsigned integer in current endian
func (d *D) TryFieldScalarU33(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(33, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU33 adds a field and reads 33 bit unsigned integer in current endian
func (d *D) FieldScalarU33(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU33(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U33", Pos: d.Pos()})
	}
	return s
}

// TryFieldU33 tries to add a field and read 33 bit unsigned integer in current endian
func (d *D) TryFieldU33(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU33(name, sms...)
	return s.ActualU(), err
}

// FieldU33 adds a field and reads 33 bit unsigned integer in current endian
func (d *D) FieldU33(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU33(name, sms...).ActualU()
}

// Reader U34

// TryU34 tries to read 34 bit unsigned integer in current endian
func (d *D) TryU34() (uint64, error) { return d.tryUEndian(34, d.Endian) }

// U34 reads 34 bit unsigned integer in current endian
func (d *D) U34() uint64 {
	v, err := d.tryUEndian(34, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U34", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU34 tries to add a field and read 34 bit unsigned integer in current endian
func (d *D) TryFieldScalarU34(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(34, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU34 adds a field and reads 34 bit unsigned integer in current endian
func (d *D) FieldScalarU34(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU34(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U34", Pos: d.Pos()})
	}
	return s
}

// TryFieldU34 tries to add a field and read 34 bit unsigned integer in current endian
func (d *D) TryFieldU34(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU34(name, sms...)
	return s.ActualU(), err
}

// FieldU34 adds a field and reads 34 bit unsigned integer in current endian
func (d *D) FieldU34(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU34(name, sms...).ActualU()
}

// Reader U35

// TryU35 tries to read 35 bit unsigned integer in current endian
func (d *D) TryU35() (uint64, error) { return d.tryUEndian(35, d.Endian) }

// U35 reads 35 bit unsigned integer in current endian
func (d *D) U35() uint64 {
	v, err := d.tryUEndian(35, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U35", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU35 tries to add a field and read 35 bit unsigned integer in current endian
func (d *D) TryFieldScalarU35(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(35, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU35 adds a field and reads 35 bit unsigned integer in current endian
func (d *D) FieldScalarU35(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU35(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U35", Pos: d.Pos()})
	}
	return s
}

// TryFieldU35 tries to add a field and read 35 bit unsigned integer in current endian
func (d *D) TryFieldU35(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU35(name, sms...)
	return s.ActualU(), err
}

// FieldU35 adds a field and reads 35 bit unsigned integer in current endian
func (d *D) FieldU35(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU35(name, sms...).ActualU()
}

// Reader U36

// TryU36 tries to read 36 bit unsigned integer in current endian
func (d *D) TryU36() (uint64, error) { return d.tryUEndian(36, d.Endian) }

// U36 reads 36 bit unsigned integer in current endian
func (d *D) U36() uint64 {
	v, err := d.tryUEndian(36, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U36", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU36 tries to add a field and read 36 bit unsigned integer in current endian
func (d *D) TryFieldScalarU36(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(36, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU36 adds a field and reads 36 bit unsigned integer in current endian
func (d *D) FieldScalarU36(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU36(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U36", Pos: d.Pos()})
	}
	return s
}

// TryFieldU36 tries to add a field and read 36 bit unsigned integer in current endian
func (d *D) TryFieldU36(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU36(name, sms...)
	return s.ActualU(), err
}

// FieldU36 adds a field and reads 36 bit unsigned integer in current endian
func (d *D) FieldU36(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU36(name, sms...).ActualU()
}

// Reader U37

// TryU37 tries to read 37 bit unsigned integer in current endian
func (d *D) TryU37() (uint64, error) { return d.tryUEndian(37, d.Endian) }

// U37 reads 37 bit unsigned integer in current endian
func (d *D) U37() uint64 {
	v, err := d.tryUEndian(37, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U37", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU37 tries to add a field and read 37 bit unsigned integer in current endian
func (d *D) TryFieldScalarU37(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(37, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU37 adds a field and reads 37 bit unsigned integer in current endian
func (d *D) FieldScalarU37(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU37(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U37", Pos: d.Pos()})
	}
	return s
}

// TryFieldU37 tries to add a field and read 37 bit unsigned integer in current endian
func (d *D) TryFieldU37(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU37(name, sms...)
	return s.ActualU(), err
}

// FieldU37 adds a field and reads 37 bit unsigned integer in current endian
func (d *D) FieldU37(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU37(name, sms...).ActualU()
}

// Reader U38

// TryU38 tries to read 38 bit unsigned integer in current endian
func (d *D) TryU38() (uint64, error) { return d.tryUEndian(38, d.Endian) }

// U38 reads 38 bit unsigned integer in current endian
func (d *D) U38() uint64 {
	v, err := d.tryUEndian(38, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U38", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU38 tries to add a field and read 38 bit unsigned integer in current endian
func (d *D) TryFieldScalarU38(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(38, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU38 adds a field and reads 38 bit unsigned integer in current endian
func (d *D) FieldScalarU38(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU38(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U38", Pos: d.Pos()})
	}
	return s
}

// TryFieldU38 tries to add a field and read 38 bit unsigned integer in current endian
func (d *D) TryFieldU38(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU38(name, sms...)
	return s.ActualU(), err
}

// FieldU38 adds a field and reads 38 bit unsigned integer in current endian
func (d *D) FieldU38(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU38(name, sms...).ActualU()
}

// Reader U39

// TryU39 tries to read 39 bit unsigned integer in current endian
func (d *D) TryU39() (uint64, error) { return d.tryUEndian(39, d.Endian) }

// U39 reads 39 bit unsigned integer in current endian
func (d *D) U39() uint64 {
	v, err := d.tryUEndian(39, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U39", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU39 tries to add a field and read 39 bit unsigned integer in current endian
func (d *D) TryFieldScalarU39(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(39, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU39 adds a field and reads 39 bit unsigned integer in current endian
func (d *D) FieldScalarU39(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU39(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U39", Pos: d.Pos()})
	}
	return s
}

// TryFieldU39 tries to add a field and read 39 bit unsigned integer in current endian
func (d *D) TryFieldU39(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU39(name, sms...)
	return s.ActualU(), err
}

// FieldU39 adds a field and reads 39 bit unsigned integer in current endian
func (d *D) FieldU39(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU39(name, sms...).ActualU()
}

// Reader U40

// TryU40 tries to read 40 bit unsigned integer in current endian
func (d *D) TryU40() (uint64, error) { return d.tryUEndian(40, d.Endian) }

// U40 reads 40 bit unsigned integer in current endian
func (d *D) U40() uint64 {
	v, err := d.tryUEndian(40, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U40", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU40 tries to add a field and read 40 bit unsigned integer in current endian
func (d *D) TryFieldScalarU40(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(40, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU40 adds a field and reads 40 bit unsigned integer in current endian
func (d *D) FieldScalarU40(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU40(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U40", Pos: d.Pos()})
	}
	return s
}

// TryFieldU40 tries to add a field and read 40 bit unsigned integer in current endian
func (d *D) TryFieldU40(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU40(name, sms...)
	return s.ActualU(), err
}

// FieldU40 adds a field and reads 40 bit unsigned integer in current endian
func (d *D) FieldU40(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU40(name, sms...).ActualU()
}

// Reader U41

// TryU41 tries to read 41 bit unsigned integer in current endian
func (d *D) TryU41() (uint64, error) { return d.tryUEndian(41, d.Endian) }

// U41 reads 41 bit unsigned integer in current endian
func (d *D) U41() uint64 {
	v, err := d.tryUEndian(41, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U41", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU41 tries to add a field and read 41 bit unsigned integer in current endian
func (d *D) TryFieldScalarU41(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(41, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU41 adds a field and reads 41 bit unsigned integer in current endian
func (d *D) FieldScalarU41(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU41(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U41", Pos: d.Pos()})
	}
	return s
}

// TryFieldU41 tries to add a field and read 41 bit unsigned integer in current endian
func (d *D) TryFieldU41(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU41(name, sms...)
	return s.ActualU(), err
}

// FieldU41 adds a field and reads 41 bit unsigned integer in current endian
func (d *D) FieldU41(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU41(name, sms...).ActualU()
}

// Reader U42

// TryU42 tries to read 42 bit unsigned integer in current endian
func (d *D) TryU42() (uint64, error) { return d.tryUEndian(42, d.Endian) }

// U42 reads 42 bit unsigned integer in current endian
func (d *D) U42() uint64 {
	v, err := d.tryUEndian(42, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U42", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU42 tries to add a field and read 42 bit unsigned integer in current endian
func (d *D) TryFieldScalarU42(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(42, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU42 adds a field and reads 42 bit unsigned integer in current endian
func (d *D) FieldScalarU42(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU42(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U42", Pos: d.Pos()})
	}
	return s
}

// TryFieldU42 tries to add a field and read 42 bit unsigned integer in current endian
func (d *D) TryFieldU42(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU42(name, sms...)
	return s.ActualU(), err
}

// FieldU42 adds a field and reads 42 bit unsigned integer in current endian
func (d *D) FieldU42(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU42(name, sms...).ActualU()
}

// Reader U43

// TryU43 tries to read 43 bit unsigned integer in current endian
func (d *D) TryU43() (uint64, error) { return d.tryUEndian(43, d.Endian) }

// U43 reads 43 bit unsigned integer in current endian
func (d *D) U43() uint64 {
	v, err := d.tryUEndian(43, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U43", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU43 tries to add a field and read 43 bit unsigned integer in current endian
func (d *D) TryFieldScalarU43(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(43, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU43 adds a field and reads 43 bit unsigned integer in current endian
func (d *D) FieldScalarU43(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU43(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U43", Pos: d.Pos()})
	}
	return s
}

// TryFieldU43 tries to add a field and read 43 bit unsigned integer in current endian
func (d *D) TryFieldU43(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU43(name, sms...)
	return s.ActualU(), err
}

// FieldU43 adds a field and reads 43 bit unsigned integer in current endian
func (d *D) FieldU43(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU43(name, sms...).ActualU()
}

// Reader U44

// TryU44 tries to read 44 bit unsigned integer in current endian
func (d *D) TryU44() (uint64, error) { return d.tryUEndian(44, d.Endian) }

// U44 reads 44 bit unsigned integer in current endian
func (d *D) U44() uint64 {
	v, err := d.tryUEndian(44, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U44", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU44 tries to add a field and read 44 bit unsigned integer in current endian
func (d *D) TryFieldScalarU44(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(44, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU44 adds a field and reads 44 bit unsigned integer in current endian
func (d *D) FieldScalarU44(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU44(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U44", Pos: d.Pos()})
	}
	return s
}

// TryFieldU44 tries to add a field and read 44 bit unsigned integer in current endian
func (d *D) TryFieldU44(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU44(name, sms...)
	return s.ActualU(), err
}

// FieldU44 adds a field and reads 44 bit unsigned integer in current endian
func (d *D) FieldU44(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU44(name, sms...).ActualU()
}

// Reader U45

// TryU45 tries to read 45 bit unsigned integer in current endian
func (d *D) TryU45() (uint64, error) { return d.tryUEndian(45, d.Endian) }

// U45 reads 45 bit unsigned integer in current endian
func (d *D) U45() uint64 {
	v, err := d.tryUEndian(45, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U45", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU45 tries to add a field and read 45 bit unsigned integer in current endian
func (d *D) TryFieldScalarU45(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(45, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU45 adds a field and reads 45 bit unsigned integer in current endian
func (d *D) FieldScalarU45(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU45(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U45", Pos: d.Pos()})
	}
	return s
}

// TryFieldU45 tries to add a field and read 45 bit unsigned integer in current endian
func (d *D) TryFieldU45(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU45(name, sms...)
	return s.ActualU(), err
}

// FieldU45 adds a field and reads 45 bit unsigned integer in current endian
func (d *D) FieldU45(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU45(name, sms...).ActualU()
}

// Reader U46

// TryU46 tries to read 46 bit unsigned integer in current endian
func (d *D) TryU46() (uint64, error) { return d.tryUEndian(46, d.Endian) }

// U46 reads 46 bit unsigned integer in current endian
func (d *D) U46() uint64 {
	v, err := d.tryUEndian(46, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U46", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU46 tries to add a field and read 46 bit unsigned integer in current endian
func (d *D) TryFieldScalarU46(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(46, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU46 adds a field and reads 46 bit unsigned integer in current endian
func (d *D) FieldScalarU46(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU46(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U46", Pos: d.Pos()})
	}
	return s
}

// TryFieldU46 tries to add a field and read 46 bit unsigned integer in current endian
func (d *D) TryFieldU46(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU46(name, sms...)
	return s.ActualU(), err
}

// FieldU46 adds a field and reads 46 bit unsigned integer in current endian
func (d *D) FieldU46(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU46(name, sms...).ActualU()
}

// Reader U47

// TryU47 tries to read 47 bit unsigned integer in current endian
func (d *D) TryU47() (uint64, error) { return d.tryUEndian(47, d.Endian) }

// U47 reads 47 bit unsigned integer in current endian
func (d *D) U47() uint64 {
	v, err := d.tryUEndian(47, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U47", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU47 tries to add a field and read 47 bit unsigned integer in current endian
func (d *D) TryFieldScalarU47(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(47, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU47 adds a field and reads 47 bit unsigned integer in current endian
func (d *D) FieldScalarU47(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU47(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U47", Pos: d.Pos()})
	}
	return s
}

// TryFieldU47 tries to add a field and read 47 bit unsigned integer in current endian
func (d *D) TryFieldU47(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU47(name, sms...)
	return s.ActualU(), err
}

// FieldU47 adds a field and reads 47 bit unsigned integer in current endian
func (d *D) FieldU47(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU47(name, sms...).ActualU()
}

// Reader U48

// TryU48 tries to read 48 bit unsigned integer in current endian
func (d *D) TryU48() (uint64, error) { return d.tryUEndian(48, d.Endian) }

// U48 reads 48 bit unsigned integer in current endian
func (d *D) U48() uint64 {
	v, err := d.tryUEndian(48, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U48", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU48 tries to add a field and read 48 bit unsigned integer in current endian
func (d *D) TryFieldScalarU48(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(48, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU48 adds a field and reads 48 bit unsigned integer in current endian
func (d *D) FieldScalarU48(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU48(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U48", Pos: d.Pos()})
	}
	return s
}

// TryFieldU48 tries to add a field and read 48 bit unsigned integer in current endian
func (d *D) TryFieldU48(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU48(name, sms...)
	return s.ActualU(), err
}

// FieldU48 adds a field and reads 48 bit unsigned integer in current endian
func (d *D) FieldU48(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU48(name, sms...).ActualU()
}

// Reader U49

// TryU49 tries to read 49 bit unsigned integer in current endian
func (d *D) TryU49() (uint64, error) { return d.tryUEndian(49, d.Endian) }

// U49 reads 49 bit unsigned integer in current endian
func (d *D) U49() uint64 {
	v, err := d.tryUEndian(49, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U49", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU49 tries to add a field and read 49 bit unsigned integer in current endian
func (d *D) TryFieldScalarU49(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(49, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU49 adds a field and reads 49 bit unsigned integer in current endian
func (d *D) FieldScalarU49(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU49(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U49", Pos: d.Pos()})
	}
	return s
}

// TryFieldU49 tries to add a field and read 49 bit unsigned integer in current endian
func (d *D) TryFieldU49(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU49(name, sms...)
	return s.ActualU(), err
}

// FieldU49 adds a field and reads 49 bit unsigned integer in current endian
func (d *D) FieldU49(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU49(name, sms...).ActualU()
}

// Reader U50

// TryU50 tries to read 50 bit unsigned integer in current endian
func (d *D) TryU50() (uint64, error) { return d.tryUEndian(50, d.Endian) }

// U50 reads 50 bit unsigned integer in current endian
func (d *D) U50() uint64 {
	v, err := d.tryUEndian(50, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U50", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU50 tries to add a field and read 50 bit unsigned integer in current endian
func (d *D) TryFieldScalarU50(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(50, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU50 adds a field and reads 50 bit unsigned integer in current endian
func (d *D) FieldScalarU50(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU50(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U50", Pos: d.Pos()})
	}
	return s
}

// TryFieldU50 tries to add a field and read 50 bit unsigned integer in current endian
func (d *D) TryFieldU50(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU50(name, sms...)
	return s.ActualU(), err
}

// FieldU50 adds a field and reads 50 bit unsigned integer in current endian
func (d *D) FieldU50(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU50(name, sms...).ActualU()
}

// Reader U51

// TryU51 tries to read 51 bit unsigned integer in current endian
func (d *D) TryU51() (uint64, error) { return d.tryUEndian(51, d.Endian) }

// U51 reads 51 bit unsigned integer in current endian
func (d *D) U51() uint64 {
	v, err := d.tryUEndian(51, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U51", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU51 tries to add a field and read 51 bit unsigned integer in current endian
func (d *D) TryFieldScalarU51(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(51, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU51 adds a field and reads 51 bit unsigned integer in current endian
func (d *D) FieldScalarU51(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU51(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U51", Pos: d.Pos()})
	}
	return s
}

// TryFieldU51 tries to add a field and read 51 bit unsigned integer in current endian
func (d *D) TryFieldU51(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU51(name, sms...)
	return s.ActualU(), err
}

// FieldU51 adds a field and reads 51 bit unsigned integer in current endian
func (d *D) FieldU51(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU51(name, sms...).ActualU()
}

// Reader U52

// TryU52 tries to read 52 bit unsigned integer in current endian
func (d *D) TryU52() (uint64, error) { return d.tryUEndian(52, d.Endian) }

// U52 reads 52 bit unsigned integer in current endian
func (d *D) U52() uint64 {
	v, err := d.tryUEndian(52, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U52", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU52 tries to add a field and read 52 bit unsigned integer in current endian
func (d *D) TryFieldScalarU52(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(52, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU52 adds a field and reads 52 bit unsigned integer in current endian
func (d *D) FieldScalarU52(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU52(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U52", Pos: d.Pos()})
	}
	return s
}

// TryFieldU52 tries to add a field and read 52 bit unsigned integer in current endian
func (d *D) TryFieldU52(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU52(name, sms...)
	return s.ActualU(), err
}

// FieldU52 adds a field and reads 52 bit unsigned integer in current endian
func (d *D) FieldU52(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU52(name, sms...).ActualU()
}

// Reader U53

// TryU53 tries to read 53 bit unsigned integer in current endian
func (d *D) TryU53() (uint64, error) { return d.tryUEndian(53, d.Endian) }

// U53 reads 53 bit unsigned integer in current endian
func (d *D) U53() uint64 {
	v, err := d.tryUEndian(53, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U53", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU53 tries to add a field and read 53 bit unsigned integer in current endian
func (d *D) TryFieldScalarU53(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(53, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU53 adds a field and reads 53 bit unsigned integer in current endian
func (d *D) FieldScalarU53(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU53(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U53", Pos: d.Pos()})
	}
	return s
}

// TryFieldU53 tries to add a field and read 53 bit unsigned integer in current endian
func (d *D) TryFieldU53(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU53(name, sms...)
	return s.ActualU(), err
}

// FieldU53 adds a field and reads 53 bit unsigned integer in current endian
func (d *D) FieldU53(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU53(name, sms...).ActualU()
}

// Reader U54

// TryU54 tries to read 54 bit unsigned integer in current endian
func (d *D) TryU54() (uint64, error) { return d.tryUEndian(54, d.Endian) }

// U54 reads 54 bit unsigned integer in current endian
func (d *D) U54() uint64 {
	v, err := d.tryUEndian(54, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U54", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU54 tries to add a field and read 54 bit unsigned integer in current endian
func (d *D) TryFieldScalarU54(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(54, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU54 adds a field and reads 54 bit unsigned integer in current endian
func (d *D) FieldScalarU54(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU54(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U54", Pos: d.Pos()})
	}
	return s
}

// TryFieldU54 tries to add a field and read 54 bit unsigned integer in current endian
func (d *D) TryFieldU54(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU54(name, sms...)
	return s.ActualU(), err
}

// FieldU54 adds a field and reads 54 bit unsigned integer in current endian
func (d *D) FieldU54(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU54(name, sms...).ActualU()
}

// Reader U55

// TryU55 tries to read 55 bit unsigned integer in current endian
func (d *D) TryU55() (uint64, error) { return d.tryUEndian(55, d.Endian) }

// U55 reads 55 bit unsigned integer in current endian
func (d *D) U55() uint64 {
	v, err := d.tryUEndian(55, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U55", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU55 tries to add a field and read 55 bit unsigned integer in current endian
func (d *D) TryFieldScalarU55(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(55, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU55 adds a field and reads 55 bit unsigned integer in current endian
func (d *D) FieldScalarU55(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU55(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U55", Pos: d.Pos()})
	}
	return s
}

// TryFieldU55 tries to add a field and read 55 bit unsigned integer in current endian
func (d *D) TryFieldU55(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU55(name, sms...)
	return s.ActualU(), err
}

// FieldU55 adds a field and reads 55 bit unsigned integer in current endian
func (d *D) FieldU55(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU55(name, sms...).ActualU()
}

// Reader U56

// TryU56 tries to read 56 bit unsigned integer in current endian
func (d *D) TryU56() (uint64, error) { return d.tryUEndian(56, d.Endian) }

// U56 reads 56 bit unsigned integer in current endian
func (d *D) U56() uint64 {
	v, err := d.tryUEndian(56, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U56", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU56 tries to add a field and read 56 bit unsigned integer in current endian
func (d *D) TryFieldScalarU56(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(56, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU56 adds a field and reads 56 bit unsigned integer in current endian
func (d *D) FieldScalarU56(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU56(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U56", Pos: d.Pos()})
	}
	return s
}

// TryFieldU56 tries to add a field and read 56 bit unsigned integer in current endian
func (d *D) TryFieldU56(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU56(name, sms...)
	return s.ActualU(), err
}

// FieldU56 adds a field and reads 56 bit unsigned integer in current endian
func (d *D) FieldU56(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU56(name, sms...).ActualU()
}

// Reader U57

// TryU57 tries to read 57 bit unsigned integer in current endian
func (d *D) TryU57() (uint64, error) { return d.tryUEndian(57, d.Endian) }

// U57 reads 57 bit unsigned integer in current endian
func (d *D) U57() uint64 {
	v, err := d.tryUEndian(57, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U57", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU57 tries to add a field and read 57 bit unsigned integer in current endian
func (d *D) TryFieldScalarU57(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(57, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU57 adds a field and reads 57 bit unsigned integer in current endian
func (d *D) FieldScalarU57(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU57(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U57", Pos: d.Pos()})
	}
	return s
}

// TryFieldU57 tries to add a field and read 57 bit unsigned integer in current endian
func (d *D) TryFieldU57(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU57(name, sms...)
	return s.ActualU(), err
}

// FieldU57 adds a field and reads 57 bit unsigned integer in current endian
func (d *D) FieldU57(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU57(name, sms...).ActualU()
}

// Reader U58

// TryU58 tries to read 58 bit unsigned integer in current endian
func (d *D) TryU58() (uint64, error) { return d.tryUEndian(58, d.Endian) }

// U58 reads 58 bit unsigned integer in current endian
func (d *D) U58() uint64 {
	v, err := d.tryUEndian(58, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U58", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU58 tries to add a field and read 58 bit unsigned integer in current endian
func (d *D) TryFieldScalarU58(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(58, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU58 adds a field and reads 58 bit unsigned integer in current endian
func (d *D) FieldScalarU58(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU58(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U58", Pos: d.Pos()})
	}
	return s
}

// TryFieldU58 tries to add a field and read 58 bit unsigned integer in current endian
func (d *D) TryFieldU58(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU58(name, sms...)
	return s.ActualU(), err
}

// FieldU58 adds a field and reads 58 bit unsigned integer in current endian
func (d *D) FieldU58(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU58(name, sms...).ActualU()
}

// Reader U59

// TryU59 tries to read 59 bit unsigned integer in current endian
func (d *D) TryU59() (uint64, error) { return d.tryUEndian(59, d.Endian) }

// U59 reads 59 bit unsigned integer in current endian
func (d *D) U59() uint64 {
	v, err := d.tryUEndian(59, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U59", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU59 tries to add a field and read 59 bit unsigned integer in current endian
func (d *D) TryFieldScalarU59(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(59, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU59 adds a field and reads 59 bit unsigned integer in current endian
func (d *D) FieldScalarU59(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU59(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U59", Pos: d.Pos()})
	}
	return s
}

// TryFieldU59 tries to add a field and read 59 bit unsigned integer in current endian
func (d *D) TryFieldU59(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU59(name, sms...)
	return s.ActualU(), err
}

// FieldU59 adds a field and reads 59 bit unsigned integer in current endian
func (d *D) FieldU59(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU59(name, sms...).ActualU()
}

// Reader U60

// TryU60 tries to read 60 bit unsigned integer in current endian
func (d *D) TryU60() (uint64, error) { return d.tryUEndian(60, d.Endian) }

// U60 reads 60 bit unsigned integer in current endian
func (d *D) U60() uint64 {
	v, err := d.tryUEndian(60, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U60", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU60 tries to add a field and read 60 bit unsigned integer in current endian
func (d *D) TryFieldScalarU60(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(60, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU60 adds a field and reads 60 bit unsigned integer in current endian
func (d *D) FieldScalarU60(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU60(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U60", Pos: d.Pos()})
	}
	return s
}

// TryFieldU60 tries to add a field and read 60 bit unsigned integer in current endian
func (d *D) TryFieldU60(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU60(name, sms...)
	return s.ActualU(), err
}

// FieldU60 adds a field and reads 60 bit unsigned integer in current endian
func (d *D) FieldU60(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU60(name, sms...).ActualU()
}

// Reader U61

// TryU61 tries to read 61 bit unsigned integer in current endian
func (d *D) TryU61() (uint64, error) { return d.tryUEndian(61, d.Endian) }

// U61 reads 61 bit unsigned integer in current endian
func (d *D) U61() uint64 {
	v, err := d.tryUEndian(61, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U61", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU61 tries to add a field and read 61 bit unsigned integer in current endian
func (d *D) TryFieldScalarU61(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(61, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU61 adds a field and reads 61 bit unsigned integer in current endian
func (d *D) FieldScalarU61(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU61(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U61", Pos: d.Pos()})
	}
	return s
}

// TryFieldU61 tries to add a field and read 61 bit unsigned integer in current endian
func (d *D) TryFieldU61(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU61(name, sms...)
	return s.ActualU(), err
}

// FieldU61 adds a field and reads 61 bit unsigned integer in current endian
func (d *D) FieldU61(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU61(name, sms...).ActualU()
}

// Reader U62

// TryU62 tries to read 62 bit unsigned integer in current endian
func (d *D) TryU62() (uint64, error) { return d.tryUEndian(62, d.Endian) }

// U62 reads 62 bit unsigned integer in current endian
func (d *D) U62() uint64 {
	v, err := d.tryUEndian(62, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U62", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU62 tries to add a field and read 62 bit unsigned integer in current endian
func (d *D) TryFieldScalarU62(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(62, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU62 adds a field and reads 62 bit unsigned integer in current endian
func (d *D) FieldScalarU62(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU62(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U62", Pos: d.Pos()})
	}
	return s
}

// TryFieldU62 tries to add a field and read 62 bit unsigned integer in current endian
func (d *D) TryFieldU62(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU62(name, sms...)
	return s.ActualU(), err
}

// FieldU62 adds a field and reads 62 bit unsigned integer in current endian
func (d *D) FieldU62(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU62(name, sms...).ActualU()
}

// Reader U63

// TryU63 tries to read 63 bit unsigned integer in current endian
func (d *D) TryU63() (uint64, error) { return d.tryUEndian(63, d.Endian) }

// U63 reads 63 bit unsigned integer in current endian
func (d *D) U63() uint64 {
	v, err := d.tryUEndian(63, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U63", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU63 tries to add a field and read 63 bit unsigned integer in current endian
func (d *D) TryFieldScalarU63(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(63, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU63 adds a field and reads 63 bit unsigned integer in current endian
func (d *D) FieldScalarU63(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU63(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U63", Pos: d.Pos()})
	}
	return s
}

// TryFieldU63 tries to add a field and read 63 bit unsigned integer in current endian
func (d *D) TryFieldU63(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU63(name, sms...)
	return s.ActualU(), err
}

// FieldU63 adds a field and reads 63 bit unsigned integer in current endian
func (d *D) FieldU63(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU63(name, sms...).ActualU()
}

// Reader U64

// TryU64 tries to read 64 bit unsigned integer in current endian
func (d *D) TryU64() (uint64, error) { return d.tryUEndian(64, d.Endian) }

// U64 reads 64 bit unsigned integer in current endian
func (d *D) U64() uint64 {
	v, err := d.tryUEndian(64, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U64", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU64 tries to add a field and read 64 bit unsigned integer in current endian
func (d *D) TryFieldScalarU64(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(64, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU64 adds a field and reads 64 bit unsigned integer in current endian
func (d *D) FieldScalarU64(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU64(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U64", Pos: d.Pos()})
	}
	return s
}

// TryFieldU64 tries to add a field and read 64 bit unsigned integer in current endian
func (d *D) TryFieldU64(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU64(name, sms...)
	return s.ActualU(), err
}

// FieldU64 adds a field and reads 64 bit unsigned integer in current endian
func (d *D) FieldU64(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU64(name, sms...).ActualU()
}

// Reader U8LE

// TryU8LE tries to read 8 bit unsigned integer in little-endian
func (d *D) TryU8LE() (uint64, error) { return d.tryUEndian(8, LittleEndian) }

// U8LE reads 8 bit unsigned integer in little-endian
func (d *D) U8LE() uint64 {
	v, err := d.tryUEndian(8, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U8LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU8LE tries to add a field and read 8 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU8LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(8, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU8LE adds a field and reads 8 bit unsigned integer in little-endian
func (d *D) FieldScalarU8LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU8LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U8LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU8LE tries to add a field and read 8 bit unsigned integer in little-endian
func (d *D) TryFieldU8LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU8LE(name, sms...)
	return s.ActualU(), err
}

// FieldU8LE adds a field and reads 8 bit unsigned integer in little-endian
func (d *D) FieldU8LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU8LE(name, sms...).ActualU()
}

// Reader U9LE

// TryU9LE tries to read 9 bit unsigned integer in little-endian
func (d *D) TryU9LE() (uint64, error) { return d.tryUEndian(9, LittleEndian) }

// U9LE reads 9 bit unsigned integer in little-endian
func (d *D) U9LE() uint64 {
	v, err := d.tryUEndian(9, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U9LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU9LE tries to add a field and read 9 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU9LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(9, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU9LE adds a field and reads 9 bit unsigned integer in little-endian
func (d *D) FieldScalarU9LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU9LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U9LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU9LE tries to add a field and read 9 bit unsigned integer in little-endian
func (d *D) TryFieldU9LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU9LE(name, sms...)
	return s.ActualU(), err
}

// FieldU9LE adds a field and reads 9 bit unsigned integer in little-endian
func (d *D) FieldU9LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU9LE(name, sms...).ActualU()
}

// Reader U10LE

// TryU10LE tries to read 10 bit unsigned integer in little-endian
func (d *D) TryU10LE() (uint64, error) { return d.tryUEndian(10, LittleEndian) }

// U10LE reads 10 bit unsigned integer in little-endian
func (d *D) U10LE() uint64 {
	v, err := d.tryUEndian(10, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U10LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU10LE tries to add a field and read 10 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU10LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(10, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU10LE adds a field and reads 10 bit unsigned integer in little-endian
func (d *D) FieldScalarU10LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU10LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U10LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU10LE tries to add a field and read 10 bit unsigned integer in little-endian
func (d *D) TryFieldU10LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU10LE(name, sms...)
	return s.ActualU(), err
}

// FieldU10LE adds a field and reads 10 bit unsigned integer in little-endian
func (d *D) FieldU10LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU10LE(name, sms...).ActualU()
}

// Reader U11LE

// TryU11LE tries to read 11 bit unsigned integer in little-endian
func (d *D) TryU11LE() (uint64, error) { return d.tryUEndian(11, LittleEndian) }

// U11LE reads 11 bit unsigned integer in little-endian
func (d *D) U11LE() uint64 {
	v, err := d.tryUEndian(11, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U11LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU11LE tries to add a field and read 11 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU11LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(11, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU11LE adds a field and reads 11 bit unsigned integer in little-endian
func (d *D) FieldScalarU11LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU11LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U11LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU11LE tries to add a field and read 11 bit unsigned integer in little-endian
func (d *D) TryFieldU11LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU11LE(name, sms...)
	return s.ActualU(), err
}

// FieldU11LE adds a field and reads 11 bit unsigned integer in little-endian
func (d *D) FieldU11LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU11LE(name, sms...).ActualU()
}

// Reader U12LE

// TryU12LE tries to read 12 bit unsigned integer in little-endian
func (d *D) TryU12LE() (uint64, error) { return d.tryUEndian(12, LittleEndian) }

// U12LE reads 12 bit unsigned integer in little-endian
func (d *D) U12LE() uint64 {
	v, err := d.tryUEndian(12, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U12LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU12LE tries to add a field and read 12 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU12LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(12, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU12LE adds a field and reads 12 bit unsigned integer in little-endian
func (d *D) FieldScalarU12LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU12LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U12LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU12LE tries to add a field and read 12 bit unsigned integer in little-endian
func (d *D) TryFieldU12LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU12LE(name, sms...)
	return s.ActualU(), err
}

// FieldU12LE adds a field and reads 12 bit unsigned integer in little-endian
func (d *D) FieldU12LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU12LE(name, sms...).ActualU()
}

// Reader U13LE

// TryU13LE tries to read 13 bit unsigned integer in little-endian
func (d *D) TryU13LE() (uint64, error) { return d.tryUEndian(13, LittleEndian) }

// U13LE reads 13 bit unsigned integer in little-endian
func (d *D) U13LE() uint64 {
	v, err := d.tryUEndian(13, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U13LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU13LE tries to add a field and read 13 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU13LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(13, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU13LE adds a field and reads 13 bit unsigned integer in little-endian
func (d *D) FieldScalarU13LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU13LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U13LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU13LE tries to add a field and read 13 bit unsigned integer in little-endian
func (d *D) TryFieldU13LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU13LE(name, sms...)
	return s.ActualU(), err
}

// FieldU13LE adds a field and reads 13 bit unsigned integer in little-endian
func (d *D) FieldU13LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU13LE(name, sms...).ActualU()
}

// Reader U14LE

// TryU14LE tries to read 14 bit unsigned integer in little-endian
func (d *D) TryU14LE() (uint64, error) { return d.tryUEndian(14, LittleEndian) }

// U14LE reads 14 bit unsigned integer in little-endian
func (d *D) U14LE() uint64 {
	v, err := d.tryUEndian(14, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U14LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU14LE tries to add a field and read 14 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU14LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(14, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU14LE adds a field and reads 14 bit unsigned integer in little-endian
func (d *D) FieldScalarU14LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU14LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U14LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU14LE tries to add a field and read 14 bit unsigned integer in little-endian
func (d *D) TryFieldU14LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU14LE(name, sms...)
	return s.ActualU(), err
}

// FieldU14LE adds a field and reads 14 bit unsigned integer in little-endian
func (d *D) FieldU14LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU14LE(name, sms...).ActualU()
}

// Reader U15LE

// TryU15LE tries to read 15 bit unsigned integer in little-endian
func (d *D) TryU15LE() (uint64, error) { return d.tryUEndian(15, LittleEndian) }

// U15LE reads 15 bit unsigned integer in little-endian
func (d *D) U15LE() uint64 {
	v, err := d.tryUEndian(15, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U15LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU15LE tries to add a field and read 15 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU15LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(15, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU15LE adds a field and reads 15 bit unsigned integer in little-endian
func (d *D) FieldScalarU15LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU15LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U15LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU15LE tries to add a field and read 15 bit unsigned integer in little-endian
func (d *D) TryFieldU15LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU15LE(name, sms...)
	return s.ActualU(), err
}

// FieldU15LE adds a field and reads 15 bit unsigned integer in little-endian
func (d *D) FieldU15LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU15LE(name, sms...).ActualU()
}

// Reader U16LE

// TryU16LE tries to read 16 bit unsigned integer in little-endian
func (d *D) TryU16LE() (uint64, error) { return d.tryUEndian(16, LittleEndian) }

// U16LE reads 16 bit unsigned integer in little-endian
func (d *D) U16LE() uint64 {
	v, err := d.tryUEndian(16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U16LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU16LE tries to add a field and read 16 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU16LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(16, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU16LE adds a field and reads 16 bit unsigned integer in little-endian
func (d *D) FieldScalarU16LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU16LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U16LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU16LE tries to add a field and read 16 bit unsigned integer in little-endian
func (d *D) TryFieldU16LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU16LE(name, sms...)
	return s.ActualU(), err
}

// FieldU16LE adds a field and reads 16 bit unsigned integer in little-endian
func (d *D) FieldU16LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU16LE(name, sms...).ActualU()
}

// Reader U17LE

// TryU17LE tries to read 17 bit unsigned integer in little-endian
func (d *D) TryU17LE() (uint64, error) { return d.tryUEndian(17, LittleEndian) }

// U17LE reads 17 bit unsigned integer in little-endian
func (d *D) U17LE() uint64 {
	v, err := d.tryUEndian(17, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U17LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU17LE tries to add a field and read 17 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU17LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(17, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU17LE adds a field and reads 17 bit unsigned integer in little-endian
func (d *D) FieldScalarU17LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU17LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U17LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU17LE tries to add a field and read 17 bit unsigned integer in little-endian
func (d *D) TryFieldU17LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU17LE(name, sms...)
	return s.ActualU(), err
}

// FieldU17LE adds a field and reads 17 bit unsigned integer in little-endian
func (d *D) FieldU17LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU17LE(name, sms...).ActualU()
}

// Reader U18LE

// TryU18LE tries to read 18 bit unsigned integer in little-endian
func (d *D) TryU18LE() (uint64, error) { return d.tryUEndian(18, LittleEndian) }

// U18LE reads 18 bit unsigned integer in little-endian
func (d *D) U18LE() uint64 {
	v, err := d.tryUEndian(18, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U18LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU18LE tries to add a field and read 18 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU18LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(18, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU18LE adds a field and reads 18 bit unsigned integer in little-endian
func (d *D) FieldScalarU18LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU18LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U18LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU18LE tries to add a field and read 18 bit unsigned integer in little-endian
func (d *D) TryFieldU18LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU18LE(name, sms...)
	return s.ActualU(), err
}

// FieldU18LE adds a field and reads 18 bit unsigned integer in little-endian
func (d *D) FieldU18LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU18LE(name, sms...).ActualU()
}

// Reader U19LE

// TryU19LE tries to read 19 bit unsigned integer in little-endian
func (d *D) TryU19LE() (uint64, error) { return d.tryUEndian(19, LittleEndian) }

// U19LE reads 19 bit unsigned integer in little-endian
func (d *D) U19LE() uint64 {
	v, err := d.tryUEndian(19, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U19LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU19LE tries to add a field and read 19 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU19LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(19, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU19LE adds a field and reads 19 bit unsigned integer in little-endian
func (d *D) FieldScalarU19LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU19LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U19LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU19LE tries to add a field and read 19 bit unsigned integer in little-endian
func (d *D) TryFieldU19LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU19LE(name, sms...)
	return s.ActualU(), err
}

// FieldU19LE adds a field and reads 19 bit unsigned integer in little-endian
func (d *D) FieldU19LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU19LE(name, sms...).ActualU()
}

// Reader U20LE

// TryU20LE tries to read 20 bit unsigned integer in little-endian
func (d *D) TryU20LE() (uint64, error) { return d.tryUEndian(20, LittleEndian) }

// U20LE reads 20 bit unsigned integer in little-endian
func (d *D) U20LE() uint64 {
	v, err := d.tryUEndian(20, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U20LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU20LE tries to add a field and read 20 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU20LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(20, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU20LE adds a field and reads 20 bit unsigned integer in little-endian
func (d *D) FieldScalarU20LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU20LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U20LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU20LE tries to add a field and read 20 bit unsigned integer in little-endian
func (d *D) TryFieldU20LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU20LE(name, sms...)
	return s.ActualU(), err
}

// FieldU20LE adds a field and reads 20 bit unsigned integer in little-endian
func (d *D) FieldU20LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU20LE(name, sms...).ActualU()
}

// Reader U21LE

// TryU21LE tries to read 21 bit unsigned integer in little-endian
func (d *D) TryU21LE() (uint64, error) { return d.tryUEndian(21, LittleEndian) }

// U21LE reads 21 bit unsigned integer in little-endian
func (d *D) U21LE() uint64 {
	v, err := d.tryUEndian(21, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U21LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU21LE tries to add a field and read 21 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU21LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(21, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU21LE adds a field and reads 21 bit unsigned integer in little-endian
func (d *D) FieldScalarU21LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU21LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U21LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU21LE tries to add a field and read 21 bit unsigned integer in little-endian
func (d *D) TryFieldU21LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU21LE(name, sms...)
	return s.ActualU(), err
}

// FieldU21LE adds a field and reads 21 bit unsigned integer in little-endian
func (d *D) FieldU21LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU21LE(name, sms...).ActualU()
}

// Reader U22LE

// TryU22LE tries to read 22 bit unsigned integer in little-endian
func (d *D) TryU22LE() (uint64, error) { return d.tryUEndian(22, LittleEndian) }

// U22LE reads 22 bit unsigned integer in little-endian
func (d *D) U22LE() uint64 {
	v, err := d.tryUEndian(22, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U22LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU22LE tries to add a field and read 22 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU22LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(22, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU22LE adds a field and reads 22 bit unsigned integer in little-endian
func (d *D) FieldScalarU22LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU22LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U22LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU22LE tries to add a field and read 22 bit unsigned integer in little-endian
func (d *D) TryFieldU22LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU22LE(name, sms...)
	return s.ActualU(), err
}

// FieldU22LE adds a field and reads 22 bit unsigned integer in little-endian
func (d *D) FieldU22LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU22LE(name, sms...).ActualU()
}

// Reader U23LE

// TryU23LE tries to read 23 bit unsigned integer in little-endian
func (d *D) TryU23LE() (uint64, error) { return d.tryUEndian(23, LittleEndian) }

// U23LE reads 23 bit unsigned integer in little-endian
func (d *D) U23LE() uint64 {
	v, err := d.tryUEndian(23, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U23LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU23LE tries to add a field and read 23 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU23LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(23, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU23LE adds a field and reads 23 bit unsigned integer in little-endian
func (d *D) FieldScalarU23LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU23LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U23LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU23LE tries to add a field and read 23 bit unsigned integer in little-endian
func (d *D) TryFieldU23LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU23LE(name, sms...)
	return s.ActualU(), err
}

// FieldU23LE adds a field and reads 23 bit unsigned integer in little-endian
func (d *D) FieldU23LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU23LE(name, sms...).ActualU()
}

// Reader U24LE

// TryU24LE tries to read 24 bit unsigned integer in little-endian
func (d *D) TryU24LE() (uint64, error) { return d.tryUEndian(24, LittleEndian) }

// U24LE reads 24 bit unsigned integer in little-endian
func (d *D) U24LE() uint64 {
	v, err := d.tryUEndian(24, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U24LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU24LE tries to add a field and read 24 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU24LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(24, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU24LE adds a field and reads 24 bit unsigned integer in little-endian
func (d *D) FieldScalarU24LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU24LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U24LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU24LE tries to add a field and read 24 bit unsigned integer in little-endian
func (d *D) TryFieldU24LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU24LE(name, sms...)
	return s.ActualU(), err
}

// FieldU24LE adds a field and reads 24 bit unsigned integer in little-endian
func (d *D) FieldU24LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU24LE(name, sms...).ActualU()
}

// Reader U25LE

// TryU25LE tries to read 25 bit unsigned integer in little-endian
func (d *D) TryU25LE() (uint64, error) { return d.tryUEndian(25, LittleEndian) }

// U25LE reads 25 bit unsigned integer in little-endian
func (d *D) U25LE() uint64 {
	v, err := d.tryUEndian(25, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U25LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU25LE tries to add a field and read 25 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU25LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(25, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU25LE adds a field and reads 25 bit unsigned integer in little-endian
func (d *D) FieldScalarU25LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU25LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U25LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU25LE tries to add a field and read 25 bit unsigned integer in little-endian
func (d *D) TryFieldU25LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU25LE(name, sms...)
	return s.ActualU(), err
}

// FieldU25LE adds a field and reads 25 bit unsigned integer in little-endian
func (d *D) FieldU25LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU25LE(name, sms...).ActualU()
}

// Reader U26LE

// TryU26LE tries to read 26 bit unsigned integer in little-endian
func (d *D) TryU26LE() (uint64, error) { return d.tryUEndian(26, LittleEndian) }

// U26LE reads 26 bit unsigned integer in little-endian
func (d *D) U26LE() uint64 {
	v, err := d.tryUEndian(26, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U26LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU26LE tries to add a field and read 26 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU26LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(26, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU26LE adds a field and reads 26 bit unsigned integer in little-endian
func (d *D) FieldScalarU26LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU26LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U26LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU26LE tries to add a field and read 26 bit unsigned integer in little-endian
func (d *D) TryFieldU26LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU26LE(name, sms...)
	return s.ActualU(), err
}

// FieldU26LE adds a field and reads 26 bit unsigned integer in little-endian
func (d *D) FieldU26LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU26LE(name, sms...).ActualU()
}

// Reader U27LE

// TryU27LE tries to read 27 bit unsigned integer in little-endian
func (d *D) TryU27LE() (uint64, error) { return d.tryUEndian(27, LittleEndian) }

// U27LE reads 27 bit unsigned integer in little-endian
func (d *D) U27LE() uint64 {
	v, err := d.tryUEndian(27, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U27LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU27LE tries to add a field and read 27 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU27LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(27, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU27LE adds a field and reads 27 bit unsigned integer in little-endian
func (d *D) FieldScalarU27LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU27LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U27LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU27LE tries to add a field and read 27 bit unsigned integer in little-endian
func (d *D) TryFieldU27LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU27LE(name, sms...)
	return s.ActualU(), err
}

// FieldU27LE adds a field and reads 27 bit unsigned integer in little-endian
func (d *D) FieldU27LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU27LE(name, sms...).ActualU()
}

// Reader U28LE

// TryU28LE tries to read 28 bit unsigned integer in little-endian
func (d *D) TryU28LE() (uint64, error) { return d.tryUEndian(28, LittleEndian) }

// U28LE reads 28 bit unsigned integer in little-endian
func (d *D) U28LE() uint64 {
	v, err := d.tryUEndian(28, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U28LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU28LE tries to add a field and read 28 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU28LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(28, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU28LE adds a field and reads 28 bit unsigned integer in little-endian
func (d *D) FieldScalarU28LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU28LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U28LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU28LE tries to add a field and read 28 bit unsigned integer in little-endian
func (d *D) TryFieldU28LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU28LE(name, sms...)
	return s.ActualU(), err
}

// FieldU28LE adds a field and reads 28 bit unsigned integer in little-endian
func (d *D) FieldU28LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU28LE(name, sms...).ActualU()
}

// Reader U29LE

// TryU29LE tries to read 29 bit unsigned integer in little-endian
func (d *D) TryU29LE() (uint64, error) { return d.tryUEndian(29, LittleEndian) }

// U29LE reads 29 bit unsigned integer in little-endian
func (d *D) U29LE() uint64 {
	v, err := d.tryUEndian(29, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U29LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU29LE tries to add a field and read 29 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU29LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(29, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU29LE adds a field and reads 29 bit unsigned integer in little-endian
func (d *D) FieldScalarU29LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU29LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U29LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU29LE tries to add a field and read 29 bit unsigned integer in little-endian
func (d *D) TryFieldU29LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU29LE(name, sms...)
	return s.ActualU(), err
}

// FieldU29LE adds a field and reads 29 bit unsigned integer in little-endian
func (d *D) FieldU29LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU29LE(name, sms...).ActualU()
}

// Reader U30LE

// TryU30LE tries to read 30 bit unsigned integer in little-endian
func (d *D) TryU30LE() (uint64, error) { return d.tryUEndian(30, LittleEndian) }

// U30LE reads 30 bit unsigned integer in little-endian
func (d *D) U30LE() uint64 {
	v, err := d.tryUEndian(30, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U30LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU30LE tries to add a field and read 30 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU30LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(30, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU30LE adds a field and reads 30 bit unsigned integer in little-endian
func (d *D) FieldScalarU30LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU30LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U30LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU30LE tries to add a field and read 30 bit unsigned integer in little-endian
func (d *D) TryFieldU30LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU30LE(name, sms...)
	return s.ActualU(), err
}

// FieldU30LE adds a field and reads 30 bit unsigned integer in little-endian
func (d *D) FieldU30LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU30LE(name, sms...).ActualU()
}

// Reader U31LE

// TryU31LE tries to read 31 bit unsigned integer in little-endian
func (d *D) TryU31LE() (uint64, error) { return d.tryUEndian(31, LittleEndian) }

// U31LE reads 31 bit unsigned integer in little-endian
func (d *D) U31LE() uint64 {
	v, err := d.tryUEndian(31, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U31LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU31LE tries to add a field and read 31 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU31LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(31, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU31LE adds a field and reads 31 bit unsigned integer in little-endian
func (d *D) FieldScalarU31LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU31LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U31LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU31LE tries to add a field and read 31 bit unsigned integer in little-endian
func (d *D) TryFieldU31LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU31LE(name, sms...)
	return s.ActualU(), err
}

// FieldU31LE adds a field and reads 31 bit unsigned integer in little-endian
func (d *D) FieldU31LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU31LE(name, sms...).ActualU()
}

// Reader U32LE

// TryU32LE tries to read 32 bit unsigned integer in little-endian
func (d *D) TryU32LE() (uint64, error) { return d.tryUEndian(32, LittleEndian) }

// U32LE reads 32 bit unsigned integer in little-endian
func (d *D) U32LE() uint64 {
	v, err := d.tryUEndian(32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U32LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU32LE tries to add a field and read 32 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU32LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(32, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU32LE adds a field and reads 32 bit unsigned integer in little-endian
func (d *D) FieldScalarU32LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU32LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U32LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU32LE tries to add a field and read 32 bit unsigned integer in little-endian
func (d *D) TryFieldU32LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU32LE(name, sms...)
	return s.ActualU(), err
}

// FieldU32LE adds a field and reads 32 bit unsigned integer in little-endian
func (d *D) FieldU32LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU32LE(name, sms...).ActualU()
}

// Reader U33LE

// TryU33LE tries to read 33 bit unsigned integer in little-endian
func (d *D) TryU33LE() (uint64, error) { return d.tryUEndian(33, LittleEndian) }

// U33LE reads 33 bit unsigned integer in little-endian
func (d *D) U33LE() uint64 {
	v, err := d.tryUEndian(33, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U33LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU33LE tries to add a field and read 33 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU33LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(33, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU33LE adds a field and reads 33 bit unsigned integer in little-endian
func (d *D) FieldScalarU33LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU33LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U33LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU33LE tries to add a field and read 33 bit unsigned integer in little-endian
func (d *D) TryFieldU33LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU33LE(name, sms...)
	return s.ActualU(), err
}

// FieldU33LE adds a field and reads 33 bit unsigned integer in little-endian
func (d *D) FieldU33LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU33LE(name, sms...).ActualU()
}

// Reader U34LE

// TryU34LE tries to read 34 bit unsigned integer in little-endian
func (d *D) TryU34LE() (uint64, error) { return d.tryUEndian(34, LittleEndian) }

// U34LE reads 34 bit unsigned integer in little-endian
func (d *D) U34LE() uint64 {
	v, err := d.tryUEndian(34, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U34LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU34LE tries to add a field and read 34 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU34LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(34, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU34LE adds a field and reads 34 bit unsigned integer in little-endian
func (d *D) FieldScalarU34LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU34LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U34LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU34LE tries to add a field and read 34 bit unsigned integer in little-endian
func (d *D) TryFieldU34LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU34LE(name, sms...)
	return s.ActualU(), err
}

// FieldU34LE adds a field and reads 34 bit unsigned integer in little-endian
func (d *D) FieldU34LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU34LE(name, sms...).ActualU()
}

// Reader U35LE

// TryU35LE tries to read 35 bit unsigned integer in little-endian
func (d *D) TryU35LE() (uint64, error) { return d.tryUEndian(35, LittleEndian) }

// U35LE reads 35 bit unsigned integer in little-endian
func (d *D) U35LE() uint64 {
	v, err := d.tryUEndian(35, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U35LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU35LE tries to add a field and read 35 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU35LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(35, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU35LE adds a field and reads 35 bit unsigned integer in little-endian
func (d *D) FieldScalarU35LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU35LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U35LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU35LE tries to add a field and read 35 bit unsigned integer in little-endian
func (d *D) TryFieldU35LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU35LE(name, sms...)
	return s.ActualU(), err
}

// FieldU35LE adds a field and reads 35 bit unsigned integer in little-endian
func (d *D) FieldU35LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU35LE(name, sms...).ActualU()
}

// Reader U36LE

// TryU36LE tries to read 36 bit unsigned integer in little-endian
func (d *D) TryU36LE() (uint64, error) { return d.tryUEndian(36, LittleEndian) }

// U36LE reads 36 bit unsigned integer in little-endian
func (d *D) U36LE() uint64 {
	v, err := d.tryUEndian(36, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U36LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU36LE tries to add a field and read 36 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU36LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(36, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU36LE adds a field and reads 36 bit unsigned integer in little-endian
func (d *D) FieldScalarU36LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU36LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U36LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU36LE tries to add a field and read 36 bit unsigned integer in little-endian
func (d *D) TryFieldU36LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU36LE(name, sms...)
	return s.ActualU(), err
}

// FieldU36LE adds a field and reads 36 bit unsigned integer in little-endian
func (d *D) FieldU36LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU36LE(name, sms...).ActualU()
}

// Reader U37LE

// TryU37LE tries to read 37 bit unsigned integer in little-endian
func (d *D) TryU37LE() (uint64, error) { return d.tryUEndian(37, LittleEndian) }

// U37LE reads 37 bit unsigned integer in little-endian
func (d *D) U37LE() uint64 {
	v, err := d.tryUEndian(37, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U37LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU37LE tries to add a field and read 37 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU37LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(37, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU37LE adds a field and reads 37 bit unsigned integer in little-endian
func (d *D) FieldScalarU37LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU37LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U37LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU37LE tries to add a field and read 37 bit unsigned integer in little-endian
func (d *D) TryFieldU37LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU37LE(name, sms...)
	return s.ActualU(), err
}

// FieldU37LE adds a field and reads 37 bit unsigned integer in little-endian
func (d *D) FieldU37LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU37LE(name, sms...).ActualU()
}

// Reader U38LE

// TryU38LE tries to read 38 bit unsigned integer in little-endian
func (d *D) TryU38LE() (uint64, error) { return d.tryUEndian(38, LittleEndian) }

// U38LE reads 38 bit unsigned integer in little-endian
func (d *D) U38LE() uint64 {
	v, err := d.tryUEndian(38, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U38LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU38LE tries to add a field and read 38 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU38LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(38, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU38LE adds a field and reads 38 bit unsigned integer in little-endian
func (d *D) FieldScalarU38LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU38LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U38LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU38LE tries to add a field and read 38 bit unsigned integer in little-endian
func (d *D) TryFieldU38LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU38LE(name, sms...)
	return s.ActualU(), err
}

// FieldU38LE adds a field and reads 38 bit unsigned integer in little-endian
func (d *D) FieldU38LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU38LE(name, sms...).ActualU()
}

// Reader U39LE

// TryU39LE tries to read 39 bit unsigned integer in little-endian
func (d *D) TryU39LE() (uint64, error) { return d.tryUEndian(39, LittleEndian) }

// U39LE reads 39 bit unsigned integer in little-endian
func (d *D) U39LE() uint64 {
	v, err := d.tryUEndian(39, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U39LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU39LE tries to add a field and read 39 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU39LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(39, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU39LE adds a field and reads 39 bit unsigned integer in little-endian
func (d *D) FieldScalarU39LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU39LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U39LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU39LE tries to add a field and read 39 bit unsigned integer in little-endian
func (d *D) TryFieldU39LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU39LE(name, sms...)
	return s.ActualU(), err
}

// FieldU39LE adds a field and reads 39 bit unsigned integer in little-endian
func (d *D) FieldU39LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU39LE(name, sms...).ActualU()
}

// Reader U40LE

// TryU40LE tries to read 40 bit unsigned integer in little-endian
func (d *D) TryU40LE() (uint64, error) { return d.tryUEndian(40, LittleEndian) }

// U40LE reads 40 bit unsigned integer in little-endian
func (d *D) U40LE() uint64 {
	v, err := d.tryUEndian(40, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U40LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU40LE tries to add a field and read 40 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU40LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(40, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU40LE adds a field and reads 40 bit unsigned integer in little-endian
func (d *D) FieldScalarU40LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU40LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U40LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU40LE tries to add a field and read 40 bit unsigned integer in little-endian
func (d *D) TryFieldU40LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU40LE(name, sms...)
	return s.ActualU(), err
}

// FieldU40LE adds a field and reads 40 bit unsigned integer in little-endian
func (d *D) FieldU40LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU40LE(name, sms...).ActualU()
}

// Reader U41LE

// TryU41LE tries to read 41 bit unsigned integer in little-endian
func (d *D) TryU41LE() (uint64, error) { return d.tryUEndian(41, LittleEndian) }

// U41LE reads 41 bit unsigned integer in little-endian
func (d *D) U41LE() uint64 {
	v, err := d.tryUEndian(41, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U41LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU41LE tries to add a field and read 41 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU41LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(41, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU41LE adds a field and reads 41 bit unsigned integer in little-endian
func (d *D) FieldScalarU41LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU41LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U41LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU41LE tries to add a field and read 41 bit unsigned integer in little-endian
func (d *D) TryFieldU41LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU41LE(name, sms...)
	return s.ActualU(), err
}

// FieldU41LE adds a field and reads 41 bit unsigned integer in little-endian
func (d *D) FieldU41LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU41LE(name, sms...).ActualU()
}

// Reader U42LE

// TryU42LE tries to read 42 bit unsigned integer in little-endian
func (d *D) TryU42LE() (uint64, error) { return d.tryUEndian(42, LittleEndian) }

// U42LE reads 42 bit unsigned integer in little-endian
func (d *D) U42LE() uint64 {
	v, err := d.tryUEndian(42, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U42LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU42LE tries to add a field and read 42 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU42LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(42, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU42LE adds a field and reads 42 bit unsigned integer in little-endian
func (d *D) FieldScalarU42LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU42LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U42LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU42LE tries to add a field and read 42 bit unsigned integer in little-endian
func (d *D) TryFieldU42LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU42LE(name, sms...)
	return s.ActualU(), err
}

// FieldU42LE adds a field and reads 42 bit unsigned integer in little-endian
func (d *D) FieldU42LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU42LE(name, sms...).ActualU()
}

// Reader U43LE

// TryU43LE tries to read 43 bit unsigned integer in little-endian
func (d *D) TryU43LE() (uint64, error) { return d.tryUEndian(43, LittleEndian) }

// U43LE reads 43 bit unsigned integer in little-endian
func (d *D) U43LE() uint64 {
	v, err := d.tryUEndian(43, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U43LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU43LE tries to add a field and read 43 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU43LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(43, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU43LE adds a field and reads 43 bit unsigned integer in little-endian
func (d *D) FieldScalarU43LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU43LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U43LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU43LE tries to add a field and read 43 bit unsigned integer in little-endian
func (d *D) TryFieldU43LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU43LE(name, sms...)
	return s.ActualU(), err
}

// FieldU43LE adds a field and reads 43 bit unsigned integer in little-endian
func (d *D) FieldU43LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU43LE(name, sms...).ActualU()
}

// Reader U44LE

// TryU44LE tries to read 44 bit unsigned integer in little-endian
func (d *D) TryU44LE() (uint64, error) { return d.tryUEndian(44, LittleEndian) }

// U44LE reads 44 bit unsigned integer in little-endian
func (d *D) U44LE() uint64 {
	v, err := d.tryUEndian(44, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U44LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU44LE tries to add a field and read 44 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU44LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(44, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU44LE adds a field and reads 44 bit unsigned integer in little-endian
func (d *D) FieldScalarU44LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU44LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U44LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU44LE tries to add a field and read 44 bit unsigned integer in little-endian
func (d *D) TryFieldU44LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU44LE(name, sms...)
	return s.ActualU(), err
}

// FieldU44LE adds a field and reads 44 bit unsigned integer in little-endian
func (d *D) FieldU44LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU44LE(name, sms...).ActualU()
}

// Reader U45LE

// TryU45LE tries to read 45 bit unsigned integer in little-endian
func (d *D) TryU45LE() (uint64, error) { return d.tryUEndian(45, LittleEndian) }

// U45LE reads 45 bit unsigned integer in little-endian
func (d *D) U45LE() uint64 {
	v, err := d.tryUEndian(45, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U45LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU45LE tries to add a field and read 45 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU45LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(45, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU45LE adds a field and reads 45 bit unsigned integer in little-endian
func (d *D) FieldScalarU45LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU45LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U45LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU45LE tries to add a field and read 45 bit unsigned integer in little-endian
func (d *D) TryFieldU45LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU45LE(name, sms...)
	return s.ActualU(), err
}

// FieldU45LE adds a field and reads 45 bit unsigned integer in little-endian
func (d *D) FieldU45LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU45LE(name, sms...).ActualU()
}

// Reader U46LE

// TryU46LE tries to read 46 bit unsigned integer in little-endian
func (d *D) TryU46LE() (uint64, error) { return d.tryUEndian(46, LittleEndian) }

// U46LE reads 46 bit unsigned integer in little-endian
func (d *D) U46LE() uint64 {
	v, err := d.tryUEndian(46, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U46LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU46LE tries to add a field and read 46 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU46LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(46, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU46LE adds a field and reads 46 bit unsigned integer in little-endian
func (d *D) FieldScalarU46LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU46LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U46LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU46LE tries to add a field and read 46 bit unsigned integer in little-endian
func (d *D) TryFieldU46LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU46LE(name, sms...)
	return s.ActualU(), err
}

// FieldU46LE adds a field and reads 46 bit unsigned integer in little-endian
func (d *D) FieldU46LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU46LE(name, sms...).ActualU()
}

// Reader U47LE

// TryU47LE tries to read 47 bit unsigned integer in little-endian
func (d *D) TryU47LE() (uint64, error) { return d.tryUEndian(47, LittleEndian) }

// U47LE reads 47 bit unsigned integer in little-endian
func (d *D) U47LE() uint64 {
	v, err := d.tryUEndian(47, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U47LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU47LE tries to add a field and read 47 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU47LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(47, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU47LE adds a field and reads 47 bit unsigned integer in little-endian
func (d *D) FieldScalarU47LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU47LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U47LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU47LE tries to add a field and read 47 bit unsigned integer in little-endian
func (d *D) TryFieldU47LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU47LE(name, sms...)
	return s.ActualU(), err
}

// FieldU47LE adds a field and reads 47 bit unsigned integer in little-endian
func (d *D) FieldU47LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU47LE(name, sms...).ActualU()
}

// Reader U48LE

// TryU48LE tries to read 48 bit unsigned integer in little-endian
func (d *D) TryU48LE() (uint64, error) { return d.tryUEndian(48, LittleEndian) }

// U48LE reads 48 bit unsigned integer in little-endian
func (d *D) U48LE() uint64 {
	v, err := d.tryUEndian(48, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U48LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU48LE tries to add a field and read 48 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU48LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(48, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU48LE adds a field and reads 48 bit unsigned integer in little-endian
func (d *D) FieldScalarU48LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU48LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U48LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU48LE tries to add a field and read 48 bit unsigned integer in little-endian
func (d *D) TryFieldU48LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU48LE(name, sms...)
	return s.ActualU(), err
}

// FieldU48LE adds a field and reads 48 bit unsigned integer in little-endian
func (d *D) FieldU48LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU48LE(name, sms...).ActualU()
}

// Reader U49LE

// TryU49LE tries to read 49 bit unsigned integer in little-endian
func (d *D) TryU49LE() (uint64, error) { return d.tryUEndian(49, LittleEndian) }

// U49LE reads 49 bit unsigned integer in little-endian
func (d *D) U49LE() uint64 {
	v, err := d.tryUEndian(49, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U49LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU49LE tries to add a field and read 49 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU49LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(49, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU49LE adds a field and reads 49 bit unsigned integer in little-endian
func (d *D) FieldScalarU49LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU49LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U49LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU49LE tries to add a field and read 49 bit unsigned integer in little-endian
func (d *D) TryFieldU49LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU49LE(name, sms...)
	return s.ActualU(), err
}

// FieldU49LE adds a field and reads 49 bit unsigned integer in little-endian
func (d *D) FieldU49LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU49LE(name, sms...).ActualU()
}

// Reader U50LE

// TryU50LE tries to read 50 bit unsigned integer in little-endian
func (d *D) TryU50LE() (uint64, error) { return d.tryUEndian(50, LittleEndian) }

// U50LE reads 50 bit unsigned integer in little-endian
func (d *D) U50LE() uint64 {
	v, err := d.tryUEndian(50, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U50LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU50LE tries to add a field and read 50 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU50LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(50, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU50LE adds a field and reads 50 bit unsigned integer in little-endian
func (d *D) FieldScalarU50LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU50LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U50LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU50LE tries to add a field and read 50 bit unsigned integer in little-endian
func (d *D) TryFieldU50LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU50LE(name, sms...)
	return s.ActualU(), err
}

// FieldU50LE adds a field and reads 50 bit unsigned integer in little-endian
func (d *D) FieldU50LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU50LE(name, sms...).ActualU()
}

// Reader U51LE

// TryU51LE tries to read 51 bit unsigned integer in little-endian
func (d *D) TryU51LE() (uint64, error) { return d.tryUEndian(51, LittleEndian) }

// U51LE reads 51 bit unsigned integer in little-endian
func (d *D) U51LE() uint64 {
	v, err := d.tryUEndian(51, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U51LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU51LE tries to add a field and read 51 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU51LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(51, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU51LE adds a field and reads 51 bit unsigned integer in little-endian
func (d *D) FieldScalarU51LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU51LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U51LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU51LE tries to add a field and read 51 bit unsigned integer in little-endian
func (d *D) TryFieldU51LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU51LE(name, sms...)
	return s.ActualU(), err
}

// FieldU51LE adds a field and reads 51 bit unsigned integer in little-endian
func (d *D) FieldU51LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU51LE(name, sms...).ActualU()
}

// Reader U52LE

// TryU52LE tries to read 52 bit unsigned integer in little-endian
func (d *D) TryU52LE() (uint64, error) { return d.tryUEndian(52, LittleEndian) }

// U52LE reads 52 bit unsigned integer in little-endian
func (d *D) U52LE() uint64 {
	v, err := d.tryUEndian(52, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U52LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU52LE tries to add a field and read 52 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU52LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(52, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU52LE adds a field and reads 52 bit unsigned integer in little-endian
func (d *D) FieldScalarU52LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU52LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U52LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU52LE tries to add a field and read 52 bit unsigned integer in little-endian
func (d *D) TryFieldU52LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU52LE(name, sms...)
	return s.ActualU(), err
}

// FieldU52LE adds a field and reads 52 bit unsigned integer in little-endian
func (d *D) FieldU52LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU52LE(name, sms...).ActualU()
}

// Reader U53LE

// TryU53LE tries to read 53 bit unsigned integer in little-endian
func (d *D) TryU53LE() (uint64, error) { return d.tryUEndian(53, LittleEndian) }

// U53LE reads 53 bit unsigned integer in little-endian
func (d *D) U53LE() uint64 {
	v, err := d.tryUEndian(53, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U53LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU53LE tries to add a field and read 53 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU53LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(53, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU53LE adds a field and reads 53 bit unsigned integer in little-endian
func (d *D) FieldScalarU53LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU53LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U53LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU53LE tries to add a field and read 53 bit unsigned integer in little-endian
func (d *D) TryFieldU53LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU53LE(name, sms...)
	return s.ActualU(), err
}

// FieldU53LE adds a field and reads 53 bit unsigned integer in little-endian
func (d *D) FieldU53LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU53LE(name, sms...).ActualU()
}

// Reader U54LE

// TryU54LE tries to read 54 bit unsigned integer in little-endian
func (d *D) TryU54LE() (uint64, error) { return d.tryUEndian(54, LittleEndian) }

// U54LE reads 54 bit unsigned integer in little-endian
func (d *D) U54LE() uint64 {
	v, err := d.tryUEndian(54, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U54LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU54LE tries to add a field and read 54 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU54LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(54, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU54LE adds a field and reads 54 bit unsigned integer in little-endian
func (d *D) FieldScalarU54LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU54LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U54LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU54LE tries to add a field and read 54 bit unsigned integer in little-endian
func (d *D) TryFieldU54LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU54LE(name, sms...)
	return s.ActualU(), err
}

// FieldU54LE adds a field and reads 54 bit unsigned integer in little-endian
func (d *D) FieldU54LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU54LE(name, sms...).ActualU()
}

// Reader U55LE

// TryU55LE tries to read 55 bit unsigned integer in little-endian
func (d *D) TryU55LE() (uint64, error) { return d.tryUEndian(55, LittleEndian) }

// U55LE reads 55 bit unsigned integer in little-endian
func (d *D) U55LE() uint64 {
	v, err := d.tryUEndian(55, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U55LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU55LE tries to add a field and read 55 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU55LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(55, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU55LE adds a field and reads 55 bit unsigned integer in little-endian
func (d *D) FieldScalarU55LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU55LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U55LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU55LE tries to add a field and read 55 bit unsigned integer in little-endian
func (d *D) TryFieldU55LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU55LE(name, sms...)
	return s.ActualU(), err
}

// FieldU55LE adds a field and reads 55 bit unsigned integer in little-endian
func (d *D) FieldU55LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU55LE(name, sms...).ActualU()
}

// Reader U56LE

// TryU56LE tries to read 56 bit unsigned integer in little-endian
func (d *D) TryU56LE() (uint64, error) { return d.tryUEndian(56, LittleEndian) }

// U56LE reads 56 bit unsigned integer in little-endian
func (d *D) U56LE() uint64 {
	v, err := d.tryUEndian(56, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U56LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU56LE tries to add a field and read 56 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU56LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(56, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU56LE adds a field and reads 56 bit unsigned integer in little-endian
func (d *D) FieldScalarU56LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU56LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U56LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU56LE tries to add a field and read 56 bit unsigned integer in little-endian
func (d *D) TryFieldU56LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU56LE(name, sms...)
	return s.ActualU(), err
}

// FieldU56LE adds a field and reads 56 bit unsigned integer in little-endian
func (d *D) FieldU56LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU56LE(name, sms...).ActualU()
}

// Reader U57LE

// TryU57LE tries to read 57 bit unsigned integer in little-endian
func (d *D) TryU57LE() (uint64, error) { return d.tryUEndian(57, LittleEndian) }

// U57LE reads 57 bit unsigned integer in little-endian
func (d *D) U57LE() uint64 {
	v, err := d.tryUEndian(57, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U57LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU57LE tries to add a field and read 57 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU57LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(57, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU57LE adds a field and reads 57 bit unsigned integer in little-endian
func (d *D) FieldScalarU57LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU57LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U57LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU57LE tries to add a field and read 57 bit unsigned integer in little-endian
func (d *D) TryFieldU57LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU57LE(name, sms...)
	return s.ActualU(), err
}

// FieldU57LE adds a field and reads 57 bit unsigned integer in little-endian
func (d *D) FieldU57LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU57LE(name, sms...).ActualU()
}

// Reader U58LE

// TryU58LE tries to read 58 bit unsigned integer in little-endian
func (d *D) TryU58LE() (uint64, error) { return d.tryUEndian(58, LittleEndian) }

// U58LE reads 58 bit unsigned integer in little-endian
func (d *D) U58LE() uint64 {
	v, err := d.tryUEndian(58, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U58LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU58LE tries to add a field and read 58 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU58LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(58, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU58LE adds a field and reads 58 bit unsigned integer in little-endian
func (d *D) FieldScalarU58LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU58LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U58LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU58LE tries to add a field and read 58 bit unsigned integer in little-endian
func (d *D) TryFieldU58LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU58LE(name, sms...)
	return s.ActualU(), err
}

// FieldU58LE adds a field and reads 58 bit unsigned integer in little-endian
func (d *D) FieldU58LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU58LE(name, sms...).ActualU()
}

// Reader U59LE

// TryU59LE tries to read 59 bit unsigned integer in little-endian
func (d *D) TryU59LE() (uint64, error) { return d.tryUEndian(59, LittleEndian) }

// U59LE reads 59 bit unsigned integer in little-endian
func (d *D) U59LE() uint64 {
	v, err := d.tryUEndian(59, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U59LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU59LE tries to add a field and read 59 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU59LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(59, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU59LE adds a field and reads 59 bit unsigned integer in little-endian
func (d *D) FieldScalarU59LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU59LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U59LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU59LE tries to add a field and read 59 bit unsigned integer in little-endian
func (d *D) TryFieldU59LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU59LE(name, sms...)
	return s.ActualU(), err
}

// FieldU59LE adds a field and reads 59 bit unsigned integer in little-endian
func (d *D) FieldU59LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU59LE(name, sms...).ActualU()
}

// Reader U60LE

// TryU60LE tries to read 60 bit unsigned integer in little-endian
func (d *D) TryU60LE() (uint64, error) { return d.tryUEndian(60, LittleEndian) }

// U60LE reads 60 bit unsigned integer in little-endian
func (d *D) U60LE() uint64 {
	v, err := d.tryUEndian(60, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U60LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU60LE tries to add a field and read 60 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU60LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(60, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU60LE adds a field and reads 60 bit unsigned integer in little-endian
func (d *D) FieldScalarU60LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU60LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U60LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU60LE tries to add a field and read 60 bit unsigned integer in little-endian
func (d *D) TryFieldU60LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU60LE(name, sms...)
	return s.ActualU(), err
}

// FieldU60LE adds a field and reads 60 bit unsigned integer in little-endian
func (d *D) FieldU60LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU60LE(name, sms...).ActualU()
}

// Reader U61LE

// TryU61LE tries to read 61 bit unsigned integer in little-endian
func (d *D) TryU61LE() (uint64, error) { return d.tryUEndian(61, LittleEndian) }

// U61LE reads 61 bit unsigned integer in little-endian
func (d *D) U61LE() uint64 {
	v, err := d.tryUEndian(61, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U61LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU61LE tries to add a field and read 61 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU61LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(61, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU61LE adds a field and reads 61 bit unsigned integer in little-endian
func (d *D) FieldScalarU61LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU61LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U61LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU61LE tries to add a field and read 61 bit unsigned integer in little-endian
func (d *D) TryFieldU61LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU61LE(name, sms...)
	return s.ActualU(), err
}

// FieldU61LE adds a field and reads 61 bit unsigned integer in little-endian
func (d *D) FieldU61LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU61LE(name, sms...).ActualU()
}

// Reader U62LE

// TryU62LE tries to read 62 bit unsigned integer in little-endian
func (d *D) TryU62LE() (uint64, error) { return d.tryUEndian(62, LittleEndian) }

// U62LE reads 62 bit unsigned integer in little-endian
func (d *D) U62LE() uint64 {
	v, err := d.tryUEndian(62, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U62LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU62LE tries to add a field and read 62 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU62LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(62, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU62LE adds a field and reads 62 bit unsigned integer in little-endian
func (d *D) FieldScalarU62LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU62LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U62LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU62LE tries to add a field and read 62 bit unsigned integer in little-endian
func (d *D) TryFieldU62LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU62LE(name, sms...)
	return s.ActualU(), err
}

// FieldU62LE adds a field and reads 62 bit unsigned integer in little-endian
func (d *D) FieldU62LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU62LE(name, sms...).ActualU()
}

// Reader U63LE

// TryU63LE tries to read 63 bit unsigned integer in little-endian
func (d *D) TryU63LE() (uint64, error) { return d.tryUEndian(63, LittleEndian) }

// U63LE reads 63 bit unsigned integer in little-endian
func (d *D) U63LE() uint64 {
	v, err := d.tryUEndian(63, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U63LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU63LE tries to add a field and read 63 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU63LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(63, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU63LE adds a field and reads 63 bit unsigned integer in little-endian
func (d *D) FieldScalarU63LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU63LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U63LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU63LE tries to add a field and read 63 bit unsigned integer in little-endian
func (d *D) TryFieldU63LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU63LE(name, sms...)
	return s.ActualU(), err
}

// FieldU63LE adds a field and reads 63 bit unsigned integer in little-endian
func (d *D) FieldU63LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU63LE(name, sms...).ActualU()
}

// Reader U64LE

// TryU64LE tries to read 64 bit unsigned integer in little-endian
func (d *D) TryU64LE() (uint64, error) { return d.tryUEndian(64, LittleEndian) }

// U64LE reads 64 bit unsigned integer in little-endian
func (d *D) U64LE() uint64 {
	v, err := d.tryUEndian(64, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U64LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU64LE tries to add a field and read 64 bit unsigned integer in little-endian
func (d *D) TryFieldScalarU64LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(64, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU64LE adds a field and reads 64 bit unsigned integer in little-endian
func (d *D) FieldScalarU64LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU64LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U64LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU64LE tries to add a field and read 64 bit unsigned integer in little-endian
func (d *D) TryFieldU64LE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU64LE(name, sms...)
	return s.ActualU(), err
}

// FieldU64LE adds a field and reads 64 bit unsigned integer in little-endian
func (d *D) FieldU64LE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU64LE(name, sms...).ActualU()
}

// Reader U8BE

// TryU8BE tries to read 8 bit unsigned integer in big-endian
func (d *D) TryU8BE() (uint64, error) { return d.tryUEndian(8, BigEndian) }

// U8BE reads 8 bit unsigned integer in big-endian
func (d *D) U8BE() uint64 {
	v, err := d.tryUEndian(8, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U8BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU8BE tries to add a field and read 8 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU8BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(8, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU8BE adds a field and reads 8 bit unsigned integer in big-endian
func (d *D) FieldScalarU8BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU8BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U8BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU8BE tries to add a field and read 8 bit unsigned integer in big-endian
func (d *D) TryFieldU8BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU8BE(name, sms...)
	return s.ActualU(), err
}

// FieldU8BE adds a field and reads 8 bit unsigned integer in big-endian
func (d *D) FieldU8BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU8BE(name, sms...).ActualU()
}

// Reader U9BE

// TryU9BE tries to read 9 bit unsigned integer in big-endian
func (d *D) TryU9BE() (uint64, error) { return d.tryUEndian(9, BigEndian) }

// U9BE reads 9 bit unsigned integer in big-endian
func (d *D) U9BE() uint64 {
	v, err := d.tryUEndian(9, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U9BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU9BE tries to add a field and read 9 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU9BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(9, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU9BE adds a field and reads 9 bit unsigned integer in big-endian
func (d *D) FieldScalarU9BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU9BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U9BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU9BE tries to add a field and read 9 bit unsigned integer in big-endian
func (d *D) TryFieldU9BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU9BE(name, sms...)
	return s.ActualU(), err
}

// FieldU9BE adds a field and reads 9 bit unsigned integer in big-endian
func (d *D) FieldU9BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU9BE(name, sms...).ActualU()
}

// Reader U10BE

// TryU10BE tries to read 10 bit unsigned integer in big-endian
func (d *D) TryU10BE() (uint64, error) { return d.tryUEndian(10, BigEndian) }

// U10BE reads 10 bit unsigned integer in big-endian
func (d *D) U10BE() uint64 {
	v, err := d.tryUEndian(10, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U10BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU10BE tries to add a field and read 10 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU10BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(10, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU10BE adds a field and reads 10 bit unsigned integer in big-endian
func (d *D) FieldScalarU10BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU10BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U10BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU10BE tries to add a field and read 10 bit unsigned integer in big-endian
func (d *D) TryFieldU10BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU10BE(name, sms...)
	return s.ActualU(), err
}

// FieldU10BE adds a field and reads 10 bit unsigned integer in big-endian
func (d *D) FieldU10BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU10BE(name, sms...).ActualU()
}

// Reader U11BE

// TryU11BE tries to read 11 bit unsigned integer in big-endian
func (d *D) TryU11BE() (uint64, error) { return d.tryUEndian(11, BigEndian) }

// U11BE reads 11 bit unsigned integer in big-endian
func (d *D) U11BE() uint64 {
	v, err := d.tryUEndian(11, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U11BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU11BE tries to add a field and read 11 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU11BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(11, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU11BE adds a field and reads 11 bit unsigned integer in big-endian
func (d *D) FieldScalarU11BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU11BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U11BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU11BE tries to add a field and read 11 bit unsigned integer in big-endian
func (d *D) TryFieldU11BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU11BE(name, sms...)
	return s.ActualU(), err
}

// FieldU11BE adds a field and reads 11 bit unsigned integer in big-endian
func (d *D) FieldU11BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU11BE(name, sms...).ActualU()
}

// Reader U12BE

// TryU12BE tries to read 12 bit unsigned integer in big-endian
func (d *D) TryU12BE() (uint64, error) { return d.tryUEndian(12, BigEndian) }

// U12BE reads 12 bit unsigned integer in big-endian
func (d *D) U12BE() uint64 {
	v, err := d.tryUEndian(12, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U12BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU12BE tries to add a field and read 12 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU12BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(12, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU12BE adds a field and reads 12 bit unsigned integer in big-endian
func (d *D) FieldScalarU12BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU12BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U12BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU12BE tries to add a field and read 12 bit unsigned integer in big-endian
func (d *D) TryFieldU12BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU12BE(name, sms...)
	return s.ActualU(), err
}

// FieldU12BE adds a field and reads 12 bit unsigned integer in big-endian
func (d *D) FieldU12BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU12BE(name, sms...).ActualU()
}

// Reader U13BE

// TryU13BE tries to read 13 bit unsigned integer in big-endian
func (d *D) TryU13BE() (uint64, error) { return d.tryUEndian(13, BigEndian) }

// U13BE reads 13 bit unsigned integer in big-endian
func (d *D) U13BE() uint64 {
	v, err := d.tryUEndian(13, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U13BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU13BE tries to add a field and read 13 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU13BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(13, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU13BE adds a field and reads 13 bit unsigned integer in big-endian
func (d *D) FieldScalarU13BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU13BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U13BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU13BE tries to add a field and read 13 bit unsigned integer in big-endian
func (d *D) TryFieldU13BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU13BE(name, sms...)
	return s.ActualU(), err
}

// FieldU13BE adds a field and reads 13 bit unsigned integer in big-endian
func (d *D) FieldU13BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU13BE(name, sms...).ActualU()
}

// Reader U14BE

// TryU14BE tries to read 14 bit unsigned integer in big-endian
func (d *D) TryU14BE() (uint64, error) { return d.tryUEndian(14, BigEndian) }

// U14BE reads 14 bit unsigned integer in big-endian
func (d *D) U14BE() uint64 {
	v, err := d.tryUEndian(14, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U14BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU14BE tries to add a field and read 14 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU14BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(14, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU14BE adds a field and reads 14 bit unsigned integer in big-endian
func (d *D) FieldScalarU14BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU14BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U14BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU14BE tries to add a field and read 14 bit unsigned integer in big-endian
func (d *D) TryFieldU14BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU14BE(name, sms...)
	return s.ActualU(), err
}

// FieldU14BE adds a field and reads 14 bit unsigned integer in big-endian
func (d *D) FieldU14BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU14BE(name, sms...).ActualU()
}

// Reader U15BE

// TryU15BE tries to read 15 bit unsigned integer in big-endian
func (d *D) TryU15BE() (uint64, error) { return d.tryUEndian(15, BigEndian) }

// U15BE reads 15 bit unsigned integer in big-endian
func (d *D) U15BE() uint64 {
	v, err := d.tryUEndian(15, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U15BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU15BE tries to add a field and read 15 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU15BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(15, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU15BE adds a field and reads 15 bit unsigned integer in big-endian
func (d *D) FieldScalarU15BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU15BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U15BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU15BE tries to add a field and read 15 bit unsigned integer in big-endian
func (d *D) TryFieldU15BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU15BE(name, sms...)
	return s.ActualU(), err
}

// FieldU15BE adds a field and reads 15 bit unsigned integer in big-endian
func (d *D) FieldU15BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU15BE(name, sms...).ActualU()
}

// Reader U16BE

// TryU16BE tries to read 16 bit unsigned integer in big-endian
func (d *D) TryU16BE() (uint64, error) { return d.tryUEndian(16, BigEndian) }

// U16BE reads 16 bit unsigned integer in big-endian
func (d *D) U16BE() uint64 {
	v, err := d.tryUEndian(16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U16BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU16BE tries to add a field and read 16 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU16BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(16, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU16BE adds a field and reads 16 bit unsigned integer in big-endian
func (d *D) FieldScalarU16BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU16BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U16BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU16BE tries to add a field and read 16 bit unsigned integer in big-endian
func (d *D) TryFieldU16BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU16BE(name, sms...)
	return s.ActualU(), err
}

// FieldU16BE adds a field and reads 16 bit unsigned integer in big-endian
func (d *D) FieldU16BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU16BE(name, sms...).ActualU()
}

// Reader U17BE

// TryU17BE tries to read 17 bit unsigned integer in big-endian
func (d *D) TryU17BE() (uint64, error) { return d.tryUEndian(17, BigEndian) }

// U17BE reads 17 bit unsigned integer in big-endian
func (d *D) U17BE() uint64 {
	v, err := d.tryUEndian(17, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U17BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU17BE tries to add a field and read 17 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU17BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(17, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU17BE adds a field and reads 17 bit unsigned integer in big-endian
func (d *D) FieldScalarU17BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU17BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U17BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU17BE tries to add a field and read 17 bit unsigned integer in big-endian
func (d *D) TryFieldU17BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU17BE(name, sms...)
	return s.ActualU(), err
}

// FieldU17BE adds a field and reads 17 bit unsigned integer in big-endian
func (d *D) FieldU17BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU17BE(name, sms...).ActualU()
}

// Reader U18BE

// TryU18BE tries to read 18 bit unsigned integer in big-endian
func (d *D) TryU18BE() (uint64, error) { return d.tryUEndian(18, BigEndian) }

// U18BE reads 18 bit unsigned integer in big-endian
func (d *D) U18BE() uint64 {
	v, err := d.tryUEndian(18, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U18BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU18BE tries to add a field and read 18 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU18BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(18, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU18BE adds a field and reads 18 bit unsigned integer in big-endian
func (d *D) FieldScalarU18BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU18BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U18BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU18BE tries to add a field and read 18 bit unsigned integer in big-endian
func (d *D) TryFieldU18BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU18BE(name, sms...)
	return s.ActualU(), err
}

// FieldU18BE adds a field and reads 18 bit unsigned integer in big-endian
func (d *D) FieldU18BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU18BE(name, sms...).ActualU()
}

// Reader U19BE

// TryU19BE tries to read 19 bit unsigned integer in big-endian
func (d *D) TryU19BE() (uint64, error) { return d.tryUEndian(19, BigEndian) }

// U19BE reads 19 bit unsigned integer in big-endian
func (d *D) U19BE() uint64 {
	v, err := d.tryUEndian(19, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U19BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU19BE tries to add a field and read 19 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU19BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(19, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU19BE adds a field and reads 19 bit unsigned integer in big-endian
func (d *D) FieldScalarU19BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU19BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U19BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU19BE tries to add a field and read 19 bit unsigned integer in big-endian
func (d *D) TryFieldU19BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU19BE(name, sms...)
	return s.ActualU(), err
}

// FieldU19BE adds a field and reads 19 bit unsigned integer in big-endian
func (d *D) FieldU19BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU19BE(name, sms...).ActualU()
}

// Reader U20BE

// TryU20BE tries to read 20 bit unsigned integer in big-endian
func (d *D) TryU20BE() (uint64, error) { return d.tryUEndian(20, BigEndian) }

// U20BE reads 20 bit unsigned integer in big-endian
func (d *D) U20BE() uint64 {
	v, err := d.tryUEndian(20, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U20BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU20BE tries to add a field and read 20 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU20BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(20, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU20BE adds a field and reads 20 bit unsigned integer in big-endian
func (d *D) FieldScalarU20BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU20BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U20BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU20BE tries to add a field and read 20 bit unsigned integer in big-endian
func (d *D) TryFieldU20BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU20BE(name, sms...)
	return s.ActualU(), err
}

// FieldU20BE adds a field and reads 20 bit unsigned integer in big-endian
func (d *D) FieldU20BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU20BE(name, sms...).ActualU()
}

// Reader U21BE

// TryU21BE tries to read 21 bit unsigned integer in big-endian
func (d *D) TryU21BE() (uint64, error) { return d.tryUEndian(21, BigEndian) }

// U21BE reads 21 bit unsigned integer in big-endian
func (d *D) U21BE() uint64 {
	v, err := d.tryUEndian(21, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U21BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU21BE tries to add a field and read 21 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU21BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(21, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU21BE adds a field and reads 21 bit unsigned integer in big-endian
func (d *D) FieldScalarU21BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU21BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U21BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU21BE tries to add a field and read 21 bit unsigned integer in big-endian
func (d *D) TryFieldU21BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU21BE(name, sms...)
	return s.ActualU(), err
}

// FieldU21BE adds a field and reads 21 bit unsigned integer in big-endian
func (d *D) FieldU21BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU21BE(name, sms...).ActualU()
}

// Reader U22BE

// TryU22BE tries to read 22 bit unsigned integer in big-endian
func (d *D) TryU22BE() (uint64, error) { return d.tryUEndian(22, BigEndian) }

// U22BE reads 22 bit unsigned integer in big-endian
func (d *D) U22BE() uint64 {
	v, err := d.tryUEndian(22, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U22BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU22BE tries to add a field and read 22 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU22BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(22, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU22BE adds a field and reads 22 bit unsigned integer in big-endian
func (d *D) FieldScalarU22BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU22BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U22BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU22BE tries to add a field and read 22 bit unsigned integer in big-endian
func (d *D) TryFieldU22BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU22BE(name, sms...)
	return s.ActualU(), err
}

// FieldU22BE adds a field and reads 22 bit unsigned integer in big-endian
func (d *D) FieldU22BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU22BE(name, sms...).ActualU()
}

// Reader U23BE

// TryU23BE tries to read 23 bit unsigned integer in big-endian
func (d *D) TryU23BE() (uint64, error) { return d.tryUEndian(23, BigEndian) }

// U23BE reads 23 bit unsigned integer in big-endian
func (d *D) U23BE() uint64 {
	v, err := d.tryUEndian(23, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U23BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU23BE tries to add a field and read 23 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU23BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(23, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU23BE adds a field and reads 23 bit unsigned integer in big-endian
func (d *D) FieldScalarU23BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU23BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U23BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU23BE tries to add a field and read 23 bit unsigned integer in big-endian
func (d *D) TryFieldU23BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU23BE(name, sms...)
	return s.ActualU(), err
}

// FieldU23BE adds a field and reads 23 bit unsigned integer in big-endian
func (d *D) FieldU23BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU23BE(name, sms...).ActualU()
}

// Reader U24BE

// TryU24BE tries to read 24 bit unsigned integer in big-endian
func (d *D) TryU24BE() (uint64, error) { return d.tryUEndian(24, BigEndian) }

// U24BE reads 24 bit unsigned integer in big-endian
func (d *D) U24BE() uint64 {
	v, err := d.tryUEndian(24, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U24BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU24BE tries to add a field and read 24 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU24BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(24, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU24BE adds a field and reads 24 bit unsigned integer in big-endian
func (d *D) FieldScalarU24BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU24BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U24BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU24BE tries to add a field and read 24 bit unsigned integer in big-endian
func (d *D) TryFieldU24BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU24BE(name, sms...)
	return s.ActualU(), err
}

// FieldU24BE adds a field and reads 24 bit unsigned integer in big-endian
func (d *D) FieldU24BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU24BE(name, sms...).ActualU()
}

// Reader U25BE

// TryU25BE tries to read 25 bit unsigned integer in big-endian
func (d *D) TryU25BE() (uint64, error) { return d.tryUEndian(25, BigEndian) }

// U25BE reads 25 bit unsigned integer in big-endian
func (d *D) U25BE() uint64 {
	v, err := d.tryUEndian(25, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U25BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU25BE tries to add a field and read 25 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU25BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(25, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU25BE adds a field and reads 25 bit unsigned integer in big-endian
func (d *D) FieldScalarU25BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU25BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U25BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU25BE tries to add a field and read 25 bit unsigned integer in big-endian
func (d *D) TryFieldU25BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU25BE(name, sms...)
	return s.ActualU(), err
}

// FieldU25BE adds a field and reads 25 bit unsigned integer in big-endian
func (d *D) FieldU25BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU25BE(name, sms...).ActualU()
}

// Reader U26BE

// TryU26BE tries to read 26 bit unsigned integer in big-endian
func (d *D) TryU26BE() (uint64, error) { return d.tryUEndian(26, BigEndian) }

// U26BE reads 26 bit unsigned integer in big-endian
func (d *D) U26BE() uint64 {
	v, err := d.tryUEndian(26, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U26BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU26BE tries to add a field and read 26 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU26BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(26, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU26BE adds a field and reads 26 bit unsigned integer in big-endian
func (d *D) FieldScalarU26BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU26BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U26BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU26BE tries to add a field and read 26 bit unsigned integer in big-endian
func (d *D) TryFieldU26BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU26BE(name, sms...)
	return s.ActualU(), err
}

// FieldU26BE adds a field and reads 26 bit unsigned integer in big-endian
func (d *D) FieldU26BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU26BE(name, sms...).ActualU()
}

// Reader U27BE

// TryU27BE tries to read 27 bit unsigned integer in big-endian
func (d *D) TryU27BE() (uint64, error) { return d.tryUEndian(27, BigEndian) }

// U27BE reads 27 bit unsigned integer in big-endian
func (d *D) U27BE() uint64 {
	v, err := d.tryUEndian(27, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U27BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU27BE tries to add a field and read 27 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU27BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(27, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU27BE adds a field and reads 27 bit unsigned integer in big-endian
func (d *D) FieldScalarU27BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU27BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U27BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU27BE tries to add a field and read 27 bit unsigned integer in big-endian
func (d *D) TryFieldU27BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU27BE(name, sms...)
	return s.ActualU(), err
}

// FieldU27BE adds a field and reads 27 bit unsigned integer in big-endian
func (d *D) FieldU27BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU27BE(name, sms...).ActualU()
}

// Reader U28BE

// TryU28BE tries to read 28 bit unsigned integer in big-endian
func (d *D) TryU28BE() (uint64, error) { return d.tryUEndian(28, BigEndian) }

// U28BE reads 28 bit unsigned integer in big-endian
func (d *D) U28BE() uint64 {
	v, err := d.tryUEndian(28, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U28BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU28BE tries to add a field and read 28 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU28BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(28, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU28BE adds a field and reads 28 bit unsigned integer in big-endian
func (d *D) FieldScalarU28BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU28BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U28BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU28BE tries to add a field and read 28 bit unsigned integer in big-endian
func (d *D) TryFieldU28BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU28BE(name, sms...)
	return s.ActualU(), err
}

// FieldU28BE adds a field and reads 28 bit unsigned integer in big-endian
func (d *D) FieldU28BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU28BE(name, sms...).ActualU()
}

// Reader U29BE

// TryU29BE tries to read 29 bit unsigned integer in big-endian
func (d *D) TryU29BE() (uint64, error) { return d.tryUEndian(29, BigEndian) }

// U29BE reads 29 bit unsigned integer in big-endian
func (d *D) U29BE() uint64 {
	v, err := d.tryUEndian(29, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U29BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU29BE tries to add a field and read 29 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU29BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(29, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU29BE adds a field and reads 29 bit unsigned integer in big-endian
func (d *D) FieldScalarU29BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU29BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U29BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU29BE tries to add a field and read 29 bit unsigned integer in big-endian
func (d *D) TryFieldU29BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU29BE(name, sms...)
	return s.ActualU(), err
}

// FieldU29BE adds a field and reads 29 bit unsigned integer in big-endian
func (d *D) FieldU29BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU29BE(name, sms...).ActualU()
}

// Reader U30BE

// TryU30BE tries to read 30 bit unsigned integer in big-endian
func (d *D) TryU30BE() (uint64, error) { return d.tryUEndian(30, BigEndian) }

// U30BE reads 30 bit unsigned integer in big-endian
func (d *D) U30BE() uint64 {
	v, err := d.tryUEndian(30, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U30BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU30BE tries to add a field and read 30 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU30BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(30, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU30BE adds a field and reads 30 bit unsigned integer in big-endian
func (d *D) FieldScalarU30BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU30BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U30BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU30BE tries to add a field and read 30 bit unsigned integer in big-endian
func (d *D) TryFieldU30BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU30BE(name, sms...)
	return s.ActualU(), err
}

// FieldU30BE adds a field and reads 30 bit unsigned integer in big-endian
func (d *D) FieldU30BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU30BE(name, sms...).ActualU()
}

// Reader U31BE

// TryU31BE tries to read 31 bit unsigned integer in big-endian
func (d *D) TryU31BE() (uint64, error) { return d.tryUEndian(31, BigEndian) }

// U31BE reads 31 bit unsigned integer in big-endian
func (d *D) U31BE() uint64 {
	v, err := d.tryUEndian(31, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U31BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU31BE tries to add a field and read 31 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU31BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(31, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU31BE adds a field and reads 31 bit unsigned integer in big-endian
func (d *D) FieldScalarU31BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU31BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U31BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU31BE tries to add a field and read 31 bit unsigned integer in big-endian
func (d *D) TryFieldU31BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU31BE(name, sms...)
	return s.ActualU(), err
}

// FieldU31BE adds a field and reads 31 bit unsigned integer in big-endian
func (d *D) FieldU31BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU31BE(name, sms...).ActualU()
}

// Reader U32BE

// TryU32BE tries to read 32 bit unsigned integer in big-endian
func (d *D) TryU32BE() (uint64, error) { return d.tryUEndian(32, BigEndian) }

// U32BE reads 32 bit unsigned integer in big-endian
func (d *D) U32BE() uint64 {
	v, err := d.tryUEndian(32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U32BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU32BE tries to add a field and read 32 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU32BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(32, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU32BE adds a field and reads 32 bit unsigned integer in big-endian
func (d *D) FieldScalarU32BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU32BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U32BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU32BE tries to add a field and read 32 bit unsigned integer in big-endian
func (d *D) TryFieldU32BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU32BE(name, sms...)
	return s.ActualU(), err
}

// FieldU32BE adds a field and reads 32 bit unsigned integer in big-endian
func (d *D) FieldU32BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU32BE(name, sms...).ActualU()
}

// Reader U33BE

// TryU33BE tries to read 33 bit unsigned integer in big-endian
func (d *D) TryU33BE() (uint64, error) { return d.tryUEndian(33, BigEndian) }

// U33BE reads 33 bit unsigned integer in big-endian
func (d *D) U33BE() uint64 {
	v, err := d.tryUEndian(33, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U33BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU33BE tries to add a field and read 33 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU33BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(33, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU33BE adds a field and reads 33 bit unsigned integer in big-endian
func (d *D) FieldScalarU33BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU33BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U33BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU33BE tries to add a field and read 33 bit unsigned integer in big-endian
func (d *D) TryFieldU33BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU33BE(name, sms...)
	return s.ActualU(), err
}

// FieldU33BE adds a field and reads 33 bit unsigned integer in big-endian
func (d *D) FieldU33BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU33BE(name, sms...).ActualU()
}

// Reader U34BE

// TryU34BE tries to read 34 bit unsigned integer in big-endian
func (d *D) TryU34BE() (uint64, error) { return d.tryUEndian(34, BigEndian) }

// U34BE reads 34 bit unsigned integer in big-endian
func (d *D) U34BE() uint64 {
	v, err := d.tryUEndian(34, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U34BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU34BE tries to add a field and read 34 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU34BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(34, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU34BE adds a field and reads 34 bit unsigned integer in big-endian
func (d *D) FieldScalarU34BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU34BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U34BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU34BE tries to add a field and read 34 bit unsigned integer in big-endian
func (d *D) TryFieldU34BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU34BE(name, sms...)
	return s.ActualU(), err
}

// FieldU34BE adds a field and reads 34 bit unsigned integer in big-endian
func (d *D) FieldU34BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU34BE(name, sms...).ActualU()
}

// Reader U35BE

// TryU35BE tries to read 35 bit unsigned integer in big-endian
func (d *D) TryU35BE() (uint64, error) { return d.tryUEndian(35, BigEndian) }

// U35BE reads 35 bit unsigned integer in big-endian
func (d *D) U35BE() uint64 {
	v, err := d.tryUEndian(35, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U35BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU35BE tries to add a field and read 35 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU35BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(35, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU35BE adds a field and reads 35 bit unsigned integer in big-endian
func (d *D) FieldScalarU35BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU35BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U35BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU35BE tries to add a field and read 35 bit unsigned integer in big-endian
func (d *D) TryFieldU35BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU35BE(name, sms...)
	return s.ActualU(), err
}

// FieldU35BE adds a field and reads 35 bit unsigned integer in big-endian
func (d *D) FieldU35BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU35BE(name, sms...).ActualU()
}

// Reader U36BE

// TryU36BE tries to read 36 bit unsigned integer in big-endian
func (d *D) TryU36BE() (uint64, error) { return d.tryUEndian(36, BigEndian) }

// U36BE reads 36 bit unsigned integer in big-endian
func (d *D) U36BE() uint64 {
	v, err := d.tryUEndian(36, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U36BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU36BE tries to add a field and read 36 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU36BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(36, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU36BE adds a field and reads 36 bit unsigned integer in big-endian
func (d *D) FieldScalarU36BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU36BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U36BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU36BE tries to add a field and read 36 bit unsigned integer in big-endian
func (d *D) TryFieldU36BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU36BE(name, sms...)
	return s.ActualU(), err
}

// FieldU36BE adds a field and reads 36 bit unsigned integer in big-endian
func (d *D) FieldU36BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU36BE(name, sms...).ActualU()
}

// Reader U37BE

// TryU37BE tries to read 37 bit unsigned integer in big-endian
func (d *D) TryU37BE() (uint64, error) { return d.tryUEndian(37, BigEndian) }

// U37BE reads 37 bit unsigned integer in big-endian
func (d *D) U37BE() uint64 {
	v, err := d.tryUEndian(37, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U37BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU37BE tries to add a field and read 37 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU37BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(37, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU37BE adds a field and reads 37 bit unsigned integer in big-endian
func (d *D) FieldScalarU37BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU37BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U37BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU37BE tries to add a field and read 37 bit unsigned integer in big-endian
func (d *D) TryFieldU37BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU37BE(name, sms...)
	return s.ActualU(), err
}

// FieldU37BE adds a field and reads 37 bit unsigned integer in big-endian
func (d *D) FieldU37BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU37BE(name, sms...).ActualU()
}

// Reader U38BE

// TryU38BE tries to read 38 bit unsigned integer in big-endian
func (d *D) TryU38BE() (uint64, error) { return d.tryUEndian(38, BigEndian) }

// U38BE reads 38 bit unsigned integer in big-endian
func (d *D) U38BE() uint64 {
	v, err := d.tryUEndian(38, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U38BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU38BE tries to add a field and read 38 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU38BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(38, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU38BE adds a field and reads 38 bit unsigned integer in big-endian
func (d *D) FieldScalarU38BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU38BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U38BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU38BE tries to add a field and read 38 bit unsigned integer in big-endian
func (d *D) TryFieldU38BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU38BE(name, sms...)
	return s.ActualU(), err
}

// FieldU38BE adds a field and reads 38 bit unsigned integer in big-endian
func (d *D) FieldU38BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU38BE(name, sms...).ActualU()
}

// Reader U39BE

// TryU39BE tries to read 39 bit unsigned integer in big-endian
func (d *D) TryU39BE() (uint64, error) { return d.tryUEndian(39, BigEndian) }

// U39BE reads 39 bit unsigned integer in big-endian
func (d *D) U39BE() uint64 {
	v, err := d.tryUEndian(39, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U39BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU39BE tries to add a field and read 39 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU39BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(39, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU39BE adds a field and reads 39 bit unsigned integer in big-endian
func (d *D) FieldScalarU39BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU39BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U39BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU39BE tries to add a field and read 39 bit unsigned integer in big-endian
func (d *D) TryFieldU39BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU39BE(name, sms...)
	return s.ActualU(), err
}

// FieldU39BE adds a field and reads 39 bit unsigned integer in big-endian
func (d *D) FieldU39BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU39BE(name, sms...).ActualU()
}

// Reader U40BE

// TryU40BE tries to read 40 bit unsigned integer in big-endian
func (d *D) TryU40BE() (uint64, error) { return d.tryUEndian(40, BigEndian) }

// U40BE reads 40 bit unsigned integer in big-endian
func (d *D) U40BE() uint64 {
	v, err := d.tryUEndian(40, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U40BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU40BE tries to add a field and read 40 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU40BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(40, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU40BE adds a field and reads 40 bit unsigned integer in big-endian
func (d *D) FieldScalarU40BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU40BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U40BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU40BE tries to add a field and read 40 bit unsigned integer in big-endian
func (d *D) TryFieldU40BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU40BE(name, sms...)
	return s.ActualU(), err
}

// FieldU40BE adds a field and reads 40 bit unsigned integer in big-endian
func (d *D) FieldU40BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU40BE(name, sms...).ActualU()
}

// Reader U41BE

// TryU41BE tries to read 41 bit unsigned integer in big-endian
func (d *D) TryU41BE() (uint64, error) { return d.tryUEndian(41, BigEndian) }

// U41BE reads 41 bit unsigned integer in big-endian
func (d *D) U41BE() uint64 {
	v, err := d.tryUEndian(41, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U41BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU41BE tries to add a field and read 41 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU41BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(41, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU41BE adds a field and reads 41 bit unsigned integer in big-endian
func (d *D) FieldScalarU41BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU41BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U41BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU41BE tries to add a field and read 41 bit unsigned integer in big-endian
func (d *D) TryFieldU41BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU41BE(name, sms...)
	return s.ActualU(), err
}

// FieldU41BE adds a field and reads 41 bit unsigned integer in big-endian
func (d *D) FieldU41BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU41BE(name, sms...).ActualU()
}

// Reader U42BE

// TryU42BE tries to read 42 bit unsigned integer in big-endian
func (d *D) TryU42BE() (uint64, error) { return d.tryUEndian(42, BigEndian) }

// U42BE reads 42 bit unsigned integer in big-endian
func (d *D) U42BE() uint64 {
	v, err := d.tryUEndian(42, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U42BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU42BE tries to add a field and read 42 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU42BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(42, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU42BE adds a field and reads 42 bit unsigned integer in big-endian
func (d *D) FieldScalarU42BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU42BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U42BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU42BE tries to add a field and read 42 bit unsigned integer in big-endian
func (d *D) TryFieldU42BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU42BE(name, sms...)
	return s.ActualU(), err
}

// FieldU42BE adds a field and reads 42 bit unsigned integer in big-endian
func (d *D) FieldU42BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU42BE(name, sms...).ActualU()
}

// Reader U43BE

// TryU43BE tries to read 43 bit unsigned integer in big-endian
func (d *D) TryU43BE() (uint64, error) { return d.tryUEndian(43, BigEndian) }

// U43BE reads 43 bit unsigned integer in big-endian
func (d *D) U43BE() uint64 {
	v, err := d.tryUEndian(43, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U43BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU43BE tries to add a field and read 43 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU43BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(43, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU43BE adds a field and reads 43 bit unsigned integer in big-endian
func (d *D) FieldScalarU43BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU43BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U43BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU43BE tries to add a field and read 43 bit unsigned integer in big-endian
func (d *D) TryFieldU43BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU43BE(name, sms...)
	return s.ActualU(), err
}

// FieldU43BE adds a field and reads 43 bit unsigned integer in big-endian
func (d *D) FieldU43BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU43BE(name, sms...).ActualU()
}

// Reader U44BE

// TryU44BE tries to read 44 bit unsigned integer in big-endian
func (d *D) TryU44BE() (uint64, error) { return d.tryUEndian(44, BigEndian) }

// U44BE reads 44 bit unsigned integer in big-endian
func (d *D) U44BE() uint64 {
	v, err := d.tryUEndian(44, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U44BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU44BE tries to add a field and read 44 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU44BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(44, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU44BE adds a field and reads 44 bit unsigned integer in big-endian
func (d *D) FieldScalarU44BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU44BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U44BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU44BE tries to add a field and read 44 bit unsigned integer in big-endian
func (d *D) TryFieldU44BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU44BE(name, sms...)
	return s.ActualU(), err
}

// FieldU44BE adds a field and reads 44 bit unsigned integer in big-endian
func (d *D) FieldU44BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU44BE(name, sms...).ActualU()
}

// Reader U45BE

// TryU45BE tries to read 45 bit unsigned integer in big-endian
func (d *D) TryU45BE() (uint64, error) { return d.tryUEndian(45, BigEndian) }

// U45BE reads 45 bit unsigned integer in big-endian
func (d *D) U45BE() uint64 {
	v, err := d.tryUEndian(45, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U45BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU45BE tries to add a field and read 45 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU45BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(45, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU45BE adds a field and reads 45 bit unsigned integer in big-endian
func (d *D) FieldScalarU45BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU45BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U45BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU45BE tries to add a field and read 45 bit unsigned integer in big-endian
func (d *D) TryFieldU45BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU45BE(name, sms...)
	return s.ActualU(), err
}

// FieldU45BE adds a field and reads 45 bit unsigned integer in big-endian
func (d *D) FieldU45BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU45BE(name, sms...).ActualU()
}

// Reader U46BE

// TryU46BE tries to read 46 bit unsigned integer in big-endian
func (d *D) TryU46BE() (uint64, error) { return d.tryUEndian(46, BigEndian) }

// U46BE reads 46 bit unsigned integer in big-endian
func (d *D) U46BE() uint64 {
	v, err := d.tryUEndian(46, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U46BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU46BE tries to add a field and read 46 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU46BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(46, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU46BE adds a field and reads 46 bit unsigned integer in big-endian
func (d *D) FieldScalarU46BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU46BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U46BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU46BE tries to add a field and read 46 bit unsigned integer in big-endian
func (d *D) TryFieldU46BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU46BE(name, sms...)
	return s.ActualU(), err
}

// FieldU46BE adds a field and reads 46 bit unsigned integer in big-endian
func (d *D) FieldU46BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU46BE(name, sms...).ActualU()
}

// Reader U47BE

// TryU47BE tries to read 47 bit unsigned integer in big-endian
func (d *D) TryU47BE() (uint64, error) { return d.tryUEndian(47, BigEndian) }

// U47BE reads 47 bit unsigned integer in big-endian
func (d *D) U47BE() uint64 {
	v, err := d.tryUEndian(47, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U47BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU47BE tries to add a field and read 47 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU47BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(47, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU47BE adds a field and reads 47 bit unsigned integer in big-endian
func (d *D) FieldScalarU47BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU47BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U47BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU47BE tries to add a field and read 47 bit unsigned integer in big-endian
func (d *D) TryFieldU47BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU47BE(name, sms...)
	return s.ActualU(), err
}

// FieldU47BE adds a field and reads 47 bit unsigned integer in big-endian
func (d *D) FieldU47BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU47BE(name, sms...).ActualU()
}

// Reader U48BE

// TryU48BE tries to read 48 bit unsigned integer in big-endian
func (d *D) TryU48BE() (uint64, error) { return d.tryUEndian(48, BigEndian) }

// U48BE reads 48 bit unsigned integer in big-endian
func (d *D) U48BE() uint64 {
	v, err := d.tryUEndian(48, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U48BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU48BE tries to add a field and read 48 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU48BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(48, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU48BE adds a field and reads 48 bit unsigned integer in big-endian
func (d *D) FieldScalarU48BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU48BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U48BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU48BE tries to add a field and read 48 bit unsigned integer in big-endian
func (d *D) TryFieldU48BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU48BE(name, sms...)
	return s.ActualU(), err
}

// FieldU48BE adds a field and reads 48 bit unsigned integer in big-endian
func (d *D) FieldU48BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU48BE(name, sms...).ActualU()
}

// Reader U49BE

// TryU49BE tries to read 49 bit unsigned integer in big-endian
func (d *D) TryU49BE() (uint64, error) { return d.tryUEndian(49, BigEndian) }

// U49BE reads 49 bit unsigned integer in big-endian
func (d *D) U49BE() uint64 {
	v, err := d.tryUEndian(49, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U49BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU49BE tries to add a field and read 49 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU49BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(49, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU49BE adds a field and reads 49 bit unsigned integer in big-endian
func (d *D) FieldScalarU49BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU49BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U49BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU49BE tries to add a field and read 49 bit unsigned integer in big-endian
func (d *D) TryFieldU49BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU49BE(name, sms...)
	return s.ActualU(), err
}

// FieldU49BE adds a field and reads 49 bit unsigned integer in big-endian
func (d *D) FieldU49BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU49BE(name, sms...).ActualU()
}

// Reader U50BE

// TryU50BE tries to read 50 bit unsigned integer in big-endian
func (d *D) TryU50BE() (uint64, error) { return d.tryUEndian(50, BigEndian) }

// U50BE reads 50 bit unsigned integer in big-endian
func (d *D) U50BE() uint64 {
	v, err := d.tryUEndian(50, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U50BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU50BE tries to add a field and read 50 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU50BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(50, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU50BE adds a field and reads 50 bit unsigned integer in big-endian
func (d *D) FieldScalarU50BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU50BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U50BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU50BE tries to add a field and read 50 bit unsigned integer in big-endian
func (d *D) TryFieldU50BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU50BE(name, sms...)
	return s.ActualU(), err
}

// FieldU50BE adds a field and reads 50 bit unsigned integer in big-endian
func (d *D) FieldU50BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU50BE(name, sms...).ActualU()
}

// Reader U51BE

// TryU51BE tries to read 51 bit unsigned integer in big-endian
func (d *D) TryU51BE() (uint64, error) { return d.tryUEndian(51, BigEndian) }

// U51BE reads 51 bit unsigned integer in big-endian
func (d *D) U51BE() uint64 {
	v, err := d.tryUEndian(51, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U51BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU51BE tries to add a field and read 51 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU51BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(51, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU51BE adds a field and reads 51 bit unsigned integer in big-endian
func (d *D) FieldScalarU51BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU51BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U51BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU51BE tries to add a field and read 51 bit unsigned integer in big-endian
func (d *D) TryFieldU51BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU51BE(name, sms...)
	return s.ActualU(), err
}

// FieldU51BE adds a field and reads 51 bit unsigned integer in big-endian
func (d *D) FieldU51BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU51BE(name, sms...).ActualU()
}

// Reader U52BE

// TryU52BE tries to read 52 bit unsigned integer in big-endian
func (d *D) TryU52BE() (uint64, error) { return d.tryUEndian(52, BigEndian) }

// U52BE reads 52 bit unsigned integer in big-endian
func (d *D) U52BE() uint64 {
	v, err := d.tryUEndian(52, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U52BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU52BE tries to add a field and read 52 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU52BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(52, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU52BE adds a field and reads 52 bit unsigned integer in big-endian
func (d *D) FieldScalarU52BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU52BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U52BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU52BE tries to add a field and read 52 bit unsigned integer in big-endian
func (d *D) TryFieldU52BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU52BE(name, sms...)
	return s.ActualU(), err
}

// FieldU52BE adds a field and reads 52 bit unsigned integer in big-endian
func (d *D) FieldU52BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU52BE(name, sms...).ActualU()
}

// Reader U53BE

// TryU53BE tries to read 53 bit unsigned integer in big-endian
func (d *D) TryU53BE() (uint64, error) { return d.tryUEndian(53, BigEndian) }

// U53BE reads 53 bit unsigned integer in big-endian
func (d *D) U53BE() uint64 {
	v, err := d.tryUEndian(53, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U53BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU53BE tries to add a field and read 53 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU53BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(53, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU53BE adds a field and reads 53 bit unsigned integer in big-endian
func (d *D) FieldScalarU53BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU53BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U53BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU53BE tries to add a field and read 53 bit unsigned integer in big-endian
func (d *D) TryFieldU53BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU53BE(name, sms...)
	return s.ActualU(), err
}

// FieldU53BE adds a field and reads 53 bit unsigned integer in big-endian
func (d *D) FieldU53BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU53BE(name, sms...).ActualU()
}

// Reader U54BE

// TryU54BE tries to read 54 bit unsigned integer in big-endian
func (d *D) TryU54BE() (uint64, error) { return d.tryUEndian(54, BigEndian) }

// U54BE reads 54 bit unsigned integer in big-endian
func (d *D) U54BE() uint64 {
	v, err := d.tryUEndian(54, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U54BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU54BE tries to add a field and read 54 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU54BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(54, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU54BE adds a field and reads 54 bit unsigned integer in big-endian
func (d *D) FieldScalarU54BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU54BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U54BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU54BE tries to add a field and read 54 bit unsigned integer in big-endian
func (d *D) TryFieldU54BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU54BE(name, sms...)
	return s.ActualU(), err
}

// FieldU54BE adds a field and reads 54 bit unsigned integer in big-endian
func (d *D) FieldU54BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU54BE(name, sms...).ActualU()
}

// Reader U55BE

// TryU55BE tries to read 55 bit unsigned integer in big-endian
func (d *D) TryU55BE() (uint64, error) { return d.tryUEndian(55, BigEndian) }

// U55BE reads 55 bit unsigned integer in big-endian
func (d *D) U55BE() uint64 {
	v, err := d.tryUEndian(55, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U55BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU55BE tries to add a field and read 55 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU55BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(55, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU55BE adds a field and reads 55 bit unsigned integer in big-endian
func (d *D) FieldScalarU55BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU55BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U55BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU55BE tries to add a field and read 55 bit unsigned integer in big-endian
func (d *D) TryFieldU55BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU55BE(name, sms...)
	return s.ActualU(), err
}

// FieldU55BE adds a field and reads 55 bit unsigned integer in big-endian
func (d *D) FieldU55BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU55BE(name, sms...).ActualU()
}

// Reader U56BE

// TryU56BE tries to read 56 bit unsigned integer in big-endian
func (d *D) TryU56BE() (uint64, error) { return d.tryUEndian(56, BigEndian) }

// U56BE reads 56 bit unsigned integer in big-endian
func (d *D) U56BE() uint64 {
	v, err := d.tryUEndian(56, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U56BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU56BE tries to add a field and read 56 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU56BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(56, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU56BE adds a field and reads 56 bit unsigned integer in big-endian
func (d *D) FieldScalarU56BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU56BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U56BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU56BE tries to add a field and read 56 bit unsigned integer in big-endian
func (d *D) TryFieldU56BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU56BE(name, sms...)
	return s.ActualU(), err
}

// FieldU56BE adds a field and reads 56 bit unsigned integer in big-endian
func (d *D) FieldU56BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU56BE(name, sms...).ActualU()
}

// Reader U57BE

// TryU57BE tries to read 57 bit unsigned integer in big-endian
func (d *D) TryU57BE() (uint64, error) { return d.tryUEndian(57, BigEndian) }

// U57BE reads 57 bit unsigned integer in big-endian
func (d *D) U57BE() uint64 {
	v, err := d.tryUEndian(57, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U57BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU57BE tries to add a field and read 57 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU57BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(57, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU57BE adds a field and reads 57 bit unsigned integer in big-endian
func (d *D) FieldScalarU57BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU57BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U57BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU57BE tries to add a field and read 57 bit unsigned integer in big-endian
func (d *D) TryFieldU57BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU57BE(name, sms...)
	return s.ActualU(), err
}

// FieldU57BE adds a field and reads 57 bit unsigned integer in big-endian
func (d *D) FieldU57BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU57BE(name, sms...).ActualU()
}

// Reader U58BE

// TryU58BE tries to read 58 bit unsigned integer in big-endian
func (d *D) TryU58BE() (uint64, error) { return d.tryUEndian(58, BigEndian) }

// U58BE reads 58 bit unsigned integer in big-endian
func (d *D) U58BE() uint64 {
	v, err := d.tryUEndian(58, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U58BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU58BE tries to add a field and read 58 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU58BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(58, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU58BE adds a field and reads 58 bit unsigned integer in big-endian
func (d *D) FieldScalarU58BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU58BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U58BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU58BE tries to add a field and read 58 bit unsigned integer in big-endian
func (d *D) TryFieldU58BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU58BE(name, sms...)
	return s.ActualU(), err
}

// FieldU58BE adds a field and reads 58 bit unsigned integer in big-endian
func (d *D) FieldU58BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU58BE(name, sms...).ActualU()
}

// Reader U59BE

// TryU59BE tries to read 59 bit unsigned integer in big-endian
func (d *D) TryU59BE() (uint64, error) { return d.tryUEndian(59, BigEndian) }

// U59BE reads 59 bit unsigned integer in big-endian
func (d *D) U59BE() uint64 {
	v, err := d.tryUEndian(59, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U59BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU59BE tries to add a field and read 59 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU59BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(59, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU59BE adds a field and reads 59 bit unsigned integer in big-endian
func (d *D) FieldScalarU59BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU59BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U59BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU59BE tries to add a field and read 59 bit unsigned integer in big-endian
func (d *D) TryFieldU59BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU59BE(name, sms...)
	return s.ActualU(), err
}

// FieldU59BE adds a field and reads 59 bit unsigned integer in big-endian
func (d *D) FieldU59BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU59BE(name, sms...).ActualU()
}

// Reader U60BE

// TryU60BE tries to read 60 bit unsigned integer in big-endian
func (d *D) TryU60BE() (uint64, error) { return d.tryUEndian(60, BigEndian) }

// U60BE reads 60 bit unsigned integer in big-endian
func (d *D) U60BE() uint64 {
	v, err := d.tryUEndian(60, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U60BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU60BE tries to add a field and read 60 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU60BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(60, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU60BE adds a field and reads 60 bit unsigned integer in big-endian
func (d *D) FieldScalarU60BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU60BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U60BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU60BE tries to add a field and read 60 bit unsigned integer in big-endian
func (d *D) TryFieldU60BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU60BE(name, sms...)
	return s.ActualU(), err
}

// FieldU60BE adds a field and reads 60 bit unsigned integer in big-endian
func (d *D) FieldU60BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU60BE(name, sms...).ActualU()
}

// Reader U61BE

// TryU61BE tries to read 61 bit unsigned integer in big-endian
func (d *D) TryU61BE() (uint64, error) { return d.tryUEndian(61, BigEndian) }

// U61BE reads 61 bit unsigned integer in big-endian
func (d *D) U61BE() uint64 {
	v, err := d.tryUEndian(61, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U61BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU61BE tries to add a field and read 61 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU61BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(61, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU61BE adds a field and reads 61 bit unsigned integer in big-endian
func (d *D) FieldScalarU61BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU61BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U61BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU61BE tries to add a field and read 61 bit unsigned integer in big-endian
func (d *D) TryFieldU61BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU61BE(name, sms...)
	return s.ActualU(), err
}

// FieldU61BE adds a field and reads 61 bit unsigned integer in big-endian
func (d *D) FieldU61BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU61BE(name, sms...).ActualU()
}

// Reader U62BE

// TryU62BE tries to read 62 bit unsigned integer in big-endian
func (d *D) TryU62BE() (uint64, error) { return d.tryUEndian(62, BigEndian) }

// U62BE reads 62 bit unsigned integer in big-endian
func (d *D) U62BE() uint64 {
	v, err := d.tryUEndian(62, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U62BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU62BE tries to add a field and read 62 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU62BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(62, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU62BE adds a field and reads 62 bit unsigned integer in big-endian
func (d *D) FieldScalarU62BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU62BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U62BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU62BE tries to add a field and read 62 bit unsigned integer in big-endian
func (d *D) TryFieldU62BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU62BE(name, sms...)
	return s.ActualU(), err
}

// FieldU62BE adds a field and reads 62 bit unsigned integer in big-endian
func (d *D) FieldU62BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU62BE(name, sms...).ActualU()
}

// Reader U63BE

// TryU63BE tries to read 63 bit unsigned integer in big-endian
func (d *D) TryU63BE() (uint64, error) { return d.tryUEndian(63, BigEndian) }

// U63BE reads 63 bit unsigned integer in big-endian
func (d *D) U63BE() uint64 {
	v, err := d.tryUEndian(63, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U63BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU63BE tries to add a field and read 63 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU63BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(63, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU63BE adds a field and reads 63 bit unsigned integer in big-endian
func (d *D) FieldScalarU63BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU63BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U63BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU63BE tries to add a field and read 63 bit unsigned integer in big-endian
func (d *D) TryFieldU63BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU63BE(name, sms...)
	return s.ActualU(), err
}

// FieldU63BE adds a field and reads 63 bit unsigned integer in big-endian
func (d *D) FieldU63BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU63BE(name, sms...).ActualU()
}

// Reader U64BE

// TryU64BE tries to read 64 bit unsigned integer in big-endian
func (d *D) TryU64BE() (uint64, error) { return d.tryUEndian(64, BigEndian) }

// U64BE reads 64 bit unsigned integer in big-endian
func (d *D) U64BE() uint64 {
	v, err := d.tryUEndian(64, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U64BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarU64BE tries to add a field and read 64 bit unsigned integer in big-endian
func (d *D) TryFieldScalarU64BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUEndian(64, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarU64BE adds a field and reads 64 bit unsigned integer in big-endian
func (d *D) FieldScalarU64BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarU64BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U64BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldU64BE tries to add a field and read 64 bit unsigned integer in big-endian
func (d *D) TryFieldU64BE(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarU64BE(name, sms...)
	return s.ActualU(), err
}

// FieldU64BE adds a field and reads 64 bit unsigned integer in big-endian
func (d *D) FieldU64BE(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarU64BE(name, sms...).ActualU()
}

// Reader S

// TryS tries to read nBits bits signed integer in current endian
func (d *D) TryS(nBits int) (int64, error) { return d.trySEndian(nBits, d.Endian) }

// S reads nBits bits signed integer in current endian
func (d *D) S(nBits int) int64 {
	v, err := d.trySEndian(nBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS tries to add a field and read nBits bits signed integer in current endian
func (d *D) TryFieldScalarS(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(nBits, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS adds a field and reads nBits bits signed integer in current endian
func (d *D) FieldScalarS(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S", Pos: d.Pos()})
	}
	return s
}

// TryFieldS tries to add a field and read nBits bits signed integer in current endian
func (d *D) TryFieldS(name string, nBits int, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS(name, nBits, sms...)
	return s.ActualS(), err
}

// FieldS adds a field and reads nBits bits signed integer in current endian
func (d *D) FieldS(name string, nBits int, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS(name, nBits, sms...).ActualS()
}

// Reader SE

// TrySE tries to read nBits signed integer in specified endian
func (d *D) TrySE(nBits int, endian Endian) (int64, error) { return d.trySEndian(nBits, endian) }

// SE reads nBits signed integer in specified endian
func (d *D) SE(nBits int, endian Endian) int64 {
	v, err := d.trySEndian(nBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "SE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarSE tries to add a field and read nBits signed integer in specified endian
func (d *D) TryFieldScalarSE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(nBits, endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarSE adds a field and reads nBits signed integer in specified endian
func (d *D) FieldScalarSE(name string, nBits int, endian Endian, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarSE(name, nBits, endian, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "SE", Pos: d.Pos()})
	}
	return s
}

// TryFieldSE tries to add a field and read nBits signed integer in specified endian
func (d *D) TryFieldSE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarSE(name, nBits, endian, sms...)
	return s.ActualS(), err
}

// FieldSE adds a field and reads nBits signed integer in specified endian
func (d *D) FieldSE(name string, nBits int, endian Endian, sms ...scalar.Mapper) int64 {
	return d.FieldScalarSE(name, nBits, endian, sms...).ActualS()
}

// Reader S1

// TryS1 tries to read 1 bit signed integer in current endian
func (d *D) TryS1() (int64, error) { return d.trySEndian(1, d.Endian) }

// S1 reads 1 bit signed integer in current endian
func (d *D) S1() int64 {
	v, err := d.trySEndian(1, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S1", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS1 tries to add a field and read 1 bit signed integer in current endian
func (d *D) TryFieldScalarS1(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(1, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS1 adds a field and reads 1 bit signed integer in current endian
func (d *D) FieldScalarS1(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS1(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S1", Pos: d.Pos()})
	}
	return s
}

// TryFieldS1 tries to add a field and read 1 bit signed integer in current endian
func (d *D) TryFieldS1(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS1(name, sms...)
	return s.ActualS(), err
}

// FieldS1 adds a field and reads 1 bit signed integer in current endian
func (d *D) FieldS1(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS1(name, sms...).ActualS()
}

// Reader S2

// TryS2 tries to read 2 bit signed integer in current endian
func (d *D) TryS2() (int64, error) { return d.trySEndian(2, d.Endian) }

// S2 reads 2 bit signed integer in current endian
func (d *D) S2() int64 {
	v, err := d.trySEndian(2, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S2", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS2 tries to add a field and read 2 bit signed integer in current endian
func (d *D) TryFieldScalarS2(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(2, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS2 adds a field and reads 2 bit signed integer in current endian
func (d *D) FieldScalarS2(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS2(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S2", Pos: d.Pos()})
	}
	return s
}

// TryFieldS2 tries to add a field and read 2 bit signed integer in current endian
func (d *D) TryFieldS2(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS2(name, sms...)
	return s.ActualS(), err
}

// FieldS2 adds a field and reads 2 bit signed integer in current endian
func (d *D) FieldS2(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS2(name, sms...).ActualS()
}

// Reader S3

// TryS3 tries to read 3 bit signed integer in current endian
func (d *D) TryS3() (int64, error) { return d.trySEndian(3, d.Endian) }

// S3 reads 3 bit signed integer in current endian
func (d *D) S3() int64 {
	v, err := d.trySEndian(3, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S3", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS3 tries to add a field and read 3 bit signed integer in current endian
func (d *D) TryFieldScalarS3(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(3, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS3 adds a field and reads 3 bit signed integer in current endian
func (d *D) FieldScalarS3(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS3(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S3", Pos: d.Pos()})
	}
	return s
}

// TryFieldS3 tries to add a field and read 3 bit signed integer in current endian
func (d *D) TryFieldS3(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS3(name, sms...)
	return s.ActualS(), err
}

// FieldS3 adds a field and reads 3 bit signed integer in current endian
func (d *D) FieldS3(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS3(name, sms...).ActualS()
}

// Reader S4

// TryS4 tries to read 4 bit signed integer in current endian
func (d *D) TryS4() (int64, error) { return d.trySEndian(4, d.Endian) }

// S4 reads 4 bit signed integer in current endian
func (d *D) S4() int64 {
	v, err := d.trySEndian(4, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S4", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS4 tries to add a field and read 4 bit signed integer in current endian
func (d *D) TryFieldScalarS4(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(4, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS4 adds a field and reads 4 bit signed integer in current endian
func (d *D) FieldScalarS4(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS4(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S4", Pos: d.Pos()})
	}
	return s
}

// TryFieldS4 tries to add a field and read 4 bit signed integer in current endian
func (d *D) TryFieldS4(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS4(name, sms...)
	return s.ActualS(), err
}

// FieldS4 adds a field and reads 4 bit signed integer in current endian
func (d *D) FieldS4(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS4(name, sms...).ActualS()
}

// Reader S5

// TryS5 tries to read 5 bit signed integer in current endian
func (d *D) TryS5() (int64, error) { return d.trySEndian(5, d.Endian) }

// S5 reads 5 bit signed integer in current endian
func (d *D) S5() int64 {
	v, err := d.trySEndian(5, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S5", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS5 tries to add a field and read 5 bit signed integer in current endian
func (d *D) TryFieldScalarS5(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(5, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS5 adds a field and reads 5 bit signed integer in current endian
func (d *D) FieldScalarS5(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS5(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S5", Pos: d.Pos()})
	}
	return s
}

// TryFieldS5 tries to add a field and read 5 bit signed integer in current endian
func (d *D) TryFieldS5(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS5(name, sms...)
	return s.ActualS(), err
}

// FieldS5 adds a field and reads 5 bit signed integer in current endian
func (d *D) FieldS5(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS5(name, sms...).ActualS()
}

// Reader S6

// TryS6 tries to read 6 bit signed integer in current endian
func (d *D) TryS6() (int64, error) { return d.trySEndian(6, d.Endian) }

// S6 reads 6 bit signed integer in current endian
func (d *D) S6() int64 {
	v, err := d.trySEndian(6, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S6", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS6 tries to add a field and read 6 bit signed integer in current endian
func (d *D) TryFieldScalarS6(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(6, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS6 adds a field and reads 6 bit signed integer in current endian
func (d *D) FieldScalarS6(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS6(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S6", Pos: d.Pos()})
	}
	return s
}

// TryFieldS6 tries to add a field and read 6 bit signed integer in current endian
func (d *D) TryFieldS6(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS6(name, sms...)
	return s.ActualS(), err
}

// FieldS6 adds a field and reads 6 bit signed integer in current endian
func (d *D) FieldS6(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS6(name, sms...).ActualS()
}

// Reader S7

// TryS7 tries to read 7 bit signed integer in current endian
func (d *D) TryS7() (int64, error) { return d.trySEndian(7, d.Endian) }

// S7 reads 7 bit signed integer in current endian
func (d *D) S7() int64 {
	v, err := d.trySEndian(7, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S7", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS7 tries to add a field and read 7 bit signed integer in current endian
func (d *D) TryFieldScalarS7(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(7, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS7 adds a field and reads 7 bit signed integer in current endian
func (d *D) FieldScalarS7(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS7(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S7", Pos: d.Pos()})
	}
	return s
}

// TryFieldS7 tries to add a field and read 7 bit signed integer in current endian
func (d *D) TryFieldS7(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS7(name, sms...)
	return s.ActualS(), err
}

// FieldS7 adds a field and reads 7 bit signed integer in current endian
func (d *D) FieldS7(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS7(name, sms...).ActualS()
}

// Reader S8

// TryS8 tries to read 8 bit signed integer in current endian
func (d *D) TryS8() (int64, error) { return d.trySEndian(8, d.Endian) }

// S8 reads 8 bit signed integer in current endian
func (d *D) S8() int64 {
	v, err := d.trySEndian(8, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S8", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS8 tries to add a field and read 8 bit signed integer in current endian
func (d *D) TryFieldScalarS8(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(8, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS8 adds a field and reads 8 bit signed integer in current endian
func (d *D) FieldScalarS8(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS8(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S8", Pos: d.Pos()})
	}
	return s
}

// TryFieldS8 tries to add a field and read 8 bit signed integer in current endian
func (d *D) TryFieldS8(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS8(name, sms...)
	return s.ActualS(), err
}

// FieldS8 adds a field and reads 8 bit signed integer in current endian
func (d *D) FieldS8(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS8(name, sms...).ActualS()
}

// Reader S9

// TryS9 tries to read 9 bit signed integer in current endian
func (d *D) TryS9() (int64, error) { return d.trySEndian(9, d.Endian) }

// S9 reads 9 bit signed integer in current endian
func (d *D) S9() int64 {
	v, err := d.trySEndian(9, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S9", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS9 tries to add a field and read 9 bit signed integer in current endian
func (d *D) TryFieldScalarS9(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(9, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS9 adds a field and reads 9 bit signed integer in current endian
func (d *D) FieldScalarS9(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS9(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S9", Pos: d.Pos()})
	}
	return s
}

// TryFieldS9 tries to add a field and read 9 bit signed integer in current endian
func (d *D) TryFieldS9(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS9(name, sms...)
	return s.ActualS(), err
}

// FieldS9 adds a field and reads 9 bit signed integer in current endian
func (d *D) FieldS9(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS9(name, sms...).ActualS()
}

// Reader S10

// TryS10 tries to read 10 bit signed integer in current endian
func (d *D) TryS10() (int64, error) { return d.trySEndian(10, d.Endian) }

// S10 reads 10 bit signed integer in current endian
func (d *D) S10() int64 {
	v, err := d.trySEndian(10, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S10", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS10 tries to add a field and read 10 bit signed integer in current endian
func (d *D) TryFieldScalarS10(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(10, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS10 adds a field and reads 10 bit signed integer in current endian
func (d *D) FieldScalarS10(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS10(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S10", Pos: d.Pos()})
	}
	return s
}

// TryFieldS10 tries to add a field and read 10 bit signed integer in current endian
func (d *D) TryFieldS10(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS10(name, sms...)
	return s.ActualS(), err
}

// FieldS10 adds a field and reads 10 bit signed integer in current endian
func (d *D) FieldS10(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS10(name, sms...).ActualS()
}

// Reader S11

// TryS11 tries to read 11 bit signed integer in current endian
func (d *D) TryS11() (int64, error) { return d.trySEndian(11, d.Endian) }

// S11 reads 11 bit signed integer in current endian
func (d *D) S11() int64 {
	v, err := d.trySEndian(11, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S11", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS11 tries to add a field and read 11 bit signed integer in current endian
func (d *D) TryFieldScalarS11(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(11, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS11 adds a field and reads 11 bit signed integer in current endian
func (d *D) FieldScalarS11(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS11(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S11", Pos: d.Pos()})
	}
	return s
}

// TryFieldS11 tries to add a field and read 11 bit signed integer in current endian
func (d *D) TryFieldS11(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS11(name, sms...)
	return s.ActualS(), err
}

// FieldS11 adds a field and reads 11 bit signed integer in current endian
func (d *D) FieldS11(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS11(name, sms...).ActualS()
}

// Reader S12

// TryS12 tries to read 12 bit signed integer in current endian
func (d *D) TryS12() (int64, error) { return d.trySEndian(12, d.Endian) }

// S12 reads 12 bit signed integer in current endian
func (d *D) S12() int64 {
	v, err := d.trySEndian(12, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S12", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS12 tries to add a field and read 12 bit signed integer in current endian
func (d *D) TryFieldScalarS12(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(12, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS12 adds a field and reads 12 bit signed integer in current endian
func (d *D) FieldScalarS12(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS12(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S12", Pos: d.Pos()})
	}
	return s
}

// TryFieldS12 tries to add a field and read 12 bit signed integer in current endian
func (d *D) TryFieldS12(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS12(name, sms...)
	return s.ActualS(), err
}

// FieldS12 adds a field and reads 12 bit signed integer in current endian
func (d *D) FieldS12(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS12(name, sms...).ActualS()
}

// Reader S13

// TryS13 tries to read 13 bit signed integer in current endian
func (d *D) TryS13() (int64, error) { return d.trySEndian(13, d.Endian) }

// S13 reads 13 bit signed integer in current endian
func (d *D) S13() int64 {
	v, err := d.trySEndian(13, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S13", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS13 tries to add a field and read 13 bit signed integer in current endian
func (d *D) TryFieldScalarS13(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(13, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS13 adds a field and reads 13 bit signed integer in current endian
func (d *D) FieldScalarS13(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS13(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S13", Pos: d.Pos()})
	}
	return s
}

// TryFieldS13 tries to add a field and read 13 bit signed integer in current endian
func (d *D) TryFieldS13(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS13(name, sms...)
	return s.ActualS(), err
}

// FieldS13 adds a field and reads 13 bit signed integer in current endian
func (d *D) FieldS13(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS13(name, sms...).ActualS()
}

// Reader S14

// TryS14 tries to read 14 bit signed integer in current endian
func (d *D) TryS14() (int64, error) { return d.trySEndian(14, d.Endian) }

// S14 reads 14 bit signed integer in current endian
func (d *D) S14() int64 {
	v, err := d.trySEndian(14, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S14", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS14 tries to add a field and read 14 bit signed integer in current endian
func (d *D) TryFieldScalarS14(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(14, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS14 adds a field and reads 14 bit signed integer in current endian
func (d *D) FieldScalarS14(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS14(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S14", Pos: d.Pos()})
	}
	return s
}

// TryFieldS14 tries to add a field and read 14 bit signed integer in current endian
func (d *D) TryFieldS14(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS14(name, sms...)
	return s.ActualS(), err
}

// FieldS14 adds a field and reads 14 bit signed integer in current endian
func (d *D) FieldS14(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS14(name, sms...).ActualS()
}

// Reader S15

// TryS15 tries to read 15 bit signed integer in current endian
func (d *D) TryS15() (int64, error) { return d.trySEndian(15, d.Endian) }

// S15 reads 15 bit signed integer in current endian
func (d *D) S15() int64 {
	v, err := d.trySEndian(15, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S15", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS15 tries to add a field and read 15 bit signed integer in current endian
func (d *D) TryFieldScalarS15(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(15, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS15 adds a field and reads 15 bit signed integer in current endian
func (d *D) FieldScalarS15(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS15(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S15", Pos: d.Pos()})
	}
	return s
}

// TryFieldS15 tries to add a field and read 15 bit signed integer in current endian
func (d *D) TryFieldS15(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS15(name, sms...)
	return s.ActualS(), err
}

// FieldS15 adds a field and reads 15 bit signed integer in current endian
func (d *D) FieldS15(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS15(name, sms...).ActualS()
}

// Reader S16

// TryS16 tries to read 16 bit signed integer in current endian
func (d *D) TryS16() (int64, error) { return d.trySEndian(16, d.Endian) }

// S16 reads 16 bit signed integer in current endian
func (d *D) S16() int64 {
	v, err := d.trySEndian(16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S16", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS16 tries to add a field and read 16 bit signed integer in current endian
func (d *D) TryFieldScalarS16(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(16, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS16 adds a field and reads 16 bit signed integer in current endian
func (d *D) FieldScalarS16(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS16(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S16", Pos: d.Pos()})
	}
	return s
}

// TryFieldS16 tries to add a field and read 16 bit signed integer in current endian
func (d *D) TryFieldS16(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS16(name, sms...)
	return s.ActualS(), err
}

// FieldS16 adds a field and reads 16 bit signed integer in current endian
func (d *D) FieldS16(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS16(name, sms...).ActualS()
}

// Reader S17

// TryS17 tries to read 17 bit signed integer in current endian
func (d *D) TryS17() (int64, error) { return d.trySEndian(17, d.Endian) }

// S17 reads 17 bit signed integer in current endian
func (d *D) S17() int64 {
	v, err := d.trySEndian(17, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S17", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS17 tries to add a field and read 17 bit signed integer in current endian
func (d *D) TryFieldScalarS17(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(17, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS17 adds a field and reads 17 bit signed integer in current endian
func (d *D) FieldScalarS17(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS17(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S17", Pos: d.Pos()})
	}
	return s
}

// TryFieldS17 tries to add a field and read 17 bit signed integer in current endian
func (d *D) TryFieldS17(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS17(name, sms...)
	return s.ActualS(), err
}

// FieldS17 adds a field and reads 17 bit signed integer in current endian
func (d *D) FieldS17(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS17(name, sms...).ActualS()
}

// Reader S18

// TryS18 tries to read 18 bit signed integer in current endian
func (d *D) TryS18() (int64, error) { return d.trySEndian(18, d.Endian) }

// S18 reads 18 bit signed integer in current endian
func (d *D) S18() int64 {
	v, err := d.trySEndian(18, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S18", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS18 tries to add a field and read 18 bit signed integer in current endian
func (d *D) TryFieldScalarS18(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(18, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS18 adds a field and reads 18 bit signed integer in current endian
func (d *D) FieldScalarS18(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS18(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S18", Pos: d.Pos()})
	}
	return s
}

// TryFieldS18 tries to add a field and read 18 bit signed integer in current endian
func (d *D) TryFieldS18(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS18(name, sms...)
	return s.ActualS(), err
}

// FieldS18 adds a field and reads 18 bit signed integer in current endian
func (d *D) FieldS18(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS18(name, sms...).ActualS()
}

// Reader S19

// TryS19 tries to read 19 bit signed integer in current endian
func (d *D) TryS19() (int64, error) { return d.trySEndian(19, d.Endian) }

// S19 reads 19 bit signed integer in current endian
func (d *D) S19() int64 {
	v, err := d.trySEndian(19, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S19", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS19 tries to add a field and read 19 bit signed integer in current endian
func (d *D) TryFieldScalarS19(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(19, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS19 adds a field and reads 19 bit signed integer in current endian
func (d *D) FieldScalarS19(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS19(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S19", Pos: d.Pos()})
	}
	return s
}

// TryFieldS19 tries to add a field and read 19 bit signed integer in current endian
func (d *D) TryFieldS19(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS19(name, sms...)
	return s.ActualS(), err
}

// FieldS19 adds a field and reads 19 bit signed integer in current endian
func (d *D) FieldS19(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS19(name, sms...).ActualS()
}

// Reader S20

// TryS20 tries to read 20 bit signed integer in current endian
func (d *D) TryS20() (int64, error) { return d.trySEndian(20, d.Endian) }

// S20 reads 20 bit signed integer in current endian
func (d *D) S20() int64 {
	v, err := d.trySEndian(20, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S20", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS20 tries to add a field and read 20 bit signed integer in current endian
func (d *D) TryFieldScalarS20(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(20, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS20 adds a field and reads 20 bit signed integer in current endian
func (d *D) FieldScalarS20(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS20(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S20", Pos: d.Pos()})
	}
	return s
}

// TryFieldS20 tries to add a field and read 20 bit signed integer in current endian
func (d *D) TryFieldS20(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS20(name, sms...)
	return s.ActualS(), err
}

// FieldS20 adds a field and reads 20 bit signed integer in current endian
func (d *D) FieldS20(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS20(name, sms...).ActualS()
}

// Reader S21

// TryS21 tries to read 21 bit signed integer in current endian
func (d *D) TryS21() (int64, error) { return d.trySEndian(21, d.Endian) }

// S21 reads 21 bit signed integer in current endian
func (d *D) S21() int64 {
	v, err := d.trySEndian(21, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S21", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS21 tries to add a field and read 21 bit signed integer in current endian
func (d *D) TryFieldScalarS21(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(21, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS21 adds a field and reads 21 bit signed integer in current endian
func (d *D) FieldScalarS21(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS21(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S21", Pos: d.Pos()})
	}
	return s
}

// TryFieldS21 tries to add a field and read 21 bit signed integer in current endian
func (d *D) TryFieldS21(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS21(name, sms...)
	return s.ActualS(), err
}

// FieldS21 adds a field and reads 21 bit signed integer in current endian
func (d *D) FieldS21(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS21(name, sms...).ActualS()
}

// Reader S22

// TryS22 tries to read 22 bit signed integer in current endian
func (d *D) TryS22() (int64, error) { return d.trySEndian(22, d.Endian) }

// S22 reads 22 bit signed integer in current endian
func (d *D) S22() int64 {
	v, err := d.trySEndian(22, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S22", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS22 tries to add a field and read 22 bit signed integer in current endian
func (d *D) TryFieldScalarS22(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(22, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS22 adds a field and reads 22 bit signed integer in current endian
func (d *D) FieldScalarS22(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS22(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S22", Pos: d.Pos()})
	}
	return s
}

// TryFieldS22 tries to add a field and read 22 bit signed integer in current endian
func (d *D) TryFieldS22(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS22(name, sms...)
	return s.ActualS(), err
}

// FieldS22 adds a field and reads 22 bit signed integer in current endian
func (d *D) FieldS22(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS22(name, sms...).ActualS()
}

// Reader S23

// TryS23 tries to read 23 bit signed integer in current endian
func (d *D) TryS23() (int64, error) { return d.trySEndian(23, d.Endian) }

// S23 reads 23 bit signed integer in current endian
func (d *D) S23() int64 {
	v, err := d.trySEndian(23, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S23", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS23 tries to add a field and read 23 bit signed integer in current endian
func (d *D) TryFieldScalarS23(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(23, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS23 adds a field and reads 23 bit signed integer in current endian
func (d *D) FieldScalarS23(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS23(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S23", Pos: d.Pos()})
	}
	return s
}

// TryFieldS23 tries to add a field and read 23 bit signed integer in current endian
func (d *D) TryFieldS23(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS23(name, sms...)
	return s.ActualS(), err
}

// FieldS23 adds a field and reads 23 bit signed integer in current endian
func (d *D) FieldS23(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS23(name, sms...).ActualS()
}

// Reader S24

// TryS24 tries to read 24 bit signed integer in current endian
func (d *D) TryS24() (int64, error) { return d.trySEndian(24, d.Endian) }

// S24 reads 24 bit signed integer in current endian
func (d *D) S24() int64 {
	v, err := d.trySEndian(24, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S24", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS24 tries to add a field and read 24 bit signed integer in current endian
func (d *D) TryFieldScalarS24(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(24, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS24 adds a field and reads 24 bit signed integer in current endian
func (d *D) FieldScalarS24(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS24(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S24", Pos: d.Pos()})
	}
	return s
}

// TryFieldS24 tries to add a field and read 24 bit signed integer in current endian
func (d *D) TryFieldS24(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS24(name, sms...)
	return s.ActualS(), err
}

// FieldS24 adds a field and reads 24 bit signed integer in current endian
func (d *D) FieldS24(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS24(name, sms...).ActualS()
}

// Reader S25

// TryS25 tries to read 25 bit signed integer in current endian
func (d *D) TryS25() (int64, error) { return d.trySEndian(25, d.Endian) }

// S25 reads 25 bit signed integer in current endian
func (d *D) S25() int64 {
	v, err := d.trySEndian(25, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S25", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS25 tries to add a field and read 25 bit signed integer in current endian
func (d *D) TryFieldScalarS25(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(25, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS25 adds a field and reads 25 bit signed integer in current endian
func (d *D) FieldScalarS25(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS25(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S25", Pos: d.Pos()})
	}
	return s
}

// TryFieldS25 tries to add a field and read 25 bit signed integer in current endian
func (d *D) TryFieldS25(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS25(name, sms...)
	return s.ActualS(), err
}

// FieldS25 adds a field and reads 25 bit signed integer in current endian
func (d *D) FieldS25(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS25(name, sms...).ActualS()
}

// Reader S26

// TryS26 tries to read 26 bit signed integer in current endian
func (d *D) TryS26() (int64, error) { return d.trySEndian(26, d.Endian) }

// S26 reads 26 bit signed integer in current endian
func (d *D) S26() int64 {
	v, err := d.trySEndian(26, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S26", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS26 tries to add a field and read 26 bit signed integer in current endian
func (d *D) TryFieldScalarS26(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(26, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS26 adds a field and reads 26 bit signed integer in current endian
func (d *D) FieldScalarS26(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS26(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S26", Pos: d.Pos()})
	}
	return s
}

// TryFieldS26 tries to add a field and read 26 bit signed integer in current endian
func (d *D) TryFieldS26(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS26(name, sms...)
	return s.ActualS(), err
}

// FieldS26 adds a field and reads 26 bit signed integer in current endian
func (d *D) FieldS26(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS26(name, sms...).ActualS()
}

// Reader S27

// TryS27 tries to read 27 bit signed integer in current endian
func (d *D) TryS27() (int64, error) { return d.trySEndian(27, d.Endian) }

// S27 reads 27 bit signed integer in current endian
func (d *D) S27() int64 {
	v, err := d.trySEndian(27, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S27", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS27 tries to add a field and read 27 bit signed integer in current endian
func (d *D) TryFieldScalarS27(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(27, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS27 adds a field and reads 27 bit signed integer in current endian
func (d *D) FieldScalarS27(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS27(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S27", Pos: d.Pos()})
	}
	return s
}

// TryFieldS27 tries to add a field and read 27 bit signed integer in current endian
func (d *D) TryFieldS27(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS27(name, sms...)
	return s.ActualS(), err
}

// FieldS27 adds a field and reads 27 bit signed integer in current endian
func (d *D) FieldS27(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS27(name, sms...).ActualS()
}

// Reader S28

// TryS28 tries to read 28 bit signed integer in current endian
func (d *D) TryS28() (int64, error) { return d.trySEndian(28, d.Endian) }

// S28 reads 28 bit signed integer in current endian
func (d *D) S28() int64 {
	v, err := d.trySEndian(28, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S28", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS28 tries to add a field and read 28 bit signed integer in current endian
func (d *D) TryFieldScalarS28(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(28, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS28 adds a field and reads 28 bit signed integer in current endian
func (d *D) FieldScalarS28(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS28(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S28", Pos: d.Pos()})
	}
	return s
}

// TryFieldS28 tries to add a field and read 28 bit signed integer in current endian
func (d *D) TryFieldS28(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS28(name, sms...)
	return s.ActualS(), err
}

// FieldS28 adds a field and reads 28 bit signed integer in current endian
func (d *D) FieldS28(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS28(name, sms...).ActualS()
}

// Reader S29

// TryS29 tries to read 29 bit signed integer in current endian
func (d *D) TryS29() (int64, error) { return d.trySEndian(29, d.Endian) }

// S29 reads 29 bit signed integer in current endian
func (d *D) S29() int64 {
	v, err := d.trySEndian(29, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S29", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS29 tries to add a field and read 29 bit signed integer in current endian
func (d *D) TryFieldScalarS29(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(29, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS29 adds a field and reads 29 bit signed integer in current endian
func (d *D) FieldScalarS29(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS29(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S29", Pos: d.Pos()})
	}
	return s
}

// TryFieldS29 tries to add a field and read 29 bit signed integer in current endian
func (d *D) TryFieldS29(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS29(name, sms...)
	return s.ActualS(), err
}

// FieldS29 adds a field and reads 29 bit signed integer in current endian
func (d *D) FieldS29(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS29(name, sms...).ActualS()
}

// Reader S30

// TryS30 tries to read 30 bit signed integer in current endian
func (d *D) TryS30() (int64, error) { return d.trySEndian(30, d.Endian) }

// S30 reads 30 bit signed integer in current endian
func (d *D) S30() int64 {
	v, err := d.trySEndian(30, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S30", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS30 tries to add a field and read 30 bit signed integer in current endian
func (d *D) TryFieldScalarS30(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(30, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS30 adds a field and reads 30 bit signed integer in current endian
func (d *D) FieldScalarS30(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS30(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S30", Pos: d.Pos()})
	}
	return s
}

// TryFieldS30 tries to add a field and read 30 bit signed integer in current endian
func (d *D) TryFieldS30(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS30(name, sms...)
	return s.ActualS(), err
}

// FieldS30 adds a field and reads 30 bit signed integer in current endian
func (d *D) FieldS30(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS30(name, sms...).ActualS()
}

// Reader S31

// TryS31 tries to read 31 bit signed integer in current endian
func (d *D) TryS31() (int64, error) { return d.trySEndian(31, d.Endian) }

// S31 reads 31 bit signed integer in current endian
func (d *D) S31() int64 {
	v, err := d.trySEndian(31, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S31", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS31 tries to add a field and read 31 bit signed integer in current endian
func (d *D) TryFieldScalarS31(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(31, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS31 adds a field and reads 31 bit signed integer in current endian
func (d *D) FieldScalarS31(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS31(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S31", Pos: d.Pos()})
	}
	return s
}

// TryFieldS31 tries to add a field and read 31 bit signed integer in current endian
func (d *D) TryFieldS31(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS31(name, sms...)
	return s.ActualS(), err
}

// FieldS31 adds a field and reads 31 bit signed integer in current endian
func (d *D) FieldS31(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS31(name, sms...).ActualS()
}

// Reader S32

// TryS32 tries to read 32 bit signed integer in current endian
func (d *D) TryS32() (int64, error) { return d.trySEndian(32, d.Endian) }

// S32 reads 32 bit signed integer in current endian
func (d *D) S32() int64 {
	v, err := d.trySEndian(32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S32", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS32 tries to add a field and read 32 bit signed integer in current endian
func (d *D) TryFieldScalarS32(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(32, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS32 adds a field and reads 32 bit signed integer in current endian
func (d *D) FieldScalarS32(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS32(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S32", Pos: d.Pos()})
	}
	return s
}

// TryFieldS32 tries to add a field and read 32 bit signed integer in current endian
func (d *D) TryFieldS32(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS32(name, sms...)
	return s.ActualS(), err
}

// FieldS32 adds a field and reads 32 bit signed integer in current endian
func (d *D) FieldS32(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS32(name, sms...).ActualS()
}

// Reader S33

// TryS33 tries to read 33 bit signed integer in current endian
func (d *D) TryS33() (int64, error) { return d.trySEndian(33, d.Endian) }

// S33 reads 33 bit signed integer in current endian
func (d *D) S33() int64 {
	v, err := d.trySEndian(33, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S33", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS33 tries to add a field and read 33 bit signed integer in current endian
func (d *D) TryFieldScalarS33(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(33, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS33 adds a field and reads 33 bit signed integer in current endian
func (d *D) FieldScalarS33(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS33(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S33", Pos: d.Pos()})
	}
	return s
}

// TryFieldS33 tries to add a field and read 33 bit signed integer in current endian
func (d *D) TryFieldS33(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS33(name, sms...)
	return s.ActualS(), err
}

// FieldS33 adds a field and reads 33 bit signed integer in current endian
func (d *D) FieldS33(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS33(name, sms...).ActualS()
}

// Reader S34

// TryS34 tries to read 34 bit signed integer in current endian
func (d *D) TryS34() (int64, error) { return d.trySEndian(34, d.Endian) }

// S34 reads 34 bit signed integer in current endian
func (d *D) S34() int64 {
	v, err := d.trySEndian(34, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S34", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS34 tries to add a field and read 34 bit signed integer in current endian
func (d *D) TryFieldScalarS34(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(34, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS34 adds a field and reads 34 bit signed integer in current endian
func (d *D) FieldScalarS34(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS34(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S34", Pos: d.Pos()})
	}
	return s
}

// TryFieldS34 tries to add a field and read 34 bit signed integer in current endian
func (d *D) TryFieldS34(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS34(name, sms...)
	return s.ActualS(), err
}

// FieldS34 adds a field and reads 34 bit signed integer in current endian
func (d *D) FieldS34(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS34(name, sms...).ActualS()
}

// Reader S35

// TryS35 tries to read 35 bit signed integer in current endian
func (d *D) TryS35() (int64, error) { return d.trySEndian(35, d.Endian) }

// S35 reads 35 bit signed integer in current endian
func (d *D) S35() int64 {
	v, err := d.trySEndian(35, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S35", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS35 tries to add a field and read 35 bit signed integer in current endian
func (d *D) TryFieldScalarS35(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(35, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS35 adds a field and reads 35 bit signed integer in current endian
func (d *D) FieldScalarS35(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS35(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S35", Pos: d.Pos()})
	}
	return s
}

// TryFieldS35 tries to add a field and read 35 bit signed integer in current endian
func (d *D) TryFieldS35(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS35(name, sms...)
	return s.ActualS(), err
}

// FieldS35 adds a field and reads 35 bit signed integer in current endian
func (d *D) FieldS35(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS35(name, sms...).ActualS()
}

// Reader S36

// TryS36 tries to read 36 bit signed integer in current endian
func (d *D) TryS36() (int64, error) { return d.trySEndian(36, d.Endian) }

// S36 reads 36 bit signed integer in current endian
func (d *D) S36() int64 {
	v, err := d.trySEndian(36, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S36", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS36 tries to add a field and read 36 bit signed integer in current endian
func (d *D) TryFieldScalarS36(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(36, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS36 adds a field and reads 36 bit signed integer in current endian
func (d *D) FieldScalarS36(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS36(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S36", Pos: d.Pos()})
	}
	return s
}

// TryFieldS36 tries to add a field and read 36 bit signed integer in current endian
func (d *D) TryFieldS36(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS36(name, sms...)
	return s.ActualS(), err
}

// FieldS36 adds a field and reads 36 bit signed integer in current endian
func (d *D) FieldS36(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS36(name, sms...).ActualS()
}

// Reader S37

// TryS37 tries to read 37 bit signed integer in current endian
func (d *D) TryS37() (int64, error) { return d.trySEndian(37, d.Endian) }

// S37 reads 37 bit signed integer in current endian
func (d *D) S37() int64 {
	v, err := d.trySEndian(37, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S37", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS37 tries to add a field and read 37 bit signed integer in current endian
func (d *D) TryFieldScalarS37(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(37, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS37 adds a field and reads 37 bit signed integer in current endian
func (d *D) FieldScalarS37(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS37(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S37", Pos: d.Pos()})
	}
	return s
}

// TryFieldS37 tries to add a field and read 37 bit signed integer in current endian
func (d *D) TryFieldS37(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS37(name, sms...)
	return s.ActualS(), err
}

// FieldS37 adds a field and reads 37 bit signed integer in current endian
func (d *D) FieldS37(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS37(name, sms...).ActualS()
}

// Reader S38

// TryS38 tries to read 38 bit signed integer in current endian
func (d *D) TryS38() (int64, error) { return d.trySEndian(38, d.Endian) }

// S38 reads 38 bit signed integer in current endian
func (d *D) S38() int64 {
	v, err := d.trySEndian(38, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S38", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS38 tries to add a field and read 38 bit signed integer in current endian
func (d *D) TryFieldScalarS38(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(38, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS38 adds a field and reads 38 bit signed integer in current endian
func (d *D) FieldScalarS38(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS38(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S38", Pos: d.Pos()})
	}
	return s
}

// TryFieldS38 tries to add a field and read 38 bit signed integer in current endian
func (d *D) TryFieldS38(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS38(name, sms...)
	return s.ActualS(), err
}

// FieldS38 adds a field and reads 38 bit signed integer in current endian
func (d *D) FieldS38(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS38(name, sms...).ActualS()
}

// Reader S39

// TryS39 tries to read 39 bit signed integer in current endian
func (d *D) TryS39() (int64, error) { return d.trySEndian(39, d.Endian) }

// S39 reads 39 bit signed integer in current endian
func (d *D) S39() int64 {
	v, err := d.trySEndian(39, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S39", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS39 tries to add a field and read 39 bit signed integer in current endian
func (d *D) TryFieldScalarS39(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(39, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS39 adds a field and reads 39 bit signed integer in current endian
func (d *D) FieldScalarS39(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS39(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S39", Pos: d.Pos()})
	}
	return s
}

// TryFieldS39 tries to add a field and read 39 bit signed integer in current endian
func (d *D) TryFieldS39(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS39(name, sms...)
	return s.ActualS(), err
}

// FieldS39 adds a field and reads 39 bit signed integer in current endian
func (d *D) FieldS39(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS39(name, sms...).ActualS()
}

// Reader S40

// TryS40 tries to read 40 bit signed integer in current endian
func (d *D) TryS40() (int64, error) { return d.trySEndian(40, d.Endian) }

// S40 reads 40 bit signed integer in current endian
func (d *D) S40() int64 {
	v, err := d.trySEndian(40, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S40", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS40 tries to add a field and read 40 bit signed integer in current endian
func (d *D) TryFieldScalarS40(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(40, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS40 adds a field and reads 40 bit signed integer in current endian
func (d *D) FieldScalarS40(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS40(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S40", Pos: d.Pos()})
	}
	return s
}

// TryFieldS40 tries to add a field and read 40 bit signed integer in current endian
func (d *D) TryFieldS40(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS40(name, sms...)
	return s.ActualS(), err
}

// FieldS40 adds a field and reads 40 bit signed integer in current endian
func (d *D) FieldS40(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS40(name, sms...).ActualS()
}

// Reader S41

// TryS41 tries to read 41 bit signed integer in current endian
func (d *D) TryS41() (int64, error) { return d.trySEndian(41, d.Endian) }

// S41 reads 41 bit signed integer in current endian
func (d *D) S41() int64 {
	v, err := d.trySEndian(41, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S41", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS41 tries to add a field and read 41 bit signed integer in current endian
func (d *D) TryFieldScalarS41(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(41, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS41 adds a field and reads 41 bit signed integer in current endian
func (d *D) FieldScalarS41(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS41(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S41", Pos: d.Pos()})
	}
	return s
}

// TryFieldS41 tries to add a field and read 41 bit signed integer in current endian
func (d *D) TryFieldS41(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS41(name, sms...)
	return s.ActualS(), err
}

// FieldS41 adds a field and reads 41 bit signed integer in current endian
func (d *D) FieldS41(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS41(name, sms...).ActualS()
}

// Reader S42

// TryS42 tries to read 42 bit signed integer in current endian
func (d *D) TryS42() (int64, error) { return d.trySEndian(42, d.Endian) }

// S42 reads 42 bit signed integer in current endian
func (d *D) S42() int64 {
	v, err := d.trySEndian(42, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S42", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS42 tries to add a field and read 42 bit signed integer in current endian
func (d *D) TryFieldScalarS42(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(42, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS42 adds a field and reads 42 bit signed integer in current endian
func (d *D) FieldScalarS42(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS42(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S42", Pos: d.Pos()})
	}
	return s
}

// TryFieldS42 tries to add a field and read 42 bit signed integer in current endian
func (d *D) TryFieldS42(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS42(name, sms...)
	return s.ActualS(), err
}

// FieldS42 adds a field and reads 42 bit signed integer in current endian
func (d *D) FieldS42(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS42(name, sms...).ActualS()
}

// Reader S43

// TryS43 tries to read 43 bit signed integer in current endian
func (d *D) TryS43() (int64, error) { return d.trySEndian(43, d.Endian) }

// S43 reads 43 bit signed integer in current endian
func (d *D) S43() int64 {
	v, err := d.trySEndian(43, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S43", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS43 tries to add a field and read 43 bit signed integer in current endian
func (d *D) TryFieldScalarS43(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(43, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS43 adds a field and reads 43 bit signed integer in current endian
func (d *D) FieldScalarS43(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS43(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S43", Pos: d.Pos()})
	}
	return s
}

// TryFieldS43 tries to add a field and read 43 bit signed integer in current endian
func (d *D) TryFieldS43(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS43(name, sms...)
	return s.ActualS(), err
}

// FieldS43 adds a field and reads 43 bit signed integer in current endian
func (d *D) FieldS43(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS43(name, sms...).ActualS()
}

// Reader S44

// TryS44 tries to read 44 bit signed integer in current endian
func (d *D) TryS44() (int64, error) { return d.trySEndian(44, d.Endian) }

// S44 reads 44 bit signed integer in current endian
func (d *D) S44() int64 {
	v, err := d.trySEndian(44, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S44", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS44 tries to add a field and read 44 bit signed integer in current endian
func (d *D) TryFieldScalarS44(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(44, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS44 adds a field and reads 44 bit signed integer in current endian
func (d *D) FieldScalarS44(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS44(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S44", Pos: d.Pos()})
	}
	return s
}

// TryFieldS44 tries to add a field and read 44 bit signed integer in current endian
func (d *D) TryFieldS44(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS44(name, sms...)
	return s.ActualS(), err
}

// FieldS44 adds a field and reads 44 bit signed integer in current endian
func (d *D) FieldS44(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS44(name, sms...).ActualS()
}

// Reader S45

// TryS45 tries to read 45 bit signed integer in current endian
func (d *D) TryS45() (int64, error) { return d.trySEndian(45, d.Endian) }

// S45 reads 45 bit signed integer in current endian
func (d *D) S45() int64 {
	v, err := d.trySEndian(45, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S45", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS45 tries to add a field and read 45 bit signed integer in current endian
func (d *D) TryFieldScalarS45(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(45, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS45 adds a field and reads 45 bit signed integer in current endian
func (d *D) FieldScalarS45(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS45(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S45", Pos: d.Pos()})
	}
	return s
}

// TryFieldS45 tries to add a field and read 45 bit signed integer in current endian
func (d *D) TryFieldS45(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS45(name, sms...)
	return s.ActualS(), err
}

// FieldS45 adds a field and reads 45 bit signed integer in current endian
func (d *D) FieldS45(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS45(name, sms...).ActualS()
}

// Reader S46

// TryS46 tries to read 46 bit signed integer in current endian
func (d *D) TryS46() (int64, error) { return d.trySEndian(46, d.Endian) }

// S46 reads 46 bit signed integer in current endian
func (d *D) S46() int64 {
	v, err := d.trySEndian(46, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S46", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS46 tries to add a field and read 46 bit signed integer in current endian
func (d *D) TryFieldScalarS46(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(46, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS46 adds a field and reads 46 bit signed integer in current endian
func (d *D) FieldScalarS46(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS46(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S46", Pos: d.Pos()})
	}
	return s
}

// TryFieldS46 tries to add a field and read 46 bit signed integer in current endian
func (d *D) TryFieldS46(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS46(name, sms...)
	return s.ActualS(), err
}

// FieldS46 adds a field and reads 46 bit signed integer in current endian
func (d *D) FieldS46(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS46(name, sms...).ActualS()
}

// Reader S47

// TryS47 tries to read 47 bit signed integer in current endian
func (d *D) TryS47() (int64, error) { return d.trySEndian(47, d.Endian) }

// S47 reads 47 bit signed integer in current endian
func (d *D) S47() int64 {
	v, err := d.trySEndian(47, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S47", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS47 tries to add a field and read 47 bit signed integer in current endian
func (d *D) TryFieldScalarS47(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(47, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS47 adds a field and reads 47 bit signed integer in current endian
func (d *D) FieldScalarS47(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS47(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S47", Pos: d.Pos()})
	}
	return s
}

// TryFieldS47 tries to add a field and read 47 bit signed integer in current endian
func (d *D) TryFieldS47(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS47(name, sms...)
	return s.ActualS(), err
}

// FieldS47 adds a field and reads 47 bit signed integer in current endian
func (d *D) FieldS47(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS47(name, sms...).ActualS()
}

// Reader S48

// TryS48 tries to read 48 bit signed integer in current endian
func (d *D) TryS48() (int64, error) { return d.trySEndian(48, d.Endian) }

// S48 reads 48 bit signed integer in current endian
func (d *D) S48() int64 {
	v, err := d.trySEndian(48, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S48", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS48 tries to add a field and read 48 bit signed integer in current endian
func (d *D) TryFieldScalarS48(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(48, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS48 adds a field and reads 48 bit signed integer in current endian
func (d *D) FieldScalarS48(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS48(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S48", Pos: d.Pos()})
	}
	return s
}

// TryFieldS48 tries to add a field and read 48 bit signed integer in current endian
func (d *D) TryFieldS48(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS48(name, sms...)
	return s.ActualS(), err
}

// FieldS48 adds a field and reads 48 bit signed integer in current endian
func (d *D) FieldS48(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS48(name, sms...).ActualS()
}

// Reader S49

// TryS49 tries to read 49 bit signed integer in current endian
func (d *D) TryS49() (int64, error) { return d.trySEndian(49, d.Endian) }

// S49 reads 49 bit signed integer in current endian
func (d *D) S49() int64 {
	v, err := d.trySEndian(49, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S49", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS49 tries to add a field and read 49 bit signed integer in current endian
func (d *D) TryFieldScalarS49(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(49, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS49 adds a field and reads 49 bit signed integer in current endian
func (d *D) FieldScalarS49(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS49(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S49", Pos: d.Pos()})
	}
	return s
}

// TryFieldS49 tries to add a field and read 49 bit signed integer in current endian
func (d *D) TryFieldS49(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS49(name, sms...)
	return s.ActualS(), err
}

// FieldS49 adds a field and reads 49 bit signed integer in current endian
func (d *D) FieldS49(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS49(name, sms...).ActualS()
}

// Reader S50

// TryS50 tries to read 50 bit signed integer in current endian
func (d *D) TryS50() (int64, error) { return d.trySEndian(50, d.Endian) }

// S50 reads 50 bit signed integer in current endian
func (d *D) S50() int64 {
	v, err := d.trySEndian(50, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S50", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS50 tries to add a field and read 50 bit signed integer in current endian
func (d *D) TryFieldScalarS50(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(50, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS50 adds a field and reads 50 bit signed integer in current endian
func (d *D) FieldScalarS50(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS50(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S50", Pos: d.Pos()})
	}
	return s
}

// TryFieldS50 tries to add a field and read 50 bit signed integer in current endian
func (d *D) TryFieldS50(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS50(name, sms...)
	return s.ActualS(), err
}

// FieldS50 adds a field and reads 50 bit signed integer in current endian
func (d *D) FieldS50(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS50(name, sms...).ActualS()
}

// Reader S51

// TryS51 tries to read 51 bit signed integer in current endian
func (d *D) TryS51() (int64, error) { return d.trySEndian(51, d.Endian) }

// S51 reads 51 bit signed integer in current endian
func (d *D) S51() int64 {
	v, err := d.trySEndian(51, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S51", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS51 tries to add a field and read 51 bit signed integer in current endian
func (d *D) TryFieldScalarS51(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(51, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS51 adds a field and reads 51 bit signed integer in current endian
func (d *D) FieldScalarS51(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS51(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S51", Pos: d.Pos()})
	}
	return s
}

// TryFieldS51 tries to add a field and read 51 bit signed integer in current endian
func (d *D) TryFieldS51(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS51(name, sms...)
	return s.ActualS(), err
}

// FieldS51 adds a field and reads 51 bit signed integer in current endian
func (d *D) FieldS51(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS51(name, sms...).ActualS()
}

// Reader S52

// TryS52 tries to read 52 bit signed integer in current endian
func (d *D) TryS52() (int64, error) { return d.trySEndian(52, d.Endian) }

// S52 reads 52 bit signed integer in current endian
func (d *D) S52() int64 {
	v, err := d.trySEndian(52, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S52", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS52 tries to add a field and read 52 bit signed integer in current endian
func (d *D) TryFieldScalarS52(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(52, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS52 adds a field and reads 52 bit signed integer in current endian
func (d *D) FieldScalarS52(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS52(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S52", Pos: d.Pos()})
	}
	return s
}

// TryFieldS52 tries to add a field and read 52 bit signed integer in current endian
func (d *D) TryFieldS52(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS52(name, sms...)
	return s.ActualS(), err
}

// FieldS52 adds a field and reads 52 bit signed integer in current endian
func (d *D) FieldS52(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS52(name, sms...).ActualS()
}

// Reader S53

// TryS53 tries to read 53 bit signed integer in current endian
func (d *D) TryS53() (int64, error) { return d.trySEndian(53, d.Endian) }

// S53 reads 53 bit signed integer in current endian
func (d *D) S53() int64 {
	v, err := d.trySEndian(53, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S53", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS53 tries to add a field and read 53 bit signed integer in current endian
func (d *D) TryFieldScalarS53(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(53, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS53 adds a field and reads 53 bit signed integer in current endian
func (d *D) FieldScalarS53(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS53(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S53", Pos: d.Pos()})
	}
	return s
}

// TryFieldS53 tries to add a field and read 53 bit signed integer in current endian
func (d *D) TryFieldS53(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS53(name, sms...)
	return s.ActualS(), err
}

// FieldS53 adds a field and reads 53 bit signed integer in current endian
func (d *D) FieldS53(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS53(name, sms...).ActualS()
}

// Reader S54

// TryS54 tries to read 54 bit signed integer in current endian
func (d *D) TryS54() (int64, error) { return d.trySEndian(54, d.Endian) }

// S54 reads 54 bit signed integer in current endian
func (d *D) S54() int64 {
	v, err := d.trySEndian(54, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S54", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS54 tries to add a field and read 54 bit signed integer in current endian
func (d *D) TryFieldScalarS54(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(54, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS54 adds a field and reads 54 bit signed integer in current endian
func (d *D) FieldScalarS54(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS54(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S54", Pos: d.Pos()})
	}
	return s
}

// TryFieldS54 tries to add a field and read 54 bit signed integer in current endian
func (d *D) TryFieldS54(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS54(name, sms...)
	return s.ActualS(), err
}

// FieldS54 adds a field and reads 54 bit signed integer in current endian
func (d *D) FieldS54(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS54(name, sms...).ActualS()
}

// Reader S55

// TryS55 tries to read 55 bit signed integer in current endian
func (d *D) TryS55() (int64, error) { return d.trySEndian(55, d.Endian) }

// S55 reads 55 bit signed integer in current endian
func (d *D) S55() int64 {
	v, err := d.trySEndian(55, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S55", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS55 tries to add a field and read 55 bit signed integer in current endian
func (d *D) TryFieldScalarS55(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(55, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS55 adds a field and reads 55 bit signed integer in current endian
func (d *D) FieldScalarS55(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS55(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S55", Pos: d.Pos()})
	}
	return s
}

// TryFieldS55 tries to add a field and read 55 bit signed integer in current endian
func (d *D) TryFieldS55(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS55(name, sms...)
	return s.ActualS(), err
}

// FieldS55 adds a field and reads 55 bit signed integer in current endian
func (d *D) FieldS55(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS55(name, sms...).ActualS()
}

// Reader S56

// TryS56 tries to read 56 bit signed integer in current endian
func (d *D) TryS56() (int64, error) { return d.trySEndian(56, d.Endian) }

// S56 reads 56 bit signed integer in current endian
func (d *D) S56() int64 {
	v, err := d.trySEndian(56, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S56", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS56 tries to add a field and read 56 bit signed integer in current endian
func (d *D) TryFieldScalarS56(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(56, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS56 adds a field and reads 56 bit signed integer in current endian
func (d *D) FieldScalarS56(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS56(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S56", Pos: d.Pos()})
	}
	return s
}

// TryFieldS56 tries to add a field and read 56 bit signed integer in current endian
func (d *D) TryFieldS56(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS56(name, sms...)
	return s.ActualS(), err
}

// FieldS56 adds a field and reads 56 bit signed integer in current endian
func (d *D) FieldS56(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS56(name, sms...).ActualS()
}

// Reader S57

// TryS57 tries to read 57 bit signed integer in current endian
func (d *D) TryS57() (int64, error) { return d.trySEndian(57, d.Endian) }

// S57 reads 57 bit signed integer in current endian
func (d *D) S57() int64 {
	v, err := d.trySEndian(57, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S57", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS57 tries to add a field and read 57 bit signed integer in current endian
func (d *D) TryFieldScalarS57(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(57, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS57 adds a field and reads 57 bit signed integer in current endian
func (d *D) FieldScalarS57(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS57(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S57", Pos: d.Pos()})
	}
	return s
}

// TryFieldS57 tries to add a field and read 57 bit signed integer in current endian
func (d *D) TryFieldS57(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS57(name, sms...)
	return s.ActualS(), err
}

// FieldS57 adds a field and reads 57 bit signed integer in current endian
func (d *D) FieldS57(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS57(name, sms...).ActualS()
}

// Reader S58

// TryS58 tries to read 58 bit signed integer in current endian
func (d *D) TryS58() (int64, error) { return d.trySEndian(58, d.Endian) }

// S58 reads 58 bit signed integer in current endian
func (d *D) S58() int64 {
	v, err := d.trySEndian(58, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S58", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS58 tries to add a field and read 58 bit signed integer in current endian
func (d *D) TryFieldScalarS58(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(58, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS58 adds a field and reads 58 bit signed integer in current endian
func (d *D) FieldScalarS58(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS58(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S58", Pos: d.Pos()})
	}
	return s
}

// TryFieldS58 tries to add a field and read 58 bit signed integer in current endian
func (d *D) TryFieldS58(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS58(name, sms...)
	return s.ActualS(), err
}

// FieldS58 adds a field and reads 58 bit signed integer in current endian
func (d *D) FieldS58(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS58(name, sms...).ActualS()
}

// Reader S59

// TryS59 tries to read 59 bit signed integer in current endian
func (d *D) TryS59() (int64, error) { return d.trySEndian(59, d.Endian) }

// S59 reads 59 bit signed integer in current endian
func (d *D) S59() int64 {
	v, err := d.trySEndian(59, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S59", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS59 tries to add a field and read 59 bit signed integer in current endian
func (d *D) TryFieldScalarS59(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(59, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS59 adds a field and reads 59 bit signed integer in current endian
func (d *D) FieldScalarS59(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS59(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S59", Pos: d.Pos()})
	}
	return s
}

// TryFieldS59 tries to add a field and read 59 bit signed integer in current endian
func (d *D) TryFieldS59(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS59(name, sms...)
	return s.ActualS(), err
}

// FieldS59 adds a field and reads 59 bit signed integer in current endian
func (d *D) FieldS59(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS59(name, sms...).ActualS()
}

// Reader S60

// TryS60 tries to read 60 bit signed integer in current endian
func (d *D) TryS60() (int64, error) { return d.trySEndian(60, d.Endian) }

// S60 reads 60 bit signed integer in current endian
func (d *D) S60() int64 {
	v, err := d.trySEndian(60, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S60", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS60 tries to add a field and read 60 bit signed integer in current endian
func (d *D) TryFieldScalarS60(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(60, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS60 adds a field and reads 60 bit signed integer in current endian
func (d *D) FieldScalarS60(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS60(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S60", Pos: d.Pos()})
	}
	return s
}

// TryFieldS60 tries to add a field and read 60 bit signed integer in current endian
func (d *D) TryFieldS60(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS60(name, sms...)
	return s.ActualS(), err
}

// FieldS60 adds a field and reads 60 bit signed integer in current endian
func (d *D) FieldS60(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS60(name, sms...).ActualS()
}

// Reader S61

// TryS61 tries to read 61 bit signed integer in current endian
func (d *D) TryS61() (int64, error) { return d.trySEndian(61, d.Endian) }

// S61 reads 61 bit signed integer in current endian
func (d *D) S61() int64 {
	v, err := d.trySEndian(61, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S61", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS61 tries to add a field and read 61 bit signed integer in current endian
func (d *D) TryFieldScalarS61(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(61, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS61 adds a field and reads 61 bit signed integer in current endian
func (d *D) FieldScalarS61(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS61(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S61", Pos: d.Pos()})
	}
	return s
}

// TryFieldS61 tries to add a field and read 61 bit signed integer in current endian
func (d *D) TryFieldS61(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS61(name, sms...)
	return s.ActualS(), err
}

// FieldS61 adds a field and reads 61 bit signed integer in current endian
func (d *D) FieldS61(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS61(name, sms...).ActualS()
}

// Reader S62

// TryS62 tries to read 62 bit signed integer in current endian
func (d *D) TryS62() (int64, error) { return d.trySEndian(62, d.Endian) }

// S62 reads 62 bit signed integer in current endian
func (d *D) S62() int64 {
	v, err := d.trySEndian(62, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S62", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS62 tries to add a field and read 62 bit signed integer in current endian
func (d *D) TryFieldScalarS62(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(62, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS62 adds a field and reads 62 bit signed integer in current endian
func (d *D) FieldScalarS62(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS62(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S62", Pos: d.Pos()})
	}
	return s
}

// TryFieldS62 tries to add a field and read 62 bit signed integer in current endian
func (d *D) TryFieldS62(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS62(name, sms...)
	return s.ActualS(), err
}

// FieldS62 adds a field and reads 62 bit signed integer in current endian
func (d *D) FieldS62(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS62(name, sms...).ActualS()
}

// Reader S63

// TryS63 tries to read 63 bit signed integer in current endian
func (d *D) TryS63() (int64, error) { return d.trySEndian(63, d.Endian) }

// S63 reads 63 bit signed integer in current endian
func (d *D) S63() int64 {
	v, err := d.trySEndian(63, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S63", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS63 tries to add a field and read 63 bit signed integer in current endian
func (d *D) TryFieldScalarS63(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(63, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS63 adds a field and reads 63 bit signed integer in current endian
func (d *D) FieldScalarS63(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS63(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S63", Pos: d.Pos()})
	}
	return s
}

// TryFieldS63 tries to add a field and read 63 bit signed integer in current endian
func (d *D) TryFieldS63(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS63(name, sms...)
	return s.ActualS(), err
}

// FieldS63 adds a field and reads 63 bit signed integer in current endian
func (d *D) FieldS63(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS63(name, sms...).ActualS()
}

// Reader S64

// TryS64 tries to read 64 bit signed integer in current endian
func (d *D) TryS64() (int64, error) { return d.trySEndian(64, d.Endian) }

// S64 reads 64 bit signed integer in current endian
func (d *D) S64() int64 {
	v, err := d.trySEndian(64, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S64", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS64 tries to add a field and read 64 bit signed integer in current endian
func (d *D) TryFieldScalarS64(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(64, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS64 adds a field and reads 64 bit signed integer in current endian
func (d *D) FieldScalarS64(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS64(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S64", Pos: d.Pos()})
	}
	return s
}

// TryFieldS64 tries to add a field and read 64 bit signed integer in current endian
func (d *D) TryFieldS64(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS64(name, sms...)
	return s.ActualS(), err
}

// FieldS64 adds a field and reads 64 bit signed integer in current endian
func (d *D) FieldS64(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS64(name, sms...).ActualS()
}

// Reader S8LE

// TryS8LE tries to read 8 bit signed integer in little-endian
func (d *D) TryS8LE() (int64, error) { return d.trySEndian(8, LittleEndian) }

// S8LE reads 8 bit signed integer in little-endian
func (d *D) S8LE() int64 {
	v, err := d.trySEndian(8, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S8LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS8LE tries to add a field and read 8 bit signed integer in little-endian
func (d *D) TryFieldScalarS8LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(8, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS8LE adds a field and reads 8 bit signed integer in little-endian
func (d *D) FieldScalarS8LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS8LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S8LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS8LE tries to add a field and read 8 bit signed integer in little-endian
func (d *D) TryFieldS8LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS8LE(name, sms...)
	return s.ActualS(), err
}

// FieldS8LE adds a field and reads 8 bit signed integer in little-endian
func (d *D) FieldS8LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS8LE(name, sms...).ActualS()
}

// Reader S9LE

// TryS9LE tries to read 9 bit signed integer in little-endian
func (d *D) TryS9LE() (int64, error) { return d.trySEndian(9, LittleEndian) }

// S9LE reads 9 bit signed integer in little-endian
func (d *D) S9LE() int64 {
	v, err := d.trySEndian(9, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S9LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS9LE tries to add a field and read 9 bit signed integer in little-endian
func (d *D) TryFieldScalarS9LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(9, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS9LE adds a field and reads 9 bit signed integer in little-endian
func (d *D) FieldScalarS9LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS9LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S9LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS9LE tries to add a field and read 9 bit signed integer in little-endian
func (d *D) TryFieldS9LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS9LE(name, sms...)
	return s.ActualS(), err
}

// FieldS9LE adds a field and reads 9 bit signed integer in little-endian
func (d *D) FieldS9LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS9LE(name, sms...).ActualS()
}

// Reader S10LE

// TryS10LE tries to read 10 bit signed integer in little-endian
func (d *D) TryS10LE() (int64, error) { return d.trySEndian(10, LittleEndian) }

// S10LE reads 10 bit signed integer in little-endian
func (d *D) S10LE() int64 {
	v, err := d.trySEndian(10, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S10LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS10LE tries to add a field and read 10 bit signed integer in little-endian
func (d *D) TryFieldScalarS10LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(10, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS10LE adds a field and reads 10 bit signed integer in little-endian
func (d *D) FieldScalarS10LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS10LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S10LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS10LE tries to add a field and read 10 bit signed integer in little-endian
func (d *D) TryFieldS10LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS10LE(name, sms...)
	return s.ActualS(), err
}

// FieldS10LE adds a field and reads 10 bit signed integer in little-endian
func (d *D) FieldS10LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS10LE(name, sms...).ActualS()
}

// Reader S11LE

// TryS11LE tries to read 11 bit signed integer in little-endian
func (d *D) TryS11LE() (int64, error) { return d.trySEndian(11, LittleEndian) }

// S11LE reads 11 bit signed integer in little-endian
func (d *D) S11LE() int64 {
	v, err := d.trySEndian(11, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S11LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS11LE tries to add a field and read 11 bit signed integer in little-endian
func (d *D) TryFieldScalarS11LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(11, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS11LE adds a field and reads 11 bit signed integer in little-endian
func (d *D) FieldScalarS11LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS11LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S11LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS11LE tries to add a field and read 11 bit signed integer in little-endian
func (d *D) TryFieldS11LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS11LE(name, sms...)
	return s.ActualS(), err
}

// FieldS11LE adds a field and reads 11 bit signed integer in little-endian
func (d *D) FieldS11LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS11LE(name, sms...).ActualS()
}

// Reader S12LE

// TryS12LE tries to read 12 bit signed integer in little-endian
func (d *D) TryS12LE() (int64, error) { return d.trySEndian(12, LittleEndian) }

// S12LE reads 12 bit signed integer in little-endian
func (d *D) S12LE() int64 {
	v, err := d.trySEndian(12, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S12LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS12LE tries to add a field and read 12 bit signed integer in little-endian
func (d *D) TryFieldScalarS12LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(12, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS12LE adds a field and reads 12 bit signed integer in little-endian
func (d *D) FieldScalarS12LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS12LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S12LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS12LE tries to add a field and read 12 bit signed integer in little-endian
func (d *D) TryFieldS12LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS12LE(name, sms...)
	return s.ActualS(), err
}

// FieldS12LE adds a field and reads 12 bit signed integer in little-endian
func (d *D) FieldS12LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS12LE(name, sms...).ActualS()
}

// Reader S13LE

// TryS13LE tries to read 13 bit signed integer in little-endian
func (d *D) TryS13LE() (int64, error) { return d.trySEndian(13, LittleEndian) }

// S13LE reads 13 bit signed integer in little-endian
func (d *D) S13LE() int64 {
	v, err := d.trySEndian(13, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S13LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS13LE tries to add a field and read 13 bit signed integer in little-endian
func (d *D) TryFieldScalarS13LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(13, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS13LE adds a field and reads 13 bit signed integer in little-endian
func (d *D) FieldScalarS13LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS13LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S13LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS13LE tries to add a field and read 13 bit signed integer in little-endian
func (d *D) TryFieldS13LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS13LE(name, sms...)
	return s.ActualS(), err
}

// FieldS13LE adds a field and reads 13 bit signed integer in little-endian
func (d *D) FieldS13LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS13LE(name, sms...).ActualS()
}

// Reader S14LE

// TryS14LE tries to read 14 bit signed integer in little-endian
func (d *D) TryS14LE() (int64, error) { return d.trySEndian(14, LittleEndian) }

// S14LE reads 14 bit signed integer in little-endian
func (d *D) S14LE() int64 {
	v, err := d.trySEndian(14, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S14LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS14LE tries to add a field and read 14 bit signed integer in little-endian
func (d *D) TryFieldScalarS14LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(14, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS14LE adds a field and reads 14 bit signed integer in little-endian
func (d *D) FieldScalarS14LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS14LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S14LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS14LE tries to add a field and read 14 bit signed integer in little-endian
func (d *D) TryFieldS14LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS14LE(name, sms...)
	return s.ActualS(), err
}

// FieldS14LE adds a field and reads 14 bit signed integer in little-endian
func (d *D) FieldS14LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS14LE(name, sms...).ActualS()
}

// Reader S15LE

// TryS15LE tries to read 15 bit signed integer in little-endian
func (d *D) TryS15LE() (int64, error) { return d.trySEndian(15, LittleEndian) }

// S15LE reads 15 bit signed integer in little-endian
func (d *D) S15LE() int64 {
	v, err := d.trySEndian(15, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S15LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS15LE tries to add a field and read 15 bit signed integer in little-endian
func (d *D) TryFieldScalarS15LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(15, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS15LE adds a field and reads 15 bit signed integer in little-endian
func (d *D) FieldScalarS15LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS15LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S15LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS15LE tries to add a field and read 15 bit signed integer in little-endian
func (d *D) TryFieldS15LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS15LE(name, sms...)
	return s.ActualS(), err
}

// FieldS15LE adds a field and reads 15 bit signed integer in little-endian
func (d *D) FieldS15LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS15LE(name, sms...).ActualS()
}

// Reader S16LE

// TryS16LE tries to read 16 bit signed integer in little-endian
func (d *D) TryS16LE() (int64, error) { return d.trySEndian(16, LittleEndian) }

// S16LE reads 16 bit signed integer in little-endian
func (d *D) S16LE() int64 {
	v, err := d.trySEndian(16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S16LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS16LE tries to add a field and read 16 bit signed integer in little-endian
func (d *D) TryFieldScalarS16LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(16, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS16LE adds a field and reads 16 bit signed integer in little-endian
func (d *D) FieldScalarS16LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS16LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S16LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS16LE tries to add a field and read 16 bit signed integer in little-endian
func (d *D) TryFieldS16LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS16LE(name, sms...)
	return s.ActualS(), err
}

// FieldS16LE adds a field and reads 16 bit signed integer in little-endian
func (d *D) FieldS16LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS16LE(name, sms...).ActualS()
}

// Reader S17LE

// TryS17LE tries to read 17 bit signed integer in little-endian
func (d *D) TryS17LE() (int64, error) { return d.trySEndian(17, LittleEndian) }

// S17LE reads 17 bit signed integer in little-endian
func (d *D) S17LE() int64 {
	v, err := d.trySEndian(17, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S17LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS17LE tries to add a field and read 17 bit signed integer in little-endian
func (d *D) TryFieldScalarS17LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(17, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS17LE adds a field and reads 17 bit signed integer in little-endian
func (d *D) FieldScalarS17LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS17LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S17LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS17LE tries to add a field and read 17 bit signed integer in little-endian
func (d *D) TryFieldS17LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS17LE(name, sms...)
	return s.ActualS(), err
}

// FieldS17LE adds a field and reads 17 bit signed integer in little-endian
func (d *D) FieldS17LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS17LE(name, sms...).ActualS()
}

// Reader S18LE

// TryS18LE tries to read 18 bit signed integer in little-endian
func (d *D) TryS18LE() (int64, error) { return d.trySEndian(18, LittleEndian) }

// S18LE reads 18 bit signed integer in little-endian
func (d *D) S18LE() int64 {
	v, err := d.trySEndian(18, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S18LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS18LE tries to add a field and read 18 bit signed integer in little-endian
func (d *D) TryFieldScalarS18LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(18, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS18LE adds a field and reads 18 bit signed integer in little-endian
func (d *D) FieldScalarS18LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS18LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S18LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS18LE tries to add a field and read 18 bit signed integer in little-endian
func (d *D) TryFieldS18LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS18LE(name, sms...)
	return s.ActualS(), err
}

// FieldS18LE adds a field and reads 18 bit signed integer in little-endian
func (d *D) FieldS18LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS18LE(name, sms...).ActualS()
}

// Reader S19LE

// TryS19LE tries to read 19 bit signed integer in little-endian
func (d *D) TryS19LE() (int64, error) { return d.trySEndian(19, LittleEndian) }

// S19LE reads 19 bit signed integer in little-endian
func (d *D) S19LE() int64 {
	v, err := d.trySEndian(19, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S19LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS19LE tries to add a field and read 19 bit signed integer in little-endian
func (d *D) TryFieldScalarS19LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(19, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS19LE adds a field and reads 19 bit signed integer in little-endian
func (d *D) FieldScalarS19LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS19LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S19LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS19LE tries to add a field and read 19 bit signed integer in little-endian
func (d *D) TryFieldS19LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS19LE(name, sms...)
	return s.ActualS(), err
}

// FieldS19LE adds a field and reads 19 bit signed integer in little-endian
func (d *D) FieldS19LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS19LE(name, sms...).ActualS()
}

// Reader S20LE

// TryS20LE tries to read 20 bit signed integer in little-endian
func (d *D) TryS20LE() (int64, error) { return d.trySEndian(20, LittleEndian) }

// S20LE reads 20 bit signed integer in little-endian
func (d *D) S20LE() int64 {
	v, err := d.trySEndian(20, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S20LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS20LE tries to add a field and read 20 bit signed integer in little-endian
func (d *D) TryFieldScalarS20LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(20, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS20LE adds a field and reads 20 bit signed integer in little-endian
func (d *D) FieldScalarS20LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS20LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S20LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS20LE tries to add a field and read 20 bit signed integer in little-endian
func (d *D) TryFieldS20LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS20LE(name, sms...)
	return s.ActualS(), err
}

// FieldS20LE adds a field and reads 20 bit signed integer in little-endian
func (d *D) FieldS20LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS20LE(name, sms...).ActualS()
}

// Reader S21LE

// TryS21LE tries to read 21 bit signed integer in little-endian
func (d *D) TryS21LE() (int64, error) { return d.trySEndian(21, LittleEndian) }

// S21LE reads 21 bit signed integer in little-endian
func (d *D) S21LE() int64 {
	v, err := d.trySEndian(21, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S21LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS21LE tries to add a field and read 21 bit signed integer in little-endian
func (d *D) TryFieldScalarS21LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(21, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS21LE adds a field and reads 21 bit signed integer in little-endian
func (d *D) FieldScalarS21LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS21LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S21LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS21LE tries to add a field and read 21 bit signed integer in little-endian
func (d *D) TryFieldS21LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS21LE(name, sms...)
	return s.ActualS(), err
}

// FieldS21LE adds a field and reads 21 bit signed integer in little-endian
func (d *D) FieldS21LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS21LE(name, sms...).ActualS()
}

// Reader S22LE

// TryS22LE tries to read 22 bit signed integer in little-endian
func (d *D) TryS22LE() (int64, error) { return d.trySEndian(22, LittleEndian) }

// S22LE reads 22 bit signed integer in little-endian
func (d *D) S22LE() int64 {
	v, err := d.trySEndian(22, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S22LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS22LE tries to add a field and read 22 bit signed integer in little-endian
func (d *D) TryFieldScalarS22LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(22, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS22LE adds a field and reads 22 bit signed integer in little-endian
func (d *D) FieldScalarS22LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS22LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S22LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS22LE tries to add a field and read 22 bit signed integer in little-endian
func (d *D) TryFieldS22LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS22LE(name, sms...)
	return s.ActualS(), err
}

// FieldS22LE adds a field and reads 22 bit signed integer in little-endian
func (d *D) FieldS22LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS22LE(name, sms...).ActualS()
}

// Reader S23LE

// TryS23LE tries to read 23 bit signed integer in little-endian
func (d *D) TryS23LE() (int64, error) { return d.trySEndian(23, LittleEndian) }

// S23LE reads 23 bit signed integer in little-endian
func (d *D) S23LE() int64 {
	v, err := d.trySEndian(23, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S23LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS23LE tries to add a field and read 23 bit signed integer in little-endian
func (d *D) TryFieldScalarS23LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(23, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS23LE adds a field and reads 23 bit signed integer in little-endian
func (d *D) FieldScalarS23LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS23LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S23LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS23LE tries to add a field and read 23 bit signed integer in little-endian
func (d *D) TryFieldS23LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS23LE(name, sms...)
	return s.ActualS(), err
}

// FieldS23LE adds a field and reads 23 bit signed integer in little-endian
func (d *D) FieldS23LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS23LE(name, sms...).ActualS()
}

// Reader S24LE

// TryS24LE tries to read 24 bit signed integer in little-endian
func (d *D) TryS24LE() (int64, error) { return d.trySEndian(24, LittleEndian) }

// S24LE reads 24 bit signed integer in little-endian
func (d *D) S24LE() int64 {
	v, err := d.trySEndian(24, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S24LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS24LE tries to add a field and read 24 bit signed integer in little-endian
func (d *D) TryFieldScalarS24LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(24, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS24LE adds a field and reads 24 bit signed integer in little-endian
func (d *D) FieldScalarS24LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS24LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S24LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS24LE tries to add a field and read 24 bit signed integer in little-endian
func (d *D) TryFieldS24LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS24LE(name, sms...)
	return s.ActualS(), err
}

// FieldS24LE adds a field and reads 24 bit signed integer in little-endian
func (d *D) FieldS24LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS24LE(name, sms...).ActualS()
}

// Reader S25LE

// TryS25LE tries to read 25 bit signed integer in little-endian
func (d *D) TryS25LE() (int64, error) { return d.trySEndian(25, LittleEndian) }

// S25LE reads 25 bit signed integer in little-endian
func (d *D) S25LE() int64 {
	v, err := d.trySEndian(25, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S25LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS25LE tries to add a field and read 25 bit signed integer in little-endian
func (d *D) TryFieldScalarS25LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(25, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS25LE adds a field and reads 25 bit signed integer in little-endian
func (d *D) FieldScalarS25LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS25LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S25LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS25LE tries to add a field and read 25 bit signed integer in little-endian
func (d *D) TryFieldS25LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS25LE(name, sms...)
	return s.ActualS(), err
}

// FieldS25LE adds a field and reads 25 bit signed integer in little-endian
func (d *D) FieldS25LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS25LE(name, sms...).ActualS()
}

// Reader S26LE

// TryS26LE tries to read 26 bit signed integer in little-endian
func (d *D) TryS26LE() (int64, error) { return d.trySEndian(26, LittleEndian) }

// S26LE reads 26 bit signed integer in little-endian
func (d *D) S26LE() int64 {
	v, err := d.trySEndian(26, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S26LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS26LE tries to add a field and read 26 bit signed integer in little-endian
func (d *D) TryFieldScalarS26LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(26, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS26LE adds a field and reads 26 bit signed integer in little-endian
func (d *D) FieldScalarS26LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS26LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S26LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS26LE tries to add a field and read 26 bit signed integer in little-endian
func (d *D) TryFieldS26LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS26LE(name, sms...)
	return s.ActualS(), err
}

// FieldS26LE adds a field and reads 26 bit signed integer in little-endian
func (d *D) FieldS26LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS26LE(name, sms...).ActualS()
}

// Reader S27LE

// TryS27LE tries to read 27 bit signed integer in little-endian
func (d *D) TryS27LE() (int64, error) { return d.trySEndian(27, LittleEndian) }

// S27LE reads 27 bit signed integer in little-endian
func (d *D) S27LE() int64 {
	v, err := d.trySEndian(27, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S27LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS27LE tries to add a field and read 27 bit signed integer in little-endian
func (d *D) TryFieldScalarS27LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(27, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS27LE adds a field and reads 27 bit signed integer in little-endian
func (d *D) FieldScalarS27LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS27LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S27LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS27LE tries to add a field and read 27 bit signed integer in little-endian
func (d *D) TryFieldS27LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS27LE(name, sms...)
	return s.ActualS(), err
}

// FieldS27LE adds a field and reads 27 bit signed integer in little-endian
func (d *D) FieldS27LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS27LE(name, sms...).ActualS()
}

// Reader S28LE

// TryS28LE tries to read 28 bit signed integer in little-endian
func (d *D) TryS28LE() (int64, error) { return d.trySEndian(28, LittleEndian) }

// S28LE reads 28 bit signed integer in little-endian
func (d *D) S28LE() int64 {
	v, err := d.trySEndian(28, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S28LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS28LE tries to add a field and read 28 bit signed integer in little-endian
func (d *D) TryFieldScalarS28LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(28, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS28LE adds a field and reads 28 bit signed integer in little-endian
func (d *D) FieldScalarS28LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS28LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S28LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS28LE tries to add a field and read 28 bit signed integer in little-endian
func (d *D) TryFieldS28LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS28LE(name, sms...)
	return s.ActualS(), err
}

// FieldS28LE adds a field and reads 28 bit signed integer in little-endian
func (d *D) FieldS28LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS28LE(name, sms...).ActualS()
}

// Reader S29LE

// TryS29LE tries to read 29 bit signed integer in little-endian
func (d *D) TryS29LE() (int64, error) { return d.trySEndian(29, LittleEndian) }

// S29LE reads 29 bit signed integer in little-endian
func (d *D) S29LE() int64 {
	v, err := d.trySEndian(29, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S29LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS29LE tries to add a field and read 29 bit signed integer in little-endian
func (d *D) TryFieldScalarS29LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(29, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS29LE adds a field and reads 29 bit signed integer in little-endian
func (d *D) FieldScalarS29LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS29LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S29LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS29LE tries to add a field and read 29 bit signed integer in little-endian
func (d *D) TryFieldS29LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS29LE(name, sms...)
	return s.ActualS(), err
}

// FieldS29LE adds a field and reads 29 bit signed integer in little-endian
func (d *D) FieldS29LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS29LE(name, sms...).ActualS()
}

// Reader S30LE

// TryS30LE tries to read 30 bit signed integer in little-endian
func (d *D) TryS30LE() (int64, error) { return d.trySEndian(30, LittleEndian) }

// S30LE reads 30 bit signed integer in little-endian
func (d *D) S30LE() int64 {
	v, err := d.trySEndian(30, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S30LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS30LE tries to add a field and read 30 bit signed integer in little-endian
func (d *D) TryFieldScalarS30LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(30, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS30LE adds a field and reads 30 bit signed integer in little-endian
func (d *D) FieldScalarS30LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS30LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S30LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS30LE tries to add a field and read 30 bit signed integer in little-endian
func (d *D) TryFieldS30LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS30LE(name, sms...)
	return s.ActualS(), err
}

// FieldS30LE adds a field and reads 30 bit signed integer in little-endian
func (d *D) FieldS30LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS30LE(name, sms...).ActualS()
}

// Reader S31LE

// TryS31LE tries to read 31 bit signed integer in little-endian
func (d *D) TryS31LE() (int64, error) { return d.trySEndian(31, LittleEndian) }

// S31LE reads 31 bit signed integer in little-endian
func (d *D) S31LE() int64 {
	v, err := d.trySEndian(31, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S31LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS31LE tries to add a field and read 31 bit signed integer in little-endian
func (d *D) TryFieldScalarS31LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(31, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS31LE adds a field and reads 31 bit signed integer in little-endian
func (d *D) FieldScalarS31LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS31LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S31LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS31LE tries to add a field and read 31 bit signed integer in little-endian
func (d *D) TryFieldS31LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS31LE(name, sms...)
	return s.ActualS(), err
}

// FieldS31LE adds a field and reads 31 bit signed integer in little-endian
func (d *D) FieldS31LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS31LE(name, sms...).ActualS()
}

// Reader S32LE

// TryS32LE tries to read 32 bit signed integer in little-endian
func (d *D) TryS32LE() (int64, error) { return d.trySEndian(32, LittleEndian) }

// S32LE reads 32 bit signed integer in little-endian
func (d *D) S32LE() int64 {
	v, err := d.trySEndian(32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S32LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS32LE tries to add a field and read 32 bit signed integer in little-endian
func (d *D) TryFieldScalarS32LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(32, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS32LE adds a field and reads 32 bit signed integer in little-endian
func (d *D) FieldScalarS32LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS32LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S32LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS32LE tries to add a field and read 32 bit signed integer in little-endian
func (d *D) TryFieldS32LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS32LE(name, sms...)
	return s.ActualS(), err
}

// FieldS32LE adds a field and reads 32 bit signed integer in little-endian
func (d *D) FieldS32LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS32LE(name, sms...).ActualS()
}

// Reader S33LE

// TryS33LE tries to read 33 bit signed integer in little-endian
func (d *D) TryS33LE() (int64, error) { return d.trySEndian(33, LittleEndian) }

// S33LE reads 33 bit signed integer in little-endian
func (d *D) S33LE() int64 {
	v, err := d.trySEndian(33, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S33LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS33LE tries to add a field and read 33 bit signed integer in little-endian
func (d *D) TryFieldScalarS33LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(33, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS33LE adds a field and reads 33 bit signed integer in little-endian
func (d *D) FieldScalarS33LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS33LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S33LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS33LE tries to add a field and read 33 bit signed integer in little-endian
func (d *D) TryFieldS33LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS33LE(name, sms...)
	return s.ActualS(), err
}

// FieldS33LE adds a field and reads 33 bit signed integer in little-endian
func (d *D) FieldS33LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS33LE(name, sms...).ActualS()
}

// Reader S34LE

// TryS34LE tries to read 34 bit signed integer in little-endian
func (d *D) TryS34LE() (int64, error) { return d.trySEndian(34, LittleEndian) }

// S34LE reads 34 bit signed integer in little-endian
func (d *D) S34LE() int64 {
	v, err := d.trySEndian(34, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S34LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS34LE tries to add a field and read 34 bit signed integer in little-endian
func (d *D) TryFieldScalarS34LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(34, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS34LE adds a field and reads 34 bit signed integer in little-endian
func (d *D) FieldScalarS34LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS34LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S34LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS34LE tries to add a field and read 34 bit signed integer in little-endian
func (d *D) TryFieldS34LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS34LE(name, sms...)
	return s.ActualS(), err
}

// FieldS34LE adds a field and reads 34 bit signed integer in little-endian
func (d *D) FieldS34LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS34LE(name, sms...).ActualS()
}

// Reader S35LE

// TryS35LE tries to read 35 bit signed integer in little-endian
func (d *D) TryS35LE() (int64, error) { return d.trySEndian(35, LittleEndian) }

// S35LE reads 35 bit signed integer in little-endian
func (d *D) S35LE() int64 {
	v, err := d.trySEndian(35, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S35LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS35LE tries to add a field and read 35 bit signed integer in little-endian
func (d *D) TryFieldScalarS35LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(35, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS35LE adds a field and reads 35 bit signed integer in little-endian
func (d *D) FieldScalarS35LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS35LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S35LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS35LE tries to add a field and read 35 bit signed integer in little-endian
func (d *D) TryFieldS35LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS35LE(name, sms...)
	return s.ActualS(), err
}

// FieldS35LE adds a field and reads 35 bit signed integer in little-endian
func (d *D) FieldS35LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS35LE(name, sms...).ActualS()
}

// Reader S36LE

// TryS36LE tries to read 36 bit signed integer in little-endian
func (d *D) TryS36LE() (int64, error) { return d.trySEndian(36, LittleEndian) }

// S36LE reads 36 bit signed integer in little-endian
func (d *D) S36LE() int64 {
	v, err := d.trySEndian(36, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S36LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS36LE tries to add a field and read 36 bit signed integer in little-endian
func (d *D) TryFieldScalarS36LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(36, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS36LE adds a field and reads 36 bit signed integer in little-endian
func (d *D) FieldScalarS36LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS36LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S36LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS36LE tries to add a field and read 36 bit signed integer in little-endian
func (d *D) TryFieldS36LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS36LE(name, sms...)
	return s.ActualS(), err
}

// FieldS36LE adds a field and reads 36 bit signed integer in little-endian
func (d *D) FieldS36LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS36LE(name, sms...).ActualS()
}

// Reader S37LE

// TryS37LE tries to read 37 bit signed integer in little-endian
func (d *D) TryS37LE() (int64, error) { return d.trySEndian(37, LittleEndian) }

// S37LE reads 37 bit signed integer in little-endian
func (d *D) S37LE() int64 {
	v, err := d.trySEndian(37, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S37LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS37LE tries to add a field and read 37 bit signed integer in little-endian
func (d *D) TryFieldScalarS37LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(37, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS37LE adds a field and reads 37 bit signed integer in little-endian
func (d *D) FieldScalarS37LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS37LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S37LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS37LE tries to add a field and read 37 bit signed integer in little-endian
func (d *D) TryFieldS37LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS37LE(name, sms...)
	return s.ActualS(), err
}

// FieldS37LE adds a field and reads 37 bit signed integer in little-endian
func (d *D) FieldS37LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS37LE(name, sms...).ActualS()
}

// Reader S38LE

// TryS38LE tries to read 38 bit signed integer in little-endian
func (d *D) TryS38LE() (int64, error) { return d.trySEndian(38, LittleEndian) }

// S38LE reads 38 bit signed integer in little-endian
func (d *D) S38LE() int64 {
	v, err := d.trySEndian(38, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S38LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS38LE tries to add a field and read 38 bit signed integer in little-endian
func (d *D) TryFieldScalarS38LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(38, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS38LE adds a field and reads 38 bit signed integer in little-endian
func (d *D) FieldScalarS38LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS38LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S38LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS38LE tries to add a field and read 38 bit signed integer in little-endian
func (d *D) TryFieldS38LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS38LE(name, sms...)
	return s.ActualS(), err
}

// FieldS38LE adds a field and reads 38 bit signed integer in little-endian
func (d *D) FieldS38LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS38LE(name, sms...).ActualS()
}

// Reader S39LE

// TryS39LE tries to read 39 bit signed integer in little-endian
func (d *D) TryS39LE() (int64, error) { return d.trySEndian(39, LittleEndian) }

// S39LE reads 39 bit signed integer in little-endian
func (d *D) S39LE() int64 {
	v, err := d.trySEndian(39, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S39LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS39LE tries to add a field and read 39 bit signed integer in little-endian
func (d *D) TryFieldScalarS39LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(39, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS39LE adds a field and reads 39 bit signed integer in little-endian
func (d *D) FieldScalarS39LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS39LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S39LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS39LE tries to add a field and read 39 bit signed integer in little-endian
func (d *D) TryFieldS39LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS39LE(name, sms...)
	return s.ActualS(), err
}

// FieldS39LE adds a field and reads 39 bit signed integer in little-endian
func (d *D) FieldS39LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS39LE(name, sms...).ActualS()
}

// Reader S40LE

// TryS40LE tries to read 40 bit signed integer in little-endian
func (d *D) TryS40LE() (int64, error) { return d.trySEndian(40, LittleEndian) }

// S40LE reads 40 bit signed integer in little-endian
func (d *D) S40LE() int64 {
	v, err := d.trySEndian(40, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S40LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS40LE tries to add a field and read 40 bit signed integer in little-endian
func (d *D) TryFieldScalarS40LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(40, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS40LE adds a field and reads 40 bit signed integer in little-endian
func (d *D) FieldScalarS40LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS40LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S40LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS40LE tries to add a field and read 40 bit signed integer in little-endian
func (d *D) TryFieldS40LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS40LE(name, sms...)
	return s.ActualS(), err
}

// FieldS40LE adds a field and reads 40 bit signed integer in little-endian
func (d *D) FieldS40LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS40LE(name, sms...).ActualS()
}

// Reader S41LE

// TryS41LE tries to read 41 bit signed integer in little-endian
func (d *D) TryS41LE() (int64, error) { return d.trySEndian(41, LittleEndian) }

// S41LE reads 41 bit signed integer in little-endian
func (d *D) S41LE() int64 {
	v, err := d.trySEndian(41, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S41LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS41LE tries to add a field and read 41 bit signed integer in little-endian
func (d *D) TryFieldScalarS41LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(41, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS41LE adds a field and reads 41 bit signed integer in little-endian
func (d *D) FieldScalarS41LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS41LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S41LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS41LE tries to add a field and read 41 bit signed integer in little-endian
func (d *D) TryFieldS41LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS41LE(name, sms...)
	return s.ActualS(), err
}

// FieldS41LE adds a field and reads 41 bit signed integer in little-endian
func (d *D) FieldS41LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS41LE(name, sms...).ActualS()
}

// Reader S42LE

// TryS42LE tries to read 42 bit signed integer in little-endian
func (d *D) TryS42LE() (int64, error) { return d.trySEndian(42, LittleEndian) }

// S42LE reads 42 bit signed integer in little-endian
func (d *D) S42LE() int64 {
	v, err := d.trySEndian(42, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S42LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS42LE tries to add a field and read 42 bit signed integer in little-endian
func (d *D) TryFieldScalarS42LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(42, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS42LE adds a field and reads 42 bit signed integer in little-endian
func (d *D) FieldScalarS42LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS42LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S42LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS42LE tries to add a field and read 42 bit signed integer in little-endian
func (d *D) TryFieldS42LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS42LE(name, sms...)
	return s.ActualS(), err
}

// FieldS42LE adds a field and reads 42 bit signed integer in little-endian
func (d *D) FieldS42LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS42LE(name, sms...).ActualS()
}

// Reader S43LE

// TryS43LE tries to read 43 bit signed integer in little-endian
func (d *D) TryS43LE() (int64, error) { return d.trySEndian(43, LittleEndian) }

// S43LE reads 43 bit signed integer in little-endian
func (d *D) S43LE() int64 {
	v, err := d.trySEndian(43, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S43LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS43LE tries to add a field and read 43 bit signed integer in little-endian
func (d *D) TryFieldScalarS43LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(43, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS43LE adds a field and reads 43 bit signed integer in little-endian
func (d *D) FieldScalarS43LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS43LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S43LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS43LE tries to add a field and read 43 bit signed integer in little-endian
func (d *D) TryFieldS43LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS43LE(name, sms...)
	return s.ActualS(), err
}

// FieldS43LE adds a field and reads 43 bit signed integer in little-endian
func (d *D) FieldS43LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS43LE(name, sms...).ActualS()
}

// Reader S44LE

// TryS44LE tries to read 44 bit signed integer in little-endian
func (d *D) TryS44LE() (int64, error) { return d.trySEndian(44, LittleEndian) }

// S44LE reads 44 bit signed integer in little-endian
func (d *D) S44LE() int64 {
	v, err := d.trySEndian(44, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S44LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS44LE tries to add a field and read 44 bit signed integer in little-endian
func (d *D) TryFieldScalarS44LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(44, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS44LE adds a field and reads 44 bit signed integer in little-endian
func (d *D) FieldScalarS44LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS44LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S44LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS44LE tries to add a field and read 44 bit signed integer in little-endian
func (d *D) TryFieldS44LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS44LE(name, sms...)
	return s.ActualS(), err
}

// FieldS44LE adds a field and reads 44 bit signed integer in little-endian
func (d *D) FieldS44LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS44LE(name, sms...).ActualS()
}

// Reader S45LE

// TryS45LE tries to read 45 bit signed integer in little-endian
func (d *D) TryS45LE() (int64, error) { return d.trySEndian(45, LittleEndian) }

// S45LE reads 45 bit signed integer in little-endian
func (d *D) S45LE() int64 {
	v, err := d.trySEndian(45, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S45LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS45LE tries to add a field and read 45 bit signed integer in little-endian
func (d *D) TryFieldScalarS45LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(45, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS45LE adds a field and reads 45 bit signed integer in little-endian
func (d *D) FieldScalarS45LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS45LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S45LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS45LE tries to add a field and read 45 bit signed integer in little-endian
func (d *D) TryFieldS45LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS45LE(name, sms...)
	return s.ActualS(), err
}

// FieldS45LE adds a field and reads 45 bit signed integer in little-endian
func (d *D) FieldS45LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS45LE(name, sms...).ActualS()
}

// Reader S46LE

// TryS46LE tries to read 46 bit signed integer in little-endian
func (d *D) TryS46LE() (int64, error) { return d.trySEndian(46, LittleEndian) }

// S46LE reads 46 bit signed integer in little-endian
func (d *D) S46LE() int64 {
	v, err := d.trySEndian(46, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S46LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS46LE tries to add a field and read 46 bit signed integer in little-endian
func (d *D) TryFieldScalarS46LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(46, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS46LE adds a field and reads 46 bit signed integer in little-endian
func (d *D) FieldScalarS46LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS46LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S46LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS46LE tries to add a field and read 46 bit signed integer in little-endian
func (d *D) TryFieldS46LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS46LE(name, sms...)
	return s.ActualS(), err
}

// FieldS46LE adds a field and reads 46 bit signed integer in little-endian
func (d *D) FieldS46LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS46LE(name, sms...).ActualS()
}

// Reader S47LE

// TryS47LE tries to read 47 bit signed integer in little-endian
func (d *D) TryS47LE() (int64, error) { return d.trySEndian(47, LittleEndian) }

// S47LE reads 47 bit signed integer in little-endian
func (d *D) S47LE() int64 {
	v, err := d.trySEndian(47, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S47LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS47LE tries to add a field and read 47 bit signed integer in little-endian
func (d *D) TryFieldScalarS47LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(47, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS47LE adds a field and reads 47 bit signed integer in little-endian
func (d *D) FieldScalarS47LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS47LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S47LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS47LE tries to add a field and read 47 bit signed integer in little-endian
func (d *D) TryFieldS47LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS47LE(name, sms...)
	return s.ActualS(), err
}

// FieldS47LE adds a field and reads 47 bit signed integer in little-endian
func (d *D) FieldS47LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS47LE(name, sms...).ActualS()
}

// Reader S48LE

// TryS48LE tries to read 48 bit signed integer in little-endian
func (d *D) TryS48LE() (int64, error) { return d.trySEndian(48, LittleEndian) }

// S48LE reads 48 bit signed integer in little-endian
func (d *D) S48LE() int64 {
	v, err := d.trySEndian(48, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S48LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS48LE tries to add a field and read 48 bit signed integer in little-endian
func (d *D) TryFieldScalarS48LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(48, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS48LE adds a field and reads 48 bit signed integer in little-endian
func (d *D) FieldScalarS48LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS48LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S48LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS48LE tries to add a field and read 48 bit signed integer in little-endian
func (d *D) TryFieldS48LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS48LE(name, sms...)
	return s.ActualS(), err
}

// FieldS48LE adds a field and reads 48 bit signed integer in little-endian
func (d *D) FieldS48LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS48LE(name, sms...).ActualS()
}

// Reader S49LE

// TryS49LE tries to read 49 bit signed integer in little-endian
func (d *D) TryS49LE() (int64, error) { return d.trySEndian(49, LittleEndian) }

// S49LE reads 49 bit signed integer in little-endian
func (d *D) S49LE() int64 {
	v, err := d.trySEndian(49, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S49LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS49LE tries to add a field and read 49 bit signed integer in little-endian
func (d *D) TryFieldScalarS49LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(49, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS49LE adds a field and reads 49 bit signed integer in little-endian
func (d *D) FieldScalarS49LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS49LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S49LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS49LE tries to add a field and read 49 bit signed integer in little-endian
func (d *D) TryFieldS49LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS49LE(name, sms...)
	return s.ActualS(), err
}

// FieldS49LE adds a field and reads 49 bit signed integer in little-endian
func (d *D) FieldS49LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS49LE(name, sms...).ActualS()
}

// Reader S50LE

// TryS50LE tries to read 50 bit signed integer in little-endian
func (d *D) TryS50LE() (int64, error) { return d.trySEndian(50, LittleEndian) }

// S50LE reads 50 bit signed integer in little-endian
func (d *D) S50LE() int64 {
	v, err := d.trySEndian(50, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S50LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS50LE tries to add a field and read 50 bit signed integer in little-endian
func (d *D) TryFieldScalarS50LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(50, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS50LE adds a field and reads 50 bit signed integer in little-endian
func (d *D) FieldScalarS50LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS50LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S50LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS50LE tries to add a field and read 50 bit signed integer in little-endian
func (d *D) TryFieldS50LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS50LE(name, sms...)
	return s.ActualS(), err
}

// FieldS50LE adds a field and reads 50 bit signed integer in little-endian
func (d *D) FieldS50LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS50LE(name, sms...).ActualS()
}

// Reader S51LE

// TryS51LE tries to read 51 bit signed integer in little-endian
func (d *D) TryS51LE() (int64, error) { return d.trySEndian(51, LittleEndian) }

// S51LE reads 51 bit signed integer in little-endian
func (d *D) S51LE() int64 {
	v, err := d.trySEndian(51, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S51LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS51LE tries to add a field and read 51 bit signed integer in little-endian
func (d *D) TryFieldScalarS51LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(51, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS51LE adds a field and reads 51 bit signed integer in little-endian
func (d *D) FieldScalarS51LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS51LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S51LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS51LE tries to add a field and read 51 bit signed integer in little-endian
func (d *D) TryFieldS51LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS51LE(name, sms...)
	return s.ActualS(), err
}

// FieldS51LE adds a field and reads 51 bit signed integer in little-endian
func (d *D) FieldS51LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS51LE(name, sms...).ActualS()
}

// Reader S52LE

// TryS52LE tries to read 52 bit signed integer in little-endian
func (d *D) TryS52LE() (int64, error) { return d.trySEndian(52, LittleEndian) }

// S52LE reads 52 bit signed integer in little-endian
func (d *D) S52LE() int64 {
	v, err := d.trySEndian(52, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S52LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS52LE tries to add a field and read 52 bit signed integer in little-endian
func (d *D) TryFieldScalarS52LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(52, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS52LE adds a field and reads 52 bit signed integer in little-endian
func (d *D) FieldScalarS52LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS52LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S52LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS52LE tries to add a field and read 52 bit signed integer in little-endian
func (d *D) TryFieldS52LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS52LE(name, sms...)
	return s.ActualS(), err
}

// FieldS52LE adds a field and reads 52 bit signed integer in little-endian
func (d *D) FieldS52LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS52LE(name, sms...).ActualS()
}

// Reader S53LE

// TryS53LE tries to read 53 bit signed integer in little-endian
func (d *D) TryS53LE() (int64, error) { return d.trySEndian(53, LittleEndian) }

// S53LE reads 53 bit signed integer in little-endian
func (d *D) S53LE() int64 {
	v, err := d.trySEndian(53, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S53LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS53LE tries to add a field and read 53 bit signed integer in little-endian
func (d *D) TryFieldScalarS53LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(53, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS53LE adds a field and reads 53 bit signed integer in little-endian
func (d *D) FieldScalarS53LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS53LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S53LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS53LE tries to add a field and read 53 bit signed integer in little-endian
func (d *D) TryFieldS53LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS53LE(name, sms...)
	return s.ActualS(), err
}

// FieldS53LE adds a field and reads 53 bit signed integer in little-endian
func (d *D) FieldS53LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS53LE(name, sms...).ActualS()
}

// Reader S54LE

// TryS54LE tries to read 54 bit signed integer in little-endian
func (d *D) TryS54LE() (int64, error) { return d.trySEndian(54, LittleEndian) }

// S54LE reads 54 bit signed integer in little-endian
func (d *D) S54LE() int64 {
	v, err := d.trySEndian(54, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S54LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS54LE tries to add a field and read 54 bit signed integer in little-endian
func (d *D) TryFieldScalarS54LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(54, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS54LE adds a field and reads 54 bit signed integer in little-endian
func (d *D) FieldScalarS54LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS54LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S54LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS54LE tries to add a field and read 54 bit signed integer in little-endian
func (d *D) TryFieldS54LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS54LE(name, sms...)
	return s.ActualS(), err
}

// FieldS54LE adds a field and reads 54 bit signed integer in little-endian
func (d *D) FieldS54LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS54LE(name, sms...).ActualS()
}

// Reader S55LE

// TryS55LE tries to read 55 bit signed integer in little-endian
func (d *D) TryS55LE() (int64, error) { return d.trySEndian(55, LittleEndian) }

// S55LE reads 55 bit signed integer in little-endian
func (d *D) S55LE() int64 {
	v, err := d.trySEndian(55, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S55LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS55LE tries to add a field and read 55 bit signed integer in little-endian
func (d *D) TryFieldScalarS55LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(55, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS55LE adds a field and reads 55 bit signed integer in little-endian
func (d *D) FieldScalarS55LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS55LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S55LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS55LE tries to add a field and read 55 bit signed integer in little-endian
func (d *D) TryFieldS55LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS55LE(name, sms...)
	return s.ActualS(), err
}

// FieldS55LE adds a field and reads 55 bit signed integer in little-endian
func (d *D) FieldS55LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS55LE(name, sms...).ActualS()
}

// Reader S56LE

// TryS56LE tries to read 56 bit signed integer in little-endian
func (d *D) TryS56LE() (int64, error) { return d.trySEndian(56, LittleEndian) }

// S56LE reads 56 bit signed integer in little-endian
func (d *D) S56LE() int64 {
	v, err := d.trySEndian(56, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S56LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS56LE tries to add a field and read 56 bit signed integer in little-endian
func (d *D) TryFieldScalarS56LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(56, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS56LE adds a field and reads 56 bit signed integer in little-endian
func (d *D) FieldScalarS56LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS56LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S56LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS56LE tries to add a field and read 56 bit signed integer in little-endian
func (d *D) TryFieldS56LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS56LE(name, sms...)
	return s.ActualS(), err
}

// FieldS56LE adds a field and reads 56 bit signed integer in little-endian
func (d *D) FieldS56LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS56LE(name, sms...).ActualS()
}

// Reader S57LE

// TryS57LE tries to read 57 bit signed integer in little-endian
func (d *D) TryS57LE() (int64, error) { return d.trySEndian(57, LittleEndian) }

// S57LE reads 57 bit signed integer in little-endian
func (d *D) S57LE() int64 {
	v, err := d.trySEndian(57, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S57LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS57LE tries to add a field and read 57 bit signed integer in little-endian
func (d *D) TryFieldScalarS57LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(57, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS57LE adds a field and reads 57 bit signed integer in little-endian
func (d *D) FieldScalarS57LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS57LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S57LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS57LE tries to add a field and read 57 bit signed integer in little-endian
func (d *D) TryFieldS57LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS57LE(name, sms...)
	return s.ActualS(), err
}

// FieldS57LE adds a field and reads 57 bit signed integer in little-endian
func (d *D) FieldS57LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS57LE(name, sms...).ActualS()
}

// Reader S58LE

// TryS58LE tries to read 58 bit signed integer in little-endian
func (d *D) TryS58LE() (int64, error) { return d.trySEndian(58, LittleEndian) }

// S58LE reads 58 bit signed integer in little-endian
func (d *D) S58LE() int64 {
	v, err := d.trySEndian(58, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S58LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS58LE tries to add a field and read 58 bit signed integer in little-endian
func (d *D) TryFieldScalarS58LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(58, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS58LE adds a field and reads 58 bit signed integer in little-endian
func (d *D) FieldScalarS58LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS58LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S58LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS58LE tries to add a field and read 58 bit signed integer in little-endian
func (d *D) TryFieldS58LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS58LE(name, sms...)
	return s.ActualS(), err
}

// FieldS58LE adds a field and reads 58 bit signed integer in little-endian
func (d *D) FieldS58LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS58LE(name, sms...).ActualS()
}

// Reader S59LE

// TryS59LE tries to read 59 bit signed integer in little-endian
func (d *D) TryS59LE() (int64, error) { return d.trySEndian(59, LittleEndian) }

// S59LE reads 59 bit signed integer in little-endian
func (d *D) S59LE() int64 {
	v, err := d.trySEndian(59, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S59LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS59LE tries to add a field and read 59 bit signed integer in little-endian
func (d *D) TryFieldScalarS59LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(59, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS59LE adds a field and reads 59 bit signed integer in little-endian
func (d *D) FieldScalarS59LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS59LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S59LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS59LE tries to add a field and read 59 bit signed integer in little-endian
func (d *D) TryFieldS59LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS59LE(name, sms...)
	return s.ActualS(), err
}

// FieldS59LE adds a field and reads 59 bit signed integer in little-endian
func (d *D) FieldS59LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS59LE(name, sms...).ActualS()
}

// Reader S60LE

// TryS60LE tries to read 60 bit signed integer in little-endian
func (d *D) TryS60LE() (int64, error) { return d.trySEndian(60, LittleEndian) }

// S60LE reads 60 bit signed integer in little-endian
func (d *D) S60LE() int64 {
	v, err := d.trySEndian(60, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S60LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS60LE tries to add a field and read 60 bit signed integer in little-endian
func (d *D) TryFieldScalarS60LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(60, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS60LE adds a field and reads 60 bit signed integer in little-endian
func (d *D) FieldScalarS60LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS60LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S60LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS60LE tries to add a field and read 60 bit signed integer in little-endian
func (d *D) TryFieldS60LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS60LE(name, sms...)
	return s.ActualS(), err
}

// FieldS60LE adds a field and reads 60 bit signed integer in little-endian
func (d *D) FieldS60LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS60LE(name, sms...).ActualS()
}

// Reader S61LE

// TryS61LE tries to read 61 bit signed integer in little-endian
func (d *D) TryS61LE() (int64, error) { return d.trySEndian(61, LittleEndian) }

// S61LE reads 61 bit signed integer in little-endian
func (d *D) S61LE() int64 {
	v, err := d.trySEndian(61, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S61LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS61LE tries to add a field and read 61 bit signed integer in little-endian
func (d *D) TryFieldScalarS61LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(61, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS61LE adds a field and reads 61 bit signed integer in little-endian
func (d *D) FieldScalarS61LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS61LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S61LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS61LE tries to add a field and read 61 bit signed integer in little-endian
func (d *D) TryFieldS61LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS61LE(name, sms...)
	return s.ActualS(), err
}

// FieldS61LE adds a field and reads 61 bit signed integer in little-endian
func (d *D) FieldS61LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS61LE(name, sms...).ActualS()
}

// Reader S62LE

// TryS62LE tries to read 62 bit signed integer in little-endian
func (d *D) TryS62LE() (int64, error) { return d.trySEndian(62, LittleEndian) }

// S62LE reads 62 bit signed integer in little-endian
func (d *D) S62LE() int64 {
	v, err := d.trySEndian(62, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S62LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS62LE tries to add a field and read 62 bit signed integer in little-endian
func (d *D) TryFieldScalarS62LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(62, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS62LE adds a field and reads 62 bit signed integer in little-endian
func (d *D) FieldScalarS62LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS62LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S62LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS62LE tries to add a field and read 62 bit signed integer in little-endian
func (d *D) TryFieldS62LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS62LE(name, sms...)
	return s.ActualS(), err
}

// FieldS62LE adds a field and reads 62 bit signed integer in little-endian
func (d *D) FieldS62LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS62LE(name, sms...).ActualS()
}

// Reader S63LE

// TryS63LE tries to read 63 bit signed integer in little-endian
func (d *D) TryS63LE() (int64, error) { return d.trySEndian(63, LittleEndian) }

// S63LE reads 63 bit signed integer in little-endian
func (d *D) S63LE() int64 {
	v, err := d.trySEndian(63, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S63LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS63LE tries to add a field and read 63 bit signed integer in little-endian
func (d *D) TryFieldScalarS63LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(63, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS63LE adds a field and reads 63 bit signed integer in little-endian
func (d *D) FieldScalarS63LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS63LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S63LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS63LE tries to add a field and read 63 bit signed integer in little-endian
func (d *D) TryFieldS63LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS63LE(name, sms...)
	return s.ActualS(), err
}

// FieldS63LE adds a field and reads 63 bit signed integer in little-endian
func (d *D) FieldS63LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS63LE(name, sms...).ActualS()
}

// Reader S64LE

// TryS64LE tries to read 64 bit signed integer in little-endian
func (d *D) TryS64LE() (int64, error) { return d.trySEndian(64, LittleEndian) }

// S64LE reads 64 bit signed integer in little-endian
func (d *D) S64LE() int64 {
	v, err := d.trySEndian(64, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S64LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS64LE tries to add a field and read 64 bit signed integer in little-endian
func (d *D) TryFieldScalarS64LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(64, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS64LE adds a field and reads 64 bit signed integer in little-endian
func (d *D) FieldScalarS64LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS64LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S64LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS64LE tries to add a field and read 64 bit signed integer in little-endian
func (d *D) TryFieldS64LE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS64LE(name, sms...)
	return s.ActualS(), err
}

// FieldS64LE adds a field and reads 64 bit signed integer in little-endian
func (d *D) FieldS64LE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS64LE(name, sms...).ActualS()
}

// Reader S8BE

// TryS8BE tries to read 8 bit signed integer in big-endian
func (d *D) TryS8BE() (int64, error) { return d.trySEndian(8, BigEndian) }

// S8BE reads 8 bit signed integer in big-endian
func (d *D) S8BE() int64 {
	v, err := d.trySEndian(8, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S8BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS8BE tries to add a field and read 8 bit signed integer in big-endian
func (d *D) TryFieldScalarS8BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(8, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS8BE adds a field and reads 8 bit signed integer in big-endian
func (d *D) FieldScalarS8BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS8BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S8BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS8BE tries to add a field and read 8 bit signed integer in big-endian
func (d *D) TryFieldS8BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS8BE(name, sms...)
	return s.ActualS(), err
}

// FieldS8BE adds a field and reads 8 bit signed integer in big-endian
func (d *D) FieldS8BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS8BE(name, sms...).ActualS()
}

// Reader S9BE

// TryS9BE tries to read 9 bit signed integer in big-endian
func (d *D) TryS9BE() (int64, error) { return d.trySEndian(9, BigEndian) }

// S9BE reads 9 bit signed integer in big-endian
func (d *D) S9BE() int64 {
	v, err := d.trySEndian(9, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S9BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS9BE tries to add a field and read 9 bit signed integer in big-endian
func (d *D) TryFieldScalarS9BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(9, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS9BE adds a field and reads 9 bit signed integer in big-endian
func (d *D) FieldScalarS9BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS9BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S9BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS9BE tries to add a field and read 9 bit signed integer in big-endian
func (d *D) TryFieldS9BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS9BE(name, sms...)
	return s.ActualS(), err
}

// FieldS9BE adds a field and reads 9 bit signed integer in big-endian
func (d *D) FieldS9BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS9BE(name, sms...).ActualS()
}

// Reader S10BE

// TryS10BE tries to read 10 bit signed integer in big-endian
func (d *D) TryS10BE() (int64, error) { return d.trySEndian(10, BigEndian) }

// S10BE reads 10 bit signed integer in big-endian
func (d *D) S10BE() int64 {
	v, err := d.trySEndian(10, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S10BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS10BE tries to add a field and read 10 bit signed integer in big-endian
func (d *D) TryFieldScalarS10BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(10, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS10BE adds a field and reads 10 bit signed integer in big-endian
func (d *D) FieldScalarS10BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS10BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S10BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS10BE tries to add a field and read 10 bit signed integer in big-endian
func (d *D) TryFieldS10BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS10BE(name, sms...)
	return s.ActualS(), err
}

// FieldS10BE adds a field and reads 10 bit signed integer in big-endian
func (d *D) FieldS10BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS10BE(name, sms...).ActualS()
}

// Reader S11BE

// TryS11BE tries to read 11 bit signed integer in big-endian
func (d *D) TryS11BE() (int64, error) { return d.trySEndian(11, BigEndian) }

// S11BE reads 11 bit signed integer in big-endian
func (d *D) S11BE() int64 {
	v, err := d.trySEndian(11, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S11BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS11BE tries to add a field and read 11 bit signed integer in big-endian
func (d *D) TryFieldScalarS11BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(11, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS11BE adds a field and reads 11 bit signed integer in big-endian
func (d *D) FieldScalarS11BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS11BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S11BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS11BE tries to add a field and read 11 bit signed integer in big-endian
func (d *D) TryFieldS11BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS11BE(name, sms...)
	return s.ActualS(), err
}

// FieldS11BE adds a field and reads 11 bit signed integer in big-endian
func (d *D) FieldS11BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS11BE(name, sms...).ActualS()
}

// Reader S12BE

// TryS12BE tries to read 12 bit signed integer in big-endian
func (d *D) TryS12BE() (int64, error) { return d.trySEndian(12, BigEndian) }

// S12BE reads 12 bit signed integer in big-endian
func (d *D) S12BE() int64 {
	v, err := d.trySEndian(12, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S12BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS12BE tries to add a field and read 12 bit signed integer in big-endian
func (d *D) TryFieldScalarS12BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(12, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS12BE adds a field and reads 12 bit signed integer in big-endian
func (d *D) FieldScalarS12BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS12BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S12BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS12BE tries to add a field and read 12 bit signed integer in big-endian
func (d *D) TryFieldS12BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS12BE(name, sms...)
	return s.ActualS(), err
}

// FieldS12BE adds a field and reads 12 bit signed integer in big-endian
func (d *D) FieldS12BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS12BE(name, sms...).ActualS()
}

// Reader S13BE

// TryS13BE tries to read 13 bit signed integer in big-endian
func (d *D) TryS13BE() (int64, error) { return d.trySEndian(13, BigEndian) }

// S13BE reads 13 bit signed integer in big-endian
func (d *D) S13BE() int64 {
	v, err := d.trySEndian(13, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S13BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS13BE tries to add a field and read 13 bit signed integer in big-endian
func (d *D) TryFieldScalarS13BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(13, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS13BE adds a field and reads 13 bit signed integer in big-endian
func (d *D) FieldScalarS13BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS13BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S13BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS13BE tries to add a field and read 13 bit signed integer in big-endian
func (d *D) TryFieldS13BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS13BE(name, sms...)
	return s.ActualS(), err
}

// FieldS13BE adds a field and reads 13 bit signed integer in big-endian
func (d *D) FieldS13BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS13BE(name, sms...).ActualS()
}

// Reader S14BE

// TryS14BE tries to read 14 bit signed integer in big-endian
func (d *D) TryS14BE() (int64, error) { return d.trySEndian(14, BigEndian) }

// S14BE reads 14 bit signed integer in big-endian
func (d *D) S14BE() int64 {
	v, err := d.trySEndian(14, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S14BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS14BE tries to add a field and read 14 bit signed integer in big-endian
func (d *D) TryFieldScalarS14BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(14, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS14BE adds a field and reads 14 bit signed integer in big-endian
func (d *D) FieldScalarS14BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS14BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S14BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS14BE tries to add a field and read 14 bit signed integer in big-endian
func (d *D) TryFieldS14BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS14BE(name, sms...)
	return s.ActualS(), err
}

// FieldS14BE adds a field and reads 14 bit signed integer in big-endian
func (d *D) FieldS14BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS14BE(name, sms...).ActualS()
}

// Reader S15BE

// TryS15BE tries to read 15 bit signed integer in big-endian
func (d *D) TryS15BE() (int64, error) { return d.trySEndian(15, BigEndian) }

// S15BE reads 15 bit signed integer in big-endian
func (d *D) S15BE() int64 {
	v, err := d.trySEndian(15, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S15BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS15BE tries to add a field and read 15 bit signed integer in big-endian
func (d *D) TryFieldScalarS15BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(15, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS15BE adds a field and reads 15 bit signed integer in big-endian
func (d *D) FieldScalarS15BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS15BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S15BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS15BE tries to add a field and read 15 bit signed integer in big-endian
func (d *D) TryFieldS15BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS15BE(name, sms...)
	return s.ActualS(), err
}

// FieldS15BE adds a field and reads 15 bit signed integer in big-endian
func (d *D) FieldS15BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS15BE(name, sms...).ActualS()
}

// Reader S16BE

// TryS16BE tries to read 16 bit signed integer in big-endian
func (d *D) TryS16BE() (int64, error) { return d.trySEndian(16, BigEndian) }

// S16BE reads 16 bit signed integer in big-endian
func (d *D) S16BE() int64 {
	v, err := d.trySEndian(16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S16BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS16BE tries to add a field and read 16 bit signed integer in big-endian
func (d *D) TryFieldScalarS16BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(16, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS16BE adds a field and reads 16 bit signed integer in big-endian
func (d *D) FieldScalarS16BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS16BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S16BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS16BE tries to add a field and read 16 bit signed integer in big-endian
func (d *D) TryFieldS16BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS16BE(name, sms...)
	return s.ActualS(), err
}

// FieldS16BE adds a field and reads 16 bit signed integer in big-endian
func (d *D) FieldS16BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS16BE(name, sms...).ActualS()
}

// Reader S17BE

// TryS17BE tries to read 17 bit signed integer in big-endian
func (d *D) TryS17BE() (int64, error) { return d.trySEndian(17, BigEndian) }

// S17BE reads 17 bit signed integer in big-endian
func (d *D) S17BE() int64 {
	v, err := d.trySEndian(17, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S17BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS17BE tries to add a field and read 17 bit signed integer in big-endian
func (d *D) TryFieldScalarS17BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(17, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS17BE adds a field and reads 17 bit signed integer in big-endian
func (d *D) FieldScalarS17BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS17BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S17BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS17BE tries to add a field and read 17 bit signed integer in big-endian
func (d *D) TryFieldS17BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS17BE(name, sms...)
	return s.ActualS(), err
}

// FieldS17BE adds a field and reads 17 bit signed integer in big-endian
func (d *D) FieldS17BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS17BE(name, sms...).ActualS()
}

// Reader S18BE

// TryS18BE tries to read 18 bit signed integer in big-endian
func (d *D) TryS18BE() (int64, error) { return d.trySEndian(18, BigEndian) }

// S18BE reads 18 bit signed integer in big-endian
func (d *D) S18BE() int64 {
	v, err := d.trySEndian(18, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S18BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS18BE tries to add a field and read 18 bit signed integer in big-endian
func (d *D) TryFieldScalarS18BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(18, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS18BE adds a field and reads 18 bit signed integer in big-endian
func (d *D) FieldScalarS18BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS18BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S18BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS18BE tries to add a field and read 18 bit signed integer in big-endian
func (d *D) TryFieldS18BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS18BE(name, sms...)
	return s.ActualS(), err
}

// FieldS18BE adds a field and reads 18 bit signed integer in big-endian
func (d *D) FieldS18BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS18BE(name, sms...).ActualS()
}

// Reader S19BE

// TryS19BE tries to read 19 bit signed integer in big-endian
func (d *D) TryS19BE() (int64, error) { return d.trySEndian(19, BigEndian) }

// S19BE reads 19 bit signed integer in big-endian
func (d *D) S19BE() int64 {
	v, err := d.trySEndian(19, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S19BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS19BE tries to add a field and read 19 bit signed integer in big-endian
func (d *D) TryFieldScalarS19BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(19, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS19BE adds a field and reads 19 bit signed integer in big-endian
func (d *D) FieldScalarS19BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS19BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S19BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS19BE tries to add a field and read 19 bit signed integer in big-endian
func (d *D) TryFieldS19BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS19BE(name, sms...)
	return s.ActualS(), err
}

// FieldS19BE adds a field and reads 19 bit signed integer in big-endian
func (d *D) FieldS19BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS19BE(name, sms...).ActualS()
}

// Reader S20BE

// TryS20BE tries to read 20 bit signed integer in big-endian
func (d *D) TryS20BE() (int64, error) { return d.trySEndian(20, BigEndian) }

// S20BE reads 20 bit signed integer in big-endian
func (d *D) S20BE() int64 {
	v, err := d.trySEndian(20, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S20BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS20BE tries to add a field and read 20 bit signed integer in big-endian
func (d *D) TryFieldScalarS20BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(20, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS20BE adds a field and reads 20 bit signed integer in big-endian
func (d *D) FieldScalarS20BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS20BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S20BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS20BE tries to add a field and read 20 bit signed integer in big-endian
func (d *D) TryFieldS20BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS20BE(name, sms...)
	return s.ActualS(), err
}

// FieldS20BE adds a field and reads 20 bit signed integer in big-endian
func (d *D) FieldS20BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS20BE(name, sms...).ActualS()
}

// Reader S21BE

// TryS21BE tries to read 21 bit signed integer in big-endian
func (d *D) TryS21BE() (int64, error) { return d.trySEndian(21, BigEndian) }

// S21BE reads 21 bit signed integer in big-endian
func (d *D) S21BE() int64 {
	v, err := d.trySEndian(21, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S21BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS21BE tries to add a field and read 21 bit signed integer in big-endian
func (d *D) TryFieldScalarS21BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(21, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS21BE adds a field and reads 21 bit signed integer in big-endian
func (d *D) FieldScalarS21BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS21BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S21BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS21BE tries to add a field and read 21 bit signed integer in big-endian
func (d *D) TryFieldS21BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS21BE(name, sms...)
	return s.ActualS(), err
}

// FieldS21BE adds a field and reads 21 bit signed integer in big-endian
func (d *D) FieldS21BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS21BE(name, sms...).ActualS()
}

// Reader S22BE

// TryS22BE tries to read 22 bit signed integer in big-endian
func (d *D) TryS22BE() (int64, error) { return d.trySEndian(22, BigEndian) }

// S22BE reads 22 bit signed integer in big-endian
func (d *D) S22BE() int64 {
	v, err := d.trySEndian(22, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S22BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS22BE tries to add a field and read 22 bit signed integer in big-endian
func (d *D) TryFieldScalarS22BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(22, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS22BE adds a field and reads 22 bit signed integer in big-endian
func (d *D) FieldScalarS22BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS22BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S22BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS22BE tries to add a field and read 22 bit signed integer in big-endian
func (d *D) TryFieldS22BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS22BE(name, sms...)
	return s.ActualS(), err
}

// FieldS22BE adds a field and reads 22 bit signed integer in big-endian
func (d *D) FieldS22BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS22BE(name, sms...).ActualS()
}

// Reader S23BE

// TryS23BE tries to read 23 bit signed integer in big-endian
func (d *D) TryS23BE() (int64, error) { return d.trySEndian(23, BigEndian) }

// S23BE reads 23 bit signed integer in big-endian
func (d *D) S23BE() int64 {
	v, err := d.trySEndian(23, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S23BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS23BE tries to add a field and read 23 bit signed integer in big-endian
func (d *D) TryFieldScalarS23BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(23, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS23BE adds a field and reads 23 bit signed integer in big-endian
func (d *D) FieldScalarS23BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS23BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S23BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS23BE tries to add a field and read 23 bit signed integer in big-endian
func (d *D) TryFieldS23BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS23BE(name, sms...)
	return s.ActualS(), err
}

// FieldS23BE adds a field and reads 23 bit signed integer in big-endian
func (d *D) FieldS23BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS23BE(name, sms...).ActualS()
}

// Reader S24BE

// TryS24BE tries to read 24 bit signed integer in big-endian
func (d *D) TryS24BE() (int64, error) { return d.trySEndian(24, BigEndian) }

// S24BE reads 24 bit signed integer in big-endian
func (d *D) S24BE() int64 {
	v, err := d.trySEndian(24, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S24BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS24BE tries to add a field and read 24 bit signed integer in big-endian
func (d *D) TryFieldScalarS24BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(24, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS24BE adds a field and reads 24 bit signed integer in big-endian
func (d *D) FieldScalarS24BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS24BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S24BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS24BE tries to add a field and read 24 bit signed integer in big-endian
func (d *D) TryFieldS24BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS24BE(name, sms...)
	return s.ActualS(), err
}

// FieldS24BE adds a field and reads 24 bit signed integer in big-endian
func (d *D) FieldS24BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS24BE(name, sms...).ActualS()
}

// Reader S25BE

// TryS25BE tries to read 25 bit signed integer in big-endian
func (d *D) TryS25BE() (int64, error) { return d.trySEndian(25, BigEndian) }

// S25BE reads 25 bit signed integer in big-endian
func (d *D) S25BE() int64 {
	v, err := d.trySEndian(25, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S25BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS25BE tries to add a field and read 25 bit signed integer in big-endian
func (d *D) TryFieldScalarS25BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(25, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS25BE adds a field and reads 25 bit signed integer in big-endian
func (d *D) FieldScalarS25BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS25BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S25BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS25BE tries to add a field and read 25 bit signed integer in big-endian
func (d *D) TryFieldS25BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS25BE(name, sms...)
	return s.ActualS(), err
}

// FieldS25BE adds a field and reads 25 bit signed integer in big-endian
func (d *D) FieldS25BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS25BE(name, sms...).ActualS()
}

// Reader S26BE

// TryS26BE tries to read 26 bit signed integer in big-endian
func (d *D) TryS26BE() (int64, error) { return d.trySEndian(26, BigEndian) }

// S26BE reads 26 bit signed integer in big-endian
func (d *D) S26BE() int64 {
	v, err := d.trySEndian(26, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S26BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS26BE tries to add a field and read 26 bit signed integer in big-endian
func (d *D) TryFieldScalarS26BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(26, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS26BE adds a field and reads 26 bit signed integer in big-endian
func (d *D) FieldScalarS26BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS26BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S26BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS26BE tries to add a field and read 26 bit signed integer in big-endian
func (d *D) TryFieldS26BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS26BE(name, sms...)
	return s.ActualS(), err
}

// FieldS26BE adds a field and reads 26 bit signed integer in big-endian
func (d *D) FieldS26BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS26BE(name, sms...).ActualS()
}

// Reader S27BE

// TryS27BE tries to read 27 bit signed integer in big-endian
func (d *D) TryS27BE() (int64, error) { return d.trySEndian(27, BigEndian) }

// S27BE reads 27 bit signed integer in big-endian
func (d *D) S27BE() int64 {
	v, err := d.trySEndian(27, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S27BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS27BE tries to add a field and read 27 bit signed integer in big-endian
func (d *D) TryFieldScalarS27BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(27, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS27BE adds a field and reads 27 bit signed integer in big-endian
func (d *D) FieldScalarS27BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS27BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S27BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS27BE tries to add a field and read 27 bit signed integer in big-endian
func (d *D) TryFieldS27BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS27BE(name, sms...)
	return s.ActualS(), err
}

// FieldS27BE adds a field and reads 27 bit signed integer in big-endian
func (d *D) FieldS27BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS27BE(name, sms...).ActualS()
}

// Reader S28BE

// TryS28BE tries to read 28 bit signed integer in big-endian
func (d *D) TryS28BE() (int64, error) { return d.trySEndian(28, BigEndian) }

// S28BE reads 28 bit signed integer in big-endian
func (d *D) S28BE() int64 {
	v, err := d.trySEndian(28, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S28BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS28BE tries to add a field and read 28 bit signed integer in big-endian
func (d *D) TryFieldScalarS28BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(28, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS28BE adds a field and reads 28 bit signed integer in big-endian
func (d *D) FieldScalarS28BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS28BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S28BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS28BE tries to add a field and read 28 bit signed integer in big-endian
func (d *D) TryFieldS28BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS28BE(name, sms...)
	return s.ActualS(), err
}

// FieldS28BE adds a field and reads 28 bit signed integer in big-endian
func (d *D) FieldS28BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS28BE(name, sms...).ActualS()
}

// Reader S29BE

// TryS29BE tries to read 29 bit signed integer in big-endian
func (d *D) TryS29BE() (int64, error) { return d.trySEndian(29, BigEndian) }

// S29BE reads 29 bit signed integer in big-endian
func (d *D) S29BE() int64 {
	v, err := d.trySEndian(29, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S29BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS29BE tries to add a field and read 29 bit signed integer in big-endian
func (d *D) TryFieldScalarS29BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(29, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS29BE adds a field and reads 29 bit signed integer in big-endian
func (d *D) FieldScalarS29BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS29BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S29BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS29BE tries to add a field and read 29 bit signed integer in big-endian
func (d *D) TryFieldS29BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS29BE(name, sms...)
	return s.ActualS(), err
}

// FieldS29BE adds a field and reads 29 bit signed integer in big-endian
func (d *D) FieldS29BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS29BE(name, sms...).ActualS()
}

// Reader S30BE

// TryS30BE tries to read 30 bit signed integer in big-endian
func (d *D) TryS30BE() (int64, error) { return d.trySEndian(30, BigEndian) }

// S30BE reads 30 bit signed integer in big-endian
func (d *D) S30BE() int64 {
	v, err := d.trySEndian(30, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S30BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS30BE tries to add a field and read 30 bit signed integer in big-endian
func (d *D) TryFieldScalarS30BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(30, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS30BE adds a field and reads 30 bit signed integer in big-endian
func (d *D) FieldScalarS30BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS30BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S30BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS30BE tries to add a field and read 30 bit signed integer in big-endian
func (d *D) TryFieldS30BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS30BE(name, sms...)
	return s.ActualS(), err
}

// FieldS30BE adds a field and reads 30 bit signed integer in big-endian
func (d *D) FieldS30BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS30BE(name, sms...).ActualS()
}

// Reader S31BE

// TryS31BE tries to read 31 bit signed integer in big-endian
func (d *D) TryS31BE() (int64, error) { return d.trySEndian(31, BigEndian) }

// S31BE reads 31 bit signed integer in big-endian
func (d *D) S31BE() int64 {
	v, err := d.trySEndian(31, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S31BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS31BE tries to add a field and read 31 bit signed integer in big-endian
func (d *D) TryFieldScalarS31BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(31, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS31BE adds a field and reads 31 bit signed integer in big-endian
func (d *D) FieldScalarS31BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS31BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S31BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS31BE tries to add a field and read 31 bit signed integer in big-endian
func (d *D) TryFieldS31BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS31BE(name, sms...)
	return s.ActualS(), err
}

// FieldS31BE adds a field and reads 31 bit signed integer in big-endian
func (d *D) FieldS31BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS31BE(name, sms...).ActualS()
}

// Reader S32BE

// TryS32BE tries to read 32 bit signed integer in big-endian
func (d *D) TryS32BE() (int64, error) { return d.trySEndian(32, BigEndian) }

// S32BE reads 32 bit signed integer in big-endian
func (d *D) S32BE() int64 {
	v, err := d.trySEndian(32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S32BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS32BE tries to add a field and read 32 bit signed integer in big-endian
func (d *D) TryFieldScalarS32BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(32, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS32BE adds a field and reads 32 bit signed integer in big-endian
func (d *D) FieldScalarS32BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS32BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S32BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS32BE tries to add a field and read 32 bit signed integer in big-endian
func (d *D) TryFieldS32BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS32BE(name, sms...)
	return s.ActualS(), err
}

// FieldS32BE adds a field and reads 32 bit signed integer in big-endian
func (d *D) FieldS32BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS32BE(name, sms...).ActualS()
}

// Reader S33BE

// TryS33BE tries to read 33 bit signed integer in big-endian
func (d *D) TryS33BE() (int64, error) { return d.trySEndian(33, BigEndian) }

// S33BE reads 33 bit signed integer in big-endian
func (d *D) S33BE() int64 {
	v, err := d.trySEndian(33, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S33BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS33BE tries to add a field and read 33 bit signed integer in big-endian
func (d *D) TryFieldScalarS33BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(33, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS33BE adds a field and reads 33 bit signed integer in big-endian
func (d *D) FieldScalarS33BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS33BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S33BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS33BE tries to add a field and read 33 bit signed integer in big-endian
func (d *D) TryFieldS33BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS33BE(name, sms...)
	return s.ActualS(), err
}

// FieldS33BE adds a field and reads 33 bit signed integer in big-endian
func (d *D) FieldS33BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS33BE(name, sms...).ActualS()
}

// Reader S34BE

// TryS34BE tries to read 34 bit signed integer in big-endian
func (d *D) TryS34BE() (int64, error) { return d.trySEndian(34, BigEndian) }

// S34BE reads 34 bit signed integer in big-endian
func (d *D) S34BE() int64 {
	v, err := d.trySEndian(34, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S34BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS34BE tries to add a field and read 34 bit signed integer in big-endian
func (d *D) TryFieldScalarS34BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(34, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS34BE adds a field and reads 34 bit signed integer in big-endian
func (d *D) FieldScalarS34BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS34BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S34BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS34BE tries to add a field and read 34 bit signed integer in big-endian
func (d *D) TryFieldS34BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS34BE(name, sms...)
	return s.ActualS(), err
}

// FieldS34BE adds a field and reads 34 bit signed integer in big-endian
func (d *D) FieldS34BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS34BE(name, sms...).ActualS()
}

// Reader S35BE

// TryS35BE tries to read 35 bit signed integer in big-endian
func (d *D) TryS35BE() (int64, error) { return d.trySEndian(35, BigEndian) }

// S35BE reads 35 bit signed integer in big-endian
func (d *D) S35BE() int64 {
	v, err := d.trySEndian(35, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S35BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS35BE tries to add a field and read 35 bit signed integer in big-endian
func (d *D) TryFieldScalarS35BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(35, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS35BE adds a field and reads 35 bit signed integer in big-endian
func (d *D) FieldScalarS35BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS35BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S35BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS35BE tries to add a field and read 35 bit signed integer in big-endian
func (d *D) TryFieldS35BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS35BE(name, sms...)
	return s.ActualS(), err
}

// FieldS35BE adds a field and reads 35 bit signed integer in big-endian
func (d *D) FieldS35BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS35BE(name, sms...).ActualS()
}

// Reader S36BE

// TryS36BE tries to read 36 bit signed integer in big-endian
func (d *D) TryS36BE() (int64, error) { return d.trySEndian(36, BigEndian) }

// S36BE reads 36 bit signed integer in big-endian
func (d *D) S36BE() int64 {
	v, err := d.trySEndian(36, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S36BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS36BE tries to add a field and read 36 bit signed integer in big-endian
func (d *D) TryFieldScalarS36BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(36, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS36BE adds a field and reads 36 bit signed integer in big-endian
func (d *D) FieldScalarS36BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS36BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S36BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS36BE tries to add a field and read 36 bit signed integer in big-endian
func (d *D) TryFieldS36BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS36BE(name, sms...)
	return s.ActualS(), err
}

// FieldS36BE adds a field and reads 36 bit signed integer in big-endian
func (d *D) FieldS36BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS36BE(name, sms...).ActualS()
}

// Reader S37BE

// TryS37BE tries to read 37 bit signed integer in big-endian
func (d *D) TryS37BE() (int64, error) { return d.trySEndian(37, BigEndian) }

// S37BE reads 37 bit signed integer in big-endian
func (d *D) S37BE() int64 {
	v, err := d.trySEndian(37, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S37BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS37BE tries to add a field and read 37 bit signed integer in big-endian
func (d *D) TryFieldScalarS37BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(37, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS37BE adds a field and reads 37 bit signed integer in big-endian
func (d *D) FieldScalarS37BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS37BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S37BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS37BE tries to add a field and read 37 bit signed integer in big-endian
func (d *D) TryFieldS37BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS37BE(name, sms...)
	return s.ActualS(), err
}

// FieldS37BE adds a field and reads 37 bit signed integer in big-endian
func (d *D) FieldS37BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS37BE(name, sms...).ActualS()
}

// Reader S38BE

// TryS38BE tries to read 38 bit signed integer in big-endian
func (d *D) TryS38BE() (int64, error) { return d.trySEndian(38, BigEndian) }

// S38BE reads 38 bit signed integer in big-endian
func (d *D) S38BE() int64 {
	v, err := d.trySEndian(38, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S38BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS38BE tries to add a field and read 38 bit signed integer in big-endian
func (d *D) TryFieldScalarS38BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(38, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS38BE adds a field and reads 38 bit signed integer in big-endian
func (d *D) FieldScalarS38BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS38BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S38BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS38BE tries to add a field and read 38 bit signed integer in big-endian
func (d *D) TryFieldS38BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS38BE(name, sms...)
	return s.ActualS(), err
}

// FieldS38BE adds a field and reads 38 bit signed integer in big-endian
func (d *D) FieldS38BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS38BE(name, sms...).ActualS()
}

// Reader S39BE

// TryS39BE tries to read 39 bit signed integer in big-endian
func (d *D) TryS39BE() (int64, error) { return d.trySEndian(39, BigEndian) }

// S39BE reads 39 bit signed integer in big-endian
func (d *D) S39BE() int64 {
	v, err := d.trySEndian(39, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S39BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS39BE tries to add a field and read 39 bit signed integer in big-endian
func (d *D) TryFieldScalarS39BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(39, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS39BE adds a field and reads 39 bit signed integer in big-endian
func (d *D) FieldScalarS39BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS39BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S39BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS39BE tries to add a field and read 39 bit signed integer in big-endian
func (d *D) TryFieldS39BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS39BE(name, sms...)
	return s.ActualS(), err
}

// FieldS39BE adds a field and reads 39 bit signed integer in big-endian
func (d *D) FieldS39BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS39BE(name, sms...).ActualS()
}

// Reader S40BE

// TryS40BE tries to read 40 bit signed integer in big-endian
func (d *D) TryS40BE() (int64, error) { return d.trySEndian(40, BigEndian) }

// S40BE reads 40 bit signed integer in big-endian
func (d *D) S40BE() int64 {
	v, err := d.trySEndian(40, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S40BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS40BE tries to add a field and read 40 bit signed integer in big-endian
func (d *D) TryFieldScalarS40BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(40, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS40BE adds a field and reads 40 bit signed integer in big-endian
func (d *D) FieldScalarS40BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS40BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S40BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS40BE tries to add a field and read 40 bit signed integer in big-endian
func (d *D) TryFieldS40BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS40BE(name, sms...)
	return s.ActualS(), err
}

// FieldS40BE adds a field and reads 40 bit signed integer in big-endian
func (d *D) FieldS40BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS40BE(name, sms...).ActualS()
}

// Reader S41BE

// TryS41BE tries to read 41 bit signed integer in big-endian
func (d *D) TryS41BE() (int64, error) { return d.trySEndian(41, BigEndian) }

// S41BE reads 41 bit signed integer in big-endian
func (d *D) S41BE() int64 {
	v, err := d.trySEndian(41, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S41BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS41BE tries to add a field and read 41 bit signed integer in big-endian
func (d *D) TryFieldScalarS41BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(41, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS41BE adds a field and reads 41 bit signed integer in big-endian
func (d *D) FieldScalarS41BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS41BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S41BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS41BE tries to add a field and read 41 bit signed integer in big-endian
func (d *D) TryFieldS41BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS41BE(name, sms...)
	return s.ActualS(), err
}

// FieldS41BE adds a field and reads 41 bit signed integer in big-endian
func (d *D) FieldS41BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS41BE(name, sms...).ActualS()
}

// Reader S42BE

// TryS42BE tries to read 42 bit signed integer in big-endian
func (d *D) TryS42BE() (int64, error) { return d.trySEndian(42, BigEndian) }

// S42BE reads 42 bit signed integer in big-endian
func (d *D) S42BE() int64 {
	v, err := d.trySEndian(42, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S42BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS42BE tries to add a field and read 42 bit signed integer in big-endian
func (d *D) TryFieldScalarS42BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(42, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS42BE adds a field and reads 42 bit signed integer in big-endian
func (d *D) FieldScalarS42BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS42BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S42BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS42BE tries to add a field and read 42 bit signed integer in big-endian
func (d *D) TryFieldS42BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS42BE(name, sms...)
	return s.ActualS(), err
}

// FieldS42BE adds a field and reads 42 bit signed integer in big-endian
func (d *D) FieldS42BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS42BE(name, sms...).ActualS()
}

// Reader S43BE

// TryS43BE tries to read 43 bit signed integer in big-endian
func (d *D) TryS43BE() (int64, error) { return d.trySEndian(43, BigEndian) }

// S43BE reads 43 bit signed integer in big-endian
func (d *D) S43BE() int64 {
	v, err := d.trySEndian(43, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S43BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS43BE tries to add a field and read 43 bit signed integer in big-endian
func (d *D) TryFieldScalarS43BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(43, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS43BE adds a field and reads 43 bit signed integer in big-endian
func (d *D) FieldScalarS43BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS43BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S43BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS43BE tries to add a field and read 43 bit signed integer in big-endian
func (d *D) TryFieldS43BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS43BE(name, sms...)
	return s.ActualS(), err
}

// FieldS43BE adds a field and reads 43 bit signed integer in big-endian
func (d *D) FieldS43BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS43BE(name, sms...).ActualS()
}

// Reader S44BE

// TryS44BE tries to read 44 bit signed integer in big-endian
func (d *D) TryS44BE() (int64, error) { return d.trySEndian(44, BigEndian) }

// S44BE reads 44 bit signed integer in big-endian
func (d *D) S44BE() int64 {
	v, err := d.trySEndian(44, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S44BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS44BE tries to add a field and read 44 bit signed integer in big-endian
func (d *D) TryFieldScalarS44BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(44, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS44BE adds a field and reads 44 bit signed integer in big-endian
func (d *D) FieldScalarS44BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS44BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S44BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS44BE tries to add a field and read 44 bit signed integer in big-endian
func (d *D) TryFieldS44BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS44BE(name, sms...)
	return s.ActualS(), err
}

// FieldS44BE adds a field and reads 44 bit signed integer in big-endian
func (d *D) FieldS44BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS44BE(name, sms...).ActualS()
}

// Reader S45BE

// TryS45BE tries to read 45 bit signed integer in big-endian
func (d *D) TryS45BE() (int64, error) { return d.trySEndian(45, BigEndian) }

// S45BE reads 45 bit signed integer in big-endian
func (d *D) S45BE() int64 {
	v, err := d.trySEndian(45, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S45BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS45BE tries to add a field and read 45 bit signed integer in big-endian
func (d *D) TryFieldScalarS45BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(45, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS45BE adds a field and reads 45 bit signed integer in big-endian
func (d *D) FieldScalarS45BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS45BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S45BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS45BE tries to add a field and read 45 bit signed integer in big-endian
func (d *D) TryFieldS45BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS45BE(name, sms...)
	return s.ActualS(), err
}

// FieldS45BE adds a field and reads 45 bit signed integer in big-endian
func (d *D) FieldS45BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS45BE(name, sms...).ActualS()
}

// Reader S46BE

// TryS46BE tries to read 46 bit signed integer in big-endian
func (d *D) TryS46BE() (int64, error) { return d.trySEndian(46, BigEndian) }

// S46BE reads 46 bit signed integer in big-endian
func (d *D) S46BE() int64 {
	v, err := d.trySEndian(46, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S46BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS46BE tries to add a field and read 46 bit signed integer in big-endian
func (d *D) TryFieldScalarS46BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(46, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS46BE adds a field and reads 46 bit signed integer in big-endian
func (d *D) FieldScalarS46BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS46BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S46BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS46BE tries to add a field and read 46 bit signed integer in big-endian
func (d *D) TryFieldS46BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS46BE(name, sms...)
	return s.ActualS(), err
}

// FieldS46BE adds a field and reads 46 bit signed integer in big-endian
func (d *D) FieldS46BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS46BE(name, sms...).ActualS()
}

// Reader S47BE

// TryS47BE tries to read 47 bit signed integer in big-endian
func (d *D) TryS47BE() (int64, error) { return d.trySEndian(47, BigEndian) }

// S47BE reads 47 bit signed integer in big-endian
func (d *D) S47BE() int64 {
	v, err := d.trySEndian(47, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S47BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS47BE tries to add a field and read 47 bit signed integer in big-endian
func (d *D) TryFieldScalarS47BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(47, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS47BE adds a field and reads 47 bit signed integer in big-endian
func (d *D) FieldScalarS47BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS47BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S47BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS47BE tries to add a field and read 47 bit signed integer in big-endian
func (d *D) TryFieldS47BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS47BE(name, sms...)
	return s.ActualS(), err
}

// FieldS47BE adds a field and reads 47 bit signed integer in big-endian
func (d *D) FieldS47BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS47BE(name, sms...).ActualS()
}

// Reader S48BE

// TryS48BE tries to read 48 bit signed integer in big-endian
func (d *D) TryS48BE() (int64, error) { return d.trySEndian(48, BigEndian) }

// S48BE reads 48 bit signed integer in big-endian
func (d *D) S48BE() int64 {
	v, err := d.trySEndian(48, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S48BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS48BE tries to add a field and read 48 bit signed integer in big-endian
func (d *D) TryFieldScalarS48BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(48, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS48BE adds a field and reads 48 bit signed integer in big-endian
func (d *D) FieldScalarS48BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS48BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S48BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS48BE tries to add a field and read 48 bit signed integer in big-endian
func (d *D) TryFieldS48BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS48BE(name, sms...)
	return s.ActualS(), err
}

// FieldS48BE adds a field and reads 48 bit signed integer in big-endian
func (d *D) FieldS48BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS48BE(name, sms...).ActualS()
}

// Reader S49BE

// TryS49BE tries to read 49 bit signed integer in big-endian
func (d *D) TryS49BE() (int64, error) { return d.trySEndian(49, BigEndian) }

// S49BE reads 49 bit signed integer in big-endian
func (d *D) S49BE() int64 {
	v, err := d.trySEndian(49, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S49BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS49BE tries to add a field and read 49 bit signed integer in big-endian
func (d *D) TryFieldScalarS49BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(49, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS49BE adds a field and reads 49 bit signed integer in big-endian
func (d *D) FieldScalarS49BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS49BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S49BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS49BE tries to add a field and read 49 bit signed integer in big-endian
func (d *D) TryFieldS49BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS49BE(name, sms...)
	return s.ActualS(), err
}

// FieldS49BE adds a field and reads 49 bit signed integer in big-endian
func (d *D) FieldS49BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS49BE(name, sms...).ActualS()
}

// Reader S50BE

// TryS50BE tries to read 50 bit signed integer in big-endian
func (d *D) TryS50BE() (int64, error) { return d.trySEndian(50, BigEndian) }

// S50BE reads 50 bit signed integer in big-endian
func (d *D) S50BE() int64 {
	v, err := d.trySEndian(50, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S50BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS50BE tries to add a field and read 50 bit signed integer in big-endian
func (d *D) TryFieldScalarS50BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(50, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS50BE adds a field and reads 50 bit signed integer in big-endian
func (d *D) FieldScalarS50BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS50BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S50BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS50BE tries to add a field and read 50 bit signed integer in big-endian
func (d *D) TryFieldS50BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS50BE(name, sms...)
	return s.ActualS(), err
}

// FieldS50BE adds a field and reads 50 bit signed integer in big-endian
func (d *D) FieldS50BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS50BE(name, sms...).ActualS()
}

// Reader S51BE

// TryS51BE tries to read 51 bit signed integer in big-endian
func (d *D) TryS51BE() (int64, error) { return d.trySEndian(51, BigEndian) }

// S51BE reads 51 bit signed integer in big-endian
func (d *D) S51BE() int64 {
	v, err := d.trySEndian(51, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S51BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS51BE tries to add a field and read 51 bit signed integer in big-endian
func (d *D) TryFieldScalarS51BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(51, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS51BE adds a field and reads 51 bit signed integer in big-endian
func (d *D) FieldScalarS51BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS51BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S51BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS51BE tries to add a field and read 51 bit signed integer in big-endian
func (d *D) TryFieldS51BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS51BE(name, sms...)
	return s.ActualS(), err
}

// FieldS51BE adds a field and reads 51 bit signed integer in big-endian
func (d *D) FieldS51BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS51BE(name, sms...).ActualS()
}

// Reader S52BE

// TryS52BE tries to read 52 bit signed integer in big-endian
func (d *D) TryS52BE() (int64, error) { return d.trySEndian(52, BigEndian) }

// S52BE reads 52 bit signed integer in big-endian
func (d *D) S52BE() int64 {
	v, err := d.trySEndian(52, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S52BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS52BE tries to add a field and read 52 bit signed integer in big-endian
func (d *D) TryFieldScalarS52BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(52, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS52BE adds a field and reads 52 bit signed integer in big-endian
func (d *D) FieldScalarS52BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS52BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S52BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS52BE tries to add a field and read 52 bit signed integer in big-endian
func (d *D) TryFieldS52BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS52BE(name, sms...)
	return s.ActualS(), err
}

// FieldS52BE adds a field and reads 52 bit signed integer in big-endian
func (d *D) FieldS52BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS52BE(name, sms...).ActualS()
}

// Reader S53BE

// TryS53BE tries to read 53 bit signed integer in big-endian
func (d *D) TryS53BE() (int64, error) { return d.trySEndian(53, BigEndian) }

// S53BE reads 53 bit signed integer in big-endian
func (d *D) S53BE() int64 {
	v, err := d.trySEndian(53, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S53BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS53BE tries to add a field and read 53 bit signed integer in big-endian
func (d *D) TryFieldScalarS53BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(53, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS53BE adds a field and reads 53 bit signed integer in big-endian
func (d *D) FieldScalarS53BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS53BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S53BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS53BE tries to add a field and read 53 bit signed integer in big-endian
func (d *D) TryFieldS53BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS53BE(name, sms...)
	return s.ActualS(), err
}

// FieldS53BE adds a field and reads 53 bit signed integer in big-endian
func (d *D) FieldS53BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS53BE(name, sms...).ActualS()
}

// Reader S54BE

// TryS54BE tries to read 54 bit signed integer in big-endian
func (d *D) TryS54BE() (int64, error) { return d.trySEndian(54, BigEndian) }

// S54BE reads 54 bit signed integer in big-endian
func (d *D) S54BE() int64 {
	v, err := d.trySEndian(54, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S54BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS54BE tries to add a field and read 54 bit signed integer in big-endian
func (d *D) TryFieldScalarS54BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(54, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS54BE adds a field and reads 54 bit signed integer in big-endian
func (d *D) FieldScalarS54BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS54BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S54BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS54BE tries to add a field and read 54 bit signed integer in big-endian
func (d *D) TryFieldS54BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS54BE(name, sms...)
	return s.ActualS(), err
}

// FieldS54BE adds a field and reads 54 bit signed integer in big-endian
func (d *D) FieldS54BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS54BE(name, sms...).ActualS()
}

// Reader S55BE

// TryS55BE tries to read 55 bit signed integer in big-endian
func (d *D) TryS55BE() (int64, error) { return d.trySEndian(55, BigEndian) }

// S55BE reads 55 bit signed integer in big-endian
func (d *D) S55BE() int64 {
	v, err := d.trySEndian(55, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S55BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS55BE tries to add a field and read 55 bit signed integer in big-endian
func (d *D) TryFieldScalarS55BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(55, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS55BE adds a field and reads 55 bit signed integer in big-endian
func (d *D) FieldScalarS55BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS55BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S55BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS55BE tries to add a field and read 55 bit signed integer in big-endian
func (d *D) TryFieldS55BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS55BE(name, sms...)
	return s.ActualS(), err
}

// FieldS55BE adds a field and reads 55 bit signed integer in big-endian
func (d *D) FieldS55BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS55BE(name, sms...).ActualS()
}

// Reader S56BE

// TryS56BE tries to read 56 bit signed integer in big-endian
func (d *D) TryS56BE() (int64, error) { return d.trySEndian(56, BigEndian) }

// S56BE reads 56 bit signed integer in big-endian
func (d *D) S56BE() int64 {
	v, err := d.trySEndian(56, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S56BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS56BE tries to add a field and read 56 bit signed integer in big-endian
func (d *D) TryFieldScalarS56BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(56, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS56BE adds a field and reads 56 bit signed integer in big-endian
func (d *D) FieldScalarS56BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS56BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S56BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS56BE tries to add a field and read 56 bit signed integer in big-endian
func (d *D) TryFieldS56BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS56BE(name, sms...)
	return s.ActualS(), err
}

// FieldS56BE adds a field and reads 56 bit signed integer in big-endian
func (d *D) FieldS56BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS56BE(name, sms...).ActualS()
}

// Reader S57BE

// TryS57BE tries to read 57 bit signed integer in big-endian
func (d *D) TryS57BE() (int64, error) { return d.trySEndian(57, BigEndian) }

// S57BE reads 57 bit signed integer in big-endian
func (d *D) S57BE() int64 {
	v, err := d.trySEndian(57, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S57BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS57BE tries to add a field and read 57 bit signed integer in big-endian
func (d *D) TryFieldScalarS57BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(57, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS57BE adds a field and reads 57 bit signed integer in big-endian
func (d *D) FieldScalarS57BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS57BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S57BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS57BE tries to add a field and read 57 bit signed integer in big-endian
func (d *D) TryFieldS57BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS57BE(name, sms...)
	return s.ActualS(), err
}

// FieldS57BE adds a field and reads 57 bit signed integer in big-endian
func (d *D) FieldS57BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS57BE(name, sms...).ActualS()
}

// Reader S58BE

// TryS58BE tries to read 58 bit signed integer in big-endian
func (d *D) TryS58BE() (int64, error) { return d.trySEndian(58, BigEndian) }

// S58BE reads 58 bit signed integer in big-endian
func (d *D) S58BE() int64 {
	v, err := d.trySEndian(58, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S58BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS58BE tries to add a field and read 58 bit signed integer in big-endian
func (d *D) TryFieldScalarS58BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(58, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS58BE adds a field and reads 58 bit signed integer in big-endian
func (d *D) FieldScalarS58BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS58BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S58BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS58BE tries to add a field and read 58 bit signed integer in big-endian
func (d *D) TryFieldS58BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS58BE(name, sms...)
	return s.ActualS(), err
}

// FieldS58BE adds a field and reads 58 bit signed integer in big-endian
func (d *D) FieldS58BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS58BE(name, sms...).ActualS()
}

// Reader S59BE

// TryS59BE tries to read 59 bit signed integer in big-endian
func (d *D) TryS59BE() (int64, error) { return d.trySEndian(59, BigEndian) }

// S59BE reads 59 bit signed integer in big-endian
func (d *D) S59BE() int64 {
	v, err := d.trySEndian(59, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S59BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS59BE tries to add a field and read 59 bit signed integer in big-endian
func (d *D) TryFieldScalarS59BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(59, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS59BE adds a field and reads 59 bit signed integer in big-endian
func (d *D) FieldScalarS59BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS59BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S59BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS59BE tries to add a field and read 59 bit signed integer in big-endian
func (d *D) TryFieldS59BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS59BE(name, sms...)
	return s.ActualS(), err
}

// FieldS59BE adds a field and reads 59 bit signed integer in big-endian
func (d *D) FieldS59BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS59BE(name, sms...).ActualS()
}

// Reader S60BE

// TryS60BE tries to read 60 bit signed integer in big-endian
func (d *D) TryS60BE() (int64, error) { return d.trySEndian(60, BigEndian) }

// S60BE reads 60 bit signed integer in big-endian
func (d *D) S60BE() int64 {
	v, err := d.trySEndian(60, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S60BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS60BE tries to add a field and read 60 bit signed integer in big-endian
func (d *D) TryFieldScalarS60BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(60, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS60BE adds a field and reads 60 bit signed integer in big-endian
func (d *D) FieldScalarS60BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS60BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S60BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS60BE tries to add a field and read 60 bit signed integer in big-endian
func (d *D) TryFieldS60BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS60BE(name, sms...)
	return s.ActualS(), err
}

// FieldS60BE adds a field and reads 60 bit signed integer in big-endian
func (d *D) FieldS60BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS60BE(name, sms...).ActualS()
}

// Reader S61BE

// TryS61BE tries to read 61 bit signed integer in big-endian
func (d *D) TryS61BE() (int64, error) { return d.trySEndian(61, BigEndian) }

// S61BE reads 61 bit signed integer in big-endian
func (d *D) S61BE() int64 {
	v, err := d.trySEndian(61, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S61BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS61BE tries to add a field and read 61 bit signed integer in big-endian
func (d *D) TryFieldScalarS61BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(61, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS61BE adds a field and reads 61 bit signed integer in big-endian
func (d *D) FieldScalarS61BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS61BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S61BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS61BE tries to add a field and read 61 bit signed integer in big-endian
func (d *D) TryFieldS61BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS61BE(name, sms...)
	return s.ActualS(), err
}

// FieldS61BE adds a field and reads 61 bit signed integer in big-endian
func (d *D) FieldS61BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS61BE(name, sms...).ActualS()
}

// Reader S62BE

// TryS62BE tries to read 62 bit signed integer in big-endian
func (d *D) TryS62BE() (int64, error) { return d.trySEndian(62, BigEndian) }

// S62BE reads 62 bit signed integer in big-endian
func (d *D) S62BE() int64 {
	v, err := d.trySEndian(62, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S62BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS62BE tries to add a field and read 62 bit signed integer in big-endian
func (d *D) TryFieldScalarS62BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(62, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS62BE adds a field and reads 62 bit signed integer in big-endian
func (d *D) FieldScalarS62BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS62BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S62BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS62BE tries to add a field and read 62 bit signed integer in big-endian
func (d *D) TryFieldS62BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS62BE(name, sms...)
	return s.ActualS(), err
}

// FieldS62BE adds a field and reads 62 bit signed integer in big-endian
func (d *D) FieldS62BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS62BE(name, sms...).ActualS()
}

// Reader S63BE

// TryS63BE tries to read 63 bit signed integer in big-endian
func (d *D) TryS63BE() (int64, error) { return d.trySEndian(63, BigEndian) }

// S63BE reads 63 bit signed integer in big-endian
func (d *D) S63BE() int64 {
	v, err := d.trySEndian(63, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S63BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS63BE tries to add a field and read 63 bit signed integer in big-endian
func (d *D) TryFieldScalarS63BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(63, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS63BE adds a field and reads 63 bit signed integer in big-endian
func (d *D) FieldScalarS63BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS63BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S63BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS63BE tries to add a field and read 63 bit signed integer in big-endian
func (d *D) TryFieldS63BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS63BE(name, sms...)
	return s.ActualS(), err
}

// FieldS63BE adds a field and reads 63 bit signed integer in big-endian
func (d *D) FieldS63BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS63BE(name, sms...).ActualS()
}

// Reader S64BE

// TryS64BE tries to read 64 bit signed integer in big-endian
func (d *D) TryS64BE() (int64, error) { return d.trySEndian(64, BigEndian) }

// S64BE reads 64 bit signed integer in big-endian
func (d *D) S64BE() int64 {
	v, err := d.trySEndian(64, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S64BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarS64BE tries to add a field and read 64 bit signed integer in big-endian
func (d *D) TryFieldScalarS64BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySEndian(64, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarS64BE adds a field and reads 64 bit signed integer in big-endian
func (d *D) FieldScalarS64BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarS64BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S64BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldS64BE tries to add a field and read 64 bit signed integer in big-endian
func (d *D) TryFieldS64BE(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarS64BE(name, sms...)
	return s.ActualS(), err
}

// FieldS64BE adds a field and reads 64 bit signed integer in big-endian
func (d *D) FieldS64BE(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarS64BE(name, sms...).ActualS()
}

// Reader UBigInt

// TryUBigInt tries to read nBits bits signed integer in current endian
func (d *D) TryUBigInt(nBits int) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, d.Endian, false)
}

// UBigInt reads nBits bits signed integer in current endian
func (d *D) UBigInt(nBits int) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, d.Endian, false)
	if err != nil {
		panic(IOError{Err: err, Op: "UBigInt", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUBigInt tries to add a field and read nBits bits signed integer in current endian
func (d *D) TryFieldScalarUBigInt(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, d.Endian, false)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUBigInt adds a field and reads nBits bits signed integer in current endian
func (d *D) FieldScalarUBigInt(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUBigInt(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UBigInt", Pos: d.Pos()})
	}
	return s
}

// TryFieldUBigInt tries to add a field and read nBits bits signed integer in current endian
func (d *D) TryFieldUBigInt(name string, nBits int, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarUBigInt(name, nBits, sms...)
	return s.ActualBigInt(), err
}

// FieldUBigInt adds a field and reads nBits bits signed integer in current endian
func (d *D) FieldUBigInt(name string, nBits int, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarUBigInt(name, nBits, sms...).ActualBigInt()
}

// Reader UBigIntE

// TryUBigIntE tries to read nBits signed integer in specified endian
func (d *D) TryUBigIntE(nBits int, endian Endian) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, endian, false)
}

// UBigIntE reads nBits signed integer in specified endian
func (d *D) UBigIntE(nBits int, endian Endian) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, endian, false)
	if err != nil {
		panic(IOError{Err: err, Op: "UBigIntE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUBigIntE tries to add a field and read nBits signed integer in specified endian
func (d *D) TryFieldScalarUBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, endian, false)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUBigIntE adds a field and reads nBits signed integer in specified endian
func (d *D) FieldScalarUBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUBigIntE(name, nBits, endian, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UBigIntE", Pos: d.Pos()})
	}
	return s
}

// TryFieldUBigIntE tries to add a field and read nBits signed integer in specified endian
func (d *D) TryFieldUBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarUBigIntE(name, nBits, endian, sms...)
	return s.ActualBigInt(), err
}

// FieldUBigIntE adds a field and reads nBits signed integer in specified endian
func (d *D) FieldUBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarUBigIntE(name, nBits, endian, sms...).ActualBigInt()
}

// Reader UBigIntLE

// TryUBigIntLE tries to read nBits bit signed integer in little-endian
func (d *D) TryUBigIntLE(nBits int) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, LittleEndian, false)
}

// UBigIntLE reads nBits bit signed integer in little-endian
func (d *D) UBigIntLE(nBits int) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, LittleEndian, false)
	if err != nil {
		panic(IOError{Err: err, Op: "UBigIntLE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUBigIntLE tries to add a field and read nBits bit signed integer in little-endian
func (d *D) TryFieldScalarUBigIntLE(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, LittleEndian, false)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUBigIntLE adds a field and reads nBits bit signed integer in little-endian
func (d *D) FieldScalarUBigIntLE(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUBigIntLE(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UBigIntLE", Pos: d.Pos()})
	}
	return s
}

// TryFieldUBigIntLE tries to add a field and read nBits bit signed integer in little-endian
func (d *D) TryFieldUBigIntLE(name string, nBits int, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarUBigIntLE(name, nBits, sms...)
	return s.ActualBigInt(), err
}

// FieldUBigIntLE adds a field and reads nBits bit signed integer in little-endian
func (d *D) FieldUBigIntLE(name string, nBits int, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarUBigIntLE(name, nBits, sms...).ActualBigInt()
}

// Reader UBigIntBE

// TryUBigIntBE tries to read nBits bit signed integer in big-endian
func (d *D) TryUBigIntBE(nBits int) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, BigEndian, false)
}

// UBigIntBE reads nBits bit signed integer in big-endian
func (d *D) UBigIntBE(nBits int) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, BigEndian, false)
	if err != nil {
		panic(IOError{Err: err, Op: "UBigIntBE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUBigIntBE tries to add a field and read nBits bit signed integer in big-endian
func (d *D) TryFieldScalarUBigIntBE(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, BigEndian, false)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUBigIntBE adds a field and reads nBits bit signed integer in big-endian
func (d *D) FieldScalarUBigIntBE(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUBigIntBE(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UBigIntBE", Pos: d.Pos()})
	}
	return s
}

// TryFieldUBigIntBE tries to add a field and read nBits bit signed integer in big-endian
func (d *D) TryFieldUBigIntBE(name string, nBits int, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarUBigIntBE(name, nBits, sms...)
	return s.ActualBigInt(), err
}

// FieldUBigIntBE adds a field and reads nBits bit signed integer in big-endian
func (d *D) FieldUBigIntBE(name string, nBits int, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarUBigIntBE(name, nBits, sms...).ActualBigInt()
}

// Reader SBigInt

// TrySBigInt tries to read nBits bits signed integer in current endian
func (d *D) TrySBigInt(nBits int) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, d.Endian, true)
}

// SBigInt reads nBits bits signed integer in current endian
func (d *D) SBigInt(nBits int) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, d.Endian, true)
	if err != nil {
		panic(IOError{Err: err, Op: "SBigInt", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarSBigInt tries to add a field and read nBits bits signed integer in current endian
func (d *D) TryFieldScalarSBigInt(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, d.Endian, true)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarSBigInt adds a field and reads nBits bits signed integer in current endian
func (d *D) FieldScalarSBigInt(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarSBigInt(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "SBigInt", Pos: d.Pos()})
	}
	return s
}

// TryFieldSBigInt tries to add a field and read nBits bits signed integer in current endian
func (d *D) TryFieldSBigInt(name string, nBits int, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarSBigInt(name, nBits, sms...)
	return s.ActualBigInt(), err
}

// FieldSBigInt adds a field and reads nBits bits signed integer in current endian
func (d *D) FieldSBigInt(name string, nBits int, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarSBigInt(name, nBits, sms...).ActualBigInt()
}

// Reader SBigIntE

// TrySBigIntE tries to read nBits signed integer in specified endian
func (d *D) TrySBigIntE(nBits int, endian Endian) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, endian, true)
}

// SBigIntE reads nBits signed integer in specified endian
func (d *D) SBigIntE(nBits int, endian Endian) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, endian, true)
	if err != nil {
		panic(IOError{Err: err, Op: "SBigIntE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarSBigIntE tries to add a field and read nBits signed integer in specified endian
func (d *D) TryFieldScalarSBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, endian, true)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarSBigIntE adds a field and reads nBits signed integer in specified endian
func (d *D) FieldScalarSBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarSBigIntE(name, nBits, endian, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "SBigIntE", Pos: d.Pos()})
	}
	return s
}

// TryFieldSBigIntE tries to add a field and read nBits signed integer in specified endian
func (d *D) TryFieldSBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarSBigIntE(name, nBits, endian, sms...)
	return s.ActualBigInt(), err
}

// FieldSBigIntE adds a field and reads nBits signed integer in specified endian
func (d *D) FieldSBigIntE(name string, nBits int, endian Endian, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarSBigIntE(name, nBits, endian, sms...).ActualBigInt()
}

// Reader SBigIntLE

// TrySBigIntLE tries to read nBits bit signed integer in little-endian
func (d *D) TrySBigIntLE(nBits int) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, LittleEndian, true)
}

// SBigIntLE reads nBits bit signed integer in little-endian
func (d *D) SBigIntLE(nBits int) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, LittleEndian, true)
	if err != nil {
		panic(IOError{Err: err, Op: "SBigIntLE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarSBigIntLE tries to add a field and read nBits bit signed integer in little-endian
func (d *D) TryFieldScalarSBigIntLE(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, LittleEndian, true)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarSBigIntLE adds a field and reads nBits bit signed integer in little-endian
func (d *D) FieldScalarSBigIntLE(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarSBigIntLE(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "SBigIntLE", Pos: d.Pos()})
	}
	return s
}

// TryFieldSBigIntLE tries to add a field and read nBits bit signed integer in little-endian
func (d *D) TryFieldSBigIntLE(name string, nBits int, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarSBigIntLE(name, nBits, sms...)
	return s.ActualBigInt(), err
}

// FieldSBigIntLE adds a field and reads nBits bit signed integer in little-endian
func (d *D) FieldSBigIntLE(name string, nBits int, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarSBigIntLE(name, nBits, sms...).ActualBigInt()
}

// Reader SBigIntBE

// TrySBigIntBE tries to read nBits bit signed integer in big-endian
func (d *D) TrySBigIntBE(nBits int) (*big.Int, error) {
	return d.tryBigIntEndianSign(nBits, BigEndian, true)
}

// SBigIntBE reads nBits bit signed integer in big-endian
func (d *D) SBigIntBE(nBits int) *big.Int {
	v, err := d.tryBigIntEndianSign(nBits, BigEndian, true)
	if err != nil {
		panic(IOError{Err: err, Op: "SBigIntBE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarSBigIntBE tries to add a field and read nBits bit signed integer in big-endian
func (d *D) TryFieldScalarSBigIntBE(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryBigIntEndianSign(nBits, BigEndian, true)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarSBigIntBE adds a field and reads nBits bit signed integer in big-endian
func (d *D) FieldScalarSBigIntBE(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarSBigIntBE(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "SBigIntBE", Pos: d.Pos()})
	}
	return s
}

// TryFieldSBigIntBE tries to add a field and read nBits bit signed integer in big-endian
func (d *D) TryFieldSBigIntBE(name string, nBits int, sms ...scalar.Mapper) (*big.Int, error) {
	s, err := d.TryFieldScalarSBigIntBE(name, nBits, sms...)
	return s.ActualBigInt(), err
}

// FieldSBigIntBE adds a field and reads nBits bit signed integer in big-endian
func (d *D) FieldSBigIntBE(name string, nBits int, sms ...scalar.Mapper) *big.Int {
	return d.FieldScalarSBigIntBE(name, nBits, sms...).ActualBigInt()
}

// Reader F

// TryF tries to read nBit IEEE 754 float in current endian
func (d *D) TryF(nBits int) (float64, error) { return d.tryFEndian(nBits, d.Endian) }

// F reads nBit IEEE 754 float in current endian
func (d *D) F(nBits int) float64 {
	v, err := d.tryFEndian(nBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF tries to add a field and read nBit IEEE 754 float in current endian
func (d *D) TryFieldScalarF(name string, nBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(nBits, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF adds a field and reads nBit IEEE 754 float in current endian
func (d *D) FieldScalarF(name string, nBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF(name, nBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F", Pos: d.Pos()})
	}
	return s
}

// TryFieldF tries to add a field and read nBit IEEE 754 float in current endian
func (d *D) TryFieldF(name string, nBits int, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF(name, nBits, sms...)
	return s.ActualF(), err
}

// FieldF adds a field and reads nBit IEEE 754 float in current endian
func (d *D) FieldF(name string, nBits int, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF(name, nBits, sms...).ActualF()
}

// Reader FE

// TryFE tries to read nBit IEEE 754 float in specified endian
func (d *D) TryFE(nBits int, endian Endian) (float64, error) { return d.tryFEndian(nBits, endian) }

// FE reads nBit IEEE 754 float in specified endian
func (d *D) FE(nBits int, endian Endian) float64 {
	v, err := d.tryFEndian(nBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFE tries to add a field and read nBit IEEE 754 float in specified endian
func (d *D) TryFieldScalarFE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(nBits, endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFE adds a field and reads nBit IEEE 754 float in specified endian
func (d *D) FieldScalarFE(name string, nBits int, endian Endian, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFE(name, nBits, endian, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFE tries to add a field and read nBit IEEE 754 float in specified endian
func (d *D) TryFieldFE(name string, nBits int, endian Endian, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFE(name, nBits, endian, sms...)
	return s.ActualF(), err
}

// FieldFE adds a field and reads nBit IEEE 754 float in specified endian
func (d *D) FieldFE(name string, nBits int, endian Endian, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFE(name, nBits, endian, sms...).ActualF()
}

// Reader F16

// TryF16 tries to read 16 bit IEEE 754 float in current endian
func (d *D) TryF16() (float64, error) { return d.tryFEndian(16, d.Endian) }

// F16 reads 16 bit IEEE 754 float in current endian
func (d *D) F16() float64 {
	v, err := d.tryFEndian(16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F16", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF16 tries to add a field and read 16 bit IEEE 754 float in current endian
func (d *D) TryFieldScalarF16(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(16, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF16 adds a field and reads 16 bit IEEE 754 float in current endian
func (d *D) FieldScalarF16(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF16(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F16", Pos: d.Pos()})
	}
	return s
}

// TryFieldF16 tries to add a field and read 16 bit IEEE 754 float in current endian
func (d *D) TryFieldF16(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF16(name, sms...)
	return s.ActualF(), err
}

// FieldF16 adds a field and reads 16 bit IEEE 754 float in current endian
func (d *D) FieldF16(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF16(name, sms...).ActualF()
}

// Reader F32

// TryF32 tries to read 32 bit IEEE 754 float in current endian
func (d *D) TryF32() (float64, error) { return d.tryFEndian(32, d.Endian) }

// F32 reads 32 bit IEEE 754 float in current endian
func (d *D) F32() float64 {
	v, err := d.tryFEndian(32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F32", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF32 tries to add a field and read 32 bit IEEE 754 float in current endian
func (d *D) TryFieldScalarF32(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(32, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF32 adds a field and reads 32 bit IEEE 754 float in current endian
func (d *D) FieldScalarF32(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF32(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F32", Pos: d.Pos()})
	}
	return s
}

// TryFieldF32 tries to add a field and read 32 bit IEEE 754 float in current endian
func (d *D) TryFieldF32(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF32(name, sms...)
	return s.ActualF(), err
}

// FieldF32 adds a field and reads 32 bit IEEE 754 float in current endian
func (d *D) FieldF32(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF32(name, sms...).ActualF()
}

// Reader F64

// TryF64 tries to read 64 bit IEEE 754 float in current endian
func (d *D) TryF64() (float64, error) { return d.tryFEndian(64, d.Endian) }

// F64 reads 64 bit IEEE 754 float in current endian
func (d *D) F64() float64 {
	v, err := d.tryFEndian(64, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F64", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF64 tries to add a field and read 64 bit IEEE 754 float in current endian
func (d *D) TryFieldScalarF64(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(64, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF64 adds a field and reads 64 bit IEEE 754 float in current endian
func (d *D) FieldScalarF64(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF64(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F64", Pos: d.Pos()})
	}
	return s
}

// TryFieldF64 tries to add a field and read 64 bit IEEE 754 float in current endian
func (d *D) TryFieldF64(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF64(name, sms...)
	return s.ActualF(), err
}

// FieldF64 adds a field and reads 64 bit IEEE 754 float in current endian
func (d *D) FieldF64(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF64(name, sms...).ActualF()
}

// Reader F16LE

// TryF16LE tries to read 16 bit IEEE 754 float in little-endian
func (d *D) TryF16LE() (float64, error) { return d.tryFEndian(16, LittleEndian) }

// F16LE reads 16 bit IEEE 754 float in little-endian
func (d *D) F16LE() float64 {
	v, err := d.tryFEndian(16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F16LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF16LE tries to add a field and read 16 bit IEEE 754 float in little-endian
func (d *D) TryFieldScalarF16LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(16, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF16LE adds a field and reads 16 bit IEEE 754 float in little-endian
func (d *D) FieldScalarF16LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF16LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F16LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldF16LE tries to add a field and read 16 bit IEEE 754 float in little-endian
func (d *D) TryFieldF16LE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF16LE(name, sms...)
	return s.ActualF(), err
}

// FieldF16LE adds a field and reads 16 bit IEEE 754 float in little-endian
func (d *D) FieldF16LE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF16LE(name, sms...).ActualF()
}

// Reader F32LE

// TryF32LE tries to read 32 bit IEEE 754 float in little-endian
func (d *D) TryF32LE() (float64, error) { return d.tryFEndian(32, LittleEndian) }

// F32LE reads 32 bit IEEE 754 float in little-endian
func (d *D) F32LE() float64 {
	v, err := d.tryFEndian(32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F32LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF32LE tries to add a field and read 32 bit IEEE 754 float in little-endian
func (d *D) TryFieldScalarF32LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(32, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF32LE adds a field and reads 32 bit IEEE 754 float in little-endian
func (d *D) FieldScalarF32LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF32LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F32LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldF32LE tries to add a field and read 32 bit IEEE 754 float in little-endian
func (d *D) TryFieldF32LE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF32LE(name, sms...)
	return s.ActualF(), err
}

// FieldF32LE adds a field and reads 32 bit IEEE 754 float in little-endian
func (d *D) FieldF32LE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF32LE(name, sms...).ActualF()
}

// Reader F64LE

// TryF64LE tries to read 64 bit IEEE 754 float in little-endian
func (d *D) TryF64LE() (float64, error) { return d.tryFEndian(64, LittleEndian) }

// F64LE reads 64 bit IEEE 754 float in little-endian
func (d *D) F64LE() float64 {
	v, err := d.tryFEndian(64, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F64LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF64LE tries to add a field and read 64 bit IEEE 754 float in little-endian
func (d *D) TryFieldScalarF64LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(64, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF64LE adds a field and reads 64 bit IEEE 754 float in little-endian
func (d *D) FieldScalarF64LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF64LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F64LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldF64LE tries to add a field and read 64 bit IEEE 754 float in little-endian
func (d *D) TryFieldF64LE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF64LE(name, sms...)
	return s.ActualF(), err
}

// FieldF64LE adds a field and reads 64 bit IEEE 754 float in little-endian
func (d *D) FieldF64LE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF64LE(name, sms...).ActualF()
}

// Reader F16BE

// TryF16BE tries to read 16 bit IEEE 754 float in big-endian
func (d *D) TryF16BE() (float64, error) { return d.tryFEndian(16, BigEndian) }

// F16BE reads 16 bit IEEE 754 float in big-endian
func (d *D) F16BE() float64 {
	v, err := d.tryFEndian(16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F16BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF16BE tries to add a field and read 16 bit IEEE 754 float in big-endian
func (d *D) TryFieldScalarF16BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(16, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF16BE adds a field and reads 16 bit IEEE 754 float in big-endian
func (d *D) FieldScalarF16BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF16BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F16BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldF16BE tries to add a field and read 16 bit IEEE 754 float in big-endian
func (d *D) TryFieldF16BE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF16BE(name, sms...)
	return s.ActualF(), err
}

// FieldF16BE adds a field and reads 16 bit IEEE 754 float in big-endian
func (d *D) FieldF16BE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF16BE(name, sms...).ActualF()
}

// Reader F32BE

// TryF32BE tries to read 32 bit IEEE 754 float in big-endian
func (d *D) TryF32BE() (float64, error) { return d.tryFEndian(32, BigEndian) }

// F32BE reads 32 bit IEEE 754 float in big-endian
func (d *D) F32BE() float64 {
	v, err := d.tryFEndian(32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F32BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF32BE tries to add a field and read 32 bit IEEE 754 float in big-endian
func (d *D) TryFieldScalarF32BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(32, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF32BE adds a field and reads 32 bit IEEE 754 float in big-endian
func (d *D) FieldScalarF32BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF32BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F32BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldF32BE tries to add a field and read 32 bit IEEE 754 float in big-endian
func (d *D) TryFieldF32BE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF32BE(name, sms...)
	return s.ActualF(), err
}

// FieldF32BE adds a field and reads 32 bit IEEE 754 float in big-endian
func (d *D) FieldF32BE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF32BE(name, sms...).ActualF()
}

// Reader F64BE

// TryF64BE tries to read 64 bit IEEE 754 float in big-endian
func (d *D) TryF64BE() (float64, error) { return d.tryFEndian(64, BigEndian) }

// F64BE reads 64 bit IEEE 754 float in big-endian
func (d *D) F64BE() float64 {
	v, err := d.tryFEndian(64, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F64BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarF64BE tries to add a field and read 64 bit IEEE 754 float in big-endian
func (d *D) TryFieldScalarF64BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFEndian(64, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarF64BE adds a field and reads 64 bit IEEE 754 float in big-endian
func (d *D) FieldScalarF64BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarF64BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F64BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldF64BE tries to add a field and read 64 bit IEEE 754 float in big-endian
func (d *D) TryFieldF64BE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarF64BE(name, sms...)
	return s.ActualF(), err
}

// FieldF64BE adds a field and reads 64 bit IEEE 754 float in big-endian
func (d *D) FieldF64BE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarF64BE(name, sms...).ActualF()
}

// Reader FP

// TryFP tries to read nBits fixed-point number in current endian
func (d *D) TryFP(nBits int, fBits int) (float64, error) {
	return d.tryFPEndian(nBits, fBits, d.Endian)
}

// FP reads nBits fixed-point number in current endian
func (d *D) FP(nBits int, fBits int) float64 {
	v, err := d.tryFPEndian(nBits, fBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP tries to add a field and read nBits fixed-point number in current endian
func (d *D) TryFieldScalarFP(name string, nBits int, fBits int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(nBits, fBits, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP adds a field and reads nBits fixed-point number in current endian
func (d *D) FieldScalarFP(name string, nBits int, fBits int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP(name, nBits, fBits, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP tries to add a field and read nBits fixed-point number in current endian
func (d *D) TryFieldFP(name string, nBits int, fBits int, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP(name, nBits, fBits, sms...)
	return s.ActualF(), err
}

// FieldFP adds a field and reads nBits fixed-point number in current endian
func (d *D) FieldFP(name string, nBits int, fBits int, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP(name, nBits, fBits, sms...).ActualF()
}

// Reader FPE

// TryFPE tries to read nBits fixed-point number in specified endian
func (d *D) TryFPE(nBits int, fBits int, endian Endian) (float64, error) {
	return d.tryFPEndian(nBits, fBits, endian)
}

// FPE reads nBits fixed-point number in specified endian
func (d *D) FPE(nBits int, fBits int, endian Endian) float64 {
	v, err := d.tryFPEndian(nBits, fBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FPE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFPE tries to add a field and read nBits fixed-point number in specified endian
func (d *D) TryFieldScalarFPE(name string, nBits int, fBits int, endian Endian, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(nBits, fBits, endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFPE adds a field and reads nBits fixed-point number in specified endian
func (d *D) FieldScalarFPE(name string, nBits int, fBits int, endian Endian, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFPE(name, nBits, fBits, endian, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FPE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFPE tries to add a field and read nBits fixed-point number in specified endian
func (d *D) TryFieldFPE(name string, nBits int, fBits int, endian Endian, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFPE(name, nBits, fBits, endian, sms...)
	return s.ActualF(), err
}

// FieldFPE adds a field and reads nBits fixed-point number in specified endian
func (d *D) FieldFPE(name string, nBits int, fBits int, endian Endian, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFPE(name, nBits, fBits, endian, sms...).ActualF()
}

// Reader FP16

// TryFP16 tries to read 16 bit fixed-point number in current endian
func (d *D) TryFP16() (float64, error) { return d.tryFPEndian(16, 8, d.Endian) }

// FP16 reads 16 bit fixed-point number in current endian
func (d *D) FP16() float64 {
	v, err := d.tryFPEndian(16, 8, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP16", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP16 tries to add a field and read 16 bit fixed-point number in current endian
func (d *D) TryFieldScalarFP16(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(16, 8, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP16 adds a field and reads 16 bit fixed-point number in current endian
func (d *D) FieldScalarFP16(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP16(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP16", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP16 tries to add a field and read 16 bit fixed-point number in current endian
func (d *D) TryFieldFP16(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP16(name, sms...)
	return s.ActualF(), err
}

// FieldFP16 adds a field and reads 16 bit fixed-point number in current endian
func (d *D) FieldFP16(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP16(name, sms...).ActualF()
}

// Reader FP32

// TryFP32 tries to read 32 bit fixed-point number in current endian
func (d *D) TryFP32() (float64, error) { return d.tryFPEndian(32, 16, d.Endian) }

// FP32 reads 32 bit fixed-point number in current endian
func (d *D) FP32() float64 {
	v, err := d.tryFPEndian(32, 16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP32", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP32 tries to add a field and read 32 bit fixed-point number in current endian
func (d *D) TryFieldScalarFP32(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(32, 16, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP32 adds a field and reads 32 bit fixed-point number in current endian
func (d *D) FieldScalarFP32(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP32(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP32", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP32 tries to add a field and read 32 bit fixed-point number in current endian
func (d *D) TryFieldFP32(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP32(name, sms...)
	return s.ActualF(), err
}

// FieldFP32 adds a field and reads 32 bit fixed-point number in current endian
func (d *D) FieldFP32(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP32(name, sms...).ActualF()
}

// Reader FP64

// TryFP64 tries to read 64 bit fixed-point number in current endian
func (d *D) TryFP64() (float64, error) { return d.tryFPEndian(64, 32, d.Endian) }

// FP64 reads 64 bit fixed-point number in current endian
func (d *D) FP64() float64 {
	v, err := d.tryFPEndian(64, 32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP64", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP64 tries to add a field and read 64 bit fixed-point number in current endian
func (d *D) TryFieldScalarFP64(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(64, 32, d.Endian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP64 adds a field and reads 64 bit fixed-point number in current endian
func (d *D) FieldScalarFP64(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP64(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP64", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP64 tries to add a field and read 64 bit fixed-point number in current endian
func (d *D) TryFieldFP64(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP64(name, sms...)
	return s.ActualF(), err
}

// FieldFP64 adds a field and reads 64 bit fixed-point number in current endian
func (d *D) FieldFP64(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP64(name, sms...).ActualF()
}

// Reader FP16LE

// TryFP16LE tries to read 16 bit fixed-point number in little-endian
func (d *D) TryFP16LE() (float64, error) { return d.tryFPEndian(16, 8, LittleEndian) }

// FP16LE reads 16 bit fixed-point number in little-endian
func (d *D) FP16LE() float64 {
	v, err := d.tryFPEndian(16, 8, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP16LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP16LE tries to add a field and read 16 bit fixed-point number in little-endian
func (d *D) TryFieldScalarFP16LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(16, 8, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP16LE adds a field and reads 16 bit fixed-point number in little-endian
func (d *D) FieldScalarFP16LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP16LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP16LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP16LE tries to add a field and read 16 bit fixed-point number in little-endian
func (d *D) TryFieldFP16LE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP16LE(name, sms...)
	return s.ActualF(), err
}

// FieldFP16LE adds a field and reads 16 bit fixed-point number in little-endian
func (d *D) FieldFP16LE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP16LE(name, sms...).ActualF()
}

// Reader FP32LE

// TryFP32LE tries to read 32 bit fixed-point number in little-endian
func (d *D) TryFP32LE() (float64, error) { return d.tryFPEndian(32, 16, LittleEndian) }

// FP32LE reads 32 bit fixed-point number in little-endian
func (d *D) FP32LE() float64 {
	v, err := d.tryFPEndian(32, 16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP32LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP32LE tries to add a field and read 32 bit fixed-point number in little-endian
func (d *D) TryFieldScalarFP32LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(32, 16, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP32LE adds a field and reads 32 bit fixed-point number in little-endian
func (d *D) FieldScalarFP32LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP32LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP32LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP32LE tries to add a field and read 32 bit fixed-point number in little-endian
func (d *D) TryFieldFP32LE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP32LE(name, sms...)
	return s.ActualF(), err
}

// FieldFP32LE adds a field and reads 32 bit fixed-point number in little-endian
func (d *D) FieldFP32LE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP32LE(name, sms...).ActualF()
}

// Reader FP64LE

// TryFP64LE tries to read 64 bit fixed-point number in little-endian
func (d *D) TryFP64LE() (float64, error) { return d.tryFPEndian(64, 32, LittleEndian) }

// FP64LE reads 64 bit fixed-point number in little-endian
func (d *D) FP64LE() float64 {
	v, err := d.tryFPEndian(64, 32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP64LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP64LE tries to add a field and read 64 bit fixed-point number in little-endian
func (d *D) TryFieldScalarFP64LE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(64, 32, LittleEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP64LE adds a field and reads 64 bit fixed-point number in little-endian
func (d *D) FieldScalarFP64LE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP64LE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP64LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP64LE tries to add a field and read 64 bit fixed-point number in little-endian
func (d *D) TryFieldFP64LE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP64LE(name, sms...)
	return s.ActualF(), err
}

// FieldFP64LE adds a field and reads 64 bit fixed-point number in little-endian
func (d *D) FieldFP64LE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP64LE(name, sms...).ActualF()
}

// Reader FP16BE

// TryFP16BE tries to read 16 bit fixed-point number in big-endian
func (d *D) TryFP16BE() (float64, error) { return d.tryFPEndian(16, 8, BigEndian) }

// FP16BE reads 16 bit fixed-point number in big-endian
func (d *D) FP16BE() float64 {
	v, err := d.tryFPEndian(16, 8, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP16BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP16BE tries to add a field and read 16 bit fixed-point number in big-endian
func (d *D) TryFieldScalarFP16BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(16, 8, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP16BE adds a field and reads 16 bit fixed-point number in big-endian
func (d *D) FieldScalarFP16BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP16BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP16BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP16BE tries to add a field and read 16 bit fixed-point number in big-endian
func (d *D) TryFieldFP16BE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP16BE(name, sms...)
	return s.ActualF(), err
}

// FieldFP16BE adds a field and reads 16 bit fixed-point number in big-endian
func (d *D) FieldFP16BE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP16BE(name, sms...).ActualF()
}

// Reader FP32BE

// TryFP32BE tries to read 32 bit fixed-point number in big-endian
func (d *D) TryFP32BE() (float64, error) { return d.tryFPEndian(32, 16, BigEndian) }

// FP32BE reads 32 bit fixed-point number in big-endian
func (d *D) FP32BE() float64 {
	v, err := d.tryFPEndian(32, 16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP32BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP32BE tries to add a field and read 32 bit fixed-point number in big-endian
func (d *D) TryFieldScalarFP32BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(32, 16, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP32BE adds a field and reads 32 bit fixed-point number in big-endian
func (d *D) FieldScalarFP32BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP32BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP32BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP32BE tries to add a field and read 32 bit fixed-point number in big-endian
func (d *D) TryFieldFP32BE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP32BE(name, sms...)
	return s.ActualF(), err
}

// FieldFP32BE adds a field and reads 32 bit fixed-point number in big-endian
func (d *D) FieldFP32BE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP32BE(name, sms...).ActualF()
}

// Reader FP64BE

// TryFP64BE tries to read 64 bit fixed-point number in big-endian
func (d *D) TryFP64BE() (float64, error) { return d.tryFPEndian(64, 32, BigEndian) }

// FP64BE reads 64 bit fixed-point number in big-endian
func (d *D) FP64BE() float64 {
	v, err := d.tryFPEndian(64, 32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP64BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarFP64BE tries to add a field and read 64 bit fixed-point number in big-endian
func (d *D) TryFieldScalarFP64BE(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryFPEndian(64, 32, BigEndian)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarFP64BE adds a field and reads 64 bit fixed-point number in big-endian
func (d *D) FieldScalarFP64BE(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarFP64BE(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP64BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldFP64BE tries to add a field and read 64 bit fixed-point number in big-endian
func (d *D) TryFieldFP64BE(name string, sms ...scalar.Mapper) (float64, error) {
	s, err := d.TryFieldScalarFP64BE(name, sms...)
	return s.ActualF(), err
}

// FieldFP64BE adds a field and reads 64 bit fixed-point number in big-endian
func (d *D) FieldFP64BE(name string, sms ...scalar.Mapper) float64 {
	return d.FieldScalarFP64BE(name, sms...).ActualF()
}

// Reader Unary

// TryUnary tries to read unary integer using ov as "one" value
func (d *D) TryUnary(ov uint64) (uint64, error) { return d.tryUnary(ov) }

// Unary reads unary integer using ov as "one" value
func (d *D) Unary(ov uint64) uint64 {
	v, err := d.tryUnary(ov)
	if err != nil {
		panic(IOError{Err: err, Op: "Unary", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUnary tries to add a field and read unary integer using ov as "one" value
func (d *D) TryFieldScalarUnary(name string, ov uint64, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryUnary(ov)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUnary adds a field and reads unary integer using ov as "one" value
func (d *D) FieldScalarUnary(name string, ov uint64, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUnary(name, ov, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Unary", Pos: d.Pos()})
	}
	return s
}

// TryFieldUnary tries to add a field and read unary integer using ov as "one" value
func (d *D) TryFieldUnary(name string, ov uint64, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarUnary(name, ov, sms...)
	return s.ActualU(), err
}

// FieldUnary adds a field and reads unary integer using ov as "one" value
func (d *D) FieldUnary(name string, ov uint64, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarUnary(name, ov, sms...).ActualU()
}

// Reader ULEB128

// TryULEB128 tries to read unsigned LEB128 integer
func (d *D) TryULEB128() (uint64, error) { return d.tryULEB128() }

// ULEB128 reads unsigned LEB128 integer
func (d *D) ULEB128() uint64 {
	v, err := d.tryULEB128()
	if err != nil {
		panic(IOError{Err: err, Op: "ULEB128", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarULEB128 tries to add a field and read unsigned LEB128 integer
func (d *D) TryFieldScalarULEB128(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryULEB128()
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarULEB128 adds a field and reads unsigned LEB128 integer
func (d *D) FieldScalarULEB128(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarULEB128(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "ULEB128", Pos: d.Pos()})
	}
	return s
}

// TryFieldULEB128 tries to add a field and read unsigned LEB128 integer
func (d *D) TryFieldULEB128(name string, sms ...scalar.Mapper) (uint64, error) {
	s, err := d.TryFieldScalarULEB128(name, sms...)
	return s.ActualU(), err
}

// FieldULEB128 adds a field and reads unsigned LEB128 integer
func (d *D) FieldULEB128(name string, sms ...scalar.Mapper) uint64 {
	return d.FieldScalarULEB128(name, sms...).ActualU()
}

// Reader SLEB128

// TrySLEB128 tries to read signed LEB128 integer
func (d *D) TrySLEB128() (int64, error) { return d.trySLEB128() }

// SLEB128 reads signed LEB128 integer
func (d *D) SLEB128() int64 {
	v, err := d.trySLEB128()
	if err != nil {
		panic(IOError{Err: err, Op: "SLEB128", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarSLEB128 tries to add a field and read signed LEB128 integer
func (d *D) TryFieldScalarSLEB128(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.trySLEB128()
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarSLEB128 adds a field and reads signed LEB128 integer
func (d *D) FieldScalarSLEB128(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarSLEB128(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "SLEB128", Pos: d.Pos()})
	}
	return s
}

// TryFieldSLEB128 tries to add a field and read signed LEB128 integer
func (d *D) TryFieldSLEB128(name string, sms ...scalar.Mapper) (int64, error) {
	s, err := d.TryFieldScalarSLEB128(name, sms...)
	return s.ActualS(), err
}

// FieldSLEB128 adds a field and reads signed LEB128 integer
func (d *D) FieldSLEB128(name string, sms ...scalar.Mapper) int64 {
	return d.FieldScalarSLEB128(name, sms...).ActualS()
}

// Reader UTF8

// TryUTF8 tries to read nBytes bytes UTF8 string
func (d *D) TryUTF8(nBytes int) (string, error) { return d.tryText(nBytes, UTF8BOM) }

// UTF8 reads nBytes bytes UTF8 string
func (d *D) UTF8(nBytes int) string {
	v, err := d.tryText(nBytes, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF8 tries to add a field and read nBytes bytes UTF8 string
func (d *D) TryFieldScalarUTF8(name string, nBytes int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryText(nBytes, UTF8BOM)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF8 adds a field and reads nBytes bytes UTF8 string
func (d *D) FieldScalarUTF8(name string, nBytes int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF8(name, nBytes, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF8 tries to add a field and read nBytes bytes UTF8 string
func (d *D) TryFieldUTF8(name string, nBytes int, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF8(name, nBytes, sms...)
	return s.ActualStr(), err
}

// FieldUTF8 adds a field and reads nBytes bytes UTF8 string
func (d *D) FieldUTF8(name string, nBytes int, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF8(name, nBytes, sms...).ActualStr()
}

// Reader UTF16

// TryUTF16 tries to read nBytes bytes UTF16 string, default big-endian and accepts BOM
func (d *D) TryUTF16(nBytes int) (string, error) { return d.tryText(nBytes, UTF16BOM) }

// UTF16 reads nBytes bytes UTF16 string, default big-endian and accepts BOM
func (d *D) UTF16(nBytes int) string {
	v, err := d.tryText(nBytes, UTF16BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF16 tries to add a field and read nBytes bytes UTF16 string, default big-endian and accepts BOM
func (d *D) TryFieldScalarUTF16(name string, nBytes int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryText(nBytes, UTF16BOM)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF16 adds a field and reads nBytes bytes UTF16 string, default big-endian and accepts BOM
func (d *D) FieldScalarUTF16(name string, nBytes int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF16(name, nBytes, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF16 tries to add a field and read nBytes bytes UTF16 string, default big-endian and accepts BOM
func (d *D) TryFieldUTF16(name string, nBytes int, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF16(name, nBytes, sms...)
	return s.ActualStr(), err
}

// FieldUTF16 adds a field and reads nBytes bytes UTF16 string, default big-endian and accepts BOM
func (d *D) FieldUTF16(name string, nBytes int, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF16(name, nBytes, sms...).ActualStr()
}

// Reader UTF16LE

// TryUTF16LE tries to read nBytes bytes UTF16 little-endian string
func (d *D) TryUTF16LE(nBytes int) (string, error) { return d.tryText(nBytes, UTF16LE) }

// UTF16LE reads nBytes bytes UTF16 little-endian string
func (d *D) UTF16LE(nBytes int) string {
	v, err := d.tryText(nBytes, UTF16LE)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16LE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF16LE tries to add a field and read nBytes bytes UTF16 little-endian string
func (d *D) TryFieldScalarUTF16LE(name string, nBytes int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryText(nBytes, UTF16LE)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF16LE adds a field and reads nBytes bytes UTF16 little-endian string
func (d *D) FieldScalarUTF16LE(name string, nBytes int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF16LE(name, nBytes, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16LE", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF16LE tries to add a field and read nBytes bytes UTF16 little-endian string
func (d *D) TryFieldUTF16LE(name string, nBytes int, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF16LE(name, nBytes, sms...)
	return s.ActualStr(), err
}

// FieldUTF16LE adds a field and reads nBytes bytes UTF16 little-endian string
func (d *D) FieldUTF16LE(name string, nBytes int, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF16LE(name, nBytes, sms...).ActualStr()
}

// Reader UTF16BE

// TryUTF16BE tries to read nBytes bytes UTF16 big-endian string
func (d *D) TryUTF16BE(nBytes int) (string, error) { return d.tryText(nBytes, UTF16BE) }

// UTF16BE reads nBytes bytes UTF16 big-endian string
func (d *D) UTF16BE(nBytes int) string {
	v, err := d.tryText(nBytes, UTF16BE)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16BE", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF16BE tries to add a field and read nBytes bytes UTF16 big-endian string
func (d *D) TryFieldScalarUTF16BE(name string, nBytes int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryText(nBytes, UTF16BE)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF16BE adds a field and reads nBytes bytes UTF16 big-endian string
func (d *D) FieldScalarUTF16BE(name string, nBytes int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF16BE(name, nBytes, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16BE", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF16BE tries to add a field and read nBytes bytes UTF16 big-endian string
func (d *D) TryFieldUTF16BE(name string, nBytes int, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF16BE(name, nBytes, sms...)
	return s.ActualStr(), err
}

// FieldUTF16BE adds a field and reads nBytes bytes UTF16 big-endian string
func (d *D) FieldUTF16BE(name string, nBytes int, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF16BE(name, nBytes, sms...).ActualStr()
}

// Reader UTF8ShortString

// TryUTF8ShortString tries to read one byte length fixed UTF8 string
func (d *D) TryUTF8ShortString() (string, error) { return d.tryTextLenPrefixed(8, -1, UTF8BOM) }

// UTF8ShortString reads one byte length fixed UTF8 string
func (d *D) UTF8ShortString() string {
	v, err := d.tryTextLenPrefixed(8, -1, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8ShortString", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF8ShortString tries to add a field and read one byte length fixed UTF8 string
func (d *D) TryFieldScalarUTF8ShortString(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryTextLenPrefixed(8, -1, UTF8BOM)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF8ShortString adds a field and reads one byte length fixed UTF8 string
func (d *D) FieldScalarUTF8ShortString(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF8ShortString(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8ShortString", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF8ShortString tries to add a field and read one byte length fixed UTF8 string
func (d *D) TryFieldUTF8ShortString(name string, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF8ShortString(name, sms...)
	return s.ActualStr(), err
}

// FieldUTF8ShortString adds a field and reads one byte length fixed UTF8 string
func (d *D) FieldUTF8ShortString(name string, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF8ShortString(name, sms...).ActualStr()
}

// Reader UTF8ShortStringFixedLen

// TryUTF8ShortStringFixedLen tries to read fixedBytes bytes long one byte length prefixed UTF8 string
func (d *D) TryUTF8ShortStringFixedLen(fixedBytes int) (string, error) {
	return d.tryTextLenPrefixed(8, fixedBytes, UTF8BOM)
}

// UTF8ShortStringFixedLen reads fixedBytes bytes long one byte length prefixed UTF8 string
func (d *D) UTF8ShortStringFixedLen(fixedBytes int) string {
	v, err := d.tryTextLenPrefixed(8, fixedBytes, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8ShortStringFixedLen", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF8ShortStringFixedLen tries to add a field and read fixedBytes bytes long one byte length prefixed UTF8 string
func (d *D) TryFieldScalarUTF8ShortStringFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryTextLenPrefixed(8, fixedBytes, UTF8BOM)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF8ShortStringFixedLen adds a field and reads fixedBytes bytes long one byte length prefixed UTF8 string
func (d *D) FieldScalarUTF8ShortStringFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF8ShortStringFixedLen(name, fixedBytes, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8ShortStringFixedLen", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF8ShortStringFixedLen tries to add a field and read fixedBytes bytes long one byte length prefixed UTF8 string
func (d *D) TryFieldUTF8ShortStringFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF8ShortStringFixedLen(name, fixedBytes, sms...)
	return s.ActualStr(), err
}

// FieldUTF8ShortStringFixedLen adds a field and reads fixedBytes bytes long one byte length prefixed UTF8 string
func (d *D) FieldUTF8ShortStringFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF8ShortStringFixedLen(name, fixedBytes, sms...).ActualStr()
}

// Reader UTF8Null

// TryUTF8Null tries to read null terminated UTF8 string
func (d *D) TryUTF8Null() (string, error) { return d.tryTextNull(1, UTF8BOM) }

// UTF8Null reads null terminated UTF8 string
func (d *D) UTF8Null() string {
	v, err := d.tryTextNull(1, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8Null", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF8Null tries to add a field and read null terminated UTF8 string
func (d *D) TryFieldScalarUTF8Null(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryTextNull(1, UTF8BOM)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF8Null adds a field and reads null terminated UTF8 string
func (d *D) FieldScalarUTF8Null(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF8Null(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8Null", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF8Null tries to add a field and read null terminated UTF8 string
func (d *D) TryFieldUTF8Null(name string, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF8Null(name, sms...)
	return s.ActualStr(), err
}

// FieldUTF8Null adds a field and reads null terminated UTF8 string
func (d *D) FieldUTF8Null(name string, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF8Null(name, sms...).ActualStr()
}

// Reader UTF16Null

// TryUTF16Null tries to read null terminated UTF16 string, default big-endian and accepts BOM
func (d *D) TryUTF16Null() (string, error) { return d.tryTextNull(2, UTF16BOM) }

// UTF16Null reads null terminated UTF16 string, default big-endian and accepts BOM
func (d *D) UTF16Null() string {
	v, err := d.tryTextNull(2, UTF16BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16Null", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF16Null tries to add a field and read null terminated UTF16 string, default big-endian and accepts BOM
func (d *D) TryFieldScalarUTF16Null(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryTextNull(2, UTF16BOM)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF16Null adds a field and reads null terminated UTF16 string, default big-endian and accepts BOM
func (d *D) FieldScalarUTF16Null(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF16Null(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16Null", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF16Null tries to add a field and read null terminated UTF16 string, default big-endian and accepts BOM
func (d *D) TryFieldUTF16Null(name string, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF16Null(name, sms...)
	return s.ActualStr(), err
}

// FieldUTF16Null adds a field and reads null terminated UTF16 string, default big-endian and accepts BOM
func (d *D) FieldUTF16Null(name string, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF16Null(name, sms...).ActualStr()
}

// Reader UTF16LENull

// TryUTF16LENull tries to read null terminated UTF16LE string
func (d *D) TryUTF16LENull() (string, error) { return d.tryTextNull(2, UTF16LE) }

// UTF16LENull reads null terminated UTF16LE string
func (d *D) UTF16LENull() string {
	v, err := d.tryTextNull(2, UTF16LE)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16LENull", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF16LENull tries to add a field and read null terminated UTF16LE string
func (d *D) TryFieldScalarUTF16LENull(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryTextNull(2, UTF16LE)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF16LENull adds a field and reads null terminated UTF16LE string
func (d *D) FieldScalarUTF16LENull(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF16LENull(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16LENull", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF16LENull tries to add a field and read null terminated UTF16LE string
func (d *D) TryFieldUTF16LENull(name string, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF16LENull(name, sms...)
	return s.ActualStr(), err
}

// FieldUTF16LENull adds a field and reads null terminated UTF16LE string
func (d *D) FieldUTF16LENull(name string, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF16LENull(name, sms...).ActualStr()
}

// Reader UTF16BENull

// TryUTF16BENull tries to read null terminated UTF16BE string
func (d *D) TryUTF16BENull() (string, error) { return d.tryTextNull(2, UTF16BE) }

// UTF16BENull reads null terminated UTF16BE string
func (d *D) UTF16BENull() string {
	v, err := d.tryTextNull(2, UTF16BE)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16BENull", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF16BENull tries to add a field and read null terminated UTF16BE string
func (d *D) TryFieldScalarUTF16BENull(name string, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryTextNull(2, UTF16BE)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF16BENull adds a field and reads null terminated UTF16BE string
func (d *D) FieldScalarUTF16BENull(name string, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF16BENull(name, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16BENull", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF16BENull tries to add a field and read null terminated UTF16BE string
func (d *D) TryFieldUTF16BENull(name string, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF16BENull(name, sms...)
	return s.ActualStr(), err
}

// FieldUTF16BENull adds a field and reads null terminated UTF16BE string
func (d *D) FieldUTF16BENull(name string, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF16BENull(name, sms...).ActualStr()
}

// Reader UTF8NullFixedLen

// TryUTF8NullFixedLen tries to read fixedBytes bytes long null terminated UTF8 string
func (d *D) TryUTF8NullFixedLen(fixedBytes int) (string, error) {
	return d.tryTextNullLen(fixedBytes, UTF8BOM)
}

// UTF8NullFixedLen reads fixedBytes bytes long null terminated UTF8 string
func (d *D) UTF8NullFixedLen(fixedBytes int) string {
	v, err := d.tryTextNullLen(fixedBytes, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8NullFixedLen", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarUTF8NullFixedLen tries to add a field and read fixedBytes bytes long null terminated UTF8 string
func (d *D) TryFieldScalarUTF8NullFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryTextNullLen(fixedBytes, UTF8BOM)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarUTF8NullFixedLen adds a field and reads fixedBytes bytes long null terminated UTF8 string
func (d *D) FieldScalarUTF8NullFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarUTF8NullFixedLen(name, fixedBytes, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8NullFixedLen", Pos: d.Pos()})
	}
	return s
}

// TryFieldUTF8NullFixedLen tries to add a field and read fixedBytes bytes long null terminated UTF8 string
func (d *D) TryFieldUTF8NullFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarUTF8NullFixedLen(name, fixedBytes, sms...)
	return s.ActualStr(), err
}

// FieldUTF8NullFixedLen adds a field and reads fixedBytes bytes long null terminated UTF8 string
func (d *D) FieldUTF8NullFixedLen(name string, fixedBytes int, sms ...scalar.Mapper) string {
	return d.FieldScalarUTF8NullFixedLen(name, fixedBytes, sms...).ActualStr()
}

// Reader Str

// TryStr tries to read nBytes bytes using encoding e
func (d *D) TryStr(nBytes int, e encoding.Encoding) (string, error) { return d.tryText(nBytes, e) }

// Str reads nBytes bytes using encoding e
func (d *D) Str(nBytes int, e encoding.Encoding) string {
	v, err := d.tryText(nBytes, e)
	if err != nil {
		panic(IOError{Err: err, Op: "Str", Pos: d.Pos()})
	}
	return v
}

// TryFieldScalarStr tries to add a field and read nBytes bytes using encoding e
func (d *D) TryFieldScalarStr(name string, nBytes int, e encoding.Encoding, sms ...scalar.Mapper) (*scalar.S, error) {
	s, err := d.TryFieldScalarFn(name, func(s scalar.S) (scalar.S, error) {
		v, err := d.tryText(nBytes, e)
		s.Actual = v
		return s, err
	}, sms...)
	if err != nil {
		return nil, err
	}
	return s, err
}

// FieldScalarStr adds a field and reads nBytes bytes using encoding e
func (d *D) FieldScalarStr(name string, nBytes int, e encoding.Encoding, sms ...scalar.Mapper) *scalar.S {
	s, err := d.TryFieldScalarStr(name, nBytes, e, sms...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Str", Pos: d.Pos()})
	}
	return s
}

// TryFieldStr tries to add a field and read nBytes bytes using encoding e
func (d *D) TryFieldStr(name string, nBytes int, e encoding.Encoding, sms ...scalar.Mapper) (string, error) {
	s, err := d.TryFieldScalarStr(name, nBytes, e, sms...)
	return s.ActualStr(), err
}

// FieldStr adds a field and reads nBytes bytes using encoding e
func (d *D) FieldStr(name string, nBytes int, e encoding.Encoding, sms ...scalar.Mapper) string {
	return d.FieldScalarStr(name, nBytes, e, sms...).ActualStr()
}
