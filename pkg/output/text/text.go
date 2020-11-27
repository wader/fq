package text

import (
	"fmt"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/pkg/decode"
	"io"
	"math"
	"strconv"
	"strings"
)

const lineBytes = int64(16)
const maxBytes = int64(16)
const addrBase = 16
const sizeBase = 10

var FieldOutput = &decode.FieldOutput{
	Name: "text",
	New:  func(v *decode.Value) decode.FieldWriter { return &FieldWriter{v: v} },
}

func digitsInBase(n int64, base int) int {
	if n == 0 {
		return 1
	}
	return int(1 + math.Floor(math.Log(float64(n))/math.Log(float64(base))))
}

func padFormatInt(i int64, base int, width int) string {
	s := strconv.FormatInt(i, base)
	p := width - len(s)
	if p > 0 {
		// TODO: something faster?
		return strings.Repeat("0", p) + s
	}
	return s
}

func bitsCeilBytes(bits int64) int64 {
	if bits&0x7 != 0 {
		return bits>>3 + 1
	}
	return bits >> 3
}

type FieldWriter struct {
	v *decode.Value
}

func (o *FieldWriter) outputValue(cw *columnwriter.Writer, v *decode.Value, depth int, rootDepth int, addrWidth int) error {
	isInArray := false
	if v.Parent != nil {
		_, isInArray = v.Parent.V.(decode.Array)
	}

	rootIndent := strings.Repeat(" ", rootDepth)
	indent := strings.Repeat("  ", depth)

	startBit := v.Range.Start
	stopBit := v.Range.Stop() - 1
	sizeBits := v.Range.Len
	lastDisplayBit := stopBit

	if sizeBits > maxBytes*8 {
		lastDisplayBit = startBit + (maxBytes*8 - 1)
		if lastDisplayBit%(lineBytes*8) != 0 {
			lastDisplayBit += (lineBytes * 8) - lastDisplayBit%(lineBytes*8) - 1
		}

		if lastDisplayBit > stopBit || stopBit-lastDisplayBit <= lineBytes*8 {
			lastDisplayBit = stopBit
		}
	}

	startByte := startBit / 8
	stopByte := stopBit / 8
	lastDisplayByte := lastDisplayBit / 8
	displaySizeBytes := lastDisplayByte - startByte + 1

	startLine := startByte / lineBytes
	startLineByteOffset := startByte % lineBytes
	startLineByte := startLine * lineBytes
	lastDisplayLine := lastDisplayByte / lineBytes

	fmt.Fprintf(cw.Columns[1], "|\n")
	fmt.Fprintf(cw.Columns[3], "|\n")
	fmt.Fprintf(cw.Columns[5], "|\n")

	switch vv := v.V.(type) {
	case decode.Struct:
		if isInArray {
			fmt.Fprintf(cw.Columns[6], "%s%s[%d]{}: ", indent, v.Parent.Name, v.Index)
		} else {
			fmt.Fprintf(cw.Columns[6], "%s%s{}: ", indent, v.Name)
		}
		fmt.Fprintf(cw.Columns[6], "%s %s fields %d (%s)\n", v, decode.BitRange(v.Range).StringByteBits(addrBase), len(vv), decode.Bits(v.Range.Len).StringByteBits(sizeBase))
	case decode.Array:
		fmt.Fprintf(cw.Columns[6], "%s%s[]: %s %s count %d (%s)\n", indent, v.Name, v, decode.BitRange(v.Range).StringByteBits(addrBase), len(vv), decode.Bits(v.Range.Len).StringByteBits(sizeBase))
	default:
		fmt.Fprintf(cw.Columns[0], "%s%s\n", rootIndent, padFormatInt(startLineByte, addrBase, addrWidth))

		color := false
		vBitBuf, err := v.BitBuf.BitBufRange(startByte*8, displaySizeBytes*8)
		if err != nil {
			return err
		}

		addrLines := lastDisplayLine - startLine + 1

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

		for i := int64(1); i < addrLines; i++ {
			fmt.Fprintf(cw.Columns[0], "%s%s\n", rootIndent, padFormatInt(startLineByte+int64(i)*lineBytes, addrBase, addrWidth))
			fmt.Fprintf(cw.Columns[1], "|\n")
			fmt.Fprintf(cw.Columns[3], "|\n")
			fmt.Fprintf(cw.Columns[5], "|\n")
		}

		if stopByte != lastDisplayByte {
			fmt.Fprintf(cw.Columns[0], "*\n")
			fmt.Fprintf(cw.Columns[1], "|\n")
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
		fmt.Fprintf(cw.Columns[6], "%s %s (%s)\n", v, decode.BitRange(v.Range).StringByteBits(addrBase), decode.Bits(v.Range.Len).StringByteBits(sizeBase))
	}

	if v.Error != nil {
		fmt.Fprintf(cw.Columns[6], "%s!%s\n", indent, v.Error)
		fmt.Fprintf(cw.Columns[1], "|\n")
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
	}

	cw.Flush()

	return nil
}

func (o *FieldWriter) Write(w io.Writer) error {
	maxAddrIndentWidth := digitsInBase(bitsCeilBytes(o.v.BitBuf.Len()), addrBase)
	o.v.WalkPreOrder(func(v *decode.Value, depth int, rootDepth int) error {
		// skip first root level
		if rootDepth > 0 {
			rootDepth--
		}

		if v.IsRoot {
			addrIndentWidth := rootDepth + digitsInBase(bitsCeilBytes(v.BitBuf.Len()), addrBase)
			if addrIndentWidth > maxAddrIndentWidth {
				maxAddrIndentWidth = addrIndentWidth
			}
		}

		return nil
	})

	cw := columnwriter.New(w, []int{maxAddrIndentWidth, 1, int(lineBytes*3) - 1, 1, int(lineBytes), 1, -1})
	return o.v.WalkPreOrder(func(v *decode.Value, depth int, rootDepth int) error {
		// skip first root level
		if rootDepth > 0 {
			rootDepth--
		}

		return o.outputValue(cw, v, depth, rootDepth, maxAddrIndentWidth-rootDepth)
	})
}
