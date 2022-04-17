package bitcoin

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var bitcoinBlockFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.BITCOIN_BLKDAT,
		Description: "Bitcoin blk.dat",
		Groups:      []string{format.PROBE},
		Dependencies: []decode.Dependency{
			{Names: []string{format.BITCOIN_BLOCK}, Group: &bitcoinBlockFormat},
		},
		DecodeFn:  decodeBlkDat,
		RootArray: true,
		RootName:  "blocks",
	})
}

func decodeBlkDat(d *decode.D, in interface{}) interface{} {
	validBlocks := 0
	for !d.End() {
		d.FieldFormat("block", bitcoinBlockFormat, nil)
		validBlocks++
	}

	if validBlocks == 0 {
		d.Fatalf("no valid blocks found")
	}

	return nil
}
