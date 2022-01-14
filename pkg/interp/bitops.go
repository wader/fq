package interp

import (
	"math/big"

	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/gojq"
)

func init() {
	functionRegisterFns = append(functionRegisterFns, func(i *Interp) []Function {
		return []Function{
			{"bnot", 0, 0, i.bnot, nil},
			{"bsl", 2, 2, i.bsl, nil},
			{"bsr", 2, 2, i.bsr, nil},
			{"band", 2, 2, i.band, nil},
			{"bor", 2, 2, i.bor, nil},
			{"bxor", 2, 2, i.bxor, nil},
		}
	})
}

func (i *Interp) bnot(c interface{}, a []interface{}) interface{} {
	switch c := c.(type) {
	case int:
		return ^c
	case *big.Int:
		return new(big.Int).Not(c)
	case gojq.JQValue:
		return i.bnot(c.JQValueToGoJQ(), a)
	default:
		return &gojqextra.UnaryTypeError{Name: "bnot", V: c}
	}
}

func (i *Interp) bsl(c interface{}, a []interface{}) interface{} {
	return gojq.BinopTypeSwitch(a[0], a[1],
		func(l, r int) bool { return false }, // TODO: can be int safe i think
		func(l, r int) interface{} { return l << r },
		func(l, r float64) interface{} { return int(l) << int(r) },
		func(l, r *big.Int) interface{} { return new(big.Int).Lsh(l, uint(r.Uint64())) },
		func(l, r string) interface{} { return &gojqextra.BinopTypeError{Name: "bsl", L: l, R: r} },
		func(l, r []interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bsl", L: l, R: r} },
		func(l, r map[string]interface{}) interface{} {
			return &gojqextra.BinopTypeError{Name: "bsl", L: l, R: r}
		},
		func(l, r interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bsl", L: l, R: r} },
	)
}

func (i *Interp) bsr(c interface{}, a []interface{}) interface{} {
	return gojq.BinopTypeSwitch(a[0], a[1],
		func(l, r int) bool { return true },
		func(l, r int) interface{} { return l >> r },
		func(l, r float64) interface{} { return int(l) >> int(r) },
		func(l, r *big.Int) interface{} { return new(big.Int).Rsh(l, uint(r.Uint64())) },
		func(l, r string) interface{} { return &gojqextra.BinopTypeError{Name: "bsr", L: l, R: r} },
		func(l, r []interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bsr", L: l, R: r} },
		func(l, r map[string]interface{}) interface{} {
			return &gojqextra.BinopTypeError{Name: "bsr", L: l, R: r}
		},
		func(l, r interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bsr", L: l, R: r} },
	)
}

func (i *Interp) band(c interface{}, a []interface{}) interface{} {
	return gojq.BinopTypeSwitch(a[0], a[1],
		func(l, r int) bool { return true },
		func(l, r int) interface{} { return l & r },
		func(l, r float64) interface{} { return int(l) & int(r) },
		func(l, r *big.Int) interface{} { return new(big.Int).And(l, r) },
		func(l, r string) interface{} { return &gojqextra.BinopTypeError{Name: "band", L: l, R: r} },
		func(l, r []interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "band", L: l, R: r} },
		func(l, r map[string]interface{}) interface{} {
			return &gojqextra.BinopTypeError{Name: "band", L: l, R: r}
		},
		func(l, r interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "band", L: l, R: r} },
	)
}

func (i *Interp) bor(c interface{}, a []interface{}) interface{} {
	return gojq.BinopTypeSwitch(a[0], a[1],
		func(l, r int) bool { return true },
		func(l, r int) interface{} { return l | r },
		func(l, r float64) interface{} { return int(l) | int(r) },
		func(l, r *big.Int) interface{} { return new(big.Int).Or(l, r) },
		func(l, r string) interface{} { return &gojqextra.BinopTypeError{Name: "bor", L: l, R: r} },
		func(l, r []interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bor", L: l, R: r} },
		func(l, r map[string]interface{}) interface{} {
			return &gojqextra.BinopTypeError{Name: "bor", L: l, R: r}
		},
		func(l, r interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bor", L: l, R: r} },
	)
}

func (i *Interp) bxor(c interface{}, a []interface{}) interface{} {
	return gojq.BinopTypeSwitch(a[0], a[1],
		func(l, r int) bool { return true },
		func(l, r int) interface{} { return l ^ r },
		func(l, r float64) interface{} { return int(l) ^ int(r) },
		func(l, r *big.Int) interface{} { return new(big.Int).Xor(l, r) },
		func(l, r string) interface{} { return &gojqextra.BinopTypeError{Name: "bxor", L: l, R: r} },
		func(l, r []interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bxor", L: l, R: r} },
		func(l, r map[string]interface{}) interface{} {
			return &gojqextra.BinopTypeError{Name: "bxor", L: l, R: r}
		},
		func(l, r interface{}) interface{} { return &gojqextra.BinopTypeError{Name: "bxor", L: l, R: r} },
	)
}
