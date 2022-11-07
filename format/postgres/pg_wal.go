package postgres

import (
	"fmt"

	"github.com/wader/fq/format/postgres/common/pg_wal/pgproee"
	"github.com/wader/fq/format/postgres/common/pg_wal/postgres"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"

	"strconv"
	"strings"
)

// partial parsing of WAL

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_WAL,
		Description: "PostgreSQL write-ahead log file",
		DecodeFn:    decodePGWAL,
		DecodeInArg: format.PostgresIn{
			Flavour: PG_FLAVOUR_POSTGRES14,
		},
		RootArray: true,
		RootName:  "pages",
	})
}

func ParseLsn(lsn string) (uint32, error) {
	// check for 0/4E394440
	str1 := lsn
	if strings.Contains(lsn, "/") {
		parts := strings.Split(lsn, "/")
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid lsn = %s", lsn)
		}
		str1 = parts[1]
	}
	// parse hex to coded file name + file offset
	r1, err := strconv.ParseInt(str1, 16, 64)
	if err != nil {
		return 0, err
	}
	return uint32(r1), err
}

func decodePGWAL(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	pgIn, ok := in.(format.PostgresIn)
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
		return postgres.DecodePGWAL(d)

	case PG_FLAVOUR_PGPROEE10,
		PG_FLAVOUR_PGPROEE11,
		PG_FLAVOUR_PGPROEE12,
		PG_FLAVOUR_PGPROEE13,
		PG_FLAVOUR_PGPROEE14:
		return pgproee.DecodePGWAL(d)

	default:
		break
	}

	return postgres.DecodePGWAL(d)
}
