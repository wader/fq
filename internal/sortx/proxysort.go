package sortx

import "sort"

type proxySort[Tae any, Ta []Tae, Tbe any, Tb []Tbe] struct {
	a  Ta
	b  Tb
	fn func(a, b Tae) bool
}

func (p proxySort[Tae, Ta, Tbe, Tb]) Len() int           { return len(p.a) }
func (p proxySort[Tae, Ta, Tbe, Tb]) Less(i, j int) bool { return p.fn(p.a[i], p.a[j]) }
func (p proxySort[Tae, Ta, Tbe, Tb]) Swap(i, j int) {
	p.a[i], p.a[j] = p.a[j], p.a[i]
	p.b[i], p.b[j] = p.b[j], p.b[i]
}

// ProxySort same as sort.Sort but is type safe, swaps an additional slice b and also does indexing
// Assumes both slices have same length.
func ProxySort[Tae any, Ta []Tae, Tbe any, Tb []Tbe](a Ta, b Tb, fn func(a, b Tae) bool) {
	sort.Sort(proxySort[Tae, Ta, Tbe, Tb]{a: a, b: b, fn: fn})
}

// ProxyStable same as sort.Proxy but is type safe, swaps an additional slice b and also does indexing
// Assumes both slices have same length.
func ProxyStable[Tae any, Ta []Tae, Tbe any, Tb []Tbe](a Ta, b Tb, fn func(a, b Tae) bool) {
	sort.Stable(proxySort[Tae, Ta, Tbe, Tb]{a: a, b: b, fn: fn})
}
