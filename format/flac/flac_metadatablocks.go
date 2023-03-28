package flac

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
)

var flacMetadatablockGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.FlacMetadatablocks,
		&decode.Format{
			Description: "FLAC metadatablocks",
			DecodeFn:    metadatablocksDecode,
			RootArray:   true,
			RootName:    "metadatablocks",
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.FlacMetadatablock}, Out: &flacMetadatablockGroup},
			},
		})
}

func metadatablocksDecode(d *decode.D) any {
	flacMetadatablocksOut := format.FlacMetadatablocksOut{}

	isLastBlock := false
	for !isLastBlock {
		_, v := d.FieldFormat("metadatablock", &flacMetadatablockGroup, nil)
		flacMetadatablockOut, ok := v.(format.FlacMetadatablockOut)
		if !ok {
			panic(fmt.Sprintf("expected FlacMetadatablocksOut, got %#+v", flacMetadatablockOut))
		}
		isLastBlock = flacMetadatablockOut.IsLastBlock
		if flacMetadatablockOut.HasStreamInfo {
			flacMetadatablocksOut.HasStreamInfo = true
			flacMetadatablocksOut.StreamInfo = flacMetadatablockOut.StreamInfo
		}
	}

	return flacMetadatablocksOut
}
