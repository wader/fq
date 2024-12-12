package interp

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
	"sync"

	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/pkg/decode"
)

type EnvFuncFn func(env *Interp) gojqx.Function

type Registry struct {
	allGroup          *decode.Group
	groups            map[string]*decode.Group
	formatResolveOnce sync.Once
	formatResolved    bool

	EnvFuncFns []EnvFuncFn

	FSs []fs.ReadDirFS
}

func NewRegistry() *Registry {
	return &Registry{
		allGroup:          &decode.Group{Name: "all"},
		groups:            map[string]*decode.Group{},
		formatResolveOnce: sync.Once{},
	}
}

func (r *Registry) Format(group *decode.Group, format *decode.Format) *decode.Format {
	if r.formatResolved {
		// for now can't change after resolved
		panic("registry already resolved")
	}

	if _, ok := r.groups[group.Name]; ok {
		panic(fmt.Sprintf("%s: format already registered", group.Name))
	}
	group.Formats = append(group.Formats, format)
	r.groups[group.Name] = group
	format.Name = group.Name

	for _, g := range format.Groups {
		r.groups[g.Name] = g
		g.Formats = append(g.Formats, format)
	}

	r.allGroup.Formats = append(r.allGroup.Formats, format)

	return format
}

func (r *Registry) FS(fs fs.ReadDirFS) {
	r.FSs = append(r.FSs, fs)
}

func (r *Registry) Func(funcFn EnvFuncFn) {
	r.EnvFuncFns = append(r.EnvFuncFns, funcFn)
}

func sortFormats(g *decode.Group) {
	slices.SortFunc(g.Formats, func(a, b *decode.Format) int {
		if a.ProbeOrder == b.ProbeOrder {
			return strings.Compare(a.Name, b.Name)
		}
		return cmp.Compare(a.ProbeOrder, b.ProbeOrder)
	})
}

func (r *Registry) resolveGroups() {
	r.formatResolveOnce.Do(func() {
		for _, g := range r.groups {
			for _, f := range g.Formats {
				for _, d := range f.Dependencies {
					if len(d.Out.Formats) != 0 {
						// already resolved
						continue
					}

					for _, dg := range d.Groups {
						d.Out.Formats = append(d.Out.Formats, dg.Formats...)
					}
					sortFormats(d.Out)
				}
			}
		}

		for _, g := range r.groups {
			sortFormats(g)
		}

		r.formatResolved = true
	})
}

func (r *Registry) Group(name string) (*decode.Group, error) {
	r.resolveGroups()
	if g, ok := r.groups[name]; ok {
		return g, nil
	}
	return nil, errors.New("format group not found")
}

func (r *Registry) MustGroup(name string) *decode.Group {
	g, err := r.Group(name)
	if err == nil {
		return g
	}
	panic(err)
}

func (r *Registry) Groups() map[string]*decode.Group {
	r.resolveGroups()
	return r.groups
}

func (r *Registry) MustAll() *decode.Group {
	r.resolveGroups()
	return r.allGroup
}
