package register

import (
	"fmt"
	"fq/pkg/decode"
)

var allFormats = map[string]*decode.Format{}

func Register(format *decode.Format) *decode.Format {
	if _, ok := allFormats[format.Name]; ok {
		panic(fmt.Sprintf("%s: already registered", format.Name))
	}
	allFormats[format.Name] = format
	return format
}

func Resolve() {
	for _, f := range allFormats {
		for _, d := range f.Deps {
			var formats []*decode.Format
			for _, dName := range d.Names {
				df, ok := allFormats[dName]
				if !ok {
					panic(fmt.Sprintf("%s: can't find dependency %s", f.Name, dName))
				}
				formats = append(formats, df)
			}
			*d.Formats = formats
		}
	}
	/*
		fmt.Printf("digraph formats {\n")
		for _, f := range allFormats {
			for _, d := range f.Deps {
				for _, dName := range d.Names {
					fmt.Printf("  %s -> %s\n", f.Name, dName)
				}
			}
		}
		fmt.Printf("}\n")
	*/
}
