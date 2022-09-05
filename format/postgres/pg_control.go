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

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_CONTROL,
		Description: "PostgreSQL control file",
		DecodeFn:    decodePgControl,
		DecodeInArg: format.PostgresIn{
			Flavour: "default",
		},
	})
}

const (
	PG_CONTROL_VERSION_11 = 1100
	PG_CONTROL_VERSION_12 = 1201
	//PG_CONTROL_VERSION_13 = 1300
	PG_CONTROL_VERSION_14 = 1300
)

const (
	PG_FLAVOUR_POSTGRES   = "postgres"
	PG_FLAVOUR_POSTGRES11 = "postgres11"
	PG_FLAVOUR_POSTGRES12 = "postgres12"
	PG_FLAVOUR_POSTGRES13 = "postgres13"
	PG_FLAVOUR_POSTGRES14 = "postgres14"
	PG_FLAVOUR_PGPRO11    = "pgpro11"
	PG_FLAVOUR_PGPRO12    = "pgpro12"
	PG_FLAVOUR_PGPRO13    = "pgpro13"
	PG_FLAVOUR_PGPRO14    = "pgpro14"
	PG_FLAVOUR_PGPROEE10  = "pgproee10"
	PG_FLAVOUR_PGPROEE11  = "pgproee11"
	PG_FLAVOUR_PGPROEE12  = "pgproee12"
	PG_FLAVOUR_PGPROEE13  = "pgproee13"
	PG_FLAVOUR_PGPROEE14  = "pgproee14"
)

func decodePgControl(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	flavour := in.(format.PostgresIn).Flavour
	switch flavour {
	case PG_FLAVOUR_POSTGRES11:
		return postgres11.DecodePgControl(d, in)
	case PG_FLAVOUR_POSTGRES12:
		return postgres12.DecodePgControl(d, in)
	case PG_FLAVOUR_POSTGRES13:
		return postgres13.DecodePgControl(d, in)
	case PG_FLAVOUR_POSTGRES14, PG_FLAVOUR_POSTGRES:
		return postgres14.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPRO11:
		return pgpro11.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPRO12:
		return pgpro12.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPRO13:
		return pgpro13.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPRO14:
		return pgpro14.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPROEE10:
		return pgproee10.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPROEE11:
		return pgproee11.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPROEE12:
		return pgproee12.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPROEE13:
		return pgproee13.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPROEE14:
		return pgproee14.DecodePgControl(d, in)
	default:
		break
	}

	return probeForDecode(d, in)
}

func probeForDecode(d *decode.D, in any) any {
	/*    0      |     8 */ // uint64 system_identifier;
	/*    8      |     4 */ // uint32 pg_control_version;
	d.U64()
	pgControlVersion := d.U32()

	// try to guess version
	switch pgControlVersion {
	case PG_CONTROL_VERSION_11:
		return postgres11.DecodePgControl(d, in)
	case PG_CONTROL_VERSION_12:
		return postgres12.DecodePgControl(d, in)
	case PG_CONTROL_VERSION_14:
		return postgres14.DecodePgControl(d, in)
	default:
		d.Fatalf("unsupported PG_CONTROL_VERSION = %d\n", pgControlVersion)
	}
	return nil
}
