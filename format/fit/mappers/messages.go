package mappers

import (
	"github.com/wader/fq/pkg/scalar"
)

type FieldDef struct {
	Name        string
	Type        string
	Unit        string
	Scale       float64
	Offset      int64
	Size        uint64
	HasSubField bool
}

type LocalFieldDef struct {
	Name             string
	Type             string
	Format           string
	Unit             string
	Scale            float64
	Offset           int64
	Size             uint64
	GlobalFieldDef   FieldDef
	GlobalMessageNo  uint64
	GlobalFieldDefNo uint64
}

type fieldDefMap map[uint64]FieldDef

func (m fieldDefMap) MapUint(s scalar.Uint) (scalar.Uint, error) {
	if t, ok := m[s.Actual]; ok {
		s.Sym = t.Name
	}
	return s, nil
}
