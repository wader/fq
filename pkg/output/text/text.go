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
	New:  func(f *decode.Field) decode.FieldWriter { return &FieldWriter{f: f} },
}

type FieldWriter struct {
	f *decode.Field
}

func (o *FieldWriter) outputValue(cw *columnwriter.Writer, v decode.Value, name string, depth int) error {

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
	fmt.Fprintf(cw.Columns[6], "%s%s: %s %s (%s)\n", indent, name, v, v.Range, decode.Bits(v.Range.Length()))

	cw.Flush()

	if stopByte != lastDisplayByte {
		fmt.Fprint(cw.Columns[0], "*\n")
		fmt.Fprintf(cw.Columns[2], "%d bytes more, ends at %x", stopByte-lastDisplayByte, stopByte)
		cw.Flush()
		// TODO: dump last line?
	}

	return nil
}

func (o *FieldWriter) output(cw *columnwriter.Writer, f *decode.Field, name string, depth int) error {
	indent := strings.Repeat("  ", depth)

	switch v := f.Value.V.(type) {
	case []*decode.Field:
		fmt.Fprintf(cw.Columns[6], "%s%s: (%s)\n", indent, name, f.Value.Desc)
		cw.Flush()
		for _, wf := range v {
			if err := o.output(cw, wf, wf.Name, depth+1); err != nil {
				return err
			}
		}
	case []decode.Value:
		for i, wv := range v {
			cw.Flush()
			switch wvf := wv.V.(type) {
			case *decode.Field:
				fmt.Fprintf(cw.Columns[6], "%s%s[%d]: (%s)\n", indent, name, i, wvf.Name)
				cw.Flush()
				if err := o.output(cw, wvf, wvf.Name, depth); err != nil {
					return err
				}
			default:
				o.outputValue(cw, wv, fmt.Sprintf("%s[%d]", name, i), depth+1)
			}
		}
	default:
		o.outputValue(cw, f.Value, name, depth+1)
	}

	if f.Error != nil {
		fmt.Fprintf(cw.Columns[6], "%s! %s\n", indent, f.Error)
		cw.Flush()
	}

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {

	cw := columnwriter.New(w, []int{8, 1, int(lineBytes*3) - 1, 1, int(lineBytes), 1, -1})

	return o.output(cw, o.f, o.f.Name, 0)
}
