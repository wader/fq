package markdown

import (
	"embed"
	"fmt"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed markdown.jq
//go:embed markdown.md
var markdownFS embed.FS

func init() {
	interp.RegisterFormat(
		format.Markdown,
		&decode.Format{
			Description: "Markdown",
			DecodeFn:    decodeMarkdown,
			Functions:   []string{"_todisplay"},
		})
	interp.RegisterFS(markdownFS)
}

func decodeMarkdown(d *decode.D) any {
	b, err := io.ReadAll(bitio.NewIOReader(d.RawLen(d.Len())))
	if err != nil {
		panic(err)
	}

	var s scalar.Any
	s.Actual = node(markdown.Parse(b, nil))
	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}

func stringSlice[T string | []byte](ss []T) []any {
	var vs []any
	for _, e := range ss {
		vs = append(vs, string(e))
	}
	return vs
}

func sliceMap[F, T any](vs []F, fn func(F) T) []T {
	ts := make([]T, len(vs))
	for i, v := range vs {
		ts[i] = fn(v)
	}
	return ts
}

func intSlice[T ~int](ss []T) []any {
	var vs []any
	for _, e := range ss {
		vs = append(vs, e)
	}
	return vs
}

func attr(v map[string]any, attr *ast.Attribute) {
	if attr == nil {
		return
	}

	v["id"] = string(attr.ID)

	var as []any
	for _, a := range attr.Attrs {
		as = append(as, string(a))
	}
	v["attrs"] = as

	var cs []any
	for _, a := range attr.Classes {
		cs = append(cs, string(a))
	}
	v["classes"] = cs
}

func leaf(v map[string]any, typ string, l ast.Leaf) {
	v["type"] = typ
	v["literal"] = string(l.Literal)

	attr(v, l.Attribute)
}

func container(v map[string]any, typ string, c ast.Container) {
	v["type"] = typ
	v["literal"] = string(c.Literal)

	var cs []any
	children := c.GetChildren()
	for _, n := range children {
		cv := node(n)
		if cv != nil {
			cs = append(cs, node(n))
		}
	}
	v["children"] = cs

	attr(v, c.Attribute)
}

func listType(t ast.ListType) []any {
	var vs []any

	if t&ast.ListTypeOrdered == ast.ListTypeOrdered {
		vs = append(vs, "ordered")
	}
	if t%ast.ListTypeOrdered == ast.ListTypeOrdered {
		vs = append(vs, "ordered")
	}
	if t%ast.ListTypeDefinition == ast.ListTypeDefinition {
		vs = append(vs, "definition")
	}
	if t%ast.ListTypeTerm == ast.ListTypeTerm {
		vs = append(vs, "term")
	}
	if t%ast.ListItemContainsBlock == ast.ListItemContainsBlock {
		vs = append(vs, "contains_block")
	}
	if t%ast.ListItemBeginningOfList == ast.ListItemBeginningOfList {
		vs = append(vs, "beginning_of_list")
	}
	if t%ast.ListItemEndOfList == ast.ListItemEndOfList {
		vs = append(vs, "end_of_list")
	}

	return vs
}

func node(n ast.Node) any {
	v := map[string]any{}

	switch n := n.(type) {
	case *ast.Text:
		if n.Leaf.Attribute == nil {
			if len(n.Leaf.Literal) > 0 {
				return string(n.Leaf.Literal)
			}
			// skip
			return nil
		}
	case *ast.Softbreak:
		leaf(v, "softbreak", n.Leaf)
	case *ast.Hardbreak:
		leaf(v, "hardbreak", n.Leaf)
	case *ast.NonBlockingSpace:
		leaf(v, "nbsp", n.Leaf)
	case *ast.Emph:
		container(v, "em", n.Container)
	case *ast.Strong:
		container(v, "strong", n.Container)
	case *ast.Del:
		container(v, "del", n.Container)
	case *ast.BlockQuote:
		container(v, "blockquote", n.Container)
	case *ast.Aside:
		container(v, "aside", n.Container)
	case *ast.Link:
		container(v, "link", n.Container)
		v["destination"] = string(n.Destination)
		v["title"] = string(n.Title)
		v["note_id"] = n.NoteID
		v["deferred_id"] = string(n.DeferredID)
		v["additional_attributes"] = stringSlice(n.AdditionalAttributes)
	case *ast.CrossReference:
		container(v, "cross_reference", n.Container)
		v["destination"] = string(n.Destination)
	case *ast.Citation:
		leaf(v, "citation", n.Leaf)
		v["destination"] = stringSlice(n.Destination)
		v["type"] = sliceMap(n.Type, func(v ast.CitationTypes) string {
			switch v {
			case ast.CitationTypeNone:
				return "none"
			case ast.CitationTypeSuppressed:
				return "suppressed"
			case ast.CitationTypeInformative:
				return "informative"
			case ast.CitationTypeNormative:
				return "normative"
			default:
				return "unknown"
			}
		})
		v["type"] = intSlice(n.Type)
		v["suffix"] = stringSlice(n.Suffix)
	case *ast.Image:
		container(v, "image", n.Container)
		v["destination"] = string(n.Destination)
		v["title"] = string(n.Title)
	case *ast.Code:
		leaf(v, "code", n.Leaf)
	case *ast.CodeBlock:
		leaf(v, "code_block", n.Leaf)
		v["is_fenced"] = n.IsFenced
		v["info"] = string(n.Info)
		if n.FenceChar != 0 {
			v["fence_char"] = string(n.FenceChar)
		}
		v["fence_length"] = n.FenceLength
		v["fence_offset"] = n.FenceOffset
	case *ast.Caption:
		container(v, "caption", n.Container)
	case *ast.CaptionFigure:
		container(v, "caption_figure", n.Container)
		v["heading_id"] = n.HeadingID
	case *ast.Document:
		container(v, "document", n.Container)
	case *ast.Paragraph:
		container(v, "paragraph", n.Container)
	case *ast.HTMLSpan:
		leaf(v, "html_span", n.Leaf)
	case *ast.HTMLBlock:
		leaf(v, "html_block", n.Leaf)
	case *ast.Heading:
		container(v, "heading", n.Container)
		v["level"] = n.Level
		v["heading_id"] = n.HeadingID
		v["is_titleblock"] = n.IsTitleblock
		v["is_special"] = n.IsSpecial
	case *ast.HorizontalRule:
		leaf(v, "hr", n.Leaf)
	case *ast.List:
		container(v, "list", n.Container)
		v["list_flags"] = listType(n.ListFlags)
		v["tight"] = n.Tight
		if n.BulletChar != 0 {
			v["bullet_char"] = string(n.BulletChar)
		}
		if n.Delimiter != 0 {
			v["delimiter"] = string(n.Delimiter)
		}
		v["start"] = n.Start
		v["ref_link"] = string(n.RefLink)
		v["is_footnotes_list"] = n.IsFootnotesList
	case *ast.ListItem:
		container(v, "list_item", n.Container)
		v["list_flags"] = listType(n.ListFlags)
		v["tight"] = n.Tight
		if n.BulletChar != 0 {
			v["bullet_char"] = string(n.BulletChar)
		}
		if n.Delimiter != 0 {
			v["delimiter"] = string(n.Delimiter)
		}
		v["ref_link"] = string(n.RefLink)
		v["is_footnotes_list"] = n.IsFootnotesList
	case *ast.Table:
		container(v, "table", n.Container)
	case *ast.TableCell:
		container(v, "table_cell", n.Container)
		v["is_header"] = n.IsHeader
		v["align"] = n.Align.String()
		v["col_span"] = n.ColSpan
	case *ast.TableHeader:
		container(v, "table_header", n.Container)
	case *ast.TableBody:
		container(v, "table_body", n.Container)
	case *ast.TableRow:
		container(v, "table_row", n.Container)
	case *ast.TableFooter:
		container(v, "table_footer", n.Container)
	case *ast.Math:
		leaf(v, "math", n.Leaf)
	case *ast.MathBlock:
		container(v, "math_block", n.Container)
	case *ast.DocumentMatter:
		container(v, "document_matter", n.Container)
		v["matter"] = func(v ast.DocumentMatters) string {
			switch v {
			case ast.DocumentMatterNone:
				return "none"
			case ast.DocumentMatterFront:
				return "front"
			case ast.DocumentMatterMain:
				return "main"
			case ast.DocumentMatterBack:
				return "back"
			default:
				return "unknown"
			}
		}(n.Matter)
	case *ast.Callout:
		leaf(v, "callout", n.Leaf)
		v["id"] = string(n.ID)
	case *ast.Index:
		leaf(v, "index", n.Leaf)
		v["primary"] = n.Primary
		v["item"] = string(n.Item)
		v["subitem"] = string(n.Subitem)
		v["id"] = n.ID
	case *ast.Subscript:
		leaf(v, "subscript", n.Leaf)
	case *ast.Superscript:
		leaf(v, "superscript", n.Leaf)
	case *ast.Footnotes:
		container(v, "footnotes", n.Container)
	default:
		panic(fmt.Sprintf("unknown node %T", node))
	}

	for k, e := range v {
		if s, ok := e.(string); ok && s == "" {
			delete(v, k)
		}
	}

	return v
}
