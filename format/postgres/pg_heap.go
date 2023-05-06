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
	interp.RegisterFormat(format.Pg_Heap, &decode.Format{
		Description: "PostgreSQL heap file",
		DecodeFn:    decodePgheap,
		DefaultInArg: format.Pg_Heap_In{
			Flavour: PG_FLAVOUR_POSTGRES14,
			Page:    0,
			Segment: 0,
		},
		RootArray: true,
		RootName:  "pages",
	})
	interp.RegisterFS(pgHeapFS)
}

func decodePgheap(d *decode.D) any {
	d.Endian = decode.LittleEndian

	var pgIn format.Pg_Heap_In
	if !d.ArgAs(&pgIn) {
		d.Fatalf("no flavour specified")
	}

	switch pgIn.Flavour {
	case PG_FLAVOUR_POSTGRES10,
		PG_FLAVOUR_POSTGRES11,
		PG_FLAVOUR_POSTGRES12,
		PG_FLAVOUR_POSTGRES13,
		PG_FLAVOUR_POSTGRES14,
		PG_FLAVOUR_POSTGRES15,
		PG_FLAVOUR_PGPRO10,
		PG_FLAVOUR_PGPRO11,
		PG_FLAVOUR_PGPRO12,
		PG_FLAVOUR_PGPRO13,
		PG_FLAVOUR_PGPRO14,
		PG_FLAVOUR_PGPRO15:
		return postgres.DecodeHeap(d, pgIn)

	case PG_FLAVOUR_PGPROEE10,
		PG_FLAVOUR_PGPROEE11,
		PG_FLAVOUR_PGPROEE12,
		PG_FLAVOUR_PGPROEE13,
		PG_FLAVOUR_PGPROEE14,
		PG_FLAVOUR_PGPROEE15:
		return pgproee.DecodeHeap(d, pgIn)

	default:
		break
	}

	return postgres.DecodeHeap(d, pgIn)
}
