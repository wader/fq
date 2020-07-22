package text

import (
	"encoding/json"
	"fmt"
	"fq/pkg/decode"
	"io"
	"strings"
)

var FieldOutput = &decode.FieldOutput{
	Name: "text",
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

func (o *FieldWriter) output(w io.Writer, f *decode.Field, depth int) error {
	indent := strings.Repeat("  ", depth)

	if (len(f.Children)) != 0 {
		fmt.Printf("%s%s: %s %s (%s) {\n", indent, f.Name, f.Value, f.Range, decode.Bits(f.Range.Length()))
		for _, c := range f.Children {
			o.output(w, c, depth+1)
		}
		fmt.Printf("%s}\n", indent)
	} else {
		fmt.Printf("%s%s: %s %s (%s)\n", indent, f.Name, f.Value, f.Range, decode.Bits(f.Range.Length()))
	}

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {
	return o.output(w, o.f, 0)
}
