package xml

import (
	"embed"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/lazyre"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
	"golang.org/x/net/html"
)

//go:embed html.jq
//go:embed html.md
var htmlFS embed.FS

func init() {
	interp.RegisterFormat(
		format.HTML,
		&decode.Format{
			Description: "HyperText Markup Language",
			ProbeOrder:  format.ProbeOrderTextFuzzy,
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeHTML,
			DefaultInArg: format.HTML_In{
				Seq:             false,
				Array:           false,
				AttributePrefix: "@",
			},
			Functions: []string{"_todisplay"},
		})
	interp.RegisterFS(htmlFS)
}

func fromHTMLToObject(n *html.Node, hi format.HTML_In) any {
	var f func(n *html.Node, seq int) any
	f = func(n *html.Node, seq int) any {
		attrs := map[string]any{}

		switch n.Type {
		case html.ElementNode:
			for _, a := range n.Attr {
				attrs[hi.AttributePrefix+a.Key] = a.Val
			}
		default:
			// skip
		}

		nNodes := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				nNodes++
			}
		}
		nSeq := -1
		if nNodes > 1 {
			nSeq = 0
		}

		var textSb *strings.Builder
		var commentSb *strings.Builder

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			switch c.Type {
			case html.ElementNode:
				if e, ok := attrs[c.Data]; ok {
					if ea, ok := e.([]any); ok {
						attrs[c.Data] = append(ea, f(c, nSeq))
					} else {
						attrs[c.Data] = []any{e, f(c, nSeq)}
					}
				} else {
					attrs[c.Data] = f(c, nSeq)
				}
				if nNodes > 1 {
					nSeq++
				}
			case html.TextNode:
				if !whitespaceRE.MatchString(c.Data) {
					if textSb == nil {
						textSb = &strings.Builder{}
					}
					textSb.WriteString(c.Data)
				}
			case html.CommentNode:
				if !whitespaceRE.MatchString(c.Data) {
					if commentSb == nil {
						commentSb = &strings.Builder{}
					}
					commentSb.WriteString(c.Data)
				}
			default:
				// skip other nodes
			}

			if textSb != nil {
				attrs["#text"] = strings.TrimSpace(textSb.String())
			}
			if commentSb != nil {
				attrs["#comment"] = strings.TrimSpace(commentSb.String())
			}
		}

		if hi.Seq && seq != -1 {
			attrs["#seq"] = seq
		}

		if len(attrs) == 0 {
			return ""
		} else if len(attrs) == 1 && attrs["#text"] != nil {
			return attrs["#text"]
		}

		return attrs
	}

	return f(n, -1)
}

func fromHTMLToArray(n *html.Node) any {
	var f func(n *html.Node) any
	f = func(n *html.Node) any {
		attrs := map[string]any{}

		switch n.Type {
		case html.ElementNode:
			for _, a := range n.Attr {
				attrs[a.Key] = a.Val
			}
		default:
			// skip
		}

		nodes := []any{}
		var textSb *strings.Builder
		var commentSb *strings.Builder

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			switch c.Type {
			case html.ElementNode:
				nodes = append(nodes, f(c))
			case html.TextNode:
				if !whitespaceRE.MatchString(c.Data) {
					if textSb == nil {
						textSb = &strings.Builder{}
					}
					textSb.WriteString(c.Data)
				}
			case html.CommentNode:
				if !whitespaceRE.MatchString(c.Data) {
					if commentSb == nil {
						commentSb = &strings.Builder{}
					}
					commentSb.WriteString(c.Data)
				}
			default:
				// skip other nodes
			}
		}

		if textSb != nil {
			attrs["#text"] = strings.TrimSpace(textSb.String())
		}
		if commentSb != nil {
			attrs["#comment"] = strings.TrimSpace(commentSb.String())
		}

		elm := []any{n.Data}
		if len(attrs) > 0 {
			elm = append(elm, attrs)
		} else {
			// make attrs null if there were none, jq allows index into null
			elm = append(elm, nil)
		}
		elm = append(elm, nodes)

		return elm
	}

	// find first element node, skip doctype etc, should be a "html" element
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			n = c
			break
		}
	}

	// should not happen
	if n == nil {
		panic("unreachable")
	}

	return f(n)
}

var htmlMagicRE = &lazyre.RE{S: `` +
	`^` + // anchor to start
	`(?i)` + // case insensitive
	`[[:graph:][:space:]]{0,64}?` + // 0-64 non-control ASCII lazily to allow comment etc
	`(?:` +
	`<\s{0,20}html|` + // <html
	// or
	`<!DOCTYPE\s{1,20}html` + // <!DOCTYPE html
	`)`,
}

func decodeHTML(d *decode.D) any {
	var hi format.HTML_In
	var pi format.Probe_In
	d.ArgAs(&hi)
	if d.ArgAs(&pi) {
		// if probing the input has to start with "<html" or "<!DOCTYPE html" this
		// is because the html parser will always succeed so we have to be careful
		if d.RE(htmlMagicRE.Must()) == nil {
			d.Fatalf("no <html> or <!DOCTYPE html> found")
		}
	}

	br := d.RawLen(d.Len())
	var r any
	var err error
	// disabled scripting means parse noscript tags etc
	n, err := html.ParseWithOptions(bitio.NewIOReader(br), html.ParseOptionEnableScripting(false))
	if err != nil {
		d.Fatalf("%s", err)
	}

	if hi.Array {
		r = fromHTMLToArray(n)
	} else {
		r = fromHTMLToObject(n, hi)
	}
	if err != nil {
		d.Fatalf("%s", err)
	}
	var s scalar.Any
	s.Actual = r

	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}
