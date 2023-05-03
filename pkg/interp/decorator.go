package interp

import (
	"fmt"
	"math/big"

	"github.com/wader/fq/internal/ansi"
	"github.com/wader/fq/pkg/bitio"
)

var PlainDecorator = Decorator{
	Column:     "|",
	ValueColor: func(v any) ansi.Code { return ansi.None },
	ByteColor:  func(b byte) ansi.Code { return ansi.None },
}

func decoratorFromOptions(opts Options) Decorator {
	d := PlainDecorator

	if opts.Unicode {
		// U+2502 Box Drawings Light Vertical
		d.Column = "│"
	}

	if opts.Color {
		colors := opts.Colors

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

		d.ValueColor = func(v any) ansi.Code {
			switch vv := v.(type) {
			case bool:
				if vv {
					return d.True
				}
				return d.False
			case string,
				bitio.Reader,
				Binary:
				return d.String
			case int, float64, int64, uint64:
				// TODO: clean up number types
				return d.Number
			case *big.Int:
				return d.Number
			case []any:
				return d.Array
			case map[string]any:
				return d.Object
			case nil:
				return d.Null

			default:
				panic(fmt.Sprintf("unreachable %v (%T)", v, v))
			}
		}
		byteDefaultColor := ansi.FromString("")
		byteColors := map[byte]ansi.Code{}
		for i := 0; i < 256; i++ {
			byteColors[byte(i)] = byteDefaultColor
		}
		for _, sr := range opts.ByteColors {
			c := ansi.FromString(sr.Value)
			for _, r := range sr.Ranges {
				for i := r[0]; i <= r[1]; i++ {
					byteColors[byte(i)] = c
				}
			}
		}
		d.ByteColor = func(b byte) ansi.Code { return byteColors[b] }
	} else {
		d.ValueColor = func(v any) ansi.Code { return ansi.None }
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

	ValueColor func(v any) ansi.Code
	ByteColor  func(b byte) ansi.Code

	Column string
}
