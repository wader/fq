package flac

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var flacMetadatablockForamt decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.FLAC_METADATABLOCKS,
		Description: "FLAC metadatablocks",
		DecodeFn:    metadatablocksDecode,
		RootArray:   true,
		RootName:    "metadatablocks",
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_METADATABLOCK}, Group: &flacMetadatablockForamt},
		},
	})
}

func metadatablocksDecode(d *decode.D, in interface{}) interface{} {
	flacMetadatablocksOut := format.FlacMetadatablocksOut{}

	isLastBlock := false
	for !isLastBlock {
		v, dv := d.FieldFormat("metadatablock", flacMetadatablockForamt, nil)
		flacMetadatablockOut, ok := dv.(format.FlacMetadatablockOut)
		if v != nil && !ok {
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
