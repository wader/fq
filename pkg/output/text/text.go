package text

import (
	"fmt"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/pkg/decode"
	"io"
	"strings"
)

const lineBytes = int64(16)
const maxBytes = int64(16)

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
	if stopBit%8 != 0 {
		stopByte++
	}
	sizeBytes := stopByte - startByte

	lastDisplayByte := stopByte
	if sizeBytes > maxBytes {
		// truncate but fill line
		// TODO: redo with max etc?
		lastDisplayByte = startByte + maxBytes
		if lastDisplayByte%lineBytes != 0 {
			lastDisplayByte += lineBytes - lastDisplayByte%lineBytes
			if lastDisplayByte > stopByte {
				lastDisplayByte = stopByte
			}
		}
	}
	displaySizeBytes := lastDisplayByte - startByte
	if displaySizeBytes == 0 {
		displaySizeBytes = 1
	}

	startLine := startByte / lineBytes
	startLineByteOffset := startByte % lineBytes
	startLineByte := startLine * lineBytes
	lastDisplayLine := lastDisplayByte / lineBytes

	addrLines := lastDisplayLine - startLine
	if lastDisplayByte%lineBytes != 0 {
		addrLines++
	}

	// log.Printf("startBit: %x\n", startBit)
	// log.Printf("stopBit: %x\n", stopBit)

	// log.Printf("startByte: %x\n", startByte)
	// log.Printf("stopByte: %x\n", stopByte)

	// log.Printf("startLineByte: %x\n", startLineByte)
	// log.Printf("stopLineByte: %x\n", stopLineByte)

	// log.Printf("addrLines: %x\n", addrLines)

	color := false

	if len(f.Children) == 0 {
		//b := f.BitBuf()

		b, _ := f.Decoder.AbsBitBuf().BitBufRange(startByte*8, displaySizeBytes*8)

		// log.Printf("truncatedStopLineByte: %#+v\n", truncatedStopLineByte)
		// log.Printf("startLineByte: %#+v\n", startLineByte)
		// log.Printf("truncatedStopLineByte - startLineByte: %#+v\n", truncatedStopLineByte-startLineByte)

		// log.Printf("addrLines: %#+v\n", addrLines)

		charToANSI := func(c byte) string {
			switch {
			case c == 0:
				return ansi.FgBrightBlack
			case c >= 32 && c <= 126, c == '\r', c == '\n', c == '\f', c == '\t', c == '\v':
				return ansi.FgWhite
			default:
				return ansi.FgBrightWhite
			}
		}

		hexpairFn := func(c byte) string {
			const hexTable = "0123456789abcdef"
			if color {
				return fmt.Sprintf("%s%c%c%s", charToANSI(c), hexTable[c>>4], hexTable[c&0xf], ansi.Reset)

			}
			return string(hexTable[c>>4]) + string(hexTable[c&0xf])
		}

		asciiFn := func(c byte) string {
			d := c
			if c < 32 || c > 126 {
				d = '.'
			}
			if color {
				return fmt.Sprintf("%s%c%s", charToANSI(c), d, ansi.Reset)
			}
			return string(d)
		}

		io.Copy(
			hexpairwriter.New(cw.Columns[2], int(lineBytes), int(startLineByteOffset), hexpairFn),
			io.LimitReader(b.Copy(), displaySizeBytes))
		io.Copy(
			asciiwriter.New(cw.Columns[4], int(lineBytes), int(startLineByteOffset), asciiFn),
			io.LimitReader(b.Copy(), displaySizeBytes))

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

	if len(f.Children) == 0 && stopByte != lastDisplayByte {
		fmt.Fprint(cw.Columns[0], "*\n")
		fmt.Fprintf(cw.Columns[2], "%d bytes more (end at %x)", stopByte-lastDisplayByte, stopByte)
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
