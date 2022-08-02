package postgres

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/flavours/pgproee11"
	"github.com/wader/fq/format/postgres/flavours/pgproee14"
	"github.com/wader/fq/format/postgres/flavours/postgres14"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PGHEAP,
		Description: "PostgreSQL heap file",
		DecodeFn:    decodePgheap,
		DecodeInArg: format.PostgresIn{
			Flavour: "default",
		},
	})
}

func decodePgheap(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	flavour := in.(format.PostgresIn).Flavour
	switch flavour {
	case PG_FLAVOUR_POSTGRES14, PG_FLAVOUR_POSTGRES:
		return postgres14.DecodeHeap(d)
	case PG_FLAVOUR_PGPROEE11:
		return pgproee11.DecodeHeap(d)
	case PG_FLAVOUR_PGPROEE14:
		return pgproee14.DecodeHeap(d)
	default:
		break
	}

	return postgres14.DecodeHeap(d)
}
