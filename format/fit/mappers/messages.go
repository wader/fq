package mappers

import (
	"github.com/wader/fq/pkg/scalar"
)

type FieldDefLookup struct {
	Name      string
	Type      string
	Formatter string
	Unit      string
	Scale     float64
	Offset    int64
}
type fieldDefMap map[uint64]FieldDefLookup

func (m fieldDefMap) MapUint(s scalar.Uint) (scalar.Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t.Name
	}
	return s, nil
}
