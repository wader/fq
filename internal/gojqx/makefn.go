package gojqx

//go:generate sh -c "cat makefn_gen.go.tmpl | go run ../../dev/tmpl.go | gofmt > makefn_gen.go"

import (
	"github.com/wader/gojq"
)

type Function struct {
	Name     string
	MinArity int
	MaxArity int
	FuncFn   func(any, []any) any
	IterFn   func(any, []any) gojq.Iter
}
