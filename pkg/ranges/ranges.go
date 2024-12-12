package ranges

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Range struct {
	Start int64
	Len   int64
}

func (r Range) Stop() int64 { return r.Start + r.Len }

func (r Range) String() string { return fmt.Sprintf("%d:%d", r.Start, r.Len) }

func (r Range) IsZero() bool { return r.Start == 0 && r.Len == 0 }

func RangeFromString(s string) Range {
	ps := strings.Split(s, ":")
	start, _ := strconv.Atoi(ps[0])
	l, _ := strconv.Atoi(ps[1])
	return Range{Start: int64(start), Len: int64(l)}
}

func SliceFromString(s string) []Range {
	rr := []Range{}
	for _, p := range strings.Split(s, " ") {
		if p == "" {
			continue
		}
		rr = append(rr, RangeFromString(p))
	}
	return rr
}

func MinMax(a, b Range) Range {
	minStart := min(a.Start, b.Start)
	maxStop := max(a.Stop(), b.Stop())
	return Range{Start: minStart, Len: maxStop - minStart}
}

// Gaps in ranges limited by total range
func Gaps(total Range, ranges []Range) []Range {
	if len(ranges) == 0 {
		return []Range{total}
	}

	slices.SortFunc(ranges, func(a, b Range) int { return cmp.Compare(a.Start, b.Start) })

	// worst case ranges+1 gaps
	merged := make([]Range, 0, len(ranges)+1)
	var madded bool
	var m Range

	for i := 0; i < len(ranges); {
		madded = false
		m = ranges[i]

		// skip empty ranges
		if m.Len == 0 {
			i++
			madded = true
			continue
		}

		j := i + 1
		for ; j < len(ranges); j++ {
			if m.Start <= ranges[j].Start && m.Stop()+1 >= ranges[j].Start {
				if ranges[j].Stop() > m.Stop() {
					m.Len = ranges[j].Stop() - m.Start
				}
			} else {
				i = j
				merged = append(merged, m)
				madded = true
				break
			}
		}

		if j >= len(ranges) {
			break
		}
	}

	if !madded {
		merged = append(merged, m)
	}

	if len(merged) == 0 {
		return []Range{total}
	}

	gaps := make([]Range, 0, len(merged))
	if merged[0].Start != total.Start {
		gaps = append(gaps, Range{Start: 0, Len: merged[0].Start})
	}
	for i := 0; i < len(merged)-1; i++ {
		gaps = append(gaps, Range{Start: merged[i].Stop(), Len: merged[i+1].Start - merged[i].Stop()})
	}
	l := merged[len(merged)-1]
	if l.Stop() != total.Stop() {
		gaps = append(gaps, Range{Start: l.Stop(), Len: total.Stop() - l.Stop()})
	}

	return gaps
}
