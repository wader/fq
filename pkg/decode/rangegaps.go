package decode

import (
	"sort"
)

func RangeGaps(total Range, ranges []Range) []Range {
	if len(ranges) == 0 {
		return []Range{total}
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})

	var merged []Range
	var madded bool
	var m Range

	for i := 0; i < len(ranges); {
		madded = false
		m = ranges[i]
		j := i + 1
		for ; j < len(ranges); j++ {
			if m.Start <= ranges[j].Start && m.Start+m.Len+1 >= ranges[j].Start {
				if ranges[j].Start+ranges[j].Len > m.Start+m.Len {
					m.Len = ranges[j].Start + ranges[j].Len - m.Start
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

	gaps := []Range{}
	if merged[0].Start != total.Start {
		gaps = append(gaps, Range{Start: 0, Len: merged[0].Start})
	}
	for i := 0; i < len(merged)-1; i++ {
		m := merged[i]
		n := merged[i+1]
		gaps = append(gaps, Range{Start: m.Start + m.Len, Len: n.Start - (m.Start + m.Len)})
	}
	l := merged[len(merged)-1]
	if l.Start+l.Len != total.Start+total.Len {
		gaps = append(gaps, Range{Start: l.Start + l.Len, Len: total.Start + total.Len - (l.Start + l.Len)})
	}

	return gaps
}
