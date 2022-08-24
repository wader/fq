package ansi

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Code struct {
	Set         []int
	Reset       []int
	SetString   string
	ResetString string
}

func MakeCode(set []int, reset []int) Code {
	setSB := &strings.Builder{}
	setSB.WriteString("\x1b[")
	for i, s := range set {
		setSB.WriteString(strconv.Itoa(s))
		if i != len(set)-1 {
			setSB.WriteString(";")
		}
	}
	setSB.WriteString("m")
	resetSB := &strings.Builder{}
	resetSB.WriteString("\x1b[")
	for i, s := range reset {
		resetSB.WriteString(strconv.Itoa(s))
		if i != len(reset)-1 {
			resetSB.WriteString(";")
		}
	}
	resetSB.WriteString("m")
	return Code{
		Set:         set,
		Reset:       reset,
		SetString:   setSB.String(),
		ResetString: resetSB.String(),
	}
}

func (c Code) Add(a Code) Code {
	return MakeCode(
		append(c.Set, a.Set...),
		append(c.Reset, a.Reset...),
	)
}

var Reset = MakeCode([]int{}, []int{})

var Black = MakeCode([]int{30}, []int{39})
var Red = MakeCode([]int{31}, []int{39})
var Green = MakeCode([]int{32}, []int{39})
var Yellow = MakeCode([]int{33}, []int{39})
var Blue = MakeCode([]int{34}, []int{39})
var Magenta = MakeCode([]int{35}, []int{39})
var Cyan = MakeCode([]int{36}, []int{39})
var White = MakeCode([]int{37}, []int{39})
var BrightBlack = MakeCode([]int{90}, []int{39})
var BrightRed = MakeCode([]int{91}, []int{39})
var BrightGreen = MakeCode([]int{92}, []int{39})
var BrightYellow = MakeCode([]int{93}, []int{39})
var BrightBlue = MakeCode([]int{94}, []int{39})
var BrightMagenta = MakeCode([]int{95}, []int{39})
var BrightCyan = MakeCode([]int{96}, []int{39})
var BrightWhite = MakeCode([]int{97}, []int{39})
var Bgblack = MakeCode([]int{40}, []int{49})
var Bgred = MakeCode([]int{41}, []int{49})
var Bggreen = MakeCode([]int{42}, []int{49})
var Bgyellow = MakeCode([]int{43}, []int{49})
var Bgblue = MakeCode([]int{44}, []int{49})
var Bgmagenta = MakeCode([]int{45}, []int{49})
var Bgcyan = MakeCode([]int{46}, []int{49})
var Bgwhite = MakeCode([]int{47}, []int{49})
var BgbrightBlack = MakeCode([]int{100}, []int{49})
var BgbrightRed = MakeCode([]int{101}, []int{49})
var BgbrightGreen = MakeCode([]int{102}, []int{49})
var BgbrightYellow = MakeCode([]int{103}, []int{49})
var BgbrightBlue = MakeCode([]int{104}, []int{49})
var BgbrightMagenta = MakeCode([]int{105}, []int{49})
var BgbrightCyan = MakeCode([]int{106}, []int{49})
var BgbrightWhite = MakeCode([]int{107}, []int{49})
var Bold = MakeCode([]int{1}, []int{22})
var Italic = MakeCode([]int{3}, []int{23})
var Underline = MakeCode([]int{4}, []int{24})
var Inverse = MakeCode([]int{7}, []int{27})

var StringToCode = map[string]Code{
	"black":           Black,
	"red":             Red,
	"green":           Green,
	"yellow":          Yellow,
	"blue":            Blue,
	"magenta":         Magenta,
	"cyan":            Cyan,
	"white":           White,
	"brightblack":     BrightBlack,
	"brightred":       BrightRed,
	"brightgreen":     BrightGreen,
	"brightyellow":    BrightYellow,
	"brightblue":      BrightBlue,
	"brightmagenta":   BrightMagenta,
	"brightcyan":      BrightCyan,
	"brightwhite":     BrightWhite,
	"bgblack":         Bgblack,
	"bgred":           Bgred,
	"bggreen":         Bggreen,
	"bgyellow":        Bgyellow,
	"bgblue":          Bgblue,
	"bgmagenta":       Bgmagenta,
	"bgcyan":          Bgcyan,
	"bgwhite":         Bgwhite,
	"bgbrightblack":   BgbrightBlack,
	"bgbrightred":     BgbrightRed,
	"bgbrightgreen":   BgbrightGreen,
	"bgbrightyellow":  BgbrightYellow,
	"bgbrightblue":    BgbrightBlue,
	"bgbrightmagenta": BgbrightMagenta,
	"bgbrightcyan":    BgbrightCyan,
	"bgbrightwhite":   BgbrightWhite,
	"bold":            Bold,
	"italic":          Italic,
	"underline":       Underline,
	"inverse":         Inverse,
}

var None Code

func FromString(s string) Code {
	c := Code{}
	for _, part := range strings.Split(s, "+") {
		pc, ok := StringToCode[part]
		if !ok {
			continue
		}
		c = c.Add(pc)
	}
	return c
}

func (c Code) Write(w io.Writer, p []byte) (int, error) {
	if c.SetString != "" {
		if n, err := w.Write([]byte(c.SetString)); err != nil {
			return n, err
		}
		if n, err := w.Write(p); err != nil {
			return n, err
		}
		if n, err := w.Write([]byte(c.ResetString)); err != nil {
			return n, err
		}
		return len(c.SetString) + len(p) + len(c.ResetString), nil
	}
	return w.Write(p)
}

func (c Code) Wrap(s string) string {
	if c.ResetString != "" {
		return c.SetString + s + c.ResetString
	}
	return s
}

type colorFormatter [3]any

func (cf colorFormatter) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		for _, s := range cf {
			if s == nil {
				continue
			}
			fmt.Fprint(state, s)
		}
	}
}

func (c Code) F(s any) fmt.Formatter {
	if c.SetString != "" {
		return colorFormatter([3]any{c.SetString, s, c.ResetString})
	}
	return colorFormatter([3]any{s})
}

type colorWriter struct {
	w io.Writer
	c Code
}

func (cw colorWriter) Write(p []byte) (int, error) {
	return cw.c.Write(cw.w, p)
}

func (c Code) W(w io.Writer) io.Writer {
	if c.SetString != "" {
		return colorWriter{w: w, c: c}
	}
	return w
}

// Len of string excluding ANSI escape sequences
func Len(s string) int {
	l := 0
	inANSI := false
	for _, c := range s {
		if inANSI {
			if c == 'm' {
				inANSI = false
			}
		} else {
			if c == '\x1b' {
				inANSI = true
			} else {
				l++
			}
		}
	}
	return l
}

// Slice string to start:stop visible characters.
// An ANSI reset is added to the end of the string.
func Slice(s string, start, stop int) string {
	l := 0
	startByte := -1
	inANSI := false
	for i, c := range s {
		if inANSI {
			if c == 'm' {
				inANSI = false
			}
		} else {
			if c == '\x1b' {
				inANSI = true
			} else {
				if startByte == -1 && l == start {
					startByte = i
					if stop == -1 {
						return s[startByte:] + "\x1b[0m"
					}
				} else if l == stop {
					return s[startByte:i] + "\x1b[0m"
				}

				l++
			}
		}
	}
	return s
}
