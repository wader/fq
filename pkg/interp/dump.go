package interp

import (
	"fmt"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"io"
	"strings"
)

type Decorator struct {
	Name   func(s string) string
	Value  func(s string) string
	Byte   func(b byte, s string) string
	Column string
}

type DisplayOptions struct {
	MaxDepth int
	Verbose  bool
	Color    bool
	Unicode  bool
	Raw      bool
	REPL     bool

	LineBytes    int
	DisplayBytes int64
	AddrBase     int
	SizeBase     int

	Decorator Decorator
}

// 0   12      34    56
// addr|hexdump|ascii|field
const (
	colAddr  = 0
	colHex   = 2
	colAscii = 4
	colField = 6
)

func dumpEx(v *decode.Value, cw *columnwriter.Writer, depth int, rootV *decode.Value, rootDepth int, addrWidth int, opts DisplayOptions) error {
	d := opts.Decorator
	// no error check as we write into buffering column
	// we check for err later for Flush()
	cprint := func(c int, a ...interface{}) {
		fmt.Fprint(cw.Columns[c], a...)
	}
	cfmt := func(c int, format string, a ...interface{}) {
		fmt.Fprintf(cw.Columns[c], format, a...)
	}

	columns := func() {
		cprint(1, d.Column)
		cprint(3, d.Column)
		cprint(5, d.Column)
	}

	if depth == 0 {
		for i := 0; i < opts.LineBytes; i++ {
			cprint(colHex, num.PadFormatInt(int64(i), opts.AddrBase, false, 2))
			if i < opts.LineBytes-1 {
				cprint(colHex, " ")
			}
		}
	}

	isInArray := false
	if v.Parent != nil {
		_, isInArray = v.Parent.V.(decode.Array)
	}

	nameV := v
	name := d.Name(nameV.Name)
	if isInArray {
		nameV = v.Parent
		name = ""
	}
	if depth == 0 {
		name = valuePath(nameV)
	}
	if isInArray {
		name += fmt.Sprintf("[%d]", v.Index)
	}

	rootIndent := strings.Repeat(" ", rootDepth)
	indent := strings.Repeat("  ", depth)

	isSimple := false

	switch vv := v.V.(type) {
	case decode.Struct:
		if depth == 0 {
			cprint(colField, name, ":")
		} else {
			cprint(colField, indent[1:], "-", name, ":")
		}
		if v.Description != "" {
			cprint(colField, " ", v.Description)
		}
	case decode.Array:
		cprint(colField, indent, name)
		cfmt(colField, "[%d]:", len(vv))
	default:
		cprint(colField, indent, name, ": ")

		if v.Symbol != "" {
			cprint(colField, d.Value(v.Symbol))
			cprint(colField, " (", previewValue(v), ")")
		} else {
			cprint(colField, d.Value(previewValue(v)))
		}
		if v.Description != "" {
			cprint(colField, fmt.Sprintf(" (%s)", v.Description))
		}

		isSimple = true
	}
	if opts.Verbose && isInArray {
		cfmt(colField, " (%s)", v.Name)
	}

	if opts.Verbose {
		cfmt(colField, " %s (%s)",
			decode.BitRange(v.Range).StringByteBits(opts.AddrBase), decode.Bits(v.Range.Len).StringByteBits(opts.SizeBase))
	}

	cprint(colField, "\n")

	if v.Error != nil {
		columns()
		cfmt(colField, "%s!%s\n", indent, v.Error)

		if opts.Verbose {
			if de, ok := v.Error.(*decode.DecodeError); ok && de.PanicStack != "" {
				ps := de.PanicStack
				for _, l := range strings.Split(ps, "\n") {
					columns()
					cfmt(colField, "%s%s\n", indent, l)
				}
			}
		}
	}

	bufferLastBit := rootV.RootBitBuf.Len() - 1
	startBit := v.Range.Start
	stopBit := v.Range.Stop() - 1
	sizeBits := v.Range.Len
	lastDisplayBit := stopBit

	if opts.DisplayBytes > 0 && sizeBits > opts.DisplayBytes*8 {
		lastDisplayBit = startBit + (opts.DisplayBytes*8 - 1)
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

	columns()

	// has length and is a simple value or a collapsed struct/array
	if v.Range.Len > 0 && (isSimple || (opts.MaxDepth != 0 && opts.MaxDepth == depth)) {
		cfmt(0, "%s%s\n",
			rootIndent, num.PadFormatInt(startLineByte, opts.AddrBase, true, addrWidth))

		vBitBuf, err := rootV.RootBitBuf.BitBufRange(startByte*8, displaySizeBits)
		if err != nil {
			return err
		}

		addrLines := lastDisplayLine - startLine + 1
		hexpairFn := func(b byte) string { return d.Byte(b, hexpairwriter.Pair(b)) }
		asciiFn := func(b byte) string { return d.Byte(b, asciiwriter.SafeASCII(b)) }

		if vBitBuf != nil {
			io.Copy(
				hexpairwriter.New(cw.Columns[colHex], opts.LineBytes, int(startLineByteOffset), hexpairFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes))
			io.Copy(
				asciiwriter.New(cw.Columns[colAscii], opts.LineBytes, int(startLineByteOffset), asciiFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes))
		}

		for i := int64(1); i < addrLines; i++ {
			lineStartByte := startLineByte + int64(i)*int64(opts.LineBytes)
			columns()
			cfmt(colAddr, "%s%s\n", rootIndent, num.PadFormatInt(lineStartByte, opts.AddrBase, true, addrWidth))
		}
		// TODO: correct? should rethink columnwriter api maybe?
		lastLineStopByte := startLineByte + int64(addrLines)*int64(opts.LineBytes) - 1
		if lastDisplayByte == bufferLastByte && lastDisplayByte != lastLineStopByte {
			// extra "|" in as EOF markers
			cfmt(colHex, "|\n")
			cfmt(colAscii, "|\n")
		}

		if stopByte != lastDisplayByte {
			isEnd := ""
			if stopBit == bufferLastBit {
				isEnd = " (end)"
			}
			columns()

			cfmt(colAddr, "%s*\n", rootIndent)
			cprint(colHex, "\n")
			// TODO: truncate if displaybytes is small?
			cfmt(colHex, "%s bytes more until %s%s",
				num.PadFormatInt(stopByte-lastDisplayByte, opts.SizeBase, true, 0),
				decode.Bits(stopBit).StringByteBits(opts.AddrBase),
				isEnd)
			// TODO: dump last line?
		}
	}

	if err := cw.Flush(); err != nil {
		return err
	}

	return nil
}

func dump(v *decode.Value, w io.Writer, opts DisplayOptions) error {
	maxAddrIndentWidth := 0
	makeWalkFn := func(fn decode.WalkFn) decode.WalkFn {
		return func(v *decode.Value, rootV *decode.Value, depth int, rootDepth int) error {
			if opts.MaxDepth != 0 && depth > opts.MaxDepth {
				return decode.ErrWalkSkipChildren
			}
			// skip first root level
			if rootDepth > 0 {
				rootDepth--
			}

			return fn(v, rootV, depth, rootDepth)
		}
	}

	v.WalkPreOrder(makeWalkFn(func(v *decode.Value, rootV *decode.Value, depth int, rootDepth int) error {
		maxAddrIndentWidth = num.MaxInt(
			maxAddrIndentWidth,
			rootDepth+num.DigitsInBase(bitio.BitsByteCount(v.Range.Stop()), true, opts.AddrBase),
		)
		return nil
	}))

	cw := columnwriter.New(w, []int{maxAddrIndentWidth, 1, opts.LineBytes*3 - 1, 1, opts.LineBytes, 1, -1})

	return v.WalkPreOrder(makeWalkFn(func(v *decode.Value, rootV *decode.Value, depth int, rootDepth int) error {
		return dumpEx(v, cw, depth, rootV, rootDepth, maxAddrIndentWidth-rootDepth, opts)
	}))
}
