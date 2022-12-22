package postgres

import (
	"embed"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/common"
	"github.com/wader/fq/format/postgres/flavours/pgpro10"
	"github.com/wader/fq/format/postgres/flavours/pgpro11"
	"github.com/wader/fq/format/postgres/flavours/pgpro12"
	"github.com/wader/fq/format/postgres/flavours/pgpro13"
	"github.com/wader/fq/format/postgres/flavours/pgpro14"
	"github.com/wader/fq/format/postgres/flavours/pgproee10"
	"github.com/wader/fq/format/postgres/flavours/pgproee11"
	"github.com/wader/fq/format/postgres/flavours/pgproee12"
	"github.com/wader/fq/format/postgres/flavours/pgproee13"
	"github.com/wader/fq/format/postgres/flavours/pgproee14"
	"github.com/wader/fq/format/postgres/flavours/postgres10"
	"github.com/wader/fq/format/postgres/flavours/postgres11"
	"github.com/wader/fq/format/postgres/flavours/postgres12"
	"github.com/wader/fq/format/postgres/flavours/postgres13"
	"github.com/wader/fq/format/postgres/flavours/postgres14"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

//go:embed pg_control.md
var pgControlFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_CONTROL,
		Description: "PostgreSQL control file",
		DecodeFn:    decodePgControl,
		DecodeInArg: format.PostgresIn{
			Flavour: "",
		},
	})
	interp.RegisterFS(pgControlFS)
}

//nolint:revive
const (
	PG_CONTROL_VERSION_10    = 1002
	PG_CONTROL_VERSION_11    = 1100
	PGPRO_CONTROL_VERSION_11 = 1200
	PG_CONTROL_VERSION_12    = 1201
	//PG_CONTROL_VERSION_13 = 1300
	PG_CONTROL_VERSION_14 = 1300
)

//nolint:revive
const (
	PG_FLAVOUR_POSTGRES10 = "postgres10"
	PG_FLAVOUR_POSTGRES11 = "postgres11"
	PG_FLAVOUR_POSTGRES12 = "postgres12"
	PG_FLAVOUR_POSTGRES13 = "postgres13"
	PG_FLAVOUR_POSTGRES14 = "postgres14"
	PG_FLAVOUR_PGPRO10    = "pgpro10"
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

	pgIn, ok := in.(format.PostgresIn)
	if !ok {
		d.Fatalf("DecodeInArg must be PostgresIn!\n")
	}

	switch pgIn.Flavour {
	case PG_FLAVOUR_POSTGRES10:
		return postgres10.DecodePgControl(d, in)
	case PG_FLAVOUR_POSTGRES11:
		return postgres11.DecodePgControl(d, in)
	case PG_FLAVOUR_POSTGRES12:
		return postgres12.DecodePgControl(d, in)
	case PG_FLAVOUR_POSTGRES13:
		return postgres13.DecodePgControl(d, in)
	case PG_FLAVOUR_POSTGRES14:
		return postgres14.DecodePgControl(d, in)
	case PG_FLAVOUR_PGPRO10:
		return pgpro10.DecodePgControl(d, in)
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
	d.SeekAbs(0)

	pgProVersion, oriVersion := common.ParsePgProVersion(uint32(pgControlVersion))

	if pgProVersion == common.PG_ORIGINAL {
		switch oriVersion {
		case PG_CONTROL_VERSION_10:
			return postgres10.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_11:
			return postgres11.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_12:
			return postgres12.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_14:
			return postgres14.DecodePgControl(d, in)
		}
	}

	if pgProVersion == common.PGPRO_STANDARD {
		switch oriVersion {
		case PG_CONTROL_VERSION_10:
			return pgpro10.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_11:
			return pgpro11.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_12:
			return pgpro12.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_14:
			return pgpro14.DecodePgControl(d, in)
		}
	}

	if pgProVersion == common.PGPRO_ENTERPRISE {
		switch oriVersion {
		case PG_CONTROL_VERSION_10:
			return pgproee10.DecodePgControl(d, in)
		case PGPRO_CONTROL_VERSION_11:
			return pgproee11.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_12:
			return pgproee12.DecodePgControl(d, in)
		case PG_CONTROL_VERSION_14:
			return pgproee14.DecodePgControl(d, in)
		}
	}

	d.Fatalf("unsupported PG_CONTROL_VERSION = %d\n", pgControlVersion)
	return nil
}