package decode

import (
	"fmt"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/pkg/bitio"
	"io"
	"math"
	"strconv"
	"strings"
)

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

type DumpOptions struct {
	MaxDepth        int
	LineBytes       int
	MaxDisplayBytes int
	AddrBase        int
	SizeBase        int
}

func (v *Value) dump(cw *columnwriter.Writer, depth int, rootDepth int, addrWidth int, opts DumpOptions) error {
	isInArray := false
	if v.Parent != nil {
		_, isInArray = v.Parent.V.(Array)
	}

	rootIndent := strings.Repeat(" ", rootDepth)
	indent := strings.Repeat("  ", depth)

	bufferLastBit := v.BitBuf.Len() - 1
	startBit := v.Range.Start
	stopBit := v.Range.Stop() - 1
	sizeBits := v.Range.Len
	lastDisplayBit := stopBit

	if sizeBits > int64(opts.MaxDisplayBytes)*8 {
		lastDisplayBit = startBit + (int64(opts.MaxDisplayBytes)*8 - 1)
		if lastDisplayBit%(int64(opts.LineBytes)*8) != 0 {
			lastDisplayBit += (int64(opts.LineBytes) * 8) - lastDisplayBit%(int64(opts.LineBytes)*8) - 1
		}

		if lastDisplayBit > stopBit || stopBit-lastDisplayBit <= int64(opts.LineBytes)*8 {
			lastDisplayBit = stopBit
		}
	}

	bufferLastByte := bufferLastBit / 8
	startByte := startBit / 8
	stopByte := stopBit / 8
	lastDisplayByte := lastDisplayBit / 8
	displaySizeBytes := lastDisplayByte - startByte + 1

	startLine := startByte / int64(opts.LineBytes)
	startLineByteOffset := startByte % int64(opts.LineBytes)
	startLineByte := startLine * int64(opts.LineBytes)
	lastDisplayLine := lastDisplayByte / int64(opts.LineBytes)

	fmt.Fprintf(cw.Columns[1], "|\n")
	fmt.Fprintf(cw.Columns[3], "|\n")
	fmt.Fprintf(cw.Columns[5], "|\n")

	switch vv := v.V.(type) {
	case Struct:
		if isInArray {
			fmt.Fprintf(cw.Columns[6], "%s%s[%d]{}: ", indent, v.Parent.Name, v.Index)
		} else {
			fmt.Fprintf(cw.Columns[6], "%s%s{}: ", indent, v.Name)
		}
		fmt.Fprintf(cw.Columns[6], "%s %s fields %d (%s)\n", v, BitRange(v.Range).StringByteBits(opts.AddrBase), len(vv), Bits(v.Range.Len).StringByteBits(opts.SizeBase))
	case Array:
		fmt.Fprintf(cw.Columns[6], "%s%s[]: %s %s count %d (%s)\n", indent, v.Name, v, BitRange(v.Range).StringByteBits(opts.AddrBase), len(vv), Bits(v.Range.Len).StringByteBits(opts.SizeBase))
	default:
		fmt.Fprintf(cw.Columns[0], "%s%s\n", rootIndent, padFormatInt(startLineByte, opts.AddrBase, addrWidth))

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
				hexpairwriter.New(cw.Columns[2], opts.LineBytes, int(startLineByteOffset), hexpairFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes))
			io.Copy(
				asciiwriter.New(cw.Columns[4], opts.LineBytes, int(startLineByteOffset), asciiFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes))
		}

		for i := int64(1); i < addrLines; i++ {
			lineStartByte := startLineByte + int64(i)*int64(opts.LineBytes)
			fmt.Fprintf(cw.Columns[0], "%s%s\n", rootIndent, padFormatInt(lineStartByte, opts.AddrBase, addrWidth))
			fmt.Fprintf(cw.Columns[1], "|\n")
			fmt.Fprintf(cw.Columns[3], "|\n")
			fmt.Fprintf(cw.Columns[5], "|\n")
		}
		// TODO: correct?
		lastLineStopByte := startLineByte + int64(addrLines)*int64(opts.LineBytes) - 1
		if lastDisplayByte == bufferLastByte && lastDisplayByte != lastLineStopByte {
			fmt.Fprintf(cw.Columns[2], "|\n")
			fmt.Fprintf(cw.Columns[4], "|\n")
		}

		if stopByte != lastDisplayByte {
			isEnd := ""
			if stopBit == bufferLastBit {
				isEnd = " (end)"
			}
			fmt.Fprintf(cw.Columns[0], "*\n")
			fmt.Fprintf(cw.Columns[1], "|\n")
			fmt.Fprint(cw.Columns[2], "\n")
			fmt.Fprintf(cw.Columns[2], "%d bytes more until %s%s", stopByte-lastDisplayByte, Bits(stopBit).StringByteBits(opts.AddrBase), isEnd)
			fmt.Fprintf(cw.Columns[3], "|\n")
			fmt.Fprintf(cw.Columns[5], "|\n")
			// TODO: dump last line?
		}

		if isInArray {
			fmt.Fprintf(cw.Columns[6], "%s%s[%d] (%s): ", indent, v.Parent.Name, v.Index, v.Name)
		} else {
			fmt.Fprintf(cw.Columns[6], "%s%s: ", indent, v.Name)
		}
		fmt.Fprintf(cw.Columns[6], "%s %s (%s)\n", v, BitRange(v.Range).StringByteBits(opts.AddrBase), Bits(v.Range.Len).StringByteBits(opts.SizeBase))
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

func (v *Value) Dump(w io.Writer, opts DumpOptions) error {
	maxAddrIndentWidth := digitsInBase(bitio.BitsByteCount(v.BitBuf.Len()), opts.AddrBase)
	makeWalkFn := func(fn WalkFn) WalkFn {
		return func(v *Value, depth int, rootDepth int) error {
			if opts.MaxDepth != 0 && depth > opts.MaxDepth {
				return ErrWalkSkip
			}
			// skip first root level
			if rootDepth > 0 {
				rootDepth--
			}

			return fn(v, depth, rootDepth)
		}
	}

	v.WalkPreOrder(makeWalkFn(func(v *Value, depth int, rootDepth int) error {
		if v.IsRoot {
			addrIndentWidth := rootDepth + digitsInBase(bitio.BitsByteCount(v.BitBuf.Len()), opts.AddrBase)
			if addrIndentWidth > maxAddrIndentWidth {
				maxAddrIndentWidth = addrIndentWidth
			}
		}
		return nil
	}))

	cw := columnwriter.New(w, []int{maxAddrIndentWidth, 1, opts.LineBytes*3 - 1, 1, opts.LineBytes, 1, -1})
	return v.WalkPreOrder(makeWalkFn(func(v *Value, depth int, rootDepth int) error {
		return v.dump(cw, depth, rootDepth, maxAddrIndentWidth-rootDepth, opts)
	}))

}
