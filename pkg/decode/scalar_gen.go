// Code below generated from scalar_gen.go.tmpl
package decode

import (
	"errors"
	"fmt"

	"github.com/wader/fq/pkg/bitio"
)

// Type BitBuf

func (s Scalar) ActualBitBuf() *bitio.Buffer {
	v, ok := s.Actual.(*bitio.Buffer)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as *bitio.Buffer", s.Actual))
	}
	return v
}
func (s Scalar) SymBitBuf() *bitio.Buffer {
	v, ok := s.Sym.(*bitio.Buffer)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as *bitio.Buffer", s.Sym))
	}
	return v
}

func (d *D) FieldBitBufScalarFn(name string, fn func(d *D) Scalar, sfns ...ScalarFn) *bitio.Buffer {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d), nil }, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "BitBuf", Pos: d.Pos()})
	}
	return v.ActualBitBuf()
}
func (d *D) FieldBitBufFn(name string, fn func(d *D) *bitio.Buffer, sfns ...ScalarFn) *bitio.Buffer {
	return d.FieldBitBufScalarFn(name, func(d *D) Scalar { return Scalar{Actual: fn(d)} }, sfns...)
}
func (d *D) TryFieldBitBufScalarFn(name string, fn func(d *D) (Scalar, error), sfns ...ScalarFn) (*bitio.Buffer, error) {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d) }, sfns...)
	if err != nil {
		return nil, err
	}
	return v.ActualBitBuf(), err
}
func (d *D) TryFieldBitBufFn(name string, fn func(d *D) (*bitio.Buffer, error), sfns ...ScalarFn) (*bitio.Buffer, error) {
	return d.TryFieldBitBufScalarFn(name, func(d *D) (Scalar, error) {
		v, err := fn(d)
		return Scalar{Actual: v}, err
	}, sfns...)
}

// Type Bool

func (s Scalar) ActualBool() bool {
	v, ok := s.Actual.(bool)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as bool", s.Actual))
	}
	return v
}
func (s Scalar) SymBool() bool {
	v, ok := s.Sym.(bool)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as bool", s.Sym))
	}
	return v
}

func (d *D) FieldBoolScalarFn(name string, fn func(d *D) Scalar, sfns ...ScalarFn) bool {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d), nil }, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Bool", Pos: d.Pos()})
	}
	return v.ActualBool()
}
func (d *D) FieldBoolFn(name string, fn func(d *D) bool, sfns ...ScalarFn) bool {
	return d.FieldBoolScalarFn(name, func(d *D) Scalar { return Scalar{Actual: fn(d)} }, sfns...)
}
func (d *D) TryFieldBoolScalarFn(name string, fn func(d *D) (Scalar, error), sfns ...ScalarFn) (bool, error) {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d) }, sfns...)
	if err != nil {
		return false, err
	}
	return v.ActualBool(), err
}
func (d *D) TryFieldBoolFn(name string, fn func(d *D) (bool, error), sfns ...ScalarFn) (bool, error) {
	return d.TryFieldBoolScalarFn(name, func(d *D) (Scalar, error) {
		v, err := fn(d)
		return Scalar{Actual: v}, err
	}, sfns...)
}

// Type F

func (s Scalar) ActualF() float64 {
	v, ok := s.Actual.(float64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as float64", s.Actual))
	}
	return v
}
func (s Scalar) SymF() float64 {
	v, ok := s.Sym.(float64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as float64", s.Sym))
	}
	return v
}

func (d *D) FieldFScalarFn(name string, fn func(d *D) Scalar, sfns ...ScalarFn) float64 {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d), nil }, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F", Pos: d.Pos()})
	}
	return v.ActualF()
}
func (d *D) FieldFFn(name string, fn func(d *D) float64, sfns ...ScalarFn) float64 {
	return d.FieldFScalarFn(name, func(d *D) Scalar { return Scalar{Actual: fn(d)} }, sfns...)
}
func (d *D) TryFieldFScalarFn(name string, fn func(d *D) (Scalar, error), sfns ...ScalarFn) (float64, error) {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d) }, sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualF(), err
}
func (d *D) TryFieldFFn(name string, fn func(d *D) (float64, error), sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFScalarFn(name, func(d *D) (Scalar, error) {
		v, err := fn(d)
		return Scalar{Actual: v}, err
	}, sfns...)
}

// Type S

func (s Scalar) ActualS() int64 {
	v, ok := s.Actual.(int64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as int64", s.Actual))
	}
	return v
}
func (s Scalar) SymS() int64 {
	v, ok := s.Sym.(int64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as int64", s.Sym))
	}
	return v
}

func (d *D) FieldSScalarFn(name string, fn func(d *D) Scalar, sfns ...ScalarFn) int64 {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d), nil }, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S", Pos: d.Pos()})
	}
	return v.ActualS()
}
func (d *D) FieldSFn(name string, fn func(d *D) int64, sfns ...ScalarFn) int64 {
	return d.FieldSScalarFn(name, func(d *D) Scalar { return Scalar{Actual: fn(d)} }, sfns...)
}
func (d *D) TryFieldSScalarFn(name string, fn func(d *D) (Scalar, error), sfns ...ScalarFn) (int64, error) {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d) }, sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualS(), err
}
func (d *D) TryFieldSFn(name string, fn func(d *D) (int64, error), sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSScalarFn(name, func(d *D) (Scalar, error) {
		v, err := fn(d)
		return Scalar{Actual: v}, err
	}, sfns...)
}

// Type Str

func (s Scalar) ActualStr() string {
	v, ok := s.Actual.(string)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as string", s.Actual))
	}
	return v
}
func (s Scalar) SymStr() string {
	v, ok := s.Sym.(string)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as string", s.Sym))
	}
	return v
}

func (d *D) FieldStrScalarFn(name string, fn func(d *D) Scalar, sfns ...ScalarFn) string {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d), nil }, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Str", Pos: d.Pos()})
	}
	return v.ActualStr()
}
func (d *D) FieldStrFn(name string, fn func(d *D) string, sfns ...ScalarFn) string {
	return d.FieldStrScalarFn(name, func(d *D) Scalar { return Scalar{Actual: fn(d)} }, sfns...)
}
func (d *D) TryFieldStrScalarFn(name string, fn func(d *D) (Scalar, error), sfns ...ScalarFn) (string, error) {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d) }, sfns...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}
func (d *D) TryFieldStrFn(name string, fn func(d *D) (string, error), sfns ...ScalarFn) (string, error) {
	return d.TryFieldStrScalarFn(name, func(d *D) (Scalar, error) {
		v, err := fn(d)
		return Scalar{Actual: v}, err
	}, sfns...)
}

// Type U

func (s Scalar) ActualU() uint64 {
	v, ok := s.Actual.(uint64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Actual %v as uint64", s.Actual))
	}
	return v
}
func (s Scalar) SymU() uint64 {
	v, ok := s.Sym.(uint64)
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v as uint64", s.Sym))
	}
	return v
}

func (d *D) FieldUScalarFn(name string, fn func(d *D) Scalar, sfns ...ScalarFn) uint64 {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d), nil }, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U", Pos: d.Pos()})
	}
	return v.ActualU()
}
func (d *D) FieldUFn(name string, fn func(d *D) uint64, sfns ...ScalarFn) uint64 {
	return d.FieldUScalarFn(name, func(d *D) Scalar { return Scalar{Actual: fn(d)} }, sfns...)
}
func (d *D) TryFieldUScalarFn(name string, fn func(d *D) (Scalar, error), sfns ...ScalarFn) (uint64, error) {
	v, err := d.TryFieldScalar(name, func(_ Scalar) (Scalar, error) { return fn(d) }, sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualU(), err
}
func (d *D) TryFieldUFn(name string, fn func(d *D) (uint64, error), sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUScalarFn(name, func(d *D) (Scalar, error) {
		v, err := fn(d)
		return Scalar{Actual: v}, err
	}, sfns...)
}

// Validate/Assert Bool
func (d *D) assertBool(assert bool, vs ...bool) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualBool()
		for _, b := range vs {
			if a == b {
				s.Description = "valid"
				return s, nil
			}
		}
		s.Description = "invalid"
		if assert && !d.Options.Force {
			return s, errors.New("failed to assert Bool")
		}
		return s, nil
	}
}

func (d *D) AssertBool(vs ...bool) func(s Scalar) (Scalar, error) {
	return d.assertBool(true, vs...)
}
func (d *D) ValidateBool(vs ...bool) func(s Scalar) (Scalar, error) {
	return d.assertBool(false, vs...)
}

// Validate/Assert F
func (d *D) assertF(assert bool, vs ...float64) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualF()
		for _, b := range vs {
			if a == b {
				s.Description = "valid"
				return s, nil
			}
		}
		s.Description = "invalid"
		if assert && !d.Options.Force {
			return s, errors.New("failed to assert F")
		}
		return s, nil
	}
}

func (d *D) AssertF(vs ...float64) func(s Scalar) (Scalar, error) {
	return d.assertF(true, vs...)
}
func (d *D) ValidateF(vs ...float64) func(s Scalar) (Scalar, error) {
	return d.assertF(false, vs...)
}

// Validate/Assert S
func (d *D) assertS(assert bool, vs ...int64) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualS()
		for _, b := range vs {
			if a == b {
				s.Description = "valid"
				return s, nil
			}
		}
		s.Description = "invalid"
		if assert && !d.Options.Force {
			return s, errors.New("failed to assert S")
		}
		return s, nil
	}
}

func (d *D) AssertS(vs ...int64) func(s Scalar) (Scalar, error) {
	return d.assertS(true, vs...)
}
func (d *D) ValidateS(vs ...int64) func(s Scalar) (Scalar, error) {
	return d.assertS(false, vs...)
}

// Validate/Assert Str
func (d *D) assertStr(assert bool, vs ...string) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualStr()
		for _, b := range vs {
			if a == b {
				s.Description = "valid"
				return s, nil
			}
		}
		s.Description = "invalid"
		if assert && !d.Options.Force {
			return s, errors.New("failed to assert Str")
		}
		return s, nil
	}
}

func (d *D) AssertStr(vs ...string) func(s Scalar) (Scalar, error) {
	return d.assertStr(true, vs...)
}
func (d *D) ValidateStr(vs ...string) func(s Scalar) (Scalar, error) {
	return d.assertStr(false, vs...)
}

// Validate/Assert U
func (d *D) assertU(assert bool, vs ...uint64) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualU()
		for _, b := range vs {
			if a == b {
				s.Description = "valid"
				return s, nil
			}
		}
		s.Description = "invalid"
		if assert && !d.Options.Force {
			return s, errors.New("failed to assert U")
		}
		return s, nil
	}
}

func (d *D) AssertU(vs ...uint64) func(s Scalar) (Scalar, error) {
	return d.assertU(true, vs...)
}
func (d *D) ValidateU(vs ...uint64) func(s Scalar) (Scalar, error) {
	return d.assertU(false, vs...)
}

// Map Bool -> Scalar

type BoolToScalar map[bool]Scalar

func (d *D) MapBoolToScalar(m BoolToScalar) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualBool()
		if ns, ok := m[a]; ok {
			ns.Actual = a
			s = ns
		}
		return s, nil
	}
}

// Map Bool -> Bool

type BoolToBool map[bool]bool

func (d *D) MapBoolToBool(m BoolToBool) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualBool()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Bool -> F

type BoolToF map[bool]float64

func (d *D) MapBoolToF(m BoolToF) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualBool()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Bool -> S

type BoolToS map[bool]int64

func (d *D) MapBoolToS(m BoolToS) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualBool()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Bool -> Str

type BoolToStr map[bool]string

func (d *D) MapBoolToStr(m BoolToStr) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualBool()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Bool -> U

type BoolToU map[bool]uint64

func (d *D) MapBoolToU(m BoolToU) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualBool()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map S -> Scalar

type SToScalar map[int64]Scalar

func (d *D) MapSToScalar(m SToScalar) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualS()
		if ns, ok := m[a]; ok {
			ns.Actual = a
			s = ns
		}
		return s, nil
	}
}

// Map S -> Bool

type SToBool map[int64]bool

func (d *D) MapSToBool(m SToBool) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualS()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map S -> F

type SToF map[int64]float64

func (d *D) MapSToF(m SToF) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualS()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map S -> S

type SToS map[int64]int64

func (d *D) MapSToS(m SToS) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualS()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map S -> Str

type SToStr map[int64]string

func (d *D) MapSToStr(m SToStr) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualS()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map S -> U

type SToU map[int64]uint64

func (d *D) MapSToU(m SToU) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualS()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Str -> Scalar

type StrToScalar map[string]Scalar

func (d *D) MapStrToScalar(m StrToScalar) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualStr()
		if ns, ok := m[a]; ok {
			ns.Actual = a
			s = ns
		}
		return s, nil
	}
}

// Map Str -> Bool

type StrToBool map[string]bool

func (d *D) MapStrToBool(m StrToBool) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualStr()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Str -> F

type StrToF map[string]float64

func (d *D) MapStrToF(m StrToF) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualStr()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Str -> S

type StrToS map[string]int64

func (d *D) MapStrToS(m StrToS) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualStr()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Str -> Str

type StrToStr map[string]string

func (d *D) MapStrToStr(m StrToStr) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualStr()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map Str -> U

type StrToU map[string]uint64

func (d *D) MapStrToU(m StrToU) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualStr()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map U -> Scalar

type UToScalar map[uint64]Scalar

func (d *D) MapUToScalar(m UToScalar) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		a := s.ActualU()
		if ns, ok := m[a]; ok {
			ns.Actual = a
			s = ns
		}
		return s, nil
	}
}

// Map U -> Bool

type UToBool map[uint64]bool

func (d *D) MapUToBool(m UToBool) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualU()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map U -> F

type UToF map[uint64]float64

func (d *D) MapUToF(m UToF) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualU()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map U -> S

type UToS map[uint64]int64

func (d *D) MapUToS(m UToS) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualU()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map U -> Str

type UToStr map[uint64]string

func (d *D) MapUToStr(m UToStr) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualU()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Map U -> U

type UToU map[uint64]uint64

func (d *D) MapUToU(m UToU) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if t, ok := m[s.ActualU()]; ok {
			s.Sym = t
		}
		return s, nil
	}
}

// Reader RawLen

func (d *D) TryRawLen(nBits int64) (*bitio.Buffer, error) { return d.tryBitBuf(nBits) }

func (d *D) ScalarRawLen(nBits int64) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryBitBuf(nBits)
		s.Actual = v
		return s, err
	}
}

func (d *D) RawLen(nBits int64) *bitio.Buffer {
	v, err := d.tryBitBuf(nBits)
	if err != nil {
		panic(IOError{Err: err, Op: "RawLen", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldRawLen(name string, nBits int64, sfns ...ScalarFn) (*bitio.Buffer, error) {
	v, err := d.TryFieldScalar(name, d.ScalarRawLen(nBits), sfns...)
	if err != nil {
		return nil, err
	}
	return v.ActualBitBuf(), err
}

func (d *D) FieldRawLen(name string, nBits int64, sfns ...ScalarFn) *bitio.Buffer {
	v, err := d.TryFieldRawLen(name, nBits, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "RawLen", Pos: d.Pos()})
	}
	return v
}

// Reader Bool

func (d *D) TryBool() (bool, error) { return d.tryBool() }

func (d *D) ScalarBool() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryBool()
		s.Actual = v
		return s, err
	}
}

func (d *D) Bool() bool {
	v, err := d.tryBool()
	if err != nil {
		panic(IOError{Err: err, Op: "Bool", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldBool(name string, sfns ...ScalarFn) (bool, error) {
	return d.TryFieldBoolFn(name, (*D).TryBool, sfns...)
}

func (d *D) FieldBool(name string, sfns ...ScalarFn) bool {
	return d.FieldBoolFn(name, (*D).Bool, sfns...)
}

// Reader UE

func (d *D) TryUE(nBits int, endian Endian) (uint64, error) { return d.tryUE(nBits, endian) }

func (d *D) ScalarUE(nBits int, endian Endian) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(nBits, endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) UE(nBits int, endian Endian) uint64 {
	v, err := d.tryUE(nBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "UE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUE(name string, nBits int, endian Endian, sfns ...ScalarFn) (uint64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUE(nBits, endian), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualU(), err
}

func (d *D) FieldUE(name string, nBits int, endian Endian, sfns ...ScalarFn) uint64 {
	v, err := d.TryFieldUE(name, nBits, endian, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UE", Pos: d.Pos()})
	}
	return v
}

// Reader U

func (d *D) TryU(nBits int) (uint64, error) { return d.tryUE(nBits, d.Endian) }

func (d *D) ScalarU(nBits int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(nBits, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U(nBits int) uint64 {
	v, err := d.tryUE(nBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU(name string, nBits int, sfns ...ScalarFn) (uint64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarU(nBits), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualU(), err
}

func (d *D) FieldU(name string, nBits int, sfns ...ScalarFn) uint64 {
	v, err := d.TryFieldU(name, nBits, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "U", Pos: d.Pos()})
	}
	return v
}

// Reader U1

func (d *D) TryU1() (uint64, error) { return d.tryUE(1, d.Endian) }

func (d *D) ScalarU1() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(1, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U1() uint64 {
	v, err := d.tryUE(1, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U1", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU1(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU1, sfns...)
}

func (d *D) FieldU1(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U1, sfns...)
}

// Reader U2

func (d *D) TryU2() (uint64, error) { return d.tryUE(2, d.Endian) }

func (d *D) ScalarU2() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(2, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U2() uint64 {
	v, err := d.tryUE(2, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U2", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU2(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU2, sfns...)
}

func (d *D) FieldU2(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U2, sfns...)
}

// Reader U3

func (d *D) TryU3() (uint64, error) { return d.tryUE(3, d.Endian) }

func (d *D) ScalarU3() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(3, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U3() uint64 {
	v, err := d.tryUE(3, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U3", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU3(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU3, sfns...)
}

func (d *D) FieldU3(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U3, sfns...)
}

// Reader U4

func (d *D) TryU4() (uint64, error) { return d.tryUE(4, d.Endian) }

func (d *D) ScalarU4() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(4, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U4() uint64 {
	v, err := d.tryUE(4, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U4", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU4(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU4, sfns...)
}

func (d *D) FieldU4(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U4, sfns...)
}

// Reader U5

func (d *D) TryU5() (uint64, error) { return d.tryUE(5, d.Endian) }

func (d *D) ScalarU5() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(5, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U5() uint64 {
	v, err := d.tryUE(5, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U5", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU5(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU5, sfns...)
}

func (d *D) FieldU5(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U5, sfns...)
}

// Reader U6

func (d *D) TryU6() (uint64, error) { return d.tryUE(6, d.Endian) }

func (d *D) ScalarU6() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(6, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U6() uint64 {
	v, err := d.tryUE(6, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U6", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU6(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU6, sfns...)
}

func (d *D) FieldU6(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U6, sfns...)
}

// Reader U7

func (d *D) TryU7() (uint64, error) { return d.tryUE(7, d.Endian) }

func (d *D) ScalarU7() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(7, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U7() uint64 {
	v, err := d.tryUE(7, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U7", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU7(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU7, sfns...)
}

func (d *D) FieldU7(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U7, sfns...)
}

// Reader U8

func (d *D) TryU8() (uint64, error) { return d.tryUE(8, d.Endian) }

func (d *D) ScalarU8() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(8, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U8() uint64 {
	v, err := d.tryUE(8, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U8", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU8(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU8, sfns...)
}

func (d *D) FieldU8(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U8, sfns...)
}

// Reader U9

func (d *D) TryU9() (uint64, error) { return d.tryUE(9, d.Endian) }

func (d *D) ScalarU9() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(9, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U9() uint64 {
	v, err := d.tryUE(9, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U9", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU9(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU9, sfns...)
}

func (d *D) FieldU9(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U9, sfns...)
}

// Reader U10

func (d *D) TryU10() (uint64, error) { return d.tryUE(10, d.Endian) }

func (d *D) ScalarU10() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(10, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U10() uint64 {
	v, err := d.tryUE(10, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U10", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU10(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU10, sfns...)
}

func (d *D) FieldU10(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U10, sfns...)
}

// Reader U11

func (d *D) TryU11() (uint64, error) { return d.tryUE(11, d.Endian) }

func (d *D) ScalarU11() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(11, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U11() uint64 {
	v, err := d.tryUE(11, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U11", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU11(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU11, sfns...)
}

func (d *D) FieldU11(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U11, sfns...)
}

// Reader U12

func (d *D) TryU12() (uint64, error) { return d.tryUE(12, d.Endian) }

func (d *D) ScalarU12() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(12, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U12() uint64 {
	v, err := d.tryUE(12, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U12", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU12(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU12, sfns...)
}

func (d *D) FieldU12(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U12, sfns...)
}

// Reader U13

func (d *D) TryU13() (uint64, error) { return d.tryUE(13, d.Endian) }

func (d *D) ScalarU13() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(13, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U13() uint64 {
	v, err := d.tryUE(13, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U13", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU13(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU13, sfns...)
}

func (d *D) FieldU13(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U13, sfns...)
}

// Reader U14

func (d *D) TryU14() (uint64, error) { return d.tryUE(14, d.Endian) }

func (d *D) ScalarU14() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(14, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U14() uint64 {
	v, err := d.tryUE(14, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U14", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU14(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU14, sfns...)
}

func (d *D) FieldU14(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U14, sfns...)
}

// Reader U15

func (d *D) TryU15() (uint64, error) { return d.tryUE(15, d.Endian) }

func (d *D) ScalarU15() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(15, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U15() uint64 {
	v, err := d.tryUE(15, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U15", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU15(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU15, sfns...)
}

func (d *D) FieldU15(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U15, sfns...)
}

// Reader U16

func (d *D) TryU16() (uint64, error) { return d.tryUE(16, d.Endian) }

func (d *D) ScalarU16() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(16, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U16() uint64 {
	v, err := d.tryUE(16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U16", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU16(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU16, sfns...)
}

func (d *D) FieldU16(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U16, sfns...)
}

// Reader U17

func (d *D) TryU17() (uint64, error) { return d.tryUE(17, d.Endian) }

func (d *D) ScalarU17() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(17, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U17() uint64 {
	v, err := d.tryUE(17, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U17", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU17(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU17, sfns...)
}

func (d *D) FieldU17(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U17, sfns...)
}

// Reader U18

func (d *D) TryU18() (uint64, error) { return d.tryUE(18, d.Endian) }

func (d *D) ScalarU18() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(18, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U18() uint64 {
	v, err := d.tryUE(18, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U18", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU18(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU18, sfns...)
}

func (d *D) FieldU18(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U18, sfns...)
}

// Reader U19

func (d *D) TryU19() (uint64, error) { return d.tryUE(19, d.Endian) }

func (d *D) ScalarU19() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(19, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U19() uint64 {
	v, err := d.tryUE(19, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U19", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU19(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU19, sfns...)
}

func (d *D) FieldU19(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U19, sfns...)
}

// Reader U20

func (d *D) TryU20() (uint64, error) { return d.tryUE(20, d.Endian) }

func (d *D) ScalarU20() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(20, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U20() uint64 {
	v, err := d.tryUE(20, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U20", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU20(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU20, sfns...)
}

func (d *D) FieldU20(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U20, sfns...)
}

// Reader U21

func (d *D) TryU21() (uint64, error) { return d.tryUE(21, d.Endian) }

func (d *D) ScalarU21() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(21, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U21() uint64 {
	v, err := d.tryUE(21, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U21", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU21(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU21, sfns...)
}

func (d *D) FieldU21(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U21, sfns...)
}

// Reader U22

func (d *D) TryU22() (uint64, error) { return d.tryUE(22, d.Endian) }

func (d *D) ScalarU22() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(22, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U22() uint64 {
	v, err := d.tryUE(22, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U22", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU22(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU22, sfns...)
}

func (d *D) FieldU22(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U22, sfns...)
}

// Reader U23

func (d *D) TryU23() (uint64, error) { return d.tryUE(23, d.Endian) }

func (d *D) ScalarU23() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(23, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U23() uint64 {
	v, err := d.tryUE(23, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U23", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU23(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU23, sfns...)
}

func (d *D) FieldU23(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U23, sfns...)
}

// Reader U24

func (d *D) TryU24() (uint64, error) { return d.tryUE(24, d.Endian) }

func (d *D) ScalarU24() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(24, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U24() uint64 {
	v, err := d.tryUE(24, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U24", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU24(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU24, sfns...)
}

func (d *D) FieldU24(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U24, sfns...)
}

// Reader U25

func (d *D) TryU25() (uint64, error) { return d.tryUE(25, d.Endian) }

func (d *D) ScalarU25() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(25, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U25() uint64 {
	v, err := d.tryUE(25, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U25", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU25(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU25, sfns...)
}

func (d *D) FieldU25(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U25, sfns...)
}

// Reader U26

func (d *D) TryU26() (uint64, error) { return d.tryUE(26, d.Endian) }

func (d *D) ScalarU26() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(26, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U26() uint64 {
	v, err := d.tryUE(26, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U26", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU26(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU26, sfns...)
}

func (d *D) FieldU26(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U26, sfns...)
}

// Reader U27

func (d *D) TryU27() (uint64, error) { return d.tryUE(27, d.Endian) }

func (d *D) ScalarU27() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(27, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U27() uint64 {
	v, err := d.tryUE(27, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U27", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU27(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU27, sfns...)
}

func (d *D) FieldU27(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U27, sfns...)
}

// Reader U28

func (d *D) TryU28() (uint64, error) { return d.tryUE(28, d.Endian) }

func (d *D) ScalarU28() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(28, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U28() uint64 {
	v, err := d.tryUE(28, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U28", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU28(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU28, sfns...)
}

func (d *D) FieldU28(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U28, sfns...)
}

// Reader U29

func (d *D) TryU29() (uint64, error) { return d.tryUE(29, d.Endian) }

func (d *D) ScalarU29() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(29, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U29() uint64 {
	v, err := d.tryUE(29, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U29", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU29(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU29, sfns...)
}

func (d *D) FieldU29(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U29, sfns...)
}

// Reader U30

func (d *D) TryU30() (uint64, error) { return d.tryUE(30, d.Endian) }

func (d *D) ScalarU30() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(30, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U30() uint64 {
	v, err := d.tryUE(30, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U30", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU30(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU30, sfns...)
}

func (d *D) FieldU30(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U30, sfns...)
}

// Reader U31

func (d *D) TryU31() (uint64, error) { return d.tryUE(31, d.Endian) }

func (d *D) ScalarU31() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(31, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U31() uint64 {
	v, err := d.tryUE(31, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U31", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU31(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU31, sfns...)
}

func (d *D) FieldU31(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U31, sfns...)
}

// Reader U32

func (d *D) TryU32() (uint64, error) { return d.tryUE(32, d.Endian) }

func (d *D) ScalarU32() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(32, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U32() uint64 {
	v, err := d.tryUE(32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U32", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU32(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU32, sfns...)
}

func (d *D) FieldU32(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U32, sfns...)
}

// Reader U33

func (d *D) TryU33() (uint64, error) { return d.tryUE(33, d.Endian) }

func (d *D) ScalarU33() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(33, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U33() uint64 {
	v, err := d.tryUE(33, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U33", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU33(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU33, sfns...)
}

func (d *D) FieldU33(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U33, sfns...)
}

// Reader U34

func (d *D) TryU34() (uint64, error) { return d.tryUE(34, d.Endian) }

func (d *D) ScalarU34() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(34, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U34() uint64 {
	v, err := d.tryUE(34, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U34", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU34(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU34, sfns...)
}

func (d *D) FieldU34(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U34, sfns...)
}

// Reader U35

func (d *D) TryU35() (uint64, error) { return d.tryUE(35, d.Endian) }

func (d *D) ScalarU35() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(35, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U35() uint64 {
	v, err := d.tryUE(35, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U35", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU35(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU35, sfns...)
}

func (d *D) FieldU35(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U35, sfns...)
}

// Reader U36

func (d *D) TryU36() (uint64, error) { return d.tryUE(36, d.Endian) }

func (d *D) ScalarU36() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(36, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U36() uint64 {
	v, err := d.tryUE(36, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U36", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU36(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU36, sfns...)
}

func (d *D) FieldU36(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U36, sfns...)
}

// Reader U37

func (d *D) TryU37() (uint64, error) { return d.tryUE(37, d.Endian) }

func (d *D) ScalarU37() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(37, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U37() uint64 {
	v, err := d.tryUE(37, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U37", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU37(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU37, sfns...)
}

func (d *D) FieldU37(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U37, sfns...)
}

// Reader U38

func (d *D) TryU38() (uint64, error) { return d.tryUE(38, d.Endian) }

func (d *D) ScalarU38() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(38, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U38() uint64 {
	v, err := d.tryUE(38, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U38", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU38(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU38, sfns...)
}

func (d *D) FieldU38(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U38, sfns...)
}

// Reader U39

func (d *D) TryU39() (uint64, error) { return d.tryUE(39, d.Endian) }

func (d *D) ScalarU39() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(39, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U39() uint64 {
	v, err := d.tryUE(39, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U39", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU39(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU39, sfns...)
}

func (d *D) FieldU39(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U39, sfns...)
}

// Reader U40

func (d *D) TryU40() (uint64, error) { return d.tryUE(40, d.Endian) }

func (d *D) ScalarU40() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(40, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U40() uint64 {
	v, err := d.tryUE(40, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U40", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU40(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU40, sfns...)
}

func (d *D) FieldU40(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U40, sfns...)
}

// Reader U41

func (d *D) TryU41() (uint64, error) { return d.tryUE(41, d.Endian) }

func (d *D) ScalarU41() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(41, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U41() uint64 {
	v, err := d.tryUE(41, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U41", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU41(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU41, sfns...)
}

func (d *D) FieldU41(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U41, sfns...)
}

// Reader U42

func (d *D) TryU42() (uint64, error) { return d.tryUE(42, d.Endian) }

func (d *D) ScalarU42() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(42, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U42() uint64 {
	v, err := d.tryUE(42, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U42", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU42(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU42, sfns...)
}

func (d *D) FieldU42(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U42, sfns...)
}

// Reader U43

func (d *D) TryU43() (uint64, error) { return d.tryUE(43, d.Endian) }

func (d *D) ScalarU43() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(43, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U43() uint64 {
	v, err := d.tryUE(43, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U43", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU43(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU43, sfns...)
}

func (d *D) FieldU43(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U43, sfns...)
}

// Reader U44

func (d *D) TryU44() (uint64, error) { return d.tryUE(44, d.Endian) }

func (d *D) ScalarU44() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(44, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U44() uint64 {
	v, err := d.tryUE(44, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U44", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU44(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU44, sfns...)
}

func (d *D) FieldU44(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U44, sfns...)
}

// Reader U45

func (d *D) TryU45() (uint64, error) { return d.tryUE(45, d.Endian) }

func (d *D) ScalarU45() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(45, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U45() uint64 {
	v, err := d.tryUE(45, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U45", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU45(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU45, sfns...)
}

func (d *D) FieldU45(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U45, sfns...)
}

// Reader U46

func (d *D) TryU46() (uint64, error) { return d.tryUE(46, d.Endian) }

func (d *D) ScalarU46() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(46, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U46() uint64 {
	v, err := d.tryUE(46, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U46", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU46(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU46, sfns...)
}

func (d *D) FieldU46(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U46, sfns...)
}

// Reader U47

func (d *D) TryU47() (uint64, error) { return d.tryUE(47, d.Endian) }

func (d *D) ScalarU47() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(47, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U47() uint64 {
	v, err := d.tryUE(47, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U47", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU47(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU47, sfns...)
}

func (d *D) FieldU47(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U47, sfns...)
}

// Reader U48

func (d *D) TryU48() (uint64, error) { return d.tryUE(48, d.Endian) }

func (d *D) ScalarU48() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(48, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U48() uint64 {
	v, err := d.tryUE(48, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U48", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU48(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU48, sfns...)
}

func (d *D) FieldU48(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U48, sfns...)
}

// Reader U49

func (d *D) TryU49() (uint64, error) { return d.tryUE(49, d.Endian) }

func (d *D) ScalarU49() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(49, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U49() uint64 {
	v, err := d.tryUE(49, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U49", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU49(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU49, sfns...)
}

func (d *D) FieldU49(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U49, sfns...)
}

// Reader U50

func (d *D) TryU50() (uint64, error) { return d.tryUE(50, d.Endian) }

func (d *D) ScalarU50() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(50, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U50() uint64 {
	v, err := d.tryUE(50, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U50", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU50(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU50, sfns...)
}

func (d *D) FieldU50(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U50, sfns...)
}

// Reader U51

func (d *D) TryU51() (uint64, error) { return d.tryUE(51, d.Endian) }

func (d *D) ScalarU51() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(51, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U51() uint64 {
	v, err := d.tryUE(51, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U51", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU51(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU51, sfns...)
}

func (d *D) FieldU51(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U51, sfns...)
}

// Reader U52

func (d *D) TryU52() (uint64, error) { return d.tryUE(52, d.Endian) }

func (d *D) ScalarU52() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(52, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U52() uint64 {
	v, err := d.tryUE(52, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U52", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU52(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU52, sfns...)
}

func (d *D) FieldU52(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U52, sfns...)
}

// Reader U53

func (d *D) TryU53() (uint64, error) { return d.tryUE(53, d.Endian) }

func (d *D) ScalarU53() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(53, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U53() uint64 {
	v, err := d.tryUE(53, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U53", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU53(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU53, sfns...)
}

func (d *D) FieldU53(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U53, sfns...)
}

// Reader U54

func (d *D) TryU54() (uint64, error) { return d.tryUE(54, d.Endian) }

func (d *D) ScalarU54() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(54, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U54() uint64 {
	v, err := d.tryUE(54, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U54", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU54(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU54, sfns...)
}

func (d *D) FieldU54(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U54, sfns...)
}

// Reader U55

func (d *D) TryU55() (uint64, error) { return d.tryUE(55, d.Endian) }

func (d *D) ScalarU55() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(55, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U55() uint64 {
	v, err := d.tryUE(55, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U55", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU55(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU55, sfns...)
}

func (d *D) FieldU55(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U55, sfns...)
}

// Reader U56

func (d *D) TryU56() (uint64, error) { return d.tryUE(56, d.Endian) }

func (d *D) ScalarU56() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(56, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U56() uint64 {
	v, err := d.tryUE(56, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U56", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU56(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU56, sfns...)
}

func (d *D) FieldU56(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U56, sfns...)
}

// Reader U57

func (d *D) TryU57() (uint64, error) { return d.tryUE(57, d.Endian) }

func (d *D) ScalarU57() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(57, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U57() uint64 {
	v, err := d.tryUE(57, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U57", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU57(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU57, sfns...)
}

func (d *D) FieldU57(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U57, sfns...)
}

// Reader U58

func (d *D) TryU58() (uint64, error) { return d.tryUE(58, d.Endian) }

func (d *D) ScalarU58() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(58, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U58() uint64 {
	v, err := d.tryUE(58, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U58", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU58(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU58, sfns...)
}

func (d *D) FieldU58(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U58, sfns...)
}

// Reader U59

func (d *D) TryU59() (uint64, error) { return d.tryUE(59, d.Endian) }

func (d *D) ScalarU59() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(59, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U59() uint64 {
	v, err := d.tryUE(59, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U59", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU59(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU59, sfns...)
}

func (d *D) FieldU59(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U59, sfns...)
}

// Reader U60

func (d *D) TryU60() (uint64, error) { return d.tryUE(60, d.Endian) }

func (d *D) ScalarU60() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(60, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U60() uint64 {
	v, err := d.tryUE(60, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U60", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU60(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU60, sfns...)
}

func (d *D) FieldU60(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U60, sfns...)
}

// Reader U61

func (d *D) TryU61() (uint64, error) { return d.tryUE(61, d.Endian) }

func (d *D) ScalarU61() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(61, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U61() uint64 {
	v, err := d.tryUE(61, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U61", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU61(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU61, sfns...)
}

func (d *D) FieldU61(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U61, sfns...)
}

// Reader U62

func (d *D) TryU62() (uint64, error) { return d.tryUE(62, d.Endian) }

func (d *D) ScalarU62() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(62, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U62() uint64 {
	v, err := d.tryUE(62, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U62", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU62(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU62, sfns...)
}

func (d *D) FieldU62(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U62, sfns...)
}

// Reader U63

func (d *D) TryU63() (uint64, error) { return d.tryUE(63, d.Endian) }

func (d *D) ScalarU63() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(63, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U63() uint64 {
	v, err := d.tryUE(63, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U63", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU63(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU63, sfns...)
}

func (d *D) FieldU63(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U63, sfns...)
}

// Reader U64

func (d *D) TryU64() (uint64, error) { return d.tryUE(64, d.Endian) }

func (d *D) ScalarU64() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(64, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U64() uint64 {
	v, err := d.tryUE(64, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "U64", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU64(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU64, sfns...)
}

func (d *D) FieldU64(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U64, sfns...)
}

// Reader U8LE

func (d *D) TryU8LE() (uint64, error) { return d.tryUE(8, LittleEndian) }

func (d *D) ScalarU8LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(8, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U8LE() uint64 {
	v, err := d.tryUE(8, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U8LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU8LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU8LE, sfns...)
}

func (d *D) FieldU8LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U8LE, sfns...)
}

// Reader U9LE

func (d *D) TryU9LE() (uint64, error) { return d.tryUE(9, LittleEndian) }

func (d *D) ScalarU9LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(9, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U9LE() uint64 {
	v, err := d.tryUE(9, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U9LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU9LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU9LE, sfns...)
}

func (d *D) FieldU9LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U9LE, sfns...)
}

// Reader U10LE

func (d *D) TryU10LE() (uint64, error) { return d.tryUE(10, LittleEndian) }

func (d *D) ScalarU10LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(10, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U10LE() uint64 {
	v, err := d.tryUE(10, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U10LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU10LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU10LE, sfns...)
}

func (d *D) FieldU10LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U10LE, sfns...)
}

// Reader U11LE

func (d *D) TryU11LE() (uint64, error) { return d.tryUE(11, LittleEndian) }

func (d *D) ScalarU11LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(11, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U11LE() uint64 {
	v, err := d.tryUE(11, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U11LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU11LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU11LE, sfns...)
}

func (d *D) FieldU11LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U11LE, sfns...)
}

// Reader U12LE

func (d *D) TryU12LE() (uint64, error) { return d.tryUE(12, LittleEndian) }

func (d *D) ScalarU12LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(12, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U12LE() uint64 {
	v, err := d.tryUE(12, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U12LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU12LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU12LE, sfns...)
}

func (d *D) FieldU12LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U12LE, sfns...)
}

// Reader U13LE

func (d *D) TryU13LE() (uint64, error) { return d.tryUE(13, LittleEndian) }

func (d *D) ScalarU13LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(13, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U13LE() uint64 {
	v, err := d.tryUE(13, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U13LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU13LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU13LE, sfns...)
}

func (d *D) FieldU13LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U13LE, sfns...)
}

// Reader U14LE

func (d *D) TryU14LE() (uint64, error) { return d.tryUE(14, LittleEndian) }

func (d *D) ScalarU14LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(14, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U14LE() uint64 {
	v, err := d.tryUE(14, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U14LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU14LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU14LE, sfns...)
}

func (d *D) FieldU14LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U14LE, sfns...)
}

// Reader U15LE

func (d *D) TryU15LE() (uint64, error) { return d.tryUE(15, LittleEndian) }

func (d *D) ScalarU15LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(15, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U15LE() uint64 {
	v, err := d.tryUE(15, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U15LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU15LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU15LE, sfns...)
}

func (d *D) FieldU15LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U15LE, sfns...)
}

// Reader U16LE

func (d *D) TryU16LE() (uint64, error) { return d.tryUE(16, LittleEndian) }

func (d *D) ScalarU16LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(16, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U16LE() uint64 {
	v, err := d.tryUE(16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U16LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU16LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU16LE, sfns...)
}

func (d *D) FieldU16LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U16LE, sfns...)
}

// Reader U17LE

func (d *D) TryU17LE() (uint64, error) { return d.tryUE(17, LittleEndian) }

func (d *D) ScalarU17LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(17, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U17LE() uint64 {
	v, err := d.tryUE(17, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U17LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU17LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU17LE, sfns...)
}

func (d *D) FieldU17LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U17LE, sfns...)
}

// Reader U18LE

func (d *D) TryU18LE() (uint64, error) { return d.tryUE(18, LittleEndian) }

func (d *D) ScalarU18LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(18, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U18LE() uint64 {
	v, err := d.tryUE(18, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U18LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU18LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU18LE, sfns...)
}

func (d *D) FieldU18LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U18LE, sfns...)
}

// Reader U19LE

func (d *D) TryU19LE() (uint64, error) { return d.tryUE(19, LittleEndian) }

func (d *D) ScalarU19LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(19, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U19LE() uint64 {
	v, err := d.tryUE(19, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U19LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU19LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU19LE, sfns...)
}

func (d *D) FieldU19LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U19LE, sfns...)
}

// Reader U20LE

func (d *D) TryU20LE() (uint64, error) { return d.tryUE(20, LittleEndian) }

func (d *D) ScalarU20LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(20, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U20LE() uint64 {
	v, err := d.tryUE(20, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U20LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU20LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU20LE, sfns...)
}

func (d *D) FieldU20LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U20LE, sfns...)
}

// Reader U21LE

func (d *D) TryU21LE() (uint64, error) { return d.tryUE(21, LittleEndian) }

func (d *D) ScalarU21LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(21, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U21LE() uint64 {
	v, err := d.tryUE(21, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U21LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU21LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU21LE, sfns...)
}

func (d *D) FieldU21LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U21LE, sfns...)
}

// Reader U22LE

func (d *D) TryU22LE() (uint64, error) { return d.tryUE(22, LittleEndian) }

func (d *D) ScalarU22LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(22, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U22LE() uint64 {
	v, err := d.tryUE(22, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U22LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU22LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU22LE, sfns...)
}

func (d *D) FieldU22LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U22LE, sfns...)
}

// Reader U23LE

func (d *D) TryU23LE() (uint64, error) { return d.tryUE(23, LittleEndian) }

func (d *D) ScalarU23LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(23, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U23LE() uint64 {
	v, err := d.tryUE(23, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U23LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU23LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU23LE, sfns...)
}

func (d *D) FieldU23LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U23LE, sfns...)
}

// Reader U24LE

func (d *D) TryU24LE() (uint64, error) { return d.tryUE(24, LittleEndian) }

func (d *D) ScalarU24LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(24, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U24LE() uint64 {
	v, err := d.tryUE(24, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U24LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU24LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU24LE, sfns...)
}

func (d *D) FieldU24LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U24LE, sfns...)
}

// Reader U25LE

func (d *D) TryU25LE() (uint64, error) { return d.tryUE(25, LittleEndian) }

func (d *D) ScalarU25LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(25, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U25LE() uint64 {
	v, err := d.tryUE(25, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U25LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU25LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU25LE, sfns...)
}

func (d *D) FieldU25LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U25LE, sfns...)
}

// Reader U26LE

func (d *D) TryU26LE() (uint64, error) { return d.tryUE(26, LittleEndian) }

func (d *D) ScalarU26LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(26, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U26LE() uint64 {
	v, err := d.tryUE(26, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U26LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU26LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU26LE, sfns...)
}

func (d *D) FieldU26LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U26LE, sfns...)
}

// Reader U27LE

func (d *D) TryU27LE() (uint64, error) { return d.tryUE(27, LittleEndian) }

func (d *D) ScalarU27LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(27, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U27LE() uint64 {
	v, err := d.tryUE(27, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U27LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU27LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU27LE, sfns...)
}

func (d *D) FieldU27LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U27LE, sfns...)
}

// Reader U28LE

func (d *D) TryU28LE() (uint64, error) { return d.tryUE(28, LittleEndian) }

func (d *D) ScalarU28LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(28, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U28LE() uint64 {
	v, err := d.tryUE(28, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U28LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU28LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU28LE, sfns...)
}

func (d *D) FieldU28LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U28LE, sfns...)
}

// Reader U29LE

func (d *D) TryU29LE() (uint64, error) { return d.tryUE(29, LittleEndian) }

func (d *D) ScalarU29LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(29, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U29LE() uint64 {
	v, err := d.tryUE(29, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U29LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU29LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU29LE, sfns...)
}

func (d *D) FieldU29LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U29LE, sfns...)
}

// Reader U30LE

func (d *D) TryU30LE() (uint64, error) { return d.tryUE(30, LittleEndian) }

func (d *D) ScalarU30LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(30, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U30LE() uint64 {
	v, err := d.tryUE(30, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U30LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU30LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU30LE, sfns...)
}

func (d *D) FieldU30LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U30LE, sfns...)
}

// Reader U31LE

func (d *D) TryU31LE() (uint64, error) { return d.tryUE(31, LittleEndian) }

func (d *D) ScalarU31LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(31, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U31LE() uint64 {
	v, err := d.tryUE(31, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U31LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU31LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU31LE, sfns...)
}

func (d *D) FieldU31LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U31LE, sfns...)
}

// Reader U32LE

func (d *D) TryU32LE() (uint64, error) { return d.tryUE(32, LittleEndian) }

func (d *D) ScalarU32LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(32, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U32LE() uint64 {
	v, err := d.tryUE(32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U32LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU32LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU32LE, sfns...)
}

func (d *D) FieldU32LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U32LE, sfns...)
}

// Reader U33LE

func (d *D) TryU33LE() (uint64, error) { return d.tryUE(33, LittleEndian) }

func (d *D) ScalarU33LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(33, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U33LE() uint64 {
	v, err := d.tryUE(33, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U33LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU33LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU33LE, sfns...)
}

func (d *D) FieldU33LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U33LE, sfns...)
}

// Reader U34LE

func (d *D) TryU34LE() (uint64, error) { return d.tryUE(34, LittleEndian) }

func (d *D) ScalarU34LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(34, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U34LE() uint64 {
	v, err := d.tryUE(34, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U34LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU34LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU34LE, sfns...)
}

func (d *D) FieldU34LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U34LE, sfns...)
}

// Reader U35LE

func (d *D) TryU35LE() (uint64, error) { return d.tryUE(35, LittleEndian) }

func (d *D) ScalarU35LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(35, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U35LE() uint64 {
	v, err := d.tryUE(35, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U35LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU35LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU35LE, sfns...)
}

func (d *D) FieldU35LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U35LE, sfns...)
}

// Reader U36LE

func (d *D) TryU36LE() (uint64, error) { return d.tryUE(36, LittleEndian) }

func (d *D) ScalarU36LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(36, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U36LE() uint64 {
	v, err := d.tryUE(36, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U36LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU36LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU36LE, sfns...)
}

func (d *D) FieldU36LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U36LE, sfns...)
}

// Reader U37LE

func (d *D) TryU37LE() (uint64, error) { return d.tryUE(37, LittleEndian) }

func (d *D) ScalarU37LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(37, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U37LE() uint64 {
	v, err := d.tryUE(37, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U37LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU37LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU37LE, sfns...)
}

func (d *D) FieldU37LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U37LE, sfns...)
}

// Reader U38LE

func (d *D) TryU38LE() (uint64, error) { return d.tryUE(38, LittleEndian) }

func (d *D) ScalarU38LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(38, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U38LE() uint64 {
	v, err := d.tryUE(38, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U38LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU38LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU38LE, sfns...)
}

func (d *D) FieldU38LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U38LE, sfns...)
}

// Reader U39LE

func (d *D) TryU39LE() (uint64, error) { return d.tryUE(39, LittleEndian) }

func (d *D) ScalarU39LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(39, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U39LE() uint64 {
	v, err := d.tryUE(39, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U39LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU39LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU39LE, sfns...)
}

func (d *D) FieldU39LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U39LE, sfns...)
}

// Reader U40LE

func (d *D) TryU40LE() (uint64, error) { return d.tryUE(40, LittleEndian) }

func (d *D) ScalarU40LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(40, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U40LE() uint64 {
	v, err := d.tryUE(40, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U40LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU40LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU40LE, sfns...)
}

func (d *D) FieldU40LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U40LE, sfns...)
}

// Reader U41LE

func (d *D) TryU41LE() (uint64, error) { return d.tryUE(41, LittleEndian) }

func (d *D) ScalarU41LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(41, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U41LE() uint64 {
	v, err := d.tryUE(41, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U41LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU41LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU41LE, sfns...)
}

func (d *D) FieldU41LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U41LE, sfns...)
}

// Reader U42LE

func (d *D) TryU42LE() (uint64, error) { return d.tryUE(42, LittleEndian) }

func (d *D) ScalarU42LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(42, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U42LE() uint64 {
	v, err := d.tryUE(42, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U42LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU42LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU42LE, sfns...)
}

func (d *D) FieldU42LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U42LE, sfns...)
}

// Reader U43LE

func (d *D) TryU43LE() (uint64, error) { return d.tryUE(43, LittleEndian) }

func (d *D) ScalarU43LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(43, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U43LE() uint64 {
	v, err := d.tryUE(43, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U43LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU43LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU43LE, sfns...)
}

func (d *D) FieldU43LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U43LE, sfns...)
}

// Reader U44LE

func (d *D) TryU44LE() (uint64, error) { return d.tryUE(44, LittleEndian) }

func (d *D) ScalarU44LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(44, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U44LE() uint64 {
	v, err := d.tryUE(44, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U44LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU44LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU44LE, sfns...)
}

func (d *D) FieldU44LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U44LE, sfns...)
}

// Reader U45LE

func (d *D) TryU45LE() (uint64, error) { return d.tryUE(45, LittleEndian) }

func (d *D) ScalarU45LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(45, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U45LE() uint64 {
	v, err := d.tryUE(45, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U45LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU45LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU45LE, sfns...)
}

func (d *D) FieldU45LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U45LE, sfns...)
}

// Reader U46LE

func (d *D) TryU46LE() (uint64, error) { return d.tryUE(46, LittleEndian) }

func (d *D) ScalarU46LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(46, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U46LE() uint64 {
	v, err := d.tryUE(46, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U46LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU46LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU46LE, sfns...)
}

func (d *D) FieldU46LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U46LE, sfns...)
}

// Reader U47LE

func (d *D) TryU47LE() (uint64, error) { return d.tryUE(47, LittleEndian) }

func (d *D) ScalarU47LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(47, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U47LE() uint64 {
	v, err := d.tryUE(47, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U47LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU47LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU47LE, sfns...)
}

func (d *D) FieldU47LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U47LE, sfns...)
}

// Reader U48LE

func (d *D) TryU48LE() (uint64, error) { return d.tryUE(48, LittleEndian) }

func (d *D) ScalarU48LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(48, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U48LE() uint64 {
	v, err := d.tryUE(48, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U48LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU48LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU48LE, sfns...)
}

func (d *D) FieldU48LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U48LE, sfns...)
}

// Reader U49LE

func (d *D) TryU49LE() (uint64, error) { return d.tryUE(49, LittleEndian) }

func (d *D) ScalarU49LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(49, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U49LE() uint64 {
	v, err := d.tryUE(49, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U49LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU49LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU49LE, sfns...)
}

func (d *D) FieldU49LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U49LE, sfns...)
}

// Reader U50LE

func (d *D) TryU50LE() (uint64, error) { return d.tryUE(50, LittleEndian) }

func (d *D) ScalarU50LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(50, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U50LE() uint64 {
	v, err := d.tryUE(50, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U50LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU50LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU50LE, sfns...)
}

func (d *D) FieldU50LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U50LE, sfns...)
}

// Reader U51LE

func (d *D) TryU51LE() (uint64, error) { return d.tryUE(51, LittleEndian) }

func (d *D) ScalarU51LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(51, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U51LE() uint64 {
	v, err := d.tryUE(51, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U51LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU51LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU51LE, sfns...)
}

func (d *D) FieldU51LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U51LE, sfns...)
}

// Reader U52LE

func (d *D) TryU52LE() (uint64, error) { return d.tryUE(52, LittleEndian) }

func (d *D) ScalarU52LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(52, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U52LE() uint64 {
	v, err := d.tryUE(52, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U52LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU52LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU52LE, sfns...)
}

func (d *D) FieldU52LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U52LE, sfns...)
}

// Reader U53LE

func (d *D) TryU53LE() (uint64, error) { return d.tryUE(53, LittleEndian) }

func (d *D) ScalarU53LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(53, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U53LE() uint64 {
	v, err := d.tryUE(53, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U53LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU53LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU53LE, sfns...)
}

func (d *D) FieldU53LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U53LE, sfns...)
}

// Reader U54LE

func (d *D) TryU54LE() (uint64, error) { return d.tryUE(54, LittleEndian) }

func (d *D) ScalarU54LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(54, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U54LE() uint64 {
	v, err := d.tryUE(54, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U54LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU54LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU54LE, sfns...)
}

func (d *D) FieldU54LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U54LE, sfns...)
}

// Reader U55LE

func (d *D) TryU55LE() (uint64, error) { return d.tryUE(55, LittleEndian) }

func (d *D) ScalarU55LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(55, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U55LE() uint64 {
	v, err := d.tryUE(55, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U55LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU55LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU55LE, sfns...)
}

func (d *D) FieldU55LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U55LE, sfns...)
}

// Reader U56LE

func (d *D) TryU56LE() (uint64, error) { return d.tryUE(56, LittleEndian) }

func (d *D) ScalarU56LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(56, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U56LE() uint64 {
	v, err := d.tryUE(56, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U56LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU56LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU56LE, sfns...)
}

func (d *D) FieldU56LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U56LE, sfns...)
}

// Reader U57LE

func (d *D) TryU57LE() (uint64, error) { return d.tryUE(57, LittleEndian) }

func (d *D) ScalarU57LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(57, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U57LE() uint64 {
	v, err := d.tryUE(57, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U57LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU57LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU57LE, sfns...)
}

func (d *D) FieldU57LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U57LE, sfns...)
}

// Reader U58LE

func (d *D) TryU58LE() (uint64, error) { return d.tryUE(58, LittleEndian) }

func (d *D) ScalarU58LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(58, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U58LE() uint64 {
	v, err := d.tryUE(58, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U58LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU58LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU58LE, sfns...)
}

func (d *D) FieldU58LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U58LE, sfns...)
}

// Reader U59LE

func (d *D) TryU59LE() (uint64, error) { return d.tryUE(59, LittleEndian) }

func (d *D) ScalarU59LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(59, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U59LE() uint64 {
	v, err := d.tryUE(59, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U59LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU59LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU59LE, sfns...)
}

func (d *D) FieldU59LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U59LE, sfns...)
}

// Reader U60LE

func (d *D) TryU60LE() (uint64, error) { return d.tryUE(60, LittleEndian) }

func (d *D) ScalarU60LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(60, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U60LE() uint64 {
	v, err := d.tryUE(60, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U60LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU60LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU60LE, sfns...)
}

func (d *D) FieldU60LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U60LE, sfns...)
}

// Reader U61LE

func (d *D) TryU61LE() (uint64, error) { return d.tryUE(61, LittleEndian) }

func (d *D) ScalarU61LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(61, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U61LE() uint64 {
	v, err := d.tryUE(61, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U61LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU61LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU61LE, sfns...)
}

func (d *D) FieldU61LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U61LE, sfns...)
}

// Reader U62LE

func (d *D) TryU62LE() (uint64, error) { return d.tryUE(62, LittleEndian) }

func (d *D) ScalarU62LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(62, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U62LE() uint64 {
	v, err := d.tryUE(62, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U62LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU62LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU62LE, sfns...)
}

func (d *D) FieldU62LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U62LE, sfns...)
}

// Reader U63LE

func (d *D) TryU63LE() (uint64, error) { return d.tryUE(63, LittleEndian) }

func (d *D) ScalarU63LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(63, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U63LE() uint64 {
	v, err := d.tryUE(63, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U63LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU63LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU63LE, sfns...)
}

func (d *D) FieldU63LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U63LE, sfns...)
}

// Reader U64LE

func (d *D) TryU64LE() (uint64, error) { return d.tryUE(64, LittleEndian) }

func (d *D) ScalarU64LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(64, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U64LE() uint64 {
	v, err := d.tryUE(64, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U64LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU64LE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU64LE, sfns...)
}

func (d *D) FieldU64LE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U64LE, sfns...)
}

// Reader U8BE

func (d *D) TryU8BE() (uint64, error) { return d.tryUE(8, BigEndian) }

func (d *D) ScalarU8BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(8, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U8BE() uint64 {
	v, err := d.tryUE(8, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U8BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU8BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU8BE, sfns...)
}

func (d *D) FieldU8BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U8BE, sfns...)
}

// Reader U9BE

func (d *D) TryU9BE() (uint64, error) { return d.tryUE(9, BigEndian) }

func (d *D) ScalarU9BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(9, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U9BE() uint64 {
	v, err := d.tryUE(9, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U9BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU9BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU9BE, sfns...)
}

func (d *D) FieldU9BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U9BE, sfns...)
}

// Reader U10BE

func (d *D) TryU10BE() (uint64, error) { return d.tryUE(10, BigEndian) }

func (d *D) ScalarU10BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(10, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U10BE() uint64 {
	v, err := d.tryUE(10, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U10BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU10BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU10BE, sfns...)
}

func (d *D) FieldU10BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U10BE, sfns...)
}

// Reader U11BE

func (d *D) TryU11BE() (uint64, error) { return d.tryUE(11, BigEndian) }

func (d *D) ScalarU11BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(11, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U11BE() uint64 {
	v, err := d.tryUE(11, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U11BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU11BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU11BE, sfns...)
}

func (d *D) FieldU11BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U11BE, sfns...)
}

// Reader U12BE

func (d *D) TryU12BE() (uint64, error) { return d.tryUE(12, BigEndian) }

func (d *D) ScalarU12BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(12, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U12BE() uint64 {
	v, err := d.tryUE(12, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U12BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU12BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU12BE, sfns...)
}

func (d *D) FieldU12BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U12BE, sfns...)
}

// Reader U13BE

func (d *D) TryU13BE() (uint64, error) { return d.tryUE(13, BigEndian) }

func (d *D) ScalarU13BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(13, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U13BE() uint64 {
	v, err := d.tryUE(13, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U13BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU13BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU13BE, sfns...)
}

func (d *D) FieldU13BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U13BE, sfns...)
}

// Reader U14BE

func (d *D) TryU14BE() (uint64, error) { return d.tryUE(14, BigEndian) }

func (d *D) ScalarU14BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(14, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U14BE() uint64 {
	v, err := d.tryUE(14, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U14BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU14BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU14BE, sfns...)
}

func (d *D) FieldU14BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U14BE, sfns...)
}

// Reader U15BE

func (d *D) TryU15BE() (uint64, error) { return d.tryUE(15, BigEndian) }

func (d *D) ScalarU15BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(15, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U15BE() uint64 {
	v, err := d.tryUE(15, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U15BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU15BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU15BE, sfns...)
}

func (d *D) FieldU15BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U15BE, sfns...)
}

// Reader U16BE

func (d *D) TryU16BE() (uint64, error) { return d.tryUE(16, BigEndian) }

func (d *D) ScalarU16BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(16, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U16BE() uint64 {
	v, err := d.tryUE(16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U16BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU16BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU16BE, sfns...)
}

func (d *D) FieldU16BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U16BE, sfns...)
}

// Reader U17BE

func (d *D) TryU17BE() (uint64, error) { return d.tryUE(17, BigEndian) }

func (d *D) ScalarU17BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(17, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U17BE() uint64 {
	v, err := d.tryUE(17, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U17BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU17BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU17BE, sfns...)
}

func (d *D) FieldU17BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U17BE, sfns...)
}

// Reader U18BE

func (d *D) TryU18BE() (uint64, error) { return d.tryUE(18, BigEndian) }

func (d *D) ScalarU18BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(18, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U18BE() uint64 {
	v, err := d.tryUE(18, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U18BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU18BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU18BE, sfns...)
}

func (d *D) FieldU18BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U18BE, sfns...)
}

// Reader U19BE

func (d *D) TryU19BE() (uint64, error) { return d.tryUE(19, BigEndian) }

func (d *D) ScalarU19BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(19, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U19BE() uint64 {
	v, err := d.tryUE(19, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U19BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU19BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU19BE, sfns...)
}

func (d *D) FieldU19BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U19BE, sfns...)
}

// Reader U20BE

func (d *D) TryU20BE() (uint64, error) { return d.tryUE(20, BigEndian) }

func (d *D) ScalarU20BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(20, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U20BE() uint64 {
	v, err := d.tryUE(20, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U20BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU20BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU20BE, sfns...)
}

func (d *D) FieldU20BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U20BE, sfns...)
}

// Reader U21BE

func (d *D) TryU21BE() (uint64, error) { return d.tryUE(21, BigEndian) }

func (d *D) ScalarU21BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(21, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U21BE() uint64 {
	v, err := d.tryUE(21, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U21BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU21BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU21BE, sfns...)
}

func (d *D) FieldU21BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U21BE, sfns...)
}

// Reader U22BE

func (d *D) TryU22BE() (uint64, error) { return d.tryUE(22, BigEndian) }

func (d *D) ScalarU22BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(22, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U22BE() uint64 {
	v, err := d.tryUE(22, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U22BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU22BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU22BE, sfns...)
}

func (d *D) FieldU22BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U22BE, sfns...)
}

// Reader U23BE

func (d *D) TryU23BE() (uint64, error) { return d.tryUE(23, BigEndian) }

func (d *D) ScalarU23BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(23, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U23BE() uint64 {
	v, err := d.tryUE(23, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U23BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU23BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU23BE, sfns...)
}

func (d *D) FieldU23BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U23BE, sfns...)
}

// Reader U24BE

func (d *D) TryU24BE() (uint64, error) { return d.tryUE(24, BigEndian) }

func (d *D) ScalarU24BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(24, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U24BE() uint64 {
	v, err := d.tryUE(24, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U24BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU24BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU24BE, sfns...)
}

func (d *D) FieldU24BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U24BE, sfns...)
}

// Reader U25BE

func (d *D) TryU25BE() (uint64, error) { return d.tryUE(25, BigEndian) }

func (d *D) ScalarU25BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(25, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U25BE() uint64 {
	v, err := d.tryUE(25, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U25BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU25BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU25BE, sfns...)
}

func (d *D) FieldU25BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U25BE, sfns...)
}

// Reader U26BE

func (d *D) TryU26BE() (uint64, error) { return d.tryUE(26, BigEndian) }

func (d *D) ScalarU26BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(26, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U26BE() uint64 {
	v, err := d.tryUE(26, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U26BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU26BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU26BE, sfns...)
}

func (d *D) FieldU26BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U26BE, sfns...)
}

// Reader U27BE

func (d *D) TryU27BE() (uint64, error) { return d.tryUE(27, BigEndian) }

func (d *D) ScalarU27BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(27, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U27BE() uint64 {
	v, err := d.tryUE(27, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U27BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU27BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU27BE, sfns...)
}

func (d *D) FieldU27BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U27BE, sfns...)
}

// Reader U28BE

func (d *D) TryU28BE() (uint64, error) { return d.tryUE(28, BigEndian) }

func (d *D) ScalarU28BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(28, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U28BE() uint64 {
	v, err := d.tryUE(28, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U28BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU28BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU28BE, sfns...)
}

func (d *D) FieldU28BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U28BE, sfns...)
}

// Reader U29BE

func (d *D) TryU29BE() (uint64, error) { return d.tryUE(29, BigEndian) }

func (d *D) ScalarU29BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(29, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U29BE() uint64 {
	v, err := d.tryUE(29, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U29BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU29BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU29BE, sfns...)
}

func (d *D) FieldU29BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U29BE, sfns...)
}

// Reader U30BE

func (d *D) TryU30BE() (uint64, error) { return d.tryUE(30, BigEndian) }

func (d *D) ScalarU30BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(30, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U30BE() uint64 {
	v, err := d.tryUE(30, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U30BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU30BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU30BE, sfns...)
}

func (d *D) FieldU30BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U30BE, sfns...)
}

// Reader U31BE

func (d *D) TryU31BE() (uint64, error) { return d.tryUE(31, BigEndian) }

func (d *D) ScalarU31BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(31, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U31BE() uint64 {
	v, err := d.tryUE(31, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U31BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU31BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU31BE, sfns...)
}

func (d *D) FieldU31BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U31BE, sfns...)
}

// Reader U32BE

func (d *D) TryU32BE() (uint64, error) { return d.tryUE(32, BigEndian) }

func (d *D) ScalarU32BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(32, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U32BE() uint64 {
	v, err := d.tryUE(32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U32BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU32BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU32BE, sfns...)
}

func (d *D) FieldU32BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U32BE, sfns...)
}

// Reader U33BE

func (d *D) TryU33BE() (uint64, error) { return d.tryUE(33, BigEndian) }

func (d *D) ScalarU33BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(33, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U33BE() uint64 {
	v, err := d.tryUE(33, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U33BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU33BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU33BE, sfns...)
}

func (d *D) FieldU33BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U33BE, sfns...)
}

// Reader U34BE

func (d *D) TryU34BE() (uint64, error) { return d.tryUE(34, BigEndian) }

func (d *D) ScalarU34BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(34, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U34BE() uint64 {
	v, err := d.tryUE(34, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U34BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU34BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU34BE, sfns...)
}

func (d *D) FieldU34BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U34BE, sfns...)
}

// Reader U35BE

func (d *D) TryU35BE() (uint64, error) { return d.tryUE(35, BigEndian) }

func (d *D) ScalarU35BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(35, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U35BE() uint64 {
	v, err := d.tryUE(35, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U35BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU35BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU35BE, sfns...)
}

func (d *D) FieldU35BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U35BE, sfns...)
}

// Reader U36BE

func (d *D) TryU36BE() (uint64, error) { return d.tryUE(36, BigEndian) }

func (d *D) ScalarU36BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(36, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U36BE() uint64 {
	v, err := d.tryUE(36, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U36BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU36BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU36BE, sfns...)
}

func (d *D) FieldU36BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U36BE, sfns...)
}

// Reader U37BE

func (d *D) TryU37BE() (uint64, error) { return d.tryUE(37, BigEndian) }

func (d *D) ScalarU37BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(37, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U37BE() uint64 {
	v, err := d.tryUE(37, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U37BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU37BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU37BE, sfns...)
}

func (d *D) FieldU37BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U37BE, sfns...)
}

// Reader U38BE

func (d *D) TryU38BE() (uint64, error) { return d.tryUE(38, BigEndian) }

func (d *D) ScalarU38BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(38, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U38BE() uint64 {
	v, err := d.tryUE(38, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U38BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU38BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU38BE, sfns...)
}

func (d *D) FieldU38BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U38BE, sfns...)
}

// Reader U39BE

func (d *D) TryU39BE() (uint64, error) { return d.tryUE(39, BigEndian) }

func (d *D) ScalarU39BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(39, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U39BE() uint64 {
	v, err := d.tryUE(39, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U39BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU39BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU39BE, sfns...)
}

func (d *D) FieldU39BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U39BE, sfns...)
}

// Reader U40BE

func (d *D) TryU40BE() (uint64, error) { return d.tryUE(40, BigEndian) }

func (d *D) ScalarU40BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(40, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U40BE() uint64 {
	v, err := d.tryUE(40, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U40BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU40BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU40BE, sfns...)
}

func (d *D) FieldU40BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U40BE, sfns...)
}

// Reader U41BE

func (d *D) TryU41BE() (uint64, error) { return d.tryUE(41, BigEndian) }

func (d *D) ScalarU41BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(41, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U41BE() uint64 {
	v, err := d.tryUE(41, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U41BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU41BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU41BE, sfns...)
}

func (d *D) FieldU41BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U41BE, sfns...)
}

// Reader U42BE

func (d *D) TryU42BE() (uint64, error) { return d.tryUE(42, BigEndian) }

func (d *D) ScalarU42BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(42, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U42BE() uint64 {
	v, err := d.tryUE(42, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U42BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU42BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU42BE, sfns...)
}

func (d *D) FieldU42BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U42BE, sfns...)
}

// Reader U43BE

func (d *D) TryU43BE() (uint64, error) { return d.tryUE(43, BigEndian) }

func (d *D) ScalarU43BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(43, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U43BE() uint64 {
	v, err := d.tryUE(43, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U43BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU43BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU43BE, sfns...)
}

func (d *D) FieldU43BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U43BE, sfns...)
}

// Reader U44BE

func (d *D) TryU44BE() (uint64, error) { return d.tryUE(44, BigEndian) }

func (d *D) ScalarU44BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(44, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U44BE() uint64 {
	v, err := d.tryUE(44, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U44BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU44BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU44BE, sfns...)
}

func (d *D) FieldU44BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U44BE, sfns...)
}

// Reader U45BE

func (d *D) TryU45BE() (uint64, error) { return d.tryUE(45, BigEndian) }

func (d *D) ScalarU45BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(45, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U45BE() uint64 {
	v, err := d.tryUE(45, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U45BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU45BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU45BE, sfns...)
}

func (d *D) FieldU45BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U45BE, sfns...)
}

// Reader U46BE

func (d *D) TryU46BE() (uint64, error) { return d.tryUE(46, BigEndian) }

func (d *D) ScalarU46BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(46, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U46BE() uint64 {
	v, err := d.tryUE(46, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U46BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU46BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU46BE, sfns...)
}

func (d *D) FieldU46BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U46BE, sfns...)
}

// Reader U47BE

func (d *D) TryU47BE() (uint64, error) { return d.tryUE(47, BigEndian) }

func (d *D) ScalarU47BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(47, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U47BE() uint64 {
	v, err := d.tryUE(47, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U47BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU47BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU47BE, sfns...)
}

func (d *D) FieldU47BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U47BE, sfns...)
}

// Reader U48BE

func (d *D) TryU48BE() (uint64, error) { return d.tryUE(48, BigEndian) }

func (d *D) ScalarU48BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(48, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U48BE() uint64 {
	v, err := d.tryUE(48, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U48BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU48BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU48BE, sfns...)
}

func (d *D) FieldU48BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U48BE, sfns...)
}

// Reader U49BE

func (d *D) TryU49BE() (uint64, error) { return d.tryUE(49, BigEndian) }

func (d *D) ScalarU49BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(49, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U49BE() uint64 {
	v, err := d.tryUE(49, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U49BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU49BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU49BE, sfns...)
}

func (d *D) FieldU49BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U49BE, sfns...)
}

// Reader U50BE

func (d *D) TryU50BE() (uint64, error) { return d.tryUE(50, BigEndian) }

func (d *D) ScalarU50BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(50, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U50BE() uint64 {
	v, err := d.tryUE(50, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U50BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU50BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU50BE, sfns...)
}

func (d *D) FieldU50BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U50BE, sfns...)
}

// Reader U51BE

func (d *D) TryU51BE() (uint64, error) { return d.tryUE(51, BigEndian) }

func (d *D) ScalarU51BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(51, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U51BE() uint64 {
	v, err := d.tryUE(51, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U51BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU51BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU51BE, sfns...)
}

func (d *D) FieldU51BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U51BE, sfns...)
}

// Reader U52BE

func (d *D) TryU52BE() (uint64, error) { return d.tryUE(52, BigEndian) }

func (d *D) ScalarU52BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(52, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U52BE() uint64 {
	v, err := d.tryUE(52, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U52BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU52BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU52BE, sfns...)
}

func (d *D) FieldU52BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U52BE, sfns...)
}

// Reader U53BE

func (d *D) TryU53BE() (uint64, error) { return d.tryUE(53, BigEndian) }

func (d *D) ScalarU53BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(53, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U53BE() uint64 {
	v, err := d.tryUE(53, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U53BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU53BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU53BE, sfns...)
}

func (d *D) FieldU53BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U53BE, sfns...)
}

// Reader U54BE

func (d *D) TryU54BE() (uint64, error) { return d.tryUE(54, BigEndian) }

func (d *D) ScalarU54BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(54, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U54BE() uint64 {
	v, err := d.tryUE(54, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U54BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU54BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU54BE, sfns...)
}

func (d *D) FieldU54BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U54BE, sfns...)
}

// Reader U55BE

func (d *D) TryU55BE() (uint64, error) { return d.tryUE(55, BigEndian) }

func (d *D) ScalarU55BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(55, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U55BE() uint64 {
	v, err := d.tryUE(55, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U55BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU55BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU55BE, sfns...)
}

func (d *D) FieldU55BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U55BE, sfns...)
}

// Reader U56BE

func (d *D) TryU56BE() (uint64, error) { return d.tryUE(56, BigEndian) }

func (d *D) ScalarU56BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(56, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U56BE() uint64 {
	v, err := d.tryUE(56, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U56BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU56BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU56BE, sfns...)
}

func (d *D) FieldU56BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U56BE, sfns...)
}

// Reader U57BE

func (d *D) TryU57BE() (uint64, error) { return d.tryUE(57, BigEndian) }

func (d *D) ScalarU57BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(57, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U57BE() uint64 {
	v, err := d.tryUE(57, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U57BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU57BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU57BE, sfns...)
}

func (d *D) FieldU57BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U57BE, sfns...)
}

// Reader U58BE

func (d *D) TryU58BE() (uint64, error) { return d.tryUE(58, BigEndian) }

func (d *D) ScalarU58BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(58, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U58BE() uint64 {
	v, err := d.tryUE(58, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U58BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU58BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU58BE, sfns...)
}

func (d *D) FieldU58BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U58BE, sfns...)
}

// Reader U59BE

func (d *D) TryU59BE() (uint64, error) { return d.tryUE(59, BigEndian) }

func (d *D) ScalarU59BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(59, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U59BE() uint64 {
	v, err := d.tryUE(59, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U59BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU59BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU59BE, sfns...)
}

func (d *D) FieldU59BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U59BE, sfns...)
}

// Reader U60BE

func (d *D) TryU60BE() (uint64, error) { return d.tryUE(60, BigEndian) }

func (d *D) ScalarU60BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(60, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U60BE() uint64 {
	v, err := d.tryUE(60, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U60BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU60BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU60BE, sfns...)
}

func (d *D) FieldU60BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U60BE, sfns...)
}

// Reader U61BE

func (d *D) TryU61BE() (uint64, error) { return d.tryUE(61, BigEndian) }

func (d *D) ScalarU61BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(61, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U61BE() uint64 {
	v, err := d.tryUE(61, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U61BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU61BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU61BE, sfns...)
}

func (d *D) FieldU61BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U61BE, sfns...)
}

// Reader U62BE

func (d *D) TryU62BE() (uint64, error) { return d.tryUE(62, BigEndian) }

func (d *D) ScalarU62BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(62, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U62BE() uint64 {
	v, err := d.tryUE(62, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U62BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU62BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU62BE, sfns...)
}

func (d *D) FieldU62BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U62BE, sfns...)
}

// Reader U63BE

func (d *D) TryU63BE() (uint64, error) { return d.tryUE(63, BigEndian) }

func (d *D) ScalarU63BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(63, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U63BE() uint64 {
	v, err := d.tryUE(63, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U63BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU63BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU63BE, sfns...)
}

func (d *D) FieldU63BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U63BE, sfns...)
}

// Reader U64BE

func (d *D) TryU64BE() (uint64, error) { return d.tryUE(64, BigEndian) }

func (d *D) ScalarU64BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUE(64, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) U64BE() uint64 {
	v, err := d.tryUE(64, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "U64BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldU64BE(name string, sfns ...ScalarFn) (uint64, error) {
	return d.TryFieldUFn(name, (*D).TryU64BE, sfns...)
}

func (d *D) FieldU64BE(name string, sfns ...ScalarFn) uint64 {
	return d.FieldUFn(name, (*D).U64BE, sfns...)
}

// Reader SE

func (d *D) TrySE(nBits int, endian Endian) (int64, error) { return d.trySE(nBits, endian) }

func (d *D) ScalarSE(nBits int, endian Endian) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(nBits, endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) SE(nBits int, endian Endian) int64 {
	v, err := d.trySE(nBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "SE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldSE(name string, nBits int, endian Endian, sfns ...ScalarFn) (int64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarSE(nBits, endian), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualS(), err
}

func (d *D) FieldSE(name string, nBits int, endian Endian, sfns ...ScalarFn) int64 {
	v, err := d.TryFieldSE(name, nBits, endian, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "SE", Pos: d.Pos()})
	}
	return v
}

// Reader S

func (d *D) TryS(nBits int) (int64, error) { return d.trySE(nBits, d.Endian) }

func (d *D) ScalarS(nBits int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(nBits, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S(nBits int) int64 {
	v, err := d.trySE(nBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS(name string, nBits int, sfns ...ScalarFn) (int64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarS(nBits), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualS(), err
}

func (d *D) FieldS(name string, nBits int, sfns ...ScalarFn) int64 {
	v, err := d.TryFieldS(name, nBits, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "S", Pos: d.Pos()})
	}
	return v
}

// Reader S1

func (d *D) TryS1() (int64, error) { return d.trySE(1, d.Endian) }

func (d *D) ScalarS1() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(1, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S1() int64 {
	v, err := d.trySE(1, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S1", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS1(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS1, sfns...)
}

func (d *D) FieldS1(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S1, sfns...)
}

// Reader S2

func (d *D) TryS2() (int64, error) { return d.trySE(2, d.Endian) }

func (d *D) ScalarS2() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(2, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S2() int64 {
	v, err := d.trySE(2, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S2", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS2(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS2, sfns...)
}

func (d *D) FieldS2(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S2, sfns...)
}

// Reader S3

func (d *D) TryS3() (int64, error) { return d.trySE(3, d.Endian) }

func (d *D) ScalarS3() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(3, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S3() int64 {
	v, err := d.trySE(3, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S3", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS3(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS3, sfns...)
}

func (d *D) FieldS3(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S3, sfns...)
}

// Reader S4

func (d *D) TryS4() (int64, error) { return d.trySE(4, d.Endian) }

func (d *D) ScalarS4() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(4, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S4() int64 {
	v, err := d.trySE(4, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S4", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS4(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS4, sfns...)
}

func (d *D) FieldS4(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S4, sfns...)
}

// Reader S5

func (d *D) TryS5() (int64, error) { return d.trySE(5, d.Endian) }

func (d *D) ScalarS5() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(5, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S5() int64 {
	v, err := d.trySE(5, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S5", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS5(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS5, sfns...)
}

func (d *D) FieldS5(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S5, sfns...)
}

// Reader S6

func (d *D) TryS6() (int64, error) { return d.trySE(6, d.Endian) }

func (d *D) ScalarS6() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(6, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S6() int64 {
	v, err := d.trySE(6, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S6", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS6(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS6, sfns...)
}

func (d *D) FieldS6(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S6, sfns...)
}

// Reader S7

func (d *D) TryS7() (int64, error) { return d.trySE(7, d.Endian) }

func (d *D) ScalarS7() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(7, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S7() int64 {
	v, err := d.trySE(7, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S7", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS7(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS7, sfns...)
}

func (d *D) FieldS7(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S7, sfns...)
}

// Reader S8

func (d *D) TryS8() (int64, error) { return d.trySE(8, d.Endian) }

func (d *D) ScalarS8() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(8, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S8() int64 {
	v, err := d.trySE(8, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S8", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS8(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS8, sfns...)
}

func (d *D) FieldS8(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S8, sfns...)
}

// Reader S9

func (d *D) TryS9() (int64, error) { return d.trySE(9, d.Endian) }

func (d *D) ScalarS9() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(9, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S9() int64 {
	v, err := d.trySE(9, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S9", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS9(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS9, sfns...)
}

func (d *D) FieldS9(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S9, sfns...)
}

// Reader S10

func (d *D) TryS10() (int64, error) { return d.trySE(10, d.Endian) }

func (d *D) ScalarS10() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(10, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S10() int64 {
	v, err := d.trySE(10, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S10", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS10(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS10, sfns...)
}

func (d *D) FieldS10(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S10, sfns...)
}

// Reader S11

func (d *D) TryS11() (int64, error) { return d.trySE(11, d.Endian) }

func (d *D) ScalarS11() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(11, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S11() int64 {
	v, err := d.trySE(11, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S11", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS11(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS11, sfns...)
}

func (d *D) FieldS11(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S11, sfns...)
}

// Reader S12

func (d *D) TryS12() (int64, error) { return d.trySE(12, d.Endian) }

func (d *D) ScalarS12() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(12, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S12() int64 {
	v, err := d.trySE(12, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S12", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS12(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS12, sfns...)
}

func (d *D) FieldS12(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S12, sfns...)
}

// Reader S13

func (d *D) TryS13() (int64, error) { return d.trySE(13, d.Endian) }

func (d *D) ScalarS13() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(13, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S13() int64 {
	v, err := d.trySE(13, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S13", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS13(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS13, sfns...)
}

func (d *D) FieldS13(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S13, sfns...)
}

// Reader S14

func (d *D) TryS14() (int64, error) { return d.trySE(14, d.Endian) }

func (d *D) ScalarS14() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(14, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S14() int64 {
	v, err := d.trySE(14, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S14", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS14(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS14, sfns...)
}

func (d *D) FieldS14(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S14, sfns...)
}

// Reader S15

func (d *D) TryS15() (int64, error) { return d.trySE(15, d.Endian) }

func (d *D) ScalarS15() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(15, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S15() int64 {
	v, err := d.trySE(15, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S15", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS15(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS15, sfns...)
}

func (d *D) FieldS15(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S15, sfns...)
}

// Reader S16

func (d *D) TryS16() (int64, error) { return d.trySE(16, d.Endian) }

func (d *D) ScalarS16() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(16, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S16() int64 {
	v, err := d.trySE(16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S16", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS16(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS16, sfns...)
}

func (d *D) FieldS16(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S16, sfns...)
}

// Reader S17

func (d *D) TryS17() (int64, error) { return d.trySE(17, d.Endian) }

func (d *D) ScalarS17() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(17, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S17() int64 {
	v, err := d.trySE(17, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S17", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS17(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS17, sfns...)
}

func (d *D) FieldS17(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S17, sfns...)
}

// Reader S18

func (d *D) TryS18() (int64, error) { return d.trySE(18, d.Endian) }

func (d *D) ScalarS18() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(18, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S18() int64 {
	v, err := d.trySE(18, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S18", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS18(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS18, sfns...)
}

func (d *D) FieldS18(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S18, sfns...)
}

// Reader S19

func (d *D) TryS19() (int64, error) { return d.trySE(19, d.Endian) }

func (d *D) ScalarS19() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(19, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S19() int64 {
	v, err := d.trySE(19, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S19", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS19(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS19, sfns...)
}

func (d *D) FieldS19(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S19, sfns...)
}

// Reader S20

func (d *D) TryS20() (int64, error) { return d.trySE(20, d.Endian) }

func (d *D) ScalarS20() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(20, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S20() int64 {
	v, err := d.trySE(20, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S20", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS20(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS20, sfns...)
}

func (d *D) FieldS20(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S20, sfns...)
}

// Reader S21

func (d *D) TryS21() (int64, error) { return d.trySE(21, d.Endian) }

func (d *D) ScalarS21() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(21, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S21() int64 {
	v, err := d.trySE(21, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S21", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS21(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS21, sfns...)
}

func (d *D) FieldS21(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S21, sfns...)
}

// Reader S22

func (d *D) TryS22() (int64, error) { return d.trySE(22, d.Endian) }

func (d *D) ScalarS22() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(22, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S22() int64 {
	v, err := d.trySE(22, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S22", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS22(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS22, sfns...)
}

func (d *D) FieldS22(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S22, sfns...)
}

// Reader S23

func (d *D) TryS23() (int64, error) { return d.trySE(23, d.Endian) }

func (d *D) ScalarS23() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(23, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S23() int64 {
	v, err := d.trySE(23, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S23", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS23(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS23, sfns...)
}

func (d *D) FieldS23(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S23, sfns...)
}

// Reader S24

func (d *D) TryS24() (int64, error) { return d.trySE(24, d.Endian) }

func (d *D) ScalarS24() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(24, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S24() int64 {
	v, err := d.trySE(24, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S24", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS24(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS24, sfns...)
}

func (d *D) FieldS24(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S24, sfns...)
}

// Reader S25

func (d *D) TryS25() (int64, error) { return d.trySE(25, d.Endian) }

func (d *D) ScalarS25() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(25, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S25() int64 {
	v, err := d.trySE(25, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S25", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS25(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS25, sfns...)
}

func (d *D) FieldS25(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S25, sfns...)
}

// Reader S26

func (d *D) TryS26() (int64, error) { return d.trySE(26, d.Endian) }

func (d *D) ScalarS26() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(26, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S26() int64 {
	v, err := d.trySE(26, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S26", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS26(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS26, sfns...)
}

func (d *D) FieldS26(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S26, sfns...)
}

// Reader S27

func (d *D) TryS27() (int64, error) { return d.trySE(27, d.Endian) }

func (d *D) ScalarS27() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(27, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S27() int64 {
	v, err := d.trySE(27, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S27", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS27(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS27, sfns...)
}

func (d *D) FieldS27(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S27, sfns...)
}

// Reader S28

func (d *D) TryS28() (int64, error) { return d.trySE(28, d.Endian) }

func (d *D) ScalarS28() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(28, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S28() int64 {
	v, err := d.trySE(28, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S28", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS28(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS28, sfns...)
}

func (d *D) FieldS28(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S28, sfns...)
}

// Reader S29

func (d *D) TryS29() (int64, error) { return d.trySE(29, d.Endian) }

func (d *D) ScalarS29() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(29, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S29() int64 {
	v, err := d.trySE(29, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S29", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS29(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS29, sfns...)
}

func (d *D) FieldS29(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S29, sfns...)
}

// Reader S30

func (d *D) TryS30() (int64, error) { return d.trySE(30, d.Endian) }

func (d *D) ScalarS30() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(30, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S30() int64 {
	v, err := d.trySE(30, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S30", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS30(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS30, sfns...)
}

func (d *D) FieldS30(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S30, sfns...)
}

// Reader S31

func (d *D) TryS31() (int64, error) { return d.trySE(31, d.Endian) }

func (d *D) ScalarS31() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(31, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S31() int64 {
	v, err := d.trySE(31, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S31", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS31(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS31, sfns...)
}

func (d *D) FieldS31(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S31, sfns...)
}

// Reader S32

func (d *D) TryS32() (int64, error) { return d.trySE(32, d.Endian) }

func (d *D) ScalarS32() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(32, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S32() int64 {
	v, err := d.trySE(32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S32", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS32(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS32, sfns...)
}

func (d *D) FieldS32(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S32, sfns...)
}

// Reader S33

func (d *D) TryS33() (int64, error) { return d.trySE(33, d.Endian) }

func (d *D) ScalarS33() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(33, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S33() int64 {
	v, err := d.trySE(33, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S33", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS33(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS33, sfns...)
}

func (d *D) FieldS33(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S33, sfns...)
}

// Reader S34

func (d *D) TryS34() (int64, error) { return d.trySE(34, d.Endian) }

func (d *D) ScalarS34() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(34, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S34() int64 {
	v, err := d.trySE(34, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S34", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS34(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS34, sfns...)
}

func (d *D) FieldS34(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S34, sfns...)
}

// Reader S35

func (d *D) TryS35() (int64, error) { return d.trySE(35, d.Endian) }

func (d *D) ScalarS35() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(35, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S35() int64 {
	v, err := d.trySE(35, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S35", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS35(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS35, sfns...)
}

func (d *D) FieldS35(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S35, sfns...)
}

// Reader S36

func (d *D) TryS36() (int64, error) { return d.trySE(36, d.Endian) }

func (d *D) ScalarS36() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(36, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S36() int64 {
	v, err := d.trySE(36, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S36", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS36(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS36, sfns...)
}

func (d *D) FieldS36(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S36, sfns...)
}

// Reader S37

func (d *D) TryS37() (int64, error) { return d.trySE(37, d.Endian) }

func (d *D) ScalarS37() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(37, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S37() int64 {
	v, err := d.trySE(37, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S37", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS37(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS37, sfns...)
}

func (d *D) FieldS37(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S37, sfns...)
}

// Reader S38

func (d *D) TryS38() (int64, error) { return d.trySE(38, d.Endian) }

func (d *D) ScalarS38() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(38, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S38() int64 {
	v, err := d.trySE(38, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S38", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS38(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS38, sfns...)
}

func (d *D) FieldS38(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S38, sfns...)
}

// Reader S39

func (d *D) TryS39() (int64, error) { return d.trySE(39, d.Endian) }

func (d *D) ScalarS39() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(39, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S39() int64 {
	v, err := d.trySE(39, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S39", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS39(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS39, sfns...)
}

func (d *D) FieldS39(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S39, sfns...)
}

// Reader S40

func (d *D) TryS40() (int64, error) { return d.trySE(40, d.Endian) }

func (d *D) ScalarS40() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(40, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S40() int64 {
	v, err := d.trySE(40, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S40", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS40(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS40, sfns...)
}

func (d *D) FieldS40(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S40, sfns...)
}

// Reader S41

func (d *D) TryS41() (int64, error) { return d.trySE(41, d.Endian) }

func (d *D) ScalarS41() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(41, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S41() int64 {
	v, err := d.trySE(41, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S41", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS41(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS41, sfns...)
}

func (d *D) FieldS41(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S41, sfns...)
}

// Reader S42

func (d *D) TryS42() (int64, error) { return d.trySE(42, d.Endian) }

func (d *D) ScalarS42() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(42, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S42() int64 {
	v, err := d.trySE(42, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S42", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS42(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS42, sfns...)
}

func (d *D) FieldS42(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S42, sfns...)
}

// Reader S43

func (d *D) TryS43() (int64, error) { return d.trySE(43, d.Endian) }

func (d *D) ScalarS43() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(43, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S43() int64 {
	v, err := d.trySE(43, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S43", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS43(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS43, sfns...)
}

func (d *D) FieldS43(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S43, sfns...)
}

// Reader S44

func (d *D) TryS44() (int64, error) { return d.trySE(44, d.Endian) }

func (d *D) ScalarS44() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(44, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S44() int64 {
	v, err := d.trySE(44, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S44", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS44(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS44, sfns...)
}

func (d *D) FieldS44(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S44, sfns...)
}

// Reader S45

func (d *D) TryS45() (int64, error) { return d.trySE(45, d.Endian) }

func (d *D) ScalarS45() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(45, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S45() int64 {
	v, err := d.trySE(45, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S45", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS45(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS45, sfns...)
}

func (d *D) FieldS45(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S45, sfns...)
}

// Reader S46

func (d *D) TryS46() (int64, error) { return d.trySE(46, d.Endian) }

func (d *D) ScalarS46() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(46, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S46() int64 {
	v, err := d.trySE(46, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S46", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS46(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS46, sfns...)
}

func (d *D) FieldS46(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S46, sfns...)
}

// Reader S47

func (d *D) TryS47() (int64, error) { return d.trySE(47, d.Endian) }

func (d *D) ScalarS47() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(47, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S47() int64 {
	v, err := d.trySE(47, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S47", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS47(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS47, sfns...)
}

func (d *D) FieldS47(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S47, sfns...)
}

// Reader S48

func (d *D) TryS48() (int64, error) { return d.trySE(48, d.Endian) }

func (d *D) ScalarS48() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(48, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S48() int64 {
	v, err := d.trySE(48, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S48", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS48(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS48, sfns...)
}

func (d *D) FieldS48(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S48, sfns...)
}

// Reader S49

func (d *D) TryS49() (int64, error) { return d.trySE(49, d.Endian) }

func (d *D) ScalarS49() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(49, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S49() int64 {
	v, err := d.trySE(49, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S49", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS49(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS49, sfns...)
}

func (d *D) FieldS49(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S49, sfns...)
}

// Reader S50

func (d *D) TryS50() (int64, error) { return d.trySE(50, d.Endian) }

func (d *D) ScalarS50() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(50, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S50() int64 {
	v, err := d.trySE(50, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S50", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS50(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS50, sfns...)
}

func (d *D) FieldS50(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S50, sfns...)
}

// Reader S51

func (d *D) TryS51() (int64, error) { return d.trySE(51, d.Endian) }

func (d *D) ScalarS51() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(51, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S51() int64 {
	v, err := d.trySE(51, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S51", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS51(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS51, sfns...)
}

func (d *D) FieldS51(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S51, sfns...)
}

// Reader S52

func (d *D) TryS52() (int64, error) { return d.trySE(52, d.Endian) }

func (d *D) ScalarS52() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(52, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S52() int64 {
	v, err := d.trySE(52, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S52", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS52(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS52, sfns...)
}

func (d *D) FieldS52(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S52, sfns...)
}

// Reader S53

func (d *D) TryS53() (int64, error) { return d.trySE(53, d.Endian) }

func (d *D) ScalarS53() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(53, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S53() int64 {
	v, err := d.trySE(53, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S53", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS53(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS53, sfns...)
}

func (d *D) FieldS53(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S53, sfns...)
}

// Reader S54

func (d *D) TryS54() (int64, error) { return d.trySE(54, d.Endian) }

func (d *D) ScalarS54() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(54, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S54() int64 {
	v, err := d.trySE(54, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S54", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS54(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS54, sfns...)
}

func (d *D) FieldS54(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S54, sfns...)
}

// Reader S55

func (d *D) TryS55() (int64, error) { return d.trySE(55, d.Endian) }

func (d *D) ScalarS55() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(55, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S55() int64 {
	v, err := d.trySE(55, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S55", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS55(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS55, sfns...)
}

func (d *D) FieldS55(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S55, sfns...)
}

// Reader S56

func (d *D) TryS56() (int64, error) { return d.trySE(56, d.Endian) }

func (d *D) ScalarS56() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(56, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S56() int64 {
	v, err := d.trySE(56, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S56", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS56(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS56, sfns...)
}

func (d *D) FieldS56(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S56, sfns...)
}

// Reader S57

func (d *D) TryS57() (int64, error) { return d.trySE(57, d.Endian) }

func (d *D) ScalarS57() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(57, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S57() int64 {
	v, err := d.trySE(57, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S57", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS57(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS57, sfns...)
}

func (d *D) FieldS57(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S57, sfns...)
}

// Reader S58

func (d *D) TryS58() (int64, error) { return d.trySE(58, d.Endian) }

func (d *D) ScalarS58() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(58, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S58() int64 {
	v, err := d.trySE(58, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S58", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS58(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS58, sfns...)
}

func (d *D) FieldS58(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S58, sfns...)
}

// Reader S59

func (d *D) TryS59() (int64, error) { return d.trySE(59, d.Endian) }

func (d *D) ScalarS59() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(59, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S59() int64 {
	v, err := d.trySE(59, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S59", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS59(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS59, sfns...)
}

func (d *D) FieldS59(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S59, sfns...)
}

// Reader S60

func (d *D) TryS60() (int64, error) { return d.trySE(60, d.Endian) }

func (d *D) ScalarS60() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(60, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S60() int64 {
	v, err := d.trySE(60, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S60", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS60(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS60, sfns...)
}

func (d *D) FieldS60(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S60, sfns...)
}

// Reader S61

func (d *D) TryS61() (int64, error) { return d.trySE(61, d.Endian) }

func (d *D) ScalarS61() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(61, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S61() int64 {
	v, err := d.trySE(61, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S61", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS61(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS61, sfns...)
}

func (d *D) FieldS61(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S61, sfns...)
}

// Reader S62

func (d *D) TryS62() (int64, error) { return d.trySE(62, d.Endian) }

func (d *D) ScalarS62() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(62, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S62() int64 {
	v, err := d.trySE(62, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S62", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS62(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS62, sfns...)
}

func (d *D) FieldS62(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S62, sfns...)
}

// Reader S63

func (d *D) TryS63() (int64, error) { return d.trySE(63, d.Endian) }

func (d *D) ScalarS63() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(63, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S63() int64 {
	v, err := d.trySE(63, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S63", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS63(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS63, sfns...)
}

func (d *D) FieldS63(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S63, sfns...)
}

// Reader S64

func (d *D) TryS64() (int64, error) { return d.trySE(64, d.Endian) }

func (d *D) ScalarS64() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(64, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S64() int64 {
	v, err := d.trySE(64, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "S64", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS64(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS64, sfns...)
}

func (d *D) FieldS64(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S64, sfns...)
}

// Reader S8LE

func (d *D) TryS8LE() (int64, error) { return d.trySE(8, LittleEndian) }

func (d *D) ScalarS8LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(8, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S8LE() int64 {
	v, err := d.trySE(8, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S8LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS8LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS8LE, sfns...)
}

func (d *D) FieldS8LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S8LE, sfns...)
}

// Reader S9LE

func (d *D) TryS9LE() (int64, error) { return d.trySE(9, LittleEndian) }

func (d *D) ScalarS9LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(9, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S9LE() int64 {
	v, err := d.trySE(9, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S9LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS9LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS9LE, sfns...)
}

func (d *D) FieldS9LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S9LE, sfns...)
}

// Reader S10LE

func (d *D) TryS10LE() (int64, error) { return d.trySE(10, LittleEndian) }

func (d *D) ScalarS10LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(10, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S10LE() int64 {
	v, err := d.trySE(10, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S10LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS10LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS10LE, sfns...)
}

func (d *D) FieldS10LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S10LE, sfns...)
}

// Reader S11LE

func (d *D) TryS11LE() (int64, error) { return d.trySE(11, LittleEndian) }

func (d *D) ScalarS11LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(11, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S11LE() int64 {
	v, err := d.trySE(11, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S11LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS11LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS11LE, sfns...)
}

func (d *D) FieldS11LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S11LE, sfns...)
}

// Reader S12LE

func (d *D) TryS12LE() (int64, error) { return d.trySE(12, LittleEndian) }

func (d *D) ScalarS12LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(12, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S12LE() int64 {
	v, err := d.trySE(12, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S12LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS12LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS12LE, sfns...)
}

func (d *D) FieldS12LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S12LE, sfns...)
}

// Reader S13LE

func (d *D) TryS13LE() (int64, error) { return d.trySE(13, LittleEndian) }

func (d *D) ScalarS13LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(13, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S13LE() int64 {
	v, err := d.trySE(13, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S13LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS13LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS13LE, sfns...)
}

func (d *D) FieldS13LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S13LE, sfns...)
}

// Reader S14LE

func (d *D) TryS14LE() (int64, error) { return d.trySE(14, LittleEndian) }

func (d *D) ScalarS14LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(14, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S14LE() int64 {
	v, err := d.trySE(14, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S14LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS14LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS14LE, sfns...)
}

func (d *D) FieldS14LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S14LE, sfns...)
}

// Reader S15LE

func (d *D) TryS15LE() (int64, error) { return d.trySE(15, LittleEndian) }

func (d *D) ScalarS15LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(15, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S15LE() int64 {
	v, err := d.trySE(15, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S15LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS15LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS15LE, sfns...)
}

func (d *D) FieldS15LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S15LE, sfns...)
}

// Reader S16LE

func (d *D) TryS16LE() (int64, error) { return d.trySE(16, LittleEndian) }

func (d *D) ScalarS16LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(16, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S16LE() int64 {
	v, err := d.trySE(16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S16LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS16LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS16LE, sfns...)
}

func (d *D) FieldS16LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S16LE, sfns...)
}

// Reader S17LE

func (d *D) TryS17LE() (int64, error) { return d.trySE(17, LittleEndian) }

func (d *D) ScalarS17LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(17, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S17LE() int64 {
	v, err := d.trySE(17, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S17LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS17LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS17LE, sfns...)
}

func (d *D) FieldS17LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S17LE, sfns...)
}

// Reader S18LE

func (d *D) TryS18LE() (int64, error) { return d.trySE(18, LittleEndian) }

func (d *D) ScalarS18LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(18, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S18LE() int64 {
	v, err := d.trySE(18, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S18LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS18LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS18LE, sfns...)
}

func (d *D) FieldS18LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S18LE, sfns...)
}

// Reader S19LE

func (d *D) TryS19LE() (int64, error) { return d.trySE(19, LittleEndian) }

func (d *D) ScalarS19LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(19, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S19LE() int64 {
	v, err := d.trySE(19, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S19LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS19LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS19LE, sfns...)
}

func (d *D) FieldS19LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S19LE, sfns...)
}

// Reader S20LE

func (d *D) TryS20LE() (int64, error) { return d.trySE(20, LittleEndian) }

func (d *D) ScalarS20LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(20, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S20LE() int64 {
	v, err := d.trySE(20, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S20LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS20LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS20LE, sfns...)
}

func (d *D) FieldS20LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S20LE, sfns...)
}

// Reader S21LE

func (d *D) TryS21LE() (int64, error) { return d.trySE(21, LittleEndian) }

func (d *D) ScalarS21LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(21, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S21LE() int64 {
	v, err := d.trySE(21, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S21LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS21LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS21LE, sfns...)
}

func (d *D) FieldS21LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S21LE, sfns...)
}

// Reader S22LE

func (d *D) TryS22LE() (int64, error) { return d.trySE(22, LittleEndian) }

func (d *D) ScalarS22LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(22, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S22LE() int64 {
	v, err := d.trySE(22, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S22LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS22LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS22LE, sfns...)
}

func (d *D) FieldS22LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S22LE, sfns...)
}

// Reader S23LE

func (d *D) TryS23LE() (int64, error) { return d.trySE(23, LittleEndian) }

func (d *D) ScalarS23LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(23, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S23LE() int64 {
	v, err := d.trySE(23, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S23LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS23LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS23LE, sfns...)
}

func (d *D) FieldS23LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S23LE, sfns...)
}

// Reader S24LE

func (d *D) TryS24LE() (int64, error) { return d.trySE(24, LittleEndian) }

func (d *D) ScalarS24LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(24, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S24LE() int64 {
	v, err := d.trySE(24, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S24LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS24LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS24LE, sfns...)
}

func (d *D) FieldS24LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S24LE, sfns...)
}

// Reader S25LE

func (d *D) TryS25LE() (int64, error) { return d.trySE(25, LittleEndian) }

func (d *D) ScalarS25LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(25, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S25LE() int64 {
	v, err := d.trySE(25, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S25LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS25LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS25LE, sfns...)
}

func (d *D) FieldS25LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S25LE, sfns...)
}

// Reader S26LE

func (d *D) TryS26LE() (int64, error) { return d.trySE(26, LittleEndian) }

func (d *D) ScalarS26LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(26, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S26LE() int64 {
	v, err := d.trySE(26, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S26LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS26LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS26LE, sfns...)
}

func (d *D) FieldS26LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S26LE, sfns...)
}

// Reader S27LE

func (d *D) TryS27LE() (int64, error) { return d.trySE(27, LittleEndian) }

func (d *D) ScalarS27LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(27, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S27LE() int64 {
	v, err := d.trySE(27, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S27LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS27LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS27LE, sfns...)
}

func (d *D) FieldS27LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S27LE, sfns...)
}

// Reader S28LE

func (d *D) TryS28LE() (int64, error) { return d.trySE(28, LittleEndian) }

func (d *D) ScalarS28LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(28, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S28LE() int64 {
	v, err := d.trySE(28, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S28LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS28LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS28LE, sfns...)
}

func (d *D) FieldS28LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S28LE, sfns...)
}

// Reader S29LE

func (d *D) TryS29LE() (int64, error) { return d.trySE(29, LittleEndian) }

func (d *D) ScalarS29LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(29, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S29LE() int64 {
	v, err := d.trySE(29, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S29LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS29LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS29LE, sfns...)
}

func (d *D) FieldS29LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S29LE, sfns...)
}

// Reader S30LE

func (d *D) TryS30LE() (int64, error) { return d.trySE(30, LittleEndian) }

func (d *D) ScalarS30LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(30, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S30LE() int64 {
	v, err := d.trySE(30, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S30LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS30LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS30LE, sfns...)
}

func (d *D) FieldS30LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S30LE, sfns...)
}

// Reader S31LE

func (d *D) TryS31LE() (int64, error) { return d.trySE(31, LittleEndian) }

func (d *D) ScalarS31LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(31, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S31LE() int64 {
	v, err := d.trySE(31, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S31LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS31LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS31LE, sfns...)
}

func (d *D) FieldS31LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S31LE, sfns...)
}

// Reader S32LE

func (d *D) TryS32LE() (int64, error) { return d.trySE(32, LittleEndian) }

func (d *D) ScalarS32LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(32, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S32LE() int64 {
	v, err := d.trySE(32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S32LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS32LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS32LE, sfns...)
}

func (d *D) FieldS32LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S32LE, sfns...)
}

// Reader S33LE

func (d *D) TryS33LE() (int64, error) { return d.trySE(33, LittleEndian) }

func (d *D) ScalarS33LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(33, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S33LE() int64 {
	v, err := d.trySE(33, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S33LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS33LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS33LE, sfns...)
}

func (d *D) FieldS33LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S33LE, sfns...)
}

// Reader S34LE

func (d *D) TryS34LE() (int64, error) { return d.trySE(34, LittleEndian) }

func (d *D) ScalarS34LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(34, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S34LE() int64 {
	v, err := d.trySE(34, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S34LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS34LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS34LE, sfns...)
}

func (d *D) FieldS34LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S34LE, sfns...)
}

// Reader S35LE

func (d *D) TryS35LE() (int64, error) { return d.trySE(35, LittleEndian) }

func (d *D) ScalarS35LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(35, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S35LE() int64 {
	v, err := d.trySE(35, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S35LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS35LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS35LE, sfns...)
}

func (d *D) FieldS35LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S35LE, sfns...)
}

// Reader S36LE

func (d *D) TryS36LE() (int64, error) { return d.trySE(36, LittleEndian) }

func (d *D) ScalarS36LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(36, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S36LE() int64 {
	v, err := d.trySE(36, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S36LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS36LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS36LE, sfns...)
}

func (d *D) FieldS36LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S36LE, sfns...)
}

// Reader S37LE

func (d *D) TryS37LE() (int64, error) { return d.trySE(37, LittleEndian) }

func (d *D) ScalarS37LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(37, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S37LE() int64 {
	v, err := d.trySE(37, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S37LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS37LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS37LE, sfns...)
}

func (d *D) FieldS37LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S37LE, sfns...)
}

// Reader S38LE

func (d *D) TryS38LE() (int64, error) { return d.trySE(38, LittleEndian) }

func (d *D) ScalarS38LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(38, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S38LE() int64 {
	v, err := d.trySE(38, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S38LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS38LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS38LE, sfns...)
}

func (d *D) FieldS38LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S38LE, sfns...)
}

// Reader S39LE

func (d *D) TryS39LE() (int64, error) { return d.trySE(39, LittleEndian) }

func (d *D) ScalarS39LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(39, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S39LE() int64 {
	v, err := d.trySE(39, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S39LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS39LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS39LE, sfns...)
}

func (d *D) FieldS39LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S39LE, sfns...)
}

// Reader S40LE

func (d *D) TryS40LE() (int64, error) { return d.trySE(40, LittleEndian) }

func (d *D) ScalarS40LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(40, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S40LE() int64 {
	v, err := d.trySE(40, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S40LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS40LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS40LE, sfns...)
}

func (d *D) FieldS40LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S40LE, sfns...)
}

// Reader S41LE

func (d *D) TryS41LE() (int64, error) { return d.trySE(41, LittleEndian) }

func (d *D) ScalarS41LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(41, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S41LE() int64 {
	v, err := d.trySE(41, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S41LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS41LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS41LE, sfns...)
}

func (d *D) FieldS41LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S41LE, sfns...)
}

// Reader S42LE

func (d *D) TryS42LE() (int64, error) { return d.trySE(42, LittleEndian) }

func (d *D) ScalarS42LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(42, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S42LE() int64 {
	v, err := d.trySE(42, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S42LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS42LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS42LE, sfns...)
}

func (d *D) FieldS42LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S42LE, sfns...)
}

// Reader S43LE

func (d *D) TryS43LE() (int64, error) { return d.trySE(43, LittleEndian) }

func (d *D) ScalarS43LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(43, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S43LE() int64 {
	v, err := d.trySE(43, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S43LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS43LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS43LE, sfns...)
}

func (d *D) FieldS43LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S43LE, sfns...)
}

// Reader S44LE

func (d *D) TryS44LE() (int64, error) { return d.trySE(44, LittleEndian) }

func (d *D) ScalarS44LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(44, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S44LE() int64 {
	v, err := d.trySE(44, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S44LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS44LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS44LE, sfns...)
}

func (d *D) FieldS44LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S44LE, sfns...)
}

// Reader S45LE

func (d *D) TryS45LE() (int64, error) { return d.trySE(45, LittleEndian) }

func (d *D) ScalarS45LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(45, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S45LE() int64 {
	v, err := d.trySE(45, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S45LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS45LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS45LE, sfns...)
}

func (d *D) FieldS45LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S45LE, sfns...)
}

// Reader S46LE

func (d *D) TryS46LE() (int64, error) { return d.trySE(46, LittleEndian) }

func (d *D) ScalarS46LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(46, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S46LE() int64 {
	v, err := d.trySE(46, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S46LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS46LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS46LE, sfns...)
}

func (d *D) FieldS46LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S46LE, sfns...)
}

// Reader S47LE

func (d *D) TryS47LE() (int64, error) { return d.trySE(47, LittleEndian) }

func (d *D) ScalarS47LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(47, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S47LE() int64 {
	v, err := d.trySE(47, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S47LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS47LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS47LE, sfns...)
}

func (d *D) FieldS47LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S47LE, sfns...)
}

// Reader S48LE

func (d *D) TryS48LE() (int64, error) { return d.trySE(48, LittleEndian) }

func (d *D) ScalarS48LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(48, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S48LE() int64 {
	v, err := d.trySE(48, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S48LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS48LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS48LE, sfns...)
}

func (d *D) FieldS48LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S48LE, sfns...)
}

// Reader S49LE

func (d *D) TryS49LE() (int64, error) { return d.trySE(49, LittleEndian) }

func (d *D) ScalarS49LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(49, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S49LE() int64 {
	v, err := d.trySE(49, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S49LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS49LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS49LE, sfns...)
}

func (d *D) FieldS49LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S49LE, sfns...)
}

// Reader S50LE

func (d *D) TryS50LE() (int64, error) { return d.trySE(50, LittleEndian) }

func (d *D) ScalarS50LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(50, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S50LE() int64 {
	v, err := d.trySE(50, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S50LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS50LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS50LE, sfns...)
}

func (d *D) FieldS50LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S50LE, sfns...)
}

// Reader S51LE

func (d *D) TryS51LE() (int64, error) { return d.trySE(51, LittleEndian) }

func (d *D) ScalarS51LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(51, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S51LE() int64 {
	v, err := d.trySE(51, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S51LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS51LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS51LE, sfns...)
}

func (d *D) FieldS51LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S51LE, sfns...)
}

// Reader S52LE

func (d *D) TryS52LE() (int64, error) { return d.trySE(52, LittleEndian) }

func (d *D) ScalarS52LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(52, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S52LE() int64 {
	v, err := d.trySE(52, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S52LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS52LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS52LE, sfns...)
}

func (d *D) FieldS52LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S52LE, sfns...)
}

// Reader S53LE

func (d *D) TryS53LE() (int64, error) { return d.trySE(53, LittleEndian) }

func (d *D) ScalarS53LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(53, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S53LE() int64 {
	v, err := d.trySE(53, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S53LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS53LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS53LE, sfns...)
}

func (d *D) FieldS53LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S53LE, sfns...)
}

// Reader S54LE

func (d *D) TryS54LE() (int64, error) { return d.trySE(54, LittleEndian) }

func (d *D) ScalarS54LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(54, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S54LE() int64 {
	v, err := d.trySE(54, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S54LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS54LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS54LE, sfns...)
}

func (d *D) FieldS54LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S54LE, sfns...)
}

// Reader S55LE

func (d *D) TryS55LE() (int64, error) { return d.trySE(55, LittleEndian) }

func (d *D) ScalarS55LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(55, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S55LE() int64 {
	v, err := d.trySE(55, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S55LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS55LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS55LE, sfns...)
}

func (d *D) FieldS55LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S55LE, sfns...)
}

// Reader S56LE

func (d *D) TryS56LE() (int64, error) { return d.trySE(56, LittleEndian) }

func (d *D) ScalarS56LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(56, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S56LE() int64 {
	v, err := d.trySE(56, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S56LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS56LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS56LE, sfns...)
}

func (d *D) FieldS56LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S56LE, sfns...)
}

// Reader S57LE

func (d *D) TryS57LE() (int64, error) { return d.trySE(57, LittleEndian) }

func (d *D) ScalarS57LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(57, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S57LE() int64 {
	v, err := d.trySE(57, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S57LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS57LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS57LE, sfns...)
}

func (d *D) FieldS57LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S57LE, sfns...)
}

// Reader S58LE

func (d *D) TryS58LE() (int64, error) { return d.trySE(58, LittleEndian) }

func (d *D) ScalarS58LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(58, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S58LE() int64 {
	v, err := d.trySE(58, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S58LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS58LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS58LE, sfns...)
}

func (d *D) FieldS58LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S58LE, sfns...)
}

// Reader S59LE

func (d *D) TryS59LE() (int64, error) { return d.trySE(59, LittleEndian) }

func (d *D) ScalarS59LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(59, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S59LE() int64 {
	v, err := d.trySE(59, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S59LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS59LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS59LE, sfns...)
}

func (d *D) FieldS59LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S59LE, sfns...)
}

// Reader S60LE

func (d *D) TryS60LE() (int64, error) { return d.trySE(60, LittleEndian) }

func (d *D) ScalarS60LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(60, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S60LE() int64 {
	v, err := d.trySE(60, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S60LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS60LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS60LE, sfns...)
}

func (d *D) FieldS60LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S60LE, sfns...)
}

// Reader S61LE

func (d *D) TryS61LE() (int64, error) { return d.trySE(61, LittleEndian) }

func (d *D) ScalarS61LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(61, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S61LE() int64 {
	v, err := d.trySE(61, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S61LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS61LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS61LE, sfns...)
}

func (d *D) FieldS61LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S61LE, sfns...)
}

// Reader S62LE

func (d *D) TryS62LE() (int64, error) { return d.trySE(62, LittleEndian) }

func (d *D) ScalarS62LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(62, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S62LE() int64 {
	v, err := d.trySE(62, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S62LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS62LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS62LE, sfns...)
}

func (d *D) FieldS62LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S62LE, sfns...)
}

// Reader S63LE

func (d *D) TryS63LE() (int64, error) { return d.trySE(63, LittleEndian) }

func (d *D) ScalarS63LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(63, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S63LE() int64 {
	v, err := d.trySE(63, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S63LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS63LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS63LE, sfns...)
}

func (d *D) FieldS63LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S63LE, sfns...)
}

// Reader S64LE

func (d *D) TryS64LE() (int64, error) { return d.trySE(64, LittleEndian) }

func (d *D) ScalarS64LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(64, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S64LE() int64 {
	v, err := d.trySE(64, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S64LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS64LE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS64LE, sfns...)
}

func (d *D) FieldS64LE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S64LE, sfns...)
}

// Reader S8BE

func (d *D) TryS8BE() (int64, error) { return d.trySE(8, BigEndian) }

func (d *D) ScalarS8BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(8, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S8BE() int64 {
	v, err := d.trySE(8, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S8BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS8BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS8BE, sfns...)
}

func (d *D) FieldS8BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S8BE, sfns...)
}

// Reader S9BE

func (d *D) TryS9BE() (int64, error) { return d.trySE(9, BigEndian) }

func (d *D) ScalarS9BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(9, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S9BE() int64 {
	v, err := d.trySE(9, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S9BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS9BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS9BE, sfns...)
}

func (d *D) FieldS9BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S9BE, sfns...)
}

// Reader S10BE

func (d *D) TryS10BE() (int64, error) { return d.trySE(10, BigEndian) }

func (d *D) ScalarS10BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(10, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S10BE() int64 {
	v, err := d.trySE(10, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S10BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS10BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS10BE, sfns...)
}

func (d *D) FieldS10BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S10BE, sfns...)
}

// Reader S11BE

func (d *D) TryS11BE() (int64, error) { return d.trySE(11, BigEndian) }

func (d *D) ScalarS11BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(11, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S11BE() int64 {
	v, err := d.trySE(11, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S11BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS11BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS11BE, sfns...)
}

func (d *D) FieldS11BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S11BE, sfns...)
}

// Reader S12BE

func (d *D) TryS12BE() (int64, error) { return d.trySE(12, BigEndian) }

func (d *D) ScalarS12BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(12, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S12BE() int64 {
	v, err := d.trySE(12, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S12BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS12BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS12BE, sfns...)
}

func (d *D) FieldS12BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S12BE, sfns...)
}

// Reader S13BE

func (d *D) TryS13BE() (int64, error) { return d.trySE(13, BigEndian) }

func (d *D) ScalarS13BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(13, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S13BE() int64 {
	v, err := d.trySE(13, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S13BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS13BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS13BE, sfns...)
}

func (d *D) FieldS13BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S13BE, sfns...)
}

// Reader S14BE

func (d *D) TryS14BE() (int64, error) { return d.trySE(14, BigEndian) }

func (d *D) ScalarS14BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(14, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S14BE() int64 {
	v, err := d.trySE(14, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S14BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS14BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS14BE, sfns...)
}

func (d *D) FieldS14BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S14BE, sfns...)
}

// Reader S15BE

func (d *D) TryS15BE() (int64, error) { return d.trySE(15, BigEndian) }

func (d *D) ScalarS15BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(15, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S15BE() int64 {
	v, err := d.trySE(15, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S15BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS15BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS15BE, sfns...)
}

func (d *D) FieldS15BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S15BE, sfns...)
}

// Reader S16BE

func (d *D) TryS16BE() (int64, error) { return d.trySE(16, BigEndian) }

func (d *D) ScalarS16BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(16, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S16BE() int64 {
	v, err := d.trySE(16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S16BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS16BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS16BE, sfns...)
}

func (d *D) FieldS16BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S16BE, sfns...)
}

// Reader S17BE

func (d *D) TryS17BE() (int64, error) { return d.trySE(17, BigEndian) }

func (d *D) ScalarS17BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(17, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S17BE() int64 {
	v, err := d.trySE(17, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S17BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS17BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS17BE, sfns...)
}

func (d *D) FieldS17BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S17BE, sfns...)
}

// Reader S18BE

func (d *D) TryS18BE() (int64, error) { return d.trySE(18, BigEndian) }

func (d *D) ScalarS18BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(18, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S18BE() int64 {
	v, err := d.trySE(18, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S18BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS18BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS18BE, sfns...)
}

func (d *D) FieldS18BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S18BE, sfns...)
}

// Reader S19BE

func (d *D) TryS19BE() (int64, error) { return d.trySE(19, BigEndian) }

func (d *D) ScalarS19BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(19, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S19BE() int64 {
	v, err := d.trySE(19, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S19BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS19BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS19BE, sfns...)
}

func (d *D) FieldS19BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S19BE, sfns...)
}

// Reader S20BE

func (d *D) TryS20BE() (int64, error) { return d.trySE(20, BigEndian) }

func (d *D) ScalarS20BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(20, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S20BE() int64 {
	v, err := d.trySE(20, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S20BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS20BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS20BE, sfns...)
}

func (d *D) FieldS20BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S20BE, sfns...)
}

// Reader S21BE

func (d *D) TryS21BE() (int64, error) { return d.trySE(21, BigEndian) }

func (d *D) ScalarS21BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(21, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S21BE() int64 {
	v, err := d.trySE(21, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S21BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS21BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS21BE, sfns...)
}

func (d *D) FieldS21BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S21BE, sfns...)
}

// Reader S22BE

func (d *D) TryS22BE() (int64, error) { return d.trySE(22, BigEndian) }

func (d *D) ScalarS22BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(22, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S22BE() int64 {
	v, err := d.trySE(22, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S22BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS22BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS22BE, sfns...)
}

func (d *D) FieldS22BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S22BE, sfns...)
}

// Reader S23BE

func (d *D) TryS23BE() (int64, error) { return d.trySE(23, BigEndian) }

func (d *D) ScalarS23BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(23, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S23BE() int64 {
	v, err := d.trySE(23, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S23BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS23BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS23BE, sfns...)
}

func (d *D) FieldS23BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S23BE, sfns...)
}

// Reader S24BE

func (d *D) TryS24BE() (int64, error) { return d.trySE(24, BigEndian) }

func (d *D) ScalarS24BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(24, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S24BE() int64 {
	v, err := d.trySE(24, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S24BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS24BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS24BE, sfns...)
}

func (d *D) FieldS24BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S24BE, sfns...)
}

// Reader S25BE

func (d *D) TryS25BE() (int64, error) { return d.trySE(25, BigEndian) }

func (d *D) ScalarS25BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(25, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S25BE() int64 {
	v, err := d.trySE(25, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S25BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS25BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS25BE, sfns...)
}

func (d *D) FieldS25BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S25BE, sfns...)
}

// Reader S26BE

func (d *D) TryS26BE() (int64, error) { return d.trySE(26, BigEndian) }

func (d *D) ScalarS26BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(26, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S26BE() int64 {
	v, err := d.trySE(26, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S26BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS26BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS26BE, sfns...)
}

func (d *D) FieldS26BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S26BE, sfns...)
}

// Reader S27BE

func (d *D) TryS27BE() (int64, error) { return d.trySE(27, BigEndian) }

func (d *D) ScalarS27BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(27, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S27BE() int64 {
	v, err := d.trySE(27, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S27BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS27BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS27BE, sfns...)
}

func (d *D) FieldS27BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S27BE, sfns...)
}

// Reader S28BE

func (d *D) TryS28BE() (int64, error) { return d.trySE(28, BigEndian) }

func (d *D) ScalarS28BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(28, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S28BE() int64 {
	v, err := d.trySE(28, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S28BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS28BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS28BE, sfns...)
}

func (d *D) FieldS28BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S28BE, sfns...)
}

// Reader S29BE

func (d *D) TryS29BE() (int64, error) { return d.trySE(29, BigEndian) }

func (d *D) ScalarS29BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(29, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S29BE() int64 {
	v, err := d.trySE(29, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S29BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS29BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS29BE, sfns...)
}

func (d *D) FieldS29BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S29BE, sfns...)
}

// Reader S30BE

func (d *D) TryS30BE() (int64, error) { return d.trySE(30, BigEndian) }

func (d *D) ScalarS30BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(30, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S30BE() int64 {
	v, err := d.trySE(30, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S30BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS30BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS30BE, sfns...)
}

func (d *D) FieldS30BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S30BE, sfns...)
}

// Reader S31BE

func (d *D) TryS31BE() (int64, error) { return d.trySE(31, BigEndian) }

func (d *D) ScalarS31BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(31, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S31BE() int64 {
	v, err := d.trySE(31, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S31BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS31BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS31BE, sfns...)
}

func (d *D) FieldS31BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S31BE, sfns...)
}

// Reader S32BE

func (d *D) TryS32BE() (int64, error) { return d.trySE(32, BigEndian) }

func (d *D) ScalarS32BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(32, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S32BE() int64 {
	v, err := d.trySE(32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S32BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS32BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS32BE, sfns...)
}

func (d *D) FieldS32BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S32BE, sfns...)
}

// Reader S33BE

func (d *D) TryS33BE() (int64, error) { return d.trySE(33, BigEndian) }

func (d *D) ScalarS33BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(33, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S33BE() int64 {
	v, err := d.trySE(33, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S33BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS33BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS33BE, sfns...)
}

func (d *D) FieldS33BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S33BE, sfns...)
}

// Reader S34BE

func (d *D) TryS34BE() (int64, error) { return d.trySE(34, BigEndian) }

func (d *D) ScalarS34BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(34, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S34BE() int64 {
	v, err := d.trySE(34, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S34BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS34BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS34BE, sfns...)
}

func (d *D) FieldS34BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S34BE, sfns...)
}

// Reader S35BE

func (d *D) TryS35BE() (int64, error) { return d.trySE(35, BigEndian) }

func (d *D) ScalarS35BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(35, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S35BE() int64 {
	v, err := d.trySE(35, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S35BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS35BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS35BE, sfns...)
}

func (d *D) FieldS35BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S35BE, sfns...)
}

// Reader S36BE

func (d *D) TryS36BE() (int64, error) { return d.trySE(36, BigEndian) }

func (d *D) ScalarS36BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(36, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S36BE() int64 {
	v, err := d.trySE(36, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S36BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS36BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS36BE, sfns...)
}

func (d *D) FieldS36BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S36BE, sfns...)
}

// Reader S37BE

func (d *D) TryS37BE() (int64, error) { return d.trySE(37, BigEndian) }

func (d *D) ScalarS37BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(37, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S37BE() int64 {
	v, err := d.trySE(37, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S37BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS37BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS37BE, sfns...)
}

func (d *D) FieldS37BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S37BE, sfns...)
}

// Reader S38BE

func (d *D) TryS38BE() (int64, error) { return d.trySE(38, BigEndian) }

func (d *D) ScalarS38BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(38, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S38BE() int64 {
	v, err := d.trySE(38, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S38BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS38BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS38BE, sfns...)
}

func (d *D) FieldS38BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S38BE, sfns...)
}

// Reader S39BE

func (d *D) TryS39BE() (int64, error) { return d.trySE(39, BigEndian) }

func (d *D) ScalarS39BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(39, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S39BE() int64 {
	v, err := d.trySE(39, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S39BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS39BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS39BE, sfns...)
}

func (d *D) FieldS39BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S39BE, sfns...)
}

// Reader S40BE

func (d *D) TryS40BE() (int64, error) { return d.trySE(40, BigEndian) }

func (d *D) ScalarS40BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(40, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S40BE() int64 {
	v, err := d.trySE(40, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S40BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS40BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS40BE, sfns...)
}

func (d *D) FieldS40BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S40BE, sfns...)
}

// Reader S41BE

func (d *D) TryS41BE() (int64, error) { return d.trySE(41, BigEndian) }

func (d *D) ScalarS41BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(41, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S41BE() int64 {
	v, err := d.trySE(41, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S41BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS41BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS41BE, sfns...)
}

func (d *D) FieldS41BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S41BE, sfns...)
}

// Reader S42BE

func (d *D) TryS42BE() (int64, error) { return d.trySE(42, BigEndian) }

func (d *D) ScalarS42BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(42, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S42BE() int64 {
	v, err := d.trySE(42, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S42BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS42BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS42BE, sfns...)
}

func (d *D) FieldS42BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S42BE, sfns...)
}

// Reader S43BE

func (d *D) TryS43BE() (int64, error) { return d.trySE(43, BigEndian) }

func (d *D) ScalarS43BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(43, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S43BE() int64 {
	v, err := d.trySE(43, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S43BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS43BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS43BE, sfns...)
}

func (d *D) FieldS43BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S43BE, sfns...)
}

// Reader S44BE

func (d *D) TryS44BE() (int64, error) { return d.trySE(44, BigEndian) }

func (d *D) ScalarS44BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(44, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S44BE() int64 {
	v, err := d.trySE(44, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S44BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS44BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS44BE, sfns...)
}

func (d *D) FieldS44BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S44BE, sfns...)
}

// Reader S45BE

func (d *D) TryS45BE() (int64, error) { return d.trySE(45, BigEndian) }

func (d *D) ScalarS45BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(45, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S45BE() int64 {
	v, err := d.trySE(45, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S45BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS45BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS45BE, sfns...)
}

func (d *D) FieldS45BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S45BE, sfns...)
}

// Reader S46BE

func (d *D) TryS46BE() (int64, error) { return d.trySE(46, BigEndian) }

func (d *D) ScalarS46BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(46, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S46BE() int64 {
	v, err := d.trySE(46, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S46BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS46BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS46BE, sfns...)
}

func (d *D) FieldS46BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S46BE, sfns...)
}

// Reader S47BE

func (d *D) TryS47BE() (int64, error) { return d.trySE(47, BigEndian) }

func (d *D) ScalarS47BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(47, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S47BE() int64 {
	v, err := d.trySE(47, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S47BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS47BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS47BE, sfns...)
}

func (d *D) FieldS47BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S47BE, sfns...)
}

// Reader S48BE

func (d *D) TryS48BE() (int64, error) { return d.trySE(48, BigEndian) }

func (d *D) ScalarS48BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(48, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S48BE() int64 {
	v, err := d.trySE(48, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S48BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS48BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS48BE, sfns...)
}

func (d *D) FieldS48BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S48BE, sfns...)
}

// Reader S49BE

func (d *D) TryS49BE() (int64, error) { return d.trySE(49, BigEndian) }

func (d *D) ScalarS49BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(49, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S49BE() int64 {
	v, err := d.trySE(49, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S49BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS49BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS49BE, sfns...)
}

func (d *D) FieldS49BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S49BE, sfns...)
}

// Reader S50BE

func (d *D) TryS50BE() (int64, error) { return d.trySE(50, BigEndian) }

func (d *D) ScalarS50BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(50, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S50BE() int64 {
	v, err := d.trySE(50, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S50BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS50BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS50BE, sfns...)
}

func (d *D) FieldS50BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S50BE, sfns...)
}

// Reader S51BE

func (d *D) TryS51BE() (int64, error) { return d.trySE(51, BigEndian) }

func (d *D) ScalarS51BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(51, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S51BE() int64 {
	v, err := d.trySE(51, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S51BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS51BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS51BE, sfns...)
}

func (d *D) FieldS51BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S51BE, sfns...)
}

// Reader S52BE

func (d *D) TryS52BE() (int64, error) { return d.trySE(52, BigEndian) }

func (d *D) ScalarS52BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(52, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S52BE() int64 {
	v, err := d.trySE(52, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S52BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS52BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS52BE, sfns...)
}

func (d *D) FieldS52BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S52BE, sfns...)
}

// Reader S53BE

func (d *D) TryS53BE() (int64, error) { return d.trySE(53, BigEndian) }

func (d *D) ScalarS53BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(53, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S53BE() int64 {
	v, err := d.trySE(53, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S53BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS53BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS53BE, sfns...)
}

func (d *D) FieldS53BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S53BE, sfns...)
}

// Reader S54BE

func (d *D) TryS54BE() (int64, error) { return d.trySE(54, BigEndian) }

func (d *D) ScalarS54BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(54, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S54BE() int64 {
	v, err := d.trySE(54, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S54BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS54BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS54BE, sfns...)
}

func (d *D) FieldS54BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S54BE, sfns...)
}

// Reader S55BE

func (d *D) TryS55BE() (int64, error) { return d.trySE(55, BigEndian) }

func (d *D) ScalarS55BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(55, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S55BE() int64 {
	v, err := d.trySE(55, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S55BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS55BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS55BE, sfns...)
}

func (d *D) FieldS55BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S55BE, sfns...)
}

// Reader S56BE

func (d *D) TryS56BE() (int64, error) { return d.trySE(56, BigEndian) }

func (d *D) ScalarS56BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(56, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S56BE() int64 {
	v, err := d.trySE(56, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S56BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS56BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS56BE, sfns...)
}

func (d *D) FieldS56BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S56BE, sfns...)
}

// Reader S57BE

func (d *D) TryS57BE() (int64, error) { return d.trySE(57, BigEndian) }

func (d *D) ScalarS57BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(57, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S57BE() int64 {
	v, err := d.trySE(57, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S57BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS57BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS57BE, sfns...)
}

func (d *D) FieldS57BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S57BE, sfns...)
}

// Reader S58BE

func (d *D) TryS58BE() (int64, error) { return d.trySE(58, BigEndian) }

func (d *D) ScalarS58BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(58, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S58BE() int64 {
	v, err := d.trySE(58, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S58BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS58BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS58BE, sfns...)
}

func (d *D) FieldS58BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S58BE, sfns...)
}

// Reader S59BE

func (d *D) TryS59BE() (int64, error) { return d.trySE(59, BigEndian) }

func (d *D) ScalarS59BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(59, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S59BE() int64 {
	v, err := d.trySE(59, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S59BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS59BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS59BE, sfns...)
}

func (d *D) FieldS59BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S59BE, sfns...)
}

// Reader S60BE

func (d *D) TryS60BE() (int64, error) { return d.trySE(60, BigEndian) }

func (d *D) ScalarS60BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(60, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S60BE() int64 {
	v, err := d.trySE(60, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S60BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS60BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS60BE, sfns...)
}

func (d *D) FieldS60BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S60BE, sfns...)
}

// Reader S61BE

func (d *D) TryS61BE() (int64, error) { return d.trySE(61, BigEndian) }

func (d *D) ScalarS61BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(61, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S61BE() int64 {
	v, err := d.trySE(61, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S61BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS61BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS61BE, sfns...)
}

func (d *D) FieldS61BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S61BE, sfns...)
}

// Reader S62BE

func (d *D) TryS62BE() (int64, error) { return d.trySE(62, BigEndian) }

func (d *D) ScalarS62BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(62, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S62BE() int64 {
	v, err := d.trySE(62, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S62BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS62BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS62BE, sfns...)
}

func (d *D) FieldS62BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S62BE, sfns...)
}

// Reader S63BE

func (d *D) TryS63BE() (int64, error) { return d.trySE(63, BigEndian) }

func (d *D) ScalarS63BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(63, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S63BE() int64 {
	v, err := d.trySE(63, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S63BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS63BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS63BE, sfns...)
}

func (d *D) FieldS63BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S63BE, sfns...)
}

// Reader S64BE

func (d *D) TryS64BE() (int64, error) { return d.trySE(64, BigEndian) }

func (d *D) ScalarS64BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.trySE(64, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) S64BE() int64 {
	v, err := d.trySE(64, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "S64BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldS64BE(name string, sfns ...ScalarFn) (int64, error) {
	return d.TryFieldSFn(name, (*D).TryS64BE, sfns...)
}

func (d *D) FieldS64BE(name string, sfns ...ScalarFn) int64 {
	return d.FieldSFn(name, (*D).S64BE, sfns...)
}

// Reader FE

func (d *D) TryFE(nBits int, endian Endian) (float64, error) { return d.tryFE(nBits, endian) }

func (d *D) ScalarFE(nBits int, endian Endian) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(nBits, endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FE(nBits int, endian Endian) float64 {
	v, err := d.tryFE(nBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFE(name string, nBits int, endian Endian, sfns ...ScalarFn) (float64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarFE(nBits, endian), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualF(), err
}

func (d *D) FieldFE(name string, nBits int, endian Endian, sfns ...ScalarFn) float64 {
	v, err := d.TryFieldFE(name, nBits, endian, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FE", Pos: d.Pos()})
	}
	return v
}

// Reader F

func (d *D) TryF(nBits int) (float64, error) { return d.tryFE(nBits, d.Endian) }

func (d *D) ScalarF(nBits int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(nBits, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F(nBits int) float64 {
	v, err := d.tryFE(nBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF(name string, nBits int, sfns ...ScalarFn) (float64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarF(nBits), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualF(), err
}

func (d *D) FieldF(name string, nBits int, sfns ...ScalarFn) float64 {
	v, err := d.TryFieldF(name, nBits, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "F", Pos: d.Pos()})
	}
	return v
}

// Reader F16

func (d *D) TryF16() (float64, error) { return d.tryFE(16, d.Endian) }

func (d *D) ScalarF16() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(16, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F16() float64 {
	v, err := d.tryFE(16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F16", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF16(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF16, sfns...)
}

func (d *D) FieldF16(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F16, sfns...)
}

// Reader F32

func (d *D) TryF32() (float64, error) { return d.tryFE(32, d.Endian) }

func (d *D) ScalarF32() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(32, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F32() float64 {
	v, err := d.tryFE(32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F32", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF32(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF32, sfns...)
}

func (d *D) FieldF32(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F32, sfns...)
}

// Reader F64

func (d *D) TryF64() (float64, error) { return d.tryFE(64, d.Endian) }

func (d *D) ScalarF64() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(64, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F64() float64 {
	v, err := d.tryFE(64, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "F64", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF64(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF64, sfns...)
}

func (d *D) FieldF64(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F64, sfns...)
}

// Reader F16LE

func (d *D) TryF16LE() (float64, error) { return d.tryFE(16, LittleEndian) }

func (d *D) ScalarF16LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(16, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F16LE() float64 {
	v, err := d.tryFE(16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F16LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF16LE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF16LE, sfns...)
}

func (d *D) FieldF16LE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F16LE, sfns...)
}

// Reader F32LE

func (d *D) TryF32LE() (float64, error) { return d.tryFE(32, LittleEndian) }

func (d *D) ScalarF32LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(32, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F32LE() float64 {
	v, err := d.tryFE(32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F32LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF32LE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF32LE, sfns...)
}

func (d *D) FieldF32LE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F32LE, sfns...)
}

// Reader F64LE

func (d *D) TryF64LE() (float64, error) { return d.tryFE(64, LittleEndian) }

func (d *D) ScalarF64LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(64, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F64LE() float64 {
	v, err := d.tryFE(64, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F64LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF64LE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF64LE, sfns...)
}

func (d *D) FieldF64LE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F64LE, sfns...)
}

// Reader F16BE

func (d *D) TryF16BE() (float64, error) { return d.tryFE(16, BigEndian) }

func (d *D) ScalarF16BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(16, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F16BE() float64 {
	v, err := d.tryFE(16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F16BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF16BE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF16BE, sfns...)
}

func (d *D) FieldF16BE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F16BE, sfns...)
}

// Reader F32BE

func (d *D) TryF32BE() (float64, error) { return d.tryFE(32, BigEndian) }

func (d *D) ScalarF32BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(32, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F32BE() float64 {
	v, err := d.tryFE(32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F32BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF32BE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF32BE, sfns...)
}

func (d *D) FieldF32BE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F32BE, sfns...)
}

// Reader F64BE

func (d *D) TryF64BE() (float64, error) { return d.tryFE(64, BigEndian) }

func (d *D) ScalarF64BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFE(64, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) F64BE() float64 {
	v, err := d.tryFE(64, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "F64BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldF64BE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryF64BE, sfns...)
}

func (d *D) FieldF64BE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).F64BE, sfns...)
}

// Reader FPE

func (d *D) TryFPE(nBits int, fBits int64, endian Endian) (float64, error) {
	return d.tryFPE(nBits, fBits, endian)
}

func (d *D) ScalarFPE(nBits int, fBits int64, endian Endian) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(nBits, fBits, endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FPE(nBits int, fBits int64, endian Endian) float64 {
	v, err := d.tryFPE(nBits, fBits, endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FPE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFPE(name string, nBits int, fBits int64, endian Endian, sfns ...ScalarFn) (float64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarFPE(nBits, fBits, endian), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualF(), err
}

func (d *D) FieldFPE(name string, nBits int, fBits int64, endian Endian, sfns ...ScalarFn) float64 {
	v, err := d.TryFieldFPE(name, nBits, fBits, endian, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FPE", Pos: d.Pos()})
	}
	return v
}

// Reader FP

func (d *D) TryFP(nBits int, fBits int64) (float64, error) { return d.tryFPE(nBits, fBits, d.Endian) }

func (d *D) ScalarFP(nBits int, fBits int64) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(nBits, fBits, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP(nBits int, fBits int64) float64 {
	v, err := d.tryFPE(nBits, fBits, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP(name string, nBits int, fBits int64, sfns ...ScalarFn) (float64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarFP(nBits, fBits), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualF(), err
}

func (d *D) FieldFP(name string, nBits int, fBits int64, sfns ...ScalarFn) float64 {
	v, err := d.TryFieldFP(name, nBits, fBits, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "FP", Pos: d.Pos()})
	}
	return v
}

// Reader FP16

func (d *D) TryFP16() (float64, error) { return d.tryFPE(16, 8, d.Endian) }

func (d *D) ScalarFP16() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(16, 8, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP16() float64 {
	v, err := d.tryFPE(16, 8, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP16", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP16(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP16, sfns...)
}

func (d *D) FieldFP16(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP16, sfns...)
}

// Reader FP32

func (d *D) TryFP32() (float64, error) { return d.tryFPE(32, 16, d.Endian) }

func (d *D) ScalarFP32() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(32, 16, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP32() float64 {
	v, err := d.tryFPE(32, 16, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP32", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP32(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP32, sfns...)
}

func (d *D) FieldFP32(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP32, sfns...)
}

// Reader FP64

func (d *D) TryFP64() (float64, error) { return d.tryFPE(64, 32, d.Endian) }

func (d *D) ScalarFP64() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(64, 32, d.Endian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP64() float64 {
	v, err := d.tryFPE(64, 32, d.Endian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP64", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP64(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP64, sfns...)
}

func (d *D) FieldFP64(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP64, sfns...)
}

// Reader FP16LE

func (d *D) TryFP16LE() (float64, error) { return d.tryFPE(16, 8, LittleEndian) }

func (d *D) ScalarFP16LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(16, 8, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP16LE() float64 {
	v, err := d.tryFPE(16, 8, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP16LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP16LE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP16LE, sfns...)
}

func (d *D) FieldFP16LE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP16LE, sfns...)
}

// Reader FP32LE

func (d *D) TryFP32LE() (float64, error) { return d.tryFPE(32, 16, LittleEndian) }

func (d *D) ScalarFP32LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(32, 16, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP32LE() float64 {
	v, err := d.tryFPE(32, 16, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP32LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP32LE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP32LE, sfns...)
}

func (d *D) FieldFP32LE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP32LE, sfns...)
}

// Reader FP64LE

func (d *D) TryFP64LE() (float64, error) { return d.tryFPE(64, 32, LittleEndian) }

func (d *D) ScalarFP64LE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(64, 32, LittleEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP64LE() float64 {
	v, err := d.tryFPE(64, 32, LittleEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP64LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP64LE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP64LE, sfns...)
}

func (d *D) FieldFP64LE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP64LE, sfns...)
}

// Reader FP16BE

func (d *D) TryFP16BE() (float64, error) { return d.tryFPE(16, 8, BigEndian) }

func (d *D) ScalarFP16BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(16, 8, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP16BE() float64 {
	v, err := d.tryFPE(16, 8, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP16BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP16BE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP16BE, sfns...)
}

func (d *D) FieldFP16BE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP16BE, sfns...)
}

// Reader FP32BE

func (d *D) TryFP32BE() (float64, error) { return d.tryFPE(32, 16, BigEndian) }

func (d *D) ScalarFP32BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(32, 16, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP32BE() float64 {
	v, err := d.tryFPE(32, 16, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP32BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP32BE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP32BE, sfns...)
}

func (d *D) FieldFP32BE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP32BE, sfns...)
}

// Reader FP64BE

func (d *D) TryFP64BE() (float64, error) { return d.tryFPE(64, 32, BigEndian) }

func (d *D) ScalarFP64BE() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryFPE(64, 32, BigEndian)
		s.Actual = v
		return s, err
	}
}

func (d *D) FP64BE() float64 {
	v, err := d.tryFPE(64, 32, BigEndian)
	if err != nil {
		panic(IOError{Err: err, Op: "FP64BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldFP64BE(name string, sfns ...ScalarFn) (float64, error) {
	return d.TryFieldFFn(name, (*D).TryFP64BE, sfns...)
}

func (d *D) FieldFP64BE(name string, sfns ...ScalarFn) float64 {
	return d.FieldFFn(name, (*D).FP64BE, sfns...)
}

// Reader Unary

func (d *D) TryUnary(ov uint64) (uint64, error) { return d.tryUnary(ov) }

func (d *D) ScalarUnary(ov uint64) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryUnary(ov)
		s.Actual = v
		return s, err
	}
}

func (d *D) Unary(ov uint64) uint64 {
	v, err := d.tryUnary(ov)
	if err != nil {
		panic(IOError{Err: err, Op: "Unary", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUnary(name string, ov uint64, sfns ...ScalarFn) (uint64, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUnary(ov), sfns...)
	if err != nil {
		return 0, err
	}
	return v.ActualU(), err
}

func (d *D) FieldUnary(name string, ov uint64, sfns ...ScalarFn) uint64 {
	v, err := d.TryFieldUnary(name, ov, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "Unary", Pos: d.Pos()})
	}
	return v
}

// Reader UTF8

func (d *D) TryUTF8(nBytes int) (string, error) { return d.tryText(nBytes, UTF8BOM) }

func (d *D) ScalarUTF8(nBytes int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryText(nBytes, UTF8BOM)
		s.Actual = v
		return s, err
	}
}

func (d *D) UTF8(nBytes int) string {
	v, err := d.tryText(nBytes, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUTF8(name string, nBytes int, sfns ...ScalarFn) (string, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUTF8(nBytes), sfns...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}

func (d *D) FieldUTF8(name string, nBytes int, sfns ...ScalarFn) string {
	v, err := d.TryFieldUTF8(name, nBytes, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8", Pos: d.Pos()})
	}
	return v
}

// Reader UTF16

func (d *D) TryUTF16(nBytes int) (string, error) { return d.tryText(nBytes, UTF16BOM) }

func (d *D) ScalarUTF16(nBytes int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryText(nBytes, UTF16BOM)
		s.Actual = v
		return s, err
	}
}

func (d *D) UTF16(nBytes int) string {
	v, err := d.tryText(nBytes, UTF16BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUTF16(name string, nBytes int, sfns ...ScalarFn) (string, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUTF16(nBytes), sfns...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}

func (d *D) FieldUTF16(name string, nBytes int, sfns ...ScalarFn) string {
	v, err := d.TryFieldUTF16(name, nBytes, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16", Pos: d.Pos()})
	}
	return v
}

// Reader UTF16LE

func (d *D) TryUTF16LE(nBytes int) (string, error) { return d.tryText(nBytes, UTF16LE) }

func (d *D) ScalarUTF16LE(nBytes int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryText(nBytes, UTF16LE)
		s.Actual = v
		return s, err
	}
}

func (d *D) UTF16LE(nBytes int) string {
	v, err := d.tryText(nBytes, UTF16LE)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16LE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUTF16LE(name string, nBytes int, sfns ...ScalarFn) (string, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUTF16LE(nBytes), sfns...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}

func (d *D) FieldUTF16LE(name string, nBytes int, sfns ...ScalarFn) string {
	v, err := d.TryFieldUTF16LE(name, nBytes, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16LE", Pos: d.Pos()})
	}
	return v
}

// Reader UTF16BE

func (d *D) TryUTF16BE(nBytes int) (string, error) { return d.tryText(nBytes, UTF16BE) }

func (d *D) ScalarUTF16BE(nBytes int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryText(nBytes, UTF16BE)
		s.Actual = v
		return s, err
	}
}

func (d *D) UTF16BE(nBytes int) string {
	v, err := d.tryText(nBytes, UTF16BE)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF16BE", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUTF16BE(name string, nBytes int, sfns ...ScalarFn) (string, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUTF16BE(nBytes), sfns...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}

func (d *D) FieldUTF16BE(name string, nBytes int, sfns ...ScalarFn) string {
	v, err := d.TryFieldUTF16BE(name, nBytes, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF16BE", Pos: d.Pos()})
	}
	return v
}

// Reader UTF8ShortString

func (d *D) TryUTF8ShortString(nBytes int) (string, error) {
	return d.tryLenPrefixedText(8, nBytes, UTF8BOM)
}

func (d *D) ScalarUTF8ShortString(nBytes int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryLenPrefixedText(8, nBytes, UTF8BOM)
		s.Actual = v
		return s, err
	}
}

func (d *D) UTF8ShortString(nBytes int) string {
	v, err := d.tryLenPrefixedText(8, nBytes, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8ShortString", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUTF8ShortString(name string, nBytes int, sfns ...ScalarFn) (string, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUTF8ShortString(nBytes), sfns...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}

func (d *D) FieldUTF8ShortString(name string, nBytes int, sfns ...ScalarFn) string {
	v, err := d.TryFieldUTF8ShortString(name, nBytes, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8ShortString", Pos: d.Pos()})
	}
	return v
}

// Reader UTF8NullTerminated

func (d *D) TryUTF8NullTerminated() (string, error) { return d.tryNullTerminatedText(1, UTF8BOM) }

func (d *D) ScalarUTF8NullTerminated() func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryNullTerminatedText(1, UTF8BOM)
		s.Actual = v
		return s, err
	}
}

func (d *D) UTF8NullTerminated() string {
	v, err := d.tryNullTerminatedText(1, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8NullTerminated", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUTF8NullTerminated(name string, sfns ...ScalarFn) (string, error) {
	return d.TryFieldStrFn(name, (*D).TryUTF8NullTerminated, sfns...)
}

func (d *D) FieldUTF8NullTerminated(name string, sfns ...ScalarFn) string {
	return d.FieldStrFn(name, (*D).UTF8NullTerminated, sfns...)
}

// Reader UTF8NullTerminatedLen

func (d *D) TryUTF8NullTerminatedLen(fixedBytes int) (string, error) {
	return d.tryNullTerminatedLenText(fixedBytes, UTF8BOM)
}

func (d *D) ScalarUTF8NullTerminatedLen(fixedBytes int) func(Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		v, err := d.tryNullTerminatedLenText(fixedBytes, UTF8BOM)
		s.Actual = v
		return s, err
	}
}

func (d *D) UTF8NullTerminatedLen(fixedBytes int) string {
	v, err := d.tryNullTerminatedLenText(fixedBytes, UTF8BOM)
	if err != nil {
		panic(IOError{Err: err, Op: "UTF8NullTerminatedLen", Pos: d.Pos()})
	}
	return v
}

func (d *D) TryFieldUTF8NullTerminatedLen(name string, fixedBytes int, sfns ...ScalarFn) (string, error) {
	v, err := d.TryFieldScalar(name, d.ScalarUTF8NullTerminatedLen(fixedBytes), sfns...)
	if err != nil {
		return "", err
	}
	return v.ActualStr(), err
}

func (d *D) FieldUTF8NullTerminatedLen(name string, fixedBytes int, sfns ...ScalarFn) string {
	v, err := d.TryFieldUTF8NullTerminatedLen(name, fixedBytes, sfns...)
	if err != nil {
		panic(IOError{Err: err, Name: name, Op: "UTF8NullTerminatedLen", Pos: d.Pos()})
	}
	return v
}
