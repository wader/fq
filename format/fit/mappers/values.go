package mappers

import (
	"math"
	"time"

	"github.com/wader/fq/pkg/scalar"
)

var epochDate = time.Date(1989, time.December, 31, 0, 0, 0, 0, time.UTC)

// Used for conversion from semicircles to decimal longitude latitude
var scConst = 180 / math.Pow(2, 31)

var invalidUint = map[string]uint64{
	"byte":    0xff,
	"enum":    0xff,
	"uint8":   0xff,
	"uint8z":  0x00,
	"uint16":  0xffff,
	"uint16z": 0x0000,
	"uint32":  0xffff_ffff,
	"uint32z": 0x0000_0000,
	"uint64":  0xffff_ffff_ffff_ffff,
	"uint64z": 0x0000_0000_0000_0000,
}

var invalidSint = map[string]int64{
	"sint8":  0x7f,
	"sint16": 0x7fff,
	"sint32": 0x7fff_ffff,
	"sint64": 0x7fff_ffff_ffff_ffff,
}

var invalidFloat = map[string]float64{
	"float32": 0xffff_ffff,
	"float64": 0xffff_ffff_ffff_ffff,
}

func GetUintFormatter(fDef LocalFieldDef) scalar.UintFn {
	return scalar.UintFn(func(s scalar.Uint) (scalar.Uint, error) {
		if s.Actual == invalidUint[fDef.Type] {
			s.Description = "invalid"
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

		if t, ok := TypeDefMap[fDef.Format]; ok {
			if u, innerOk := t[s.Actual]; innerOk {
				s.Sym = u.Name
			}
		}

		switch fDef.Format {
		case "date_time",
			"local_date_time":
			s.Description = epochDate.Add(time.Duration(s.Actual) * time.Second).Format(time.RFC3339)
		default:
			s.Description = fDef.Unit
		}

		return s, nil
	})
}

func GetSintFormatter(fDef LocalFieldDef) scalar.SintFn {
	return scalar.SintFn(func(s scalar.Sint) (scalar.Sint, error) {
		if s.Actual == invalidSint[fDef.Type] {
			s.Description = "invalid"
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

func GetFloatFormatter(fDef LocalFieldDef) scalar.FltFn {
	return scalar.FltFn(func(s scalar.Flt) (scalar.Flt, error) {
		if s.Actual == invalidFloat[fDef.Type] {
			s.Description = "invalid"
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
