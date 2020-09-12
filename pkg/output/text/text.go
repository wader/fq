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

const lineBytes = int64(16)
const maxBytes = int64(32)

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
	rangeBytes := truncatedStopByte - startByte + 1

	startLineByte := (startByte / lineBytes) * lineBytes
	startLineByteOffset := startByte % lineBytes
	stopLineByte := (stopByte / lineBytes) * lineBytes
	truncatedStopLineByte := truncatedStopByte / lineBytes * lineBytes
	var addrLines int64 = 1

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
			hexpairwriter.New(cw.Columns[2], int(lineBytes), int(startLineByteOffset)),
			io.LimitReader(b.Copy(), rangeBytes))
		io.Copy(
			asciiwriter.New(cw.Columns[4], int(lineBytes), int(startLineByteOffset)),
			io.LimitReader(b.Copy(), rangeBytes))

		// fmt.Fprintf(cw.Columns[2], "%s", hexpairs(b, lineBytes, startLineByteOffset))
		// fmt.Fprintf(cw.Columns[4], "%s", printable(b, startLineByteOffset))
	}

	for i := int64(0); i < addrLines; i++ {
		fmt.Fprintf(cw.Columns[0], "%.8x\n", startLineByte+i*lineBytes)
		fmt.Fprintf(cw.Columns[1], "\n")
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
	}

	indent := strings.Repeat("  ", depth)
	nameIndex := ""
	if f.Index > 0 {
		nameIndex = fmt.Sprintf("[%d]", f.Index)
	}
	fmt.Fprintf(cw.Columns[6], "%s%s%s: %s %s (%s)\n", indent, f.Name, nameIndex, f.Value, f.Range, decode.Bits(f.Range.Length()))

	cw.Flush()

	if len(f.Children) == 0 && stopLineByte != truncatedStopLineByte {
		fmt.Fprint(cw.Columns[0], "*\n")
		fmt.Fprintf(cw.Columns[2], "%d byte (end at %x)", stopByte-truncatedStopLineByte, stopByte)
		cw.Flush()
		// TODO: dump last line?
	}

	for _, c := range f.Children {
		if err := o.output(cw, c, depth+1); err != nil {
			return err
		}
	}

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {

	cw := columnwriter.New(w, []int{8, 1, int(lineBytes*3) - 1, 1, int(lineBytes), 1, -1})

	return o.output(cw, o.f, 0)
}
