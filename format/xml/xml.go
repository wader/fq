package xml

// object mode inspired by https://www.xml.com/pub/a/2006/05/31/converting-between-xml-and-json.html

// TODO: keep <?xml>? root #desc?
// TODO: xml default indent?

import (
	"bytes"
	"embed"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/wader/fq/format"
	"github.com/wader/fq/internal/gojqex"
	"github.com/wader/fq/internal/proxysort"
	"github.com/wader/fq/internal/stringsex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

//go:embed xml.jq
var xmlFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.XML,
		Description: "Extensible Markup Language",
		ProbeOrder:  format.ProbeOrderTextFuzzy,
		Groups:      []string{format.PROBE},
		DecodeFn:    decodeXML,
		DecodeInArg: format.XMLIn{
			Seq:   false,
			Array: false,
		},
		Functions: []string{"_todisplay"},
	})
	interp.RegisterFS(xmlFS)
	interp.RegisterFunc1("toxml", toXML)
	interp.RegisterFunc0("fromxmlentities", func(_ *interp.Interp, c string) any {
		return html.UnescapeString(c)
	})
	interp.RegisterFunc0("toxmlentities", func(_ *interp.Interp, c string) any {
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

func elmName(space, local string) string {
	if space == "" {
		return local
	}
	return space + ":" + local
}

// xmlNNStack is used to undo namespace url resolving, space is url not the "alias" name
type xmlNNStack []xmlNS

func (nss xmlNNStack) lookup(name xml.Name) string {
	for i := len(nss) - 1; i >= 0; i-- {
		ns := nss[i]
		if name.Space == ns.url {
			return ns.name
		}
	}
	return ""
}

func (nss xmlNNStack) push(name string, url string) xmlNNStack {
	n := append([]xmlNS{}, nss...)
	n = append(n, xmlNS{name: name, url: url})
	return xmlNNStack(n)
}

func fromXMLToArray(n xmlNode) any {
	var f func(n xmlNode, nss xmlNNStack) []any
	f = func(n xmlNode, nss xmlNNStack) []any {
		attrs := map[string]any{}
		for _, a := range n.Attrs {
			local, space := a.Name.Local, a.Name.Space
			name := local
			if space != "" {
				if space == "xmlns" {
					nss = nss.push(local, a.Value)
				} else {
					space = nss.lookup(a.Name)
				}
				name = space + ":" + local
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

		local, space := n.XMLName.Local, n.XMLName.Space
		if space != "" {
			space = nss.lookup(n.XMLName)
		}
		// only add if ns is found and not default ns
		name := elmName(space, local)
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

func fromXMLToObject(n xmlNode, xi format.XMLIn) any {
	var f func(n xmlNode, seq int, nss xmlNNStack) any
	f = func(n xmlNode, seq int, nss xmlNNStack) any {
		attrs := map[string]any{}

		for _, a := range n.Attrs {
			local, space := a.Name.Local, a.Name.Space
			name := local
			if space != "" {
				if space == "xmlns" {
					nss = nss.push(local, a.Value)
				} else {
					space = nss.lookup(a.Name)
				}
				name = space + ":" + local
			}
			attrs["-"+name] = a.Value
		}

		for i, nn := range n.Nodes {
			nSeq := i
			if len(n.Nodes) == 1 {
				nSeq = -1
			}
			local, space := nn.XMLName.Local, nn.XMLName.Space
			if space != "" {
				space = nss.lookup(nn.XMLName)
			}
			// only add if ns is found and not default ns
			name := elmName(space, local)
			if e, ok := attrs[name]; ok {
				if ea, ok := e.([]any); ok {
					attrs[name] = append(ea, f(nn, nSeq, nss))
				} else {
					attrs[name] = []any{e, f(nn, nSeq, nss)}
				}
			} else {
				attrs[name] = f(nn, nSeq, nss)
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

		if len(attrs) == 0 {
			return ""
		} else if len(attrs) == 1 && attrs["#text"] != nil {
			return attrs["#text"]
		}

		return attrs
	}

	return map[string]any{
		n.XMLName.Local: f(n, -1, nil),
	}
}

func decodeXML(d *decode.D, in any) any {
	xi, _ := in.(format.XMLIn)

	br := d.RawLen(d.Len())
	var r any
	var err error

	xd := xml.NewDecoder(bitio.NewIOReader(br))
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
	if err != nil {
		d.Fatalf("%s", err)
	}
	var s scalar.S
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
				d.Fatalf("root element has trailing non-whitespace %q", stringsex.TrimN(string(t), 50, "..."))
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

type ToXMLOpts struct {
	Indent int
}

func toXMLFromObject(c any, opts ToXMLOpts) any {
	var f func(name string, content any) (xmlNode, int)
	f = func(name string, content any) (xmlNode, int) {
		n := xmlNode{
			XMLName: xml.Name{Local: name},
		}

		seq := -1
		var orderSeqs []int
		var orderNames []string

		switch v := content.(type) {
		case string:
			n.Chardata = []byte(v)
		case map[string]any:
			for k, v := range v {
				switch {
				case k == "#seq":
					seq, _ = strconv.Atoi(v.(string))
				case k == "#text":
					s, _ := v.(string)
					n.Chardata = []byte(s)
				case k == "#comment":
					s, _ := v.(string)
					n.Comment = []byte(s)
				case strings.HasPrefix(k, "-"):
					s, _ := v.(string)
					n.Attrs = append(n.Attrs, xml.Attr{
						Name:  xml.Name{Local: k[1:]},
						Value: s,
					})
				default:
					switch v := v.(type) {
					case []any:
						if len(v) > 0 {
							for _, c := range v {
								nn, nseq := f(k, c)
								n.Nodes = append(n.Nodes, nn)
								orderNames = append(orderNames, k)
								orderSeqs = append(orderSeqs, nseq)
							}
						} else {
							nn, nseq := f(k, "")
							n.Nodes = append(n.Nodes, nn)
							orderNames = append(orderNames, k)
							orderSeqs = append(orderSeqs, nseq)
						}
					default:
						nn, nseq := f(k, v)
						n.Nodes = append(n.Nodes, nn)
						orderNames = append(orderNames, k)
						orderSeqs = append(orderSeqs, nseq)
					}
				}
			}
		}

		// if one #seq was found, assume all have them, otherwise sort by name
		if len(orderSeqs) > 0 && orderSeqs[0] != -1 {
			proxysort.Sort(orderSeqs, n.Nodes, func(ss []int, i, j int) bool { return ss[i] < ss[j] })
		} else {
			proxysort.Sort(orderNames, n.Nodes, func(ss []string, i, j int) bool { return ss[i] < ss[j] })
		}

		sort.Slice(n.Attrs, func(i, j int) bool {
			a, b := n.Attrs[i].Name, n.Attrs[j].Name
			return a.Space < b.Space || a.Local < b.Local
		})

		return n, seq
	}

	n, _ := f("doc", c)
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
			XMLName: xml.Name{Local: name},
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
					Name:  xml.Name{Local: k},
					Value: s,
				})
			}
		}

		sort.Slice(n.Attrs, func(i, j int) bool {
			a, b := n.Attrs[i].Name, n.Attrs[j].Name
			return a.Space < b.Space || a.Local < b.Local
		})

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
		return gojqex.FuncTypeError{Name: "toxml", V: c}
	}
	n, ok := f(ca)
	if !ok {
		// TODO: better error
		return gojqex.FuncTypeError{Name: "toxml", V: c}
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
	if v, ok := gojqex.Cast[map[string]any](c); ok {
		return toXMLFromObject(gojqex.NormalizeToStrings(v), opts)
	} else if v, ok := gojqex.Cast[[]any](c); ok {
		return toXMLFromArray(gojqex.NormalizeToStrings(v), opts)
	}
	return gojqex.FuncTypeError{Name: "toxml", V: c}
}
