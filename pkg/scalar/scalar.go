package scalar

// TODO: d.Pos check err
// TODO: fn for actual?
// TODO: dsl import .? own scalar package?
// TODO: better IOError op names

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/pkg/bitio"
)

//go:generate sh -c "cat scalar_gen.go.tmpl | go run ../../dev/tmpl.go ../decode/types.json | gofmt > scalar_gen.go"

type Scalarable interface {
	ScalarActual() any
	ScalarValue() any
	ScalarSym() any
	ScalarDescription() string
	ScalarFlags() Flags
	ScalarDisplayFormat() DisplayFormat
}

type DisplayFormat int

const (
	NumberDecimal DisplayFormat = iota
	NumberBinary
	NumberOctal
	NumberHex
)

func (df DisplayFormat) FormatBase() int {
	switch df {
	case NumberDecimal:
		return 10
	case NumberBinary:
		return 2
	case NumberOctal:
		return 8
	case NumberHex:
		return 16
	default:
		return 0
	}
}

const (
	FlagGap Flags = 1 << iota
	FlagSynthetic
)

type Flags uint

func (f Flags) IsGap() bool       { return f&FlagGap != 0 }
func (f Flags) IsSynthetic() bool { return f&FlagSynthetic != 0 }

// TODO: todos
// rename raw?
// crc
//
//
// d.FieldU2("emphasis").
//   MapSym().
//   Actual
//   Sym
//   Value()
//   Scalar()

var UintBin = UintFn(func(s Uint) (Uint, error) { s.DisplayFormat = NumberBinary; return s, nil })
var UintOct = UintFn(func(s Uint) (Uint, error) { s.DisplayFormat = NumberOctal; return s, nil })
var UintDec = UintFn(func(s Uint) (Uint, error) { s.DisplayFormat = NumberDecimal; return s, nil })
var UintHex = UintFn(func(s Uint) (Uint, error) { s.DisplayFormat = NumberHex; return s, nil })
var SintBin = SintFn(func(s Sint) (Sint, error) { s.DisplayFormat = NumberBinary; return s, nil })
var SintOct = SintFn(func(s Sint) (Sint, error) { s.DisplayFormat = NumberOctal; return s, nil })
var SintDec = SintFn(func(s Sint) (Sint, error) { s.DisplayFormat = NumberDecimal; return s, nil })
var SintHex = SintFn(func(s Sint) (Sint, error) { s.DisplayFormat = NumberHex; return s, nil })
var BigIntBin = BigIntFn(func(s BigInt) (BigInt, error) { s.DisplayFormat = NumberBinary; return s, nil })
var BigIntOct = BigIntFn(func(s BigInt) (BigInt, error) { s.DisplayFormat = NumberOctal; return s, nil })
var BigIntDec = BigIntFn(func(s BigInt) (BigInt, error) { s.DisplayFormat = NumberDecimal; return s, nil })
var BigIntHex = BigIntFn(func(s BigInt) (BigInt, error) { s.DisplayFormat = NumberHex; return s, nil })

func UintActualAdd(n int) UintActualFn {
	// TODO: use math.Add/Sub?
	return UintActualFn(func(a uint64) uint64 { return uint64(int64(a) + int64(n)) })
}

func SintActualAdd(n int) SintActualFn {
	return SintActualFn(func(a int64) int64 { return a + int64(n) })
}

func StrActualTrim(cutset string) StrActualFn {
	return StrActualFn(func(a string) string { return strings.Trim(a, cutset) })
}

var ActualTrimSpace = StrActualFn(strings.TrimSpace)

func strMapToSym(fn func(s string) (any, error), try bool) StrMapper {
	return StrFn(func(s Str) (Str, error) {
		ts := strings.TrimSpace(s.Actual)
		n, err := fn(ts)
		if err != nil {
			if try {
				return s, nil
			}
			return s, err
		}
		s.Sym = n
		return s, nil
	})
}

func TryStrSymParseUint(base int) StrMapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseUint(s, base, 64) }, true)
}

func TryStrSymParseInt(base int) StrMapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseInt(s, base, 64) }, true)
}

func TryStrSymParseFloat(base int) StrMapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseFloat(s, base) }, true)
}

func StrSymParseUint(base int) StrMapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseUint(s, base, 64) }, false)
}

func StrSymParseInt(base int) StrMapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseInt(s, base, 64) }, false)
}

func StrSymParseFloat(base int) StrMapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseFloat(s, base) }, false)
}

type URangeEntry struct {
	Range [2]uint64
	S     Uint
}

// UintRangeToScalar maps uint64 ranges to a scalar, first in range is chosen
type UintRangeToScalar []URangeEntry

func (rs UintRangeToScalar) MapUint(s Uint) (Uint, error) {
	n := s.Actual
	for _, re := range rs {
		if n >= re.Range[0] && n <= re.Range[1] {
			ns := re.S
			ns.Actual = s.Actual
			s = ns
			break
		}
	}
	return s, nil
}

// SRangeToScalar maps ranges to a scalar, first in range is chosen
type SRangeEntry struct {
	Range [2]int64
	S     Sint
}

// SRangeToScalar maps sint64 ranges to a scalar, first in range is chosen
type SRangeToScalar []SRangeEntry

func (rs SRangeToScalar) MapSint(s Sint) (Sint, error) {
	n := s.Actual
	for _, re := range rs {
		if n >= re.Range[0] && n <= re.Range[1] {
			ns := re.S
			ns.Actual = s.Actual
			s = ns
			break
		}
	}
	return s, nil
}

func RawSym(s BitBuf, nBytes int, fn func(b []byte) string) (BitBuf, error) {
	br := s.Actual
	brLen, err := bitiox.Len(br)
	if err != nil {
		return BitBuf{}, err
	}
	if nBytes < 0 {
		nBytes = int(brLen) / 8
		if brLen%8 != 0 {
			nBytes++
		}
	}
	if brLen < int64(nBytes)*8 {
		return s, nil
	}
	// TODO: shared somehow?
	b := make([]byte, nBytes)
	if _, err := br.ReadBitsAt(b, int64(nBytes)*8, 0); err != nil {
		return s, err
	}

	s.Sym = fn(b[0:nBytes])

	return s, nil
}

var RawUUID = BitBufFn(func(s BitBuf) (BitBuf, error) {
	return RawSym(s, -1, func(b []byte) string {
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	})
})

var RawHex = BitBufFn(func(s BitBuf) (BitBuf, error) {
	return RawSym(s, -1, func(b []byte) string { return fmt.Sprintf("%x", b) })
})

type RawBytesMap []struct {
	Bytes  []byte
	Scalar BitBuf
}

func (m RawBytesMap) MapBitBuf(s BitBuf) (BitBuf, error) {
	rc, err := bitio.CloneReader(s.Actual)
	if err != nil {
		return s, err
	}
	bb := &bytes.Buffer{}
	if _, err := bitiox.CopyBits(bb, rc); err != nil {
		return s, err
	}
	for _, bs := range m {
		if bytes.Equal(bb.Bytes(), bs.Bytes) {
			ns := bs.Scalar
			ns.Actual = s.Actual
			s = ns
			break
		}
	}
	return s, nil
}

var unixTimeEpochDate = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

func UintActualDateDescription(epoch time.Time, unit time.Duration, format string) UintFn {
	return UintFn(func(s Uint) (Uint, error) {
		s.Description = epoch.Add(time.Duration(s.Actual) * unit).Format(format)
		return s, nil
	})
}

func UintActualUnixTimeDescription(unit time.Duration, format string) UintFn {
	return UintActualDateDescription(unixTimeEpochDate, unit, format)
}

func SintActualDateDescription(epoch time.Time, unit time.Duration, format string) SintFn {
	return SintFn(func(s Sint) (Sint, error) {
		s.Description = epoch.Add(time.Duration(s.Actual) * unit).Format(format)
		return s, nil
	})
}

func SintActualUnixTimeDescription(unit time.Duration, format string) SintFn {
	return SintActualDateDescription(unixTimeEpochDate, unit, format)
}

func FltActualDateDescription(epoch time.Time, unit time.Duration, format string) FltFn {
	return FltFn(func(s Flt) (Flt, error) {
		s.Description = epoch.Add(time.Duration(s.Actual) * unit).Format(format)
		return s, nil
	})
}

func FltActualUnixTimeDescription(unit time.Duration, format string) FltFn {
	return FltActualDateDescription(unixTimeEpochDate, unit, format)
}
