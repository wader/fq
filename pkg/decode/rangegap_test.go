package decode_test

import (
	"fq/pkg/decode"
	"log"
	"testing"
)

func Test(t *testing.T) {

	l := decode.RangeGaps(10, []decode.Range{{Start: 1, Len: 1}, {Start: 1, Len: 4}, {Start: 7, Len: 1}})
	log.Printf("l: %#+v\n", l)

}
