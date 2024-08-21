package riff

// http://midi.teragonaudio.com/tech/aiff.htm

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/pkg/decode"
	"github.com/wader/fq/pkg/interp"
	"github.com/wader/fq/pkg/scalar"
)

func init() {
	interp.RegisterFormat(
		format.AIFF,
		&decode.Format{
			ProbeOrder:  format.ProbeOrderBinFuzzy,
			Description: "Audio Interchange File Format",
			Groups:      []*decode.Group{format.Probe},
			DecodeFn:    aiffDecode,
		})
}

const aiffRiffType = "AIFF"

// pstring:
// > Pascal-style string, a one-byte count followed by that many text bytes. The total number of bytes in this data type should be even.
// > A pad byte can be added to the end of the text to accomplish this. This pad byte is not reflected in the count.
func aiffPString(d *decode.D) string {
	l := d.U8()
	pad := (l + 1) % 2
	s := d.UTF8(int(l + pad))
	return s[0 : l+1-pad]
}

func aiffDecode(d *decode.D) any {
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
			case "FORM":
				riffType = d.FieldUTF8("format", 4, d.StrAssert(aiffRiffType))
				return true, nil
			case "COMT":
				numComments := d.FieldU16("num_comments")
				d.FieldArray("comments", func(d *decode.D) {
					for i := 0; i < int(numComments); i++ {
						d.FieldStruct("comment", func(d *decode.D) {
							d.FieldU32("timestamp")
							d.FieldU16("marker_id")
							count := d.FieldU16("count")
							pad := count % 2
							d.FieldUTF8("text", int(count))
							if pad != 0 {
								d.FieldRawLen("pad", int64(pad)*8)
							}
						})
					}
				})
				return false, nil
			case "COMM":
				d.FieldU16("num_channels")
				d.FieldU32("num_sample_frames")
				d.FieldU16("sample_size")
				// TODO: support big float?
				d.FieldF80("sample_rate")
				return false, nil
			case "SSND":
				d.FieldU32("offset")
				d.FieldU32("block_size")
				d.FieldRawLen("data", d.BitsLeft())
				return false, nil
			case "MARK":
				numMarkers := d.FieldU16("num_markers")
				d.FieldArray("markers", func(d *decode.D) {
					for i := 0; i < int(numMarkers); i++ {
						d.FieldStruct("marker", func(d *decode.D) {
							d.FieldU16("id")
							d.FieldU32("position")
							d.FieldStrFn("name", aiffPString)
						})
					}
				})
				return false, nil
			default:
				d.FieldRawLen("data", d.BitsLeft())
				return false, nil
			}
		},
	)

	if riffType != aiffRiffType {
		d.Errorf("wrong or no AIFF riff type found (%s)", riffType)
	}

	return nil
}
