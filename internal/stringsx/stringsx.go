package stringsx

func TrimN(s string, n int, suffix string) string {
	sl, tl := len(s), len(suffix)
	if sl+tl <= n {
		return s
	}
	return s[0:max(n-tl, 0)] + suffix
}
