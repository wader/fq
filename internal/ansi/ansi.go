package ansi

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
