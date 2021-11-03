package flac

import (
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var flacMetadatablockForamt []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.FLAC_METADATABLOCKS,
		Description: "FLAC metadatablocks",
		DecodeFn:    metadatablocksDecode,
		RootArray:   true,
		RootName:    "metadatablocks",
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &flacMetadatablockForamt},
		},
	})
}

func metadatablocksDecode(d *decode.D, in interface{}) interface{} {
	flacMetadatablocksOut := format.FlacMetadatablocksOut{}

	isLastBlock := false
	for !isLastBlock {
		_, flacMetadatablockOutAny := d.FieldFormat("metadatablock", flacMetadatablockForamt, nil)
		flacMetadatablockOut, ok := flacMetadatablockOutAny.(format.FlacMetadatablockOut)
		if !ok {
			d.Invalid(fmt.Sprintf("expected FlacMetadatablocksOut, got %#+v", flacMetadatablockOut))
		}
		isLastBlock = flacMetadatablockOut.IsLastBlock
		if flacMetadatablockOut.HasStreamInfo {
			flacMetadatablocksOut.HasStreamInfo = true
			flacMetadatablocksOut.StreamInfo = flacMetadatablockOut.StreamInfo
		}
	}

	return flacMetadatablocksOut
}
