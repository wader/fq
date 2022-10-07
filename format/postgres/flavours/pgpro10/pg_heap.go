package pgpro10

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/flavours/postgres14/common14"
	"github.com/wader/fq/pkg/decode"
)

func DecodeHeap(d *decode.D, args format.PostgresHeapIn) any {
	heap := &common14.Heap{
		Args:                 args,
		DecodePageHeaderData: common14.DecodePageHeader,
	}
	return common14.DecodeHeap(heap, d)
}
