package interp

import (
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
)

type previewNode struct {
	name     string
	pos      int64
	count    int
	previews map[string]struct{}
	nodes    map[string]*previewNode
}

func previewValue(v *decode.Value) string {
	switch vv := v.V.(type) {
	case decode.Array:
		return "[]"
	case decode.Struct:
		return v.Description
	case bool:
		if vv {
			return "true"
		}
		return "false"
	case int64:
		// TODO: DisplayFormat is weird
		return num.PadFormatInt(vv, decode.DisplayFormatToBase(v.DisplayFormat), true, 0)
	case uint64:
		return num.PadFormatUint(vv, decode.DisplayFormatToBase(v.DisplayFormat), true, 0)
	case float64:
		// TODO: float32? better truncated to significant digits?
		return strconv.FormatFloat(vv, 'g', -1, 64)
	case string:
		if len(vv) > 50 {
			return fmt.Sprintf("%q", vv[0:50]) + "..."
		}
		return fmt.Sprintf("%q", vv)
	case []byte:
		if len(vv) > 16 {
			return hex.EncodeToString(vv[0:16]) + "..."
		}
		return hex.EncodeToString(vv)
	case *bitio.Buffer:
		vvLen := vv.Len()
		if vvLen > 16*8 {
			bs, _ := vv.BytesRange(0, 16)
			return hex.EncodeToString(bs) + "..."
		}
		bs, _ := vv.BytesRange(0, int(bitio.BitsByteCount(vvLen)))
		return hex.EncodeToString(bs)
	case nil:
		return "none"

	case []interface{}:
		// TODO: remove?
		return "json []"
	case map[string]interface{}:
		// TODO: remove?
		return "json {}"
	default:
		panic("unreachable")
	}
}

func preview(v *decode.Value, w io.Writer, _ Options) error {
	switch v.V.(type) {
	case decode.Struct:
		sn := &previewNode{name: ".", count: -1, previews: map[string]struct{}{}, nodes: map[string]*previewNode{}}
		err := previewEx(v, sn, 0)
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
		fmt.Fprintln(w, previewValue(v))
	}
	return nil
}

func previewEx(v *decode.Value, n *previewNode, depth int) error {

	if _, ok := v.V.(decode.Array); !ok {
		p := previewValue(v)
		if p != "" {
			if _, ok := n.previews[p]; !ok {
				n.previews[p] = struct{}{}
			}
		}
	}

	switch vv := v.V.(type) {
	case decode.Struct:
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
				if a, ok := sv.V.(decode.Array); ok {
					sn.count = len(a)
				}
				n.nodes[sv.Name] = sn
			}

			err := previewEx(sv, sn, depth+1)
			if err != nil {
				return err
			}
		}
	case decode.Array:
		for _, av := range vv {
			err := previewEx(av, n, depth+1)
			_ = n
			if err != nil {
				return err
			}
		}
	}

	return nil
}
