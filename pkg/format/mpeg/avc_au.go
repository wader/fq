package mpeg

// ISO/IEC 14496-15, 5.3.3.1.2 Syntax

import (
	"fq/pkg/decode"
	"fq/pkg/format"
)

var avcSampleNALFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.MPEG_AVC_AU,
		Description: "H.264/AVC access unit",
		DecodeFn:    avcDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.MPEG_AVC_NALU}, Formats: &avcSampleNALFormat},
		},
	})
}

func avcDecode(d *decode.D, in interface{}) interface{} {
	avcIn, ok := in.(format.AvcIn)
	if !ok {
		d.Invalid("avcIn required")
	}

	d.FieldArrayFn("sample", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldStructFn("nal", func(d *decode.D) {
				l := d.FieldU("length", int(avcIn.LengthSize)*8)
				d.FieldDecodeLen("nal", int64(l)*8, avcSampleNALFormat)
			})
		}
	})

	return nil
}
