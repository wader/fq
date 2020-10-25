package text

import (
	"fmt"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/pkg/bitbuf"
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

func (o *FieldWriter) outputValue(cw *columnwriter.Writer, v *decode.Value, index int, depth int) error {

	startBit := v.Range.Start
	stopBit := v.Range.Stop

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

	var b *bitbuf.Buffer
	if v.BitBuf != nil {
		b, _ = v.BitBuf.BitBufRange(0, 0)
	}

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

	if b != nil {
		io.Copy(
			hexpairwriter.New(cw.Columns[2], int(lineBytes), int(startLineByteOffset), hexpairFn),
			io.LimitReader(b.Copy(), displaySizeBytes))
		io.Copy(
			asciiwriter.New(cw.Columns[4], int(lineBytes), int(startLineByteOffset), asciiFn),
			io.LimitReader(b.Copy(), displaySizeBytes))
	}

	// fmt.Fprintf(cw.Columns[2], "%s", hexpairs(b, lineBytes, startLineByteOffset))
	// fmt.Fprintf(cw.Columns[4], "%s", printable(b, startLineByteOffset))

	for i := 0; i < addrLines; i++ {
		fmt.Fprintf(cw.Columns[0], "%.8x\n", startLineByte+int64(i)*lineBytes)
		fmt.Fprintf(cw.Columns[1], "\n")
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
	}

	indent := strings.Repeat("  ", depth)

	switch v.Parent.V.(type) {
	case decode.Struct:
		fmt.Fprintf(cw.Columns[6], "%s%s: %s %s (%s)\n", indent, v.Name, v, v.Range, decode.Bits(v.Range.Length()))
	case decode.Array:
		fmt.Fprintf(cw.Columns[6], "%s%s[%d] (%s): %s %s (%s)\n", indent, v.Parent.Name, index, v.Name, v, v.Range, decode.Bits(v.Range.Length()))
	default:
		panic("invalid parent")
	}

	cw.Flush()

	if stopByte != lastDisplayByte {
		fmt.Fprint(cw.Columns[0], "*\n")
		fmt.Fprintf(cw.Columns[2], "%d bytes more, ends at %x", stopByte-lastDisplayByte, stopByte)
		cw.Flush()
		// TODO: dump last line?
	}

	return nil
}

func (o *FieldWriter) write(cw *columnwriter.Writer, v *decode.Value) error {
	var walkFn func(v *decode.Value, index int, depth int) error
	walkFn = func(v *decode.Value, index int, depth int) error {
		indent := strings.Repeat("  ", depth)

		switch vv := v.V.(type) {
		case decode.Struct:
			fmt.Fprintf(cw.Columns[6], "%s%s: %s %s\n", indent, v.Name, v, v.Range)
			cw.Flush()
			for _, wv := range vv {
				walkFn(wv, -1, depth+1)
			}
		case decode.Array:
			// fmt.Fprintf(cw.Columns[6], "%s%s: %s\n", indent, name, v)
			cw.Flush()
			for i, wv := range vv {
				walkFn(wv, i, depth)
			}
		default:
			o.outputValue(cw, v, index, depth)
		}
		return nil
	}
	return walkFn(v, -1, 0)
}

func (o *FieldWriter) Write(w io.Writer) error {
	cw := columnwriter.New(w, []int{8, 1, int(lineBytes*3) - 1, 1, int(lineBytes), 1, -1})
	return o.write(cw, o.v)
}
