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
	New:  func(v *decode.Value) decode.FieldWriter { return &FieldWriter{v: v} },
}

type FieldWriter struct {
	v *decode.Value
}

func (o *FieldWriter) outputValue(cw *columnwriter.Writer, v *decode.Value, depth int) error {
	isInArray := false
	if v.Parent != nil {
		_, isInArray = v.Parent.V.(decode.Array)
	}

	indent := strings.Repeat("  ", depth)

	switch v.V.(type) {
	case decode.Struct:
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
		if isInArray {
			fmt.Fprintf(cw.Columns[6], "%s%s[%d]{}: ", indent, v.Parent.Name, v.Index)
		} else {
			fmt.Fprintf(cw.Columns[6], "%s%s{}: ", indent, v.Name)
		}
		fmt.Fprintf(cw.Columns[6], "%s %s (%s)\n", v, v.Range, decode.Bits(v.Range.Length()))
	case decode.Array:
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
		fmt.Fprintf(cw.Columns[6], "%s%s[]: %s\n", indent, v.Name, v.Range)
	default:
		absRange := v.Range

		startBit := absRange.Start
		stopBit := absRange.Stop

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
			}

			if lastDisplayByte > stopByte || stopByte-lastDisplayByte <= lineBytes {
				lastDisplayByte = stopByte
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

		addrLines := 1

		// log.Printf("startBit: %x\n", startBit)
		// log.Printf("stopBit: %x\n", stopBit)

		// log.Printf("addrLines: %x\n", addrLines)

		color := false

		//b := f.BitBuf()

		//f.Value

		//b := f.BitBuf()
		// TODO: abs bitbuf
		//b, _ := v.BitBuf.BitBufRange(startByte*8, displaySizeBytes*8)

		absBitBuf := v.AbsBitBuf()
		vBitBuf, _ := absBitBuf.BitBufRange(startByte*8, displaySizeBytes*8)

		addrLines = int(lastDisplayLine - startLine)
		if lastDisplayByte%lineBytes != 0 {
			addrLines++
		}

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

		if vBitBuf != nil {
			io.Copy(
				hexpairwriter.New(cw.Columns[2], int(lineBytes), int(startLineByteOffset), hexpairFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes))
			io.Copy(
				asciiwriter.New(cw.Columns[4], int(lineBytes), int(startLineByteOffset), asciiFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes))
		}

		// fmt.Fprintf(cw.Columns[2], "%s", hexpairs(b, lineBytes, startLineByteOffset))
		// fmt.Fprintf(cw.Columns[4], "%s", printable(b, startLineByteOffset))

		for i := 0; i < addrLines; i++ {
			fmt.Fprintf(cw.Columns[0], "%.8x\n", startLineByte+int64(i)*lineBytes)
			fmt.Fprintf(cw.Columns[1], "\n")
			fmt.Fprintf(cw.Columns[3], "|\n")
			fmt.Fprintf(cw.Columns[5], "|\n")
		}

		if stopByte != lastDisplayByte {
			fmt.Fprint(cw.Columns[0], "*\n")
			fmt.Fprint(cw.Columns[2], "\n")
			fmt.Fprintf(cw.Columns[2], "%d bytes more, ends at %x", stopByte-lastDisplayByte, stopByte)
			fmt.Fprintf(cw.Columns[3], "|\n")
			fmt.Fprintf(cw.Columns[5], "|\n")
			// TODO: dump last line?
		}

		if isInArray {
			fmt.Fprintf(cw.Columns[6], "%s%s[%d] (%s): ", indent, v.Parent.Name, v.Index, v.Name)
		} else {
			fmt.Fprintf(cw.Columns[6], "%s%s: ", indent, v.Name)
		}
		fmt.Fprintf(cw.Columns[6], "%s %s (%s)\n", v, v.Range, decode.Bits(v.Range.Length()))
	}

	if v.Error != nil {
		fmt.Fprintf(cw.Columns[6], "%s!%s\n", indent, v.Error)
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
	}

	cw.Flush()

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {
	cw := columnwriter.New(w, []int{8, 1, int(lineBytes*3) - 1, 1, int(lineBytes), 1, -1})
	return o.v.WalkPreOrder(func(v *decode.Value, depth int) error {
		return o.outputValue(cw, v, depth)
	})
}
