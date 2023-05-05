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
	interp.RegisterFormat(format.Pg_BTree, &decode.Format{
		Description: "PostgreSQL btree index file",
		DecodeFn:    decodePgBTree,
		DefaultInArg: format.Pg_BTree_In{
			Page: 0,
		},
		RootArray: true,
		RootName:  "pages",
	})
	interp.RegisterFS(pgBTreeFS)
}

func decodePgBTree(d *decode.D) any {
	d.Endian = decode.LittleEndian
	var pgIn format.Pg_BTree_In
	if !d.ArgAs(&pgIn) {
		d.Fatalf("no page specified")
	}
	postgres.DecodePgBTree(d, pgIn.Page)
	return nil
}
