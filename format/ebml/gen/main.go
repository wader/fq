package main

// TODO: cleanup this mess

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Schema struct {
	Elements []Element `xml:"element"`
}

// <element name="FieldOrder" path="\Segment\Tracks\TrackEntry\Video\FieldOrder" id="0x9D" type="uinteger" minver="4" range="0-14" default="2" minOccurs="1" maxOccurs="1">
// 	<documentation lang="en" purpose="definition">Declare the field ordering of the video. If FlagInterlaced is not set to 1, this Element MUST be ignored.</documentation>
// 	<restriction>
// 		<enum value="0" label="progressive"/>
// 		<enum value="1" label="tff">
// 			<documentation lang="en" purpose="definition">Top field displayed first. Top field stored first.</documentation>
// 		</enum>
// 		<enum value="2" label="undetermined"/>
// 		<enum value="6" label="bff">
// 			<documentation lang="en" purpose="definition">Bottom field displayed first. Bottom field stored first.</documentation>
// 		</enum>
// 		<enum value="9" label="bff(swapped)">
// 			<documentation lang="en" purpose="definition">Top field displayed first. Fields are interleaved in storage with the top line of the top field stored first.</documentation>
// 		</enum>
// 		<enum value="14" label="tff(swapped)">
// 			<documentation lang="en" purpose="definition">Bottom field displayed first. Fields are interleaved in storage with the top line of the top field stored first.</documentation>
// 		</enum>
// 	</restriction>
// 	<extension webm="0"/>
// 	<extension cppname="VideoFieldOrder"/>
// </element>

type Enum struct {
	Value          string          `xml:"value,attr"`
	Label          string          `xml:"label,attr"`
	Documentations []Documentation `xml:"documentation"`
}

type Element struct {
	Name           string          `xml:"name,attr"`
	Path           string          `xml:"path,attr"`
	ID             string          `xml:"id,attr"`
	Type           string          `xml:"type,attr"`
	Range          string          `xml:"range,attr"`
	Default        string          `xml:"default,attr"`
	MinOccurs      string          `xml:"minOccurs,attr"`
	MaxOccurs      string          `xml:"maxOccurs,attr"`
	Length         string          `xml:"length,attr"`
	Documentations []Documentation `xml:"documentation"`
	Enums          []Enum          `xml:"restriction>enum"`
}

// <documentation lang="en" purpose="definition">A randomly generated unique ID to identify the Segment amongst many others (128 bits).</documentation>
type Documentation struct {
	Purpose string `xml:"purpose,attr"`
	Value   string `xml:",cdata"`
}

var symLowerRE = regexp.MustCompile(`[^a-z0-9]+`)
var newLineRE = regexp.MustCompile(`\n`)
var doubleParanRE = regexp.MustCompile(`\(\(.+?\)\)`)
var refsRE = regexp.MustCompile(`\[@[?!](.+?)\]`)
var longParanRE = regexp.MustCompile(`\(.{20,}?\)`)
var whitespaceRE = regexp.MustCompile(`\s+`)
var quotesRE = regexp.MustCompile("`")

func findDefintion(docs []Documentation) (string, bool) {
	for _, d := range docs {
		if d.Purpose == "definition" {
			s := d.Value
			s = doubleParanRE.ReplaceAllLiteralString(s, "")
			s = longParanRE.ReplaceAllLiteralString(s, "")
			s = refsRE.ReplaceAllString(s, "$1")
			s = whitespaceRE.ReplaceAllLiteralString(s, " ")
			s = quotesRE.ReplaceAllLiteralString(s, "")
			s = strings.TrimRight(s, " .")

			if i := strings.IndexAny(s, ".,;"); i != -1 {
				s = s[0:i]
			}

			return s, true
		}
	}
	return "", false
}

func title(s string) string {
	if len(s) <= 1 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func symLower(s string) string {
	s = strings.ToLower(s)
	return symLowerRE.ReplaceAllStringFunc(s, func(s string) string { return "_" })
}

var camelToSnakeRe = regexp.MustCompile(`[[:lower:]][[:upper:]]`)

// "AaaBbb" -> "aaa_bbb"
func camelToSnake(s string) string {
	return strings.ToLower(camelToSnakeRe.ReplaceAllStringFunc(s, func(s string) string {
		return s[0:1] + "_" + s[1:2]
	}))
}

func main() {
	xmlPath := os.Args[1]
	r, err := os.Open(xmlPath)
	if err != nil {
		panic(err)
	}
	pkgName := os.Args[2]
	ebmlPkgPath := os.Args[3]
	root := os.Args[4]

	fmt.Printf("// Code below generated from %s\n", xmlPath)
	fmt.Printf("package %s\n", pkgName)
	fmt.Printf("import (\n")
	fmt.Printf("  %q\n", ebmlPkgPath)
	fmt.Printf(")\n")

	fmt.Printf("var RootElement = &ebml.Master{\n")
	fmt.Printf("  ElementType: ebml.ElementType{\n")
	fmt.Printf("    ID: RootID,\n")
	fmt.Printf("    ParentID: -1,\n")
	fmt.Printf("    Name: \"\",\n")
	fmt.Printf("  },\n")
	fmt.Printf("  Master: map[ebml.ID]ebml.Element{\n")
	fmt.Printf("    ebml.HeaderID: ebml.Header,\n")
	fmt.Printf("    %sID: %sElement,\n", root, root)
	fmt.Printf("   },\n")
	fmt.Printf("}\n")

	xd := xml.NewDecoder(r)
	var es Schema
	if err := xd.Decode(&es); err != nil {
		panic(err)
	}

	fmt.Println("const (")
	fmt.Printf("  RootID = ebml.RootID\n")
	for _, e := range es.Elements {
		fmt.Printf("  %sID = %s\n", e.Name, strings.ToLower(e.ID))
	}
	fmt.Println(")")

	var names []string
	names = append(names, "Root")

	for _, e := range es.Elements {
		var children []Element
		for _, c := range es.Elements {
			suffix := strings.TrimPrefix(c.Path, e.Path+"\\")
			if suffix == "" || strings.Contains(suffix, `\`) {
				continue
			}
			children = append(children, c)
		}
		if len(children) == 0 {
			continue
		}

		var parent Element
		parentPath := e.Path[0:strings.LastIndex(e.Path, `\`)]
		for _, c := range es.Elements {
			if c.Path == parentPath {
				parent = c
				break
			}
		}

		names = append(names, e.Name)
		fmt.Printf("var %sElement = &ebml.Master{\n", e.Name)
		fmt.Printf("  ElementType: ebml.ElementType{\n")
		fmt.Printf("    ID: %sID,\n", e.Name)
		if parent.Name != "" {
			fmt.Printf("    ParentID: %sID,\n", parent.Name)
		} else {
			fmt.Printf("    ParentID: RootID,\n")
		}
		fmt.Printf("    Name: %q,\n", camelToSnake(e.Name))
		if def, defOk := findDefintion(e.Documentations); defOk {
			fmt.Printf("    Definition: %q,\n", newLineRE.ReplaceAllString(def, " "))
		}
		fmt.Printf("  },\n")
		fmt.Printf("  Master: map[ebml.ID]ebml.Element{\n")
		for _, c := range children {
			fmt.Printf("    %sID: %sElement,\n", c.Name, c.Name)
		}
		fmt.Printf("  },\n")
		fmt.Printf("}\n")

		for _, c := range children {
			if c.Type == "master" {
				continue
			}

			typ := c.Type
			switch typ {
			case "utf-8":
				typ = "UTF8"
			}

			enumType := "struct{}"
			switch typ {
			case "integer":
				enumType = "int64"
			case "uinteger":
				enumType = "uint64"
			case "string":
				enumType = "string"
			}

			names = append(names, c.Name)
			fmt.Printf("var %sElement = &ebml.%s{\n", c.Name, title(typ))
			fmt.Printf("  ElementType: ebml.ElementType{\n")
			fmt.Printf("    ID: %sID,\n", c.Name)
			fmt.Printf("    ParentID: %sID,\n", e.Name)
			fmt.Printf("    Name: %q,\n", camelToSnake(c.Name))
			def, defOk := findDefintion(c.Documentations)
			if defOk {
				fmt.Printf("  Definition: %q,\n", newLineRE.ReplaceAllString(def, " "))
			}
			fmt.Printf("  },\n")
			if len(c.Enums) > 0 {
				fmt.Printf("  Enums: map[%s]ebml.Enum{\n", enumType)

				// matroska.xml has dup keys (e.g. PARTS)
				enumDups := map[string]bool{}

				for _, e := range c.Enums {
					if _, ok := enumDups[e.Value]; ok {
						continue
					}
					enumDups[e.Value] = true

					enumDef, enumDefOk := findDefintion(e.Documentations)

					switch c.Type {
					case "integer":
						n, _ := strconv.ParseInt(e.Value, 0, 64)
						fmt.Printf("    %d:{", n)
					case "uinteger":
						n, _ := strconv.ParseUint(e.Value, 0, 64)
						fmt.Printf("    %d:{", n)
					case "string":
						fmt.Printf("    %q:{", e.Value)
					}

					labelOk := !strings.ContainsAny(e.Label, "()")

					if labelOk {
						fmt.Printf("      Name: %q,", symLower(e.Label))
					}

					if enumDefOk {
						fmt.Printf("      Description: %q,", newLineRE.ReplaceAllString(enumDef, " "))
					} else if !labelOk {
						fmt.Printf("      Description: %q,", newLineRE.ReplaceAllString(e.Label, " "))
					}

					fmt.Printf("    },\n")
				}
				fmt.Printf("  },\n")
			}

			fmt.Printf("}\n")
		}
		fmt.Printf("\n")
	}

	fmt.Println("var IDToElement = map[ebml.ID]ebml.Element{")
	for _, n := range names {
		fmt.Printf("  %sID: %sElement,\n", n, n)
	}

	fmt.Println("}")
}
