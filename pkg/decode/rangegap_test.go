package decode_test

import (
	"fmt"
	"fq/pkg/decode"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestRangeGaps(t *testing.T) {
	r := func(s string) decode.Range {
		ps := strings.Split(s, ",")
		start, _ := strconv.Atoi(ps[0])
		len, _ := strconv.Atoi(ps[1])
		return decode.Range{Start: int64(start), Len: int64(len)}
	}
	rs := func(s string) []decode.Range {
		rr := []decode.Range{}
		for _, p := range strings.Split(s, " ") {
			if p == "" {
				continue
			}
			rr = append(rr, r(p))
		}
		return rr
	}

	testCases := []struct {
		total    decode.Range
		ranges   []decode.Range
		expected []decode.Range
	}{
		{r("0,0"), rs(""), rs("0,0")},
		{r("0,10"), rs(""), rs("0,10")},

		{r("0,10"), rs("0,10"), rs("")},

		{r("0,10"), rs("1,9"), rs("0,1")},
		{r("0,10"), rs("0,9"), rs("9,1")},

		{r("0,10"), rs("1,1 8,1"), rs("0,1 2,6 9,1")},
		{r("0,10"), rs("1,1 2,5 8,1"), rs("0,1 9,1")},
		{r("0,10"), rs("1,1 2,8 8,2"), rs("0,1")},
		{r("0,10"), rs("0,4 2,8 8,2"), rs("")},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("%v_%v_%v", tC.total, tC.ranges, tC.expected), func(t *testing.T) {
			actual := decode.RangeGaps(tC.total, tC.ranges)
			if !reflect.DeepEqual(tC.expected, actual) {
				t.Errorf("expected %v, got %v", tC.expected, actual)
			}
		})
	}
}
