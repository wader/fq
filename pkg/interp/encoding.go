package interp

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1" //nolint: gosec
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/wader/fq/internal/colorjson"
	"github.com/wader/fq/internal/gojqextra"
	"github.com/wader/fq/internal/mapstruct"
	"github.com/wader/fq/internal/proxysort"
	"github.com/wader/fq/pkg/bitio"

	"golang.org/x/crypto/md4" //nolint: staticcheck
	"golang.org/x/crypto/sha3"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"gopkg.in/yaml.v3"
)

// TODO: Fn1, Fn2 etc?
// TODO: struct arg, own reflect code? no need for refs etc

// TODO: xml default indent?
// TODO: query dup key
// TODO: walk tostring tests
// TODO: tocrc*, to/fromutf*
// TODO: keep <?xml>? root #desc?
// TODO: yaml type eval? walk eval?
// TODO: error messages, typeof?
// TODO: map struct gojq.JQValue

// convert to gojq compatible values
func norm(v any) any {
	switch v := v.(type) {
	case map[string]any:
		for i, e := range v {
			v[i] = norm(e)
		}
		return v
	case map[any]any:
		// for gopkg.in/yaml.v2
		vm := map[string]any{}
		for i, e := range v {
			switch i := i.(type) {
			case string:
				vm[i] = norm(e)
			case int:
				vm[strconv.Itoa(i)] = norm(e)
			}
		}
		return vm
	case []map[string]any:
		var vs []any
		for _, e := range v {
			vs = append(vs, norm(e))
		}
		return vs
	case []any:
		for i, e := range v {
			v[i] = norm(e)
		}
		return v
	default:
		v, _ = gojqextra.ToGoJQValue(v)
		return v
	}
}

func addFunc[Tc any](name string, fn func(c Tc) any) {
	if name[0] != '_' {
		panic(fmt.Sprintf("invalid addFunc name %q", name))
	}
	functionRegisterFns = append(
		functionRegisterFns,
		func(i *Interp) []Function {
			return []Function{{
				name, 0, 0, func(c any, a []any) any {
					cv, ok := gojqextra.CastFn[Tc](c, mapstruct.ToStruct)
					if !ok {
						return gojqextra.FuncTypeError{Name: name[1:], V: c}
					}

					return fn(cv)
				},
				nil,
			}}
		})
}

func addFunc1[Tc any, Ta0 any](name string, fn func(c Tc, a0 Ta0) any) {
	if name[0] != '_' {
		panic(fmt.Sprintf("invalid addFunc name %q", name))
	}
	functionRegisterFns = append(
		functionRegisterFns,
		func(i *Interp) []Function {
			return []Function{{
				name, 1, 1, func(c any, a []any) any {
					cv, ok := gojqextra.CastFn[Tc](c, mapstruct.ToStruct)
					if !ok {
						return gojqextra.FuncTypeError{Name: name[1:], V: c}
					}
					a0, ok := gojqextra.CastFn[Ta0](a[0], mapstruct.ToStruct)
					if !ok {
						return gojqextra.FuncArgTypeError{Name: name[1:], ArgName: "first", V: c}
					}

					return fn(cv, a0)
				},
				nil,
			}}
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

// xmlNNStack is used to undo namespace url resolving, space is url not the "alias" name
type xmlNNStack []xmlNS

func (nss xmlNNStack) lookup(name xml.Name) string {
	for i := len(nss) - 1; i >= 0; i-- {
		ns := nss[i]
		if name.Space == ns.url {
			return ns.name
		}
	}
	return name.Space
}

func (nss xmlNNStack) push(name string, url string) xmlNNStack {
	n := append([]xmlNS{}, nss...)
	n = append(n, xmlNS{name: name, url: url})
	return xmlNNStack(n)
}

func init() {
	type ToJSONOpts struct {
		Indent int
	}
	addFunc1("_tojson", func(c any, opts ToJSONOpts) any {
		// TODO: share
		cj := colorjson.NewEncoder(
			false,
			false,
			opts.Indent,
			func(v any) any {
				if v, ok := toValue(nil, v); ok {
					return v
				}
				panic(fmt.Sprintf("toValue not a JQValue value: %#v %T", v, v))
			},
			colorjson.Colors{},
		)
		bb := &bytes.Buffer{}
		if err := cj.Marshal(c, bb); err != nil {
			return err
		}
		return bb.String()
	})

	type FromXMLOpts struct {
		Seq   bool
		Array bool
	}
	fromXMLObject := func(s string, opts FromXMLOpts) any {
		d := xml.NewDecoder(strings.NewReader(s))
		d.Strict = false
		var n xmlNode
		if err := d.Decode(&n); err != nil {
			return err
		}

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
				name := local
				if space != "" {
					space = nss.lookup(nn.XMLName)
					name = space + ":" + name
				}
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

			if opts.Seq && seq != -1 {
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
	fromXMLArray := func(s string) any {
		d := xml.NewDecoder(strings.NewReader(s))
		d.Strict = false
		var n xmlNode
		if err := d.Decode(&n); err != nil {
			return err
		}

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

			name, space := n.XMLName.Local, n.XMLName.Space
			if space != "" {
				space = nss.lookup(n.XMLName)
				name = space + ":" + name
			}
			elm := []any{name}
			if len(attrs) > 0 {
				elm = append(elm, attrs)
			}
			if len(nodes) > 0 {
				elm = append(elm, nodes)
			}

			return elm
		}

		return f(n, nil)
	}
	addFunc1("_fromxml", func(s string, opts FromXMLOpts) any {
		if opts.Array {
			return fromXMLArray(s)
		}
		return fromXMLObject(s, opts)
	})

	type ToXMLOpts struct {
		Indent int
	}
	toXMLObject := func(c map[string]any, opts ToXMLOpts) any {
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
	toXMLArray := func(c []any, opts ToXMLOpts) any {
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

		n, ok := f(c)
		if !ok {
			// TODO: error
			return nil
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
	addFunc1("_toxml", func(c any, opts ToXMLOpts) any {
		switch c := c.(type) {
		case map[string]any:
			return toXMLObject(c, opts)
		case []any:
			return toXMLArray(c, opts)
		default:
			return gojqextra.FuncTypeError{Name: "toxml", V: c}
		}
	})

	type FromHTMLOpts struct {
		Seq   bool
		Array bool
	}
	fromHTMLObject := func(s string, opts FromHTMLOpts) any {
		doc, err := html.Parse(bytes.NewBuffer([]byte(s)))
		if err != nil {
			return err
		}

		var f func(n *html.Node, seq int) any
		f = func(n *html.Node, seq int) any {
			attrs := map[string]any{}

			switch n.Type {
			case html.ElementNode:
				for _, a := range n.Attr {
					attrs["-"+a.Key] = a.Val
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

			if opts.Seq && seq != -1 {
				attrs["#seq"] = seq
			}

			if len(attrs) == 0 {
				return ""
			} else if len(attrs) == 1 && attrs["#text"] != nil {
				return attrs["#text"]
			}

			return attrs
		}

		return f(doc, -1)
	}
	fromHTMLArray := func(s string) any {
		doc, err := html.Parse(bytes.NewBuffer([]byte(s)))
		if err != nil {
			return err
		}

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
			}
			if len(nodes) > 0 {
				elm = append(elm, nodes)
			}

			return elm
		}

		return f(doc.FirstChild)
	}
	addFunc1("_fromhtml", func(s string, opts FromHTMLOpts) any {
		if opts.Array {
			return fromHTMLArray(s)
		}
		return fromHTMLObject(s, opts)
	})

	addFunc("_fromyaml", func(s string) any {
		var t any
		if err := yaml.Unmarshal([]byte(s), &t); err != nil {
			return err
		}
		return norm(t)
	})

	addFunc("_toyaml", func(c any) any {
		b, err := yaml.Marshal(norm(c))
		if err != nil {
			return err
		}
		return string(b)
	})

	addFunc("_fromtoml", func(s string) any {
		var t any
		if err := toml.Unmarshal([]byte(s), &t); err != nil {
			return err
		}
		return norm(t)
	})

	addFunc("_totoml", func(c map[string]any) any {
		b := &bytes.Buffer{}
		if err := toml.NewEncoder(b).Encode(norm(c)); err != nil {
			return err
		}
		return b.String()
	})

	type FromCSVOpts struct {
		Comma   string
		Comment string
	}
	addFunc1("_fromcsv", func(s string, opts FromCSVOpts) any {
		var rvs []any
		r := csv.NewReader(strings.NewReader(s))
		r.TrimLeadingSpace = true
		r.LazyQuotes = true
		if opts.Comma != "" {
			r.Comma = rune(opts.Comma[0])
		}
		if opts.Comment != "" {
			r.Comment = rune(opts.Comment[0])
		}
		for {
			r, err := r.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			var vs []any
			for _, s := range r {
				vs = append(vs, s)
			}
			rvs = append(rvs, vs)
		}
		return rvs
	})

	type ToCSVOpts struct {
		Comma string
	}
	addFunc1("_tocsv", func(c []any, opts ToCSVOpts) any {
		b := &bytes.Buffer{}
		w := csv.NewWriter(b)
		if opts.Comma != "" {
			w.Comma = rune(opts.Comma[0])
		}
		for _, row := range c {
			var ss []string
			rs, ok := row.([]any)
			if !ok {
				return fmt.Errorf("expected row to be an array, got %s", gojqextra.TypeErrorPreview(row))
			}
			for _, v := range rs {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("expected record to be a string, got %s", gojqextra.TypeErrorPreview(v))
				}
				ss = append(ss, s)
			}
			if err := w.Write(ss); err != nil {
				return err
			}
		}
		w.Flush()

		return b.String()
	})

	addFunc("_fromhex", func(s string) any {
		b, err := hex.DecodeString(s)
		if err != nil {
			return err
		}
		bb, err := newBinaryFromBitReader(bitio.NewBitReader(b, -1), 8, 0)
		if err != nil {
			return err
		}
		return bb
	})
	addFunc("_tohex", func(c any) any {
		br, err := toBitReader(c)
		if err != nil {
			return err
		}
		buf := &bytes.Buffer{}
		if _, err := io.Copy(hex.NewEncoder(buf), bitio.NewIOReader(br)); err != nil {
			return err
		}
		return buf.String()
	})

	// TODO: other encodings and share?
	base64Encoding := func(enc string) *base64.Encoding {
		switch enc {
		case "url":
			return base64.URLEncoding
		case "rawstd":
			return base64.RawStdEncoding
		case "rawurl":
			return base64.RawURLEncoding
		default:
			return base64.StdEncoding
		}
	}
	type FromBase64Opts struct {
		Encoding string
	}
	addFunc1("_frombase64", func(s string, opts FromBase64Opts) any {
		b, err := base64Encoding(opts.Encoding).DecodeString(s)
		if err != nil {
			return err
		}
		bin, err := newBinaryFromBitReader(bitio.NewBitReader(b, -1), 8, 0)
		if err != nil {
			return err
		}
		return bin
	})
	type ToBase64Opts struct {
		Encoding string
	}
	addFunc1("_tobase64", func(c any, opts ToBase64Opts) any {
		br, err := toBitReader(c)
		if err != nil {
			return err
		}
		bb := &bytes.Buffer{}
		wc := base64.NewEncoder(base64Encoding(opts.Encoding), bb)
		if _, err := io.Copy(wc, bitio.NewIOReader(br)); err != nil {
			return err
		}
		wc.Close()
		return bb.String()
	})

	addFunc("_fromxmlentities", func(s string) any {
		return html.UnescapeString(s)
	})
	addFunc("_toxmlentities", func(s string) any {
		return html.EscapeString(s)
	})

	addFunc("_fromurlencode", func(s string) any {
		u, _ := url.QueryUnescape(s)
		return u
	})
	addFunc("_tourlencode", func(s string) any {
		return url.QueryEscape(s)
	})

	addFunc("_fromurlpath", func(s string) any {
		u, _ := url.PathUnescape(s)
		return u
	})
	addFunc("_tourlpath", func(s string) any {
		return url.PathEscape(s)
	})

	fromURLValues := func(q url.Values) any {
		qm := map[string]any{}
		for k, v := range q {
			if len(v) > 1 {
				vm := []any{}
				for _, v := range v {
					vm = append(vm, v)
				}
				qm[k] = vm
			} else {
				qm[k] = v[0]
			}
		}

		return qm
	}
	addFunc("_fromurlquery", func(s string) any {
		q, err := url.ParseQuery(s)
		if err != nil {
			return err
		}
		return fromURLValues(q)
	})
	toURLValues := func(c map[string]any) url.Values {
		qv := url.Values{}
		for k, v := range c {
			if va, ok := gojqextra.Cast[[]any](v); ok {
				var ss []string
				for _, s := range va {
					if s, ok := gojqextra.Cast[string](s); ok {
						ss = append(ss, s)
					}
				}
				qv[k] = ss
			} else if vs, ok := gojqextra.Cast[string](v); ok {
				qv[k] = []string{vs}
			}
		}
		return qv
	}
	addFunc("_tourlquery", func(c map[string]any) any {
		return toURLValues(c).Encode()
	})

	addFunc("_fromurl", func(s string) any {
		u, err := url.Parse(s)
		if err != nil {
			return err
		}

		m := map[string]any{}
		if u.Scheme != "" {
			m["scheme"] = u.Scheme
		}
		if u.User != nil {
			um := map[string]any{
				"username": u.User.Username(),
			}
			if p, ok := u.User.Password(); ok {
				um["password"] = p
			}
			m["user"] = um
		}
		if u.Host != "" {
			m["host"] = u.Host
		}
		if u.Path != "" {
			m["path"] = u.Path
		}
		if u.RawPath != "" {
			m["rawpath"] = u.RawPath
		}
		if u.RawQuery != "" {
			m["rawquery"] = u.RawQuery
			m["query"] = fromURLValues(u.Query())
		}
		if u.Fragment != "" {
			m["fragment"] = u.Fragment
		}
		return m
	})
	addFunc("_tourl", func(c map[string]any) any {
		str := func(v any) string { s, _ := gojqextra.Cast[string](v); return s }
		u := url.URL{
			Scheme:   str(c["scheme"]),
			Host:     str(c["host"]),
			Path:     str(c["path"]),
			Fragment: str(c["fragment"]),
		}

		if um, ok := gojqextra.Cast[map[string]any](c["user"]); ok {
			username, password := str(um["username"]), str(um["password"])
			if username != "" {
				if password == "" {
					u.User = url.User(username)
				} else {
					u.User = url.UserPassword(username, password)
				}
			}
		}
		if s, ok := gojqextra.Cast[string](c["rawquery"]); ok {
			u.RawQuery = s
		}
		if qm, ok := gojqextra.Cast[map[string]any](c["query"]); ok {
			u.RawQuery = toURLValues(qm).Encode()
		}

		return u.String()
	})

	hashFn := func(s string) hash.Hash {
		switch s {
		case "md4":
			return md4.New()
		case "md5":
			return md5.New()
		case "sha1":
			return sha1.New()
		case "sha256":
			return sha256.New()
		case "sha512":
			return sha512.New()
		case "sha3_224":
			return sha3.New224()
		case "sha3_256":
			return sha3.New256()
		case "sha3_384":
			return sha3.New384()
		case "sha3_512":
			return sha3.New512()
		default:
			return nil
		}
	}
	type ToHashOpts struct {
		Name string
	}
	addFunc1("_tohash", func(c any, opts ToHashOpts) any {
		inBR, err := toBitReader(c)
		if err != nil {
			return err
		}

		h := hashFn(opts.Name)
		if h == nil {
			return fmt.Errorf("unknown hash function %s", opts.Name)
		}
		if _, err := io.Copy(h, bitio.NewIOReader(inBR)); err != nil {
			return err
		}

		outBR := bitio.NewBitReader(h.Sum(nil), -1)

		bb, err := newBinaryFromBitReader(outBR, 8, 0)
		if err != nil {
			return err
		}
		return bb
	})

	strEncodingFn := func(s string) encoding.Encoding {
		switch s {
		case "UTF8":
			return unicode.UTF8
		case "UTF16":
			return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		case "UTF16LE":
			return unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		case "UTF16BE":
			return unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
		case "CodePage037":
			return charmap.CodePage037
		case "CodePage437":
			return charmap.CodePage437
		case "CodePage850":
			return charmap.CodePage850
		case "CodePage852":
			return charmap.CodePage852
		case "CodePage855":
			return charmap.CodePage855
		case "CodePage858":
			return charmap.CodePage858
		case "CodePage860":
			return charmap.CodePage860
		case "CodePage862":
			return charmap.CodePage862
		case "CodePage863":
			return charmap.CodePage863
		case "CodePage865":
			return charmap.CodePage865
		case "CodePage866":
			return charmap.CodePage866
		case "CodePage1047":
			return charmap.CodePage1047
		case "CodePage1140":
			return charmap.CodePage1140
		case "ISO8859_1":
			return charmap.ISO8859_1
		case "ISO8859_2":
			return charmap.ISO8859_2
		case "ISO8859_3":
			return charmap.ISO8859_3
		case "ISO8859_4":
			return charmap.ISO8859_4
		case "ISO8859_5":
			return charmap.ISO8859_5
		case "ISO8859_6":
			return charmap.ISO8859_6
		case "ISO8859_6E":
			return charmap.ISO8859_6E
		case "ISO8859_6I":
			return charmap.ISO8859_6I
		case "ISO8859_7":
			return charmap.ISO8859_7
		case "ISO8859_8":
			return charmap.ISO8859_8
		case "ISO8859_8E":
			return charmap.ISO8859_8E
		case "ISO8859_8I":
			return charmap.ISO8859_8I
		case "ISO8859_9":
			return charmap.ISO8859_9
		case "ISO8859_10":
			return charmap.ISO8859_10
		case "ISO8859_13":
			return charmap.ISO8859_13
		case "ISO8859_14":
			return charmap.ISO8859_14
		case "ISO8859_15":
			return charmap.ISO8859_15
		case "ISO8859_16":
			return charmap.ISO8859_16
		case "KOI8R":
			return charmap.KOI8R
		case "KOI8U":
			return charmap.KOI8U
		case "Macintosh":
			return charmap.Macintosh
		case "MacintoshCyrillic":
			return charmap.MacintoshCyrillic
		case "Windows874":
			return charmap.Windows874
		case "Windows1250":
			return charmap.Windows1250
		case "Windows1251":
			return charmap.Windows1251
		case "Windows1252":
			return charmap.Windows1252
		case "Windows1253":
			return charmap.Windows1253
		case "Windows1254":
			return charmap.Windows1254
		case "Windows1255":
			return charmap.Windows1255
		case "Windows1256":
			return charmap.Windows1256
		case "Windows1257":
			return charmap.Windows1257
		case "Windows1258":
			return charmap.Windows1258
		case "XUserDefined":
			return charmap.XUserDefined
		default:
			return nil
		}
	}
	type ToStrEncodingOpts struct {
		Encoding string
	}
	addFunc1("_tostrencoding", func(c string, opts ToStrEncodingOpts) any {
		h := strEncodingFn(opts.Encoding)
		if h == nil {
			return fmt.Errorf("unknown string encoding %s", opts.Encoding)
		}

		bb := &bytes.Buffer{}
		if _, err := io.Copy(h.NewEncoder().Writer(bb), strings.NewReader(c)); err != nil {
			return err
		}
		outBR := bitio.NewBitReader(bb.Bytes(), -1)
		bin, err := newBinaryFromBitReader(outBR, 8, 0)
		if err != nil {
			return err
		}

		return bin
	})
	type FromStrEncodingOpts struct {
		Encoding string
	}
	addFunc1("_fromstrencoding", func(c any, opts FromStrEncodingOpts) any {
		inBR, err := toBitReader(c)
		if err != nil {
			return err
		}
		h := strEncodingFn(opts.Encoding)
		if h == nil {
			return fmt.Errorf("unknown string encoding %s", opts.Encoding)
		}

		bb := &bytes.Buffer{}
		if _, err := io.Copy(bb, h.NewDecoder().Reader(bitio.NewIOReader(inBR))); err != nil {

			return err
		}

		return bb.String()
	})
}
