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
	"github.com/wader/fq/format/postgres/flavours/pgproee15"
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
	interp.RegisterFormat(format.Pg_Control, &decode.Format{
		Description: "PostgreSQL control file",
		DecodeFn:    decodePgControl,
		DefaultInArg: format.Pg_Control_In{
			Flavour: "",
		},
	})
	interp.RegisterFS(pgControlFS)
}

const (
	PG_CONTROL_VERSION_10    = 1002
	PG_CONTROL_VERSION_11    = 1100
	PGPRO_CONTROL_VERSION_11 = 1200
	PG_CONTROL_VERSION_12    = 1201
	//PG_CONTROL_VERSION_13 = 1300
	PG_CONTROL_VERSION_14 = 1300
	//PG_CONTROL_VERSION_15 = 1300
)

const (
	PG_FLAVOUR_POSTGRES10 = "postgres10"
	PG_FLAVOUR_POSTGRES11 = "postgres11"
	PG_FLAVOUR_POSTGRES12 = "postgres12"
	PG_FLAVOUR_POSTGRES13 = "postgres13"
	PG_FLAVOUR_POSTGRES14 = "postgres14"
	PG_FLAVOUR_POSTGRES15 = "postgres15"
	PG_FLAVOUR_PGPRO10    = "pgpro10"
	PG_FLAVOUR_PGPRO11    = "pgpro11"
	PG_FLAVOUR_PGPRO12    = "pgpro12"
	PG_FLAVOUR_PGPRO13    = "pgpro13"
	PG_FLAVOUR_PGPRO14    = "pgpro14"
	PG_FLAVOUR_PGPRO15    = "pgpro15"
	PG_FLAVOUR_PGPROEE10  = "pgproee10"
	PG_FLAVOUR_PGPROEE11  = "pgproee11"
	PG_FLAVOUR_PGPROEE12  = "pgproee12"
	PG_FLAVOUR_PGPROEE13  = "pgproee13"
	PG_FLAVOUR_PGPROEE14  = "pgproee14"
	PG_FLAVOUR_PGPROEE15  = "pgproee15"
)

func decodePgControl(d *decode.D) any {
	d.Endian = decode.LittleEndian

	var pgIn format.Pg_Control_In
	if !d.ArgAs(&pgIn) {
		d.Fatalf("no flavour specified")
	}

	switch pgIn.Flavour {
	case PG_FLAVOUR_POSTGRES10:
		return postgres10.DecodePgControl(d)
	case PG_FLAVOUR_POSTGRES11:
		return postgres11.DecodePgControl(d)
	case PG_FLAVOUR_POSTGRES12:
		return postgres12.DecodePgControl(d)
	case PG_FLAVOUR_POSTGRES13:
		return postgres13.DecodePgControl(d)
	case PG_FLAVOUR_POSTGRES14, PG_FLAVOUR_POSTGRES15, PG_FLAVOUR_PGPRO15:
		return postgres14.DecodePgControl(d)
	case PG_FLAVOUR_PGPRO10:
		return pgpro10.DecodePgControl(d)
	case PG_FLAVOUR_PGPRO11:
		return pgpro11.DecodePgControl(d)
	case PG_FLAVOUR_PGPRO12:
		return pgpro12.DecodePgControl(d)
	case PG_FLAVOUR_PGPRO13:
		return pgpro13.DecodePgControl(d)
	case PG_FLAVOUR_PGPRO14:
		return pgpro14.DecodePgControl(d)
	case PG_FLAVOUR_PGPROEE10:
		return pgproee10.DecodePgControl(d)
	case PG_FLAVOUR_PGPROEE11:
		return pgproee11.DecodePgControl(d)
	case PG_FLAVOUR_PGPROEE12:
		return pgproee12.DecodePgControl(d)
	case PG_FLAVOUR_PGPROEE13:
		return pgproee13.DecodePgControl(d)
	case PG_FLAVOUR_PGPROEE14:
		return pgproee14.DecodePgControl(d)
	case PG_FLAVOUR_PGPROEE15:
		return pgproee15.DecodePgControl(d)
	default:
		break
	}

	return probeForDecode(d)
}

func probeForDecode(d *decode.D) any {
	/*    0      |     8 */ // uint64 system_identifier;
	/*    8      |     4 */ // uint32 pg_control_version;
	d.U64()
	pgControlVersion := d.U32()
	d.SeekAbs(0)

	pgProVersion, oriVersion := common.ParsePgProVersion(uint32(pgControlVersion))

	if pgProVersion == common.PG_ORIGINAL {
		switch oriVersion {
		case PG_CONTROL_VERSION_10:
			return postgres10.DecodePgControl(d)
		case PG_CONTROL_VERSION_11:
			return postgres11.DecodePgControl(d)
		case PG_CONTROL_VERSION_12:
			return postgres12.DecodePgControl(d)
		case PG_CONTROL_VERSION_14:
			return postgres14.DecodePgControl(d)
		}
	}

	if pgProVersion == common.PGPRO_STANDARD {
		switch oriVersion {
		case PG_CONTROL_VERSION_10:
			return pgpro10.DecodePgControl(d)
		case PG_CONTROL_VERSION_11:
			return pgpro11.DecodePgControl(d)
		case PG_CONTROL_VERSION_12:
			return pgpro12.DecodePgControl(d)
		case PG_CONTROL_VERSION_14:
			return pgpro14.DecodePgControl(d)
		}
	}

	if pgProVersion == common.PGPRO_ENTERPRISE {
		switch oriVersion {
		case PG_CONTROL_VERSION_10:
			return pgproee10.DecodePgControl(d)
		case PGPRO_CONTROL_VERSION_11:
			return pgproee11.DecodePgControl(d)
		case PG_CONTROL_VERSION_12:
			return pgproee12.DecodePgControl(d)
		case PG_CONTROL_VERSION_14:
			return pgproee14.DecodePgControl(d)
		}
	}

	d.Fatalf("unsupported PG_CONTROL_VERSION = %d", pgControlVersion)
	return nil
}
