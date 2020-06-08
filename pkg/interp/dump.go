package interp

import (
	"errors"
	"fmt"
	"fq/internal/asciiwriter"
	"fq/internal/columnwriter"
	"fq/internal/hexpairwriter"
	"fq/internal/num"
	"fq/pkg/bitio"
	"fq/pkg/decode"
	"io"
	"strconv"
	"strings"
)

// 0   12      34    56
// addr|hexdump|ascii|field
const (
	colAddr  = 0
	colHex   = 2
	colAscii = 4
	colField = 6
)

func dumpEx(v *decode.Value, cw *columnwriter.Writer, depth int, rootV *decode.Value, rootDepth int, addrWidth int, opts Options) error {
	deco := opts.Decorator
	// no error check as we write into buffering column
	// we check for err later for Flush()
	cprint := func(c int, a ...interface{}) {
		fmt.Fprint(cw.Columns[c], a...)
	}
	cfmt := func(c int, format string, a ...interface{}) {
		fmt.Fprintf(cw.Columns[c], format, a...)
	}

	columns := func() {
		cprint(1, deco.Column, "\n")
		cprint(3, deco.Column, "\n")
		cprint(5, deco.Column, "\n")
	}

	var header string
	if depth == 0 {
		for i := 0; i < opts.LineBytes; i++ {
			header += num.PadFormatInt(int64(i), opts.AddrBase, false, 2)
			if i < opts.LineBytes-1 {
				header += " "
			}
		}
	}

	isInArray := false
	if v.Parent != nil {
		_, isInArray = v.Parent.V.(decode.Array)
	}

	nameV := v
	name := nameV.Name
	if isInArray {
		nameV = v.Parent
		name = ""
	}
	if depth == 0 {
		name = valuePathDecorated(nameV, deco)
	} else {
		name = deco.ObjectKey.Wrap(name)
	}

	rootIndent := strings.Repeat(" ", rootDepth)
	indent := strings.Repeat("  ", depth)

	isSimple := false

	if depth == 0 {
		switch v.V.(type) {
		case decode.Struct:
			cfmt(colHex, "%s", deco.DumpHeader.F(header))
		case decode.Array:
			cfmt(colHex, "%s", deco.DumpHeader.F(header))
		default:
			columns()
			cfmt(colHex, "%s", deco.DumpHeader.F(header))
			cw.Flush()
		}
	}

	cfmt(colField, "%s%s", indent, name)
	if isInArray {
		cfmt(colField, "%s%s%s", deco.Index.F("["), deco.Number.F(strconv.Itoa(v.Index)), deco.Index.F("]"))
	}
	cprint(colField, ":")
	if opts.Verbose && isInArray {
		cfmt(colField, " %s", v.Name)
	}

	switch v.V.(type) {
	case decode.Struct, map[string]interface{}:
		cfmt(colField, " %s", deco.Object.F("{}"))
		if v.Description != "" {
			cfmt(colField, " %s", deco.Value.F(v.Description))
		}
		if v.Format != nil {
			cfmt(colField, " (%s)", deco.Value.F(v.Format.Name))
		}
	case decode.Array, []interface{}:
		var l int
		switch vv := v.V.(type) {
		case decode.Array:
			l = len(vv)
		case []interface{}:
			l = len(vv)
		default:
			panic("unreachable")
		}

		cfmt(colField, " %s%s%s", deco.Index.F("["), deco.Number.F(strconv.Itoa(l)), deco.Index.F("]"))
		if v.Description != "" {
			cfmt(colField, " %s", deco.Value.F(v.Description))
		}
		if v.Format != nil {
			cfmt(colField, " (%s)", deco.Value.F(v.Format.Name))
		}
	default:
		if v.Symbol != "" {
			cfmt(colField, " %s", deco.Value.F(v.Symbol))
			cfmt(colField, " (%s)", deco.ValueColor(v).F(previewValue(v)))
		} else {
			cfmt(colField, " %s", deco.ValueColor(v).F(previewValue(v)))
		}
		// TODO: similar to struct/array?
		if v.Description != "" {
			cfmt(colField, fmt.Sprintf(" (%s)", deco.Value.F(v.Description)))
		}

		isSimple = true
	}

	if opts.Verbose {
		cfmt(colField, " %s (%s)",
			num.BitRange(v.Range).StringByteBits(opts.AddrBase), num.Bits(v.Range.Len).StringByteBits(opts.SizeBase))
	}

	cprint(colField, "\n")

	if v.Err != nil {
		var printErrs func(depth int, err error)
		printErrs = func(depth int, err error) {
			indent := strings.Repeat("  ", depth)

			var formatErr decode.FormatError
			var decodeFormatsErr decode.DecodeFormatsError

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
				case errors.Is(formatErr.Err, decode.DecodeFormatsError{}):
					printErrs(depth+1, formatErr.Err)
				}
			case errors.As(err, &decodeFormatsErr):
				cfmt(colField, "%s  %s\n", indent, err)
				for _, e := range decodeFormatsErr.Errs {
					printErrs(depth+1, e)
				}
			default:
				columns()
				cfmt(colField, "%s!%s\n", indent, deco.Error.F(v.Err.Error()))
			}
		}

		printErrs(depth, v.Err)
	}

	bufferLastBit := rootV.RootBitBuf.Len() - 1
	startBit := v.Range.Start
	stopBit := v.Range.Stop() - 1
	sizeBits := v.Range.Len
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

	// log.Printf("v: %#+v\n", v)
	// log.Printf("isSimple: %#+v\n", isSimple)
	// log.Printf("opts: %#+v\n", opts)
	// log.Printf("depth: %#+v\n", depth)

	// has length and is a simple value or a collapsed struct/array
	if v.Range.Len > 0 && (isSimple || (opts.Depth != 0 && opts.Depth == depth)) {
		cfmt(colAddr, "%s%s\n",
			rootIndent, deco.DumpAddr.F(num.PadFormatInt(startLineByte, opts.AddrBase, true, addrWidth)))

		vBitBuf, err := rootV.RootBitBuf.BitBufRange(startByte*8, displaySizeBits)
		if err != nil {
			return err
		}

		addrLines := lastDisplayLine - startLine + 1
		hexpairFn := func(b byte) string { return deco.ByteColor(b).Wrap(hexpairwriter.Pair(b)) }
		asciiFn := func(b byte) string { return deco.ByteColor(b).Wrap(asciiwriter.SafeASCII(b)) }

		if vBitBuf != nil {
			if _, err := io.Copy(
				hexpairwriter.New(cw.Columns[colHex], opts.LineBytes, int(startLineByteOffset), hexpairFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes)); err != nil {
				return err
			}
			if _, err := io.Copy(
				asciiwriter.New(cw.Columns[colAscii], opts.LineBytes, int(startLineByteOffset), asciiFn),
				io.LimitReader(vBitBuf.Copy(), displaySizeBytes)); err != nil {
				return err
			}
		}

		for i := int64(1); i < addrLines; i++ {
			lineStartByte := startLineByte + i*int64(opts.LineBytes)
			columns()
			cfmt(colAddr, "%s%s\n", rootIndent, deco.DumpAddr.F(num.PadFormatInt(lineStartByte, opts.AddrBase, true, addrWidth)))
		}
		// TODO: correct? should rethink columnwriter api maybe?
		lastLineStopByte := startLineByte + addrLines*int64(opts.LineBytes) - 1
		if lastDisplayByte == bufferLastByte && lastDisplayByte != lastLineStopByte {
			// extra "|" in as EOF markers
			cfmt(colHex, "%s\n", deco.Column)
			cfmt(colAscii, "%s\n", deco.Column)
		}

		if stopByte != lastDisplayByte {
			isEnd := ""
			if stopBit == bufferLastBit {
				isEnd = " (end)"
			}
			columns()

			cfmt(colAddr, "%s%s\n", rootIndent, deco.DumpAddr.F("*"))
			cprint(colHex, "\n")
			// TODO: truncate if displaybytes is small?
			cfmt(colHex, "until %s%s (%s)",
				num.Bits(stopBit).StringByteBits(opts.AddrBase),
				isEnd,
				num.PadFormatInt(bitio.BitsByteCount(sizeBits), opts.SizeBase, true, 0))
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
			// skip first root level
			if rootDepth > 0 {
				rootDepth--
			}

			return fn(v, rootV, depth, rootDepth)
		}
	}

	_ = v.WalkPreOrder(makeWalkFn(func(v *decode.Value, rootV *decode.Value, depth int, rootDepth int) error {
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
