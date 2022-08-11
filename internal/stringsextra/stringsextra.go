package stringsextra

import "github.com/wader/fq/internal/mathextra"

func TrimN(s string, n int, suffix string) string {
	sl, tl := len(s), len(suffix)
	if sl+tl <= n {
		return s
	}
	return s[0:mathextra.MaxInt(n-tl, 0)] + suffix
}
