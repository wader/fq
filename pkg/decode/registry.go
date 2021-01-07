package decode

import (
	"errors"
	"fmt"
	"sync"
)

type Registry struct {
	Groups      map[string][]*Format
	resolveOnce sync.Once
}

func NewRegistry() *Registry {
	return &Registry{
		Groups:      map[string][]*Format{},
		resolveOnce: sync.Once{},
	}
}

func (r *Registry) register(groupName string, format *Format, single bool) *Format {
	formats, ok := r.Groups[groupName]
	if ok {
		if !single {
			panic(fmt.Sprintf("%s: format already registered", groupName))
		}
	} else {
		formats = []*Format{}
	}

	// prepend to allow override
	r.Groups[groupName] = append([]*Format{format}, formats...)

	return format
}

func (r *Registry) resolve() error {
	for _, fs := range r.Groups {
		for _, f := range fs {
			for _, d := range f.Dependencies {
				var formats []*Format
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

	return nil
}

func (r *Registry) MustRegister(format *Format) *Format {
	r.register(format.Name, format, false)
	for _, g := range format.Groups {
		r.register(g, format, true)
	}
	r.register("all", format, true)

	return format
}

func (r *Registry) Group(name string) ([]*Format, error) {
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

func (r *Registry) MustGroup(name string) []*Format {
	if g, err := r.Group(name); err == nil {
		return g
	} else {
		panic(err)
	}
}

func (r *Registry) MustAll() []*Format {
	return r.MustGroup("all")
}
