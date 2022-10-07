package postgres

import (
	"embed"
	"github.com/wader/fq/format/postgres/common/pg_heap/pgproee"
	"github.com/wader/fq/format/postgres/common/pg_heap/postgres"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

// TO DO
// oom kill on 1 GB file

//go:embed pg_heap.md
var pgHeapFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_HEAP,
		Description: "PostgreSQL heap file",
		DecodeFn:    decodePgheap,
		DecodeInArg: format.PostgresHeapIn{
			Flavour:       PG_FLAVOUR_POSTGRES14,
			PageNumber:    0,
			SegmentNumber: 0,
		},
		RootArray: true,
		RootName:  "pages",
	})
	interp.RegisterFS(pgHeapFS)
}

func decodePgheap(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	pgIn, ok := in.(format.PostgresHeapIn)
	if !ok {
		d.Fatalf("DecodeInArg must be PostgresIn!\n")
	}

	switch pgIn.Flavour {
	case PG_FLAVOUR_POSTGRES10,
		PG_FLAVOUR_POSTGRES11,
		PG_FLAVOUR_POSTGRES12,
		PG_FLAVOUR_POSTGRES13,
		PG_FLAVOUR_POSTGRES14,
		PG_FLAVOUR_PGPRO10,
		PG_FLAVOUR_PGPRO11,
		PG_FLAVOUR_PGPRO12,
		PG_FLAVOUR_PGPRO13,
		PG_FLAVOUR_PGPRO14:
		return postgres.DecodeHeap(d, pgIn)

	case PG_FLAVOUR_PGPROEE10,
		PG_FLAVOUR_PGPROEE11,
		PG_FLAVOUR_PGPROEE12,
		PG_FLAVOUR_PGPROEE13,
		PG_FLAVOUR_PGPROEE14:
		return pgproee.DecodeHeap(d, pgIn)

	default:
		break
	}

	return postgres.DecodeHeap(d, pgIn)
}
