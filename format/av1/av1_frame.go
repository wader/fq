package av1

// matroska "Low Overhead Bitstream Format syntax" format
// "Each Block contains one Temporal Unit containing one or more OBUs. Each OBU stored in the Block MUST contain its header and its payload."
// "The OBUs in the Block follow the [Low Overhead Bitstream Format syntax]. They MUST have the [obu_has_size_field] set to 1 except for the last OBU in the frame, for which [obu_has_size_field] MAY be set to 0, in which case it is assumed to fill the remainder of the frame."

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var obuFormat decode.Group

func init() {
	registry.MustRegister(decode.Format{
		Name:        format.AV1_FRAME,
		Description: "AV1 frame",
		DecodeFn:    frameDecode,
		RootArray:   true,
		RootName:    "frame",
		Dependencies: []decode.Dependency{
			{Names: []string{format.AV1_OBU}, Group: &obuFormat},
		},
	})
}

func frameDecode(d *decode.D, in interface{}) interface{} {
	for d.NotEnd() {
		d.FieldFormat("obu", obuFormat, nil)
	}

	return nil
}
