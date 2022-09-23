package postgres

import (
	"embed"
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/postgres/flavours/postgres14"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

//go:embed pg_btree.md
var pgBTreeFS embed.FS

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.PG_BTREE,
		Description: "PostgreSQL btree index file",
		DecodeFn:    decodePgBTree,
		DecodeInArg: format.PostgresIn{
			Flavour: PG_FLAVOUR_POSTGRES14,
		},
		RootArray: true,
		RootName:  "pages",
	})
	interp.RegisterFS(pgBTreeFS)
}

func decodePgBTree(d *decode.D, in any) any {
	d.Endian = decode.LittleEndian

	pgIn, ok := in.(format.PostgresIn)
	if !ok {
		d.Fatalf("DecodeInArg must be PostgresIn!\n")
	}

	switch pgIn.Flavour {
	case PG_FLAVOUR_POSTGRES14:
		return postgres14.DecodePgBTree(d)
	}

	return nil
}
