package decode

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Field struct {
	Name     string
	Range    Range
	Value    Value
	Decoder  Decoder
	Children []*Field
}

var lookupRe = regexp.MustCompile(`^([\w_]*)(?:\[(\d+)\])?$`)

func (f *Field) Eval(exp string) (*Field, error) {
	lf := f.Lookup(exp)
	if lf == nil {
		return lf, fmt.Errorf("not found")
	}

	return lf, nil
}

func (f *Field) Lookup(path string) *Field {
	if path == "" {
		return f
	}

	parts := strings.SplitN(path, ".", 2)
	first := parts[0]
	rest := ""
	if len(parts) > 1 {
		rest = parts[1]
	}

	index := 0
	firstSM := lookupRe.FindStringSubmatch(first)
	name := firstSM[1]
	indexStr := firstSM[2]
	if indexStr != "" {
		index, _ = strconv.Atoi(indexStr)
	}

	var indexC = 0
	for _, c := range f.Children {
		if name != "" && c.Name != name {
			continue
		}

		if indexC != index {
			indexC++
			continue
		}

		return c.Lookup(rest)
	}

	return nil
}

func (f *Field) Sort() {
	if len(f.Children) == 0 {
		return
	}

	sort.Slice(f.Children, func(i, j int) bool {
		return f.Children[i].Range.Start < f.Children[j].Range.Start
	})

	for _, fc := range f.Children {
		if fc.Value.Type == TypeDecoder {
			// already sorted
			continue
		}
		fc.Sort()
	}
}
