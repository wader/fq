package interp

import (
	"fq/internal/ansi"
	"fq/pkg/decode"
	"strconv"
	"strings"
)

type stringRanges struct {
	rs [][2]int
	s  string
}

func atoi(s string) int {
	n, _ := strconv.ParseUint(s, 0, 64)
	return int(n)
}

// 0-255=brightwhite,0=brightblack,32-126:9-13=white
func parseCSVRangeMap(s string) []stringRanges {
	var srs []stringRanges

	for _, stringRangesStr := range strings.Split(s, ",") {
		stringRangesStr = strings.TrimSpace(stringRangesStr)
		var rs [][2]int

		stringRangesParts := strings.Split(stringRangesStr, "=")
		if len(stringRangesParts) != 2 {
			continue
		}

		for _, rangeStr := range strings.Split(stringRangesParts[0], ":") {
			rangeStr = strings.TrimSpace(rangeStr)
			rangeStrParts := strings.SplitN(rangeStr, "-", 2)
			start := 0
			stop := 0

			if len(rangeStrParts) == 1 {
				start = atoi(rangeStrParts[0])
				stop = start
			} else {
				start = atoi(rangeStrParts[0])
				stop = atoi(rangeStrParts[1])
			}

			rs = append(rs, [2]int{start, stop})
		}

		srs = append(srs, stringRanges{rs: rs, s: stringRangesParts[1]})
	}

	return srs
}

// key=value,a=b,.. -> {"key": "value", "a": "b", ...}
func parseCSVStringMap(s string) map[string]string {
	m := map[string]string{}

	for _, stringKVStr := range strings.Split(s, ",") {
		stringKVStr = strings.TrimSpace(stringKVStr)
		stringKVParts := strings.Split(stringKVStr, "=")
		if len(stringKVParts) != 2 {
			continue
		}

		m[strings.TrimSpace(stringKVParts[0])] = strings.TrimSpace(stringKVParts[1])
	}

	return m
}

var PlainDecorator = Decorator{
	Column:     "|",
	ValueColor: func(v *decode.Value) ansi.Code { return ansi.None },
	ByteColor:  func(b byte) ansi.Code { return ansi.None },
}

func decoratorFromOptions(opts Options) Decorator {
	d := PlainDecorator

	if opts.Unicode {
		// U+2502 Box Drawings Light Vertical
		d.Column = "│"
	}

	if opts.Color {
		colors := parseCSVStringMap(opts.Colors)

		d.Null = ansi.FromString(colors["null"])
		d.False = ansi.FromString(colors["false"])
		d.True = ansi.FromString(colors["true"])
		d.Number = ansi.FromString(colors["number"])
		d.String = ansi.FromString(colors["string"])
		d.ObjectKey = ansi.FromString(colors["objectkey"])
		d.Array = ansi.FromString(colors["array"])
		d.Object = ansi.FromString(colors["object"])

		d.Index = ansi.FromString(colors["index"])
		d.Value = ansi.FromString(colors["value"])

		d.DumpHeader = ansi.FromString(colors["dumpheader"])
		d.DumpAddr = ansi.FromString(colors["dumpaddr"])

		d.Error = ansi.FromString(colors["error"])

		d.ValueColor = func(v *decode.Value) ansi.Code {
			switch vv := v.V.(type) {
			case decode.Array:
				return d.Array
			case decode.Struct:
				return d.Object
			case bool:
				if vv {
					return d.True
				}
				return d.False
			case string:
				return d.String
			case nil:
				return d.Null
			case int, float64, int64, uint64:
				// TODO: clean up number types
				return d.Number
			default:
				// TODO: error?
				return d.Value
			}
		}
		byteDefaultColor := ansi.FromString("")
		byteColors := map[byte]ansi.Code{}
		for i := 0; i < 256; i++ {
			byteColors[byte(i)] = byteDefaultColor
		}
		for _, sr := range parseCSVRangeMap(opts.ByteColors) {
			c := ansi.FromString(sr.s)
			for _, r := range sr.rs {
				for i := r[0]; i <= r[1]; i++ {
					byteColors[byte(i)] = c
				}
			}
		}
		d.ByteColor = func(b byte) ansi.Code { return byteColors[b] }
	} else {
		d.ValueColor = func(v *decode.Value) ansi.Code { return ansi.None }
		d.ByteColor = func(b byte) ansi.Code { return ansi.None }
	}

	return d
}

type Decorator struct {
	Null      ansi.Code
	False     ansi.Code
	True      ansi.Code
	Number    ansi.Code
	String    ansi.Code
	ObjectKey ansi.Code
	Array     ansi.Code
	Object    ansi.Code

	Index ansi.Code
	Value ansi.Code

	DumpHeader ansi.Code
	DumpAddr   ansi.Code

	Error ansi.Code

	ValueColor func(v *decode.Value) ansi.Code
	ByteColor  func(b byte) ansi.Code

	Column string
}
