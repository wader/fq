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

// 0-255:brightwhite,0:brightblack,32-126+9-13:white
func ParseStringRanges(s string) []stringRanges {
	var srs []stringRanges

	for _, stringRangesStr := range strings.Split(s, ",") {
		var rs [][2]int

		stringRangesParts := strings.Split(stringRangesStr, ":")
		if len(stringRangesParts) != 2 {
			continue
		}

		for _, rangeStr := range strings.Split(stringRangesParts[0], "+") {
			var err error
			rangeStrParts := strings.SplitN(rangeStr, "-", 2)
			start := 0
			stop := 0

			if len(rangeStrParts) == 1 {
				start, err = strconv.Atoi(rangeStrParts[0])
				if err != nil {
					continue
				}
				stop = start
			} else {
				start, err = strconv.Atoi(rangeStrParts[0])
				if err != nil {
					continue
				}
				stop, err = strconv.Atoi(rangeStrParts[1])
				if err != nil {
					continue
				}
			}

			rs = append(rs, [2]int{start, stop})
		}

		srs = append(srs, stringRanges{rs: rs, s: stringRangesParts[1]})
	}

	return srs
}

func decoratorFromDumpOptions(opts DisplayOptions) Decorator {
	colStr := "|"
	if opts.Unicode {
		colStr = "\xe2\x94\x82"
	}

	deco := Decorator{
		Column: colStr,
	}

	if opts.Color {
		deco.Null = ansi.FromString(opts.Colors["null"])
		deco.False = ansi.FromString(opts.Colors["false"])
		deco.True = ansi.FromString(opts.Colors["true"])
		deco.Number = ansi.FromString(opts.Colors["number"])
		deco.String = ansi.FromString(opts.Colors["string"])
		deco.ObjectKey = ansi.FromString(opts.Colors["objectkey"])
		deco.Array = ansi.FromString(opts.Colors["array"])
		deco.Object = ansi.FromString(opts.Colors["object"])

		deco.Index = ansi.FromString(opts.Colors["index"])

		deco.Value = ansi.FromString(opts.Colors["value"])
		deco.Frame = ansi.FromString(opts.Colors["frame"])

		deco.Error = ansi.FromString(opts.Colors["error"])

		deco.ValueColor = func(v *decode.Value) ansi.Color {
			switch vv := v.V.(type) {
			case decode.Array:
				return deco.Array
			case decode.Struct:
				return deco.Object
			case bool:
				if vv {
					return deco.True
				}
				return deco.False
			case string:
				return deco.String
			case nil:
				return deco.Null
			case int, float64, int64, uint64:
				// TODO: clean up number types
				return deco.Number
			default:
				// TODO: error?
				return deco.Value
			}
		}

		byteColors := map[byte]ansi.Color{}
		byteDefaultColor := ansi.FromString("")
		for i := 0; i < 256; i++ {
			byteColors[byte(i)] = byteDefaultColor
		}
		for _, sr := range ParseStringRanges(opts.Colors["bytes"]) {
			c := ansi.FromString(sr.s)
			for _, r := range sr.rs {
				for i := r[0]; i <= r[1]; i++ {
					byteColors[byte(i)] = c
				}
			}
		}
		deco.ByteColor = func(b byte) ansi.Color { return byteColors[b] }
	} else {
		deco.ValueColor = func(v *decode.Value) ansi.Color { return ansi.FromString("") }
		deco.ByteColor = func(b byte) ansi.Color { return ansi.FromString("") }
	}

	return deco
}

type Decorator struct {
	Null      ansi.Color
	False     ansi.Color
	True      ansi.Color
	Number    ansi.Color
	String    ansi.Color
	ObjectKey ansi.Color
	Array     ansi.Color
	Object    ansi.Color

	Index ansi.Color

	Name  ansi.Color
	Value ansi.Color
	Frame ansi.Color

	Error ansi.Color

	ValueColor func(v *decode.Value) ansi.Color

	ByteColor func(b byte) ansi.Color

	Column string
}
