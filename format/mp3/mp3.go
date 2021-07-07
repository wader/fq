package mp3

// TODO: vbri
// TODO: resync on garbage? between id3v2 and first frame for example
// TODO: mime audio/mpeg

import (
	"fq/format"
	"fq/format/registry"
	"fq/pkg/decode"
)

var headerFormat []*decode.Format
var footerFormat []*decode.Format
var mp3Frame []*decode.Format

func init() {
	registry.MustRegister(&decode.Format{
		Name:        format.MP3,
		ProbeOrder:  10, // after most others
		Description: "MP3 file",
		Groups:      []string{format.PROBE},
		DecodeFn:    mp3Decode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.ID3V2}, Formats: &headerFormat},
			{Names: []string{format.ID3V1, format.ID3V11, format.APEV2}, Formats: &footerFormat},
			{Names: []string{format.MP3_FRAME}, Formats: &mp3Frame},
		},
	})
}

func mp3Decode(d *decode.D, in interface{}) interface{} {
	// there are mp3s files in the wild with multiple headers, two id3v2 tags etc
	d.FieldArrayFn("headers", func(d *decode.D) {
		for d.NotEnd() {
			if dv, _, _ := d.FieldTryDecode("header", headerFormat); dv == nil {
				return
			}
		}
	})

	validFrames := 0
	foundInvalid := false
	d.FieldArrayFn("frames", func(d *decode.D) {
		for d.NotEnd() {
			startFindSync := d.Pos()
			syncLen, err := d.TryPeekFind(16, 8, func(v uint64) bool { return v&0b1111_1111_1110_0000 == 0b1111_1111_1110_0000 }, d.BitsLeft())
			if err != nil {
				break
			}
			if syncLen > 0 {
				d.SeekRel(syncLen)
			}

			if dv, _, _ := d.FieldTryDecode("frame", mp3Frame); dv == nil {
				foundInvalid = true
				d.SeekAbs(startFindSync)
				break
			}
			validFrames++
		}
	})
	// TODO: better validate
	if validFrames == 0 || (validFrames < 2 && foundInvalid) {
		d.Invalid("no frames found")
	}

	d.FieldArrayFn("footers", func(d *decode.D) {
		for d.NotEnd() {
			if dv, _, _ := d.FieldTryDecode("footer", footerFormat); dv == nil {
				return
			}
		}
	})

	return nil
}
