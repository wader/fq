package rangegap

import "sort"

func Find(start int64, stop int64, ranges [][2]int64) [][2]int64 {
	if len(ranges) == 0 {
		return [][2]int64{{start, stop}}
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})

	var merged [][2]int64
	var nadded bool
	var nstart int64
	var nstop int64

	for i := 0; i < len(ranges); {
		nadded = false
		nstart = ranges[i][0]
		nstop = ranges[i][1]
		j := i + 1
		for ; j < len(ranges); j++ {
			cstart := ranges[j][0]
			if nstart <= cstart && nstop+1 >= cstart {
				if ranges[j][1] > nstop {
					nstop = ranges[j][1]
				}
			} else {
				i = j
				merged = append(merged, [2]int64{nstart, nstop})
				break
			}
		}

		if j >= len(ranges) {
			break
		}

	}

	if !nadded {
		merged = append(merged, [2]int64{nstart, nstop})
	}

	var gaps [][2]int64
	if start < merged[0][0] {
		gaps = append(gaps, [2]int64{start, merged[0][0] - 1})
	}
	for i := 0; i < len(merged)-1; i++ {
		gaps = append(gaps, [2]int64{merged[i][1] + 1, merged[i+1][0] - 1})
	}
	if stop > merged[len(merged)-1][1] {
		gaps = append(gaps, [2]int64{merged[len(merged)-1][1] + 1, stop})
	}

	return gaps
}
