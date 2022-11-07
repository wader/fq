package common

import (
	"github.com/wader/fq/pkg/scalar"
)

//nolint:revive
const (
	PageSize                 = 8192
	FirstNormalTransactionID = 3

	LP_UNUSED   = 0 /* unused (should always have lp_len=0) */
	LP_NORMAL   = 1 /* used (should always have lp_len>0) */
	LP_REDIRECT = 2 /* HOT redirect (should have lp_len=0) */
	LP_DEAD     = 3
)

func TransactionIDIsNormal(xid uint64) bool {
	return xid >= FirstNormalTransactionID
}

type lpFlagsMapper struct{}

func (m lpFlagsMapper) MapScalar(s scalar.S) (scalar.S, error) {
	switch s.ActualU() {
	case LP_UNUSED:
		s.Sym = "LP_UNUSED"
	case LP_NORMAL:
		s.Sym = "LP_NORMAL"
	case LP_REDIRECT:
		s.Sym = "LP_REDIRECT"
	case LP_DEAD:
		s.Sym = "LP_DEAD"
	}
	return s, nil
}

var LpFlagsMapper = lpFlagsMapper{}

type Mask struct {
	Mask uint64
}

func (m Mask) MapScalar(s scalar.S) (scalar.S, error) {
	m1 := s.ActualU()
	v := IsMaskSet(m1, m.Mask)
	s.Actual = v
	return s, nil
}
