package decode

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

type Registry struct {
	allGroups   map[string][]*Format
	resolveOnce sync.Once
}

func NewRegistry() *Registry {
	return &Registry{
		allGroups:   map[string][]*Format{},
		resolveOnce: sync.Once{},
	}
}

func (r *Registry) register(groupName string, format *Format, single bool) *Format {
	formats, ok := r.allGroups[groupName]
	if ok {
		if !single {
			panic(fmt.Sprintf("%s: already registered", groupName))
		}
	} else {
		formats = []*Format{}
	}

	// prepend to allow override
	r.allGroups[groupName] = append([]*Format{format}, formats...)

	return format
}

func (r *Registry) resolve() error {
	for _, fs := range r.allGroups {
		for _, f := range fs {
			for _, d := range f.Deps {
				var formats []*Format
				for _, dName := range d.Names {
					df, ok := r.allGroups[dName]
					if !ok {
						return fmt.Errorf("%s: can't find dependency %s", f.Name, dName)
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

	if g, ok := r.allGroups[name]; ok {
		return g, nil
	}
	return nil, errors.New("no such group")
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

func (r *Registry) Dot(w io.Writer) {
	formatSeen := map[string]struct{}{}

	fmt.Fprintf(w, "digraph formats {\n")
	for groupName, fs := range r.allGroups {
		if groupName == "all" {
			continue
		}
		for _, f := range fs {
			if groupName != f.Name {
				fmt.Fprintf(w, "  %s -> %s\n", groupName, f.Name)
			}

			if _, ok := formatSeen[f.Name]; ok {
				continue
			}
			formatSeen[f.Name] = struct{}{}

			for _, d := range f.Deps {

				for _, dName := range d.Names {
					fmt.Fprintf(w, "  %s -> %s\n", f.Name, dName)
				}
			}
		}
	}
	fmt.Fprintf(w, "}\n")
}
