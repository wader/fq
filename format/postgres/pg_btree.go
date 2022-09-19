package postgres

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/flavours/postgres14"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_BTREE,
		Description: "PostgreSQL btree index file",
		DecodeFn:    decodePgBTree,
		DecodeInArg: format.PostgresIn{
			Flavour: "default",
		},
		RootArray: true,
		RootName:  "pages",
	})
}

func decodePgBTree(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	pgIn, ok := in.(format.PostgresIn)
	if !ok {
		d.Fatalf("DecodeInArg must be PostgresIn!\n")
	}

	switch pgIn.Flavour {
	case PG_FLAVOUR_POSTGRES14, PG_FLAVOUR_POSTGRES:
		return postgres14.DecodePgBTree(d)
	}

	return nil
}
