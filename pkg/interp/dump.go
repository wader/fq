package interp

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/wader/fq/internal/ansi"
	"github.com/wader/fq/internal/asciiwriter"
	"github.com/wader/fq/internal/binwriter"
	"github.com/wader/fq/internal/bitioex"
	"github.com/wader/fq/internal/columnwriter"
	"github.com/wader/fq/internal/hexpairwriter"
	"github.com/wader/fq/internal/mathex"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/scalar"
)

// TODO: refactor this
// move more things to jq?
// select columns?
// smart line wrap instead of truncate?
// binary/octal/... dump?

// 0   12      34    56
// addr|hexdump|ascii|field
const (
	colAddr  = 0
	colHex   = 2
	colASCII = 4
	colField = 6
)

const rootIndentWidth = 2
const treeIndentWidth = 2

func isCompound(v *decode.Value) bool {
	switch v.V.(type) {
	case *decode.Compound:
		return true
	default:
		return false
	}
}

type dumpCtx struct {
	opts        *Options
	buf         []byte
	cw          *columnwriter.Writer
	hexHeader   string
	asciiHeader string
}

func indentStr(n int) string {
	const spaces = "                                                                "
	for n > len(spaces) {
		return strings.Repeat(" ", n)
	}
	return spaces[0:n]
}

func dumpEx(v *decode.Value, ctx *dumpCtx, depth int, rootV *decode.Value, rootDepth int, addrWidth int) error {
	opts := ctx.opts
	cw := ctx.cw
	buf := ctx.buf
	deco := ctx.opts.Decorator

	// no error check as we write into buffering column
	// we check for err later for Flush()
	cprint := func(c int, a ...any) {
		fmt.Fprint(cw.Columns[c], a...)
	}
	// cfmt: column i fmt.fprintf
	cfmt := func(c int, format string, a ...any) {
		fmt.Fprintf(cw.Columns[c], format, a...)
	}

	isInArray := false
	isCompound := isCompound(v)
	inArrayLen := 0
	if v.Parent != nil {
		if dc, ok := v.Parent.V.(*decode.Compound); ok {
			isInArray = dc.IsArray
			inArrayLen = len(dc.Children)
		}
	}

	nameV := v
	name := nameV.Name
	if isInArray {
		nameV = v.Parent
		name = ""
	}
	if depth == 0 {
		name = valuePathExprDecorated(nameV, deco)
	} else {
		name = deco.ObjectKey.Wrap(name)
	}

	rootIndent := indentStr(rootIndentWidth * rootDepth)
	indent := indentStr(treeIndentWidth * depth)

	if opts.ArrayTruncate != 0 && depth != 0 && isInArray && v.Index >= opts.ArrayTruncate {
		cfmt(colField, "%s%s%s:%s%s: ...",
			indent,
			deco.Index.F("["),
			deco.Number.F(strconv.Itoa(v.Index)),
			deco.Number.F(strconv.Itoa(inArrayLen)),
			deco.Index.F("]"),
		)
		cw.Flush()
		return decode.ErrWalkBreak
	}

	innerRange := v.InnerRange()
	willDisplayData := innerRange.Len > 0 && (!isCompound || (opts.Depth != 0 && opts.Depth == depth))

	// show address bar on root, nested root and format change
	if depth == 0 || v.IsRoot || v.Format != nil {
		// write header: 00 01 02 03 04
		cfmt(colHex, "%s", deco.DumpHeader.F(ctx.hexHeader))
		// write header: 012345
		cfmt(colASCII, "%s", deco.DumpHeader.F(ctx.asciiHeader))

		if willDisplayData {
			cw.Flush()
		}
	}

	cfmt(colField, "%s%s", indent, name)
	if isInArray {
		cfmt(colField, "%s%s%s", deco.Index.F("["), deco.Number.F(strconv.Itoa(v.Index)), deco.Index.F("]"))
	}

	var desc string

	switch vv := v.V.(type) {
	case *decode.Compound:
		if vv.IsArray {
			cfmt(colField, "%s%s:%s%s", deco.Index.F("["), deco.Number.F("0"), deco.Number.F(strconv.Itoa(len(vv.Children))), deco.Index.F("]"))
		} else {
			cfmt(colField, "%s", deco.Object.F("{}"))
		}
		cprint(colField, ":")
		desc = vv.Description

	case Scalarable:
		cprint(colField, ":")
		actual := vv.ScalarActual()
		sym := vv.ScalarSym()
		df := vv.ScalarDisplayFormat()
		if sym == nil {
			cfmt(colField, " %s", deco.ValueColor(actual).F(previewValue(actual, df)))
		} else {
			cfmt(colField, " %s", deco.ValueColor(sym).F(previewValue(sym, scalar.NumberDecimal)))
			cfmt(colField, " (%s)", deco.ValueColor(actual).F(previewValue(actual, df)))
		}
		desc = vv.ScalarDescription()

	default:
		panic(fmt.Sprintf("unreachable vv %#+v", vv))
	}

	if isCompound {
		if isInArray {
			cfmt(colField, " %s", v.Name)
		}
		if desc != "" {
			cfmt(colField, " %s", deco.Value.F(desc))
		}
	} else {
		if opts.Verbose && isInArray {
			cfmt(colField, " %s", v.Name)
		}
		if desc != "" {
			cfmt(colField, " (%s)", deco.Value.F(desc))
		}
	}

	if v.Format != nil {
		cfmt(colField, " (%s)", deco.Value.F(v.Format.Name))
	}
	valueErr := v.Err

	if opts.Verbose {
		cfmt(colField, " %s (%s)",
			mathex.BitRange(innerRange).StringByteBits(opts.Addrbase), mathex.Bits(innerRange.Len).StringByteBits(opts.Sizebase))
	}

	cprint(colField, "\n")

	// --------------------------------------------------
	// Error handling
	// --------------------------------------------------

	if valueErr != nil {
		var printErrs func(depth int, err error)
		printErrs = func(depth int, err error) {
			indent := indentStr(treeIndentWidth * depth)

			var formatErr decode.FormatError
			var decodeFormatsErr decode.FormatsError

			switch {
			case errors.As(err, &formatErr):
				cfmt(colField, "%s  %s: %s: %s\n", indent, deco.Error.F("error"), formatErr.Format.Name, formatErr.Err.Error())

				if opts.Verbose {
					for _, f := range formatErr.Stacktrace.Frames() {
						cfmt(colField, "%s    %s\n", indent, f.Function)
						cfmt(colField, "%s      %s:%d\n", indent, f.File, f.Line)
					}
				}
				switch {
				case errors.Is(formatErr.Err, decode.FormatsError{}):
					printErrs(depth+1, formatErr.Err)
				}
			case errors.As(err, &decodeFormatsErr):
				cfmt(colField, "%s  %s\n", indent, err)
				for _, e := range decodeFormatsErr.Errs {
					printErrs(depth+1, e)
				}
			default:
				cfmt(colField, "%s!%s\n", indent, deco.Error.F(err.Error()))
			}
		}

		printErrs(depth, valueErr)
	}

	// --------------------------------------------------
	// For a given field, compute various helper variables
	// --------------------------------------------------

	rootBitLen, err := bitioex.Len(rootV.RootReader)
	if err != nil {
		return err
	}

	bufferLastBit := rootBitLen - 1
	startBit := innerRange.Start     // field's start bit index (for entire file)
	stopBit := innerRange.Stop() - 1 // field's end bit index (for entire file); inclusive
	sizeBits := innerRange.Len       // field's bit length (1, 8, 16, 32, ...)

	// determine lastDisplayBit:
	// sometimes the field's bit length overflows the max width of a line;
	// cut off the overflow in such cases.
	lastDisplayBit := stopBit
	displayBits := int64(opts.DisplayBytes) * 8
	lineBits := int64(opts.LineBytes) * 8
	if opts.DisplayBytes > 0 && sizeBits > displayBits {
		lastDisplayBit = startBit + (displayBits - 1)
		if lastDisplayBit%lineBits != 0 {
			lastDisplayBit += lineBits - lastDisplayBit%lineBits - 1
		}

		if lastDisplayBit > stopBit || stopBit-lastDisplayBit <= lineBits {
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
	if opts.Base == 2 && displaySizeBits > stopBit-startBit {
		displaySizeBits = stopBit - startBit + 1 // TODO: -1 hmm
	}

	startLine := startByte / int64(opts.LineBytes)
	startLineByteOffset := startByte % int64(opts.LineBytes)
	startLineBitOffset := startBit % int64(opts.LineBytes*8)

	startLineByte := startLine * int64(opts.LineBytes)
	lastDisplayLine := lastDisplayByte / int64(opts.LineBytes)

	// --------------------------------------------------
	// Output Data
	// --------------------------------------------------

	// has length and is not compound or a collapsed struct/array (max depth)
	if willDisplayData {
		// write address: 0x00012 (example)
		cfmt(colAddr, "%s%s\n",
			rootIndent, deco.DumpAddr.F(mathex.PadFormatInt(startLineByte, opts.Addrbase, true, addrWidth)))

		vBR1, err := bitioex.Range(rootV.RootReader, startByte*8, displaySizeBits)
		if err != nil {
			return err
		}

		addrLines := lastDisplayLine - startLine + 1
		hexpairFn := func(b byte) string { return deco.ByteColor(b).Wrap(hexpairwriter.Pair(b)) }
		binFn := func(b byte) string { return deco.ByteColor(b).Wrap(string("01"[int(b)])) }
		asciiFn := func(b byte) string { return deco.ByteColor(b).Wrap(asciiwriter.SafeASCII(b)) }

		switch opts.Base {
		case 16:
			// write hex: 89 50 4e 47 0d 0a 1a 0a ...
			hexBR, err := bitio.CloneReadSeeker(vBR1)
			if err != nil {
				return err
			}
			if _, err := bitioex.CopyBitsBuffer(
				hexpairwriter.New(cw.Columns[colHex], opts.LineBytes, int(startLineByteOffset), hexpairFn),
				hexBR,
				buf); err != nil {
				return err
			}
		case 2:
			// write bits: 100010010101000...
			vBR2, err := bitioex.Range(rootV.RootReader, startBit, displaySizeBits)
			if err != nil {
				return err
			}
			hexBR, err := bitio.CloneReadSeeker(vBR2)
			if err != nil {
				return err
			}
			if _, err := bitio.CopyBuffer(
				binwriter.New(cw.Columns[colHex], opts.LineBytes*8, int(startLineBitOffset), binFn),
				hexBR,
				buf); err != nil {
				return err
			}
		}
		// write ascii: .PNG.........IHDR...
		asciiBR, err := bitio.CloneReadSeeker(vBR1)
		if err != nil {
			return err
		}
		if _, err := bitioex.CopyBitsBuffer(
			asciiwriter.New(cw.Columns[colASCII], opts.LineBytes, int(startLineByteOffset), asciiFn),
			asciiBR,
			buf); err != nil {
			return err
		}

		for i := int64(1); i < addrLines; i++ {
			lineStartByte := startLineByte + i*int64(opts.LineBytes)
			cfmt(colAddr, "%s%s\n", rootIndent, deco.DumpAddr.F(mathex.PadFormatInt(lineStartByte, opts.Addrbase, true, addrWidth)))
		}
		// TODO: correct? should rethink columnwriter api maybe?
		lastLineStopByte := startLineByte + addrLines*int64(opts.LineBytes) - 1
		if lastDisplayByte == bufferLastByte && lastDisplayByte != lastLineStopByte {
			// extra "|" as end markers
			cfmt(colHex, "%s\n", deco.Column)
			cfmt(colASCII, "%s\n", deco.Column)
		}

		if stopByte != lastDisplayByte {
			isEnd := ""
			if stopBit == bufferLastBit {
				isEnd = " (end)"
			}

			cfmt(colAddr, "%s%s\n", rootIndent, deco.DumpAddr.F("*"))
			cprint(colHex, "\n")
			// TODO: truncate if display_bytes is small?
			cfmt(colHex, "until %s%s (%s)",
				mathex.Bits(stopBit).StringByteBits(opts.Addrbase),
				isEnd,
				mathex.PadFormatInt(bitio.BitsByteCount(sizeBits), opts.Sizebase, true, 0))
			// TODO: dump last line?
		}
	}

	if err := cw.Flush(); err != nil {
		return err
	}

	return nil
}

func dump(v *decode.Value, w io.Writer, opts *Options) error {
	maxAddrIndentWidth := 0
	makeWalkFn := func(fn decode.WalkFn) decode.WalkFn {
		return func(v *decode.Value, rootV *decode.Value, depth int, rootDepth int) error {
			if opts.Depth != 0 && depth > opts.Depth {
				return decode.ErrWalkSkipChildren
			}

			return fn(v, rootV, depth, rootDepth)
		}
	}

	_ = v.WalkPreOrder(makeWalkFn(func(v *decode.Value, _ *decode.Value, _ int, rootDepth int) error {
		maxAddrIndentWidth = mathex.Max(
			maxAddrIndentWidth,
			rootIndentWidth*rootDepth+mathex.DigitsInBase(bitio.BitsByteCount(v.InnerRange().Stop()), true, opts.Addrbase),
		)
		return nil
	}))

	var displayLenFn func(s string) int
	var displayTruncateFn func(s string, start, stop int) string
	if opts.Color {
		displayLenFn = ansi.Len
		displayTruncateFn = ansi.Slice
	}

	addrColumnWidth := maxAddrIndentWidth
	var hexColumnWidth int
	switch opts.Base {
	case 16:
		hexColumnWidth = opts.LineBytes*3 - 1
	case 2:
		hexColumnWidth = opts.LineBytes * 8
	}
	asciiColumnWidth := opts.LineBytes
	treeColumnWidth := -1
	// TODO: set with and truncate/wrap properly
	// if opts.Width != 0 {
	// 	treeColumnWidth = mathex.Max(0, opts.Width-(addrColumnWidth+hexColumnWidth+asciiColumnWidth+3 /* bars */))
	// }

	cw := columnwriter.New(
		w,
		&columnwriter.MultiLineColumn{Width: addrColumnWidth, LenFn: displayLenFn, SliceFn: displayTruncateFn},
		columnwriter.BarColumn(opts.Decorator.Column),
		&columnwriter.MultiLineColumn{Width: hexColumnWidth, LenFn: displayLenFn, SliceFn: displayTruncateFn},
		columnwriter.BarColumn(opts.Decorator.Column),
		&columnwriter.MultiLineColumn{Width: asciiColumnWidth, LenFn: displayLenFn, SliceFn: displayTruncateFn},
		columnwriter.BarColumn(opts.Decorator.Column),
		&columnwriter.MultiLineColumn{Width: treeColumnWidth, Wrap: false, LenFn: displayLenFn, SliceFn: displayTruncateFn},
	)

	buf := make([]byte, 32*1024)

	var hexHeader string
	var asciiHeader string
	var spaceLength int
	switch opts.Base {
	case 16:
		spaceLength = 1
	case 2:
		spaceLength = 8 - 2 // TODO: adapt for wider screens
	}
	for i := 0; i < opts.LineBytes; i++ {
		s := mathex.PadFormatInt(int64(i), opts.Addrbase, false, 2)
		hexHeader += s
		if spaceLength > 1 || i < opts.LineBytes-1 {
			hexHeader += strings.Repeat(" ", spaceLength)
		}
		asciiHeader += s[len(s)-1:]
	}

	ctx := &dumpCtx{
		opts:        opts,
		buf:         buf,
		cw:          cw,
		hexHeader:   hexHeader,
		asciiHeader: asciiHeader,
	}

	return v.WalkPreOrder(makeWalkFn(func(v *decode.Value, rootV *decode.Value, depth int, rootDepth int) error {
		return dumpEx(v, ctx, depth, rootV, rootDepth, maxAddrIndentWidth-rootDepth)
	}))
}

func hexdump(w io.Writer, bv Binary, opts *Options) error {
	br, err := bitioex.Range(bv.br, bv.r.Start, bv.r.Len)
	if err != nil {
		return err
	}

	// TODO: hack
	opts.Verbose = true
	return dump(
		&decode.Value{
			// TODO: hack
			V:          &scalar.BitBuf{Actual: br},
			Range:      bv.r,
			RootReader: bv.br,
		},
		w,
		opts,
	)
}
