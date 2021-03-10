package ansi

import (
	"fmt"
	"io"
	"strings"
)

const FgBlack = "\x1b[30m"
const FgRed = "\x1b[31m"
const FgGreen = "\x1b[32m"
const FgYellow = "\x1b[33m"
const FgBlue = "\x1b[34m"
const FgMagenta = "\x1b[35m"
const FgCyan = "\x1b[36m"
const FgWhite = "\x1b[37m"
const FgBrightBlack = "\x1b[90m"
const FgBrightRed = "\x1b[91m"
const FgBrightGreen = "\x1b[92m"
const FgBrightYellow = "\x1b[93m"
const FgBrightBlue = "\x1b[94m"
const FgBrightMagenta = "\x1b[95m"
const FgBrightCyan = "\x1b[96m"
const FgBrightWhite = "\x1b[97m"
const BgBlack = "\x1b[40m"
const BgRed = "\x1b[41m"
const BgGreen = "\x1b[42m"
const BgYellow = "\x1b[43m"
const BgBlue = "\x1b[44m"
const BgMagenta = "\x1b[45m"
const BgCyan = "\x1b[46m"
const BgWhite = "\x1b[47m"
const BgBrightBlack = "\x1b[100m"
const BgBrightRed = "\x1b[101m"
const BgBrightGreen = "\x1b[102m"
const BgBrightYellow = "\x1b[103m"
const BgBrightBlue = "\x1b[104m"
const BgBrightMagenta = "\x1b[105m"
const BgBrightCyan = "\x1b[106m"
const BgBrightWhite = "\x1b[107m"
const Reset = "\x1b[0m"

var Foreground = map[string]string{
	"black":         FgBlack,
	"red":           FgRed,
	"green":         FgGreen,
	"yellow":        FgYellow,
	"blue":          FgBlue,
	"magenta":       FgMagenta,
	"cyan":          FgCyan,
	"white":         FgWhite,
	"brightblack":   FgBrightBlack,
	"brightred":     FgBrightRed,
	"brightgreen":   FgBrightGreen,
	"brightyellow":  FgBrightYellow,
	"brightblue":    FgBrightBlue,
	"brightmagenta": FgBrightMagenta,
	"brightcyan":    FgBrightCyan,
	"brightwhite":   FgBrightWhite,
}

var Background = map[string]string{
	"black":         BgBlack,
	"red":           BgRed,
	"green":         BgGreen,
	"yellow":        BgYellow,
	"blue":          BgBlue,
	"magenta":       BgMagenta,
	"cyan":          BgCyan,
	"white":         BgWhite,
	"brightblack":   BgBrightBlack,
	"brightred":     BgBrightRed,
	"brightgreen":   BgBrightGreen,
	"brightyellow":  BgBrightYellow,
	"brightblue":    BgBrightBlue,
	"brightmagenta": BgBrightMagenta,
	"brightcyan":    BgBrightCyan,
	"brightwhite":   BgBrightWhite,
}

type Color string

func FromString(s string) Color {
	parts := strings.SplitN(s, ":", 2)
	fg := ""
	bg := ""
	if len(parts) > 0 {
		fg = Foreground[parts[0]]
	}
	if len(parts) > 1 {
		bg = Foreground[parts[1]]
	}
	return Color(fg + bg)
}

func (c Color) Write(w io.Writer, p []byte) (int, error) {
	if c != "" {
		if n, err := w.Write([]byte(c)); err != nil {
			return n, err
		}
		if n, err := w.Write(p); err != nil {
			return n, err
		}
		if n, err := w.Write([]byte(Reset)); err != nil {
			return n, err
		}
		return len(c) + len(p) + len(Reset), nil
	}
	return w.Write(p)
}

func (c Color) Wrap(s string) string {
	if c != "" {
		return string(c) + s + Reset
	}
	return s
}

type colorFormatter [3]string

func (cf colorFormatter) Format(state fmt.State, verb rune) {
	switch verb {
	case 's':
		switch len(cf) {
		case 1:
			fmt.Fprint(state, cf[0])
		case 3:
			fmt.Fprint(state, cf[0], cf[1], cf[2])
		default:
			panic("unreachable")
		}
	}
}

func (c Color) F(s string) fmt.Formatter {
	if c != "" {
		return colorFormatter([3]string{string(c), s, Reset})
	}
	return colorFormatter([3]string{s})
}

type colorWriter struct {
	w io.Writer
	c Color
}

func (cw colorWriter) Write(p []byte) (int, error) {
	return cw.c.Write(cw.w, p)
}

func (c Color) W(w io.Writer) io.Writer {
	if c != "" {
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
