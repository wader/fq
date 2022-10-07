package pgproee

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/common/pg_heap/postgres"
	"github.com/wader/fq/pkg/decode"
)

func DecodeHeap(d *decode.D, args format.PostgresHeapIn) any {
	heap := &postgres.Heap{
		Args:                 args,
		DecodePageHeaderData: DecodePageHeaderData,
		DecodePageSpecial:    DecodePageSpecial,
	}
	return postgres.Decode(heap, d)
}
