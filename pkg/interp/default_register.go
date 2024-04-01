package interp

import (
	"io/fs"

	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/gojq"
)

// DefaultRegistry global registry used by formats
var DefaultRegistry = NewRegistry()

func RegisterFormat(group *decode.Group, format *decode.Format) {
	DefaultRegistry.Format(group, format)
}

func RegisterFS(fs fs.ReadDirFS) {
	DefaultRegistry.FS(fs)
}

func RegisterFunc0[Tc any](name string, fn func(e *Interp, c Tc) any) {
	DefaultRegistry.Func(gojqx.Func0(name, fn))
}

func RegisterFunc1[Tc any, Ta0 any](name string, fn func(e *Interp, c Tc, a0 Ta0) any) {
	DefaultRegistry.Func(gojqx.Func1(name, fn))
}

func RegisterFunc2[Tc any, Ta0 any, Ta1 any](name string, fn func(e *Interp, c Tc, a0 Ta0, a1 Ta1) any) {
	DefaultRegistry.Func(gojqx.Func2(name, fn))
}

func RegisterIter0[Tc any](name string, fn func(e *Interp, c Tc) gojq.Iter) {
	DefaultRegistry.Func(gojqx.Iter0(name, fn))
}

func RegisterIter1[Tc any, Ta0 any](name string, fn func(e *Interp, c Tc, a0 Ta0) gojq.Iter) {
	DefaultRegistry.Func(gojqx.Iter1(name, fn))
}

func RegisterIter2[Tc any, Ta0 any, Ta1 any](name string, fn func(e *Interp, c Tc, a0 Ta0, a1 Ta1) gojq.Iter) {
	DefaultRegistry.Func(gojqx.Iter2(name, fn))
}
