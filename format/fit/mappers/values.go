package mappers

import (
	"fmt"
	"math"

	"github.com/wader/fq/pkg/scalar"
)

// Convertion from semicircles to decimal longitude latitude
var scConst = float64(180 / math.Pow(2, 31))

func GetUintFormatter(formatter string, unit string, scale float64, offset int64) scalar.UintFn {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		if scale != 0.0 && offset != 0 {
			s.Sym = (float64(s.Actual) / scale) - float64(offset)
		} else {
			if scale != 0.0 {
				s.Sym = float64(s.Actual) / scale
			}
			if offset != 0 {
				s.Sym = int64(s.Actual) - (offset)
			}
		}

		s.Description = unit
		if t, ok := TypeDefMap[formatter]; ok {
			if u, innerok := t[s.Actual]; innerok {
				s.Sym = u.Name
			}
		}
		return s, nil
	})
}

func GetSintFormatter(formatter string, unit string, scale float64, offset int64) scalar.SintFn {
	return scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
		if unit == "semicircles" {
			s.Sym = fmt.Sprintf("%f", float64(s.Actual)*scConst)
		} else {
			if scale != 0.0 && offset != 0 {
				s.Sym = (float64(s.Actual) / scale) - float64(offset)
			} else {
				if scale != 0.0 {
					s.Sym = float64(s.Actual) / scale
				}
				if offset != 0 {
					s.Sym = int64(s.Actual) - (offset)
				}
			}

			s.Description = unit
		}
		return s, nil
	})
}
