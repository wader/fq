package flac

// TODO: reuse samples buffer

import (
	"bytes"
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
		Description: "Free Lossless Audio Codec file",
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
	var flacFrameIn format.FlacFrameIn
	var framesNDecodedSamples uint64

	d.FieldArrayFn("metadatablocks", func(d *decode.D) {
		for {
			_, v := d.FieldDecode("metadatablock", flacMetadatablockFormat)
			flacMetadatablockOut, ok := v.(format.FlacMetadatablockOut)
			if !ok {
				panic(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
			}
			if flacMetadatablockOut.HasStreamInfo {
				streamInfo = flacMetadatablockOut.StreamInfo
				flacFrameIn = format.FlacFrameIn{
					StreamInfo:   streamInfo,
					NSamplesLeft: streamInfo.TotalSamplesInStream,
				}
			}
			if flacMetadatablockOut.IsLastBlock {
				return
			}
		}
	})

	md5Samples := md5.New()
	d.FieldArrayFn("frames", func(d *decode.D) {
		for d.NotEnd() {
			// flac frame might need some fields from stream info to decode
			_, v := d.FieldDecode("frame", flacFrameFormat, decode.FormatOptions{InArg: flacFrameIn})
			ffo, ok := v.(*format.FlacFrameOut)
			if !ok {
				panic(fmt.Sprintf("expected FlacFrameOut got %#+v", v))
			}

			decode.MustCopy(md5Samples, bytes.NewReader(ffo.SamplesBuf))
			framesNDecodedSamples += uint64(ffo.NDecodedSamples)
			// 0 total samples means unknown
			if streamInfo.TotalSamplesInStream > 0 {
				flacFrameIn.NSamplesLeft -= ffo.NSteamSamples
			}
		}
	})

	_ = streamInfo
	// if streamInfo.D != nil {
	// 	md5Value := streamInfo.D.FieldGet("md5")
	// 	d.FieldChecksumRange("md5_calculated", md5Value.Range.Start, md5Value.Range.Len, md5Samples.Sum(nil), decode.BigEndian)
	// }
	d.FieldValueBytes("md5_calculated", md5Samples.Sum(nil), "")
	d.FieldValueU("decoded_samples", framesNDecodedSamples, "")

	return nil
}
