package bitcoin

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var bitcoinBlockGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.Bitcoin_Blkdat,
		&decode.Format{
			Description: "Bitcoin blk.dat",
			Groups:      []*decode.Group{format.Probe},
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.Bitcoin_Block}, Out: &bitcoinBlockGroup},
			},
			DecodeFn:  decodeBlkDat,
			RootArray: true,
			RootName:  "blocks",
		})
}

func decodeBlkDat(d *decode.D) any {
	validBlocks := 0
	for !d.End() {
		d.FieldFormat("block", &bitcoinBlockGroup, format.Bitcoin_Block_In{HasHeader: true})
		validBlocks++
	}

	if validBlocks == 0 {
		d.Fatalf("no valid blocks found")
	}

	return nil
}
