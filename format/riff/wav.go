package riff

// http://soundfile.sapp.org/doc/WaveFormat/
// https://github.com/FFmpeg/FFmpeg/blob/master/libavformat/wavdec.c
// https://tech.ebu.ch/docs/tech/tech3285.pdf
// http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/WAVE/WAVE.html
// TODO: audio/wav
// TODO: default little endian

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

var wavHeaderFormat decode.Group
var wavFooterFormat decode.Group

func init() {
	interp.RegisterFormat(decode.Format{
		Name:        format.WAV,
		ProbeOrder:  format.ProbeOrderBinFuzzy, // after most others (overlap some with webp)
		Description: "WAV file",
		Groups:      []string{format.PROBE},
		DecodeFn:    wavDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ID3V2}, Group: &wavHeaderFormat},
			{Names: []string{format.ID3V1, format.ID3V11}, Group: &wavFooterFormat},
		},
	})
}

const (
	formatExtensible = 0xfffe
)

var (
	subFormatPCMBytes  = [16]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}
	subFormatIEEEFloat = [16]byte{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x00, 0x00, 0xaa, 0x00, 0x38, 0x9b, 0x71}
)

const wavRiffType = "WAVE"

var subFormatNames = scalar.BytesToScalar{
	{Bytes: subFormatPCMBytes[:], Scalar: scalar.S{Sym: "pcm"}},
	{Bytes: subFormatIEEEFloat[:], Scalar: scalar.S{Sym: "ieee_float"}},
}

func wavDecode(d *decode.D, _ any) any {
	d.Endian = decode.LittleEndian

	// there are wav files in the wild with id3v2 header id3v1 footer
	_, _, _ = d.TryFieldFormat("header", wavHeaderFormat, nil)

	var riffType string
	riffDecode(
		d,
		nil,
		func(d *decode.D, path path) (string, int64) {
			id := d.FieldUTF8("id", 4, chunkIDDescriptions)

			const restOfFileLen = 0xffffffff
			size := int64(d.FieldUScalarFn("size", func(d *decode.D) scalar.S {
				l := d.U32()
				if l == restOfFileLen {
					return scalar.S{Actual: l, ActualDisplay: scalar.NumberHex, Description: "Rest of file"}
				}
				return scalar.S{Actual: l, ActualDisplay: scalar.NumberDecimal}
			}))

			if size == restOfFileLen {
				size = d.BitsLeft() / 8
			}

			return id, size
		},
		func(d *decode.D, id string, path path) (bool, any) {
			switch id {
			case "RIFF":
				riffType = d.FieldUTF8("format", 4, d.AssertStr(wavRiffType))
				return true, nil

			case "LIST":
				typ := d.FieldUTF8("type", 4)
				switch typ {
				case "strl":
					return true, &aviStrl{}
				}

				return true, nil

			case "fmt ":
				audioFormat := d.FieldU16("audio_format", format.WAVTagNames)
				d.FieldU16("num_channels")
				d.FieldU32("sample_rate")
				d.FieldU32("byte_rate")
				d.FieldU16("block_align")
				d.FieldU16("bits_per_sample")

				if audioFormat == formatExtensible && d.BitsLeft() > 0 {
					d.FieldU16("extension_size")
					d.FieldU16("valid_bits_per_sample")
					d.FieldU32("channel_mask")
					d.FieldRawLen("sub_format", 16*8, subFormatNames)
				}
				return false, nil

			case "data":
				d.FieldRawLen("samples", d.BitsLeft())
				return false, nil

			case "fact":
				d.FieldU32("sample_length")
				return false, nil

			case "smpl":
				d.FieldU32("manufacturer")
				d.FieldU32("product")
				d.FieldU32("sample_period")
				d.FieldU32("midi_unity_note")
				d.FieldU32("midi_pitch_fraction")
				d.FieldU32("smpte_format")
				d.FieldU32("smpte_offset")
				numSampleLoops := int(d.FieldU32("number_of_sample_loops"))
				samplerDataBytes := int(d.FieldU32("sampler_data_bytes"))
				d.FieldArray("samples_loops", func(d *decode.D) {
					for i := 0; i < numSampleLoops; i++ {
						d.FieldStruct("sample_loop", func(d *decode.D) {
							d.FieldUTF8("id", 4)
							d.FieldU32("type", scalar.UToSymStr{
								0: "forward",
								1: "forward_backward",
								2: "backward",
							})
							d.FieldU32("start")
							d.FieldU32("end")
							d.FieldU32("fraction")
							d.FieldU32("number_of_times")
						})
					}
				})
				d.FieldRawLen("sampler_data", int64(samplerDataBytes)*8)
				return false, nil

			default:
				if riffIsStringChunkID(id) {
					d.FieldUTF8NullFixedLen("value", int(d.BitsLeft())/8)
					return false, nil
				}

				d.FieldRawLen("data", d.BitsLeft())
				return false, nil
			}
		},
	)

	if riffType != wavRiffType {
		d.Errorf("wrong or no AVI riff type found (%s)", riffType)
	}
	_, _, _ = d.TryFieldFormat("footer", wavFooterFormat, nil)

	return nil
}
