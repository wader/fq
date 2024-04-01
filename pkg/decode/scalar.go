package decode

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/scalar"
)

func bitBufIsZero(s scalar.BitBuf, isValidate bool) (scalar.BitBuf, error) {
	br := s.Actual

	isZero := true
	// TODO: shared
	b := make([]byte, 32*1024)
	bLen := int64(len(b)) * 8
	brLen, err := bitiox.Len(br)
	if err != nil {
		return scalar.BitBuf{}, err
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

func (d *D) BitBufIsZero() scalar.BitBufMapper {
	return scalar.BitBufFn(func(s scalar.BitBuf) (scalar.BitBuf, error) {
		return bitBufIsZero(s, false)
	})
}

func (d *D) BitBufValidateIsZero() scalar.BitBufMapper {
	return scalar.BitBufFn(func(s scalar.BitBuf) (scalar.BitBuf, error) {
		return bitBufIsZero(s, !d.Options.Force)
	})
}

// TODO: generate?
func assertBitBuf(s scalar.BitBuf, isErr bool, bss ...[]byte) (scalar.BitBuf, error) {
	bb := &bytes.Buffer{}
	if _, err := bitiox.CopyBits(bb, s.Actual); err != nil {
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

func (d *D) AssertBitBuf(bss ...[]byte) scalar.BitBufMapper {
	return scalar.BitBufFn(func(s scalar.BitBuf) (scalar.BitBuf, error) {
		return assertBitBuf(s, !d.Options.Force, bss...)
	})
}

func (d *D) ValidateBitBuf(bss ...[]byte) scalar.BitBufMapper {
	return scalar.BitBufFn(func(s scalar.BitBuf) (scalar.BitBuf, error) {
		return assertBitBuf(s, false, bss...)
	})
}

func UintAssertBytes(s scalar.Uint, isErr bool, endian Endian, bss ...[]byte) (scalar.Uint, error) {
	var bo binary.ByteOrder
	switch endian {
	case LittleEndian:
		bo = binary.BigEndian
	case BigEndian:
		bo = binary.BigEndian
	default:
		panic("invalid endian")
	}

	au := s.Actual
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

func (d *D) UintAssertBytes(bss ...[]byte) scalar.UintMapper {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		return UintAssertBytes(s, true, d.Endian, bss...)
	})
}
func (d *D) UintValidateBytes(bss ...[]byte) scalar.UintMapper {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		return UintAssertBytes(s, false, d.Endian, bss...)
	})
}
func (d *D) AssertULEBytes(bss ...[]byte) scalar.UintMapper {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		return UintAssertBytes(s, true, LittleEndian, bss...)
	})
}
func (d *D) UintValidateLEBytes(bss ...[]byte) scalar.UintMapper {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		return UintAssertBytes(s, false, LittleEndian, bss...)
	})
}
func (d *D) AssertUBEBytes(bss ...[]byte) scalar.UintMapper {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		return UintAssertBytes(s, true, BigEndian, bss...)
	})
}
func (d *D) UintValidateBEBytes(bss ...[]byte) scalar.UintMapper {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		return UintAssertBytes(s, false, BigEndian, bss...)
	})
}
