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

var wavHeaderGroup decode.Group
var wavFooterGroup decode.Group

func init() {
	interp.RegisterFormat(
		format.WAV,
		&decode.Format{
			ProbeOrder:  format.ProbeOrderBinFuzzy, // after most others (overlap some with webp)
			Description: "WAV file",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    wavDecode,
			Dependencies: []decode.Dependency{
				{Groups: []*decode.Group{format.ID3v2}, Out: &wavHeaderGroup},
				{Groups: []*decode.Group{format.ID3v1, format.ID3v11}, Out: &wavFooterGroup},
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

var subFormatNames = scalar.RawBytesMap{
	{Bytes: subFormatPCMBytes[:], Scalar: scalar.BitBuf{Sym: "pcm"}},
	{Bytes: subFormatIEEEFloat[:], Scalar: scalar.BitBuf{Sym: "ieee_float"}},
}

func wavDecode(d *decode.D) any {
	d.Endian = decode.LittleEndian

	// there are wav files in the wild with id3v2 header id3v1 footer
	_, _, _ = d.TryFieldFormat("header", &wavHeaderGroup, nil)

	var riffType string
	riffDecode(
		d,
		nil,
		func(d *decode.D, path path) (string, int64) {
			id := d.FieldUTF8("id", 4, scalar.ActualTrimSpace, chunkIDDescriptions)

			const restOfFileLen = 0xffffffff
			size := int64(d.FieldScalarUintFn("size", func(d *decode.D) scalar.Uint {
				l := d.U32()
				if l == restOfFileLen {
					return scalar.Uint{Actual: l, DisplayFormat: scalar.NumberHex, Description: "Rest of file"}
				}
				return scalar.Uint{Actual: l, DisplayFormat: scalar.NumberDecimal}
			}).Actual)

			if size == restOfFileLen {
				size = d.BitsLeft() / 8
			}

			return id, size
		},
		func(d *decode.D, id string, path path) (bool, any) {
			switch id {
			case "RIFF":
				riffType = d.FieldUTF8("format", 4, d.StrAssert(wavRiffType))
				return true, nil

			case "LIST":
				d.FieldUTF8("type", 4)
				return true, nil

			case "fmt":
				audioFormat := d.FieldU16("audio_format", format.WAVTagNames)
				d.FieldU16("num_channels")
				d.FieldU32("sample_rate")
				d.FieldU32("byte_rate")
				d.FieldU16("block_align")
				d.FieldU16("bits_per_sample")

				if d.BitsLeft() > 0 {
					if audioFormat == formatExtensible {
						d.FieldU16("extension_size")
						d.FieldU16("valid_bits_per_sample")
						d.FieldU32("channel_mask")
						d.FieldRawLen("sub_format", 16*8, subFormatNames)
					} else {
						cbSize := d.FieldU16("cb_size")
						d.FieldRawLen("unknown", int64(cbSize)*8)
					}
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
							d.FieldU32("type", scalar.UintMapSymStr{
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

			case "bext":
				d.FieldUTF8NullFixedLen("description", 256)
				d.FieldUTF8NullFixedLen("originator", 32)
				d.FieldUTF8NullFixedLen("originator_reference", 32)
				d.FieldUTF8NullFixedLen("originator_date", 10)
				d.FieldUTF8NullFixedLen("originator_time", 8)
				d.FieldU32("time_reference_low")
				d.FieldU32("time_reference_high")
				d.FieldU16("version")
				d.FieldRawLen("umid", 64*8)
				d.FieldU16("loudness_value")
				d.FieldU16("loudness_range")
				d.FieldU16("max_true_peak_level")
				d.FieldU16("max_momentary_loudness")
				d.FieldU16("max_short_term_loudness")
				d.FieldRawLen("reserved", 180*8)
				d.FieldRawLen("coding_history", d.BitsLeft())
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
		d.Errorf("wrong or no WAV riff type found (%s)", riffType)
	}
	_, _, _ = d.TryFieldFormat("footer", &wavFooterGroup, nil)

	return nil
}
