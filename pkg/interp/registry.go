package interp

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"sync"

	"github.com/wader/fq/internal/gojqex"
	"github.com/wader/fq/pkg/decode"
)

type EnvFuncFn func(env *Interp) gojqex.Function

type Registry struct {
	FormatGroups      map[string]decode.Group
	formatResolveOnce sync.Once
	formatResolved    bool

	EnvFuncFns []EnvFuncFn

	FSs []fs.ReadDirFS
}

func NewRegistry() *Registry {
	return &Registry{
		FormatGroups:      map[string]decode.Group{},
		formatResolveOnce: sync.Once{},
	}
}

func (r *Registry) format(groupName string, format decode.Format, single bool) {
	if r.formatResolved {
		// for now can't change after resolved
		panic("registry already resolved")
	}

	group, ok := r.FormatGroups[groupName]
	if ok {
		if !single {
			panic(fmt.Sprintf("%s: format already registered", groupName))
		}
	} else {
		group = decode.Group{}
	}

	r.FormatGroups[groupName] = append(group, format)
}

func (r *Registry) Format(format decode.Format) decode.Format {
	r.format(format.Name, format, false)
	for _, g := range format.Groups {
		r.format(g, format, true)
	}
	r.format("all", format, true)

	return format
}

func (r *Registry) FS(fs fs.ReadDirFS) {
	r.FSs = append(r.FSs, fs)
}

func (r *Registry) Func(funcFn EnvFuncFn) {
	r.EnvFuncFns = append(r.EnvFuncFns, funcFn)
}

func sortFormats(g decode.Group) {
	sort.Slice(g, func(i, j int) bool {
		if g[i].ProbeOrder == g[j].ProbeOrder {
			return g[i].Name < g[j].Name
		}
		return g[i].ProbeOrder < g[j].ProbeOrder
	})
}

func (r *Registry) resolveFormats() error {
	for _, fs := range r.FormatGroups {
		for _, f := range fs {
			for _, d := range f.Dependencies {
				var group decode.Group
				for _, dName := range d.Names {
					df, ok := r.FormatGroups[dName]
					if !ok {
						return fmt.Errorf("%s: can't find format dependency %s", f.Name, dName)
					}
					group = append(group, df...)
				}

				sortFormats(group)
				*d.Group = group
			}
		}
	}

	for _, fs := range r.FormatGroups {
		sortFormats(fs)
	}

	r.formatResolved = true

	return nil
}

func (r *Registry) FormatGroup(name string) (decode.Group, error) {
	r.formatResolveOnce.Do(func() {
		if err := r.resolveFormats(); err != nil {
			panic(err)
		}
	})

	if g, ok := r.FormatGroups[name]; ok {
		return g, nil
	}
	return nil, errors.New("format group not found")
}

func (r *Registry) MustFormatGroup(name string) decode.Group {
	g, err := r.FormatGroup(name)
	if err == nil {
		return g
	}
	panic(err)
}

func (r *Registry) MustAll() decode.Group {
	return r.MustFormatGroup("all")
}
