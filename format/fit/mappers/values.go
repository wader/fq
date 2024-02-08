package mappers

import (
	"math"

	"github.com/wader/fq/pkg/scalar"
)

// Used for conversion from semicircles to decimal longitude latitude
var scConst = 180 / math.Pow(2, 31)

var invalidUint = map[string]uint64{
	"byte":    0xFF,
	"enum":    0xFF,
	"uint8":   0xFF,
	"uint8z":  0x00,
	"uint16":  0xFFFF,
	"uint16z": 0x0000,
	"uint32":  0xFFFFFFFF,
	"uint32z": 0x00000000,
	"uint64":  0xFFFFFFFFFFFFFFFF,
	"uint64z": 0x0000000000000000,
}

var invalidSint = map[string]int64{
	"sint8":  0x7F,
	"sint16": 0x7FFF,
	"sint32": 0x7FFFFFFF,
	"sint64": 0x7FFFFFFFFFFFFFFF,
}

var invalidFloat = map[string]float64{
	"float32": 0xFFFFFFFF,
	"float64": 0xFFFFFFFFFFFFFFFF,
}

func GetUintFormatter(fDef FieldDef) scalar.UintFn {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		if s.Actual == invalidUint[fDef.Type] {
			s.Sym = "[invalid]"
			return s, nil
		}
		if fDef.Scale != 0.0 && fDef.Offset != 0 {
			s.Sym = (float64(s.Actual) / fDef.Scale) - float64(fDef.Offset)
		} else {
			if fDef.Scale != 0.0 {
				s.Sym = float64(s.Actual) / fDef.Scale
			}
			if fDef.Offset != 0 {
				s.Sym = int64(s.Actual) - fDef.Offset
			}
		}

		s.Description = fDef.Unit
		if t, ok := TypeDefMap[fDef.Format]; ok {
			if u, innerok := t[s.Actual]; innerok {
				s.Sym = u.Name
			}
		}
		return s, nil
	})
}

func GetSintFormatter(fDef FieldDef) scalar.SintFn {
	return scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
		if s.Actual == invalidSint[fDef.Type] {
			s.Sym = "[invalid]"
			return s, nil
		}
		if fDef.Unit == "semicircles" {
			s.Sym = float64(s.Actual) * scConst
		} else {
			if fDef.Scale != 0.0 && fDef.Offset != 0 {
				s.Sym = (float64(s.Actual) / fDef.Scale) - float64(fDef.Offset)
			} else {
				if fDef.Scale != 0.0 {
					s.Sym = float64(s.Actual) / fDef.Scale
				}
				if fDef.Offset != 0 {
					s.Sym = s.Actual - fDef.Offset
				}
			}

			s.Description = fDef.Unit
		}
		return s, nil
	})
}

func GetFloatFormatter(fDef FieldDef) scalar.FltFn {
	return scalar.FltFn(func(s scalar.Flt) (scalar.Flt, error) {
		if s.Actual == invalidFloat[fDef.Type] {
			s.Sym = "[invalid]"
			return s, nil
		}
		if fDef.Scale != 0.0 && fDef.Offset != 0 {
			s.Sym = (s.Actual / fDef.Scale) - float64(fDef.Offset)
		} else {
			if fDef.Scale != 0.0 {
				s.Sym = s.Actual / fDef.Scale
			}
			if fDef.Offset != 0 {
				s.Sym = s.Actual - float64(fDef.Offset)
			}
		}

		s.Description = fDef.Unit
		return s, nil
	})
}
