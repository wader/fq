package stringsex

import "github.com/wader/fq/internal/mathex"

func TrimN(s string, n int, suffix string) string {
	sl, tl := len(s), len(suffix)
	if sl+tl <= n {
		return s
	}
	return s[0:mathex.Max(n-tl, 0)] + suffix
}
