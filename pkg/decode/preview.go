package decode

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type previewNode struct {
	name     string
	pos      int64
	count    int
	previews map[string]struct{}
	nodes    map[string]*previewNode
}

func (v *Value) Preview(w io.Writer) error {
	switch v.V.(type) {
	case Struct:
		sn := &previewNode{name: ".", count: -1, previews: map[string]struct{}{}, nodes: map[string]*previewNode{}}
		err := v.preview(sn, 0)
		if err != nil {
			return err
		}

		var printFn func(s *previewNode, depth int)
		printFn = func(s *previewNode, depth int) {
			indent := strings.Repeat("  ", depth)
			fmt.Fprint(w, indent)
			fmt.Fprint(w, s.name)
			if s.count != -1 {
				fmt.Fprintf(w, "[%d]", s.count)
			}

			var sortedPreviews []string
			for p := range s.previews {
				sortedPreviews = append(sortedPreviews, p)
			}
			sort.Strings(sortedPreviews)
			if len(sortedPreviews) > 10 {
				sortedPreviews = sortedPreviews[0:10]
				sortedPreviews = append(sortedPreviews, "...")
			}
			if len(sortedPreviews) > 0 {
				fmt.Fprintf(w, " (%s)", strings.Join(sortedPreviews, ","))
			}
			fmt.Fprintln(w)

			var sortedNodes []*previewNode
			for _, n := range s.nodes {
				sortedNodes = append(sortedNodes, n)
			}
			// sort.Slice(sorted, func(i, j int) bool {
			// 	return sorted[i].name < sorted[j].name
			// })
			sort.Slice(sortedNodes, func(i, j int) bool {
				return sortedNodes[i].pos < sortedNodes[j].pos
			})

			for _, n := range sortedNodes {
				printFn(n, depth+1)
			}
		}

		printFn(sn, 0)
	default:
		fmt.Fprintln(w, v.PreviewString())
	}
	return nil
}

func (v *Value) preview(n *previewNode, depth int) error {

	if _, ok := v.V.(Array); !ok {
		p := v.PreviewString()
		if p != "" {
			if _, ok := n.previews[p]; !ok {
				n.previews[p] = struct{}{}
			}
		}
	}

	switch vv := v.V.(type) {
	case Struct:
		for _, sv := range vv {
			sn, ok := n.nodes[sv.Name]
			if !ok {
				sn = &previewNode{
					name:     sv.Name,
					pos:      sv.Range.Start,
					count:    -1,
					previews: map[string]struct{}{},
					nodes:    map[string]*previewNode{},
				}
				if a, ok := sv.V.(Array); ok {
					sn.count = len(a)
				}
				n.nodes[sv.Name] = sn
			}

			err := sv.preview(sn, depth+1)
			if err != nil {
				return nil
			}
		}
	case Array:
		for _, av := range vv {
			err := av.preview(n, depth+1)
			_ = n
			if err != nil {
				return err
			}
		}
	}

	return nil
}
