package format

import (
	"errors"
	"fmt"
	"fq/pkg/decode"
	"io"
	"sync"
)

// TODO: some struct so it's more flexible

var allGroups = map[string][]*decode.Format{}

var resolveOnce = sync.Once{}

func register(groupName string, format *decode.Format, single bool) *decode.Format {
	formats, ok := allGroups[groupName]
	if ok {
		if !single {
			panic(fmt.Sprintf("%s: already registered", groupName))
		}
	} else {
		formats = []*decode.Format{}
	}
	allGroups[groupName] = append(formats, format)

	return format
}

func resolve() error {
	for _, fs := range allGroups {
		for _, f := range fs {
			for _, d := range f.Deps {
				var formats []*decode.Format
				for _, dName := range d.Names {
					df, ok := allGroups[dName]
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

func MustRegister(format *decode.Format) *decode.Format {
	register(format.Name, format, false)
	for _, g := range format.Groups {
		register(g, format, true)
	}
	if !format.SkipProbe {
		register("probeable", format, true)
	}
	register("all", format, true)

	return format
}

func Group(name string) ([]*decode.Format, error) {
	resolveOnce.Do(func() {
		if err := resolve(); err != nil {
			panic(err)
		}
	})

	if g, ok := allGroups[name]; ok {
		return g, nil
	}
	return nil, errors.New("no such group")
}

func MustGroup(name string) []*decode.Format {
	if g, err := Group(name); err == nil {
		return g
	} else {
		panic(err)
	}
}

func MustAll(name string) []*decode.Format {
	return MustGroup("all")
}

// func MustAll() []*decode.Format {
// 	// TODO: only once
// 	if err := Resolve(); err != nil {
// 		panic(err)
// 	}

// 	formatSeen := map[string]struct{}{}
// 	var formats []*decode.Format

// 	for _, fs := range allGroups {
// 		for _, f := range fs {
// 			if _, ok := formatSeen[f.Name]; ok {
// 				continue
// 			}
// 			formatSeen[f.Name] = struct{}{}
// 			formats = append(formats, f)

// 		}
// 	}

// 	return formats
// }

func Dot(w io.Writer, formats []*decode.Format) {
	formatSeen := map[string]struct{}{}

	fmt.Fprintf(w, "digraph formats {\n")

	for _, f := range formats {
		// if groupName != f.Name {
		// 	fmt.Fprintf(w, "  %s -> %s\n", groupName, f.Name)
		// }

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

	fmt.Fprintf(w, "}\n")
}
