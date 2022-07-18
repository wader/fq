// Code below generated from scalar_gen.go.tmpl
package scalar

import (
	"fmt"
	"math/big"

	"github.com/wader/fq/pkg/bitio"
)

// Type BigInt

// ActualBigInt asserts actual value is a BigInt and returns it
func (s S) ActualBigInt() *big.Int {
	v, ok := s.Actual.(*big.Int)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as *big.Int", s.Actual))
	}
	return v
}

// ActualBigIntFn map actual BigInt using f
type ActualBigIntFn func(a *big.Int) *big.Int

func (fn ActualBigIntFn) MapScalar(s S) (S, error) {
	s.Actual = fn(s.ActualBigInt())
	return s, nil
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s S) SymBigInt() *big.Int {
	v, ok := s.Sym.(*big.Int)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as *big.Int", s.Sym))
	}
	return v
}

// SymBigIntFn map sym BigInt using f
type SymBigIntFn func(a *big.Int) *big.Int

func (f SymBigIntFn) MapScalar(s S) (S, error) {
	s.Sym = f(s.SymBigInt())
	return s, nil
}

// Type BitBuf

// ActualBitBuf asserts actual value is a BitBuf and returns it
func (s S) ActualBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.Actual.(bitio.ReaderAtSeeker)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as bitio.ReaderAtSeeker", s.Actual))
	}
	return v
}

// ActualBitBufFn map actual BitBuf using f
type ActualBitBufFn func(a bitio.ReaderAtSeeker) bitio.ReaderAtSeeker

func (fn ActualBitBufFn) MapScalar(s S) (S, error) {
	s.Actual = fn(s.ActualBitBuf())
	return s, nil
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s S) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as bitio.ReaderAtSeeker", s.Sym))
	}
	return v
}

// SymBitBufFn map sym BitBuf using f
type SymBitBufFn func(a bitio.ReaderAtSeeker) bitio.ReaderAtSeeker

func (f SymBitBufFn) MapScalar(s S) (S, error) {
	s.Sym = f(s.SymBitBuf())
	return s, nil
}

// Type Bool

// ActualBool asserts actual value is a Bool and returns it
func (s S) ActualBool() bool {
	v, ok := s.Actual.(bool)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as bool", s.Actual))
	}
	return v
}

// ActualBoolFn map actual Bool using f
type ActualBoolFn func(a bool) bool

func (fn ActualBoolFn) MapScalar(s S) (S, error) {
	s.Actual = fn(s.ActualBool())
	return s, nil
}

// SymBool asserts symbolic value is a Bool and returns it
func (s S) SymBool() bool {
	v, ok := s.Sym.(bool)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as bool", s.Sym))
	}
	return v
}

// SymBoolFn map sym Bool using f
type SymBoolFn func(a bool) bool

func (f SymBoolFn) MapScalar(s S) (S, error) {
	s.Sym = f(s.SymBool())
	return s, nil
}

// Type F

// ActualF asserts actual value is a F and returns it
func (s S) ActualF() float64 {
	v, ok := s.Actual.(float64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as float64", s.Actual))
	}
	return v
}

// ActualFFn map actual F using f
type ActualFFn func(a float64) float64

func (fn ActualFFn) MapScalar(s S) (S, error) {
	s.Actual = fn(s.ActualF())
	return s, nil
}

// SymF asserts symbolic value is a F and returns it
func (s S) SymF() float64 {
	v, ok := s.Sym.(float64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as float64", s.Sym))
	}
	return v
}

// SymFFn map sym F using f
type SymFFn func(a float64) float64

func (f SymFFn) MapScalar(s S) (S, error) {
	s.Sym = f(s.SymF())
	return s, nil
}

// Type S

// ActualS asserts actual value is a S and returns it
func (s S) ActualS() int64 {
	v, ok := s.Actual.(int64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as int64", s.Actual))
	}
	return v
}

// ActualSFn map actual S using f
type ActualSFn func(a int64) int64

func (fn ActualSFn) MapScalar(s S) (S, error) {
	s.Actual = fn(s.ActualS())
	return s, nil
}

// SymS asserts symbolic value is a S and returns it
func (s S) SymS() int64 {
	v, ok := s.Sym.(int64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as int64", s.Sym))
	}
	return v
}

// SymSFn map sym S using f
type SymSFn func(a int64) int64

func (f SymSFn) MapScalar(s S) (S, error) {
	s.Sym = f(s.SymS())
	return s, nil
}

// Type Str

// ActualStr asserts actual value is a Str and returns it
func (s S) ActualStr() string {
	v, ok := s.Actual.(string)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as string", s.Actual))
	}
	return v
}

// ActualStrFn map actual Str using f
type ActualStrFn func(a string) string

func (fn ActualStrFn) MapScalar(s S) (S, error) {
	s.Actual = fn(s.ActualStr())
	return s, nil
}

// SymStr asserts symbolic value is a Str and returns it
func (s S) SymStr() string {
	v, ok := s.Sym.(string)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as string", s.Sym))
	}
	return v
}

// SymStrFn map sym Str using f
type SymStrFn func(a string) string

func (f SymStrFn) MapScalar(s S) (S, error) {
	s.Sym = f(s.SymStr())
	return s, nil
}

// Type U

// ActualU asserts actual value is a U and returns it
func (s S) ActualU() uint64 {
	v, ok := s.Actual.(uint64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as uint64", s.Actual))
	}
	return v
}

// ActualUFn map actual U using f
type ActualUFn func(a uint64) uint64

func (fn ActualUFn) MapScalar(s S) (S, error) {
	s.Actual = fn(s.ActualU())
	return s, nil
}

// SymU asserts symbolic value is a U and returns it
func (s S) SymU() uint64 {
	v, ok := s.Sym.(uint64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as uint64", s.Sym))
	}
	return v
}

// SymUFn map sym U using f
type SymUFn func(a uint64) uint64

func (f SymUFn) MapScalar(s S) (S, error) {
	s.Sym = f(s.SymU())
	return s, nil
}

// Map Bool -> Scalar
type BoolToScalar map[bool]S

func (m BoolToScalar) MapScalar(s S) (S, error) {
	a := s.ActualBool()
	if ns, ok := m[a]; ok {
		ns.Actual = a
		s = ns
	}
	return s, nil
}

// Map Bool -> Description
type BoolToDescription map[bool]string

func (m BoolToDescription) MapScalar(s S) (S, error) {
	a := s.ActualBool()
	if d, ok := m[a]; ok {
		s.Description = d
	}
	return s, nil
}

// Map Bool -> Sym Bool
type BoolToSymBool map[bool]bool

func (m BoolToSymBool) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualBool()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool -> Sym F
type BoolToSymF map[bool]float64

func (m BoolToSymF) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualBool()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool -> Sym S
type BoolToSymS map[bool]int64

func (m BoolToSymS) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualBool()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool -> Sym Str
type BoolToSymStr map[bool]string

func (m BoolToSymStr) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualBool()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool -> Sym U
type BoolToSymU map[bool]uint64

func (m BoolToSymU) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualBool()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map S -> Scalar
type SToScalar map[int64]S

func (m SToScalar) MapScalar(s S) (S, error) {
	a := s.ActualS()
	if ns, ok := m[a]; ok {
		ns.Actual = a
		s = ns
	}
	return s, nil
}

// Map S -> Description
type SToDescription map[int64]string

func (m SToDescription) MapScalar(s S) (S, error) {
	a := s.ActualS()
	if d, ok := m[a]; ok {
		s.Description = d
	}
	return s, nil
}

// Map S -> Sym Bool
type SToSymBool map[int64]bool

func (m SToSymBool) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualS()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map S -> Sym F
type SToSymF map[int64]float64

func (m SToSymF) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualS()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map S -> Sym S
type SToSymS map[int64]int64

func (m SToSymS) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualS()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map S -> Sym Str
type SToSymStr map[int64]string

func (m SToSymStr) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualS()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map S -> Sym U
type SToSymU map[int64]uint64

func (m SToSymU) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualS()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str -> Scalar
type StrToScalar map[string]S

func (m StrToScalar) MapScalar(s S) (S, error) {
	a := s.ActualStr()
	if ns, ok := m[a]; ok {
		ns.Actual = a
		s = ns
	}
	return s, nil
}

// Map Str -> Description
type StrToDescription map[string]string

func (m StrToDescription) MapScalar(s S) (S, error) {
	a := s.ActualStr()
	if d, ok := m[a]; ok {
		s.Description = d
	}
	return s, nil
}

// Map Str -> Sym Bool
type StrToSymBool map[string]bool

func (m StrToSymBool) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualStr()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str -> Sym F
type StrToSymF map[string]float64

func (m StrToSymF) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualStr()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str -> Sym S
type StrToSymS map[string]int64

func (m StrToSymS) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualStr()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str -> Sym Str
type StrToSymStr map[string]string

func (m StrToSymStr) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualStr()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str -> Sym U
type StrToSymU map[string]uint64

func (m StrToSymU) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualStr()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map U -> Scalar
type UToScalar map[uint64]S

func (m UToScalar) MapScalar(s S) (S, error) {
	a := s.ActualU()
	if ns, ok := m[a]; ok {
		ns.Actual = a
		s = ns
	}
	return s, nil
}

// Map U -> Description
type UToDescription map[uint64]string

func (m UToDescription) MapScalar(s S) (S, error) {
	a := s.ActualU()
	if d, ok := m[a]; ok {
		s.Description = d
	}
	return s, nil
}

// Map U -> Sym Bool
type UToSymBool map[uint64]bool

func (m UToSymBool) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualU()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map U -> Sym F
type UToSymF map[uint64]float64

func (m UToSymF) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualU()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map U -> Sym S
type UToSymS map[uint64]int64

func (m UToSymS) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualU()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map U -> Sym Str
type UToSymStr map[uint64]string

func (m UToSymStr) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualU()]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map U -> Sym U
type UToSymU map[uint64]uint64

func (m UToSymU) MapScalar(s S) (S, error) {
	if t, ok := m[s.ActualU()]; ok {
		s.Sym = t
	}
	return s, nil
}
