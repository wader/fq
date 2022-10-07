package pgpro10

import (
	"github.com/wader/fq/format"
	postgres2 "github.com/wader/fq/format/postgres/common/pg_heap/postgres"
	"github.com/wader/fq/pkg/decode"
)

func DecodeHeap(d *decode.D, args format.PostgresHeapIn) any {
	heap := &postgres2.Heap{
		Args:                 args,
		DecodePageHeaderData: postgres2.DecodePageHeader,
	}
	return postgres2.DecodeHeap(heap, d)
}
