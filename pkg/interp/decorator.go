package interp

import (
	"fq/internal/ansi"
	"fq/pkg/decode"
)

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
		for i := 0; i < 256; i++ {
			byteColors[byte(i)] = ansi.FromString(opts.ByteColors[byte(i)])
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
