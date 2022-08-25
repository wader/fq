package decode

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/wader/fq/internal/bitioex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/scalar"
)

func bitBufIsZero(s scalar.S, isValidate bool) (scalar.S, error) {
	br := s.ActualBitBuf()

	isZero := true
	// TODO: shared
	b := make([]byte, 32*1024)
	bLen := int64(len(b)) * 8
	brLen, err := bitioex.Len(br)
	if err != nil {
		return scalar.S{}, err
	}
	brLeft := brLen
	brPos := int64(0)

	for brLeft > 0 {
		rl := brLeft
		if brLeft > bLen {
			rl = bLen
		}
		// zero last byte if uneven read
		if rl%8 != 0 {
			b[rl/8] = 0
		}

		n, err := bitio.ReadAtFull(br, b, rl, brPos)
		if err != nil {
			return s, err
		}
		nb := int(bitio.BitsByteCount(n))

		for i := 0; i < nb; i++ {
			if b[i] != 0 {
				isZero = false
				break
			}
		}

		brLeft -= n
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
	bb := &bytes.Buffer{}
	if _, err := bitioex.CopyBits(bb, s.ActualBitBuf()); err != nil {
		return s, err
	}
	for _, bs := range bss {
		if bytes.Equal(bb.Bytes(), bs) {
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
