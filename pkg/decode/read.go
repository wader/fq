package decode

import (
	"bytes"
	"fmt"
	"math"
	"math/big"

	"github.com/wader/fq/internal/mathextra"
	"github.com/wader/fq/pkg/bitio"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

func (d *D) tryBitBuf(nBits int64) (*bitio.Buffer, error) {
	return d.bitBuf.BitBufLen(nBits)
}

func (d *D) tryUEndian(nBits int, endian Endian) (uint64, error) {
	if nBits < 0 {
		return 0, fmt.Errorf("tryUEndian nBits must be >= 0 (%d)", nBits)
	}
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}

	return n, nil
}

func (d *D) trySEndian(nBits int, endian Endian) (int64, error) {
	if nBits < 0 {
		return 0, fmt.Errorf("trySEndian nBits must be >= 0 (%d)", nBits)
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
func reverseBytes(a []byte) {
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
	_, err := bitio.ReadFull(d.bitBuf, buf, nBits)
	if err != nil {
		return nil, err
	}

	if endian == LittleEndian {
		reverseBytes(buf)
	}

	n := new(big.Int)
	if sign {
		mathextra.BigIntSetBytesSigned(n, buf)
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
	n, err := d.bits(nBits)
	if err != nil {
		return 0, err
	}
	if endian == LittleEndian {
		n = bitio.Uint64ReverseBytes(nBits, n)
	}
	switch nBits {
	case 16:
		return float64(mathextra.Float16(uint16(n)).Float32()), nil
	case 32:
		return float64(math.Float32frombits(uint32(n))), nil
	case 64:
		return math.Float64frombits(n), nil
	default:
		return 0, fmt.Errorf("unsupported float size %d", nBits)
	}
}

func (d *D) tryFPEndian(nBits int, fBits int, endian Endian) (float64, error) {
	if nBits < 0 {
		return 0, fmt.Errorf("tryFPEndian nBits must be >= 0 (%d)", nBits)
	}
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
	if nBytes < 0 {
		return "", fmt.Errorf("tryText nBytes must be >= 0 (%d)", nBytes)
	}
	bytesLeft := d.BitsLeft() / 8
	if int64(nBytes) > bytesLeft {
		return "", fmt.Errorf("tryText nBytes %d outside buffer, %d bytes left", nBytes, bytesLeft)
	}

	bs, err := d.bitBuf.BytesLen(nBytes)
	if err != nil {
		return "", err
	}
	return e.NewDecoder().String(string(bs))
}

// read length prefixed text (ex pascal short string)
// lBits length prefix
// fixedBytes if != -1 read nBytes but trim to length
//nolint:unparam
func (d *D) tryTextLenPrefixed(lenBits int, fixedBytes int, e encoding.Encoding) (string, error) {
	if lenBits < 0 {
		return "", fmt.Errorf("tryTextLenPrefixed lenBits must be >= 0 (%d)", lenBits)
	}
	if fixedBytes < 0 {
		return "", fmt.Errorf("tryTextLenPrefixed fixedBytes must be >= 0 (%d)", fixedBytes)
	}
	bytesLeft := d.BitsLeft() / 8
	if int64(fixedBytes) > bytesLeft {
		return "", fmt.Errorf("tryTextLenPrefixed fixedBytes %d outside, %d bytes left", fixedBytes, bytesLeft)
	}

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

func (d *D) tryTextNull(nullBytes int, e encoding.Encoding) (string, error) {
	if nullBytes < 1 {
		return "", fmt.Errorf("tryTextNull nullBytes must be >= 1 (%d)", nullBytes)
	}

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

func (d *D) tryTextNullLen(fixedBytes int, e encoding.Encoding) (string, error) {
	if fixedBytes < 0 {
		return "", fmt.Errorf("tryTextNullLen fixedBytes must be >= 0 (%d)", fixedBytes)
	}
	bytesLeft := d.BitsLeft() / 8
	if int64(fixedBytes) > bytesLeft {
		return "", fmt.Errorf("tryTextNullLen fixedBytes %d outside, %d bytes left", fixedBytes, bytesLeft)
	}

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
