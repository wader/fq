package interp

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/wader/fq/internal/ansi"
	"github.com/wader/fq/internal/asciiwriter"
	"github.com/wader/fq/internal/bitiox"
	"github.com/wader/fq/internal/columnwriter"
	"github.com/wader/fq/internal/hexpairwriter"
	"github.com/wader/fq/internal/mathx"
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
		cfmt(colHex, "%s", deco.DumpHeader.F(ctx.hexHeader))
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
	isSynthetic := false

	switch vv := v.V.(type) {
	case *decode.Compound:
		if vv.IsArray {
			cfmt(colField, "%s%s:%s%s", deco.Index.F("["), deco.Number.F("0"), deco.Number.F(strconv.Itoa(len(vv.Children))), deco.Index.F("]"))
		} else {
			cfmt(colField, "%s", deco.Object.F("{}"))
		}
		cprint(colField, ":")
		desc = vv.Description

	case scalar.Scalarable:
		cprint(colField, ":")
		actual := vv.ScalarActual()
		sym := vv.ScalarSym()
		df := vv.ScalarDisplayFormat()
		if sym == nil {
			cfmt(colField, " %s", deco.ValueColor(actual).F(previewValue(actual, df, opts)))
		} else {
			cfmt(colField, " %s", deco.ValueColor(sym).F(previewValue(sym, scalar.NumberDecimal, opts)))
			cfmt(colField, " (%s)", deco.ValueColor(actual).F(previewValue(actual, df, opts)))
		}
		desc = vv.ScalarDescription()
		isSynthetic = vv.ScalarFlags().IsSynthetic()
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

	if opts.Verbose && !isSynthetic {
		cfmt(colField, " %s (%s)",
			mathx.BitRange(innerRange).StringByteBits(opts.Addrbase), mathx.Bits(innerRange.Len).StringByteBits(opts.Sizebase))
	}

	cprint(colField, "\n")

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

	rootBitLen, err := bitiox.Len(rootV.RootReader)
	if err != nil {
		return err
	}

	bufferLastBit := rootBitLen - 1
	startBit := innerRange.Start
	stopBit := innerRange.Stop() - 1
	sizeBits := innerRange.Len
	lastDisplayBit := stopBit

	if opts.DisplayBytes > 0 && sizeBits > int64(opts.DisplayBytes)*8 {
		lastDisplayBit = startBit + (int64(opts.DisplayBytes)*8 - 1)
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

	// has length and is not compound or a collapsed struct/array (max depth)
	if willDisplayData {
		cfmt(colAddr, "%s%s\n",
			rootIndent, deco.DumpAddr.F(mathx.PadFormatInt(startLineByte, opts.Addrbase, true, addrWidth)))

		vBR, err := bitiox.Range(rootV.RootReader, startByte*8, displaySizeBits)
		if err != nil {
			return err
		}

		addrLines := lastDisplayLine - startLine + 1
		hexpairFn := func(b byte) string { return deco.ByteColor(b).Wrap(hexpairwriter.Pair(b)) }
		asciiFn := func(b byte) string { return deco.ByteColor(b).Wrap(asciiwriter.SafeASCII(b)) }

		hexBR, err := bitio.CloneReadSeeker(vBR)
		if err != nil {
			return err
		}
		if _, err := bitiox.CopyBitsBuffer(
			hexpairwriter.New(cw.Columns[colHex], opts.LineBytes, int(startLineByteOffset), hexpairFn),
			hexBR,
			buf); err != nil {
			return err
		}

		asciiBR, err := bitio.CloneReadSeeker(vBR)
		if err != nil {
			return err
		}
		if _, err := bitiox.CopyBitsBuffer(
			asciiwriter.New(cw.Columns[colASCII], opts.LineBytes, int(startLineByteOffset), asciiFn),
			asciiBR,
			buf); err != nil {
			return err
		}

		for i := int64(1); i < addrLines; i++ {
			lineStartByte := startLineByte + i*int64(opts.LineBytes)
			cfmt(colAddr, "%s%s\n", rootIndent, deco.DumpAddr.F(mathx.PadFormatInt(lineStartByte, opts.Addrbase, true, addrWidth)))
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
				mathx.Bits(stopBit).StringByteBits(opts.Addrbase),
				isEnd,
				mathx.PadFormatInt(bitio.BitsByteCount(sizeBits), opts.Sizebase, true, 0))
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
		maxAddrIndentWidth = max(
			maxAddrIndentWidth,
			rootIndentWidth*rootDepth+mathx.DigitsInBase(bitio.BitsByteCount(v.InnerRange().Stop()), true, opts.Addrbase),
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
	hexColumnWidth := opts.LineBytes*3 - 1
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
	for i := 0; i < opts.LineBytes; i++ {
		s := mathx.PadFormatInt(int64(i), opts.Addrbase, false, 2)
		hexHeader += s
		if i < opts.LineBytes-1 {
			hexHeader += " "
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
	br, err := bitiox.Range(bv.br, bv.r.Start, bv.r.Len)
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
