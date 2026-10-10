package interp

import (
	"bytes"
	"fmt"

	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/pkg/bitio"
)

func init() {
	RegisterFunc1("_binary_indices", (*Interp)._binaryIndices)
	RegisterFunc1("_binary_starts_with", (*Interp)._binaryStartsWith)
	RegisterFunc1("_binary_ends_with", (*Interp)._binaryEndsWith)
	RegisterFunc1("_binary_ascii_case", (*Interp)._binaryASCIICase)
	RegisterFunc1("_binary_needle_length", (*Interp)._binaryNeedleLength)
}

// The string functions reach a binary through a coercion that replaces every
// byte it cannot read as UTF-8, so a needle and a haystack that differ only in
// those bytes compare equal. These work on the bytes themselves.
//
// A needle is a string, taken as its UTF-8 bytes, or a binary, taken as its
// bytes. Anything else is an error, which is what the string versions do.
func needleBytes(v any) ([]byte, error) {
	switch vv := v.(type) {
	case string:
		return []byte(vv), nil
	case ToBinary:
		return subjectBytes(vv)
	default:
		return nil, fmt.Errorf("a needle applied to a binary must be a string or a binary, got %v", v)
	}
}

// The bytes of a binary, read the way the overloads already here read them.
// A binary whose length is not a whole number of units carries padding, and
// toReader is what puts it in front, so going through it is what keeps these
// functions and match reporting the same offsets over the same subject.
func subjectBytes(c any) ([]byte, error) {
	bv, err := toBinary(c)
	if err != nil {
		return nil, err
	}
	br, err := bv.toReader()
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if _, err := bitiox.CopyBits(buf, br); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// every offset the needle begins at, counted in bytes and allowed to overlap,
// which is what the string version does
func (i *Interp) _binaryIndices(c any, needle any) any {
	hay, err := subjectBytes(c)
	if err != nil {
		return err
	}
	nee, err := needleBytes(needle)
	if err != nil {
		return err
	}
	vs := []any{}
	if len(nee) == 0 {
		return vs
	}
	for off := 0; off+len(nee) <= len(hay); off++ {
		if bytes.Equal(hay[off:off+len(nee)], nee) {
			vs = append(vs, off)
		}
	}
	return vs
}

func (i *Interp) _binaryStartsWith(c any, needle any) any {
	hay, err := subjectBytes(c)
	if err != nil {
		return err
	}
	nee, err := needleBytes(needle)
	if err != nil {
		return err
	}
	return bytes.HasPrefix(hay, nee)
}

func (i *Interp) _binaryEndsWith(c any, needle any) any {
	hay, err := subjectBytes(c)
	if err != nil {
		return err
	}
	nee, err := needleBytes(needle)
	if err != nil {
		return err
	}
	return bytes.HasSuffix(hay, nee)
}

// the length of a needle in bytes, so the jq side can slice a prefix or a
// suffix off without measuring it a second way
func (i *Interp) _binaryNeedleLength(c any, needle any) any {
	nee, err := needleBytes(needle)
	if err != nil {
		return err
	}
	return len(nee)
}

// Only the twenty six letters move. Every other byte is left alone, including
// the ones that are not ASCII at all, which is where the coercion used to turn
// one byte into the three of a replacement character.
func (i *Interp) _binaryASCIICase(c any, upper any) any {
	up, ok := upper.(bool)
	if !ok {
		return fmt.Errorf("ascii case wants a boolean")
	}
	b, err := subjectBytes(c)
	if err != nil {
		return err
	}
	o := make([]byte, len(b))
	for i, v := range b {
		switch {
		case up && v >= 'a' && v <= 'z':
			v -= 'a' - 'A'
		case !up && v >= 'A' && v <= 'Z':
			v += 'a' - 'A'
		}
		o[i] = v
	}
	bin, err := NewBinaryFromBitReader(bitio.NewBitReader(o, -1), 8, 0)
	if err != nil {
		return err
	}
	return bin
}
