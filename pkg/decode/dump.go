package decode

import (
	"fmt"
	"fq/internal/ansi"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"fq/pkg/bitio"
	"io"
	"strings"
)

type DumpOptions struct {
	MaxDepth int
	Verbose  bool

	LineBytes       int
	MaxDisplayBytes int64
	AddrBase        int
	SizeBase        int
}

func (v *Value) dump(cw *columnwriter.Writer, depth int, rootV *Value, rootDepth int, addrWidth int, opts DumpOptions) error {
	isInArray := false
	if v.Parent != nil {
		_, isInArray = v.Parent.V.(Array)
	}

	nameV := v
	name := nameV.Name
	if isInArray {
		nameV = v.Parent
		name = ""
	}
	if depth == 0 {
		name = nameV.Path()
	}
	if isInArray {
		name += fmt.Sprintf("[%d]", v.Index)
	}

	rootIndent := strings.Repeat(" ", rootDepth)
	indent := strings.Repeat("  ", depth)

	isField := false

	fmt.Fprint(cw.Columns[6], indent, name)

	switch vv := v.V.(type) {
	case Struct:
		fmt.Fprint(cw.Columns[6], "{}:")
		if v.Description != "" {
			fmt.Fprint(cw.Columns[6], " ", v.Description)
		}
	case Array:
		fmt.Fprintf(cw.Columns[6], "[%d]:", len(vv))
	default:
		fmt.Fprintf(cw.Columns[6], ": %s", v)
		isField = true
	}
	if opts.Verbose {
		fmt.Fprintf(cw.Columns[6], " (%s)", v.Name)
	}

	if opts.Verbose {
		fmt.Fprintf(cw.Columns[6], " %s (%s)",
			BitRange(v.Range).StringByteBits(opts.AddrBase), Bits(v.Range.Len).StringByteBits(opts.SizeBase))
	}

	if v.Error != nil {
		fmt.Fprintf(cw.Columns[6], "%s!%s\n", indent, v.Error)
		fmt.Fprintf(cw.Columns[1], "|\n")
		fmt.Fprintf(cw.Columns[3], "|\n")
		fmt.Fprintf(cw.Columns[5], "|\n")
	}

	bufferLastBit := rootV.RootBitBuf.Len() - 1
	startBit := v.Range.Start
	stopBit := v.Range.Stop() - 1
	sizeBits := v.Range.Len
	lastDisplayBit := stopBit

	if opts.MaxDisplayBytes > 0 && sizeBits > opts.MaxDisplayBytes*8 {
		lastDisplayBit = startBit + (opts.MaxDisplayBytes*8 - 1)
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
	displaySizeBits := displaySizeBytes * 8
	maxDisplaySizeBits := bufferLastBit - startByte*8 + 1
	if sizeBits == 0 {
		displaySizeBits = 0
	}
	if displaySizeBits > maxDisplaySizeBits {
		displaySizeBits = maxDisplaySizeBits
	}

	startLine := startByte / int64(opts.LineBytes)
	startLineByteOffset := startByte % int64(opts.LineBytes)
	startLineByte := startLine * int64(opts.LineBytes)
	lastDisplayLine := lastDisplayByte / int64(opts.LineBytes)

	fmt.Fprintf(cw.Columns[1], "|\n")
	fmt.Fprintf(cw.Columns[3], "|\n")
	fmt.Fprintf(cw.Columns[5], "|\n")

	if isField || (opts.MaxDepth != 0 && opts.MaxDepth == depth) {
		fmt.Fprintf(cw.Columns[0], "%s%s\n",
			rootIndent, num.PadFormatInt(startLineByte, opts.AddrBase, addrWidth))

		color := false
		vBitBuf, err := rootV.RootBitBuf.BitBufRange(startByte*8, displaySizeBits)
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
			s := hexpairwriter.Pair(c)
			if color {
				return fmt.Sprintf("%s%s%s", charToANSI(c), s, ansi.Reset)

			}
			return s
		}
		asciiFn := func(c byte) string {
			if color {
				return fmt.Sprintf("%s%s%s", charToANSI(c), asciiwriter.SafeASCII(c), ansi.Reset)
			}
			return asciiwriter.SafeASCII(c)
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
			fmt.Fprintf(cw.Columns[0], "%s%s\n", rootIndent, num.PadFormatInt(lineStartByte, opts.AddrBase, addrWidth))
			fmt.Fprintf(cw.Columns[1], "|\n")
			fmt.Fprintf(cw.Columns[3], "|\n")
			fmt.Fprintf(cw.Columns[5], "|\n")
		}
		// TODO: correct? should rethink columnwriter api maybe?
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
			fmt.Fprintf(cw.Columns[0], "%s*\n", rootIndent)
			fmt.Fprintf(cw.Columns[1], "|\n")
			fmt.Fprint(cw.Columns[2], "\n")
			fmt.Fprintf(cw.Columns[2], "%d bytes more until %s%s",
				stopByte-lastDisplayByte, Bits(stopBit).StringByteBits(opts.AddrBase), isEnd)
			fmt.Fprintf(cw.Columns[3], "|\n")
			fmt.Fprintf(cw.Columns[5], "|\n")
			// TODO: dump last line?
		}
	}

	cw.Flush()

	return nil
}

func (v *Value) Dump(w io.Writer, opts DumpOptions) error {
	maxAddrIndentWidth := num.DigitsInBase(bitio.BitsByteCount(v.RootBitBuf.Len()), opts.AddrBase)
	makeWalkFn := func(fn WalkFn) WalkFn {
		return func(v *Value, rootV *Value, depth int, rootDepth int) error {
			if opts.MaxDepth != 0 && depth > opts.MaxDepth {
				return ErrWalkSkipChildren
			}
			// skip first root level
			if rootDepth > 0 {
				rootDepth--
			}

			return fn(v, rootV, depth, rootDepth)
		}
	}

	v.WalkPreOrder(makeWalkFn(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		if v.IsRoot {
			addrIndentWidth := rootDepth + num.DigitsInBase(bitio.BitsByteCount(v.RootBitBuf.Len()), opts.AddrBase)
			if addrIndentWidth > maxAddrIndentWidth {
				maxAddrIndentWidth = addrIndentWidth
			}
		}
		return nil
	}))

	cw := columnwriter.New(w, []int{maxAddrIndentWidth, 1, opts.LineBytes*3 - 1, 1, opts.LineBytes, 1, -1})
	return v.WalkPreOrder(makeWalkFn(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		return v.dump(cw, depth, rootV, rootDepth, maxAddrIndentWidth-rootDepth, opts)
	}))

}
