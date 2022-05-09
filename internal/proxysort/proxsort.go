// proxysort sorts one slice and but swaps index in two
// assumes both slices have same length.
package proxysort

import "sort"

type proxySort[Tae any, Ta []Tae, Tbe any, Tb []Tbe] struct {
	a  Ta
	b  Tb
	fn func(Ta []Tae, i, j int) bool
}

func (p proxySort[Tae, Ta, Tbe, Tb]) Len() int           { return len(p.a) }
func (p proxySort[Tae, Ta, Tbe, Tb]) Less(i, j int) bool { return p.fn(p.a, i, j) }
func (p proxySort[Tae, Ta, Tbe, Tb]) Swap(i, j int) {
	p.a[i], p.a[j] = p.a[j], p.a[i]
	p.b[i], p.b[j] = p.b[j], p.b[i]
}

func Sort[Tae any, Ta []Tae, Tbe any, Tb []Tbe](a Ta, b Tb, fn func(Ta []Tae, i, j int) bool) {
	sort.Sort(proxySort[Tae, Ta, Tbe, Tb]{a: a, b: b, fn: fn})
}

func Stable[Tae any, Ta []Tae, Tbe any, Tb []Tbe](a Ta, b Tb, fn func(Ta []Tae, i, j int) bool) {
	sort.Stable(proxySort[Tae, Ta, Tbe, Tb]{a: a, b: b, fn: fn})
}
