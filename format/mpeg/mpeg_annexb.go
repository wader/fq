package mpeg

import (
	"github.com/wader/fq/format"
	"github.com/wader/fq/format/registry"
	"github.com/wader/fq/pkg/decode"
)

var annexBAVCNALUFormat []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.MPEG_ANNEXB,
		Description: "H.264/AVC Annex B",
		Groups:      []string{format.PROBE},
		DecodeFn:    annexBDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.AVC_NALU}, Formats: &annexBAVCNALUFormat},
		},
	})
}

func annexBFindStartCode(d *decode.D) (int64, int64, error) {
	offset, v, err := d.TryPeekFind(32, 8, func(v uint64) bool {
		return annexBDecodeStartCodeLen(v) > 0
	}, d.BitsLeft())
	return offset, annexBDecodeStartCodeLen(v), err
}

func annexBDecodeStartCodeLen(v uint64) int64 {
	switch {
	case v == 0x00_00_00_01:
		return 4 * 8
	case v&0xff_ff_ff_00 == 0x00_00_01_00:
		return 3 * 8
	default:
		return 0
	}
}

func annexBDecode(d *decode.D, in interface{}) interface{} {
	currentOffset, currentPrefixLen, err := annexBFindStartCode(d)
	// TODO: really restrict to 0?
	if err != nil || currentOffset != 0 {
		d.Invalid("could not find start code (first)")
	}

	// TODO: root array
	d.FieldArrayFn("nalus", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldBitBufLen("start_code", currentPrefixLen)

			nextOffset, nextPrefixLen, err := annexBFindStartCode(d)
			if err != nil {
				nextOffset = d.Len() - d.Pos()
			}

			naluLen := nextOffset
			d.FieldFormatLen("nalu", naluLen, avcNALUFormat)

			currentPrefixLen = nextPrefixLen
		}
	})

	return nil
}
