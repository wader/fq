package decode

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/scalar"
)

func bitBufIsZero(s scalar.S, isValidate bool) (scalar.S, error) {
	bb, ok := s.Actual.(*bitio.Buffer)
	if !ok {
		return s, nil
	}

	isZero := true
	// TODO: shared
	b := make([]byte, 32*1024)
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
		if isValidate {
			return s, errors.New("validate is zero failed")
		}
	}

	return s, nil
}

func (d *D) BitBufIsZero() scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return bitBufIsZero(s, false)
	})
}

func (d *D) BitBufValidateIsZero() scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return bitBufIsZero(s, !d.Options.Force)
	})
}

// TODO: generate?
func assertBitBuf(s scalar.S, isErr bool, bss ...[]byte) (scalar.S, error) {
	ab, err := s.ActualBitBuf().Bytes()
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
	if isErr {
		return s, errors.New("failed to validate raw")
	}
	return s, nil
}

func (d *D) AssertBitBuf(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertBitBuf(s, !d.Options.Force, bss...)
	})
}

func (d *D) ValidateBitBuf(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertBitBuf(s, false, bss...)
	})
}

func assertUBytes(s scalar.S, isErr bool, endian Endian, bss ...[]byte) (scalar.S, error) {
	var bo binary.ByteOrder
	switch endian {
	case LittleEndian:
		bo = binary.BigEndian
	case BigEndian:
		bo = binary.BigEndian
	default:
		panic("invalid endian")
	}

	au := s.ActualU()
	for _, bs := range bss {
		var bu uint64
		switch len(bs) {
		case 1:
			bu = uint64(bs[0])
		case 2:
			bu = uint64(bo.Uint16(bs))
		case 4:
			bu = uint64(bo.Uint32(bs))
		case 8:
			bu = bo.Uint64(bs)
		default:
			panic("invalid bs length")
		}

		if au == bu {
			s.Description = "valid"
			return s, nil
		}
	}
	s.Description = "invalid"
	if isErr {
		return s, errors.New("failed to validate raw")
	}
	return s, nil
}

func (d *D) AssertUBytes(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertUBytes(s, true, d.Endian, bss...)
	})
}
func (d *D) ValidateUBytes(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertUBytes(s, false, d.Endian, bss...)
	})
}
func (d *D) AssertULEBytes(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertUBytes(s, true, LittleEndian, bss...)
	})
}
func (d *D) ValidateULEBytes(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertUBytes(s, false, LittleEndian, bss...)
	})
}
func (d *D) AssertUBEBytes(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertUBytes(s, true, BigEndian, bss...)
	})
}
func (d *D) ValidateUBEBytes(bss ...[]byte) scalar.Mapper {
	return scalar.Fn(func(s scalar.S) (scalar.S, error) {
		return assertUBytes(s, false, BigEndian, bss...)
	})
}
