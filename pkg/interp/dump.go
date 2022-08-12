package interp

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/wader/fq/internal/ansi"
	"github.com/wader/fq/internal/asciiwriter"
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
	opts        Options
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

	columns := func() {
		cprint(1, deco.Column, "\n")
		cprint(3, deco.Column, "\n")
		cprint(5, deco.Column, "\n")
	}

	isInArray := false
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
		columns()
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

	// show address bar on root, nested root and format change
	if depth == 0 || v.IsRoot || v.Format != nil {
		if !isCompound(v) {
			columns()
		}
		cfmt(colHex, "%s", deco.DumpHeader.F(ctx.hexHeader))
		cfmt(colASCII, "%s", deco.DumpHeader.F(ctx.asciiHeader))
		if !isCompound(v) {
			cw.Flush()
		}
	}

	cfmt(colField, "%s%s", indent, name)
	if isInArray {
		cfmt(colField, "%s%s%s", deco.Index.F("["), deco.Number.F(strconv.Itoa(v.Index)), deco.Index.F("]"))
	}

	var valueErr error

	// TODO: cleanup map[string]any []any or json format
	// dump should use some internal interface instead?
	switch vv := v.V.(type) {
	case *decode.Compound:
		if vv.IsArray {
			cfmt(colField, "%s%s:%s%s", deco.Index.F("["), deco.Number.F("0"), deco.Number.F(strconv.Itoa(len(vv.Children))), deco.Index.F("]"))
		} else {
			cfmt(colField, "%s", deco.Object.F("{}"))
		}
		cprint(colField, ":")
		if isInArray {
			cfmt(colField, " %s", v.Name)
		}
		if vv.Description != "" {
			cfmt(colField, " %s", deco.Value.F(vv.Description))
		}
	case *scalar.S:
		switch av := vv.Actual.(type) {
		case map[string]any:
			cfmt(colField, ": %s", deco.Object.F("{}"))
		case []any:
			// TODO: format?
			cfmt(colField, ": %s%s:%s%s", deco.Index.F("["), deco.Number.F("0"), deco.Number.F(strconv.Itoa(len(av))), deco.Index.F("]"))
		default:
			cprint(colField, ":")
			if vv.Sym == nil {
				cfmt(colField, " %s", deco.ValueColor(vv.Actual).F(previewValue(vv.Actual, vv.ActualDisplay)))
			} else {
				cfmt(colField, " %s", deco.ValueColor(vv.Sym).F(previewValue(vv.Sym, vv.SymDisplay)))
				cfmt(colField, " (%s)", deco.ValueColor(vv.Actual).F(previewValue(vv.Actual, vv.ActualDisplay)))
			}
		}

		if opts.Verbose && isInArray {
			cfmt(colField, " %s", v.Name)
		}
		if vv.Description != "" {
			cfmt(colField, " (%s)", deco.Value.F(vv.Description))
		}
	default:
		panic(fmt.Sprintf("unreachable vv %#+v", vv))
	}

	if v.Format != nil {
		cfmt(colField, " (%s)", deco.Value.F(v.Format.Name))
	}
	valueErr = v.Err

	innerRange := v.InnerRange()

	if opts.Verbose {
		cfmt(colField, " %s (%s)",
			mathex.BitRange(innerRange).StringByteBits(opts.Addrbase), mathex.Bits(innerRange.Len).StringByteBits(opts.Sizebase))
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
				columns()
				cfmt(colField, "%s  %s: %s: %s\n", indent, deco.Error.F("error"), formatErr.Format.Name, formatErr.Err.Error())

				if opts.Verbose {
					for _, f := range formatErr.Stacktrace.Frames() {
						columns()
						cfmt(colField, "%s    %s\n", indent, f.Function)
						columns()
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
				columns()
				cfmt(colField, "%s!%s\n", indent, deco.Error.F(err.Error()))
			}
		}

		printErrs(depth, valueErr)
	}

	rootBitLen, err := bitioex.Len(rootV.RootReader)
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

	columns()

	// has length and is not compound or a collapsed struct/array (max depth)
	if innerRange.Len > 0 && (!isCompound(v) || (opts.Depth != 0 && opts.Depth == depth)) {
		cfmt(colAddr, "%s%s\n",
			rootIndent, deco.DumpAddr.F(mathex.PadFormatInt(startLineByte, opts.Addrbase, true, addrWidth)))

		vBR, err := bitioex.Range(rootV.RootReader, startByte*8, displaySizeBits)
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
		if _, err := bitioex.CopyBitsBuffer(
			hexpairwriter.New(cw.Columns[colHex], opts.LineBytes, int(startLineByteOffset), hexpairFn),
			hexBR,
			buf); err != nil {
			return err
		}

		asciiBR, err := bitio.CloneReadSeeker(vBR)
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
			columns()
			cfmt(colAddr, "%s%s\n", rootIndent, deco.DumpAddr.F(mathex.PadFormatInt(lineStartByte, opts.Addrbase, true, addrWidth)))
		}
		// TODO: correct? should rethink columnwriter api maybe?
		lastLineStopByte := startLineByte + addrLines*int64(opts.LineBytes) - 1
		if lastDisplayByte == bufferLastByte && lastDisplayByte != lastLineStopByte {
			// extra "|" in as EOF markers
			cfmt(colHex, "%s\n", deco.Column)
			cfmt(colASCII, "%s\n", deco.Column)
		}

		if stopByte != lastDisplayByte {
			isEnd := ""
			if stopBit == bufferLastBit {
				isEnd = " (end)"
			}
			columns()

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

func dump(v *decode.Value, w io.Writer, opts Options) error {
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

	cw := columnwriter.New(
		w,
		[]int{
			maxAddrIndentWidth,
			1,
			opts.LineBytes*3 - 1,
			1,
			opts.LineBytes,
			1,
			-1,
		})
	buf := make([]byte, 32*1024)

	if opts.Color {
		cw.DisplayLenFn = ansi.Len
		cw.DisplayTruncateFn = ansi.Truncate
	}

	var hexHeader string
	var asciiHeader string
	for i := 0; i < opts.LineBytes; i++ {
		s := mathex.PadFormatInt(int64(i), opts.Addrbase, false, 2)
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

func hexdump(w io.Writer, bv Binary, opts Options) error {
	br, err := bitioex.Range(bv.br, bv.r.Start, bv.r.Len)
	if err != nil {
		return err
	}

	// TODO: hack
	opts.Verbose = true
	return dump(
		&decode.Value{
			// TODO: hack
			V:          &scalar.S{Actual: br},
			Range:      bv.r,
			RootReader: bv.br,
		},
		w,
		opts,
	)
}
