package sortex

import "sort"

// Slice same as sort.Slice but type safe and also indexes for you
func Slice[T any](s []T, less func(a, b T) bool) {
	sort.Slice(s, func(i, j int) bool { return less(s[i], s[j]) })
}

// SliceStable same as sort.SliceStable but type safe and also indexes for you
func SliceStable[T any](s []T, less func(a, b T) bool) {
	sort.SliceStable(s, func(i, j int) bool { return less(s[i], s[j]) })
}
