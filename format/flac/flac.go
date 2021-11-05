package flac

// https://xiph.org/flac/format.html
// https://wiki.hydrogenaud.io/index.php?title=FLAC_decoder_testbench
// TODO:
// 16 - Part 6 of Ladybug Castle (partition order 8 with escape codes)
// 32 - Part 5 of The Four of Us Are Dying (partition order 8 with escape codes)

import (
	"bytes"
	"crypto/md5"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/internal/num"
	"github.com/wader/fq/pkg/decode"
)

var flacMetadatablockFormat []*decode.Format
var flacFrameFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.FLAC,
		Description: "Free Lossless Audio Codec file",
		Groups:      []string{format.PROBE},
		DecodeFn:    flacDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_METADATABLOCKS}, Formats: &flacMetadatablockFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &flacFrameFormat},
		},
	})
}

func flacDecode(d *decode.D, in interface{}) interface{} {
	d.FieldUTF8("magic", 4, d.AssertStr("fLaC"))

	var streamInfo format.FlacStreamInfo
	var flacFrameIn format.FlacFrameIn
	var framesNDecodedSamples uint64
	var streamTotalSamples uint64
	var streamDecodedSamples uint64

	_, v := d.FieldFormat("metadatablocks", flacMetadatablockFormat, nil)
	flacMetadatablockOut, ok := v.(format.FlacMetadatablocksOut)
	if !ok {
		panic(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
	}
	if flacMetadatablockOut.HasStreamInfo {
		streamInfo = flacMetadatablockOut.StreamInfo
		streamTotalSamples = streamInfo.TotalSamplesInStream
		flacFrameIn = format.FlacFrameIn{StreamInfo: streamInfo}
	}

	md5Samples := md5.New()
	d.FieldArray("frames", func(d *decode.D) {
		for d.NotEnd() {
			// flac frame might need some fields from stream info to decode
			_, v := d.FieldFormat("frame", flacFrameFormat, flacFrameIn)
			ffo, ok := v.(format.FlacFrameOut)
			if !ok {
				panic(fmt.Sprintf("expected FlacFrameOut got %#+v", v))
			}

			samplesInFrame := ffo.Samples
			if streamTotalSamples > 0 {
				samplesInFrame = num.MinUInt64(streamTotalSamples-streamDecodedSamples, ffo.Samples)
			}
			frameStreamSamplesBuf := ffo.SamplesBuf[0 : samplesInFrame*uint64(ffo.Channels*ffo.BitsPerSample/8)]
			framesNDecodedSamples += ffo.Samples

			decode.MustCopy(d, md5Samples, bytes.NewReader(frameStreamSamplesBuf))
			streamDecodedSamples += ffo.Samples

			// reuse buffer if possible
			flacFrameIn.SamplesBuf = ffo.SamplesBuf
		}
	})

	d.FieldValueRaw("md5_calculated", md5Samples.Sum(nil), d.ValidateRaw(streamInfo.MD5), d.RawHex)
	d.FieldValueU("decoded_samples", framesNDecodedSamples)

	return nil
}
