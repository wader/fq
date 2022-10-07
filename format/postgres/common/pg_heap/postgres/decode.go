package postgres

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
)

func DecodeHeap(d *decode.D, args format.PostgresHeapIn) any {
	heap := &Heap{
		Args:                 args,
		DecodePageHeaderData: DecodePageHeader,
	}
	return Decode(heap, d)
}
