package bitcoin

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var bitcoinBlockFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
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

func decodeBlkDat(d *decode.D, in any) any {
	validBlocks := 0
	for !d.End() {
		d.FieldFormat("block", bitcoinBlockFormat, format.BitCoinBlockIn{HasHeader: true})
		validBlocks++
	}

	if validBlocks == 0 {
		d.Fatalf("no valid blocks found")
	}

	return nil
}
