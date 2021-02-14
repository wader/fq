package decode

import (
	"fmt"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"fq/pkg/bitio"
	"io"
	"strings"
)

type Decorator struct {
	Name   func(s string) string
	Value  func(s string) string
	Byte   func(b byte, s string) string
	Column string
}

type DumpOptions struct {
	MaxDepth int
	Verbose  bool
	Color    bool
	Unicode  bool

	LineBytes       int
	MaxDisplayBytes int64
	AddrBase        int
	SizeBase        int

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

func (v *Value) dump(cw *columnwriter.Writer, depth int, rootV *Value, rootDepth int, addrWidth int, opts DumpOptions) error {
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

	isInArray := false
	if v.Parent != nil {
		_, isInArray = v.Parent.V.(Array)
	}

	nameV := v
	name := d.Name(nameV.Name)
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

	isSimple := false

	switch vv := v.V.(type) {
	case Struct:
		if depth == 0 {
			cprint(colField, name, ":")
		} else {
			cprint(colField, indent[1:], "-", name, ":")
		}
		if v.Description != "" {
			cprint(colField, " ", v.Description)
		}
	case Array:
		cprint(colField, indent, name)
		cfmt(colField, "[%d]:", len(vv))
	default:
		cprint(colField, indent, name, ": ", d.Value(v.String()))
		isSimple = true
	}
	if opts.Verbose && isInArray {
		cfmt(colField, " (%s)", v.Name)
	}

	if opts.Verbose {
		cfmt(colField, " %s (%s)",
			BitRange(v.Range).StringByteBits(opts.AddrBase), Bits(v.Range.Len).StringByteBits(opts.SizeBase))
	}

	cprint(colField, "\n")

	if v.Error != nil {
		columns()
		cfmt(colField, "%s!%s\n", indent, v.Error)

		if opts.Verbose {
			if de, ok := v.Error.(*DecodeError); ok && de.PanicStack != "" {
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
				Bits(stopBit).StringByteBits(opts.AddrBase),
				isEnd)
			// TODO: dump last line?
		}
	}

	if err := cw.Flush(); err != nil {
		return err
	}

	return nil
}

func (v *Value) Dump(w io.Writer, opts DumpOptions) error {
	maxAddrIndentWidth := 0
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
		maxAddrIndentWidth = num.MaxInt(
			maxAddrIndentWidth,
			rootDepth+num.DigitsInBase(bitio.BitsByteCount(v.Range.Stop()), true, opts.AddrBase),
		)
		return nil
	}))

	cw := columnwriter.New(w, []int{maxAddrIndentWidth, 1, opts.LineBytes*3 - 1, 1, opts.LineBytes, 1, -1})
	return v.WalkPreOrder(makeWalkFn(func(v *Value, rootV *Value, depth int, rootDepth int) error {
		return v.dump(cw, depth, rootV, rootDepth, maxAddrIndentWidth-rootDepth, opts)
	}))
}
