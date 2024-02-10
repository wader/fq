package mappers

import (
	"github.com/wader/fq/pkg/scalar"
)

type TypeDefLookup struct {
	Name string
	Type string
}

type typeDefMap map[uint64]TypeDefLookup

func (m typeDefMap) MapUint(s scalar.Uint) (scalar.Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t.Name
	}
	return s, nil
}
