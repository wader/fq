package flac

// TODO: 24 bit picture length truncate warning
// TODO: reuse samples buffer

import (
	"crypto/md5"
	"fq/pkg/decode"
	"fq/pkg/format"
)

var metadatablockFormat []*decode.Format
var frameFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.FLAC,
		Description: "Free lossless audio codec",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"audio/x-flac"},
		DecodeFn:    flacDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &metadatablockFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &frameFormat},
		},
	})
}

func flacDecode(d *decode.D) interface{} {
	d.FieldValidateUTF8("magic", "fLaC")

	md5Samples := md5.New()

	var flacFrameIn *format.FlacFrameIn
	var streamInfo format.FlacMetadatablockStreamInfo

	dv := d.Decode(metadatablockFormat)
	if dv, ok := dv.(*format.FlacMetadatablockOut); ok {
		streamInfo = dv.StreamInfo
		flacFrameIn = &format.FlacFrameIn{StreamInfo: dv.StreamInfo}
	}

	d.FieldArrayFn("frame", func(d *decode.D) {
		for d.NotEnd() {
			// flac frame might need some fields from stream info to decode
			_, dv := d.FieldDecode("frame", frameFormat, decode.FormatOptions{InArg: flacFrameIn})
			if dv, ok := dv.(*format.FlacFrameOut); ok {
				if _, err := md5Samples.Write(dv.SamplesBuf); err != nil {
					panic(err)
				}
			}
		}
	})

	_ = streamInfo
	// if streamInfo.D != nil {
	// 	md5Value := streamInfo.D.FieldGet("md5")
	// 	d.FieldChecksumRange("md5_calculated", md5Value.Range.Start, md5Value.Range.Len, md5Samples.Sum(nil), decode.BigEndian)
	// }

	return nil
}
