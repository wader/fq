package interp

import (
	"fq/internal/ansi"
)

func decoratorFromDumpOptions(opts DisplayOptions) Decorator {
	colStr := "|"
	if opts.Unicode {
		colStr = "\xe2\x94\x82"
	}
	// nameFn := func(s string) string { return s }
	// valueFn := func(s string) string { return s }
	// frameFn := func(s string) string { return s }
	// byteFn := func(b byte, s string) string { return s }
	// column := colStr + "\n"
	// if opts.Color {
	// 	nameFn = func(s string) string { return ansi.FgBrightBlue + s + ansi.Reset }
	// 	valueFn = func(s string) string { return ansi.FgBrightCyan + s + ansi.Reset }
	// 	frameFn = func(s string) string { return ansi.FgYellow + s + ansi.Reset }
	// 	byteFn = func(b byte, s string) string {
	// 		switch {
	// 		case b == 0:
	// 			return ansi.FgBrightBlack + s + ansi.Reset
	// 		case b >= 32 && b <= 126, b == '\r', b == '\n', b == '\f', b == '\t', b == '\v':
	// 			return ansi.FgWhite + s + ansi.Reset
	// 		default:
	// 			return ansi.FgBrightWhite + s + ansi.Reset
	// 		}
	// 	}
	// 	column = ansi.FgWhite + colStr + ansi.Reset + "\n"
	// }

	return Decorator{
		Reset:  ansi.FromString(opts.Color["reset"]),
		Null:   ansi.FromString(opts.Color["null"]),
		False:  ansi.FromString(opts.Color["false"]),
		True:   ansi.FromString(opts.Color["true"]),
		Number: ansi.FromString(opts.Color["number"]),
		String: ansi.FromString(opts.Color["string"]),
		Key:    ansi.FromString(opts.Color["key"]),
		Array:  ansi.FromString(opts.Color["array"]),
		Object: ansi.FromString(opts.Color["object"]),

		Name:  ansi.FromString(opts.Color["name"]),
		Frame: ansi.FromString(opts.Color["frame"]),

		ByteColor: func(b byte) ansi.Color { return ansi.FromString("white") },

		Column: colStr,
	}
}

// resetColor     = newColor("0")    // Reset
// nullColor      = newColor("90")   // Bright black
// falseColor     = newColor("33")   // Yellow
// trueColor      = newColor("33")   // Yellow
// numberColor    = newColor("36")   // Cyan
// stringColor    = newColor("32")   // Green
// objectKeyColor = newColor("34;1") // Bold Blue
// arrayColor     = []byte(nil)      // No color
// objectColor    = []byte(nil)      // No color

type Decorator struct {
	Reset  ansi.Color // Reset
	Null   ansi.Color // Bright black
	False  ansi.Color // Yellow
	True   ansi.Color // Yellow
	Number ansi.Color // Cyan
	String ansi.Color // Green
	Key    ansi.Color // Bold Blue
	Array  ansi.Color // White
	Object ansi.Color // White

	Name  ansi.Color
	Frame ansi.Color

	ByteColor func(b byte) ansi.Color

	Column string
}
