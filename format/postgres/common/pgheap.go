package common

import (
	"github.com/wader/fq/pkg/scalar"
)

const (
	HeapPageSize = 8192
)

type lpOffMapper struct{}

func (m lpOffMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := s.ActualU() & 0x7fff
	s.Actual = v
	return s, nil
}

var LpOffMapper = lpOffMapper{}

type lpFlagsMapper struct{}

func (m lpFlagsMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := (s.ActualU() >> 15) & 0x3
	s.Actual = v
	return s, nil
}

var LpFlagsMapper = lpFlagsMapper{}

type lpLenMapper struct{}

func (m lpLenMapper) MapScalar(s scalar.S) (scalar.S, error) {
	v := (s.ActualU() >> 17) & 0x7fff
	s.Actual = v
	return s, nil
}

var LpLenMapper = lpLenMapper{}

type Mask struct {
	Mask uint64
}

func (m Mask) MapScalar(s scalar.S) (scalar.S, error) {
	m1 := s.ActualU()
	v := IsMaskSet(m1, m.Mask)
	s.Actual = v
	return s, nil
}
