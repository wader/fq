package xml

// object mode inspired by https://www.xml.com/pub/a/2006/05/31/converting-between-xml-and-json.html

// TODO: keep <?xml>? root #desc?
// TODO: refactor to share more code
// TODO: rewrite ns stack

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqx"
	"github.com/wader/fq/internal/sortx"
	"github.com/wader/fq/internal/stringsx"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed xml.jq
//go:embed xml.md
var xmlFS embed.FS

func init() {
	interp.RegisterFormat(
		format.XML,
		&decode.Format{
			Description: "Extensible Markup Language",
			ProbeOrder:  format.ProbeOrderTextFuzzy,
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    decodeXML,
			DefaultInArg: format.XML_In{
				Seq:             false,
				Array:           false,
				AttributePrefix: "@",
			},
			Functions: []string{"_todisplay"},
		})
	interp.RegisterFS(xmlFS)
	interp.RegisterFunc1("to_xml", toXML)
	interp.RegisterFunc0("from_xmlentities", func(_ *interp.Interp, c string) any {
		return html.UnescapeString(c)
	})
	interp.RegisterFunc0("to_xmlentities", func(_ *interp.Interp, c string) any {
		return html.EscapeString(c)
	})
}

var whitespaceRE = regexp.MustCompile(`^\s*$`)

type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",attr"`
	Chardata []byte     `xml:",chardata"`
	Comment  []byte     `xml:",comment"`
	Nodes    []xmlNode  `xml:",any"`
}

func (n *xmlNode) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node xmlNode
	return d.DecodeElement((*node)(n), &start)
}

type xmlNS struct {
	name string
	url  string
}

// TODO: nss pop? attr not stack?

// xmlNNStack is used to undo namespace url resolving, space is url not the "alias" name
type xmlNNStack []xmlNS

func (nss xmlNNStack) lookup(name xml.Name) string {
	var s string
	for i := len(nss) - 1; i >= 0; i-- {
		ns := nss[i]
		if name.Space == ns.url {
			// first found or is default namespace
			if s == "" || ns.name == "" {
				s = ns.name
			}
			if s == "" {
				break
			}
		}
	}
	return s
}

func (nss xmlNNStack) push(name string, url string) xmlNNStack {
	n := append([]xmlNS{}, nss...)
	n = append(n, xmlNS{name: name, url: url})
	return xmlNNStack(n)
}

func elmName(space, local string) string {
	if space == "" {
		return local
	}
	return space + ":" + local
}

func fromXMLToObject(n xmlNode, xi format.XML_In) any {
	var f func(n xmlNode, seq int, nss xmlNNStack) (string, any)
	f = func(n xmlNode, seq int, nss xmlNNStack) (string, any) {
		attrs := map[string]any{}

		for _, a := range n.Attrs {
			local, space := a.Name.Local, a.Name.Space
			if space == "xmlns" {
				nss = nss.push(local, a.Value)
			} else if local == "xmlns" {
				// track default namespace
				nss = nss.push("", a.Value)
			}
		}

		for _, a := range n.Attrs {
			local, space := a.Name.Local, a.Name.Space
			name := local
			if space != "" {
				if space == "xmlns" {
					// nop
				} else {
					space = nss.lookup(a.Name)
				}
				name = elmName(space, local)
			}
			attrs[xi.AttributePrefix+name] = a.Value
		}

		for i, nn := range n.Nodes {
			nSeq := i
			if len(n.Nodes) == 1 {
				nSeq = -1
			}

			nname, naddrs := f(nn, nSeq, nss)

			if e, ok := attrs[nname]; ok {
				if ea, ok := e.([]any); ok {
					attrs[nname] = append(ea, naddrs)
				} else {
					attrs[nname] = []any{e, naddrs}
				}
			} else {
				attrs[nname] = naddrs
			}
		}

		if xi.Seq && seq != -1 {
			attrs["#seq"] = seq
		}
		if attrs["#text"] == nil && !whitespaceRE.Match(n.Chardata) {
			attrs["#text"] = strings.TrimSpace(string(n.Chardata))
		}
		if attrs["#comment"] == nil && !whitespaceRE.Match(n.Comment) {
			attrs["#comment"] = strings.TrimSpace(string(n.Comment))
		}

		local, space := n.XMLName.Local, n.XMLName.Space
		if space != "" {
			space = nss.lookup(n.XMLName)
		}
		name := elmName(space, local)

		if len(attrs) == 0 {
			return name, ""
		} else if len(attrs) == 1 && attrs["#text"] != nil {
			return name, attrs["#text"]
		}

		return name, attrs
	}

	name, attrs := f(n, -1, nil)
	return map[string]any{name: attrs}
}

func fromXMLToArray(n xmlNode) any {
	var f func(n xmlNode, nss xmlNNStack) []any
	f = func(n xmlNode, nss xmlNNStack) []any {
		attrs := map[string]any{}

		for _, a := range n.Attrs {
			local, space := a.Name.Local, a.Name.Space
			if space == "xmlns" {
				nss = nss.push(local, a.Value)
			} else if local == "xmlns" {
				// track default namespace
				nss = nss.push("", a.Value)
			}
		}

		for _, a := range n.Attrs {
			local, space := a.Name.Local, a.Name.Space
			name := local
			if space != "" {
				if space == "xmlns" {
					// nop
				} else {
					space = nss.lookup(a.Name)
				}
				name = elmName(space, local)
			}
			attrs[name] = a.Value
		}

		if attrs["#text"] == nil && !whitespaceRE.Match(n.Chardata) {
			attrs["#text"] = strings.TrimSpace(string(n.Chardata))
		}
		if attrs["#comment"] == nil && !whitespaceRE.Match(n.Comment) {
			attrs["#comment"] = strings.TrimSpace(string(n.Comment))
		}

		nodes := []any{}
		for _, c := range n.Nodes {
			nodes = append(nodes, f(c, nss))
		}

		name := elmName(nss.lookup(n.XMLName), n.XMLName.Local)

		elm := []any{name}
		if len(attrs) > 0 {
			elm = append(elm, attrs)
		} else {
			// make attrs null if there were none, jq allows index into null
			elm = append(elm, nil)
		}
		elm = append(elm, nodes)

		return elm
	}

	return f(n, nil)
}

// from golang encoding/xml, copyright 2009 The Go Authors
// the Char production of https://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xD7FF ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

func decodeXMLSeekFirstValidRune(br io.ReadSeeker) error {
	buf := bufio.NewReader(br)
	r, sz, err := buf.ReadRune()
	if err != nil {
		return err
	}
	if _, err := br.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if r == utf8.RuneError && sz == 1 {
		return fmt.Errorf("invalid UTF-8")
	}
	if !isInCharacterRange(r) {
		return fmt.Errorf("illegal character code %U", r)
	}

	return nil
}

func decodeXML(d *decode.D) any {
	var xi format.XML_In
	d.ArgAs(&xi)

	bbr := d.RawLen(d.Len())
	var r any

	br := bitio.NewIOReadSeeker(bbr)

	// this reimplements same xml rune range validation as ecoding/xml but fails faster
	if err := decodeXMLSeekFirstValidRune(br); err != nil {
		d.Fatalf("%s", err)
	}

	xd := xml.NewDecoder(br)

	xd.Strict = false
	var n xmlNode
	if err := xd.Decode(&n); err != nil {
		d.Fatalf("%s", err)
	}

	if xi.Array {
		r = fromXMLToArray(n)
	} else {
		r = fromXMLToObject(n, xi)
	}
	var s scalar.Any
	s.Actual = r

	switch s.Actual.(type) {
	case map[string]any,
		[]any:
	default:
		d.Fatalf("root not object or array")
	}

	// continue decode to end and make sure there is only things we want to ignore
	for {
		d.SeekAbs(xd.InputOffset() * 8)
		t, err := xd.Token()
		if errors.Is(err, io.EOF) {
			break
		}

		switch t := t.(type) {
		case xml.CharData:
			if !whitespaceRE.Match([]byte(t)) {
				d.Fatalf("root element has trailing non-whitespace %q", stringsx.TrimN(string(t), 50, "..."))
			}
			// ignore trailing whitespace
		case xml.ProcInst:
			// ignore trailing process instructions <?elm?>
		case xml.StartElement:
			d.Fatalf("root element has trailing element <%s>", elmName(t.Name.Space, t.Name.Local))
		default:
			d.Fatalf("root element has trailing data")
		}
	}

	d.Value.V = &s
	d.Value.Range.Len = d.Len()

	return nil
}

func xmlNameFromStr(s string) xml.Name {
	return xml.Name{Local: s}
}

func xmlNameSort(a, b xml.Name) int {
	if a.Space != b.Space {
		if a.Space == "" {
			return 1
		}
		return strings.Compare(a.Space, b.Space)
	}
	return strings.Compare(a.Local, b.Local)
}

type ToXMLOpts struct {
	Indent          int
	AttributePrefix string `default:"@"`
}

func toXMLFromObject(c any, opts ToXMLOpts) any {
	var f func(name string, content any) (xmlNode, int, bool)
	f = func(name string, content any) (xmlNode, int, bool) {
		n := xmlNode{
			XMLName: xml.Name{Local: name},
		}

		hasSeq := false
		seq := 0
		orderHasSeq := false
		var orderSeqs []int
		var orderNames []string

		switch v := content.(type) {
		case string:
			n.Chardata = []byte(v)
		case map[string]any:
			for k, v := range v {
				switch {
				case k == "#seq":
					hasSeq = true
					seq, _ = strconv.Atoi(v.(string))
				case k == "#text":
					s, _ := v.(string)
					n.Chardata = []byte(s)
				case k == "#comment":
					s, _ := v.(string)
					n.Comment = []byte(s)
				case strings.HasPrefix(k, opts.AttributePrefix):
					s, _ := v.(string)
					a := xml.Attr{
						Name:  xmlNameFromStr(k[1:]),
						Value: s,
					}
					n.Attrs = append(n.Attrs, a)
				default:
					switch v := v.(type) {
					case []any:
						if len(v) > 0 {
							for _, c := range v {
								nn, nseq, nHasSeq := f(k, c)
								n.Nodes = append(n.Nodes, nn)
								orderNames = append(orderNames, k)
								orderSeqs = append(orderSeqs, nseq)
								orderHasSeq = orderHasSeq || nHasSeq
							}
						} else {
							nn, nseq, nHasSeq := f(k, "")
							n.Nodes = append(n.Nodes, nn)
							orderNames = append(orderNames, k)
							orderSeqs = append(orderSeqs, nseq)
							orderHasSeq = orderHasSeq || nHasSeq
						}
					default:
						nn, nseq, nHasSeq := f(k, v)
						n.Nodes = append(n.Nodes, nn)
						orderNames = append(orderNames, k)
						orderSeqs = append(orderSeqs, nseq)
						orderHasSeq = orderHasSeq || nHasSeq
					}
				}
			}
		}

		// if one #seq was found, assume all have them, otherwise sort by name
		if orderHasSeq {
			sortx.ProxySort(orderSeqs, n.Nodes, func(a, b int) bool { return a < b })
		} else {
			sortx.ProxySort(orderNames, n.Nodes, func(a, b string) bool { return a < b })
		}

		slices.SortFunc(n.Attrs, func(a, b xml.Attr) int { return xmlNameSort(a.Name, b.Name) })

		return n, seq, hasSeq
	}

	n, _, _ := f("doc", c)
	if len(n.Nodes) == 1 && len(n.Attrs) == 0 && n.Comment == nil && n.Chardata == nil {
		n = n.Nodes[0]
	}

	bb := &bytes.Buffer{}
	e := xml.NewEncoder(bb)
	e.Indent("", strings.Repeat(" ", opts.Indent))
	if err := e.Encode(n); err != nil {
		return err
	}
	if err := e.Flush(); err != nil {
		return err
	}

	return bb.String()
}

// ["elm", {attrs}, [children]] -> <elm attrs...>children...</elm>
func toXMLFromArray(c any, opts ToXMLOpts) any {
	var f func(elm []any) (xmlNode, bool)
	f = func(elm []any) (xmlNode, bool) {
		var name string
		var attrs map[string]any
		var children []any

		for _, v := range elm {
			switch v := v.(type) {
			case string:
				if name == "" {
					name = v
				}
			case map[string]any:
				if attrs == nil {
					attrs = v
				}
			case []any:
				if children == nil {
					children = v
				}
			}
		}

		if name == "" {
			return xmlNode{}, false
		}

		n := xmlNode{
			XMLName: xmlNameFromStr(name),
		}

		for k, v := range attrs {
			switch k {
			case "#comment":
				s, _ := v.(string)
				n.Comment = []byte(s)
			case "#text":
				s, _ := v.(string)
				n.Chardata = []byte(s)
			default:
				s, _ := v.(string)
				n.Attrs = append(n.Attrs, xml.Attr{
					Name:  xmlNameFromStr(k),
					Value: s,
				})
			}
		}

		slices.SortFunc(n.Attrs, func(a, b xml.Attr) int { return xmlNameSort(a.Name, b.Name) })

		for _, c := range children {
			c, ok := c.([]any)
			if !ok {
				continue
			}
			if cn, ok := f(c); ok {
				n.Nodes = append(n.Nodes, cn)
			}
		}

		return n, true
	}

	ca, ok := c.([]any)
	if !ok {
		return gojqx.FuncTypeError{Name: "to_xml", V: c}
	}
	n, ok := f(ca)
	if !ok {
		// TODO: better error
		return gojqx.FuncTypeError{Name: "to_xml", V: c}
	}
	bb := &bytes.Buffer{}
	e := xml.NewEncoder(bb)
	e.Indent("", strings.Repeat(" ", opts.Indent))
	if err := e.Encode(n); err != nil {
		return err
	}
	if err := e.Flush(); err != nil {
		return err
	}

	return bb.String()
}

func toXML(_ *interp.Interp, c any, opts ToXMLOpts) any {
	if v, ok := gojqx.Cast[map[string]any](c); ok {
		return toXMLFromObject(gojqx.NormalizeToStrings(v), opts)
	} else if v, ok := gojqx.Cast[[]any](c); ok {
		return toXMLFromArray(gojqx.NormalizeToStrings(v), opts)
	}
	return gojqx.FuncTypeError{Name: "to_xml", V: c}
}
