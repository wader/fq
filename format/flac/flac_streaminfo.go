package flac

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.FLAC_STREAMINFO,
		Description: "FLAC streaminfo",
		DecodeFn:    streaminfoDecode,
	})
}

func streaminfoDecode(d *decode.D, in interface{}) interface{} {
	d.FieldU16("minimum_block_size")
	d.FieldU16("maximum_block_size")
	d.FieldU24("minimum_frame_size")
	d.FieldU24("maximum_frame_size")
	sampleRate := d.FieldU("sample_rate", 20)
	// <3> (number of channels)-1. FLAC supports from 1 to 8 channels
	d.FieldU3("channels", d.UAdd(1))
	// <5> (bits per sample)-1. FLAC supports from 4 to 32 bits per sample. Currently the reference encoder and decoders only support up to 24 bits per sample.
	bitPerSample := d.FieldU5("bits_per_sample", d.UAdd(1))
	totalSamplesInStream := d.FieldU("total_samples_in_stream", 36)
	md5BB := d.FieldRawLen("md5", 16*8, d.RawHex)
	md5b, err := md5BB.Bytes()
	if err != nil {
		d.IOPanic(err)
	}

	return format.FlacStreaminfoOut{
		StreamInfo: format.FlacStreamInfo{
			SampleRate:           sampleRate,
			BitPerSample:         bitPerSample,
			TotalSamplesInStream: totalSamplesInStream,
			MD5:                  md5b,
		},
	}
}
