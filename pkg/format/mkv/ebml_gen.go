// https://raw.githubusercontent.com/cellar-wg/matroska-specification/aa2144a58b661baf54b99bab41113d66b0f5ff62/ebml_matroska.xml

// +build ignore

package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type EBMLSchema struct {
	Elements []element `xml:"element"`
}

// <element name="EBMLMaxIDLength" path="\EBML\EBMLMaxIDLength" id="0x42F2" type="uinteger" range="4" default="4" minOccurs="1" maxOccurs="1"/>
type element struct {
	Name          string          `xml:"name,attr"`
	Path          string          `xml:"path,attr"`
	ID            string          `xml:"id,attr"`
	Type          string          `xml:"type,attr"`
	Range         string          `xml:"range,attr"`
	Default       string          `xml:"default,attr"`
	MinOccurs     string          `xml:"minOccurs,attr"`
	MaxOccurs     string          `xml:"maxOccurs,attr"`
	Length        string          `xml:"length,attr"`
	Documentation []documentation `xml:"documentation"`
}

// <documentation lang="en" purpose="definition">A randomly generated unique ID to identify the Segment amongst many others (128 bits).</documentation>
type documentation struct {
	Purpose string `xml:"purpose,attr"`
	Value   string `xml:",cdata"`
}

func main() {
	r, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	xd := xml.NewDecoder(r)
	var es EBMLSchema
	xd.Decode(&es)

	for _, e := range es.Elements {
		var children []element
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

		fmt.Printf("var mkv%s = ebmlTag{\n", e.Name)
		for _, c := range children {
			id := strings.ToLower(c.ID[2:])
			var def string
			for _, d := range c.Documentation {
				if d.Purpose == "definition" {
					def = strings.TrimSpace(d.Value)
					break
				}
			}
			_ = def
			extra := ""
			typ := c.Type
			switch typ {
			case "master":
				extra = ", tag: mkv" + c.Name
			case "utf-8":
				typ = "UTF8"
			}

			fmt.Printf("\t0x%s: {name: %q, typ: ebml%s%s},\n", id, c.Name, strings.Title(typ), extra)

		}
		fmt.Printf("}\n")
		fmt.Printf("\n")
	}

}
