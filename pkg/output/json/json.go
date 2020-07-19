package json

import (
	"encoding/json"
	"fmt"
	"fq/pkg/decode"
	"io"
	"strings"
)

var FieldOutput = &decode.FieldOutput{
	Name: "json",
	New:  func(f *decode.Field) decode.FieldWriter { return &FieldWriter{f: f} },
}

type prefixPrinter struct {
	w      io.Writer
	prefix string
}

func (pp prefixPrinter) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(pp.w, pp.prefix+format+"\n", a...)
}

type FieldWriter struct {
	f *decode.Field
}

func jsonEscape(v interface{}) string {
	e, _ := json.Marshal(v)
	return string(e)
}

func jsonField(k string, v interface{}) string {
	return fmt.Sprintf("%s: %s", jsonEscape(k), jsonEscape(v))
}

func (o *FieldWriter) output(w io.Writer, f *decode.Field, depth int) {
	p := prefixPrinter{w: w, prefix: strings.Repeat("  ", depth)}

	p.Printf("%s: {", jsonEscape(f.Name))
	p.Printf(`  "value": %s`, jsonEscape(f.Value))

	if (len(f.Children)) != 0 {
		p.Printf(`  "fields": {`)
		for _, c := range f.Children {
			o.output(w, c, depth+1)
		}
		p.Printf("  }\n")
	}

	p.Printf("}")
}

func (o *FieldWriter) Write(w io.Writer) {
	o.output(w, o.f, 0)
}
