package text

import (
	"fmt"
	"fq/internal/columnwriter"
	"fq/pkg/decode"
	"io"
	"strings"
)

const lineBytes = 16

const hextable = "0123456789abcdef"

func hexpairs(b []byte, width int, offset int) string {
	if len(b) == 0 {
		return ""
	}
	t := offset + len(b)
	s := make([]byte, t*3-1)

	for i := 0; i < t; i++ {
		if i < offset {
			s[i*3+0] = ' '
			s[i*3+1] = ' '
		} else {
			v := b[i-offset]
			s[i*3+0] = hextable[v>>4]
			s[i*3+1] = hextable[v&0xf]
		}
		if i != t-1 {
			if i > 0 && i%width == width-1 {
				s[i*3+2] = '\n'
			} else {
				s[i*3+2] = ' '
			}
		}
	}

	return string(s[0 : t*3-1])
}

func printable(b []byte, offset int) string {
	t := offset + len(b)
	s := make([]byte, t)

	for i := 0; i < t; i++ {
		if i < offset {
			s[i] = ' '
		} else {
			v := b[i-offset]
			if v < 32 || v > 126 {
				s[i] = '.'
			} else {
				s[i] = v
			}
		}
	}

	return string(s)
}

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

	startLineByte := int((startByte / lineBytes) * lineBytes)
	startLineByteOffset := int(startByte % lineBytes)
	stopLineByte := int((stopByte / lineBytes) * lineBytes)
	addrLines := 1

	// log.Printf("startBit: %x\n", startBit)
	// log.Printf("stopBit: %x\n", stopBit)

	// log.Printf("startByte: %x\n", startByte)
	// log.Printf("stopByte: %x\n", stopByte)

	// log.Printf("startLineByte: %x\n", startLineByte)
	// log.Printf("stopLineByte: %x\n", stopLineByte)

	// log.Printf("addrLines: %x\n", addrLines)

	if len(f.Children) == 0 {
		b, err := f.BitBuf().BytesBitRange(0, f.Range.Length(), 0)
		if err != nil {
			return err
		}
		addrLines = ((stopLineByte - startLineByte) / lineBytes) + 1

		fmt.Fprintf(cw.Columns[2], "%s", hexpairs(b, lineBytes, startLineByteOffset))
		fmt.Fprintf(cw.Columns[4], "%s", printable(b, startLineByteOffset))
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
	cw.Columns[1].Wrap = true
	cw.Columns[2].Wrap = true
	cw.Columns[4].Wrap = true

	return o.output(cw, o.f, 0)
}
