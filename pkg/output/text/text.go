package text

import (
	"fmt"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/pkg/decode"
	"io"
	"strings"
)

const lineBytes = 16
const maxBytes = 32

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

	if startBit != stopBit {
		stopBit--
	}

	startByte := startBit / 8
	stopByte := stopBit / 8
	truncatedStopByte := stopByte
	if truncatedStopByte-startByte > maxBytes {
		truncatedStopByte = startByte + maxBytes
	}
	rangeBytes := truncatedStopByte - startByte

	startLineByte := int((startByte / lineBytes) * lineBytes)
	startLineByteOffset := int(startByte % lineBytes)
	stopLineByte := int((stopByte / lineBytes) * lineBytes)
	truncatedStopLineByte := int((truncatedStopByte / lineBytes) * lineBytes)
	addrLines := 1

	// log.Printf("startBit: %x\n", startBit)
	// log.Printf("stopBit: %x\n", stopBit)

	// log.Printf("startByte: %x\n", startByte)
	// log.Printf("stopByte: %x\n", stopByte)

	// log.Printf("startLineByte: %x\n", startLineByte)
	// log.Printf("stopLineByte: %x\n", stopLineByte)

	// log.Printf("addrLines: %x\n", addrLines)

	if len(f.Children) == 0 {
		b := f.BitBuf()
		addrLines = ((truncatedStopLineByte - startLineByte) / lineBytes) + 1

		io.Copy(
			hexpairwriter.New(cw.Columns[2], lineBytes, startLineByteOffset),
			io.LimitReader(b.Copy(), rangeBytes))
		io.Copy(
			asciiwriter.New(cw.Columns[4], lineBytes, startLineByteOffset),
			io.LimitReader(b.Copy(), rangeBytes))

		// fmt.Fprintf(cw.Columns[2], "%s", hexpairs(b, lineBytes, startLineByteOffset))
		// fmt.Fprintf(cw.Columns[4], "%s", printable(b, startLineByteOffset))
	}

	for i := 0; i < addrLines; i++ {
		fmt.Fprintf(cw.Columns[0], "%.8x\n", int(startLineByte)+i*lineBytes)
		fmt.Fprintf(cw.Columns[1], "\n")
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

	if len(f.Children) == 0 && stopLineByte != truncatedStopLineByte {
		fmt.Fprint(cw.Columns[0], "*\n")
		fmt.Fprintf(cw.Columns[2], "%d byte skipped", stopLineByte-truncatedStopLineByte)
		cw.Flush()
		// TODO: dump last line?
		fmt.Fprintf(cw.Columns[0], "%.8x\n", int(stopLineByte))
		fmt.Fprintf(cw.Columns[2], "...")
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
		cw.Flush()
	}

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
	cw := columnwriter.New(w, []int{8, 1, lineBytes*3 - 1, 1, lineBytes, 1, -1})
	// cw.Columns[1].Wrap = true
	// cw.Columns[2].Wrap = true
	// cw.Columns[4].Wrap = true

	return o.output(cw, o.f, 0)
}
