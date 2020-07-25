package text

import (
	"encoding/hex"
	"fmt"
	"fq/internal/columnwriter"
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

func (o *FieldWriter) output(cw *columnwriter.Writer, f *decode.Field, depth int) error {
	indent := strings.Repeat("  ", depth)

	if (len(f.Children)) != 0 {
		fmt.Fprintf(cw, "%s%s: %s %s (%s) {\n", indent, f.Name, f.Value, f.Range, decode.Bits(f.Range.Length()))
		for _, c := range f.Children {
			o.output(cw, c, depth+1)
		}
		fmt.Fprintf(cw, "%s}\n", indent)

	} else {
		b, err := f.BitBuf().BytesBitRange(0, f.Range.Length(), 0)
		if err != nil {
			return err
		}
		h := hex.Dumper(cw)
		h.Write(b)
		h.Close()

		cw.Next()

		fmt.Fprintf(cw, "%s%s: %s %s (%s)\n", indent, f.Name, f.Value, f.Range, decode.Bits(f.Range.Length()))
		cw.Next()

	}

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {
	cw := columnwriter.New(w, []int{80, -1})

	return o.output(cw, o.f, 0)
}
