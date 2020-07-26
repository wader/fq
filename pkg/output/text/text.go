package text

import (
	"fmt"
	"fq/internal/columnwriter"
	"fq/internal/hexdump"
	"fq/pkg/decode"
	"io"
	"strings"
)

var FieldOutput = &decode.FieldOutput{
	Name: "text",
	New:  func(f *decode.Field) decode.FieldWriter { return &FieldWriter{f: f} },
}

type FieldWriter struct {
	f *decode.Field
}

func (o *FieldWriter) output(cw *columnwriter.Writer, f *decode.Field, depth int) error {
	indent := strings.Repeat("  ", depth)

	if f.Value.Type != decode.TypeDecoder {
		b, err := f.BitBuf().BytesBitRange(0, f.Range.Length(), 0)
		if err != nil {
			return err
		}
		start := f.Decoder.AbsPos(f.Range.Start)
		h := hexdump.Dumper(start/8, cw.Column(0))
		h.Write(b)
		h.Close()
	}

	fmt.Fprintf(cw.Column(1), "%s%s: %s %s (%s)\n", indent, f.Name, f.Value, f.Range, decode.Bits(f.Range.Length()))
	// if f.Children != nil {
	// 	fmt.Fprint(cw.Column(1), " {")
	// }
	// fmt.Fprint(cw.Column(1), "\n")

	cw.Flush()

	for _, c := range f.Children {
		o.output(cw, c, depth+1)
	}

	// if f.Children != nil {
	// 	fmt.Fprintf(cw.Column(1), "%s}\n", indent)
	// }

	//cw.Flush()

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {
	cw := columnwriter.New(w, []int{78, -1})

	return o.output(cw, o.f, 0)
}
