package av1

// matroska "Low Overhead Bitstream Format syntax" format
// "Each Block contains one Temporal Unit containing one or more OBUs. Each OBU stored in the Block MUST contain its header and its payload."
// "The OBUs in the Block follow the [Low Overhead Bitstream Format syntax]. They MUST have the [obu_has_size_field] set to 1 except for the last OBU in the frame, for which [obu_has_size_field] MAY be set to 0, in which case it is assumed to fill the remainder of the frame."

import (
	"fq/format"
	"fq/pkg/decode"
)

var obuFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.AV1_FRAME,
		Description: "AV1 frame",
		DecodeFn:    frameDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AV1_OBU}, Formats: &obuFormat},
		},
	})
}

func frameDecode(d *decode.D, in interface{}) interface{} {
	d.FieldArrayFn("frame", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldDecode("obu", obuFormat)
		}
	})

	return nil
}
