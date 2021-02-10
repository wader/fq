package av1

// matroska "Low Overhead Bitstream Format syntax" format
// "Each Block contains one Temporal Unit containing one or more OBUs. Each OBU stored in the Block MUST contain its header and its payload."
// "The OBUs in the Block follow the [Low Overhead Bitstream Format syntax]. They MUST have the [obu_has_size_field] set to 1 except for the last OBU in the frame, for which [obu_has_size_field] MAY be set to 0, in which case it is assumed to fill the remainder of the frame."

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.AV1_FRAME,
		Description: "AV1 frame",
		DecodeFn:    frameDecode,
	})
}

func frameDecode(d *decode.D, in interface{}) interface{} {
	d.FieldArrayFn("frame", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStructFn("obu", func(d *decode.D) {
				obuDecode(d, nil)
			})
		}
	})

	return nil
}
