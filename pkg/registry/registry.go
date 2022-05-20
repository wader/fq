package registry

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/wader/fq/pkg/decode"
)

type Registry struct {
	Groups      map[string]decode.Group
	resolveOnce sync.Once
	resolved    bool
}

func New() *Registry {
	return &Registry{
		Groups:      map[string]decode.Group{},
		resolveOnce: sync.Once{},
	}
}

func (r *Registry) register(groupName string, format decode.Format, single bool) decode.Format {
	if r.resolved {
		// for now can't change after resolved
		panic("registry already resolved")
	}

	group, ok := r.Groups[groupName]
	if ok {
		if !single {
			panic(fmt.Sprintf("%s: format already registered", groupName))
		}
	} else {
		group = decode.Group{}
	}

	r.Groups[groupName] = append(group, format)

	return format
}

func (r *Registry) MustRegister(format decode.Format) decode.Format {
	r.register(format.Name, format, false)
	for _, g := range format.Groups {
		r.register(g, format, true)
	}
	r.register("all", format, true)

	return format
}

func sortFormats(g decode.Group) {
	sort.Slice(g, func(i, j int) bool {
		if g[i].ProbeOrder == g[j].ProbeOrder {
			return g[i].Name < g[j].Name
		}
		return g[i].ProbeOrder < g[j].ProbeOrder
	})
}

func (r *Registry) resolve() error {
	for _, fs := range r.Groups {
		for _, f := range fs {
			for _, d := range f.Dependencies {
				var group decode.Group
				for _, dName := range d.Names {
					df, ok := r.Groups[dName]
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

	for _, fs := range r.Groups {
		sortFormats(fs)
	}

	r.resolved = true

	return nil
}

func (r *Registry) Group(name string) (decode.Group, error) {
	r.resolveOnce.Do(func() {
		if err := r.resolve(); err != nil {
			panic(err)
		}
	})

	if g, ok := r.Groups[name]; ok {
		return g, nil
	}
	return nil, errors.New("format group not found")
}

func (r *Registry) MustGroup(name string) decode.Group {
	g, err := r.Group(name)
	if err == nil {
		return g
	}
	panic(err)
}

func (r *Registry) MustAll() decode.Group {
	return r.MustGroup("all")
}
