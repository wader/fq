package flac

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var flacMetadatablockFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.FLAC_METADATABLOCKS,
		Description: "FLAC metadatablocks",
		DecodeFn:    metadatablocksDecode,
		RootArray:   true,
		RootName:    "metadatablocks",
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_METADATABLOCK}, Group: &flacMetadatablockFormat},
		},
	})
}

func metadatablocksDecode(d *decode.D, in any) any {
	flacMetadatablocksOut := format.FlacMetadatablocksOut{}

	isLastBlock := false
	for !isLastBlock {
		dv, v := d.FieldFormat("metadatablock", flacMetadatablockFormat, nil)
		flacMetadatablockOut, ok := v.(format.FlacMetadatablockOut)
		if dv != nil && !ok {
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
