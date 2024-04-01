package interp

import (
	"math/big"

	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/gojq"
)

func init() {
	RegisterFunc0("bnot", (*Interp).bnot)
	RegisterFunc2("bsl", (*Interp).bsl)
	RegisterFunc2("bsr", (*Interp).bsr)
	RegisterFunc2("band", (*Interp).band)
	RegisterFunc2("bor", (*Interp).bor)
	RegisterFunc2("bxor", (*Interp).bxor)
}

func (i *Interp) bnot(c any) any {
	switch c := c.(type) {
	case int:
		return ^c
	case *big.Int:
		return new(big.Int).Not(c)
	case gojq.JQValue:
		return i.bnot(c.JQValueToGoJQ())
	default:
		return &gojqx.UnaryTypeError{Name: "bnot", V: c}
	}
}

func (i *Interp) bsl(c any, a any, b any) any {
	return gojq.BinopTypeSwitch(a, b,
		func(l, r int) any {
			if v := l << r; v>>r == l {
				return v
			}
			return new(big.Int).Lsh(big.NewInt(int64(l)), uint(r))
		},
		func(l, r float64) any { return int(l) << int(r) },
		func(l, r *big.Int) any { return new(big.Int).Lsh(l, uint(r.Uint64())) },
		func(l, r string) any { return &gojqx.BinopTypeError{Name: "bsl", L: l, R: r} },
		func(l, r []any) any { return &gojqx.BinopTypeError{Name: "bsl", L: l, R: r} },
		func(l, r map[string]any) any {
			return &gojqx.BinopTypeError{Name: "bsl", L: l, R: r}
		},
		func(l, r any) any { return &gojqx.BinopTypeError{Name: "bsl", L: l, R: r} },
	)
}

func (i *Interp) bsr(c any, a any, b any) any {
	return gojq.BinopTypeSwitch(a, b,
		func(l, r int) any { return l >> r },
		func(l, r float64) any { return int(l) >> int(r) },
		func(l, r *big.Int) any { return new(big.Int).Rsh(l, uint(r.Uint64())) },
		func(l, r string) any { return &gojqx.BinopTypeError{Name: "bsr", L: l, R: r} },
		func(l, r []any) any { return &gojqx.BinopTypeError{Name: "bsr", L: l, R: r} },
		func(l, r map[string]any) any {
			return &gojqx.BinopTypeError{Name: "bsr", L: l, R: r}
		},
		func(l, r any) any { return &gojqx.BinopTypeError{Name: "bsr", L: l, R: r} },
	)
}

func (i *Interp) band(c any, a any, b any) any {
	return gojq.BinopTypeSwitch(a, b,
		func(l, r int) any { return l & r },
		func(l, r float64) any { return int(l) & int(r) },
		func(l, r *big.Int) any { return new(big.Int).And(l, r) },
		func(l, r string) any { return &gojqx.BinopTypeError{Name: "band", L: l, R: r} },
		func(l, r []any) any { return &gojqx.BinopTypeError{Name: "band", L: l, R: r} },
		func(l, r map[string]any) any {
			return &gojqx.BinopTypeError{Name: "band", L: l, R: r}
		},
		func(l, r any) any { return &gojqx.BinopTypeError{Name: "band", L: l, R: r} },
	)
}

func (i *Interp) bor(c any, a any, b any) any {
	return gojq.BinopTypeSwitch(a, b,
		func(l, r int) any { return l | r },
		func(l, r float64) any { return int(l) | int(r) },
		func(l, r *big.Int) any { return new(big.Int).Or(l, r) },
		func(l, r string) any { return &gojqx.BinopTypeError{Name: "bor", L: l, R: r} },
		func(l, r []any) any { return &gojqx.BinopTypeError{Name: "bor", L: l, R: r} },
		func(l, r map[string]any) any {
			return &gojqx.BinopTypeError{Name: "bor", L: l, R: r}
		},
		func(l, r any) any { return &gojqx.BinopTypeError{Name: "bor", L: l, R: r} },
	)
}

func (i *Interp) bxor(c any, a any, b any) any {
	return gojq.BinopTypeSwitch(a, b,
		func(l, r int) any { return l ^ r },
		func(l, r float64) any { return int(l) ^ int(r) },
		func(l, r *big.Int) any { return new(big.Int).Xor(l, r) },
		func(l, r string) any { return &gojqx.BinopTypeError{Name: "bxor", L: l, R: r} },
		func(l, r []any) any { return &gojqx.BinopTypeError{Name: "bxor", L: l, R: r} },
		func(l, r map[string]any) any {
			return &gojqx.BinopTypeError{Name: "bxor", L: l, R: r}
		},
		func(l, r any) any { return &gojqx.BinopTypeError{Name: "bxor", L: l, R: r} },
	)
}
