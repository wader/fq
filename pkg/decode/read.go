package decode

import (
	"bytes"
	"fmt"
	"math"

	"github.com/wader/fq/pkg/bitio"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

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

func (d *D) tryFPE(nBits int, fBits int, endian Endian) (float64, error) {
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
//nolint:unparam
func (d *D) tryTextLenPrefixed(lenBits int, fixedBytes int, e encoding.Encoding) (string, error) {
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
