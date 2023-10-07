// Code below generated from scalar_gen.go.tmpl
package scalar

import (
	"fmt"
	"math/big"

	"github.com/wader/fq/pkg/bitio"
)

// Type Any
// does not use embedding for common fields as it works poorly with struct literals
type Any struct {
	Sym         any
	Description string
	Flags       Flags
	Actual      any
}

// interp.Scalarable
func (s Any) ScalarActual() any { return s.Actual }
func (s Any) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s Any) ScalarSym() any                     { return s.Sym }
func (s Any) ScalarDescription() string          { return s.Description }
func (s Any) ScalarFlags() Flags                 { return s.Flags }
func (s Any) ScalarDisplayFormat() DisplayFormat { return 0 }

func AnyActual(v any) AnyMapper {
	return AnyFn(func(s Any) (Any, error) { s.Actual = v; return s, nil })
}
func AnySym(v any) AnyMapper {
	return AnyFn(func(s Any) (Any, error) { s.Sym = v; return s, nil })
}
func AnyDescription(v string) AnyMapper {
	return AnyFn(func(s Any) (Any, error) { s.Description = v; return s, nil })
}

type AnyMapper interface {
	MapAny(Any) (Any, error)
}

// AnyFn map actual Any using f
type AnyFn func(s Any) (Any, error)

func (fn AnyFn) MapAny(s Any) (Any, error) {
	return fn(s)
}

// AnyActualFn map actual Any using f
type AnyActualFn func(a any) any

// TODO: error?
func (fn AnyActualFn) MapAny(s Any) (Any, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// AnySymFn map sym Any using f
type AnySymFn func(a any) any

func (f AnySymFn) MapAny(s Any) (Any, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// AnyDescriptionFn map sym Any using f
type AnyDescriptionFn func(a string) string

func (f AnyDescriptionFn) MapAny(s Any) (Any, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s Any) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s Any) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s Any) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s Any) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s Any) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s Any) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s Any) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s Any) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s Any) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s Any) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s Any) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s Any) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s Any) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s Any) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s Any) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s Any) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as any", s.Sym, s.Sym))
	}
	return v
}

// Type BigInt
// does not use embedding for common fields as it works poorly with struct literals
type BigInt struct {
	Sym           any
	Description   string
	Flags         Flags
	Actual        *big.Int
	DisplayFormat DisplayFormat
}

// interp.Scalarable
func (s BigInt) ScalarActual() any { return s.Actual }
func (s BigInt) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s BigInt) ScalarSym() any                     { return s.Sym }
func (s BigInt) ScalarDescription() string          { return s.Description }
func (s BigInt) ScalarFlags() Flags                 { return s.Flags }
func (s BigInt) ScalarDisplayFormat() DisplayFormat { return s.DisplayFormat }

func BigIntActual(v *big.Int) BigIntMapper {
	return BigIntFn(func(s BigInt) (BigInt, error) { s.Actual = v; return s, nil })
}
func BigIntSym(v any) BigIntMapper {
	return BigIntFn(func(s BigInt) (BigInt, error) { s.Sym = v; return s, nil })
}
func BigIntDescription(v string) BigIntMapper {
	return BigIntFn(func(s BigInt) (BigInt, error) { s.Description = v; return s, nil })
}

type BigIntMapper interface {
	MapBigInt(BigInt) (BigInt, error)
}

// BigIntFn map actual BigInt using f
type BigIntFn func(s BigInt) (BigInt, error)

func (fn BigIntFn) MapBigInt(s BigInt) (BigInt, error) {
	return fn(s)
}

// BigIntActualFn map actual BigInt using f
type BigIntActualFn func(a *big.Int) *big.Int

// TODO: error?
func (fn BigIntActualFn) MapBigInt(s BigInt) (BigInt, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// BigIntSymFn map sym BigInt using f
type BigIntSymFn func(a any) any

func (f BigIntSymFn) MapBigInt(s BigInt) (BigInt, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// BigIntDescriptionFn map sym BigInt using f
type BigIntDescriptionFn func(a string) string

func (f BigIntDescriptionFn) MapBigInt(s BigInt) (BigInt, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s BigInt) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s BigInt) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s BigInt) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s BigInt) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s BigInt) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s BigInt) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s BigInt) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s BigInt) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s BigInt) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s BigInt) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s BigInt) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s BigInt) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s BigInt) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s BigInt) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s BigInt) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s BigInt) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as *big.Int", s.Sym, s.Sym))
	}
	return v
}

// Type BitBuf
// does not use embedding for common fields as it works poorly with struct literals
type BitBuf struct {
	Sym         any
	Description string
	Flags       Flags
	Actual      bitio.ReaderAtSeeker
}

// interp.Scalarable
func (s BitBuf) ScalarActual() any { return s.Actual }
func (s BitBuf) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s BitBuf) ScalarSym() any                     { return s.Sym }
func (s BitBuf) ScalarDescription() string          { return s.Description }
func (s BitBuf) ScalarFlags() Flags                 { return s.Flags }
func (s BitBuf) ScalarDisplayFormat() DisplayFormat { return 0 }

func BitBufActual(v bitio.ReaderAtSeeker) BitBufMapper {
	return BitBufFn(func(s BitBuf) (BitBuf, error) { s.Actual = v; return s, nil })
}
func BitBufSym(v any) BitBufMapper {
	return BitBufFn(func(s BitBuf) (BitBuf, error) { s.Sym = v; return s, nil })
}
func BitBufDescription(v string) BitBufMapper {
	return BitBufFn(func(s BitBuf) (BitBuf, error) { s.Description = v; return s, nil })
}

type BitBufMapper interface {
	MapBitBuf(BitBuf) (BitBuf, error)
}

// BitBufFn map actual BitBuf using f
type BitBufFn func(s BitBuf) (BitBuf, error)

func (fn BitBufFn) MapBitBuf(s BitBuf) (BitBuf, error) {
	return fn(s)
}

// BitBufActualFn map actual BitBuf using f
type BitBufActualFn func(a bitio.ReaderAtSeeker) bitio.ReaderAtSeeker

// TODO: error?
func (fn BitBufActualFn) MapBitBuf(s BitBuf) (BitBuf, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// BitBufSymFn map sym BitBuf using f
type BitBufSymFn func(a any) any

func (f BitBufSymFn) MapBitBuf(s BitBuf) (BitBuf, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// BitBufDescriptionFn map sym BitBuf using f
type BitBufDescriptionFn func(a string) string

func (f BitBufDescriptionFn) MapBitBuf(s BitBuf) (BitBuf, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s BitBuf) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s BitBuf) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s BitBuf) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s BitBuf) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s BitBuf) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s BitBuf) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s BitBuf) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s BitBuf) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s BitBuf) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s BitBuf) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s BitBuf) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s BitBuf) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s BitBuf) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s BitBuf) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s BitBuf) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s BitBuf) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bitio.ReaderAtSeeker", s.Sym, s.Sym))
	}
	return v
}

// Type Bool
// does not use embedding for common fields as it works poorly with struct literals
type Bool struct {
	Sym         any
	Description string
	Flags       Flags
	Actual      bool
}

// interp.Scalarable
func (s Bool) ScalarActual() any { return s.Actual }
func (s Bool) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s Bool) ScalarSym() any                     { return s.Sym }
func (s Bool) ScalarDescription() string          { return s.Description }
func (s Bool) ScalarFlags() Flags                 { return s.Flags }
func (s Bool) ScalarDisplayFormat() DisplayFormat { return 0 }

func BoolActual(v bool) BoolMapper {
	return BoolFn(func(s Bool) (Bool, error) { s.Actual = v; return s, nil })
}
func BoolSym(v any) BoolMapper {
	return BoolFn(func(s Bool) (Bool, error) { s.Sym = v; return s, nil })
}
func BoolDescription(v string) BoolMapper {
	return BoolFn(func(s Bool) (Bool, error) { s.Description = v; return s, nil })
}

type BoolMapper interface {
	MapBool(Bool) (Bool, error)
}

// BoolFn map actual Bool using f
type BoolFn func(s Bool) (Bool, error)

func (fn BoolFn) MapBool(s Bool) (Bool, error) {
	return fn(s)
}

// BoolActualFn map actual Bool using f
type BoolActualFn func(a bool) bool

// TODO: error?
func (fn BoolActualFn) MapBool(s Bool) (Bool, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// BoolSymFn map sym Bool using f
type BoolSymFn func(a any) any

func (f BoolSymFn) MapBool(s Bool) (Bool, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// BoolDescriptionFn map sym Bool using f
type BoolDescriptionFn func(a string) string

func (f BoolDescriptionFn) MapBool(s Bool) (Bool, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s Bool) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s Bool) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s Bool) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s Bool) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s Bool) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s Bool) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s Bool) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s Bool) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s Bool) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s Bool) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s Bool) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s Bool) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s Bool) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s Bool) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s Bool) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s Bool) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as bool", s.Sym, s.Sym))
	}
	return v
}

// Type Flt
// does not use embedding for common fields as it works poorly with struct literals
type Flt struct {
	Sym         any
	Description string
	Flags       Flags
	Actual      float64
}

// interp.Scalarable
func (s Flt) ScalarActual() any { return s.Actual }
func (s Flt) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s Flt) ScalarSym() any                     { return s.Sym }
func (s Flt) ScalarDescription() string          { return s.Description }
func (s Flt) ScalarFlags() Flags                 { return s.Flags }
func (s Flt) ScalarDisplayFormat() DisplayFormat { return 0 }

func FltActual(v float64) FltMapper {
	return FltFn(func(s Flt) (Flt, error) { s.Actual = v; return s, nil })
}
func FltSym(v any) FltMapper {
	return FltFn(func(s Flt) (Flt, error) { s.Sym = v; return s, nil })
}
func FltDescription(v string) FltMapper {
	return FltFn(func(s Flt) (Flt, error) { s.Description = v; return s, nil })
}

type FltMapper interface {
	MapFlt(Flt) (Flt, error)
}

// FltFn map actual Flt using f
type FltFn func(s Flt) (Flt, error)

func (fn FltFn) MapFlt(s Flt) (Flt, error) {
	return fn(s)
}

// FltActualFn map actual Flt using f
type FltActualFn func(a float64) float64

// TODO: error?
func (fn FltActualFn) MapFlt(s Flt) (Flt, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// FltSymFn map sym Flt using f
type FltSymFn func(a any) any

func (f FltSymFn) MapFlt(s Flt) (Flt, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// FltDescriptionFn map sym Flt using f
type FltDescriptionFn func(a string) string

func (f FltDescriptionFn) MapFlt(s Flt) (Flt, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s Flt) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s Flt) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s Flt) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s Flt) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s Flt) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s Flt) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s Flt) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s Flt) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s Flt) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s Flt) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s Flt) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s Flt) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s Flt) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s Flt) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s Flt) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s Flt) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as float64", s.Sym, s.Sym))
	}
	return v
}

// Type Sint
// does not use embedding for common fields as it works poorly with struct literals
type Sint struct {
	Sym           any
	Description   string
	Flags         Flags
	Actual        int64
	DisplayFormat DisplayFormat
}

// interp.Scalarable
func (s Sint) ScalarActual() any { return s.Actual }
func (s Sint) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s Sint) ScalarSym() any                     { return s.Sym }
func (s Sint) ScalarDescription() string          { return s.Description }
func (s Sint) ScalarFlags() Flags                 { return s.Flags }
func (s Sint) ScalarDisplayFormat() DisplayFormat { return s.DisplayFormat }

func SintActual(v int64) SintMapper {
	return SintFn(func(s Sint) (Sint, error) { s.Actual = v; return s, nil })
}
func SintSym(v any) SintMapper {
	return SintFn(func(s Sint) (Sint, error) { s.Sym = v; return s, nil })
}
func SintDescription(v string) SintMapper {
	return SintFn(func(s Sint) (Sint, error) { s.Description = v; return s, nil })
}

type SintMapper interface {
	MapSint(Sint) (Sint, error)
}

// SintFn map actual Sint using f
type SintFn func(s Sint) (Sint, error)

func (fn SintFn) MapSint(s Sint) (Sint, error) {
	return fn(s)
}

// SintActualFn map actual Sint using f
type SintActualFn func(a int64) int64

// TODO: error?
func (fn SintActualFn) MapSint(s Sint) (Sint, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// SintSymFn map sym Sint using f
type SintSymFn func(a any) any

func (f SintSymFn) MapSint(s Sint) (Sint, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// SintDescriptionFn map sym Sint using f
type SintDescriptionFn func(a string) string

func (f SintDescriptionFn) MapSint(s Sint) (Sint, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s Sint) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s Sint) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s Sint) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s Sint) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s Sint) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s Sint) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s Sint) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s Sint) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s Sint) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s Sint) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s Sint) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s Sint) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s Sint) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s Sint) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s Sint) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s Sint) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as int64", s.Sym, s.Sym))
	}
	return v
}

// Type Str
// does not use embedding for common fields as it works poorly with struct literals
type Str struct {
	Sym         any
	Description string
	Flags       Flags
	Actual      string
}

// interp.Scalarable
func (s Str) ScalarActual() any { return s.Actual }
func (s Str) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s Str) ScalarSym() any                     { return s.Sym }
func (s Str) ScalarDescription() string          { return s.Description }
func (s Str) ScalarFlags() Flags                 { return s.Flags }
func (s Str) ScalarDisplayFormat() DisplayFormat { return 0 }

func StrActual(v string) StrMapper {
	return StrFn(func(s Str) (Str, error) { s.Actual = v; return s, nil })
}
func StrSym(v any) StrMapper {
	return StrFn(func(s Str) (Str, error) { s.Sym = v; return s, nil })
}
func StrDescription(v string) StrMapper {
	return StrFn(func(s Str) (Str, error) { s.Description = v; return s, nil })
}

type StrMapper interface {
	MapStr(Str) (Str, error)
}

// StrFn map actual Str using f
type StrFn func(s Str) (Str, error)

func (fn StrFn) MapStr(s Str) (Str, error) {
	return fn(s)
}

// StrActualFn map actual Str using f
type StrActualFn func(a string) string

// TODO: error?
func (fn StrActualFn) MapStr(s Str) (Str, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// StrSymFn map sym Str using f
type StrSymFn func(a any) any

func (f StrSymFn) MapStr(s Str) (Str, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// StrDescriptionFn map sym Str using f
type StrDescriptionFn func(a string) string

func (f StrDescriptionFn) MapStr(s Str) (Str, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s Str) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s Str) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s Str) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s Str) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s Str) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s Str) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s Str) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s Str) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s Str) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s Str) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s Str) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s Str) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s Str) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s Str) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s Str) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s Str) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as string", s.Sym, s.Sym))
	}
	return v
}

// Type Uint
// does not use embedding for common fields as it works poorly with struct literals
type Uint struct {
	Sym           any
	Description   string
	Flags         Flags
	Actual        uint64
	DisplayFormat DisplayFormat
}

// interp.Scalarable
func (s Uint) ScalarActual() any { return s.Actual }
func (s Uint) ScalarValue() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}
func (s Uint) ScalarSym() any                     { return s.Sym }
func (s Uint) ScalarDescription() string          { return s.Description }
func (s Uint) ScalarFlags() Flags                 { return s.Flags }
func (s Uint) ScalarDisplayFormat() DisplayFormat { return s.DisplayFormat }

func UintActual(v uint64) UintMapper {
	return UintFn(func(s Uint) (Uint, error) { s.Actual = v; return s, nil })
}
func UintSym(v any) UintMapper {
	return UintFn(func(s Uint) (Uint, error) { s.Sym = v; return s, nil })
}
func UintDescription(v string) UintMapper {
	return UintFn(func(s Uint) (Uint, error) { s.Description = v; return s, nil })
}

type UintMapper interface {
	MapUint(Uint) (Uint, error)
}

// UintFn map actual Uint using f
type UintFn func(s Uint) (Uint, error)

func (fn UintFn) MapUint(s Uint) (Uint, error) {
	return fn(s)
}

// UintActualFn map actual Uint using f
type UintActualFn func(a uint64) uint64

// TODO: error?
func (fn UintActualFn) MapUint(s Uint) (Uint, error) {
	s.Actual = fn(s.Actual)
	return s, nil
}

// UintSymFn map sym Uint using f
type UintSymFn func(a any) any

func (f UintSymFn) MapUint(s Uint) (Uint, error) {
	s.Sym = f(s.Sym)
	return s, nil
}

// UintDescriptionFn map sym Uint using f
type UintDescriptionFn func(a string) string

func (f UintDescriptionFn) MapUint(s Uint) (Uint, error) {
	s.Description = f(s.Description)
	return s, nil
}

// TrySymAny try assert symbolic value is a Any and return result
func (s Uint) TrySymAny() (any, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(any)
	return v, ok
}

// SymAny asserts symbolic value is a Any and returns it
func (s Uint) SymAny() any {
	v, ok := s.TrySymAny()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBigInt try assert symbolic value is a BigInt and return result
func (s Uint) TrySymBigInt() (*big.Int, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(*big.Int)
	return v, ok
}

// SymBigInt asserts symbolic value is a BigInt and returns it
func (s Uint) SymBigInt() *big.Int {
	v, ok := s.TrySymBigInt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBitBuf try assert symbolic value is a BitBuf and return result
func (s Uint) TrySymBitBuf() (bitio.ReaderAtSeeker, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bitio.ReaderAtSeeker)
	return v, ok
}

// SymBitBuf asserts symbolic value is a BitBuf and returns it
func (s Uint) SymBitBuf() bitio.ReaderAtSeeker {
	v, ok := s.TrySymBitBuf()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// TrySymBool try assert symbolic value is a Bool and return result
func (s Uint) TrySymBool() (bool, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(bool)
	return v, ok
}

// SymBool asserts symbolic value is a Bool and returns it
func (s Uint) SymBool() bool {
	v, ok := s.TrySymBool()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// TrySymFlt try assert symbolic value is a Flt and return result
func (s Uint) TrySymFlt() (float64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(float64)
	return v, ok
}

// SymFlt asserts symbolic value is a Flt and returns it
func (s Uint) SymFlt() float64 {
	v, ok := s.TrySymFlt()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// TrySymSint try assert symbolic value is a Sint and return result
func (s Uint) TrySymSint() (int64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(int64)
	return v, ok
}

// SymSint asserts symbolic value is a Sint and returns it
func (s Uint) SymSint() int64 {
	v, ok := s.TrySymSint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// TrySymStr try assert symbolic value is a Str and return result
func (s Uint) TrySymStr() (string, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(string)
	return v, ok
}

// SymStr asserts symbolic value is a Str and returns it
func (s Uint) SymStr() string {
	v, ok := s.TrySymStr()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// TrySymUint try assert symbolic value is a Uint and return result
func (s Uint) TrySymUint() (uint64, bool) {
	//nolint:gosimple,nolintlint
	v, ok := s.Sym.(uint64)
	return v, ok
}

// SymUint asserts symbolic value is a Uint and returns it
func (s Uint) SymUint() uint64 {
	v, ok := s.TrySymUint()
	if !ok {
		panic(fmt.Sprintf("failed to type assert s.Sym %v (%T) as uint64", s.Sym, s.Sym))
	}
	return v
}

// Map Bool
type BoolMap map[bool]Bool

func (m BoolMap) MapBool(s Bool) (Bool, error) {
	if ns, ok := m[s.Actual]; ok {
		ns.Actual = s.Actual
		return ns, nil
	}
	return s, nil
}

// Map Bool description
type BoolMapDescription map[bool]string

func (m BoolMapDescription) MapBool(s Bool) (Bool, error) {
	if d, ok := m[s.Actual]; ok {
		s.Description = d
	}
	return s, nil
}

// Map Bool sym Bool
type BoolMapSymBool map[bool]bool

func (m BoolMapSymBool) MapBool(s Bool) (Bool, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool sym Flt
type BoolMapSymFlt map[bool]float64

func (m BoolMapSymFlt) MapBool(s Bool) (Bool, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool sym Sint
type BoolMapSymSint map[bool]int64

func (m BoolMapSymSint) MapBool(s Bool) (Bool, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool sym Str
type BoolMapSymStr map[bool]string

func (m BoolMapSymStr) MapBool(s Bool) (Bool, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Bool sym Uint
type BoolMapSymUint map[bool]uint64

func (m BoolMapSymUint) MapBool(s Bool) (Bool, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Sint
type SintMap map[int64]Sint

func (m SintMap) MapSint(s Sint) (Sint, error) {
	if ns, ok := m[s.Actual]; ok {
		ns.Actual = s.Actual
		return ns, nil
	}
	return s, nil
}

// Map Sint description
type SintMapDescription map[int64]string

func (m SintMapDescription) MapSint(s Sint) (Sint, error) {
	if d, ok := m[s.Actual]; ok {
		s.Description = d
	}
	return s, nil
}

// Map Sint sym Bool
type SintMapSymBool map[int64]bool

func (m SintMapSymBool) MapSint(s Sint) (Sint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Sint sym Flt
type SintMapSymFlt map[int64]float64

func (m SintMapSymFlt) MapSint(s Sint) (Sint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Sint sym Sint
type SintMapSymSint map[int64]int64

func (m SintMapSymSint) MapSint(s Sint) (Sint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Sint sym Str
type SintMapSymStr map[int64]string

func (m SintMapSymStr) MapSint(s Sint) (Sint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Sint sym Uint
type SintMapSymUint map[int64]uint64

func (m SintMapSymUint) MapSint(s Sint) (Sint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str
type StrMap map[string]Str

func (m StrMap) MapStr(s Str) (Str, error) {
	if ns, ok := m[s.Actual]; ok {
		ns.Actual = s.Actual
		return ns, nil
	}
	return s, nil
}

// Map Str description
type StrMapDescription map[string]string

func (m StrMapDescription) MapStr(s Str) (Str, error) {
	if d, ok := m[s.Actual]; ok {
		s.Description = d
	}
	return s, nil
}

// Map Str sym Bool
type StrMapSymBool map[string]bool

func (m StrMapSymBool) MapStr(s Str) (Str, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str sym Flt
type StrMapSymFlt map[string]float64

func (m StrMapSymFlt) MapStr(s Str) (Str, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str sym Sint
type StrMapSymSint map[string]int64

func (m StrMapSymSint) MapStr(s Str) (Str, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str sym Str
type StrMapSymStr map[string]string

func (m StrMapSymStr) MapStr(s Str) (Str, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Str sym Uint
type StrMapSymUint map[string]uint64

func (m StrMapSymUint) MapStr(s Str) (Str, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Uint
type UintMap map[uint64]Uint

func (m UintMap) MapUint(s Uint) (Uint, error) {
	if ns, ok := m[s.Actual]; ok {
		ns.Actual = s.Actual
		return ns, nil
	}
	return s, nil
}

// Map Uint description
type UintMapDescription map[uint64]string

func (m UintMapDescription) MapUint(s Uint) (Uint, error) {
	if d, ok := m[s.Actual]; ok {
		s.Description = d
	}
	return s, nil
}

// Map Uint sym Bool
type UintMapSymBool map[uint64]bool

func (m UintMapSymBool) MapUint(s Uint) (Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Uint sym Flt
type UintMapSymFlt map[uint64]float64

func (m UintMapSymFlt) MapUint(s Uint) (Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Uint sym Sint
type UintMapSymSint map[uint64]int64

func (m UintMapSymSint) MapUint(s Uint) (Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Uint sym Str
type UintMapSymStr map[uint64]string

func (m UintMapSymStr) MapUint(s Uint) (Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}

// Map Uint sym Uint
type UintMapSymUint map[uint64]uint64

func (m UintMapSymUint) MapUint(s Uint) (Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t
	}
	return s, nil
}
