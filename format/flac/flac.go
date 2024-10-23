package flac

// https://xiph.org/flac/format.html
// https://wiki.hydrogenaud.io/index.php?title=FLAC_decoder_testbench

import (
	"bytes"
	"crypto/md5"
	"fmt"

	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/bitio"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var flacMetadatablocksGroup decode.Group
var flacFrameGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.FLAC,
		&decode.Format{
			Description: "Free Lossless Audio Codec file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    flacDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.FLAC_Metadatablocks}, Out: &flacMetadatablocksGroup},
				{Groups: []*decode.Group{format.FLAC_Frame}, Out: &flacFrameGroup},
			},
		})
}

func flacDecode(d *decode.D) any {
	d.FieldUTF8("magic", 4, d.StrAssert("fLaC"))

	var streamInfo format.FLAC_Stream_Info
	var flacFrameIn format.FLAC_Frame_In
	var framesNDecodedSamples uint64
	var streamTotalSamples uint64
	var streamDecodedSamples uint64

	_, v := d.FieldFormat("metadatablocks", &flacMetadatablocksGroup, nil)
	flacMetadatablockOut, ok := v.(format.FLAC_Metadatablocks_Out)
	if !ok {
		panic(fmt.Sprintf("expected FlacMetadatablockOut got %#+v", v))
	}
	if flacMetadatablockOut.HasStreamInfo {
		streamInfo = flacMetadatablockOut.StreamInfo
		streamTotalSamples = streamInfo.TotalSamplesInStream
		flacFrameIn = format.FLAC_Frame_In{BitsPerSample: int(streamInfo.BitsPerSample)}
	}

	md5Samples := md5.New()
	d.FieldArray("frames", func(d *decode.D) {
		for d.NotEnd() {
			// flac frame might need some fields from stream info to decode
			_, v := d.FieldFormat("frame", &flacFrameGroup, flacFrameIn)
			ffo, ok := v.(format.FLAC_Frame_Out)
			if !ok {
				panic(fmt.Sprintf("expected FlacFrameOut got %#+v", v))
			}

			samplesInFrame := ffo.Samples
			if streamTotalSamples > 0 {
				samplesInFrame = min(streamTotalSamples-streamDecodedSamples, ffo.Samples)
			}
			frameStreamSamplesBuf := ffo.SamplesBuf[0 : samplesInFrame*uint64(ffo.Channels*ffo.BitsPerSample/8)]
			framesNDecodedSamples += ffo.Samples

			d.Copy(md5Samples, bytes.NewReader(frameStreamSamplesBuf))
			streamDecodedSamples += ffo.Samples

			// reuse buffer if possible
			flacFrameIn.SamplesBuf = ffo.SamplesBuf
		}
	})

	md5CalcValue := d.FieldRootBitBuf("md5_calculated", bitio.NewBitReader(md5Samples.Sum(nil), -1))
	_ = md5CalcValue.TryBitBufScalarFn(d.ValidateBitBuf(streamInfo.MD5), scalar.RawHex)
	d.FieldValueUint("decoded_samples", framesNDecodedSamples)

	return nil
}
