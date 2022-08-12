package main

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

func findDefintion(docs []Documentation) (string, bool) {
	for _, d := range docs {
		if d.Purpose == "definition" {
			return strings.TrimSpace(d.Value), true
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

var symLowerRE = regexp.MustCompile(`[^a-z0-9]+`)

var newLineRE = regexp.MustCompile(`\n`)

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
	scalarPkgPath := os.Args[4]
	root := os.Args[5]

	fmt.Printf("// Code below generated from %s\n", xmlPath)
	fmt.Printf("//nolint:revive\n")
	fmt.Printf("package %s\n", pkgName)
	fmt.Printf("import (\n")
	fmt.Printf("\t%q\n", ebmlPkgPath)
	fmt.Printf("\t%q\n", scalarPkgPath)
	fmt.Printf(")\n")

	fmt.Printf("var Root = ebml.Tag{\n")
	fmt.Printf("\tebml.HeaderID: {Name: \"ebml\", Type: ebml.Master, Tag: ebml.Header},\n")
	fmt.Printf("\t%sID: {Name: \"%s\", Type: ebml.Master, Tag: %s},\n", root, camelToSnake(root), root)
	fmt.Printf("}\n")

	xd := xml.NewDecoder(r)
	var es Schema
	if err := xd.Decode(&es); err != nil {
		panic(err)
	}

	fmt.Println("const (")
	for _, e := range es.Elements {
		fmt.Printf("\t%sID = %s\n", e.Name, strings.ToLower(e.ID))
	}
	fmt.Println(")")

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

		fmt.Printf("var %s = ebml.Tag{\n", e.Name)
		for _, c := range children {
			def, defOk := findDefintion(c.Documentations)
			extra := ""
			typ := c.Type
			switch typ {
			case "master":
				extra = ", Tag: " + c.Name
			case "utf-8":
				typ = "UTF8"
			}

			fmt.Printf("\t%sID: {\n", c.Name)
			fmt.Printf("\t\tName: %q,\n", camelToSnake(c.Name))
			if defOk {
				fmt.Printf("\t\tDefinition: %q,\n", newLineRE.ReplaceAllString(def, " "))
			}
			fmt.Printf("\t\tType: ebml.%s%s,\n", title(typ), extra)
			if len(c.Enums) > 0 {
				switch c.Type {
				case "integer":
					fmt.Printf("\t\tIntegerEnums: scalar.SToScalar{\n")
				case "uinteger":
					fmt.Printf("\t\tUintegerEnums: scalar.UToScalar{\n")
				case "string":
					fmt.Printf("\t\tStringEnums: scalar.StrToScalar{\n")
				}

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
						n, _ := strconv.ParseInt(e.Value, 10, 64)
						fmt.Printf("\t\t\t%d:{\n", n)
					case "uinteger":
						n, _ := strconv.ParseUint(e.Value, 10, 64)
						fmt.Printf("\t\t\t%d:{\n", n)
					case "string":
						fmt.Printf("\t\t\t%q:{\n", e.Value)
					}

					labelOk := !strings.ContainsAny(e.Label, "()")

					if labelOk {
						fmt.Printf("\t\t\t\tSym: %q,\n", symLower(e.Label))
					}

					if enumDefOk {
						fmt.Printf("\t\t\t\tDescription: %q,\n", newLineRE.ReplaceAllString(enumDef, " "))
					} else if !labelOk {
						fmt.Printf("\t\t\t\tDescription: %q,\n", newLineRE.ReplaceAllString(e.Label, " "))
					}

					fmt.Printf("\t\t\t},\n")
				}
				fmt.Printf("\t\t},\n")
			}
			fmt.Printf("\t},\n")
		}
		fmt.Printf("}\n")
		fmt.Printf("\n")
	}

}
