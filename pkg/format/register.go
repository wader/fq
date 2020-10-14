package format

import (
	"fmt"
	"fq/pkg/decode"
	"io"
)

// TODO: some struct so it's more flexible

var allGroups = map[string][]*decode.Format{}

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

func MustAll() []*decode.Format {
	// TODO: only once
	if err := Resolve(); err != nil {
		panic(err)
	}

	formatSeen := map[string]struct{}{}
	var formats []*decode.Format

	for _, fs := range allGroups {
		for _, f := range fs {
			if _, ok := formatSeen[f.Name]; ok {
				continue
			}
			formatSeen[f.Name] = struct{}{}
			formats = append(formats, f)

		}
	}

	return formats
}

func Resolve() error {
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

func Dot(w io.Writer) {
	formatSeen := map[string]struct{}{}

	fmt.Fprintf(w, "digraph formats {\n")
	for groupName, fs := range allGroups {
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
