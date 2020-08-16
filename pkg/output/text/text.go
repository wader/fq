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

	startBit := f.Decoder.AbsPos(f.Range.Start)
	stopBit := f.Decoder.AbsPos(f.Range.Stop)

	startByte := startBit / 8
	stopByte := stopBit / 8
	if startBit%8 != 0 {
		stopByte++
	}

	lineAddrByte := int((startByte >> 4) << 4)
	lineAddrByteOffset := int(startByte & 0xf)

	addrLines := 1

	if len(f.Children) == 0 {
		b, err := f.BitBuf().BytesBitRange(0, f.Range.Length(), 0)
		if err != nil {
			return err
		}

		addrStopBytes := lineAddrByteOffset + int(stopByte) - int(startByte)
		addrLines = addrStopBytes / 16
		if addrStopBytes%16 != 0 {
			addrLines++
		}

		fmt.Fprintf(cw.Columns[2], "%s", hexdump.Hexpairs(lineAddrByteOffset, b))
		fmt.Fprintf(cw.Columns[4], "%s", hexdump.Printable(lineAddrByteOffset, b))
	}

	for i := 0; i < addrLines; i++ {
		fmt.Fprintf(cw.Columns[0], "%.8x\n", int(lineAddrByte)+i*16)
		fmt.Fprintf(cw.Columns[1], "|\n")
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
	}

	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(cw.Columns[6], "%s%s: %s %s (%s)\n", indent, f.Name, f.Value, f.Range, decode.Bits(f.Range.Length()))
	// if f.Children != nil {
	// 	fmt.Fprint(cw.Column(1), " {")
	// }
	// fmt.Fprint(cw.Column(1), "\n")

	cw.Flush()

	for _, c := range f.Children {
		if err := o.output(cw, c, depth+1); err != nil {
			return err
		}
	}

	// if f.Children != nil {
	// 	fmt.Fprintf(cw.Column(1), "%s}\n", indent)
	// }

	//cw.Flush()

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {
	cw := columnwriter.New(w, []int{8, 1, 47, 1, 16, 1, -1})
	cw.Columns[1].Wrap = true
	cw.Columns[2].Wrap = true
	cw.Columns[4].Wrap = true

	return o.output(cw, o.f, 0)
}
