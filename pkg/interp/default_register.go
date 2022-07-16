package interp

import (
	"io/fs"

	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/gojq"
)

// DefaultRegister global registry used by formats
var DefaultRegister = NewRegistry()

func RegisterFormat(format decode.Format) {
	DefaultRegister.Format(format)
}

func RegisterFS(fs fs.ReadDirFS) {
	DefaultRegister.FS(fs)
}

func RegisterFunc0[Tc any](name string, fn func(e *Interp, c Tc) any) {
	DefaultRegister.Func(gojqextra.Func0(name, fn))
}

func RegisterFunc1[Tc any, Ta0 any](name string, fn func(e *Interp, c Tc, a0 Ta0) any) {
	DefaultRegister.Func(gojqextra.Func1(name, fn))
}

func RegisterFunc2[Tc any, Ta0 any, Ta1 any](name string, fn func(e *Interp, c Tc, a0 Ta0, a1 Ta1) any) {
	DefaultRegister.Func(gojqextra.Func2(name, fn))
}

func RegisterIter0[Tc any](name string, fn func(e *Interp, c Tc) gojq.Iter) {
	DefaultRegister.Func(gojqextra.Iter0(name, fn))
}

func RegisterIter1[Tc any, Ta0 any](name string, fn func(e *Interp, c Tc, a0 Ta0) gojq.Iter) {
	DefaultRegister.Func(gojqextra.Iter1(name, fn))
}

func RegisterIter2[Tc any, Ta0 any, Ta1 any](name string, fn func(e *Interp, c Tc, a0 Ta0, a1 Ta1) gojq.Iter) {
	DefaultRegister.Func(gojqextra.Iter2(name, fn))
}
