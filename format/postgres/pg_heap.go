package postgres

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/flavours/pgpro11"
	"github.com/wader/fq/format/postgres/flavours/pgpro12"
	"github.com/wader/fq/format/postgres/flavours/pgpro13"
	"github.com/wader/fq/format/postgres/flavours/pgpro14"
	"github.com/wader/fq/format/postgres/flavours/pgproee10"
	"github.com/wader/fq/format/postgres/flavours/pgproee11"
	"github.com/wader/fq/format/postgres/flavours/pgproee12"
	"github.com/wader/fq/format/postgres/flavours/pgproee13"
	"github.com/wader/fq/format/postgres/flavours/pgproee14"
	"github.com/wader/fq/format/postgres/flavours/postgres11"
	"github.com/wader/fq/format/postgres/flavours/postgres12"
	"github.com/wader/fq/format/postgres/flavours/postgres13"
	"github.com/wader/fq/format/postgres/flavours/postgres14"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

// TO DO
// oom kill on 1 GB file

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_HEAP,
		Description: "PostgreSQL heap file",
		DecodeFn:    decodePgheap,
		DecodeInArg: format.PostgresIn{
			Flavour: "default",
		},
		RootArray: true,
		RootName:  "Pages",
	})
}

func decodePgheap(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	flavour := in.(format.PostgresIn).Flavour
	switch flavour {
	case PG_FLAVOUR_POSTGRES11:
		return postgres11.DecodeHeap(d)
	case PG_FLAVOUR_POSTGRES12:
		return postgres12.DecodeHeap(d)
	case PG_FLAVOUR_POSTGRES13:
		return postgres13.DecodeHeap(d)
	case PG_FLAVOUR_POSTGRES14, PG_FLAVOUR_POSTGRES:
		return postgres14.DecodeHeap(d)
	case PG_FLAVOUR_PGPROEE10:
		return pgproee10.DecodeHeap(d)
	case PG_FLAVOUR_PGPROEE11:
		return pgproee11.DecodeHeap(d)
	case PG_FLAVOUR_PGPROEE12:
		return pgproee12.DecodeHeap(d)
	case PG_FLAVOUR_PGPROEE13:
		return pgproee13.DecodeHeap(d)
	case PG_FLAVOUR_PGPROEE14:
		return pgproee14.DecodeHeap(d)
	case PG_FLAVOUR_PGPRO11:
		return pgpro11.DecodeHeap(d)
	case PG_FLAVOUR_PGPRO12:
		return pgpro12.DecodeHeap(d)
	case PG_FLAVOUR_PGPRO13:
		return pgpro13.DecodeHeap(d)
	case PG_FLAVOUR_PGPRO14:
		return pgpro14.DecodeHeap(d)

	default:
		break
	}

	return postgres14.DecodeHeap(d)
}
