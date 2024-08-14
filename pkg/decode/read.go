package decode

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"github.com/wader/fq/internal/mathx"
	"github.com/wader/fq/pkg/bitio"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

func (d *D) tryBitBuf(nBits int64) (bitio.ReaderAtSeeker, error) {
	return d.TryBitBufLen(nBits)
}

func (d *D) tryUEndian(nBits int, endian Endian) (uint64, error) {
	if nBits < 0 {
		return 0, fmt.Errorf("tryUEndian nBits must be >= 0 (%d)", nBits)
	}
	n, err := d.TryUintBits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.ReverseBytes64(nBits, n)
	}

	return n, nil
}

func (d *D) trySEndian(nBits int, endian Endian) (int64, error) {
	if nBits < 1 {
		return 0, fmt.Errorf("trySEndian nBits must be >= 1 (%d)", nBits)
	}
	n, err := d.tryUEndian(nBits, endian)
	if err != nil {
		return 0, err
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

// from https://github.com/golang/go/wiki/SliceTricks#reversing
func ReverseBytes(a []byte) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

func (d *D) tryBigIntEndianSign(nBits int, endian Endian, sign bool) (*big.Int, error) {
	if nBits < 0 {
		return nil, fmt.Errorf("tryBigIntEndianSign nBits must be >= 0 (%d)", nBits)
	}
	b := int(bitio.BitsByteCount(int64(nBits)))
	buf := d.SharedReadBuf(b)[0:b]
	_, err := bitio.ReadFull(d.bitBuf, buf, int64(nBits))
	if err != nil {
		return nil, err
	}

	if endian == LittleEndian {
		ReverseBytes(buf)
	}

	n := new(big.Int)
	if sign {
		mathx.BigIntSetBytesSigned(n, buf)
	} else {
		n.SetBytes(buf)
	}
	n.Rsh(n, uint((8-nBits%8)%8))

	return n, nil
}

func (d *D) tryFEndian(nBits int, endian Endian) (float64, error) {
	if nBits < 0 {
		return 0, fmt.Errorf("tryFEndian nBits must be >= 0 (%d)", nBits)
	}
	b, err := d.TryBits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		ReverseBytes(b)
	}
	switch nBits {
	case 16:
		return float64(mathx.Float16(binary.BigEndian.Uint16(b)).Float32()), nil
	case 32:
		return float64(math.Float32frombits(binary.BigEndian.Uint32(b))), nil
	case 64:
		return math.Float64frombits(binary.BigEndian.Uint64(b)), nil
	case 80:
		return mathx.NewFloat80FromBytes(b).Float64(), nil
	default:
		return 0, fmt.Errorf("unsupported float size %d", nBits)
	}
}

func (d *D) tryFPEndian(nBits int, fBits int, endian Endian) (float64, error) {
	if nBits < 0 {
		return 0, fmt.Errorf("tryFPEndian nBits must be >= 0 (%d)", nBits)
	}
	n, err := d.TryUintBits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.ReverseBytes64(nBits, n)
	}
	return float64(n) / float64(uint64(1<<fBits)), nil
}

var UTF8BOM = unicode.UTF8BOM
var UTF16BOM = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
var UTF16BE = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
var UTF16LE = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)

func (d *D) tryText(nBytes int, e encoding.Encoding) (string, error) {
	if nBytes < 0 {
		return "", fmt.Errorf("tryText nBytes must be >= 0 (%d)", nBytes)
	}
	bytesLeft := d.BitsLeft() / 8
	if int64(nBytes) > bytesLeft {
		return "", fmt.Errorf("tryText nBytes %d outside buffer, %d bytes left", nBytes, bytesLeft)
	}

	bs, err := d.TryBytesLen(nBytes)
	if err != nil {
		return "", err
	}
	return e.NewDecoder().String(string(bs))
}

// read length prefixed text (ex pascal short string)
// lenBytes length prefix
// fixedBytes if != -1 read nBytes but trim to length
//
//nolint:unparam
func (d *D) tryTextLenPrefixed(prefixLenBytes int, fixedBytes int, e encoding.Encoding) (string, error) {
	if prefixLenBytes < 0 {
		return "", fmt.Errorf("tryTextLenPrefixed lenBytes must be >= 0 (%d)", prefixLenBytes)
	}
	bytesLeft := d.BitsLeft() / 8
	if int64(fixedBytes) > bytesLeft {
		return "", fmt.Errorf("tryTextLenPrefixed fixedBytes %d outside, %d bytes left", fixedBytes, bytesLeft)
	}

	p := d.Pos()
	lenBytes, err := d.TryUintBits(prefixLenBytes * 8)
	if err != nil {
		return "", err
	}

	readBytes := int(lenBytes)
	if fixedBytes != -1 {
		// TODO: error?
		readBytes = fixedBytes - prefixLenBytes
		lenBytes = min(lenBytes, uint64(readBytes))
	}

	bs, err := d.TryBytesLen(readBytes)
	if err != nil {
		d.SeekAbs(p)
		return "", err
	}
	return e.NewDecoder().String(string(bs[0:lenBytes]))
}

func (d *D) tryTextNull(charBytes int, e encoding.Encoding) (string, error) {
	if charBytes < 1 {
		return "", fmt.Errorf("tryTextNull charBytes must be >= 1 (%d)", charBytes)
	}

	p := d.Pos()
	peekBits, _, err := d.TryPeekFind(charBytes*8, int64(charBytes)*8, -1, func(v uint64) bool { return v == 0 })
	if err != nil {
		return "", err
	}
	n := (int(peekBits) / 8) + charBytes
	bs, err := d.TryBytesLen(n)
	if err != nil {
		d.SeekAbs(p)
		return "", err
	}

	return e.NewDecoder().String(string(bs[0 : n-charBytes]))
}

func (d *D) tryTextNullLen(fixedBytes int, e encoding.Encoding) (string, error) {
	if fixedBytes < 0 {
		return "", fmt.Errorf("tryTextNullLen fixedBytes must be >= 0 (%d)", fixedBytes)
	}
	bytesLeft := d.BitsLeft() / 8
	if int64(fixedBytes) > bytesLeft {
		return "", fmt.Errorf("tryTextNullLen fixedBytes %d outside, %d bytes left", fixedBytes, bytesLeft)
	}

	bs, err := d.TryBytesLen(fixedBytes)
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
		b, err := d.TryUintBits(1)
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
	n, err := d.TryUintBits(1)
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

// Unsigned LEB128, also known as "Base 128 Varint".
//
// Description from wasm spec:
//
//	uN ::= n:byte          => n                     (if n < 2^7 && n < 2^N)
//	       n:byte m:u(N-7) => 2^7 * m + (n - 2^7)   (if n >= 2^7 && N > 7)
//
// Varint description:
// https://protobuf.dev/programming-guides/encoding/#varints
func (d *D) tryULEB128() (uint64, error) {
	var result uint64
	var shift uint

	for {
		b := d.U8()
		if shift >= 63 && b != 0 {
			return 0, fmt.Errorf("overflow when reading unsigned leb128, shift %d >= 63", shift)
		}
		result |= (b & 0b01111111) << shift
		if b&0b10000000 == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}

// Signed LEB128, description from wasm spec
//
//	sN ::= n:byte          => n                     (if n < 2^6 && n < 2^(N-1))
//	       n:byte          => n - 2^7               (if 2^6 <= n < 2^7 && n >= 2^7 - 2^(N-1))
//	       n:byte m:s(N-7) => 2^7 * m + (n - 2^7)   (if n >= 2^7 && N > 7)
func (d *D) trySLEB128() (int64, error) {
	const n = 64
	var result int64
	var shift uint
	var b byte

	for {
		b = byte(d.U8())
		if shift == 63 && b != 0 && b != 0x7f {
			return 0, fmt.Errorf("overflow when reading signed leb128, shift %d >= 63", shift)
		}

		result |= int64(b&0x7f) << shift
		shift += 7

		if b&0x80 == 0 {
			break
		}
	}

	if shift < n && (b&0x40) == 0x40 {
		result |= -1 << shift
	}

	return result, nil
}
