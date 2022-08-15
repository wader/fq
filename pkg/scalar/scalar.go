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

	"github.com/wader/fq/internal/bitioex"
	"github.com/wader/fq/pkg/bitio"
)

//go:generate sh -c "cat scalar_gen.go.tmpl | go run ../../dev/tmpl.go ../decode/types.json | gofmt > scalar_gen.go"

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

type S struct {
	Actual        any // nil, int, int64, uint64, float64, string, bool, []byte, *bit.Int, bitio.BitReaderAtSeeker,
	ActualDisplay DisplayFormat
	Sym           any
	SymDisplay    DisplayFormat
	Description   string
	Unknown       bool
}

func (s S) Value() any {
	if s.Sym != nil {
		return s.Sym
	}
	return s.Actual
}

type Mapper interface {
	MapScalar(S) (S, error)
}

type Fn func(S) (S, error)

func (fn Fn) MapScalar(s S) (S, error) {
	return fn(s)
}

var ActualBin = Fn(func(s S) (S, error) { s.ActualDisplay = NumberBinary; return s, nil })
var ActualOct = Fn(func(s S) (S, error) { s.ActualDisplay = NumberOctal; return s, nil })
var ActualDec = Fn(func(s S) (S, error) { s.ActualDisplay = NumberDecimal; return s, nil })
var ActualHex = Fn(func(s S) (S, error) { s.ActualDisplay = NumberHex; return s, nil })

var SymBin = Fn(func(s S) (S, error) { s.SymDisplay = NumberBinary; return s, nil })
var SymOct = Fn(func(s S) (S, error) { s.SymDisplay = NumberOctal; return s, nil })
var SymDec = Fn(func(s S) (S, error) { s.SymDisplay = NumberDecimal; return s, nil })
var SymHex = Fn(func(s S) (S, error) { s.SymDisplay = NumberHex; return s, nil })

func Actual(v any) Mapper {
	return Fn(func(s S) (S, error) { s.Actual = v; return s, nil })
}
func Sym(v any) Mapper {
	return Fn(func(s S) (S, error) { s.Sym = v; return s, nil })
}
func Description(v string) Mapper {
	return Fn(func(s S) (S, error) { s.Description = v; return s, nil })
}

func ActualUAdd(n int) ActualUFn {
	// TODO: use math.Add/Sub?
	return ActualUFn(func(a uint64) uint64 { return uint64(int64(a) + int64(n)) })
}

func ActualSAdd(n int) ActualSFn {
	return ActualSFn(func(a int64) int64 { return a + int64(n) })
}

func ActualTrim(cutset string) ActualStrFn {
	return ActualStrFn(func(a string) string { return strings.Trim(a, cutset) })
}

var ActualTrimSpace = ActualStrFn(strings.TrimSpace)

func strMapToSym(fn func(s string) (any, error), try bool) Mapper {
	return Fn(func(s S) (S, error) {
		ts := strings.TrimSpace(s.ActualStr())
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

func TrySymUParseUint(base int) Mapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseUint(s, base, 64) }, true)
}

func TrySymSParseInt(base int) Mapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseInt(s, base, 64) }, true)
}

func TrySymFParseFloat(base int) Mapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseFloat(s, base) }, true)
}

func SymUParseUint(base int) Mapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseUint(s, base, 64) }, false)
}

func SymSParseInt(base int) Mapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseInt(s, base, 64) }, false)
}

func SymFParseFloat(base int) Mapper {
	return strMapToSym(func(s string) (any, error) { return strconv.ParseFloat(s, base) }, false)
}

type URangeEntry struct {
	Range [2]uint64
	S     S
}

// URangeToScalar maps uint64 ranges to a scalar, first in range is chosen
type URangeToScalar []URangeEntry

func (rs URangeToScalar) MapScalar(s S) (S, error) {
	n := s.ActualU()
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
	S     S
}

// SRangeToScalar maps sint64 ranges to a scalar, first in range is chosen
type SRangeToScalar []SRangeEntry

func (rs SRangeToScalar) MapScalar(s S) (S, error) {
	n := s.ActualS()
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

func RawSym(s S, nBytes int, fn func(b []byte) string) (S, error) {
	br := s.ActualBitBuf()
	brLen, err := bitioex.Len(br)
	if err != nil {
		return S{}, err
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

var RawUUID = Fn(func(s S) (S, error) {
	return RawSym(s, -1, func(b []byte) string {
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	})
})

var RawHex = Fn(func(s S) (S, error) {
	return RawSym(s, -1, func(b []byte) string { return fmt.Sprintf("%x", b) })
})

type BytesToScalar []struct {
	Bytes  []byte
	Scalar S
}

func (m BytesToScalar) MapScalar(s S) (S, error) {
	rc, err := bitio.CloneReader(s.ActualBitBuf())
	if err != nil {
		return s, err
	}
	bb := &bytes.Buffer{}
	if _, err := bitioex.CopyBits(bb, rc); err != nil {
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

// TODO: nicer api, use generics

var unixTimeEpochDate = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

func DescriptionActualUTime(epoch time.Time, format string) Mapper {
	return Fn(func(s S) (S, error) {
		s.Description = epoch.Add(time.Second * time.Duration(s.ActualU())).Format(format)
		return s, nil
	})
}

func DescriptionSymUTime(epoch time.Time, format string) Mapper {
	return Fn(func(s S) (S, error) {
		s.Description = epoch.Add(time.Second * time.Duration(s.SymU())).Format(format)
		return s, nil
	})
}

var DescriptionActualUUnixTime = DescriptionActualUTime(unixTimeEpochDate, time.RFC3339)
var DescriptionSymUUnixTime = DescriptionSymUTime(unixTimeEpochDate, time.RFC3339)

func DescriptionActualSTime(epoch time.Time, format string) Mapper {
	return Fn(func(s S) (S, error) {
		s.Description = epoch.Add(time.Second * time.Duration(s.ActualS())).Format(format)
		return s, nil
	})
}

func DescriptionSymSTime(epoch time.Time, format string) Mapper {
	return Fn(func(s S) (S, error) {
		s.Description = epoch.Add(time.Second * time.Duration(s.SymS())).Format(format)
		return s, nil
	})
}

var DescriptionActualSUnixTime = DescriptionActualSTime(unixTimeEpochDate, time.RFC3339)
var DescriptionSymSUnixTime = DescriptionSymSTime(unixTimeEpochDate, time.RFC3339)
