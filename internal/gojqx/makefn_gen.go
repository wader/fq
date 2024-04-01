// Generated from makefn_gen.go.tmpl
package gojqx

import (
	"github.com/wader/fq/internal/mapstruct"
	"github.com/wader/gojq"
)

func Func0[Tenv any, Tc any](name string, fn func(e Tenv, c Tc) any) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 0, MaxArity: 0}
		f.FuncFn = func(c any, a []any) any {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return FuncTypeError{Name: name, V: c}
			}

			return fn(env, cv)
		}
		return f
	}
}

func Func1[Tenv any, Tc any, Ta0 any](name string, fn func(e Tenv, c Tc, a0 Ta0) any) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 1, MaxArity: 1}
		f.FuncFn = func(c any, a []any) any {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return FuncTypeError{Name: name, V: c}
			}
			a0, ok := CastFn[Ta0](a[0], mapstruct.ToStruct)
			if !ok {
				return FuncArgTypeError{Name: name, ArgName: "first", V: a[0]}
			}

			return fn(env, cv, a0)
		}
		return f
	}
}

func Func2[Tenv any, Tc any, Ta0 any, Ta1 any](name string, fn func(e Tenv, c Tc, a0 Ta0, a1 Ta1) any) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 2, MaxArity: 2}
		f.FuncFn = func(c any, a []any) any {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return FuncTypeError{Name: name, V: c}
			}
			a0, ok := CastFn[Ta0](a[0], mapstruct.ToStruct)
			if !ok {
				return FuncArgTypeError{Name: name, ArgName: "first", V: a[0]}
			}
			a1, ok := CastFn[Ta1](a[1], mapstruct.ToStruct)
			if !ok {
				return FuncArgTypeError{Name: name, ArgName: "second", V: a[1]}
			}

			return fn(env, cv, a0, a1)
		}
		return f
	}
}

func Func3[Tenv any, Tc any, Ta0 any, Ta1 any, Ta2 any](name string, fn func(e Tenv, c Tc, a0 Ta0, a1 Ta1, a2 Ta2) any) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 3, MaxArity: 3}
		f.FuncFn = func(c any, a []any) any {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return FuncTypeError{Name: name, V: c}
			}
			a0, ok := CastFn[Ta0](a[0], mapstruct.ToStruct)
			if !ok {
				return FuncArgTypeError{Name: name, ArgName: "first", V: a[0]}
			}
			a1, ok := CastFn[Ta1](a[1], mapstruct.ToStruct)
			if !ok {
				return FuncArgTypeError{Name: name, ArgName: "second", V: a[1]}
			}
			a2, ok := CastFn[Ta2](a[2], mapstruct.ToStruct)
			if !ok {
				return FuncArgTypeError{Name: name, ArgName: "third", V: a[2]}
			}

			return fn(env, cv, a0, a1, a2)
		}
		return f
	}
}

func Iter0[Tenv any, Tc any](name string, fn func(e Tenv, c Tc) gojq.Iter) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 0, MaxArity: 0}
		f.IterFn = func(c any, a []any) gojq.Iter {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncTypeError{Name: name, V: c})
			}

			return fn(env, cv)
		}
		return f
	}
}

func Iter1[Tenv any, Tc any, Ta0 any](name string, fn func(e Tenv, c Tc, a0 Ta0) gojq.Iter) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 1, MaxArity: 1}
		f.IterFn = func(c any, a []any) gojq.Iter {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncTypeError{Name: name, V: c})
			}
			a0, ok := CastFn[Ta0](a[0], mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncArgTypeError{Name: name, ArgName: "first", V: a[0]})
			}

			return fn(env, cv, a0)
		}
		return f
	}
}

func Iter2[Tenv any, Tc any, Ta0 any, Ta1 any](name string, fn func(e Tenv, c Tc, a0 Ta0, a1 Ta1) gojq.Iter) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 2, MaxArity: 2}
		f.IterFn = func(c any, a []any) gojq.Iter {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncTypeError{Name: name, V: c})
			}
			a0, ok := CastFn[Ta0](a[0], mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncArgTypeError{Name: name, ArgName: "first", V: a[0]})
			}
			a1, ok := CastFn[Ta1](a[1], mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncArgTypeError{Name: name, ArgName: "second", V: a[1]})
			}

			return fn(env, cv, a0, a1)
		}
		return f
	}
}

func Iter3[Tenv any, Tc any, Ta0 any, Ta1 any, Ta2 any](name string, fn func(e Tenv, c Tc, a0 Ta0, a1 Ta1, a2 Ta2) gojq.Iter) func(env Tenv) Function {
	return func(env Tenv) Function {
		f := Function{Name: name, MinArity: 3, MaxArity: 3}
		f.IterFn = func(c any, a []any) gojq.Iter {
			cv, ok := CastFn[Tc](c, mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncTypeError{Name: name, V: c})
			}
			a0, ok := CastFn[Ta0](a[0], mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncArgTypeError{Name: name, ArgName: "first", V: a[0]})
			}
			a1, ok := CastFn[Ta1](a[1], mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncArgTypeError{Name: name, ArgName: "second", V: a[1]})
			}
			a2, ok := CastFn[Ta2](a[2], mapstruct.ToStruct)
			if !ok {
				return gojq.NewIter(FuncArgTypeError{Name: name, ArgName: "third", V: a[2]})
			}

			return fn(env, cv, a0, a1, a2)
		}
		return f
	}
}
