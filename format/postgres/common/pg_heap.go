package common

import (
	"github.com/wader/fq/pkg/scalar"
)

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

func (m lpFlagsMapper) MapUint(s scalar.Uint) (scalar.Uint, error) {
	switch s.Actual {
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

func (m Mask) MapUint(s scalar.Uint) (scalar.Uint, error) {
	m1 := s.Actual
	v := IsMaskSet(m1, m.Mask)
	s.Actual = v
	return s, nil
}
