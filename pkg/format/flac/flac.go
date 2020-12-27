package flac

// TODO: 24 bit picture length truncate warning
// TODO: reuse samples buffer

import (
	"crypto/md5"
	"fmt"
	"fq/pkg/decode"
	"fq/pkg/format"
)

var flacMetadatablockFormat []*decode.Format
var flacFrameFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.FLAC,
		Description: "Free lossless audio codec",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"audio/x-flac"},
		DecodeFn:    flacDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &flacMetadatablockFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &flacFrameFormat},
		},
	})
}

func flacDecode(d *decode.D, in interface{}) interface{} {
	d.FieldValidateUTF8("magic", "fLaC")

	var streamInfo format.FlacMetadatablockStreamInfo
	var flacFrameIn *format.FlacFrameIn

	d.FieldArrayFn("metadatablocks", func(d *decode.D) {
		for {
			_, dv := d.FieldDecode("metadatablock", flacMetadatablockFormat)
			flacMetadatablockOut, _ := dv.(*format.FlacMetadatablockOut)
			if flacMetadatablockOut == nil {
				d.Invalid(fmt.Sprintf("expected FlacMetadatablockOut got %v", dv))
			}
			if flacMetadatablockOut.StreamInfo != nil {
				streamInfo = *flacMetadatablockOut.StreamInfo
				flacFrameIn = &format.FlacFrameIn{StreamInfo: streamInfo}
			}
			if flacMetadatablockOut.LastBlock {
				return
			}
		}
	})

	md5Samples := md5.New()
	d.FieldArrayFn("frame", func(d *decode.D) {
		for d.NotEnd() {
			// flac frame might need some fields from stream info to decode
			_, dv := d.FieldDecode("frame", flacFrameFormat, decode.FormatOptions{InArg: flacFrameIn})
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
