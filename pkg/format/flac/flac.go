package flac

// TODO: 24 bit picture length truncate warning
// TODO: reuse samples buffer

import (
	"crypto/md5"
	"fq/pkg/decode"
	"fq/pkg/format"
	"log"
)

var metadatablockFormat []*decode.Format
var frameFormat []*decode.Format

func init() {
	format.MustRegister(&decode.Format{
		Name:        format.FLAC,
		Description: "Free lossless audio codec",
		Groups:      []string{format.PROBE},
		MIMEs:       []string{"audio/x-flac"},
		DecodeFn:    flacDecode,
		Dependencies: []decode.Dependency{
			{Names: []string{format.FLAC_METADATABLOCK}, Formats: &metadatablockFormat},
			{Names: []string{format.FLAC_FRAME}, Formats: &frameFormat},
		},
	})
}

func flacDecode(d *decode.D) interface{} {
	d.FieldValidateUTF8("magic", "fLaC")

	md5Samples := md5.New()

	dv := metadatablockFormat[0].DecodeFn(d)
	si, ok := dv.(*streamInfo)
	if !ok {
		d.Invalid("failed to decode metadatablock")
	}

	d.FieldArrayFn("frame", func(d *decode.D) {
		for d.NotEnd() {
			d.FieldDecode("frame", frameFormat)
		}
	})

	if si.d != nil {
		md5Value := si.d.FieldGet("md5")
		log.Printf("md5Value: %#+v\n", md5Value)
		si.d.FieldChecksumRange("md5_", md5Value.Range.Start, md5Value.Range.Len, md5Samples.Sum(nil), decode.BigEndian)
	}

	return nil
}
