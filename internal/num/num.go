package num

import (
	"math"
	"strconv"
	"strings"
)

func DigitsInBase(n int64, base int) int {
	if n == 0 {
		return 1
	}
	return int(1 + math.Floor(math.Log(float64(n))/math.Log(float64(base))))
}

func PadFormatInt(i int64, base int, width int) string {
	s := strconv.FormatInt(i, base)
	p := width - len(s)
	if p > 0 {
		// TODO: something faster?
		return strings.Repeat("0", p) + s
	}
	return s
}
