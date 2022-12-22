package postgres

import (
	"embed"

	"github.com/wader/fq/format/postgres/common/pg_btree/postgres"

	"github.com/wader/fq/format"
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
		DecodeInArg: format.PostgresBTreeIn{
			Page: 0,
		},
		RootArray: true,
		RootName:  "pages",
	})
	interp.RegisterFS(pgBTreeFS)
}

func decodePgBTree(d *decode.D, in any) any {
	pgIn, ok := in.(format.PostgresBTreeIn)
	if !ok {
		d.Fatalf("DecodeInArg must be PostgresBTreeIn!\n")
	}
	return postgres.DecodePgBTree(d, pgIn)
}
