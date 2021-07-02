package registry

import (
	"errors"
	"fmt"
	"fq/pkg/decode"
	"sort"
	"sync"
)

type Registry struct {
	Groups      map[string][]*decode.Format
	resolveOnce sync.Once
}

func New() *Registry {
	return &Registry{
		Groups:      map[string][]*decode.Format{},
		resolveOnce: sync.Once{},
	}
}

func (r *Registry) register(groupName string, format *decode.Format, single bool) *decode.Format {
	formats, ok := r.Groups[groupName]
	if ok {
		if !single {
			panic(fmt.Sprintf("%s: format already registered", groupName))
		}
	} else {
		formats = []*decode.Format{}
	}

	r.Groups[groupName] = append(formats, format)

	return format
}

func (r *Registry) resolve() error {
	for _, fs := range r.Groups {
		for _, f := range fs {
			for _, d := range f.Dependencies {
				var formats []*decode.Format
				for _, dName := range d.Names {
					df, ok := r.Groups[dName]
					if !ok {
						return fmt.Errorf("%s: can't find format dependency %s", f.Name, dName)
					}
					formats = append(formats, df...)
				}

				*d.Formats = formats
			}
		}
	}

	for _, fs := range r.Groups {
		sort.Slice(fs, func(i, j int) bool {
			if fs[i].ProbeWeight == fs[j].ProbeWeight {
				return fs[i].Name < fs[j].Name
			}
			return fs[i].ProbeWeight < fs[j].ProbeWeight
		})
	}

	return nil
}

func (r *Registry) MustRegister(format *decode.Format) *decode.Format {
	r.register(format.Name, format, false)
	for _, g := range format.Groups {
		r.register(g, format, true)
	}
	r.register("all", format, true)

	return format
}

func (r *Registry) Group(name string) ([]*decode.Format, error) {
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

func (r *Registry) MustGroup(name string) []*decode.Format {
	if g, err := r.Group(name); err == nil {
		return g
	} else {
		panic(err)
	}
}

func (r *Registry) MustAll() []*decode.Format {
	return r.MustGroup("all")
}
