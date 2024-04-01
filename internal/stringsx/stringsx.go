package stringsx

import "github.com/wader/fq/internal/mathx"

func TrimN(s string, n int, suffix string) string {
	sl, tl := len(s), len(suffix)
	if sl+tl <= n {
		return s
	}
	return s[0:mathx.Max(n-tl, 0)] + suffix
}
