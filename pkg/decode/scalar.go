package decode

// TODO: d.Pos check err
// TODO: fn for actual?
// TODO: dsl import .? own scalar package?
// TODO: better IOError op names

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/ranges"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

//go:generate sh -c "cat scalar_gen.go.tmpl | go run ../../dev/tmpl.go scalar_gen.go.json | gofmt > scalar_gen.go"

type ScalarFn func(Scalar) (Scalar, error)

func (d *D) Bin(s Scalar) (Scalar, error) { s.DisplayFormat = NumberBinary; return s, nil }
func (d *D) Oct(s Scalar) (Scalar, error) { s.DisplayFormat = NumberOctal; return s, nil }
func (d *D) Dec(s Scalar) (Scalar, error) { s.DisplayFormat = NumberDecimal; return s, nil }
func (d *D) Hex(s Scalar) (Scalar, error) { s.DisplayFormat = NumberHex; return s, nil }

func (d *D) Actual(a interface{}) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) { s.Actual = a; return s, nil }
}

func (d *D) Sym(sym interface{}) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) { s.Sym = sym; return s, nil }
}

func (d *D) Description(desc string) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) { s.Description = desc; return s, nil }
}

func (d *D) UAdd(n int) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if v, ok := s.Actual.(uint64); ok {
			// TODO: use math.Add/Sub?
			s.Actual = uint64(int64(v) + int64(n))
		}
		return s, nil
	}
}

// TODO: nicer api?
func (d *D) Trim(cutset string) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		if v, ok := s.Actual.(string); ok {
			s.Actual = strings.Trim(v, cutset)
		}
		return s, nil
	}
}

func (d *D) TrimSpace(s Scalar) (Scalar, error) {
	if v, ok := s.Actual.(string); ok {
		s.Actual = strings.TrimSpace(v)
	}
	return s, nil
}

func (d *D) RawSym(s Scalar, nBytes int, fn func(b []byte) string) (Scalar, error) {
	bb, ok := s.Actual.(*bitio.Buffer)
	if !ok {
		return s, nil
	}
	bbLen := bb.Len()
	if nBytes < 0 {
		nBytes = int(bbLen) / 8
		if bbLen%8 != 0 {
			nBytes++
		}
	}
	if bbLen < int64(nBytes)*8 {
		return s, nil
	}
	b := d.SharedReadBuf(nBytes)
	if _, err := bb.ReadBitsAt(b, nBytes*8, 0); err != nil {
		return s, err
	}

	s.Sym = fn(b[0:nBytes])

	return s, nil
}

func (d *D) RawUUID(s Scalar) (Scalar, error) {
	const uuidLen = 16
	return d.RawSym(s, uuidLen, func(b []byte) string {
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	})
}

func (d *D) RawHex(s Scalar) (Scalar, error) {
	return d.RawSym(s, -1, func(b []byte) string { return fmt.Sprintf("%x", b) })
}

func (d *D) RawHexReverse(s Scalar) (Scalar, error) {
	return d.RawSym(s, -1, func(b []byte) string {
		return fmt.Sprintf("%x", bitio.ReverseBytes(append([]byte{}, b...)))
	})
}

func (d *D) MapURangeToScalar(rm map[[2]uint64]Scalar) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		n, ok := s.Actual.(uint64)
		if !ok {
			return s, nil
		}
		for r, rs := range rm {
			if n >= r[0] && n <= r[1] {
				ns := rs
				ns.Actual = s.Actual
				s = ns
				break
			}
		}
		return s, nil
	}
}

func (d *D) MapSRangeToScalar(rm map[[2]int64]Scalar) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		n, ok := s.Actual.(int64)
		if !ok {
			return s, nil
		}
		for r, rs := range rm {
			if n >= r[0] && n <= r[1] {
				ns := rs
				ns.Actual = s.Actual
				s = ns
				break
			}
		}
		return s, nil
	}
}

type BytesToScalar []struct {
	Bytes  []byte
	Scalar Scalar
}

func (d *D) MapRawToScalar(btss BytesToScalar) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		// TODO: check type assert?
		ab, err := s.Actual.(*bitio.Buffer).Bytes()
		if err != nil {
			return s, err
		}
		for _, bs := range btss {
			if bytes.Equal(ab, bs.Bytes) {
				ns := bs.Scalar
				ns.Actual = s.Actual
				break
			}
		}
		return s, nil
	}
}

func (d *D) bitBufIsZero(s Scalar, err bool) (Scalar, error) {
	bb, ok := s.Actual.(*bitio.Buffer)
	if !ok {
		return s, nil
	}

	isZero := true
	b := d.SharedReadBuf(32 * 1024)
	bLen := len(b) * 8
	bbLeft := int(bb.Len())
	bbPos := int64(0)

	for bbLeft > 0 {
		rl := bbLeft
		if bbLeft > bLen {
			rl = bLen
		}
		// zero last byte if uneven read
		if rl%8 != 0 {
			b[rl/8] = 0
		}

		n, err := bitio.ReadAtFull(bb, b, rl, bbPos)
		if err != nil {
			return s, err
		}
		nb := int(bitio.BitsByteCount(int64(n)))

		for i := 0; i < nb; i++ {
			if b[i] != 0 {
				isZero = false
				break
			}
		}

		bbLeft -= n
	}

	if isZero {
		s.Description = "all zero"
	} else {
		s.Description = "all not zero"
		if err {
			return s, errors.New("validate is zero failed")
		}
	}

	return s, nil
}

func (d *D) BitBufIsZero(s Scalar) (Scalar, error) {
	return d.bitBufIsZero(s, false)
}

func (d *D) BitBufValidateIsZero(s Scalar) (Scalar, error) {
	return d.bitBufIsZero(s, true)
}

// TODO: generate?
func (d *D) assertRaw(assert bool, bss ...[]byte) func(s Scalar) (Scalar, error) {
	return func(s Scalar) (Scalar, error) {
		// TODO: check type assert?
		ab, err := s.Actual.(*bitio.Buffer).Bytes()
		if err != nil {
			return s, err
		}
		for _, bs := range bss {
			if bytes.Equal(ab, bs) {
				s.Description = "valid"
				return s, nil
			}
		}
		s.Description = "invalid"
		if assert {
			return s, errors.New("failed to validate raw")
		}
		return s, nil
	}
}

func (d *D) AssertRaw(bss ...[]byte) func(s Scalar) (Scalar, error) {
	return d.assertRaw(true, bss...)
}
func (d *D) ValidateRaw(bss ...[]byte) func(s Scalar) (Scalar, error) {
	return d.assertRaw(false, bss...)
}

func (d *D) TryFieldValue(name string, fn func() (*Value, error)) (*Value, error) {
	start := d.Pos()
	v, err := fn()
	stop := d.Pos()
	v.Name = name
	v.RootBitBuf = d.bitBuf
	v.Range = ranges.Range{Start: start, Len: stop - start}
	if err != nil {
		return nil, err
	}
	d.AddChild(v)

	return v, err
}

func (d *D) FieldValue(name string, fn func() *Value) *Value {
	v, err := d.TryFieldValue(name, func() (*Value, error) { return fn(), nil })
	if err != nil {
		panic(err)
	}
	return v
}

// looks a bit weird to force at least one ScalarFn arg
func (d *D) TryFieldScalar(name string, sfn ScalarFn, sfns ...ScalarFn) (Scalar, error) {
	v, err := d.TryFieldValue(name, func() (*Value, error) {
		s, err := sfn(Scalar{})
		if err != nil {
			return &Value{V: s}, err
		}
		for _, sfn := range sfns {
			s, err = sfn(s)
			if err != nil {
				return &Value{V: s}, err
			}
		}
		return &Value{V: s}, nil
	})
	if err != nil {
		return Scalar{}, err
	}
	return v.V.(Scalar), nil
}

func (d *D) FieldScalar(name string, sfn ScalarFn, sfns ...ScalarFn) Scalar {
	v, err := d.TryFieldScalar(name, sfn, sfns...)
	if err != nil {
		panic(err)
	}
	return v
}

func (d *D) tryUE(nBits int, endian Endian) (uint64, error) {
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}

	return n, nil
}

func (d *D) tryBitBuf(nBits int64) (*bitio.Buffer, error) {
	return d.bitBuf.BitBufLen(nBits)
}

func (d *D) trySE(nBits int, endian Endian) (int64, error) {
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	if nBits == 0 {
		return 0, nil
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}
	var s int64
	if n&(1<<(nBits-1)) > 0 {
		// two's complement
		s = -int64((^n & ((1 << nBits) - 1)) + 1)
	} else {
		s = int64(n)
	}

	return s, nil
}

func (d *D) tryFE(nBits int, endian Endian) (float64, error) {
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}
	switch nBits {
	case 32:
		return float64(math.Float32frombits(uint32(n))), nil
	case 64:
		return math.Float64frombits(n), nil
	default:
		return 0, fmt.Errorf("unsupported float size %d", nBits)
	}
}

func (d *D) tryFPE(nBits int, fBits int64, endian Endian) (float64, error) {
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}
	return float64(n) / float64(uint64(1<<fBits)), nil
}

var UTF8BOM = unicode.UTF8BOM
var UTF16BOM = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
var UTF16BE = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
var UTF16LE = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

func (d *D) tryText(nBytes int, e encoding.Encoding) (string, error) {
	bs, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		return "", err
	}
	return e.NewDecoder().String(string(bs))
}

// read length prefixed text (ex pascal short string)
// lBits length prefix
// fixedBytes if != -1 read nBytes but trim to length
func (d *D) tryLenPrefixedText(lenBits int, fixedBytes int, e encoding.Encoding) (string, error) {
	p := d.Pos()
	l, err := d.bits(lenBits)
	if err != nil {
		return "", err
	}

	n := int(l)
	if fixedBytes != -1 {
		n = fixedBytes - 1
		// TODO: error?
		if l > uint64(n) {
			l = uint64(n)
		}
	}

	bs, err := d.bitBuf.BytesLen(n)
	if err != nil {
		d.SeekAbs(p)
		return "", err
	}
	return e.NewDecoder().String(string(bs[0:l]))
}

func (d *D) tryNullTerminatedText(nullBytes int, e encoding.Encoding) (string, error) {
	p := d.Pos()
	peekBits, _, err := d.TryPeekFind(nullBytes*8, 8, -1, func(v uint64) bool { return v == 0 })
	if err != nil {
		return "", err
	}
	n := (int(peekBits) / 8) + nullBytes
	bs, err := d.bitBuf.BytesLen(n)
	if err != nil {
		d.SeekAbs(p)
		return "", err
	}

	return e.NewDecoder().String(string(bs[0 : n-nullBytes]))
}

func (d *D) tryNullTerminatedLenText(fixedBytes int, e encoding.Encoding) (string, error) {
	bs, err := d.bitBuf.BytesLen(fixedBytes)
	if err != nil {
		return "", err
	}
	nullIndex := bytes.IndexByte(bs, 0)
	if nullIndex != -1 {
		bs = bs[:nullIndex]
	}

	return e.NewDecoder().String(string(bs))
}

// ov is what to treat as 1
func (d *D) tryUnary(ov uint64) (uint64, error) {
	p := d.Pos()
	var n uint64
	for {
		b, err := d.bits(1)
		if err != nil {
			d.SeekAbs(p)
			return 0, err
		}
		if b != ov {
			break
		}
		n++
	}
	return n, nil
}

func (d *D) tryBool() (bool, error) {
	n, err := d.bits(1)
	if err != nil {
		return false, err
	}
	return n == 1, nil
}
