package json

import (
	"encoding/json"
	"fmt"
	"fq/internal/decode"
	"io"
	"strings"
)

type prefixPrinter struct {
	w      io.Writer
	prefix string
}

func (pp prefixPrinter) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(pp.w, pp.prefix+format+"\n", a...)
}

func New(w io.Writer) *Output {
	return &Output{w: w}
}

type Output struct {
	w io.Writer
}

func jsonEscape(v interface{}) string {
	e, _ := json.Marshal(v)
	return string(e)
}

func jsonField(k string, v interface{}) string {
	return fmt.Sprintf("%s: %s", jsonEscape(k), jsonEscape(v))
}

func (o *Output) output(f *decode.Field, depth int) {
	p := prefixPrinter{w: o.w, prefix: strings.Repeat("  ", depth)}

	p.Printf("%s: {", jsonEscape(f.Name))
	p.Printf(`  "value": %s`, jsonEscape(f.Value))

	if (len(f.Children)) != 0 {
		p.Printf(`  "fields": {`)
		for _, c := range f.Children {
			o.output(c, depth+1)
		}
		p.Printf("  }\n")
	}

	p.Printf("}")
}

func (o *Output) Output(f *decode.Field) {
	o.output(f, 0)
}
